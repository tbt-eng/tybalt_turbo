<script lang="ts">
  import { ClientResponseError } from "pocketbase";
  import { pb } from "$lib/pocketbase";
  import { shortDate } from "$lib/utilities";
  import DSCallout from "$lib/components/DSCallout.svelte";
  import DSActionButton from "$lib/components/DSActionButton.svelte";
  import DSPopover from "$lib/components/DSPopover.svelte";
  import DSDateInput from "$lib/components/DSDateInput.svelte";
  import DsTextInput from "$lib/components/DSTextInput.svelte";

  // Moves every active project on an earlier revision of a rate sheet to this
  // revision, from one date. Only admins can do it, so the page shows this only
  // to them. The server checks every rule, and refuses when the number of
  // projects to move has changed since the confirmation showed it.

  type Rollout = {
    effective_date: string;
    move: number;
    closed: number;
    cancelled: number;
    later_period: number;
  };

  let { id, revision, effectiveDate }: { id: string; revision: number; effectiveDate: string } =
    $props();

  const today = new Date().toISOString().slice(0, 10);

  // The counts on the revision's effective date decide whether to offer the
  // move. The form has its own counts, for the date in it.
  let offer = $state<Rollout | null>(null);
  let counts = $state<Rollout | null>(null);
  let moved = $state<number | null>(null);
  let show = $state(false);
  let date = $state("");
  let reason = $state("");
  let errors = $state({} as Record<string, { message: string }>);
  let globalError = $state<string | null>(null);
  let submitting = $state(false);

  function preview(on?: string) {
    return pb.send<Rollout>(`/api/rate_sheets/${id}/rollout`, {
      method: "GET",
      query: on ? { effective_date: on } : {},
    });
  }

  // A sheet that is not the latest revision, or whose effective date is still
  // to come, has nothing to offer.
  async function load() {
    try {
      offer = await preview();
    } catch {
      offer = null;
    }
  }

  $effect(() => {
    if (id) load();
  });

  function showError(error: unknown) {
    const response = error instanceof ClientResponseError ? error.data : {};
    errors = (response?.data as Record<string, { message: string }> | undefined) ?? {};
    globalError = response?.message ?? "The projects could not be moved. Try again.";
    // A refusal because the projects changed carries the new counts.
    if (typeof response?.move === "number") counts = response as Rollout;
  }

  // Count again when the date changes. A late answer for an earlier date is
  // ignored.
  $effect(() => {
    if (!show || !date) return;
    const requested = date;
    errors = {};
    globalError = null;
    preview(requested)
      .then((result) => {
        if (requested === date) counts = result;
      })
      .catch((error: unknown) => {
        if (requested !== date) return;
        counts = null;
        showError(error);
      });
  });

  function open() {
    date = effectiveDate;
    reason = "";
    counts = offer;
    moved = null;
    show = true;
  }

  async function submit() {
    if (submitting || !counts) return;
    errors = {};
    globalError = null;
    submitting = true;
    try {
      const result = await pb.send<Rollout>(`/api/rate_sheets/${id}/rollout`, {
        method: "POST",
        body: { effective_date: date, reason, expected_count: counts.move },
      });
      moved = result.move;
      show = false;
      await load();
    } catch (error: unknown) {
      showError(error);
    } finally {
      submitting = false;
    }
  }

  function projects(n: number) {
    return `${n} project${n === 1 ? "" : "s"}`;
  }

  // The active projects on an earlier revision. One skipped on the default
  // date for a later period can still move from a later date, so it counts.
  const movable = $derived(offer ? offer.move + offer.later_period : 0);

  const skipped = $derived(
    counts
      ? [
          counts.closed ? `${counts.closed} closed` : "",
          counts.cancelled ? `${counts.cancelled} cancelled` : "",
          counts.later_period ? `${counts.later_period} with a period from this date or later` : "",
        ].filter((text) => text !== "")
      : [],
  );
</script>

{#if moved !== null || movable > 0}
  <div class="mb-6 flex flex-col gap-2">
    {#if moved !== null}
      <DSCallout tone="success">Moved {projects(moved)} to rev. {revision}.</DSCallout>
    {/if}
    {#if movable > 0}
      <DSCallout tone="info">
        {movable === 1 ? "1 active project uses" : `${movable} active projects use`} an earlier revision
        of this rate sheet.
        {#snippet actions()}
          <DSActionButton variant="secondary" action={open}>Move projects</DSActionButton>
        {/snippet}
      </DSCallout>
    {/if}
  </div>
{/if}

{#if show}
  <DSPopover
    bind:show
    title="Move Projects to rev. {revision}"
    error={globalError}
    {submitting}
    submitLabel="Move projects"
    submitDisabled={!counts || counts.move === 0 || reason.trim() === ""}
    onSubmit={submit}
  >
    <label class="flex flex-col gap-1">
      <span>Starts on</span>
      <DSDateInput bind:value={date} min={effectiveDate} max={today} disabled={submitting} />
    </label>
    {#if counts}
      <p>{projects(counts.move)} will move.</p>
      {#if skipped.length > 0}
        <p>Skipped: {skipped.join(", ")}.</p>
      {/if}
    {/if}
    <DsTextInput
      bind:value={reason}
      {errors}
      fieldName="reason"
      uiName="Reason"
      placeholder="Why the projects are moving"
      disabled={submitting}
    />
    {#if date}
      <DSCallout tone="warning">
        Time from {shortDate(date, true)} on will be priced by rev. {revision}. Earlier time keeps
        its price.
      </DSCallout>
    {/if}
  </DSPopover>
{/if}
