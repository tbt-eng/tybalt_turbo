<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/deploy_groups.md","sha256":"1a788e80c60ed72d0152c03941f9f29068f62b802cf69909771c838a4c9ee139","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/deploy_groups.md

## Nobody can list deploy groups through the API

Lines 7-9

Asked as: list deploy groups; can't see deploy groups; fetch all deploy groups via API; deploy groups api list

The deploy_groups collection has no list rule. Nobody can list its records through the API, except a PocketBase superuser. Turbo's own server code can still read them.

## Nobody can view one deploy group through the API

Lines 11-13

Asked as: view one deploy group; open a single deploy group; deploy group details not loading; get deploy group by id

The deploy_groups collection has no view rule. Nobody can view a single record through the API, except a PocketBase superuser. Turbo's own server code can still read it.

## Nobody can create deploy groups through the API

Lines 15-17

Asked as: add a deploy group; create new deploy group; can't make deploy group; deploy group creation blocked

The deploy_groups collection has no create rule. Nobody can create a record through the API, except a PocketBase superuser. Turbo's own server code can still create records.

## Nobody can update deploy groups through the API

Lines 19-21

Asked as: edit deploy group; change deploy group settings; can't update deploy group; deploy group changes not saved

The deploy_groups collection has no update rule. Nobody can update a record through the API, except a PocketBase superuser. Turbo's own server code can still update records.

## Nobody can delete deploy groups through the API

Lines 23-25

Asked as: remove deploy group; delete deploy group; deploy group won't go away; get rid of deploy group

The deploy_groups collection has no delete rule. Nobody can delete a record through the API, except a PocketBase superuser. Turbo's own server code can still delete records.

## Deploy group name is required, up to 100 characters

Lines 32-32

Asked as: deploy group missing name; blank deploy group name; name too long deploy group; save fails no name

Each deploy group needs a name (name) of at most 100 characters. A record without a name is not valid.

## Two deploy groups cannot share a name

Lines 38-40

Asked as: duplicate deploy group name; name already taken; same deploy group twice; deploy group name exists

A unique index (idx_deploy_groups_name) on the name field means two deploy groups cannot have the same name.
