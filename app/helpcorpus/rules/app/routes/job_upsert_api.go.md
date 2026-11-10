<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_upsert_api.go","sha256":"8fe9b4fee5e4b04f00e0e7b17c8962e0b6bc3c9b7122244037da40743c2c69ea","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_upsert_api.go

## Allocation errors name the row by position

Lines 37-52

When an allocation fails basic validation, Turbo reports "allocation validation failed". The error is attached to the field allocation_division_N, where N is the row's position counting from 0. The editor uses this to show which row is wrong.

## Job save requires at least one division allocation

Lines 54-66

Saving or creating a job with an empty allocations list is refused with 'at least one allocation is required' (code required, field allocations).

## Each allocation needs a division

Lines 69-72

An allocation row with no division is refused with 'division is required'. The error names the row's division field.

## Allocation hours cannot be negative

Lines 73-75

An allocation with hours below 0 is refused with 'hours must be >= 0' (code invalid_hours). Zero hours is allowed.

## A division cannot be used twice in a job's allocations

Lines 76-79

If a division appears on more than one allocation row, Turbo refuses the later row with 'division is already selected on row N', where N is the first row that used it.

## Allocation change detection

Lines 84-109

Turbo treats allocations as changed when the row count differs, a division is added or missing, or the hours for a division differ. Duplicate existing rows for one division also count as changed, so the save normalizes them.

## Job save requires sign-in and a job id

Lines 114-122

Updating a job without a signed-in user is refused as 'unauthorized'. A request with no job id is refused with 'missing job id'.

## Job save rejects an unreadable request body

Lines 124-127

When saving a job, Turbo refuses a request body that is not valid JSON and shows "invalid JSON body". Nothing is changed.

## Allocation checks run before the job is looked up

Lines 129-136

When updating a job, Turbo checks the allocation list (at least one row, a division on each row, hours of 0 or more, no repeated division) before it loads the job or checks who the user is. A bad allocation list is refused with "allocation validation failed", even if the job does not exist or the user may not edit it.

## Updating a job that does not exist

Lines 142-149

If the job id is not found, the update is refused with code job_not_found and the message 'job not found: ...'.

## Who may edit a job and what managers may change

Lines 152-164

A job claim holder may change any project field. A manager or alternate manager without the job claim may save allocations only, or may do a full edit only on a proposal (number starts with P). Anyone else gets 'only job claim holders may change project fields; assigned managers may save allocations only' (job_edit_forbidden).

## Allocation-only save skips full job validation

Lines 166-171

When the request has no job object, Turbo saves only the allocations. It does not validate or save the job record, so old projects with missing billing data and their approved fields stay unchanged. An empty job object counts as a full job save and is validated.

## Allocation-only save needs an active user

Lines 172-180

In an allocation-only save, an inactive user is refused with 'an active user is required' (inactive_user).

## Allocation-only save blocked when job editing is disabled

Lines 181-189

In an allocation-only save, Turbo checks the jobs editing switch. If editing is disabled, it refuses with the error ErrJobsEditingDisabled, which the full job path checks in ProcessJobCore.

## Cancelled proposals' allocations cannot be changed

Lines 190-193

An allocation-only save on a proposal whose status is Cancelled is refused with 'cancelled proposals cannot be modified' (cancelled_proposal_immutable).

## Job number cannot be changed through save

Lines 196-201

If the request's job object has a number that differs from the current number, Turbo refuses with 'use Change Number to renumber a job'. Renumbering must use the Change Number action.

## Job update ignores system fields and unknown fields

Lines 203-215

On update, Turbo ignores id, collectionId, collectionName, created, updated and divisions in the job object. It also ignores any key that is not a field of the jobs collection.

## Full job update runs job validation rules

Lines 217-229

A full job update (not allocation-only) runs the job business rules (ProcessJobCore) before saving. A failure returns that error with its own status.

## Allocation divisions must exist and be active

Lines 231-260

On update, each allocation's division must exist, or the save is refused with 'selected division was not found' (invalid_division). It must also be active, or the save is refused with 'selected division is inactive' (division_not_active).

## Missing hours on stored allocations count as zero

Lines 262-276

When Turbo compares the submitted allocations with the job's stored allocations, a stored row with no hours counts as 0 hours. A submitted 0 for that division therefore counts as no change.

## Unchanged allocation-only save does nothing

Lines 275-281

If an allocation-only save has the same allocations as the current ones, Turbo returns success without replacing them. Allocation IDs and any pending review stay as they are.

## Changed allocations mark the job as not imported

Lines 275-316

When the allocations change, Turbo marks the job as modified (not imported) for writeback tracking. This applies to allocation-only saves too.

## Full job save errors

Lines 282-306

On a full update, field validation failures are returned as 'validation failed' with a message for each field. Other save failures return 'error saving job: ...' (error_saving_job).

## Allocations are replaced as a set

Lines 318-356

Saving deletes all of the job's existing allocations and inserts the submitted ones in the same transaction. The full job update also saves the job in that transaction, and a failure rolls back both.

## A failed job save changes nothing

Lines 357-374

Turbo saves the job fields and replaces the allocations in one transaction. If any step fails, the save returns an error message and a code, or field-level errors. A success returns the job id.

## Only job claim holders can create jobs

Lines 381-398

Creating a job without sign-in is refused as 'unauthorized'. A user without the job claim is refused with 'you are not authorized to create jobs'.

## Job creation validates allocations first

Lines 400-407

Job creation applies the same allocation rules as update: at least one row, a division on each row, hours of 0 or more, and no repeated division. A failure is refused before the job is created.

## Job creation runs job rules and generates the number

Lines 421-468

Creating a job runs the job business rules (ProcessJobCore), which generate the job number. It ignores id, collectionId, collectionName, created, updated and divisions in the job object. A failure returns that error. A save failure returns 'validation failed' with field messages, or 'error creating job: ...' (error_creating_job).

## Job creation checks allocation divisions

Lines 471-526

On creation, each allocation's division must exist ('selected division was not found') and be active ('selected division is inactive'). The job and allocations are created in one transaction, so a failed check leaves no job.

## A failed job creation leaves no job or allocations

Lines 527-545

Turbo creates a new job and its allocations in one transaction. If job rules, saving, division checks or allocation saves fail, the error is returned and nothing is created. A success returns the new job's id.

## Proposal validation check needs sign-in and a proposal

Lines 554-590

The proposal readiness check requires a signed-in user and a proposal id ('missing proposal id'). An unknown id returns valid false with 'proposal not found'. A job whose number does not start with P returns valid false with 'the specified job is not a proposal'.

## Proposal readiness check changes nothing

Lines 592-606

The check validates the proposal without saving it. It returns valid true, or valid false with errors for each field.

## Proposal needs a valid location Plus Code

Lines 615-622

A proposal with no location, or a location that is not a valid Plus Code, fails readiness with 'location (Plus Code) is required'.

## Proposal needs branch and client

Lines 624-638

A proposal without a branch fails readiness with 'branch is required'. A proposal without a client fails with 'client is required'.

## Proposal contact is required and must belong to the client

Lines 640-664

A proposal with no contact fails with 'contact is required'. If the contact is not found, it fails with 'specified contact not found'. If the contact belongs to a different client, it fails with 'contact must belong to the selected client'. The client match is only checked when the proposal has a client.

## Proposal manager must be set and active

Lines 666-687

A proposal without a manager fails readiness with 'manager is required'. An inactive manager fails with 'manager must be an active user'. If the status lookup fails, the message is 'failed to verify manager active status'.

## Alternate manager must be active when set

Lines 689-704

If a proposal has an alternate manager, that person must be active, or readiness fails with 'alternate manager must be an active user'. If the lookup fails, the message is 'failed to verify alternate manager active status'. An empty alternate manager is not checked.

## Proposal needs opening and submission due dates

Lines 706-722

A proposal without an opening date fails readiness with 'proposal_opening_date is required'. A proposal without a submission due date fails with 'proposal_submission_due_date is required'.

## Submission due date cannot precede opening date

Lines 724-732

When both dates are set, a due date before the opening date fails with 'proposal submission due date must be on or after opening date'. The same date is allowed.

## Proposal needs a value or time and materials

Lines 734-744

A proposal fails readiness unless its proposal value is above 0 or it is marked time and materials. Turbo checks this whatever the proposal's status, so the proposal is ready to be awarded. The message is 'proposal must have a proposal value or be marked as time and materials to create a referencing project'.

## Proposal description minimum length

Lines 746-753

A proposal whose description is shorter than 3 characters fails readiness with 'description must be at least 3 characters'.
