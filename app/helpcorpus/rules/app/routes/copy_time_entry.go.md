<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/copy_time_entry.go","sha256":"4de9733a71aee2762e0bf8923528371cae9982133602a0aed1382f4db609e43d","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/copy_time_entry.go

## Copy errors with an invalid date or save failure

Lines 22-34

If the original's date is not in YYYY-MM-DD format, the week ending cannot be generated, or the save fails, the copy fails as a server error. Turbo shows the message of the error, such as "invalid date format on original record", "error generating week ending" or "error saving copied record".

## Copy moves the entry forward one day

Lines 45-70

The copy action creates a new time entry dated one day after the original. Turbo calculates the week ending (week_ending) of the new date. It does not use the original's week ending.

## Copy duplicates fields except identity and date fields

Lines 50-61

The new time entry takes every field value from the original except id, created, updated, date, week_ending and tsid. The new entry is not bundled into a time sheet because tsid is not copied. The owner (uid) is set to the signed-in user.

## Copied entry goes through normal time entry processing

Lines 72-84

Before saving, Turbo runs the time entry processing hook (ProcessTimeEntry) on the copy. If the hook rejects it, the copy fails with the hook's error and nothing is saved.

## Copying a time entry requires time editing to be allowed

Lines 105-110

Before Turbo copies a time entry, it calls requireTimeEditing and requireTimeClaim. The copy stops with their error if time editing is not allowed or the user lacks the time claim. This file does not show those checks' messages or conditions.

## Copy is all-or-nothing

Lines 118-157

The copy runs in a single transaction. If any step fails, no new time entry is kept.

## Copying an unknown time entry is refused

Lines 121-129

If the time entry id does not exist, Turbo refuses the copy. It shows the message "time entry <id> not found".

## Only the owner can copy a time entry

Lines 131-138

Turbo refuses to copy a time entry whose owner (uid) is not the signed-in user. It shows the message "you are not the owner of this time entry".

## Bundled time entries cannot be copied

Lines 140-147

Turbo refuses to copy a time entry that already belongs to a time sheet (its tsid is not empty). It shows the message "cannot copy a time entry that is already bundled into a time sheet".

## Copy refusal shows an error message and code

Lines 159-169

When Turbo refuses to copy a time entry, it shows the error message and a code. The codes are record_not_found, unauthorized, entry_already_bundled, invalid_date, invalid_week_ending and error_saving_record.

## Copy shows time entry validation errors from the hooks

Lines 170-173

The copy goes through the normal time entry processing (ProcessTimeEntry). If that processing refuses the copied entry with a hook error (HookError), Turbo shows that error as it is. It does not replace it with a copy-specific message.

## Successful copy confirmation

Lines 177-180

After a successful copy, Turbo shows the message "Time entry copied to tomorrow" and returns the id of the new record.
