import { pb } from "$lib/pocketbase";
import type { PageLoad } from "./$types";
import { error } from "@sveltejs/kit";

export const load: PageLoad = async ({ params }) => {
  const rateSheetId = params.id;

  try {
    // Fetch rate sheet record
    const rateSheet = await pb.collection("rate_sheets").getOne(rateSheetId);

    // Fetch all rate_sheet_entries for this sheet (expanded with role relation)
    const entries = await pb.collection("rate_sheet_entries").getFullList({
      filter: `rate_sheet="${rateSheetId}"`,
      expand: "role",
    });

    // Earlier names of the sheet, newest first
    const { former_names: formerNames } = await pb.send<{ former_names: string[] }>(
      `/api/rate_sheets/${rateSheetId}/former_names`,
      { method: "GET" },
    );

    return {
      rateSheet,
      entries,
      formerNames,
    };
  } catch (err) {
    console.error(`loading rate sheet for print: ${err}`);
    throw error(404, `Failed to load rate sheet: ${err}`);
  }
};
