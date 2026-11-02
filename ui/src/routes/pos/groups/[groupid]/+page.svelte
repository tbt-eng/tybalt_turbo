<script lang="ts">
  import DSPageHeader from "$lib/components/DSPageHeader.svelte";
  import DsLabel from "$lib/components/DsLabel.svelte";
  import PurchaseOrdersList from "$lib/components/PurchaseOrdersList.svelte";
  import { fetchPoGroup, poGroupTitle } from "$lib/poGroups";
  import { fetchGroupPOs } from "$lib/poVisibility";
  import { formatCurrencyAmount, shortDate } from "$lib/utilities";
  import type { PageData } from "./$types";

  let { data }: { data: PageData } = $props();

  // Refreshed with the list whenever one of the group's purchase orders
  // changes, so the counts stay current.
  let group = $derived(data.group);

  const createdBy = $derived(
    [`created ${shortDate(group.created, true)}`, group.creator_name]
      .filter((part) => part)
      .join(" by "),
  );

  async function refreshGroup() {
    const [summary, items] = await Promise.all([
      fetchPoGroup(data.group.id),
      fetchGroupPOs(data.group.id),
    ]);
    group = summary;
    return items;
  }
</script>

<div class="space-y-2 p-4">
  <DSPageHeader title={poGroupTitle(group)}>
    {#snippet badges()}
      {#if group.unapproved_count > 0}
        <DsLabel tone="warning">{group.unapproved_count} awaiting approval</DsLabel>
      {:else}
        <DsLabel tone="success">All approved</DsLabel>
      {/if}
    {/snippet}
    {#snippet context()}
      {group.count} purchase orders, {formatCurrencyAmount(group.total, group.currency_code)} · {createdBy}
    {/snippet}
  </DSPageHeader>
</div>

<PurchaseOrdersList
  {data}
  showOwner={true}
  showRemaining={true}
  refreshItems={refreshGroup}
  shouldRefreshOnEvent={(record) => record.po_group === data.group.id}
/>
