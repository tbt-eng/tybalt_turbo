import { pb } from "$lib/pocketbase";
import { error } from "@sveltejs/kit";
import type { InvoiceWithClient } from "$lib/invoices";
import type { PageLoad } from "./$types";

export const load: PageLoad = async () => {
  try {
    const invoices = await pb.collection("invoices").getFullList<InvoiceWithClient>({
      sort: "-created",
      expand: "client",
    });
    return { invoices };
  } catch (err) {
    throw error(403, `Failed to load invoices: ${err}`);
  }
};
