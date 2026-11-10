<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/divisions.md","sha256":"606729f06aa60520aadb651ab15b54bb355dae612f7b5eb944ba012437926786","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/divisions.md

## Signed-in staff can list divisions

Lines 7-11

Asked as: list divisions; see all divisions; divisions api list; who can read divisions

Any signed-in person can list the records of the divisions collection through the API. A request with no signed-in person is refused.

## Signed-in staff can view one division

Lines 13-17

Asked as: view one division; open a division record; division details api; single division lookup

Any signed-in person can view a single division record through the API. A request with no signed-in person is refused.

## Nobody can create a division through the API

Lines 19-21

Asked as: add a new division; create division api; can't make a division; division creation denied

The divisions collection has no create rule, so nobody can create a division through the API, except a PocketBase superuser. Turbo's own server code can still create divisions.

## Nobody can update a division through the API

Lines 23-25

Asked as: edit a division; change division name; division update forbidden; modify division record

The divisions collection has no update rule, so nobody can change a division through the API, except a PocketBase superuser. Turbo's own server code can still update divisions.

## Nobody can delete a division through the API

Lines 27-29

Asked as: remove a division; delete division; division can't be deleted; drop a division

The divisions collection has no delete rule, so nobody can delete a division through the API, except a PocketBase superuser. Turbo's own server code can still delete divisions.

## Division code must be unique

Lines 35-46

Asked as: duplicate division code; division code already taken; division save fails; division name too short

Two divisions can't have the same code (unique index on code). A division needs a code of at least 1 character (code) and a name of at least 2 characters (name).
