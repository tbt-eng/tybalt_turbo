<script lang="ts">
  import { goto } from "$app/navigation";
  import { resolve } from "$app/paths";
  import DSHelp from "$lib/components/DSHelp.svelte";
  import RequestEditor from "$lib/components/jobs/RequestEditor.svelte";
  let { data } = $props();
</script>

<!-- The same padding and title row as the other editors, such as JobsEditor. -->
<div class="flex w-full flex-col gap-2 p-2">
  <div class="flex w-full items-center gap-2">
    <h1 class="min-w-0 flex-1 text-xl font-bold text-neutral-800">New Project Request</h1>
    <DSHelp doc="project-requests#create-a-request" label="Help" />
  </div>
  <p>Save a draft, upload a signed PDF that authorizes the project, then submit the request.</p>
  <RequestEditor
    request={data.request}
    options={data.options}
    oncancel={() => goto(resolve("/jobs/list"))}
  />
</div>
