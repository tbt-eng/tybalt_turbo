<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/computersLoad.ts","sha256":"76edfb4b67a5cc18b047e33cc179bae0f1297f747801980804d0eb0157b522e9","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/computersLoad.ts

## Computers pages need the computers or admin claim

Lines 9-16

Asked as: computers page blocked; can't open computers page; computers claim required; access denied computers; no permission computers

Turbo stops a person from opening a Computers page unless they hold the computers or admin claim (checked with canViewComputers). Otherwise it refuses with the message "The computers or admin claim is required to see computers." The server also refuses the data.

## Computer sign-in and raw login lists load all matching rows

Lines 18-26

Asked as: sign-in list empty; raw logins not loading; all logins showing; old sign-ins missing; search returns everything

The sign-ins and rejected reports (raw logins) lists load every row that matches the search text in the URL (q). Turbo trims the search text, and an empty search loads all rows. The server sends all rows because it deletes rows after 30 days. A person without the computers or admin claim cannot load these lists.

## Computer list search runs after typing pauses for 300 ms

Lines 28-51

Asked as: search lags while typing; search fires too early; typing delay in search; spaces trigger search; search not updating after typing

On the sign-ins and raw logins lists, Turbo runs a new search 300 milliseconds after typing stops. A change of only leading or trailing spaces loads nothing, because Turbo compares the trimmed text with the last search. The search changes the URL without adding a history entry and keeps focus and scroll position.
