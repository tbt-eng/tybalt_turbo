<script lang="ts">
  import { pb } from "$lib/pocketbase";
  import DSPicker from "$lib/components/DSPicker.svelte";
  import DsTextInput from "$lib/components/DSTextInput.svelte";
  import DSCallout from "$lib/components/DSCallout.svelte";
  import DSHelp from "$lib/components/DSHelp.svelte";
  import DSActionButton from "$lib/components/DSActionButton.svelte";
  import RateSheetRatesTable from "./RateSheetRatesTable.svelte";
  import {
    type ProposedRates,
    proposedRatesPayload,
    rateChanged,
    rateMissing,
    withEveryRole,
    withSourceRates,
  } from "$lib/proposedRates";
  import RateSheetMatchCallout from "./RateSheetMatchCallout.svelte";
  import { rateSheets } from "$lib/stores/rateSheets";
  import { rateRoles } from "$lib/stores/rateRoles";
  import type { RateSheetComparison } from "$lib/rateSheetMatches";
  import { rateSheetCodeRevision, rateSheetLabel } from "$lib/rateSheetLabel";
  import { shortDate } from "$lib/utilities";

  // Proposed rates in a project request: the requester starts from an active
  // sheet, names the new sheet, says why it is needed and changes the rates.
  // The server keeps them as a draft that only this request uses.
  // onUse chooses an active sheet in place of the proposal. blocked is true
  // while the page shows the sheets that refused the last save, so the same
  // sheets are not listed twice. errors holds the field errors of the last
  // save, under the server's keys (proposed_rates.name and so on), and the
  // errors of single rates under rateErrorKey.
  let {
    proposal = $bindable(),
    onUse,
    blocked = false,
    errors = {},
  }: {
    proposal: ProposedRates;
    onUse?: (id: string) => void;
    blocked?: boolean;
    errors?: Record<string, { message: string }>;
  } = $props();

  rateSheets.init();
  rateRoles.init();
  let loadError = $state("");
  // The source whose rates the rows show, so that a change of source reloads
  // them and a saved proposal keeps its own rates.
  let loadedSource = $state(proposal.based_on);

  const source = $derived($rateSheets.items.find((sheet) => sheet.id === proposal.based_on));
  const sourceName = $derived(source ? rateSheetCodeRevision(source) : "");

  // The rates of the source, once loaded, for the source rates column.
  let sourceRates = $state<{
    source: string;
    rates: { role: string; rate: number; overtime_rate: number }[];
  } | null>(null);

  // Load the rates of the source. A new source replaces the rates with its
  // rates; a saved proposal keeps its own.
  $effect(() => {
    const source = proposal.based_on;
    if (!source || sourceRates?.source === source) return;
    let current = true;
    loadError = "";
    pb.collection("rate_sheet_entries")
      .getFullList({ filter: pb.filter("rate_sheet = {:source}", { source }), expand: "role" })
      .then((entries) => {
        if (!current) return;
        const rows = entries
          .map((entry) => ({
            role: entry.role,
            roleName: entry.expand?.role?.name ?? entry.role,
            rate: entry.rate,
            overtime_rate: entry.overtime_rate,
            source_rate: entry.rate,
            source_overtime_rate: entry.overtime_rate,
          }))
          .sort((a, b) => a.roleName.localeCompare(b.roleName));
        sourceRates = {
          source,
          rates: entries.map(({ role, rate, overtime_rate }) => ({ role, rate, overtime_rate })),
        };
        if (source === loadedSource) return;
        loadedSource = source;
        proposal.rows = rows;
      })
      .catch(() => {
        if (current) loadError = "Turbo could not load the rates of that sheet.";
      });
    return () => {
      current = false;
    };
  });

  // The rows always cover every rate role: a role with no known rate gets an
  // empty row to fill in. This runs once the rows belong to the chosen source,
  // for rows loaded from a new source and for a saved proposal alike. It
  // writes only when a role is missing, so it settles after one pass.
  $effect(() => {
    if (!proposal.based_on || loadedSource !== proposal.based_on) return;
    const rows = withEveryRole(proposal.rows, $rateRoles.items);
    if (rows) proposal.rows = rows;
  });

  // Each row shows the source's rates for its role, also the empty row of a
  // role that a saved proposal has no rate for. Only the source rates column
  // changes: a saved proposal keeps its own rates. It writes only when a row
  // lacks source rates that the source has, so it settles after one pass.
  $effect(() => {
    const source = proposal.based_on;
    if (!source || loadedSource !== source || sourceRates?.source !== source) return;
    const rows = withSourceRates(proposal.rows, sourceRates.rates);
    if (rows) proposal.rows = rows;
  });

  // Compare the rates with the active sheets as the requester changes them,
  // so that a sheet that already fits is found before a new one is proposed.
  // Rates with no change equal the source sheet, which needs no notice, and a
  // role with no rate is left out of the comparison as it is of a save. The
  // server leaves the source out of the close sheets, since the table already
  // compares the rates with it.
  const noComparison: RateSheetComparison = { exact: [], near: [] };
  let comparison = $state<RateSheetComparison>(noComparison);
  let compareSequence = 0;
  $effect(() => {
    // The reason is left out so that typing it does not compare again.
    const body = JSON.stringify({ ...proposedRatesPayload(proposal), reason: "" });
    const sequence = ++compareSequence;
    if (!proposal.rows.some((row) => !rateMissing(row) && rateChanged(row))) {
      comparison = noComparison;
      return;
    }
    const timer = setTimeout(async () => {
      try {
        const response = await pb.send("/api/job_requests/compare_rates", {
          method: "POST",
          body,
          requestKey: null,
        });
        if (sequence !== compareSequence) return;
        comparison = { exact: response.exact ?? [], near: response.near ?? [] };
      } catch (error) {
        console.error("Failed to compare the proposed rates:", error);
        if (sequence === compareSequence) comparison = noComparison;
      }
    }, 400);
    return () => clearTimeout(timer);
  });
</script>

<div class="flex flex-col gap-2">
  {#if $rateSheets.index !== null}
    <DSPicker
      bind:value={proposal.based_on}
      items={$rateSheets.items.filter((sheet) => sheet.active)}
      index={$rateSheets.index}
      {errors}
      fieldName="proposed_rates.based_on"
      uiName="Start from"
      empty
    >
      {#snippet option(row)}{rateSheetLabel(row)} · {shortDate(row.effective_date, true)}{/snippet}
      {#snippet help()}
        <DSHelp doc="proposed-rates#propose-new-rates" iconOnly />
      {/snippet}
    </DSPicker>
  {/if}
  <DsTextInput bind:value={proposal.name} {errors} fieldName="proposed_rates.name" uiName="Name" />
  <DsTextInput
    bind:value={proposal.reason}
    {errors}
    fieldName="proposed_rates.reason"
    uiName="Reason"
    rows={2}
    maxlength={500}
    placeholder="Why the client needs these rates, such as a contract schedule"
    helperText="Required to save the proposed rates."
  />
  {#if loadError}<DSCallout tone="danger">{loadError}</DSCallout>{/if}
  {#if proposal.rows.length > 0}
    <RateSheetRatesTable bind:rows={proposal.rows} editable {sourceName} {errors} />
  {/if}
  {#if comparison.exact.length > 0 && !blocked}
    <RateSheetMatchCallout
      tone="warning"
      title="Same rates as an active sheet"
      intro="Turbo won't save these rates. Use that sheet, or change a rate."
      matches={comparison.exact}
    >
      {#snippet actions()}
        {#each comparison.exact as match (match.id)}
          <DSActionButton action={() => onUse?.(match.id)}
            >Use {rateSheetCodeRevision(match)}</DSActionButton
          >
        {/each}
      {/snippet}
    </RateSheetMatchCallout>
  {/if}
  {#if comparison.near.length > 0}
    <RateSheetMatchCallout
      tone="info"
      title="Close to active sheets"
      intro="Use one of these sheets if it fits the client's contract."
      matches={comparison.near}
    >
      {#snippet actions()}
        {#each comparison.near as match (match.id)}
          <DSActionButton action={() => onUse?.(match.id)}
            >Use {rateSheetCodeRevision(match)}</DSActionButton
          >
        {/each}
      {/snippet}
    </RateSheetMatchCallout>
  {/if}
</div>
