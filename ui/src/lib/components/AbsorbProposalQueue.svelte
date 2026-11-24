<script lang="ts">
  import { resolve } from "$app/paths";
  import { pb } from "$lib/pocketbase";
  import { globalStore } from "$lib/stores/global";
  import { getApiErrorMessage } from "$lib/errors";
  import { shortDate } from "$lib/utilities";
  import {
    absorbProposalCollectionLabel,
    absorbProposalReviewHref,
    absorbProposalStatusLabel,
    canAbsorbDirectly,
    type AbsorbProposalRow,
  } from "$lib/absorbProposals";
  import DSActionButton from "./DSActionButton.svelte";
  import DSCallout from "./DSCallout.svelte";
  import DsLabel from "./DsLabel.svelte";
  import DsList from "./DSList.svelte";
  import DSPopover from "./DSPopover.svelte";
  import DSStatusLabel from "./DSStatusLabel.svelte";
  import DSTabBar from "./DSTabBar.svelte";
  import DSTextInput from "./DSTextInput.svelte";

  // The absorb proposal queue. Absorb holders see every proposal and review
  // or decline the pending ones; everyone else sees their own and may
  // withdraw a pending one. A requester may delete their withdrawn proposal.
  let {
    status,
    items,
    viewerId,
    onchange,
  }: {
    status: "pending" | "decided";
    items: AbsorbProposalRow[];
    viewerId: string;
    // Reloads the items after a change.
    onchange: () => Promise<void>;
  } = $props();

  const pendingTab = $derived(status === "pending");
  const canReview = $derived(canAbsorbDirectly($globalStore.claims));

  let searchTerm = $state("");

  let declineRow = $state<AbsorbProposalRow | null>(null);
  let showDecline = $state(false);
  let declineComment = $state("");
  let declining = $state(false);
  let declineError = $state<string | null>(null);

  let deleteRow = $state<AbsorbProposalRow | null>(null);
  let showDelete = $state(false);
  let deleting = $state(false);
  let deleteError = $state<string | null>(null);
  // Deleted proposals leave the list before the reload returns.
  let deletedIds = $state<string[]>([]);
  const visibleItems = $derived(items.filter((row) => !deletedIds.includes(row.id)));

  function recordNames(row: AbsorbProposalRow) {
    return row.records.map((record) => record.record_name).join(", ");
  }

  function rowTitle(row: AbsorbProposalRow) {
    return row.client_name ? `${row.target_name} · ${row.client_name}` : row.target_name;
  }

  async function refresh() {
    await globalStore.refreshAttentionCounts();
    await onchange();
  }

  function openDecline(row: AbsorbProposalRow) {
    declineRow = row;
    declineComment = "";
    declineError = null;
    showDecline = true;
  }

  async function decline() {
    if (!declineRow) return;
    declining = true;
    declineError = null;
    try {
      await pb.send(`/api/absorb_proposals/${declineRow.id}/decline`, {
        method: "POST",
        body: { comment: declineComment },
      });
      showDecline = false;
      await refresh();
    } catch (error) {
      declineError = getApiErrorMessage(error, "Decline failed");
    } finally {
      declining = false;
    }
  }

  function openDelete(row: AbsorbProposalRow) {
    deleteRow = row;
    deleteError = null;
    showDelete = true;
  }

  async function deleteProposal() {
    if (!deleteRow) return;
    const id = deleteRow.id;
    deleting = true;
    deleteError = null;
    try {
      await pb.send(`/api/absorb_proposals/${id}`, { method: "DELETE" });
      deletedIds = [...deletedIds, id];
      showDelete = false;
      await onchange();
    } catch (error) {
      deleteError = getApiErrorMessage(error, "Delete failed");
    } finally {
      deleting = false;
    }
  }

  async function withdraw(row: AbsorbProposalRow) {
    try {
      await pb.send(`/api/absorb_proposals/${row.id}/withdraw`, { method: "POST" });
      await refresh();
    } catch (error) {
      globalStore.addError(getApiErrorMessage(error, "Withdraw failed"));
    }
  }
</script>

<div>
  <DSTabBar
    tabs={[
      { label: "Pending", href: resolve("/absorb/proposals"), active: pendingTab },
      {
        label: "Decided",
        href: `${resolve("/absorb/proposals")}?status=decided`,
        active: !pendingTab,
      },
    ]}
  />
  <DsList
    items={visibleItems}
    search={true}
    pagination={!pendingTab}
    bind:searchTerm
    inListHeader={pendingTab ? "Absorb Proposals Pending" : "Decided Absorb Proposals"}
    searchText={(row) =>
      [rowTitle(row), recordNames(row), row.requester_name, row.reason].join(" ")}
    emptyMessage={searchTerm.trim() !== ""
      ? "No matching proposals."
      : pendingTab
        ? "No absorb proposals are waiting."
        : "No decided absorb proposals."}
  >
    {#snippet anchor(row: AbsorbProposalRow)}
      {absorbProposalCollectionLabel(row.collection_name)}
    {/snippet}
    {#snippet headline(row: AbsorbProposalRow)}
      {rowTitle(row)}
    {/snippet}
    {#snippet byline(row: AbsorbProposalRow)}
      <span class="flex flex-wrap items-center gap-2">
        <DSStatusLabel status={absorbProposalStatusLabel(row.status)} kind="proposal" />
        {#if row.in_progress}
          <DsLabel tone="info">Absorbed, waiting for commit</DsLabel>
        {/if}
      </span>
    {/snippet}
    {#snippet line1(row: AbsorbProposalRow)}
      Absorb: {recordNames(row)}
    {/snippet}
    {#snippet line2(row: AbsorbProposalRow)}
      {row.reason}
    {/snippet}
    {#snippet line3(row: AbsorbProposalRow)}
      Proposed by {row.requester_name} on {shortDate(row.created.slice(0, 10), true)}
      {#if row.decided_at}
        · {absorbProposalStatusLabel(row.status)}
        {#if row.decided_by_name}by {row.decided_by_name}{/if}
        on {shortDate(row.decided_at.slice(0, 10), true)}
      {/if}
      {#if row.decision_comment}
        · {row.decision_comment}
      {/if}
    {/snippet}
    {#snippet actions(row: AbsorbProposalRow)}
      {#if row.status === "pending" && !row.in_progress}
        {#if canReview}
          <DSActionButton
            action={absorbProposalReviewHref(row)}
            icon="mdi:merge"
            title="Review"
            color="yellow"
          />
          <DSActionButton
            action={() => openDecline(row)}
            icon="mdi:cancel"
            title="Decline"
            color="orange"
          />
        {/if}
        {#if row.requester === viewerId}
          <DSActionButton
            action={() => withdraw(row)}
            icon="mdi:close"
            title="Withdraw"
            color="neutral"
          />
        {/if}
      {:else if row.status === "withdrawn" && row.requester === viewerId}
        <DSActionButton
          action={() => openDelete(row)}
          icon="mdi:delete"
          title="Delete"
          color="red"
        />
      {/if}
    {/snippet}
  </DsList>
</div>

<DSPopover
  bind:show={showDecline}
  title="Decline Proposal"
  subtitle={declineRow ? rowTitle(declineRow) : ""}
  error={declineError}
  submitting={declining}
  submitLabel="Decline"
  submitDisabled={declineComment.trim() === ""}
  onSubmit={decline}
>
  <DSTextInput
    bind:value={declineComment}
    errors={{}}
    fieldName="comment"
    uiName="Why"
    placeholder="This goes to the requester"
  />
</DSPopover>

<DSPopover
  bind:show={showDelete}
  title="Delete Proposal?"
  subtitle={deleteRow ? rowTitle(deleteRow) : ""}
  error={deleteError}
  submitting={deleting}
  danger
  submitLabel="Delete"
  cancelLabel="Keep it"
  onSubmit={deleteProposal}
>
  <p>Turbo removes this withdrawn proposal. The records it names don't change.</p>
  <DSCallout tone="danger" title="This is permanent">You can't undo this.</DSCallout>
</DSPopover>
