// Duplicate rate sheets and the absorb of one into another, as the API
// returns them (GET /api/rate_sheets/duplicates and
// GET /api/rate_sheets/{id}/absorb_candidates). The server decides which
// sheets can be absorbed and says why not; the UI only shows it.

export interface RateSheetAbsorbSheet {
  id: string;
  code: string;
  name: string;
  revision: number;
  effective_date: string;
  active: boolean;
  job_count: number;
  request_count: number;
  // The number of revisions of the sheet's code.
  revisions: number;
}

export interface RateSheetDuplicate extends RateSheetAbsorbSheet {
  // The group's target: the sheet with the earliest effective date.
  target: boolean;
  // The sheet can absorb at least one other sheet of its group.
  can_absorb: boolean;
  // Why the group's target can't absorb this sheet.
  reason?: string;
  // The number of sheets that an absorb into the target deletes: 1, or
  // every revision of the name.
  absorbs?: number;
}

export interface RateSheetDuplicateGroup {
  sheets: RateSheetDuplicate[];
  // Why no sheet of the group can be absorbed into its target.
  reason?: string;
}

export interface RateSheetAbsorbCandidate extends RateSheetAbsorbSheet {
  // Why the target can't absorb this sheet; absent when it can.
  reason?: string;
  absorbs?: number;
}

// Absorbing rate sheets needs admin, or both absorb and rate_sheet_revise.
// The server checks the same claims.
export function canAbsorbRateSheets(claims: string[]): boolean {
  return (
    claims.includes("admin") || (claims.includes("absorb") && claims.includes("rate_sheet_revise"))
  );
}

// An absorbed sheet as the absorb_actions record stores it: every field of
// the sheet, with its entries. The Absorb Actions page names each one.
export interface AbsorbedRateSheet {
  id: string;
  code: string;
  name: string;
  revision: number;
}

// The sheets in an absorb action's absorbed_records, by code and revision.
// Anything that is not a stored sheet is left out.
export function absorbedRateSheets(records: unknown): AbsorbedRateSheet[] {
  if (!Array.isArray(records)) return [];
  return records
    .filter(
      (record): record is Record<string, unknown> =>
        !!record && typeof record === "object" && typeof record.id === "string",
    )
    .map((record) => ({
      id: record.id as string,
      code: typeof record.code === "string" ? record.code : "",
      name: typeof record.name === "string" ? record.name : "",
      revision: typeof record.revision === "number" ? record.revision : 0,
    }))
    .sort((a, b) => a.code.localeCompare(b.code) || a.revision - b.revision);
}
