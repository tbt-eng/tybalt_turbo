<script lang="ts">
  import { resolve } from "$app/paths";
  import DSToggle from "$lib/components/DSToggle.svelte";
  import DsActionButton from "$lib/components/DSActionButton.svelte";
  import DsLabel from "$lib/components/DsLabel.svelte";
  import DsList from "$lib/components/DSList.svelte";
  import type { WorkRecordListRow } from "$lib/svelte-types";
  import { globalStore } from "$lib/stores/global";
  import { allWorkRecords, myWorkRecords } from "$lib/stores/workRecords";
  import { shortDate } from "$lib/utilities";

  const scopeOptions = [
    { id: "mine", label: "Mine / Created" },
    { id: "all", label: "All" },
  ];

  let scope = $state<"mine" | "all">("mine");

  const canViewAll = $derived($globalStore.user_claims_summary.claims.includes("report"));
  const activeStore = $derived(scope === "all" ? $allWorkRecords : $myWorkRecords);
  const items = $derived(activeStore.items);
  const loading = $derived(activeStore.loading);
  const error = $derived(activeStore.error ?? "");
  const listHeader = $derived(scope === "all" ? "All Work Records" : "Mine / Created");

  function linkedEntryHref(item: WorkRecordListRow) {
    if (!item.linked_time_entry_id) {
      return "";
    }
    if (item.linked_timesheet_id) {
      return `/time/sheets/${item.linked_timesheet_id}/details?entry=${item.linked_time_entry_id}`;
    }
    return `/time/entries/${item.linked_time_entry_id}/details`;
  }

  $effect(() => {
    const reportViewer = canViewAll;
    const currentScope = scope;
    void myWorkRecords.refreshIfStale();
    if (reportViewer) {
      void allWorkRecords.refreshIfStale();
    }
    if (!reportViewer && currentScope === "all") {
      scope = "mine";
    }
  });
</script>

{#if error !== ""}
  <div class="rounded-sm border border-red-300 bg-red-50 p-3 text-sm text-red-700">{error}</div>
{/if}

<DsList {items} search={true} inListHeader={listHeader} searchPlaceholder="search work records...">
  {#snippet searchBarExtra()}
    <div class="flex items-center gap-2 max-[639px]:w-full max-[639px]:flex-wrap">
      {#if canViewAll}
        <DSToggle bind:value={scope} options={scopeOptions} ariaLabel="Work record scope filter" />
      {/if}
      {#if loading}
        <span class="text-sm text-neutral-600">Loading...</span>
      {/if}
    </div>
  {/snippet}

  {#snippet anchor(item: WorkRecordListRow)}
    <a
      href={resolve(`/time/work-records/${encodeURIComponent(item.id)}`)}
      class="text-blue-600 hover:underline"
    >
      {shortDate(item.date)}
    </a>
  {/snippet}

  {#snippet headline(item: WorkRecordListRow)}
    <a
      href={resolve(`/time/work-records/${encodeURIComponent(item.id)}`)}
      class="text-blue-600 hover:underline"
    >
      {item.work_record}
    </a>
  {/snippet}

  {#snippet byline(item: WorkRecordListRow)}
    {#if item.approved}
      <DsLabel color="green">Approved</DsLabel>
    {:else}
      <DsLabel color="orange">Open</DsLabel>
    {/if}
    {#if item.linked_time_entry_id}
      <DsLabel color="cyan">Linked</DsLabel>
    {/if}
  {/snippet}

  {#snippet line1(item: WorkRecordListRow)}
    <span class="flex flex-wrap items-center gap-1">
      {#if item.owner_name}
        <span>{item.owner_name}</span>
      {/if}
      {#if item.job_number}
        <span>{item.job_number}</span>
      {/if}
      {#if item.type_name}
        <DsLabel color="teal">{item.type_name}</DsLabel>
      {/if}
      {#if item.location}
        <span class="opacity-60">{item.location}</span>
      {/if}
    </span>
  {/snippet}

  {#snippet line2(item: WorkRecordListRow)}
    <span class="opacity-70">{item.work_description}</span>
  {/snippet}

  {#snippet line3(item: WorkRecordListRow)}
    {#if item.linked_time_entry_id}
      <span class="opacity-60">
        Linked time entry {item.linked_time_entry_date}
        {#if item.linked_timesheet_id}
          on a bundled timesheet
        {:else}
          as a standalone entry
        {/if}
      </span>
    {/if}
  {/snippet}

  {#snippet actions(item: WorkRecordListRow)}
    {#if item.can_create_time_entry}
      <DsActionButton
        action={`/time/entries/add/${item.id}`}
        icon="mdi:add-bold"
        title="Create Time Entry"
        color="green"
      />
    {:else if item.linked_time_entry_id}
      <DsActionButton
        action={linkedEntryHref(item)}
        icon="mdi:open-in-new"
        title="Open Time Entry"
      />
    {/if}
  {/snippet}
</DsList>
