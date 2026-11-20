<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/computer_logins.md","sha256":"2d741452129fbee231f3eac7d9e273a574670e6bfe19c4e1858c7a918830b046","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/computer_logins.md

## Who can list computer login records

Lines 7-14

Asked as: list computer logins; who can see login records; computer login list refused; view all machine sign-ins

Only a signed-in person who holds the computers claim or the admin claim can list computer login records (computer_logins) through the API. Everyone else is refused.

## Who can view one computer login record

Lines 16-23

Asked as: view one login record; open a single computer login; can't see computer sign-in details; single login entry access denied

Only a signed-in person who holds the computers claim or the admin claim can view a single computer login record (computer_logins) through the API. Everyone else is refused.

## Nobody can create computer logins through the API

Lines 25-27

Asked as: add a computer login; create login record via API; insert computer sign-in; API login record creation blocked

Nobody can create a computer login record (computer_logins) through the API, except a PocketBase superuser. Turbo's own server code can still create them.

## Nobody can update computer logins through the API

Lines 29-31

Asked as: edit computer login record; change a login entry; update sign-in record API; login record won't save changes

Nobody can update a computer login record (computer_logins) through the API, except a PocketBase superuser. Turbo's own server code can still update them.

## Nobody can delete computer logins through the API

Lines 33-35

Asked as: remove computer login entry; delete machine sign-ins; erase login history; can't delete login records

Nobody can delete a computer login record (computer_logins) through the API, except a PocketBase superuser. Turbo's own server code can still delete them.

## Deleting a computer deletes its login records

Lines 42-42

Asked as: computer deleted login history gone; login records after removing computer; what happens to logins when computer deleted; delete computer also removes logins

Each computer login record (computer_logins) must belong to one computer (computer). When that computer is deleted, Turbo also deletes all its login records (cascadeDelete).
