<script lang="ts">
  import DsActionButton from "$lib/components/DSActionButton.svelte";
  import DsLabel from "$lib/components/DsLabel.svelte";
  import DsList from "$lib/components/DSList.svelte";
  import type { PageData } from "./$types";
  import type { WorkRecordListRow } from "$lib/svelte-types";
  import { shortDate } from "$lib/utilities";

  let { data }: { data: PageData } = $props();

  function linkedEntryHref(item: WorkRecordListRow) {
    if (!item.linked_time_entry_id) {
      return "";
    }
    if (item.linked_timesheet_id) {
      return `/time/sheets/${item.linked_timesheet_id}/details?entry=${item.linked_time_entry_id}`;
    }
    return `/time/entries/${item.linked_time_entry_id}/details`;
  }
</script>

<div class="mx-auto p-4">
  <h1 class="mb-4 text-2xl font-bold">Work Records</h1>

  <DsList items={data.items} search={true} inListHeader="Work Records">
    {#snippet anchor(item: WorkRecordListRow)}
      <a
        href={`/time/work-records/${encodeURIComponent(item.work_record)}`}
        class="text-blue-600 hover:underline"
      >
        {shortDate(item.date)}
      </a>
    {/snippet}

    {#snippet headline(item: WorkRecordListRow)}
      <a
        href={`/time/work-records/${encodeURIComponent(item.work_record)}`}
        class="text-blue-600 hover:underline"
      >
        {item.work_record}
      </a>
    {/snippet}

    {#snippet byline(item: WorkRecordListRow)}
      {#if item.approved}
        <DsLabel color="green">Approved</DsLabel>
      {:else}
        <DsLabel color="orange">Open</DsLabel>
      {/if}
      {#if item.linked_time_entry_id}
        <DsLabel color="cyan">Linked</DsLabel>
      {/if}
    {/snippet}

    {#snippet line1(item: WorkRecordListRow)}
      <span class="flex flex-wrap items-center gap-1">
        {#if item.job_number}
          <span>{item.job_number}</span>
        {/if}
        {#if item.type_name}
          <DsLabel color="teal">{item.type_name}</DsLabel>
        {/if}
        {#if item.location}
          <span class="opacity-60">{item.location}</span>
        {/if}
      </span>
    {/snippet}

    {#snippet line2(item: WorkRecordListRow)}
      <span class="opacity-70">{item.work_description}</span>
    {/snippet}

    {#snippet line3(item: WorkRecordListRow)}
      {#if item.linked_time_entry_id}
        <span class="opacity-60">
          Linked time entry {item.linked_time_entry_date}
          {#if item.linked_timesheet_id}
            on a bundled timesheet
          {:else}
            as a standalone entry
          {/if}
        </span>
      {/if}
    {/snippet}

    {#snippet actions(item: WorkRecordListRow)}
      {#if item.can_create_time_entry}
        <DsActionButton
          action={`/time/entries/add/${item.id}`}
          icon="mdi:add-bold"
          title="Create Time Entry"
          color="green"
        />
      {:else if item.linked_time_entry_id}
        <DsActionButton action={linkedEntryHref(item)} icon="mdi:open-in-new" title="Open Time Entry" />
      {/if}
    {/snippet}
  </DsList>
</div>
