// Shared branches collection store.
//
// Why this exists:
// Several editors previously loaded branches with ad-hoc `getFullList` calls.
// Centralizing through `createCollectionStore` gives one source of truth for
// initial load, realtime updates, and search index behavior.
import type { BranchesResponse } from "$lib/pocketbase-types";
import { createCollectionStore } from "./collectionStore";
import { pb } from "$lib/pocketbase";

export const branches = createCollectionStore<BranchesResponse>(
  "branches",
  {
    sort: "name",
    requestKey: "branches",
  },
  {
    fields: ["name", "code"],
    storeFields: ["id", "name", "code"],
  },
  async (item) => {
    const fullRecord = await pb.collection("branches").getOne<BranchesResponse>(item.id);
    branches.update((state) => ({
      ...state,
      items: [...state.items, fullRecord],
      index: state.index?.add(fullRecord) || state.index,
    }));
  },
  async (item) => {
    const fullRecord = await pb.collection("branches").getOne<BranchesResponse>(item.id);
    branches.update((state) => ({
      ...state,
      items: state.items.map((i) => (i.id === item.id ? fullRecord : i)),
      index: state.index?.replace(fullRecord) || state.index,
    }));
  },
);

// The branches a user holding claimIds may choose: those with no
// allowed_claims and those allowing one of the claims. The record's current
// branch is always kept so a value already on the record is never hidden.
export function branchesAllowedFor(
  items: BranchesResponse[],
  claimIds: string[],
  currentBranch = "",
): BranchesResponse[] {
  return items.filter(
    (branch) =>
      branch.id === currentBranch ||
      !branch.allowed_claims?.length ||
      branch.allowed_claims.some((claimId) => claimIds.includes(claimId)),
  );
}
