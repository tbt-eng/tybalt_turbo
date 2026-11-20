<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/computer_raw_logins.md","sha256":"b53dc11118c9bcb47be265ca04f82b08ead8e08c0a9963ec5c258bc086428a91","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/computer_raw_logins.md

## Who can list computer raw login records

Lines 7-14

Asked as: can't see login list; who can list logins; computer logins missing from api; list computer login records

A signed-in person who holds the computers claim or the admin claim can list computer_raw_logins records through the API. Everyone else cannot list them.

## Who can view one computer raw login record

Lines 16-23

Asked as: view one login record; can't open a login entry; single computer login details; who can see a login

A signed-in person who holds the computers claim or the admin claim can view a single computer_raw_logins record through the API. Everyone else cannot view one.

## Nobody can create computer raw login records through the API

Lines 25-27

Asked as: can't add login record; create computer login via api; api login entry won't save; insert login record

No one can create a computer_raw_logins record through the API, except a PocketBase superuser. Turbo's own server code can still create them.

## Nobody can update computer raw login records through the API

Lines 29-31

Asked as: edit login record; update computer login api; change a login entry; api won't modify logins

No one can update a computer_raw_logins record through the API, except a PocketBase superuser. Turbo's own server code can still update them.

## Nobody can delete computer raw login records through the API

Lines 33-35

Asked as: remove login record; delete computer login api; erase a login entry; api can't delete logins

No one can delete a computer_raw_logins record through the API, except a PocketBase superuser. Turbo's own server code can still delete them.
