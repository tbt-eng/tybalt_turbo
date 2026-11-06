<script lang="ts">
  import DSHelp from "$lib/components/DSHelp.svelte";
  import RequestEditor from "$lib/components/jobs/RequestEditor.svelte";
  let { data } = $props();
</script>

<div class="mx-auto max-w-5xl space-y-4 p-4">
  <div class="flex flex-wrap items-center justify-between gap-2">
    <h1 class="text-2xl font-bold">New Project Request</h1>
    <DSHelp doc="project-requests#create-a-request" label="Help" />
  </div>
  <p>Save a draft, upload a signed PDF that authorizes the project, then submit the request.</p>
  <RequestEditor request={data.request} options={data.options} />
</div>
