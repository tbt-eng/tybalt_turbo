<script lang="ts">
  import DsActionButton from "$lib/components/DSActionButton.svelte";
  import DsList from "$lib/components/DSList.svelte";
  import { pb } from "$lib/pocketbase";
  import type { PaginatedJobWorkRecordsResponse, WorkRecordListRow } from "$lib/svelte-types";
  import { downloadCSV } from "$lib/utilities";

  interface JobTimeEntry {
    id: string;
    date: string;
    hours: number;
    description: string;
    work_record: string;
    week_ending: string;
    tsid: string;
    division_code: string;
    time_type_code: string;
    surname: string;
    given_name: string;
    category_name: string;
  }

  interface $$Props {
    summary: Record<string, any>;
    items: JobTimeEntry[];
    listLoading: boolean;
    loadMore: () => void;
    page: number;
    totalPages: number;
    jobId: string;
    approvedWorkRecordsActive: boolean;
  }

  const { summary, items, listLoading, loadMore, page, totalPages, jobId, approvedWorkRecordsActive } = $props();
  let approvedWorkRecords = $state([] as WorkRecordListRow[]);
  let approvedWorkRecordsPage = $state(1);
  let approvedWorkRecordsTotalPages = $state(1);
  let approvedWorkRecordsLoading = $state(false);
  let approvedWorkRecordsError = $state("");
  let approvedWorkRecordsLoaded = $state(false);

  async function fetchFullReport() {
    const url = `${pb.baseUrl}/api/jobs/${jobId}/time/full_report`;
    const fileName = `job_full_time_report_${jobId}.csv`;
    await downloadCSV(url, fileName);
  }

  async function fetchApprovedWorkRecords(pageToLoad: number, reset = false) {
    if (approvedWorkRecordsLoading) return;

    if (!reset && approvedWorkRecordsLoaded && pageToLoad > approvedWorkRecordsTotalPages) {
      return;
    }

    approvedWorkRecordsLoading = true;
    approvedWorkRecordsError = "";
    try {
      const response = (await pb.send(
        `/api/jobs/${jobId}/work-records?page=${pageToLoad}&limit=50`,
        {
          method: "GET",
        },
      )) as PaginatedJobWorkRecordsResponse;

      approvedWorkRecords = reset ? response.data : [...approvedWorkRecords, ...response.data];
      approvedWorkRecordsPage = response.page;
      approvedWorkRecordsTotalPages = response.total_pages || 1;
      approvedWorkRecordsLoaded = true;
    } catch (error: any) {
      approvedWorkRecordsError =
        error?.data?.message ?? error?.message ?? "Failed to load approved work records.";
      if (reset) {
        approvedWorkRecords = [];
        approvedWorkRecordsPage = 1;
        approvedWorkRecordsTotalPages = 1;
      }
      console.error("Failed to load approved work records for job", error);
    } finally {
      approvedWorkRecordsLoading = false;
    }
  }

  function loadMoreApprovedWorkRecords() {
    if (approvedWorkRecordsLoading || approvedWorkRecordsPage >= approvedWorkRecordsTotalPages) {
      return;
    }
    void fetchApprovedWorkRecords(approvedWorkRecordsPage + 1, false);
  }

  $effect(() => {
    if (!approvedWorkRecordsActive || approvedWorkRecordsLoaded) {
      return;
    }

    approvedWorkRecordsPage = 1;
    approvedWorkRecordsTotalPages = 1;
    approvedWorkRecords = [];
    void fetchApprovedWorkRecords(1, true);
  });
</script>

<div class="px-4">
  <div class="mb-4 rounded-sm border border-neutral-200 bg-white p-3">
    <h3 class="mb-2 font-semibold">Approved Work Records</h3>

    {#if approvedWorkRecordsError !== ""}
      <div class="rounded-sm border border-red-300 bg-red-50 p-3 text-sm text-red-700">
        {approvedWorkRecordsError}
      </div>
    {:else if approvedWorkRecordsLoading && approvedWorkRecords.length === 0}
      <div class="text-sm text-neutral-600">Loading approved work records…</div>
    {:else if approvedWorkRecords.length === 0}
      <div class="text-sm text-neutral-600">No approved work records found.</div>
    {:else}
      <DsList items={approvedWorkRecords} search={false} inListHeader="Approved Work Records">
        {#snippet anchor(item: WorkRecordListRow)}
          <a href={`/time/work-records/${item.id}`} class="text-blue-600 hover:underline">{item.work_record}</a>
        {/snippet}
        {#snippet headline(item: WorkRecordListRow)}{item.owner_name}{/snippet}
        {#snippet byline(item: WorkRecordListRow)}{item.date}{/snippet}
        {#snippet line1(item: WorkRecordListRow)}{item.type_name} · {item.location}{/snippet}
        {#snippet line2(item: WorkRecordListRow)}{item.work_description}{/snippet}
      </DsList>

      {#if approvedWorkRecordsPage < approvedWorkRecordsTotalPages}
        <div class="mt-4 text-center">
          <DsActionButton action={loadMoreApprovedWorkRecords} color="blue" loading={approvedWorkRecordsLoading}>
            Load More Approved Work Records
          </DsActionButton>
        </div>
      {/if}
    {/if}
  </div>

  <div class="space-y-1">
    <div><span class="font-semibold">Total Hours:</span> {summary.total_hours ?? 0}</div>
    {#if summary.earliest_entry}
      <div>
        <span class="font-semibold">Date Range:</span>
        {summary.earliest_entry} – {summary.latest_entry}
      </div>
    {/if}
  </div>
</div>

<div class="mt-2 px-4">
  <button type="button" onclick={fetchFullReport} class="text-blue-600 hover:underline">
    Full Report (CSV)
  </button>
</div>

{#if listLoading && items.length === 0}
  <div class="px-4">Loading…</div>
{:else if items.length === 0}
  <div class="px-4">No entries found.</div>
{:else}
  <div class="w-full overflow-hidden">
    <DsList {items} search={false} inListHeader="Time Entries">
      {#snippet anchor(item: JobTimeEntry)}{item.date}{/snippet}
      {#snippet headline(item: JobTimeEntry)}{item.hours}{/snippet}
      {#snippet byline(item: JobTimeEntry)}{item.given_name} {item.surname}{/snippet}
      {#snippet line1(item: JobTimeEntry)}
        <span class="font-bold">{item.division_code}</span>
        <span class="font-bold">{item.time_type_code}</span>
        {item.description}
      {/snippet}
    </DsList>
    {#if page < totalPages}
      <div class="mt-4 text-center">
        <button
          class="rounded-sm bg-blue-600 px-4 py-2 text-white"
          onclick={loadMore}
          disabled={listLoading}
        >
          {listLoading ? "Loading…" : "Load More"}
        </button>
      </div>
    {/if}
  </div>
{/if}
