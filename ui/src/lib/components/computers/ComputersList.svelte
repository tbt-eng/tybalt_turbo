<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { resolve } from "$app/paths";
  import DSActionButton from "$lib/components/DSActionButton.svelte";
  import DSHelp from "$lib/components/DSHelp.svelte";
  import DSList from "$lib/components/DSList.svelte";
  import DsLabel from "$lib/components/DsLabel.svelte";
  import ComputersTabs from "./ComputersTabs.svelte";
  import {
    computerAssignment,
    computerSearchText,
    lastUserLabel,
    localDay,
    modelLabel,
    type ComputerRow,
  } from "$lib/computers";
  import { getApiErrorMessage } from "$lib/errors";
  import { pb } from "$lib/pocketbase";
  import { globalStore } from "$lib/stores/global";
  import { shortDate } from "$lib/utilities";

  // The Computers list and the Retired list. Both load every computer of the
  // list and page it here.
  let { items, retired }: { items: ComputerRow[]; retired: boolean } = $props();

  let busy = $state("");

  function day(datetime: string): string {
    const local = localDay(datetime);
    return local ? shortDate(local, true) : "";
  }

  async function act(computer: ComputerRow, path: "assign" | "unretire", failure: string) {
    busy = computer.id;
    try {
      await pb.send(`/api/computers/${computer.id}/${path}`, {
        method: "POST",
        // Assign refuses if another person signed in since the page loaded.
        body: path === "assign" ? { user: computer.last_user } : undefined,
      });
      await invalidateAll();
    } catch (error) {
      globalStore.addError(getApiErrorMessage(error, failure));
    } finally {
      busy = "";
    }
  }
</script>

<div>
  <ComputersTabs active={retired ? "retired" : "list"} />
  <DSList
    {items}
    search={true}
    pagination={true}
    searchText={computerSearchText}
    inListHeader={retired ? "Retired Computers" : "Computers"}
    emptyMessage={retired ? "No computers are retired." : "No computers."}
  >
    {#snippet searchBarExtra()}
      {#if retired}
        <DSHelp doc="computers#retired-computers" label="Help" />
      {:else}
        <DSHelp doc="computers#find-a-computer" label="Help" />
      {/if}
    {/snippet}
    {#snippet anchor(computer: ComputerRow)}
      <a class="link" href={resolve(`/computers/${computer.id}/details`)}>
        {computer.computer_name || computer.serial}
      </a>
    {/snippet}
    {#snippet headline(computer: ComputerRow)}
      {modelLabel(computer)}
    {/snippet}
    {#snippet byline(computer: ComputerRow)}
      <span class="flex flex-wrap items-center gap-2">
        {#if !computer.computer_name.includes(computer.serial)}
          <span>{computer.serial}</span>
        {/if}
        {#if retired}
          <DsLabel tone="neutral" size="sm">Retired {day(computer.retired)}</DsLabel>
        {/if}
      </span>
    {/snippet}
    {#snippet line1(computer: ComputerRow)}
      {[
        computer.last_login && `Last sign-in ${day(computer.last_login)}`,
        lastUserLabel(computer),
        computer.os_version && `Windows ${computer.os_version}`,
      ]
        .filter(Boolean)
        .join(" · ")}
    {/snippet}
    {#snippet line2(computer: ComputerRow)}
      {@const assignment = computerAssignment(computer)}
      {#if retired}
        <!-- A retired computer needs no assignment, so nothing is a warning. -->
        {#if computer.assigned_to}
          Assigned to {computer.assigned_to_name}
        {/if}
      {:else if assignment.state === "last_user"}
        Assigned {day(computer.assigned_time)}
      {:else if assignment.state === "other"}
        <DsLabel tone="warning" size="sm">Assigned to {computer.assigned_to_name}</DsLabel>
      {:else}
        <DsLabel tone="warning" size="sm">Not assigned</DsLabel>
      {/if}
    {/snippet}
    {#snippet actions(computer: ComputerRow)}
      {#if retired}
        <DSActionButton
          action={() => act(computer, "unretire", "Un-retire failed")}
          loading={busy === computer.id}
          color="blue"
        >
          Un-retire
        </DSActionButton>
      {:else if computerAssignment(computer).canAssign}
        <DSActionButton
          action={() => act(computer, "assign", "Assign failed")}
          loading={busy === computer.id}
          color="green"
          title={`Assign to ${lastUserLabel(computer)}`}
        >
          Assign
        </DSActionButton>
      {/if}
    {/snippet}
  </DSList>
</div>
