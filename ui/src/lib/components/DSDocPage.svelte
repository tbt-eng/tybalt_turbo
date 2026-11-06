<script lang="ts">
  import DSCallout from "$lib/components/DSCallout.svelte";
  import DSDoc from "$lib/components/DSDoc.svelte";
  import DSPageHeader from "$lib/components/DSPageHeader.svelte";
  import DSSection from "$lib/components/DSSection.svelte";
  import { commitUrl, type DocPage } from "$lib/docsFormat";
  import { appConfig, jobRequestsEnabled } from "$lib/stores/appConfig";

  // One full page of the user docs.
  let { page }: { page: DocPage } = $props();
  // Staff can read this page before the flag goes on.
  const requestsOff = $derived(
    page.slug === "project-requests" && $appConfig.initialized && !$jobRequestsEnabled,
  );
</script>

<div class="mx-auto space-y-4 p-4">
  <div class="space-y-3">
    <DSPageHeader title={page.title} keepInPrint />
    {#if requestsOff}
      <DSCallout tone="info">
        Project requests aren't on yet. Create projects the way you do now. The steps to change a
        project's value apply now.
      </DSCallout>
    {/if}
    <DSDoc blocks={page.intro} />
  </div>
  {#each page.sections as section (section.id)}
    <div id={section.id} class="scroll-mt-4">
      <DSSection title={section.title}>
        {#snippet titleNote()}
          {#if section.commit}
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
