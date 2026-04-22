<script lang="ts">
  import WorkRecordEditor from "$lib/components/WorkRecordEditor.svelte";
  import type { PageData } from "./$types";

  let { data }: { data: PageData } = $props();
</script>

{#key data.id}
  <WorkRecordEditor
    item={data.item}
    editing={data.editing}
    id={data.id}
    detail={data.detail}
    copyContext={null}
    cancelAction="/time/work-records"
  />
{/key}
