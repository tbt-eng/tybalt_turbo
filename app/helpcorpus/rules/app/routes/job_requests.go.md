<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_requests.go","sha256":"38a4f2777aa7a6cf96c4ceda623caa65f512b3e83c273bfae24148fc209822f1","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_requests.go

## Abandoned requests are removed after the abandonment period

Lines 27-27

Asked as: old draft deleted; request gone after weeks; stale request removed automatically; unused request disappeared

A project request counts as abandoned when its last edit (last_edited) is the abandonment limit or older. The limit is JOB_REQUEST_ABANDONMENT_DAYS days. A request edited exactly that long ago counts as abandoned.

## Request fields editable only from a fixed list

Lines 30-30

Asked as: can't change this field; field won't save; field locked on request; why won't it let me edit

A save with any field outside the allowed business fields is refused with 'Field is not editable: <field>'.

## Only active staff can act on project requests

Lines 102-114

Asked as: can't act on request; not active staff; inactive account cannot approve; request button not working

Turbo refuses any request action from a user who is not an active staff member. The message is "An active staff member is required." A user who is not signed in sees "Sign in to continue."

## Claim check also requires an active staff member

Lines 102-127

Asked as: claim check fails on request; missing permission to act; sign in message on request

When Turbo checks a claim for a project request action, it first requires an active staff member and refuses with 'An active staff member is required.' A signed-in person without the claim gets 'The <claim> claim is required.' A request without sign-in gets 'Sign in to continue.'

## Missing claim refusal names the claim

Lines 115-127

Asked as: claim required error; don't have permission; missing permission message

When a person lacks a required claim, Turbo refuses with 'The <claim> claim is required.' Turbo checks that the person is active before it checks the claim.

## Project requests can be turned off

Lines 128-137

Asked as: project requests turned off; can't save request at all; request feature disabled

When project requests are not required (IsJobRequestApprovalRequired is false), Turbo refuses to save, submit, recall, reject, approve or create a request. The message is "Project requests are disabled."

## Abandoned request definition

Lines 138-140

Asked as: request is too old; how old before abandoned; last edited too long ago

A request is abandoned when its last edit is at least JOB_REQUEST_ABANDONMENT_DAYS days old. A request edited exactly that many days ago counts as abandoned.

## Who can view a project request

Lines 141-172

Asked as: who can see request; can't open request; request hidden from me; view someone else's request

The creator of a request can view it. Anyone with the job claim can view any request. A project authorization approver can view requests that name them as approver. An admin can view only abandoned requests. Anyone else sees "You cannot view this request."

## Missing request gives not found message

Lines 159-172

Asked as: request not found; can't find my request; request missing; blank request page

Opening a request that does not exist gives 'Project request not found.' A request the person may not view gives 'You cannot view this request.'

## Only the creator edits, and only in some states

Lines 173-188

Asked as: can't edit request; edit button locked; only creator edits; why can't I change request

Only the creator can edit a request, and only when its status is Draft, Recalled or Rejected. Others see "Only the creator can edit this request." In any other status the creator sees "Recall this request before editing it."

## Stale request version is refused

Lines 200-207

Asked as: someone changed request; request out of date; conflict saving request; reload request changed

Turbo compares the version the user sends with the request's saved version (version). If they differ, the action is refused with "This request changed. Reload it before continuing."

## Editing a request clears approval and issues no token

Lines 208-245

Asked as: edit resets approval; approval gone after edit; pdf back to pending; editing removes approval

Each business edit of a request raises its version by 1 and sets last edited to now. It clears the submission token and any rejection. A Rejected request returns to Draft. It also resets the attached authorization PDF to Pending, clears its approval, copies the new project value to it, and raises its version.

## Editing a request bumps version and clears approvals

Lines 222-245

Asked as: editing clears approval; rejected goes back to draft; approval reset after change

A real edit raises the version by 1 and sets last edited to now. It clears the submission token and any rejection, and turns a Rejected request back into a Draft. It also resets the attachment's approval to Pending and copies the project value to it.

## Project value must be a whole non-negative amount

Lines 246-261

Asked as: project value decimal; cents in project value; negative project amount; value not whole dollars

A saved project value must be a whole number of CAD dollars. It must be 0 or more, finite, and at most 9007199254740991. Otherwise the message is "Project value must be a whole CAD dollar amount."

## Only listed fields can be edited

Lines 250-258

Asked as: can't save that field; field not allowed; field not editable

A request save accepts only the business fields (requestBusinessFields). Any other field is refused with "Field is not editable: " followed by the field name.

## Request notes length limit

Lines 262-273

Asked as: notes too short; note length error; notes character limit; client note text size

Notes are trimmed. They must be text, and must be blank or 10 to 1000 characters. Otherwise the message is "Notes must be 10 to 1000 characters, or blank." The limit applies because the notes become a client note when the job is created.

## Fee basis must be chosen

Lines 274-279

Asked as: fee basis not chosen; time and materials choice; pick hourly or fixed fee

The time and materials value must be true or false, otherwise the message is "Select a fee basis." Setting it marks the fee basis as chosen (time_and_materials_set). Submission needs this mark.

## Draft checks on project status, approver, parent and proposal

Lines 287-326

Asked as: project status error on draft; approver limit too low draft; parent or proposal wrong

A saved request's project status must be Active or Closed. If an approver is set, their project limit must cover the project value (ProjectAuthorizationLimit). A parent must be a top-level project (number like 25-123) with the request's client. A proposal must be a P-number job for the same client, with status In Progress, Submitted or Awarded.

## Project status must be Active or Closed

Lines 288-290

Asked as: project status wrong; status must be active; can't save project status

A request saves only when its project status is Active or Closed. Otherwise Turbo says 'Select Active or Closed for the project status.'

## Chosen approver must have enough project limit

Lines 294-298

Asked as: approver limit too small; approver can't cover value; approver limit check

When a request has an approver, Turbo checks the approver's project authorization limit against the project value (ProjectAuthorizationLimit). Approve and submit run this check. Creation skips it.

## Parent project rules for a subjob request

Lines 299-310

Asked as: parent project error; subjob parent invalid; wrong parent client; top level project needed

A parent must exist ('Select a valid parent project.'), must have a top-level number like 25-123 ('Only a top-level project can be a parent.'), and must have the request's client ('The subjob must use the parent project's client.').

## Linked proposal rules

Lines 311-350

Asked as: proposal not valid; linked proposal status; proposal must be awarded; proposal client mismatch

A linked proposal must exist, have a number starting with P and belong to the request's client ('Select a proposal for this client.'). Its status must be In Progress, Submitted or Awarded. Submission needs it Awarded: 'Set the linked proposal to Awarded before continuing.'

## Submission requires a complete request

Lines 327-350

Asked as: can't submit request; submit fails required fields; approve incomplete request; missing info before submit

Submit and approve need every required field: description, location, client, contact, manager, branch, rate sheet, award date, invoicing information, approver and fn_agreement_answer. Otherwise the message is "Required field: " and the field name. They also need a chosen fee basis and a positive project value. A linked proposal must be Awarded, or the message is "Set the linked proposal to Awarded before continuing."

## Required fields before submitting a request

Lines 331-338

Asked as: required field missing; what fields to fill in; fields needed before submitting

A request needs description, location, client, contact, manager, branch, rate sheet, award date, invoicing information, approver and fn agreement answer, or Turbo says 'Required field: <field>'. A fee basis must be chosen: 'Select a fee basis.'

## Allocation hours must be whole and in active divisions

Lines 361-407

Asked as: allocation hours fraction; division not active; invalid division on request; allocations missing

Allocation hours must be whole numbers. A division must exist ('Select a valid division.') and be active ('Select an active division.'). A complete request needs allocations.

## Authorization PDF is required

Lines 371-381

Asked as: no authorization pdf; upload pdf error; need attachment to submit

A request must have an uploaded authorization PDF with a document hash before it is submitted, approved or created. Otherwise the message is "Upload an authorization PDF."

## Allocation rules for requests

Lines 383-407

Asked as: allocation error; division hours invalid; draft without allocations

Division hours must be whole numbers, finite, and each division must exist and be active. Messages include "Allocation hours must be whole numbers.", "Select a valid division." and "Select an active division." A completed request needs allocations, but a draft save may leave them out.

## Creating a request needs the time claim and starts as a draft

Lines 421-435

Asked as: who can create request; new request status; need time claim to create

Only a user with the time claim can create a request. The new request has status Draft, project status Active and version 1.

## Choose a rate sheet or propose rates, not both

Lines 459-485

Asked as: rate sheet and proposed rates; both rate sheet and custom rates; rates conflict

When a save includes proposed rates and also picks a different rate sheet, Turbo refuses it with "Choose a rate sheet or propose new rates, not both." When a user picks another sheet, the request's own draft rates are dropped. Proposed rates make the request's rate sheet the draft.

## Unchanged save does nothing

Lines 500-509

Asked as: save did nothing; no change saved; request not updated when unchanged

When an existing request is saved with no change to its fields, allocations or proposed rates, Turbo does not save it. Its version and last edited date stay the same.

## Who can reject a request

Lines 570-588

Asked as: who can reject request; job holder reject; approver reject request; reject not allowed

The selected approver (project_authorization claim) can reject while the request is Awaiting approval or Awaiting creation. A job holder (job claim) can reject only once it is Awaiting creation. Others see "Only the selected approver, or a job holder once the request is approved, can reject this request."

## Who can reject a request and the reason needed

Lines 570-715

Asked as: rejection reason required; enter reason to reject; reject needs reason

The selected approver (project_authorization claim) can reject while the request is Awaiting approval or Awaiting creation. A job holder can reject only when Awaiting creation. Others get 'Only the selected approver, or a job holder once the request is approved, can reject this request.' A reason is required: 'Enter a rejection reason.' Rejection records who and when and clears the token.

## Checklist must be fully confirmed

Lines 590-617

Asked as: checklist not confirmed; checklist changed reload; unanswered checklist questions; no checklist saved

Approve and create need every question of the current checklist (approval or creation) answered true. The checklist version sent must match the current one, or the message is "The checklist changed. Reload the request and review it again." No saved checklist gives "No <scope> checklist is saved." The wrong number of answers is refused. An unconfirmed question gives "Confirm: " and its key.

## Submit a request

Lines 644-671

Asked as: submit my request; send request for approval; edit rejected before submit

Submit needs the time claim and an editable request, and the creator is the only one who can edit. A Rejected request must be edited first, with the message "Edit the rejected request before submission." A successful submit sets status to Awaiting approval, records the submission time, issues a new submission token and clears earlier approvals.

## Recall a request

Lines 672-715

Asked as: recall my request; take back submission; withdraw request

Only the creator can recall, with the message "Only the creator can recall this request." The request must be Awaiting approval or Awaiting creation, or the message is "This action requires a submitted or approved request." Recall sets status to Recalled, clears the submission token and approvals.

## Rejection needs a reason

Lines 690-715

Asked as: reject reason blank; reject without reason; rejection note required

Reject needs a non-blank reason, or the message is "Enter a rejection reason." The request and its attachment become Rejected with the rejector, time and reason recorded. The submission token and approvals are cleared.

## Approve a request

Lines 716-764

Asked as: approve request button; approver sign off; no submission token

Only the selected approver can approve, and only a submitted request (status Awaiting approval, not rejected). The request must have a submission token, or the message is "This request has no submission token. Recall and submit it again." The approver's project limit must cover the value. Approval moves the request to Awaiting creation and records the limit as the approved maximum value.

## Approval must match what the approver reviewed

Lines 740-746

Asked as: reviewed data changed; pdf changed after review; review token mismatch

Approve checks that the review token and document hash from the approver's page match the current data and PDF. If they differ, the message is "The reviewed data changed. Reload and review it again."

## Create a job from an approved request

Lines 765-821

Asked as: create job from request; job creation stale approval; approval incomplete

Only a user with the job claim can create the job, and only when the request is Awaiting creation. The approval must still match: same version, project value and review token, approver and approval limit at least the value. Otherwise the message is "The approval is stale. Recall the request for a new approval." or "The saved approval is incomplete. Recall the request for a new approval." The approver's authority limit is not rechecked at creation.

## What job creation does with the request

Lines 850-922

Asked as: what happens to request when job created; request deleted after job; job from request fields

Creation builds a job from the request's business fields, without the approver. The job status is the project status, and the authorizing document is PA. The job's fn_agreement is true when the answer is Yes. Allocations and the authorization PDF move to the job, with approval sequence 1. Then Turbo deletes the request.

## Request notes become the job's first client note

Lines 924-945

Asked as: request notes client note; notes copied to job; notes after job created

On job creation, non-blank request notes become a client note on the job and client. The note keeps the request creator as author. Blank notes make no note.

## Who can delete a request

Lines 946-993

Asked as: delete my request; can't delete request; remove abandoned request; admin delete request

A creator can delete a request that they can still edit (Draft, Recalled or Rejected). Otherwise only an admin can delete it, and only when it is abandoned. A non-abandoned request gives "Only abandoned requests can be deleted." Deleting also removes its attachment and allocations.
