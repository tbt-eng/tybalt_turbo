import { pb } from "$lib/pocketbase";
import type { PageLoad } from "./$types";
import { error } from "@sveltejs/kit";
import type { RateSheetJob } from "$lib/components/rateSheets/RateSheetJobs.svelte";

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

    // Fetch all rate_roles (to compute missing roles)
    const allRoles = await pb.collection("rate_roles").getFullList({ sort: "name" });

    // Earlier names and codes of the sheet, newest first
    const { former_names: formerNames, former_codes: formerCodes } = await pb.send<{
      former_names: string[];
      former_codes: string[];
    }>(`/api/rate_sheets/${rateSheetId}/former_names`, { method: "GET" });

    // The jobs of this revision, now or in their rate sheet history
    const jobs = await pb.send<RateSheetJob[]>(`/api/rate_sheets/${rateSheetId}/jobs`, {
      method: "GET",
    });

    return {
      rateSheet,
      entries,
      formerNames,
      formerCodes: formerCodes ?? [],
      allRoles,
      jobs,
    };
  } catch (err) {
    console.error(`loading rate sheet details: ${err}`);
    throw error(404, `Failed to load rate sheet details: ${err}`);
  }
};
