<script lang="ts">
  import DSHelp from "$lib/components/DSHelp.svelte";
  import RateSheetLabel from "$lib/components/rateSheets/RateSheetLabel.svelte";
  import type { RateSheet } from "./timeSummary";
  let { rateSheet }: { rateSheet?: RateSheet } = $props();
</script>

<DSHelp title="How time values are calculated">
  {#if rateSheet?.id}
    <p>
      Current rate sheet: <RateSheetLabel sheet={rateSheet} />.
    </p>
  {/if}
  <p>
    Each entry uses the regular rate for its role from the rate sheet in effect on its date, as set
    in the job's rate sheet history. If no job rate matches, a staff default charge-out rate greater
    than zero is used.
  </p>
  <p>
    Estimated amounts include staff defaults. Partial amounts exclude hours that cannot be priced.
    If no hours can be priced, the value is shown as —.
  </p>
  <p>
    Includes committed and uncommitted entries in the selected dates. Excludes meal hours, time
    amendments, and overtime billing.
  </p>
  <p>
    Shares use the full value for the job and date range. They are unavailable when that value is
    incomplete or zero.
  </p>
  <p>
    Click column headings to sort and table values to filter. The total follows the visible rows;
    the CSV always includes the full date range and all pricing details.
  </p>
  <p>
    Changes to the rate sheet history, a sheet's rates, or staff defaults recalculate past values.
  </p>
</DSHelp>
