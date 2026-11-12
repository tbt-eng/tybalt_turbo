<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/rateSheetMatches.ts","sha256":"3b5fca12a2aacf74b5640918b0106d02710aad6d214dce30e5cbb4cb2a58d0e1","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/rateSheetMatches.ts

## Rate sheet match details and exact versus near matches

Lines 5-17

Asked as: rate sheet match details; exact vs near match; rate sheet differing roles; what does exact match mean; compare rate sheets results

Each match shows the active sheet's id, name, revision, effective date (effective_date) and the number of roles whose rates differ (differing_roles). A differing_roles of 0 means an exact match. A comparison (POST /api/rate_sheets/compare) returns two lists: exact matches (exact) and near matches (near).

## Duplicate rate sheet error on create or activation

Lines 19-33

Asked as: duplicate rate sheet error; can't save rate sheet; rate sheet already exists; activation blocked duplicate; rate sheet error on entries

When saving a new rate sheet or activating one is blocked with the error code duplicate_rate_sheet, Turbo reads the list of matching active sheets from the error. A create reports the error on the field "entries". An activation reports it on the field "active". Any other error is not treated as a duplicate.

## Duplicate error without a match list shows no matches

Lines 30-32

Asked as: duplicate error no matches; blocked but no matching sheets; empty match list rate sheet; duplicate with no details; rate sheet error missing matches

When the duplicate_rate_sheet error has no list of matches (matches), Turbo still treats the save or activation as blocked as a duplicate. It returns an empty list of matches instead of nothing.
