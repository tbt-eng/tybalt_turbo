<script lang="ts">
  import type { Snippet } from "svelte";
  import { untrack } from "svelte";

  let {
    title,
    collapsed = true,
    subtle = false,
    children,
    headerActions,
  }: {
    title: string;
    collapsed?: boolean;
    // A small grey title, for a detail inside a list row.
    subtle?: boolean;
    children?: Snippet<[]>;
    headerActions?: Snippet<[boolean]>;
  } = $props();

  let isOpen = $state(untrack(() => !collapsed));
  let isCollapsed = $derived(!isOpen);

  $effect(() => {
    isOpen = !collapsed;
  });
</script>

<details bind:open={isOpen} class="space-y-2">
  <summary
    class={subtle
      ? "cursor-pointer text-xs text-neutral-600 hover:text-neutral-900"
      : "cursor-pointer font-semibold text-neutral-700"}>{title}</summary
  >
  {#if !isCollapsed}
    {#if headerActions}
      <div class="flex min-h-9 items-center justify-end">
        {@render headerActions?.(isCollapsed)}
      </div>
    {/if}
    <div class="space-y-2">
      {@render children?.()}
    </div>
  {/if}
</details>
