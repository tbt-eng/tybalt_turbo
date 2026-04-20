<script lang="ts">
  import DsActionButton from "$lib/components/DSActionButton.svelte";
  import DsLabel from "$lib/components/DsLabel.svelte";
  import { pb } from "$lib/pocketbase";
  import { invalidateAll } from "$app/navigation";
  import { globalStore } from "$lib/stores/global";
  import { getApiErrorMessage } from "$lib/errors";
  import type { PageData } from "./$types";
  import { shortDate } from "$lib/utilities";
  import Icon from "@iconify/svelte";

  let { data }: { data: PageData } = $props();
  let noteText = $state("");
  let savingNote = $state(false);
  let noteGuidanceOpen = $state(false);

  function linkedEntryHref() {
    if (!data.item.linked_time_entry_id) {
      return "/time/work-records";
    }
    if (data.item.linked_timesheet_id) {
      return `/time/sheets/${data.item.linked_timesheet_id}/details?entry=${data.item.linked_time_entry_id}`;
    }
    return `/time/entries/${data.item.linked_time_entry_id}/details`;
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
      <h1 class="text-2xl font-bold">Work Record {data.item.work_record}</h1>
      <div class="mt-1 flex flex-wrap items-center gap-2">
        <span>{shortDate(data.item.date, true)}</span>
        {#if data.item.approved}
          <DsLabel color="green">Approved</DsLabel>
        {:else}
          <DsLabel color="orange">Open</DsLabel>
        {/if}
      </div>
    </div>
    <div class="flex flex-wrap gap-2">
      <DsActionButton action="/time/work-records">Back to Work Records</DsActionButton>
      {#if data.item.can_edit}
        <DsActionButton action={`/time/work-records/${data.item.id}/edit`} icon="mdi:pencil" title="Edit" color="blue" />
      {/if}
      {#if data.item.can_copy_to_tomorrow}
        <DsActionButton
          action={`/time/work-records/add?copy_from=${encodeURIComponent(data.item.id)}`}
          icon="mdi:content-copy"
          title="Copy to Tomorrow"
          color="green"
        />
      {/if}
      {#if data.item.can_create_time_entry}
        <DsActionButton
          action={`/time/entries/add/${data.item.id}`}
          icon="mdi:add-bold"
          title="Create Time Entry"
          color="green"
        />
      {:else if data.item.linked_time_entry_id}
        <DsActionButton action={linkedEntryHref()} icon="mdi:open-in-new" title="Open Time Entry" />
      {/if}
    </div>
  </div>

  <div class="grid gap-4 md:grid-cols-2">
    <div class="rounded-lg bg-neutral-100 p-4">
      <h2 class="mb-2 text-lg font-semibold">Summary</h2>
      <div class="space-y-2 text-sm">
        <div><span class="font-medium">Owner:</span> {data.item.owner_name || "None"}</div>
        <div><span class="font-medium">Job:</span> {data.item.job_number || "None"}</div>
        <div><span class="font-medium">Type:</span> {data.item.type_name || "None"}</div>
        <div><span class="font-medium">Location:</span> {data.item.location || "None"}</div>
        <div><span class="font-medium">Report To:</span> {data.item.report_to || "None"}</div>
        <div><span class="font-medium">On Site:</span> {data.item.hours_on_site} hours</div>
        <div><span class="font-medium">Travel:</span> {data.item.hours_travel_time} hours</div>
        <div><span class="font-medium">Distance:</span> {data.item.distance_travelled_km} km</div>
        <div>
          <span class="font-medium">Total:</span>
          {data.item.hours_on_site + data.item.hours_travel_time} hours
        </div>
        <div><span class="font-medium">Vehicle:</span> {data.item.vehicle_type || "None"}</div>
        <div><span class="font-medium">Passenger:</span> {data.item.is_passenger ? "Yes" : "No"}</div>
        <div>
          <span class="font-medium">Company Unit:</span>
          {data.item.company_vehicle_unit_number || "None"}
        </div>
        <div><span class="font-medium">Field Book:</span> {data.item.field_book_number || "None"}</div>
        <div>
          <span class="font-medium">Field Book Page:</span>
          {data.item.field_book_page_number || "None"}
        </div>
        <div><span class="font-medium">Sub Contractor:</span> {data.item.sub_contractor || "None"}</div>
        <div><span class="font-medium">Equipment:</span> {data.item.equipment || "None"}</div>
        <div><span class="font-medium">Supplies:</span> {data.item.supplies || "None"}</div>
      </div>
    </div>

    <div class="rounded-lg bg-neutral-100 p-4">
      <h2 class="mb-2 text-lg font-semibold">Linked Time Entry</h2>
      {#if data.item.linked_time_entry_id}
        <div class="space-y-2 text-sm">
          <div><span class="font-medium">Date:</span> {shortDate(data.item.linked_time_entry_date, true)}</div>
          <div>
            <span class="font-medium">Description:</span>
            {data.item.linked_time_entry_description || "None"}
          </div>
          {#if data.item.linked_timesheet_id}
            <div>
              <span class="font-medium">Timesheet Week Ending:</span>
              {shortDate(data.item.linked_time_sheet_week_ending, true)}
            </div>
          {:else}
            <div>Currently a standalone time entry.</div>
          {/if}
        </div>
      {:else}
        <div class="text-sm text-neutral-600">
          No time entry is linked yet. Start the linked entry from this work record.
        </div>
      {/if}
    </div>
  </div>

  <div class="mt-4 rounded-lg bg-neutral-100 p-4">
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

  <div class="mt-4 rounded-lg bg-neutral-100 p-4">
    <h2 class="mb-2 text-lg font-semibold">Work Description</h2>
    <p class="whitespace-pre-wrap text-sm text-neutral-800">{data.item.work_description}</p>
  </div>

  <div class="mt-4 rounded-lg bg-neutral-100 p-4">
    <h2 class="mb-2 text-lg font-semibold">Consumables</h2>
    {#if data.item.consumables.length === 0}
      <p class="text-sm text-neutral-600">No consumables for this work record type.</p>
    {:else}
      <div class="grid gap-3 md:grid-cols-2">
        {#each data.item.consumables as consumable (consumable.consumable_id)}
          <div class="rounded-sm border border-neutral-200 bg-white p-3">
            <div class="font-medium">
              {consumable.consumable_name}
              {#if consumable.unit_label}
                <span class="text-sm text-neutral-500">({consumable.unit_label})</span>
              {/if}
            </div>
            <div class="mt-1 text-sm text-neutral-700">
              {#if consumable.input_kind === "boolean"}
                {consumable.boolean_selected ? "Selected" : "Not selected"}
              {:else}
                {consumable.quantity_number > 0 ? consumable.quantity_number : "No quantity recorded"}
              {/if}
            </div>
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
          <DsActionButton action={addNote} color="green" loading={savingNote}>Add Note</DsActionButton>
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
            <p class="mt-2 whitespace-pre-wrap text-sm">{note.note}</p>
          </div>
        {/each}
      </div>
    {/if}
  </div>
</div>
