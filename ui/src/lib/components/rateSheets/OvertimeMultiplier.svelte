<script lang="ts" generics="T extends { rate: number | null; overtime_rate: number | null }">
  import DsActionButton from "$lib/components/DSActionButton.svelte";
  import { withOvertimeMultiplier } from "$lib/proposedRates";

  // A multiplier and a button that set every overtime rate in rows to the rate
  // times the multiplier. It sits above the Overtime Rate column of an
  // editable rates table. A row with no rate stays as it is.
  let { rows = $bindable() }: { rows: T[] } = $props();

  // An empty input binds null.
  let value = $state<number | null>(1.3);
</script>

<div class="flex items-center gap-1">
  <input
    type="number"
    min="0.01"
    step="0.01"
    bind:value
    class="w-16 rounded-sm border border-neutral-300 px-2 py-1 text-right text-sm"
    title="Overtime multiplier"
    aria-label="Overtime multiplier"
  />
  <DsActionButton
    color="neutral"
    title="Apply multiplier to all overtime rates"
    action={() => (rows = withOvertimeMultiplier(rows, value))}
  >
    &times;
  </DsActionButton>
</div>
