<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/client_notes.go","sha256":"1f352d8db16ea0c31bc5b87d4747d692565b1404aa98cf61505b9725287023c5","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/client_notes.go

## Client note creation requires a signed-in user

Lines 25-36

Asked as: note without login; create client note not signed in; authentication required client note; anonymous client note error

Turbo refuses to create a client note when the request has no signed-in user. The message is "authentication required to create client note" (code authentication_required on uid, "authentication required").

## Client note author is set to the signed-in user

Lines 38-38

Asked as: wrong author on client note; note shows someone else's name; change note author; author overwritten on note

When a client note is created, Turbo sets its author (uid) to the signed-in user. Any uid sent with the request is replaced.

## Client note requires a client

Lines 48-59

Asked as: client note no client; save note without client; client missing from note; note needs a customer

Turbo refuses a client note that has no client. The message is "client is required" (code required on client).

## Client note needs a job or the not-applicable flag

Lines 61-76

Asked as: note without job; job required for note; can't save note no job; mark note not applicable

A client note with no job must be marked job not applicable (job_not_applicable). Otherwise Turbo refuses it with "job is required unless marked not applicable" on job and "job must be selected or marked not applicable" on job_not_applicable.

## Note marked not applicable skips job checks

Lines 78-80

Asked as: skip job on note; note job not applicable; no job checks on note; job N/A note saves

When a client note has no job and is marked job not applicable, Turbo accepts it without checking the job or job_status_changed_to. A job_status_changed_to value sent with it is not checked.

## Client note job must exist

Lines 82-94

Asked as: job not found on note; note job missing; wrong job number note; job doesn't exist note

Turbo refuses a client note whose job cannot be found. The message is "job not found" (code not_found on job).

## Client note job must belong to the client

Lines 96-110

Asked as: job wrong client note; note job different client; job owner mismatch note; job not for this client

The note's job must have the selected client as its client, or as its job owner (job_owner). Otherwise Turbo refuses it with "job must belong to the selected client" (code job_client_mismatch).

## Job status change on a note is limited by job type

Lines 115-158

Asked as: change job status from note; note status change refused; proposal status on note; project status note error

When a client note sets job_status_changed_to, the value must be allowed for the job type. A job whose number starts with P is a proposal, and the allowed values are "No Bid" and "Cancelled". Any other job is a project, and no value is allowed. A value that is not allowed is refused with "job_status_changed_to value '<value>' is not valid for proposals" (or "projects"), code invalid_for_job_type.

## Empty job status change is not checked

Lines 116-116

Asked as: blank job status note; no status change note; empty status on note; note saves without status

When job_status_changed_to is empty, Turbo skips the status check and accepts the note once the client and job checks pass.
