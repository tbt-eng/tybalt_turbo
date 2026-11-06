import type { PurchaseOrdersResponse } from "$lib/pocketbase-types";
import type { PageLoad } from "./$types";
import { fetchVisiblePOs, isPORevision } from "$lib/poVisibility";
export const load: PageLoad = async () => {
  try {
    const result = await fetchVisiblePOs("active");
    return {
      items: result,
      // Matches the active scope, so realtime events for purchase orders the
      // caller can see but that aren't Active or under revision stay off the list.
      createdItemIsVisible: (record: PurchaseOrdersResponse) =>
        record.status === "Active" || isPORevision(record),
      realtime_source: "visible" as const,
    };
  } catch (error) {
    console.error(`loading data: ${error}`);
    return { items: [], realtime_source: "visible" as const };
  }
};
