<script lang="ts">
  import DsActionButton from "$lib/components/DSActionButton.svelte";
  import DsLabel from "$lib/components/DsLabel.svelte";
  import { pb } from "$lib/pocketbase";
  import { goto } from "$app/navigation";
  import { globalStore } from "$lib/stores/global";
  import { getApiErrorMessage } from "$lib/errors";
  import type { PageData } from "./$types";
  import { shortDate } from "$lib/utilities";

  let { data }: { data: PageData } = $props();

  function linkedEntryHref() {
    if (!data.item.linked_time_entry_id) {
      return "/time/work-records";
    }
    if (data.item.linked_timesheet_id) {
      return `/time/sheets/${data.item.linked_timesheet_id}/details?entry=${data.item.linked_time_entry_id}`;
    }
    return `/time/entries/${data.item.linked_time_entry_id}/details`;
  }

  async function copyToTomorrow() {
    try {
      const result = (await pb.send(
        `/api/work_records/${encodeURIComponent(data.item.id)}/copy_to_tomorrow`,
        {
          method: "POST",
        },
      )) as { id: string };
      goto(`/time/work-records/${result.id}`);
    } catch (error: any) {
      globalStore.addError(getApiErrorMessage(error, "Copy to tomorrow failed"));
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
      {#if data.item.can_copy_to_tomorrow}
        <DsActionButton
          action={() => copyToTomorrow()}
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
        <div><span class="font-medium">Job:</span> {data.item.job_number || "None"}</div>
        <div><span class="font-medium">Type:</span> {data.item.type_name || "None"}</div>
        <div><span class="font-medium">Location:</span> {data.item.location || "None"}</div>
        <div><span class="font-medium">On Site:</span> {data.item.hours_on_site} hours</div>
        <div><span class="font-medium">Travel:</span> {data.item.hours_travel_time} hours</div>
        <div>
          <span class="font-medium">Total:</span>
          {data.item.hours_on_site + data.item.hours_travel_time} hours
        </div>
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
    <h2 class="mb-2 text-lg font-semibold">Work Description</h2>
    <p class="whitespace-pre-wrap text-sm text-neutral-800">{data.item.work_description}</p>
  </div>
</div>
