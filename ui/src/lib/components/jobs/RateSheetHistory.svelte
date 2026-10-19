<script lang="ts">
  import { resolve } from "$app/paths";
  import { ClientResponseError } from "pocketbase";
  import { pb } from "$lib/pocketbase";
  import { rateSheets, type RateSheetResponse } from "$lib/stores/rateSheets";
  import { shortDate } from "$lib/utilities";
  import DSSection from "$lib/components/DSSection.svelte";
  import DSList from "$lib/components/DSList.svelte";
  import DSActionButton from "$lib/components/DSActionButton.svelte";
  import DSPopover from "$lib/components/DSPopover.svelte";
  import DSCallout from "$lib/components/DSCallout.svelte";
  import DsSelector from "$lib/components/DSSelector.svelte";
  import DSDateInput from "$lib/components/DSDateInput.svelte";
  import DsTextInput from "$lib/components/DSTextInput.svelte";

  // A project's rate sheet history: which rate sheet prices its time from which
  // date. Each entry is priced by the period its date falls in, so a change
  // here reprices that period's past time. The server checks every rule; the
  // form only collects the change.

  type Period = {
    id: string;
    rate_sheet: {
      id: string;
      name: string;
      revision: number;
      effective_date: string;
      active: boolean;
    };
    starts: string;
    through: string;
    reason: string;
    created_by: string;
    updated_by: string;
  };
  type History = { can_revise: boolean; can_set_first: boolean; periods: Period[] };

  let { jobId, onchanged }: { jobId: string; onchanged?: () => void | Promise<void> } = $props();

  rateSheets.init();

  let history = $state<History | null>(null);
  let loadError = $state(false);

  // The open form: the first sheet of a project with none, a new dated period,
  // a change to a period, or a delete.
  type Mode = "first" | "add" | "edit" | "delete";
  let mode = $state<Mode | null>(null);
  let editing = $state<Period | null>(null);
  let sheet = $state("");
  let date = $state("");
  let reason = $state("");
  let errors = $state({} as Record<string, { message: string }>);
  let globalError = $state<string | null>(null);
  let submitting = $state(false);
  let show = $state(false);

  const today = new Date().toISOString().slice(0, 10);

  async function load() {
    loadError = false;
    try {
      history = await pb.send<History>(`/api/jobs/${jobId}/rate_sheets`, { method: "GET" });
    } catch {
      loadError = true;
    }
  }

  $effect(() => {
    if (jobId) load();
  });

  // Active sheets can be chosen. A period keeps its sheet after the sheet is
  // deactivated, so the current choice is offered too.
  const sheetOptions = $derived(
    $rateSheets.items.filter(
      (s: RateSheetResponse) => s.active || (editing !== null && s.id === editing.rate_sheet.id),
    ),
  );

  const isOpening = $derived(mode === "first" || (mode === "edit" && editing?.starts === ""));
  const firstDated = $derived(history?.periods.find((p) => p.starts !== "")?.starts ?? "");

  function open(next: Mode, period: Period | null = null) {
    mode = next;
    editing = period;
    sheet = next === "add" || next === "first" ? "" : (period?.rate_sheet.id ?? "");
    date = next === "add" ? today : (period?.starts ?? "");
    reason = "";
    errors = {};
    globalError = null;
    show = true;
  }

  function close() {
    show = false;
    mode = null;
    editing = null;
  }

  function periodDates(period: Period) {
    const from = period.starts ? `From ${shortDate(period.starts, true)}` : "From the start";
    return period.through ? `${from} to ${shortDate(period.through, true)}` : `${from} onward`;
  }

  async function submit() {
    if (submitting || mode === null) return;
    errors = {};
    globalError = null;
    submitting = true;
    try {
      let result: { periods: Period[] };
      if (mode === "delete") {
        result = await pb.send(`/api/jobs/${jobId}/rate_sheets/${editing?.id}`, {
          method: "DELETE",
        });
      } else if (mode === "edit") {
        const body: Record<string, string> = { rate_sheet: sheet, reason };
        if (editing?.starts) body.effective_date = date;
        result = await pb.send(`/api/jobs/${jobId}/rate_sheets/${editing?.id}`, {
          method: "PATCH",
          body,
        });
      } else {
        result = await pb.send(`/api/jobs/${jobId}/rate_sheets`, {
          method: "POST",
          body: { rate_sheet: sheet, effective_date: mode === "add" ? date : "", reason },
        });
      }
      if (history) history = { ...history, periods: result.periods };
      close();
      await load();
      await onchanged?.();
    } catch (error: unknown) {
      const response = error instanceof ClientResponseError ? error.data : {};
      errors = (response?.data as Record<string, { message: string }> | undefined) ?? {};
      globalError = response?.message ?? "The rate sheet history could not be changed. Try again.";
    } finally {
      submitting = false;
    }
  }

  const titles: Record<Mode, string> = {
    first: "Set Rate Sheet",
    add: "Start a New Rate Sheet",
    edit: "Change Rate Sheet Period",
    delete: "Delete Rate Sheet Period",
  };
</script>

{#if history && (history.periods.length > 0 || history.can_set_first)}
  <DSSection title="Rate sheets">
    {#snippet actions()}
      {#if history?.can_set_first}
        <DSActionButton variant="secondary" action={() => open("first")}
          >Set rate sheet</DSActionButton
        >
      {:else if history?.can_revise}
        <DSActionButton variant="secondary" action={() => open("add")}>New period</DSActionButton>
      {/if}
    {/snippet}
    {#if history.periods.length}
      <DSList items={history.periods}>
        {#snippet headline(period: Period)}
          <a class="link" href={resolve("/rate-sheets/[id]/details", { id: period.rate_sheet.id })}
            >{period.rate_sheet.name} (rev. {period.rate_sheet.revision})</a
          >
        {/snippet}
        {#snippet line1(period: Period)}
          {periodDates(period)}
        {/snippet}
        {#snippet line2(period: Period)}
          {#if period.reason}{period.reason}{/if}
          {#if period.updated_by || period.created_by}
            {period.reason ? " · " : ""}Set by {period.updated_by || period.created_by}
          {/if}
        {/snippet}
        {#snippet actions(period: Period)}
          {#if history?.can_revise}
            <DSActionButton variant="chip" action={() => open("edit", period)}
              >Change</DSActionButton
            >
            {#if period.starts}
              <DSActionButton variant="chip" action={() => open("delete", period)}
                >Delete</DSActionButton
              >
            {/if}
          {/if}
        {/snippet}
      </DSList>
    {:else}
      <DSCallout tone="warning">
        This project has no rate sheet. Its time is priced at employee default charge-out rates.
      </DSCallout>
    {/if}
  </DSSection>
{:else if loadError}
  <DSCallout tone="danger" title="Could not load the rate sheet history">
    {#snippet actions()}
      <DSActionButton variant="secondary" action={load}>Try again</DSActionButton>
    {/snippet}
  </DSCallout>
{/if}

{#if mode !== null}
  <DSPopover
    bind:show
    title={titles[mode]}
    error={globalError}
    {submitting}
    submitLabel={mode === "delete" ? "Delete" : "Save"}
    danger={mode === "delete"}
    onSubmit={submit}
    onCancel={close}
  >
    {#if mode === "delete" && editing}
      <p>
        Time from {shortDate(editing.starts, true)} will be priced by the period before it. Past values
        change.
      </p>
    {:else}
      <DsSelector
        bind:value={sheet}
        items={sheetOptions}
        {errors}
        fieldName="rate_sheet"
        uiName="Rate sheet"
        disabled={submitting}
      >
        {#snippet optionTemplate(item)}{item.name} (rev. {item.revision}){/snippet}
      </DsSelector>
      {#if !isOpening}
        <label class="flex flex-col gap-1">
          <span>Starts on</span>
          <DSDateInput bind:value={date} max={today} disabled={submitting} />
        </label>
      {/if}
      {#if mode !== "first"}
        <DsTextInput
          bind:value={reason}
          {errors}
          fieldName="reason"
          uiName="Reason"
          placeholder="Why the rate sheet is changing"
          disabled={submitting}
        />
      {/if}
      {#if mode === "first"}
        <DSCallout tone="info">
          All of this project's time, past and future, will be priced by this rate sheet.
        </DSCallout>
      {:else if isOpening}
        <DSCallout tone="warning">
          All time {firstDated ? `before ${shortDate(firstDated, true)}` : "on this project"} will be
          repriced by the chosen sheet.
        </DSCallout>
      {:else if date && date < today}
        <DSCallout tone="warning">
          Time from {shortDate(date, true)} onward will be repriced by the chosen sheet.
        </DSCallout>
      {/if}
    {/if}
  </DSPopover>
{/if}
