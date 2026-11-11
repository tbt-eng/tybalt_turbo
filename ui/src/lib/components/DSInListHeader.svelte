<script lang="ts">
  import type { Snippet } from "svelte";
  let {
    value,
    snippet,
    extra,
  }: {
    value: string;
    snippet?: Snippet<[string]>;
    // Sits at the far right of the header, such as page help on a list with
    // no search bar or tab bar.
    extra?: Snippet;
  } = $props();
</script>

<li
  class="inlistheader col-span-3 border-b border-blue-300 bg-blue-200 px-2 text-neutral-700 {extra
    ? 'flex flex-wrap items-center justify-between gap-x-2'
    : ''}"
>
  {#if snippet !== undefined}
    {@render snippet(value)}
  {:else}
    {value}
  {/if}
  {#if extra !== undefined}
    {@render extra()}
  {/if}
</li>
