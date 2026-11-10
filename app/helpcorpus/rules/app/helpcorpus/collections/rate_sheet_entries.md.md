<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/rate_sheet_entries.md","sha256":"4e42f6ee61ebc58d85a69aa6fc8771c3919436472add9de066be01b16c7c40ee","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/rate_sheet_entries.md

## Who can list rate sheet entries

Lines 7-11

Asked as: who can see rate sheet entries; list rate entries access; can I browse rates; who may read rate lines

Any signed-in person can list rate sheet entries (rate_sheet_entries). Anyone who is not signed in cannot.

## Who can view one rate sheet entry

Lines 13-17

Asked as: open a single rate entry; view one rate line; who can see a rate; rate entry details access

Any signed-in person can view a single rate sheet entry (rate_sheet_entries). Anyone who is not signed in cannot.

## Only procurement can create rate sheet entries

Lines 19-24

Asked as: add a new rate entry; cannot create rate; why can't I add rate; need job claim for rates

A signed-in person can create a rate sheet entry (rate_sheet_entries) only if they hold the job claim (user_claims_via_uid.cid.name = 'job'). Staff without that claim cannot create one.

## Nobody can update rate sheet entries through the API

Lines 26-28

Asked as: edit a rate entry; change rate via api; rate not saving update; can't modify my rates

Nobody can update a rate sheet entry (rate_sheet_entries) through the API, because there is no update rule. Only a PocketBase superuser or Turbo's own server code can change one.

## Nobody can delete rate sheet entries through the API

Lines 30-40

Asked as: remove a rate entry; delete rate line; rate entry will not delete; get rid of a rate

Nobody can delete a rate sheet entry (rate_sheet_entries) through the API, because there is no delete rule. Only a PocketBase superuser or Turbo's own code can delete one. Deleting its rate sheet (rate_sheets) also deletes its entries (cascadeDelete).

## Rate and overtime rate minimums

Lines 38-42

Asked as: rate must be whole number; overtime rate decimals; rate value rejected; minimum rate amount

Each entry needs a role, a rate sheet, a rate and an overtime rate. The rate (rate) must be a whole number of at least 1. The overtime rate (overtime_rate) must be at least 1 and may have decimals.

## Each role appears once per rate sheet

Lines 46-50

Asked as: duplicate role on rate sheet; same role twice; role already has rate; rate sheet role conflict

A rate sheet cannot have two entries for the same role. A unique index on role and rate sheet (role, rate_sheet) rejects a second entry.
