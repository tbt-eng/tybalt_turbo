<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/proposedRates.ts","sha256":"151b1694ea2cb2bda4959355411aa28b2ce6aecf67b55f945284b82a29b2fec2","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/proposedRates.ts

## A role has no rate if rate or overtime rate is empty

Lines 29-36

Asked as: rate missing; role has no rate; empty overtime rate; blank rate field; zero rate counts as missing

In a proposed rate sheet, a role counts as missing its rate (rateMissing) when its rate is empty or its overtime rate is empty. Empty means null, undefined or an empty string. A rate of 0 is not empty.

## Overtime multiplier sets overtime rate from rate, rounded to the cent

Lines 38-52

Asked as: overtime multiplier not applied; overtime rate from multiplier; multiplier rounding to cent; empty row stays empty; overtime rate calculation

The overtime multiplier sets each row's overtime rate to the rate times the multiplier, rounded to the cent (withOvertimeMultiplier). A row with no rate, or a rate that is not a finite number, keeps its current overtime rate, so an empty row stays empty. The proposed rates editor and the rate sheet copy page both use it.

## Invalid overtime multiplier changes nothing

Lines 45-46

Asked as: zero overtime multiplier; negative multiplier ignored; multiplier does nothing; overtime multiplier invalid; multiplier left blank

If the overtime multiplier is empty, not a finite number, or zero or negative, Turbo changes no overtime rates. Only a positive number is applied.

## When a proposed rate counts as changed

Lines 54-62

Asked as: why rate shows as changed; rate differs from source; changed flag on role; overtime rate changed; no source rate changed

A proposed rate counts as changed (rateChanged) when the source sheet has no rate for the role. It also counts as changed when the rate differs from the source rate, or the overtime rate differs from the source overtime rate.

## Saving a request leaves out roles with no rate

Lines 64-80

Asked as: cannot save request incomplete; save without all rates; roles missing from saved request; request saved but rate gone; incomplete proposal save

When a project request is saved, the proposed rates (proposed_rates) include only roles that have both a rate and an overtime rate. Roles missing either rate are left out, so an incomplete proposal can still be saved. Submission is what names the roles that have no rate. The name, reason and based_on are sent as entered.

## Editor adds an empty row for every role lacking one

Lines 82-102

Asked as: missing role row added; role not in rates editor; new empty row appeared; rate sheet missing role; rates sorted by role name

The rates editor adds an empty row (null rate and overtime rate) for each rate role that has no row, then sorts all rows by role name (withEveryRole). A proposal needs a rate for every role, and the source sheet or an older saved proposal can lack one. When every role already has a row, nothing changes (null is returned).

## Refused save shows errors on name, reason and based on

Lines 104-122

Asked as: save error not showing; server error message missing; error only on name field; request save refused errors; reason field error shown

After a refused request save, the rates editor shows its own error message only for the proposed_rates name, reason and based_on fields. Server errors on other fields are not shown there. Only errors that have a message are kept.
