<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/rate_sheet_code_changes.md","sha256":"5cdd0bcb9f8586aec0d5fed42cb4c07597d6a6169cdc6b47756fda948d612598","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/rate_sheet_code_changes.md

## Who can list rate sheet code change records

Lines 9-11

Asked as: list rate sheet code changes; see all code change history; who can view code changes; browse rate code change log; show rate sheet code log

Any signed-in person can list rate sheet code change records (rate_sheet_code_changes). Someone who is not signed in cannot.

## Who can view one rate sheet code change record

Lines 15-17

Asked as: view one code change record; open a rate code change; see single rate sheet change; look up a code change; details of code change entry

Any signed-in person can view a single rate sheet code change record (rate_sheet_code_changes). Someone who is not signed in cannot.

## Nobody can create rate sheet code change records through the API

Lines 19-21

Asked as: cannot add code change record; create rate code change via api; add new rate sheet code change; api says forbidden on code change; insert code change record

Nobody can create a rate sheet code change record (rate_sheet_code_changes) through the API, except a PocketBase superuser. Turbo's own server code can still create them.

## Nobody can update rate sheet code change records through the API

Lines 23-25

Asked as: cannot edit code change record; change rate code change via api; update rate sheet code change; api rejects code change edit; modify code change entry

Nobody can update a rate sheet code change record (rate_sheet_code_changes) through the API, except a PocketBase superuser. Turbo's own server code can still update them.

## Nobody can delete rate sheet code change records through the API

Lines 27-29

Asked as: cannot remove code change record; delete rate code change api; erase rate sheet code change; api refuses code change delete; get rid of code change entry

Nobody can delete a rate sheet code change record (rate_sheet_code_changes) through the API, except a PocketBase superuser. Turbo's own server code can still delete them.

## Required fields of a rate sheet code change record

Lines 33-39

Asked as: code change missing old code; rate code change required fields; which fields are needed for code change; rate sheet link optional; code change needs new code

Each rate sheet code change record requires an old code (old_code) and a new code (new_code). The link to a rate sheet (rate_sheet) is optional and can point to at most one rate sheet.
