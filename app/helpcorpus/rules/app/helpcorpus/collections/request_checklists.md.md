<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/request_checklists.md","sha256":"723a507dfe1bb7b032a33c45d98dd2b1dbcf9a15dd0d605fee503f7be496685a","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/request_checklists.md

## Nobody can list request checklists through the API

Lines 7-9

Asked as: list checklists api; can't see request checklists; fetch all checklists; api returns no checklists

The request_checklists collection has no list rule. Nobody can list its records through the API, except a PocketBase superuser. Turbo's own server code can still read them.

## Nobody can view one request checklist through the API

Lines 11-13

Asked as: view one checklist api; open single checklist record; api checklist forbidden; get checklist by id

The request_checklists collection has no view rule. Nobody can view a single record through the API, except a PocketBase superuser. Turbo's own server code can still read it.

## Nobody can create request checklists through the API

Lines 15-17

Asked as: create checklist through api; add new checklist api; api post checklist denied; insert checklist record

The request_checklists collection has no create rule. Nobody can create a record through the API, except a PocketBase superuser. Turbo's own server code can still create records.

## Nobody can update request checklists through the API

Lines 19-21

Asked as: edit checklist api; update checklist via api; api change checklist blocked; modify checklist record

The request_checklists collection has no update rule. Nobody can update a record through the API, except a PocketBase superuser. Turbo's own server code can still update records.

## Nobody can delete request checklists through the API

Lines 23-25

Asked as: delete checklist api; remove checklist via api; api checklist deletion not allowed; erase checklist record

The request_checklists collection has no delete rule. Nobody can delete a record through the API, except a PocketBase superuser. Turbo's own server code can still delete records.

## Checklist version is a whole number of at least 1

Lines 32-32

Asked as: checklist version zero; version must be number; invalid checklist version; decimal version checklist

The version field of a request checklist is required. It must be a whole number (onlyInt) and at least 1.

## Checklist questions are required

Lines 33-33

Asked as: checklist questions missing; blank checklist questions; save checklist without questions; empty questions field

The questions field of a request checklist is required. It holds JSON.

## Checklist scope allows only approval or creation

Lines 36-36

Asked as: checklist scope invalid; wrong checklist scope; scope value rejected; checklist type other than approval

The scope field of a request checklist is required. It holds exactly one value: "approval" or "creation".

## Each scope can have each checklist version only once

Lines 38-42

Asked as: duplicate checklist version; same version twice; checklist version already exists; repeat scope and version

A unique index (idx_request_checklists_scope_version) covers scope and version together. Two request checklists cannot have the same scope and the same version. The same version number can exist once for "approval" and once for "creation".
