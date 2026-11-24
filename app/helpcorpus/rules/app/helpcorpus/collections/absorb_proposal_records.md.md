<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/absorb_proposal_records.md","sha256":"dc376465a85a9b149706c2598cb815f2acb2b16bf44e908267ac3442eddf644b","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/absorb_proposal_records.md

## Nobody can list absorb proposal records through the API

Lines 7-9

Asked as: list absorb proposal records API; can't see absorb proposals via API; API returns no proposal records; who can read absorb records

The list rule for absorb_proposal_records has no filter, so nobody can list these records through the API. Only a PocketBase superuser can. Turbo's own server code can still read them.

## Nobody can view one absorb proposal record through the API

Lines 11-13

Asked as: view single absorb record API; open one absorb proposal record; API 404 absorb record; fetch one proposal record

The view rule for absorb_proposal_records has no filter, so nobody can view a single record through the API. Only a PocketBase superuser can. Turbo's own server code can still read them.

## Nobody can create absorb proposal records through the API

Lines 15-17

Asked as: create absorb record through API; can't add absorb proposal record; API rejects new absorb record; post absorb proposal via API

The create rule for absorb_proposal_records has no filter, so nobody can create these records through the API. Only a PocketBase superuser can. Turbo's own server code can still create them.

## Nobody can update absorb proposal records through the API

Lines 19-21

Asked as: edit absorb proposal record API; update absorb record fails; API won't change proposal record; modify absorb record via API

The update rule for absorb_proposal_records has no filter, so nobody can update these records through the API. Only a PocketBase superuser can. Turbo's own server code can still update them.

## Nobody can delete absorb proposal records through the API

Lines 23-25

Asked as: delete absorb proposal record; absorb record vanished after proposal deleted; remove absorb record API; deleting proposal removes its records

The delete rule for absorb_proposal_records has no filter, so nobody can delete these records through the API. Only a PocketBase superuser can. Turbo's own server code can still delete them. When the parent absorb proposal (proposal) is deleted, its records are deleted with it (cascadeDelete).

## Record collection must be clients, contacts or vendors

Lines 31-34

Asked as: absorb record wrong collection type; record type not allowed absorb; absorb record missing proposal; absorb record validation error record id

Each absorb proposal record (collection_name) must be one of three values: clients, client_contacts or vendors. It also needs a proposal and a record_id. The record_name is optional.

## A record can be in only one open absorb proposal

Lines 44-44

Asked as: client already in absorb proposal; can't add client to second proposal; vendor already absorbing duplicate; contact in two open proposals

Two absorb proposal records cannot both be open (open = 1) for the same collection_name and record_id. A client, contact or vendor can therefore be in only one open absorb proposal at a time. Closed records are not limited.
