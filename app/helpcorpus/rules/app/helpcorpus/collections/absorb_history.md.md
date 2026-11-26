<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/absorb_history.md","sha256":"8aaa6dc77e11340b53239878e24cc27ab4f8aa0bdca3a0d2463d826062921ebf","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/absorb_history.md

## Nobody can list absorb history records through the API

Lines 7-9

Asked as: list absorb history; can't see absorb history list; api access absorb history; read all absorb records

The absorb history collection (absorb_history) has no list rule. No one can list its records through the API, except a PocketBase superuser. Turbo's own server code can still read them.

## Nobody can view one absorb history record through the API

Lines 11-13

Asked as: view one absorb record; open single absorb history; fetch absorb history entry; absorb record not visible

The absorb history collection (absorb_history) has no view rule. No one can view a single record through the API, except a PocketBase superuser. Turbo's own server code can still read it.

## Nobody can create absorb history records through the API

Lines 15-17

Asked as: create absorb history record; add absorb entry via api; insert absorb history; post absorb record

The absorb history collection (absorb_history) has no create rule. No one can create a record through the API, except a PocketBase superuser. Turbo's own server code can still create records.

## Nobody can update absorb history records through the API

Lines 19-21

Asked as: edit absorb history; change absorb record; update absorb entry api; modify absorb history

The absorb history collection (absorb_history) has no update rule. No one can change a record through the API, except a PocketBase superuser. Turbo's own server code can still update records.

## Nobody can delete absorb history records through the API

Lines 23-25

Asked as: delete absorb history; remove absorb record; erase absorb entry; drop absorb history row

The absorb history collection (absorb_history) has no delete rule. No one can delete a record through the API, except a PocketBase superuser. Turbo's own server code can still delete records.

## Absorb history record fields and size limits

Lines 31-36

Asked as: absorb history fields; absorb record size limit; json field too big; required absorb record columns

Each absorb history record must have a collection name (collection_name), a target ID (target_id), the absorbed records (absorbed_records) and the updated references (updated_references). The two JSON fields (absorbed_records, updated_references) are each limited to 2,000,000 bytes (maxSize). The record ID is 15 characters of lowercase letters and digits.
