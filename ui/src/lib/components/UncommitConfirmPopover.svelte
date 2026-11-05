<script lang="ts">
  import { expenseUncommitReasonText } from "$lib/expenseUncommitRequests";
  import type { ExpenseUncommitRequestRow } from "$lib/svelte-types";
  import DSPopover from "./DSPopover.svelte";

  let {
    show = $bindable(),
    recordLabel,
    pendingRequest,
    submitting = false,
    error = null,
    onSubmit,
    onCancel,
  }: {
    show: boolean;
    recordLabel: string;
    // Confirming approves this request and does what it asked for.
    pendingRequest?: ExpenseUncommitRequestRow;
    submitting?: boolean;
    error?: string | null;
    onSubmit: () => void;
    onCancel: () => void;
  } = $props();

  const rejects = $derived(pendingRequest?.outcome === "uncommit_and_reject");
</script>

<DSPopover
  bind:show
  title={rejects ? "Confirm Uncommit and Reject" : "Confirm Uncommit"}
  subtitle={rejects
    ? `As ${pendingRequest?.requester_name} requested, this will uncommit and reject the ${recordLabel}.`
    : `This will return the ${recordLabel} to an approved, uncommitted state.`}
  {error}
  {submitting}
  submitLabel={rejects ? "Uncommit and Reject" : "Confirm Uncommit"}
  {onSubmit}
  {onCancel}
>
  <div class="rounded-sm border border-red-300 bg-red-50 p-3 text-red-900">
    <p class="font-semibold">
      Please confirm with accounting and any other report consumers before continuing.
    </p>
    <p class="mt-2 text-sm">
      Uncommitting this {recordLabel} will change downstream reports, so make sure anyone relying on those
      numbers is aligned first.
    </p>
    <p class="mt-2 text-sm">
      If this {recordLabel} has already written back to a legacy system, that writeback must be cleaned
      up manually after the uncommit.
    </p>
  </div>

  <div class="rounded-sm border border-amber-300 bg-amber-50 p-3 text-amber-900">
    <p class="font-semibold">Review what will happen if you continue:</p>
    <ul class="mt-2 list-disc space-y-1 pl-5 text-sm">
      <li>The current committed state will be cleared.</li>
      {#if rejects && pendingRequest}
        <li>
          The {recordLabel} will be rejected and go back to its owner, with this reason: {expenseUncommitReasonText(
            pendingRequest,
          )}
        </li>
      {:else}
        <li>The {recordLabel} will return to an approved, uncommitted state.</li>
      {/if}
      {#if pendingRequest}
        <li>{pendingRequest.requester_name}'s request will be marked approved.</li>
      {/if}
      <li>No uncommit request is sent until you confirm.</li>
    </ul>
  </div>
</DSPopover>
