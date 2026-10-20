import { pb } from "$lib/pocketbase";
import { error } from "@sveltejs/kit";
import { loadInvoiceLines, type InvoiceWithClient } from "$lib/invoices";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ params }) => {
  try {
    const [invoice, lines] = await Promise.all([
      pb.collection("invoices").getOne<InvoiceWithClient>(params.iid, { expand: "client" }),
      loadInvoiceLines(params.iid),
    ]);
    return { invoice, lines };
  } catch (err) {
    throw error(404, `Failed to load the invoice: ${err}`);
  }
};
