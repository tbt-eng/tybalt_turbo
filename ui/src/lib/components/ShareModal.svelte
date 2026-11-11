<script lang="ts">
  // Shares a time sheet with another manager. Work on a person's time sheet may
  // have been done under another manager. The person's own manager approves the
  // time sheet and shares it so that the other manager can see the hours worked
  // for them. Each viewer is a time_sheet_reviewers row. A viewer finds the time
  // sheet under Shared with Me and can read it, but can't approve or reject it.
  //
  // Only the approver of the time sheet can add or remove a viewer. The create
  // and delete rules of time_sheet_reviewers enforce this, and pages show the
  // button that opens this modal only when canShareTimesheet() in
  // $lib/timesheets/actions.ts is true.
  import DsActionButton from "./DSActionButton.svelte";
  import DSModal from "./DSModal.svelte";
  import DSSection from "./DSSection.svelte";
  import DSSelector from "./DSSelector.svelte";
  import { pb } from "$lib/pocketbase";
  import { managers } from "$lib/stores/managers";
  import { globalStore } from "$lib/stores/global";
  import { getApiErrorMessage } from "$lib/errors";
  import { shortDate } from "$lib/utilities";

  type Reviewer = {
    id: string;
    time_sheet: string;
    reviewer: string;
    reviewed: string;
    surname: string;
    given_name: string;
  };

  type ReviewSummary = {
    owner_name?: string;
    week_ending?: string;
    reviewers?: Reviewer[];
  };

  // initialize the stores, noop if already initialized
  managers.init();

  // onchange gets the time sheet id and its new viewer count after a viewer is
  // added or removed, so the page can update its "Shared (N)" badge.
  let { onchange }: { onchange?: (timeSheetId: string, count: number) => void } = $props();

  const viewerId = pb.authStore.record?.id ?? "";

  let show = $state(false);
  let itemId = $state("");
  let newViewer = $state("");
  let reviewers = $state([] as Reviewer[]);
  let ownerName = $state("");
  let weekEnding = $state("");

  // Managers who are not viewers yet. The signed-in approver is left out.
  const choices = $derived(
    $managers.items.filter((m) => m.id !== viewerId && !reviewers.some((r) => r.reviewer === m.id)),
  );

  function closeModal() {
    show = false;
    itemId = "";
    newViewer = "";
    reviewers = [];
    ownerName = "";
    weekEnding = "";
  }

  async function reloadReviewers() {
    try {
      const summary: ReviewSummary = await pb.send("/api/time_sheets/" + itemId + "/reviewers", {
        method: "GET",
      });
      reviewers = summary.reviewers || [];
      ownerName = summary.owner_name || "";
      weekEnding = summary.week_ending || "";
    } catch (error) {
      globalStore.addError(getApiErrorMessage(error, "Failed to load viewers"));
      reviewers = []; // Ensure we have valid state even on error
    }
  }

  export async function openModal(id: string) {
    show = true;
    itemId = id;
    await reloadReviewers();
  }

  async function addViewer() {
    if (!newViewer) {
      globalStore.addError("Select a manager to add as a viewer");
      return;
    }

    try {
      await pb.collection("time_sheet_reviewers").create({
        time_sheet: itemId,
        reviewer: newViewer,
      });
      newViewer = "";
      await reloadReviewers();
      onchange?.(itemId, reviewers.length);
    } catch (error) {
      globalStore.addError(getApiErrorMessage(error, "Failed to add viewer"));
    }
  }

  async function deleteViewer(reviewerRecordId: string) {
    try {
      await pb.collection("time_sheet_reviewers").delete(reviewerRecordId);
      await reloadReviewers();
      onchange?.(itemId, reviewers.length);
    } catch (error) {
      globalStore.addError(getApiErrorMessage(error, "Failed to remove viewer"));
    }
  }
</script>

<DSModal
  bind:show
  title={ownerName ? `Share ${ownerName}'s time sheet` : "Share time sheet"}
  subtitle={weekEnding ? `Week ending ${shortDate(weekEnding, true)}` : ""}
  size="sm"
  onclose={closeModal}
>
  <div class="flex flex-col gap-4">
    {#if reviewers.length > 0}
      <DSSection title="Viewers">
        <ul class="flex flex-col gap-2">
          {#each reviewers as reviewer (reviewer.id)}
            <li class="flex items-center justify-between gap-2">
              <span>{reviewer.surname}, {reviewer.given_name}</span>
              <DsActionButton
                action={() => deleteViewer(reviewer.id)}
                icon="mdi:delete"
                color="red"
                title="Remove Viewer"
              />
            </li>
          {/each}
        </ul>
      </DSSection>
    {/if}
    <DSSelector
      bind:value={() => newViewer, (value) => (newViewer = String(value))}
      items={choices}
      errors={{}}
      fieldName="reviewer"
      uiName="Manager"
    >
      {#snippet optionTemplate(m)}{m.surname}, {m.given_name}{/snippet}
      {#snippet trailing()}
        <DsActionButton
          action={addViewer}
          icon="feather:plus-circle"
          color="green"
          title="Add Viewer"
        />
      {/snippet}
    </DSSelector>
  </div>
  {#snippet footer()}
    <DsActionButton variant="secondary" action={closeModal}>Close</DsActionButton>
  {/snippet}
</DSModal>
