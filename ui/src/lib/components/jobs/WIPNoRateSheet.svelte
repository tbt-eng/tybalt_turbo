<script lang="ts">
  import WIPTimeDetails from "./WIPTimeDetails.svelte";
  import DSHelp from "$lib/components/DSHelp.svelte";
  import type { JobWIP } from "./wip";
  let { data, compact = false }: { data: JobWIP; compact?: boolean } = $props();
</script>

{#if data.no_rate_sheet}
  <DSHelp
    label="No rate sheet · staff rates used where available"
    title="Time uses staff defaults"
    warning
    symbol="NR"
    iconOnly={compact}
  >
    <WIPTimeDetails {data} />
  </DSHelp>
{/if}
