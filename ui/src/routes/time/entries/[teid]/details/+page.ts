import { pb } from "$lib/pocketbase";
import type { TimeEntriesResponse } from "$lib/pocketbase-types";
import type { PageLoad } from "./$types";
import { error } from "@sveltejs/kit";

export const load: PageLoad = async ({ params }) => {
  try {
    const item = await pb.collection("time_entries").getOne<TimeEntriesResponse>(params.teid, {
      expand:
        "job,time_type,division,category,role,work_record_subject_id,work_record_subject_id.work_record",
    });

    return { item };
  } catch (err) {
    console.error(`loading time entry details: ${err}`);
    throw error(404, "Failed to load time entry details");
  }
};
