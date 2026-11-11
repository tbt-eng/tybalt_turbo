<script lang="ts">
  import DSHelp from "$lib/components/DSHelp.svelte";
  import DsList from "$lib/components/DSList.svelte";
  import DsLabel from "$lib/components/DsLabel.svelte";
  import DsActionButton from "$lib/components/DSActionButton.svelte";
  import ShareModal from "$lib/components/ShareModal.svelte";
  import TimesheetSharedBadge from "$lib/components/TimesheetSharedBadge.svelte";
  import { pb } from "$lib/pocketbase";
  import { canShareTimesheet } from "$lib/timesheets/actions";
  import type { PageData } from "./$types";
  import type { TimeSheetTallyQueryRow } from "$lib/utilities";
  import { shortDate } from "$lib/utilities";
  import { untrack } from "svelte";

  let { data }: { data: PageData } = $props();
  let items: TimeSheetTallyQueryRow[] = $state(untrack(() => data.items));
  let shareModal: ShareModal | undefined = $state();
  const viewerId = pb.authStore.record?.id ?? "";

  function setSharedCount(id: string, count: number) {
    items = items.map((ts) => (ts.id === id ? { ...ts, shared_reviewer_count: count } : ts));
  }
</script>

<DsList items={items as TimeSheetTallyQueryRow[]} search={true} inListHeader="Approved Time Sheets">
  {#snippet searchBarExtra()}
    <DSHelp doc="review-time#see-the-time-sheets-you-approved" label="Help" />
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
    {#if tally.approved !== ""}
      <div class="flex items-center gap-1">
        <DsLabel color="green">Approved</DsLabel>
        <span>on {shortDate(tally.approved, true)}</span>
      </div>
    {/if}
  {/snippet}
  {#snippet line3(tally: TimeSheetTallyQueryRow)}
    {#if tally.committed !== ""}
      <div class="mt-1 flex items-center gap-1">
        <DsLabel color="blue">Committed</DsLabel>
        <span>on {shortDate(tally.committed, true)} by {tally.committer_name}</span>
      </div>
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
  {/snippet}
</DsList>

<ShareModal bind:this={shareModal} onchange={setSharedCount} />
