<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_upsert_api.go","sha256":"d5eab159cbe5b6ead90411a9ddff32756169bd47cf6037f1f31fc2742ac89d8f","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_upsert_api.go

## Allocation errors point to the row that failed

Lines 37-52

Asked as: allocation error on wrong row; which row failed allocation; allocation validation failed; error on division field; allocation row error location

When an allocation row fails validation, Turbo reports the error on that row's division field (allocation_division_N, where N is the row's position starting at 0). The message is "allocation validation failed". The error carries a code and a text for the row, such as required, invalid_hours or duplicate_division.

## Job allocations need at least one row

Lines 54-66

Asked as: job without allocations; no allocation rows; save job with no allocations; at least one allocation; job needs an allocation

When a job is created or saved, Turbo refuses a request with no allocations. The message is "at least one allocation is required" (field allocations, code required).

## Each allocation needs a division and non-negative hours

Lines 69-75

Asked as: division is required; allocation missing division; negative hours error; hours must be zero or more; minus hours not accepted

Turbo refuses an allocation row with no division, with the message "division is required". It refuses negative hours with "hours must be >= 0". Zero hours are allowed.

## A division can appear only once per job

Lines 76-79

Asked as: same division twice; duplicate division on job; division already selected; can't add division twice; repeated division allocation

Turbo refuses a second allocation row with the same division. The message is "division is already selected on row N", where N is the first row that used it.

## Changed allocations mark the job as not imported

Lines 84-109

Asked as: job marked not imported; job shows modified after allocation change; writeback flag allocations changed; job flagged as modified; allocation change not imported

When the allocations differ from the existing ones, Turbo marks the job as modified (not imported) for writeback tracking. Different row counts, a different division, different hours, or duplicate existing rows count as different.

## Saving a job needs a signed-in user

Lines 114-117

Asked as: not signed in saving job; unauthorized job save; job save no user; can't save job logged out; save fails user missing

Saving a job's fields and allocations fails with 'unauthorized' when the request has no signed-in user. Creating a job and checking a proposal fail the same way.

## Job save needs a job id

Lines 119-122

Asked as: missing job id save; job id not sent; save job without id; no job id error

Saving a job's fields and allocations fails with 'missing job id' when the request has no job id.

## Job save rejects an unreadable request body

Lines 124-127

Asked as: invalid JSON body; bad request body job; malformed job request; job save body not parsed; broken json saving job

When a job save request body is not valid JSON, Turbo refuses it with "invalid JSON body". It changes nothing.

## Allocation rows are checked before the job is touched

Lines 129-136

Asked as: bad allocation before save; job unchanged after allocation error; invalid allocations check first; allocation checked before saving; allocation error job untouched

On a job save, Turbo checks the allocation rows (at least one, a division on each row, hours of 0 or more, no repeated division) before it opens the save. If a row is invalid, it returns the allocation error and does not change the job or its allocations.

## Job save is all or nothing

Lines 140-347

Asked as: partial job save; job half saved; rollback job save failure; save failed job changes lost; inactive division blocks save

Turbo saves the job and replaces its allocations in one transaction. If any step fails, such as job validation, an inactive division or an allocation that cannot be saved, nothing is saved. The job and its old allocations stay as they were.

## Saving an unknown job fails

Lines 142-149

Asked as: job not found saving; save nonexistent job; unknown job id error; job_not_found code; job does not exist save

When the job does not exist, Turbo refuses the save with the code job_not_found.

## Who may save a job or its allocations

Lines 152-164

Asked as: manager cannot edit job fields; who can edit job; can't change project fields; alternate manager edit permission; only procurement edits job

A procurement user (the job claim) may save all job fields and allocations. An assigned manager or alternate manager without the claim may save only allocations, or a full job edit on a proposal (number starting with P). Anyone else gets 'only job claim holders may change project fields; assigned managers may save allocations only'.

## Allocation-only save respects the jobs editing switch

Lines 166-180

Asked as: allocations only save disabled; editing switched off allocations; jobs editing disabled error; can't save allocations editing off; ErrJobsEditingDisabled

A save with no job object (allocations only) skips full job validation and does not save the job record. It still fails when job editing is switched off, with the error ErrJobsEditingDisabled.

## Cancelled proposals cannot have allocations changed

Lines 181-184

Asked as: cancelled proposal allocations; change allocations on cancelled; cancelled proposal locked; can't modify cancelled proposal; cancelled proposal edit blocked

An allocations-only save on a proposal whose status is Cancelled fails with 'cancelled proposals cannot be modified'.

## Job number cannot change on save

Lines 187-192

Asked as: change job number; renumber job on save; job number won't change; edit job number error; change number action

A save that sends a number different from the job's current number fails with 'use Change Number to renumber a job'. Only the Change Number action renumbers a job.

## Save ignores unknown and system job fields

Lines 194-206

Asked as: unknown field job save; system fields ignored; collectionId error saving; extra fields in job data; created updated fields rejected

When saving a job, Turbo ignores id, collectionId, collectionName, created, updated and divisions. It also ignores any field the jobs collection does not have.

## Full job save runs job validation; allocation-only does not

Lines 208-220

Asked as: job validation skipped allocations; full save validation; allocations only skips validation; save stopped by validation; business rules on job save

A full job save (a job object is sent, even an empty one) runs the job validation and business rules (ProcessJobCore) and then saves the job. If validation fails, the save stops with that error. An allocations-only save skips both.

## Allocation divisions must exist and be active

Lines 222-251

Asked as: division not found allocation; inactive division allocation; can't use inactive division; division does not exist; selected division not active

Turbo refuses an allocation whose division is not found, with the message "selected division was not found". It refuses an inactive division with "selected division is inactive". This applies when a job is created and when it is saved.

## Unchanged allocations are a no-op

Lines 253-272

Asked as: same allocations no change; save unchanged allocations; allocations identical nothing saved; review kept after save; resave same allocations

An allocations-only save whose allocations match the existing ones (same divisions and hours) changes nothing. It keeps the allocation records and any pending review.

## Job save validation errors are reported per field

Lines 273-297

Asked as: validation failed field messages; error saving job generic; which field is invalid saving; job save field errors; save error message per field

If saving the job fails on field validation, Turbo returns 'validation failed' with each field's message. Other save failures return 'error saving job'.

## Saving allocations replaces the old set

Lines 309-346

Asked as: replace allocations; old allocations deleted; allocations overwritten on save; new allocation rows replace old; allocation save removes existing

Saving deletes all of the job's existing allocations and inserts the submitted rows, all in one transaction. A failure undoes the whole save.

## Successful job save returns the job id

Lines 365-365

Asked as: what does save return; job save response id; save succeeded no id; allocations unchanged success response; returned job id after save

When a job save succeeds, Turbo answers with the id of the job that it saved. A save that changes nothing, because only allocations were sent and they are unchanged, also succeeds and returns the id.

## Job creation needs sign-in

Lines 372-375

Asked as: create job not signed in; new job unauthorized; create job logged out; can't create job no user; sign in to create job

Creating a job fails with 'unauthorized' when the request has no signed-in user.

## Only procurement can create jobs

Lines 382-389

Asked as: who can create jobs; not authorized to create jobs; can't create job as manager; procurement create job only; create job permission denied

Creating a job requires the job claim (procurement). Anyone else gets 'you are not authorized to create jobs'.

## Job creation validates, numbers and saves the job with its allocations

Lines 403-518

Asked as: new job numbering; job number generated on create; create job with allocations; job creation validation; new job saved with allocations

Creating a job runs the job validation (ProcessJobCore), which generates the job number, and then saves the job. Turbo then checks each allocation division exists and is active, and creates the allocations. All of it happens in one transaction.

## Job creation is all or nothing and returns the new id

Lines 403-536

Asked as: job not created partial; create job rollback; new job id returned; failed create leaves nothing; create job transaction failed

Turbo creates the job and its allocations in one transaction. If validation, the job save or any allocation fails, no job is created. When it succeeds, Turbo answers with the id of the new job.

## Create ignores system job fields

Lines 412-422

Asked as: create job extra fields; system fields on new job; divisions field create error; ignored fields new job; id sent when creating job

When creating a job, Turbo ignores id, collectionId, collectionName, created, updated and divisions in the job data.

## Job creation reports field validation errors

Lines 436-459

Asked as: error creating job; new job field errors; create job validation message; why job creation failed fields; create job error text

If saving the new job fails on field validation, Turbo returns 'validation failed' with each field's message. Other failures return 'error creating job'.

## Proposal check needs sign-in and an id

Lines 545-553

Asked as: proposal check not signed in; check proposal missing id; proposal readiness unauthorized; missing proposal id; check proposal no user

Checking a proposal for project readiness fails with 'unauthorized' without a signed-in user and 'missing proposal id' without an id.

## Proposal check only accepts proposals

Lines 556-581

Asked as: proposal not found check; job is not a proposal; number not starting with P; check on non proposal job; proposal check invalid job

The check returns valid false with 'proposal not found' when the job does not exist. It returns 'the specified job is not a proposal' when the number does not start with P.

## Proposal check result

Lines 583-597

Asked as: proposal check result; does proposal check change anything; proposal check errors list; is proposal ready result; proposal valid false list

The proposal check changes nothing. It returns valid true, or valid false with a list of field errors.

## Proposal needs a valid location Plus Code

Lines 606-613

Asked as: location required proposal; plus code missing; proposal location invalid; plus code format proposal; proposal can't create project location

A proposal is not ready for project creation without a location that matches the Plus Code pattern. The message is 'location (Plus Code) is required'.

## Proposal needs a branch and a client

Lines 615-629

Asked as: proposal missing branch; proposal missing client; branch required proposal; client required proposal; proposal not ready no client

A proposal is not ready without a branch ('branch is required') and a client ('client is required').

## Proposal contact must exist and belong to the client

Lines 631-655

Asked as: proposal contact missing; contact not belong to client; contact not found proposal; proposal contact error; wrong contact for client

A proposal needs a contact ('contact is required'). The contact must exist ('specified contact not found') and belong to the proposal's client ('contact must belong to the selected client'). The ownership check runs only when a client is set.

## Proposal manager must be active

Lines 657-678

Asked as: proposal manager missing; manager inactive proposal; manager not active user; proposal manager check failed; proposal owner inactive

A proposal needs a manager ('manager is required'), and the manager must be an active user ('manager must be an active user'). If the status lookup fails, the message is 'failed to verify manager active status'.

## Proposal alternate manager must be active if set

Lines 680-695

Asked as: alternate manager inactive; second manager not active; proposal backup manager error; alternate manager check failed; secondary manager inactive

An alternate manager is optional on a proposal. If set, the person must be active ('alternate manager must be an active user'). If the lookup fails, the message is 'failed to verify alternate manager active status'.

## Proposal needs opening and submission due dates

Lines 697-713

Asked as: proposal dates missing; opening date required; submission due date missing; proposal no due date; proposal no opening date

A proposal needs an opening date ('proposal_opening_date is required') and a submission due date ('proposal_submission_due_date is required').

## Proposal due date cannot precede opening date

Lines 715-723

Asked as: due date before opening; submission date earlier than opening; proposal dates wrong order; due date same as opening; proposal date order error

When both dates are set, a due date before the opening date fails with 'proposal submission due date must be on or after opening date'. The same date is allowed.

## Proposal needs a value or time and materials

Lines 725-735

Asked as: proposal no value; proposal value zero not ready; time and materials proposal; proposal can't create project value; proposal needs value or T&M

For project creation readiness, a proposal needs a proposal value above 0 or time and materials turned on. Otherwise the message is 'proposal must have a proposal value or be marked as time and materials to create a referencing project'. Turbo checks this whatever the proposal's status.

## Proposal description needs 3 characters

Lines 737-744

Asked as: proposal description too short; description characters minimum; proposal description error; short description proposal; description not ready project

A proposal is not ready when its description is shorter than 3 characters. The message is 'description must be at least 3 characters'.
