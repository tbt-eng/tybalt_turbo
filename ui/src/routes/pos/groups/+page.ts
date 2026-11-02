import type { PageLoad } from "./$types";
import { fetchPoGroups, type PoGroupApiError, type PoGroupSummary } from "$lib/poGroups";

export const load: PageLoad = async () => {
  try {
    return { groups: await fetchPoGroups(), loadError: "" };
  } catch (caught) {
    const error = caught as PoGroupApiError;
    return {
      groups: [] as PoGroupSummary[],
      loadError: error.response?.message ?? "Couldn't load purchase order groups.",
    };
  }
};
