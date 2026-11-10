<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/timesheet_details.go","sha256":"432b3ba17e26b8efa6539951a69e095ca2995f9cf099758ceeac20d9394e2de0","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/timesheet_details.go

## Time sheet details require a signed-in user

Lines 20-23

Turbo refuses the time sheet details request when the user is not signed in. The message is "unauthorized".

## Time sheet details need a time sheet id

Lines 25-28

Turbo refuses the time sheet details request when no id is given. The message is "id is required".

## Missing or hidden time sheet gives the same message

Lines 30-41

When the time sheet does not exist, Turbo answers "time sheet not found or not authorized". It gives the same message when the user may not view the time sheet. A user cannot tell whether a hidden time sheet exists.

## Details response lists entries newest first with related data

Lines 43-52

The details response holds the time sheet and all its time entries, sorted by date with the newest first. Each entry includes its job, time type, division, category and role.

## Details response includes approver info and reviewer count

Lines 54-69

The details response also gives the approver information from getTimesheetApproverInfo. It gives the number of shared reviewers (sharedReviewerCount), counted from the time sheet's reviewer records.

## Shared reviewer count counts all reviewer records

Lines 73-90

The shared reviewer count is the number of reviewer records for the time sheet. It does not check the time sheet's state.

## Time sheet owner can always view details

Lines 97-99

The person whose time sheet it is (uid) can view its details in any state, submitted or not.

## Admin views only committed time sheets

Lines 101-107

An admin (the admin claim) can view another person's time sheet details only when the time sheet has been committed (committed is set). An uncommitted time sheet stays hidden from an admin unless another rule applies.

## Approver views only submitted time sheets

Lines 109-111

The approver of a time sheet (approver) can view its details only when the time sheet is submitted. Before submission, or after a recall, the approver cannot view it.

## Shared reviewers view submitted time sheets only

Lines 113-122

A person added as a reviewer of a time sheet (a time_sheet_reviewers record) can view its details only when the time sheet is submitted.

## Committer views submitted and approved time sheets

Lines 124-130

A committer (the commit claim) can view another person's time sheet details only when it is submitted and approved (approved is set). A time sheet that is not yet approved is hidden from a committer.

## Report staff view committed time sheets

Lines 132-138

A person with the report claim can view time sheet details only when the time sheet is committed (committed is set).

## Everyone else cannot view time sheet details

Lines 140-140

A user who matches none of the owner, admin, approver, reviewer, committer or report rules cannot view the time sheet details.
