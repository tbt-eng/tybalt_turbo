// Absorb suggestions: buckets of probable duplicate clients or vendors, best
// first, for the users who absorb them. The server finds them on each request
// (app/routes/absorb_suggestions.go, matching in app/absorb/suggest).

import { resolve } from "$app/paths";
import { pb } from "$lib/pocketbase";

export type AbsorbSuggestionCollection = "clients" | "vendors";
export type AbsorbSuggestionLevel = "strict" | "normal" | "loose";

export const absorbSuggestionCollections: AbsorbSuggestionCollection[] = ["clients", "vendors"];

export const absorbSuggestionLevels: { id: AbsorbSuggestionLevel; label: string }[] = [
  { id: "strict", label: "Strict" },
  { id: "normal", label: "Normal" },
  { id: "loose", label: "Loose" },
];

export type AbsorbSuggestionMember = {
  id: string;
  name: string;
  alias: string;
  // Jobs for a client (as client or owner); POs and expenses for a vendor.
  references: number;
  // The match against the suggested target, from 0 to 1 (1 for the target).
  score: number;
  // An open absorb proposal that names the record, or "".
  proposal: string;
};

export type AbsorbSuggestionBucket = {
  key: string;
  target_id: string;
  confidence: number;
  moved_references: number;
  rank: number;
  members: AbsorbSuggestionMember[];
};

// The collection's absorb action that waits for a commit or an undo. The
// server doesn't record who made it.
export type AbsorbSuggestionPendingAction = {
  id: string;
  target_id: string;
  target_name: string;
  absorbed: string[];
  created: string;
};

export type AbsorbSuggestionsResponse = {
  collection: AbsorbSuggestionCollection;
  level: AbsorbSuggestionLevel;
  buckets: AbsorbSuggestionBucket[];
  // Buckets hidden because someone marked them "Not duplicates".
  dismissed: number;
  pending_action: AbsorbSuggestionPendingAction | null;
};

export function isAbsorbSuggestionCollection(value: unknown): value is AbsorbSuggestionCollection {
  return value === "clients" || value === "vendors";
}

// The suggestions page on a collection's tab. The tab is in the URL, so a
// reload keeps it.
export function absorbSuggestionsHref(collection: AbsorbSuggestionCollection): string {
  return `${resolve("/absorb/suggestions")}?tab=${collection}`;
}

export async function fetchAbsorbSuggestions(
  collection: AbsorbSuggestionCollection,
  level: AbsorbSuggestionLevel,
): Promise<AbsorbSuggestionsResponse> {
  return (await pb.send(`/api/absorb_suggestions/${collection}?level=${level}`, {
    method: "GET",
  })) as AbsorbSuggestionsResponse;
}

// Marks the records "Not duplicates" for everyone. Their bucket stays hidden
// until a new record joins it.
export async function dismissAbsorbSuggestion(
  collection: AbsorbSuggestionCollection,
  recordIds: string[],
): Promise<void> {
  await pb.send(`/api/absorb_suggestions/${collection}/dismiss`, {
    method: "POST",
    body: { record_ids: recordIds },
  });
}

// Absorbs the records into the target with the normal absorb route.
export async function absorbIntoTarget(
  collection: AbsorbSuggestionCollection,
  targetId: string,
  recordIds: string[],
): Promise<void> {
  await pb.send(`/api/${collection}/${targetId}/absorb`, {
    method: "POST",
    body: { ids_to_absorb: recordIds },
  });
}

// The page of a record, to check it in a new tab.
export function absorbSuggestionRecordHref(
  collection: AbsorbSuggestionCollection,
  id: string,
): string {
  return `/${collection}/${id}/details`;
}
