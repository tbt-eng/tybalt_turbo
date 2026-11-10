<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/copy_timesheet_entries.go","sha256":"5fe54fe4df2423dc04d7cc5ac0f09bbd60007118582b5bf422e5158a70b507e9","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/copy_timesheet_entries.go

## Copy to next week needs time editing and time claim

Lines 25-30

Before copying a time sheet's entries to the next week, Turbo runs requireTimeEditing and requireTimeClaim for the user. If either check fails, the copy does not run. The code for these checks is in other files, so their rules are not stated here.

## Time sheet to copy must exist

Lines 38-46

When a user copies entries from a time sheet that Turbo cannot find, Turbo refuses with "time sheet <id> not found".

## Only the owner can copy a time sheet's entries

Lines 48-54

A user who does not own the source time sheet cannot copy its entries. Turbo refuses with "you are not the owner of this time sheet".

## Target week is source week ending plus 7 days

Lines 56-64

Turbo sets the target week ending (week_ending) to the source time sheet's week ending plus 7 days. If the source week ending cannot be read, Turbo refuses with "invalid week_ending on source time sheet" and the parse error.

## Copy refused if next week's time sheet exists

Lines 66-83

If the user already has a time sheet for the target week, Turbo refuses the copy with "a time sheet already exists for the next week".

## Copy refused if next week already has time entries

Lines 85-102

If the user already has time entries for the target week, even without a time sheet, Turbo refuses the copy with "time entries already exist for the next week".

## Cannot copy from an empty time sheet

Lines 104-121

If the source time sheet has no time entries for the user, Turbo refuses with "there are no time entries to copy from this time sheet". Turbo reads the entries ordered by date.

## Each entry is copied 7 days later

Lines 123-129

Turbo copies every source entry with a 7-day date offset using copyTimeEntryWithDayOffset. The copy runs inside one transaction. If any entry fails, Turbo returns that error and the whole copy fails. The rules for what each copy keeps are in another file.

## Copied entry must land in the target week

Lines 130-136

If a copied entry's week ending is not the target week ending, Turbo refuses with "a source time entry date does not match the source time sheet week". The transaction fails, so Turbo keeps no copies.

## Successful copy response

Lines 140-166

After a successful copy, Turbo replies "Time sheet entries copied to next week". The reply includes the number of entries copied, the new record IDs and the target week ending. Turbo creates only the entries and not a time sheet for the new week.

## Copy errors without a status default to bad request

Lines 149-164

Turbo shows a copy error as its message and its code. If no status was set for the error, Turbo treats it as a bad request. Errors from entry hooks (HookError) keep their own status.
