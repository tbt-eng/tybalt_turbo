<script module lang="ts">
  export interface InstantSearchHit {
    id: string;
    href: string;
    title: string;
    // Where the hit is, such as the page that holds a section.
    context?: string;
    // The text around the match. The parts with `match` show as marked.
    excerpt: { text: string; match: boolean }[];
  }
</script>

<script lang="ts">
  import Icon from "@iconify/svelte";
  import { flushSync, type Snippet } from "svelte";
  import DSActionButton from "./DSActionButton.svelte";

  // A search bar above one area of the app. While it holds a query, the hits
  // replace the content below it. `search` runs on each keystroke, so it must
  // answer from memory.
  let {
    label,
    search,
    emptyMessage = "Nothing matches. Try other words.",
    scroller = $bindable(),
    children,
  }: {
    // The placeholder, and the name that a screen reader says.
    label: string;
    search: (query: string) => InstantSearchHit[];
    emptyMessage?: string;
    // The element that scrolls the hits or the content.
    scroller?: HTMLDivElement;
    children: Snippet;
  } = $props();

  const minQueryLength = 2;

  let query = $state("");
  let active = $state(0);
  let list: HTMLUListElement | undefined = $state();

  const hits = $derived(query.trim().length >= minQueryLength ? search(query) : null);

  // The hits start at the top. The content returns to where it was.
  let searching = false;
  let contentScrollTop = 0;
  function showContent(restoreScroll: boolean) {
    searching = false;
    flushSync();
    if (restoreScroll && scroller) scroller.scrollTop = contentScrollTop;
  }

  function oninput(event: Event & { currentTarget: HTMLInputElement }) {
    if (!searching && scroller) contentScrollTop = scroller.scrollTop;
    query = event.currentTarget.value;
    active = 0;
    if (query.trim().length < minQueryLength) {
      if (searching) showContent(true);
      return;
    }
    searching = true;
    flushSync();
    if (scroller) scroller.scrollTop = 0;
  }

  // A hit's link needs the content back before the browser looks for its target.
  function clear(restoreScroll: boolean) {
    query = "";
    active = 0;
    showContent(restoreScroll);
  }

  const activeLink = () => list?.querySelector<HTMLAnchorElement>("a[data-active]");

  function onkeydown(event: KeyboardEvent) {
    if (event.key === "Escape") {
      clear(true);
    } else if (!hits?.length) {
      return;
    } else if (event.key === "ArrowDown" || event.key === "ArrowUp") {
      const step = event.key === "ArrowDown" ? 1 : -1;
      active = Math.min(Math.max(active + step, 0), hits.length - 1);
      flushSync();
      activeLink()?.scrollIntoView({ block: "nearest" });
    } else if (event.key === "Enter") {
      activeLink()?.click();
    } else {
      return;
    }
    event.preventDefault();
  }
</script>

<div class="flex min-h-0 min-w-0 flex-1 flex-col">
  <div role="search" data-print-hide class="border-b border-neutral-200 bg-white p-3">
    <div
      class="flex items-center gap-2 rounded-sm border border-neutral-400 bg-white px-3 py-2 focus-within:border-blue-700 focus-within:ring-1 focus-within:ring-blue-700"
    >
      <Icon icon="mdi:magnify" width="24px" class="shrink-0 text-neutral-500" />
      <input
        type="text"
        enterkeyhint="search"
        autocomplete="off"
        spellcheck="false"
        aria-label={label}
        placeholder={label}
        class="min-w-0 flex-1 bg-transparent text-lg outline-none placeholder:text-neutral-500"
        value={query}
        {oninput}
        {onkeydown}
      />
      {#if query}
        <DSActionButton variant="plain" action={() => clear(true)}>Clear</DSActionButton>
      {/if}
    </div>
  </div>
  <div bind:this={scroller} class="min-h-0 flex-1 overflow-y-auto">
    {#if hits}
      <div class="space-y-2 p-4">
        <p class="text-sm text-neutral-600" role="status">
          {#if hits.length === 0}
            {emptyMessage}
          {:else}
            {hits.length} {hits.length === 1 ? "result" : "results"}
          {/if}
        </p>
        <ul bind:this={list} class="space-y-1">
          {#each hits as hit, index (hit.id)}
            <li>
              <!-- eslint-disable svelte/no-navigation-without-resolve -- callers pass any href -->
              <a
                href={hit.href}
                data-active={index === active ? "" : undefined}
                onclick={() => clear(false)}
                onmouseenter={() => (active = index)}
                class={`block rounded-sm px-3 py-2 ${index === active ? "bg-neutral-200" : ""}`}
              >
                <span class="flex flex-wrap items-baseline gap-x-2">
                  <span class="font-semibold text-blue-700">{hit.title}</span>
                  {#if hit.context}
                    <span class="text-sm text-neutral-600">{hit.context}</span>
                  {/if}
                </span>
                <span class="block text-sm text-neutral-700">
                  {#each hit.excerpt as part, partIndex (partIndex)}{#if part.match}<mark
                        class="bg-yellow-200 text-inherit">{part.text}</mark
                      >{:else}{part.text}{/if}{/each}
                </span>
              </a>
              <!-- eslint-enable svelte/no-navigation-without-resolve -->
            </li>
          {/each}
        </ul>
      </div>
    {:else}
      {@render children()}
    {/if}
  </div>
</div>
