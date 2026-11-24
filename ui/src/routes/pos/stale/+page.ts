import type { PurchaseOrdersResponse } from "$lib/pocketbase-types";
import type { PageLoad } from "./$types";
import { fetchVisiblePOs } from "$lib/poVisibility";
import { staleCutoff } from "$lib/poAge";
import { error } from "@sveltejs/kit";
import { get } from "svelte/store";
import { globalStore } from "$lib/stores/global";
export const load: PageLoad = async () => {
  await globalStore.ensureClaims();
  if (!get(globalStore).claims.includes("payables_admin"))
    error(403, "Only payables staff can list stale purchase orders.");
  const staleDate = staleCutoff();
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
