<script module lang="ts">
  // A list of this many items or fewer is a dropdown; a longer one is a search.
  export const pickerDropdownLimit = 10;
</script>

<script lang="ts" generics="T extends { id: string }">
  import type { Snippet } from "svelte";
  import type MiniSearch from "minisearch";
  import type { SearchResult } from "minisearch";
  import DsSelector from "./DSSelector.svelte";
  import DsAutoComplete from "./DSAutoComplete.svelte";

  // Picks one record from a list: a DSSelector for a short list, a
  // DSAutoComplete for a long one. items is the list the user may choose from,
  // already filtered; index is the search index that holds them, needed for a
  // long list only, and the search offers only the items. option renders a
  // record in either mode. With empty, the dropdown starts with a
  // "-- select --" choice so the value can be cleared, as the search's Clear
  // button does. A value that is not in items, such as a sheet deactivated
  // since it was chosen, still shows as the current choice.
  let {
    value = $bindable(),
    items,
    index = null,
    option,
    errors,
    fieldName,
    uiName,
    disabled = false,
    empty = false,
    help,
  }: {
    value: string;
    items: T[];
    index?: MiniSearch<T> | null;
    option: Snippet<[T]>;
    errors: Record<string, { message: string }>;
    fieldName: string;
    uiName: string;
    disabled?: boolean;
    empty?: boolean;
    help?: Snippet;
  } = $props();

  const byId = $derived(new Map(items.map((item) => [item.id, item])));

  // The current choice when it is not in items, from the search index.
  const kept = $derived.by(() => {
    if (!value || byId.has(value)) return [];
    const stored = index?.getStoredFields(value);
    return stored ? [{ ...stored, id: value } as unknown as T] : [];
  });

  const choices = $derived([...(empty ? [{ id: "" } as T] : []), ...kept, ...items]);

  const offered = (result: T) => byId.has(result.id);
  const record = (result: SearchResult) => byId.get(result.id) ?? (result as unknown as T);
</script>

{#if items.length <= pickerDropdownLimit}
  <DsSelector bind:value items={choices} {errors} {fieldName} {uiName} {disabled} {help}>
    {#snippet optionTemplate(item)}
      {#if item.id === ""}-- select --{:else}{@render option(item)}{/if}
    {/snippet}
  </DsSelector>
{:else if index !== null}
  <DsAutoComplete
    bind:value
    {index}
    filter={offered}
    {errors}
    {fieldName}
    {uiName}
    {disabled}
    {help}
  >
    {#snippet resultTemplate(result)}{@render option(record(result))}{/snippet}
  </DsAutoComplete>
{/if}
