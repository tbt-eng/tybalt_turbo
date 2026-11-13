<script lang="ts">
  import type { AdminProfilesEditPageData } from "$lib/svelte-types";
  import AdminProfilesEditor from "$lib/components/AdminProfilesEditor.svelte";
  import DSHelp from "$lib/components/DSHelp.svelte";
  let { data }: { data: AdminProfilesEditPageData } = $props();
</script>

<div class="mx-auto space-y-4 p-4">
  <div class="flex items-center gap-2">
    <h1 class="min-w-0 flex-1 text-2xl font-bold">
      {#if data.item.surname || data.item.given_name}
        {data.item.given_name} {data.item.surname}
      {:else}
        {data.item.job_title || "-"}
      {/if}
    </h1>
    <DSHelp doc="admin-profiles#change-an-admin-profile" label="Help" />
  </div>

  <AdminProfilesEditor {data} />
</div>
