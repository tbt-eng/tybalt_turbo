<script lang="ts">
  import DSList from "$lib/components/DSList.svelte";
  import DSActionButton from "$lib/components/DSActionButton.svelte";
  import DSHelp from "$lib/components/DSHelp.svelte";
  import DsLabel from "$lib/components/DsLabel.svelte";
  import { jobs, type JobApiResponse } from "$lib/stores/jobs";
  import { globalStore } from "$lib/stores/global";
  import { resolve } from "$app/paths";

  // Active projects and subjobs that have no invoicing profile yet. The
  // profile is optional, so procurement fills these in before the project is
  // invoiced. The jobs store already holds every job, so this view filters it.
  jobs.init();

  const items = $derived(
    ($jobs.items as JobApiResponse[]).filter(
      (job) => !job.number?.startsWith("P") && job.status === "Active" && job.no_invoicing_profile,
    ),
  );
  const searchText = (job: JobApiResponse) =>
    [job.number, job.description, job.client, job.client_alias, job.manager, job.branch].join(" ");
</script>

<DSList
  {items}
  search
  {searchText}
  searchPlaceholder="Filter by number, description, client or manager"
  inListHeader="Active projects without an invoicing profile"
  emptyMessage={$jobs.loading ? "Loading projects…" : "Every active project has a profile."}
>
  {#snippet searchBarExtra()}
    <DSHelp doc="clients#find-projects-without-a-profile" label="Help" />
  {/snippet}
  {#snippet anchor(job: JobApiResponse)}
    <a class="link font-bold" href={resolve("/jobs/[jid]/details", { jid: job.id })}>{job.number}</a
    >
  {/snippet}
  {#snippet headline(job: JobApiResponse)}{job.description}{/snippet}
  {#snippet byline(job: JobApiResponse)}{job.client}{/snippet}
  {#snippet line1(job: JobApiResponse)}{#if job.branch}<DsLabel color="neutral"
        >{job.branch}</DsLabel
      >{/if}{#if job.manager}<DsLabel color="purple">{job.manager}</DsLabel>{/if}{/snippet}
  {#snippet actions(job: JobApiResponse)}
    {#if $globalStore.claims.includes("job")}
      <DSActionButton
        action={resolve("/jobs/[jid]/edit", { jid: job.id })}
        icon="mdi:pencil"
        title="Edit"
        color="blue"
      />
    {/if}
  {/snippet}
</DSList>
