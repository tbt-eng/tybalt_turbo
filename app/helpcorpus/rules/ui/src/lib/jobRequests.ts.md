<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/jobRequests.ts","sha256":"dc623e169494fb069fa1eaaeb02e4862132dc3d91e73a323200af6d038949ce9","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/jobRequests.ts

## Print check lists missing project request details

Lines 1-95

Asked as: can't print project request; what's missing before printing; print blocked missing fields; list of missing items

Before printing a project request, Turbo builds a list of missing items. The list can contain: Client name, Description, Valid Plus Code location, Project contact name and email, Billing contact name and email, Positive whole-dollar project value, Fee basis, Rate sheet and Client billing rates. The server checks the rules again on each action.

## Contact email shows a note when the client has none

Lines 23-26

Asked as: no email address showing; contact email blank message; why says no email; contact has no email

A contact's email address shows as written when it exists. If the contact has no email and is marked no email (no_email), Turbo shows "No email address". Otherwise it shows nothing.

## Billing address falls back to the client's address

Lines 30-43

Asked as: billing address blank; invoice address missing; address taken from client; empty billing address

The billing address of a project request is the invoicing profile's address when that has any address part. If the invoicing profile has none, Turbo uses the client's address. An address joins the parts (address, city, province/state, postal code, country) with commas and skips empty parts.

## Project value shows as whole Canadian dollars or Not provided

Lines 44-52

Asked as: project value says not provided; dollar amount blank; why value not showing; value shows zero

A project value shows as Canadian dollars with no cents only when it is a positive whole number. Zero, negative, fractional and non-number values show "Not provided".

## Location must be a valid Plus Code

Lines 59-63

Asked as: location rejected; plus code not accepted; site location wrong format; location format error

A project request's location must be a Plus Code, or Turbo lists "Valid Plus Code location" as missing. It has 8 characters, a plus sign, then 2 or 3 characters, from the Plus Code alphabet (2-9, C, F, G, H, J, M, P, Q, R, V, W, X), in either case. Blank space around it is ignored.

## Contact email required unless the contact has no email

Lines 64-71

Asked as: contact email rejected; billing contact incomplete; invalid contact email; contact name missing

The project contact and the billing contact each need a name, or Turbo lists "Project contact name" or "Billing contact name". Each also needs a valid email (text@text.text, no spaces), or Turbo lists "... email". A contact marked no email (no_email) skips the email check.

## Project value must be a positive whole number

Lines 72-73

Asked as: cents in project value; decimal project value error; value must be whole dollars; project value not accepted

A project request's project value must be a whole number above zero. Otherwise Turbo lists "Positive whole-dollar project value" as missing.

## Fee basis must be explicitly chosen

Lines 74-75

Asked as: fee basis not set; time and materials not chosen; billing type missing; can't pick fee type

A project request needs its fee basis chosen (time_and_materials_set is true) and the time-and-materials value must be true or false. Otherwise Turbo lists "Fee basis" as missing.

## Rate sheet readiness depends on whether rates are proposed

Lines 76-83

Asked as: rate sheet error; rate sheet not active; role missing a rate; rate sheet not ready

For proposed rates, the rate sheet is ready when no role lacks a rate (rate_sheet_missing_roles is empty). For a chosen existing rate sheet, it must be active. A request with no rate sheet, or one that is not ready, gets "Rate sheet" in the missing list.

## Time-and-materials needs complete client billing rates

Lines 84-93

Asked as: billing rates missing overtime; no rates on time and materials; rate entries incomplete; client rates not complete

Only for a time-and-materials request, Turbo checks the rate sheet entries. It lists "Client billing rates" as missing when there are no entries or any entry has no valid rate or overtime rate.

## Job creation queue visibility

Lines 96-103

Asked as: can't see creation queue; job creation queue missing; procurement queue not showing; no queue for job requests

A person sees the job creation queue only when project requests are enabled, they hold the job claim (procurement), and the server has sent the count for "/job-requests?queue=creation".
