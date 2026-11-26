// Absorb history: the committed absorbs, with who absorbed and who committed.
// Turbo keeps each for 90 days (app/routes/absorb_history.go).

import { pb } from "$lib/pocketbase";

export type AbsorbHistoryRow = {
  id: string;
  collection_name: string;
  target_id: string;
  // The name of the kept record when the absorb was committed.
  target_name: string;
  // The names of the absorbed records.
  absorbed: string[];
  references_moved: number;
  proposal: string;
  absorbed_by: string;
  absorbed_by_name: string;
  absorbed_at: string;
  committed_by: string;
  committed_by_name: string;
  committed_at: string;
};

// One absorb with the saved copies of the absorbed records and the
// references that moved, for troubleshooting.
export type AbsorbHistoryDetail = AbsorbHistoryRow & {
  absorbed_records: unknown;
  updated_references: unknown;
};

export async function fetchAbsorbHistory(): Promise<AbsorbHistoryRow[]> {
  return (await pb.send("/api/absorb_history", { method: "GET" })) as AbsorbHistoryRow[];
}

export async function fetchAbsorbHistoryDetail(id: string): Promise<AbsorbHistoryDetail> {
  return (await pb.send(`/api/absorb_history/${id}`, { method: "GET" })) as AbsorbHistoryDetail;
}

export function absorbCollectionLabel(name: string): string {
  if (name === "client_contacts") return "Client Contacts";
  if (name === "clients") return "Clients";
  if (name === "vendors") return "Vendors";
  if (name === "rate_sheets") return "Rate Sheets";
  return name;
}
