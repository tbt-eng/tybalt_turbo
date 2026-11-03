<script lang="ts">
  import { resolve } from "$app/paths";
  import type { WorkflowRecord } from "$lib/clientEditors";
  import DSList from "$lib/components/DSList.svelte";
  import DSCallout from "$lib/components/DSCallout.svelte";
  import DSActionButton from "$lib/components/DSActionButton.svelte";
  import { JobValue, valueCAD } from "$lib/jobValue.svelte";
  import { shortDate } from "$lib/utilities";
  import JobValueReview from "./JobValueReview.svelte";
  // The signed-in approver's pending value changes.
  type Change = WorkflowRecord & { id: string };
  const value = new JobValue(undefined, true);
  $effect(() => {
    value.reload();
  });
  // The change whose approval or rejection dialog is open.
  let reviewing = $state<WorkflowRecord | null>(null),
    decision = $state<"approve" | "reject" | null>(null);
  function open(item: WorkflowRecord, name: "approve" | "reject") {
    value.error = "";
    reviewing = item;
    decision = name;
  }
</script>

<div class="space-y-4">
  {#if value.error && !decision}<DSCallout tone="danger">{value.error}</DSCallout>{/if}

  <DSList
    items={value.items as Change[]}
    search={true}
    pagination={true}
    inListHeader="Value Changes Pending My Approval"
    emptyMessage={value.loading
      ? "Loading value changes…"
      : "No value changes are waiting for your approval."}
  >
    {#snippet anchor(item)}
      <a class="link" href={resolve(`/jobs/${item.job}/details`)}>{item.job_number || item.job}</a>
    {/snippet}
    {#snippet headline(item)}
      {item.job_description || "Job details"}
    {/snippet}
    {#snippet byline(item)}
      {valueCAD(item.effective_project_value)} → {valueCAD(item.project_value)}
    {/snippet}
    {#snippet line1(item)}
      {[item.related?.client?.name, item.related?.manager?.name].filter(Boolean).join(" · ")}
    {/snippet}
    {#snippet line2(item)}
      <span class="flex flex-wrap items-center gap-2">
        Proposed by {item.uploader_name || item.uploader} on {shortDate(item.uploaded, true)}
        <DSActionButton variant="plain" icon="mdi:file-pdf-box" action={() => value.pdf(item)}
          >Authorization PDF</DSActionButton
        >
      </span>
    {/snippet}
    {#snippet actions(item)}
      <DSActionButton
        action={() => open(item, "approve")}
        icon="mdi:approve"
        title="Approve"
        color="green"
      />
      <DSActionButton
        action={() => open(item, "reject")}
        icon="mdi:cancel"
        title="Reject"
        color="orange"
      />
    {/snippet}
  </DSList>
</div>

{#if reviewing}
  <JobValueReview {value} item={reviewing} bind:decision />
{/if}
