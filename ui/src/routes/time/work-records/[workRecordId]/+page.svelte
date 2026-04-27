<script lang="ts">
  import DsActionButton from "$lib/components/DSActionButton.svelte";
  import { pb } from "$lib/pocketbase";
  import { invalidateAll } from "$app/navigation";
  import { globalStore } from "$lib/stores/global";
  import { getApiErrorMessage } from "$lib/errors";
  import type { PageData } from "./$types";
  import { pocketBaseFileHref, shortDate } from "$lib/utilities";
  import DsFileLink from "$lib/components/DsFileLink.svelte";
  import Icon from "@iconify/svelte";

  let { data }: { data: PageData } = $props();
  let noteText = $state("");
  let savingNote = $state(false);
  let noteGuidanceOpen = $state(false);

  const activeConsumables = $derived(
    data.item.consumables.filter((c) =>
      c.input_kind === "boolean" ? c.boolean_selected : c.quantity_number > 0,
    ),
  );

  function linkedEntryHref(worker: (typeof data.item.subjects)[number]) {
    if (!worker.linked_time_entry_id) {
      return "/time/work-records";
    }
    if (worker.linked_timesheet_id) {
      return `/time/sheets/${worker.linked_timesheet_id}/details?entry=${worker.linked_time_entry_id}`;
    }
    return `/time/entries/${worker.linked_time_entry_id}/details`;
  }

  function workerVehicle(worker: (typeof data.item.subjects)[number]) {
    return [
      worker.vehicle_type || "None",
      worker.is_passenger ? "Passenger" : "",
      worker.company_vehicle_unit_number ? `Unit ${worker.company_vehicle_unit_number}` : "",
      `${worker.distance_travelled_km || 0} km`,
    ]
      .filter(Boolean)
      .join(" / ");
  }

  async function addNote() {
    if (savingNote || noteText.trim() === "") return;
    savingNote = true;
    try {
      await pb.send(`/api/work_records/${encodeURIComponent(data.item.id)}/notes`, {
        method: "POST",
        body: JSON.stringify({
          note: noteText.trim(),
        }),
        headers: {
          "Content-Type": "application/json",
        },
      });
      noteText = "";
      await invalidateAll();
    } catch (error: any) {
      globalStore.addError(getApiErrorMessage(error, "Adding work record note failed"));
    } finally {
      savingNote = false;
    }
  }
</script>

<div class="mx-auto max-w-4xl p-4">
  <div class="mb-4 flex flex-wrap items-center justify-between gap-3">
    <div>
      <h1 class="text-2xl font-bold">{data.item.work_record}</h1>
      <div class="mt-1 flex flex-wrap items-center gap-2">
        <span>{shortDate(data.item.date, true)}</span>
        {#if data.item.approved}
          <Icon icon="mdi:lock" width="18" class="text-neutral-500" aria-label="Locked" />
        {/if}
      </div>
    </div>
    <div class="flex flex-wrap gap-2">
      <DsActionButton action="/time/work-records">Back to Work Records</DsActionButton>
      {#if data.item.can_edit}
        <DsActionButton
          action={`/time/work-records/${data.item.id}/edit`}
          icon="mdi:pencil"
          title="Edit"
          color="blue"
        />
      {/if}
      <a
        href={`/time/work-records/${data.item.id}/print`}
        target="_blank"
        class="inline-flex items-center justify-center rounded-sm border border-neutral-300 bg-neutral-100 px-3 py-2 text-sm font-medium text-neutral-800 hover:bg-neutral-200"
      >
        <Icon icon="mdi:printer" width="18px" class="mr-1" />
        Print
      </a>
      {#if data.item.can_copy_to_tomorrow}
        <DsActionButton
          action={`/time/work-records/add?copy_from=${encodeURIComponent(data.item.id)}`}
          icon="mdi:content-copy"
          title="Copy to Tomorrow"
          color="green"
        />
      {/if}
    </div>
  </div>

  <div
    class="summary-grid grid auto-cols-fr grid-flow-col gap-4"
    style="--mobile-rows: repeat(7, auto); --desktop-rows: repeat(4, auto)"
  >
    <div class="rounded-lg bg-neutral-100 p-4">
      <h2 class="mb-2 text-lg font-semibold">Summary</h2>
      <div class="space-y-2 text-sm">
        <div><span class="font-medium">Creator:</span> {data.item.creator_name || data.item.creator || "None"}</div>
        <div><span class="font-medium">Job:</span> {data.item.job_number || "None"}</div>
        <div><span class="font-medium">Type:</span> {data.item.type_name || "None"}</div>
        <div><span class="font-medium">Location:</span> {data.item.location || "None"}</div>
        <div><span class="font-medium">Report To:</span> {data.item.report_to || "None"}</div>
        <div><span class="font-medium">Lane km:</span> {data.item.lane_kms || "None"}</div>
        <div>
          <span class="font-medium">Sub Contractor:</span>
          {data.item.sub_contractor || "None"}
        </div>
      </div>
    </div>

    <div class="rounded-lg bg-neutral-100 p-4">
      <h2 class="mb-2 text-lg font-semibold">Workers</h2>
      <div class="space-y-3">
        {#each data.item.subjects as worker (worker.id)}
          <div class="rounded-sm border border-neutral-200 bg-white p-3 text-sm">
            <div class="flex flex-wrap items-start justify-between gap-3">
              <div class="min-w-0">
                <div class="font-semibold">{worker.worker_name || worker.uid}</div>
                <div class="text-neutral-600">
                  {worker.hours_on_site} site / {worker.hours_travel_time} travel /
                  {worker.hours_on_site + worker.hours_travel_time} total hours
                </div>
                <div class="text-neutral-600">{workerVehicle(worker)}</div>
              </div>
              <div class="flex shrink-0 flex-wrap items-center gap-2">
                {#if worker.approved}
                  <Icon icon="mdi:lock" width="16" class="text-neutral-500" aria-label="Locked" />
                {/if}
                {#if worker.linked_time_entry_id}
                  <DsActionButton
                    action={linkedEntryHref(worker)}
                    icon="mdi:open-in-new"
                    title="Open Time Entry"
                  />
                {:else if worker.can_create_time_entry}
                  <DsActionButton
                    action={`/time/entries/add/${data.item.id}?worker=${encodeURIComponent(worker.id)}`}
                    icon="mdi:add-bold"
                    title="Create Time Entry"
                    color="green"
                  />
                {/if}
              </div>
            </div>
            {#if worker.linked_time_entry_id}
              <div class="mt-2 text-xs text-neutral-600">
                Time entry:
                <a href={linkedEntryHref(worker)} class="text-blue-600 hover:underline">
                  {shortDate(worker.linked_time_entry_date, true)}
                </a>
                {#if worker.linked_timesheet_id}
                  / week ending {shortDate(worker.linked_timesheet_week_ending, true)}
                {/if}
              </div>
            {/if}
          </div>
        {/each}
      </div>
    </div>

    <div class="rounded-lg bg-neutral-100 p-4">
      <h2 class="mb-2 text-lg font-semibold">Field Book</h2>
      <div class="space-y-2 text-sm">
        <div>
          <span class="font-medium">Number:</span>
          {data.item.field_book_number || "None"}
        </div>
        <div>
          <span class="font-medium">Page:</span>
          {data.item.field_book_page_number || "None"}
        </div>
      </div>
    </div>

    <div class="rounded-lg bg-neutral-100 p-4">
      <h2 class="mb-2 text-lg font-semibold">Equipment</h2>
      <p class="text-sm">{data.item.equipment || "None"}</p>
    </div>

    <div class="rounded-lg bg-neutral-100 p-4">
      <h2 class="mb-2 text-lg font-semibold">Supplies</h2>
      <p class="text-sm">{data.item.supplies || "None"}</p>
    </div>

    <div class="rounded-lg bg-neutral-100 p-4">
      <h2 class="mb-2 text-lg font-semibold">Lineage</h2>
      <div class="space-y-2 text-sm">
        <div>
          <span class="font-medium">Root Work Record:</span>
          {data.item.root_work_record || data.item.work_record}
        </div>
        <div>
          <span class="font-medium">Copied From:</span>
          {data.item.parent_work_record || "This is the root record"}
        </div>
      </div>
    </div>
  </div>

  <div class="mt-4 rounded-lg bg-neutral-100 p-4">
    <h2 class="mb-2 text-lg font-semibold">Work Description</h2>
    <p class="text-sm whitespace-pre-wrap text-neutral-800">{data.item.work_description}</p>
  </div>

  <div class="mt-4 rounded-lg bg-neutral-100 p-4">
    <h2 class="mb-2 text-lg font-semibold">Paper Copy PDF</h2>
    {#if data.item.attachment}
      <a
        href={pocketBaseFileHref("work_records", data.item.id, data.item.attachment)}
        target="_blank"
        class="text-blue-600 hover:underline"
      >
        <DsFileLink filename={data.item.attachment} />
      </a>
    {:else}
      <p class="text-sm text-neutral-600">No paper copy PDF attached.</p>
    {/if}
  </div>

  <div class="mt-4 rounded-lg bg-neutral-100 p-4">
    <h2 class="mb-2 text-lg font-semibold">Consumables</h2>
    {#if activeConsumables.length === 0}
      <p class="text-sm text-neutral-600">No consumables recorded.</p>
    {:else}
      <div class="grid grid-cols-2 gap-3 md:grid-cols-3">
        {#each activeConsumables as consumable (consumable.consumable_id)}
          <div
            class="flex items-center justify-between gap-2 rounded-sm border border-neutral-200 bg-white p-3"
          >
            <span class="min-w-0 text-sm font-medium">
              {consumable.consumable_name}
            </span>
            <span class="shrink-0 text-sm text-neutral-700">
              {#if consumable.input_kind === "boolean"}
                <Icon icon="mdi:check" width="16px" class="text-green-600" />
              {:else}
                {consumable.quantity_number}{#if consumable.unit_label}&nbsp;{consumable.unit_label}{/if}
              {/if}
            </span>
          </div>
        {/each}
      </div>
    {/if}
  </div>

  <div class="mt-4 rounded-lg bg-neutral-100 p-4">
    <div class="mb-2 flex items-center justify-between gap-2">
      <h2 class="text-lg font-semibold">Notes</h2>
      <button
        type="button"
        class="inline-flex items-center gap-1 text-sm text-blue-600 hover:underline"
        onclick={() => (noteGuidanceOpen = !noteGuidanceOpen)}
        aria-label="Note guidance"
        aria-expanded={noteGuidanceOpen}
      >
        <Icon icon="mdi:information-outline" width="16px" />
        Note Guidance
      </button>
    </div>

    {#if noteGuidanceOpen}
      <div class="mb-3 rounded-sm border border-blue-300 bg-blue-50 p-3 text-sm text-blue-900">
        Be accurate. Do not swear. Do not write anything you would not want a client to see.
      </div>
    {/if}

    {#if data.item.can_add_note}
      <div class="mb-4 flex flex-col gap-2">
        <textarea
          bind:value={noteText}
          rows="4"
          class="w-full rounded-sm border border-neutral-300 px-2 py-1"
          placeholder="Add a note"
        ></textarea>
        <div>
          <DsActionButton action={addNote} color="green" loading={savingNote}
            >Add Note</DsActionButton
          >
        </div>
      </div>
    {/if}

    {#if data.item.notes.length === 0}
      <p class="text-sm text-neutral-600">No notes yet.</p>
    {:else}
      <div class="space-y-3">
        {#each data.item.notes as note (note.id)}
          <div class="rounded-sm border border-neutral-200 bg-white p-3">
            <div class="text-sm font-medium">{note.author_name || note.uid}</div>
            <div class="text-xs text-neutral-500">{note.created}</div>
            <p class="mt-2 text-sm whitespace-pre-wrap">{note.note}</p>
          </div>
        {/each}
      </div>
    {/if}
  </div>
</div>

<style>
  .summary-grid {
    grid-template-rows: var(--mobile-rows);
  }
  @media (min-width: 768px) {
    .summary-grid {
      grid-template-rows: var(--desktop-rows);
    }
  }
</style>
