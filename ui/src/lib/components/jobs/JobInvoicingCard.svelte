<script lang="ts">
  import DSCard from "$lib/components/DSCard.svelte";
  import DSCallout from "$lib/components/DSCallout.svelte";
  import DSField from "$lib/components/DSField.svelte";
  import DSFieldList from "$lib/components/DSFieldList.svelte";
  import InvoicingInstructions from "$lib/components/clients/InvoicingInstructions.svelte";
  import { globalStore } from "$lib/stores/global";
  import { resolve } from "$app/paths";
  import type { JobDetails } from "./types";

  // Who receives a project's invoices and how, from the client invoicing
  // profile the project uses. The profile is optional. Procurement (the job
  // claim) sees that a project has none, so it can set one before the project
  // is invoiced; for everyone else the empty card is hidden.
  let { job }: { job: JobDetails } = $props();

  const jobHolder = $derived($globalStore.claims.includes("job"));
  const profileHref = $derived(
    job.invoicing_information && job.client?.id
      ? `${resolve("/clients/[cid]/details", { cid: job.client.id })}?show_invoicing=${job.invoicing_information}#invoicing`
      : "",
  );
</script>

{#if job.invoicing_information || jobHolder}
  <DSCard title="Invoicing">
    {#snippet actions()}
      {#if profileHref}
        <!-- eslint-disable-next-line svelte/no-navigation-without-resolve -- query and hash on a client page -->
        <a class="link text-sm" href={profileHref}>Open profile</a>
      {/if}
    {/snippet}
    {#if job.invoicing_information}
      <DSFieldList>
        <DSField label="Profile" value={job.invoicing_profile_name} />
        {#if job.invoice_contact_name || job.invoice_contact_email}
          <DSField label="Send invoices to">
            {job.invoice_contact_name}
            {#if job.invoice_contact_name && job.invoice_contact_email}·{/if}
            {#if job.invoice_contact_email}
              <a class="link" href={`mailto:${job.invoice_contact_email}`}
                >{job.invoice_contact_email}</a
              >
            {/if}
          </DSField>
        {/if}
      </DSFieldList>
      <InvoicingInstructions text={job.invoice_instructions || ""} hideWhenEmpty />
    {:else}
      <DSCallout tone="warning">No invoicing profile. Edit the job to select one.</DSCallout>
    {/if}
  </DSCard>
{/if}
