<script lang="ts">
  import DsActionButton from "$lib/components/DSActionButton.svelte";
  import DsList from "$lib/components/DSList.svelte";
  import { pb } from "$lib/pocketbase";
  import type { PaginatedJobWorkRecordsResponse, WorkRecordListRow } from "$lib/svelte-types";

  interface $$Props {
    jobId: string;
    active: boolean;
  }

  let { jobId, active } = $props();
  let records = $state([] as WorkRecordListRow[]);
  let page = $state(1);
  let total = $state(0);
  let totalPages = $state(1);
  let loading = $state(false);
  let error = $state("");
  let loaded = $state(false);
  let loadedJobId = $state("");

  function reset() {
    records = [];
    page = 1;
    total = 0;
    totalPages = 1;
    loading = false;
    error = "";
    loaded = false;
  }

  async function fetchRecords(pageToLoad: number, resetList = false) {
    if (loading) return;

    if (!resetList && loaded && pageToLoad > totalPages) {
      return;
    }

    loading = true;
    error = "";
    try {
      const response = (await pb.send(
        `/api/jobs/${jobId}/work-records?page=${pageToLoad}&limit=50`,
        {
          method: "GET",
        },
      )) as PaginatedJobWorkRecordsResponse;

      records = resetList ? response.data : [...records, ...response.data];
      page = response.page;
      total = response.total;
      totalPages = response.total_pages || 1;
      loaded = true;
    } catch (e: any) {
      error = e?.data?.message ?? e?.message ?? "Failed to load work records.";
      if (resetList) {
        records = [];
        page = 1;
        total = 0;
        totalPages = 1;
      }
      console.error("Failed to load work records for job", e);
    } finally {
      loading = false;
    }
  }

  function loadMore() {
    if (loading || page >= totalPages) {
      return;
    }
    void fetchRecords(page + 1, false);
  }

  $effect(() => {
    if (loadedJobId === jobId) {
      return;
    }

    loadedJobId = jobId;
    reset();
  });

  $effect(() => {
    if (!active || loaded) {
      return;
    }

    void fetchRecords(1, true);
  });
</script>

{#if loaded}
  <div class="px-4">
    <div class="mb-3 text-sm text-neutral-600">
      {total} approved work record{total === 1 ? "" : "s"}
    </div>
  </div>
{/if}

{#if error !== ""}
  <div class="px-4">
    <div class="rounded-sm border border-red-300 bg-red-50 p-3 text-sm text-red-700">
      {error}
    </div>
  </div>
{:else if loading && records.length === 0}
  <div class="px-4">Loading…</div>
{:else if records.length === 0}
  <div class="px-4">No approved work records found.</div>
{:else}
  <div class="w-full overflow-hidden">
    <DsList items={records} search={false} inListHeader="Work Records">
      {#snippet anchor(item: WorkRecordListRow)}
        <a href={`/time/work-records/${item.id}`} class="text-blue-600 hover:underline">{item.work_record}</a>
      {/snippet}
      {#snippet headline(item: WorkRecordListRow)}{item.owner_name}{/snippet}
      {#snippet byline(item: WorkRecordListRow)}{item.date}{/snippet}
      {#snippet line1(item: WorkRecordListRow)}{item.type_name} · {item.location}{/snippet}
      {#snippet line2(item: WorkRecordListRow)}{item.work_description}{/snippet}
    </DsList>

    {#if page < totalPages}
      <div class="mt-4 text-center">
        <DsActionButton action={loadMore} color="blue" loading={loading}>
          Load More
        </DsActionButton>
      </div>
    {/if}
  </div>
{/if}
