import type { PageLoad } from "./$types";
import { pb } from "$lib/pocketbase";
import type { AbsorbProposalRow } from "$lib/absorbProposals";

export const load: PageLoad = async ({ url }) => {
  const status = url.searchParams.get("status") === "decided" ? "decided" : "pending";
  const items = (await pb.send(`/api/absorb_proposals?status=${status}`, {
    method: "GET",
  })) as AbsorbProposalRow[];
  return { status, items } as const;
};
