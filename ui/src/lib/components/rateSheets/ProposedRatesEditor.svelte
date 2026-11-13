<script lang="ts" module>
  import type { ProposedRate } from "./RateSheetRatesTable.svelte";

  // The proposed rates of a project request, as the request editor holds them.
  export interface ProposedRates {
    name: string;
    reason: string;
    based_on: string;
    rows: ProposedRate[];
  }

  export function emptyProposedRates(): ProposedRates {
    return { name: "", reason: "", based_on: "", rows: [] };
  }

  // The proposed_rates object of a request save.
  export function proposedRatesPayload(proposal: ProposedRates) {
    return {
      name: proposal.name,
      reason: proposal.reason,
      based_on: proposal.based_on,
      entries: proposal.rows.map((row) => ({
        role: row.role,
        rate: Number(row.rate),
        overtime_rate: Number(row.overtime_rate),
      })),
    };
  }
</script>

<script lang="ts">
  import { pb } from "$lib/pocketbase";
  import DsAutoComplete from "$lib/components/DSAutoComplete.svelte";
  import DsTextInput from "$lib/components/DSTextInput.svelte";
  import DSCallout from "$lib/components/DSCallout.svelte";
  import DSHelp from "$lib/components/DSHelp.svelte";
  import RateSheetRatesTable from "./RateSheetRatesTable.svelte";
  import { rateSheets } from "$lib/stores/rateSheets";

  // Proposed rates in a project request: the requester starts from an active
  // sheet, names the new sheet, says why it is needed and changes the rates.
  // The server keeps them as a draft that only this request uses.
  let { proposal = $bindable() }: { proposal: ProposedRates } = $props();

  rateSheets.init();
  const noErrors: Record<string, { message: string }> = {};
  let loadError = $state("");
  // The source whose rates the rows show, so that a change of source reloads
  // them and a saved proposal keeps its own rates.
  let loadedSource = $state(proposal.based_on);

  const sourceName = $derived(
    $rateSheets.items.find((sheet) => sheet.id === proposal.based_on)?.name ?? "",
  );

  // A new source replaces the rates with its rates.
  $effect(() => {
    const source = proposal.based_on;
    if (!source || source === loadedSource) return;
    let current = true;
    loadError = "";
    pb.collection("rate_sheet_entries")
      .getFullList({ filter: pb.filter("rate_sheet = {:source}", { source }), expand: "role" })
      .then((entries) => {
        if (!current) return;
        loadedSource = source;
        proposal.rows = entries
          .map((entry) => ({
            role: entry.role,
            roleName: entry.expand?.role?.name ?? entry.role,
            rate: entry.rate,
            overtime_rate: entry.overtime_rate,
            source_rate: entry.rate,
            source_overtime_rate: entry.overtime_rate,
          }))
          .sort((a, b) => a.roleName.localeCompare(b.roleName));
      })
      .catch(() => {
        if (current) loadError = "Turbo could not load the rates of that sheet.";
      });
    return () => {
      current = false;
    };
  });
</script>

<div class="flex flex-col gap-2">
  {#if $rateSheets.index !== null}
    <DsAutoComplete
      bind:value={proposal.based_on}
      index={$rateSheets.index}
      filter={(sheet) => sheet.active}
      errors={noErrors}
      fieldName="proposed_rates_based_on"
      uiName="Start from"
    >
      {#snippet resultTemplate(row)}{row.name} (rev. {row.revision}){/snippet}
      {#snippet help()}
        <DSHelp doc="proposed-rates#propose-new-rates" iconOnly />
      {/snippet}
    </DsAutoComplete>
  {/if}
  <DsTextInput
    bind:value={proposal.name}
    errors={noErrors}
    fieldName="proposed_rates_name"
    uiName="Name"
  />
  <label class="flex flex-col gap-1">
    Reason
    <textarea
      rows="2"
      maxlength="500"
      class="rounded-sm border border-neutral-300 p-2"
      placeholder="Why the client needs these rates, such as a contract schedule"
      bind:value={proposal.reason}></textarea>
  </label>
  {#if loadError}<DSCallout tone="danger">{loadError}</DSCallout>{/if}
  {#if proposal.rows.length > 0}
    <RateSheetRatesTable bind:rows={proposal.rows} editable {sourceName} />
  {/if}
</div>
