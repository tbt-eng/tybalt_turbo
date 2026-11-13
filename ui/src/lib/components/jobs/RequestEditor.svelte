<script lang="ts">
  import { pb } from "$lib/pocketbase";
  import { goto, invalidateAll, beforeNavigate, replaceState } from "$app/navigation";
  import { fullName } from "$lib/jobRequests";
  import { workflowError, type WorkflowRecord } from "$lib/clientEditors";
  import { untrack } from "svelte";
  import MiniSearch from "minisearch";
  import DSLocationPicker from "$lib/components/DSLocationPicker.svelte";
  import DSActionButton from "$lib/components/DSActionButton.svelte";
  import DsAutoComplete from "$lib/components/DSAutoComplete.svelte";
  import DSCallout from "$lib/components/DSCallout.svelte";
  import DSHelp from "$lib/components/DSHelp.svelte";
  import DsDateInput from "$lib/components/DSDateInput.svelte";
  import DsSelector from "$lib/components/DSSelector.svelte";
  import DsTextInput from "$lib/components/DSTextInput.svelte";
  import DSToggle from "$lib/components/DSToggle.svelte";
  import { globalStore } from "$lib/stores/global";
  import { clients } from "$lib/stores/clients";
  import { divisions } from "$lib/stores/divisions";
  import { jobs } from "$lib/stores/jobs";
  import { rateSheets } from "$lib/stores/rateSheets";
  import { formatJobLabel } from "$lib/utilities";
  import InvoicingProfileEmptyState from "$lib/components/clients/InvoicingProfileEmptyState.svelte";
  import { canCreateInvoicingProfile, clientLabel } from "$lib/clientWorkspace";
  import {
    type ClientWorkflowTarget,
    clearClientWorkflowDraft,
    keepClientWorkflowDraft,
    readClientWorkflowDraft,
  } from "$lib/clientWorkflowNavigation";
  let {
    request,
    allocations: initialAllocations = [],
    options,
    dirty = $bindable(false),
    onsaved,
  }: {
    request: WorkflowRecord;
    allocations?: WorkflowRecord[];
    options: WorkflowRecord;
    dirty?: boolean;
    onsaved?: () => Promise<void>;
  } = $props();
  clients.init();
  jobs.init();
  divisions.init();
  rateSheets.init();
  const restored = untrack(() =>
    readClientWorkflowDraft<{
      item: WorkflowRecord;
      allocations: WorkflowRecord[];
      fee: string;
      saved: string;
    }>(`request:${request.id || "new"}`, pb.authStore.record?.id || ""),
  );
  let item = $state<WorkflowRecord>(
    untrack(() => ({
      project_award_date: "",
      project_completion_date: "",
      ...(restored?.item || request),
    })),
  );
  let allocations = $state<WorkflowRecord[]>(
    untrack(() =>
      (restored?.allocations || initialAllocations).map((row) => ({
        division: row.division,
        hours: row.hours,
      })),
    ),
  );
  let fee = $state(
    untrack(
      () =>
        restored?.fee ??
        (request.time_and_materials_set ? (request.time_and_materials ? "tm" : "fixed") : ""),
    ),
  );
  let saved = $state(untrack(() => restored?.saved || JSON.stringify({ item, allocations, fee })));
  let openingClient = false;
  let saving = $state(false),
    error = $state("");
  let approvers = $state<WorkflowRecord[]>([]),
    approversLoaded = $state(false),
    approverError = $state("");
  const contacts = $derived(
    (options.contacts || []).filter((record: WorkflowRecord) => record.client === item.client),
  );
  const profiles = $derived(
    (options.invoicing_profiles || []).filter(
      (record: WorkflowRecord) => record.client === item.client,
    ),
  );
  const proposals = $derived(
    (options.proposals || []).filter(
      (record: WorkflowRecord) =>
        record.client === item.client &&
        ["In Progress", "Submitted", "Awarded"].includes(record.status),
    ),
  );
  const proposalIds = $derived(new Set(proposals.map((record: WorkflowRecord) => record.id)));
  // From the request options, not the profiles store: a requester without the
  // job claim can only list their own profile and their direct reports.
  const managerIndex = $derived.by(() => {
    const index = new MiniSearch<WorkflowRecord>({
      fields: ["surname", "given_name", "name"],
      storeFields: ["surname", "given_name", "name"],
      searchOptions: { combineWith: "AND" },
    });
    index.addAll(options.users || []);
    return index;
  });
  const activeDivisionIds = $derived(
    new Set(
      (options.divisions || [])
        .filter((record: WorkflowRecord) => record.active)
        .map((record: WorkflowRecord) => record.id),
    ),
  );
  const approverItems = $derived([
    { id: "", name: "Select approver" },
    ...(item.approver && !approvers.some((row) => row.id === item.approver)
      ? [{ id: item.approver, name: "Current approver — check authority" }]
      : []),
    ...approvers.map((row) => ({
      id: row.id,
      name: `${row.name} — up to $${row.max_value.toLocaleString("en-CA")}`,
    })),
  ]);
  const noFieldErrors: Record<string, { message: string }> = {};
  const feeOptions = [
    { id: "fixed", label: "Fixed fee" },
    { id: "tm", label: "Time and materials" },
  ];
  const answerOptions = [
    { id: "Yes", label: "Yes" },
    { id: "No", label: "No" },
  ];
  const statusOptions = [
    { id: "Active", label: "Active" },
    { id: "Closed", label: "Closed" },
  ];
  function managerLabel(row: WorkflowRecord) {
    if (row.surname) return `${row.surname}, ${row.given_name}`;
    return row.name || "Inactive person";
  }
  function invoicingProfileLabel(row: WorkflowRecord) {
    const name =
      row.name ||
      fullName(
        (options.contacts || []).find((contact: WorkflowRecord) => contact.id === row.contact),
      ) ||
      row.id;
    return row.invoicing_instructions ? `${name} — ${row.invoicing_instructions}` : name;
  }
  // The approver list depends on the project value, so say why it is empty.
  const approverHint = $derived.by(() => {
    const value = item.project_value;
    if (!(Number.isInteger(value) && value > 0))
      return "Enter the project value in whole dollars to list approvers.";
    if (approversLoaded && approvers.length === 0)
      return "No approver has a limit that covers this project value.";
    return "";
  });
  $effect(() => {
    dirty = JSON.stringify({ item, allocations, fee }) !== saved;
  });
  $effect(() => {
    const value = item.project_value;
    let current = true;
    approvers = [];
    approversLoaded = false;
    if (Number.isInteger(value) && value > 0) {
      pb.send(`/api/project_authorization/approvers?value=${value}`, {
        method: "GET",
        requestKey: null,
      })
        .then((response) => {
          if (current) {
            approvers = response.items;
            approversLoaded = true;
            approverError = "";
          }
        })
        .catch((e) => {
          if (current) approverError = workflowError(e);
        });
    }
    return () => {
      current = false;
    };
  });
  beforeNavigate(({ cancel }) => {
    if (openingClient) return;
    if (dirty && !saving && !window.confirm("Discard unsaved request changes?")) cancel();
    else clearClientWorkflowDraft();
  });
  async function manageClient(target: ClientWorkflowTarget = "workspace", event?: MouseEvent) {
    if (event && (event.metaKey || event.ctrlKey || event.shiftKey || event.altKey)) return;
    event?.preventDefault();
    try {
      const href = keepClientWorkflowDraft(
        `request:${request.id || "new"}`,
        pb.authStore.record?.id || "",
        $state.snapshot({ item, allocations, fee, saved }),
        item.client,
        item.contact || "",
        target,
        (url) => replaceState(url, {}),
      );
      openingClient = true;
      await goto(href);
    } catch (e) {
      error = workflowError(e);
    } finally {
      openingClient = false;
    }
  }
  function setClient(client: string) {
    if (client === item.client) return;
    item.client = client;
    item.contact = "";
    item.invoicing_information = "";
    item.proposal = "";
  }
  export async function save(): Promise<boolean> {
    if (saving) return false;
    saving = true;
    error = "";
    const fields = [
      "description",
      "notes",
      "client",
      "contact",
      "manager",
      "alternate_manager",
      "branch",
      "location",
      "job_owner",
      "parent",
      "proposal",
      "project_status",
      "project_value",
      "project_award_date",
      "project_completion_date",
      "client_po",
      "client_reference_number",
      "fn_agreement_answer",
      "invoicing_information",
      "rate_sheet",
      "approver",
    ];
    const payload = Object.fromEntries(
      fields.filter((field) => item[field] !== undefined).map((field) => [field, item[field]]),
    );
    if (fee) {
      payload.time_and_materials = fee === "tm";
    }
    try {
      const result = await pb.send(item.id ? `/api/job_requests/${item.id}` : "/api/job_requests", {
        method: item.id ? "PATCH" : "POST",
        body: { version: item.version, request: payload, allocations },
      });
      saved = JSON.stringify({ item, allocations, fee });
      dirty = false;
      clearClientWorkflowDraft();
      if (!item.id) await goto(`/job-requests/${result.request?.id ?? result.id}`);
      else if (onsaved) await onsaved();
      else await invalidateAll();
      return true;
    } catch (e) {
      error = workflowError(e);
      return false;
    } finally {
      saving = false;
    }
  }
</script>

<form
  class="flex flex-col gap-4"
  onsubmit={(event) => {
    event.preventDefault();
    save();
  }}
>
  {#if error}<div role="alert"><DSCallout tone="danger">{error}</DSCallout></div>{/if}
  <fieldset disabled={saving} class="flex flex-col gap-2">
    <span class="flex w-full items-center gap-2">
      <label for="project_award_date">Project Award Date</label>
      <DsDateInput
        class="flex-1 md:flex-none"
        name="project_award_date"
        bind:value={item.project_award_date}
      />
    </span>
    <span class="flex w-full items-center gap-2">
      <label for="project_completion_date">Project Completion Date</label>
      <DsDateInput
        class="flex-1 md:flex-none"
        name="project_completion_date"
        bind:value={item.project_completion_date}
      />
    </span>
    <DSLocationPicker bind:value={item.location} errors={{}} fieldName="location" />
    <DsTextInput
      bind:value={item.description}
      errors={noFieldErrors}
      fieldName="description"
      uiName="Description"
    />
    {#if $clients.index !== null}
      <DsAutoComplete
        bind:value={() => item.client, setClient}
        index={$clients.index}
        errors={noFieldErrors}
        fieldName="client"
        uiName="Client"
        disabled={!!item.parent}
      >
        {#snippet resultTemplate(row)}{clientLabel({ name: row.name, alias: row.alias })}{/snippet}
      </DsAutoComplete>
    {/if}
    <DsSelector
      bind:value={item.contact}
      items={[{ id: "" }, ...contacts]}
      errors={noFieldErrors}
      fieldName="contact"
      uiName="Client Contact"
      disabled={!item.client}
    >
      {#snippet optionTemplate(row: WorkflowRecord)}{row.id
          ? `${fullName(row)}${row.email ? ` — ${row.email}` : ""}`
          : item.client
            ? "Select a contact"
            : "Select a client first"}{/snippet}
    </DsSelector>
    {#if item.client && !profiles.length}
      <div class="flex flex-col gap-1">
        <h3 class="text-sm font-semibold">Invoicing profile</h3>
        <InvoicingProfileEmptyState
          canCreate={canCreateInvoicingProfile($globalStore.claims)}
          canManageContacts={$globalStore.claims.includes("job")}
          hasContacts={contacts.length > 0}
          oncreate={() => manageClient("invoicing")}
          onaddcontact={() => manageClient("contact")}
        />
      </div>
    {:else}
      <DsSelector
        bind:value={item.invoicing_information}
        items={[{ id: "" }, ...profiles]}
        errors={noFieldErrors}
        fieldName="invoicing_information"
        uiName="Invoicing profile"
        disabled={!item.client}
      >
        {#snippet optionTemplate(row: WorkflowRecord)}{row.id
            ? invoicingProfileLabel(row)
            : item.client
              ? "Select an invoicing profile"
              : "Select a client first"}{/snippet}
      </DsSelector>
      {#if item.client}<a
          class="link"
          href={`/clients/${item.client}/details?project_contact=${encodeURIComponent(item.contact || "")}`}
          onclick={(event) => manageClient("workspace", event)}
          >Manage client contacts and invoicing</a
        >{/if}
    {/if}
    <DsAutoComplete
      bind:value={item.manager}
      index={managerIndex}
      errors={noFieldErrors}
      fieldName="manager"
      uiName="Manager"
    >
      {#snippet resultTemplate(row)}{managerLabel(row)}{/snippet}
    </DsAutoComplete>
    <DsAutoComplete
      bind:value={item.alternate_manager}
      index={managerIndex}
      errors={noFieldErrors}
      fieldName="alternate_manager"
      uiName="Alternate Manager"
    >
      {#snippet resultTemplate(row)}{managerLabel(row)}{/snippet}
    </DsAutoComplete>
    <div class="flex w-full items-center gap-2">
      <span id="request-fn-agreement">First Nations agreement</span>
      <DSToggle
        bind:value={item.fn_agreement_answer}
        options={answerOptions}
        ariaLabelledBy="request-fn-agreement"
      />
    </div>
    <div class="flex w-full items-center gap-2">
      <span id="request-status">Status</span>
      <DSToggle
        bind:value={item.project_status}
        options={statusOptions}
        ariaLabelledBy="request-status"
      />
    </div>
    <DsTextInput
      bind:value={item.project_value}
      type="number"
      min={1}
      step={1}
      errors={noFieldErrors}
      fieldName="project_value"
      uiName="Project value (CAD, before tax)"
    />
    <div class="flex w-full items-center gap-2">
      <span id="request-fee-basis">Fee basis</span>
      <DSToggle bind:value={fee} options={feeOptions} ariaLabelledBy="request-fee-basis" />
    </div>
    <DsTextInput
      bind:value={item.client_po}
      errors={noFieldErrors}
      fieldName="client_po"
      uiName="Client PO"
    />
    <DsTextInput
      bind:value={item.client_reference_number}
      errors={noFieldErrors}
      fieldName="client_reference_number"
      uiName="Client Reference Number"
    />
    {#if $rateSheets.index !== null}
      <DsAutoComplete
        bind:value={item.rate_sheet}
        index={$rateSheets.index}
        filter={(sheet) => sheet.active}
        errors={noFieldErrors}
        fieldName="rate_sheet"
        uiName="Rate Sheet"
      >
        {#snippet resultTemplate(row)}{row.name} (rev. {row.revision}) - {row.effective_date}{/snippet}
        {#snippet help()}
          <DSHelp doc="rate-sheets#choose-a-rate-sheet-for-a-project" iconOnly />
        {/snippet}
      </DsAutoComplete>
    {/if}
    {#if $jobs.index !== null}
      <DsAutoComplete
        bind:value={item.proposal}
        index={$jobs.index}
        filter={(job) => proposalIds.has(job.id)}
        errors={noFieldErrors}
        fieldName="proposal"
        uiName="Proposal"
        disabled={!item.client}
      >
        {#snippet resultTemplate(row)}{formatJobLabel(row)}{/snippet}
      </DsAutoComplete>
    {/if}
  </fieldset>
  <fieldset disabled={saving} class="flex flex-col gap-2">
    <legend class="font-semibold">Divisions</legend>
    {#if $divisions.index !== null}
      {#each allocations as allocation, index}
        <div class="flex flex-wrap items-center gap-2">
          <div class="min-w-[280px] flex-1">
            <DsAutoComplete
              bind:value={allocation.division}
              index={$divisions.index}
              filter={(division) => activeDivisionIds.has(division.id)}
              errors={noFieldErrors}
              fieldName={`allocation_division_${index}`}
              uiName="Division"
            >
              {#snippet resultTemplate(row)}{row.code} - {row.name}{/snippet}
            </DsAutoComplete>
          </div>
          <div class="w-56">
            <DsTextInput
              bind:value={allocation.hours}
              type="number"
              step={1}
              errors={noFieldErrors}
              fieldName={`allocation_hours_${index}`}
              uiName="Hours"
            />
          </div>
          <DSActionButton action={() => (allocations = allocations.filter((_, i) => i !== index))}
            >Remove</DSActionButton
          >
        </div>
      {/each}
    {/if}
    <div>
      <DSActionButton
        action={() => (allocations = [...allocations, { division: "", hours: 0 }])}
        icon="feather:plus-circle"
        color="green">Add division</DSActionButton
      >
    </div>
  </fieldset>
  <fieldset disabled={saving} class="flex flex-col gap-2">
    {#if $clients.index !== null}
      <DsAutoComplete
        bind:value={item.job_owner}
        index={$clients.index}
        errors={noFieldErrors}
        fieldName="job_owner"
        uiName="Job Owner"
      >
        {#snippet resultTemplate(row)}{clientLabel({ name: row.name, alias: row.alias })}{/snippet}
      </DsAutoComplete>
    {/if}
    <DsSelector
      bind:value={item.branch}
      items={[{ id: "", name: "Select a branch" }, ...(options.branches || [])]}
      errors={noFieldErrors}
      fieldName="branch"
      uiName="Branch"
    >
      {#snippet optionTemplate(row: WorkflowRecord)}{row.code ?? row.name}{/snippet}
    </DsSelector>
    <label class="flex flex-col gap-1">
      Notes (optional)
      <textarea
        rows="3"
        maxlength="1000"
        class="rounded-sm border border-neutral-300 p-2"
        bind:value={item.notes}></textarea>
      <span class="text-xs text-neutral-600"
        >For example, a property's legal description. You can include the notes when printing the
        Project Authorization. They become the job's first note when the job is created. 10 to 1000
        characters.</span
      >
    </label>
    <DsSelector
      bind:value={item.approver}
      items={approverItems}
      errors={noFieldErrors}
      fieldName="approver"
      uiName="Selected approver"
      helperText={approverHint}
    >
      {#snippet optionTemplate(row)}{row.name}{/snippet}
    </DsSelector>
    {#if approverError}<div role="alert">
        <DSCallout tone="danger">{approverError}</DSCallout>
      </div>{/if}
  </fieldset>
  <p class="text-sm text-neutral-600">
    Drafts can be incomplete. Submission checks all required data, the PDF, allocations, and
    approval authority.
  </p>
  <div>
    <DSActionButton variant="primary" type="submit" disabled={saving}
      >{saving ? "Saving…" : "Save request"}</DSActionButton
    >
  </div>
</form>
