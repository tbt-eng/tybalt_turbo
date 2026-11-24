<script lang="ts">
  import DSHelp from "$lib/components/DSHelp.svelte";
  import DsActionButton from "$lib/components/DSActionButton.svelte";
  import PurchaseOrdersList from "$lib/components/PurchaseOrdersList.svelte";
  import { fetchVisiblePOs, type VisiblePurchaseOrderResponse } from "$lib/poVisibility";
  import type { PurchaseOrdersResponse } from "$lib/pocketbase-types";
  import { onMount } from "svelte";
  import { EXPIRING_WITHIN_DAYS, expiringCutoff } from "$lib/poAge";

  interface PurchaseOrdersListData {
    items?: VisiblePurchaseOrderResponse[];
    createdItemIsVisible?: (record: PurchaseOrdersResponse) => boolean;
    realtime_source?: "visible" | "pending";
  }

  let days = $state(EXPIRING_WITHIN_DAYS);
  let loading = $state(true);
  let errorMessage = $state<string | null>(null);
  let expiringBefore = $state("");
  let listKey = $state(0);
  let listData = $state<PurchaseOrdersListData>({
    items: undefined,
    realtime_source: "visible",
  });

  async function load(): Promise<void> {
    const safeDays = Number.isFinite(days) && days > 0 ? Math.floor(days) : EXPIRING_WITHIN_DAYS;
    days = safeDays;
    expiringBefore = expiringCutoff(safeDays);

    loading = true;
    errorMessage = null;
    try {
      const before = expiringBefore;
      const result = await fetchVisiblePOs("expiring", before);
      listData = {
        items: result,
        // Matches the expiring scope, so realtime events for other purchase
        // orders the caller can see stay off the list.
        createdItemIsVisible: (record) =>
          record.status === "Active" &&
          record.type === "Recurring" &&
          record.end_date !== "" &&
          record.end_date <= before,
        realtime_source: "visible",
      };
      listKey += 1;
    } catch (e: any) {
      errorMessage = e?.response?.message ?? "Failed to load expiring purchase orders";
      listData = {
        items: [],
        realtime_source: "visible",
      };
      listKey += 1;
    } finally {
      loading = false;
    }
  }

  onMount(() => {
    void load();
  });
</script>

<div class="flex items-center gap-x-2 bg-neutral-200 p-2">
  <input
    type="number"
    min="1"
    bind:value={days}
    class="w-24 rounded-sm border border-neutral-300 px-2 py-1 text-base"
    title="Days ahead"
    onkeydown={(e) => e.key === "Enter" && load()}
  />
  <DsActionButton action={load} icon="mdi:magnify" title="Load" color="yellow" />
</div>

{#if loading}
  <div class="p-4">Loading…</div>
{:else if errorMessage}
  <div class="p-4 text-red-600">{errorMessage}</div>
{:else}
  {#key listKey}
    <PurchaseOrdersList
      inListHeader={`Expiring Recurring Purchase Orders (end date on or before ${expiringBefore})`}
      data={listData}
      showRemaining={true}
      ageLabels={false}
    >
      {#snippet searchBarExtra()}
        <DSHelp doc="find-a-purchase-order#recurring-pos-that-end-soon" label="Help" />
      {/snippet}
    </PurchaseOrdersList>
  {/key}
{/if}
