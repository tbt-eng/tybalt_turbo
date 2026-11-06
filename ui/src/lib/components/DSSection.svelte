<script lang="ts">
  import type { Snippet } from "svelte";

  // A page-level section: an h2 heading with optional actions, and no box.
  let {
    title,
    titleNote,
    actions,
    children,
  }: {
    title: string;
    // Small text directly after the heading.
    titleNote?: Snippet;
    actions?: Snippet;
    children: Snippet;
  } = $props();
</script>

<section class="space-y-2">
  <div class="flex flex-wrap items-center justify-between gap-2">
    <div class="flex flex-wrap items-baseline gap-2">
      <h2 class="text-lg font-semibold">{title}</h2>
      {#if titleNote}
        <span class="text-sm italic">{@render titleNote()}</span>
      {/if}
    </div>
    {#if actions}
      <div class="flex flex-wrap items-center gap-2">{@render actions()}</div>
    {/if}
  </div>
  {@render children()}
</section>
