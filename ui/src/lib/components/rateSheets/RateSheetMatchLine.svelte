<script lang="ts">
  import {
    rateSheetMatchGroups,
    type RateSheetEntryMatch,
    type RateSheetMatchGroup,
  } from "$lib/rateSheetLabel";

  // Why a search found a rate sheet: "Jobs 91-8901 · 91-8902 · 91-8903 +5
  // more", then "Client Lakehead Surveys, Inc." (rateSheetMatchGroups). Each
  // label and "+N more" is muted, so the job numbers and client names stand
  // out. The groups wrap, with a gap between them. It renders nothing when no
  // job number or client name matched. The Rate Sheets list and the rate sheet
  // pickers (RateSheetOption.svelte) both show it.
  let {
    sheet,
    matches,
  }: {
    sheet: { job_numbers?: string | null; client_names?: string | null };
    // The rule of the search that found the sheet: rateSheetListMatch or
    // rateSheetPickerMatch.
    matches: RateSheetEntryMatch;
  } = $props();

  const groups: RateSheetMatchGroup[] = $derived(rateSheetMatchGroups(sheet, matches));
</script>

{#if groups.length > 0}
  <span class="flex flex-wrap gap-x-4 text-sm">
    {#each groups as group (group.label)}
      <span>
        <span class="text-neutral-500">{group.label}</span>
        <span class="text-neutral-800">{group.shown.join(" · ")}</span>
        {#if group.more > 0}
          <span class="whitespace-nowrap text-neutral-500">+{group.more} more</span>
        {/if}
      </span>
    {/each}
  </span>
{/if}
