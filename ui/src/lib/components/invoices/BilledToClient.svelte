<script lang="ts">
  import { resolve } from "$app/paths";
  import { pb } from "$lib/pocketbase";
  import DSCard from "$lib/components/DSCard.svelte";
  import DSFieldList from "$lib/components/DSFieldList.svelte";
  import DSField from "$lib/components/DSField.svelte";
  import DSActionButton from "$lib/components/DSActionButton.svelte";
  import DSPopover from "$lib/components/DSPopover.svelte";
  import DSSelector from "$lib/components/DSSelector.svelte";
  import DSTextInput from "$lib/components/DSTextInput.svelte";
  import {
    formatCurrency,
    formatCurrencyAmount,
    saveErrorMessage,
    shortDate,
  } from "$lib/utilities";
  import { billCost, loadCostInvoice, type CostInvoice } from "$lib/invoices";
  import type { InvoicesResponse } from "$lib/pocketbase-types";

  // The client invoice an expense or PO was billed back on. Anyone who can
  // open the record sees it; an invoicing holder can bill an unbilled,
  // eligible record on one of the client's issued invoices. Shows nothing
  // when the record is not billed and the viewer cannot bill it.
  let { kind, id }: { kind: "expense" | "po"; id: string } = $props();

  let view = $state<CostInvoice | null>(null);
  let invoices = $state([] as InvoicesResponse[]);
  let invoiceId = $state("");
  let amount = $state("" as string | number);
  let showBill = $state(false);
  let error = $state<string | null>(null);
  let submitting = $state(false);

  async function load() {
    try {
      view = await loadCostInvoice(kind, id);
    } catch {
      view = null;
    }
  }

  $effect(() => {
    load();
  });

  async function openBill() {
    if (!view) return;
    error = null;
    invoiceId = "";
    amount = "";
    showBill = true;
    try {
      invoices = await pb.collection("invoices").getFullList({
        filter: pb.filter("client = {:client} && status = 'Issued'", { client: view.client_id }),
        sort: "-date,-number",
      });
    } catch (err) {
      error = saveErrorMessage(err, "Could not load the client's invoices.");
    }
  }

  async function bill() {
    submitting = true;
    error = null;
    try {
      await billCost(invoiceId, kind, id, Number(amount));
      showBill = false;
      await load();
    } catch (err) {
      error = saveErrorMessage(err, "The cost could not be billed.");
    } finally {
      submitting = false;
    }
  }
</script>

{#if view && (view.line || view.can_bill)}
  <DSCard title="Billed to client">
    {#if view.line}
      <DSFieldList>
        <DSField
          label="Invoice"
          value={view.line.invoice_number}
          href={view.can_open_invoice
            ? resolve(`/invoices/${view.line.invoice_id}/details`)
            : undefined}
        />
        <DSField
          label="Invoice date"
          value={view.line.invoice_date ? shortDate(view.line.invoice_date, true) : ""}
        />
        <DSField label="Amount billed" value={formatCurrency(view.line.amount)} />
      </DSFieldList>
    {:else}
      <div class="flex flex-col items-start gap-2">
        <p>Not billed to the client yet.</p>
        <DSActionButton variant="primary" action={openBill}>Bill on invoice</DSActionButton>
      </div>
    {/if}
  </DSCard>

  <DSPopover
    bind:show={showBill}
    title="Bill on invoice"
    subtitle="Pick the client's issued invoice, and enter what was billed for this cost in Canadian dollars, before tax."
    submitLabel="Bill"
    {error}
    {submitting}
    submitDisabled={!invoiceId || !(Number(amount) > 0)}
    onSubmit={bill}
  >
    {#if invoices.length === 0 && !error}
      <p>The client has no issued invoices yet. Issue the invoice first.</p>
    {:else}
      <DSSelector
        bind:value={invoiceId}
        items={invoices}
        errors={{}}
        fieldName="invoice"
        uiName="Invoice"
      >
        {#snippet optionTemplate(invoice)}
          {invoice.number} · {invoice.date ? shortDate(invoice.date, true) : ""}
        {/snippet}
      </DSSelector>
      <DSFieldList>
        <DSField label="Cost" value={formatCurrencyAmount(view.cost_total, view.currency_code)} />
      </DSFieldList>
      <DSTextInput
        bind:value={amount}
        type="number"
        step={0.01}
        min={0}
        errors={{}}
        fieldName="amount"
        uiName="Amount billed"
      />
    {/if}
  </DSPopover>
{/if}
