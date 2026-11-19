<script lang="ts">
  import DsLabel from "$lib/components/DsLabel.svelte";
  import { formatCurrency } from "$lib/utilities";
  import {
    type ProposedRate,
    type RateField,
    rateChanged,
    missingRateLabel,
    rateErrorKey,
    rateMissing,
  } from "$lib/proposedRates";
  import OvertimeMultiplier from "./OvertimeMultiplier.svelte";

  // The rates of a proposed rate sheet beside the rates of the sheet it
  // started from, with each changed role marked. With editable, the rate and
  // the overtime rate of each role are inputs, with a multiplier above the
  // overtime rates. errors holds the errors of the last save under
  // rateErrorKey, and each shows under the rate it belongs to. Read only, a
  // role with no rate shows no amount, and the page says what to do.
  // With find, the table is the input of Find by rates: a rate and an
  // optional overtime rate for each role, with no source rates, no summary,
  // no multiplier and no labels. find implies editable. With find, Enter in a
  // rate does not submit the form around the table; it calls onEnter, such as
  // a search at once.
  let {
    rows = $bindable(),
    editable = false,
    find = false,
    sourceName = "",
    errors = {},
    onEnter,
  }: {
    rows: ProposedRate[];
    editable?: boolean;
    find?: boolean;
    sourceName?: string;
    errors?: Record<string, { message: string }>;
    onEnter?: () => void;
  } = $props();

  function findKeydown(event: KeyboardEvent) {
    if (event.key !== "Enter") return;
    event.preventDefault();
    onEnter?.();
  }

  const inputs = $derived(editable || find);

  const errorOf = (row: ProposedRate, field: RateField) =>
    errors[rateErrorKey(row.role, field)]?.message;

  // A role that lacks either rate is counted apart, not as a change.
  const missingCount = $derived(rows.filter(rateMissing).length);
  const changedCount = $derived(rows.filter((row) => !rateMissing(row) && rateChanged(row)).length);
</script>

<div class="flex flex-col gap-2">
  {#if !find}<p class="text-sm text-neutral-600">
      {changedCount === 0
        ? "No rate differs from the source sheet."
        : `${changedCount} ${changedCount === 1 ? "role differs" : "roles differ"} from ${sourceName || "the source sheet"}.`}
      {#if missingCount > 0 && editable}
        {`${missingCount} ${missingCount === 1 ? "role is missing a rate" : "roles are missing rates"}. Enter ${missingCount === 1 ? "its" : "their"} rates before you submit.`}
      {/if}
    </p>{/if}
  <div class="overflow-x-auto">
    <table class="w-full border-collapse text-sm">
      <thead>
        {#if editable && !find}
          <tr>
            <th></th>
            <th></th>
            <th class="px-2 py-1 text-left"><OvertimeMultiplier bind:rows /></th>
            <th></th>
          </tr>
        {/if}
        <tr class="border-b border-neutral-300 text-left">
          <th class="px-2 py-1">Role</th>
          <th class="px-2 py-1">Rate</th>
          <th class="px-2 py-1">Overtime Rate</th>
          {#if !find}<th class="px-2 py-1">Source rates</th>{/if}
        </tr>
      </thead>
      <tbody>
        {#each rows as row, i (row.role)}
          <tr class="border-b border-neutral-200">
            <td class="px-2 py-1">
              <span class="flex items-center gap-2">
                {row.roleName}
                {#if !find}
                  {#if rateMissing(row)}<DsLabel tone="danger" size="sm"
                      >{missingRateLabel(row)}</DsLabel
                    >
                  {:else if rateChanged(row)}<DsLabel tone="warning" size="sm">Changed</DsLabel
                    >{/if}
                {/if}
              </span>
              {#if errorOf(row, "role")}
                <div class="mt-1 text-xs text-red-600">{errorOf(row, "role")}</div>
              {/if}
            </td>
            <td class="px-2 py-1">
              {#if inputs}
                <input
                  type="number"
                  min="1"
                  step="1"
                  aria-label={`Rate for ${row.roleName}`}
                  aria-invalid={errorOf(row, "rate") ? true : undefined}
                  class="{find ? 'w-20' : 'w-24'} rounded-sm border px-1 {errorOf(row, 'rate')
                    ? 'border-red-500'
                    : 'border-neutral-300'}"
                  onkeydown={find ? findKeydown : undefined}
                  bind:value={rows[i].rate}
                />
                {#if errorOf(row, "rate")}
                  <div class="mt-1 text-xs text-red-600">{errorOf(row, "rate")}</div>
                {/if}
              {:else if row.rate !== null}
                {formatCurrency(row.rate)}
              {/if}
            </td>
            <td class="px-2 py-1">
              {#if inputs}
                <input
                  type="number"
                  min="1"
                  step="0.01"
                  placeholder={find ? "Optional" : undefined}
                  aria-label={`Overtime rate for ${row.roleName}`}
                  aria-invalid={errorOf(row, "overtime_rate") ? true : undefined}
                  class="{find ? 'w-20' : 'w-24'} rounded-sm border px-1 {errorOf(
                    row,
                    'overtime_rate',
                  )
                    ? 'border-red-500'
                    : 'border-neutral-300'}"
                  onkeydown={find ? findKeydown : undefined}
                  bind:value={rows[i].overtime_rate}
                />
                {#if errorOf(row, "overtime_rate")}
                  <div class="mt-1 text-xs text-red-600">{errorOf(row, "overtime_rate")}</div>
                {/if}
              {:else if row.overtime_rate !== null}
                {formatCurrency(row.overtime_rate)}
              {/if}
            </td>
            {#if !find}
              <td class="px-2 py-1 text-neutral-600">
                {#if row.source_rate !== undefined}
                  {formatCurrency(row.source_rate)} / {formatCurrency(row.source_overtime_rate)}
                {/if}
              </td>
            {/if}
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
</div>
