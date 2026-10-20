<script lang="ts">
  import InvoiceEditor from "$lib/components/invoices/InvoiceEditor.svelte";
  import DSCallout from "$lib/components/DSCallout.svelte";
  import type { PageData } from "./$types";

  let { data }: { data: PageData } = $props();
</script>

{#if data.invoice.status === "Draft"}
  <InvoiceEditor
    id={data.invoice.id}
    item={{
      client: data.invoice.client,
      number: data.invoice.number,
      date: data.invoice.date,
      total: data.invoice.total || "",
    }}
  />
{:else}
  <div class="p-2">
    <DSCallout tone="warning">An issued invoice cannot be edited.</DSCallout>
  </div>
{/if}
