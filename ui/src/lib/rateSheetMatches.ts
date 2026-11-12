// Matches between a rate sheet and the active sheets, as the API returns them
// (POST /api/rate_sheets/compare, and the duplicate_rate_sheet error of a
// create or an activation).

export interface RateSheetMatch {
  id: string;
  name: string;
  revision: number;
  effective_date: string;
  // 0 for an exact match.
  differing_roles: number;
}

export interface RateSheetComparison {
  exact: RateSheetMatch[];
  near: RateSheetMatch[];
}

// Returns the matching sheets from a pb.send error when the save or the
// activation was blocked because an active sheet has the same rates, or null
// for any other error. field is "entries" for a create, "active" for an
// activation, "entry" for a corrected rate on an active sheet.
export function duplicateRateSheetMatches(
  error: unknown,
  field: "entries" | "active" | "entry",
): RateSheetMatch[] | null {
  const fieldError = (
    error as { data?: { data?: Record<string, { code?: string; data?: unknown }> } }
  )?.data?.data?.[field];
  if (fieldError?.code !== "duplicate_rate_sheet") return null;
  const matches = (fieldError.data as { matches?: RateSheetMatch[] } | undefined)?.matches;
  return Array.isArray(matches) ? matches : [];
}
