<script lang="ts">
  import { afterNavigate } from "$app/navigation";
  import { resolve } from "$app/paths";
  import { page } from "$app/state";
  import DSInstantSearch from "$lib/components/DSInstantSearch.svelte";
  import DSSideNav from "$lib/components/DSSideNav.svelte";
  import { searchDocs } from "$lib/docs";
  import { helpBubble } from "$lib/stores/helpBubble";
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

  // This column scrolls by itself, so a new page needs to start at the top.
  let content: HTMLDivElement | undefined = $state();
  afterNavigate(({ to }) => {
    if (content && !to?.url.hash) content.scrollTop = 0;
  });
</script>

<div class="flex h-full flex-col md:flex-row">
  <DSSideNav title="Help" {items} />
  <DSInstantSearch label="Search Help" {search} bind:scroller={content}>
    {@render children()}
  </DSInstantSearch>
</div>
