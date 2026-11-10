<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_requests.go","sha256":"9025509bb920756e6ae9c60d6a138f3f5c7b2c124deb6b968e1501545345e9f4","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_requests.go

## Fields a project request can store

Lines 30-30

A project request keeps only these business fields: project_status, approver, client, job_owner, contact, manager, alternate_manager, branch, rate_sheet, parent, proposal, invoicing_information, description, notes, location, client_po, client_reference_number, project_award_date, project_completion_date, project_value, time_and_materials and fn_agreement_answer. Turbo refuses any other field on a save.

## Only signed-in users reach project request routes

Lines 46-87

All project request and job attachment routes require a signed-in user (the users collection). Validation errors on these routes show as 'Check the request fields.'

## Only active staff can act on project requests

Lines 93-105

Turbo refuses a signed-in user who is not an active staff member with "An active staff member is required.". A user who is not signed in gets "Sign in to continue.".

## Claim-restricted actions name the missing claim

Lines 106-118

When an action needs a claim and the user lacks it, Turbo says "The <claim> claim is required.". The user must also be active staff.

## Project request actions need the feature enabled

Lines 119-128

Turbo refuses project request actions with "Project requests are disabled." when project request approval is not enabled. The check uses utilities.IsJobRequestApprovalRequired (requestEnabled). Saving, submitting, approving, creating and similar actions all run this check.

## Request is abandoned after a set number of days

Lines 129-131

A request is abandoned when its last_edited date is not later than JOB_REQUEST_ABANDONMENT_DAYS days ago. A request edited exactly at the cutoff counts as abandoned.

## Who can view a project request

Lines 132-163

The creator can view a request. A job holder (the job claim) can view any request. A project authorization approver (the project_authorization claim) can view only a request that names them as approver. An admin can view only an abandoned request. Everyone else gets "You cannot view this request.". A missing request gives "Project request not found.".

## Only the creator can edit, only in editable states

Lines 164-179

Only the creator can edit a request, otherwise "Only the creator can edit this request.". Editing is allowed in Draft, Recalled or Rejected status. In any other status Turbo says "Recall this request before editing it.".

## Stale request version is refused

Lines 191-198

Save and workflow actions compare the version the user sent with the request's current version. A mismatch gives "This request changed. Reload it before continuing.".

## Editing a request resets approval and submission

Lines 199-236

A business edit raises the request version by 1 and updates last_edited. It clears the submission token and any rejection details. A Rejected request returns to Draft. The attached authorization PDF's approval is cleared, its status goes back to Pending, and its project value is copied from the request.

## Project value must be a whole non-negative number

Lines 237-252

A request's project value must be a whole CAD dollar amount of 0 or more, and no larger than 9007199254740991. Otherwise Turbo says "Project value must be a whole CAD dollar amount.".

## Only listed fields are editable on a request

Lines 241-248

Turbo refuses any field outside the business field list with "Field is not editable: <field>".

## Request notes length and type

Lines 253-264

Notes must be text. Turbo trims them. They must be blank or 10 to 1000 characters, otherwise "Notes must be 10 to 1000 characters, or blank.".

## Fee basis must be true or false

Lines 265-270

The time_and_materials value must be a boolean, otherwise "Select a fee basis.". Setting it marks the fee basis as chosen.

## Draft requires Active or Closed project status

Lines 279-281

Saving a request requires project status Active or Closed, otherwise "Select Active or Closed for the project status.".

## Approver's project limit is checked on save

Lines 285-289

When a request names an approver, Turbo checks that approver's project authorization limit against the project value using utilities.ProjectAuthorizationLimit. A failure is returned as the error from that function. Creation of the job skips this check.

## Parent project rules for a subjob

Lines 290-301

A parent must exist, otherwise "Select a valid parent project.". Its number must look like 25-123 or 25-1234, otherwise "Only a top-level project can be a parent.". The request's client must match the parent's client, otherwise "The subjob must use the parent project's client.".

## Linked proposal rules on a request

Lines 302-315

A proposal must exist, otherwise "Select a valid proposal.". Its number must start with P and its client must match the request's client, otherwise "Select a proposal for this client.". Its status must be In Progress, Submitted or Awarded, otherwise "Select an In Progress, Submitted, or Awarded proposal.".

## Required fields before submitting a request

Lines 322-332

Submit, approve and create need these fields filled: description, location, client, contact, manager, branch, rate sheet, project award date, invoicing information, approver and FN agreement answer. A blank one gives "Required field: <field>". A fee basis must be chosen, otherwise "Select a fee basis.". Project value must be above 0, otherwise "Project value must be positive.".

## Linked proposal must be Awarded to continue

Lines 333-341

For submit, approve and create, a linked proposal must have status Awarded. Otherwise Turbo says "Set the linked proposal to Awarded before continuing.".

## Request allocations must be valid

Lines 348-398

Allocation hours must be whole numbers, otherwise "Allocation hours must be whole numbers.". Invalid numbers give "Invalid allocation hours.". Each division must exist ("Select a valid division.") and be active ("Select an active division."). A complete request must have allocations that pass validateAllocationInputs.

## Authorization PDF required before submitting

Lines 362-371

A request needs an uploaded authorization PDF with a document hash. Otherwise Turbo says "Upload an authorization PDF.". It also checks the file with ensureAuthorizationFile.

## Creating a request needs the time claim

Lines 413-426

Creating a new request requires the time claim. It starts as Draft with project status Active and version 1.

## Saving an existing request checks view, edit rights, then version

Lines 427-438

When a person saves an existing project request, Turbo first loads it and checks that the person can view it. Then it checks that the person may edit it in its current state. Last, it checks that the version sent matches the saved version.

## A fee basis change counts as a request change

Lines 439-443

Before it applies a save, Turbo records every business field of the request and whether the fee basis is set (time_and_materials_set). After the save it compares the two, and any difference counts as a change to the request.

## Saving a request re-checks the draft rules

Lines 444-449

After Turbo applies the submitted fields to a request, it re-checks the draft rules. These cover project status, client references, the approver's project limit, the parent project and the linked proposal. A failed check refuses the save.

## Allocations replace the old set when changed

Lines 454-501

When the allocations sent differ from the saved ones, Turbo deletes the old allocations and saves the new ones. When no allocations are sent, the old ones stay.

## Saving with no changes does nothing

Lines 464-473

When an existing request is saved with no change to its fields or allocations, Turbo does not save, so the version and last edited date stay the same.

## Which checklist a request action answers

Lines 512-518

A request in Awaiting creation answers the creation checklist. In any other status the approval checklist applies.

## Who can reject a request

Lines 520-538

Rejection is possible only in Awaiting approval or Awaiting creation. The selected approver with the project_authorization claim can reject in both. A job holder (the job claim) can reject only in Awaiting creation.

## Checklist must be fully confirmed and current

Lines 540-567

Every question of the current checklist must be answered true. If none is saved Turbo says "No <scope> checklist is saved.". A version below 1 or above the current gives "Checklist version N is not supported.". An older version gives "The checklist changed. Reload the request and review it again.". A wrong number of answers gives "Use only the N checklist answers for version V.". An answer that is not true gives "Confirm: <key>".

## Every request action needs the feature on, view access and current version

Lines 568-592

Submit, recall, reject, approve and create all require project requests to be enabled. If the request does not exist, Turbo says "Project request not found.". If the person cannot view it, Turbo says "You cannot view this request.". The version sent must match the saved version, or the action fails with "This request changed. Reload it before continuing.".

## Submitting a request

Lines 594-618

Submitting needs the time claim and an editable request, and only the creator can submit. A Rejected request gives "Edit the rejected request before submission.". The request must be complete. Submit sets status Awaiting approval, records the submitted time, issues a new 32-character submission token, and resets the PDF approval.

## Recalling a request

Lines 619-628

Only the creator can recall, otherwise "Only the creator can recall this request.". The request must be in Awaiting approval or Awaiting creation, otherwise "This action requires a submitted or approved request.". Recall sets status Recalled, clears rejection details and the submission token, and resets the PDF approval.

## Rejecting a request

Lines 629-659

Rejecting needs a non-blank reason, otherwise "Enter a rejection reason.". Unauthorized users get "Only the selected approver, or a job holder once the request is approved, can reject this request.". Reject sets status Rejected and records the time, rejector and trimmed reason. It clears the submission token and marks the PDF Rejected.

## Approving a request

Lines 660-676

Only a submitted request (Awaiting approval, not rejected) can be approved, otherwise "Only a submitted request can be approved.". A missing submission token gives "This request has no submission token. Recall and submit it again.". Only the selected approver can approve, otherwise "Only the selected approver can approve this request.". The approver's project limit must cover the project value.

## Approval requires fresh review and sets status

Lines 677-699

Approval needs the current approval checklist and a complete request. The review token and document hash must match, otherwise "The reviewed data changed. Reload and review it again.". Approval marks the PDF Approved with approver, time, project value, the approver's limit, request version and checklist version. The request moves to Awaiting creation.

## Creating a job from an approved request

Lines 700-719

Creating needs the job claim and status Awaiting creation, otherwise "This request is not awaiting creation.". A missing token gives "This request has no submission token. Recall and submit it again.". It needs the creation checklist and a complete request. The approver's limit is not rechecked at creation.

## Stale or incomplete approval blocks creation

Lines 724-736

Creation fails with "The approval is stale. Recall the request for a new approval." when the PDF is not Approved, its version, value or review token differ, or the caller's token or hash differ. It fails with "The saved approval is incomplete. Recall the request for a new approval." when the approver differs, the approval time is missing, the approved limit is below the project value, or a rejection exists. A missing saved checklist gives "The saved checklist is invalid.".

## Creating a job records the creation checklist version

Lines 737-750

When a job holder creates a job from an approved request, Turbo stores the creation checklist version they confirmed on the authorization PDF record (creation_checklist_version). It then creates the job and queues a "created" notice. The response returns the new job's id as both id and job_id.

## Request changes are saved before notices go out

Lines 752-759

For every action except create, Turbo saves the authorization PDF record, if there is one, and then the request. Notices are sent only after the whole transaction succeeds, so a failed action sends none.

## Which request actions send a notice

Lines 760-772

Submit, approve and reject each queue a notice (submitted, approved, rejected). Recall sends none. A reject counts as a withdrawn approval only when the request is Awaiting creation and the person rejecting is the selected approver. When a job holder rejects, it does not count as withdrawn.

## What job creation copies and sets

Lines 784-800

The new job copies the request's business fields except the approver. Its status is the request's project status. FN agreement is true only when the answer is Yes. Its authorizing document is set to "PA". It records the request creator, created and submitted dates, who created it, and the completion time.

## Job creation moves allocations, PDF and note

Lines 805-846

Request allocations become job time allocations and are deleted from the request. The authorization PDF moves to the job as its initial approval with approval sequence 1. Non-blank notes become the job's first client note, written by the request creator. The request is then deleted.

## Request notes become the job's first client note

Lines 848-869

When Turbo creates a job from a request, the request's notes become a client note on the job's client and job. The note's author is the request creator, not the person who creates the job. If the notes are blank after trimming, Turbo creates no note.

## Only admins delete abandoned requests

Lines 870-908

Deleting a request needs the admin claim. A request not yet abandoned gives "Only abandoned requests can be deleted.". Deletion also removes its attachment and allocations.
