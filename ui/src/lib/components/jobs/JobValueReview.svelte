<script lang="ts">
  import type { WorkflowRecord } from "$lib/clientEditors";
  import DSPopover from "$lib/components/DSPopover.svelte";
  import DSFieldList from "$lib/components/DSFieldList.svelte";
  import DSField from "$lib/components/DSField.svelte";
  import DSActionButton from "$lib/components/DSActionButton.svelte";
  import DsCheck from "$lib/components/DsCheck.svelte";
  import { JobValue, valueCAD } from "$lib/jobValue.svelte";
  import { shortDate } from "$lib/utilities";
  // The selected approver's decision on one pending value change. Set decision
  // to open the approval or rejection dialog; each instance keeps its own
  // confirmation and rejection reason.
  let {
    value,
    item,
    decision = $bindable(null),
  }: {
    value: JobValue;
    item: WorkflowRecord;
    decision?: "approve" | "reject" | null;
  } = $props();
  // The revision under review. A correction, such as a replacement PDF, saves a
  // new version, which clears the confirmation and reason given for the earlier
  // one, so approval always follows a review of the data it sends. Opening a
  // dialog clears them too.
  const revision = $derived(`${item.id}:${item.version}:${decision}`);
  let reviewed = $derived.by(() => {
      void revision;
      return false;
    }),
    reason = $derived.by(() => {
      void revision;
      return "";
    });
  const noErrors = {};
  async function decide(name: "approve" | "reject") {
    if (await value.act(item, name, name === "reject" ? reason : "")) decision = null;
  }
  function close() {
    value.error = "";
    decision = null;
  }
  const job = $derived(item.job_number || "this job");
</script>

<DSPopover
  show={decision === "approve"}
  title="Approve value change"
  subtitle={`Set ${job} to ${valueCAD(item.project_value)}.`}
  error={value.error || null}
  submitting={value.busy}
  submitLabel="Approve value"
  submitDisabled={!reviewed}
  onSubmit={() => decide("approve")}
  onCancel={close}
>
  <DSFieldList>
    <DSField label="Current value" value={valueCAD(item.effective_project_value)} />
    <DSField label="Proposed value" value={valueCAD(item.project_value)} />
    <DSField
      label="Proposed by"
      value={`${item.uploader_name || item.uploader} on ${shortDate(item.uploaded, true)}`}
    />
  </DSFieldList>
  <div>
    <DSActionButton variant="plain" icon="mdi:file-pdf-box" action={() => value.pdf(item)}
      >Open authorization PDF</DSActionButton
    >
  </div>
  <DsCheck
    bind:value={reviewed}
    errors={noErrors}
    fieldName="reviewed"
    uiName="I reviewed this PDF and its full proposed value."
  />
</DSPopover>

<DSPopover
  show={decision === "reject"}
  title="Reject value change"
  subtitle={`The value of ${job} stays ${valueCAD(item.effective_project_value)}.`}
  error={value.error || null}
  submitting={value.busy}
  submitLabel="Reject"
  danger
  submitDisabled={!reason.trim()}
  onSubmit={() => decide("reject")}
  onCancel={close}
>
  <label class="flex flex-col gap-1">
    <span>Rejection reason</span>
    <textarea
      bind:value={reason}
      rows={3}
      class="rounded-sm border border-neutral-300 p-2"
      disabled={value.busy}></textarea>
  </label>
</DSPopover>
