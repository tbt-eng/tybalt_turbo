<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/client_notes.go","sha256":"1f352d8db16ea0c31bc5b87d4747d692565b1404aa98cf61505b9725287023c5","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/client_notes.go

## Client note creation requires a signed-in user

Lines 25-36

Turbo refuses to create a client note when the request has no signed-in user. The message is "authentication required to create client note" (code authentication_required on uid, "authentication required").

## Client note author is set to the signed-in user

Lines 38-38

When a client note is created, Turbo sets its author (uid) to the signed-in user. Any uid sent with the request is replaced.

## Client note requires a client

Lines 48-59

Turbo refuses a client note that has no client. The message is "client is required" (code required on client).

## Client note needs a job or the not-applicable flag

Lines 61-76

A client note with no job must be marked job not applicable (job_not_applicable). Otherwise Turbo refuses it with "job is required unless marked not applicable" on job and "job must be selected or marked not applicable" on job_not_applicable.

## Note marked not applicable skips job checks

Lines 78-80

When a client note has no job and is marked job not applicable, Turbo accepts it without checking the job or job_status_changed_to. A job_status_changed_to value sent with it is not checked.

## Client note job must exist

Lines 82-94

Turbo refuses a client note whose job cannot be found. The message is "job not found" (code not_found on job).

## Client note job must belong to the client

Lines 96-110

The note's job must have the selected client as its client, or as its job owner (job_owner). Otherwise Turbo refuses it with "job must belong to the selected client" (code job_client_mismatch).

## Job status change on a note is limited by job type

Lines 115-158

When a client note sets job_status_changed_to, the value must be allowed for the job type. A job whose number starts with P is a proposal, and the allowed values are "No Bid" and "Cancelled". Any other job is a project, and no value is allowed. A value that is not allowed is refused with "job_status_changed_to value '<value>' is not valid for proposals" (or "projects"), code invalid_for_job_type.

## Empty job status change is not checked

Lines 116-116

When job_status_changed_to is empty, Turbo skips the status check and accepts the note once the client and job checks pass.
