<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/proposedRates.ts","sha256":"777369b975880955fc649513c4d6e7f7e59bc67c260ce2827f18d7a119473774","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/proposedRates.ts

## A proposed rate with an empty rate or overtime rate counts as missing

Lines 29-35

Asked as: rate is empty; missing rate; no overtime rate; blank rate field; rate not filled in

In a project request's proposed rates, a role has no rate (rateMissing) when its rate or its overtime rate is empty. Empty means null, undefined or an empty string. A rate of 0 is not empty.

## When a proposed rate counts as changed from the source sheet

Lines 37-45

Asked as: rate differs from source; rate changed from sheet; not matching rate sheet; overtime rate changed; role has no source rate

A role's proposed rate counts as changed (rateChanged) when the source rate sheet has no rate for the role. It also counts as changed when the rate or the overtime rate differs from the source sheet's. A role is unchanged only when both rates equal the source rates.

## Saving a project request leaves out roles that have no rate

Lines 47-63

Asked as: cannot save incomplete proposal; save without all rates; roles without rate left out; draft saved missing roles; submit missing rate roles

When a project request is saved, the proposed rates send only the roles that have both a rate and an overtime rate. Roles with a missing rate are left out, so an incomplete proposal can still be saved. Submission names the roles that have no rate. The name, the reason and the rate sheet it is based on (based_on) are always sent.

## The rates editor adds an empty row for every role without one

Lines 65-85

Asked as: role missing from rates table; empty row added for role; role not in rate sheet; add row for every role; rates editor blank row

The rates editor adds an empty row, with no rate or overtime rate, for each rate role that has no row. This happens when the source sheet, or a proposal saved before the role was added, lacks the role. The rows are sorted by role name. When every role already has a row, nothing is added (withEveryRole returns null).

## Only three rates editor fields show their own error after a refused save

Lines 87-105

Asked as: save refused error not showing; no error on reason field; rates editor error message; based on error missing; server rejected save no message

After the server refuses a project request save, the rates editor shows its own error only for the name, the reason and the based-on fields (proposed_rates.name, proposed_rates.reason, proposed_rates.based_on). Only errors that have a message are shown. Other field errors are not shown by the editor.
