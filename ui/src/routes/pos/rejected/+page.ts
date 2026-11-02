import { pb } from "$lib/pocketbase";
import type { PurchaseOrdersResponse } from "$lib/pocketbase-types";
import type { PageLoad } from "./$types";
import { fetchVisiblePOs } from "$lib/poVisibility";

export const load: PageLoad = async () => {
  try {
    const currentUserId = pb.authStore.record?.id ?? "";
    const result = await fetchVisiblePOs("rejected");
    return {
      items: result,
      // Matches the rejected scope, so realtime events for other purchase
      // orders the caller can see stay off the list.
      createdItemIsVisible: (record: PurchaseOrdersResponse) =>
        record.status === "Unapproved" && record.rejected !== "" && record.uid === currentUserId,
      realtime_source: "visible" as const,
    };
  } catch (error) {
    console.error(`loading data: ${error}`);
    return { items: [], realtime_source: "visible" as const };
  }
};
