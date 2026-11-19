<script lang="ts">
  import { pb } from "$lib/pocketbase";
  import DSActionButton from "$lib/components/DSActionButton.svelte";
  import DSCallout from "$lib/components/DSCallout.svelte";
  import DSList from "$lib/components/DSList.svelte";
  import RateSheetLabel from "./RateSheetLabel.svelte";
  import RateSheetRatesTable from "./RateSheetRatesTable.svelte";
  import { rateRoles } from "$lib/stores/rateRoles";
  import { type ProposedRate, withEveryRole } from "$lib/proposedRates";
  import {
    type FoundRate,
    type FoundRateSheet,
    findRatesErrors,
    findRatesKey,
    findResultMessage,
    findStatusText,
    foundRates,
    rateEntered,
  } from "$lib/findRates";
  import { countLabel } from "$lib/clientWorkspace";
  import { shortDate } from "$lib/utilities";

  // Find by rates in a project request: a rate and an optional overtime rate
  // for each role, and the active sheets that have every rate entered. It
  // stores nothing. onUse chooses a sheet; onPropose carries the rates into
  // proposed rates. Enter in a rate searches at once and never submits the
  // request form. rows is bindable so the rates stay while the requester
  // looks at another choice.
  let {
    rows = $bindable(),
    onUse,
    onPropose,
  }: {
    rows: ProposedRate[];
    onUse: (id: string) => void;
    onPropose: (rates: FoundRate[]) => void;
  } = $props();

  rateRoles.init();

  // A row for every rate role, in role order. It writes only when a role is
  // missing, so it settles after one pass.
  $effect(() => {
    const every = withEveryRole(rows, $rateRoles.items);
    if (every) rows = every;
  });

  const errors = $derived(findRatesErrors(rows));
  const key = $derived(findRatesKey(rows));

  // The result is stored with the key of the rates it describes and shown only
  // while the rates still have that key, as in the proposed rates comparison.
  // The search runs 400 ms after the last change, or at once on Enter in a
  // rate.
  let stored = $state<{ key: string; matches: FoundRateSheet[] | null; failed: boolean }>({
    key: "",
    matches: null,
    failed: false,
  });
  const current = $derived(stored.key === key ? stored : null);
  let sequence = 0;
  let searchNow: (() => void) | null = null;
  $effect(() => {
    const body = key;
    const mine = ++sequence;
    searchNow = null;
    if (!body) return;
    const search = async () => {
      searchNow = null;
      try {
        const response = await pb.send("/api/job_requests/find_rates", {
          method: "POST",
          body,
          requestKey: null,
        });
        if (mine === sequence)
          stored = { key: body, matches: response.matches ?? [], failed: false };
      } catch (error) {
        console.error("Failed to find rate sheets by rates:", error);
        if (mine === sequence) stored = { key: body, matches: null, failed: true };
      }
    };
    const timer = setTimeout(search, 400);
    searchNow = () => {
      clearTimeout(timer);
      search();
    };
    return () => clearTimeout(timer);
  });

  const propose = () => onPropose(foundRates(rows));

  const entered = $derived(rows.some(rateEntered));
  const invalid = $derived(Object.keys(errors).length > 0);
  const count = $derived(current?.matches ? current.matches.length : null);
  const statusText = $derived(
    findStatusText({ empty: !entered, invalid, failed: !!current?.failed, count }),
  );
  // What a screen reader hears: only a result or a failure, not each step.
  const announced = $derived(
    entered && !invalid && (current?.failed || count !== null) ? statusText : "",
  );

  // Below lg the results are under the table, so the status line scrolls to
  // them.
  let results = $state<HTMLElement | null>(null);
  const showResults = () => results?.scrollIntoView({ behavior: "smooth", block: "start" });
</script>

<!-- Below lg: a status line of fixed height under the choice of rates, then the
  table, then the results. From lg up: the table and, beside it, the results in
  a column that stays in view while the page scrolls. -->
<div class="flex flex-col gap-2">
  <p class="sr-only" aria-live="polite">{announced}</p>
  <div class="flex min-h-8 items-center lg:hidden">
    {#if count !== null && entered && !invalid}
      <DSActionButton variant="plain" action={showResults}>{statusText}</DSActionButton>
    {:else}
      <p class="text-sm text-neutral-600" aria-hidden="true">{statusText}</p>
    {/if}
  </div>
  <div class="grid gap-4 lg:grid-cols-[minmax(0,30rem)_minmax(0,1fr)] lg:items-start">
    <RateSheetRatesTable bind:rows find {errors} onEnter={() => searchNow?.()} />
    <div bind:this={results} class="flex scroll-mt-4 flex-col gap-2 lg:sticky lg:top-4">
      {#if !entered}
        <div class="max-lg:hidden">
          <DSCallout tone="info"
            >Enter the rates you know. Turbo lists the active sheets that have every one of them.
            The overtime rate is optional.</DSCallout
          >
        </div>
      {:else if invalid}
        <p class="text-sm text-neutral-600 max-lg:hidden">Correct the rates marked in the table.</p>
      {:else if current?.failed}
        <DSCallout tone="danger">Turbo could not search the rate sheets. Try again.</DSCallout>
      {:else if current?.matches && current.matches.length === 0}
        <DSCallout tone="info" title={findResultMessage(0)}>
          <p>Propose new rates, starting from the sheet that is closest to the contract.</p>
          {#snippet actions()}
            <DSActionButton variant="primary" action={propose}
              >Propose new rates with these rates</DSActionButton
            >
          {/snippet}
        </DSCallout>
      {:else if current?.matches}
        <DSList items={current.matches} inListHeader={findResultMessage(current.matches.length)}>
          {#snippet headline(sheet)}
            <RateSheetLabel {sheet} copy />
          {/snippet}
          {#snippet byline(sheet)}
            {shortDate(sheet.effective_date, true)}
          {/snippet}
          {#snippet line1(sheet)}
            {sheet.job_count > 0 ? countLabel(sheet.job_count, "job") : ""}
          {/snippet}
          {#snippet actions(sheet)}
            <DSActionButton action={() => onUse(sheet.id)}>Use</DSActionButton>
          {/snippet}
        </DSList>
        <div>
          <DSActionButton action={propose}>Propose new rates with these rates</DSActionButton>
        </div>
      {:else}
        <p class="text-sm text-neutral-600 max-lg:hidden">Searching…</p>
      {/if}
    </div>
  </div>
</div>
