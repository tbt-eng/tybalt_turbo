<script lang="ts">
  import DSDoc from "$lib/components/DSDoc.svelte";
  import DSPageHeader from "$lib/components/DSPageHeader.svelte";
  import DSSection from "$lib/components/DSSection.svelte";
  import DsLabel from "$lib/components/DsLabel.svelte";
  import { commitUrl, type DocPage } from "$lib/docsFormat";
  import { whatsNewSeenBefore } from "$lib/stores/helpBubble";
  import { newWhatsNewWeeks } from "$lib/whatsNewCommit";

  // One full page of the user docs.
  let { page }: { page: DocPage } = $props();

  const newWeeks = $derived(
    $whatsNewSeenBefore === undefined
      ? new Set<string>()
      : newWhatsNewWeeks(
          page.sections.flatMap(({ commit }) => (commit ? [commit] : [])),
          $whatsNewSeenBefore,
        ),
  );
</script>

<div class="mx-auto space-y-4 p-4">
  <div class="space-y-3">
    <DSPageHeader title={page.title} keepInPrint />
    <DSDoc blocks={page.intro} />
  </div>
  {#each page.sections as section (section.id)}
    <div id={section.id} class="scroll-mt-4">
      <DSSection title={section.title}>
        {#snippet titleNote()}
          {#if section.commit}
            {#if newWeeks.has(section.commit)}<DsLabel tone="info" size="sm">New</DsLabel>{/if}
            <!-- eslint-disable-next-line svelte/no-navigation-without-resolve -- A link to GitHub. -->
            <a class="link" href={commitUrl(section.commit)} target="_blank" rel="noopener"
              >to {section.commit}</a
            >
          {/if}
        {/snippet}
        <DSDoc blocks={section.blocks} />
      </DSSection>
    </div>
  {/each}
</div>
