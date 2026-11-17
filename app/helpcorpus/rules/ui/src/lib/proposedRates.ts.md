<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/proposedRates.ts","sha256":"32d8727e9ba3ea551a3dd14b97d44e8325d5a8058612b63fdaa796f0a7452ebc","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/proposedRates.ts

## A role has no rate if either rate is empty

Lines 29-36

Asked as: role has no rate; missing overtime rate; rate is empty; role without rate; rate missing for role

In the proposed rates of a project request, a role counts as missing its rate when its rate or its overtime rate is empty (null, undefined or an empty string). A role needs both rates to count as having a rate (rateMissing).

## Overtime multiplier sets overtime rate from rate, rounded to the cent

Lines 38-52

Asked as: overtime multiplier; multiply overtime rate; apply overtime multiplier; overtime rate not calculated; overtime rounding to cent

Applying an overtime multiplier sets each row's overtime rate to the rate times the multiplier, rounded to the cent. A row with no rate keeps its overtime rate, so an empty row stays empty. A multiplier that is empty, not a number or not above 0 changes nothing. The proposed rates editor and the rate sheet copy page both use this (withOvertimeMultiplier).

## When a proposed rate counts as changed

Lines 54-62

Asked as: rate changed from source; rate differs from sheet; changed proposed rate; highlight changed rate; rate different from rate sheet

A proposed rate counts as changed when the source sheet has no rate for the role, or when the rate or the overtime rate differs from the source sheet's (rateChanged).

## Saving a request leaves out roles with no rate

Lines 64-82

Asked as: save request without all rates; incomplete proposal save; role dropped when saving; unsaved role rates; submit request missing rates

When a project request is saved, the proposed rates include only roles that have both a rate and an overtime rate (entries). A role missing either rate is left out, so an incomplete proposal can still be saved. Submitting the request names the roles that have no rate.

## Editor refuses a save when a role has only one of its two rates

Lines 94-110

Asked as: cannot save rate; enter the rate error; only one rate filled; overtime rate missing error; save refused rate row

The request editor refuses a save when a role has a rate but no overtime rate, or an overtime rate but no rate. It shows "Enter the rate." or "Enter the overtime rate." on that role's row. A role with both rates empty gets no error.

## Editor adds an empty row for each rate role that lacks one

Lines 112-132

Asked as: role missing from rates table; no row for role; add empty row for role; role not listed in rates; rates table missing role

A proposal needs a rate for every role. The editor adds an empty row, with no rates, for each rate role that has no row, for instance when the source sheet lacks the role or the proposal was saved before the role was added. The rows are sorted by role name. Nothing changes when every role already has a row (withEveryRole).

## Source rates fill in only rows that have none

Lines 134-151

Asked as: source rate not filled in; fill rates from sheet; rate sheet rates not copied; empty row no rates; copy source rates to row

The editor adds the source sheet's rate and overtime rate to each row that has no source rates and whose role is on the source sheet. It never changes the rate or overtime rate typed in the row. A saved proposal has source rates only for roles it has a rate for, so the empty row of another role gets them this way (withSourceRates).

## Server errors on a rate show on the role's row

Lines 153-212

Asked as: error on rate row; rate error not shown; server error on rate; which role has the error; rate save error location

When the server refuses a save of a project request, an error on one rate of a sent entry (key proposed_rates.entries.N.rate, overtime_rate or role) appears on that role's row in the rates table. The entry is found by its place in the save, not its place in the table. Errors on name, reason and based_on show on their own fields.

## Row error messages are reworded for the rates table

Lines 172-182

Asked as: overtime_rate message wording; entry overtime_rate error; odd rate error message; rate error text unclear; raw error message rates

The server's message on one rate is shown in the table's words. The leading "entry " is dropped, "overtime_rate" becomes "overtime rate", the first letter is capitalised and a full stop is added when the message has no ending punctuation. For example "entry overtime_rate must be at least 1" shows as "Overtime rate must be at least 1."

## General error message leaves out errors shown on rows

Lines 214-225

Asked as: duplicate error message; error shown twice; general error repeats row error; save error message cluttered; extra error at top of form

In the general error message of a refused save, errors already shown on a role's row are removed, because their keys name a place in the save. The other errors stay.
