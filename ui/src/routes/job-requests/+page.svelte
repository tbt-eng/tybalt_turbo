<script lang="ts">
  import DSHelp from "$lib/components/DSHelp.svelte";
  import { resolve } from "$app/paths";
  import DSList from "$lib/components/DSList.svelte";
  import DSTabBar, { type TabItem } from "$lib/components/DSTabBar.svelte";
  import DSCallout from "$lib/components/DSCallout.svelte";
  import DSStatusLabel from "$lib/components/DSStatusLabel.svelte";
  import DsLabel from "$lib/components/DsLabel.svelte";
  import { globalStore } from "$lib/stores/global";
  import { jobRequestsEnabled } from "$lib/stores/appConfig";
  import { wholeCAD, canViewJobCreateQueue } from "$lib/jobRequests";
  import { shortDate } from "$lib/utilities";
  let { data } = $props();
  type RequestRow = {
    id: string;
    submitted?: string;
    client_name?: string;
    description?: string;
    project_value: number;
    status?: string;
    parent?: string;
    rejection_reason?: string;
  };
  const titles: Record<string, string> = {
    mine: "My Project Requests",
    all: "All Project Requests",
    approval: "Project Requests Pending My Approval",
    creation: "Job Create Queue",
    abandoned: "Abandoned Project Requests",
  };
  const claims = $derived($globalStore.claims);
  // The lists this person may open. The nav hides them while requests are
  // disabled, so they stay reachable here.
  const tabs = $derived(
    (
      [
        ["mine", "Mine", true],
        ["all", "All", claims.includes("job")],
        ["approval", "Approvals", $jobRequestsEnabled && claims.includes("project_authorization")],
        [
          "creation",
          "Job create queue",
          canViewJobCreateQueue(claims, $jobRequestsEnabled, $globalStore.attentionCounts.counts),
        ],
        ["abandoned", "Abandoned", claims.includes("admin")],
      ] as const
    )
      .filter(([, , visible]) => visible)
      .map(([queue, label]): TabItem => ({
        label,
        href: `?queue=${queue}`,
        active: data.queue === queue,
      })),
  );
</script>

{#snippet anchor(item: RequestRow)}
  <a class="link" href={resolve(`/job-requests/${item.id}`)}
    >{item.submitted ? shortDate(item.submitted, true) : "Draft"}</a
  >
{/snippet}
{#snippet headline(item: RequestRow)}
  {item.client_name || item.description || "Untitled request"}
{/snippet}
{#snippet byline(item: RequestRow)}
  <span class="flex flex-wrap items-center gap-2">
    {#if item.project_value > 0}{wholeCAD(item.project_value)}{/if}
    {#if item.status}<DSStatusLabel status={item.status} kind="request" />{/if}
    {#if item.parent}<DsLabel color="cyan">Subjob</DsLabel>{/if}
  </span>
{/snippet}
{#snippet line1(item: RequestRow)}
  {#if item.client_name}{item.description}{/if}
{/snippet}
{#snippet line2(item: RequestRow)}
  {#if item.status === "Rejected" && item.rejection_reason}
    <DsLabel tone="danger">{item.rejection_reason}</DsLabel>
  {/if}
{/snippet}

{#snippet help()}
  {#if data.queue === "approval"}
    <DSHelp doc="project-requests#approve-a-request" label="Help" />
  {:else if data.queue === "creation"}
    <DSHelp doc="project-requests#create-the-job" label="Help" />
  {:else if data.queue === "mine"}
    <DSHelp doc="project-requests#create-a-request" label="Help" />
  {:else}
    <DSHelp doc="project-requests#request-statuses" label="Help" />
  {/if}
{/snippet}

<div class="space-y-4">
  {#if tabs.length > 1}<DSTabBar {tabs} rightContent={help} />{/if}
  {#if !$jobRequestsEnabled}
    <DSCallout tone="warning">
      Request processing is disabled. Saved requests remain available.
    </DSCallout>
  {/if}
  {#if data.queue === "abandoned"}
    <DSCallout tone="info">
      These requests have no saved edits in the last 60 days. Review them with the creator before
      you delete them.
    </DSCallout>
  {/if}
  {#if data.queue === "creation"}
    <DSCallout tone="info">
      Review each approved request and its PDF. Check existing jobs and pending requests before you
      create the job.
    </DSCallout>
  {/if}

  <DSList
    items={data.items}
    search={true}
    pagination={true}
    inListHeader={titles[data.queue] || titles.mine}
    emptyMessage="No requests in this list."
    {anchor}
    {headline}
    {byline}
    {line1}
    {line2}
  >
    {#snippet searchBarExtra()}
      {#if tabs.length <= 1}{@render help()}{/if}
    {/snippet}
  </DSList>

  {#if data.approvedItems.length > 0}
    <DSList
      items={data.approvedItems}
      inListHeader="Approved By Me, Awaiting Job Creation"
      {anchor}
      {headline}
      {byline}
      {line1}
      {line2}
    />
  {/if}
</div>
