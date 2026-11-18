<script lang="ts">
  import DSHelp from "$lib/components/DSHelp.svelte";
  import DSList from "$lib/components/DSList.svelte";
  import DSSection from "$lib/components/DSSection.svelte";
  import DsLabel from "$lib/components/DsLabel.svelte";
  import RateSheetLabel from "$lib/components/rateSheets/RateSheetLabel.svelte";
  import { rateSheetUsage, type WorkspaceRateSheet } from "$lib/clientWorkspace";
  import { rateSheetLabel } from "$lib/rateSheetLabel";
  import { shortDate } from "$lib/utilities";

  // The rate sheet revisions that the client's jobs use now or used in the
  // past, sheets with current jobs first. Each row links to its revision.
  let { sheets }: { sheets: WorkspaceRateSheet[] } = $props();
</script>

<DSSection title={`Rate sheets (${sheets.length})`}>
  {#snippet actions()}
    <DSHelp doc="rate-sheet-periods#see-the-rate-sheets-of-a-client" iconOnly />
  {/snippet}
  <DSList
    items={sheets}
    search={sheets.length > 0}
    searchText={(sheet) => [rateSheetLabel(sheet), sheet.effective_date].join(" ")}
    searchPlaceholder="Filter rate sheets"
    pagination={true}
    emptyMessage={sheets.length
      ? "No rate sheets match this filter."
      : "No jobs of this client use a rate sheet."}
  >
    {#snippet headline(sheet)}
      <RateSheetLabel {sheet} />
    {/snippet}
    {#snippet byline(sheet)}
      {shortDate(sheet.effective_date, true)}
    {/snippet}
    {#snippet line1(sheet)}
      {rateSheetUsage(sheet)}
    {/snippet}
    {#snippet actions(sheet)}
      <DsLabel size="sm" tone={sheet.active ? "success" : "neutral"}>
        {sheet.active ? "Active" : "Inactive"}
      </DsLabel>
    {/snippet}
  </DSList>
</DSSection>
