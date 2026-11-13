<script lang="ts">
  import DsList from "$lib/components/DSList.svelte";
  import DsActionButton from "$lib/components/DSActionButton.svelte";
  import DsLabel from "$lib/components/DsLabel.svelte";
  import DSCallout from "$lib/components/DSCallout.svelte";
  import DSHelp from "$lib/components/DSHelp.svelte";
  import { globalStore } from "$lib/stores/global";
  import { shortDate } from "$lib/utilities";
  import {
    canAbsorbRateSheets,
    type RateSheetDuplicate,
    type RateSheetDuplicateGroup,
  } from "$lib/rateSheetAbsorb";
  import { resolve } from "$app/paths";
  import type { PageData } from "./$types";

  let { data }: { data: PageData } = $props();

  // One row per sheet. group is the heading of the sheet's group: its target,
  // which the server puts first.
  type DuplicateRow = RateSheetDuplicate & { group: string };

  function heading({ sheets }: RateSheetDuplicateGroup): string {
    return `Same rates as ${sheets[0].name} (rev. ${sheets[0].revision})`;
  }

  const rows = $derived<DuplicateRow[]>(
    data.groups.flatMap((group) =>
      group.sheets.map((sheet) => ({ ...sheet, group: heading(group) })),
    ),
  );

  // Why no sheet of a group can be absorbed, by the group's heading.
  const groupReasons = $derived(
    new Map(data.groups.map((group) => [heading(group), group.reason])),
  );

  // The server refuses an absorb by anyone else.
  const canAbsorb = $derived(canAbsorbRateSheets($globalStore.claims));
</script>

{#if rows.length === 0}
  <div class="flex flex-col gap-2 p-4">
    <div class="flex justify-end">
      <DSHelp doc="rate-sheet-absorb#find-duplicate-rate-sheets" label="Help" />
    </div>
    <DSCallout tone="info">No rate sheets have identical rates.</DSCallout>
  </div>
{:else}
  <DsList items={rows} search={true} inListHeader="Duplicate Rate Sheets" groupField="group">
    {#snippet searchBarExtra()}
      <DSHelp doc="rate-sheet-absorb#find-duplicate-rate-sheets" label="Help" />
    {/snippet}

    {#snippet headline({ id, name }: DuplicateRow)}
      <a href={resolve("/rate-sheets/[id]/details", { id })} class="link">{name}</a>
    {/snippet}

    {#snippet byline({ effective_date, revision, job_count }: DuplicateRow)}
      rev. {revision} • {shortDate(effective_date, true)} • {job_count} job{job_count === 1
        ? ""
        : "s"}
    {/snippet}

    {#snippet line1({ active, target, absorbs, revisions }: DuplicateRow)}
      <div class="flex flex-wrap gap-1">
        {#if active}
          <DsLabel tone="success" size="sm">Active</DsLabel>
        {:else}
          <DsLabel tone="neutral" size="sm">Inactive</DsLabel>
        {/if}
        {#if target}
          <DsLabel tone="info" size="sm">Target</DsLabel>
        {/if}
        {#if absorbs && absorbs > 1}
          <DsLabel tone="neutral" size="sm">Absorbs all {revisions} revisions</DsLabel>
        {/if}
      </div>
    {/snippet}

    {#snippet line2({ reason }: DuplicateRow)}
      {#if reason}
        Can't be absorbed into the target. {reason}
      {/if}
    {/snippet}

    {#snippet actions({ id, can_absorb }: DuplicateRow)}
      {#if canAbsorb && can_absorb}
        <DsActionButton
          action={`/rate-sheets/${id}/absorb`}
          icon="mdi:merge"
          title="Absorb identical sheets"
          color="yellow"
        />
      {/if}
    {/snippet}

    {#snippet groupFooter(group: string)}
      {#if groupReasons.get(group)}
        <div class="col-span-3 p-2">
          <DSCallout tone="warning">{groupReasons.get(group)}</DSCallout>
        </div>
      {/if}
    {/snippet}
  </DsList>
{/if}
