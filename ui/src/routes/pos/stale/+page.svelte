<script lang="ts">
  import type { PageData } from "./$types";
  import PurchaseOrdersList from "$lib/components/PurchaseOrdersList.svelte";
  import { localIsoDate } from "$lib/localDate";
  let { data }: { data: PageData } = $props();
  const staleDate = localIsoDate(new Date(Date.now() - 30 * 24 * 60 * 60 * 1000));
</script>

<PurchaseOrdersList
  inListHeader={`Stale Purchase Orders (Active before ${staleDate})`}
  {data}
  showRemaining={true}
/>
