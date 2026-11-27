<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/project_workflow_guards.go","sha256":"487e2e3287c68b415bf35fb313aabd770e2ac9252e1a22f5c5bf10d497bd6083","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/project_workflow_guards.go

## Workflow actions run the strict job validation

Lines 81-86

Asked as: action saves job wrong check; workflow save strict validation; why job save fails on action

A project workflow action saves the job through the strict job validation (ProcessJobCoreStrict) and tags the save with the action's name. The project guards use that name to decide which changes the action may make. Turbo clears the tag after the save.

## Actions that need a claim refuse users without it

Lines 90-106

Asked as: authentication is required error; need to sign in; missing claim error; not allowed to approve project

When a workflow action needs a claim and nobody is signed in, Turbo refuses with 'authentication is required'. When the user holds none of the needed claims, Turbo refuses with 'a user with the <claim> or <claim> claim is required'. One matching claim is enough.

## Project writes need the job claim, except value approval and restore

Lines 107-116

Asked as: cannot save project without permission; need job permission to save; project save permission denied

Saving a project (not a proposal) requires a user with the job claim. Approving or restoring a project value is exempt from this claim check. Proposals skip these project checks.

## Project value must be whole dollars

Lines 117-119

Asked as: project value has cents; cents in project amount; decimal project value error; round project value

Turbo refuses a project save when the project value has cents. It says "project value must use whole dollars".

## Projects must come from an approved request when approval is required

Lines 120-128

Asked as: cannot create new project; project creation blocked; need approved request for project; new project not allowed

When project request approval is required (IsJobRequestApprovalRequired), a new project can be created only by the create-job workflow step. Otherwise Turbo says "create projects through an approved project request".

## Project client changes only through a client merge

Lines 129-132

Asked as: cannot change project client; change client on project; client locked on project; switch project client

On a project (not a proposal), a normal edit that changes the client is refused with 'project client can change only through a client merge'. Only the merge action passes this check.

## Project value changes need approval

Lines 133-135

Asked as: cannot change project value; edit project amount refused; project value locked; value change blocked

On a project, an edit that changes the project value is refused with 'project value changes require approval'. Only value approval and value restore can change it.

## A user who is tied to project records cannot be deleted

Lines 142-156

Asked as: cannot delete user; delete user error; user tied to projects; remove user fails

Turbo refuses to delete a user who is the creator, approver or rejector of a project request, the uploader, approver or rejector of an attachment, or the request creator or creator of a job. The message is 'user identity is retained in project authorization records'.

## A user who manages a job cannot be deleted

Lines 157-194

Asked as: cannot delete manager; user manages a job; delete user who manages job; manager deletion blocked

Turbo refuses to delete a user who is the manager of any job, project or proposal. The message is 'user is the manager of one or more jobs'. The check runs before the delete, so the manager is never blanked.

## Request checklist versions must be sequential

Lines 198-214

Asked as: checklist version wrong; next checklist version number; request checklist numbering error; checklist version rejected

A new request checklist must have a version number equal to the latest version for its scope plus one. Otherwise Turbo refuses with 'the next <scope> request checklist version is N'. Turbo also refuses a checklist whose questions are not valid.

## Saved request checklists cannot be edited or deleted

Lines 215-219

Asked as: cannot edit saved checklist; change saved checklist; delete checklist refused; update checklist locked

Turbo refuses to update or delete a saved request checklist with 'a saved request checklist cannot change; add the next version'.

## Project requests, allocations and attachments need workflow actions

Lines 221-249

Asked as: cannot edit project request directly; use workflow for attachments; project request save refused; delete attachment outside workflow

Creating or updating a project request, request allocation or job attachment outside the project workflow actions is refused with 'use the project workflow actions'. Deleting one outside the workflow is refused with 'use the project workflow delete action'.

## Attachment must belong to exactly one request or job

Lines 227-231

Asked as: attachment linked to both; attachment has no request or job; attachment link error; attachment belongs to nothing

A job attachment must link to either a request or a job, not both and not neither. Otherwise Turbo refuses with 'an attachment must belong to one request or job'.

## Approved project evidence cannot be edited

Lines 232-234

Asked as: cannot edit approved evidence; approved attachment locked; change approved attachment; evidence edit refused

An attachment that belongs to a job and has status Approved cannot be edited. Turbo refuses with 'approved project evidence cannot be edited'.

## Attachment values must be whole dollars

Lines 235-240

Asked as: attachment value has cents; attachment amount not whole dollars; previous value cents error; approval max cents

On a job attachment, project_value, previous_project_value and approval_max_value must be whole dollars. Otherwise Turbo refuses with '<field> must use whole dollars'.

## Initial project approval cannot be deleted

Lines 250-253

Asked as: cannot delete first approval; remove initial approval; initial approval locked; delete original approval

Turbo refuses to delete the initial approval attachment of a job, with 'initial project approval cannot be deleted'.

## Approved value changes are deleted newest first

Lines 254-263

Asked as: cannot delete approved change; delete older approval first; pending changes block delete; remove approved value change

An approved attachment cannot be deleted while the job has pending (not Approved) attachments or approved ones with a higher approval sequence. Turbo refuses with 'delete pending or later project changes first'.

## Job number cannot switch between project and proposal

Lines 270-274

Asked as: cannot change job number type; switch project to proposal; proposal to project number; job type change blocked

Editing a job so its number changes from project type to proposal type, or the reverse, is refused with 'a job number cannot switch between project and proposal'.

## Project writes must use the validated job action

Lines 275-283

Asked as: project save outside workflow; save project directly refused; validated job action needed; project edit not validated

Saving a project outside a workflow action is refused with 'project writes require the validated job action', unless it was validated and unchanged since. Proposals skip these project guards.

## Project value must be whole dollars at save

Lines 285-287

Asked as: project value cents on save; save fails cents; project amount decimal save; whole dollar error on save

When Turbo saves a project, a project value with cents is refused with 'project value must use whole dollars'.

## Only the create-job action makes projects when approval is required

Lines 294-302

Asked as: create project outside workflow; new project from job screen; project creation blocked by approval; make project without request

When project request approval is required, creating a project outside the create-job action is refused with 'create projects through an approved project request'.

## Project client is immutable except merge or client change

Lines 303-308

Asked as: project client locked; client cannot change on save; change client save error; client change refused

On save, changing a project's client is refused with 'project client is immutable', except by a client merge or the change-client action.

## Project value changes at save need approval actions

Lines 309-311

Asked as: value change on save refused; project amount change blocked; save changes project value; value changed without approval

On save, changing a project's value is refused with 'project value changes require approval', except by value approval or restore.

## Set rate sheet changes only through rate sheet history

Lines 312-316

Asked as: cannot change rate sheet; rate sheet locked on project; edit rate sheet refused; set rate sheet history

Once a project has a rate sheet, changing it in a job edit is refused with 'change a project's rate sheet from its rate sheet history'. Setting the first rate sheet is a normal edit.

## Project creation history is protected

Lines 318-322

Asked as: cannot change creation history; project creator locked; edit created by date; creation details protected

Changing request_creator, request_created, request_submitted, created_by or creation_completed on a project is refused with 'project creation history is protected', except in the create-job action.

## Job with attachments cannot be deleted

Lines 327-339

Asked as: cannot delete job with attachments; job delete blocked attachments; remove job attachments first; delete job error attachments

Turbo refuses to delete a job that still has project attachments, with 'delete permitted project attachments before deleting this job'.

## Approval limit must be positive whole dollars

Lines 340-345

Asked as: approval limit invalid; max value error; authorization limit too large; limit must be positive

A project authorization limit (max_value) must be a whole number above 0 and no more than 9007199254740991. Otherwise Turbo refuses with 'approval limit must be positive whole dollars no greater than 9007199254740991'.

## Approval limits attach only to project_authorization claim

Lines 346-360

Asked as: properties require assignment; approval limit on wrong claim; authorization property error; limit only for project authorization

Authorization properties can be saved only for a user claim assignment of the project_authorization claim. Otherwise Turbo refuses with 'properties require a project_authorization assignment'.

## Project request is checked as the future job

Lines 363-374

Asked as: project request check fails; request job preview wrong; project request job fields; project status on request

Turbo checks a project request by building the job it would create, without assigning a job number. It copies the request's client, owner, contact, managers, branch, rate sheet, parent, proposal, dates, value and other job fields to that job. It sets the job's status from the request's project status (project_status) and its authorizing document to "PA" (Project Authorization).

## Project request validation: open proposal ignored

Lines 375-386

Asked as: proposal in progress blocks request; open proposal request error; submit request before proposal awarded; proposal not awarded check

When Turbo validates a project request, a linked proposal whose status is In Progress or Submitted is left out of the check. A request can be submitted and approved before the proposal is Awarded.

## First Nations agreement answer is required

Lines 387-390

Asked as: First Nations answer missing; FN agreement required; project request First Nations; indigenous agreement question

A project request must have a First Nations agreement answer of Yes or No. Otherwise Turbo refuses with 'First Nations agreement answer is required'. The answer sets the job's fn_agreement.

## Request validation checks draft rates or active rate sheet

Lines 391-403

Asked as: rate sheet not active; draft rates check request; request rates invalid; inactive rate sheet error

If the request has proposed rates, Turbo validates that draft. Otherwise Turbo checks the request's rate sheet is active.

## Request branch must exist and be active

Lines 404-410

Asked as: branch missing on request; request branch inactive; no branch on project request; branch not active error

A project request needs a branch that exists, else 'branch is required'. If the branch is inactive, Turbo refuses with 'branch must be active'.

## Value restore only to latest approved change's previous value

Lines 416-425

Asked as: cannot restore project value; restore previous value error; undo value change refused; restore value wrong amount

Restoring a project value works only when the latest approved attachment (by approval sequence) is not the initial approval. The job's new value must equal its stored previous value, and its current value must equal its approved value. Otherwise Turbo refuses with 'restore only the latest approved change's stored previous value'.

## Pending value changes block restoration

Lines 426-432

Asked as: cannot restore value; pending changes stop restore; restore blocked by pending; delete pending before restore

Turbo refuses a value restore while the job has any attachment that is not Approved. The message is 'delete pending project value changes before restoration'.
