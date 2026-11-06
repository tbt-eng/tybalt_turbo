<script lang="ts">
  // A second navigation column for one area of the app, beside the main menu.
  // On a phone the links wrap above the content.
  export interface SideNavItem {
    label: string;
    href: string;
    active?: boolean;
  }
  let { title, items }: { title: string; items: SideNavItem[] } = $props();
</script>

<nav
  data-print-hide
  aria-label={title}
  class="shrink-0 border-b border-neutral-200 bg-neutral-100 md:w-48 md:overflow-y-auto md:border-r md:border-b-0"
>
  <h2 class="px-3 pt-4 pb-2 text-lg font-semibold">{title}</h2>
  <ul class="flex flex-wrap gap-1 px-1 pb-2 md:flex-col md:flex-nowrap">
    {#each items as item (item.href)}
      <li>
        <!-- eslint-disable svelte/no-navigation-without-resolve -- callers pass any href -->
        <a
          href={item.href}
          aria-current={item.active ? "page" : undefined}
          class={`block rounded-sm px-2 py-1 ${
            item.active
              ? "bg-white font-semibold text-blue-700 shadow-sm"
              : "text-neutral-700 hover:bg-neutral-200"
          }`}
        >
          {item.label}
        </a>
        <!-- eslint-enable svelte/no-navigation-without-resolve -->
      </li>
    {/each}
  </ul>
</nav>
