<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/absorb_actions.md","sha256":"7092721f7a58f36713d70ee3b5a60d32015b23c606aa025525b3fd5a924afbe6","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/absorb_actions.md

## Who can list and view absorb action records

Lines 7-17

Asked as: who can view absorb records; absorb action access rights; can't see absorb actions list; admin see rate_sheets absorb

Staff with the absorb claim can list and view absorb action records (absorb_actions). Staff with the admin claim can also list and view them, but only the records whose collection_name is 'rate_sheets'. Other people cannot.

## Nobody can create absorb action records through the API

Lines 19-21

Asked as: create absorb record via API; add absorb action manually; API cannot create absorb; insert absorb action record

No rule allows creating an absorb action record (absorb_actions) through the API. Only a PocketBase superuser can. Turbo's own server code can still create them.

## Nobody can update absorb action records through the API

Lines 23-25

Asked as: edit absorb action record; update absorb action API; can't change absorb record; absorb action edit forbidden

No rule allows updating an absorb action record (absorb_actions) through the API. Only a PocketBase superuser can. Turbo's own server code can still update them.

## Who can delete absorb action records

Lines 27-31

Asked as: delete absorb action record; remove absorb entry; who may delete absorb; admin delete rate_sheets absorb

Staff with the absorb claim can delete absorb action records (absorb_actions). Staff with the admin claim can delete only the records whose collection_name is 'rate_sheets'. Other people cannot.

## Required fields of an absorb action record

Lines 35-44

Asked as: absorb record required fields; absorb action missing field; absorbed_records size limit; target_id absorb required

Each absorb action record needs a collection_name, a target_id, absorbed_records and updated_references. The proposal field is optional. The two JSON fields (absorbed_records, updated_references) are each limited to 2,000,000 bytes (maxSize).

## Only one absorb action record per collection

Lines 46-50

Asked as: duplicate absorb record collection; two absorb entries same collection; unique collection_name absorb error; absorb record already exists

A unique index on collection_name means that two absorb action records cannot name the same collection. Each collection has at most one record.
