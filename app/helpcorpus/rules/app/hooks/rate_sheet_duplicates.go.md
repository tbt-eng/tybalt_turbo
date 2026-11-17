<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/rate_sheet_duplicates.go","sha256":"338d1935bbea0b1a187ebe5ef3377c317beae56eaef8e111d263fcce331d87c8","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/rate_sheet_duplicates.go

## Near-match limits for similar rate sheets

Lines 16-21

Asked as: similar sheets not shown; why only three matches; too many differences hidden; sheet differs in many roles

Callers that show "differs from X in N roles" use two constants (RateSheetNearMatchMaxRoles and RateSheetNearMatchLimit). A rate sheet that differs in more than 5 roles is not shown as close. At most 3 near matches are shown.

## Differing roles counts changed and one-sided roles

Lines 31-41

Asked as: differs in N roles count; roles count looks wrong; missing role counted as difference; what does differs in mean

For a rate sheet compared with a set of entries, the differing roles (DifferingRoles) count every role whose rate or overtime rate differs. They also count every role that is on one side only. An exact match has 0.

## Comparison result separates exact matches from near matches

Lines 43-49

Asked as: which sheets are identical; list of similar rate sheets; exact versus near matches; comparison results list

The comparison (CompareWithActiveRateSheets) returns two lists. Exact (Exact) holds the active sheets with identical entries. Near (Near) holds every other active sheet that was compared, closest first.

## Closest near matches are trimmed and sorted

Lines 51-65

Asked as: why matches list cut short; only closest sheets shown; sheet dropped from similar list; limit on similar sheets

The closest-matches helper (Closest) returns at most the given limit of near matches, closest first. It stops at the first sheet that differs in more roles than the maximum. A sheet that differs in exactly the maximum number of roles is still included.

## Sheet with no entries compares as empty

Lines 67-69

Asked as: rate sheet with no roles; empty rate sheet duplicate check; blank rate sheet compared; sheet without any entries

An active rate sheet with no entries is compared as a sheet with no roles. Its empty role row is ignored.

## Rates compare in whole cents

Lines 80-92

Asked as: rates flagged as same wrongly; decimal rate shows as different; rounding causes duplicate warning; cents rounding rate difference

Turbo rounds each rate and overtime rate to whole cents before it compares them. Float rounding on the way in does not make two rates differ.

## Identical rates key ignores order and state

Lines 103-120

Asked as: same rates different order; rates key; same roles listed differently; identical rates different row order

Two sets of entries share a rates key (RateSheetRatesKey) when they have the same roles and the same rate and overtime rate in cents for each role. The key does not depend on the order of the entries. It works for sheets of any state, active or inactive. It is an empty string for no entries.

## Comparison covers only active rate sheets

Lines 122-149

Asked as: inactive rate sheet not flagged; old sheet not matched; revision matched to earlier revision; why old rates not checked

Turbo compares entries only with active rate sheets (CompareWithActiveRateSheets). It never compares with inactive sheets. It leaves out every revision of the excluded code, so a revision is not compared with its own earlier revisions. A sheet with no code yet passes an empty exclude code.

## Failure to load active sheets stops the comparison

Lines 150-152

Asked as: duplicate check fails to load; rate sheet check error; no matches returned on error; cannot read active rate sheets

If Turbo cannot read the active rate sheets, the comparison returns the error and no matches.

## Comparison query reads a sheet's code, name, revision and effective date

Lines 154-175

Asked as: which existing sheet matched; match shows no name; sheet code and revision shown; effective date of matching sheet

Each match carries the sheet's code, name, revision and effective date, so that a person can identify the existing sheet.

## Roles with a blank name are ignored when comparing

Lines 176-178

Asked as: blank role row ignored; role with no name; empty role line in rates; unnamed role in rate sheet

A row with an empty role adds no entry to the sheet. A sheet with no entries therefore compares as having no roles.

## Exact and near match definition

Lines 181-199

Asked as: when are rates identical; what counts as same rates; how close is near match; sheet counts as duplicate

An active sheet is an exact match when it has the same set of roles and the same rate and overtime rate in cents for every role. The differing-roles count adds the roles whose rate or overtime rate differ and the roles on one side only. A count of 0 is an exact match. Any other count is a near match.

## Match ordering

Lines 206-217

Asked as: order of similar sheets; which match listed first; closest sheet sorted first; matches sorted by code

Exact and near matches are sorted by the number of differing roles, then by code, then by revision. The closest sheet comes first.

## Duplicate rate sheet error

Lines 233-253

Asked as: a rate sheet with identical rates is already active; rates identical to active sheet; duplicate rate sheet error; use existing sheet instead

When a rate sheet has rates identical to an active sheet, Turbo refuses with "a rate sheet with identical rates is already active". The field error (duplicate_rate_sheet) says "the rates are identical to active rate sheet <label>; use that sheet instead". It names the first match and lists every match so the UI can link to the existing sheet.
