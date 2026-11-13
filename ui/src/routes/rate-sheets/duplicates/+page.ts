import { pb } from "$lib/pocketbase";
import type { PageLoad } from "./$types";
import { error } from "@sveltejs/kit";
import type { RateSheetDuplicateGroup } from "$lib/rateSheetAbsorb";

export const load: PageLoad = async () => {
  try {
    const { groups } = await pb.send<{ groups: RateSheetDuplicateGroup[] }>(
      "/api/rate_sheets/duplicates",
      { method: "GET" },
    );
    return { groups };
  } catch (err) {
    throw error(500, `Failed to load duplicate rate sheets: ${err}`);
  }
};
