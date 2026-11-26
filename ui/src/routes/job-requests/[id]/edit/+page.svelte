<script lang="ts">
  import { goto, invalidateAll } from "$app/navigation";
  import { resolve } from "$app/paths";
  import DSCallout from "$lib/components/DSCallout.svelte";
  import DSHelp from "$lib/components/DSHelp.svelte";
  import DSPageHeader from "$lib/components/DSPageHeader.svelte";
  import RequestEditor from "$lib/components/jobs/RequestEditor.svelte";
  import { globalStore } from "$lib/stores/global";
  let { data } = $props();
  const detail = $derived(data.detail);
  // The request's own page, which shows it read only.
  const back = () => goto(resolve("/job-requests/[id]", { id: detail.request.id }));
  // A save returns to the request's page, where the PDF and Submit are. The
  // editor marks the form saved first, so its unsaved-changes check lets the
  // navigation through.
  async function saved() {
    await globalStore.refreshAttentionCounts();
    await back();
  }
  // Save and keep editing stays here. The reload brings the new version, and
  // the form starts again from it.
  async function keepEditing() {
    await Promise.all([invalidateAll(), globalStore.refreshAttentionCounts()]);
  }
</script>

<!-- The same padding as a new request. Help sits at the right of the title row. -->
<div class="flex w-full flex-col gap-2 p-2">
  <DSPageHeader
    title="Edit Project Request"
    subtitle={detail.request.description || "Untitled request"}
  >
    {#snippet actions()}
      <DSHelp doc="project-requests#change-a-request" label="Help" />
    {/snippet}
  </DSPageHeader>
  {#if detail.request.rejection_reason}
    <DSCallout tone="danger" title="Rejected"
      >Rejection reason: {detail.request.rejection_reason}</DSCallout
    >
  {/if}
  <!-- A link to another request reuses this page, and a save that stays here
       changes the version. The form starts again for both. -->
  {#key `${detail.request.id}:${detail.request.version}`}
    <RequestEditor
      request={detail.request}
      allocations={detail.allocations}
      options={data.options}
      related={detail.related}
      oncancel={back}
      onsaved={saved}
      onkeepediting={keepEditing}
    />
  {/key}
</div>
