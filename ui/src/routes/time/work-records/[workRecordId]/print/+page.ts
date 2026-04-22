import { pb } from "$lib/pocketbase";
import type { WorkRecordDetailsResponse } from "$lib/svelte-types";
import type { PageLoad } from "./$types";
import { error } from "@sveltejs/kit";

export const load: PageLoad = async ({ params }) => {
  try {
    const item = (await pb.send(`/api/work_records/${encodeURIComponent(params.workRecordId)}`, {
      method: "GET",
    })) as WorkRecordDetailsResponse;

    return { item };
  } catch (err) {
    console.error(`loading printable work record: ${err}`);
    throw error(404, `Failed to load printable work record: ${err}`);
  }
};
