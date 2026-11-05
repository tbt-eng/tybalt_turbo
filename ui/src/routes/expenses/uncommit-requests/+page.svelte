<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { resolve } from "$app/paths";
  import { pb } from "$lib/pocketbase";
  import { globalStore } from "$lib/stores/global";
  import { expensesEditingEnabled } from "$lib/stores/appConfig";
  import { getApiErrorMessage } from "$lib/errors";
  import { allExpensePaymentTypeOptions } from "$lib/expensePaymentTypes";
  import {
    daysAgoLabel,
    daysSince,
    expenseUncommitBlockLabel,
    expenseUncommitBlockText,
    expenseUncommitOutcomeLabel,
    expenseUncommitReasonText,
    expenseUncommitStatusLabel,
    isPayrollPaymentType,
  } from "$lib/expenseUncommitRequests";
  import { formatCurrencyAmount, shortDate } from "$lib/utilities";
  import type { ExpenseUncommitRequestRow } from "$lib/svelte-types";
  import DSActionButton from "$lib/components/DSActionButton.svelte";
  import DSCallout from "$lib/components/DSCallout.svelte";
  import DsCheck from "$lib/components/DsCheck.svelte";
  import DsEditingDisabledBanner from "$lib/components/DsEditingDisabledBanner.svelte";
  import DsLabel from "$lib/components/DsLabel.svelte";
  import DsList from "$lib/components/DSList.svelte";
  import DSPopover from "$lib/components/DSPopover.svelte";
  import DSStatusLabel from "$lib/components/DSStatusLabel.svelte";
  import DSTabBar from "$lib/components/DSTabBar.svelte";
  import DSTextInput from "$lib/components/DSTextInput.svelte";
  import type { PageData } from "./$types";

  // Past this the queue marks an expense as old, though it can still be approved.
  const OLD_COMMIT_DAYS = 30;

  let { data }: { data: PageData } = $props();

  const viewerId = pb.authStore.record?.id ?? "";
  const pendingTab = $derived(data.status === "pending");
  const canDecide = $derived($globalStore.claims.includes("admin") && $expensesEditingEnabled);
  const approvable = $derived(data.items.filter((row) => !row.blocked));

  // Approve all acts on the whole queue, so it is off while a search hides rows.
  let searchTerm = $state("");
  let selected = $state<Record<string, boolean>>({});
  const selectedIds = $derived(approvable.filter((row) => selected[row.id]).map((row) => row.id));

  let showApprove = $state(false);
  let approveIds = $state<string[]>([]);
  let approving = $state(false);
  let approveError = $state<string | null>(null);
  let failures = $state<{ id: string; label: string; message: string }[]>([]);
  const approveRows = $derived(data.items.filter((row) => approveIds.includes(row.id)));
  const rejectCount = $derived(
    approveRows.filter((row) => row.outcome === "uncommit_and_reject").length,
  );

  let declineRow = $state<ExpenseUncommitRequestRow | null>(null);
  let showDecline = $state(false);
  let declineComment = $state("");
  let declining = $state(false);
  let declineError = $state<string | null>(null);

  function paymentTypeLabel(paymentType: string) {
    return (
      allExpensePaymentTypeOptions.find((option) => option.id === paymentType)?.name ?? paymentType
    );
  }

  function rowTitle(row: ExpenseUncommitRequestRow) {
    const title = row.po_number ? `PO ${row.po_number}` : row.description || "No PO";
    return row.vendor_name ? `${title} · ${row.vendor_name}` : title;
  }

  function openApprove(ids: string[]) {
    approveIds = ids;
    approveError = null;
    showApprove = true;
  }

  async function refresh() {
    await globalStore.refreshAttentionCounts();
    await invalidateAll();
  }

  async function approve() {
    approving = true;
    approveError = null;
    try {
      const response = (await pb.send("/api/expense_uncommit_requests/approve", {
        method: "POST",
        body: { ids: approveIds },
      })) as { results: { id: string; ok: boolean; message?: string }[] };
      failures = response.results
        .filter((result) => !result.ok)
        .map((result) => {
          const row = data.items.find((candidate) => candidate.id === result.id);
          return {
            id: result.id,
            label: row
              ? `${rowTitle(row)}, ${formatCurrencyAmount(row.total, row.currency_code)}`
              : result.id,
            message: result.message ?? "Approval failed",
          };
        });
      selected = {};
      showApprove = false;
      await refresh();
    } catch (error) {
      approveError = getApiErrorMessage(error, "Approval failed");
    } finally {
      approving = false;
    }
  }

  function openDecline(row: ExpenseUncommitRequestRow) {
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
      await pb.send(`/api/expense_uncommit_requests/${declineRow.id}/decline`, {
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

  async function cancelRequest(row: ExpenseUncommitRequestRow) {
    try {
      await pb.send(`/api/expense_uncommit_requests/${row.id}/cancel`, { method: "POST" });
      await refresh();
    } catch (error) {
      globalStore.addError(getApiErrorMessage(error, "Cancel failed"));
    }
  }
</script>

<svelte:head>
  <title>Uncommit Requests</title>
</svelte:head>

<div class="space-y-4">
  {#if !$expensesEditingEnabled}
    <DsEditingDisabledBanner message="Expense editing is disabled, so requests can't be changed." />
  {/if}

  {#if failures.length > 0}
    <DSCallout tone="danger" title="These requests weren't approved">
      <ul class="list-disc pl-5">
        {#each failures as failure (failure.id)}
          <li>{failure.label}: {failure.message}</li>
        {/each}
      </ul>
    </DSCallout>
  {/if}

  <div>
    <DSTabBar
      tabs={[
        {
          label: "Pending",
          href: resolve("/expenses/uncommit-requests"),
          active: pendingTab,
        },
        {
          label: "Decided",
          href: `${resolve("/expenses/uncommit-requests")}?status=decided`,
          active: !pendingTab,
        },
      ]}
    />
    <DsList
      items={data.items}
      search={true}
      pagination={!pendingTab}
      bind:searchTerm
      inListHeader={pendingTab ? "Uncommit Requests Pending" : "Decided Uncommit Requests"}
      searchText={(row) =>
        [rowTitle(row), row.employee_name, row.requester_name, row.reason_name, row.note].join(" ")}
      emptyMessage={searchTerm.trim() !== ""
        ? "No matching requests."
        : pendingTab
          ? "No uncommit requests are waiting."
          : "No uncommit requests yet."}
    >
      {#snippet searchBarExtra()}
        {#if pendingTab && canDecide}
          <div class="flex items-center gap-2 max-[639px]:w-full max-[639px]:flex-wrap">
            <DSActionButton
              color="green"
              disabled={selectedIds.length === 0}
              action={() => openApprove(selectedIds)}
            >
              Approve selected ({selectedIds.length})
            </DSActionButton>
            <DSActionButton
              color="green"
              disabled={approvable.length === 0 || searchTerm.trim() !== ""}
              title={searchTerm.trim() !== "" ? "Clear the search to approve all" : undefined}
              action={() => openApprove(approvable.map((row) => row.id))}
            >
              Approve all ({approvable.length})
            </DSActionButton>
          </div>
        {/if}
      {/snippet}
      {#snippet anchor(row: ExpenseUncommitRequestRow)}
        {@const label = row.po_number || shortDate(row.committed.slice(0, 10), true)}
        <span class="flex items-center gap-2 whitespace-nowrap">
          {#if pendingTab && canDecide}
            <span>
              <DsCheck
                bind:value={selected[row.id]}
                errors={{}}
                fieldName={`select-${row.id}`}
                uiName={`Select ${rowTitle(row)}`}
                disabled={!!row.blocked}
                hideLabel={true}
              />
            </span>
          {/if}
          {#if row.expense}
            <a class="link" href={resolve(`/expenses/${row.expense}/details`)}>{label}</a>
          {:else}
            {label}
          {/if}
        </span>
      {/snippet}
      {#snippet headline(row: ExpenseUncommitRequestRow)}
        {row.vendor_name || row.description || row.employee_name}
      {/snippet}
      {#snippet byline(row: ExpenseUncommitRequestRow)}
        <span class="flex flex-wrap items-center gap-2">
          {formatCurrencyAmount(row.total, row.currency_code)}
          {#if !pendingTab}
            <DSStatusLabel status={expenseUncommitStatusLabel(row.status)} kind="uncommit" />
          {/if}
          <DsLabel tone={row.outcome === "uncommit_and_reject" ? "warning" : "neutral"}>
            {expenseUncommitOutcomeLabel(row.outcome)}
          </DsLabel>
          {#if row.blocked}
            <DsLabel tone="danger">{expenseUncommitBlockLabel(row.blocked)}</DsLabel>
          {:else if pendingTab && row.committed && daysSince(row.committed) > OLD_COMMIT_DAYS}
            <DsLabel tone="warning">Committed over {OLD_COMMIT_DAYS} days ago</DsLabel>
          {/if}
        </span>
      {/snippet}
      {#snippet line1(row: ExpenseUncommitRequestRow)}
        {[
          row.vendor_name ? row.description : "",
          row.employee_name,
          paymentTypeLabel(row.payment_type),
        ]
          .filter(Boolean)
          .join(" · ")}
      {/snippet}
      {#snippet line2(row: ExpenseUncommitRequestRow)}
        {expenseUncommitReasonText(row)}
      {/snippet}
      {#snippet line3(row: ExpenseUncommitRequestRow)}
        {#if row.blocked}
          <div>
            {expenseUncommitBlockText(row.blocked)} Confirm first, then use Uncommit on the expense.
          </div>
        {/if}
        {#if pendingTab}
          Requested by {row.requester_name}
          {daysAgoLabel(row.created)}
          {#if row.committed}
            · Committed {shortDate(row.committed.slice(0, 10), true)}
            ({daysAgoLabel(row.committed)})
          {/if}
          {#if row.pay_period_ending && isPayrollPaymentType(row.payment_type)}
            · Pay period {shortDate(row.pay_period_ending, true)}
          {/if}
        {:else}
          Requested by {row.requester_name}
          {#if row.decided_at}
            · {expenseUncommitStatusLabel(row.status)} by {row.decided_by_name} on {shortDate(
              row.decided_at.slice(0, 10),
              true,
            )}
          {/if}
          {#if row.decision_comment}
            · {row.decision_comment}
          {/if}
        {/if}
      {/snippet}
      {#snippet actions(row: ExpenseUncommitRequestRow)}
        {#if pendingTab && canDecide}
          <DSActionButton
            disabled={!!row.blocked}
            action={() => openApprove([row.id])}
            icon="mdi:approve"
            title={row.blocked ? expenseUncommitBlockLabel(row.blocked) : "Approve"}
            color="green"
          />
          <DSActionButton
            action={() => openDecline(row)}
            icon="mdi:cancel"
            title="Decline"
            color="orange"
          />
        {/if}
        {#if pendingTab && row.requester === viewerId && $expensesEditingEnabled}
          <DSActionButton
            action={() => cancelRequest(row)}
            icon="mdi:close"
            title="Cancel request"
            color="neutral"
          />
        {/if}
      {/snippet}
    </DsList>
  </div>
</div>

<DSPopover
  bind:show={showApprove}
  title={approveIds.length === 1 ? "Approve Request" : `Approve ${approveIds.length} Requests`}
  error={approveError}
  submitting={approving}
  submitLabel="Approve"
  onSubmit={approve}
>
  <DSCallout tone="warning" title="Review what will happen">
    <ul class="list-disc pl-5">
      <li>{rejectCount} will be uncommitted and rejected.</li>
      <li>
        {approveIds.length - rejectCount} will be uncommitted and go back to the commit queue.
      </li>
      <li>Reports change, and any legacy writeback must be cleaned up by hand.</li>
    </ul>
  </DSCallout>
</DSPopover>

<DSPopover
  bind:show={showDecline}
  title="Decline Request"
  subtitle={declineRow ? rowTitle(declineRow) : ""}
  error={declineError}
  submitting={declining}
  submitLabel="Decline"
  submitDisabled={declineComment.trim().length < 4}
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
