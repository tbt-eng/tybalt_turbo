<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/project_workflow_guards.go","sha256":"ebf43f2af92a96d9361de5e750a44366790e64745f8fd875d5b360cc43657da5","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/project_workflow_guards.go

## Workflow transitions still get normal job validation

Lines 82-86

Asked as: workflow skips validation; strict checks on workflow; job validation during approval; workflow fails validation

A controlled project workflow transition (ProcessJobWorkflow) marks the job with its operation and then runs the strict job validation (ProcessJobCoreStrict). The operation mark is cleared afterwards.

## Project workflow actions require a signed-in user

Lines 90-93

Asked as: not signed in error; authentication is required; must log in first; no user signed in

When a project workflow check needs a user and no user is signed in, Turbo refuses with the message "authentication is required".

## Workflow actions need an active user with a required claim

Lines 94-108

Asked as: inactive user blocked; missing claim for workflow; cannot act on project; user not active

Turbo refuses a project workflow action unless the user is active and holds at least one of the listed claims. The message is "an active user with the <claims> claim is required", and the claims are joined with " or ". An inactive user is refused even when they hold the claim.

## Project writes need an active user with the job claim

Lines 111-120

Asked as: cannot save project; missing job claim; procurement permission needed; can't write project

Writing a project (not a proposal) requires an active user who holds the job claim (procurement). Otherwise Turbo refuses with "an active user with the job claim is required". Approving or restoring a project value skips this claim check. Proposals skip all these project checks.

## Value approval and restoration skip the job claim check

Lines 115-120

Asked as: approve value without claim; restore value no claim; approver needs no procurement; approving doesn't need job claim

When the operation is approve_value or restore_value, the job claim check is skipped. The user who approves or restores a project value does not need the job claim (ProjectWorkflowMutationApproveValue, ProjectWorkflowMutationRestoreValue).

## Project value must be whole dollars

Lines 121-123

Asked as: project value has cents; value with decimals; dollars and cents not allowed; project amount whole dollars

Turbo refuses a project value (project_value) that has cents. The message is "project value must use whole dollars". Proposals are exempt in the validation step.

## New projects need an approved project request when approval is required

Lines 124-132

Asked as: create project without request; new project blocked; project request approval needed; why can't I create project

When the job request approval setting is on, a new project can be created only through the create-project step of an approved project request. Any other attempt is refused with "create projects through an approved project request".

## Project client changes only through merge or change-client

Lines 133-136

Asked as: change project client; edit client on project; move project to another client; client is locked

An existing project's client cannot be edited directly. The validation step says "project client can change only through a client merge". The write guard says "project client is immutable" and also lets the change-client action through, which moves a project and its subjobs to another client with a recorded reason.

## Project value changes require approval

Lines 137-139

Asked as: change project amount; edit project value directly; value change blocked; update project total

Changing the value of an existing project is refused unless the change comes from approving or restoring a value. The message is "project value changes require approval".

## Cannot delete a user tied to project authorization records

Lines 146-160

Asked as: delete user error; remove user with requests; user has project history; can't delete staff member

Turbo refuses to delete a user who appears as creator, approver or rejector on project requests, as uploader, approver or rejector on job attachments, or as request creator or creator of a job. The message is "user identity is retained in project authorization records".

## Cannot delete a user who manages a job

Lines 191-198

Asked as: user is a manager; delete job manager; remove staff managing job; reassign job manager first

Turbo refuses to delete a user who is the manager of any job, project or proposal. The message is "user is the manager of one or more jobs". Deleting through the records API gives a generic delete error instead.

## Request checklist versions must be sequential

Lines 202-213

Asked as: checklist version number; next checklist version; request checklist numbering; checklist version wrong

A new request checklist must have a version number one higher than the latest version for its scope. Otherwise Turbo says "the next <scope> request checklist version is N". The scope check runs only when the scope field exists.

## Request checklist questions must be valid

Lines 214-217

Asked as: checklist question error; bad checklist questions; questions won't save checklist; checklist validation message

A new request checklist is refused if its questions fail validation, and Turbo shows the validation error message.

## Saved request checklists cannot be edited or deleted

Lines 219-223

Asked as: edit saved checklist; delete checklist; change published checklist; make new checklist version

Updating or deleting a saved request checklist is refused with "a saved request checklist cannot change; add the next version".

## Project requests, allocations and attachments only change through workflow actions

Lines 225-230

Asked as: cannot edit project request; edit allocation directly; attachment update blocked; use workflow to delete

Creating or updating a project request, request allocation or job attachment outside the project workflow is refused with "use the project workflow actions". Deleting outside the workflow is refused with "use the project workflow delete action".

## Attachment must belong to exactly one request or job

Lines 231-235

Asked as: attachment linked to both; attachment no request or job; attach to project and job; attachment link error

A job attachment must link to either a project request or a job, not both and not neither. Otherwise Turbo says "an attachment must belong to one request or job".

## Approved project evidence cannot be edited

Lines 236-238

Asked as: edit approved attachment; approved document locked; change approved evidence; modify approved file

An approved attachment already linked to a job cannot be edited. Turbo says "approved project evidence cannot be edited".

## Attachment dollar amounts must be whole

Lines 239-244

Asked as: attachment amount has cents; approval max decimals; previous value cents; attachment dollars must be whole

The project value, previous project value and approval max value on an attachment must be whole dollars. Turbo says "<field> must use whole dollars".

## Initial project approval cannot be deleted

Lines 250-257

Asked as: delete first approval; remove initial approval; delete original approval attachment; initial approval locked

Deleting the initial approved attachment of a job is refused with "initial project approval cannot be deleted".

## Approved value history must be deleted newest first

Lines 258-267

Asked as: delete old approval; pending changes block deletion; remove value change order; delete attachment with later changes

An approved attachment on a job cannot be deleted while the job has pending attachments or approved ones with a later approval sequence. Turbo says "delete pending or later project changes first".

## Job number cannot switch between project and proposal

Lines 274-278

Asked as: change job to proposal; convert project to proposal; job number type change; switch project and proposal number

Editing a job so its number changes from project type to proposal type, or the reverse, is refused with "a job number cannot switch between project and proposal".

## Project writes must go through validated job action

Lines 279-287

Asked as: project save blocked; write project directly; validated job action needed; project not saved

A project write is refused unless it comes from a workflow action or a validated job save. Turbo says "project writes require the validated job action". Proposals skip these project guards.

## Saving a project requires the validated job action

Lines 282-287

Asked as: save project after validation; validation voided by edit; changed field after validate; validation already used

A project is saved only inside a workflow action or right after Turbo has validated that exact record. A project write that is neither is refused with "project writes require the validated job action". Changing any field after validation voids the validation. Each validation counts once.

## Project value must be whole dollars at save

Lines 289-291

Asked as: save project with cents; project value save fails; decimal project value save; round value before saving

At save, a project whose value has cents is refused with "project value must use whole dollars".

## Restoring a project value needs historical validation

Lines 292-296

Asked as: restore value check; restore fails historical; restore project value validation; historical restore error

A value restore runs the historical restoration check before saving.

## Project creation requires approved request when approval is on

Lines 298-306

Asked as: project save no approved request; create job not allowed; new project save error; approval required for new project

When approval of project requests is required, a new project is refused with "create projects through an approved project request" unless the operation is create_job. This check runs when the project is saved.

## Project client is immutable except merge or change-client

Lines 307-311

Asked as: client cannot change on save; project client locked; save changed client error; client change rejected on save

Changing the client of an existing project at save is refused with "project client is immutable". Only a client merge or the change-client action may do it.

## Project value changes need approval at save

Lines 313-315

Asked as: save changed project value; value change save blocked; amount changed on save; approve value first

Changing the value of an existing project at save is refused with "project value changes require approval". Only approve_value and restore_value operations may change it.

## Rate sheet changes go through rate sheet history

Lines 316-320

Asked as: change rate sheet; rate sheet locked; edit project rate sheet; rate sheet history required

Once a project has a rate sheet, changing it directly is refused with "change a project's rate sheet from its rate sheet history". Setting the first rate sheet is a normal edit.

## Project creation history is protected

Lines 322-326

Asked as: creation date changed; change created by; change request creator; history fields locked

The request creator, request created, request submitted, created by and creation completed fields of a project change only during the create-project step. Otherwise Turbo says "project creation history is protected".

## Job with attachments cannot be deleted

Lines 331-343

Asked as: delete job with attachments; cannot delete job files; job has attachments; remove job attachments first

Deleting a job that still has project attachments is refused with "delete permitted project attachments before deleting this job".

## Approval limit must be positive whole dollars

Lines 345-349

Asked as: approval limit error; max value invalid; authorization limit number; limit must be positive

A project authorization limit (max_value) must be a positive whole number no greater than 9007199254740991. Otherwise Turbo says "approval limit must be positive whole dollars no greater than 9007199254740991".

## Approval limits attach only to project_authorization claim

Lines 350-361

Asked as: claim properties error; approval limit on claim; properties need assignment; limit on wrong claim

Properties on a user claim assignment are refused unless the claim is project_authorization. Turbo says "properties require a project_authorization assignment".

## Project request data is checked as a future job

Lines 368-384

Asked as: submit project request error; First Nations agreement required; inactive rate sheet submit; request submission validation

Submitting a project request validates it as a future project with the PA authorizing document, without assigning a number. The First Nations agreement answer must be Yes or No, else Turbo says "First Nations agreement answer is required". The rate sheet must be active.

## Project request branch must exist and be active

Lines 386-394

Asked as: branch missing on request; inactive branch request; pick branch for request; branch error on request

A project request needs an existing branch, else "branch is required". An inactive branch is refused with "branch must be active". Then the full job validation runs.

## Only the latest approved change can be restored

Lines 399-407

Asked as: restore older value; cannot restore initial approval; restore previous value error; stored previous value mismatch

Restoring a project value works only for the latest approved change, never the initial approval. The job's value must match what that change replaced. Otherwise Turbo says "restore only the latest approved change's stored previous value".

## Pending value changes block restoration

Lines 408-414

Asked as: pending value changes restore; restore blocked by pending; unapproved changes stop restore; delete pending before restore

A value restore is refused while the job has any non-approved attachment. Turbo says "delete pending project value changes before restoration".
