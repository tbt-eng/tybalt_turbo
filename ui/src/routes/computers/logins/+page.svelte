<script lang="ts">
  import { untrack } from "svelte";
  import DSHelp from "$lib/components/DSHelp.svelte";
  import ComputerLoginsList from "$lib/components/computers/ComputerLoginsList.svelte";
  import ComputersTabs from "$lib/components/computers/ComputersTabs.svelte";
  import { computerListNavigator } from "$lib/computersLoad";
  import type { PageData } from "./$types";

  let { data }: { data: PageData } = $props();

  // The search box keeps what the person types; the URL follows it.
  let searchTerm = $state(untrack(() => data.q));
  const navigate = computerListNavigator(
    "/computers/logins",
    untrack(() => data.q),
  );
</script>

<svelte:head>
  <title>Logins</title>
</svelte:head>

<div>
  <ComputersTabs active="logins" />
  <ComputerLoginsList items={data.items} bind:searchTerm onSearch={navigate}>
    {#snippet searchBarExtra()}
      <DSHelp doc="computers#logins" label="Help" />
    {/snippet}
  </ComputerLoginsList>
</div>
