<script lang="ts">
  import { resolve } from "$app/paths";
  import type { Snippet } from "svelte";
  import DSList from "$lib/components/DSList.svelte";
  import DsLabel from "$lib/components/DsLabel.svelte";
  import {
    localDay,
    localTime,
    loginReportedName,
    loginTurboUserName,
    type ComputerListState,
    type ComputerLoginRow,
  } from "$lib/computers";
  import { shortDate } from "$lib/utilities";

  // Logins, newest first: the Logins page, which pages on the server, and
  // the newest sign-ins on the page of one computer, which names no computer.
  let {
    items,
    showComputer = true,
    totalItems,
    page = 1,
    perPage = 20,
    searchTerm = $bindable(""),
    onStateChange,
    searchBarExtra,
  }: {
    items: ComputerLoginRow[];
    showComputer?: boolean;
    // Set for the server-paged list, with search and page controls.
    totalItems?: number;
    page?: number;
    perPage?: number;
    searchTerm?: string;
    onStateChange?: (state: ComputerListState) => void;
    searchBarExtra?: Snippet;
  } = $props();

  const paged = $derived(totalItems !== undefined);

  function day(datetime: string): string {
    const local = localDay(datetime);
    return local ? shortDate(local, true) : "";
  }
</script>

<DSList
  {items}
  search={paged}
  pagination={paged}
  {totalItems}
  {page}
  {perPage}
  bind:searchTerm
  onStateChange={(state) =>
    onStateChange?.({ q: state.searchTerm, page: state.page, perPage: state.perPage })}
  searchPlaceholder="search logins..."
  emptyMessage={searchTerm ? "No logins match the search." : "No logins."}
  {searchBarExtra}
>
  {#snippet anchor(login: ComputerLoginRow)}
    {day(login.created)}
  {/snippet}
  {#snippet headline(login: ComputerLoginRow)}
    {#if showComputer}
      <a class="link" href={resolve(`/computers/${login.computer}/details`)}>
        {login.computer_name || login.serial}
      </a>
    {:else}
      {loginReportedName(login)}
    {/if}
  {/snippet}
  {#snippet byline(login: ComputerLoginRow)}
    {localTime(login.created)}
  {/snippet}
  {#snippet line1(login: ComputerLoginRow)}
    {@const turboUser = loginTurboUserName(login)}
    <span class="flex flex-wrap items-center gap-2">
      {#if showComputer}
        <span>{loginReportedName(login)}</span>
      {/if}
      {#if turboUser}
        <span>Turbo user {turboUser}</span>
      {/if}
      {#if !login.user}
        <DsLabel tone="neutral" size="sm">No Turbo user matched</DsLabel>
      {/if}
    </span>
  {/snippet}
</DSList>
