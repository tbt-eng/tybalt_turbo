// Who may change rate sheets, by claim. Each function mirrors the server's
// check for the route or collection rule named beside it, so the pages offer
// a control only to the users that the server allows.

// A new sheet or a copy (POST /api/rate_sheets with revision 0), and a new
// entry on an incomplete sheet (the rate_sheet_entries create rule).
export function canCreateRateSheets(claims: string[]): boolean {
  return claims.includes("job");
}

// Activate and deactivate (POST /api/rate_sheets/{id}/activate and
// /deactivate).
export function canActivateRateSheets(claims: string[]): boolean {
  return claims.includes("job");
}

// A revision (POST /api/rate_sheets with a revision above 0) and a rename
// (POST /api/rate_sheets/{id}/rename).
export function canReviseRateSheets(claims: string[]): boolean {
  return claims.includes("admin") || claims.includes("rate_sheet_revise");
}

// A correction of an entry's rates (PUT /api/rate_sheet_entries/{id}) and a
// move of projects to a revision (/api/rate_sheets/{id}/rollout). The server
// also refuses a rollout by an inactive user.
export function canCorrectRateSheets(claims: string[]): boolean {
  return claims.includes("admin");
}

// A delete of an inactive sheet that nothing uses (DELETE
// /api/rate_sheets/{id}). The server checks the sheet.
export function canDeleteRateSheets(claims: string[]): boolean {
  return claims.includes("admin");
}
