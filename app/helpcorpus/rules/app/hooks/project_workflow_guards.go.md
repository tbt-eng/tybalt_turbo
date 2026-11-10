<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/project_workflow_guards.go","sha256":"66f3e628dd525d9b5cf6ef428f063ab5ba6a3bda03e018fec2fd75485c632ae4","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/project_workflow_guards.go

## Job workflow actions run normal validation in a controlled transition

Lines 78-83

A job workflow action (ProcessJobWorkflow) marks the job with the named mutation, then runs the strict job validation (ProcessJobCoreStrict). The mark is cleared when the action ends. Later guards use the mark to allow the controlled changes.

## Workflow actions need an active user with the right claim

Lines 87-107

When Turbo checks an actor's claims, it refuses a missing actor with 'authentication is required'. It also refuses an inactive user, or one without any of the listed claims. The message is 'an active user with the <claims joined by " or "> claim is required'. (activeClaim)

## Project writes need an active user with the job claim

Lines 108-117

When a project job is written, Turbo requires an active user who has the job claim (job). Otherwise it refuses with 'an active user with the job claim is required', or 'authentication is required' when no user is signed in. Approving or restoring a project value skips this check. Proposals skip it too.

## Project value must be whole dollars

Lines 118-120

Turbo refuses a project whose value (project_value) has cents. The message is 'project value must use whole dollars'. Proposals are not checked.

## New projects need an approved project request when required

Lines 121-129

When the setting that requires job request approval is on, Turbo refuses to create a project directly. The message is 'create projects through an approved project request'. Only the create-job workflow step may create it.

## Project client changes only through merge or change-client

Lines 130-133

On an existing project, a changed client is refused unless the action is a client merge. The validateProjectJobWrite check says 'project client can change only through a client merge'. The jobs guard also allows the change-client action and says 'project client is immutable' otherwise.

## Project value changes need approval

Lines 134-136

On an existing project, changing the project value is refused unless the action is approve-value or restore-value. The message is 'project value changes require approval'.

## Cannot delete a user kept in project authorization records

Lines 143-157

Turbo refuses to delete a user who is the creator, approver or rejector on a project request, the uploader, approver or rejector on an attachment, or the request creator or creator of a job. The message is 'user identity is retained in project authorization records'.

## Cannot delete a user who manages a job

Lines 188-195

Turbo refuses to delete a user who is the manager of any job, project or proposal. The message is 'user is the manager of one or more jobs'. A delete through the records API shows PocketBase's generic delete error instead.

## New request checklist must be the next version

Lines 199-215

A new request checklist must have a version number one higher than the latest version for its scope. Otherwise Turbo refuses with 'the next <scope> request checklist version is N'. The check is skipped when the scope field does not exist yet. Turbo also refuses a checklist whose questions are invalid, giving that error's message.

## Saved request checklists cannot be edited or deleted

Lines 216-220

Turbo refuses to update or delete a saved request checklist. The message is 'a saved request checklist cannot change; add the next version'.

## Project requests, allocations and attachments need workflow actions

Lines 222-250

Creating or updating a project request, request allocation or job attachment outside a project workflow action is refused. The message is 'use the project workflow actions'. Deleting one outside the workflow is refused with 'use the project workflow delete action'.

## Attachment must belong to exactly one request or job

Lines 228-232

A job attachment must link to either a request or a job, not both and not neither. Otherwise Turbo refuses with 'an attachment must belong to one request or job'.

## Approved project evidence cannot be edited

Lines 233-235

Turbo refuses to edit an attachment that is linked to a job and already has status Approved. The message is 'approved project evidence cannot be edited'.

## Attachment values must be whole dollars

Lines 236-241

An attachment's project_value, previous_project_value and approval_max_value must be whole dollars. Otherwise Turbo refuses with '<field> must use whole dollars'.

## Initial project approval cannot be deleted

Lines 251-254

Turbo refuses to delete an approved job attachment that is the initial approval (initial_approval). The message is 'initial project approval cannot be deleted'.

## Delete later or pending project changes first

Lines 255-263

Turbo refuses to delete an approved attachment while the job has pending (non-approved) attachments or approved ones with a larger approval sequence. The message is 'delete pending or later project changes first'.

## Job number cannot switch between project and proposal

Lines 271-275

Turbo refuses to change a job's number so that its type changes between project and proposal. The message is 'a job number cannot switch between project and proposal'.

## Proposals skip the project guards

Lines 276-278

For a proposal, the jobs guard passes the write on after the number-type check. The project rules below do not apply to proposals.

## Project writes require the validated job action

Lines 279-284

A project write is refused unless it comes through a workflow action or carries a validation mark matching the record's current contents. The message is 'project writes require the validated job action'. The mark is cleared after each write.

## Guard enforces whole-dollar project value

Lines 286-288

The jobs guard refuses a project value with cents, with 'project value must use whole dollars'.

## Value restore must use latest approval's previous value

Lines 289-293

A restore-value action is refused unless the latest approved change (highest approval sequence) is not the initial approval. Its stored previous value must equal the restored value, and its value must equal the job's current value. Otherwise the message is 'restore only the latest approved change's stored previous value'. Pending changes are also refused with 'delete pending project value changes before restoration'.

## Jobs guard blocks direct project creation when approval required

Lines 295-303

With job request approval required, the jobs guard refuses creating a project unless the action is create-job. The message is 'create projects through an approved project request'.

## Jobs guard blocks client and value changes

Lines 304-313

On an existing project, a client change is refused unless the action is merge or change-client ('project client is immutable'). A value change is refused unless the action is approve-value or restore-value ('project value changes require approval').

## Project creation history is protected

Lines 314-318

Turbo refuses to change a project's request creator, request created, request submitted, created by or creation completed fields except in the create-job action. The message is 'project creation history is protected'.

## Cannot delete job with project attachments

Lines 323-335

Turbo refuses to delete a job that still has attachments. The message is 'delete permitted project attachments before deleting this job'.

## Approval limit must be positive whole dollars

Lines 336-341

A project authorization limit (max_value) must be above 0, a whole number and at most 9007199254740991. Otherwise Turbo refuses with 'approval limit must be positive whole dollars no greater than 9007199254740991'.

## Approval limit needs a project_authorization assignment

Lines 342-356

Project authorization properties can attach only to a user claim whose claim is project_authorization. Otherwise Turbo refuses with 'properties require a project_authorization assignment'.

## Project request is checked as a future job

Lines 359-387

Before a project request is accepted, Turbo copies its fields into a new job. The job has the request's project status and the authorizing document set to "PA". It then runs the normal job validation with the rate sheet and branch checks. The job number is not assigned.

## Project request validates as a future job

Lines 360-387

Turbo builds a trial job from the request's fields. It sets the status from project_status and the authorizing document to 'PA', without assigning a number. It checks that the rate sheet is active, then runs the full job validation.

## Project request needs a First Nations agreement answer

Lines 371-374

Validating a project request requires fn_agreement_answer to be 'Yes' or 'No'. Otherwise Turbo refuses with 'First Nations agreement answer is required'.

## Project request branch must exist and be active

Lines 378-384

A project request must have a branch that exists, or Turbo refuses with 'branch is required'. If branches have an active flag and the branch is inactive, it refuses with 'branch must be active'.

## Only the latest approved value change can be restored

Lines 389-399

Restoring a project's value (restore_value) is refused unless the job's newest approved attachment (highest approval_sequence) is not the initial approval, and its stored previous project value equals the job's new project value. Its project value must also equal the job's current value. Otherwise Turbo says "restore only the latest approved change's stored previous value".

## Pending value changes block restoring a project value

Lines 400-406

Restoring a project's previous value is refused while the job has any attachment that is not Approved. Turbo says "delete pending project value changes before restoration".
