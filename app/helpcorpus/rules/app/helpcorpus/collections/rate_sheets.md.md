<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/rate_sheets.md","sha256":"475a9e0c5d82a958846d903865ff0adb2996197727f2a8aed1699d9878dda3ec","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/rate_sheets.md

## Who can list or view rate sheets

Lines 7-17

Asked as: see rate sheets; view rate sheet list; who can open rate sheets; rate sheet access

Any signed-in person can list rate sheets and view one rate sheet. A request with no signed-in person is refused.

## Who can create a rate sheet

Lines 19-27

Asked as: make a new rate sheet; cannot create rate sheet; create revision permission; procurement new revision

A signed-in person can create a rate sheet if they hold the rate sheet revise claim (rate_sheet_revise) or the admin claim. Procurement (the job claim) can also create one, but only when the request sends revision 0. A job-claim holder without the other claims cannot create a revision above 0.

## Nobody can update or delete rate sheets through the API

Lines 29-35

Asked as: edit rate sheet via api; delete rate sheet api; change saved rate sheet; rate sheet cannot be modified

No rule allows an update or a delete of a rate sheet through the API. Only a PocketBase superuser or Turbo's own server code can do either.

## Rate sheet date must be in YYYY-MM-DD format

Lines 39-44

Asked as: rate sheet date format; effective date wrong; rate sheet save error date; revision must be number

A rate sheet needs an effective date (effective_date) in the form YYYY-MM-DD, such as 2025-01-31. It also needs a name. The revision (revision) must be a whole number.

## Name and revision must be unique together

Lines 49-53

Asked as: duplicate rate sheet name; same rate sheet twice; rate sheet already exists; name taken revision

Two rate sheets cannot have the same name and revision together. The same name can appear again with a different revision.
