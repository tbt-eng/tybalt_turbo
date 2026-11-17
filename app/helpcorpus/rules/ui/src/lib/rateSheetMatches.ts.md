<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/rateSheetMatches.ts","sha256":"69f179ea2c09e4e618d8b8d3bc3bd48c788f617a3a3d98085b6ad90bf55482bb","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/rateSheetMatches.ts

## Each matching rate sheet shows its code, revision and differing roles

Lines 5-18

Asked as: matching rate sheet details; rate sheet revision shown; which rate sheet matches; differing roles count; exact or near match

Each match the API returns for a duplicate rate sheet carries the sheet's id, code, name, revision and effective date. It also carries the number of differing roles (differing_roles), which is 0 for an exact match. A comparison groups the active sheets into exact matches and near matches.

## Duplicate rate sheet error is recognised for four actions

Lines 20-32

Asked as: rate sheet already exists; duplicate rate sheet error; can't save rate sheet; rate sheet blocked duplicate; proposed rates rejected duplicate

The UI treats an error as a duplicate rate sheet block only when the error's code for the checked field is "duplicate_rate_sheet". The field is "entries" for a create, "active" for an activation, "entry" for a corrected rate on an active sheet, and "proposed_rates" for the proposed rates of a project request. For any other error, the UI finds no matching sheets (null).

## Duplicate block with no match list shows no matching sheets

Lines 33-34

Asked as: duplicate error but no matches; rate sheet blocked no list; empty matching sheets; activation blocked nothing listed; duplicate with no details

When the error has the code "duplicate_rate_sheet" but carries no list of matches, the UI returns an empty list. It does not return null. The save or activation is still treated as blocked by a duplicate, but no matching sheets are listed.
