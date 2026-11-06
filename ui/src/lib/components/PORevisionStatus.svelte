<script lang="ts">
  import DSActionButton from "./DSActionButton.svelte";
  import DSField from "./DSField.svelte";
  import DSFieldList from "./DSFieldList.svelte";
  import DSModal from "./DSModal.svelte";
  import {
    fetchVisiblePO,
    isPORevision,
    poRevisionWaitingOn,
    type VisiblePurchaseOrderResponse,
  } from "$lib/poVisibility";
  import { formatCurrencyAmount, shortDate } from "$lib/utilities";

  // The status of a purchase order under revision. Clicking it explains what is
  // happening and who it is waiting on. It loads the purchase order itself so
  // that any list can show it, whatever fields that list has.
  let { poId, showDetailsLink = true }: { poId: string; showDetailsLink?: boolean } = $props();

  let show = $state(false);
  let po = $state<VisiblePurchaseOrderResponse | null>(null);
  let loadError = $state(false);

  // The page behind the popup can be out of date, so the purchase order may have
  // left revision since it loaded.
  const underRevision = $derived(po !== null && isPORevision(po));
  const waitingOn = $derived(po ? poRevisionWaitingOn(po) : null);

  async function open() {
    show = true;
    po = null;
    loadError = false;
    try {
      po = await fetchVisiblePO(poId);
    } catch (error) {
      console.error(`loading purchase order revision: ${error}`);
      loadError = true;
    }
  }
</script>

<DSActionButton color="yellow" title="See what is happening with this revision" action={open}>
  Under revision
</DSActionButton>

<DSModal
  bind:show
  title={po && underRevision ? `${po.po_number} is under revision` : "Under revision"}
  size="sm"
>
  {#if loadError}
    <p>This purchase order could not be loaded. It may have changed since this page opened.</p>
  {:else if po && !underRevision}
    <p>{po.po_number} is no longer under revision. Its status is {po.status}.</p>
  {:else if po && waitingOn}
    <div class="flex flex-col gap-4">
      {#if po.rejected !== ""}
        <p>
          {po.rejector_name || "An approver"} rejected this revision on {shortDate(po.rejected)}. No
          expenses can be added to it unless it is revised again and approved.
        </p>
      {:else}
        <p>
          {po.uid_name || "Its creator"} revised the total and sent it back for approval. It keeps its
          number, but no expenses can be added to it until it is approved again.
        </p>
      {/if}
      <DSFieldList>
        <DSField label="Waiting on" value={waitingOn.name} />
        <DSField label="Needed" value={waitingOn.reason} />
        <DSField
          label={po.rejected !== "" ? "Rejected total" : "Revised total"}
          value={formatCurrencyAmount(po.approval_total, po.currency_code)}
        />
      </DSFieldList>
    </div>
  {:else}
    <p>Loading…</p>
  {/if}
  {#snippet footer()}
    <DSActionButton variant="secondary" action={() => (show = false)}>Close</DSActionButton>
    {#if showDetailsLink}
      <DSActionButton variant="primary" action={`/pos/${poId}/details`}>
        View purchase order
      </DSActionButton>
    {/if}
  {/snippet}
</DSModal>
