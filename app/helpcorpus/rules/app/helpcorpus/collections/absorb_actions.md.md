<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/absorb_actions.md","sha256":"10165e164625cf016db94502c4b047176d72bc8e5b080838855e193d0b8f932a","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/absorb_actions.md

## Who can list or view absorb action records

Lines 7-17

Asked as: who can see absorb records; absorb access view list; admin can only see rate_sheets; absorb history visibility

Staff with absorb access (the 'absorb' claim) can list and view absorb action records (absorb_actions). An admin (the 'admin' claim) can list and view only the records whose collection_name is 'rate_sheets'.

## Nobody can create absorb action records through the API

Lines 19-21

Asked as: cannot create absorb record; api create absorb action fails; insert absorb record forbidden; why can't I add absorb entry

Nobody can create an absorb action record (absorb_actions) through the API, except a PocketBase superuser. Turbo's own server code can still create them.

## Nobody can update absorb action records through the API

Lines 23-25

Asked as: cannot edit absorb action; update absorb record forbidden; api change absorb entry blocked; absorb record edit not allowed

Nobody can update an absorb action record (absorb_actions) through the API, except a PocketBase superuser. Turbo's own server code can still update them.

## Who can delete absorb action records

Lines 27-31

Asked as: delete absorb record; remove absorb action entry; admin cannot delete absorb; who can erase absorb history

Staff with absorb access (the 'absorb' claim) can delete absorb action records (absorb_actions). An admin (the 'admin' claim) can delete only the records whose collection_name is 'rate_sheets'.

## Required fields and size limits on absorb actions

Lines 35-45

Asked as: absorb record validation error; required fields absorb action; absorb json too large; absorbed_by optional field

An absorb action record requires collection_name, target_id, absorbed_records and updated_references. The two JSON fields (absorbed_records, updated_references) are limited to 2000000 bytes each. The proposal and absorbed_by (one user) fields are optional.

## Only one absorb action record for each collection

Lines 49-51

Asked as: duplicate absorb record collection; absorb already exists for collection; second absorb entry same collection; unique absorb collection_name error

A unique index on collection_name means that two absorb action records cannot have the same collection_name. Each collection (such as rate_sheets) has at most one record.
