<script lang="ts">
  import DsLabel from "$lib/components/DsLabel.svelte";
  import { formatCurrency } from "$lib/utilities";
  import { type ProposedRate, rateChanged, rateMissing } from "$lib/proposedRates";

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

  // A role with no rate is counted apart, not as a change.
  const missingCount = $derived(rows.filter(rateMissing).length);
  const changedCount = $derived(rows.filter((row) => !rateMissing(row) && rateChanged(row)).length);
</script>

<div class="flex flex-col gap-2">
  <p class="text-sm text-neutral-600">
    {changedCount === 0
      ? "No rate differs from the source sheet."
      : `${changedCount} ${changedCount === 1 ? "role differs" : "roles differ"} from ${sourceName || "the source sheet"}.`}
    {#if missingCount > 0}
      {`${missingCount} ${missingCount === 1 ? "role has" : "roles have"} no rate. Enter ${missingCount === 1 ? "its" : "their"} rates before you submit.`}
    {/if}
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
                {#if rateMissing(row)}<DsLabel tone="danger" size="sm">No rate</DsLabel>
                {:else if rateChanged(row)}<DsLabel tone="warning" size="sm">Changed</DsLabel>{/if}
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
