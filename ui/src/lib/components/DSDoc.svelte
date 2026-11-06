<script lang="ts">
  import DSFlow from "$lib/components/DSFlow.svelte";
  import type { DocBlock } from "$lib/docsFormat";

  // Renders parsed blocks from the user docs in src/docs.
  let { blocks }: { blocks: DocBlock[] } = $props();
</script>

<div class="space-y-2">
  {#each blocks as block, index (index)}
    {#if block.kind === "flow"}
      <DSFlow steps={block.steps} notes={block.notes} />
    {:else}
      <div
        class="space-y-2 [&_li]:mt-1 [&_ol]:list-decimal [&_ol]:space-y-1 [&_ol]:pl-5 [&_table]:w-full [&_table]:text-left [&_table]:text-sm [&_td]:border-t [&_td]:border-neutral-200 [&_td]:py-1 [&_td]:pr-3 [&_th]:py-1 [&_th]:pr-3 [&_th]:font-semibold [&_ul]:list-disc [&_ul]:space-y-1 [&_ul]:pl-5"
      >
        <!-- eslint-disable-next-line svelte/no-at-html-tags -- Repo-authored Markdown. -->
        {@html block.html}
      </div>
    {/if}
  {/each}
</div>
