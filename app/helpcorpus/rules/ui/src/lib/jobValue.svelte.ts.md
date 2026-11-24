<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/jobValue.svelte.ts","sha256":"6671d6a011c0a6fe091b5be56e766e8dcea03d39f459593ecd9e7603ef52f2e0","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/jobValue.svelte.ts

## Project value shown in whole Canadian dollars

Lines 9-14

Asked as: project value in dollars; value shows $0; value cents rounding; CAD value display

Turbo shows a project value in Canadian dollars (CAD) with no cents, rounded to the whole dollar. A value that is not a number shows as $0.

## Approved-by-me list covers 30 days, newest first

Lines 19-36

Asked as: my approvals last 30 days; approved by me list; approvals history newest first; list of jobs I approved

An approver's list of their own approvals (queue approved_by_me) covers the last 30 days and is sorted newest first. Each entry is either a request that became a job (new_job) or a later value change (value_change). The date is the job creation date for a new job and the approval date for a value change.

## Approved values are ordered by approval sequence

Lines 54-59

Asked as: approved values order; value history sequence; which approved value is listed first; approved changes list order

A job's approved project values are listed newest first, by approval sequence (approval_sequence). Only changes with status Approved are in this list.

## Current value comes from the server's mark

Lines 60-62

Asked as: wrong current project value; job shows old value; which value is current; latest value missing from history

The current project value is the approved change that the server marks as current. Turbo does not take the newest one, because a viewer who may see only some approvals can have the latest one left out of the history.

## A job has one pending or rejected value change

Lines 63-64

Asked as: pending value change on job; rejected value change still showing; only one value change per job; stuck value change on job

The job's in-progress value change is the first change whose status is not Approved, so a pending or rejected one. A job has at most one.

## Approval queue and job history load different lists

Lines 66-79

Asked as: approval queue empty for job; job page missing approvals; why waiting approvals list differs; value changes not on job page

The approver's queue loads the changes that wait for approval (queue=approval). A job page loads the value changes of that one job.

## Approve, reject and delete send the displayed version

Lines 93-113

Asked as: approve wrong version; reject sends old version; version changed after approving; delete sent stale document

When a person approves, rejects or deletes a value change, Turbo sends the version and the PDF hash (document_hash) that the page shows. A decision can change the version even when the PDF hash stays the same. A rejection also sends the reason and the job data that the person saw (review_token).

## Value change actions refresh the page

Lines 93-122

Asked as: page not updating after approve; approve button no refresh; reject failed no change shown; history not reloading after delete

After a successful approve, reject or delete, Turbo reloads the history, the attention counts and the page data. If the action fails, Turbo shows the error and does not refresh.

## Value change PDF download file name

Lines 124-143

Asked as: value PDF file name; download value pdf named; PDF will not open; cannot load this PDF

A downloaded value change PDF is named with the job number (or "project" if none), the word value, and the approval date. If the change is not approved yet, the name uses the upload date. Without download, the PDF opens in a new tab. If the PDF cannot load, Turbo shows "Cannot load this PDF."

## Missing job details that block a value change

Lines 146-167

Asked as: why cannot fill value change form; missing job details before value change; job completion date required value; value change form blocked

Before a value change form is filled in, Turbo lists the job details that are missing: an invoicing profile, division hours, a location and a branch. It also lists a project completion date, but only when the job status is Active. A missing authorizing document is not listed, because the server treats the project as PA and approval records it. The server checks the job in full when the change is saved.
