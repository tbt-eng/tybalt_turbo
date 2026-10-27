<script lang="ts">
  import type { PageData } from "./$types";
  import PurchaseOrdersList from "$lib/components/PurchaseOrdersList.svelte";
  import DsCheck from "$lib/components/DsCheck.svelte";
  import type { VisiblePurchaseOrderResponse } from "$lib/poVisibility";
  import { localIsoDate } from "$lib/localDate";
  let { data }: { data: PageData } = $props();
  const staleDate = localIsoDate(new Date(Date.now() - 30 * 24 * 60 * 60 * 1000));

  // Expense POs are paid out of pocket, so the creator is only reimbursed once
  // they submit an expense against the PO.
  let expenseOnly = $state(false);
  const paymentTypeFilter = $derived((po: VisiblePurchaseOrderResponse) =>
    expenseOnly ? po.payment_type === "Expense" : true,
  );
</script>

<PurchaseOrdersList
  inListHeader={`Stale Purchase Orders (Active before ${staleDate})`}
  {data}
  showRemaining={true}
  filter={paymentTypeFilter}
>
  {#snippet searchBarExtra()}
    <div class="flex items-center gap-2 max-[639px]:w-full max-[639px]:flex-wrap">
      <DsCheck
        bind:value={expenseOnly}
        errors={{}}
        fieldName="expense_only"
        uiName="Expense only"
      />
    </div>
  {/snippet}
</PurchaseOrdersList>
