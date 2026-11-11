<script lang="ts">
  import DSHelp from "$lib/components/DSHelp.svelte";
  import DsList from "$lib/components/DSList.svelte";
  import DsActionButton from "$lib/components/DSActionButton.svelte";
  import ShareModal from "$lib/components/ShareModal.svelte";
  import TimesheetSharedBadge from "$lib/components/TimesheetSharedBadge.svelte";
  import TimesheetRejection from "$lib/components/TimesheetRejection.svelte";
  import { canShareTimesheet } from "$lib/timesheets/actions";
  import { pb } from "$lib/pocketbase";
  import type { PageData } from "./$types";
  import type { TimeSheetTallyQueryRow } from "$lib/utilities";
  import { shortDate } from "$lib/utilities";
  import { globalStore } from "$lib/stores/global";
  import { untrack } from "svelte";
  import { getApiErrorMessage } from "$lib/errors";

  let { data }: { data: PageData } = $props();
  let items: TimeSheetTallyQueryRow[] = $state(untrack(() => data.items));
  let shareModal: ShareModal | undefined = $state();
  const viewerId = pb.authStore.record?.id ?? "";

  function setSharedCount(id: string, count: number) {
    items = items.map((ts) => (ts.id === id ? { ...ts, shared_reviewer_count: count } : ts));
  }

  async function approve(id: string) {
    try {
      await pb.send(`/api/time_sheets/${id}/approve`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
      });
      await globalStore.refreshAttentionCounts();
      // remove from list after approving
      items = items.filter((ts: TimeSheetTallyQueryRow) => ts.id !== id);
    } catch (error: any) {
      globalStore.addError(getApiErrorMessage(error, "Approve failed"));
    }
  }
</script>

<DsList items={items as TimeSheetTallyQueryRow[]} search={true} inListHeader="Pending Time Sheets">
  {#snippet searchBarExtra()}
    <DSHelp doc="time#approve-a-time-sheet" label="Help" />
  {/snippet}
  {#snippet anchor({ id, week_ending }: TimeSheetTallyQueryRow)}
    <a class="font-bold hover:underline" href={`/time/sheets/${id}/details`}
      >{shortDate(week_ending, true)}</a
    >
  {/snippet}
  {#snippet headline(tally: TimeSheetTallyQueryRow)}
    {tally.given_name} {tally.surname}
  {/snippet}
  {#snippet byline(tally: TimeSheetTallyQueryRow)}
    {#if tally.work_total_hours > 0}
      <span>{tally.work_total_hours} hours worked</span>
    {:else}
      <span>no work</span>
    {/if}
  {/snippet}
  {#snippet line1(tally: TimeSheetTallyQueryRow)}
    <span class="flex flex-wrap items-center gap-1">
      <span>{tally.non_work_total_hours} hours off</span>
      <TimesheetSharedBadge count={tally.shared_reviewer_count} />
    </span>
  {/snippet}
  {#snippet line2(tally: TimeSheetTallyQueryRow)}
    <TimesheetRejection
      rejected={tally.rejected}
      rejectorName={tally.rejector_name}
      reason={tally.rejection_reason}
    />
    {#if tally.committed !== ""}
      <span>Committed on {shortDate(tally.committed, true)} by {tally.committer_name}</span>
    {/if}
  {/snippet}
  {#snippet actions(tally: TimeSheetTallyQueryRow)}
    {#if canShareTimesheet(tally, viewerId)}
      <DsActionButton
        action={() => shareModal?.openModal(tally.id)}
        icon="mdi:ios-share"
        title="Share with another manager"
        color="purple"
      />
    {/if}
    <DsActionButton
      action={`/time/sheets/${tally.id}/details`}
      icon="mdi:details"
      title="Details"
      color="purple"
    />
  {/snippet}
</DsList>

<ShareModal bind:this={shareModal} onchange={setSharedCount} />
