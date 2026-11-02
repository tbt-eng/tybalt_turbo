import type { PageLoad } from "./$types";
import { error } from "@sveltejs/kit";
import { fetchGroupPOs } from "$lib/poVisibility";
import { fetchPoGroup, type PoGroupApiError } from "$lib/poGroups";

export const load: PageLoad = async ({ params }) => {
  try {
    const [group, items] = await Promise.all([
      fetchPoGroup(params.groupid),
      fetchGroupPOs(params.groupid),
    ]);
    // The page refreshes the group and its list itself when one of its
    // purchase orders changes.
    return { group, items, realtime_source: "none" as const };
  } catch (caught) {
    const loadError = caught as PoGroupApiError;
    throw error(
      loadError.status === 404 ? 404 : 500,
      loadError.response?.message ?? "Couldn't load the purchase order group.",
    );
  }
};
