// Absorb proposals: anyone proposes that duplicate clients, contacts or
// vendors be absorbed, and an absorb holder reviews the proposal. The server
// routes are under /api/absorb_proposals (app/routes/absorb_proposals.go).

// The collections whose absorb can be proposed.
export type AbsorbProposalCollection = "clients" | "client_contacts" | "vendors";

export type AbsorbProposalStatus = "pending" | "accepted" | "declined" | "withdrawn" | "obsolete";

export type AbsorbProposalRecord = {
  record_id: string;
  // A copy of the name: a committed absorb deletes the record.
  record_name: string;
};

export type AbsorbProposalRow = {
  id: string;
  collection_name: AbsorbProposalCollection;
  target_id: string;
  target_name: string;
  // For contacts, the client they belong to.
  client_id: string;
  client_name: string;
  requester: string;
  requester_name: string;
  reason: string;
  status: AbsorbProposalStatus;
  decided_at: string;
  decided_by_name: string;
  decision_comment: string;
  created: string;
  // An absorb holder has absorbed the records and not yet committed.
  in_progress: boolean;
  records: AbsorbProposalRecord[];
};

// Whether the user absorbs, or only proposes an absorb.
export function canAbsorbDirectly(claims: string[]): boolean {
  return claims.includes("absorb");
}

export function absorbProposalCollectionLabel(collection: string): string {
  if (collection === "clients") return "Client";
  if (collection === "client_contacts") return "Contact";
  if (collection === "vendors") return "Vendor";
  return collection;
}

export function absorbProposalStatusLabel(status: AbsorbProposalStatus): string {
  return status.charAt(0).toUpperCase() + status.slice(1);
}

// The absorb page of a proposal's target, where a holder reviews it.
export function absorbProposalReviewHref(row: AbsorbProposalRow): string {
  const query = `?proposal=${encodeURIComponent(row.id)}`;
  if (row.collection_name === "client_contacts")
    return `/clients/${row.client_id}/contacts/${row.target_id}/absorb${query}`;
  if (row.collection_name === "vendors") return `/vendors/${row.target_id}/absorb${query}`;
  return `/clients/${row.target_id}/absorb${query}`;
}
