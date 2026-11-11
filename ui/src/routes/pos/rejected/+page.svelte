<script lang="ts">
  import DSHelp from "$lib/components/DSHelp.svelte";
  import type { PageData } from "./$types";
  import PurchaseOrdersList from "$lib/components/PurchaseOrdersList.svelte";
  let { data }: { data: PageData } = $props();
</script>

<PurchaseOrdersList inListHeader="Rejected Purchase Orders (Mine)" {data}>
  {#snippet searchBarExtra()}
    <DSHelp doc="find-a-purchase-order#fix-a-rejected-po" label="Help" />
  {/snippet}
</PurchaseOrdersList>
