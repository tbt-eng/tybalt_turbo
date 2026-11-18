<script lang="ts">
  import DsList from "$lib/components/DSList.svelte";
  import DsActionButton from "$lib/components/DSActionButton.svelte";
  import DSHelp from "$lib/components/DSHelp.svelte";
  import RateSheetLabel from "$lib/components/rateSheets/RateSheetLabel.svelte";
  import RateSheetMatchLine from "$lib/components/rateSheets/RateSheetMatchLine.svelte";
  import {
    rateSheetListMatch,
    rateSheetSearchPlaceholder,
    rateSheetSearchText,
  } from "$lib/rateSheetLabel";
  import { canCreateRateSheets, canReviseRateSheets } from "$lib/rateSheetPermissions";
  import { globalStore } from "$lib/stores/global";
  import type { RateSheetResponse } from "$lib/stores/rateSheets";
  import { rateSheets } from "$lib/stores/rateSheets";
  import { shortDate } from "$lib/utilities";

  // initialize the store, noop if already initialized
  rateSheets.init();

  // Revise and Use as Template show only to the users that the server allows.
  const canRevise = $derived(canReviseRateSheets($globalStore.claims));
  const canCreate = $derived(canCreateRateSheets($globalStore.claims));

  let searchTerm = $state("");
  const matches = $derived(rateSheetListMatch(searchTerm));
</script>

{#if $rateSheets.items.length > 0 || $rateSheets.initialized}
  <DsList
    items={$rateSheets.items}
    search={true}
    bind:searchTerm
    searchText={rateSheetSearchText}
    searchPlaceholder={rateSheetSearchPlaceholder}
    inListHeader="Rate Sheets"
  >
    {#snippet searchBarExtra()}
      <DsActionButton
        action="/rate-sheets/duplicates"
        icon="mdi:content-duplicate"
        title="Sheets with identical rates"
        color="blue">Duplicates</DsActionButton
      >
      <DSHelp doc="rate-sheets#find-a-rate-sheet" label="Help" />
    {/snippet}

    {#snippet headline(sheet: RateSheetResponse)}
      <RateSheetLabel {sheet} />
    {/snippet}

    {#snippet byline({ effective_date, job_count }: RateSheetResponse)}
      <span class="text-neutral-500">
        {shortDate(effective_date, true)} • {job_count} job{job_count === 1 ? "" : "s"}
      </span>
    {/snippet}

    {#snippet line1({ active }: RateSheetResponse)}
      {#if active}
        <span
          class="inline-flex items-center rounded-full bg-green-100 px-2 py-0.5 text-xs font-medium text-green-800"
        >
          Active
        </span>
      {:else}
        <span
          class="inline-flex items-center rounded-full bg-neutral-100 px-2 py-0.5 text-xs font-medium text-neutral-600"
        >
          Inactive
        </span>
      {/if}
    {/snippet}

    {#snippet line2(sheet: RateSheetResponse)}
      <RateSheetMatchLine {sheet} {matches} />
    {/snippet}

    {#snippet actions({ id, active }: RateSheetResponse)}
      {#if active && canRevise}
        <DsActionButton
          action={`/rate-sheets/copy?revise=${id}`}
          icon="mdi:file-replace-outline"
          title="Revise"
          color="blue"
        />
      {/if}
      {#if canCreate}
        <DsActionButton
          action={`/rate-sheets/copy?from=${id}`}
          icon="mdi:content-copy"
          title="Use as Template"
          color="green"
        />
      {/if}
    {/snippet}
  </DsList>
{:else if $rateSheets.loading}
  <div class="p-4 text-neutral-500">Loading rate sheets...</div>
{/if}
