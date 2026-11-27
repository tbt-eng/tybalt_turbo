<script lang="ts">
  import { resolve } from "$app/paths";
  import type { Snippet } from "svelte";
  import DSList from "$lib/components/DSList.svelte";
  import DsLabel from "$lib/components/DsLabel.svelte";
  import {
    localDay,
    localTime,
    loginReportedName,
    loginSearchText,
    loginTurboUserName,
    type ComputerLoginRow,
  } from "$lib/computers";
  import { shortDate } from "$lib/utilities";

  // Logins, newest first: the Logins page, which searches on the server, and
  // every sign-in on the page of one computer, which names no computer and
  // has no search.
  let {
    items,
    showComputer = true,
    searchTerm = $bindable(""),
    onSearch,
    searchBarExtra,
  }: {
    items: ComputerLoginRow[];
    showComputer?: boolean;
    searchTerm?: string;
    // Set for the Logins page: it shows the search, and each change of the
    // search goes here so the page can load the matches.
    onSearch?: (q: string) => void;
    searchBarExtra?: Snippet;
  } = $props();

  function day(datetime: string): string {
    const local = localDay(datetime);
    return local ? shortDate(local, true) : "";
  }
</script>

<DSList
  {items}
  search={onSearch !== undefined}
  searchText={loginSearchText}
  bind:searchTerm
  onStateChange={(state) => onSearch?.(state.searchTerm)}
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
