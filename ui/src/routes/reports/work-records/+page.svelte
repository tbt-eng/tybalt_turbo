<script lang="ts">
  import DsActionButton from "$lib/components/DSActionButton.svelte";
  import DsAutoComplete from "$lib/components/DSAutoComplete.svelte";
  import DsDateInput from "$lib/components/DSDateInput.svelte";
  import DsList from "$lib/components/DSList.svelte";
  import DsSelector from "$lib/components/DSSelector.svelte";
  import { pb } from "$lib/pocketbase";
  import type {
    WorkRecordConsumablesReportResponse,
    WorkRecordConsumableReportRow,
    WorkRecordListRow,
    WorkRecordSearchResponse,
  } from "$lib/svelte-types";
  import { jobs } from "$lib/stores/jobs";
  import { profiles } from "$lib/stores/profiles";
  import { workRecordTypes } from "$lib/stores/workRecordTypes";

  jobs.init();
  profiles.init();
  workRecordTypes.init();

  type Tab = "records" | "unlinked" | "pending" | "consumables";

  const tabs: Array<{ id: Tab; name: string }> = [
    { id: "records", name: "Records" },
    { id: "unlinked", name: "Unlinked" },
    { id: "pending", name: "Not On Approved Timesheet" },
    { id: "consumables", name: "Consumables" },
  ];

  const pageSize = 100;

  function emptyFilters() {
    return {
      uid: "",
      job: "",
      type: "",
      date_from: "",
      date_to: "",
    };
  }

  let activeTab = $state<Tab>("records");
  let listItems = $state([] as WorkRecordListRow[]);
  let consumableItems = $state([] as WorkRecordConsumableReportRow[]);
  let total = $state(0);
  let loading = $state(false);
  let error = $state("");
  let filters = $state(emptyFilters());
  let appliedFilters = $state(emptyFilters());
  let offset = $state(0);
  let consumablesHasMore = $state(false);
  let requestVersion = 0;

  const hasMore = $derived(activeTab === "consumables" ? consumablesHasMore : listItems.length < total);
  const rowSummary = $derived(
    activeTab === "consumables"
      ? `${consumableItems.length}${hasMore ? "+" : ""} rows loaded`
      : total === 0
        ? "0 rows"
        : `${listItems.length} of ${total} rows`,
  );

  function endpointForTab(tab: Tab) {
    switch (tab) {
      case "records":
        return "/api/work_records";
      case "unlinked":
        return "/api/reports/work_records/unlinked";
      case "pending":
        return "/api/reports/work_records/pending_approval";
      case "consumables":
        return "/api/reports/work_records/consumables";
    }
  }

  async function loadReport(
    tab: Tab,
    currentFilters: ReturnType<typeof emptyFilters>,
    nextOffset = 0,
    append = false,
  ) {
    const nextVersion = ++requestVersion;
    loading = true;
    error = "";
    try {
      const params = new URLSearchParams();
      for (const [key, value] of Object.entries(currentFilters)) {
        if (value !== "") params.set(key, value);
      }
      if (tab === "records") {
        params.set("scope", "all");
      }
      params.set("limit", `${pageSize}`);
      params.set("offset", `${nextOffset}`);

      const response = await pb.send(`${endpointForTab(tab)}?${params.toString()}`, {
        method: "GET",
      });

      if (nextVersion !== requestVersion) return;

      if (tab === "consumables") {
        const typed = response as WorkRecordConsumablesReportResponse;
        consumableItems = append ? [...consumableItems, ...typed.items] : typed.items;
        listItems = [];
        total = consumableItems.length;
        consumablesHasMore = typed.has_more;
      } else {
        const typed = response as WorkRecordSearchResponse;
        listItems = append ? [...listItems, ...typed.items] : typed.items;
        consumableItems = [];
        total = typed.total;
        consumablesHasMore = false;
      }
      offset = nextOffset;
    } catch (err: any) {
      if (nextVersion !== requestVersion) return;
      error = err?.data?.message ?? err?.message ?? "Failed to load work record report";
      if (!append) {
        listItems = [];
        consumableItems = [];
        total = 0;
        offset = 0;
        consumablesHasMore = false;
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
    void loadReport(activeTab, appliedFilters, offset + pageSize, true);
  }

  $effect(() => {
    const currentTab = activeTab;
    const currentFilters = appliedFilters;
    void loadReport(currentTab, currentFilters, 0, false);
  });
</script>

<div class="mx-auto max-w-6xl p-4">
  <div class="mb-4 flex items-center justify-between gap-2">
    <h1 class="text-2xl font-bold">Work Record Reports</h1>
    <DsActionButton action="/time/work-records" color="blue">Open Work Records</DsActionButton>
  </div>

  <div class="mb-4 flex flex-wrap gap-2">
    {#each tabs as tab}
      <DsActionButton
        action={() => {
          activeTab = tab.id;
        }}
        color={activeTab === tab.id ? "green" : undefined}
      >
        {tab.name}
      </DsActionButton>
    {/each}
  </div>

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
      <label for="report_date_from">Date From</label>
      <DsDateInput id="report_date_from" bind:value={filters.date_from} />
    </div>

    <div class="flex flex-col gap-1">
      <label for="report_date_to">Date To</label>
      <DsDateInput id="report_date_to" bind:value={filters.date_to} />
    </div>

    <div class="md:col-span-2 flex gap-2">
      <DsActionButton action={applyFilters} color="green" loading={loading}>Run Report</DsActionButton>
      <DsActionButton action={clearFilters}>Clear</DsActionButton>
    </div>
  </div>

  {#if error !== ""}
    <div class="mt-4 rounded-sm border border-red-300 bg-red-50 p-3 text-sm text-red-700">{error}</div>
  {/if}

  <div class="mt-4 text-sm text-neutral-600">
    {rowSummary}
    {#if loading}
      <span class="ml-2">Loading...</span>
    {/if}
  </div>

  {#if activeTab === "consumables"}
    <div class="mt-4">
      <DsList items={consumableItems} search={true} inListHeader="Consumables">
        {#snippet anchor(item: WorkRecordConsumableReportRow)}
          <a href={`/time/work-records/${item.work_record_id}`} class="text-blue-600 hover:underline">
            {item.work_record}
          </a>
        {/snippet}
        {#snippet headline(item: WorkRecordConsumableReportRow)}{item.consumable_name}{/snippet}
        {#snippet byline(item: WorkRecordConsumableReportRow)}{item.owner_name} · {item.job_number}{/snippet}
        {#snippet line1(item: WorkRecordConsumableReportRow)}
          {item.date} · {item.type_name}
        {/snippet}
        {#snippet line2(item: WorkRecordConsumableReportRow)}
          {#if item.input_kind === "boolean"}
            {item.boolean_selected ? "Selected" : "Not selected"}
          {:else}
            {item.quantity_number} {item.unit_label}
          {/if}
        {/snippet}
      </DsList>

      {#if hasMore}
        <div class="mt-4 text-center">
          <DsActionButton action={loadMore} color="blue" loading={loading}>Load More Consumables</DsActionButton>
        </div>
      {/if}
    </div>
  {:else}
    <div class="mt-4">
      <DsList items={listItems} search={true} inListHeader="Work Records">
        {#snippet anchor(item: WorkRecordListRow)}
          <a href={`/time/work-records/${item.id}`} class="text-blue-600 hover:underline">{item.work_record}</a>
        {/snippet}
        {#snippet headline(item: WorkRecordListRow)}{item.owner_name}{/snippet}
        {#snippet byline(item: WorkRecordListRow)}{item.job_number}{/snippet}
        {#snippet line1(item: WorkRecordListRow)}
          {item.date} · {item.type_name} · {item.location}
        {/snippet}
        {#snippet line2(item: WorkRecordListRow)}{item.work_description}{/snippet}
      </DsList>

      {#if hasMore}
        <div class="mt-4 text-center">
          <DsActionButton action={loadMore} color="blue" loading={loading}>Load More Work Records</DsActionButton>
        </div>
      {/if}
    </div>
  {/if}
</div>
