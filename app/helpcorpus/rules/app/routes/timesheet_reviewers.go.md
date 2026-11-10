<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/timesheet_reviewers.go","sha256":"55c777a1dc9d59871eaea90bf820e2624f3dd240177d945af58ac4ac433560fb","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/timesheet_reviewers.go

## Who can see a time sheet's reviewers

Lines 10-30

Asked as: who can see reviewers; reviewer list empty; can't see who reviews my timesheet; hidden reviewers on timesheet

Turbo lists the reviewers of a time sheet only to the time sheet's owner (uid), its approver, or one of its reviewers. For anyone else the reviewer list comes back empty, with no refusal message.

## Reviewer name may be blank without a profile

Lines 16-20

Asked as: reviewer name blank; reviewer shows no name; missing reviewer name; reviewer without profile

Each reviewer's surname and given name come from the reviewer's profile. A reviewer who has no profile is still listed, with the name fields empty.

## Reviewer list order: reviewed first

Lines 29-29

Asked as: reviewer list order; sort reviewers; reviewers sorted wrong; which reviewer listed first

Turbo sorts the reviewers of a time sheet by their reviewed value, highest first (ORDER BY reviewed DESC).

## Reviewer list fetch failure message

Lines 62-66

Asked as: failed to fetch reviewers; reviewers won't load; reviewer list error; can't load reviewers

If Turbo cannot fetch the reviewers, it shows the error "failed to fetch reviewers: " followed by the cause.

## Owner name falls back to the user account name

Lines 69-77

Asked as: owner name wrong; timesheet owner name missing; owner name shows username; owner name no profile

The owner name on the reviewer summary is the owner's profile given name and surname joined by a space. If the owner has no profile, Turbo uses the name on the user account (users.name).

## Summary shows owner and week ending to any caller

Lines 69-100

Asked as: timesheet owner name shown; week ending shows for anyone; summary visible without access; who owns this timesheet

The summary always carries the owner's name and the time sheet's week ending. Turbo looks these up by time sheet id only, with no check that the caller is the owner, approver or a reviewer. A caller who is not allowed to see reviewers still gets the owner name and week ending, with an empty reviewer list.

## Time sheet metadata fetch failure message

Lines 85-91

Asked as: failed to fetch timesheet metadata; timesheet details won't load; week ending missing; timesheet not found error

If Turbo cannot read the time sheet's owner name and week ending (for example, the time sheet does not exist), it shows the error "failed to fetch time sheet metadata: " followed by the cause.
