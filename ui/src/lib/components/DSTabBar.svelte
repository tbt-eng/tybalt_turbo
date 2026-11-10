<script lang="ts">
  import type { Snippet } from "svelte";
  export interface TabItem {
    label: string;
    href: string;
    active?: boolean;
    // For a tab that changes the view and not the page. A click calls this and
    // does not follow href.
    onselect?: () => void;
  }
  let { tabs, rightContent }: { tabs: TabItem[]; rightContent?: Snippet } = $props();
</script>

<div class="flex flex-wrap items-center gap-x-4 gap-y-1 rounded-t bg-neutral-100 px-2 py-1">
  {#each tabs as t (t.href + t.label)}
    <!-- eslint-disable svelte/no-navigation-without-resolve -- callers pass any href -->
    <a
      href={t.href}
      aria-current={t.active ? "page" : undefined}
      onclick={(event) => {
        if (!t.onselect) return;
        event.preventDefault();
        t.onselect();
      }}
      class={`rounded-t px-3 py-1 ${
        t.active
          ? "bg-white font-semibold text-blue-700 shadow-sm"
          : "text-neutral-700 hover:bg-neutral-200"
      }`}
    >
      {t.label}
    </a>
    <!-- eslint-enable svelte/no-navigation-without-resolve -->
  {/each}
  {#if rightContent}
    <div class="ml-auto">
      {@render rightContent()}
    </div>
  {/if}
</div>
