import type { PurchaseOrdersResponse } from "$lib/pocketbase-types";
import type { PageLoad } from "./$types";
import { fetchVisiblePOs } from "$lib/poVisibility";
export const load: PageLoad = async () => {
  // 30 days ago as YYYY-MM-DD
  const staleDate = new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString().replace("T", " ");
  try {
    const result = await fetchVisiblePOs("stale", staleDate);
    return {
      items: result,
      // Matches the stale scope, so realtime events for other purchase orders
      // the caller can see stay off the list.
      createdItemIsVisible: (record: PurchaseOrdersResponse) =>
        record.status === "Active" &&
        ((record.second_approval !== "" && record.second_approval < staleDate) ||
          record.approved < staleDate),
      realtime_source: "visible" as const,
    };
  } catch (error) {
    console.error(`loading data: ${error}`);
    return { items: [], realtime_source: "visible" as const };
  }
};
