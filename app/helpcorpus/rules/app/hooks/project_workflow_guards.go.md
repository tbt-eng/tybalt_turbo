<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/project_workflow_guards.go","sha256":"7a37bd136ecaedd1f74249013b19c86eb0f51494a12ea3e9387b2aa9e106f05f","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/project_workflow_guards.go

## Controlled job workflow transitions still get full job validation

Lines 81-86

Asked as: workflow action validation; job action tag; strict job validation on workflow; why job validation runs on approve

When Turbo changes a job through a workflow action (create job, approve value, restore value, merge, change client, rate sheet), it tags the job with that action and runs the strict job validation (ProcessJobCoreStrict). The action tag lets the guards allow the matching change, such as a value change on approve_value. The normal validation still applies, and the tag is cleared afterwards.

## Authentication required for claim-restricted workflow actions

Lines 90-109

Asked as: not signed in; no user logged in; authentication required error; inactive user cannot act

If no signed-in user is given, Turbo refuses with 'authentication is required'. An inactive user, or a user without any of the required claims, gets 'an active user with the <claims> claim is required'.

## Claim check names the accepted claims in its refusal

Lines 99-108

Asked as: missing claim error; claim required message; which claim do I need; or claim refused

When the actor is inactive or lacks every listed claim, Turbo refuses with "an active user with the X or Y claim is required", listing the claims joined by "or". The actor needs only one of the listed claims. Turbo checks the actor's active status as well as the claim.

## Project writes need an active user with the job claim

Lines 111-120

Asked as: cannot save project without job claim; procurement permission needed; project edit not allowed; need procurement access

When a project (not a proposal) is written, Turbo requires an active user with the job claim (procurement). Otherwise it refuses with 'an active user with the job claim is required'. Approving or restoring a project value skips this check. Proposals skip it too.

## Project value must be whole dollars

Lines 121-123

Asked as: cents in project value; decimal project amount; project value rounding error; dollars and cents not allowed

Turbo refuses a project's project_value that has cents, with 'project value must use whole dollars'. Proposals are not checked here.

## New projects need an approved project request when approval is required

Lines 124-132

Asked as: cannot create project directly; new project blocked; project request approval needed; create project without request

When job request approval is required (IsJobRequestApprovalRequired), Turbo refuses to create a project directly with 'create projects through an approved project request'. Only the create-job step of the project workflow may create it.

## Project client cannot be changed by a normal edit

Lines 133-136

Asked as: cannot change project client; client locked on project; wrong client on project; move project to another client

Turbo refuses a change to a project's client with 'project client is immutable' (and 'project client can change only through a client merge' in the validation helper). Only a client merge may change it, and the change-client action may as well in the guard on save.

## Project value changes require approval

Lines 137-139

Asked as: cannot edit project amount; changing project value blocked; project value locked; edit budget without approval

Turbo refuses to change a project's value with 'project value changes require approval'. Only approving a value change or restoring a previous value may change it.

## Cannot delete a user who appears in project authorization records

Lines 146-160

Asked as: cannot delete user; remove user account failed; user still on request; delete staff member error

Turbo refuses to delete a user who is the creator, approver or rejector of a project request, the uploader, approver or rejector of an attachment, or the request creator or creator of a job. The message is 'user identity is retained in project authorization records'.

## Cannot delete a user who manages a job

Lines 161-199

Asked as: cannot delete manager; user manages a job; deleting user from records api fails; remove manager from job first

Turbo refuses to delete a user who is the manager of any job, project or proposal, with 'user is the manager of one or more jobs'. A delete through the records API shows PocketBase's generic delete error instead.

## Request checklist versions must be sequential

Lines 200-218

Asked as: checklist version wrong; next checklist number; request checklist numbering; version out of sequence

A new request checklist for a scope must have a version one higher than the latest version of that scope. Otherwise Turbo refuses with 'the next <scope> request checklist version is N'. Turbo also refuses a checklist whose questions are invalid, giving the error message.

## Saved request checklists cannot be edited or deleted

Lines 219-223

Asked as: cannot edit saved checklist; change finished checklist; delete checklist blocked; update checklist after saving

Turbo refuses to update or delete a saved request checklist with 'a saved request checklist cannot change; add the next version'.

## Project requests, allocations and attachments need workflow actions

Lines 224-253

Asked as: cannot edit project request directly; allocation edit blocked; attachment save refused; use workflow button instead

Creating or updating a project request, request allocation or job attachment outside the project workflow actions is refused with 'use the project workflow actions'. Deleting one outside the workflow is refused with 'use the project workflow delete action'.

## Attachment must belong to exactly one request or job

Lines 231-235

Asked as: attachment linked to two things; attachment has no job; attachment request and job both; orphan attachment

A job attachment must have either a request or a job, not both and not neither. Otherwise Turbo refuses with 'an attachment must belong to one request or job'.

## Approved project evidence cannot be edited

Lines 236-238

Asked as: cannot edit approved attachment; approved evidence locked; change approved document; edit approved file

Turbo refuses to edit an existing attachment of a job whose status is Approved, with 'approved project evidence cannot be edited'.

## Attachment values must be whole dollars

Lines 239-244

Asked as: cents on attachment amount; attachment value decimals; previous value cents; approval max cents

An attachment's project_value, previous_project_value and approval_max_value must be whole dollars. Otherwise Turbo refuses with '<field> must use whole dollars'.

## Initial project approval cannot be deleted

Lines 254-257

Asked as: cannot delete first approval; remove initial approval; original approval locked; delete original approval attachment

Turbo refuses to delete an approved job attachment that is the initial approval, with 'initial project approval cannot be deleted'.

## Delete later or pending project changes first

Lines 258-267

Asked as: cannot delete approved attachment; later change blocks delete; pending changes block removal; delete newer approval first

Turbo refuses to delete an approved project attachment while the job has any non-approved attachment or an approved one with a higher approval sequence. The message is 'delete pending or later project changes first'.

## Job number cannot switch between project and proposal

Lines 274-278

Asked as: change project to proposal; job type switch; number changed to proposal; convert proposal to project

Turbo refuses to edit a job so that its number changes between project and proposal type, with 'a job number cannot switch between project and proposal'.

## Project writes must use the validated job action

Lines 279-288

Asked as: project save refused; save project blocked; project write rejected; must use job action

A project save outside a workflow action must match a validated state. Otherwise Turbo refuses with 'project writes require the validated job action'. Proposals skip these guards.

## Project job saves are refused outside the workflow

Lines 282-287

Asked as: cannot save project directly; project job save refused; edit project record fails; validated state stale

When a project job is saved with no workflow operation and no matching validated state, Turbo refuses with "project writes require the validated job action". The validated state counts only if the record is unchanged since validation. Turbo clears that state after each save.

## Project value changes in the job guard need approval

Lines 289-315

Asked as: project value edit refused; change value on save; value cents on job; job amount change blocked

A save that changes an existing project's value is refused with "project value changes require approval". Only the approve-value and restore-value workflow actions can change it. The job guard also refuses a value that is not whole dollars with "project value must use whole dollars".

## Value restoration only to the latest approved previous value

Lines 292-296

Asked as: restore previous value fails; undo value change; revert project amount; restore old budget

Restoring a project value works only when the latest approved change (by approval sequence) is not the initial approval, and its previous value is the restored value and its value is the current one. Otherwise Turbo refuses with 'restore only the latest approved change's stored previous value'.

## New project creation needs an approved request when approval is required

Lines 298-306

Asked as: new project needs request; approval setting blocks project; project creation refused setting; no approved request for project

When the setting that requires job request approval is on, creating a new project outside the create-job workflow is refused with "create projects through an approved project request". This check also runs in the job guard.

## Project job client change refused unless merge or change-client

Lines 307-312

Asked as: change project client refused; client merge needed; switch client on job; project owner client change

A save that changes an existing project's client is refused with "project client is immutable". The client merge and the change-client workflow actions can change it.

## Project rate sheet changes go through rate sheet history

Lines 316-320

Asked as: cannot change project rate sheet; rate sheet edit blocked; update project rates; swap rate sheet

Once a project has a rate sheet, a normal edit cannot change it. Turbo refuses with "change a project's rate sheet from its rate sheet history". Setting the first rate sheet is a normal edit.

## Project creation history is protected

Lines 322-326

Asked as: creation history changed; change created by; edit project creator; project audit fields locked

Turbo refuses a change to a project's request_creator, request_created, request_submitted, created_by or creation_completed, with 'project creation history is protected'. Only the create-job step may set them.

## Job with attachments cannot be deleted

Lines 331-339

Asked as: cannot delete job; job delete blocked; remove job with files; job still has attachments

Turbo refuses to delete a job that has any attachments, with 'delete permitted project attachments before deleting this job'.

## Approval limit must be positive whole dollars

Lines 344-349

Asked as: approval limit too high; max value invalid; authorization limit cents; negative approval limit

A project authorization limit (max_value) must be a positive whole number of dollars of at most 9007199254740991. Otherwise Turbo refuses with 'approval limit must be positive whole dollars no greater than 9007199254740991'.

## Approval limits attach only to project_authorization claim

Lines 350-361

Asked as: approval limit not saved; limit on wrong claim; properties not saving claim; set limit for user claim

Turbo saves approval limit properties only for a user claim assignment of the project_authorization claim. Otherwise it refuses with 'properties require a project_authorization assignment'.

## Project request branch must exist and be active

Lines 374-404

Asked as: branch missing on request; inactive branch request; project request branch error; choose active branch

A project request's branch must exist, otherwise Turbo says 'branch is required'. An inactive branch gets 'branch must be active'. Turbo then runs full job validation on the future job as a project with the Project Authorization (PA) document.

## Project request needs a First Nations agreement answer

Lines 379-382

Asked as: first nations question required; agreement answer missing; FN answer blank; request check fails first nations

When a project request is checked, its First Nations agreement answer must be Yes or No. Otherwise Turbo refuses with 'First Nations agreement answer is required'.

## Project request rate sheet check uses the draft

Lines 383-395

Asked as: rate sheet not active; proposed rates check; draft rates error; no active rate sheet on request

If the request has proposed rates, Turbo validates the draft rates instead of requiring an active rate sheet. Without proposed rates, the rate sheet must be active.

## Pending value changes block restoration

Lines 418-425

Asked as: cannot restore value; pending change blocks restore; restore refused attachments; unapproved changes exist

Turbo refuses to restore a project value while the job has any non-approved attachment, with 'delete pending project value changes before restoration'.
