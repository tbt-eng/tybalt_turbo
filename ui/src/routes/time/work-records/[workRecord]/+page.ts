import { pb } from "$lib/pocketbase";
import type { PageLoad } from "./$types";
import { error } from "@sveltejs/kit";
import type { WorkRecordDetailsResponse } from "$lib/svelte-types";

export const load: PageLoad = async ({ params }) => {
  try {
    const item = (await pb.send(`/api/work_records/${encodeURIComponent(params.workRecord)}`, {
      method: "GET",
    })) as WorkRecordDetailsResponse;

    return {
      item,
    };
  } catch (err) {
    console.error(`loading work record details: ${err}`);
    throw error(404, `Failed to load work record details: ${err}`);
  }
};
