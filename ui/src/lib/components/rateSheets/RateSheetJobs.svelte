<script lang="ts" module>
  // A job that uses or used one revision of a rate sheet, from
  // GET /api/rate_sheets/{id}/jobs, sorted by client name and job number.
  export interface RateSheetJob {
    id: string;
    number: string;
    description: string;
    client_id: string;
    client_name: string;
    // True when the job uses the revision now, false when only a period of
    // its rate sheet history did.
    current: boolean;
  }
</script>

<script lang="ts">
  import { resolve } from "$app/paths";
  import DSList from "$lib/components/DSList.svelte";
  import DSSection from "$lib/components/DSSection.svelte";
  import DsLabel from "$lib/components/DsLabel.svelte";

  // The jobs and clients of one revision. The list is paged, so a sheet with
  // thousands of jobs shows one page at a time, and its search finds a job or
  // a client among them.
  let { jobs }: { jobs: RateSheetJob[] } = $props();

  const clientCount = $derived(new Set(jobs.map((job) => job.client_id)).size);
</script>

<DSSection title="Jobs and clients">
  {#snippet titleNote()}
    {#if jobs.length > 0}
      {jobs.length} job{jobs.length === 1 ? "" : "s"} · {clientCount} client{clientCount === 1
        ? ""
        : "s"}
    {/if}
  {/snippet}
  <DSList
    items={jobs}
    search={jobs.length > 0}
    searchText={(job) => [job.number, job.client_name, job.description].join(" ")}
    pagination={true}
    stripeKey={(job) => job.client_id}
    emptyMessage="No jobs use this revision, now or in the past."
  >
    {#snippet anchor(job)}
      <a class="link" href={resolve(`/jobs/${job.id}/details`)}>{job.number}</a>
    {/snippet}
    {#snippet headline(job)}
      {job.client_name}
    {/snippet}
    {#snippet byline(job)}
      {job.description}
    {/snippet}
    {#snippet actions(job)}
      <DsLabel size="sm" tone={job.current ? "success" : "neutral"}>
        {job.current ? "Current" : "Past"}
      </DsLabel>
    {/snippet}
  </DSList>
</DSSection>
