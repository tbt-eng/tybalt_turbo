<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/rate_sheet_name_changes.md","sha256":"26150999c44429e16a293e58e93eb42d9243b56b41c14219e368ebc932bbaedd","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/rate_sheet_name_changes.md

## Who can list rate sheet name change records

Lines 9-11

Asked as: list rate sheet renames; see rename history; who can see name changes; browse rate sheet name log

Any signed-in person can list the records of rate sheet name changes (rate_sheet_name_changes). A request without a signed-in user cannot.

## Who can view one rate sheet name change record

Lines 15-17

Asked as: view one rename record; open a rate sheet rename; see details of a name change; check a rate sheet rename

Any signed-in person can view one record of a rate sheet name change (rate_sheet_name_changes). A request without a signed-in user cannot.

## Nobody can create rate sheet name change records through the API

Lines 21-21

Asked as: cannot add rename record; create rate sheet rename via api; manually log a name change; api rename insert forbidden

No one can create a rate sheet name change record (rate_sheet_name_changes) through the API, except a PocketBase superuser. Turbo's own server code can still create them.

## Nobody can update rate sheet name change records through the API

Lines 25-25

Asked as: edit a rename record; change a logged rename; fix rate sheet rename entry; api update rename blocked

No one can update a rate sheet name change record (rate_sheet_name_changes) through the API, except a PocketBase superuser. Turbo's own server code can still update them.

## Nobody can delete rate sheet name change records through the API

Lines 29-29

Asked as: remove a rename record; delete name change log; erase rate sheet rename; api delete rename not allowed

No one can delete a rate sheet name change record (rate_sheet_name_changes) through the API, except a PocketBase superuser. Turbo's own server code can still delete them.

## Fields required on a rate sheet name change record

Lines 35-40

Asked as: required fields rename record; what to fill in rename; missing rate sheet name change field; renamed_by required; rename record validation error

Each rate sheet name change record requires the rate sheet (rate_sheet), the old name (old_name), the new name (new_name) and the person who renamed it (renamed_by, a user). Turbo sets the created time (created) by itself when the record is created.
