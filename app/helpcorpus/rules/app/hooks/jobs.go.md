<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/jobs.go","sha256":"e3677dc94a488a674889619653dcbce48a6d30ccd3b74842d2126429300628b6","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/jobs.go

## Project validation level depends on proposal or retirement

Lines 24-36

Asked as: why full validation on save; which checks run when saving; strict validation edit; validation for project edits

Turbo applies full validation to every job save, except a proposal save or a status-only project retirement. Strict callers force full validation even for status-only edits (ProcessJobCoreStrict).

## Retiring a project skips full validation

Lines 29-31

Asked as: closing project skips checks; cancel project without billing check; retire project quickly; closed project no billing needed

A save on an existing project that changes only the status to Closed or Cancelled gets the relaxed validation. The relaxed path skips the billing, contact and manager checks. Setting a project Active does not qualify, nor does any edit of another field. Proposals are also excluded. The authorizing document is ignored when checking for other changes.

## Checks run in a fixed order on every job save

Lines 40-86

Asked as: order of save checks; which check runs first; why save refused first; sequence of job checks

Turbo checks in this order: job editing switch, project write check, cancelled proposal lock, outstanding balance permission, rate sheet, field cleaning, then full job validation. The first failing check refuses the save.

## Job editing can be switched off

Lines 42-48

Asked as: cannot edit jobs at all; jobs locked for editing; job editing turned off; can't save any job

When job editing is disabled in the app configuration, Turbo refuses every job create and update. It returns the error utilities.ErrJobsEditingDisabled (defined in another file, so its message is not shown here).

## Project write workflow check runs before other job rules

Lines 49-51

Asked as: project write workflow refused; why project write blocked; save blocked before balance check; validate project job write

Turbo runs the project write check (validateProjectJobWrite) right after the editing switch and before the cancelled-proposal, balance and rate sheet checks. A refusal from that check stops the save. This file does not show what that check requires.

## Cancelled proposals cannot be changed

Lines 53-67

Asked as: cancelled proposal can't edit; change cancelled proposal; proposal locked after cancel; modify cancelled proposal

Turbo refuses any edit to a proposal whose saved status is Cancelled. The message is "cancelled proposals cannot be modified" (the detail reads "this proposal has been cancelled and cannot be modified").

## Job number is assigned only on create

Lines 88-135

Asked as: job number changes on edit; number reassigned when updating; why job number not changing; new job number assigned

A new job gets its number after it passes validation. A job with a parent gets a sub-job number based on the parent's number. A job with no parent gets a top-level number for its type (project or proposal). An update never changes the number.

## Sub-job number failure message

Lines 98-126

Asked as: error generating sub-job number; sub-job number failed; can't create sub-job number; job number error on save

When Turbo cannot generate a sub-job number, it refuses the save with 'error generating sub-job number'. For a top-level job the message is 'error generating job number'.

## Edited imported jobs are marked as changed

Lines 140-144

Asked as: imported job edited; legacy system sync job; job no longer imported; edit imported job flag

When an existing job changes in any field, Turbo marks it as no longer imported, so the legacy system gets the change written back.

## Outstanding balance date is set automatically

Lines 150-171

Asked as: outstanding balance date changed; balance date not updating; why balance date reset; balance date set to today

On a project, Turbo sets the outstanding balance date (outstanding_balance_date) to today when the balance changes. On create, a nonzero balance counts as a change. When the balance is unchanged on an update, the original date is kept.

## Who may change the outstanding balance

Lines 173-245

Asked as: can't edit outstanding balance; insufficient permissions balance; payables admin balance change; authentication required balance

Changing the outstanding balance (outstanding_balance) on an existing job needs the job claim (jobs claim) or the payables admin claim (payables_admin). Without either, Turbo says "insufficient permissions to edit outstanding balance". If the request has no signed-in user, it says "authentication required to edit outstanding balance". A save that leaves the balance unchanged needs no claim, and neither does creating a job.

## Active projects need a completion date

Lines 247-262

Asked as: project completion date is required; active project no end date; blank completion date active; why active project refused

A project with status Active and no completion date (project_completion_date) is refused with "project completion date is required". Closed and Cancelled projects and proposals can leave it blank.

## Completion date must be valid and not before award date

Lines 263-295

Asked as: completion before award date; invalid completion date; end date before start date; project dates in wrong order

A completion date must be in YYYY-MM-DD format, otherwise Turbo says "invalid project completion date". When the project has an award date, the completion date must be on or after it. The same day is allowed. Otherwise Turbo says "project completion date must be on or after project award date". A project with no award date skips this comparison. An invalid award date is refused with "invalid project award date".

## Date format errors on completion and award dates

Lines 263-284

Asked as: date format wrong project; invalid project award date; completion date format error; date not YYYY-MM-DD

A project completion date not in YYYY-MM-DD is refused with 'invalid project completion date'. An award date not in that format is refused with 'invalid project award date'. The comparison runs only if an award date exists.

## Completion date on the award date is allowed

Lines 285-293

Asked as: completion same day as award; completion equals award date; end date same as award; can completion match award

A completion date equal to the award date passes. Only an earlier date is refused with 'project completion date must be on or after project award date'.

## Proposal due date must not precede opening date

Lines 301-332

Asked as: due date before opening; submission due date wrong; proposal dates reversed; opening date after due

A proposal's submission due date must be on or after its opening date. The same day is allowed. Otherwise Turbo says "proposal submission due date must be on or after opening date". Dates in the wrong format are refused with "invalid proposal opening date format" or "invalid proposal submission due date format".

## Proposal date errors

Lines 301-330

Asked as: bad proposal date; proposal date format error; due date invalid format; opening date error

A bad opening date gives 'invalid proposal opening date format'. A bad due date gives 'invalid proposal submission due date format'. A due date equal to the opening date passes.

## Proposals are recognized by their dates on create

Lines 340-349

Asked as: why saved as proposal; job created as proposal; proposal dates on create; treated as proposal

While a job is being created, Turbo treats it as a proposal when the project award date is empty and both proposal dates are filled. On update it uses the original job number instead.

## Retirement means only status changes to Closed or Cancelled

Lines 351-361

Asked as: what counts as retiring project; retire project by status change; setting project closed only; project retirement definition

A project retirement is a save of an existing project that changes nothing except status (authorizing_document is ignored) and sets status to Closed or Cancelled. Setting a project Active is not a retirement. New jobs and proposals never count.

## Job fields Turbo sets automatically by type

Lines 370-400

Asked as: proposal value reset to zero; authorizing document changed to PA; why PA set on project; fields disappeared on save

On a project, Turbo sets the proposal value (proposal_value) to 0 and the authorizing document to PA on every save, so older projects with PO or Unauthorized become PA. On a proposal, Turbo clears the authorizing document, client PO, client reference number, rate sheet, outstanding balance and date, proposal link and project value. Turbo also trims spaces from the client PO.

## Client PO is trimmed

Lines 374-377

Asked as: client PO spaces; trim PO number; PO extra spaces removed; client PO whitespace

Turbo trims leading and trailing spaces from a job's client PO (client_po) when saving.

## Proposal fields cleared and project fields forced

Lines 379-398

Asked as: proposal fields disappeared; project value cleared; rate sheet removed from proposal; fields wiped on save

For a project, Turbo sets proposal value to 0 and authorizing document to PA. For a proposal, Turbo clears authorizing document, client PO, client reference number, rate sheet, outstanding balance and its date, proposal reference and project value.

## Sub-job number format and no nested sub-jobs

Lines 427-450

Asked as: sub-job of sub-job; nested sub-jobs not allowed; parent job number format; can't add child to sub-job

Only a base job number (optionally starting with P, 2-digit year, dash, 3 or 4 digits) can be a parent. A sub-job number (such as 25-0123-01) is not a base number, so a sub-job cannot have sub-jobs.

## Job type comes from the number or dates

Lines 453-500

Asked as: job type wrong; project or proposal decided how; job type changed unexpectedly; type from dates

For an existing job, the type comes from its number: a number starting with P is a proposal, any other is a project. A sub-job suffix is ignored. For a new job, a job with a linked proposal is a project. Otherwise the job needs either only a project award date (project) or both proposal dates and no award date (proposal).

## Proposal type comes from the P prefix

Lines 453-464

Asked as: P prefix proposal; why number starts with P; proposal number letter P; what P means in job number

Turbo treats a job as a proposal when its number starts with P, after removing any sub-job suffix. Any other number is a project.

## New job needs exactly one date configuration

Lines 480-497

Asked as: exactly one of project or proposal; date configuration required; new job dates error; both award and proposal dates

A new job must have either a project award date alone or both proposal dates with no award date. Anything else is refused with "exactly one of project or proposal date configurations is required".

## Status-only retirement needs no billing data

Lines 507-515

Asked as: closing project without invoicing; no billing needed to close; imported project missing billing; status only save no contact

On a project update that does not force full validation, Turbo returns right after the completion date check. It skips the invoicing profile, project value, contact and manager checks, because many imported projects lack them.

## Proposal status-only changes

Lines 517-560

Asked as: proposal status change refused; value requirement for status; No Bid status save; proposal status only edit

When only the status of a proposal changes, Turbo refuses Active and Closed and checks the value requirement. It also checks the comment for No Bid or Cancelled, and skips the other full validation. Turbo clears fields on proposals before the change check, so an imported proposal may still be fully validated.

## Proposal status change without other edits is limited

Lines 517-535

Asked as: invalid status for proposal; set proposal active; submitted awarded needs value; proposal closed refused

When only the status of a proposal changes, Turbo refuses Active or Closed with 'invalid status for proposal'. It also requires a value (or time and materials) for Submitted, Awarded and Not Awarded.

## Status-only proposal update skips other checks

Lines 517-561

Asked as: proposal status save skips checks; cleanJob balance date effect; why proposal full validation; status update proposal date check

When a proposal save changes only status, Turbo returns after the status checks and skips the date, location, branch, contact and manager checks. Normalization by cleanJob (for example clearing an old balance date) can make it count as a full edit.

## No Bid or Cancelled proposal needs a comment first

Lines 536-557

Asked as: comment required for status; no bid needs note; cancel proposal needs note; a comment is required to set this status

Moving a proposal to No Bid or Cancelled needs a client note that records a change to that status. Without it Turbo says "a comment is required to set this status". The check runs only when the status actually changes and not on create.

## Date rules by job type on create and edit

Lines 563-647

Asked as: invalid dates for project; invalid dates for proposal; project needs award date; proposal dates required

On create, a project needs an award date and no proposal dates ("invalid dates for project"). A proposal needs both proposal dates and no award date ("invalid dates for proposal"). On update, these rules apply only when a date field changed, and an empty field falls back to its original value.

## Location is required

Lines 649-659

Asked as: invalid or missing location; plus code location; location blank on job; job location format

Every job needs a location in Plus Code format. A blank or malformed location is refused with "invalid or missing location".

## Branch is required

Lines 661-670

Asked as: branch is required; no branch on job; missing branch on save; pick a branch

A job without a branch is refused with "branch is required".

## Contact and manager are required

Lines 672-714

Asked as: contact is required; manager is required; imported job missing contact; no manager on job

On creates and full edits, a job without a contact or manager is refused with "contact is required" or "manager is required". A status-only update is exempt. Turbo does not fill these in for imported jobs, so such a job needs both before any other edit succeeds.

## Allowed statuses for projects

Lines 718-728

Asked as: invalid status for project; project status awarded; can't set project submitted; project status options

A project may be Active, Closed or Cancelled. Turbo refuses Awarded, Not Awarded, Submitted, In Progress and No Bid with "invalid status for project".

## Allowed statuses for proposals

Lines 732-754

Asked as: new proposal status; proposal active status refused; proposal status list; must start in progress

A proposal may be In Progress, Submitted, Awarded, Not Awarded, Cancelled or No Bid. Turbo refuses Active and Closed with "invalid status for proposal". A new proposal must start as In Progress or Submitted, otherwise Turbo says "new proposals must start as In Progress or Submitted".

## New proposals start In Progress or Submitted

Lines 744-754

Asked as: new proposals must start; create proposal as awarded; starting status for proposal; new proposal status refused

A new proposal can only have status In Progress or Submitted. Other statuses are refused with 'new proposals must start as In Progress or Submitted'.

## Comment check on full proposal edits

Lines 758-784

Asked as: full edit proposal note; no bid note needed on edit; client note for proposal cancel; proposal comment missing

On a full edit of a proposal, moving it to No Bid or Cancelled also needs a matching client note. The check runs only when the status differs from the saved one.

## Linked proposal must be Awarded

Lines 788-841

Asked as: linked proposal not awarded; proposal must be awarded; referenced job not a proposal; invalid referenced proposal

A job that links a proposal needs a proposal that exists, is a proposal and has status Awarded. If it is In Progress or Submitted, Turbo says "proposal must be awarded". If it is Not Awarded, Cancelled or No Bid, Turbo says "proposal has invalid status". A job that does not point to a proposal gets "referenced job is not a proposal", and a missing proposal gets "invalid referenced proposal".

## Proposal reference messages by status

Lines 789-841

Asked as: project proposal reference error; proposal reference status; referenced proposal missing; project linked proposal refused

A project that references a proposal must point to a job with a proposal number. In Progress or Submitted gives 'proposal must be awarded'. Not Awarded, Cancelled or No Bid gives 'proposal has invalid status'. A missing job gives 'invalid referenced proposal'; a non-proposal gives 'referenced job is not a proposal'.

## Sub-job parent rules

Lines 843-889

Asked as: invalid parent job; parent missing client; sub-job client mismatch; proposals cannot have sub-jobs

A parent must exist ("invalid parent") and must not be a proposal ("proposals cannot have sub-jobs"). The parent must have a client ("parent missing client"). A sub-job's client must match its parent's ("sub-job must have same client as parent"). Turbo writes the parent's client onto the sub-job.

## Contact must belong to the job's client

Lines 891-923

Asked as: contact from wrong client; contact requires a client; invalid contact on job; contact not belong to client

A job's contact needs a client ("contact requires a client"). The contact must exist ("invalid contact") and belong to that client ("contact does not belong to client").

## Job number cannot change type on update

Lines 929-941

Asked as: job number type conflict; change project to proposal number; number type mismatch; can't change job type

On update, the type implied by the saved number must match the derived type. Otherwise Turbo says "job number type conflict".

## Managers must be active users

Lines 943-1226

Asked as: manager not active; inactive manager on job; alternate manager must be active; manager user disabled

The manager and the alternate manager, if set, must be active users. Otherwise Turbo says "manager must be an active user" or "alternate manager must be an active user".

## Projects need an invoicing profile

Lines 951-984

Asked as: invoicing profile required; no billing profile on project; invoicing information missing; profile not for client

A project save needs an invoicing profile (invoicing_information). Otherwise Turbo says "invoicing profile is required". A status-only change to Closed or Cancelled is exempt. Setting an older project Active, or any other edit, needs one. Turbo also checks that the profile and contact belong to the job's client.

## Status-only check ignores authorizing document

Lines 986-992

Asked as: authorizing document status change; status only check PA; why authorizing doc ignored; status change only document

When Turbo decides whether only the status changed, it ignores the authorizing document, because Turbo sets that itself.

## Comment lookup is by target status

Lines 994-1004

Asked as: failed to check client notes; client note lookup failed; comment check error; note target status

The required comment is a client note on the job whose status-changed-to equals the new status. If Turbo cannot check, it refuses with 'failed to check for client notes'.

## Active or Closed project needs a value above zero

Lines 1010-1041

Asked as: project value must be above zero; zero value project active; value greater than zero needed; time and materials project value

A project with status Active or Closed needs a project value greater than zero, including time and materials projects. Otherwise Turbo says "project value must be greater than zero". An existing project whose stored value is 0 and whose status does not change can still be edited. A save marked as a value restore skips this check.

## Value restore workflow skips the project value rule

Lines 1014-1016

Asked as: restore project value; restore value workflow; value check skipped restore; restore value save

A project save that is a restore-value workflow operation (ProjectWorkflowMutationRestoreValue) skips the project value check.

## Project value of zero is allowed when unchanged

Lines 1025-1029

Asked as: project value zero edit; zero value project still editable; old project value 0; edit zero value project

An existing project with a stored project value of 0 can be edited while its value stays 0 and its status does not change. Moving it to Active or Closed still needs a value above zero.

## Proposal value or time and materials required

Lines 1042-1055

Asked as: proposal value required; time and materials proposal; proposal value or T&M; submitted proposal no value

A proposal with status Submitted, Awarded or Not Awarded needs a proposal value above zero or the time and materials flag. Otherwise Turbo says "proposal value or time and materials required".

## Proposal value of zero allowed with time and materials

Lines 1043-1054

Asked as: zero value proposal time and materials; T&M flag on proposal; proposal zero value refused; time and materials checkbox

A proposal with status Submitted, Awarded or Not Awarded needs a proposal value above zero unless time and materials is on. Otherwise Turbo refuses with 'proposal value or time and materials required'.

## Rate sheet required and active for projects

Lines 1065-1176

Asked as: rate sheet required for new project; rate sheet inactive; invalid rate sheet; rate sheet not found

A new project needs a rate sheet, otherwise Turbo says "rate sheet is required for new projects". A rate sheet that is set or changed must exist ("invalid rate sheet") and be active ("rate sheet must be active"). Proposals skip these checks. An unchanged rate sheet is not rechecked.

## Rate sheet rule exempts proposals and unchanged sheets

Lines 1065-1075

Asked as: proposal rate sheet error; unchanged rate sheet recheck; rate sheet only checked when changed; why rate sheet not checked

Proposals skip rate sheet validation. For an existing project, Turbo checks the rate sheet exists ('invalid rate sheet') and is active only when it changed.

## Changing or clearing a rate sheet

Lines 1099-1141

Asked as: rate sheet cannot be cleared; change rate sheet permission; rate_sheet_revise claim; remove rate sheet from project

A project's rate sheet cannot be cleared once set. Turbo says "rate sheet cannot be cleared". Changing it to another sheet needs the rate sheet revise claim (rate_sheet_revise). Otherwise Turbo says "insufficient permissions to change rate sheet".

## Manager check failure message

Lines 1181-1213

Asked as: failed to check manager; manager active check failed; alternate manager check error; manager status lookup failed

If Turbo cannot check whether the manager or alternate manager is active, it refuses with 'failed to check manager active status' or 'failed to check alternate manager active status'.

## Top-level job number generation

Lines 1248-1293

Asked as: next top-level job number; no unused job number; job number 9999 limit; how job numbers assigned

A new job without a parent gets the next number for the current two-digit year, as YY-NNNN for a project or PYY-NNNN for a proposal. Turbo takes the highest existing sequence across the 3-digit and 4-digit formats and adds one. It writes new numbers with 4 digits. Sub-jobs are not counted. If no unused number remains up to 9999, number generation fails.

## Next top-level job number is highest plus one, per year

Lines 1254-1293

Asked as: unable to generate unique job number; job number sequence; next number for year; job number 26-0015

Turbo numbers a new job as the 2-digit current year, a dash, and a 4-digit sequence (such as 26-0015). A proposal starts with P. The sequence is one more than the highest existing 3- or 4-digit number for that year and type, or 0001 if none. Sub-jobs are not counted. It skips numbers already used, and fails with 'unable to generate a unique job number' past 9999.

## Sub-job number generation

Lines 1295-1331

Asked as: maximum sub-jobs 99; sub-job suffix; too many sub-jobs; sub-job numbering

A sub-job number is its parent's number plus a two-digit suffix. The suffix is the highest existing suffix plus one, e.g. 25-0123-01. Only a top-level job can be a parent. Beyond 99 sub-jobs, Turbo reports "maximum number of sub-jobs reached (99)". Turbo reads at most 100 existing children when it looks for the highest suffix.
