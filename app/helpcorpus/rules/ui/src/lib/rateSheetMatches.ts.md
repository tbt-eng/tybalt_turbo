<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/rateSheetMatches.ts","sha256":"63021466b6ad964bc241d2fd2186bbbf2e31c80cd1731f0aad62467aa1239427","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/rateSheetMatches.ts

## Rate sheet comparison separates exact and near matches

Lines 5-17

Asked as: exact vs near match; rate sheet comparison results; differing roles; similar rate sheets list; what does near match mean

A rate sheet comparison (POST /api/rate_sheets/compare) returns two lists of active sheets: exact matches and near matches. An exact match has differing_roles of 0. A near match is any other sheet that is returned, and it carries its number of differing roles.

## Duplicate rate sheet error is recognised on create, activation or rate correction

Lines 19-30

Asked as: can't save rate sheet; rate sheet save blocked; duplicate rate error; activation blocked same rates; rates already exist

The Turbo UI treats a save or an activation as blocked when the error has the code duplicate_rate_sheet. The code must be on the field "entries" for a create, "active" for an activation, or "entry" for a corrected rate on an active sheet. This means an active rate sheet already has the same rates. Any other error gives no matches (null).

## Blocked duplicate rate sheet lists matching active sheets, or an empty list

Lines 31-32

Asked as: which active sheets clash; blocked by existing sheet; matching sheets list empty; rate sheet conflict details; same rates as another sheet

When the duplicate_rate_sheet error is found, the UI takes the list of matching active sheets from the error's data (matches). Each match has an id, name, revision, effective date and a count of differing roles (differing_roles, 0 for an exact match). If the error has no valid list of matches, the UI uses an empty list and does not return null.
