<script lang="ts">
  import { resolve } from "$app/paths";
  import DSPageHeader from "$lib/components/DSPageHeader.svelte";
  import DSActionButton from "$lib/components/DSActionButton.svelte";
  import DSList from "$lib/components/DSList.svelte";
  import DSStatusLabel from "$lib/components/DSStatusLabel.svelte";
  import { globalStore } from "$lib/stores/global";
  import { formatCurrency, shortDate } from "$lib/utilities";
  import type { InvoiceWithClient } from "$lib/invoices";
  import type { PageData } from "./$types";

  let { data }: { data: PageData } = $props();

  const searchText = (invoice: InvoiceWithClient) =>
    `${invoice.number} ${invoice.expand?.client?.name ?? ""} ${invoice.status}`;
</script>

<div class="flex w-full flex-col gap-4 p-2">
  <DSPageHeader title="Invoices" subtitle="Client invoices, and the costs billed back on them">
    {#snippet actions()}
      {#if $globalStore.claims.includes("invoicing")}
        <DSActionButton action="/invoices/add" variant="primary">New invoice</DSActionButton>
      {/if}
    {/snippet}
  </DSPageHeader>

  <DSList
    items={data.invoices}
    search={true}
    {searchText}
    searchPlaceholder="search by number or client..."
    emptyMessage="No invoices yet."
  >
    {#snippet headline(invoice: InvoiceWithClient)}
      <a class="link" href={resolve(`/invoices/${invoice.id}/details`)}>
        {invoice.number || "Draft without a number"}
      </a>
      · {invoice.expand?.client?.name ?? ""}
    {/snippet}
    {#snippet byline(invoice: InvoiceWithClient)}
      <DSStatusLabel status={invoice.status} kind="invoice" />
    {/snippet}
    {#snippet line1(invoice: InvoiceWithClient)}
      {#if invoice.date}{shortDate(invoice.date, true)} ·{/if}
      {invoice.total ? formatCurrency(invoice.total) : "no total yet"}
    {/snippet}
  </DSList>
</div>
