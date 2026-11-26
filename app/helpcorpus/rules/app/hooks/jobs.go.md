<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/jobs.go","sha256":"0d87555c56f6882be1171e3dfde64ca45972ffdcf524eb76fe2bc49f71558370","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/jobs.go

## Only project retirement skips full validation for projects

Lines 20-36

Asked as: close project without validation error; retire project status only; cancel project skips checks; why project save skips validation; mark project closed quickly

A save that changes only a project's status to Closed or Cancelled counts as a retirement (isProjectRetirement) and skips the full job validation. Setting a project to Active never counts as a retirement. A save that also changes any other field (besides authorizing_document) gets full validation. The strict path (ProcessJobCoreStrict) always applies full validation, even to status-only edits.

## Job editing can be switched off

Lines 41-48

Asked as: cannot save job at all; job editing turned off; jobs are read only; can't create or update jobs; editing disabled message

When job editing is disabled in the app configuration, Turbo refuses every job create and update. It returns the error utilities.ErrJobsEditingDisabled. The setting is read by utilities.IsJobsEditingEnabled.

## Cancelled proposals cannot be changed

Lines 53-67

Asked as: cannot edit cancelled proposal; proposal locked after cancelling; change a cancelled proposal; cancelled proposal won't save; reopen cancelled proposal

Turbo refuses any change to a proposal whose saved status is Cancelled. The message is "cancelled proposals cannot be modified".

## Job numbers are assigned automatically

Lines 88-126

Asked as: how job numbers get assigned; new job has no number; automatic job number; unable to generate unique number; job number ran out

Turbo assigns a number on create. A top-level job gets the current two-digit year and the next 4-digit sequence, such as 26-0015, or P26-0015 for a proposal. The sequence follows the highest existing 3- or 4-digit number for that year. Turbo fails with "unable to generate a unique job number" if the sequence passes 9999.

## Sub-job number errors show as number-generation errors

Lines 96-126

Asked as: sub-job number error; error generating job number; number generation failed; cannot create sub-job number; job numbering error on save

When Turbo cannot generate a sub-job or top-level job number, it refuses with 'error generating sub-job number' or 'error generating job number'.

## Edited jobs are marked as not imported

Lines 140-144

Asked as: job flagged as not imported; job goes back to legacy system; edited job resent to legacy; imported flag cleared after edit; why job is re-exported

When an existing job is updated and any field changes, Turbo marks it as no longer imported. The job is then written back to the legacy system.

## Outstanding balance date is set automatically for projects

Lines 150-171

Asked as: outstanding balance date wrong; balance date changed automatically; why balance date updated today; project balance date not saved; balance date keeps old value

On a project, when the outstanding balance changes, Turbo sets the outstanding balance date to today. On create, any non-zero balance counts as a change. If the balance does not change on an update, the saved date is kept.

## Only jobs or payables_admin claim can change outstanding balance

Lines 173-244

Asked as: cannot change outstanding balance; no permission to edit balance; balance edit denied; need payables admin for balance; not logged in balance edit

Changing a job's outstanding balance needs the jobs claim (job) or the payables_admin claim. Without either, Turbo refuses with "insufficient permissions to edit outstanding balance". Without a signed-in user, Turbo says "authentication required to edit outstanding balance". Saving the same balance needs no claim. Creating a job skips this check.

## Active project needs a completion date

Lines 247-262

Asked as: active project completion date missing; project completion date required; active project won't save no end date; blank completion date error; estimated completion date for active

An Active project must have a project completion date. Otherwise Turbo says "project completion date is required" (an estimate is acceptable). Closed and Cancelled projects may leave it blank.

## Completion date must not be before award date

Lines 263-295

Asked as: completion date before award date; project end date earlier than award; completion date invalid compared to award; end date must follow award date; award date and completion date conflict

A project completion date must be a valid YYYY-MM-DD date. It may equal the project award date but not fall before it. Otherwise Turbo says "project completion date must be on or after project award date". The comparison runs only when the project has an award date, because legacy projects may have none.

## Proposal submission due date must not precede opening date

Lines 298-332

Asked as: proposal due date before opening; submission date earlier than opening; due date invalid for proposal; opening date after due date; proposal dates out of order

A proposal's submission due date may equal its opening date but not fall before it. Otherwise Turbo says "proposal submission due date must be on or after opening date". Both dates must be in YYYY-MM-DD format.

## Retiring a project skips full validation

Lines 351-361

Asked as: retire project without full checks; closing project skips invoicing check; project cancelled no contact needed; retiring project status change; why closing project needs no manager

When a save only changes an existing project's status to Closed or Cancelled, Turbo skips the full job validation. Retiring needs no invoicing profile, project value, contact or manager. Setting a project to Active does not qualify and gets full validation. The completion-date check on the new status still runs. Changing authorizing_document does not count as another change.

## Projects and proposals are cleaned on every save

Lines 363-400

Asked as: proposal fields cleared on save; project value reset to zero; authorizing document changed to PA; client PO spaces removed; why fields blank after saving

On a project, Turbo sets proposal value to 0 and authorizing document to PA. On a proposal, Turbo blanks authorizing document, client PO, client reference number, rate sheet, outstanding balance (set to 0), outstanding balance date, linked proposal, and project value (set to 0). Turbo trims spaces from client PO.

## Only top-level jobs can be parents of sub-jobs

Lines 418-450

Asked as: cannot add sub-job to sub-job; invalid base number; sub-job under sub-job; numbering under a sub-job; parent number not allowed

A sub-job number can be created only under a base job number: [P]YY-NNN or [P]YY-NNNN. A job that is itself a sub-job (such as 25-0123-01) cannot be a parent. Turbo refuses to number a sub-job under it with 'invalid base number'.

## Proposals are recognized by a leading P in the number

Lines 452-464

Asked as: how proposal is identified; what P prefix means; job starts with P; proposal number format; is this a proposal or project

Turbo treats a job as a proposal when its number starts with P, including a sub-job suffix such as P25-0123-01. Every other number is a project.

## New job type comes from its dates

Lines 478-500

Asked as: exactly one date configuration required; job type from dates; why job is project or proposal; new job type wrong; project or proposal dates mix

A new job is a project when it links a proposal, or when it has a project award date and no proposal dates. It is a proposal when it has both proposal dates and no award date. Any other mix gives "exactly one of project or proposal date configurations is required". An existing job's type comes from its number: a P prefix means proposal.

## Proposal status-only change skips full validation

Lines 507-561

Asked as: proposal status change only error; changing proposal status to active; invalid status for proposal change; status only edit proposal blocked; proposal status update fails

When only the status of a proposal changes, Turbo does not run the full validation. Active and Closed are refused with "invalid status for proposal". The value requirement still applies. No Bid or Cancelled needs a comment first. Fields that Turbo blanks on save can make the update count as a full edit.

## Date rules on create

Lines 563-594

Asked as: invalid dates for project; invalid dates for proposal; new project date error; new proposal dates rejected; create project date combination wrong

A new project must have a project award date and no proposal dates, or Turbo says "invalid dates for project". A new proposal must have both proposal dates and no award date, or Turbo says "invalid dates for proposal".

## Date rules on update only apply when a date changes

Lines 595-647

Asked as: date check on update; changing award date fails; edit proposal dates error; date combination when updating; project dates changed refused

On update, Turbo checks the date combination only if an award, opening or due date changed. It uses the new value, or the saved value where the new one is blank. A project must keep an award date and no proposal dates. A proposal must keep both proposal dates, no award date, and a due date that is not before the opening date.

## Location is required

Lines 649-659

Asked as: invalid or missing location; location plus code error; job location not accepted; plus code wrong; no location on job

Every job needs a location that is a valid Plus Code. Otherwise Turbo says "invalid or missing location".

## Branch is required

Lines 661-670

Asked as: branch missing on job; choose branch error; job has no branch; branch required message; cannot save without branch

Every job needs a branch. Otherwise Turbo says "branch is required".

## Contact and manager are required on full edits

Lines 672-714

Asked as: contact required on job; manager required error; who is the manager needed; job needs contact and manager; full edit missing contact

Creates and full edits need a contact and a manager. Turbo says "contact is required" or "manager is required". Status-only updates and project retirements skip this. Turbo never fills these fields in.

## Allowed statuses for projects

Lines 717-728

Asked as: invalid status for project; project status not allowed; awarded status on project; submitted project status refused; which statuses can a project have

A project may be Active, Closed or Cancelled. Awarded, Not Awarded, Submitted, In Progress and No Bid give "invalid status for project".

## Allowed statuses for proposals

Lines 732-754

Asked as: invalid status for proposal; proposal can't be active; new proposal start status; new proposals must start in progress; proposal closed status not allowed

A proposal may be In Progress, Submitted, Awarded, Not Awarded, Cancelled or No Bid. Active and Closed give "invalid status for proposal". A new proposal must start as In Progress or Submitted, or Turbo says "new proposals must start as In Progress or Submitted".

## No Bid or Cancelled proposal needs a comment first

Lines 758-784

Asked as: comment required to set status; no bid needs note; cancel proposal needs comment; client note before no bid; status change needs comment

Setting a proposal to No Bid or Cancelled needs a client note whose status-change-to matches that status. Otherwise Turbo says "a comment is required to set this status". The check runs only when the status changes, and not on create.

## Project linked to a proposal needs an Awarded proposal

Lines 788-841

Asked as: proposal must be awarded; linked proposal not awarded; project linked to proposal status; proposal has invalid status; link project to proposal error

A project's linked proposal must exist and be a proposal. If it is In Progress or Submitted, Turbo says "proposal must be awarded". If it is Not Awarded, Cancelled or No Bid, Turbo says "proposal has invalid status".

## Sub-job parent rules

Lines 843-889

Asked as: parent job rules for sub-job; proposals cannot have sub-jobs; sub-job client mismatch; different client than parent; sub-job parent missing

A sub-job's parent must exist and must not be a proposal ("proposals cannot have sub-jobs"). The parent must have a client. A different client gives "sub-job must have same client as parent". Turbo sets the sub-job's client to the parent's.

## Contact must belong to the job's client

Lines 891-927

Asked as: contact does not belong to client; contact requires client; wrong contact for client; client contact mismatch; invoicing profile client check

A job with a contact needs a client ("contact requires a client"). The contact must exist and belong to that client, or Turbo says "contact does not belong to client". Contact and invoicing profile are checked against the client by utilities.ValidateClientReferences. The invoicing profile is optional.

## Job type check on update uses the existing number

Lines 928-942

Asked as: job number type conflict; changed job type on update; cannot switch project to proposal; number type doesn't match save; job type change refused

On update, Turbo takes the job type from the existing job number. A number that starts with P is a proposal, and any other number is a project. A save that implies a different type is refused with 'job number type conflict'.

## Manager and alternate manager must be active

Lines 943-946

Asked as: manager not active user; alternate manager inactive; manager must be active; inactive user as manager; alternate manager must be active

The manager and the alternate manager must be active users. Otherwise Turbo says "manager must be an active user" or "alternate manager must be an active user".

## Invoicing profile is optional but must match the client

Lines 951-957

Asked as: invoicing profile must match client; optional invoicing profile error; invoicing profile wrong client; profile not for this client; remove invoicing profile

A job may have no invoicing profile. When one is set, Turbo checks through ValidateClientReferences that the contact and invoicing profile belong to the job's client.

## Proposals can use time and materials instead of a value

Lines 979-1030

Asked as: proposal value or time materials; proposal needs value for submitted; time and materials proposal required; proposal value missing on status change; awarded proposal zero value

A proposal set to Submitted, Awarded or Not Awarded needs a proposal value above zero or time and materials marked. Otherwise Turbo refuses with 'proposal value or time and materials required'. The error appears on the status field for a status-only change and on proposal_value for a full edit.

## Project value must be positive for Active or Closed

Lines 983-1014

Asked as: project value must be above zero; active project zero value; closed project value required; time and materials project value error; project value greater than zero

A project that is Active or Closed needs a project value above zero, including Time and Materials projects. Otherwise Turbo says "project value must be greater than zero". An older project whose saved value and new value are both zero can be edited if the status does not change. A value restore from the project workflow skips the check.

## Zero-value legacy projects can be edited without a new value

Lines 986-1002

Asked as: zero value project save; legacy project with zero value; edit old project value zero; restore project value; value zero edit blocked

An existing project with a stored project value of 0 can be saved with its value still 0 and its status unchanged, without the 'project value must be greater than zero' error. A status change to Active or Closed still needs a value above zero. Restoring a value through the project workflow (ProjectWorkflowMutationRestoreValue) skips this check.

## Submitted, Awarded or Not Awarded proposal needs value

Lines 1015-1028

Asked as: proposal value needed for awarded; submitted proposal without value; time and materials not ticked; not awarded proposal value missing; proposal value or T and M

A proposal with status Submitted, Awarded or Not Awarded needs a proposal value above zero or Time and Materials turned on. Otherwise Turbo says "proposal value or time and materials required".

## New project needs an active rate sheet

Lines 1039-1059

Asked as: rate sheet required new project; rate sheet must be active; no rate sheet on new project; active rate sheet needed; new project rate sheet error

A new project must have a rate sheet, or Turbo says "rate sheet is required for new projects". The rate sheet must exist and be active, or Turbo says "rate sheet must be active". Proposals skip this.

## Proposals do not use rate sheets

Lines 1043-1046

Asked as: proposal rate sheet cleared; why proposal has no rate sheet; rate sheet not used on proposals; proposal rate sheet removed; rate sheet checks skipped proposal

Turbo skips the rate sheet checks for proposals. It clears the rate sheet on proposals.

## Rate sheet cannot change through a job save

Lines 1073-1098

Asked as: rate sheet cannot be cleared; change rate sheet from history; cannot change rate sheet on save; rate sheet locked on project; set rate sheet on project

Once a project has a rate sheet, a job save cannot clear it ("rate sheet cannot be cleared"). A job save also cannot change it ("change the rate sheet from the rate sheet history"). A project with no rate sheet can be given one by any job editor.

## Rate sheet setter is recorded on the job

Lines 1100-1103

Asked as: who set the rate sheet; rate sheet setter recorded; first rate sheet user; rate sheet actor; record who chose rate sheet

When a job save sets or changes the rate sheet and the user is signed in, Turbo records that user as the person who set the first rate sheet (rateSheetActorKey).

## Rate sheet must exist when set on a new project

Lines 1105-1138

Asked as: invalid rate sheet; rate sheet not found; inactive rate sheet chosen; rate sheet must be active error; setting rate sheet on project fails

When a project sets a rate sheet, Turbo looks it up. It refuses with 'invalid rate sheet' if the sheet is not found. It refuses with 'rate sheet must be active' if the sheet is inactive. Turbo checks only when the rate sheet is being set or changed.

## Alternate manager is optional but must be active

Lines 1140-1188

Asked as: alternate manager optional check; inactive alternate manager; alternate manager not active; manager inactive refused; backup manager must be active

The alternate manager is optional. When one is set, Turbo checks that the user is active. If not, it refuses with 'alternate manager must be an active user'. An inactive manager is refused with 'manager must be an active user'.

## Top-level job numbers count up within the year

Lines 1190-1255

Asked as: next job number for year; how top-level numbers count up; new project number sequence; job number highest plus one; out of job numbers this year

A new top-level job gets the next number for the current two-digit year: YY-NNNN for a project or PYY-NNNN for a proposal. Turbo takes the highest existing 3- or 4-digit sequence for that prefix and adds 1. It always writes new numbers with four digits. Sub-jobs are not counted. If no unused number remains up to 9999, number generation fails with an error.

## Sub-job numbers follow the parent

Lines 1257-1294

Asked as: sub-job number suffix; how sub-job numbers are made; parent already has suffix; too many sub-jobs; sub-job limit 99

A sub-job number is the parent's number plus a two-digit suffix, one above the highest existing suffix. A parent that already has a suffix cannot have sub-jobs. A parent allows at most 99 sub-jobs.
