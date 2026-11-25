<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/absorb_dismissals.md","sha256":"1e93c936f74ff0a1430efadc46e8b3b3ae0d7e08bbc99b5ad95d9b7e86537b7d","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/absorb_dismissals.md

## Nobody can list absorb dismissals through the API

Lines 7-9

Asked as: list absorb dismissals; can't see dismissals list; api returns empty dismissals; who can read dismissals

The absorb_dismissals collection has no list rule. Nobody can list its records through the API, except a PocketBase superuser. Turbo's own server code can still read them.

## Nobody can view one absorb dismissal through the API

Lines 11-13

Asked as: view single dismissal; open one absorb dismissal; dismissal record not visible; api get dismissal forbidden

The absorb_dismissals collection has no view rule. Nobody can view a single record through the API, except a PocketBase superuser. Turbo's own server code can still read them.

## Nobody can create absorb dismissals through the API

Lines 15-17

Asked as: create dismissal via api; api rejects new dismissal; can't add dismissal record; post dismissal forbidden

The absorb_dismissals collection has no create rule. Nobody can create a record through the API, except a PocketBase superuser. Turbo's own server code can still create them.

## Nobody can update or delete absorb dismissals through the API

Lines 19-25

Asked as: edit or remove dismissal; change dismissal through api; delete dismissal record; update dismissal denied

The absorb_dismissals collection has no update rule and no delete rule. Nobody can change or delete a record through the API, except a PocketBase superuser. Turbo's own server code can still do both.

## Absorb dismissal covers only clients or vendors

Lines 29-33

Asked as: dismissal wrong collection type; only clients or vendors; record_ids too many; dismissal validation error

The collection field (collection_name) of an absorb dismissal is required and accepts only one value: "clients" or "vendors". The list of dismissed record IDs (record_ids) is required and has a maximum size of 100000.
