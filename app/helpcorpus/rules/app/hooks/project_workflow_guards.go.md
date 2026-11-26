<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/project_workflow_guards.go","sha256":"c15c9696b51125fafbfdbeb2a103cd2b8067061fd8a60e8cf6d2df7747062fc5","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/project_workflow_guards.go

## Project actions need an active user with the right claim

Lines 90-110

Asked as: cannot save project, not signed in; inactive user blocked from project; missing permission for project action; authentication is required error

The activeClaim check refuses an action when nobody is signed in, with the message "authentication is required". It also refuses when the user is inactive or lacks every listed claim. The message is then "an active user with the <claims> claim is required", and the claims are joined with " or ".

## Project writes need an active user with the job claim

Lines 111-120

Asked as: cannot save project without procurement access; need job permission to edit project; save project error permission

When someone saves a project, Turbo requires an active user with the job claim (procurement). Otherwise it refuses with 'an active user with the job claim is required'. Approving or restoring a project value skips this check. Proposals skip it too.

## Proposals skip the project write checks

Lines 111-114

Asked as: proposal saving blocked; proposal wrongly checked as project; save proposal permission error

When the job is a proposal, validateProjectJobWrite returns at once. A proposal does not need the job claim, the whole-dollar value, the approval rule or the client and value change rules from this function.

## Approving or restoring a value needs no job claim

Lines 115-120

Asked as: approve value without job permission; restore value permission problem; approving project value blocked

The job claim is required for a project write. The exception is a value approval or a value restoration (approve_value, restore_value). Those two operations skip the active-user and job-claim check.

## Project value must be whole dollars

Lines 121-123

Asked as: project value has cents; project amount with decimals; cents not allowed in value; round project value to dollars

Turbo refuses a project whose project value (project_value) has cents. The message is 'project value must use whole dollars'. Proposals are not checked here.

## Projects need an approved project request when approval is required

Lines 124-132

Asked as: cannot create project directly; new project blocked by approval setting; create project without request; project request approval required

When project request approval is switched on (IsJobRequestApprovalRequired), Turbo refuses to create a project directly. The message is 'create projects through an approved project request'. Only the create-job workflow step may create it.

## Project client changes only by client merge or change-client action

Lines 133-136

Asked as: change project client; switch client on project; project client won't change; merge clients on project

On an existing project, a change of client is refused. Only a client merge or the change-client action can change it. The write validation says 'project client can change only through a client merge'. The job save guard says 'project client is immutable'.

## Project value changes need approval

Lines 137-139

Asked as: cannot edit project value; change project amount blocked; project value locked; update value needs approval

On an existing project, Turbo refuses a changed project value unless the change comes from approving or restoring a value. The message is 'project value changes require approval'.

## Deleting a user is refused if tied to project records

Lines 146-160

Asked as: cannot delete user account; remove user who created requests; user has project records; delete user fails

Turbo refuses to delete a user who appears in project authorization records. These are the creator, approver or rejector on project requests, the uploader, approver or rejector on attachments, and the request creator or creator on jobs. The message is 'user identity is retained in project authorization records'.

## Deleting a user who manages a job is refused

Lines 191-198

Asked as: delete user who manages job; user is job manager; cannot remove job manager; delete user generic error

Turbo refuses to delete a user who is the manager of any job, project or proposal. The message is 'user is the manager of one or more jobs'. A delete through the records API shows PocketBase's generic delete error instead.

## Request checklist versions must be the next number

Lines 202-218

Asked as: checklist version wrong number; next checklist version; new checklist version rejected; request checklist version number

A new request checklist must have a version one higher than the latest version for its scope. Otherwise Turbo refuses with 'the next <scope> request checklist version is N'. Turbo also refuses a checklist whose questions are invalid, giving the error message. Migrations that save the first versions before the scope field exists skip the version check.

## Request checklist version check is skipped without a scope field

Lines 203-216

Asked as: checklist version check skipped; old checklist migration error; checklist questions invalid

The next-version rule applies only when the request_checklists collection has a scope field. A migration saves the first versions before that field exists. The checklist's questions are also validated, and an error from that check is shown as given.

## Saved request checklists cannot be edited or deleted

Lines 219-223

Asked as: edit saved checklist; delete request checklist; changing old checklist blocked; update saved checklist

Turbo refuses to update or delete a saved request checklist. The message is 'a saved request checklist cannot change; add the next version'.

## Project requests, allocations and attachments need workflow actions

Lines 225-252

Asked as: cannot edit project request directly; use workflow for allocations; delete attachment outside workflow; project request edit blocked

Turbo refuses creating or updating job requests, request allocations and job attachments outside the project workflow actions. The message is 'use the project workflow actions'. Deleting them outside the workflow gives 'use the project workflow delete action'.

## Attachment must belong to exactly one request or job

Lines 231-235

Asked as: attachment linked to two jobs; attachment with no job; attachment must have one link; attachment belongs to request and job

A project attachment must link to a request or to a job, never both and never neither. Otherwise Turbo refuses with 'an attachment must belong to one request or job'.

## Approved project evidence cannot be edited

Lines 236-238

Asked as: cannot edit approved attachment; approved evidence locked; change approved document blocked

Once an attachment on a job is Approved, Turbo refuses any edit. The message is 'approved project evidence cannot be edited'.

## Attachment values must be whole dollars

Lines 239-244

Asked as: attachment value has cents; previous value decimals; approval max cents error; attachment amount not whole dollars

Turbo refuses an attachment where project_value, previous_project_value or approval_max_value has cents. The message is '<field> must use whole dollars'.

## Initial project approval cannot be deleted

Lines 254-257

Asked as: cannot delete initial approval; remove first approval attachment; delete original approval blocked

Turbo refuses to delete the initial approval attachment of a job, with the message 'initial project approval cannot be deleted'.

## Approved value changes are deleted newest first

Lines 258-267

Asked as: delete approved change order; pending changes block deletion; delete newest project change first; cannot delete older approval

To delete an approved project attachment, nothing may be pending or not Approved on that job, and no later approved change may exist (higher approval_sequence). Otherwise Turbo says 'delete pending or later project changes first'.

## A job number cannot switch between project and proposal

Lines 274-278

Asked as: change job number to proposal; turn project into proposal; switch job number type; project proposal number change

Turbo refuses to change a job's number so that it turns a project into a proposal or the reverse. The message is 'a job number cannot switch between project and proposal'.

## Project saves must go through the validated job action

Lines 279-287

Asked as: project save blocked outside workflow; project write needs validation; save project directly error

Turbo refuses to save a project outside a workflow action unless the validated job action ran on the same data. The message is 'project writes require the validated job action'. Proposals skip these guards.

## Proposals pass the project save guard

Lines 279-280

Asked as: proposal save guard; proposal skips project checks; save proposal changes value

When the job is a proposal, the save guard goes straight to the next step. The project rules on value, client, rate sheet and creation history do not apply. The job number check on switching type still applies first.

## Job save needs a workflow action or a validated job

Lines 282-287

Asked as: project save validation error; save project after edit; validated state mismatch; save needs revalidation

Turbo refuses a project save with "project writes require the validated job action". This happens when no workflow operation is set and the record's validated state does not match its current field values. The validated state is cleared after each save, so each save needs a new validation.

## Project job save re-checks whole dollars

Lines 289-291

Asked as: job save cents error; project value decimals on save; whole dollars on job save

At the job save, Turbo refuses a project value with cents. The message is 'project value must use whole dollars'.

## Value restore only for latest approved change

Lines 292-296

Asked as: restore old project value; restore previous value rejected; only latest change can restore; restore value wrong amount

A restore of a project value is allowed only for the latest approved change (by approval sequence) that is not the initial approval. The restored value must be that change's stored previous value, and the current value must be its project value. Otherwise Turbo says 'restore only the latest approved change's stored previous value'.

## New project is refused outside an approved project request

Lines 298-306

Asked as: cannot add new project; create project blocked; new project needs approved request

When the setting that requires approval of job requests is on, a new project can be created only by the create_job operation. Any other new project is refused with "create projects through an approved project request". The save guard applies this check to new projects.

## Project client is immutable except merge or change client

Lines 307-312

Asked as: change project client blocked; client locked on project; project client immutable

On an existing project, the save guard refuses a changed client with "project client is immutable". It allows the change only for a client merge or the change-client action.

## Project value change on save needs approval

Lines 313-315

Asked as: project value save blocked; change value on save refused; value change needs approval save

On an existing project, the save guard refuses a changed project value with "project value changes require approval". Only the approve_value and restore_value operations may change it.

## Project rate sheet changes only through rate sheet history

Lines 316-320

Asked as: cannot change rate sheet; edit project rate sheet; rate sheet locked on project; update rate sheet blocked

Once a project has a rate sheet, Turbo refuses to change it in a normal edit. The message is 'change a project's rate sheet from its rate sheet history'. Setting the first rate sheet is a normal edit.

## Project creation history is protected

Lines 322-326

Asked as: cannot change creation history; edit project created by; request dates locked; creation fields protected

Turbo refuses to change the request creator, request created, request submitted, created by or creation completed fields on a project, except while creating the job from a request. The message is 'project creation history is protected'.

## Job with attachments cannot be deleted

Lines 331-343

Asked as: cannot delete job with attachments; delete job blocked by files; remove attachments before job

Turbo refuses to delete a job that still has attachments. The message is 'delete permitted project attachments before deleting this job'.

## Approval limit must be positive whole dollars

Lines 344-349

Asked as: approval limit invalid; max value must be positive; authorization limit too large; approval limit whole dollars

A project authorization limit (max_value) must be above 0, a whole number of dollars, and at most 9007199254740991. Otherwise Turbo refuses with 'approval limit must be positive whole dollars no greater than 9007199254740991'.

## Approval limits only on project_authorization claim

Lines 350-364

Asked as: approval limit on wrong claim; limit property rejected; project authorization limit only

Turbo saves approval limit properties only for a user claim assignment of the project_authorization claim. Otherwise it says 'properties require a project_authorization assignment'.

## Create-job request data is validated as a future project

Lines 368-417

Asked as: create job from request fails; request validation error job; future project validation failed

ValidateProjectRequestData builds a draft job from the request's fields. It sets the job status from the request's project_status. It sets the authorizing document to "PA" and sets First Nations agreement to true only when the answer is "Yes". It then runs the full job validation as a create.

## Open proposals are ignored when checking a project request

Lines 379-390

Asked as: proposal in progress blocks request; open proposal request check; proposal not awarded request

When a project request names a proposal that is In Progress or Submitted, Turbo leaves the proposal out of the request check. The proposal must be Awarded only when the job is created.

## Project request needs a First Nations agreement answer

Lines 391-394

Asked as: First Nations answer missing; First Nations agreement required; project request yes no error

Before checking a project request, Turbo requires the First Nations agreement answer to be exactly Yes or No. Otherwise it refuses with 'First Nations agreement answer is required'.

## Request rate sheet draft or active sheet is checked

Lines 395-407

Asked as: rate sheet not active; proposed rates check fails; request rate sheet error

If a project request has proposed rates, Turbo validates that draft. Otherwise it requires an active rate sheet (validateRateSheetIsActive).

## Project request branch must exist and be active

Lines 408-414

Asked as: branch missing on request; inactive branch selected; request branch error

Turbo refuses a project request with no valid branch ('branch is required'). It also refuses an inactive branch ('branch must be active').

## Restore refused while value changes are pending

Lines 430-437

Asked as: cannot restore project value; pending changes block restore; restore refused unapproved attachments

Turbo refuses to restore a project value while the job has any attachment that is not Approved. The message is 'delete pending project value changes before restoration'.
