<script lang="ts">
  import ComputersList from "$lib/components/computers/ComputersList.svelte";
  import type { PageData } from "./$types";

  let { data }: { data: PageData } = $props();
</script>

<svelte:head>
  <title>Computers</title>
</svelte:head>

<ComputersList items={data.items} retired={false} />
