<script lang="ts">
  import { resolve } from "$app/paths";
  import { pb } from "$lib/pocketbase";
  import DSCallout from "$lib/components/DSCallout.svelte";
  import DSCollapsible from "$lib/components/DSCollapsible.svelte";
  import DsLabel from "$lib/components/DsLabel.svelte";
  import DSStatusLabel from "$lib/components/DSStatusLabel.svelte";
  import { countLabel, type WorkspaceUsage, type WorkspaceUsageRow } from "$lib/clientWorkspace";

  // The jobs and requests that use one contact or one invoicing profile of a
  // client. It shows the counts and loads the list when a person opens it.
  let {
    clientId,
    kind,
    recordId,
    jobCount,
    requestCount,
  }: {
    clientId: string;
    kind: "contacts" | "invoicing_profiles";
    recordId: string;
    jobCount: number;
    requestCount: number;
  } = $props();

  const title = $derived(
    [countLabel(jobCount, "job"), requestCount > 0 ? countLabel(requestCount, "request") : ""]
      .filter(Boolean)
      .join(" · "),
  );

  // Loaded on the first open and kept, so closing and opening does not reload.
  let usage: Promise<WorkspaceUsage> | undefined;
  function load() {
    usage ??= pb.send<WorkspaceUsage>(`/api/clients/${clientId}/${kind}/${recordId}/usage`, {
      method: "GET",
      requestKey: null,
    });
    return usage;
  }
</script>

{#snippet roles(row: WorkspaceUsageRow)}
  <!-- An invoicing profile has one role, so only a contact shows them. -->
  {#if kind === "contacts"}
    {#if row.project_contact}<DsLabel size="sm" color="blue">Project contact</DsLabel>{/if}
    {#if row.invoice_recipient}<DsLabel size="sm" color="purple">Invoice recipient</DsLabel>{/if}
  {/if}
{/snippet}

{#if jobCount + requestCount === 0}
  <span class="text-xs text-neutral-600">{title}</span>
{:else}
  <DSCollapsible {title} subtle>
    {#await load()}
      <p class="text-xs text-neutral-600">Loading…</p>
    {:then result}
      <ul class="space-y-1 text-sm" aria-label="Jobs and requests">
        {#each result.jobs as job (job.id)}
          <li class="flex items-center gap-2">
            <a class="link shrink-0" href={resolve("/jobs/[jid]/details", { jid: job.id })}
              >{job.number}</a
            >
            <span class="min-w-0 flex-1 truncate" title={job.description}>{job.description}</span>
            {#if job.number.startsWith("P")}<DsLabel size="sm">Proposal</DsLabel>{/if}
            <DSStatusLabel status={job.status} kind="job" size="sm" />
            {@render roles(job)}
          </li>
        {/each}
        {#each result.requests as request (request.id)}
          <li class="flex items-center gap-2">
            <a class="link shrink-0" href={resolve("/job-requests/[id]", { id: request.id })}
              >Project request</a
            >
            <span class="min-w-0 flex-1 truncate" title={request.description}
              >{request.description}</span
            >
            <DSStatusLabel status={request.status} kind="request" size="sm" />
            {@render roles(request)}
          </li>
        {/each}
      </ul>
      {#if result.hidden_requests > 0}
        <p class="text-xs text-neutral-600">
          {countLabel(result.hidden_requests, "request")} that you cannot open
          {result.hidden_requests === 1 ? "is" : "are"} not listed.
        </p>
      {/if}
    {:catch}
      <DSCallout tone="danger">The list did not load. Close it and open it again.</DSCallout>
    {/await}
  </DSCollapsible>
{/if}
