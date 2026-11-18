<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/proposedRates.ts","sha256":"3b479e89dc71406fcb8bf6ff92edbb80323fb20a45552eac7b7284a5acb23b2d","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/proposedRates.ts

## A role has no rate if either rate is empty

Lines 37-40

Asked as: role has no rate; rate left blank; empty rate field; role missing rate; incomplete rates row

In the proposed rates table, a role counts as having no rate when its rate or its overtime rate is empty (null, undefined or an empty string). The role needs both rates to count as complete (rateMissing).

## Missing-rate label: No overtime rate or No rate

Lines 42-48

Asked as: no overtime rate label; why says no rate; missing rate tag; overtime rate blank label; role labelled no rate

The table labels a role "No overtime rate" when only its overtime rate is empty. It labels the role "No rate" when the rate is empty, or when both rates are empty. A role with both rates has no label (missingRateLabel).

## Overtime multiplier sets overtime rate from the rate

Lines 50-64

Asked as: overtime multiplier not working; apply multiplier to rates; multiply rates for overtime; rate times multiplier; multiplier fills overtime rate

Applying an overtime multiplier sets each row's overtime rate to the rate times the multiplier, rounded to the cent. A row with no rate keeps its overtime rate, so an empty row stays empty. A multiplier that is empty, not a number or not above 0 changes nothing. The proposed rates editor and the rate sheet copy page both use it (withOvertimeMultiplier).

## Multiplier changes list leaves out unchanged rows

Lines 66-74

Asked as: multiplier list missing rows; rows not changed by multiplier; unchanged overtime rows hidden; multiplier skipped row; refused save after multiplier

The list of rows that the multiplier changes leaves out any row whose overtime rate stays the same. Such a row keeps its state, such as a refused save (overtimeMultiplierChanges).

## When a rate counts as changed from the source sheet

Lines 76-84

Asked as: rate differs from source sheet; why marked changed; rate highlighted as changed; no source rate for role; overtime differs from sheet

A role's rate counts as changed when the source sheet has no rate for the role. It also counts as changed when the rate or the overtime rate differs from the source sheet's (rateChanged).

## Saving a request leaves out roles with no rate

Lines 86-104

Asked as: save request with missing rate; save incomplete proposal; roles without rate dropped; request saved without some roles; payload missing rates

When a project request is saved, the proposed rates payload carries the name, reason and based_on, and only the roles that have both rates. A role with no rate is left out, so an incomplete proposal can still be saved. Submission names the roles that have no rate (proposedRatesPayload).

## Comparison with active sheets depends on sheet and rates only

Lines 106-113

Asked as: comparison not updating; active sheets comparison stale; typing name changes comparison; rate change comparison refresh; compare with active sheets

The live comparison with the active rate sheets reads the sheet the proposal is based on (based_on) and the rates that a save would send. Typing the name or the reason keeps the comparison. Changing a rate makes it stale (comparisonKey).

## A stale comparison is never shown

Lines 115-120

Asked as: old comparison showing; comparison wrong after edit; stale comparison shown; comparison still loading old result; comparison from previous rates

The editor shows the stored comparison only when it was made for the current rates. Otherwise it shows none, even while the new comparison is loading (comparisonFor).

## Editor refuses a save when a role has only one rate

Lines 132-148

Asked as: cannot save one rate; enter the rate error; enter overtime rate error; save refused missing rate; only one rate entered

The request editor refuses a save when a role has a rate but no overtime rate, and shows "Enter the overtime rate." on that role. It shows "Enter the rate." when the role has an overtime rate but no rate. A role with both rates empty gets no error and the save leaves it out. The refusal stops a typed rate from being lost (incompleteRateErrors).

## Editor adds an empty row for each missing role

Lines 150-170

Asked as: role missing from proposal; role not in rates table; empty row for every role; added role missing rows; source sheet lacks role

A proposal needs a rate for every rate role. The editor adds an empty row for each role that has no row, sorted by role name. It adds nothing when every role already has a row. The source sheet, or a proposal saved before a role was added, can lack a role (withEveryRole).

## Source rates fill empty rows without changing rates

Lines 172-189

Asked as: fill rates from source sheet; empty rows copy sheet rates; source rates not filled; prefill rates from sheet; rates copied without overwriting

For a row that has no source rates, the editor fills in the source sheet's rate and overtime rate when the source sheet has the role. It never changes the proposed rates. A row that already has source rates, or a role the source sheet lacks, stays as it is (withSourceRates).

## Review of proposed rates shows saved and missing roles

Lines 203-220

Asked as: review shows missing roles; review rates list empty; roles without rate in review; review table blank rows; saved rates review order

The review of a request lists each saved rate and an empty row for each role with no rate, sorted by role name. Each row shows the source sheet's rates where it has them. The rates table marks an empty row "No rate" (reviewRows).

## Server rate errors are placed on the role's row

Lines 222-281

Asked as: error on wrong role row; server rate error location; overtime rate error message; save error next to rate; which role has the error

When a request save is refused, an error on an entry's rate, overtime rate or role goes on that role's row in the rates table. The server numbers entries by their place in the save, not in the table, so Turbo maps the number back to the role. The message is reworded: the "entry " prefix is dropped, "overtime_rate" becomes "overtime rate", the first letter is capitalised and a full stop is added. For example, "Overtime rate must be at least 1." Errors on name, reason and based_on show on their own fields (proposedRatesFieldErrors).

## General error message leaves out errors shown on fields

Lines 283-298

Asked as: duplicate error message; general error repeats field error; error shown twice on save; top error message clutter; save error text repeated

The general error message for a refused save leaves out the errors that the editor already shows on the name, reason and based_on fields or on a rate row. The other errors stay in it (withoutShownErrors).
