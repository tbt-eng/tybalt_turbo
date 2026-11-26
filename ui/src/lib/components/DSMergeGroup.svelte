<script lang="ts" module>
  // One record in a DSMergeGroup.
  export type MergeRow = {
    id: string;
    label: string;
    // Muted text after the label, such as an alias.
    detail?: string;
    // A number in the right column, such as how many records use it.
    count?: number;
    // Opens the record in a new tab, so the group keeps its state.
    href?: string;
    // Why the row can't be picked; shown as a warning label.
    locked?: string;
  };
</script>

<script lang="ts">
  import Icon from "@iconify/svelte";
  import type { Snippet } from "svelte";
  import DsLabel from "./DsLabel.svelte";

  // A dense group of records where one is kept (a radio) and others are
  // ticked to merge into it (a checkbox on every other row). The whole row,
  // past the radio, is the click target for its tick. Choosing another record
  // to keep clears its tick, and the old one becomes an unticked row. The
  // absorb suggestions page uses it for each bucket of duplicates.
  //
  // Keyboard: a row takes focus (the first row of a group is a tab stop).
  // Up and Down move to the row above or below, and on into the next or the
  // previous group on the page. Space ticks or unticks the row.
  let {
    id,
    title,
    subtitle = "",
    rows,
    keepId = $bindable(),
    picked = $bindable(),
    suggestedId,
    keepLabel = "Keep",
    pickLabel = "Absorb",
    countLabel = "",
    disabled = false,
    actions,
    notice,
  }: {
    // The DOM id of the group, which takes focus (tabindex -1).
    id: string;
    title: string;
    subtitle?: string;
    rows: MergeRow[];
    keepId: string;
    picked: string[];
    // Marked "Suggested" in its row.
    suggestedId?: string;
    keepLabel?: string;
    pickLabel?: string;
    countLabel?: string;
    // No record can be kept or ticked.
    disabled?: boolean;
    // Buttons, at the right of the title row.
    actions?: Snippet;
    // A callout under the rows, such as an error.
    notice?: Snippet;
  } = $props();

  const titleId = $derived(`${id}-title`);

  function keep(rowId: string) {
    keepId = rowId;
    picked = picked.filter((pickedId) => pickedId !== rowId);
  }

  function toggle(rowId: string, on: boolean) {
    picked = on ? [...picked.filter((p) => p !== rowId), rowId] : picked.filter((p) => p !== rowId);
  }

  // Every row of every group on the page, in order, so the arrow keys go
  // from one group into the next.
  function moveFocus(from: HTMLElement, step: 1 | -1) {
    const all = Array.from(document.querySelectorAll<HTMLElement>("[data-merge-row]"));
    all[all.indexOf(from) + step]?.focus();
  }

  function onRowKeydown(event: KeyboardEvent, row: MergeRow, canPick: boolean) {
    if (event.altKey || event.ctrlKey || event.metaKey) return;
    const rowElement = event.currentTarget as HTMLElement;
    if (event.key === "ArrowDown" || event.key === "ArrowUp") {
      // Also stops the arrows from changing the kept record when a radio
      // has focus.
      event.preventDefault();
      moveFocus(rowElement, event.key === "ArrowDown" ? 1 : -1);
      return;
    }
    // A control in the row (the radio, the checkbox, the link) handles its
    // own Space.
    if (event.key === " " && event.target === rowElement) {
      event.preventDefault();
      if (canPick) toggle(row.id, !picked.includes(row.id));
    }
  }

  function rowClass(row: MergeRow): string {
    if (row.id === keepId) return "bg-white font-semibold";
    if (picked.includes(row.id)) return "bg-amber-100";
    return "";
  }
</script>

<section
  {id}
  tabindex="-1"
  aria-labelledby={titleId}
  class="rounded-sm bg-neutral-100 p-2 outline-blue-500 focus:outline-2"
>
  <div class="flex flex-wrap items-center gap-x-3 gap-y-1 pb-1">
    <h3 id={titleId} class="font-semibold">{title}</h3>
    {#if subtitle}<span class="text-sm text-neutral-600">{subtitle}</span>{/if}
    {#if actions}
      <div class="ml-auto flex flex-wrap items-center gap-2">{@render actions()}</div>
    {/if}
  </div>
  <table class="w-full border-collapse text-sm">
    <thead>
      <tr class="text-left text-xs text-neutral-600">
        <th class="w-12 px-2 py-0.5 font-normal">{keepLabel}</th>
        <th class="w-14 px-2 py-0.5 font-normal">{pickLabel}</th>
        <th class="px-2 py-0.5 font-normal"><span class="sr-only">Record</span></th>
        <th class="w-24 px-2 py-0.5 text-right font-normal">{countLabel}</th>
      </tr>
    </thead>
    <tbody>
      {#each rows as row, index (row.id)}
        {@const checkId = `${id}-pick-${row.id}`}
        {@const isKept = row.id === keepId}
        {@const canPick = !isKept && !row.locked && !disabled}
        <tr
          data-merge-row
          tabindex={index === 0 ? 0 : -1}
          aria-label={`${row.label}${isKept ? ", kept" : picked.includes(row.id) ? ", ticked" : ""}`}
          class="border-t border-neutral-200 -outline-offset-2 outline-blue-500 focus-within:outline-2 {rowClass(
            row,
          )}"
          onkeydown={(event) => onRowKeydown(event, row, canPick)}
        >
          <td class="px-2">
            <input
              type="radio"
              class="h-4 w-4 align-middle"
              name={`${id}-keep`}
              checked={isKept}
              {disabled}
              aria-label={`${keepLabel} ${row.label}`}
              onchange={() => keep(row.id)}
            />
          </td>
          <td class="px-2">
            {#if !isKept}
              <input
                type="checkbox"
                class="h-4 w-4 align-middle {canPick ? '' : 'cursor-not-allowed opacity-50'}"
                id={checkId}
                aria-label={`${pickLabel} ${row.label}`}
                checked={picked.includes(row.id)}
                disabled={!canPick}
                onchange={(event) => toggle(row.id, event.currentTarget.checked)}
              />
            {/if}
          </td>
          <td class="p-0">
            <label
              for={isKept ? undefined : checkId}
              class="flex flex-wrap items-center gap-x-2 px-2 py-1 {canPick
                ? 'cursor-pointer'
                : ''}"
            >
              <span>{row.label}</span>
              {#if row.detail}<span class="font-normal text-neutral-600">{row.detail}</span>{/if}
              {#if row.href}
                <!-- eslint-disable svelte/no-navigation-without-resolve -- callers pass any href -->
                <a
                  class="link"
                  href={row.href}
                  target="_blank"
                  rel="noopener"
                  title={`Open ${row.label} in a new tab`}
                  aria-label={`Open ${row.label} in a new tab`}
                >
                  <Icon icon="mdi:open-in-new" width="14px" />
                </a>
                <!-- eslint-enable svelte/no-navigation-without-resolve -->
              {/if}
              {#if isKept}<DsLabel tone="info" size="sm">Kept</DsLabel>{/if}
              {#if row.id === suggestedId && !isKept}
                <DsLabel tone="neutral" size="sm">Suggested</DsLabel>
              {/if}
              {#if row.locked}<DsLabel tone="warning" size="sm">{row.locked}</DsLabel>{/if}
            </label>
          </td>
          <td class="p-0 text-right">
            <label
              for={isKept ? undefined : checkId}
              class="block px-2 py-1 tabular-nums {canPick ? 'cursor-pointer' : ''}"
            >
              {row.count ?? ""}
            </label>
          </td>
        </tr>
      {/each}
    </tbody>
  </table>
  {#if notice}
    <div class="pt-2">{@render notice()}</div>
  {/if}
</section>
