<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/rate_sheet_name_changes.md","sha256":"d038740a5c70cff44f15580fab5a6654ca1b6bb8e7f407af755de13286154286","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/rate_sheet_name_changes.md

## Who can list or view rate sheet name change records

Lines 7-17

Asked as: who can see rename history; view rate sheet renames; not signed in access; list rate sheet name log

Any signed-in person can list and view records of rate sheet name changes (rate_sheet_name_changes). A request without a sign-in cannot.

## Nobody creates rate sheet name changes through the API

Lines 19-21

Asked as: add rename record via API; create rate sheet rename entry; insert name change manually; api post rate sheet name

Nobody can create a rate sheet name change record (rate_sheet_name_changes) through the API, except a PocketBase superuser. Turbo's own server code can still create one.

## Nobody updates rate sheet name changes through the API

Lines 23-25

Asked as: edit a rename record; change logged rate sheet name; fix rename history entry; update name change via API

Nobody can update a rate sheet name change record (rate_sheet_name_changes) through the API, except a PocketBase superuser. Turbo's own server code can still update one.

## Nobody deletes rate sheet name changes through the API

Lines 27-29

Asked as: remove rename log entry; delete rate sheet rename record; erase name change history; api delete rename

Nobody can delete a rate sheet name change record (rate_sheet_name_changes) through the API, except a PocketBase superuser. Turbo's own server code can still delete one.

## Required fields of a rate sheet name change

Lines 31-40

Asked as: rename record missing fields; which fields are required for rename; who renamed the rate sheet; old name and new name fields

A rate sheet name change record must name its rate sheet (rate_sheet) and the person who renamed it (renamed_by). The old name (old_name) and new name (new_name) are optional text. Turbo sets the created time (created) when the record is made.
