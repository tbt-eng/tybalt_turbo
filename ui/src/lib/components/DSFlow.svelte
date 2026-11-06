<script lang="ts">
  import Icon from "@iconify/svelte";
  import type { FlowStep } from "$lib/docsFormat";

  // Cards that show who acts and what happens next. Docs write these as a
  // `flow` block; see src/docs/README.md. Step details and notes are HTML from
  // the repo's own Markdown. The layout follows the space available, because a
  // help popover is narrow on any screen.
  let { steps, notes = [] }: { steps: FlowStep[]; notes?: string[] } = $props();
</script>

<div class="@container space-y-2">
  <ol class="flex flex-col gap-1 @2xl:flex-row @2xl:items-stretch">
    {#each steps as step, index (index)}
      <li class="contents">
        <div
          class="rounded-sm border border-l-4 border-neutral-200 border-l-blue-700 p-3 @2xl:flex-1"
        >
          <div class="text-xs font-semibold text-blue-700 uppercase">{step.actor}</div>
          <div class="font-semibold">{step.step}</div>
          <!-- eslint-disable-next-line svelte/no-at-html-tags -- Repo-authored Markdown. -->
          <div class="text-sm text-neutral-600">{@html step.detail}</div>
        </div>
        {#if index < steps.length - 1}
          <div class="flex shrink-0 justify-center self-center text-neutral-400" aria-hidden="true">
            <Icon icon="mdi:arrow-down" width="16" class="@2xl:hidden" />
            <Icon icon="mdi:arrow-right" width="16" class="hidden @2xl:block" />
          </div>
        {/if}
      </li>
    {/each}
  </ol>
  {#each notes as note, index (index)}
    <!-- eslint-disable-next-line svelte/no-at-html-tags -- Repo-authored Markdown. -->
    <p class="text-sm text-neutral-600">{@html note}</p>
  {/each}
</div>
