<script lang="ts">
  import DSCallout from "$lib/components/DSCallout.svelte";
  import DSPopover from "$lib/components/DSPopover.svelte";
  import DsCheck from "$lib/components/DsCheck.svelte";
  import { pb } from "$lib/pocketbase";
  import { ClientResponseError } from "pocketbase";
  import type { UnlinkedSubjob } from "./types";

  // Links a project to the jobs numbered as its subjobs that are not linked to
  // it, as legacy subjobs were imported. The admin confirms the subjobs by
  // number, and the server refuses if the jobs it would adopt have changed.
  let {
    show = $bindable(),
    job,
    onchanged,
  }: {
    show: boolean;
    job: { id: string; number: string; unlinked_subjobs?: UnlinkedSubjob[] | null };
    onchanged: () => void | Promise<void>;
  } = $props();

  let confirmed = $state(false);
  let submitting = $state(false);
  let errors = $state({} as Record<string, { message: string }>);
  let globalError = $state<string | null>(null);

  const adoptable = $derived((job.unlinked_subjobs ?? []).filter((subjob) => !subjob.blocker));
  const blocked = $derived((job.unlinked_subjobs ?? []).filter((subjob) => subjob.blocker));
  const adoptableNumbers = $derived(adoptable.map((subjob) => subjob.number).join(", "));

  function close() {
    if (submitting) return;
    show = false;
  }

  async function submit() {
    if (submitting || !confirmed || adoptable.length === 0) return;
    errors = {};
    globalError = null;
    submitting = true;
    try {
      await pb.send(`/api/jobs/${job.id}/adopt-subjobs`, {
        method: "POST",
        body: { subjobs: adoptable.map((subjob) => subjob.id) },
      });
      show = false;
      await onchanged();
    } catch (error: unknown) {
      const response = error instanceof ClientResponseError ? error.data : {};
      errors = (response?.data as Record<string, { message: string }> | undefined) ?? {};
      globalError = response?.message ?? "The subjobs could not be adopted. Try again.";
      // The unlinked subjobs changed since the page loaded. Reload them to
      // confirm the current list instead.
      if (errors.subjobs) {
        globalError = null;
        confirmed = false;
        await onchanged();
      }
    } finally {
      submitting = false;
    }
  }
</script>

<DSPopover
  bind:show
  title="Adopt Subjobs"
  subtitle={`Link project ${job.number} to the jobs numbered as its subjobs that are not linked to it. Only the link changes.`}
  error={globalError}
  {submitting}
  submitLabel="Adopt Subjobs"
  submitDisabled={!confirmed || adoptable.length === 0}
  onSubmit={submit}
  onCancel={close}
>
  {#if adoptable.length}
    <DsCheck
      bind:value={confirmed}
      {errors}
      fieldName="subjobs"
      uiName={`Make ${adoptableNumbers} ${adoptable.length === 1 ? "a subjob" : "subjobs"} of ${job.number}.`}
      disabled={submitting}
    />
  {:else}
    <DSCallout tone="info">None of these jobs can be adopted yet.</DSCallout>
  {/if}

  {#if blocked.length}
    <DSCallout tone="warning" title="Not adopted">
      <ul class="flex flex-col gap-1">
        {#each blocked as subjob (subjob.id)}
          <li>{subjob.number}: {subjob.blocker}</li>
        {/each}
      </ul>
    </DSCallout>
  {/if}
</DSPopover>
