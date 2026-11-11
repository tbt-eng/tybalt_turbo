<script lang="ts">
  import { flushSync, onDestroy } from "svelte";
  import { afterNavigate, goto } from "$app/navigation";
  import { resolve } from "$app/paths";
  import { page } from "$app/state";
  import DSChat, { type ChatTurnView } from "$lib/components/DSChat.svelte";
  import DSInstantSearch from "$lib/components/DSInstantSearch.svelte";
  import DSSideNav from "$lib/components/DSSideNav.svelte";
  import DSTabBar from "$lib/components/DSTabBar.svelte";
  import { searchDocs } from "$lib/docs";
  import { answerText, sourceLinks } from "$lib/helpChatFormat";
  import { helpChatEnabled } from "$lib/stores/appConfig";
  import { helpBubble } from "$lib/stores/helpBubble";
  import { helpChat, helpChatBusy } from "$lib/stores/helpChat";
  let { data, children } = $props();

  const currentSlug = $derived(page.params.slug ?? "start-here");
  const items = $derived(
    data.pages.map(({ slug, label, indent }) => ({
      label,
      indent,
      href: resolve(`/docs/${slug}`),
      active: slug === currentSlug,
      dotLabel: slug === "whats-new" && $helpBubble ? "New in Help" : undefined,
    })),
  );

  // A hit opens its section, or the top of the page for a page's intro.
  const search = (query: string) =>
    searchDocs(query).map((hit) => ({
      id: hit.id,
      href: resolve(`/docs/${hit.slug}`) + (hit.sectionId ? `#${hit.sectionId}` : ""),
      title: hit.section || hit.page,
      context: hit.section ? hit.page : undefined,
      excerpt: hit.excerpt,
    }));

  // Ask: a question goes to the server, and the conversation shows in place of
  // the docs page until the user closes it or opens a page.
  const docTitles = $derived(
    Object.fromEntries(data.pages.map(({ slug, label }) => [slug, label])),
  );
  const tones = { error: "danger", cancelled: "info" } as const;
  const turns: ChatTurnView[] = $derived(
    $helpChat.turns.map((turn) => ({
      key: turn.key,
      question: turn.question,
      html: turn.html,
      running: turn.status === "running",
      activity: turn.activity,
      message: turn.message,
      tone: tones[turn.status as keyof typeof tones] ?? "warning",
      sources: sourceLinks(turn.sources, docTitles),
      rating: turn.rating,
      noteSent: turn.noteSent,
      canRate: turn.id !== "" && (turn.status === "answered" || turn.status === "gap"),
    })),
  );

  const ask = (question: string) => helpChat.ask(question, page.url.pathname);

  // Once there is a conversation, two tabs change between it and the docs.
  const tabs = $derived(
    $helpChat.turns.length === 0
      ? []
      : [
          {
            label: "Docs",
            href: "#docs",
            active: !$helpChat.open,
            onselect: () => helpChat.show(false),
          },
          {
            label: "Ask",
            href: "#ask",
            active: $helpChat.open,
            onselect: () => helpChat.show(true),
          },
        ],
  );

  // A docs link from an answer shows in the Docs tab. The page must be back
  // before the app looks for the link's section. The conversation stays.
  function openDoc(href: string) {
    helpChat.show(false);
    flushSync();
    // eslint-disable-next-line svelte/no-navigation-without-resolve -- A docs path that the server sent.
    goto(href);
  }

  function copy(key: number) {
    navigator.clipboard?.writeText(answerText(helpChat.copied(key))).catch(() => {});
  }

  // The model is ready sooner when it has seen the shared prompt.
  $effect(() => {
    if ($helpChatEnabled) helpChat.warm();
  });
  onDestroy(() => helpChat.close());

  // This column scrolls by itself, so a new page needs to start at the top.
  let content: HTMLDivElement | undefined = $state();
  afterNavigate(({ to }) => {
    // A page that the user opens shows in the Docs tab.
    helpChat.show(false);
    if (content && !to?.url.hash) content.scrollTop = 0;
  });
</script>

<!--
  A wide screen has the list of pages in a column at the left, and the tabs above
  the page or the conversation. A phone has no room for a column, so the list of
  pages is under the tabs, and only the Docs tab shows it: the pages are not part
  of a conversation. The tabs stay in one place when the user changes tab.
-->
<div class="flex h-full flex-col md:flex-row">
  <div class={["max-md:order-2 max-md:shrink-0 md:contents", $helpChat.open && "max-md:hidden"]}>
    <DSSideNav title="Help" {items} />
  </div>
  <div class="flex min-h-0 min-w-0 flex-1 flex-col max-md:contents">
    {#if tabs.length > 0}
      <div class="max-md:order-1 max-md:shrink-0">
        <DSTabBar {tabs} />
      </div>
    {/if}
    <div class="flex min-h-0 min-w-0 flex-1 flex-col max-md:order-3">
      {#if $helpChat.open}
        <!-- The conversation has its own box, at the bottom, for the next question. -->
        <DSChat
          {turns}
          notice="Don't share private info like passwords. Turbo records chats. AI can make mistakes."
          closeLabel="End chat"
          onask={ask}
          onrate={helpChat.rate}
          onnote={helpChat.sendNote}
          onopen={helpChat.opened}
          onnavigate={openDoc}
          oncopy={copy}
          onstop={helpChat.stop}
          onclose={helpChat.close}
        />
      {:else}
        <DSInstantSearch
          label={$helpChatEnabled ? "Search Help, or ask a question" : "Search Help"}
          {search}
          actionLabel={$helpChatEnabled && !$helpChatBusy ? "Ask Help" : undefined}
          onaction={ask}
          bind:scroller={content}
        >
          {@render children()}
        </DSInstantSearch>
      {/if}
    </div>
  </div>
</div>
