<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { untrack } from "svelte";
  import DSActionButton from "$lib/components/DSActionButton.svelte";
  import DSCallout from "$lib/components/DSCallout.svelte";
  import DSCodeBlock from "$lib/components/DSCodeBlock.svelte";
  import DSField from "$lib/components/DSField.svelte";
  import DSFieldList from "$lib/components/DSFieldList.svelte";
  import DSHelp from "$lib/components/DSHelp.svelte";
  import DSList from "$lib/components/DSList.svelte";
  import DSModal from "$lib/components/DSModal.svelte";
  import DSPopover from "$lib/components/DSPopover.svelte";
  import ComputersTabs from "$lib/components/computers/ComputersTabs.svelte";
  import {
    localDay,
    localTime,
    rawLoginReportText,
    rawLoginSearchText,
    type ComputerRawLoginRow,
  } from "$lib/computers";
  import { computerListNavigator } from "$lib/computersLoad";
  import { getApiErrorMessage } from "$lib/errors";
  import { pb } from "$lib/pocketbase";
  import { formatDateTime, shortDate } from "$lib/utilities";
  import type { PageData } from "./$types";

  let { data }: { data: PageData } = $props();

  // The search box keeps what the person types; the URL follows it.
  let searchTerm = $state(untrack(() => data.q));
  const navigate = computerListNavigator(
    "/computers/raw_logins",
    untrack(() => data.q),
  );

  // The report that the modal shows. It is untrusted JSON from the network,
  // so DSCodeBlock shows it as text.
  let showReport = $state(false);
  let reportRow = $state<ComputerRawLoginRow | null>(null);
  let reportLoading = $state(false);
  let reportError = $state<string | null>(null);

  let showDelete = $state(false);
  let deleteRow = $state<ComputerRawLoginRow | null>(null);
  let deleting = $state(false);
  let deleteError = $state<string | null>(null);

  function day(datetime: string): string {
    const local = localDay(datetime);
    return local ? shortDate(local, true) : "";
  }

  function computerLabel(row: ComputerRawLoginRow): string {
    return row.computer_name || "Unknown computer";
  }

  async function openReport(row: ComputerRawLoginRow) {
    reportRow = row;
    reportError = null;
    reportLoading = true;
    showReport = true;
    try {
      const loaded = (await pb.send(`/api/computers/raw_logins/${row.id}`, {
        method: "GET",
      })) as ComputerRawLoginRow;
      // Keep it only if the modal still shows this row.
      if (reportRow?.id === row.id) reportRow = loaded;
    } catch (error) {
      if (reportRow?.id === row.id) {
        reportError = getApiErrorMessage(error, "The report failed to load");
      }
    } finally {
      if (reportRow?.id === row.id) reportLoading = false;
    }
  }

  function openDelete(row: ComputerRawLoginRow) {
    deleteRow = row;
    deleteError = null;
    showDelete = true;
  }

  async function confirmDelete() {
    if (!deleteRow) return;
    deleting = true;
    try {
      await pb.send(`/api/computers/raw_logins/${deleteRow.id}`, { method: "DELETE" });
      showDelete = false;
      await invalidateAll();
    } catch (error) {
      deleteError = getApiErrorMessage(error, "Delete failed");
    } finally {
      deleting = false;
    }
  }
</script>

<svelte:head>
  <title>Raw Logins</title>
</svelte:head>

<div>
  <ComputersTabs active="raw_logins" />
  <DSList
    items={data.items}
    search={true}
    searchText={rawLoginSearchText}
    bind:searchTerm
    onStateChange={(state) => navigate(state.searchTerm)}
    searchPlaceholder="search raw logins..."
    emptyMessage={searchTerm ? "No raw logins match the search." : "No raw logins."}
  >
    {#snippet searchBarExtra()}
      <DSHelp doc="computers#raw-logins" label="Help" />
    {/snippet}
    {#snippet anchor(row: ComputerRawLoginRow)}
      {day(row.created)}
    {/snippet}
    {#snippet headline(row: ComputerRawLoginRow)}
      {computerLabel(row)}
    {/snippet}
    {#snippet byline(row: ComputerRawLoginRow)}
      {localTime(row.created)}
    {/snippet}
    {#snippet line1(row: ComputerRawLoginRow)}
      {row.error}
    {/snippet}
    {#snippet actions(row: ComputerRawLoginRow)}
      <DSActionButton action={() => openReport(row)} color="blue">Report</DSActionButton>
      <DSActionButton action={() => openDelete(row)} icon="mdi:delete" title="Delete" color="red" />
    {/snippet}
  </DSList>
</div>

<DSModal
  bind:show={showReport}
  title="Raw Login"
  subtitle={reportRow ? computerLabel(reportRow) : ""}
  size="lg"
>
  {#if reportRow}
    <div class="flex flex-col gap-4">
      <DSFieldList>
        <DSField label="Received" value={formatDateTime(reportRow.created.replace(" ", "T"))} />
        <DSField label="Error" value={reportRow.error} />
      </DSFieldList>
      {#if reportError}
        <DSCallout tone="danger">{reportError}</DSCallout>
      {:else if reportLoading}
        <p>Loading the report…</p>
      {:else}
        <DSCodeBlock text={rawLoginReportText(reportRow.report)} label="Copy the report" />
      {/if}
    </div>
  {/if}
</DSModal>

<DSPopover
  bind:show={showDelete}
  title="Delete Raw Login"
  subtitle={deleteRow ? computerLabel(deleteRow) : ""}
  error={deleteError}
  submitting={deleting}
  submitLabel="Delete"
  danger={true}
  onSubmit={confirmDelete}
>
  <p>Turbo deletes the report. You can't undo this.</p>
</DSPopover>
