<script lang="ts">
  import AbsorbEditor from "$lib/components/AbsorbEditor.svelte";
  import DSCallout from "$lib/components/DSCallout.svelte";
  import DSHelp from "$lib/components/DSHelp.svelte";
  import type { BaseSystemFields } from "$lib/pocketbase-types";
  import { shortDate } from "$lib/utilities";
  import type { RateSheetAbsorbCandidate } from "$lib/rateSheetAbsorb";
  import type { PageData } from "./$types";

  let { data }: { data: PageData } = $props();

  // AbsorbEditor takes PocketBase records; these rows come from a custom
  // route, so they carry only the id of the system fields.
  type AbsorbRow = RateSheetAbsorbCandidate & BaseSystemFields;

  // The editor finds the target among its records and offers the rest.
  const records = $derived([
    data.target,
    ...data.candidates.filter((c) => !c.reason),
  ] as unknown as AbsorbRow[]);
  const refused = $derived(data.candidates.filter((c) => c.reason));

  function describe(item: AbsorbRow): string {
    const parts = [
      `${item.name} rev. ${item.revision}`,
      shortDate(item.effective_date, true),
      `${item.job_count} job${item.job_count === 1 ? "" : "s"}`,
      item.active ? "Active" : "Inactive",
    ];
    if (item.absorbs && item.absorbs > 1) parts.push(`absorbs all ${item.revisions} revisions`);
    return parts.join(" • ");
  }
</script>

<div class="flex flex-col gap-2 p-4">
  <div class="flex justify-end">
    <DSHelp doc="rate-sheet-absorb#absorb-a-duplicate-rate-sheet" label="Help" />
  </div>
  {#if !data.target.active}
    <DSCallout tone="warning">
      This sheet is inactive. It can absorb only inactive sheets, and new projects can't select it.
    </DSCallout>
  {/if}
</div>

<AbsorbEditor
  collectionName="rate_sheets"
  targetRecordId={data.target.id}
  availableRecords={records}
>
  <!-- The editor renders this inside a select option too, so it is text only. -->
  {#snippet recordSnippet(item: AbsorbRow)}
    {describe(item)}
  {/snippet}
</AbsorbEditor>

{#if refused.length > 0}
  <div class="p-4">
    <DSCallout tone="info" title="Identical sheets that this sheet can't absorb">
      <ul class="flex flex-col gap-1">
        {#each refused as sheet (sheet.id)}
          <li>{sheet.name} rev. {sheet.revision}: {sheet.reason}</li>
        {/each}
      </ul>
    </DSCallout>
  </div>
{/if}
