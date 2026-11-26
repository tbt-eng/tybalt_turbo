<script lang="ts">
  import { resolve } from "$app/paths";
  import DSTabBar from "$lib/components/DSTabBar.svelte";
  import DsList from "$lib/components/DSList.svelte";
  import DSModal from "$lib/components/DSModal.svelte";
  import DSCallout from "$lib/components/DSCallout.svelte";
  import DSCodeBlock from "$lib/components/DSCodeBlock.svelte";
  import DSFieldList from "$lib/components/DSFieldList.svelte";
  import DSField from "$lib/components/DSField.svelte";
  import DSHelp from "$lib/components/DSHelp.svelte";
  import DsActionButton from "$lib/components/DSActionButton.svelte";
  import { getApiErrorMessage } from "$lib/errors";
  import { shortDate } from "$lib/utilities";
  import {
    absorbCollectionLabel,
    fetchAbsorbHistoryDetail,
    type AbsorbHistoryDetail,
    type AbsorbHistoryRow,
  } from "$lib/absorbHistory";
  import type { PageData } from "./$types";

  let { data }: { data: PageData } = $props();

  let searchTerm = $state("");
  let showDetail = $state(false);
  let detail = $state<AbsorbHistoryDetail | null>(null);
  let detailError = $state("");

  function day(value: string) {
    return value ? shortDate(value.slice(0, 10), true) : "";
  }

  // "by Ann on 2026 Nov 25", with either part left out when it is unknown.
  function byOn(name: string, date: string) {
    return [name ? `by ${name}` : "", day(date) ? `on ${day(date)}` : ""].filter(Boolean).join(" ");
  }

  function kept(row: AbsorbHistoryRow) {
    return row.target_name || row.target_id;
  }

  async function open(row: AbsorbHistoryRow) {
    detail = null;
    detailError = "";
    showDetail = true;
    try {
      detail = await fetchAbsorbHistoryDetail(row.id);
    } catch (error) {
      detailError = getApiErrorMessage(error, "The absorb could not be loaded");
    }
  }
</script>

<svelte:head>
  <title>Absorb History</title>
</svelte:head>

<div>
  <DSTabBar
    tabs={[
      { label: "Pending", href: resolve("/absorb/actions"), active: false },
      { label: "History", href: resolve("/absorb/history"), active: true },
    ]}
  >
    {#snippet rightContent()}
      <DSHelp doc="absorb-history#find-an-absorb" label="Help" />
    {/snippet}
  </DSTabBar>
  <DsList
    items={data.items}
    search={true}
    pagination={true}
    bind:searchTerm
    inListHeader="Committed Absorbs"
    searchText={(row) =>
      [
        absorbCollectionLabel(row.collection_name),
        kept(row),
        row.absorbed.join(" "),
        row.absorbed_by_name,
        row.committed_by_name,
      ].join(" ")}
    emptyMessage={searchTerm.trim() !== ""
      ? "No matching absorbs."
      : "No absorb was committed in the last 90 days."}
  >
    {#snippet anchor(row: AbsorbHistoryRow)}
      {absorbCollectionLabel(row.collection_name)}
    {/snippet}
    {#snippet headline(row: AbsorbHistoryRow)}
      {kept(row)}
    {/snippet}
    {#snippet line1(row: AbsorbHistoryRow)}
      Absorbed: {row.absorbed.join(", ")}
    {/snippet}
    {#snippet line2(row: AbsorbHistoryRow)}
      {row.references_moved}
      {row.references_moved === 1 ? "reference" : "references"} moved
    {/snippet}
    {#snippet line3(row: AbsorbHistoryRow)}
      Absorbed {byOn(row.absorbed_by_name, row.absorbed_at)} · Committed {byOn(
        row.committed_by_name,
        row.committed_at,
      )}
    {/snippet}
    {#snippet actions(row: AbsorbHistoryRow)}
      <DsActionButton action={() => open(row)} variant="secondary">Details</DsActionButton>
    {/snippet}
  </DsList>
</div>

<DSModal bind:show={showDetail} title="Committed absorb" size="lg">
  {#if detailError}
    <DSCallout tone="danger">{detailError}</DSCallout>
  {:else if !detail}
    <p>Loading...</p>
  {:else}
    <div class="flex flex-col gap-3">
      <DSFieldList>
        <DSField label="Type" value={absorbCollectionLabel(detail.collection_name)} />
        <DSField label="Kept" value={`${kept(detail)} (${detail.target_id})`} />
        <DSField label="Absorbed" value={byOn(detail.absorbed_by_name, detail.absorbed_at)} />
        <DSField label="Committed" value={byOn(detail.committed_by_name, detail.committed_at)} />
        <DSField label="References moved" value={String(detail.references_moved)} />
      </DSFieldList>
      <DSField label="Absorbed records, as they were">
        <DSCodeBlock
          text={JSON.stringify(detail.absorbed_records, null, 2)}
          label="Copy the absorbed records"
        />
      </DSField>
      <DSField label="References moved, by table, field and row (the value is the absorbed record)">
        <DSCodeBlock
          text={JSON.stringify(detail.updated_references, null, 2)}
          label="Copy the references moved"
        />
      </DSField>
    </div>
  {/if}
</DSModal>
