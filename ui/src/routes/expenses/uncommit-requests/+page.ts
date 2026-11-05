import type { PageLoad } from "./$types";
import { pb } from "$lib/pocketbase";
import type { ExpenseUncommitRequestRow } from "$lib/svelte-types";

export const load: PageLoad = async ({ url }) => {
  const status = url.searchParams.get("status") === "decided" ? "decided" : "pending";
  const items = (await pb.send(`/api/expense_uncommit_requests?status=${status}`, {
    method: "GET",
  })) as ExpenseUncommitRequestRow[];
  return { status, items };
};
