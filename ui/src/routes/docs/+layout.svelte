<script lang="ts">
  import { afterNavigate } from "$app/navigation";
  import { resolve } from "$app/paths";
  import { page } from "$app/state";
  import DSSideNav from "$lib/components/DSSideNav.svelte";
  let { data, children } = $props();

  const currentSlug = $derived(page.params.slug ?? "start-here");
  const items = $derived(
    data.pages.map(({ slug, title }) => ({
      label: title,
      href: resolve(`/docs/${slug}`),
      active: slug === currentSlug,
    })),
  );

  // This column scrolls by itself, so a new page needs to start at the top.
  let content: HTMLDivElement;
  afterNavigate(({ to }) => {
    if (!to?.url.hash) content.scrollTop = 0;
  });
</script>

<div class="flex h-full flex-col md:flex-row">
  <DSSideNav title="Help" {items} />
  <div bind:this={content} class="min-h-0 min-w-0 flex-1 overflow-y-auto">
    {@render children()}
  </div>
</div>
