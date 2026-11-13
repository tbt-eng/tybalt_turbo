<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/rateSheetMatches.ts","sha256":"97dc813f6a81210084cfb709baba4f59622e47293f25e1ff244a1445400bea81","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/rateSheetMatches.ts

## Exact and near rate sheet matches are told apart by differing roles

Lines 5-17

Asked as: exact match vs near match; rate sheet almost identical; differing roles count; what does differing_roles mean; 0 differing roles

A rate sheet comparison separates exact matches from near matches. Each match gives the number of roles that differ (differing_roles), which is 0 for an exact match.

## Duplicate rate sheet error returns the matching active sheets

Lines 19-33

Asked as: duplicate rate sheet error; rate sheet already exists; which sheets have the same rates; matching active sheets; save failed same rates

When a save or activation fails with the error code duplicate_rate_sheet, an active rate sheet already has the same rates. Turbo then shows the matching sheets from the error: id, name, revision, effective date and number of differing roles (differing_roles). Any other error does not give matches.

## Duplicate rate sheet check covers four actions

Lines 19-31

Asked as: cannot activate rate sheet; blocked because rates already used; project request rates duplicate; correcting rate on active sheet blocked; new rate sheet rejected same rates

Turbo can block four actions because an active rate sheet has the same rates. These are creating a rate sheet (field entries), activating one (active), correcting a rate on an active sheet (entry), and the proposed rates of a project request (proposed_rates). The error is read only from the field that matches the action.

## Duplicate error without a match list shows no sheets

Lines 32-33

Asked as: blocked duplicate but no sheet named; duplicate error no list; activation blocked without matches; empty match list rate sheet; duplicate with no sheet shown

If the duplicate_rate_sheet error has no list of matches, or the matches are not a list, Turbo treats the match list as empty. The save or activation is still recognized as blocked as a duplicate, but no sheets are named.
