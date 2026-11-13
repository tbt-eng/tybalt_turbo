<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/rate_sheet_entries.md","sha256":"f2af76734bb0b87c326f7719f038b981b7038686b2c51f584c76204b54e28a42","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/rate_sheet_entries.md

## Who can list or view rate sheet entries

Lines 9-17

Asked as: cannot see rate entries; rate sheet entries hidden; why can't I view rates; missing rates in list

Any signed-in person can list and view a rate sheet entry (rate_sheet_entries), but only when its rate sheet does not belong to a project request (rate_sheet.request is empty). Entries of a rate sheet that belongs to a project request are hidden from the API for list and view.

## Who can create a rate sheet entry

Lines 21-26

Asked as: add a rate; create rate entry; who can add rates; requester adding rate

Only procurement (the job claim) can create a rate sheet entry through the API. The rate sheet in the request must not belong to a project request (rate_sheet.request must be empty). A requester changes a draft through the project request instead.

## Nobody can update or delete rate sheet entries through the API

Lines 28-34

Asked as: edit a rate; change rate after saving; delete rate entry; fix wrong rate

A rate sheet entry (rate_sheet_entries) has no update rule and no delete rule. Nobody can update or delete one through the API, except a PocketBase superuser. Turbo's own server code can still do both. A new rate needs a new revision of the rate sheet.

## Rate and overtime rate minimums on an entry

Lines 40-44

Asked as: rate must be whole number; overtime rate decimals; rate validation error; minimum rate value

Every rate sheet entry needs a role, a rate sheet, a rate and an overtime rate. The rate must be a whole number of at least 1 (min 1, onlyInt). The overtime rate must be at least 1 and can have decimals.

## Deleting a rate sheet deletes its entries

Lines 42-42

Asked as: rate sheet removed rates gone; deleting sheet removes rates; rates vanish with sheet; what happens when sheet deleted

The rate sheet field of an entry is a cascade-delete relation. When a rate sheet is deleted, Turbo deletes all of its entries.

## One entry per role in a rate sheet

Lines 50-52

Asked as: duplicate role on rate sheet; role already has rate; two rates same role; rate for role rejected

A rate sheet can have only one entry for each role. A unique index on role and rate sheet rejects a second entry with the same pair.
