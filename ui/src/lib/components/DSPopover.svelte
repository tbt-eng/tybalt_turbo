<script lang="ts">
  import type { Snippet } from "svelte";
  import DSActionButton from "./DSActionButton.svelte";
  import DSModal from "./DSModal.svelte";

  // A confirm or small form dialog: DSModal with a Cancel and submit footer.

  let {
    show = $bindable(),
    title,
    subtitle = "",
    error = null,
    submitting = false,
    submitLabel = "Submit",
    cancelLabel = "Cancel",
    danger = false,
    onSubmit,
    onCancel,
    submitDisabled = false,
    children,
  }: {
    show: boolean;
    title: string;
    subtitle?: string;
    error?: string | null;
    submitting?: boolean;
    submitLabel?: string;
    // The button that closes the dialog and does nothing. Name it for what
    // stays, such as "Keep it", when the submit button itself says Cancel.
    cancelLabel?: string;
    // A destructive action shows a red confirm button.
    danger?: boolean;
    onSubmit: () => void;
    onCancel?: () => void;
    // Keeps the submit button disabled until the form can be submitted.
    submitDisabled?: boolean;
    children: Snippet;
  } = $props();

  function handleCancel() {
    if (onCancel) {
      onCancel();
    } else {
      show = false;
    }
  }
</script>

<DSModal bind:show {title} {subtitle} size="sm" dismissible={!submitting} onclose={handleCancel}>
  <div class="flex flex-col gap-4">
    {@render children()}
  </div>

  {#if error}
    <p class="mt-4 text-sm text-red-600">{error}</p>
  {/if}

  {#snippet footer()}
    <DSActionButton variant="secondary" action={handleCancel} disabled={submitting}>
      {cancelLabel}
    </DSActionButton>
    <DSActionButton
      variant={danger ? "danger" : "primary"}
      action={onSubmit}
      disabled={submitting || submitDisabled}
    >
      {submitting ? "Saving..." : submitLabel}
    </DSActionButton>
  {/snippet}
</DSModal>
