<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/absorbHistory.ts","sha256":"1322cd299b822c3303fb6887482f45e9f40da1c8dc586b3ae41faae3995949ff","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/absorbHistory.ts

## Absorb history is kept for 90 days

Lines 1-3

Asked as: how long is absorb history kept; absorb history retention period; when does absorb history get deleted; old absorb records expire; absorb log 90 days

Absorb history lists the committed absorbs, with who absorbed and who committed each one. Turbo keeps each entry for 90 days. The retention itself is set in app/routes/absorb_history.go, not in this file.

## Absorb history shows saved details of each absorb

Lines 6-37

Asked as: who merged the records; absorb details who committed it; what references were moved; merged clients saved copies; absorb history detail view

Each history entry shows the kept record's name when the absorb was committed (target_name). It also shows the names of the absorbed records and the number of references moved (references_moved). It records who absorbed and when, and who committed and when. The detail view of one absorb adds saved copies of the absorbed records and the references that moved, for troubleshooting.

## Absorb history labels for record types

Lines 39-45

Asked as: absorb history collection names; client contacts label in absorb history; rate sheets shown as; vendors name in absorb log; record type label absorb

Absorb history shows these labels for record types: client_contacts as "Client Contacts", clients as "Clients", vendors as "Vendors" and rate_sheets as "Rate Sheets". Any other collection name is shown as it is.
