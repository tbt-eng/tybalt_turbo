<script lang="ts">
  // The second step before a purchase order is closed or cancelled. Both are
  // final, so the dialog says what happens and that it can't be undone, and
  // nothing is sent until the person confirms. The owner of the purchase order
  // and payables staff use the same dialog.
  import { pb } from "$lib/pocketbase";
  import { globalStore } from "$lib/stores/global";
  import DSCallout from "$lib/components/DSCallout.svelte";
  import DSPopover from "$lib/components/DSPopover.svelte";

  let {
    show = $bindable(),
    action,
    po,
    ondone,
  }: {
    show: boolean;
    action: "close" | "cancel";
    po: { id: string; po_number: string };
    // Called after the server closed or cancelled the purchase order.
    ondone?: () => void | Promise<void>;
  } = $props();

  let submitting = $state(false);
  let error = $state<string | null>(null);
  let wasOpen = false;

  $effect(() => {
    if (show && !wasOpen) error = null;
    wasOpen = show;
  });

  const cancelling = $derived(action === "cancel");
  const name = $derived(po.po_number ? `PO ${po.po_number}` : "this purchase order");

  async function confirm() {
    submitting = true;
    error = null;
    try {
      await pb.send(`/api/purchase_orders/${po.id}/${action}`, { method: "POST" });
      await globalStore.refreshAttentionCounts();
      show = false;
      await ondone?.();
    } catch (e) {
      // The server says why, such as expenses that still await approval.
      const failure = e as { response?: { message?: string }; message?: string };
      error =
        failure?.response?.message ??
        failure?.message ??
        (cancelling ? "Cancel failed" : "Close failed");
    } finally {
      submitting = false;
    }
  }
</script>

<DSPopover
  bind:show
  title={cancelling ? `Cancel ${name}?` : `Close ${name}?`}
  {error}
  {submitting}
  danger
  submitLabel={cancelling ? "Cancel Purchase Order" : "Close Purchase Order"}
  cancelLabel={cancelling ? "Keep it" : "Keep it open"}
  onSubmit={confirm}
>
  {#if cancelling}
    <p>
      Turbo cancels this purchase order. No expense can use it, and its number can't be used again.
    </p>
  {:else}
    <p>Turbo closes this purchase order. No more expenses can use it.</p>
  {/if}
  <p>Turbo records your name and the date.</p>
  <DSCallout tone="danger" title="This is permanent">
    You can't undo this. A {cancelling ? "cancelled" : "closed"} purchase order can't be made active again.
  </DSCallout>
</DSPopover>
