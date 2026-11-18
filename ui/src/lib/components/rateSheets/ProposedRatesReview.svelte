<script lang="ts">
  import DSSection from "$lib/components/DSSection.svelte";
  import DSFieldList from "$lib/components/DSFieldList.svelte";
  import DSField from "$lib/components/DSField.svelte";
  import DSHelp from "$lib/components/DSHelp.svelte";
  import DsCheck from "$lib/components/DsCheck.svelte";
  import { shortDate } from "$lib/utilities";
  import RateSheetRatesTable from "./RateSheetRatesTable.svelte";
  import RateSheetLabel from "./RateSheetLabel.svelte";
  import { rateSheetCodeRevision } from "$lib/rateSheetLabel";
  import DSCallout from "$lib/components/DSCallout.svelte";
  import { type ProposedRate, rateMissing, reviewRows } from "$lib/proposedRates";
  import type { WorkflowRecord } from "$lib/clientEditors";

  // A project request's proposed rates, read only, as the approver and
  // procurement review them: the reason, and each rate beside the rate of the
  // sheet it started from, with changes marked. related is the request
  // detail's related records. With canApprove, the approver states here that
  // they approve the new sheet and its rates; the statement is apart from the
  // checklist. Afterwards the section shows who approved.
  let {
    related,
    canApprove = false,
    approved = $bindable(false),
  }: { related: WorkflowRecord; canApprove?: boolean; approved?: boolean } = $props();

  const source = $derived(related.rate_sheet_source);
  // Every role: the saved rates, and each role with no rate, marked No rate.
  // Submission needs a rate for every role, but a role can be added after
  // submission.
  const rows = $derived<ProposedRate[]>(
    reviewRows(related.rate_sheet_entries ?? [], related.rate_sheet_missing_rates ?? []),
  );
  const missingCount = $derived(rows.filter(rateMissing).length);
  const sheet = $derived(related.rate_sheet);
  const approver = $derived(related.rate_sheet_rates_approver);
  const approvedBy = $derived(
    approver
      ? [approver.name, approvedDate(related.rate_sheet?.rates_approved)]
          .filter(Boolean)
          .join(" · ")
      : "",
  );

  // The server stores the time in UTC; show the viewer's local day.
  function approvedDate(stored: string | undefined): string {
    const time = new Date((stored ?? "").replace(" ", "T"));
    if (Number.isNaN(time.getTime())) return "";
    const day = [time.getFullYear(), time.getMonth() + 1, time.getDate()]
      .map((part) => String(part).padStart(2, "0"))
      .join("-");
    return shortDate(day, true);
  }
</script>

<DSSection title="Proposed rates">
  {#snippet actions()}
    <DSHelp doc="proposed-rates#approve-proposed-rates" label="Help" />
  {/snippet}
  <DSFieldList>
    {#if sheet?.code}
      <!-- Proposed rates have no page until the job is created. -->
      <DSField label="Rate sheet"><RateSheetLabel {sheet} link={false} /></DSField>
    {/if}
    {#if source}
      <DSField label="Started from"><RateSheetLabel sheet={source} /></DSField>
    {/if}
    <DSField label="Reason" value={related.rate_sheet?.reason} />
    <DSField label="Rates approved by" value={approvedBy} />
  </DSFieldList>
  {#if missingCount > 0}
    <DSCallout
      tone="warning"
      title={`${missingCount} ${missingCount === 1 ? "role has" : "roles have"} no rate`}
    >
      Approval and job creation fail until every role has a rate. After a recall or a rejection, the
      requester can enter the rates.
    </DSCallout>
  {/if}
  <RateSheetRatesTable {rows} sourceName={source ? rateSheetCodeRevision(source) : ""} />
  {#if canApprove}
    <DsCheck
      bind:value={approved}
      errors={{}}
      fieldName="rates_approved"
      uiName="I approve the creation of this new rate sheet and its rates."
    />
  {/if}
</DSSection>
