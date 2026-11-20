<script lang="ts">
  import ComputersList from "$lib/components/computers/ComputersList.svelte";
  import type { PageData } from "./$types";

  let { data }: { data: PageData } = $props();
</script>

<svelte:head>
  <title>Retired Computers</title>
</svelte:head>

<ComputersList items={data.items} retired={true} />
