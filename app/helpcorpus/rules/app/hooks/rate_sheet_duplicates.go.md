<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/rate_sheet_duplicates.go","sha256":"9f1aa5164abb516a324652c1222a669bcb4a15a9ccfca6538c650789cdc4bfa7","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/rate_sheet_duplicates.go

## Near-match limits for similar rate sheets

Lines 16-21

Asked as: similar rate sheets; how many near matches; differs in roles; close rate sheet limit

Callers that show "differs from X in N roles" use two limits. A rate sheet that differs in more than 5 roles (RateSheetNearMatchMaxRoles) is not shown as close. At most 3 near matches (RateSheetNearMatchLimit) are shown.

## Closest near matches: trimming rule

Lines 50-64

Asked as: closest rate sheets order; trim near matches list; too many similar sheets; which sheets are closest

Closest returns the near matches closest first. It stops at the given limit on the number of matches. It also stops at the first sheet that differs in more roles than the maximum. A sheet that differs in exactly the maximum number of roles is still included.

## Rates are compared in whole cents

Lines 78-90

Asked as: rates rounding to cents; small rate difference; fraction of a cent; rates look same but differ

Turbo rounds each rate and overtime rate to whole cents before it compares them. A difference smaller than a rounding step in the cents does not make two rates differ.

## Identical-rates key for rate sheets of any state

Lines 101-118

Asked as: same rates key; compare rate entries order; inactive sheet same rates; empty rate sheet key

RateSheetRatesKey gives two sets of entries the same key only if they have the same roles and the same rate and overtime rate in cents for each role. The order of the entries does not matter. It works on active and inactive sheets alike. It is an empty string for no entries. If a role appears twice in the entries, the last one counts.

## Comparison covers only active rate sheets, excluding one name

Lines 120-146

Asked as: compare with active sheets; why inactive sheet not matched; skip my earlier revisions; revision matching itself

CompareWithActiveRateSheets compares a set of entries with every active rate sheet. Inactive sheets are never compared. Sheets with the excluded name are skipped, so a revision is not compared with its own earlier revisions.

## What counts as an exact match or a differing role

Lines 151-195

Asked as: exact match rate sheet; what counts as differing role; missing role on one side; empty rate sheet differs

A sheet is an exact match only when it has the same set of roles and the same rate and overtime rate (in cents) for every role. A role counts as differing if its rates differ, or if it is missing on one side. A sheet with no entries has no roles, so it differs in every wanted role.

## Order of exact and near matches

Lines 197-213

Asked as: sort matching rate sheets; order of duplicate matches; ties sorted by name; which match listed first

Matches are sorted by number of differing roles, fewest first. Ties go by sheet name, then by revision, lowest first.

## Duplicate rate sheet error message

Lines 229-249

Asked as: rate sheet already exists; duplicate rates refused; identical rates error; cannot save same rates

When a rate sheet's rates are identical to an active sheet, Turbo refuses with "a rate sheet with identical rates is already active". The error on the field reads: the rates are identical to active rate sheet "<name>" (rev. <n>); use that sheet instead. The name and revision are those of the first match. The error carries every match (code duplicate_rate_sheet), so the UI can link to the existing sheets.
