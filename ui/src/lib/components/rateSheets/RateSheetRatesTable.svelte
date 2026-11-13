<script lang="ts" module>
  // One role's rates on a proposed rate sheet, with the rates of the sheet
  // that the proposal started from when they are known.
  export interface ProposedRate {
    role: string;
    roleName: string;
    rate: number;
    overtime_rate: number;
    source_rate?: number;
    source_overtime_rate?: number;
  }

  // A rate is changed when it differs from the source sheet's, or when the
  // source sheet has no rate for the role.
  export function rateChanged(row: ProposedRate): boolean {
    return (
      row.source_rate === undefined ||
      Number(row.rate) !== row.source_rate ||
      Number(row.overtime_rate) !== row.source_overtime_rate
    );
  }
</script>

<script lang="ts">
  import DsLabel from "$lib/components/DsLabel.svelte";
  import { formatCurrency } from "$lib/utilities";

  // The rates of a proposed rate sheet beside the rates of the sheet it
  // started from, with each changed role marked. With editable, the rate and
  // the overtime rate of each role are inputs.
  let {
    rows = $bindable(),
    editable = false,
    sourceName = "",
  }: {
    rows: ProposedRate[];
    editable?: boolean;
    sourceName?: string;
  } = $props();

  const changedCount = $derived(rows.filter(rateChanged).length);
</script>

<div class="flex flex-col gap-2">
  <p class="text-sm text-neutral-600">
    {changedCount === 0
      ? "No rate differs from the source sheet."
      : `${changedCount} role${changedCount === 1 ? "" : "s"} differ from ${sourceName || "the source sheet"}.`}
  </p>
  <div class="overflow-x-auto">
    <table class="w-full border-collapse text-sm">
      <thead>
        <tr class="border-b border-neutral-300 text-left">
          <th class="px-2 py-1">Role</th>
          <th class="px-2 py-1">Rate</th>
          <th class="px-2 py-1">Overtime Rate</th>
          <th class="px-2 py-1">Source rates</th>
        </tr>
      </thead>
      <tbody>
        {#each rows as row, i (row.role)}
          <tr class="border-b border-neutral-200">
            <td class="px-2 py-1">
              <span class="flex items-center gap-2">
                {row.roleName}
                {#if rateChanged(row)}<DsLabel tone="warning" size="sm">Changed</DsLabel>{/if}
              </span>
            </td>
            <td class="px-2 py-1">
              {#if editable}
                <input
                  type="number"
                  min="1"
                  step="1"
                  aria-label={`Rate for ${row.roleName}`}
                  class="w-24 rounded-sm border border-neutral-300 px-1"
                  bind:value={rows[i].rate}
                />
              {:else}
                {formatCurrency(row.rate)}
              {/if}
            </td>
            <td class="px-2 py-1">
              {#if editable}
                <input
                  type="number"
                  min="1"
                  step="0.01"
                  aria-label={`Overtime rate for ${row.roleName}`}
                  class="w-24 rounded-sm border border-neutral-300 px-1"
                  bind:value={rows[i].overtime_rate}
                />
              {:else}
                {formatCurrency(row.overtime_rate)}
              {/if}
            </td>
            <td class="px-2 py-1 text-neutral-600">
              {#if row.source_rate !== undefined}
                {formatCurrency(row.source_rate)} / {formatCurrency(row.source_overtime_rate)}
              {/if}
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
</div>
