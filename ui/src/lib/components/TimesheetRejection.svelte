<script lang="ts">
  import DsLabel from "./DsLabel.svelte";
  import { shortDate } from "$lib/utilities";

  // A rejected time sheet in a list row: who rejected it, when, and why. It
  // matches the rejection on the time sheet details page. Shows nothing when the
  // time sheet is not rejected.
  let {
    rejected = "",
    rejectorName = "",
    reason = "",
  }: { rejected?: string; rejectorName?: string; reason?: string } = $props();
</script>

{#if rejected !== ""}
  <div>
    <div class="flex flex-wrap items-center gap-1">
      <DsLabel tone="danger">Rejected</DsLabel>
      {#if rejectorName}
        <span>by {rejectorName}</span>
      {/if}
      <span>on {shortDate(rejected, true)}</span>
    </div>
    {#if reason}
      <div class="text-red-600 italic">{reason}</div>
    {/if}
  </div>
{/if}
