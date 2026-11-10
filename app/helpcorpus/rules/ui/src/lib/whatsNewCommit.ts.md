<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/whatsNewCommit.ts","sha256":"9a73733c3342158bdacbd58b18a596bfc56dd2c7cc661835022907c6f9f4d74c","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/whatsNewCommit.ts

## What's new finds the latest commit from the first week heading

Lines 3-5

Turbo reads the first heading in the What's new page that matches "## ... (to <commit>)". The commit is 8 to 40 lowercase hex characters (a-f, 0-9). It takes this as the latest commit (latestWhatsNewCommit). If no heading matches, it returns no commit (null).

## First visit to What's new marks only the newest week

Lines 11-11

When this browser has not seen any commit before (seen is null), Turbo marks only the newest week as new. Older weeks are not marked.

## Weeks newer than the last seen commit are marked new

Lines 12-13

When the last seen commit is on the page, Turbo marks as new only the weeks listed before it, which are the newer ones. The week of the seen commit is not marked. If the seen commit is the newest, no week is marked.

## A seen commit missing from the page marks every week new

Lines 12-13

When the commit this browser saw last is no longer on the What's new page, Turbo marks every week on the page as new.
