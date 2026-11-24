<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { pb } from "$lib/pocketbase";
  import AbsorbProposalQueue from "$lib/components/AbsorbProposalQueue.svelte";
  import type { PageData } from "./$types";

  let { data }: { data: PageData } = $props();
</script>

<svelte:head>
  <title>Absorb Proposals</title>
</svelte:head>

<AbsorbProposalQueue
  status={data.status}
  items={data.items}
  viewerId={pb.authStore.record?.id ?? ""}
  onchange={invalidateAll}
/>
