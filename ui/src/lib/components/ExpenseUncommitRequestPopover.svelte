<script lang="ts">
  import { untrack } from "svelte";
  import { pb } from "$lib/pocketbase";
  import { getApiErrorMessage } from "$lib/errors";
  import { expenseUncommitOutcomeOptions } from "$lib/expenseUncommitRequests";
  import type { ExpenseUncommitOutcome, ExpenseUncommitReason } from "$lib/svelte-types";
  import DSPopover from "./DSPopover.svelte";
  import DSSelector from "./DSSelector.svelte";
  import DSTextInput from "./DSTextInput.svelte";

  let {
    show = $bindable(),
    expenseId,
    onRequested,
  }: {
    show: boolean;
    expenseId: string;
    onRequested: () => void | Promise<void>;
  } = $props();

  let reasons = $state<ExpenseUncommitReason[]>([]);
  let reasonId = $state("");
  let outcome = $state<ExpenseUncommitOutcome>("uncommit_and_reject");
  let note = $state("");
  let error = $state<string | null>(null);
  let submitting = $state(false);

  const reason = $derived(reasons.find((candidate) => candidate.id === reasonId));
  const canSubmit = $derived(
    reason !== undefined && (!reason.note_required || note.trim().length > 0),
  );

  async function loadReasons() {
    try {
      reasons = await pb.send("/api/expense_uncommit_requests/reasons", { method: "GET" });
    } catch (err) {
      error = getApiErrorMessage(err, "Failed to load uncommit reasons");
    }
  }

  $effect(() => {
    if (show) untrack(() => void loadReasons());
  });

  // Each reason suggests an outcome, which the requester can still change.
  $effect(() => {
    const suggested = reason?.default_outcome;
    if (suggested) untrack(() => (outcome = suggested));
  });

  function close() {
    show = false;
    reasonId = "";
    note = "";
    error = null;
  }

  async function submit() {
    submitting = true;
    error = null;
    try {
      await pb.send(`/api/expenses/${expenseId}/uncommit_requests`, {
        method: "POST",
        body: { reason: reasonId, outcome, note: note.trim() },
      });
      close();
      await onRequested();
    } catch (err) {
      error = getApiErrorMessage(err, "Failed to request the uncommit");
    } finally {
      submitting = false;
    }
  }
</script>

<DSPopover
  bind:show
  title="Request Uncommit"
  subtitle="Dean will see this in the uncommit request queue."
  {error}
  {submitting}
  submitLabel="Request"
  submitDisabled={!canSubmit}
  onSubmit={submit}
  onCancel={close}
>
  <DSSelector bind:value={reasonId} items={reasons} errors={{}} fieldName="reason" uiName="Reason">
    {#snippet optionTemplate(item)}
      {item.name}
    {/snippet}
  </DSSelector>
  <DSTextInput
    bind:value={note}
    errors={{}}
    fieldName="note"
    uiName={reason?.note_required ? "Note (required)" : "Note"}
    placeholder={reason?.note_prompt ?? ""}
  />
  <DSSelector
    bind:value={outcome}
    items={expenseUncommitOutcomeOptions}
    errors={{}}
    fieldName="outcome"
    uiName="Result"
    helperText={outcome === "uncommit_and_reject"
      ? "The expense goes back to its owner to correct."
      : "The expense goes back to the commit queue as it is."}
  >
    {#snippet optionTemplate(item)}
      {item.name}
    {/snippet}
  </DSSelector>
</DSPopover>
