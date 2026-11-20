<script lang="ts">
  import { untrack } from "svelte";
  import DSHelp from "$lib/components/DSHelp.svelte";
  import ComputerLoginsList from "$lib/components/computers/ComputerLoginsList.svelte";
  import ComputersTabs from "$lib/components/computers/ComputersTabs.svelte";
  import { computerListNavigator } from "$lib/computersLoad";
  import type { PageData } from "./$types";

  let { data }: { data: PageData } = $props();

  // The search box keeps what the person types; the URL follows it.
  let searchTerm = $state(untrack(() => data.state.q));
  const navigate = computerListNavigator(
    "/computers/logins",
    untrack(() => data.state),
  );
</script>

<svelte:head>
  <title>Logins</title>
</svelte:head>

<div>
  <ComputersTabs active="logins" />
  <ComputerLoginsList
    items={data.result.items}
    totalItems={data.result.totalItems}
    page={data.result.page}
    perPage={data.result.perPage}
    bind:searchTerm
    onStateChange={navigate}
  >
    {#snippet searchBarExtra()}
      <DSHelp doc="computers#logins" label="Help" />
    {/snippet}
  </ComputerLoginsList>
</div>
