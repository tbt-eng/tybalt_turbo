<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/computersLoad.ts","sha256":"b993b74697340c1d2371c2e01ee81f0f0ac060d209c00094733ad586d79e1898","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/computersLoad.ts

## Computers pages need the computers or admin claim

Lines 9-16

Asked as: computers page access denied; no permission for computers; can't open computers page; computers claim missing; need admin claim to see computers

Turbo stops a Computers page for a person who has neither the computers claim nor the admin claim. It shows the message "The computers or admin claim is required to see computers." The server also refuses the data (canViewComputers, requireComputersAccess).

## Sign-ins and rejected reports lists load every matching row

Lines 18-27

Asked as: sign-ins list not all rows; sign-ins not paged; rejected reports show everything; old sign-ins missing 30 days; why only one rejected report

The sign-ins list and the rejected reports (raw logins) list load all rows that match the search in the URL (q). Turbo does not page them. Sign-ins are deleted after 30 days, and Turbo keeps one rejected report per computer name. The person needs the computers or admin claim to load them.

## Search text is trimmed and an empty search loads everything

Lines 24-25

Asked as: empty search shows everything; search box cleared list; trailing spaces in search; blank search results reset; clear search to see all

Turbo trims spaces from the search text (q) before it asks the server. When the trimmed text is empty, Turbo sends no search and the list shows all rows.

## Computer list search waits 300 ms after typing stops

Lines 29-52

Asked as: search lags while typing; list reloads too slow; back button after searching; search adds history entries; typing spaces reloads list

On the sign-ins and rejected reports lists, Turbo changes the URL search 300 ms after the person stops typing, and the list then reloads. It replaces the history entry, so each search does not add a Back step. A change that only adds or removes spaces loads nothing new.
