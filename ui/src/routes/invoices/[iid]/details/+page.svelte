<script lang="ts">
  import { resolve } from "$app/paths";
  import { goto, invalidateAll } from "$app/navigation";
  import DSPageHeader from "$lib/components/DSPageHeader.svelte";
  import DSActionButton from "$lib/components/DSActionButton.svelte";
  import DSCardGrid from "$lib/components/DSCardGrid.svelte";
  import DSCard from "$lib/components/DSCard.svelte";
  import DSFieldList from "$lib/components/DSFieldList.svelte";
  import DSField from "$lib/components/DSField.svelte";
  import DSSection from "$lib/components/DSSection.svelte";
  import DSList from "$lib/components/DSList.svelte";
  import DSPopover from "$lib/components/DSPopover.svelte";
  import DSMenuItem from "$lib/components/DSMenuItem.svelte";
  import DSStatusLabel from "$lib/components/DSStatusLabel.svelte";
  import AddBilledCostPopover from "$lib/components/invoices/AddBilledCostPopover.svelte";
  import { pb } from "$lib/pocketbase";
  import { globalStore } from "$lib/stores/global";
  import {
    formatCurrency,
    formatCurrencyAmount,
    formatDateTime,
    saveErrorMessage,
    shortDate,
  } from "$lib/utilities";
  import { costLabel, issueInvoice, type InvoiceLine } from "$lib/invoices";
  import type { PageData } from "./$types";

  let { data }: { data: PageData } = $props();

  const invoice = $derived(data.invoice);
  const isDraft = $derived(invoice.status === "Draft");
  const canEdit = $derived($globalStore.claims.includes("invoicing"));
  const billed = $derived(
    Math.round(data.lines.reduce((sum, line) => sum + line.amount * 100, 0)) / 100,
  );

  let showIssue = $state(false);
  let showDelete = $state(false);
  let showAdd = $state(false);
  let removing = $state<InvoiceLine | null>(null);
  let showRemove = $state(false);
  let actionError = $state<string | null>(null);
  let submitting = $state(false);

  async function run(action: () => Promise<unknown>, fallback: string, after?: () => void) {
    submitting = true;
    actionError = null;
    try {
      await action();
      after?.();
      await invalidateAll();
    } catch (err) {
      actionError = saveErrorMessage(err, fallback);
    } finally {
      submitting = false;
    }
  }

  function openPopover(set: () => void) {
    actionError = null;
    set();
  }
</script>

<div class="flex w-full flex-col gap-4 p-2">
  <DSPageHeader
    title={invoice.number ? `Invoice ${invoice.number}` : "Draft invoice"}
    subtitle={invoice.expand?.client?.name ?? ""}
  >
    {#snippet badges()}
      <DSStatusLabel status={invoice.status} kind="invoice" />
    {/snippet}
    {#snippet actions()}
      {#if canEdit && isDraft}
        <DSActionButton variant="secondary" action={`/invoices/${invoice.id}/edit`}>
          Edit
        </DSActionButton>
        <DSActionButton variant="primary" action={() => openPopover(() => (showIssue = true))}>
          Issue
        </DSActionButton>
      {:else if canEdit}
        <DSActionButton variant="primary" action={() => openPopover(() => (showAdd = true))}>
          Add billed-back cost
        </DSActionButton>
      {/if}
    {/snippet}
    {#snippet menu()}
      {#if canEdit && isDraft}
        <DSMenuItem
          label="Delete draft"
          icon="mdi:delete"
          tone="danger"
          action={() => openPopover(() => (showDelete = true))}
        />
      {/if}
    {/snippet}
  </DSPageHeader>

  <DSCardGrid>
    <DSCard title="Invoice">
      <DSFieldList>
        <DSField
          label="Client"
          value={invoice.expand?.client?.name}
          href={`/clients/${invoice.client}/details`}
        />
        <DSField label="Number" value={invoice.number} />
        <DSField label="Date" value={invoice.date ? shortDate(invoice.date, true) : ""} />
        <DSField label="Total" value={invoice.total ? formatCurrency(invoice.total) : ""} />
        <DSField
          label="Issued"
          value={invoice.issued_at ? formatDateTime(invoice.issued_at) : ""}
        />
      </DSFieldList>
    </DSCard>
    {#if !isDraft}
      <DSCard title="Billed-back costs">
        <DSFieldList>
          <DSField label="Billed on this invoice" value={formatCurrency(billed)} />
          <DSField label="Costs" value={data.lines.length} showEmpty />
        </DSFieldList>
      </DSCard>
    {/if}
  </DSCardGrid>

  {#if !isDraft}
    <DSSection title="Billed-back costs">
      {#if data.lines.length === 0}
        <p>No costs billed on this invoice yet.</p>
      {/if}
      <DSList items={data.lines}>
        {#snippet headline(line: InvoiceLine)}
          {#if line.can_open}
            <a
              class="link"
              href={line.kind === "po"
                ? resolve(`/pos/${line.cost_id}/details`)
                : resolve(`/expenses/${line.cost_id}/details`)}>{costLabel(line)}</a
            >
          {:else}
            {costLabel(line)}
          {/if}
          · {formatCurrency(line.amount)} billed
        {/snippet}
        {#snippet line1(line: InvoiceLine)}
          {line.job_number}
          {#if line.vendor_name}· {line.vendor_name}{/if}
          {#if line.date}· {shortDate(line.date, true)}{/if}
          · cost {formatCurrencyAmount(line.cost_total, line.currency_code)}
        {/snippet}
        {#snippet line2(line: InvoiceLine)}
          Added by {line.creator_name || "unknown"} {formatDateTime(line.created)}
        {/snippet}
        {#snippet actions(line: InvoiceLine)}
          {#if canEdit}
            <DSActionButton
              icon="mdi:delete"
              title="Remove"
              color="red"
              action={() =>
                openPopover(() => {
                  removing = line;
                  showRemove = true;
                })}
            />
          {/if}
        {/snippet}
      </DSList>
    </DSSection>
  {:else}
    <p>Costs are billed back once the invoice is issued. Issue it when it is sent to the client.</p>
  {/if}
</div>

<DSPopover
  bind:show={showIssue}
  title="Issue this invoice?"
  subtitle="Issue it once it is sent to the client. An issued invoice cannot be edited or deleted."
  submitLabel="Issue"
  error={actionError}
  {submitting}
  onSubmit={() =>
    run(
      () => issueInvoice(invoice.id),
      "Could not issue the invoice.",
      () => (showIssue = false),
    )}
>
  <DSFieldList>
    <DSField label="Number" value={invoice.number} showEmpty />
    <DSField label="Date" value={invoice.date ? shortDate(invoice.date, true) : ""} showEmpty />
    <DSField label="Total" value={invoice.total ? formatCurrency(invoice.total) : ""} showEmpty />
  </DSFieldList>
</DSPopover>

<DSPopover
  bind:show={showDelete}
  title="Delete this draft?"
  subtitle="Its number becomes free for another invoice."
  submitLabel="Delete"
  danger
  error={actionError}
  {submitting}
  onSubmit={() =>
    run(
      () => pb.collection("invoices").delete(invoice.id),
      "Could not delete the draft.",
      () => goto(resolve("/invoices/list")),
    )}
>
  <p>{invoice.number ? `Invoice ${invoice.number}` : "This draft"} will be deleted.</p>
</DSPopover>

<DSPopover
  bind:show={showRemove}
  title="Remove this cost?"
  subtitle="It can then be billed again, on this or another invoice."
  submitLabel="Remove"
  danger
  error={actionError}
  {submitting}
  onSubmit={() =>
    run(
      () => pb.collection("invoice_lines").delete(removing?.id ?? ""),
      "Could not remove the cost.",
      () => (showRemove = false),
    )}
>
  <p>{removing ? costLabel(removing) : ""}</p>
</DSPopover>

{#if !isDraft && canEdit}
  <AddBilledCostPopover
    bind:show={showAdd}
    invoiceId={invoice.id}
    onAdded={() => invalidateAll()}
  />
{/if}
