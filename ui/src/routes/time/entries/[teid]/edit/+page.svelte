<script lang="ts">
  import type { PageData } from "./$types";
  import DSHelp from "$lib/components/DSHelp.svelte";
  import TimeEntriesEditor from "$lib/components/TimeEntriesEditor.svelte";
  let { data }: { data: PageData } = $props();
</script>

<div class="flex justify-end px-2 pt-2">
  <DSHelp doc="time#enter-time" label="Help" />
</div>
<TimeEntriesEditor {data} />
