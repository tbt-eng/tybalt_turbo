<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_set_status.go","sha256":"e72144cb5e59a9033cc9e596b46f89d915fd614372364445147feee8ae780d44","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_set_status.go

## Request needs a signed-in user and a job id

Lines 26-34

Turbo refuses a request from a user who is not signed in with "unauthorized". It refuses a request with no job id with "missing job id".

## Only Cancelled or No Bid can be set here

Lines 41-47

The proposal status endpoint accepts only the statuses "Cancelled" and "No Bid". Any other status is refused with "this endpoint only supports setting status to Cancelled or No Bid".

## A comment is required to cancel or mark No Bid

Lines 49-54

Turbo refuses the status change when the comment is empty or only spaces. The message is "a comment is required when setting status to" followed by the requested status (Cancelled or No Bid).

## Note and status change succeed or fail together

Lines 58-132

The client note and the status update happen in one transaction. If either save fails, Turbo saves neither. The errors are "error creating client note" and "error updating job status".

## Missing proposal gives job not found

Lines 60-64

When no job exists with the given id, Turbo refuses with "job not found".

## Only proposals can use this status endpoint

Lines 66-70

Turbo refuses the change when the job number does not start with P. The message is "this endpoint is only for proposals".

## A cancelled proposal cannot be changed again

Lines 72-77

When the proposal's status is already Cancelled, Turbo refuses the change with "this proposal is already cancelled". This also applies when the requested status is Cancelled or No Bid. A proposal with No Bid status can still be changed to Cancelled.

## Who may cancel a proposal or mark it No Bid

Lines 79-91

Holders of the job claim may set the status. The proposal's manager (manager) or alternate manager (alternate_manager) may also set it. Anyone else gets "you are not authorized to update this job".

## Proposal must have a client to change status

Lines 100-104

Turbo refuses the change when the proposal has no client. The message is "this proposal has no client; a client is required to add a note". The authorization check runs before this one.

## Status change records a client note

Lines 106-122

Turbo creates a client note for the proposal's client with the trimmed comment, the user who made the change, and the new status (job_status_changed_to). It then sets the proposal's status to Cancelled or No Bid.
