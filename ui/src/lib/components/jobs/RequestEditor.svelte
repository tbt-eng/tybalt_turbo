<script lang="ts">
  import DSRequiredLegend from "$lib/components/DSRequiredLegend.svelte";
  import DSRequiredMark from "$lib/components/DSRequiredMark.svelte";
  import { pb } from "$lib/pocketbase";
  import { goto, invalidateAll, beforeNavigate, replaceState } from "$app/navigation";
  import { resolve } from "$app/paths";
  import { copiedAllocations, fullName, hasAllocations } from "$lib/jobRequests";
  import { workflowError, type WorkflowRecord } from "$lib/clientEditors";
  import { untrack } from "svelte";
  import MiniSearch from "minisearch";
  import DSLocationPicker from "$lib/components/DSLocationPicker.svelte";
  import DSActionButton from "$lib/components/DSActionButton.svelte";
  import DsAutoComplete from "$lib/components/DSAutoComplete.svelte";
  import DSCallout from "$lib/components/DSCallout.svelte";
  import DSHelp from "$lib/components/DSHelp.svelte";
  import DSPicker from "$lib/components/DSPicker.svelte";
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
  import { rateSheetCodeRevision, rateSheetSearchPlaceholder } from "$lib/rateSheetLabel";
  import ClientCreatePopover from "$lib/components/clients/ClientCreatePopover.svelte";
  import InvoicingProfileEmptyState from "$lib/components/clients/InvoicingProfileEmptyState.svelte";
  import ProposedRatesEditor from "$lib/components/rateSheets/ProposedRatesEditor.svelte";
  import RateSheetOption from "$lib/components/rateSheets/RateSheetOption.svelte";
  import RateSheetFinder from "$lib/components/rateSheets/RateSheetFinder.svelte";
  import {
    type FoundRate,
    carriedFindRates,
    changedCarriedRates,
    withCarriedRates,
  } from "$lib/findRates";
  import {
    type ProposedRate,
    type ProposedRates,
    emptyProposedRates,
    incompleteRateErrors,
    proposedRatesPayload,
    proposedRatesFieldErrors,
    withoutShownErrors,
  } from "$lib/proposedRates";
  import RateSheetMatchCallout from "$lib/components/rateSheets/RateSheetMatchCallout.svelte";
  import { duplicateRateSheetMatches, type RateSheetMatch } from "$lib/rateSheetMatches";
  import { canCreateInvoicingProfile, clientLabel } from "$lib/clientWorkspace";
  import {
    type ClientWorkflowTarget,
    clearClientWorkflowDraft,
    keepClientWorkflowDraft,
    mergeClientWorkflowDraft,
    readClientWorkflowDraft,
    workflowDraftSnapshot,
  } from "$lib/clientWorkflowNavigation";
  let {
    request,
    allocations: initialAllocations = [],
    options,
    related = {},
    dirty = $bindable(false),
    onsaved,
    onkeepediting,
    oncancel,
  }: {
    request: WorkflowRecord;
    allocations?: WorkflowRecord[];
    options: WorkflowRecord;
    // The request detail's related records, which hold proposed rates.
    related?: WorkflowRecord;
    dirty?: boolean;
    onsaved?: () => Promise<void>;
    // Reloads a saved request for Save and keep editing, which stays on the
    // form. A new request has the button without it: it opens its edit page.
    onkeepediting?: () => Promise<void>;
    // Shows Cancel beside Save, as the other editors do.
    oncancel?: () => void;
  } = $props();
  clients.init();
  jobs.init();
  divisions.init();
  rateSheets.init();
  type FormState = {
    item: WorkflowRecord;
    allocations: WorkflowRecord[];
    fee: string;
    rateMode: string;
    proposal: ProposedRates;
  };
  // The rates editor adds an empty row for each role with no rate. Those rows
  // are left out, so they do not make the request unsaved until a rate is
  // typed in one.
  const comparable = (part: string, value: unknown) => {
    if (part !== "proposal") return value;
    const rates = value as ProposedRates;
    return {
      ...rates,
      rows: rates.rows.filter((row) => row.rate !== null || row.overtime_rate !== null),
    };
  };
  // The request as saved on the server.
  const loaded: FormState = untrack(() => ({
    item: { project_award_date: "", project_completion_date: "", ...request },
    allocations: initialAllocations.map((row) => ({ division: row.division, hours: row.hours })),
    fee: request.time_and_materials_set ? (request.time_and_materials ? "tm" : "fixed") : "",
    rateMode: related.rate_sheet_proposed ? "propose" : "choose",
    proposal: related.rate_sheet_proposed
      ? {
          code: related.rate_sheet?.code ?? "",
          name: related.rate_sheet?.name ?? "",
          reason: related.rate_sheet?.reason ?? "",
          based_on: related.rate_sheet?.based_on ?? "",
          rows: (related.rate_sheet_entries ?? [])
            .map((entry: WorkflowRecord) => ({
              role: entry.role_id,
              roleName: entry.role_name,
              rate: entry.rate,
              overtime_rate: entry.overtime_rate,
              source_rate: entry.source_rate,
              source_overtime_rate: entry.source_overtime_rate,
            }))
            .sort((a: ProposedRate, b: ProposedRate) => a.roleName.localeCompare(b.roleName)),
        }
      : emptyProposedRates(),
  }));
  // A draft kept while the user was in the client workspace keeps what the
  // user changed. The rest, and the version, come from the request as loaded
  // now, since an absorb there may have changed it.
  const draft = untrack(() =>
    readClientWorkflowDraft<FormState & { saved: string }>(
      `request:${request.id || "new"}`,
      pb.authStore.record?.id || "",
    ),
  );
  const start = draft
    ? mergeClientWorkflowDraft<FormState>(draft, draft.saved, loaded, ["item"], comparable)
    : loaded;
  // The request either chooses an active rate sheet or proposes new rates.
  // Find by rates is a tool between the two, not a state of the request: it
  // ends with Use on a sheet (choose) or with a switch to Propose new rates,
  // which carries the rates into the proposal (propose). While it is open,
  // the request keeps the state it had before, findFrom, which a save would
  // not send, so a save is refused.
  const rateModeOptions = [
    { id: "choose", label: "Choose a sheet" },
    { id: "find", label: "Find by rates" },
    { id: "propose", label: "Propose new rates" },
  ];
  let rateMode = $state(start.rateMode);
  let findFrom = $state(untrack(() => rateMode));
  // The state of the request: choose or propose, also while Find by rates is
  // open.
  const requestRateMode = () => (rateMode === "find" ? findFrom : rateMode);
  // The rates entered in Find by rates. They stay while the requester looks
  // at another choice, until the editor closes.
  let findRows = $state<ProposedRate[]>([]);
  // The rates carried from Find by rates into the proposal. Each starting
  // sheet chosen in the proposal takes them in place of its own rates for
  // those roles, until the requester leaves Propose new rates.
  let carried = $state<FoundRate[]>([]);
  // The rates carried at the last switch to Propose new rates, so that a
  // switch replaces only the rates changed in Find by rates since.
  let lastCarried: FoundRate[] = [];
  let proposal = $state<ProposedRates>(start.proposal);
  // Active sheets with the same rates, from the last save that Turbo refused,
  // with the rates that the save sent. The notice shows only while the rates
  // are those: after a change, the editor's live comparison takes over.
  let refusal = $state<{ matches: RateSheetMatch[]; rates: string } | null>(null);
  const rateMatches = $derived(
    refusal && JSON.stringify(proposedRatesPayload(proposal).entries) === refusal.rates
      ? refusal.matches
      : null,
  );
  // Field errors of the proposed rates from the last save that Turbo refused.
  let proposalErrors = $state<Record<string, { message: string }>>({});
  let item = $state<WorkflowRecord>(start.item);
  let allocations = $state<WorkflowRecord[]>(start.allocations);
  let fee = $state(start.fee);
  const snapshot = () =>
    workflowDraftSnapshot(
      { item, allocations, fee, rateMode: requestRateMode(), proposal },
      comparable,
    );
  // The baseline is the request as saved, also after a draft is merged, so the
  // form is unsaved exactly when it differs from the server copy.
  let saved = $state(workflowDraftSnapshot(loaded, comparable));
  let openingClient = false;
  let form: HTMLFormElement | undefined = $state();
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
  // The approver list depends on the project value alone, so the dropdown
  // says so until there is one, as Client Contact does for the client.
  const hasProjectValue = $derived(Number.isInteger(item.project_value) && item.project_value > 0);
  const approverItems = $derived([
    {
      id: "",
      name: !hasProjectValue
        ? "Enter the project value first"
        : approversLoaded && approvers.length === 0
          ? "No approver covers this value"
          : "Select approver",
    },
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
    if (!hasProjectValue) return "Enter the project value in whole dollars to list approvers.";
    if (approversLoaded && approvers.length === 0)
      return "No approver has a limit that covers this project value.";
    return "";
  });
  $effect(() => {
    dirty = snapshot() !== saved;
  });
  // The request's draft is not in the Rate Sheet list, so choosing a sheet
  // starts from an empty choice.
  function setRateMode(mode: string) {
    if (mode === rateMode) return;
    if (mode === "find") findFrom = rateMode;
    rateMode = mode;
    refusal = null;
    if (mode === "propose") carryFoundRates();
    else carried = [];
    if (mode === "choose" && related.rate_sheet_proposed && item.rate_sheet === request.rate_sheet)
      item.rate_sheet = "";
  }
  // A switch to Propose new rates, from Find by rates or from Choose a sheet,
  // carries the rates entered in Find by rates that have no error. With a
  // starting sheet already chosen, the rates changed since the last switch
  // replace the proposal's rates for their roles at once; otherwise the
  // editor puts all of them in place when the requester chooses one. With
  // nothing entered, nothing changes.
  function carryFoundRates() {
    const rates = carriedFindRates(findRows);
    const changed = changedCarriedRates(rates, lastCarried);
    lastCarried = rates;
    carried = rates;
    if (proposal.based_on) proposal.rows = withCarriedRates(proposal.rows, changed);
  }
  // Choose a sheet that Find by rates found. The requester saves the request.
  function useFoundSheet(id: string) {
    rateMode = "choose";
    item.rate_sheet = id;
    refusal = null;
    carried = [];
  }
  // Choose an active sheet that has the proposed rates, or rates close to
  // them, instead.
  function useMatchingSheet(id: string) {
    rateMode = "choose";
    item.rate_sheet = id;
    refusal = null;
    save();
  }
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
        $state.snapshot({ item, allocations, fee, rateMode: requestRateMode(), proposal, saved }),
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
  // Only procurement (the job claim) can create a client or a contact.
  const hasJobClaim = $derived($globalStore.claims.includes("job"));
  let showAddClient = $state(false);
  // The field that takes the new client: the client or the job owner.
  let newClientField = $state<"client" | "job_owner">("client");
  function openAddClient(field: "client" | "job_owner") {
    newClientField = field;
    showAddClient = true;
  }
  function clientCreated(client: { id: string }) {
    if (newClientField === "client") setClient(client.id);
    else item.job_owner = client.id;
  }
  function setClient(client: string) {
    if (client === item.client) return;
    item.client = client;
    item.contact = "";
    item.invoicing_information = "";
    item.proposal = "";
  }
  // A chosen proposal fills the divisions only while none is chosen. A
  // later choice, or a division chosen before the reply, discards the
  // reply. If the fetch fails, the requester adds the divisions by hand.
  let proposalFetch = 0;
  function setProposal(id: string) {
    if (id === item.proposal) return;
    item.proposal = id;
    const ticket = ++proposalFetch;
    if (!id || hasAllocations(allocations)) return;
    pb.send<WorkflowRecord>(`/api/jobs/${id}/details`, { method: "GET", requestKey: null })
      .then((source) => {
        if (ticket !== proposalFetch || item.proposal !== id || hasAllocations(allocations)) return;
        const rows = copiedAllocations(source.allocations, options.divisions);
        if (rows.length) allocations = rows;
      })
      .catch(() => {});
  }
  // stay keeps the form open after the save, for Save and keep editing.
  export async function save(stay = false): Promise<boolean> {
    if (saving) return false;
    error = "";
    refusal = null;
    proposalErrors = {};
    // Find by rates only searches. The requester chooses what it leads to.
    if (rateMode === "find") {
      error =
        "Not saved. Find by rates only searches. Select Use on a sheet, or choose Propose new rates.";
      return false;
    }
    const proposing = rateMode === "propose";
    // A save leaves out a role with no rate, so a role with only one of its
    // rates would lose it. Refuse the save and mark the empty rate instead.
    if (proposing) {
      const incomplete = incompleteRateErrors(proposal.rows);
      if (Object.keys(incomplete).length > 0) {
        proposalErrors = incomplete;
        error = "Not saved. Enter both rates of each role marked below, or leave both empty.";
        return false;
      }
    }
    saving = true;
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
      // Proposed rates replace the choice of a sheet.
      ...(proposing ? [] : ["rate_sheet"]),
      "approver",
    ];
    const payload = Object.fromEntries(
      fields.filter((field) => item[field] !== undefined).map((field) => [field, item[field]]),
    );
    if (fee) {
      payload.time_and_materials = fee === "tm";
    }
    // Kept so that an error on one of the entries sent goes on its row.
    const proposedRates = proposing ? proposedRatesPayload(proposal) : undefined;
    try {
      const result = await pb.send(item.id ? `/api/job_requests/${item.id}` : "/api/job_requests", {
        method: item.id ? "PATCH" : "POST",
        body: {
          version: item.version,
          request: payload,
          allocations,
          ...(proposedRates ? { proposed_rates: proposedRates } : {}),
        },
      });
      saved = snapshot();
      dirty = false;
      clearClientWorkflowDraft();
      if (!item.id) {
        const id = result.request?.id ?? result.id;
        await goto(
          stay ? resolve("/job-requests/[id]/edit", { id }) : resolve("/job-requests/[id]", { id }),
        );
      } else if (stay && onkeepediting) await onkeepediting();
      else if (onsaved) await onsaved();
      else await invalidateAll();
      return true;
    } catch (e) {
      const matches = duplicateRateSheetMatches(e, "proposed_rates");
      refusal =
        matches?.length && proposedRates
          ? { matches, rates: JSON.stringify(proposedRates.entries) }
          : null;
      proposalErrors = proposedRatesFieldErrors(e, proposedRates?.entries);
      // The callout beside the rates names the matching sheet. An error on one
      // rate shows on its row and an error on a field of the proposed rates on
      // that field, so the message here leaves them out.
      if (!matches?.length) error = workflowError(withoutShownErrors(e, proposedRates?.entries));
      return false;
    } finally {
      saving = false;
    }
  }
</script>

<!-- The same spacing and mobile text size as the other editors. -->
<form
  bind:this={form}
  class="flex w-full flex-col gap-2 max-lg:[&_button]:text-base max-lg:[&_input]:text-base max-lg:[&_label]:text-base max-lg:[&_select]:text-base max-lg:[&_textarea]:text-base"
  onsubmit={(event) => {
    event.preventDefault();
    save();
  }}
>
  {#if error}<div role="alert"><DSCallout tone="danger">{error}</DSCallout></div>{/if}
  <!-- A draft saves with fields missing, so the asterisks mark what the
    server checks when the request is submitted. -->
  <DSRequiredLegend label="Required to submit" />
  <fieldset disabled={saving} class="flex flex-col gap-2">
    <span class="flex w-full items-center gap-2">
      <label for="project_award_date">Project Award Date<DSRequiredMark /></label>
      <DsDateInput
        class="flex-1 md:flex-none"
        name="project_award_date"
        required
        bind:value={item.project_award_date}
      />
    </span>
    <span class="flex w-full items-center gap-2">
      <label for="project_completion_date"
        >Project Completion Date{#if item.project_status === "Active"}<DSRequiredMark />{/if}</label
      >
      <DsDateInput
        class="flex-1 md:flex-none"
        name="project_completion_date"
        required={item.project_status === "Active"}
        bind:value={item.project_completion_date}
      />
    </span>
    <div class="flex w-full flex-col gap-1">
      <span>Location<DSRequiredMark /></span>
      <DSLocationPicker bind:value={item.location} errors={{}} fieldName="location" />
    </div>
    <DsTextInput
      bind:value={item.description}
      errors={noFieldErrors}
      fieldName="description"
      uiName="Description"
      required
    />
    {#if $clients.index !== null}
      <div class="flex w-full items-end gap-1">
        <div class="flex-1">
          <DsAutoComplete
            bind:value={() => item.client, setClient}
            index={$clients.index}
            errors={noFieldErrors}
            fieldName="client"
            uiName="Client"
            required
            disabled={!!item.parent}
          >
            {#snippet resultTemplate(row)}{clientLabel({
                name: row.name,
                alias: row.alias,
              })}{/snippet}
          </DsAutoComplete>
        </div>
        {#if !item.client && !item.parent && hasJobClaim}
          <DSActionButton
            action={() => openAddClient("client")}
            icon="feather:plus-circle"
            color="green"
            title="Add new client"
          />
        {/if}
      </div>
    {/if}
    <div class="flex w-full items-end gap-1">
      <div class="flex-1">
        <DsSelector
          bind:value={item.contact}
          items={[{ id: "" }, ...contacts]}
          errors={noFieldErrors}
          fieldName="contact"
          fill
          uiName="Client Contact"
          required
          disabled={!item.client}
        >
          {#snippet optionTemplate(row: WorkflowRecord)}{row.id
              ? `${fullName(row)}${row.email ? ` — ${row.email}` : ""}`
              : item.client
                ? "Select a contact"
                : "Select a client first"}{/snippet}
        </DsSelector>
      </div>
      {#if item.client && !item.contact && hasJobClaim}
        <DSActionButton
          action={() => manageClient("contact")}
          icon="feather:plus-circle"
          color="green"
          title="Add new contact"
        />
      {/if}
    </div>
    {#if item.client && !profiles.length}
      <div class="flex flex-col gap-1">
        <h3 class="text-sm font-semibold">Invoicing profile<DSRequiredMark /></h3>
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
        fill
        uiName="Invoicing profile"
        required
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
      required
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
      <span id="request-fn-agreement">First Nations agreement<DSRequiredMark /></span>
      <DSToggle
        bind:value={item.fn_agreement_answer}
        options={answerOptions}
        ariaLabelledBy="request-fn-agreement"
        required
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
      stacked
      uiName="Project value (CAD, before tax)"
      required
    />
    <div class="flex w-full items-center gap-2">
      <span id="request-fee-basis">Fee basis<DSRequiredMark /></span>
      <DSToggle bind:value={fee} options={feeOptions} ariaLabelledBy="request-fee-basis" required />
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
    <div class="flex w-full flex-wrap items-center gap-2">
      <span id="request-rate-mode">Rates</span>
      <!-- Below lg the toggle takes its own full-width line, so the help sits
        beside the label instead of on a line of its own. -->
      <div class="max-lg:order-2 max-lg:w-full">
        <DSToggle
          bind:value={() => rateMode, setRateMode}
          options={rateModeOptions}
          ariaLabelledBy="request-rate-mode"
          narrowGrid
        />
      </div>
      {#if rateMode === "find"}
        <span class="max-lg:order-1">
          <DSHelp doc="proposed-rates#find-a-sheet-by-rates" iconOnly />
        </span>
      {/if}
    </div>
    {#if rateMode === "propose"}
      <ProposedRatesEditor
        bind:proposal
        {carried}
        onUse={useMatchingSheet}
        blocked={!!rateMatches?.length}
        errors={proposalErrors}
      />
    {:else if rateMode === "find"}
      <RateSheetFinder bind:rows={findRows} onUse={useFoundSheet} />
    {:else if $rateSheets.index !== null}
      <DSPicker
        bind:value={item.rate_sheet}
        items={$rateSheets.items.filter((sheet) => sheet.active)}
        index={$rateSheets.index}
        errors={noFieldErrors}
        fieldName="rate_sheet"
        uiName="Rate Sheet"
        required
        placeholder={rateSheetSearchPlaceholder}
        empty
      >
        {#snippet option(row, result)}<RateSheetOption sheet={row} {result} />{/snippet}
        {#snippet help()}
          <DSHelp doc="rate-sheets#choose-a-rate-sheet-for-a-project" iconOnly />
        {/snippet}
      </DSPicker>
    {/if}
    {#if rateMatches?.length}
      {@const match = rateMatches[0]}
      <RateSheetMatchCallout
        tone="danger"
        title="Not saved"
        intro="An active sheet has these rates. Use that sheet, or change a rate."
        matches={rateMatches}
      >
        {#snippet actions()}
          <DSActionButton action={() => useMatchingSheet(match.id)}
            >Use {rateSheetCodeRevision(match)}</DSActionButton
          >
        {/snippet}
      </RateSheetMatchCallout>
    {/if}
    <DSCallout tone="info"
      >The rate sheet stays with the project for its whole life. It must match the client's
      contract.</DSCallout
    >
    {#if $jobs.index !== null}
      <DsAutoComplete
        bind:value={() => item.proposal, setProposal}
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
    <legend class="font-semibold">Divisions<DSRequiredMark /></legend>
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
      <div class="flex w-full items-end gap-1">
        <div class="flex-1">
          <DsAutoComplete
            bind:value={item.job_owner}
            index={$clients.index}
            errors={noFieldErrors}
            fieldName="job_owner"
            uiName="Job Owner"
          >
            {#snippet resultTemplate(row)}{clientLabel({
                name: row.name,
                alias: row.alias,
              })}{/snippet}
          </DsAutoComplete>
        </div>
        {#if !item.job_owner && hasJobClaim}
          <DSActionButton
            action={() => openAddClient("job_owner")}
            icon="feather:plus-circle"
            color="green"
            title="Add new client as job owner"
          />
        {/if}
      </div>
    {/if}
    <DsSelector
      bind:value={item.branch}
      items={[{ id: "", name: "Select a branch" }, ...(options.branches || [])]}
      errors={noFieldErrors}
      fieldName="branch"
      fill
      uiName="Branch"
      required
    >
      {#snippet optionTemplate(row: WorkflowRecord)}{row.code ?? row.name}{/snippet}
    </DsSelector>
    <label class="flex flex-col gap-1">
      Notes
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
      fill
      uiName="Selected approver"
      required
      disabled={!hasProjectValue && !item.approver}
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
  <span class="flex w-full gap-2">
    <DSActionButton type="submit" disabled={saving}
      >{saving ? "Saving…" : "Save request"}</DSActionButton
    >
    {#if onkeepediting || !item.id}
      <!-- The same field checks as Save request, which the form runs on submit. -->
      <DSActionButton
        disabled={saving}
        action={() => {
          if (form?.reportValidity()) save(true);
        }}>Save and keep editing</DSActionButton
      >
    {/if}
    {#if oncancel}<DSActionButton action={oncancel}>Cancel</DSActionButton>{/if}
  </span>
</form>
<ClientCreatePopover bind:show={showAddClient} onCreated={clientCreated} />
