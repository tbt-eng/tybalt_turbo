<script lang="ts">
  import DSSection from "$lib/components/DSSection.svelte";
  import DSFieldList from "$lib/components/DSFieldList.svelte";
  import DSField from "$lib/components/DSField.svelte";
  import DSHelp from "$lib/components/DSHelp.svelte";
  import DsCheck from "$lib/components/DsCheck.svelte";
  import { shortDate } from "$lib/utilities";
  import RateSheetRatesTable from "./RateSheetRatesTable.svelte";
  import type { ProposedRate } from "$lib/proposedRates";
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
  const rows = $derived<ProposedRate[]>(
    (related.rate_sheet_entries ?? [])
      .map((entry: WorkflowRecord) => ({
        role: entry.role_id,
        roleName: entry.role_name,
        rate: entry.rate,
        overtime_rate: entry.overtime_rate,
        source_rate: entry.source_rate,
        source_overtime_rate: entry.source_overtime_rate,
      }))
      .sort((a: ProposedRate, b: ProposedRate) => a.roleName.localeCompare(b.roleName)),
  );
  const sourceLabel = $derived(source ? `${source.name} (rev. ${source.revision})` : "");
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
    <DSField label="Name" value={related.rate_sheet?.name} />
    <DSField label="Started from" value={sourceLabel} />
    <DSField label="Reason" value={related.rate_sheet?.reason} />
    <DSField label="Rates approved by" value={approvedBy} />
  </DSFieldList>
  <RateSheetRatesTable {rows} sourceName={source?.name ?? ""} />
  {#if canApprove}
    <DsCheck
      bind:value={approved}
      errors={{}}
      fieldName="rates_approved"
      uiName="I approve the creation of this new rate sheet and its rates."
    />
  {/if}
</DSSection>
