<script lang="ts">
  import DsActionButton from "$lib/components/DSActionButton.svelte";
  import DsAutoComplete from "$lib/components/DSAutoComplete.svelte";
  import DsDateInput from "$lib/components/DSDateInput.svelte";
  import DsLabel from "$lib/components/DsLabel.svelte";
  import DsList from "$lib/components/DSList.svelte";
  import DsSelector from "$lib/components/DSSelector.svelte";
  import { pb } from "$lib/pocketbase";
  import type { WorkRecordListRow, WorkRecordSearchResponse } from "$lib/svelte-types";
  import { globalStore } from "$lib/stores/global";
  import { jobs } from "$lib/stores/jobs";
  import { profiles } from "$lib/stores/profiles";
  import { workRecordTypes } from "$lib/stores/workRecordTypes";
  import { shortDate } from "$lib/utilities";

  jobs.init();
  profiles.init();
  workRecordTypes.init();

  type WorkRecordFilters = {
    uid: string;
    job: string;
    type: string;
    date_from: string;
    date_to: string;
    linked: string;
    approved: string;
  };

  const booleanOptions = [
    { id: "", name: "Any" },
    { id: "true", name: "Yes" },
    { id: "false", name: "No" },
  ];

  function emptyFilters(): WorkRecordFilters {
    return {
      uid: "",
      job: "",
      type: "",
      date_from: "",
      date_to: "",
      linked: "",
      approved: "",
    };
  }

  const pageSize = 100;

  let items = $state([] as WorkRecordListRow[]);
  let total = $state(0);
  let loading = $state(false);
  let error = $state("");
  let scope = $state<"mine" | "all">("mine");
  let filters = $state(emptyFilters());
  let appliedFilters = $state(emptyFilters());
  let offset = $state(0);
  let requestVersion = 0;

  const canViewAll = $derived($globalStore.user_claims_summary.claims.includes("report"));
  const canCreateWorkRecords = $derived(
    $globalStore.user_claims_summary.claims.includes("time") ||
      $globalStore.user_claims_summary.claims.includes("work_record"),
  );
  const scopeSummary = $derived(
    scope === "all"
      ? "Showing all work records."
      : "Showing records where you are the worker or creator.",
  );
  const listHeader = $derived(scope === "all" ? "All Work Records" : "Mine / Created");
  const hasMore = $derived(items.length < total);

  async function loadRecords(
    currentScope: "mine" | "all",
    currentFilters: WorkRecordFilters,
    nextOffset = 0,
    append = false,
  ) {
    const nextVersion = ++requestVersion;
    loading = true;
    error = "";

    try {
      const params = new URLSearchParams();
      params.set("scope", currentScope);
      params.set("limit", `${pageSize}`);
      params.set("offset", `${nextOffset}`);

      for (const [key, value] of Object.entries(currentFilters)) {
        if (value !== "") {
          params.set(key, value);
        }
      }

      const response = (await pb.send(`/api/work_records?${params.toString()}`, {
        method: "GET",
      })) as WorkRecordSearchResponse;

      if (nextVersion !== requestVersion) return;

      items = append ? [...items, ...response.items] : response.items;
      total = response.total;
      offset = response.offset;
    } catch (err: any) {
      if (nextVersion !== requestVersion) return;
      error = err?.data?.message ?? err?.message ?? "Failed to load work records";
      if (!append) {
        items = [];
        total = 0;
        offset = 0;
      }
    } finally {
      if (nextVersion === requestVersion) {
        loading = false;
      }
    }
  }

  function applyFilters() {
    appliedFilters = { ...filters };
  }

  function clearFilters() {
    const cleared = emptyFilters();
    filters = cleared;
    appliedFilters = { ...cleared };
  }

  function loadMore() {
    if (loading || !hasMore) return;
    void loadRecords(scope, appliedFilters, offset + pageSize, true);
  }

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
    const currentScope = scope;
    const currentFilters = appliedFilters;
    const reportViewer = canViewAll;

    if (!reportViewer && currentScope === "all") {
      scope = "mine";
      return;
    }

    void loadRecords(currentScope, currentFilters, 0, false);
  });
</script>

<div class="mx-auto max-w-6xl p-4">
  <div class="mb-4 flex flex-wrap items-start justify-between gap-3">
    <div>
      <h1 class="text-2xl font-bold">Work Records</h1>
      <p class="text-sm text-neutral-600">{scopeSummary}</p>
    </div>

    <div class="flex flex-wrap items-center justify-end gap-2">
      {#if canViewAll}
        <div class="inline-flex overflow-hidden rounded-sm border border-neutral-300 bg-white">
          <button
            type="button"
            class={`px-3 py-2 text-sm font-medium ${scope === "mine" ? "bg-neutral-900 text-white" : "text-neutral-700"}`}
            onclick={() => (scope = "mine")}
          >
            Mine / Created
          </button>
          <button
            type="button"
            class={`px-3 py-2 text-sm font-medium ${scope === "all" ? "bg-neutral-900 text-white" : "text-neutral-700"}`}
            onclick={() => (scope = "all")}
          >
            All
          </button>
        </div>
      {/if}

      {#if canCreateWorkRecords}
        <DsActionButton action="/time/work-records/add" icon="feather:plus-circle" title="New Work Record" color="green" />
      {/if}
    </div>
  </div>

  {#if canViewAll}
    <div class="grid gap-3 rounded-sm bg-neutral-100 p-4 md:grid-cols-2">
      {#if $profiles.index !== null}
        <DsAutoComplete bind:value={filters.uid} index={$profiles.index} errors={{}} fieldName="uid" uiName="Employee">
          {#snippet resultTemplate(item)}{item.surname}, {item.given_name}{/snippet}
        </DsAutoComplete>
      {/if}

      {#if $jobs.index !== null}
        <DsAutoComplete bind:value={filters.job} index={$jobs.index} errors={{}} fieldName="job" uiName="Job">
          {#snippet resultTemplate(item)}{item.number} - {item.description}{/snippet}
        </DsAutoComplete>
      {/if}

      {#if $workRecordTypes.index !== null}
        <DsAutoComplete bind:value={filters.type} index={$workRecordTypes.index} errors={{}} fieldName="type" uiName="Type">
          {#snippet resultTemplate(item)}{item.name}{/snippet}
        </DsAutoComplete>
      {/if}

      <div class="flex flex-col gap-1">
        <label for="date_from">Date From</label>
        <DsDateInput id="date_from" bind:value={filters.date_from} />
      </div>

      <div class="flex flex-col gap-1">
        <label for="date_to">Date To</label>
        <DsDateInput id="date_to" bind:value={filters.date_to} />
      </div>

      <DsSelector bind:value={filters.linked} items={booleanOptions} errors={{}} fieldName="linked" uiName="Linked">
        {#snippet optionTemplate(item)}{item.name}{/snippet}
      </DsSelector>

      <DsSelector
        bind:value={filters.approved}
        items={booleanOptions}
        errors={{}}
        fieldName="approved"
        uiName="Approved"
      >
        {#snippet optionTemplate(item)}{item.name}{/snippet}
      </DsSelector>

      <div class="md:col-span-2 flex gap-2">
        <DsActionButton action={applyFilters} color="green" loading={loading}>Search</DsActionButton>
        <DsActionButton action={clearFilters}>Clear</DsActionButton>
      </div>
    </div>
  {/if}

  {#if error !== ""}
    <div class="mt-4 rounded-sm border border-red-300 bg-red-50 p-3 text-sm text-red-700">{error}</div>
  {/if}

  <div class="mt-4 text-sm text-neutral-600">
    {#if total === 0}
      0 records
    {:else}
      {items.length} of {total} records
    {/if}
    {#if loading}
      <span class="ml-2">Loading...</span>
    {/if}
  </div>

  <div class="mt-4">
    <DsList items={items} search={true} inListHeader={listHeader}>
      {#snippet anchor(item: WorkRecordListRow)}
        <a
          href={`/time/work-records/${encodeURIComponent(item.work_record)}`}
          class="text-blue-600 hover:underline"
        >
          {shortDate(item.date)}
        </a>
      {/snippet}

      {#snippet headline(item: WorkRecordListRow)}
        <a
          href={`/time/work-records/${encodeURIComponent(item.work_record)}`}
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
          <DsActionButton action={linkedEntryHref(item)} icon="mdi:open-in-new" title="Open Time Entry" />
        {/if}
      {/snippet}
    </DsList>

    {#if hasMore}
      <div class="mt-4 text-center">
        <DsActionButton action={loadMore} color="blue" loading={loading}>Load More Work Records</DsActionButton>
      </div>
    {/if}
  </div>
</div>
