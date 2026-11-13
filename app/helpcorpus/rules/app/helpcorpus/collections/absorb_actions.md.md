<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/absorb_actions.md","sha256":"35e6e16bcdeea9b6dd373bbc0dc086b5a8fb79874a32bb0269c9fc5c25070bf9","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/absorb_actions.md

## Who can list or view absorb action records

Lines 7-17

Asked as: who can see absorb records; list absorb actions; can't view absorb action; rate sheets records visible

Staff with the absorb claim can list and view any absorb action record (absorb_actions). Staff with the admin claim can list and view only the records whose collection_name is 'rate_sheets'. Everyone else gets no records.

## Nobody can create absorb action records through the API

Lines 19-21

Asked as: create absorb action api; add absorb record via api; post absorb action; api says forbidden to create absorb

Nobody can create an absorb action record (absorb_actions) through the API, except a PocketBase superuser. Turbo's own server code can still create them.

## Nobody can update absorb action records through the API

Lines 23-25

Asked as: edit absorb action; update absorb record api; change absorb action record; api refuses absorb update

Nobody can update an absorb action record (absorb_actions) through the API, except a PocketBase superuser. Turbo's own server code can still update them.

## Who can delete absorb action records

Lines 27-31

Asked as: delete absorb action; remove absorb record; who can delete absorb records; can't delete absorb action

Staff with the absorb claim can delete any absorb action record (absorb_actions). Staff with the admin claim can delete only the records whose collection_name is 'rate_sheets'. Everyone else cannot delete.

## Absorb action record fields and size limits

Lines 33-43

Asked as: absorb record required fields; absorb action too large; absorbed_records size limit; updated_references maximum bytes

Each absorb action record requires collection_name, target_id, absorbed_records and updated_references. The two JSON fields (absorbed_records, updated_references) each hold at most 2,000,000 bytes (maxSize 2000000).

## Only one absorb action record for each collection

Lines 45-49

Asked as: duplicate absorb record collection; second absorb action same collection; unique collection_name absorb; absorb record already exists

A unique index on collection_name allows only one absorb action record (absorb_actions) for each collection. A second record for the same collection is refused.
