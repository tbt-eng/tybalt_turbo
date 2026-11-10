<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/branches.md","sha256":"dae9fed6dafaf376d6c8d02ed0cb0c123a62e522d5ff43829911bf0eb1f9be21","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/branches.md

## Who can list or view branches

Lines 9-17

Asked as: who can see branches; branch list access; view branch records; branches visible to staff

Any signed-in person can list branches and view a single branch record. A request with no signed-in user cannot (list and view rule: @request.auth.id != "").

## Nobody can create, update or delete branches through the API

Lines 19-29

Asked as: cannot edit branch; add a new branch; delete branch from api; branch changes not saved

The branches collection has no create, update or delete rule. Nobody can create, update or delete a branch through the API except a PocketBase superuser. Turbo's own server code can still change branches.

## Branch code and name must be unique and at least 2 characters

Lines 35-48

Asked as: duplicate branch code; branch name already taken; branch code too short; branch name rules

Each branch needs a code (code) and a name (name), and each must have at least 2 characters. Two branches cannot have the same code, and two branches cannot have the same name (unique indexes).
