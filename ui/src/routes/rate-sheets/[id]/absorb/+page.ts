import { pb } from "$lib/pocketbase";
import type { PageLoad } from "./$types";
import { error } from "@sveltejs/kit";
import type { RateSheetAbsorbCandidate, RateSheetAbsorbSheet } from "$lib/rateSheetAbsorb";

// The target and the sheets with identical rates. The server says why the
// target can't absorb a sheet; the page offers only the others.
export const load: PageLoad = async ({ params }) => {
  try {
    const result = await pb.send<{
      target: RateSheetAbsorbSheet;
      candidates: RateSheetAbsorbCandidate[];
    }>(`/api/rate_sheets/${params.id}/absorb_candidates`, { method: "GET" });
    return result;
  } catch (err) {
    throw error(404, `Failed to load rate sheet: ${err}`);
  }
};
