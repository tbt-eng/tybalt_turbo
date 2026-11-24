<script module>
  let idCounter = $state(0);
</script>

<script lang="ts" generics="T extends SelectOption">
  import type { Snippet } from "svelte";
  import type { SelectOption } from "$lib/pocketbase-types";
  import DsActionButton from "./DSActionButton.svelte";
  import DSRequiredMark from "./DSRequiredMark.svelte";
  // get an id for this instance from the counter in the module context then
  // increment it so the next instance gets a different id
  const thisId = idCounter;
  idCounter += 1;

  let {
    value = $bindable(),
    items = [] as T[],
    errors,
    fieldName,
    uiName,
    clear = false,
    optionTemplate,
    trailing,
    disabled = false,
    fill = false,
    helperText = "",
    help,
    required = false,
  }: {
    value: string | number;
    items: T[];
    errors: Record<string, { message: string }>;
    fieldName: string;
    uiName: string;
    clear?: boolean;
    optionTemplate: Snippet<[T]>;
    trailing?: Snippet<[]>;
    disabled?: boolean;
    // The dropdown takes the rest of the row, not only the width of its text.
    fill?: boolean;
    // A line under the field explaining it.
    helperText?: string;
    // Field help (DSHelp) at the far right of the field row.
    help?: Snippet;
    // A red asterisk after the label; the server still validates.
    required?: boolean;
  } = $props();

  function clearValue() {
    value = "";
  }
</script>

<div
  class="flex w-full gap-2 {helperText ? 'flex-col' : ''} {errors[fieldName] !== undefined
    ? 'bg-red-200'
    : ''}"
>
  <span class="flex w-full gap-2" class:items-center={help !== undefined}>
    <label for={`select-input-${thisId}`}
      >{uiName}{#if required}<DSRequiredMark />{/if}</label
    >
    <select
      id={`select-input-${thisId}`}
      name={fieldName}
      bind:value
      class="max-w-full min-w-0 rounded-sm border border-neutral-300 px-1 {fill
        ? 'flex-1'
        : ''} {disabled ? 'opacity-50' : ''}"
      aria-required={required || undefined}
      {disabled}
    >
      {#each items as item (item.id)}
        <option value={item.id}>{@render optionTemplate(item)}</option>
      {/each}
    </select>
    {#if trailing !== undefined}
      {@render trailing()}
    {/if}
    {#if clear === true && value !== undefined && value !== "" && !disabled}
      <DsActionButton action={clearValue} title="Clear">Clear</DsActionButton>
    {/if}
    {#if help}
      <span class="ml-auto shrink-0">{@render help()}</span>
    {/if}
  </span>
  {#if errors[fieldName] !== undefined}
    <span class="text-red-600">{errors[fieldName].message}</span>
  {/if}
  {#if helperText}
    <div class="text-sm text-neutral-600">{helperText}</div>
  {/if}
</div>
