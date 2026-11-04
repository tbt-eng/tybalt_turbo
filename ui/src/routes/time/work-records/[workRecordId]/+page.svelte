<script lang="ts">
  import DsActionButton from "$lib/components/DSActionButton.svelte";
  import DSPageHeader from "$lib/components/DSPageHeader.svelte";
  import DSCardGrid from "$lib/components/DSCardGrid.svelte";
  import DSCard from "$lib/components/DSCard.svelte";
  import DSFieldList from "$lib/components/DSFieldList.svelte";
  import DSField from "$lib/components/DSField.svelte";
  import DSCallout from "$lib/components/DSCallout.svelte";
  import { pb } from "$lib/pocketbase";
  import { invalidateAll } from "$app/navigation";
  import { globalStore } from "$lib/stores/global";
  import { getApiErrorMessage } from "$lib/errors";
  import type { PageData } from "./$types";
  import { formatDateTime, pocketBaseFileHref, shortDate } from "$lib/utilities";
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

<div class="mx-auto max-w-4xl space-y-4 p-4">
  <DSPageHeader title={data.item.work_record}>
    {#snippet badges()}
      {#if data.item.approved}
        <Icon icon="mdi:lock" width="18" class="text-neutral-500" aria-label="Locked" />
      {/if}
    {/snippet}
    {#snippet context()}
      {shortDate(data.item.date, true)}
    {/snippet}
    {#snippet actions()}
      <DsActionButton action="/time/work-records">Back to Work Records</DsActionButton>
      {#if data.item.can_edit}
        <DsActionButton
          action={`/time/work-records/${data.item.id}/edit`}
          icon="mdi:pencil"
          title="Edit"
          color="blue"
        />
      {/if}
      <DsActionButton
        action={() => window.open(`/time/work-records/${data.item.id}/print`, "_blank")}
        icon="mdi:printer"
        title="Print"
        color="neutral"
      />
      {#if data.item.can_copy_to_tomorrow}
        <DsActionButton
          action={`/time/work-records/add?copy_from=${encodeURIComponent(data.item.id)}`}
          icon="mdi:content-copy"
          title="Copy to Tomorrow"
          color="green"
        />
      {/if}
    {/snippet}
  </DSPageHeader>

  <DSCardGrid>
    <DSCard title="Summary">
      <DSFieldList>
        <DSField label="Creator" value={data.item.creator_name || data.item.creator} />
        <DSField label="Job" value={data.item.job_number} />
        <DSField label="Type" value={data.item.type_name} />
        <DSField label="Location" value={data.item.location} />
        <DSField label="Report to" value={data.item.report_to} />
        <DSField label="Lane km" value={data.item.lane_kms || null} />
        <DSField label="Subcontractor" value={data.item.sub_contractor} />
      </DSFieldList>
    </DSCard>

    <DSCard title="Workers" span={2}>
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
                <a href={linkedEntryHref(worker)} class="link">
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
    </DSCard>

    {#if data.item.field_book_number || data.item.field_book_page_number}
      <DSCard title="Field Book">
        <DSFieldList>
          <DSField label="Number" value={data.item.field_book_number} />
          <DSField label="Page" value={data.item.field_book_page_number || null} />
        </DSFieldList>
      </DSCard>
    {/if}

    {#if data.item.equipment || data.item.supplies}
      <DSCard title="Equipment and Supplies">
        <DSFieldList>
          <DSField label="Equipment" value={data.item.equipment} />
          <DSField label="Supplies" value={data.item.supplies} />
        </DSFieldList>
      </DSCard>
    {/if}

    <DSCard title="Lineage">
      <DSFieldList>
        <DSField
          label="Root work record"
          value={data.item.root_work_record || data.item.work_record}
        />
        <DSField
          label="Copied from"
          value={data.item.parent_work_record || "This is the root record"}
        />
      </DSFieldList>
    </DSCard>

    {#if data.item.attachment}
      <DSCard title="Paper Copy PDF">
        <a
          href={pocketBaseFileHref("work_records", data.item.id, data.item.attachment)}
          target="_blank"
          class="link"
        >
          <DsFileLink filename={data.item.attachment} />
        </a>
      </DSCard>
    {/if}

    <DSCard title="Work Description" span="full">
      <p class="text-sm whitespace-pre-wrap text-neutral-800">{data.item.work_description}</p>
    </DSCard>

    {#if activeConsumables.length > 0}
      <DSCard title="Consumables" span="full">
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
      </DSCard>
    {/if}

    <DSCard title="Notes" span="full">
      {#snippet actions()}
        <DsActionButton
          variant="plain"
          icon="mdi:information-outline"
          action={() => (noteGuidanceOpen = !noteGuidanceOpen)}
        >
          Note Guidance
        </DsActionButton>
      {/snippet}

      {#if noteGuidanceOpen}
        <DSCallout tone="info" icon={false}>
          Be accurate. Do not swear. Do not write anything you would not want a client to see.
        </DSCallout>
      {/if}

      {#if data.item.can_add_note}
        <div class="flex flex-col gap-2">
          <textarea
            bind:value={noteText}
            rows="4"
            class="w-full rounded-sm border border-neutral-300 px-2 py-1"
            placeholder="Add a note"></textarea>
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
              <div class="text-xs text-neutral-500">{formatDateTime(note.created)}</div>
              <p class="mt-2 text-sm whitespace-pre-wrap">{note.note}</p>
            </div>
          {/each}
        </div>
      {/if}
    </DSCard>
  </DSCardGrid>
</div>
