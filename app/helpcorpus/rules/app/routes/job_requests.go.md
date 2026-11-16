<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_requests.go","sha256":"69e31b091c7ba9cc32374cbf37f3ed6aa6f5a20ae082eec66da53277402b4e1c","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_requests.go

## Abandonment window for project requests uses a day count

Lines 27-27

Asked as: request not edited for days; old project request expires; when request auto closes; stale request cutoff days

Turbo treats a project request as abandoned when its last edit (last_edited) is not after now minus JOB_REQUEST_ABANDONMENT_DAYS days. A request edited exactly on that boundary counts as abandoned. The day count comes from utilities.JOB_REQUEST_ABANDONMENT_DAYS, defined in another file.

## Project request routes require a signed-in user

Lines 55-84

Asked as: request page not logged in; need to log in; api request login required; attachment route sign in

Every project request route and every job attachment route under /api/job_requests and /api/job_attachments needs a signed-in user. The project authorization approvers list route needs sign-in too.

## Which actions a project request supports

Lines 59-72

Asked as: what can I do with request; request available actions; buttons on project request; what request supports

A project request supports list, options, save (create and update), compare rates, detail, print, duplicates, document upload and removal, and the actions submit, recall, reject, approve and create. It also supports delete.

## Validation failures on a request give a generic message

Lines 86-97

Asked as: request form error message; why won't my request save; field errors on request; attachment validation message

When field validation fails on a project request or attachment route, Turbo answers "Check the request fields." with the field errors.

## Only active staff can use project requests

Lines 103-115

Asked as: inactive account can't use requests; why can't I access requests; deactivated staff request access; not active staff member

A person who is not signed in sees "Sign in to continue." An inactive staff member sees "An active staff member is required."

## Request routes need sign-in

Lines 103-106

Asked as: not signed in request; session expired request page; logged out project request; request page sign in

Every project request route needs a signed-in user. An unauthenticated caller gets "Sign in to continue."

## Claim required for a workflow action

Lines 116-128

Asked as: missing permission for request action; don't have the claim; no permission to approve; claim needed for request

A project request action that needs a claim refuses a person without it. The message is "The <claim> claim is required." An inactive person gets "An active staff member is required." first.

## Project requests can be switched off

Lines 129-138

Asked as: project requests turned off; request feature disabled; can't save request feature off; project requests not required

When project requests are not required (IsJobRequestApprovalRequired is false), Turbo refuses to save, act on or delete-by-creator-edit a project request. It shows "Project requests are disabled."

## Abandoned request definition

Lines 139-141

Asked as: when is request abandoned; how old before abandoned; definition of abandoned request; abandoned request meaning

A request is abandoned when its last edit is JOB_REQUEST_ABANDONMENT_DAYS days old or older. A request edited exactly that long ago counts as abandoned.

## Who can view a project request

Lines 142-173

Asked as: who can see my request; can't open someone's request; request access denied; who can view request

The creator can view a request. Procurement (the job claim) can view any request. A project authorization holder (the project_authorization claim) can view a request only when selected as its approver. An admin can view a request only when it is abandoned. Anyone else sees "You cannot view this request."

## Missing request is not found

Lines 160-164

Asked as: request id not found; request disappeared; project request missing; request doesn't exist

If the request does not exist, Turbo says "Project request not found."

## Only creator edits, and only in some states

Lines 174-189

Asked as: can't edit my request; edit button locked; editing someone else's request; only creator edits request

Only the creator can edit a request. Otherwise Turbo says "Only the creator can edit this request." A request can be edited only when it is Draft, Recalled or Rejected. In any other status Turbo says "Recall this request before editing it."

## A project request cannot be edited without recall

Lines 184-188

Asked as: why can't I edit request; request locked for editing; edit blocked by status; can't change request now

A creator who edits a request that is not Draft, Recalled or Rejected gets "Recall this request before editing it." Only these three statuses can be edited.

## Stale request version is refused

Lines 201-208

Asked as: someone else changed request; version conflict on save; request was updated already; reload request before saving

Saves and actions must carry the request's current version. Otherwise Turbo says "This request changed. Reload it before continuing."

## Editing a request resets its approvals

Lines 209-246

Asked as: what happens after I edit; edit resets approval; approval cleared after edit; edit sends request back

A saved edit raises the request version, updates last edited, and clears the submission token. A Rejected request returns to Draft, and its rejection details are cleared. Any authorization PDF record is reset to Pending, with its approval cleared and its project value copied from the request.

## Project value must be a whole dollar amount

Lines 247-262

Asked as: project value cents error; value must be whole dollars; project value not accepted; project value decimal

Project value must be a whole number, zero or more, up to 9007199254740991. Otherwise Turbo says "Project value must be a whole CAD dollar amount."

## Only listed fields are editable

Lines 251-258

Asked as: field can't be changed; read only request field; not editable field error; can't edit this field

A request accepts only its business fields. Any other field gives "Field is not editable: " and the field name.

## Request notes limits

Lines 263-274

Asked as: notes too short; request notes character limit; notes validation error; notes length rule

Notes must be text. Turbo trims them. They must be blank or 10 to 1000 characters, otherwise Turbo says "Notes must be 10 to 1000 characters, or blank."

## Fee basis must be chosen

Lines 275-280

Asked as: fee basis required; time and materials choice; pick time and materials; fee basis not selected

Time and materials must be true or false, otherwise Turbo says "Select a fee basis." Setting it marks the fee basis as chosen.

## Project status must be Active or Closed

Lines 289-291

Asked as: project status active closed; status dropdown error; starting status of request; project status invalid

A request's project status must be Active or Closed. Otherwise Turbo says "Select Active or Closed for the project status." A new request starts as Active.

## Client references and approver limit checked at draft

Lines 292-294

Asked as: client reference check; client references error; contact client references request; client references on save

On every save, Turbo checks the client references of the request.

## Approver's limit is checked when saving

Lines 295-299

Asked as: approver limit not enough; approver can't cover value; approver limit save error; approver project limit low

When a request has an approver, saving checks that the approver's project limit covers the project value (ProjectAuthorizationLimit). Job creation skips this check.

## Parent project rules for a subjob

Lines 300-311

Asked as: subjob parent project; parent project invalid; subjob client mismatch; parent must be top level

A parent must exist, or Turbo says "Select a valid parent project." It must be top-level, with a number like 25-123, or Turbo says "Only a top-level project can be a parent." The request's client must match the parent's client, or Turbo says "The subjob must use the parent project's client."

## Linked proposal rules

Lines 312-325

Asked as: proposal link invalid; linked proposal wrong status; proposal must match client; proposal number starts P

A linked proposal must exist ("Select a valid proposal."). Its number must start with P and its client must match the request's ("Select a proposal for this client."). Its status must be In Progress, Submitted or Awarded ("Select an In Progress, Submitted, or Awarded proposal.").

## Fields required to submit

Lines 332-342

Asked as: required fields to submit; why can't I submit; missing fields before submit; fields needed to approve

To submit or approve, these must be filled: description, location, client, contact, manager, branch, rate sheet, project award date, invoicing information, approver and fn agreement answer. Otherwise Turbo says "Required field: " and the field. A fee basis must be chosen, and project value must be above zero ("Project value must be positive.").

## Linked proposal must be Awarded

Lines 343-351

Asked as: linked proposal not awarded; proposal must be awarded; awarded proposal needed to create; proposal status blocks submit

To submit, approve or create, a linked proposal must be Awarded. Otherwise Turbo says "Set the linked proposal to Awarded before continuing."

## Request data needs contact emails, allocations and hook checks

Lines 352-371

Asked as: contact emails before submit; allocations needed to submit; submit data validation fails; checks before submitting request

Before a request is submitted, approved or created, Turbo checks its client references, contact emails, division allocations and the project request data rules (ValidateProjectRequestData). Allocations are required at this point.

## Authorization PDF required

Lines 372-382

Asked as: upload authorization PDF; no PDF attached to request; missing authorization document; submit needs PDF

A request needs an uploaded authorization PDF with a document hash to be submitted, approved or created. Otherwise Turbo says "Upload an authorization PDF."

## Allocation hours and divisions

Lines 384-408

Asked as: allocation hours invalid; allocation decimal hours error; division not valid allocation; inactive division allocation

Allocation hours must be whole numbers ("Allocation hours must be whole numbers.") and valid numbers ("Invalid allocation hours."). Each division must exist ("Select a valid division.") and be active ("Select an active division."). A complete request must have allocations. When saving, an empty allocation list is allowed.

## Creating a request needs the time claim

Lines 423-436

Asked as: who can create request; need time claim to create; create request permission; can't create new request

Creating a request needs the time claim. A new request starts as Draft, Active and version 1.

## New request starts as an Active draft

Lines 431-436

Asked as: new request default status; starting status of request; first version of request; new request starts draft

A new project request starts as a Draft with project status Active, version 1, and the creator set to the person who saves it.

## Editing a request needs the creator and an editable state

Lines 437-448

Asked as: who can update request; only creator can update; update request permission; edit requires creator and draft

When updating an existing project request, Turbo first loads it only if the person can view it. Then it requires that the person is the creator and that the status is Draft, Recalled or Rejected.

## Request save refuses a stale version

Lines 445-447

Asked as: version mismatch saving request; save rejected old version; stale request version save; saved version out of date

When a person updates a project request, the version sent must equal the saved version. Otherwise Turbo refuses with 'This request changed. Reload it before continuing.'

## Save without changes writes nothing

Lines 449-510

Asked as: save without changes; no change save request; nothing changed on update; save doesn't bump version

When an update to an existing project request changes no business field, allocation or proposed rates, Turbo saves nothing. The version and last-edited date stay as they were.

## Request fields are checked in a set order on save

Lines 454-500

Asked as: order of checks on save; why rate sheet error after save; save validation order; save checks sequence

On save, Turbo applies the submitted fields, then checks the draft rules (project status, client references, approver limit, parent, proposal). Only then does it check the rate sheet choice and the allocations.

## Rate sheet or proposed rates, not both

Lines 460-486

Asked as: rate sheet and proposed rates; choose rate sheet or proposed rates; can't have both rates; rate sheet conflict

A request cannot both choose a rate sheet and propose new rates. Turbo says "Choose a rate sheet or propose new rates, not both." Choosing another sheet drops the request's own proposed-rates draft.

## Chosen rate sheet replaces a request's own proposed rates draft

Lines 466-486

Asked as: change rate sheet deletes proposed rates; proposed rates lost; switching rate sheet removes draft; proposed rates draft gone

When a person picks a different rate sheet, Turbo deletes the proposed rates draft that the request owns. When proposed rates are sent, the request's rate sheet is set to that draft.

## Unchanged save does nothing

Lines 487-552

Asked as: saving unchanged request; no version change on save; save with no edits; unchanged request save

Saving an existing request with no change to its fields, allocations or proposed rates does not save or bump the version. Allocations, when changed, are replaced in full.

## Request allocations are replaced when they change

Lines 495-552

Asked as: allocations changed on save; allocation rows deleted; empty allocations on save; allocations replaced when edited

When the allocations sent differ from those saved, Turbo deletes the old allocation rows and saves the new ones. Allocations are validated on save even when the list is not required, and omitting them leaves the saved ones untouched.

## Checklist scope depends on request status

Lines 563-569

Asked as: checklist for awaiting creation; which checklist shows; approval or creation checklist; checklist depends on status

A request with status Awaiting creation uses the creation checklist. Any other status uses the approval checklist.

## Who can reject a request

Lines 571-589

Asked as: who can reject request; reject permission for approver; job holder reject request; can't reject this request

The selected approver (with the project_authorization claim) can reject while the request is Awaiting approval or Awaiting creation. A job holder (job claim) can reject only once it is Awaiting creation. Others see "Only the selected approver, or a job holder once the request is approved, can reject this request."

## Checklist must be fully confirmed

Lines 591-618

Asked as: checklist questions not answered; tick all checklist boxes; checklist confirm error; checklist version changed

Approval and creation each need every question of the current checklist answered true, with the current version. Missing checklist gives "No <scope> checklist is saved." An old version gives "The checklist changed. Reload the request and review it again." Wrong answer count gives "Use only the N checklist answers for version V." An unticked answer gives "Confirm: " and the question.

## Submitting a request

Lines 645-672

Asked as: how to submit request; send request for approval; submit button not working; submit request steps

Submitting needs the time claim, and the request must be editable. A Rejected request must be edited first ("Edit the rejected request before submission."). Submit sets status to Awaiting approval, records submitted time, issues a new submission token, and clears earlier approvals.

## Submit refuses a rejected request that was not edited

Lines 652-654

Asked as: rejected request submit fails; resubmit rejected request; edit rejected before submit; can't submit rejected request

Submitting a request with status Rejected is refused with "Edit the rejected request before submission." An edit moves a Rejected request back to Draft.

## Recall by the creator

Lines 673-682

Asked as: recall my request; take back submitted request; withdraw request from approval; undo submit request

Only the creator can recall, and only an Awaiting approval or Awaiting creation request. Otherwise Turbo says "This action requires a submitted or approved request." or "Only the creator can recall this request." Recall sets status to Recalled and ends the submission token and approvals.

## Rejection needs a reason

Lines 683-716

Asked as: reject reason required; rejection needs comment; blank rejection reason; why reject fails

Rejecting needs a non-blank reason, otherwise Turbo says "Enter a rejection reason." Reject sets status to Rejected, records rejector, time and reason, clears the submission token and marks the PDF Rejected.

## Approving a request

Lines 717-765

Asked as: approve project request; approver can't approve; approve button request; who can approve request

Only the selected approver can approve ("Only the selected approver can approve this request."). The request must be Awaiting approval and not rejected ("Only a submitted request can be approved."). It needs a submission token, or Turbo says "This request has no submission token. Recall and submit it again." The approver's project limit must cover the value. Approval sets status to Awaiting creation and records the approver's limit.

## Approval needs a submission token

Lines 724-726

Asked as: no submission token; approve without token; submission token missing; submission token error

Approving or creating from a request that has no submission token is refused with "This request has no submission token. Recall and submit it again."

## Review must match what was approved

Lines 741-757

Asked as: reviewed data changed; review token mismatch; approve needs proposed rates; data changed before approval

Approving needs the review token and document hash of the data the approver saw. Otherwise Turbo says "The reviewed data changed. Reload and review it again." Approval also requires approving proposed rates if any (RatesApproved).

## Creating the job from a request

Lines 766-792

Asked as: create job from request; turn request into job; make job button; job creation from request

Creating needs the job claim and status Awaiting creation ("This request is not awaiting creation."). It needs a submission token. It does not re-check the approver's limit. It makes a job from the request, using status project_status and authorizing document PA.

## Stale or incomplete approval blocks creation

Lines 790-802

Asked as: approval stale creation fails; can't create job approval; saved approval incomplete; approval checklist invalid

Creation fails with "The approval is stale. Recall the request for a new approval." if the PDF is not Approved, or version, value or review token changed. It fails with "The saved approval is incomplete. Recall the request for a new approval." if the approver, approval time or approval limit does not cover the value, or any rejection is present. It fails with "The saved checklist is invalid." if the approval checklist version is gone.

## What job creation transfers

Lines 803-923

Asked as: what happens when job created; what job creation copies; request transferred to job; job creation transfers

Creating a job activates approved proposed rates, sets fn_agreement true only when the answer is Yes, and copies the request's allocations to the job. The authorization PDF moves to the job as its first approval. The notes become the job's first client note by the request creator. The request is deleted.

## Project request fields that are copied to the job

Lines 862-868

Asked as: fields copied to job; job gets request fields; fn agreement set on job; job status from request

Job creation copies each project request business field (requestBusinessFields) that the job has, except the approver, to the new job. It sets the job status from the project status (project_status). It sets the First Nations agreement flag (fn_agreement) to true only when the answer (fn_agreement_answer) is "Yes".

## Created job note and request cleanup

Lines 925-946

Asked as: notes become client note; request notes job note; request deleted after job; notes copied to job

When the request has notes, job creation makes them the first client note of the new job, authored by the request creator. Blank notes make no note. The request, once the job exists, is deleted.

## Deleting a request

Lines 947-995

Asked as: delete project request; who can remove request; only abandoned deleted; remove request and allocations

The creator can delete a request that they can still edit. Otherwise only an admin can delete it, and only if abandoned. Others see "Only abandoned requests can be deleted." Deleting removes the PDF record and allocations.
