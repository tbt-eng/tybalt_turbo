<script lang="ts">
  import DSPageHeader from "$lib/components/DSPageHeader.svelte";
  import DSCallout from "$lib/components/DSCallout.svelte";
  import DsActionButton from "$lib/components/DSActionButton.svelte";
  import DsAutoComplete from "$lib/components/DSAutoComplete.svelte";
  import DsDateInput from "$lib/components/DSDateInput.svelte";
  import DsList from "$lib/components/DSList.svelte";
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
    { id: "pending", name: "Not On Approved Time Sheet" },
    { id: "consumables", name: "Consumables" },
  ];

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
  let loading = $state(false);
  let error = $state("");
  let filters = $state(emptyFilters());
  let appliedFilters = $state(emptyFilters());
  let requestVersion = 0;

  const rowSummary = $derived(
    activeTab === "consumables" ? `${consumableItems.length} rows` : `${listItems.length} rows`,
  );

  function endpointForTab(tab: Tab) {
    switch (tab) {
      case "records":
        return "/api/reports/work_records/search";
      case "unlinked":
        return "/api/reports/work_records/unlinked";
      case "pending":
        return "/api/reports/work_records/pending_approval";
      case "consumables":
        return "/api/reports/work_records/consumables";
    }
  }

  async function loadReport(tab: Tab, currentFilters: ReturnType<typeof emptyFilters>) {
    const nextVersion = ++requestVersion;
    loading = true;
    error = "";
    try {
      const params = new URLSearchParams();
      for (const [key, value] of Object.entries(currentFilters)) {
        if (value !== "") params.set(key, value);
      }
      const response = await pb.send(`${endpointForTab(tab)}?${params.toString()}`, {
        method: "GET",
      });

      if (nextVersion !== requestVersion) return;

      if (tab === "consumables") {
        const typed = response as WorkRecordConsumablesReportResponse;
        consumableItems = typed.items;
        listItems = [];
      } else {
        const typed = response as WorkRecordSearchResponse;
        listItems = typed.items;
        consumableItems = [];
      }
    } catch (err: any) {
      if (nextVersion !== requestVersion) return;
      error = err?.data?.message ?? err?.message ?? "Failed to load work record report";
      listItems = [];
      consumableItems = [];
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

  $effect(() => {
    const currentTab = activeTab;
    const currentFilters = appliedFilters;
    void loadReport(currentTab, currentFilters);
  });
</script>

<div class="mx-auto max-w-6xl p-4">
  <div class="mb-4">
    <DSPageHeader title="Work Record Reports">
      {#snippet actions()}
        <DsActionButton action="/time/work-records" color="blue">Open Work Records</DsActionButton>
      {/snippet}
    </DSPageHeader>
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
      <DsAutoComplete
        bind:value={filters.uid}
        index={$profiles.index}
        idField="uid"
        errors={{}}
        fieldName="uid"
        uiName="Staff Member"
      >
        {#snippet resultTemplate(item)}{item.surname}, {item.given_name}{/snippet}
      </DsAutoComplete>
    {/if}

    {#if $jobs.index !== null}
      <DsAutoComplete
        bind:value={filters.job}
        index={$jobs.index}
        errors={{}}
        fieldName="job"
        uiName="Job"
      >
        {#snippet resultTemplate(item)}{item.number} - {item.description}{/snippet}
      </DsAutoComplete>
    {/if}

    {#if $workRecordTypes.index !== null}
      <DsAutoComplete
        bind:value={filters.type}
        index={$workRecordTypes.index}
        errors={{}}
        fieldName="type"
        uiName="Type"
      >
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

    <div class="flex gap-2 md:col-span-2">
      <DsActionButton action={applyFilters} color="green" {loading}>Run Report</DsActionButton>
      <DsActionButton action={clearFilters}>Clear</DsActionButton>
    </div>
  </div>

  {#if error !== ""}
    <div class="mt-4">
      <DSCallout tone="danger">{error}</DSCallout>
    </div>
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
          <a href={`/time/work-records/${item.work_record_id}`} class="link">
            {item.work_record}
          </a>
        {/snippet}
        {#snippet headline(item: WorkRecordConsumableReportRow)}{item.consumable_name}{/snippet}
        {#snippet byline(item: WorkRecordConsumableReportRow)}{item.creator_name || item.creator} · {item.job_number}{/snippet}
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
    </div>
  {:else}
    <div class="mt-4">
      <DsList items={listItems} search={true} inListHeader="Work Records">
        {#snippet anchor(item: WorkRecordListRow)}
          <a href={`/time/work-records/${item.id}`} class="link">{item.work_record}</a>
        {/snippet}
        {#snippet headline(item: WorkRecordListRow)}{item.creator_name || item.creator}{/snippet}
        {#snippet byline(item: WorkRecordListRow)}{item.job_number}{/snippet}
        {#snippet line1(item: WorkRecordListRow)}
          {item.date} · {item.type_name} · {item.location}
        {/snippet}
        {#snippet line2(item: WorkRecordListRow)}{item.work_description}{/snippet}
      </DsList>
    </div>
  {/if}
</div>
