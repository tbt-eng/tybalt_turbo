// Duplicate rate sheets and the absorb of one into another, as the API
// returns them (GET /api/rate_sheets/duplicates and
// GET /api/rate_sheets/{id}/absorb_candidates). The server decides which
// sheets can be absorbed and says why not; the UI only shows it.

export interface RateSheetAbsorbSheet {
  id: string;
  name: string;
  revision: number;
  effective_date: string;
  active: boolean;
  job_count: number;
  request_count: number;
  // The number of revisions of the sheet's name.
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
