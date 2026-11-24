<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/absorb_proposals.md","sha256":"dcd69345c70b0ee4c2805bd7ed1013eedb04ebf7fca0d5455194d0694a2e68bd","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/absorb_proposals.md

## Nobody can list absorb proposals through the API

Lines 7-9

Asked as: cannot list proposals; see absorb proposals list; empty list API; why no proposals showing

The list rule for absorb proposals (absorb_proposals) is empty. Nobody can list the records through the API, except a PocketBase superuser. Turbo's own server code can still read them.

## Nobody can view one absorb proposal through the API

Lines 11-13

Asked as: cannot open one proposal; view single proposal API; proposal record forbidden; cannot see proposal details

The view rule for absorb proposals (absorb_proposals) is empty. Nobody can view a single record through the API, except a PocketBase superuser. Turbo's own server code can still read them.

## Nobody can create absorb proposals through the API

Lines 15-17

Asked as: cannot create proposal; new absorb proposal forbidden; API create proposal denied; how to add proposal

The create rule for absorb proposals (absorb_proposals) is empty. Nobody can create a record through the API, except a PocketBase superuser. Turbo's own server code can still create them.

## Nobody can update absorb proposals through the API

Lines 19-21

Asked as: cannot edit proposal; update proposal denied; change proposal through API; proposal read only

The update rule for absorb proposals (absorb_proposals) is empty. Nobody can change a record through the API, except a PocketBase superuser. Turbo's own server code can still update them.

## Nobody can delete absorb proposals through the API

Lines 23-25

Asked as: cannot delete proposal; remove proposal API; delete absorb proposal forbidden; proposal won't delete

The delete rule for absorb proposals (absorb_proposals) is empty. Nobody can delete a record through the API, except a PocketBase superuser. Turbo's own server code can still delete them.

## Absorb proposal collection type and status values

Lines 31-39

Asked as: proposal status options; what kinds of proposal; merge clients or vendors; proposal type contacts vendors

An absorb proposal (absorb_proposals) is about one of three kinds of record: clients, contacts (client_contacts) or vendors (collection_name). Its status (status) is one of pending, accepted, declined, withdrawn or obsolete.

## Absorb proposal reason and decision comment limits

Lines 38-42

Asked as: reason too long; proposal text character limit; decision comment length; how long can comment be

An absorb proposal needs a reason (reason) of at most 1000 characters. The decision comment (decision_comment) has the same 1000-character limit and is optional.
