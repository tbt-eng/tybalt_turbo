<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/timesheet_details.go","sha256":"432b3ba17e26b8efa6539951a69e095ca2995f9cf099758ceeac20d9394e2de0","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/timesheet_details.go

## Time sheet details require a signed-in user

Lines 20-23

Asked as: unauthorized error; not signed in; session expired; login required for timesheet

Turbo refuses the time sheet details request when the user is not signed in. The message is "unauthorized".

## Time sheet details need a time sheet id

Lines 25-28

Asked as: id is required; missing timesheet id; no timesheet id; timesheet id blank

Turbo refuses the time sheet details request when no id is given. The message is "id is required".

## Missing or hidden time sheet gives the same message

Lines 30-41

Asked as: time sheet not found; timesheet not found or not authorized; can't find my timesheet; timesheet does not exist

When the time sheet does not exist, Turbo answers "time sheet not found or not authorized". It gives the same message when the user may not view the time sheet. A user cannot tell whether a hidden time sheet exists.

## Details response lists entries newest first with related data

Lines 43-52

Asked as: timesheet entries order; entries newest first; list of time entries; timesheet details entries

The details response holds the time sheet and all its time entries, sorted by date with the newest first. Each entry includes its job, time type, division, category and role.

## Details response includes approver info and reviewer count

Lines 54-69

Asked as: approver info on timesheet; shared reviewer count; how many reviewers; reviewer count missing

The details response also gives the approver information from getTimesheetApproverInfo. It gives the number of shared reviewers (sharedReviewerCount), counted from the time sheet's reviewer records.

## Shared reviewer count counts all reviewer records

Lines 73-90

Asked as: reviewer count includes drafts; count reviewers regardless of state; reviewers counted on draft; number of reviewer records

The shared reviewer count is the number of reviewer records for the time sheet. It does not check the time sheet's state.

## Time sheet owner can always view details

Lines 97-99

Asked as: view my own timesheet; see my timesheet any state; owner view draft timesheet; my unsubmitted timesheet

The person whose time sheet it is (uid) can view its details in any state, submitted or not.

## Admin views only committed time sheets

Lines 101-107

Asked as: admin cannot see timesheet; admin can't open uncommitted timesheet; admin view someone's timesheet; why admin locked out timesheet

An admin (the admin claim) can view another person's time sheet details only when the time sheet has been committed (committed is set). An uncommitted time sheet stays hidden from an admin unless another rule applies.

## Approver views only submitted time sheets

Lines 109-111

Asked as: approver can't see timesheet; approver not seeing timesheet; manager can't open timesheet; timesheet after recall hidden

The approver of a time sheet (approver) can view its details only when the time sheet is submitted. Before submission, or after a recall, the approver cannot view it.

## Shared reviewers view submitted time sheets only

Lines 113-122

Asked as: reviewer can't open timesheet; shared reviewer no access; colleague added as reviewer cannot see; reviewer access before submit

A person added as a reviewer of a time sheet (a time_sheet_reviewers record) can view its details only when the time sheet is submitted.

## Committer views submitted and approved time sheets

Lines 124-130

Asked as: committer can't see timesheet; commit role timesheet hidden; timesheet not approved yet hidden; committer access unapproved timesheet

A committer (the commit claim) can view another person's time sheet details only when it is submitted and approved (approved is set). A time sheet that is not yet approved is hidden from a committer.

## Report staff view committed time sheets

Lines 132-138

Asked as: report user can't see timesheet; reporting access timesheet; report role cannot open timesheet; finance can't view timesheet

A person with the report claim can view time sheet details only when the time sheet is committed (committed is set).

## Everyone else cannot view time sheet details

Lines 140-140

Asked as: no access to timesheet; why can't I view this timesheet; permission denied timesheet; not allowed to see timesheet

A user who matches none of the owner, admin, approver, reviewer, committer or report rules cannot view the time sheet details.
