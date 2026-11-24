import type { PageLoad } from "./$types";
import { loadClientMerge } from "$lib/clientEditorLoad";
import { loadAbsorbProposal } from "$lib/absorbProposals";

// ?proposal=<id> opens the editor to review that absorb proposal.
export const load: PageLoad = async ({ params, url }) => ({
  ...(await loadClientMerge(params.cid, params.kid, url)),
  proposal: await loadAbsorbProposal(url),
});
