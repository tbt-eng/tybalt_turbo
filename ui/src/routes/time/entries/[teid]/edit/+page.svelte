<script lang="ts">
  import type { PageData } from "./$types";
  import DSHelp from "$lib/components/DSHelp.svelte";
  import TimeEntriesEditor from "$lib/components/TimeEntriesEditor.svelte";
  let { data }: { data: PageData } = $props();
</script>

<TimeEntriesEditor {data}>
  {#snippet help()}
    <DSHelp doc="time#enter-time" label="Help" />
  {/snippet}
</TimeEntriesEditor>
