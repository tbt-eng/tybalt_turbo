<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_requests.go","sha256":"933e7828385a956237b22c6c8928072981fe30f478e4f4e15da5d47a0a17612c","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_requests.go

## Abandonment window for project requests

Lines 27-27

Asked as: request went stale; old draft expired; when does request expire; abandoned after how many days; not edited for days

A project request counts as abandoned when its last edit (last_edited) is JOB_REQUEST_ABANDONMENT_DAYS days old or older. A request edited exactly that many days ago is abandoned, because the check is 'not after'.

## Only listed fields can be edited

Lines 30-30

Asked as: field not editable; cannot change this field; save rejected field; which fields can I edit; locked field on request

A save may set only the business fields (project status, approver, client, job owner, contact, manager, alternate manager, branch, rate sheet, parent, proposal, invoicing information, description, notes, location, client PO, client reference number, award date, completion date, project value, fee basis and First Nations agreement answer). Any other field gets "Field is not editable: <field>".

## Project request routes need sign-in

Lines 57-59

Asked as: not logged in; session expired request; sign in error; need to log in; unauthenticated request page

All project request routes and attachment routes require a signed-in user. Turbo refuses an unauthenticated user with 'Sign in to continue.'

## Authorization attachment approve and reject routes

Lines 71-86

Asked as: approve attachment route; reject attachment; authorization file routes; which routes for approvals; attachment download delete

Authorization attachments (job_attachments) have their own routes to list, save, download, delete, approve and reject. Project requests have routes to submit, recall, reject, approve and create.

## Validation errors show a generic message

Lines 95-98

Asked as: validation failed; fields have errors; what is wrong with my form; form error message; request won't validate

When fields fail validation, Turbo says 'Check the request fields.' and lists the failing fields.

## Only active staff can use project requests

Lines 106-118

Asked as: not active staff; why can't I access requests; inactive employee blocked; access denied project requests; staff account disabled

Turbo refuses a person who is not an active staff member with "An active staff member is required." A missing sign-in gets "Sign in to continue."

## Missing claim refusal message

Lines 119-131

Asked as: missing permission claim; don't have the claim; permission needed for action; lack access right; claim required error

When an action needs a claim and the person lacks it, Turbo says 'The <claim> claim is required.' The person must also be active staff.

## Project requests can be switched off

Lines 132-141

Asked as: project requests turned off; feature disabled; cannot submit request at all; requests switched off; approval setting disabled

When project requests are not required (IsJobRequestApprovalRequired is false), Turbo refuses to save, submit, recall, approve, reject, create or edit a request. It shows "Project requests are disabled." Deleting a request does not check this setting.

## Abandoned request definition

Lines 142-144

Asked as: what makes request abandoned; when is request considered abandoned; abandoned request meaning; how old before abandoned; stale request definition

A request is abandoned when its last edit (last_edited) is JOB_REQUEST_ABANDONMENT_DAYS days old or older. A request edited exactly that many days ago counts as abandoned.

## Who can view a project request

Lines 145-176

Asked as: cannot see request; who can view request; request hidden from me; view permission denied; why can't I open this

The creator of a request can view it. Anyone with the job claim (procurement) can view any request. A person with the project_authorization claim can view a request only when they are its selected approver. An admin can view a request only when it is abandoned. Anyone else sees "You cannot view this request."

## Request not found message

Lines 163-167

Asked as: request does not exist; missing project request; bad request id; can't find request; no such request

When a project request id does not exist, Turbo says 'Project request not found.'

## Only creator edits, in Draft, Recalled or Rejected

Lines 177-192

Asked as: cannot edit submitted request; who can edit request; edit locked after submit; not the creator editing; can't change approved request

Only the creator can edit a request, and only while its status is Draft, Recalled or Rejected. Others get "Only the creator can edit this request." A creator who edits a Submitted or approved request sees "Recall this request before editing it."

## Edit refusal messages

Lines 184-191

Asked as: edit blocked message; someone else's request edit; recall before editing; why can't I edit my request; edit refused in this status

A person other than the creator who edits a request sees 'Only the creator can edit this request.' The creator who edits a request in any other status sees 'Recall this request before editing it.'

## Stale version is refused

Lines 204-211

Asked as: request changed underneath me; someone else edited it; out of date request; reload request version; version mismatch save

Every save or action carries the request version. When it differs from the saved version, Turbo refuses with "This request changed. Reload it before continuing."

## Clearing an attachment approval

Lines 212-220

Asked as: clear attachment approval; remove approval from file; reset authorization approval; undo approval of attachment; attachment back to pending

Clearing an authorization attachment's approval empties the approver, approval date, review token and rejection details. It resets the approval limit, versions and sequence to 0 and sets the status to Pending.

## Editing bumps version and clears approval

Lines 226-249

Asked as: editing resets approval; rejected goes back to draft; edit clears submission; why status changed to draft; approval lost after edit

Each edit raises the request version and updates its last edit time. It clears the submission token and any rejection. A Rejected request returns to Draft. It resets the attachment's approval to Pending and copies the project value.

## Project value must be whole dollars

Lines 250-265

Asked as: project value decimal; cents in project value; value not whole dollars; project value zero submit; dollar amount invalid

Project value must be a whole CAD dollar amount of zero or more. Otherwise Turbo says "Project value must be a whole CAD dollar amount." A request may be saved with a value of 0, but submission needs a positive value.

## Unknown request fields are refused

Lines 254-262

Asked as: unknown field saved; extra field in save; field not allowed on request; save with bad field name; invalid field name

Saving a request with a field outside the editable list fails with 'Field is not editable: <field>'.

## Request notes limits

Lines 266-277

Asked as: notes too short; notes length limit; note character count; notes text rules; notes over 1000

Notes are trimmed. They must be text, and 10 to 1000 characters, or blank. Otherwise Turbo says "Notes must be 10 to 1000 characters, or blank."

## Fee basis must be chosen

Lines 278-283

Asked as: fee basis not chosen; time and materials missing; choose billing basis; fee basis required to submit; t and m selection

Time and materials must be a true or false choice, otherwise "Select a fee basis." Setting it records that a fee basis was chosen. Submission refuses a request without that choice.

## Draft validation of status, parent and proposal

Lines 291-314

Asked as: parent project invalid; subjob setup error; wrong project status choice; approver limit too low for value; draft status parent check

Project status must be Active or Closed ("Select Active or Closed for the project status."). A chosen approver's project limit is checked against the project value (ProjectAuthorizationLimit). A parent must exist, must be a top-level project with a number like 25-123, and must have the same client as the request. Otherwise Turbo says "Select a valid parent project.", "Only a top-level project can be a parent." or "The subjob must use the parent project's client."

## Project status must be Active or Closed

Lines 292-294

Asked as: project status options; status must be active or closed; new request status; pick status for project; invalid project status

A request must have a project status of Active or Closed. Otherwise Turbo says 'Select Active or Closed for the project status.' A new request starts as Active.

## Approver limit checked on save

Lines 298-302

Asked as: approver limit too low; approver cannot approve this value; authorization limit exceeded; approver's limit on save; project value above limit

When a request has an approver, saving checks that approver's project authorization limit against the project value (ProjectAuthorizationLimit). Job creation skips this check.

## Parent project rules for a subjob

Lines 303-314

Asked as: subjob parent error; parent project not found; child job client mismatch; top level project required; parent number format

A parent must be an existing job, or Turbo says 'Select a valid parent project.' Its number must look like 25-123, or Turbo says 'Only a top-level project can be a parent.' The client must match the parent's, or Turbo says 'The subjob must use the parent project's client.'

## Proposal must be for the client and open

Lines 315-328

Asked as: proposal wrong client; proposal not open; proposal status not allowed; pick a proposal; proposal must start with P

A proposal must exist, start with P and belong to the request's client ("Select a proposal for this client."). Its status must be In Progress, Submitted or Awarded, else "Proposal <number> is <status>. Select an In Progress, Submitted, or Awarded proposal."

## Job creation needs an Awarded proposal

Lines 332-357

Asked as: cannot create job yet; proposal not awarded; award proposal first; job creation blocked by proposal; proposal still in progress

A request with a proposal can be submitted and approved while the proposal is In Progress or Submitted. Creating the job needs the proposal Awarded, else "Proposal <number> is <status>. Award it before you create the job." Turbo does not award the proposal itself. The refusal carries the code proposal_not_awarded and the proposal id.

## Required fields for submission

Lines 362-372

Asked as: cannot submit request; missing fields before submit; what to fill before approval; required field error; submit needs more info

Before submit or approval, these fields must be filled: description, location, client, contact, manager, branch, rate sheet, award date, approver and First Nations agreement answer. Otherwise Turbo says "Required field: <field>". It also needs a fee basis and a project value above 0 ("Project value must be positive.").

## Contacts and divisions are validated

Lines 375-392

Asked as: contact email invalid; division data wrong on submit; client reference error submitting; submission data checks; bad contact details

Submission validates client references, contact emails, the allocations and the project request data before it continues.

## Request allocations must be whole hours in active divisions

Lines 385-391

Asked as: division allocation invalid; inactive division selected; fractional hours allocation; division not valid; allocation division error

A complete request needs division allocations. Hours must be whole numbers ("Allocation hours must be whole numbers."). Each division must exist and be active ("Select a valid division.", "Select an active division.").

## Allocation hours and divisions

Lines 385-431

Asked as: allocation hours decimal; division allocation required; allocations needed to submit; hours must be whole; allocation list problem

Request allocations must have valid hours, and whole hours: 'Allocation hours must be whole numbers.' The division must exist ('Select a valid division.') and be active ('Select an active division.'). Submission requires allocations, but saving a draft does not.

## Authorization PDF required

Lines 395-405

Asked as: missing authorization pdf; upload signed document; authorization file missing; no pdf attached; need PDF to submit

A complete request needs an uploaded authorization PDF with its hash. Otherwise Turbo says "Upload an authorization PDF."

## Allocations optional on draft save

Lines 407-410

Asked as: save draft without allocations; empty allocations allowed; allocations required for submit; draft allocation list empty; skip allocations on save

On a draft save, an empty allocation list passes validation. A submission requires allocations.

## Saving a request is refused when project requests are disabled

Lines 432-442

Asked as: cannot save request; save blocked by setting; requests disabled on save; saving a request fails; project request save off

Creating or editing a project request through the save route first checks that project requests are enabled (IsJobRequestApprovalRequired). When they are off, Turbo refuses with "Project requests are disabled."

## A new project request starts as an Active Draft at version 1

Lines 443-459

Asked as: new request defaults; initial request status; what state is new request; default version of request; creating request starts draft

A new project request is created with status Draft, project status Active and version 1. Turbo sets the creator to the person who saves it and stamps last edited with the current time.

## Creating a request needs time claim

Lines 445-459

Asked as: who can create request; need time claim to create; creating request permission; time claim for new request; new request access

A new request is saved as a Draft with Active project status, version 1 and the signed-in person as creator. Creating one needs the time claim.

## Creating a request needs the time claim

Lines 446-459

Asked as: time claim missing; time claim error message; cannot start new request; create request permission error; no time permission

Creating a request requires the time claim, and the claim error reads "The time claim is required." A new request starts as Draft with project status Active and version 1.

## Editing a request checks visibility, editability and version

Lines 460-471

Asked as: cannot save my edit; save refused existing request; edit not allowed visibility; version check on save; saving changed request fails

When a person saves an existing project request, Turbo loads it only if the person can view it. Turbo then requires that they be the creator and that the request be Draft, Recalled or Rejected. Finally it requires the version sent to match the saved version.

## Save detects changes in business fields, rates and allocations

Lines 472-482

Asked as: what counts as an edit; change detection on request; rates change counted; is request modified; fee basis change tracked

Before applying the edit, Turbo remembers the business fields and the fee basis flag (time_and_materials_set). After applying, it compares each one. A change in any field, in the allocations or in the proposed rates counts as an edit.

## Choose a rate sheet or propose rates, not both

Lines 483-509

Asked as: rate sheet and proposed rates; propose rates and choose sheet; both rate options error; drop proposed rates; rate sheet conflict

A request that proposes new rates cannot also choose a different rate sheet. Turbo says "Choose a rate sheet or propose new rates, not both." Choosing an active rate sheet drops the request's draft of proposed rates.

## Allocations are checked only when the save sends them

Lines 510-523

Asked as: allocations not saved; allocation list ignored; empty allocations on save; allocation change detection; allocations checked when sent

Allocations are compared with the saved ones only when the save includes an allocations list. A sent list is validated on save, and an empty list is accepted. Turbo counts the allocations as changed only when they differ from the saved ones.

## Saving without changes does nothing

Lines 524-533

Asked as: save did nothing; no change saved; version unchanged after save; nothing to save; save without edits

Saving an existing request with no changes to fields, allocations or proposed rates saves nothing and keeps the version.

## Saving creates the request first, then its proposed rates

Lines 534-555

Asked as: proposed rates not saved; order of saving rates; rates saved after request; draft rate sheet deleted; new request rate saving

On a new request, Turbo saves the request first. It then saves any proposed rates, because the draft refers to the request, and saves the request again. On an existing request, Turbo saves the proposed rates first and then marks the request edited. A dropped draft rate sheet is deleted after the save.

## Saving replaces allocations

Lines 556-575

Asked as: allocations replaced; old allocations removed; allocation update overwrites; new allocations deleted old; change division hours

When allocations change, Turbo deletes the old request allocations and saves the new ones.

## Which checklist applies

Lines 586-592

Asked as: which checklist to answer; creation checklist or approval; checklist for status; awaiting creation checklist; what checklist appears

A request Awaiting creation answers the creation checklist. Any other status answers the approval checklist.

## Who can reject a request

Lines 594-612

Asked as: cannot reject request; who can reject; reject button not allowed; only approver rejects; job holder reject rights

Rejecting works only while the request is Awaiting approval or Awaiting creation. The selected approver with the project_authorization claim can reject in both states. A person with the job claim can reject only once the request is Awaiting creation. Others see "Only the selected approver, or a job holder once the request is approved, can reject this request."

## Checklist must be fully confirmed

Lines 614-641

Asked as: checklist not confirmed; confirm checklist question; checklist changed reload; approval checklist unanswered; answer all checklist items

Approval and creation each need the current checklist of that step. Turbo refuses when no checklist is saved, the version is unsupported, or the answer count differs. A page with an older checklist version gets "The checklist changed. Reload the request and review it again." Every answer must be true, else "Confirm: <question>".

## Request actions require enabled requests, visibility and current version

Lines 642-666

Asked as: action refused on request; submit recall approve blocked; version mismatch action; cannot act on request; action permission check

Every action (submit, recall, reject, approve, create) first checks that project requests are enabled. The person must be able to view the request, and the version sent must match the saved version. Otherwise the action is refused.

## Submitting a request

Lines 668-695

Asked as: how to submit request; submit button error; rejected request cannot submit; submission token issued; send request for approval

Submitting needs the time claim and an editable request. A Rejected request must be edited first ("Edit the rejected request before submission."). Submit sets status Awaiting approval, records the submit time and issues a new submission token. It clears earlier approval of the attachment and of proposed rates.

## Recall a request

Lines 696-705

Asked as: withdraw request; take back submitted request; recall my request; only creator recall; unsubmit request

Only the creator can recall ("Only the creator can recall this request."). Recall works only on a request that is Awaiting approval or Awaiting creation ("This action requires a submitted or approved request."). It sets status Recalled, ends the submission token and clears approvals.

## Rejection needs a reason

Lines 706-739

Asked as: reject without reason; rejection reason required; why reject needs text; blank rejection note; enter reason to reject

A rejection needs a non-blank reason ("Enter a rejection reason."). It sets status Rejected and records the time, the rejector and the reason. It ends the submission token and rejects the attachment.

## Approving a request

Lines 740-756

Asked as: approve request not allowed; only approver approves; approve needs submission token; approval refused status; cannot approve my request

Only a submitted, not rejected request (Awaiting approval) can be approved ("Only a submitted request can be approved."). It needs a submission token, else "This request has no submission token. Recall and submit it again." Only the selected approver can approve ("Only the selected approver can approve this request."), and the approver's project limit must cover the value.

## Approval checks limit, checklist and review data

Lines 753-770

Asked as: approval limit too low; approve blocked by checklist; approval validation failed; reviewed data changed reload; approver limit must cover

The approver's project limit must cover the value. They must confirm the approval checklist, and the request must pass full validation. If the review token or document hash differs, Turbo says 'The reviewed data changed. Reload and review it again.'

## Approval checks the reviewed data

Lines 757-788

Asked as: review token missing; document hash mismatch; approval data changed; approve current data; reviewed document changed

The approver must send the review token and document hash of the current data. Otherwise Turbo says "The reviewed data changed. Reload and review it again." Approval then records the approver, the time, the limit (approval_max_value), the request version and the checklist version. It approves any proposed rates (rates_approved) and sets status Awaiting creation.

## Approval records and moves the request on

Lines 771-788

Asked as: what approval does; approved request next step; attachment approved record; request awaiting creation after approval; approval saves review token

Approval marks the attachment Approved with the approver, time, value, the approver's limit (approval_max_value), request version and checklist version. It approves proposed rates, saves a fresh review token and sets the request to Awaiting creation.

## Creating the job from a request

Lines 789-805

Asked as: create job from request; make job from project request; job creation steps; request to job conversion; turn request into job

Creating a job needs the job claim and a request that is Awaiting creation ("This request is not awaiting creation."). It needs a submission token. It creates the job with the request's business fields, the project status as job status and authorizing document PA, and it records the creator and the request times. The job's First Nations agreement is true only when the answer is Yes.

## Creating the job needs job claim and state

Lines 789-804

Asked as: cannot create job; job creation not allowed; need job claim to create; request not ready for job; creation confirmation needed

Creating a job needs the job claim and a request Awaiting creation ('This request is not awaiting creation.'). It also needs a submission token. The person must confirm the creation checklist.

## Creation skips the approver limit check

Lines 805-832

Asked as: limit rechecked at creation; approver limit at job create; creation refused limit lower; job creation value limit; limit check when creating

At creation, the approver's project limit is not rechecked against the value. Turbo compares the saved approval limit with the project value, and refuses if the limit is lower.

## Stale or incomplete approval blocks creation

Lines 806-839

Asked as: approval is stale; saved approval incomplete; checklist invalid on create; recall for new approval; approval missing when creating

Creation refuses with "The approval is stale. Recall the request for a new approval." when the approval is not Approved, or the request version, project value, review token or document hash differs. It refuses with "The saved approval is incomplete. Recall the request for a new approval." when the approver or approval time is missing, the approval limit is below the value, or there is a rejection. It says "The saved checklist is invalid." for an unknown approval checklist version.

## What job creation moves and removes

Lines 840-960

Asked as: what happens when job created; rates renamed on creation; allocations moved to job; request deleted after job; authorization pdf attached to job

Creation activates approved proposed rates and may rename them (rate_sheet_name). It copies the request's allocations to the job and deletes them from the request. It attaches the authorization PDF to the job as its initial approval, number 1. It makes non-blank notes the job's first client note, authored by the request creator. It then deletes the request.

## Creation records the creation checklist and renames rates

Lines 840-846

Asked as: rename proposed rate sheet; creation checklist version; rate name blank; job creator renames rates; approve proposed rates before job

Creation records the creation checklist version. It requires approval of any proposed rates. The job creator can rename the proposed rate sheet, and an empty name keeps the requester's.

## Notifications after actions

Lines 869-875

Asked as: notifications not sent; email after submit; approval notice sent; withdrawn approval notice; who gets notified

Submit, approve and reject send notices after saving. Creation sends a 'created' notice. A rejection by the approver after approval counts as a withdrawn approval, while a job holder's rejection does not.

## What job creation copies and removes

Lines 888-960

Asked as: fields copied to job; job gets request fields; fn agreement on job; what carries over to job; creation copies data

Creation activates approved proposed rates and copies the business fields except the approver. The job gets the request's project status as its status, fn_agreement set from 'Yes', and authorizing document PA. It records the request creator, created and submitted times, and the creating user. Allocations move to the job, the attachment moves to the job as its first approval (sequence 1) and the request is deleted.

## Request notes become a client note

Lines 962-983

Asked as: notes copied to job; request notes client note; note on new job; client note author; notes become job note

If the request notes are not blank, creation makes a client note on the new job. The note keeps the request creator as its author.

## Who can delete a request

Lines 984-1030

Asked as: delete project request; cannot delete request; only abandoned deleted; remove my request; admin delete abandoned request

The creator can delete a request they can still edit (Draft, Recalled or Rejected). Anyone else needs the admin claim and the request must be abandoned ("Only abandoned requests can be deleted."). Deleting also removes its attachment and allocations.
