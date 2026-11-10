<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/unbundle_timesheet.go","sha256":"7e0f420573701224c1b8da958bb1ca8799b85fdb8630274df222909dc64f3573","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/unbundle_timesheet.go

## Only staff with the time claim can recall a time sheet

Lines 20-22

Turbo runs a claim check (requireTimeClaim) before it recalls (unbundles) a time sheet. If the person lacks the time claim, the check returns its error and the time sheet stays unchanged. This file does not show the claim check's message.

## Recall refuses a time sheet that is not found

Lines 32-39

When the time sheet id in the request does not match a time sheet, Turbo refuses the recall. The message is "error fetching time sheet: ..." followed by the underlying error.

## Only the owner can recall a time sheet

Lines 41-44

Turbo refuses the recall when the time sheet's owner (uid) is not the person asking. The message is "time sheet does not belong to user". This applies to every status, and no one else can recall another person's time sheet.

## An approved time sheet must be rejected before recall

Lines 46-51

A time sheet that has an approved date and no rejected date cannot be recalled. Turbo refuses with "approved time sheets must be rejected before being unbundled". An approved time sheet that has also been rejected passes this check. A submitted time sheet that is not yet approved passes this check as well.

## A committed time sheet cannot be recalled

Lines 53-56

When a time sheet has a committed date, Turbo refuses the recall with "committed time sheets cannot be unbundled". Rejection does not lift this block.

## Recall frees the time entries and deletes the time sheet

Lines 58-81

A successful recall clears the time sheet link (tsid) on the owner's time entries that point to this time sheet. Turbo then deletes the time sheet record. The entries stay and can be edited again. If saving an entry or deleting the sheet fails, Turbo shows "error updating time entry: ..." or "error deleting time sheet: ...". All the steps run in one transaction, so a failure rolls back every change.

## Recall reports success or the error message

Lines 82-86

When the recall succeeds, Turbo shows the message "Time sheet unbundled successfully". When any step above fails, Turbo shows the error text of that step.
