<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/proposedRates.ts","sha256":"a84c798e80032234484e75d1198a96febf501355a2928b9aa9aeccb0999b6140","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/proposedRates.ts

## A role has no rate when either rate is empty

Lines 33-40

Asked as: role shows no rate; rate left empty; overtime rate blank; role has no rate; rate missing

In the proposed rates table, a role counts as having no rate (rateMissing) when its rate or its overtime rate is empty. An empty value is null, undefined or an empty string. A role needs both rates to count as having a rate.

## Label for a role with a missing rate

Lines 42-48

Asked as: no overtime rate label; why does it say no overtime rate; rate label on role; missing rate label; no rate tag

The rates table labels a role "No overtime rate" when the rate is entered and only the overtime rate is empty. In every other case where a rate is missing, it shows "No rate", including when both are empty. A role with both rates has no label.

## Overtime multiplier sets overtime rate from the rate

Lines 50-64

Asked as: overtime multiplier; apply multiplier to overtime; overtime rate from rate; multiply overtime rate; set overtime by multiplier

Applying an overtime multiplier sets each row's overtime rate to the rate times the multiplier, rounded to the cent. A row with no rate keeps its overtime rate, so an empty row stays empty. A multiplier that is empty, not a number or not above 0 changes nothing. The proposed rates editor and the rate sheet copy page both use it.

## When a proposed rate counts as changed

Lines 66-74

Asked as: rate marked as changed; why is rate changed; differs from source sheet; changed rate highlight; rate changed from sheet

A role's rate counts as changed (rateChanged) when the source sheet has no rate for the role. It also counts as changed when the rate or the overtime rate differs from the source sheet's.

## Saving proposed rates leaves out roles with no rate

Lines 76-94

Asked as: save with missing rates; save incomplete proposal; roles left out on save; submit missing rate roles; save request without rates

A save of a project request sends only the roles that have both rates. A role with a missing rate is left out, so an incomplete proposal can still be saved. Submission names the roles that have no rate. The save sends the name, reason and based_on fields with the rates.

## Editor refuses a save when a role has only one rate

Lines 106-122

Asked as: cannot save one rate; enter the rate error; enter overtime rate error; only one rate entered; save blocked rate error

The request editor refuses a save while a role has one of its two rates but not the other. It shows "Enter the rate." on the rate when only the overtime rate is entered. It shows "Enter the overtime rate." when only the rate is entered. A role with both rates empty gives no error.

## Editor adds an empty row for each missing role

Lines 124-144

Asked as: role missing from rates table; empty row for role; add missing role row; role not listed in editor; rows missing for roles

The editor adds an empty row, with no rate and no overtime rate, for each rate role that has no row. This covers a source sheet, or a proposal saved before a role was added, that lacks the role. The rows are sorted by role name. Nothing changes when every role has a row.

## Source rates fill rows that have none

Lines 146-163

Asked as: source rate not showing; empty row shows sheet rate; fill rate from source sheet; source sheet rate in row; rates copied from sheet

The editor adds the source sheet's rates to each row that has no source rate, so an empty row for another role shows the source sheet's rates. It changes only the source rates and never the entered rates. A row with no matching source entry stays as it is.

## Refused save shows server errors on the rates table

Lines 165-224

Asked as: server error on rate row; save error next to rate; error shown on wrong row; rate error after refused save; save errors on fields

After a refused save, errors on name, reason and based_on show on their own fields. An error on one rate of a sent entry goes on that role's row. The entry is matched by its place in the save, not in the table. The row message drops the "entry " prefix, writes "overtime rate" with a space, and ends with a capital letter and a full stop.

## General error message leaves out errors shown in the editor

Lines 226-241

Asked as: duplicate error message on save; general error message repeated; error shown twice on save; top error message rates; too many errors on save

The general error message of a refused save leaves out the errors that the rates editor already shows on its fields and on the rate rows. The other errors stay.
