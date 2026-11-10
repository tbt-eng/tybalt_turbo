<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/absorb_actions.md","sha256":"02d8bc3e65e268f4f4504c2f10e3947e23945a93d5ac8b9591eeeeafb1c275a4","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/absorb_actions.md

## Who can list and view absorb action records

Lines 7-17

Asked as: cannot see absorb records; who can view absorb actions; absorb list access denied; need absorb permission to view

Only a person who holds the absorb claim can list absorb_actions records and view a single one through the API. The filter checks that one of the signed-in person's claims is named 'absorb'.

## Nobody can create absorb action records through the API

Lines 19-21

Asked as: cannot create absorb action; add absorb record via API; api create absorb fails; new absorb action not saved

The create rule for absorb_actions is empty, so nobody can create a record through the API, except a PocketBase superuser. Turbo's own server code can still create records.

## Nobody can update absorb action records through the API

Lines 23-25

Asked as: cannot edit absorb action; update absorb record api; absorb record change refused; modify absorb action not allowed

The update rule for absorb_actions is empty, so nobody can update a record through the API, except a PocketBase superuser. Turbo's own server code can still update records.

## Only the absorb claim can delete absorb action records

Lines 27-31

Asked as: cannot delete absorb record; remove absorb action api; absorb deletion not permitted; who can erase absorb entries

Only a person who holds the absorb claim can delete an absorb_actions record through the API. The filter checks that one of the signed-in person's claims is named 'absorb'.

## Absorb action field requirements and size limits

Lines 35-43

Asked as: absorb record validation error; absorbed records too large; required fields absorb action; absorb json size limit

Each absorb_actions record requires collection_name, target_id, absorbed_records (absorbed records) and updated_references (updated references). The two JSON fields (absorbed_records, updated_references) hold at most 2,000,000 in size each. The id is 15 characters of lowercase letters and digits.

## Only one absorb action record per collection name

Lines 47-49

Asked as: duplicate absorb collection name; same collection absorbed twice; absorb record already exists; unique collection name error

A unique index on collection_name means two absorb_actions records cannot have the same collection_name.
