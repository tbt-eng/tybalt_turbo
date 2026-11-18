<script lang="ts">
  import type { SearchResult } from "minisearch";
  import RateSheetMatchLine from "./RateSheetMatchLine.svelte";
  import { rateSheetLabel, rateSheetPickerMatch } from "$lib/rateSheetLabel";
  import type { RateSheetResponse } from "$lib/stores/rateSheets";
  import { shortDate } from "$lib/utilities";

  // A rate sheet as a DSPicker option: its label and effective date. In a
  // search result, a line under it names the jobs and clients that the search
  // matched (RateSheetMatchLine).
  let { sheet, result }: { sheet: RateSheetResponse; result: SearchResult | undefined } = $props();
</script>

{rateSheetLabel(sheet)} · {shortDate(sheet.effective_date, true)}
<RateSheetMatchLine {sheet} matches={rateSheetPickerMatch(sheet, result)} />
