<script lang="ts">
  import type { PageData } from "./$types";
  import { pb } from "$lib/pocketbase";
  import { globalStore } from "$lib/stores/global";
  import DsActionButton from "$lib/components/DSActionButton.svelte";
  import DsLabel from "$lib/components/DsLabel.svelte";
  import DSField from "$lib/components/DSField.svelte";
  import DSFieldList from "$lib/components/DSFieldList.svelte";
  import DSPageHeader from "$lib/components/DSPageHeader.svelte";
  import DSCard from "$lib/components/DSCard.svelte";
  import DSCardGrid from "$lib/components/DSCardGrid.svelte";
  import DSCallout from "$lib/components/DSCallout.svelte";
  import DSStatusLabel from "$lib/components/DSStatusLabel.svelte";
  import VendorInvoiceParts from "$lib/components/expenses/VendorInvoiceParts.svelte";
  import {
    loadVendorInvoice,
    type VendorInvoiceView,
  } from "$lib/components/expenses/vendorInvoice";
  import StoredFileHashRepairPopover from "$lib/components/StoredFileHashRepairPopover.svelte";
  import RejectModal from "$lib/components/RejectModal.svelte";
  import UncommitConfirmPopover from "$lib/components/UncommitConfirmPopover.svelte";
  import ExpenseUncommitRequestPopover from "$lib/components/ExpenseUncommitRequestPopover.svelte";
  import {
    expenseUncommitBlockText,
    expenseUncommitOutcomeLabel,
    expenseUncommitReasonText,
  } from "$lib/expenseUncommitRequests";
  import { advanceExpenseApprovalQueue } from "$lib/expenseApprovalQueue";
  import { expenseTitle } from "$lib/expenseTitle";
  import { getApiErrorMessage } from "$lib/errors";
  import {
    formatCurrencyAmount,
    formatCurrencyEquivalent,
    openExpenseAttachment,
    shortDate,
    trimmedOrEmpty,
  } from "$lib/utilities";
  import { goto } from "$app/navigation";
  import { resolve } from "$app/paths";
  import { expensesEditingEnabled } from "$lib/stores/appConfig";

  export let data: PageData;
  const viewerId = pb.authStore.record?.id ?? "";
  let rejectModal: RejectModal;
  let expense = data.expense;
  // Moving through the approval queue reuses this component, so follow data.
  $: expense = data.expense;
  let hasCommitAccess = false;
  let hasAdminAccess = false;
  let hasBookKeeperAccess = false;
  let isOwner = false;
  let isApprover = false;
  let showUncommitConfirm = false;
  let showHashRepairPopover = false;
  let uncommitSubmitting = false;
  let uncommitError: string | null = null;
  let showUncommitRequest = false;

  $: uncommitRequest = expense.uncommit_request;
  $: holdsCommitClaim = $globalStore.claims.includes("commit");
  $: holdsAdminClaim = $globalStore.claims.includes("admin");
  $: hasCommitAccess = $globalStore.showAllUi || $globalStore.claims.includes("commit");
  $: hasAdminAccess = $globalStore.showAllUi || $globalStore.claims.includes("admin");
  $: hasBookKeeperAccess = $globalStore.claims.includes("book_keeper");
  $: isOwner = expense.creator === viewerId;
  $: isApprover = expense.approver === viewerId;
  $: hasAttachmentRepairTarget = Boolean(expense.attachment || expense.attachment_document);

  // The vendor invoice and every expense that shares it.
  let invoiceView: VendorInvoiceView | null = null;
  async function loadInvoice(id: string, number: string) {
    invoiceView = null;
    if (!number) return;
    try {
      invoiceView = await loadVendorInvoice(id);
    } catch {
      invoiceView = null;
    }
  }
  $: loadInvoice(expense.id, expense.vendor_invoice_number);

  async function refreshExpense() {
    try {
      expense = await pb.send(`/api/expenses/details/${expense.id}`, { method: "GET" });
    } catch (error: any) {
      globalStore.addError(error?.response?.error || "Failed to refresh expense details");
    }
  }

  async function submitExpense() {
    try {
      await pb.send(`/api/expenses/${expense.id}/submit`, { method: "POST" });
      await globalStore.refreshAttentionCounts();
      await refreshExpense();
    } catch (error: any) {
      globalStore.addError(error?.response?.error || "Submit failed");
    }
  }

  async function recallExpense() {
    try {
      await pb.send(`/api/expenses/${expense.id}/recall`, { method: "POST" });
      await globalStore.refreshAttentionCounts();
      await refreshExpense();
    } catch (error: any) {
      globalStore.addError(error?.response?.error || "Recall failed");
    }
  }

  async function approveExpense() {
    try {
      await pb.send(`/api/expenses/${expense.id}/approve`, { method: "POST" });
      await globalStore.refreshAttentionCounts();
    } catch (error: any) {
      globalStore.addError(error?.response?.error || "Approve failed");
      return;
    }
    await moveToNextPending("approved");
  }

  // Only an approver working through their pending queue moves on after a
  // reject; a commit holder rejecting an approved expense stays here.
  async function handleRejected() {
    if (isApprover && expense.approved === "") {
      await moveToNextPending("rejected");
    } else {
      await refreshExpense();
    }
  }

  async function moveToNextPending(verb: "approved" | "rejected") {
    try {
      await advanceExpenseApprovalQueue(expense, verb);
    } catch (error) {
      globalStore.addError(getApiErrorMessage(error, "Failed to load the next pending expense"));
      await refreshExpense();
    }
  }

  async function commitExpense() {
    try {
      await pb.send(`/api/expenses/${expense.id}/commit`, { method: "POST" });
      await globalStore.refreshAttentionCounts();
      await goto(resolve("/reports/expense/queue"));
    } catch (error: any) {
      globalStore.addError(error?.response?.error || "Commit failed");
    }
  }

  async function uncommitExpense() {
    uncommitSubmitting = true;
    uncommitError = null;
    try {
      await pb.send(`/api/expenses/${expense.id}/uncommit`, { method: "POST" });
      showUncommitConfirm = false;
      await globalStore.refreshAttentionCounts();
      await refreshExpense();
    } catch (error: any) {
      uncommitError = error?.response?.error || "Uncommit failed";
    } finally {
      uncommitSubmitting = false;
    }
  }

  async function deleteExpense() {
    try {
      await pb.collection("expenses").delete(expense.id);
      goto(resolve("/expenses/list"));
    } catch (error: any) {
      globalStore.addError(error?.response?.message || "Delete failed");
    }
  }

  function openRejectModal() {
    rejectModal?.openModal(expense.id);
  }

  function openUncommitConfirm() {
    uncommitError = null;
    showUncommitConfirm = true;
  }

  function closeUncommitConfirm() {
    uncommitError = null;
    showUncommitConfirm = false;
  }

  async function uncommitRequested() {
    await globalStore.refreshAttentionCounts();
    await refreshExpense();
  }

  async function cancelUncommitRequest(id: string) {
    try {
      await pb.send(`/api/expense_uncommit_requests/${id}/cancel`, { method: "POST" });
      await uncommitRequested();
    } catch (error) {
      globalStore.addError(getApiErrorMessage(error, "Cancel failed"));
    }
  }

  function openHashRepairPopover() {
    showHashRepairPopover = true;
  }

  function closeHashRepairPopover() {
    showHashRepairPopover = false;
  }

  function poDiffers(expenseVal: string, poVal: string): boolean {
    return !!expense.purchase_order && !!poVal && expenseVal !== poVal;
  }

  function paymentTypeLabel(pt: string) {
    switch (pt) {
      case "CorporateCreditCard":
        return "Corporate Credit Card";
      case "OnAccount":
        return "On Account";
      case "PersonalReimbursement":
        return "Personal Reimbursement";
      default:
        return pt;
    }
  }

  function personLabel(name: string, uid: string, fallback: string): string {
    return trimmedOrEmpty(name) || trimmedOrEmpty(uid) || fallback;
  }

  // Expenses have no status field; derive the one the header shows.
  $: status = expense.committed
    ? "Committed"
    : expense.rejected
      ? "Rejected"
      : expense.approved
        ? "Approved"
        : expense.submitted
          ? "Submitted"
          : "Draft";
  $: pendingMyApproval =
    isApprover &&
    expense.submitted &&
    expense.approved === "" &&
    expense.rejected === "" &&
    expense.committed === "";
  $: hasSettledCadTotal = Number(expense.settled_total ?? 0) > 0;
  $: showIndicativeCadEquivalent = expense.currency_code !== "CAD" && !hasSettledCadTotal;
</script>

<!-- A field value that differs from the linked PO's shows the PO's value. -->
{#snippet poDiff(text: string)}
  <DsLabel tone="warning" size="sm">PO: {text}</DsLabel>
{/snippet}

<div class="mx-auto space-y-4 p-4">
  <DSPageHeader title={expenseTitle(expense)} subtitle={expense.description}>
    {#snippet badges()}
      <DSStatusLabel {status} kind="expense" />
    {/snippet}
    {#snippet context()}
      {expense.uid_name} · {shortDate(expense.date)}
      {#if pendingMyApproval}
        · Pending your approval
        {#if data.leftToReview > 1}
          · {data.leftToReview} left to review, including this one
        {:else if data.leftToReview === 1}
          · the last one
        {/if}
      {/if}
    {/snippet}
    {#snippet actions()}
      {#if $expensesEditingEnabled}
        {#if hasAdminAccess && expense.committed !== ""}
          <DsActionButton
            action={() => openUncommitConfirm()}
            icon="mdi:undo"
            title="Uncommit"
            color="orange"
          />
        {/if}
        <!-- An admin holder has Uncommit above, so a request would only ask themselves. -->
        {#if uncommitRequest?.can_request && !holdsAdminClaim}
          <DsActionButton
            action={() => (showUncommitRequest = true)}
            icon="mdi:undo"
            title="Request uncommit"
            color="orange"
          />
        {/if}

        <!-- Owner draft actions: an unsubmitted expense can still be edited and submitted. -->
        {#if isOwner && !expense.submitted}
          <DsActionButton
            action={`/expenses/${expense.id}/edit`}
            icon="mdi:pencil"
            title="Edit"
            color="blue"
          />
          <DsActionButton
            action={() => submitExpense()}
            icon="mdi:send"
            title="Submit"
            color="blue"
          />
        {/if}

        <!-- Owner recall: allowed for submitted-but-unapproved expenses and rejected expenses until committed. -->
        {#if isOwner && ((expense.submitted && expense.approved === "") || expense.rejected !== "") && expense.committed === ""}
          <DsActionButton
            action={() => recallExpense()}
            icon="mdi:rewind"
            title="Recall"
            color="orange"
          />
        {/if}

        <!-- Owner delete: only available while the expense is still a draft. -->
        {#if isOwner && !expense.submitted && expense.committed === ""}
          <DsActionButton
            action={() => deleteExpense()}
            icon="mdi:delete"
            title="Delete"
            color="red"
          />
        {/if}

        <!-- Approver approve: only the assigned approver can approve a submitted, unrejected, uncommitted expense. -->
        {#if pendingMyApproval}
          <DsActionButton
            action={() => approveExpense()}
            icon="mdi:approve"
            title="Approve"
            color="green"
          />
        {/if}

        <!-- Reject behavior: approvers may reject any submitted, uncommitted expense; commit holders may reject approved ones. -->
        {#if (isApprover || (hasCommitAccess && expense.approved !== "")) && expense.submitted && expense.rejected === "" && expense.committed === ""}
          <DsActionButton
            action={() => openRejectModal()}
            icon="mdi:cancel"
            title="Reject"
            color="orange"
          />
        {/if}

        <!-- Commit behavior: commit holders can commit only once the submitted expense has been approved. -->
        {#if hasCommitAccess && expense.submitted && expense.approved !== "" && expense.rejected === "" && expense.committed === ""}
          <DsActionButton
            action={() => commitExpense()}
            icon="mdi:check-all"
            title="Commit"
            color="green"
          />
        {/if}
      {/if}
    {/snippet}
  </DSPageHeader>

  {#if uncommitRequest?.pending}
    <DSCallout tone="info" title="Uncommit requested">
      {uncommitRequest.pending.requester_name} asked for
      {expenseUncommitOutcomeLabel(uncommitRequest.pending.outcome).toLowerCase()}.
      {expenseUncommitReasonText(uncommitRequest.pending)}.
      <a class="link" href={resolve("/expenses/uncommit-requests")}>View the queue</a>
      {#snippet actions()}
        {#if uncommitRequest?.pending?.requester === viewerId && $expensesEditingEnabled}
          <DsActionButton action={() => cancelUncommitRequest(uncommitRequest.pending.id)}>
            Cancel request
          </DsActionButton>
        {/if}
      {/snippet}
    </DSCallout>
  {:else if uncommitRequest?.blocked && holdsCommitClaim}
    <DSCallout tone="warning" title="Uncommit can't be requested here">
      {expenseUncommitBlockText(uncommitRequest.blocked)} Message Dean on Teams.
    </DSCallout>
  {/if}

  {#if expense.po_owner_uid_mismatch}
    <DSCallout tone="warning" title="UID mismatch with linked purchase order">
      This expense was submitted by
      {personLabel(expense.uid_name, expense.uid, "unknown submitter")}, but the linked PO
      {#if expense.purchase_order_number}
        {expense.purchase_order_number}
      {/if}
      is owned by
      {personLabel(expense.po_uid_name, expense.po_uid, "unknown PO owner")}.
    </DSCallout>
  {/if}

  {#if expense.rejected}
    <DSCallout
      tone="danger"
      title={`Rejected ${shortDate(expense.rejected, true)} by ${expense.rejector_name}`}
    >
      {expense.rejection_reason}
    </DSCallout>
  {/if}

  <DSCardGrid>
    <DSCard title="Expense">
      <DSFieldList>
        <DSField label="Date" value={shortDate(expense.date, true)} />
        {#if expense.creator && expense.creator !== expense.uid}
          <DSField label="Submitted for" value={expense.uid_name} />
          <DSField label="Entered by" value={expense.creator_name} />
        {:else}
          <DSField label="Submitted by" value={expense.uid_name} />
        {/if}
        <!-- The header shows the description; repeat it here only to flag a PO mismatch. -->
        {#if poDiffers(expense.description, expense.po_description)}
          <DSField label="Description">
            {expense.description}
            {@render poDiff(expense.po_description)}
          </DSField>
        {/if}
        <DSField label="Total">
          {formatCurrencyAmount(expense.total, expense.currency_code)}
          {#if expense.purchase_order && expense.po_total && expense.total !== expense.po_total}
            {@render poDiff(
              formatCurrencyAmount(
                expense.po_total,
                expense.po_currency_code || expense.currency_code,
              ),
            )}
          {/if}
        </DSField>
        {#if showIndicativeCadEquivalent}
          <DSField
            label="Indicative CAD equivalent"
            value={formatCurrencyEquivalent(
              expense.total * expense.currency_rate,
              expense.currency_rate,
              expense.currency_rate_date,
            )}
          />
        {/if}
        {#if hasSettledCadTotal}
          <DSField
            label="Settled total"
            value={formatCurrencyAmount(expense.settled_total, "CAD")}
          />
        {/if}
        {#if expense.payment_type}
          <DSField label="Payment type">
            {paymentTypeLabel(expense.payment_type)}
            {#if poDiffers(expense.payment_type, expense.po_payment_type)}
              {@render poDiff(paymentTypeLabel(expense.po_payment_type))}
            {/if}
          </DSField>
        {/if}
        {#if expense.cc_last_4_digits}
          <DSField label="Card" value={`****${expense.cc_last_4_digits}`} />
        {/if}
        {#if expense.payment_type === "Mileage"}
          <DSField label="Distance" value={`${expense.distance} km`} />
        {/if}
        {#if expense.allowance_types && expense.allowance_types !== "[]"}
          <DSField label="Allowance types" value={expense.allowance_types} />
        {/if}
      </DSFieldList>
    </DSCard>

    <DSCard title="Coding">
      <DSFieldList>
        {#if expense.vendor}
          <DSField label="Vendor" hint={trimmedOrEmpty(expense.vendor_alias) || undefined}>
            <a class="link" href={resolve(`/vendors/${expense.vendor}/details`)}
              >{expense.vendor_name}</a
            >
            {#if poDiffers(expense.vendor, expense.po_vendor)}
              {@render poDiff(
                trimmedOrEmpty(expense.po_vendor_alias)
                  ? `${expense.po_vendor_name} (${trimmedOrEmpty(expense.po_vendor_alias)})`
                  : expense.po_vendor_name,
              )}
            {/if}
          </DSField>
        {/if}
        {#if expense.purchase_order}
          <DSField
            label="Purchase order"
            value={expense.purchase_order_number}
            href={resolve(`/pos/${expense.purchase_order}/details`)}
          />
        {/if}
        {#if expense.job}
          <DSField label="Job" hint={expense.job_description || undefined}>
            <a class="link" href={resolve(`/jobs/${expense.job}/details`)}>{expense.job_number}</a>
            {#if poDiffers(expense.job, expense.po_job)}
              {@render poDiff(
                expense.po_job_description
                  ? `${expense.po_job_number} — ${expense.po_job_description}`
                  : expense.po_job_number,
              )}
            {/if}
          </DSField>
        {/if}
        <DSField label="Client" value={expense.client_name} />
        {#if expense.division_code}
          <DSField label="Division">
            {expense.division_code} — {expense.division_name}
            {#if poDiffers(expense.division, expense.po_division)}
              {@render poDiff(`${expense.po_division_code} — ${expense.po_division_name}`)}
            {/if}
          </DSField>
        {/if}
        {#if expense.category_name}
          <DSField label="Category">
            {expense.category_name}
            {#if poDiffers(expense.category, expense.po_category)}
              {@render poDiff(expense.po_category_name)}
            {/if}
          </DSField>
        {/if}
        {#if expense.kind_name}
          <DSField label="Kind">
            {expense.kind_name}
            {#if poDiffers(expense.kind, expense.po_kind)}
              {@render poDiff(expense.po_kind_name)}
            {/if}
          </DSField>
        {/if}
        {#if expense.branch_name}
          <DSField label="Branch">
            {expense.branch_name}
            {#if poDiffers(expense.branch_name, expense.po_branch_name)}
              {@render poDiff(expense.po_branch_name)}
            {/if}
          </DSField>
        {/if}
      </DSFieldList>
    </DSCard>

    <DSCard title="Status">
      <DSFieldList>
        <DSField label="Submitted" value={expense.submitted ? "Yes" : "No"} />
        {#if expense.approved}
          <DSField
            label="Approved"
            value={`${shortDate(expense.approved, true)} by ${expense.approver_name}`}
          />
        {/if}
        {#if expense.committed}
          <DSField
            label="Committed"
            value={shortDate(expense.committed, true)}
            hint={expense.committed_week_ending
              ? `week ending ${shortDate(expense.committed_week_ending, true)}`
              : undefined}
          />
        {/if}
        {#if expense.pay_period_ending}
          <DSField label="Pay period ending" value={shortDate(expense.pay_period_ending, true)} />
        {/if}
      </DSFieldList>
    </DSCard>

    {#if expense.attachment || expense.attachment_missing_reason || (hasAdminAccess && expense.attachment_document) || expense.vendor_invoice_number}
      <DSCard title="Attachment" span="full">
        {#if expense.attachment}
          <div class="flex flex-wrap items-center gap-2">
            <DsActionButton
              action={() => openExpenseAttachment(expense.id, expense.attachment)}
              icon="mdi:download"
            >
              Download
            </DsActionButton>
            {#if hasBookKeeperAccess}
              <DsActionButton
                action={`/pos/search?source_expense=${encodeURIComponent(expense.id)}`}
                icon="mdi:file-document-plus-outline"
                title="Create another expense with this attachment"
                color="green"
              />
            {/if}
            {#if hasAdminAccess}
              <DsActionButton
                variant="plain"
                action={openHashRepairPopover}
                title="Audit attachment hash"
              >
                {expense.attachment_hash ? expense.attachment_hash.slice(0, 8) : "No hash"}
              </DsActionButton>
            {:else if expense.attachment_hash}
              <DsLabel size="sm">{expense.attachment_hash.slice(0, 8)}</DsLabel>
            {/if}
            {#if expense.attachment_hash}
              <DsActionButton
                variant="plain"
                icon="mdi:content-copy"
                title="Copy full hash"
                action={() => navigator.clipboard.writeText(expense.attachment_hash)}
              />
            {/if}
          </div>
        {:else if expense.attachment_missing_reason}
          <DSCallout tone="warning" title="Missing historical attachment">
            {expense.attachment_missing_reason}
            {#snippet actions()}
              {#if hasAdminAccess}
                <DsActionButton variant="secondary" action={openHashRepairPopover}
                  >Review</DsActionButton
                >
              {/if}
            {/snippet}
          </DSCallout>
        {:else if hasAdminAccess && expense.attachment_document}
          <DSCallout tone="danger" title="Attachment document issue">
            The linked document has no downloadable attachment.
            {#snippet actions()}
              <DsActionButton variant="secondary" action={openHashRepairPopover}
                >Review</DsActionButton
              >
            {/snippet}
          </DSCallout>
        {/if}

        {#if expense.vendor_invoice_number}
          <!-- Vendor invoice details are stored on the attached document. -->
          <DSFieldList columns={2}>
            <DSField label="Vendor invoice" value={expense.vendor_invoice_number} />
            <DSField
              label="Invoice date"
              value={expense.vendor_invoice_date
                ? shortDate(expense.vendor_invoice_date, true)
                : ""}
            />
            <DSField
              label="Invoice total"
              value={formatCurrencyAmount(expense.vendor_invoice_total, expense.currency_code)}
            />
          </DSFieldList>
          {#if invoiceView && invoiceView.parts.length > 0}
            <VendorInvoiceParts invoice={invoiceView} currentExpenseId={expense.id} />
          {/if}
        {/if}
      </DSCard>
    {/if}
  </DSCardGrid>

  <RejectModal collectionName="expenses" bind:this={rejectModal} on:refresh={handleRejected} />
  <ExpenseUncommitRequestPopover
    bind:show={showUncommitRequest}
    expenseId={expense.id}
    onRequested={uncommitRequested}
  />
  <UncommitConfirmPopover
    bind:show={showUncommitConfirm}
    recordLabel="expense"
    pendingRequest={uncommitRequest?.pending}
    submitting={uncommitSubmitting}
    error={uncommitError}
    onSubmit={uncommitExpense}
    onCancel={closeUncommitConfirm}
  />
  <StoredFileHashRepairPopover
    show={showHashRepairPopover}
    recordId={expense.id}
    hasAttachment={hasAttachmentRepairTarget}
    currentHash={expense.attachment_hash}
    currentUpdated={expense.updated}
    missingReason={expense.attachment_missing_reason}
    auditPath={`/api/expenses/${expense.id}/attachment_hash/audit`}
    replacePath={`/api/expenses/${expense.id}/attachment_hash/replace`}
    markMissingPath={`/api/expenses/${expense.id}/attachment_missing/mark`}
    onClose={closeHashRepairPopover}
    onRepaired={refreshExpense}
  />
</div>
