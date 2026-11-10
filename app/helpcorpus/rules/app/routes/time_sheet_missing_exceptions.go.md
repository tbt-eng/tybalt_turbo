<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/time_sheet_missing_exceptions.go","sha256":"8c55e55fa24fc6724e0b3ac52fab175ad7fc9c603adfdc55139e38490b3989bd","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/time_sheet_missing_exceptions.go

## Who may view time tracking lists

Lines 20-47

A person with the report claim or the commit claim (a committer) may view the time tracking list of ignored time sheets. Anyone else gets the message "you are not authorized to view time tracking".

## Only committers may change missing time sheet exceptions

Lines 49-67

Only a committer (the commit claim) may mark a missing time sheet as ignored or remove that mark. Anyone else, including a person with only the report claim, gets "you are not authorized to update missing time sheet exceptions".

## Week ending is required and must be a valid date

Lines 69-95

A missing time sheet exception request needs a week ending. If it is empty, Turbo refuses with "invalid week ending" and the detail "weekEnding is required". If the value fails the format check (isValidWeekEnding), Turbo refuses with "invalid week ending" and the detail "weekEnding must be in YYYY-MM-DD format".

## User must be given and must exist

Lines 97-125

Marking or unmarking an exception needs a user id. Turbo trims spaces from the id first. If it is empty, Turbo refuses with "invalid uid" and the detail "uid is required". If no user has that id, Turbo refuses with "user not found".

## Ignoring a missing time sheet

Lines 127-175

A committer can mark a user's missing time sheet for a week ending as ignored. Turbo records who marked it (created_by) and returns ignored as true. The checks run in this order: commit claim, week ending, user.

## Marking twice is harmless

Lines 149-175

If a user and week ending are already marked as ignored, marking them again changes nothing and raises no error. Turbo still returns ignored as true.

## Removing an ignored mark

Lines 179-215

A committer can remove the ignored mark for a user and week ending. Turbo returns ignored as false. Removing a mark that does not exist raises no error. The checks run in this order: commit claim, week ending, user.

## List of users whose missing time sheet is ignored

Lines 217-254

For a week ending, Turbo lists users who have an ignored mark, have no time sheet for that week, and have time_sheet_expected set to 1 on their admin profile. A user with no admin profile counts as 0 and is left out. Turbo sorts the list by surname, then given name. A missing name or email shows as an empty string.

## Ignored list needs a valid week ending

Lines 223-226

The list of ignored missing time sheets needs a valid week ending in YYYY-MM-DD format. It refuses an empty or badly formatted value with "invalid week ending".
