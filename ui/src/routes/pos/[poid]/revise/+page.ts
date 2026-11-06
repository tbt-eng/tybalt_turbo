import { PurchaseOrdersStatusOptions } from "$lib/pocketbase-types";
import { pb } from "$lib/pocketbase";
import type { PageLoad } from "./$types";
import type { PurchaseOrdersPageData } from "$lib/svelte-types";
import { redirect } from "@sveltejs/kit";

// Revising sends an Active purchase order back for approval with a new total.
// The server checks the remaining conditions (creator, no expenses, no child
// POs) when the revision is submitted.
export const load: PageLoad<PurchaseOrdersPageData> = async ({ params }) => {
  const item = await pb.collection("purchase_orders").getOne(params.poid);
  if (
    item.legacy_manual_entry ||
    item.status !== PurchaseOrdersStatusOptions.Active ||
    item.uid !== pb.authStore.record?.id
  ) {
    throw redirect(303, `/pos/${params.poid}/details`);
  }
  return { item, editing: true, id: params.poid };
};
