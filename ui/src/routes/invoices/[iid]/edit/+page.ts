import { pb } from "$lib/pocketbase";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ params }) => {
  try {
    const invoice = await pb.collection("invoices").getOne(params.iid);
    return { invoice };
  } catch (err) {
    throw error(404, `Failed to load the invoice: ${err}`);
  }
};
