<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/jobs.go","sha256":"e3677dc94a488a674889619653dcbce48a6d30ccd3b74842d2126429300628b6","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/jobs.go

## Project validation level depends on proposal or retirement

Lines 24-36

Turbo applies full validation to every job save, except a proposal save or a status-only project retirement. Strict callers force full validation even for status-only edits (ProcessJobCoreStrict).

## Retiring a project skips full validation

Lines 29-31

A save on an existing project that changes only the status to Closed or Cancelled gets the relaxed validation. The relaxed path skips the billing, contact and manager checks. Setting a project Active does not qualify, nor does any edit of another field. Proposals are also excluded. The authorizing document is ignored when checking for other changes.

## Checks run in a fixed order on every job save

Lines 40-86

Turbo checks in this order: job editing switch, project write check, cancelled proposal lock, outstanding balance permission, rate sheet, field cleaning, then full job validation. The first failing check refuses the save.

## Job editing can be switched off

Lines 42-48

When job editing is disabled in the app configuration, Turbo refuses every job create and update. It returns the error utilities.ErrJobsEditingDisabled (defined in another file, so its message is not shown here).

## Project write workflow check runs before other job rules

Lines 49-51

Turbo runs the project write check (validateProjectJobWrite) right after the editing switch and before the cancelled-proposal, balance and rate sheet checks. A refusal from that check stops the save. This file does not show what that check requires.

## Cancelled proposals cannot be changed

Lines 53-67

Turbo refuses any edit to a proposal whose saved status is Cancelled. The message is "cancelled proposals cannot be modified" (the detail reads "this proposal has been cancelled and cannot be modified").

## Job number is assigned only on create

Lines 88-135

A new job gets its number after it passes validation. A job with a parent gets a sub-job number based on the parent's number. A job with no parent gets a top-level number for its type (project or proposal). An update never changes the number.

## Sub-job number failure message

Lines 98-126

When Turbo cannot generate a sub-job number, it refuses the save with 'error generating sub-job number'. For a top-level job the message is 'error generating job number'.

## Edited imported jobs are marked as changed

Lines 140-144

When an existing job changes in any field, Turbo marks it as no longer imported, so the legacy system gets the change written back.

## Outstanding balance date is set automatically

Lines 150-171

On a project, Turbo sets the outstanding balance date (outstanding_balance_date) to today when the balance changes. On create, a nonzero balance counts as a change. When the balance is unchanged on an update, the original date is kept.

## Who may change the outstanding balance

Lines 173-245

Changing the outstanding balance (outstanding_balance) on an existing job needs the job claim (jobs claim) or the payables admin claim (payables_admin). Without either, Turbo says "insufficient permissions to edit outstanding balance". If the request has no signed-in user, it says "authentication required to edit outstanding balance". A save that leaves the balance unchanged needs no claim, and neither does creating a job.

## Active projects need a completion date

Lines 247-262

A project with status Active and no completion date (project_completion_date) is refused with "project completion date is required". Closed and Cancelled projects and proposals can leave it blank.

## Completion date must be valid and not before award date

Lines 263-295

A completion date must be in YYYY-MM-DD format, otherwise Turbo says "invalid project completion date". When the project has an award date, the completion date must be on or after it. The same day is allowed. Otherwise Turbo says "project completion date must be on or after project award date". A project with no award date skips this comparison. An invalid award date is refused with "invalid project award date".

## Date format errors on completion and award dates

Lines 263-284

A project completion date not in YYYY-MM-DD is refused with 'invalid project completion date'. An award date not in that format is refused with 'invalid project award date'. The comparison runs only if an award date exists.

## Completion date on the award date is allowed

Lines 285-293

A completion date equal to the award date passes. Only an earlier date is refused with 'project completion date must be on or after project award date'.

## Proposal due date must not precede opening date

Lines 301-332

A proposal's submission due date must be on or after its opening date. The same day is allowed. Otherwise Turbo says "proposal submission due date must be on or after opening date". Dates in the wrong format are refused with "invalid proposal opening date format" or "invalid proposal submission due date format".

## Proposal date errors

Lines 301-330

A bad opening date gives 'invalid proposal opening date format'. A bad due date gives 'invalid proposal submission due date format'. A due date equal to the opening date passes.

## Proposals are recognized by their dates on create

Lines 340-349

While a job is being created, Turbo treats it as a proposal when the project award date is empty and both proposal dates are filled. On update it uses the original job number instead.

## Retirement means only status changes to Closed or Cancelled

Lines 351-361

A project retirement is a save of an existing project that changes nothing except status (authorizing_document is ignored) and sets status to Closed or Cancelled. Setting a project Active is not a retirement. New jobs and proposals never count.

## Job fields Turbo sets automatically by type

Lines 370-400

On a project, Turbo sets the proposal value (proposal_value) to 0 and the authorizing document to PA on every save, so older projects with PO or Unauthorized become PA. On a proposal, Turbo clears the authorizing document, client PO, client reference number, rate sheet, outstanding balance and date, proposal link and project value. Turbo also trims spaces from the client PO.

## Client PO is trimmed

Lines 374-377

Turbo trims leading and trailing spaces from a job's client PO (client_po) when saving.

## Proposal fields cleared and project fields forced

Lines 379-398

For a project, Turbo sets proposal value to 0 and authorizing document to PA. For a proposal, Turbo clears authorizing document, client PO, client reference number, rate sheet, outstanding balance and its date, proposal reference and project value.

## Sub-job number format and no nested sub-jobs

Lines 427-450

Only a base job number (optionally starting with P, 2-digit year, dash, 3 or 4 digits) can be a parent. A sub-job number (such as 25-0123-01) is not a base number, so a sub-job cannot have sub-jobs.

## Job type comes from the number or dates

Lines 453-500

For an existing job, the type comes from its number: a number starting with P is a proposal, any other is a project. A sub-job suffix is ignored. For a new job, a job with a linked proposal is a project. Otherwise the job needs either only a project award date (project) or both proposal dates and no award date (proposal).

## Proposal type comes from the P prefix

Lines 453-464

Turbo treats a job as a proposal when its number starts with P, after removing any sub-job suffix. Any other number is a project.

## New job needs exactly one date configuration

Lines 480-497

A new job must have either a project award date alone or both proposal dates with no award date. Anything else is refused with "exactly one of project or proposal date configurations is required".

## Status-only retirement needs no billing data

Lines 507-515

On a project update that does not force full validation, Turbo returns right after the completion date check. It skips the invoicing profile, project value, contact and manager checks, because many imported projects lack them.

## Proposal status-only changes

Lines 517-560

When only the status of a proposal changes, Turbo refuses Active and Closed and checks the value requirement. It also checks the comment for No Bid or Cancelled, and skips the other full validation. Turbo clears fields on proposals before the change check, so an imported proposal may still be fully validated.

## Proposal status change without other edits is limited

Lines 517-535

When only the status of a proposal changes, Turbo refuses Active or Closed with 'invalid status for proposal'. It also requires a value (or time and materials) for Submitted, Awarded and Not Awarded.

## Status-only proposal update skips other checks

Lines 517-561

When a proposal save changes only status, Turbo returns after the status checks and skips the date, location, branch, contact and manager checks. Normalization by cleanJob (for example clearing an old balance date) can make it count as a full edit.

## No Bid or Cancelled proposal needs a comment first

Lines 536-557

Moving a proposal to No Bid or Cancelled needs a client note that records a change to that status. Without it Turbo says "a comment is required to set this status". The check runs only when the status actually changes and not on create.

## Date rules by job type on create and edit

Lines 563-647

On create, a project needs an award date and no proposal dates ("invalid dates for project"). A proposal needs both proposal dates and no award date ("invalid dates for proposal"). On update, these rules apply only when a date field changed, and an empty field falls back to its original value.

## Location is required

Lines 649-659

Every job needs a location in Plus Code format. A blank or malformed location is refused with "invalid or missing location".

## Branch is required

Lines 661-670

A job without a branch is refused with "branch is required".

## Contact and manager are required

Lines 672-714

On creates and full edits, a job without a contact or manager is refused with "contact is required" or "manager is required". A status-only update is exempt. Turbo does not fill these in for imported jobs, so such a job needs both before any other edit succeeds.

## Allowed statuses for projects

Lines 718-728

A project may be Active, Closed or Cancelled. Turbo refuses Awarded, Not Awarded, Submitted, In Progress and No Bid with "invalid status for project".

## Allowed statuses for proposals

Lines 732-754

A proposal may be In Progress, Submitted, Awarded, Not Awarded, Cancelled or No Bid. Turbo refuses Active and Closed with "invalid status for proposal". A new proposal must start as In Progress or Submitted, otherwise Turbo says "new proposals must start as In Progress or Submitted".

## New proposals start In Progress or Submitted

Lines 744-754

A new proposal can only have status In Progress or Submitted. Other statuses are refused with 'new proposals must start as In Progress or Submitted'.

## Comment check on full proposal edits

Lines 758-784

On a full edit of a proposal, moving it to No Bid or Cancelled also needs a matching client note. The check runs only when the status differs from the saved one.

## Linked proposal must be Awarded

Lines 788-841

A job that links a proposal needs a proposal that exists, is a proposal and has status Awarded. If it is In Progress or Submitted, Turbo says "proposal must be awarded". If it is Not Awarded, Cancelled or No Bid, Turbo says "proposal has invalid status". A job that does not point to a proposal gets "referenced job is not a proposal", and a missing proposal gets "invalid referenced proposal".

## Proposal reference messages by status

Lines 789-841

A project that references a proposal must point to a job with a proposal number. In Progress or Submitted gives 'proposal must be awarded'. Not Awarded, Cancelled or No Bid gives 'proposal has invalid status'. A missing job gives 'invalid referenced proposal'; a non-proposal gives 'referenced job is not a proposal'.

## Sub-job parent rules

Lines 843-889

A parent must exist ("invalid parent") and must not be a proposal ("proposals cannot have sub-jobs"). The parent must have a client ("parent missing client"). A sub-job's client must match its parent's ("sub-job must have same client as parent"). Turbo writes the parent's client onto the sub-job.

## Contact must belong to the job's client

Lines 891-923

A job's contact needs a client ("contact requires a client"). The contact must exist ("invalid contact") and belong to that client ("contact does not belong to client").

## Job number cannot change type on update

Lines 929-941

On update, the type implied by the saved number must match the derived type. Otherwise Turbo says "job number type conflict".

## Managers must be active users

Lines 943-1226

The manager and the alternate manager, if set, must be active users. Otherwise Turbo says "manager must be an active user" or "alternate manager must be an active user".

## Projects need an invoicing profile

Lines 951-984

A project save needs an invoicing profile (invoicing_information). Otherwise Turbo says "invoicing profile is required". A status-only change to Closed or Cancelled is exempt. Setting an older project Active, or any other edit, needs one. Turbo also checks that the profile and contact belong to the job's client.

## Status-only check ignores authorizing document

Lines 986-992

When Turbo decides whether only the status changed, it ignores the authorizing document, because Turbo sets that itself.

## Comment lookup is by target status

Lines 994-1004

The required comment is a client note on the job whose status-changed-to equals the new status. If Turbo cannot check, it refuses with 'failed to check for client notes'.

## Active or Closed project needs a value above zero

Lines 1010-1041

A project with status Active or Closed needs a project value greater than zero, including time and materials projects. Otherwise Turbo says "project value must be greater than zero". An existing project whose stored value is 0 and whose status does not change can still be edited. A save marked as a value restore skips this check.

## Value restore workflow skips the project value rule

Lines 1014-1016

A project save that is a restore-value workflow operation (ProjectWorkflowMutationRestoreValue) skips the project value check.

## Project value of zero is allowed when unchanged

Lines 1025-1029

An existing project with a stored project value of 0 can be edited while its value stays 0 and its status does not change. Moving it to Active or Closed still needs a value above zero.

## Proposal value or time and materials required

Lines 1042-1055

A proposal with status Submitted, Awarded or Not Awarded needs a proposal value above zero or the time and materials flag. Otherwise Turbo says "proposal value or time and materials required".

## Proposal value of zero allowed with time and materials

Lines 1043-1054

A proposal with status Submitted, Awarded or Not Awarded needs a proposal value above zero unless time and materials is on. Otherwise Turbo refuses with 'proposal value or time and materials required'.

## Rate sheet required and active for projects

Lines 1065-1176

A new project needs a rate sheet, otherwise Turbo says "rate sheet is required for new projects". A rate sheet that is set or changed must exist ("invalid rate sheet") and be active ("rate sheet must be active"). Proposals skip these checks. An unchanged rate sheet is not rechecked.

## Rate sheet rule exempts proposals and unchanged sheets

Lines 1065-1075

Proposals skip rate sheet validation. For an existing project, Turbo checks the rate sheet exists ('invalid rate sheet') and is active only when it changed.

## Changing or clearing a rate sheet

Lines 1099-1141

A project's rate sheet cannot be cleared once set. Turbo says "rate sheet cannot be cleared". Changing it to another sheet needs the rate sheet revise claim (rate_sheet_revise). Otherwise Turbo says "insufficient permissions to change rate sheet".

## Manager check failure message

Lines 1181-1213

If Turbo cannot check whether the manager or alternate manager is active, it refuses with 'failed to check manager active status' or 'failed to check alternate manager active status'.

## Top-level job number generation

Lines 1248-1293

A new job without a parent gets the next number for the current two-digit year, as YY-NNNN for a project or PYY-NNNN for a proposal. Turbo takes the highest existing sequence across the 3-digit and 4-digit formats and adds one. It writes new numbers with 4 digits. Sub-jobs are not counted. If no unused number remains up to 9999, number generation fails.

## Next top-level job number is highest plus one, per year

Lines 1254-1293

Turbo numbers a new job as the 2-digit current year, a dash, and a 4-digit sequence (such as 26-0015). A proposal starts with P. The sequence is one more than the highest existing 3- or 4-digit number for that year and type, or 0001 if none. Sub-jobs are not counted. It skips numbers already used, and fails with 'unable to generate a unique job number' past 9999.

## Sub-job number generation

Lines 1295-1331

A sub-job number is its parent's number plus a two-digit suffix. The suffix is the highest existing suffix plus one, e.g. 25-0123-01. Only a top-level job can be a parent. Beyond 99 sub-jobs, Turbo reports "maximum number of sub-jobs reached (99)". Turbo reads at most 100 existing children when it looks for the highest suffix.
