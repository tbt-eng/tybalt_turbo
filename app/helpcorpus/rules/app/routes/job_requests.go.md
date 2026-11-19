<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_requests.go","sha256":"e0821991b31576183a6897b7951ae30f5a85de793d83397d82ebc9e298460827","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_requests.go

## Fields a project request holds and copies to the job

Lines 30-30

Asked as: fields copied to job; which fields carry over; what gets copied to job; editable request fields

A project request holds these business fields: project status, approver, client, job owner, contact, manager, alternate manager, branch, rate sheet, parent, proposal, invoicing information, description, notes, location, client PO, client reference number, project award date, project completion date, project value, fee basis (time_and_materials) and the First Nations agreement answer (fn_agreement_answer). A person can edit only these fields. Job creation copies them to the job, except the approver.

## Project request and attachment endpoints need sign-in

Lines 55-85

Asked as: anonymous call refused; not logged in to api; need to sign in; signed out access denied

Every project request endpoint requires a signed-in staff user. The job attachment endpoints and the project authorization approvers list need one too. Turbo refuses an anonymous call.

## Project request workflow actions available

Lines 68-73

Asked as: request action buttons; what buttons does request have; request endpoints list; submit or recall options

A project request has five workflow actions: submit, recall, reject, approve and create. Each is a separate POST on the request. The request also has save, delete, document upload and removal, duplicates, rate comparison and rate finding endpoints.

## Field validation errors show a generic message

Lines 87-98

Asked as: check the request fields error; which fields failed validation; generic form error message; list of invalid fields

When a project request or attachment call fails field validation, Turbo shows 'Check the request fields.' and lists the failing fields. Errors that carry their own message and status keep that message.

## Only active staff can use project request actions

Lines 104-116

Asked as: not active staff error; inactive account cannot act; account not active message; not signed in message

A person who is not an active staff member gets "An active staff member is required." A person who is not signed in gets "Sign in to continue."

## Claim check also needs an active staff member

Lines 104-129

Asked as: missing claim error; claim required message; inactive user permission check; no permission for action

When Turbo checks a claim for a project request action, the person must first be active staff. A signed-out person gets 'Sign in to continue.' An inactive person gets 'An active staff member is required.' An active person without the claim gets 'The <claim> claim is required.'

## Project requests can be switched off

Lines 130-139

Asked as: project requests disabled; requests switched off; cannot save requests turned off; approval not required setting

When project request approval is not required (IsJobRequestApprovalRequired), Turbo refuses saving, submitting, recalling, rejecting, approving, creating and deleting a project request. The message is "Project requests are disabled."

## Abandoned request definition

Lines 140-142

Asked as: when request is abandoned; old unedited request; stale draft after days; request untouched for days

A project request is abandoned when its last edit is 'JOB_REQUEST_ABANDONMENT_DAYS' days old or older. The comparison includes the exact cutoff moment.

## Who can view a project request

Lines 143-174

Asked as: cannot view this request; who can see request; view request denied; admin can see abandoned

The creator can view a project request. Anyone with the job claim (procurement) can view any request. A project authorization approver can view requests that name them as approver. An admin can view only abandoned requests. Anyone else gets "You cannot view this request."

## Only the creator edits, and only in editable states

Lines 175-190

Asked as: only creator can edit; cannot edit someone else's request; recall before editing; request locked for editing

Only the creator can edit a project request, and only when its status is Draft, Recalled or Rejected. Others get "Only the creator can edit this request." In any other status the message is "Recall this request before editing it."

## Stale request version is refused

Lines 202-209

Asked as: request changed reload; someone else edited it; outdated version error; version mismatch on save

If the version sent with an edit or action differs from the saved request version, Turbo refuses it with "This request changed. Reload it before continuing."

## Editing a request resets approval and token

Lines 210-247

Asked as: edit resets approval; editing clears submission; rejected goes back to draft; edit clears pdf approval

A saved business edit raises the request version by 1 and sets last edited to now. It clears the submission token and any rejection details. A Rejected request goes back to Draft. It also resets the authorization PDF's approval to Pending and updates its project value.

## Project value must be a whole dollar amount

Lines 248-263

Asked as: project value decimal error; value must be whole number; cents not allowed in value; negative project value

Project value on a request must be a whole CAD dollar number of 0 or more, no larger than 9007199254740991. Otherwise Turbo says "Project value must be a whole CAD dollar amount." A draft can hold 0. Submitting requires more than 0.

## Only listed request fields are editable

Lines 252-260

Asked as: field not editable error; cannot change this field; field refused on save

A save with a field outside the allowed business fields is refused with "Field is not editable: " and the field name.

## Request notes length limit

Lines 264-275

Asked as: notes too short; notes character limit; notes must be text; notes length error

Request notes are trimmed. They must be 10 to 1000 characters, or blank. Otherwise Turbo says "Notes must be 10 to 1000 characters, or blank." Non-text notes get "Notes must be text."

## Fee basis must be chosen

Lines 276-281

Asked as: select fee basis; time and materials choice; fee basis required; no fee basis chosen

The time and materials field must be true or false, or Turbo says "Select a fee basis." Saving it marks the fee basis as chosen. Submitting without a choice gives the same message.

## Project status must be Active or Closed

Lines 290-292

Asked as: project status must be; active or closed only; status error on request; select active status

A request's project status must be Active or Closed. Otherwise Turbo says "Select Active or Closed for the project status."

## Approver's project limit is checked on save

Lines 296-300

Asked as: approver limit exceeded; approver cannot approve this amount; over approver spending limit; approver value cap

When a request names an approver, saving checks that approver's project authorization limit against the project value (ProjectAuthorizationLimit). Approval and submission check it too. Job creation skips this check.

## Parent project rules for subjob requests

Lines 301-312

Asked as: parent project invalid; subjob client mismatch; parent must be top level; subproject parent number error

A parent must be an existing project whose number looks like 25-123 or 25-1234 (top-level). A subjob must use the parent's client. The messages are "Select a valid parent project.", "Only a top-level project can be a parent." and "The subjob must use the parent project's client."

## Linked proposal rules when saving

Lines 313-326

Asked as: linked proposal not valid; proposal client mismatch; proposal status wrong; proposal must start with P

A linked proposal must be a job whose number starts with P and has the request's client. Its status must be In Progress, Submitted or Awarded. Otherwise Turbo asks for a valid proposal for this client or an In Progress, Submitted, or Awarded proposal.

## Required fields before submit, approve or create

Lines 333-343

Asked as: required field missing on submit; cannot submit blank fields; missing field before approve; project value must be positive

Submitting, approving or creating needs description, location, client, contact, manager, branch, rate sheet, award date, invoicing information, approver and fn agreement answer. A blank one gives "Required field: " and the field name. Project value must be above 0 ("Project value must be positive.").

## Linked proposal must be Awarded to continue

Lines 344-352

Asked as: proposal not awarded; set proposal to awarded; linked proposal blocks submit; proposal status stops submission

A request linked to a proposal can be submitted, approved or created only when the proposal's status is Awarded. Otherwise Turbo says "Set the linked proposal to Awarded before continuing."

## Request details checked before submit, approve and create

Lines 353-372

Asked as: contact email invalid; client reference check failed; request data validation failed; submit blocked by checks

Before a project request is submitted, approved or created, Turbo checks the contact emails (ValidateRequestContactEmails), the client references and the project request data (ValidateProjectRequestData). These checks are defined in other files. A failure blocks the action.

## Authorization PDF required

Lines 373-384

Asked as: no pdf uploaded; upload authorization form; missing pdf on submit; authorization document missing

Submitting, approving or creating needs an uploaded authorization PDF with a document hash. Otherwise Turbo says "Upload an authorization PDF."

## Division allocation rules

Lines 385-409

Asked as: allocation hours decimal; invalid division allocation; inactive division selected; allocation must be whole hours

Allocations must be whole hours ("Allocation hours must be whole numbers."), on an existing division ("Select a valid division.") that is active ("Select an active division."). Submitting requires validated allocations. Saving a draft with no allocations sent skips this check.

## Saving a request: one transaction, creator and defaults

Lines 410-437

Asked as: who created the request; new request default status; saving new request creator; requests disabled on save

A new project request starts as a Draft with project status Active and version 1. Turbo records the person who creates it as its creator. Turbo refuses to save when project requests are disabled.

## Creating a project request needs the time claim

Lines 424-437

Asked as: cannot create request; need time claim; no permission to create; create request permission

Only a person with the time claim can create a project request or submit one. A new request starts as Draft with project status Active and version 1.

## Editing an existing request needs visibility, edit rights and current version

Lines 438-449

Asked as: cannot save request changes; save fails cannot view; edit existing request permission; save refused for request

When a person saves changes to an existing request, Turbo loads the request and checks that they can view it. It then checks that they can edit it and that the version they send is the current one. Any failure refuses the save.

## Rate sheet choice or proposed rates, not both

Lines 461-487

Asked as: rate sheet and proposed rates; choose one rate option; proposed rates conflict; both rate sheet and rates

A request either chooses an active rate sheet or proposes new rates. Doing both gives "Choose a rate sheet or propose new rates, not both." Choosing a different sheet deletes the request's own draft of proposed rates.

## Allocation changes are validated only when allocations are sent

Lines 488-502

Asked as: allocations not checked on save; empty allocations on draft; allocation check skipped

On save, Turbo checks division allocations only when the person sends an allocation list. An empty list is allowed on a draft save. A list that matches the saved allocations counts as no change.

## Unchanged save does nothing

Lines 503-511

Asked as: save did not change version; nothing changed on save; last edited date unchanged; save without changes

When an existing request is saved with no change to its fields, allocations or proposed rates, Turbo does not save it. The version and last edited date stay the same.

## Saving a new request also saves proposed rates

Lines 512-519

Asked as: proposed rates on new request; saving rates twice; rates draft saved after request

When a new request is created with proposed rates, Turbo saves the request first, then saves the proposed rates draft, then saves the request again. This is because the draft refers to the request.

## Editing a request deletes a dropped rates draft

Lines 528-533

Asked as: rates draft removed; proposed rates deleted; draft rates disappeared; deleted draft after choosing sheet

After a save, Turbo deletes the request's draft of proposed rates when the person chose an active rate sheet in its place.

## Saving replaces allocations

Lines 534-553

Asked as: allocations replaced on save; old allocations deleted; allocation list overwritten; new allocations replace old

When the allocations sent differ from the saved ones, Turbo deletes the saved allocations and creates the new ones.

## Which checklist a request action answers

Lines 564-570

Asked as: which checklist applies; awaiting creation checklist; approval checklist questions; checklist for request status

A request in 'Awaiting creation' status answers the creation checklist. A request in any other status answers the approval checklist.

## Who can reject and reason required

Lines 572-590

Asked as: who can reject request; rejection reason missing; job holder rejects request; reject reason required

Rejecting works in Awaiting approval or Awaiting creation. The selected approver (with the project_authorization claim) can reject, and a job holder can reject once it is Awaiting creation. Anyone else gets "Only the selected approver, or a job holder once the request is approved, can reject this request." A reason is required ("Enter a rejection reason.").

## Approval needs a complete checklist

Lines 592-619

Asked as: checklist not answered; checklist version not supported; confirm checklist question; checklist changed reload

Approving or creating needs the current checklist version, with every question answered true. Messages: "No <scope> checklist is saved.", "Checklist version N is not supported.", "Use only the N checklist answers for version V." and "Confirm: " with the question key. An out-of-date version gets "The checklist changed. Reload the request and review it again."

## Actions need the current request version

Lines 620-644

Asked as: version mismatch on action; request changed before approve; reload before submitting; old version on action

Submit, recall, reject, approve and create each load the request, and the person must be able to view it. The version they send must match the request's version. Otherwise Turbo says 'This request changed. Reload it before continuing.'

## Submitting a request

Lines 646-673

Asked as: submit request steps; cannot submit rejected request; submission token created; edit rejected before submit

Submitting needs the time claim, an editable request and complete data. It sets status Awaiting approval, records the submitted time and issues a new random submission token. It clears earlier approval of the PDF and of proposed rates. A Rejected request must be edited first ("Edit the rejected request before submission.").

## Recalling a request

Lines 674-710

Asked as: recall request; only creator recalls; pull back submitted request; withdraw submitted request

Only the creator can recall a request ("Only the creator can recall this request."). It must be Awaiting approval or Awaiting creation, otherwise "This action requires a submitted or approved request." Recall sets status Recalled, and clears the submission token and approvals.

## Rejection effects

Lines 698-717

Asked as: what happens when rejected; rejected request status; reject clears token; pdf marked rejected

Rejecting sets status Rejected and records the time, rejector and trimmed reason. It clears the submission token and marks the authorization PDF Rejected. The creator must edit and resubmit.

## Approving a request

Lines 718-766

Asked as: only approver can approve; no submission token; approve awaiting approval; approval fails for request

Only the selected approver can approve ("Only the selected approver can approve this request."). The request must be Awaiting approval and not rejected ("Only a submitted request can be approved."). Without a submission token Turbo says "This request has no submission token. Recall and submit it again." Approval moves it to Awaiting creation.

## Approval records the reviewed data

Lines 742-765

Asked as: reviewed data changed; review token mismatch; approval records approver limit; data changed since review

Approval needs the review token and the document hash to match the current data and PDF. Otherwise Turbo says "The reviewed data changed. Reload and review it again." Approval records the approver, the time, the approver's limit (approval_max_value), the request version and the checklist version.

## Creating a job from an approved request

Lines 767-784

Asked as: create job from request; not awaiting creation; job holder creates job; creation checklist required

Only a job holder (the job claim) can create the job, and only when the request is Awaiting creation ("This request is not awaiting creation."). A missing submission token gives "This request has no submission token. Recall and submit it again." A creation checklist must be completed.

## Stale or incomplete approval blocks creation

Lines 791-803

Asked as: approval is stale; saved approval incomplete; approval expired before create; saved checklist invalid

Creation is refused with "The approval is stale. Recall the request for a new approval." if the PDF is not Approved, the version, value or review token differs, or the document hash differs. It says "The saved approval is incomplete. Recall the request for a new approval." if the approver, approval time or limit does not cover the value, or a rejection exists. A missing approval checklist gives "The saved checklist is invalid."

## Creating a job from a request

Lines 804-823

Asked as: job creation rate sheet name; proposed rates approval on create; creation checklist version saved; approval check on create

When a job holder (the job claim) creates the job, Turbo records an approval-checklist check. It saves the creation checklist version on the authorization PDF record. It requires approval of any proposed rates and applies the name the job creator gives the rate sheet. It then creates the job.

## Notices are sent after the action succeeds

Lines 825-851

Asked as: email after approval; submitted notice sent; withdrawal notice sent; rejection email sent

After an approved submit, approve or reject, Turbo queues a notice (submitted, approved or rejected). Turbo sends it after the transaction commits. A rejection by the selected approver after approval counts as a withdrawal of the approval.

## What job creation produces

Lines 852-924

Asked as: job fields after creation; fn agreement true; allocations move to job; request deleted after job

Creating a job copies the request's business fields except the approver. The job gets the request's project status as its status, and the fn agreement is true only when the answer is "Yes". The authorizing document is set to PA. The job records the request's creator, created and submitted times, the person who created it and the completion time. Allocations move to the job and the request is deleted.

## Authorization PDF moves to the job as first approval

Lines 903-912

Asked as: authorization pdf attached to job; pdf moves to job; first approval sequence; initial approval on job

On job creation the authorization PDF leaves the request and attaches to the new job. It is marked as the initial approval, with approval sequence 1.

## Request notes become a client note

Lines 926-947

Asked as: request notes client note; notes copied to job; first note on job; notes become job note

When creating the job, non-blank request notes become the job's first client note, written by the request creator. Blank notes make no note.

## Deleting a project request

Lines 948-989

Asked as: delete project request; cannot delete request; only admin deletes abandoned; delete draft request

The creator can delete a request they can still edit (Draft, Recalled or Rejected). Otherwise only an admin can delete it, and only when abandoned ("Only abandoned requests can be deleted."). Deleting also deletes its PDF and allocations.
