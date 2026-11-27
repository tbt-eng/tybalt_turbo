<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_requests.go","sha256":"71f60d032387998249ed615664a1c9e5ccb72df860be0897ecb847ecca529d40","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_requests.go

## Who can reach project request routes

Lines 55-100

Asked as: not signed in to requests; request page login; check the request fields; request page won't open

Every project request route requires a signed-in user. A validation failure on a request returns the message "Check the request fields.". A request refused by a Turbo workflow check returns that check's message.

## Project requests can be switched off

Lines 119-128

Asked as: requests turned off; can't save request; project requests disabled; can't delete my request

Saving, submitting, recalling, rejecting, approving, creating and editing a project request all fail with "Project requests are disabled." when the company setting that requires project request approval is off. Deleting a request is also refused with this message when the creator deletes their own editable request; an admin deleting an abandoned request is not blocked by it.

## Abandoned request definition

Lines 129-131

Asked as: when is a request abandoned; how many days abandoned; old unedited request; abandonment cutoff

A request is abandoned when its last edit is at or before the cutoff of JOB_REQUEST_ABANDONMENT_DAYS days ago. A request edited exactly at the cutoff counts as abandoned.

## Who can view a project request

Lines 132-160

Asked as: can't see request; cannot view this request; why can't I open request; request not visible to me

The creator can always view their request. Anyone with the job claim (procurement) can view any request. An approver (project_authorization claim) can view only requests where they are the selected approver. An admin can view only abandoned requests. Anyone else gets "You cannot view this request."

## Request owner can view abandoned or any own request

Lines 132-146

Asked as: can't open my request; open abandoned request; admin can't open request; request won't open for me

A request's creator can always open it. Staff with the job claim can open any request. The selected approver (the project_authorization claim) can open requests assigned to them. An admin can open a request only when it is abandoned.

## Only the creator edits, and only in some states

Lines 161-173

Asked as: can't edit request; edit button greyed out; why can't I change request; recall before editing

Only the creator can edit a request, otherwise Turbo says "Only the creator can edit this request." Editing is allowed when the status is Draft, Recalled or Rejected. In any other status Turbo says "Recall this request before editing it."

## Stale request version is refused

Lines 185-192

Asked as: request changed reload; version mismatch on save; someone else edited request; save conflict request

Saves and actions carry the request's version number. If it differs from the saved version, Turbo refuses with "This request changed. Reload it before continuing."

## Editing a request resets its approval

Lines 193-230

Asked as: edit resets approval; approval went pending; rejected back to draft; changing request lost approval

A saved edit raises the version, updates last edited, clears the submission token and any rejection, and turns a Rejected request back into Draft. It also resets the authorization PDF's approval to Pending and copies the project value to it.

## Project value must be a whole CAD amount

Lines 231-246

Asked as: project value has cents; value must be whole dollars; decimal project value error; zero value submit

Project value must be a number with no cents, zero or more, and at most 9007199254740991. Otherwise Turbo says "Project value must be a whole CAD dollar amount." Zero can be saved as a draft, but submitting needs a positive value.

## Only listed fields can be edited on a request

Lines 235-243

Asked as: field is not editable; can't change that field; field locked on request; field refused on save

A save that includes any field outside the allowed business fields is refused with "Field is not editable: <field>".

## Request notes length limit

Lines 247-258

Asked as: notes too short; notes character limit; notes must be text; notes length error

Notes are trimmed. They must be text, and 10 to 1000 characters, or blank. Otherwise Turbo says "Notes must be text." or "Notes must be 10 to 1000 characters, or blank."

## Fee basis must be true or false

Lines 259-264

Asked as: select a fee basis; time and materials choice; fee basis yes or no; fee basis error

The fee basis (time_and_materials) must be a yes/no value, otherwise Turbo says "Select a fee basis." Setting it records that a fee basis was chosen, and submitting needs that choice.

## Draft checks on project status, approver and parent

Lines 272-295

Asked as: parent project error; subjob client mismatch; select active or closed; approver limit too low

Project status must be Active or Closed, otherwise Turbo says "Select Active or Closed for the project status." When an approver is chosen, Turbo checks that their project limit covers the project value. A parent must be an existing top-level project (number like 25-123), the same client as the request. Messages: "Select a valid parent project.", "Only a top-level project can be a parent.", "The subjob must use the parent project's client."

## Linked proposal must belong to the client and be open

Lines 296-309

Asked as: proposal not valid; select proposal for client; proposal status wrong; linked proposal error

A linked proposal must exist, have a number starting with P and have the request's client, otherwise Turbo says "Select a valid proposal." or "Select a proposal for this client." Its status must be In Progress, Submitted or Awarded, otherwise Turbo says "Proposal X is Y. Select an In Progress, Submitted, or Awarded proposal."

## Job creation needs an Awarded proposal

Lines 313-338

Asked as: can't create job; proposal not awarded; award proposal before job; job creation blocked

A request can be submitted and approved while its proposal is In Progress or Submitted, but creating the job requires the proposal to be Awarded. Otherwise Turbo says "Proposal X is Y. Award it before you create the job." with code proposal_not_awarded. Turbo never awards the proposal itself.

## Required fields to submit or approve a request

Lines 339-387

Asked as: can't submit request; required field missing; upload authorization pdf; submit request fails

Submit and approve need description, location, client, contact, manager, branch, rate sheet, project award date, approver and fn agreement answer, otherwise "Required field: <field>". Also needed: a fee basis ("Select a fee basis."), a positive project value ("Project value must be positive."), valid client references and contact emails, valid allocations, and an uploaded authorization PDF ("Upload an authorization PDF.").

## Allocation rules on a request

Lines 388-412

Asked as: allocation hours decimal; invalid division; inactive division; allocation hours error

Division hours must be whole numbers ("Allocation hours must be whole numbers."). Each division must exist ("Select a valid division.") and be active ("Select an active division."). Invalid hours give "Invalid allocation hours." Allocations are required to submit, but optional when saving a draft.

## Creating a request needs the time claim and starts as a Draft

Lines 427-440

Asked as: can't create request; new request permission; need time access to create; new request status draft

Creating a project request requires the time claim. A new request starts with status Draft, project status Active and version 1, and the signed-in user is its creator.

## Choose a rate sheet or propose rates, not both

Lines 464-490

Asked as: rate sheet and proposed rates; both rates chosen; proposed rates error; choose one rate option

If a save proposes new rates and also chooses a different rate sheet, Turbo says "Choose a rate sheet or propose new rates, not both." Choosing an active sheet instead drops the request's own draft rates.

## Saving with no changes does nothing

Lines 505-514

Asked as: save does nothing; version not changing; nothing saved no changes; unchanged request save

Updating an existing request with no changes to its fields, allocations or proposed rates saves nothing, and the version does not change.

## Who can reject a request

Lines 575-593

Asked as: who can reject request; reject button missing; can't reject request; job holder reject

The selected approver can reject while the request is Awaiting approval or Awaiting creation. A job holder (job claim) can reject only once it is Awaiting creation. Others get "Only the selected approver, or a job holder once the request is approved, can reject this request."

## Approval checklist must be fully confirmed

Lines 595-622

Asked as: checklist not confirmed; approval checklist error; checklist version wrong; confirm checklist question

Approving or creating requires every checklist question to be answered true and the checklist version to be current. Otherwise: "No <scope> checklist is saved.", "Checklist version N is not supported.", "The checklist changed. Reload the request and review it again.", "Use only the N checklist answers for version V." or "Confirm: <question>".

## Submitting a request

Lines 649-676

Asked as: send request for approval; submit request fails; rejected request submit; submission needs time access

Submitting needs the time claim, an editable request and complete data, with the creator as the user. A Rejected request must be edited first ("Edit the rejected request before submission."). Submit sets status Awaiting approval, records the time, clears earlier approvals and issues a new submission token. It does not change version or last edited.

## Recalling a request

Lines 677-717

Asked as: pull back request; withdraw submitted request; recall not allowed; cancel submission

Only the creator can recall ("Only the creator can recall this request."), and only an Awaiting approval or Awaiting creation request, otherwise "This action requires a submitted or approved request." Recall sets status Recalled, ends the submission token and clears approvals.

## Rejecting a request needs a reason

Lines 687-717

Asked as: reject without reason; enter rejection reason; rejection reason required; withdraw approval

Rejection needs a non-blank reason, otherwise "Enter a rejection reason." It sets status Rejected, records rejector, time and reason, ends the submission token and marks the authorization PDF Rejected. An approver who rejects after approval counts as withdrawing it.

## Approving a request

Lines 718-766

Asked as: approve request fails; approver cannot approve; no submission token; approve button not working

Only the selected approver can approve, and only an Awaiting approval request that is not rejected ("Only a submitted request can be approved.", "Only the selected approver can approve this request."). A request with no submission token gives "This request has no submission token. Recall and submit it again." The approver's project limit must cover the value. Approval moves the request to Awaiting creation and records the approver, time, limit and checklist version.

## Approval fails if reviewed data changed

Lines 742-748

Asked as: reviewed data changed; review token mismatch; document changed before approval; reload and review again

If the review token or document hash the approver sends does not match the current request and PDF, Turbo says "The reviewed data changed. Reload and review it again."

## Creating the job from an approved request

Lines 767-790

Asked as: create job from request; job creation fails; not awaiting creation; make job button

Creating the job needs the job claim and an Awaiting creation request with a submission token, plus the creation checklist. The authority check on the approver's limit is skipped at this step. Otherwise: "This request is not awaiting creation."

## Stale or incomplete approval blocks creation

Lines 796-825

Asked as: approval is stale; approval incomplete; saved checklist invalid; page loaded too old

Creation is refused with "The approval is stale. Recall the request for a new approval." if the approval is not Approved, or the version, project value or document hash differs. If only the page is old, Turbo says "The request changed after this page loaded. Reload the page, then create the job." An approval with a different approver, no approved time, a limit below the project value or any rejection gives "The saved approval is incomplete. Recall the request for a new approval." An invalid approval checklist gives "The saved checklist is invalid."

## What job creation copies and deletes

Lines 874-946

Asked as: what happens creating job; request deleted after job; fields copied to job; fn agreement job

The new job copies the request's business fields except the approver. Its status is the request's project status, and fn agreement is true only when the answer is "Yes". Authorizing document is PA. Allocations move to the job. The authorization PDF moves to the job as its first approval. The request's notes become the job's first client note, kept under the creator, and the request is deleted.

## Blank notes make no client note

Lines 951-955

Asked as: no client note created; blank notes job; empty notes no note; notes missing on job

If the request's notes are blank after trimming, job creation makes no client note.

## Deleting a project request

Lines 977-1008

Asked as: delete project request; only admin can delete; delete abandoned request; can't delete request

The creator can delete a request they can still edit. For any other request only an admin can delete it, and only if it is abandoned, otherwise "Only abandoned requests can be deleted." Deleting also deletes its authorization PDF and allocations.
