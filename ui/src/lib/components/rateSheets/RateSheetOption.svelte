<script lang="ts">
  import type { SearchResult } from "minisearch";
  import { rateSheetLabel, rateSheetMatchLine, rateSheetPickerMatch } from "$lib/rateSheetLabel";
  import type { RateSheetResponse } from "$lib/stores/rateSheets";
  import { shortDate } from "$lib/utilities";

  // A rate sheet as a DSPicker option: its label and effective date. In a
  // search result, a line under it names the jobs and clients that the search
  // matched (rateSheetMatchLine).
  let { sheet, result }: { sheet: RateSheetResponse; result: SearchResult | undefined } = $props();

  const line = $derived(rateSheetMatchLine(sheet, rateSheetPickerMatch(sheet, result)));
</script>

{rateSheetLabel(sheet)} · {shortDate(sheet.effective_date, true)}
{#if line}<span class="block text-sm text-neutral-600">{line}</span>{/if}
