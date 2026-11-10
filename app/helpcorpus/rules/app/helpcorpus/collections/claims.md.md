<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/claims.md","sha256":"c404d118dd19aaf88be08e4157941ea70cef3382c27fa6c97dba8680639994ac","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/claims.md

## Who can list or view claims

Lines 7-17

Asked as: who can see claims; list claims through api; view claim record; claims not showing for signed out user

Any signed-in person can list the claims records and view one record through the API. A request without a signed-in person cannot (list and view rule: @request.auth.id != "").

## Nobody can create, update or delete claims through the API

Lines 19-29

Asked as: cannot edit claim via api; create claim api forbidden; delete claim not allowed; change a claim record

The claims collection has no create, update or delete rule. Nobody can create, change or delete a claim through the API, except a PocketBase superuser. Turbo's own server code can still do these actions.

## Claim field limits: id, name and description

Lines 33-37

Asked as: claim id length; claim name too short; claim description required; valid characters for claim id

A claim's id (id) is required and has exactly 15 characters, using only lowercase letters and digits. Its name (name) is required and needs at least 2 characters. Its description (description) is required and needs at least 5 characters.

## Two claims cannot have the same name

Lines 43-45

Asked as: duplicate claim name; claim name already taken; same claim name twice; unique claim name error

A unique index on the claim name (name) means that two claims cannot share a name.
