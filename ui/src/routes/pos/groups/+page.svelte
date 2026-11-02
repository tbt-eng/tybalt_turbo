<script lang="ts">
  import { resolve } from "$app/paths";
  import DSCallout from "$lib/components/DSCallout.svelte";
  import DSList from "$lib/components/DSList.svelte";
  import DsLabel from "$lib/components/DsLabel.svelte";
  import { poGroupTitle, type PoGroupSummary } from "$lib/poGroups";
  import { formatCurrencyAmount, shortDate } from "$lib/utilities";
  import type { PageData } from "./$types";

  let { data }: { data: PageData } = $props();
</script>

{#if data.loadError}
  <DSCallout tone="danger">{data.loadError}</DSCallout>
{/if}

<DSList
  items={data.groups}
  search={true}
  searchText={(group: PoGroupSummary) =>
    [poGroupTitle(group), group.vendor_name, group.creator_name].join(" ")}
  inListHeader="PO Groups"
  emptyMessage="No purchase order groups yet."
>
  {#snippet headline(group)}
    <a class="link" href={resolve(`/pos/groups/${group.id}`)}>{poGroupTitle(group)}</a>
  {/snippet}
  {#snippet byline(group)}
    {group.count} purchase orders, {formatCurrencyAmount(group.total, group.currency_code)}
  {/snippet}
  {#snippet line1(group)}
    {[`Created ${shortDate(group.created, true)}`, group.creator_name]
      .filter((part) => part)
      .join(" by ")}
  {/snippet}
  {#snippet line2(group)}
    {#if group.unapproved_count > 0}
      <DsLabel tone="warning" size="sm">{group.unapproved_count} awaiting approval</DsLabel>
    {:else}
      <DsLabel tone="success" size="sm">All approved</DsLabel>
    {/if}
  {/snippet}
</DSList>
