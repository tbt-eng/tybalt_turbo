<script lang="ts">
  import DSSection from "$lib/components/DSSection.svelte";
  import DSFieldList from "$lib/components/DSFieldList.svelte";
  import DSField from "$lib/components/DSField.svelte";
  import DSHelp from "$lib/components/DSHelp.svelte";
  import RateSheetRatesTable, { type ProposedRate } from "./RateSheetRatesTable.svelte";
  import type { WorkflowRecord } from "$lib/clientEditors";

  // A project request's proposed rates, read only, as the approver and
  // procurement review them: the reason, and each rate beside the rate of the
  // sheet it started from, with changes marked. related is the request
  // detail's related records.
  let { related }: { related: WorkflowRecord } = $props();

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
</script>

<DSSection title="Proposed rates">
  {#snippet actions()}
    <DSHelp doc="proposed-rates#approve-proposed-rates" label="Help" />
  {/snippet}
  <DSFieldList>
    <DSField label="Name" value={related.rate_sheet?.name} />
    <DSField label="Started from" value={sourceLabel} />
    <DSField label="Reason" value={related.rate_sheet?.reason} />
  </DSFieldList>
  <RateSheetRatesTable {rows} sourceName={source?.name ?? ""} />
</DSSection>
