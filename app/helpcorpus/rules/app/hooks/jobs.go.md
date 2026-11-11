<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/jobs.go","sha256":"3f38ac0e1c6dbbc77925e50763b33db9ba3aac32771b79994bb185db33d049da","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/jobs.go

## Job saves go through one set of business rules

Lines 20-22

Asked as: job save rules; who can save jobs; job save permissions; checks applied when saving

Turbo applies the job rules (ProcessJob, ProcessJobCore) to every job save made through a request. The rules use the saving user's record (the auth record) to check permissions such as editing the outstanding balance.

## Which saves get full validation

Lines 27-36

Asked as: why project fully validated; project save validation; strict project checks; which project saves are checked

Turbo fully validates every project save, except a save that only sets an existing project to Closed or Cancelled (a retirement). Setting a project Active is not a retirement. Proposal saves get the relaxed status-only check when only the status changes. The strict path (ProcessJobCoreStrict) always validates fully.

## Job editing can be switched off

Lines 40-48

Asked as: can't edit jobs; job editing disabled; jobs locked down; editing jobs turned off

When job editing is disabled in the app configuration, Turbo refuses every job create or update. It returns the error utilities.ErrJobsEditingDisabled. The setting comes from IsJobsEditingEnabled.

## Cancelled proposals cannot be changed

Lines 53-67

Asked as: edit cancelled proposal; cancelled proposal locked; can't change cancelled proposal; proposal already cancelled

Turbo refuses any edit to a proposal whose saved status is Cancelled. It says "cancelled proposals cannot be modified" (this proposal has been cancelled and cannot be modified).

## Cancelled proposal check uses the original number and status

Lines 53-58

Asked as: reopen cancelled proposal; uncancel proposal; cancelled proposal status change; restore cancelled proposal

Turbo decides that a proposal is cancelled from the saved record, not from the incoming change. A proposal already Cancelled is refused even if the edit would change the status.

## New job numbers are generated

Lines 88-126

Asked as: how job numbers are assigned; next job number; job number sequence; sub-job numbering scheme

On create, Turbo assigns the number. A top-level job gets the next number for the current two-digit year in the form YY-NNNN, with a P prefix for proposals. The next number is one above the highest existing 3- or 4-digit sequence for that year. A sub-job gets its parent's number plus a two-digit suffix, one above the highest existing suffix. Turbo allows at most 99 sub-jobs per parent. A sub-job parent must be a top-level number, so a sub-job cannot have sub-jobs.

## Number generation failure is reported

Lines 100-126

Asked as: job number error; error generating job number; can't get job number; sub-job number failed

When Turbo cannot generate a job number, it refuses the save with "error generating job number" for top-level jobs, or "error generating sub-job number" for sub-jobs.

## Edited jobs lose imported flag

Lines 140-144

Asked as: job no longer imported; legacy system modified flag; imported job edited; why job marked modified

When an update changes any field of a job, Turbo marks the job as no longer imported. The legacy system then receives the job as modified.

## Outstanding balance date is set automatically

Lines 150-171

Asked as: balance date not updating; outstanding date set automatically; when balance date changes; outstanding balance date

On a project, Turbo sets the outstanding balance date (outstanding_balance_date) to today when the balance changes. On create, any nonzero balance counts as a change. If the balance is unchanged on an update, the old date is kept.

## Who may change a job's outstanding balance

Lines 173-245

Asked as: can't edit outstanding balance; outstanding balance permission denied; who can change balance; balance edit not allowed

Changing the outstanding balance of an existing job needs the job claim or the payables_admin claim. Without either, Turbo says "insufficient permissions to edit outstanding balance". Without a signed-in user, it says "authentication required to edit outstanding balance". Creating a job skips this check, and so does leaving the balance unchanged.

## Active projects need a completion date

Lines 247-262

Asked as: active project no end date; completion date required; project end date missing; active project needs finish date

An Active project must have a project completion date. Otherwise Turbo says "project completion date is required" (an estimate is acceptable). Closed and Cancelled projects, and proposals, may leave it blank.

## Completion date cannot precede award date

Lines 263-294

Asked as: completion date before award; end date earlier than award; award date after completion; project finish before start

A project completion date must be a valid YYYY-MM-DD date. It must be on or after the project award date, so the same day is allowed. Turbo compares the two only when an award date exists, since legacy projects may have none. Otherwise it says "project completion date must be on or after project award date".

## Proposal due date must not precede opening date

Lines 298-332

Asked as: due date before opening; submission date before open date; proposal dates wrong order; proposal due date error

A proposal's submission due date must be on or after its opening date, so the same day is allowed. Both dates must be in YYYY-MM-DD format. Otherwise Turbo says "proposal submission due date must be on or after opening date".

## Retiring a project skips billing and contact checks

Lines 351-361

Asked as: close project without contact; cancel project no invoicing; close project quickly; retire project without checks

A save that only sets an existing project to Closed or Cancelled skips the later field checks. It does not need an invoicing profile, a project value, a contact or a manager. Active and Closed still need a project completion date check on the new status.

## Proposals lose project-only fields on save

Lines 370-398

Asked as: proposal fields cleared; why client PO disappeared; proposal project value cleared; proposal lost values on save

On every save Turbo trims whitespace from the client PO. On a proposal it clears authorizing document, client PO, client reference number, rate sheet, outstanding balance and its date, proposal link and project value.

## Projects get fixed authorizing document and no proposal value

Lines 379-387

Asked as: authorizing document changed to PA; PO became PA; why authorizing doc is PA; project document reset

On every project save, Turbo sets the authorizing document (authorizing_document) to PA and the proposal value to 0. Older projects with PO or Unauthorized become PA on their next save.

## Job number format and sub-job eligibility

Lines 418-450

Asked as: invalid job number format; sub-job of sub-job; job number pattern; which numbers can be parents

Only a base job number can be a parent: an optional P, 2-digit year, a dash and 3 or 4 digits, such as 25-123 or P25-0123. A number with a sub-job suffix, such as 25-0123-01, is not a base number. Sub-jobs cannot have their own sub-jobs.

## Job type comes from the number prefix

Lines 452-464

Asked as: proposal or project by number; meaning of P prefix; job type wrong; job type from number

Turbo treats a job as a proposal when its number starts with P, after removing any sub-job suffix. Every other number is a project. On update, the type comes from the existing job number.

## Job type is set by the dates on create

Lines 478-500

Asked as: not a project or proposal; exactly one date configuration; date configurations required; job type from dates

On create, a job with a proposal link, or with only a project award date, is a project. A job with both proposal dates and no award date is a proposal. Any other combination is refused with "exactly one of project or proposal date configurations is required". After creation, the job number fixes the type: a leading P means proposal.

## New job with a proposal reference is a project

Lines 480-497

Asked as: job linked to proposal is project; create project from proposal; proposal reference project type; project created from proposal

On create, a job that names a proposal (proposal) is a project, whatever dates it has. Without a proposal, a project has an award date and no proposal dates, and a proposal has both proposal dates and no award date. Any other date mix is refused with "exactly one of project or proposal date configurations is required".

## Status-only edit of a project skips full validation

Lines 507-515

Asked as: limited checks on project save; why location not checked; status only project save; project save skips checks

When a save of an existing project is not full validation, the date and completion checks run and then validation ends. Location, branch, contact, manager, status list and other checks do not run.

## Proposal status-only change is checked separately

Lines 517-560

Asked as: change proposal status only; proposal status update checks; status change proposal validation; proposal status-only save

When only the status of a proposal changes, Turbo checks the status list, the value requirement, and the comment for No Bid or Cancelled, then skips the other checks. A proposal save where another field changes gets full validation.

## Allowed statuses for proposals

Lines 522-532

Asked as: invalid status for proposal; proposal can't be active; which proposal statuses; proposal status options

A proposal may be In Progress, Submitted, Awarded, Not Awarded, Cancelled or No Bid. Turbo refuses Active and Closed with "invalid status for proposal". A new proposal must start as In Progress or Submitted.

## Comment required for No Bid or Cancelled proposals

Lines 533-557

Asked as: no bid needs comment; cancel proposal comment required; comment to set status; why comment needed proposal

Before a proposal moves to No Bid or Cancelled, the job needs a client note whose status change matches that status. Otherwise Turbo says "a comment is required to set this status". This applies when the status actually changes.

## Date rules on create

Lines 563-594

Asked as: invalid dates for project; new project dates wrong; new proposal dates error; create project date rules

A new project needs a project award date and must have both proposal dates empty ("invalid dates for project"). A new proposal needs both proposal dates and no award date ("invalid dates for proposal").

## Date rules on update

Lines 595-647

Asked as: can't change award date; update project dates error; proposal dates changed error; editing project dates rejected

On update, Turbo checks dates only if the award, opening or due date changed. It then uses the new values, or the old ones where the new are blank. A project must keep an award date and no proposal dates. A proposal must keep both proposal dates and no award date.

## Location is required

Lines 649-659

Asked as: plus code location; invalid location job; missing location; location not accepted

Every job needs a location as a valid Plus Code. Otherwise Turbo says "invalid or missing location".

## Branch is required

Lines 661-670

Asked as: no branch on job; branch missing error; pick a branch; job branch needed

A job without a branch is refused with "branch is required".

## Contact and manager are required

Lines 672-714

Asked as: job needs contact; no manager on job; contact missing on save; manager missing error

A fully validated job must have a client contact and a manager. Turbo says "contact is required" or "manager is required". Turbo does not fill them in. A status-only retirement of a project, or a status-only proposal change, does not need them.

## Allowed statuses for projects

Lines 716-728

Asked as: invalid status for project; project in progress not allowed; project status options; project can't be submitted

A project may be Active, Closed or Cancelled. Turbo refuses In Progress, Submitted, Awarded, Not Awarded and No Bid with "invalid status for project".

## A project's value is checked against its status when Active or Closed

Lines 729-731

Asked as: project value zero error; active project value above zero; value greater than zero; project value must be positive

When a project's status is Active or Closed, Turbo checks that the project value is greater than zero. If it is zero or less, Turbo refuses with "project value must be greater than zero". Other project statuses skip this check, and the field named in the error depends on the caller (project_value on a full edit).

## Proposals cannot be Active or Closed

Lines 732-743

Asked as: proposal set active; proposal closed not allowed; active proposal error; proposal can't be closed

On a full edit or create, a proposal with status Active or Closed is refused with "invalid status for proposal". The message says proposals may be In Progress, Submitted, Awarded, Not Awarded, Cancelled or No Bid.

## New proposals must start In Progress or Submitted

Lines 744-754

Asked as: new proposal starting status; new proposals must start; proposal starts as awarded; first status of proposal

A new proposal with any other status is refused with "new proposals must start as In Progress or Submitted".

## Proposal value is checked on full edits

Lines 755-757

Asked as: proposal value required submitted; T&M required proposal; awarded proposal no value; proposal value needed on save

On a full edit or create, a proposal set Submitted, Awarded or Not Awarded needs a proposal value above zero or time and materials. Otherwise Turbo refuses with "proposal value or time and materials required".

## Comment check on full proposal edits

Lines 758-785

Asked as: comment on proposal no bid edit; full edit proposal comment; proposal edit needs comment; comment check proposal update

On a full proposal edit, the same comment is required when the status changes to No Bid or Cancelled. A new proposal is not checked.

## Referenced proposal must be Awarded

Lines 788-841

Asked as: proposal not awarded linked project; referenced proposal error; link project to proposal; proposal must be awarded

A project that points to a proposal needs a proposal that exists, is a proposal, and has status Awarded. In Progress or Submitted gives "proposal must be awarded". Not Awarded, Cancelled or No Bid gives "proposal has invalid status". A missing job gives "invalid referenced proposal", and a non-proposal gives "referenced job is not a proposal".

## Sub-job parent rules

Lines 843-889

Asked as: sub-job parent client; sub-job different client; parent proposal sub-job; sub-job parent rules

A sub-job's parent must exist and must not be a proposal ("proposals cannot have sub-jobs"). The parent must have a client. A different client is refused with "sub-job must have same client as parent". Turbo writes the parent's client onto the sub-job.

## Contact must belong to the client

Lines 891-923

Asked as: contact wrong client; contact not for client; contact requires client; contact from another client

A contact needs a client ("contact requires a client"). The contact must exist and belong to that client ("contact does not belong to client").

## Job number cannot change type

Lines 929-941

Asked as: job number type conflict; change proposal to project number; number type mismatch; renumber job type

On update, the type implied by the existing number must match the derived type. Otherwise Turbo says "job number type conflict".

## Projects need an invoicing profile

Lines 956-984

Asked as: no invoicing profile; profile required for project; invoicing profile error; billing profile missing

A project must name an invoicing profile. Otherwise Turbo says "invoicing profile is required". A status-only change to Closed or Cancelled is exempt. Setting a project Active still needs one. Turbo also checks that the contact and profile belong to the client.

## Status-only change ignores authorizing document

Lines 986-992

Asked as: authorizing document status change; PA change counts as status; status only ignores document; document change not status

When Turbo decides whether an update changes only the status, it ignores the authorizing document field.

## Project value required for Active or Closed

Lines 1010-1041

Asked as: active project zero value; time and materials project value; project value needed for active; closed project value zero

A project that is Active or Closed needs a project value greater than zero, including time and materials. Otherwise Turbo says "project value must be greater than zero". An existing project with a stored zero value may still be edited if its status stays the same. A value-restore workflow operation skips the check.

## Zero-value legacy projects can be edited without a value

Lines 1013-1029

Asked as: old project zero value edit; legacy project no value; stored zero value project; edit project without value

An existing project whose stored project value is 0, and still 0, can be edited if its status does not change. A status change, such as to Closed, still needs a value above zero. A project value restore through the project workflow skips this value check.

## New projects require a value for Active or Closed in the value check

Lines 1025-1041

Asked as: new project active zero value; create project value error; new active project no value; new closed project value

A new project set Active or Closed with a project value of zero or less is refused with "project value must be greater than zero". The zero-value exception applies only to existing projects.

## Proposal value or time and materials required

Lines 1042-1057

Asked as: proposal value required; submitted proposal no value; turn on time and materials; proposal needs value or T&M

A proposal that is Submitted, Awarded or Not Awarded needs a proposal value above zero or time and materials turned on. Otherwise Turbo says "proposal value or time and materials required".

## New projects need an active rate sheet

Lines 1066-1165

Asked as: rate sheet required; inactive rate sheet; invalid rate sheet project; new project rate sheet

A new project must have a rate sheet, or Turbo says "rate sheet is required for new projects". A rate sheet set on a save must exist ("invalid rate sheet") and be active ("rate sheet must be active"). Proposals skip this check.

## Rate sheet cannot change on a job save once set

Lines 1100-1130

Asked as: clear rate sheet; change rate sheet; rate sheet history; rate sheet locked

Once a project has a rate sheet, a job save cannot clear it ("rate sheet cannot be cleared"). A job save cannot change it either ("change the rate sheet from the rate sheet history"). A project with no rate sheet can be given one. Turbo records who set the first rate sheet.

## Managers must be active users

Lines 1167-1215

Asked as: manager inactive; alternate manager not active; manager must be active user; inactive user as manager

The manager and alternate manager must be active users. Otherwise Turbo says "manager must be an active user" or "alternate manager must be an active user".

## Top-level job number generation

Lines 1237-1282

Asked as: next project number; job number 9999 full; how new job number built; top-level number sequence

A new top-level job gets the current 2-digit year, a dash and a 4-digit sequence, one more than the highest sequence for that year. Proposals start with P. Both 3-digit legacy and 4-digit numbers count toward the highest. Turbo skips numbers already in use and fails if none is free up to 9999.

## Sub-job number generation

Lines 1284-1321

Asked as: sub-job number error; more than 99 sub-jobs; sub-job suffix; sub-job number already exists

A new sub-job gets its parent's number plus a 2-digit suffix: one more than the highest existing suffix among the parent's sub-jobs, starting at 01. Turbo refuses with "error generating sub-job number" at more than 99 sub-jobs, or if the number already exists. It reads at most 100 sub-jobs of a parent when it finds the highest suffix.
