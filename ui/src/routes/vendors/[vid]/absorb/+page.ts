import type { PageLoad } from "./$types";
import { loadAbsorbProposal } from "$lib/absorbProposals";

// ?proposal=<id> opens the editor to review that absorb proposal.
export const load: PageLoad = async ({ url }) => ({ proposal: await loadAbsorbProposal(url) });
