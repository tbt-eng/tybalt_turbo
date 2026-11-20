<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/computersLoad.ts","sha256":"57fd271c1bf6432e52ec5152c88e49b3870d929c2aeb54c057b4adf5ffde8dd7","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/computersLoad.ts

## Computers pages need the computers or admin claim

Lines 15-22

Asked as: computers page blocked; no access to computers; computers claim required; can't see computers page; admin claim needed

Turbo stops a Computers page for a person who has neither the computers claim nor the admin claim. The page shows the message "The computers or admin claim is required to see computers." The server also refuses the data (the check is canViewComputers).

## Computer sign-in and rejected report lists load one page at a time

Lines 24-35

Asked as: sign-ins list paging; rejected logins next page; page past the end; computer sign-ins pagination; raw logins page size

The sign-ins list and the rejected reports (raw logins) list load from the server one page at a time. The page, the page size (perPage) and the search text (q) come from the URL. The access check runs before the list loads. The server moves a page past the end back to the last page.

## Typing in a computer list search waits 300 ms

Lines 37-58

Asked as: search box lag; computer list search slow; sign-ins search delay; raw logins search reloads; typing search wait

When the search text of the sign-ins list or the raw logins list changes, Turbo waits 300 milliseconds after the last keystroke before it reloads the list. A change of page or page size reloads at once. Turbo trims spaces from the search text to see whether it changed, so a change of only leading or trailing spaces counts as no change and reloads at once.
