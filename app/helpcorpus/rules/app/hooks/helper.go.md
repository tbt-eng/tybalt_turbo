<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/helper.go","sha256":"4676b385bf925f6254769181f6e58d5286923079046ca4e3998b674de2a85aaf","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/helper.go

## Uploaded file names have # and ? replaced

Lines 22-33

Turbo replaces each '#' and '?' in an uploaded file name with '_' (NormalizePendingFileNames). The stored file name then has no URL fragment or query characters.

## Attachment hash is a SHA256 of the uploaded file

Lines 38-48

When one file is uploaded to a field, Turbo calculates its SHA256 hash and returns it as a hex string. When no file is uploaded for the field, Turbo calculates no hash and raises no error.

## Only one file allowed per upload field

Lines 50-63

Turbo refuses a request that uploads more than one file to a single field. The error for that field is "too many files uploaded for field <field>" (code too_many_files).

## File open or hash failure is refused

Lines 68-99

If Turbo cannot open the uploaded file, it refuses with "error opening file for field <field>". If it cannot read the file to calculate the hash, it refuses with "error calculating hash".

## Division must exist and be active

Lines 107-140

Turbo refuses a division that cannot be found, with "specified division could not be found". It refuses an inactive division with "specified division is inactive". An empty division id passes without any check.

## Job and division pair must be allocated

Lines 159-162

Turbo checks that the division is in the job's time allocations (job_time_allocations). If the job or the division is empty or only spaces, Turbo skips the check.

## Allocation lookup failure gives a job error

Lines 165-172

If Turbo cannot read the job's allocations, it refuses with "Error checking job allocations" on the job field (code job_allocation_error).

## Job with no division allocations is refused

Lines 174-179

If the job has no division allocations, Turbo refuses with "This job has no division allocations configured". The error attaches to the job field.

## Division not allocated to the job is refused

Lines 181-195

If the division is not among the job's allocated divisions, Turbo refuses with "Division <code> is not allocated to this job". The message uses the division code, or the division id when the code is not available. The error attaches to the division field.

## Allocation check runs on create or job/division change only

Lines 213-228

Turbo checks the job/division allocation on every new record. On an update, it checks only when the job or the division differs from the saved value, ignoring surrounding spaces. Other edits, such as comments or approvals, skip the check, so old records with bad pairs still save. If the original record cannot be loaded, Turbo checks.

## Awarded proposal needs a referencing project to charge time

Lines 231-249

For an Awarded proposal, Turbo looks at jobs that reference it and have a project award date. If there are none, it refuses with "Proposal is Awarded but has no referencing project yet; you cannot charge time to the proposal".

## Awarded proposal accepts time for under 30 days

Lines 251-270

Turbo uses the earliest award date among the referencing projects. Time is allowed while the proposal's age is under 30 days (30*24 hours). At 30 days or more, Turbo refuses with "Proposal was awarded more than 30 days in the past; use the referencing project instead".

## Bad or future project award date is refused

Lines 253-261

If the earliest project award date is not a valid YYYY-MM-DD date, Turbo refuses with "Referenced project has an invalid project award date". If the date is in the future, it refuses with "Referenced project has a future project award date".

## Projects must be Active to charge time

Lines 282-301

For a job that is not a proposal, Turbo refuses time unless the status is Active. The message is "Job status must be Active". A missing job gives "invalid job reference".

## Proposal statuses that allow time tracking

Lines 286-295

For a proposal (found by its number), time can be charged when its status is In Progress, Submitted or Active (legacy). An Awarded proposal follows the 30-day window. Any other status is refused with "Proposal status must be In Progress, Submitted, or recently Awarded".
