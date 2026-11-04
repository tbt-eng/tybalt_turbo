<script lang="ts">
  import DsActionButton from "$lib/components/DSActionButton.svelte";
  import DsLabel from "$lib/components/DsLabel.svelte";
  import DSPageHeader from "$lib/components/DSPageHeader.svelte";
  import DSCardGrid from "$lib/components/DSCardGrid.svelte";
  import DSCard from "$lib/components/DSCard.svelte";
  import DSFieldList from "$lib/components/DSFieldList.svelte";
  import DSField from "$lib/components/DSField.svelte";
  import type { PageData } from "./$types";
  import { shortDate } from "$lib/utilities";

  let { data }: { data: PageData } = $props();

  const linkedWorkRecord = $derived.by(
    () => (data.item.expand?.work_record_subject_id as any)?.expand?.work_record,
  );
  const workRecordNumber = $derived.by(
    () => data.item.work_record || linkedWorkRecord?.number || "",
  );
</script>

<div class="mx-auto max-w-4xl space-y-4 p-4">
  <DSPageHeader
    title={shortDate(data.item.date, true)}
    subtitle={data.item.expand?.time_type?.name}
  >
    {#snippet badges()}
      {#if workRecordNumber}
        <DsLabel color="cyan">Work Record {workRecordNumber}</DsLabel>
      {/if}
    {/snippet}
    {#snippet context()}
      {data.item.tsid === "" ? "Standalone time entry" : "Bundled into a timesheet"}
    {/snippet}
    {#snippet actions()}
      <DsActionButton action="/time/entries/list">Back to Entries</DsActionButton>
      {#if data.item.tsid === ""}
        <DsActionButton action={`/time/entries/${data.item.id}/edit`} color="blue"
          >Edit</DsActionButton
        >
      {/if}
      {#if data.item.tsid !== ""}
        <DsActionButton action={`/time/sheets/${data.item.tsid}/details?entry=${data.item.id}`}>
          Open Timesheet
        </DsActionButton>
      {/if}
      {#if data.item.work_record_subject_id && linkedWorkRecord?.id}
        <DsActionButton action={`/time/work-records/${linkedWorkRecord.id}`}
          >Open Work Record</DsActionButton
        >
      {/if}
    {/snippet}
  </DSPageHeader>

  <DSCardGrid>
    <DSCard title="Summary">
      <DSFieldList>
        <DSField label="Division" value={data.item.expand?.division?.name} />
        <DSField label="Hours" value={data.item.hours} />
        <DSField label="Meals hours" value={data.item.meals_hours || null} />
        <DSField label="Week ending" value={shortDate(data.item.week_ending, true)} />
      </DSFieldList>
    </DSCard>

    {#if data.item.expand?.job}
      <DSCard title="Job">
        <DSFieldList>
          <DSField label="Job" value={data.item.expand?.job?.number} />
          <DSField label="Category" value={data.item.expand?.category?.name} />
          <DSField label="Role" value={data.item.expand?.role?.name} />
        </DSFieldList>
      </DSCard>
    {/if}

    {#if data.item.description}
      <DSCard title="Description" span="full">
        <p class="text-sm whitespace-pre-wrap text-neutral-800">{data.item.description}</p>
      </DSCard>
    {/if}
  </DSCardGrid>
</div>
