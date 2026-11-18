<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/rateSheetAbsorb.ts","sha256":"8cf6988ec6a6812253c5c48e7d80bbe7bbd563cad7328d4565439208cf7ee7e4","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/rateSheetAbsorb.ts

## Server decides which duplicate sheets can be absorbed and why not

Lines 1-35

Asked as: why can't this sheet be absorbed; absorb blocked reason; server says no absorb; duplicate sheet can't merge

The server decides which sheets can be absorbed. When the target can't absorb a sheet, the server gives a reason (reason), and the UI only shows it. A duplicate group can also have a reason that no sheet in it can be absorbed into its target. A sheet can absorb others when can_absorb is true.

## Duplicate rate sheet target is the earliest effective date

Lines 19-22

Asked as: which rate sheet is the main one; earliest effective date duplicate; which duplicate survives; target rate sheet choice

In a group of duplicate rate sheets, the target is the sheet with the earliest effective date (target). An absorb moves the other sheets into this target.

## Absorbing into the target deletes one sheet or all revisions of a name

Lines 26-28

Asked as: how many sheets get deleted; absorb deletes all revisions; does absorb remove every revision; delete one sheet or all

An absorb into the target deletes either 1 sheet or every revision of the name (absorbs). The server decides which. The UI only shows the count.

## Absorb candidates carry a reason only when blocked

Lines 37-41

Asked as: absorb button greyed out; candidate can't be absorbed; why is this candidate blocked; no reason means absorbable

Each absorb candidate for a target rate sheet has a reason only when the target can't absorb it. When the reason is absent, the target can absorb that candidate.

## Who may absorb rate sheets

Lines 43-49

Asked as: who can absorb rate sheets; need admin to absorb; absorb permission missing; can't absorb, no access; revise claim needed for absorb

A person with the admin claim may absorb rate sheets. A person with both the absorb claim and the rate_sheet_revise claim may also absorb them. A person with only one of those two claims may not. The server checks the same claims (canAbsorbRateSheets).

## Absorb Actions page lists absorbed sheets by code and revision

Lines 60-76

Asked as: absorb actions list order; sorted by code and revision; missing code shows blank; absorbed sheets not showing; revision shows 0

The Absorb Actions page lists the rate sheets of an absorb action's absorbed records. It sorts them by code, then by revision. Records that are not stored sheets with a text id are left out. A missing code or name shows as empty, and a missing revision shows as 0.
