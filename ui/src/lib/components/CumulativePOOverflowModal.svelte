<script lang="ts">
  import DsActionButton from "./DSActionButton.svelte";
  import { fade } from "svelte/transition";
  import { goto } from "$app/navigation";
  import { formatCurrencyAmount } from "$lib/utilities";

  let show = $state(false);
  let overflowData = $state<{
    parent_po: string;
    parent_po_number: string;
    overflow_amount: number;
    po_total: number;
    // The PO's currency; the PO total and the overflow are both in it.
    currency_code: string;
  } | null>(null);

  function closeModal() {
    show = false;
    overflowData = null;
  }

  export function openModal(data: {
    parent_po: string;
    parent_po_number: string;
    overflow_amount: number;
    po_total: number;
    // The PO's currency; the PO total and the overflow are both in it.
    currency_code: string;
  }) {
    show = true;
    overflowData = data;
  }

  function handleCreateChild() {
    if (!overflowData) return;
    goto(`/pos/${overflowData.parent_po}/add-child`);
  }
</script>

{#if show}
  <div
    class="fixed inset-0 z-50 overflow-x-hidden overflow-y-auto"
    transition:fade={{ duration: 200 }}
  >
    <!-- Backdrop/overlay -->
    <div class="fixed inset-0 bg-black/80"></div>
    <!-- Modal content -->
    <div
      class="relative z-50 mx-auto my-20 flex w-11/12 flex-col rounded-lg bg-neutral-800 p-4 text-neutral-300 md:w-3/5"
    >
      <div class="flex items-start justify-between">
        <h1>Expenses Exceed Purchase Order Total</h1>
      </div>
      <div class="my-4 flex flex-col gap-4">
        {#if overflowData}
          <p>
            This expense (together with others that may exist against this PO) exceeds total of
            purchase order {overflowData.parent_po_number} by {formatCurrencyAmount(
              overflowData.overflow_amount,
              overflowData.currency_code,
            )}.
          </p>
          <p>
            Would you like to create a child PO? Set its total to the full expense, because one
            expense can't use two purchase orders.
          </p>
          <div class="flex gap-4">
            <DsActionButton action={handleCreateChild} color="green">
              Create Child PO
            </DsActionButton>
            <DsActionButton action={closeModal} color="red">Cancel</DsActionButton>
          </div>
        {/if}
      </div>
    </div>
  </div>
{/if}
