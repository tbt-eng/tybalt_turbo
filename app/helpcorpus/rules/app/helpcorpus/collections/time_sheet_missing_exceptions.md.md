<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/time_sheet_missing_exceptions.md","sha256":"19d05b21dcca5e31ae101e25633c7eb97803be5194121adfa3ea75659719010f","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/time_sheet_missing_exceptions.md

## Nobody can list missing time sheet exceptions through the API

Lines 7-9

Asked as: list missing timesheet exceptions; can't see exceptions in API; get all timesheet exceptions; API returns nothing for exceptions

The time_sheet_missing_exceptions collection has no list rule, so nobody can list its records through the API. Only a PocketBase superuser can. Turbo's own server code can still read them.

## Nobody can view one missing time sheet exception through the API

Lines 11-13

Asked as: view one timesheet exception; open a single exception record; fetch exception by id; can't read exception via API

The time_sheet_missing_exceptions collection has no view rule, so nobody can view a single record through the API. Only a PocketBase superuser can. Turbo's own server code can still read it.

## Nobody can create missing time sheet exceptions through the API

Lines 15-17

Asked as: create timesheet exception; add missing timesheet exception; API rejects new exception; post exception record

The time_sheet_missing_exceptions collection has no create rule, so nobody can create a record through the API. Only a PocketBase superuser can. Turbo's own server code can still create records.

## Nobody can update missing time sheet exceptions through the API

Lines 19-21

Asked as: edit timesheet exception; change exception via API; update exception record; modify missing timesheet entry

The time_sheet_missing_exceptions collection has no update rule, so nobody can update a record through the API. Only a PocketBase superuser can. Turbo's own server code can still update records.

## Nobody can delete missing time sheet exceptions through the API

Lines 23-32

Asked as: delete timesheet exception; remove exception record; deleting user removes exceptions; erase missing timesheet exception

The time_sheet_missing_exceptions collection has no delete rule, so nobody can delete a record through the API. Only a PocketBase superuser can. Turbo's own server code can still delete records. Deleting the user (uid) also deletes that person's exceptions (cascadeDelete).

## Each exception needs a user and a creator

Lines 32-34

Asked as: exception missing user; created_by required; exception needs creator; save exception without user

A missing time sheet exception requires one user (uid) and one creator (created_by), each a single record from users.

## Week ending must be a date in YYYY-MM-DD form

Lines 33-33

Asked as: week ending date format; wrong date for week ending; exception date invalid; week ending validation error

The week ending (week_ending) of a missing time sheet exception is required. It must match the pattern of four digits, two digits and two digits separated by hyphens, such as 2025-01-31.

## One exception per user and week ending

Lines 40-42

Asked as: duplicate timesheet exception; exception already exists for week; same week exception twice; two exceptions same user

Two missing time sheet exceptions cannot have the same user (uid) and week ending (week_ending). A unique index enforces this.
