<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_requests.go","sha256":"8e17c63a4a723bb9012e15207a8e66af6c5bd8e505ff92d84cdd5491dd609db3","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_requests.go

## Creator can rename proposed rates when the job is created

Lines 47-49

Asked as: rename rate sheet; rate sheet name lost; proposed rates wrong name; name my rates when creating job

A job creator can send a rate sheet name (rate_sheet_name) when creating the job from a request. Turbo then renames the proposed rates to it (renameProposedRates). An empty name keeps the name the requester gave.

## Project request and attachment routes need a signed-in user

Lines 53-55

Asked as: cannot open job requests; not logged in attachments; sign in to see requests; attachments page access denied

Every project request route (/api/job_requests), the authorization approvers list and every job attachment route (/api/job_attachments) require a signed-in user (users auth). A person who is not signed in cannot use them.

## Request actions are limited to five

Lines 65-67

Asked as: which request actions exist; can't find request button; submit recall reject buttons; request action limits

The request routes allow submit, recall, reject, approve and create. Each is a POST on a single request.

## Attachment routes: approve and reject actions exist

Lines 70-79

Asked as: approve attachment PDF; reject authorization attachment; attachment approve button missing; attachment versus request approval

Job attachments (the authorization PDF record) have their own approve and reject actions. Both are separate from the project request's own approve and reject actions. Other rules describe the attachment's checks.

## Action errors show a message to the person

Lines 82-93

Asked as: error message not showing; why did my save fail; vague error on request; what does that error mean

Workflow errors are returned with their message. Field validation errors return the message "Check the request fields."

## Only active staff can use project requests

Lines 99-111

Asked as: not active staff error; cannot submit request as contractor; inactive account request blocked; sign in to continue message

Turbo refuses a project request action from a person who is not active staff. It shows "An active staff member is required." A person who is not signed in sees "Sign in to continue."

## Claim required messages

Lines 112-124

Asked as: claim required message; missing time claim; need job claim to create; why do I lack permission

When an action needs a claim that the person lacks, Turbo shows "The <claim> claim is required." Saving a new request needs the time claim. Creating a job from a request needs the job claim.

## Project requests can be turned off

Lines 125-134

Asked as: project requests disabled; cannot save request turned off; request feature switched off; requests not required right now

When project requests are not required (utilities.IsJobRequestApprovalRequired is false), Turbo refuses to save, submit, recall, reject, approve or create a project request. It shows "Project requests are disabled." Deleting a request does not run this check.

## Abandoned request after set days without edits

Lines 135-137

Asked as: request abandoned after no edits; old draft request expired; how long until request abandoned; stale draft days limit

A project request is abandoned when its last edit is not newer than JOB_REQUEST_ABANDONMENT_DAYS days ago. A request edited exactly that long ago counts as abandoned.

## Who can view a project request

Lines 138-169

Asked as: cannot see request; who can view request; request hidden from me; view someone else's request

The creator can view their request. Anyone with the job claim can view any request. The selected approver with the project authorization claim can view the request. An admin can view only an abandoned request. Anyone else sees "You cannot view this request."

## Request not found

Lines 156-160

Asked as: request missing or gone; project request not found; wrong request id link; deleted request cannot open

When the request id does not exist, Turbo shows "Project request not found."

## Who can edit a project request and when

Lines 170-185

Asked as: cannot edit request; request locked for editing; why can't I change request; edit button greyed out

Only the creator can edit a request, and Turbo shows "Only the creator can edit this request." to others. The request must be Draft, Recalled or Rejected. In any other status Turbo shows "Recall this request before editing it."

## Stale request version is refused

Lines 197-204

Asked as: request changed reload; someone else edited request; version mismatch on save; conflicting edits request

When the version a person sends differs from the saved version, Turbo refuses with "This request changed. Reload it before continuing." This applies to saving and to every action.

## Editing a request resets approval and rejection

Lines 205-242

Asked as: edit clears approval; changing request resets approval; rejected request back to draft; editing wipes attachment approval

Each business edit adds 1 to the request version, updates its last edited time and clears its submission token. A Rejected request goes back to Draft and its rejection details are cleared. The authorization attachment's approval is cleared and its status returns to Pending.

## Project value must be a whole dollar amount

Lines 243-258

Asked as: project value decimal error; cents in project value; value too large; whole dollars only value

The project value on a request must be a whole, non-negative number of CAD dollars, no larger than 9007199254740991. Otherwise Turbo shows "Project value must be a whole CAD dollar amount."

## Only listed fields are editable on a request

Lines 247-254

Asked as: field not editable; cannot change this field; which request fields editable; read only request field

Turbo accepts only the request's business fields. Any other field is refused with "Field is not editable: <field>".

## Request notes length

Lines 259-270

Asked as: notes too short; notes length limit; request notes character count; notes not text error

Turbo trims the notes of a request. The notes must be blank or 10 to 1000 characters long. Otherwise Turbo shows "Notes must be 10 to 1000 characters, or blank." Notes that are not text show "Notes must be text."

## Fee basis must be chosen

Lines 271-276

Asked as: fee basis missing; time and materials choice; select fee basis; cannot submit without fee basis

Setting time and materials on a request requires true or false, or Turbo shows "Select a fee basis." The request counts as having a fee basis only after this choice. Submitting without it shows the same message.

## Project status must be Active or Closed

Lines 285-287

Asked as: project status options; status must be active closed; project status error saving

Saving a request requires project status Active or Closed. Otherwise Turbo shows "Select Active or Closed for the project status."

## Approver's project limit checked on save

Lines 291-295

Asked as: approver limit too low; approver authorization limit; approver cannot cover value; approver spending limit exceeded

When a request has an approver, saving or submitting checks that the approver's project authorization limit covers the project value (utilities.ProjectAuthorizationLimit). Creating a job skips this check.

## Parent project rules for a subjob

Lines 296-307

Asked as: parent project invalid; subjob client mismatch; top level project parent; job number format parent

A parent must be an existing job whose number looks like 25-123 or 25-1234, a top-level project. Otherwise Turbo shows "Select a valid parent project." or "Only a top-level project can be a parent." The request's client must match the parent's client, or Turbo shows "The subjob must use the parent project's client."

## Linked proposal rules

Lines 308-321

Asked as: linked proposal invalid; proposal wrong client; proposal status not allowed; choose proposal for client

A linked proposal must be a job whose number starts with P and that has the same client. Its status must be In Progress, Submitted or Awarded. Otherwise Turbo shows "Select a proposal for this client." or "Select an In Progress, Submitted, or Awarded proposal."

## Required fields before submit or approval

Lines 324-338

Asked as: required field missing submit; cannot submit incomplete request; project value must be positive; what fields to fill in

A complete request needs description, location, client, contact, manager, branch, rate sheet, award date, invoicing information, approver and the fn agreement answer. A missing one shows "Required field: <field>". The project value must be above zero, or Turbo shows "Project value must be positive."

## Linked proposal must be Awarded to continue

Lines 339-347

Asked as: proposal not awarded yet; linked proposal blocks submit; set proposal awarded; cannot continue with proposal

A request with a linked proposal can be submitted, approved or used to create a job only when the proposal is Awarded. Otherwise Turbo shows "Set the linked proposal to Awarded before continuing."

## Allocation hours rules

Lines 354-404

Asked as: division allocation hours; allocation decimal hours; invalid division on request; inactive division selected

Division allocation hours on a request must be whole numbers, or Turbo shows "Allocation hours must be whole numbers." Each division must exist and be active, or Turbo shows "Select a valid division." or "Select an active division." A complete request must have allocations.

## Authorization PDF required

Lines 368-378

Asked as: missing authorization PDF; upload authorization document; no PDF attached request; cannot submit without pdf

A complete request needs an uploaded authorization PDF with its hash. Otherwise Turbo shows "Upload an authorization PDF."

## Unreadable save body is refused

Lines 407-410

Asked as: invalid request data; save request failed body; request save gives bad data

When Turbo cannot read the body of a request to create or edit a project request, it refuses with 'Invalid request data.' No change is made.

## Only users with the time claim can start a project request

Lines 418-422

Asked as: cannot start new request; new request time claim; no permission to create request; project requests turned off creating

Creating a new project request requires the time claim. A person without it gets the message "The time claim is required." Turbo also refuses it when project requests are disabled.

## New request defaults

Lines 419-432

Asked as: default request status; new request starting values; what does new request start as

A new request starts as a Draft with project status Active and version 1. Its creator is the person who saves it.

## Editing someone else's request: view check comes before edit check

Lines 433-444

Asked as: edit someone else's request; check order saving request; view then edit then version

When a person saves changes to an existing request, Turbo first checks that they can view it. Then it checks that they can edit it. Last it checks that the version they sent matches the saved one.

## Rate sheet or proposed rates, not both

Lines 456-482

Asked as: rate sheet and proposed rates; both rate sheet and new rates; choose one rate option

A request can choose an active rate sheet or propose new rates, but not both. Choosing another rate sheet drops the request's draft proposed rates. Otherwise Turbo shows "Choose a rate sheet or propose new rates, not both."

## Proposed rates select the request's own draft sheet

Lines 465-474

Asked as: proposed rates draft sheet; rate sheet conflicts proposed; choose rate sheet not both

When a request has proposed rates and a draft, Turbo sets the request's rate sheet (rate_sheet) to that draft. Turbo gives the message "Choose a rate sheet or propose new rates, not both." when the request also chooses a different sheet.

## Chosen rate sheet replaces a draft of proposed rates

Lines 475-482

Asked as: chosen rate sheet deletes draft; draft rates disappeared; switching rate sheet removes draft

When a request chooses an active rate sheet other than its own proposed-rates draft, Turbo deletes the draft after saving. The check on the chosen sheet runs first (checkChosenRateSheet).

## Allocations are replaced as a whole when they change

Lines 483-496

Asked as: allocations not saving; allocations replaced changed; old allocations removed; saved allocations stay same

When the allocations sent differ from the saved ones, Turbo deletes the old request allocations and saves the new ones. When no allocations are sent, the saved ones stay unchanged.

## Saving with no changes does nothing

Lines 497-506

Asked as: save with no changes; version unchanged saving; nothing changed on save

Saving an existing request with no changes to its fields, allocations or proposed rates does not change the request or its version.

## Saving a new request with proposed rates saves it twice

Lines 507-514

Asked as: new request saved twice; proposed rates new request save; why request saved two times

When a new request has proposed rates, Turbo saves the request first, then the proposed-rates draft, and then saves the request again so that it refers to the draft.

## Editing an existing request counts as an edit

Lines 515-522

Asked as: editing counts as edit; edit clears submission token; edit clears rejection

Saving a changed existing request goes through the edit step (markRequestEdited). It raises the version and the last-edited time, and it clears the submission token and any rejection.

## Request checklist depends on request status

Lines 559-565

Asked as: checklist for awaiting creation; which checklist applies; approval checklist status

A request in the status Awaiting creation uses the creation checklist. In any other status the approval checklist applies.

## Who can reject a request

Lines 567-585

Asked as: who can reject request; cannot reject request; job holder reject rights; only approver can reject

The selected approver with the project authorization claim can reject while the request is Awaiting approval or Awaiting creation. A job holder can reject only once it is Awaiting creation. Anyone else sees "Only the selected approver, or a job holder once the request is approved, can reject this request."

## Checklist must be fully confirmed

Lines 587-614

Asked as: checklist not confirmed; confirm checklist questions; outdated checklist message; no checklist saved

Approval and creation need the current checklist version, with every question answered true and no extra answers. A missing checklist shows "No <scope> checklist is saved." An unsupported version is refused, and an outdated one shows "The checklist changed. Reload the request and review it again." An unconfirmed question shows "Confirm: <question>".

## Unreadable action body is refused

Lines 617-620

Asked as: invalid action data; action request body unreadable; submit approve fails bad data

When Turbo cannot read the body of a submit, recall, reject, approve or create action on a project request, it refuses with 'Invalid action data.'

## Every request action checks the feature and visibility first

Lines 624-633

Asked as: action fails feature check; cannot act on request; action blocked visibility

Before submit, recall, reject, approve or create, Turbo checks that project requests are enabled and that the person can view the request. It then checks that the version sent matches the saved version. A failure on any of these stops the action.

## Request detail fields cannot be changed by a stale page

Lines 629-633

Asked as: stale page action failed; old page approve request; request changed before approving

A save or an action based on an old version of a request is refused with "This request changed. Reload it before continuing." This applies to every submit, recall, reject, approve and create action, which compare the version sent with the saved version.

## Submitting a project request

Lines 641-665

Asked as: how to submit request; submit button not working; send request for approval; submitting resets attachment

Submitting needs the time claim, and only the creator can submit an editable request. A Rejected request shows "Edit the rejected request before submission." Submitting sets the status to Awaiting approval, records the submit time and issues a new submission token. It resets the attachment's approval.

## Rejected request needs editing before resubmission

Lines 648-650

Asked as: cannot submit rejected request; resubmit rejected request; rejected request resubmission

Submitting a request in the Rejected status is refused with "Edit the rejected request before submission." Editing a rejected request returns it to Draft.

## Recall a submitted request

Lines 666-706

Asked as: withdraw submitted request; take back request; recall request for changes; undo submit

Only the creator can recall a request, and only when it is Awaiting approval or Awaiting creation. Otherwise Turbo shows "This action requires a submitted or approved request." or "Only the creator can recall this request." Recall sets the status to Recalled, clears the submission token and the attachment's approval.

## Recall and reject only apply to submitted or approved requests

Lines 667-669

Asked as: recall reject not allowed; request status not submitted; action needs submitted request

Recall and reject work only when the request is Awaiting approval or Awaiting creation. Otherwise Turbo says "This action requires a submitted or approved request." Both actions end the submission by clearing its token.

## Rejection needs a reason

Lines 686-706

Asked as: reject without reason; rejection reason required; send request back with reason

Rejecting needs a reason that is not blank. Otherwise Turbo shows "Enter a rejection reason." Rejection sets the status to Rejected, records who and when, clears the submission token and marks the attachment Rejected.

## Approving a project request

Lines 707-719

Asked as: approve button not working; who can approve request; no submission token approve; only approver approves

Only the selected approver can approve, and only a submitted request that is Awaiting approval and not rejected. Turbo shows "Only a submitted request can be approved." or "Only the selected approver can approve this request." A request without a submission token shows "This request has no submission token. Recall and submit it again."

## Approval checks and result

Lines 720-746

Asked as: approval fails review changed; reviewed data changed; approval result status; approval limit check approve

Approval needs the approver's project limit to cover the value and a checklist. The review token and document hash must match the current request, or Turbo shows "The reviewed data changed. Reload and review it again." Approval marks the attachment Approved with the limit, version and checklist version, and sets the request to Awaiting creation.

## Creating a job from an approved request

Lines 747-766

Asked as: create job from request; request not awaiting creation; create job button fails; job creation no token

Creating needs the job claim and a request that is Awaiting creation, or Turbo shows "This request is not awaiting creation." A missing submission token shows "This request has no submission token. Recall and submit it again." Authority is not rechecked at creation.

## Stale or incomplete approval blocks creation

Lines 767-783

Asked as: approval stale creation blocked; saved approval incomplete; cannot create job approval; checklist invalid creating job

Creation is refused with "The approval is stale. Recall the request for a new approval." when the approval, its version, value or review token no longer match. It is refused with "The saved approval is incomplete. Recall the request for a new approval." when the approver, approved time or approval limit fails. The limit must cover the project value. It shows "The saved checklist is invalid." when the approval checklist version does not exist.

## What job creation copies and sets

Lines 784-856

Asked as: what copies to job; job fields from request; FN agreement on job; authorizing document PA

Creating a job copies the request's business fields except the approver. The job status is the request's project status. The job's FN agreement is true only when the answer is Yes, and its authorizing document is PA. The creator's name, request dates and the creating person are recorded. Approved proposed rates become an active rate sheet, which the creator can rename.

## Rejection notice marks a withdrawn approval

Lines 810-816

Asked as: approval withdrawn notice; rejection notice approver; notice when approver rejects

When the selected approver rejects a request that is Awaiting creation, the notice is marked as an approval withdrawn. A job holder's rejection is not.

## Job creation moves allocations, attachment and notes

Lines 860-900

Asked as: allocations move to job; attachment moves to job; request deleted after job created; first approval sequence

Creating a job moves the request's division allocations to the job and deletes them from the request. The authorization attachment moves to the job as its first approval, with sequence 1. Non-blank notes become the job's first client note, by the request creator. The request is then deleted.

## Request notes become a client note on the new job

Lines 903-924

Asked as: notes become client note; request notes missing on job; client note from request

When Turbo creates a job from a request, it makes the request's notes the job's first client note, with the request's creator as the author. Blank notes make no note.

## Deleting a project request

Lines 925-966

Asked as: delete project request; cannot delete request; only abandoned requests delete; remove old request

The creator can delete a request they can still edit. Anyone else needs the admin claim, and the request must be abandoned. Otherwise Turbo shows "Only abandoned requests can be deleted." Deleting also removes its attachment and allocations.
