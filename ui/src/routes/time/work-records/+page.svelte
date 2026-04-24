<script lang="ts">
  import { resolve } from "$app/paths";
  import DSToggle from "$lib/components/DSToggle.svelte";
  import DsActionButton from "$lib/components/DSActionButton.svelte";
  import DsLabel from "$lib/components/DsLabel.svelte";
  import DsList from "$lib/components/DSList.svelte";
  import DsFileLink from "$lib/components/DsFileLink.svelte";
  import type { WorkRecordListRow, WorkRecordSubjectResponse } from "$lib/svelte-types";
  import { authStore } from "$lib/stores/auth";
  import { globalStore } from "$lib/stores/global";
  import { allWorkRecords, myWorkRecords } from "$lib/stores/workRecords";
  import { pocketBaseFileHref, shortDate } from "$lib/utilities";

  const scopeOptions = [
    { id: "mine", label: "Mine / Created" },
    { id: "all", label: "All" },
  ];

  let scope = $state<"mine" | "all">("mine");

  const canViewAll = $derived($globalStore.user_claims_summary.claims.includes("report"));
  const activeStore = $derived(scope === "all" ? $allWorkRecords : $myWorkRecords);
  const items = $derived(activeStore.items);
  const loading = $derived(activeStore.loading);
  const error = $derived(activeStore.error ?? "");
  const listHeader = $derived(scope === "all" ? "All Work Records" : "Mine / Created");
  const currentUserId = $derived($authStore?.model?.id ?? "");

  function linkedEntryHref(subject: WorkRecordSubjectResponse) {
    if (!subject.linked_time_entry_id) {
      return "";
    }
    if (subject.linked_timesheet_id) {
      return `/time/sheets/${subject.linked_timesheet_id}/details?entry=${subject.linked_time_entry_id}`;
    }
    return `/time/entries/${subject.linked_time_entry_id}/details`;
  }

  function mySubject(item: WorkRecordListRow) {
    return item.subjects?.find((subject) => subject.uid === currentUserId);
  }

  function subjectName(subject: WorkRecordSubjectResponse) {
    return subject.owner_name || subject.uid;
  }

  function subjectHours(subject: WorkRecordSubjectResponse) {
    const hours = subject.hours_on_site + subject.hours_travel_time;
    return Number.isInteger(hours) ? hours.toString() : hours.toFixed(2).replace(/0+$/, "").replace(/\.$/, "");
  }

  $effect(() => {
    const reportViewer = canViewAll;
    const currentScope = scope;
    void myWorkRecords.refreshIfStale();
    if (reportViewer) {
      void allWorkRecords.refreshIfStale();
    }
    if (!reportViewer && currentScope === "all") {
      scope = "mine";
    }
  });
</script>

{#if error !== ""}
  <div class="rounded-sm border border-red-300 bg-red-50 p-3 text-sm text-red-700">{error}</div>
{/if}

<DsList {items} search={true} inListHeader={listHeader} searchPlaceholder="search work records...">
  {#snippet searchBarExtra()}
    <div class="flex items-center gap-2 max-[639px]:w-full max-[639px]:flex-wrap">
      {#if canViewAll}
        <DSToggle bind:value={scope} options={scopeOptions} ariaLabel="Work record scope filter" />
      {/if}
      {#if loading}
        <span class="text-sm text-neutral-600">Loading...</span>
      {/if}
    </div>
  {/snippet}

  {#snippet anchor(item: WorkRecordListRow)}
    <a
      href={resolve(`/time/work-records/${encodeURIComponent(item.id)}`)}
      class="text-blue-600 hover:underline"
    >
      {shortDate(item.date)}
    </a>
  {/snippet}

  {#snippet headline(item: WorkRecordListRow)}
    <a
      href={resolve(`/time/work-records/${encodeURIComponent(item.id)}`)}
      class="text-blue-600 hover:underline"
    >
      {item.work_record}
    </a>
  {/snippet}

  {#snippet byline(item: WorkRecordListRow)}
    {#if item.creator === currentUserId}
      <DsLabel color="purple">Creator</DsLabel>
    {/if}
    {#if item.approved}
      <DsLabel color="green">Approved</DsLabel>
    {:else}
      <DsLabel color="orange">Open</DsLabel>
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
    <span class="flex flex-wrap items-center gap-2">
      <span class="opacity-70">{item.work_description}</span>
      {#if item.attachment}
        <a
          href={pocketBaseFileHref("work_records", item.id, item.attachment)}
          target="_blank"
          class="text-blue-600 hover:underline"
        >
          <DsFileLink filename={item.attachment} />
        </a>
      {/if}
    </span>
  {/snippet}

  {#snippet line3(item: WorkRecordListRow)}
    <span class="flex flex-wrap items-center gap-2">
      {#each item.subjects ?? [] as subject (subject.id)}
        <span class="inline-flex flex-wrap items-center gap-1 rounded-sm border border-neutral-300 bg-white/60 px-2 py-1 text-sm">
          <span class:font-semibold={subject.uid === currentUserId}>{subjectName(subject)}</span>
          <span class="text-neutral-600">{subjectHours(subject)}h</span>
          {#if subject.uid === currentUserId}
            <DsLabel color="cyan">You</DsLabel>
          {/if}
          {#if subject.approved}
            <DsLabel color="green">Approved</DsLabel>
          {:else if subject.linked_time_entry_id}
            <DsLabel color="blue">Linked</DsLabel>
          {:else}
            <DsLabel color="orange">Open</DsLabel>
          {/if}
        </span>
      {/each}
      {#if (item.subjects?.length ?? 0) === 0}
        <span class="text-sm text-neutral-600">No workers</span>
      {/if}
    </span>
  {/snippet}

  {#snippet actions(item: WorkRecordListRow)}
    {@const subject = mySubject(item)}
    {#if subject?.can_create_time_entry}
      <DsActionButton
        action={`/time/entries/add/${item.id}?worker=${encodeURIComponent(subject.id)}`}
        icon="mdi:add-bold"
        title="Create Time Entry"
        color="green"
      />
    {:else if subject?.linked_time_entry_id}
      <DsActionButton
        action={linkedEntryHref(subject)}
        icon="mdi:open-in-new"
        title="Open Time Entry"
      />
    {/if}
  {/snippet}
</DsList>
