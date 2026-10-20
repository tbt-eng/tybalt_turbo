<script lang="ts">
  import { resolve } from "$app/paths";
  import DSList from "$lib/components/DSList.svelte";
  import DSStatusLabel from "$lib/components/DSStatusLabel.svelte";
  import DSFieldList from "$lib/components/DSFieldList.svelte";
  import DSField from "$lib/components/DSField.svelte";
  import { formatCurrencyAmount, formatDateTime, shortDate } from "$lib/utilities";
  import { partStatus, type VendorInvoicePart, type VendorInvoiceView } from "./vendorInvoice";

  // Every expense on a vendor invoice, so the people splitting or approving it
  // can see which expenses share it and how much of it is accounted for. Parts
  // the viewer may not open are listed without a link.
  let {
    invoice,
    currentExpenseId = "",
  }: {
    invoice: VendorInvoiceView;
    // The expense being viewed or edited, marked in the list.
    currentExpenseId?: string;
  } = $props();

  const money = (amount: number) => formatCurrencyAmount(amount, invoice.currency_code);
  const remaining = $derived(
    Math.round((invoice.vendor_invoice_total - invoice.allocated) * 100) / 100,
  );
</script>

<div class="flex w-full flex-col gap-2">
  <DSFieldList>
    <!-- An invoice without details has no total to account against yet. -->
    <DSField
      label="Accounted for"
      value={invoice.vendor_invoice_total > 0
        ? `${money(invoice.allocated)} of ${money(invoice.vendor_invoice_total)}`
        : `${money(invoice.allocated)} on expenses so far`}
      hint={invoice.vendor_invoice_total > 0 && remaining > 0
        ? `${money(remaining)} not yet on an expense`
        : undefined}
    />
    {#if invoice.recorded_by_name}
      <DSField
        label="Invoice details set by"
        value={invoice.recorded_by_name}
        hint={invoice.vendor_invoice_recorded
          ? formatDateTime(invoice.vendor_invoice_recorded)
          : undefined}
      />
    {/if}
  </DSFieldList>

  <DSList items={invoice.parts} inListHeader="Expenses on this invoice">
    {#snippet headline(part: VendorInvoicePart)}
      {#if part.can_open && part.id !== currentExpenseId}
        <a class="link" href={resolve(`/expenses/${part.id}/details`)}
          >{part.po_number ? `PO ${part.po_number}` : "No PO"}</a
        >
      {:else}
        {part.po_number ? `PO ${part.po_number}` : "No PO"}
      {/if}
      · {money(part.total)}
      {#if part.id === currentExpenseId}(this expense){/if}
    {/snippet}
    {#snippet byline(part: VendorInvoicePart)}
      <DSStatusLabel status={partStatus(part)} kind="expense" />
    {/snippet}
    {#snippet line1(part: VendorInvoicePart)}
      {shortDate(part.date, true)} · for {part.for_name || "unknown"}
      {#if part.entered_by_name && part.entered_by_name !== part.for_name}
        · entered by {part.entered_by_name}
      {/if}
    {/snippet}
  </DSList>
</div>
