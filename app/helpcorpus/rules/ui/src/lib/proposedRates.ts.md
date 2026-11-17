<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/proposedRates.ts","sha256":"a386bd8ef519afc7be0b9b7d817ec7141a7dde48d545ce76e7f0b7e563c355bc","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/proposedRates.ts

## A new proposed rates form starts empty

Lines 29-31

Asked as: new proposed rates form blank; empty rates form on start; new rate request empty fields

A new proposed rates form has an empty description (name), an empty reason, an empty based-on sheet (based_on) and no rate rows.

## A role has no rate if either rate is empty

Lines 34-40

Asked as: role missing rate; rate shows blank for role; overtime rate empty role

In the proposed rates table, a role counts as missing its rate when its rate or its overtime rate is empty. An empty value is null, undefined or an empty string.

## Overtime multiplier sets overtime to rate times multiplier

Lines 42-56

Asked as: overtime multiplier not working; set overtime from multiplier; overtime calculated from rate

The overtime multiplier sets each row's overtime rate to the rate times the multiplier, rounded to the cent. The proposed rates editor and the rate sheet copy page both use it.

## Multiplier ignored if empty, not a number, or not positive

Lines 49-54

Asked as: multiplier does nothing; multiplier typed but no change; multiplier negative or zero

A multiplier that is empty, not a number, zero or negative changes nothing. A row with an empty or non-numeric rate keeps its overtime rate, so an empty row stays empty.

## A rate counts as changed against the source sheet

Lines 58-66

Asked as: rate differs from source sheet; highlighted rate changed; rate changed from original sheet

A role's rate counts as changed when the source sheet has no rate for the role. It also counts as changed when the rate or the overtime rate differs from the source sheet's value.

## Saving a request leaves out roles with no rate

Lines 68-86

Asked as: save incomplete proposal; save without all roles; saving request missing rates

A save of a project request sends the description, reason and based-on sheet, plus an entry for each role that has both rates. A role missing either rate is left out, so an incomplete proposal can still be saved. Submission names the roles that have no rate.

## Editor refuses a save when a role has only one rate

Lines 98-114

Asked as: cannot save only one rate; enter the rate message; overtime rate required error

The request editor refuses a save when a role has one of its two rates and not the other. It shows "Enter the rate." on the rate when only the overtime rate is filled. It shows "Enter the overtime rate." on the overtime rate when only the rate is filled.

## A role with both rates empty gives no error

Lines 98-114

Asked as: both rates blank no error; empty role save error; role left blank on save

A role with both its rate and its overtime rate empty gets no error from the editor, and the save leaves that role out.

## Editor adds an empty row for each missing role

Lines 116-136

Asked as: missing role row added; role not in rate table; add rows for roles

The editor adds an empty row for each rate role that has no row, because a proposal needs a rate for every role. The source sheet, or a proposal saved before a role was added, can lack a role. The rows are sorted by role name. When every role has a row, nothing is added.

## Source rates fill only rows that have none

Lines 138-155

Asked as: fill rates from source sheet; source sheet rates not copied; autofill rates from based on

The editor adds the source sheet's rate and overtime rate to each row that has no source rate and whose role is on the source sheet. It never changes the rates the requester entered. A row for a role the source sheet lacks gets nothing.

## Refused save shows errors on the description, reason and based-on fields

Lines 157-161

Asked as: no error under name field; error next to reason box; based on field error message

After a refused save, the description (name), reason and based-on sheet (based_on) fields each show their own server error message, when the server gives one.

## Rate errors from the server go on the role's row

Lines 171-216

Asked as: rate error on wrong row; server rate message on role; overtime rate must be at least

When the server refuses a rate in a save, the editor shows the error on that role's row. It finds the role by the place of the entry in the save, not in the table. It rewrites the message in table words, such as "Overtime rate must be at least 1." for "entry overtime_rate must be at least 1". The message starts with a capital letter and ends with a full stop.

## General error message leaves out errors shown on rows

Lines 218-229

Asked as: duplicate error at top; general error shows row errors; top of form error repeated

The general error message for a refused save leaves out the errors that the rates table shows on a row, because their keys name a place in the save. The other errors stay.
