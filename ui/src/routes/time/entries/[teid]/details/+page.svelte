<script lang="ts">
  import DsActionButton from "$lib/components/DSActionButton.svelte";
  import DsLabel from "$lib/components/DsLabel.svelte";
  import type { PageData } from "./$types";
  import { shortDate } from "$lib/utilities";

  let { data }: { data: PageData } = $props();

  const linkedWorkRecord = $derived.by(
    () => (data.item.expand?.work_record_subject_id as any)?.expand?.work_record,
  );
  const workRecordNumber = $derived.by(() => data.item.work_record || linkedWorkRecord?.number || "");
</script>

<div class="mx-auto max-w-4xl p-4">
  <div class="mb-4 flex flex-wrap items-center justify-between gap-3">
    <div>
      <h1 class="text-2xl font-bold">Time Entry Details</h1>
      <div class="mt-1 flex flex-wrap items-center gap-2">
        <span>{shortDate(data.item.date, true)}</span>
        {#if workRecordNumber}
          <DsLabel color="cyan">Work Record {workRecordNumber}</DsLabel>
        {/if}
      </div>
    </div>
    <div class="flex flex-wrap gap-2">
      <DsActionButton action="/time/entries/list">Back to Entries</DsActionButton>
      {#if data.item.tsid === ""}
        <DsActionButton action={`/time/entries/${data.item.id}/edit`} color="blue">Edit</DsActionButton>
      {/if}
      {#if data.item.tsid !== ""}
        <DsActionButton action={`/time/sheets/${data.item.tsid}/details?entry=${data.item.id}`}>
          Open Timesheet
        </DsActionButton>
      {/if}
      {#if data.item.work_record_subject_id && linkedWorkRecord?.id}
        <DsActionButton action={`/time/work-records/${linkedWorkRecord.id}`}>Open Work Record</DsActionButton>
      {/if}
    </div>
  </div>

  <div class="grid gap-4 md:grid-cols-2">
    <div class="rounded-lg bg-neutral-100 p-4">
      <h2 class="mb-2 text-lg font-semibold">Summary</h2>
      <div class="space-y-2 text-sm">
        <div><span class="font-medium">Time Type:</span> {data.item.expand?.time_type?.name || "None"}</div>
        <div><span class="font-medium">Division:</span> {data.item.expand?.division?.name || "None"}</div>
        <div><span class="font-medium">Hours:</span> {data.item.hours}</div>
        <div><span class="font-medium">Meals Hours:</span> {data.item.meals_hours}</div>
        <div><span class="font-medium">Week Ending:</span> {shortDate(data.item.week_ending, true)}</div>
      </div>
    </div>

    <div class="rounded-lg bg-neutral-100 p-4">
      <h2 class="mb-2 text-lg font-semibold">Job Context</h2>
      <div class="space-y-2 text-sm">
        <div><span class="font-medium">Job:</span> {data.item.expand?.job?.number || "None"}</div>
        <div><span class="font-medium">Category:</span> {data.item.expand?.category?.name || "None"}</div>
        <div><span class="font-medium">Role:</span> {data.item.expand?.role?.name || "None"}</div>
        {#if data.item.tsid === ""}
          <div>Standalone time entry.</div>
        {:else}
          <div>Bundled into a timesheet.</div>
        {/if}
      </div>
    </div>
  </div>

  <div class="mt-4 rounded-lg bg-neutral-100 p-4">
    <h2 class="mb-2 text-lg font-semibold">Description</h2>
    <p class="whitespace-pre-wrap text-sm text-neutral-800">{data.item.description || "None"}</p>
  </div>
</div>
