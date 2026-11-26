<script lang="ts">
  import DSHelp from "$lib/components/DSHelp.svelte";
  import DSActionButton from "$lib/components/DSActionButton.svelte";
  import DSPopover from "$lib/components/DSPopover.svelte";
  import DSTextInput from "$lib/components/DSTextInput.svelte";
  import { pb } from "$lib/pocketbase";
  import { goto, invalidateAll } from "$app/navigation";
  import { globalStore } from "$lib/stores/global";
  import { jobRequestsEnabled, jobsEditingEnabled } from "$lib/stores/appConfig";
  import { JobsStatusOptions } from "$lib/pocketbase-types";
  import DSLocationPicker from "$lib/components/DSLocationPicker.svelte";
  import ProposedRatesReview from "$lib/components/rateSheets/ProposedRatesReview.svelte";
  import {
    emptyChecklist,
    contactEmail,
    fullName,
    requestNextStep,
    wholeCAD,
  } from "$lib/jobRequests";
  import { workflowError, type WorkflowRecord } from "$lib/clientEditors";
  import { clientLabel } from "$lib/clientWorkspace";
  import { rateSheetLabel } from "$lib/rateSheetLabel";
  import DSSection from "$lib/components/DSSection.svelte";
  import DSPageHeader from "$lib/components/DSPageHeader.svelte";
  import DSCallout from "$lib/components/DSCallout.svelte";
  import DSCard from "$lib/components/DSCard.svelte";
  import DSCardGrid from "$lib/components/DSCardGrid.svelte";
  import DSField from "$lib/components/DSField.svelte";
  import DSFieldList from "$lib/components/DSFieldList.svelte";
  import DsCheck from "$lib/components/DsCheck.svelte";
  import DSList from "$lib/components/DSList.svelte";
  import DSStatusLabel from "$lib/components/DSStatusLabel.svelte";
  import DsLabel from "$lib/components/DsLabel.svelte";
  import { shortDate } from "$lib/utilities";
  import { resolve } from "$app/paths";
  // The duplicate-work review (GET /api/job_requests/{id}/duplicates). The
  // server orders each list by tier, so the groups appear exact matches first.
  type DuplicateTier = "exact" | "similar" | "client";
  type DuplicateJob = {
    id: string;
    number: string;
    description: string;
    status: string;
    manager_name: string;
    parent_number: string;
    reasons: string[];
    tier: DuplicateTier;
    last_activity: string;
  };
  type DuplicateRequest = {
    id: string;
    description: string;
    status: string;
    reasons: string[];
    tier: DuplicateTier;
  };
  type RequestDuplicates = {
    jobs: DuplicateJob[];
    hidden_client_jobs: number;
    requests: DuplicateRequest[];
  };
  const duplicateTierLabels: Record<string, string> = {
    exact: "Exact matches",
    similar: "Same client, similar description",
    client: "The client's other recent open projects",
  };
  const duplicateJobSearchText = (job: DuplicateJob) =>
    [job.number, job.description, job.manager_name, ...job.reasons].join(" ");
  // Outside the exact matches, the group heading already says "same client".
  const shownReasons = (row: { tier: DuplicateTier; reasons: string[] }) =>
    row.tier === "exact" ? row.reasons : row.reasons.filter((reason) => reason !== "Same client");
  let duplicateFilter = $state("");
  let { data } = $props();
  const duplicates = $derived<RequestDuplicates>({
    jobs: data.duplicates?.jobs ?? [],
    hidden_client_jobs: data.duplicates?.hidden_client_jobs ?? 0,
    requests: data.duplicates?.requests ?? [],
  });
  // Procurement scans the whole list, so say plainly what it leaves out.
  const hiddenClientJobsNote = $derived.by(() => {
    const hidden = duplicates.hidden_client_jobs;
    if (hidden === 0) return "";
    return hidden === 1
      ? "1 other project of the client is not shown because it is closed, cancelled, or has no award and no time in the last year."
      : `${hidden} other projects of the client are not shown because they are closed, cancelled, or have no award and no time in the last year.`;
  });
  let busy = $state(false),
    error = $state(""),
    // A refused workflow action shows beside the workflow buttons.
    actionError = $state(""),
    confirmingDelete = $state(false),
    confirmingRemovePDF = $state(false),
    // Notes may be internal, so the printout leaves them out unless chosen.
    printNotes = $state(false);
  let checklist = $state<Record<string, boolean>>({}),
    rejectionReason = $state(""),
    // Procurement can correct the name of proposed rates as they create the job.
    rateSheetName = $state(""),
    // The approver approves proposed rates beside the rates, not in the checklist.
    ratesApproved = $state(false);
  let fileInput: HTMLInputElement | undefined = $state();
  const detail = $derived(data.detail),
    request = $derived(detail.request),
    permission = $derived(detail.permissions),
    related = $derived(detail.related);
  const viewerId = pb.authStore.record?.id ?? "";
  // Every viewer sees who acts next, so a page without buttons still explains
  // where the request stands.
  const nextStep = $derived(
    requestNextStep({
      status: request.status || request.state,
      creatorName: fullName(related.creator),
      approverName: fullName(related.approver),
      viewerIsCreator: viewerId !== "" && request.creator === viewerId,
      viewerIsApprover: viewerId !== "" && request.approver === viewerId,
      canCreate: permission.can_create,
    }),
  );
  // The Help link for a viewer who can still change the request's path, such
  // as a creator who can recall it or edit a rejected one. A viewer with no
  // action reads the statuses instead.
  const canActOnRequest = $derived(
    permission.can_edit || permission.can_recall || permission.can_reject,
  );
  // The request is changed on its edit page, so no action here can meet
  // unsaved changes.
  const canEdit = $derived(permission.can_edit && $jobRequestsEnabled);
  const editHref = $derived(resolve("/job-requests/[id]/edit", { id: request.id }));
  const questions = $derived(detail.checklist.questions);
  const allChecked = $derived(
    questions.length > 0 && questions.every((question) => checklist[question.key]),
  );
  // The page is reused when a link opens another request. Ask again there,
  // but keep the choice when this request reloads, as an upload does.
  const requestId = $derived(request.id);
  $effect(() => {
    requestId;
    printNotes = false;
  });
  $effect(() => {
    // The server digest includes submission_token. Resubmission changes it even
    // when the request version/PDF are unchanged; require fresh answers. Every
    // reload resets them, so a job holder who also approved answers again.
    detail.review_token;
    request.version;
    checklist = emptyChecklist(detail.checklist.questions);
    ratesApproved = false;
  });
  // The name starts as the saved name. A reload gives related new objects, so
  // the effect depends only on the saved name and the request id, which are
  // strings: it resets the typed name only when either of them changes.
  const savedRateSheetName = $derived(
    related.rate_sheet_proposed ? (related.rate_sheet?.name ?? "") : "",
  );
  $effect(() => {
    requestId;
    rateSheetName = savedRateSheetName;
  });
  async function refresh() {
    await Promise.all([invalidateAll(), globalStore.refreshAttentionCounts()]);
  }
  // A request can be submitted and approved while its proposal is In Progress
  // or Submitted. The job is created only from an Awarded proposal, and only
  // procurement can award it, here or on the proposal.
  const proposal = $derived(related.proposal as WorkflowRecord | undefined);
  const proposalOpen = $derived(
    proposal?.status === JobsStatusOptions["In Progress"] ||
      proposal?.status === JobsStatusOptions.Submitted,
  );
  // Only procurement can act on the status, so only procurement gets the
  // notice. Everyone sees the status on the Proposal field.
  const jobHolder = $derived($globalStore.claims.includes("job"));
  const canAward = $derived(jobHolder && $jobsEditingEnabled && $jobRequestsEnabled);
  // The award needs a value or Time and Materials. Say so before the attempt.
  const needsValue =
    "The proposal needs a Proposal Value or Time and Materials before it can be awarded.";
  const proposalNeedsValue = $derived(
    !!proposal && !(Number(proposal.proposal_value) > 0) && !proposal.time_and_materials,
  );
  // Why an award failed, in plain words: the known reason, else what the
  // server says about each field, else its message.
  function awardFailure(e: unknown): string {
    const failure = e as { data?: WorkflowRecord; response?: WorkflowRecord };
    const body = failure?.response ?? failure?.data;
    const fields = Object.values((body?.data ?? {}) as Record<string, WorkflowRecord>);
    if (fields.some((field) => field?.code === "value_required_for_status")) return needsValue;
    const reasons = fields.map((field) => field?.message).filter(Boolean);
    return reasons.length ? reasons.join(" ") : workflowError(e);
  }
  let awardError = $state(""),
    // Shown after a successful award.
    awardedNotice = $state(""),
    // The proposal of a failed award, for the link to edit it.
    awardFailedId = $state(""),
    // Set when creation answers proposal_not_awarded: the proposal to award.
    confirmAwardId = $state(""),
    confirmingAward = $state(false);
  // Award through the jobs API as the signed-in user, so the proposal's own
  // rules apply. Returns false and shows the server's reason on failure.
  async function awardProposal(id: string): Promise<boolean> {
    awardError = "";
    awardFailedId = "";
    awardedNotice = "";
    try {
      await pb.collection("jobs").update(id, { status: JobsStatusOptions.Awarded });
      awardedNotice = `Proposal ${proposal?.number} is now Awarded.`;
      return true;
    } catch (e) {
      awardError = awardFailure(e);
      awardFailedId = id;
      return false;
    }
  }
  async function award() {
    if (!proposal) return;
    busy = true;
    error = "";
    try {
      if (await awardProposal(proposal.id)) await refresh();
    } finally {
      busy = false;
    }
  }
  // The award leaves the review digest unchanged, so the same answers and
  // digest create the job.
  async function awardAndCreate() {
    busy = true;
    const awarded = await awardProposal(confirmAwardId);
    busy = false;
    confirmingAward = false;
    if (!awarded) return;
    await action("create");
    // A creation refused for another reason leaves the page open: show the
    // proposal as Awarded.
    if (actionError) await refresh();
  }
  // The proposal id of a proposal_not_awarded refusal, in the shape the jobs
  // hook uses: data.proposal.code and data.proposal.data.proposal_id.
  type FieldErrors = { data?: { proposal?: { code?: string; data?: { proposal_id?: unknown } } } };
  function proposalNotAwarded(e: unknown): string {
    const failure = e as { data?: FieldErrors; response?: FieldErrors };
    const field = (failure?.response ?? failure?.data)?.data?.proposal;
    const id = field?.data?.proposal_id;
    return field?.code === "proposal_not_awarded" && typeof id === "string" ? id : "";
  }
  async function action(name: string) {
    busy = true;
    error = "";
    actionError = "";
    awardError = "";
    try {
      const result = await pb.send(
        `/api/job_requests/${request.id}${name === "delete" ? "" : `/${name}`}`,
        {
          method: name === "delete" ? "DELETE" : "POST",
          body: {
            version: request.version,
            // Echo the reviewed digest, never reconstruct it from current data.
            // It binds this page's checklist to the server's submission_token;
            // an old page must fail after recall and resubmission.
            review_token: detail.review_token,
            document_hash: detail.attachment?.document_hash,
            checklist_version: detail.checklist.version,
            checklist_answers: checklist,
            rejection_reason: rejectionReason,
            ...(name === "approve" && related.rate_sheet_proposed
              ? { rates_approved: ratesApproved }
              : {}),
            ...(name === "create" && related.rate_sheet_proposed
              ? { rate_sheet_name: rateSheetName }
              : {}),
          },
        },
      );
      if (name === "create") {
        await globalStore.refreshAttentionCounts();
        await goto(`/jobs/${result.job?.id ?? result.id ?? result.job_id}/details`);
      } else if (name === "delete") {
        await globalStore.refreshAttentionCounts();
        // A creator deleting their own request returns to their list.
        await goto(permission.can_edit ? "/job-requests" : "/job-requests?queue=abandoned");
      } else await refresh();
    } catch (e) {
      const unawarded = name === "create" && canAward ? proposalNotAwarded(e) : "";
      if (unawarded) {
        confirmAwardId = unawarded;
        confirmingAward = true;
      }
      // Delete is confirmed in a dialog, which shows its error. The other
      // actions show theirs in the Workflow actions section.
      else if (name === "delete") error = workflowError(e);
      else actionError = workflowError(e);
    } finally {
      busy = false;
    }
  }
  async function upload(event: Event) {
    const input = event.currentTarget as HTMLInputElement,
      file = input.files?.[0];
    if (!file) return;
    busy = true;
    error = "";
    try {
      const body = new FormData();
      body.append("document", file);
      body.append("version", String(request.version));
      await pb.send(`/api/job_requests/${request.id}/document`, { method: "POST", body });
      await refresh();
    } catch (e) {
      error = workflowError(e);
    } finally {
      busy = false;
      input.value = "";
    }
  }
  // The server deletes the attachment and its file outright.
  async function removePDF() {
    busy = true;
    error = "";
    try {
      await pb.send(`/api/job_requests/${request.id}/document`, {
        method: "DELETE",
        body: { version: request.version },
      });
      confirmingRemovePDF = false;
      await refresh();
    } catch (e) {
      error = workflowError(e);
    } finally {
      busy = false;
    }
  }
  async function openPDF() {
    busy = true;
    error = "";
    try {
      const response = await fetch(
        `${pb.baseURL}/api/job_attachments/${detail.attachment?.id}/document`,
        { headers: { Authorization: pb.authStore.token } },
      );
      if (!response.ok) throw new Error("Cannot load the authorization PDF.");
      const url = URL.createObjectURL(await response.blob());
      window.open(url, "_blank");
      setTimeout(() => URL.revokeObjectURL(url), 60000);
    } catch (e) {
      error = workflowError(e);
    } finally {
      busy = false;
    }
  }
  // A reviewer must see what a request still lacks, so a field the request
  // needs shows "Not provided" when it is empty. Other empty fields are hidden.
  const missing = "Not provided";
  const withEmail = (person?: WorkflowRecord) =>
    fullName(person) ? [fullName(person), contactEmail(person)].filter(Boolean).join(" · ") : "";
  const feeBasis = $derived(
    request.time_and_materials_set
      ? request.time_and_materials
        ? "Time and materials"
        : "Fixed fee"
      : "",
  );
  async function print() {
    const printPage = { id: request.id };
    // eslint-disable-next-line svelte/no-navigation-without-resolve -- query on a resolved path
    if (printNotes) await goto(`${resolve("/job-requests/[id]/print", printPage)}?notes=1`);
    else await goto(resolve("/job-requests/[id]/print", printPage));
  }
</script>

<!-- The same padding and header as the other detail pages, such as a job. -->
<div class="mx-auto space-y-4 p-4">
  <DSPageHeader title="Project Request" subtitle={request.description || "Untitled request"}>
    {#snippet badges()}
      <DSStatusLabel kind="request" status={request.status || request.state} />
    {/snippet}
    {#snippet context()}
      Version {request.version} · Requested by {fullName(related.creator) || request.creator} · Selected
      approver: {fullName(related.approver) || "Not selected"}
    {/snippet}
    {#snippet actions()}
      {#if request.notes}
        <div>
          <DsCheck
            bind:value={printNotes}
            errors={{}}
            fieldName="print_notes"
            uiName="Include notes"
          />
        </div>
      {/if}
      <DSActionButton action={print}>Print Project Authorization</DSActionButton>
      {#if canEdit}
        <DSActionButton action={editHref}>Edit request</DSActionButton>
      {/if}
    {/snippet}
  </DSPageHeader>
  {#if error}<div role="alert"><DSCallout tone="danger">{error}</DSCallout></div>{/if}
  {#if !$jobRequestsEnabled}
    <DSCallout tone="warning">Request processing is disabled. Printing remains available.</DSCallout
    >
  {/if}
  {#snippet awardAction()}
    {#if proposal && proposalNeedsValue}
      <DSActionButton action={resolve("/jobs/[jid]/edit", { jid: proposal.id })}
        >Edit the proposal</DSActionButton
      >
    {:else}
      <DSActionButton disabled={busy} action={award}>Award proposal</DSActionButton>
    {/if}
  {/snippet}
  {#if proposal && proposalOpen && jobHolder}
    <DSCallout
      tone={proposalNeedsValue ? "warning" : "info"}
      title={`Proposal ${proposal.number} is ${proposal.status}`}
      actions={canAward ? awardAction : undefined}
    >
      {#if proposalNeedsValue}
        {needsValue}
      {:else}
        Award it before you create the job.
      {/if}
    </DSCallout>
  {/if}
  {#if awardedNotice}<div role="status">
      <DSCallout tone="success">{awardedNotice}</DSCallout>
    </div>{/if}
  {#if awardError}<div role="alert">
      <DSCallout tone="danger" title="Cannot award the proposal">
        <p>{awardError}</p>
        <p>
          <a class="link" href={resolve("/jobs/[jid]/edit", { jid: awardFailedId })}
            >Edit the proposal</a
          >
        </p>
      </DSCallout>
    </div>{/if}
  {#if request.rejection_reason}
    <DSCallout tone="danger" title="Rejected">
      Rejection reason: {request.rejection_reason}. The creator must make a correction before
      submission.
    </DSCallout>
  {/if}
  <DSCardGrid>
    <DSCard title="Project">
      <DSFieldList>
        <DSField label="Client" value={related.client ? clientLabel(related.client) : missing} />
        <DSField label="Project contact" value={withEmail(related.contact) || missing} />
        <DSField label="Manager" value={fullName(related.manager) || missing} />
        <DSField label="Alternate manager" value={fullName(related.alternate_manager)} />
        <DSField label="Branch" value={related.branch?.name || missing} />
        <DSField
          label="Job owner"
          value={related.job_owner ? clientLabel(related.job_owner) : ""}
        />
        <!-- A location that is set shows on the map in its own card. -->
        {#if !request.location}<DSField label="Location" value={missing} />{/if}
        <DSField label="Intended project status" value={request.project_status} />
        <DSField
          label="Award date"
          value={request.project_award_date ? shortDate(request.project_award_date, true) : missing}
        />
        <DSField
          label="Completion date"
          value={request.project_completion_date
            ? shortDate(request.project_completion_date, true)
            : missing}
        />
        <DSField
          label="First Nations agreement"
          value={request.fn_agreement_answer || "Not answered"}
        />
        {#if related.parent}
          <DSField
            label="Parent"
            value={`${related.parent.number} — ${related.parent.description}`}
            href={resolve("/jobs/[jid]/details", { jid: related.parent.id })}
          />
        {/if}
        {#if related.proposal}
          <DSField
            label="Proposal"
            value={`${related.proposal.number} — ${related.proposal.description}`}
            href={resolve("/jobs/[jid]/details", { jid: related.proposal.id })}
            hint={related.proposal.status}
          />
        {/if}
        {#if request.notes}
          <DSField label="Notes">
            <span class="whitespace-pre-wrap">{request.notes}</span>
          </DSField>
        {/if}
      </DSFieldList>
    </DSCard>
    <DSCard title="Value and rates">
      <DSFieldList>
        <DSField label="Project value" value={wholeCAD(request.project_value) || missing} />
        <DSField label="Fee basis" value={feeBasis || missing} />
        <DSField
          label="Rate sheet"
          value={related.rate_sheet ? rateSheetLabel(related.rate_sheet) : missing}
          hint={related.rate_sheet_proposed ? "Proposed rates" : undefined}
        />
        <DSField label="Client PO" value={request.client_po} />
        <DSField label="Client reference" value={request.client_reference_number} />
      </DSFieldList>
    </DSCard>
    <DSCard title="Invoicing">
      <DSFieldList>
        <DSField label="Invoicing profile" value={related.invoicing_information?.name} />
        <DSField label="Bill to" value={related.bill_to_name || missing} />
        <DSField label="Billing contact" value={withEmail(related.billing_contact) || missing} />
        {#if related.invoicing_information?.invoicing_instructions}
          <DSField label="Invoicing instructions">
            <span class="whitespace-pre-wrap"
              >{related.invoicing_information.invoicing_instructions}</span
            >
          </DSField>
        {/if}
      </DSFieldList>
    </DSCard>
    <DSCard title="Division hours">
      <DSFieldList>
        {#each detail.allocations as allocation (allocation.id)}
          <DSField
            label={[allocation.expand?.division?.code, allocation.expand?.division?.name]
              .filter(Boolean)
              .join(" — ") || "Unknown division"}
            value={`${allocation.hours} h`}
          />
        {:else}
          <DSField label="Divisions" value={missing} />
        {/each}
      </DSFieldList>
    </DSCard>
    {#if request.location}
      <DSCard title="Location">
        <DSLocationPicker
          value={request.location}
          errors={{}}
          fieldName="location"
          disabled
          readonly
          square
        />
      </DSCard>
    {/if}
  </DSCardGrid>
  {#if related.rate_sheet_proposed}<ProposedRatesReview
      {related}
      canApprove={permission.can_approve}
      bind:approved={ratesApproved}
    />{/if}
  <DSSection title="Authorization evidence">
    {#snippet actions()}
      <DSHelp doc="project-requests#upload-and-submit" label="Help" />
    {/snippet}
    {#if detail.attachment}
      <div>
        <DSActionButton variant="plain" disabled={busy} action={openPDF}
          >Open authorization PDF</DSActionButton
        >
      </div>
      <p class="text-xs text-neutral-500" title={detail.attachment.document_hash}>
        SHA-256 (first 7): {detail.attachment.document_hash.slice(0, 7)}
      </p>{:else}<p>No authorization PDF uploaded.</p>{/if}
    {#if canEdit}
      <!-- The button opens the file dialog of the hidden input. -->
      <div class="flex flex-wrap gap-2">
        <DSActionButton variant="secondary" disabled={busy} action={() => fileInput?.click()}
          >{detail.attachment ? "Replace PDF" : "Upload PDF"}</DSActionButton
        >
        {#if detail.attachment}<DSActionButton
            variant="secondary"
            disabled={busy}
            action={() => {
              error = "";
              confirmingRemovePDF = true;
            }}>Remove PDF</DSActionButton
          >{/if}
        <input
          bind:this={fileInput}
          class="hidden"
          type="file"
          accept="application/pdf,.pdf"
          aria-label={detail.attachment ? "Replace PDF" : "Upload PDF"}
          disabled={busy}
          onchange={upload}
        />
      </div>
      <p class="text-sm text-neutral-600">
        Upload the signed copy, including the rate appendix for T&amp;M work, or other sufficient
        authorization evidence.
      </p>
      <DSPopover
        bind:show={confirmingRemovePDF}
        title="Remove PDF"
        danger
        submitLabel="Remove"
        submitting={busy}
        {error}
        onSubmit={removePDF}
      >
        <p>Turbo deletes this PDF. You can't undo this. Upload a PDF before you submit.</p>
      </DSPopover>{/if}
  </DSSection>
  {#if permission.can_create}
    <DSSection title="Duplicate-work review">
      <p>
        Possible matches require review. A match alone does not prove that the work is duplicated.
      </p>
      <p>
        The list shows exact matches of any age or status (same proposal, parent, client PO or
        client reference) and the client's open projects with an award or time in the last year.
        {hiddenClientJobsNote}
      </p>
      <DSList
        items={duplicates.jobs}
        search
        searchText={duplicateJobSearchText}
        bind:searchTerm={duplicateFilter}
        searchPlaceholder="Filter by number, description, manager or reason"
        groupField="tier"
        emptyMessage={duplicateFilter
          ? "No listed job matches this filter."
          : "No possible matching jobs found."}
      >
        {#snippet groupHeader(tier: string)}{duplicateTierLabels[tier] ?? tier}{/snippet}
        {#snippet headline(job: DuplicateJob)}<a
            class="link"
            href={resolve("/jobs/[jid]/details", { jid: job.id })}>{job.number}</a
          >{/snippet}
        {#snippet byline(job: DuplicateJob)}
          <span class="flex flex-wrap items-center gap-2">
            <span>{job.description}</span>
            <DSStatusLabel kind="job" status={job.status} size="sm" />
            {#if job.manager_name}<span>{job.manager_name}</span>{/if}
            {#if job.parent_number}<span>Parent {job.parent_number}</span>{/if}
            {#if job.last_activity}<span>Last activity {shortDate(job.last_activity, true)}</span
              >{/if}
            {#each shownReasons(job) as reason (reason)}<DsLabel size="sm">{reason}</DsLabel>{/each}
          </span>
        {/snippet}
      </DSList>
      <DSList
        items={duplicates.requests}
        inListHeader="Other pending requests"
        emptyMessage="No possible matching requests found."
      >
        {#snippet headline(other: DuplicateRequest)}<a
            class="link"
            href={resolve("/job-requests/[id]", { id: other.id })}>{other.id}</a
          >{/snippet}
        {#snippet byline(other: DuplicateRequest)}
          <span class="flex flex-wrap items-center gap-2">
            <span>{other.description}</span>
            <DSStatusLabel kind="request" status={other.status} size="sm" />
            {#each other.reasons as reason (reason)}<DsLabel size="sm">{reason}</DsLabel>{/each}
          </span>
        {/snippet}
      </DSList>
    </DSSection>
  {/if}
  {#snippet checklistAnswers()}
    {#each questions as question (question.key)}<label class="flex items-start gap-2"
        ><input
          type="checkbox"
          bind:checked={checklist[question.key]}
          class="mt-1"
        />{question.text}</label
      >{/each}
  {/snippet}
  {#if $jobRequestsEnabled}
    <DSSection title="Workflow actions">
      {#snippet actions()}
        {#if permission.can_approve}
          <DSHelp doc="project-requests#approve-a-request" label="Help" />
        {:else if permission.can_create}
          <DSHelp doc="project-requests#create-the-job" label="Help" />
        {:else if permission.can_submit}
          <DSHelp doc="project-requests#upload-and-submit" label="Help" />
        {:else if canActOnRequest}
          <DSHelp doc="project-requests#change-a-request" label="Help" />
        {:else}
          <DSHelp doc="project-requests#request-statuses" label="Help" />
        {/if}
      {/snippet}
      {#if nextStep}<p>{nextStep}</p>{/if}
      {#if detail.approval_blocked_reason}
        <DSCallout tone="warning">{detail.approval_blocked_reason}</DSCallout>
      {/if}
      {#if actionError}<div role="alert">
          <DSCallout tone="danger">{actionError}</DSCallout>
        </div>{/if}
      {#if permission.can_approve}<div class="space-y-3">
          {@render checklistAnswers()}
          <div class="flex">
            <DSActionButton
              color="green"
              disabled={busy || !allChecked || (related.rate_sheet_proposed && !ratesApproved)}
              action={() => action("approve")}>Approve request</DSActionButton
            >
          </div>
        </div>{/if}
      {#if permission.can_create}<div class="space-y-3">
          {@render checklistAnswers()}
          {#if related.rate_sheet_proposed}
            <!-- The server keeps the requester's name when this is empty. -->
            <DSTextInput
              bind:value={rateSheetName}
              errors={{}}
              fieldName="rate_sheet_name"
              uiName="Rate sheet name"
              placeholder="Optional description"
              helperText={`Optional. The new sheet's code is ${related.rate_sheet?.code}. Leave the name empty to keep the requester's name.`}
            />
          {/if}
          <div class="flex">
            <DSActionButton
              color="green"
              disabled={busy || !allChecked}
              action={() => action("create")}>Create Job</DSActionButton
            >
          </div>
        </div>{/if}
      {#if permission.can_submit || permission.can_recall}
        <div class="flex flex-wrap gap-2">
          {#if permission.can_submit}
            <DSActionButton color="blue" disabled={busy} action={() => action("submit")}
              >Submit</DSActionButton
            >
          {/if}
          {#if permission.can_recall}
            <DSActionButton disabled={busy} action={() => action("recall")}
              >Recall for editing</DSActionButton
            >
          {/if}
        </div>
      {/if}
      {#if permission.can_reject}<div class="space-y-2">
          <DSTextInput
            bind:value={rejectionReason}
            errors={{}}
            fieldName="rejection_reason"
            uiName="Rejection reason"
            rows={3}
          />
          <div class="flex">
            <DSActionButton
              color="red"
              disabled={busy || !rejectionReason.trim()}
              action={() => action("reject")}>Reject</DSActionButton
            >
          </div>
        </div>{/if}
    </DSSection>
  {/if}
  <DSPopover
    bind:show={confirmingAward}
    title="Award the proposal"
    submitLabel="Award and create"
    submitting={busy}
    onSubmit={awardAndCreate}
  >
    <p>
      Proposal {proposal?.number} is not marked as Awarded. Set it to Awarded and continue? Turbo then
      creates the job.
    </p>
  </DSPopover>
  {#if permission.can_delete}
    {@const deleteLabel = permission.can_edit ? "Delete request" : "Delete abandoned request"}
    <DSActionButton color="red" disabled={busy} action={() => (confirmingDelete = true)}>
      {deleteLabel}
    </DSActionButton>
    <DSPopover
      bind:show={confirmingDelete}
      title={deleteLabel}
      danger
      submitLabel="Delete"
      submitting={busy}
      {error}
      onSubmit={() => action("delete")}
    >
      <p>Turbo removes this request, its PDF and its division hours. You can't undo this.</p>
    </DSPopover>
  {/if}
</div>
