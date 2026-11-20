<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_requests.go","sha256":"1e21462ebcb83942bd1a8d074f83366e4f32ce9317ab1f0ab51e050a425936e2","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_requests.go

## Project request routes need a signed-in user

Lines 55-98

Asked as: not logged in project request; request page no access; can't reach job request; session expired attachment upload

All project request and job attachment routes require a signed-in staff user. Turbo returns errors from these routes as readable messages. A validation failure shows 'Check the request fields.'

## Only active staff can use project requests

Lines 104-116

Asked as: not active staff member; my account is inactive; can't use project requests; I'm signed out

A person who is not an active staff member is refused with "An active staff member is required." A person who is not signed in sees "Sign in to continue."

## Claim required message for workflow actions

Lines 104-129

Asked as: missing claim permission; why does action fail; don't have the right access; workflow button refused

When an action needs a claim that the person lacks, Turbo refuses it with 'The <claim> claim is required.' An inactive user first gets 'An active staff member is required.' A person who is not signed in gets 'Sign in to continue.'

## Project requests can be switched off

Lines 130-139

Asked as: project requests turned off; can't save project request; request feature disabled; buttons greyed out for requests

Every project request action (save, submit, recall, reject, approve, create) is refused with "Project requests are disabled." when Turbo's project request approval setting is off (utilities.IsJobRequestApprovalRequired). Deleting is different: it checks the setting only when it asks whether the creator can still edit.

## When a project request counts as abandoned

Lines 140-142

Asked as: request too old; request gone stale; when does request expire; old unedited request

A request is abandoned when its last edit is JOB_REQUEST_ABANDONMENT_DAYS days old or older (the value comes from utilities). A request edited exactly that many days ago counts as abandoned. Abandonment lets an admin view and delete the request.

## Who can see a project request

Lines 143-174

Asked as: can't see my request; who can view request; request missing from list; why can't I open request

The creator can always see a request. Anyone with the job claim (procurement) can see every request. A person with the project authorization claim sees only requests where they are the selected approver. An admin sees only abandoned requests. Everyone else sees "You cannot view this request." A missing request gives "Project request not found."

## Submit and recall actions need the right state and creator

Lines 161-174

Asked as: can't submit my request; recall button not working; save returns old request; request not found on submit

Submit and recall are separate steps from saving. A save returns the updated request. An action on a missing request gives 'Project request not found.' A person who cannot view the request gets 'You cannot view this request.'

## Only the creator can edit, and only in some states

Lines 175-190

Asked as: can't edit someone's request; locked request editing; edit button disabled; someone else's request change

Only the creator can edit a project request: others see "Only the creator can edit this request." Editing is allowed in Draft, Recalled or Rejected status. In any other status the message is "Recall this request before editing it."

## Stale version is refused

Lines 202-209

Asked as: request was changed by someone; version mismatch save; someone else edited it; save conflict reload

Save, submit, recall, reject, approve and create are refused with "This request changed. Reload it before continuing." when the version sent is not the request's current version.

## Clearing an attachment's approval

Lines 210-218

Asked as: approval cleared after edit; pdf approval disappeared; approver removed from attachment; why did approval reset

When an edit, submit, recall or rejection clears the authorization PDF's approval, Turbo empties its approver, approved date, review token and rejection details. It resets the max value, request version, checklist version and approval sequence to 0. It sets the status back to Pending.

## Editing a request resets its approval state

Lines 224-247

Asked as: editing wiped my approval; approval reset after change; request back to pending; edit undid approval

A business edit adds 1 to the version, updates the last edited time, and clears the submission token. A Rejected request goes back to Draft and its rejection details are cleared. Any authorization PDF record is reset to Pending, its approval is cleared, and it takes the new project value.

## Project value must be a whole dollar amount

Lines 248-263

Asked as: project value rejected; cents in project value; decimal project amount; value too large

Project value must be a whole CAD dollar amount, zero or more, no larger than 9007199254740991. Otherwise Turbo says "Project value must be a whole CAD dollar amount." A draft can hold 0. Submitting needs a value above 0.

## Only listed fields can be edited

Lines 252-260

Asked as: can't change that field; field locked in request; uneditable field error; field not editable message

A save that includes a field outside the editable business fields is refused with "Field is not editable: " and the field name.

## Notes length limit

Lines 264-275

Asked as: notes too short; notes character limit; request notes error; notes length problem

Request notes must be text. After trimming, they must be 10 to 1000 characters, or blank. Otherwise Turbo says "Notes must be 10 to 1000 characters, or blank." The notes become the client note on the new job.

## Fee basis must be chosen

Lines 276-281

Asked as: fee basis missing; time and materials choice; can't pick fee basis; submit fails no basis

Setting time and materials needs a true or false value, otherwise Turbo says "Select a fee basis." Submitting without a chosen fee basis gives the same message.

## Project status must be Active or Closed

Lines 290-292

Asked as: project status error; status must be active; closed project status; can't pick project status

A request's project status must be Active or Closed, otherwise Turbo says "Select Active or Closed for the project status."

## Approver's limit is checked on save

Lines 296-300

Asked as: approver limit too low; approver can't cover value; save fails approver limit; approver spending limit

When a request has an approver, saving and submitting check that the approver's project limit covers the project value (utilities.ProjectAuthorizationLimit). Creating the job skips this check.

## Parent project rules for a subjob

Lines 301-312

Asked as: subjob parent problem; parent project not valid; subjob different client; top level parent number

A parent must exist, or Turbo says "Select a valid parent project." The parent's number must look like 25-123, with two digits, a dash and three or four digits, or Turbo says "Only a top-level project can be a parent." The request's client must match the parent's, or Turbo says "The subjob must use the parent project's client."

## Linked proposal rules

Lines 313-326

Asked as: proposal link not working; linked proposal wrong; proposal status not allowed; proposal client mismatch

A linked proposal must exist, or Turbo says "Select a valid proposal." Its number must start with P and its client must match the request, or Turbo says "Select a proposal for this client." Its status must be In Progress, Submitted or Awarded, or Turbo says "Select an In Progress, Submitted, or Awarded proposal."

## Required fields to submit

Lines 333-343

Asked as: can't submit missing fields; required field blank; project value zero submit; fill in before submitting

Submit, approve and create need these fields filled: description, location, client, contact, manager, branch, rate sheet, project award date, invoicing information, approver and fn agreement answer. A blank one gives "Required field: " and the field name. A project value of 0 or less gives "Project value must be positive."

## Proposal must be Awarded before submit

Lines 344-352

Asked as: proposal not awarded yet; procurement must award; linked proposal blocks submit; can't submit until awarded

Submitting, approving or creating from a request with a linked proposal needs the proposal to be Awarded. Otherwise Turbo says "Procurement must award the linked proposal before this request can continue."

## Request validation checks client references and contact emails

Lines 353-372

Asked as: client reference check failed; contact email invalid; submit validation errors; invoicing info rejected

Before a request can be submitted or approved, Turbo validates its client, contact and invoicing references (ValidateClientReferences) and its contact emails (ValidateRequestContactEmails). It also runs the project request data checks (ValidateProjectRequestData). The rules for those checks are in other files.

## Authorization PDF required

Lines 373-382

Asked as: no authorization pdf; upload pdf to submit; missing signed form; authorization document needed

Submit, approve and create need an uploaded authorization PDF with a stored hash. Otherwise Turbo says "Upload an authorization PDF."

## Division allocation rules

Lines 385-409

Asked as: division hours decimal; division not active; allocation error on save; pick a division

A submitted request needs at least one division allocation. Hours must be whole numbers, or Turbo says "Allocation hours must be whole numbers." Each division must exist ("Select a valid division.") and be active ("Select an active division."). A save may leave allocations out. If it sends some, the old ones are replaced only when they differ.

## Creating a request needs the time claim

Lines 424-437

Asked as: can't start new request; no time claim; create request button missing; new request not allowed

Starting a new project request needs the time claim, otherwise Turbo says "The time claim is required." The new request starts as a Draft with project status Active and version 1.

## New project request starts as Active draft

Lines 428-437

Asked as: new request defaults; what status new request; who becomes creator; first version of request

A new project request starts as a Draft with project status Active and version 1. Turbo sets the person who saves it as the creator and stamps the last edited time.

## Rate sheet or proposed rates, not both

Lines 461-487

Asked as: rate sheet and new rates; both rates chosen error; proposed rates conflict; draft rates deleted

A request either chooses a rate sheet or proposes new rates. Doing both gives "Choose a rate sheet or propose new rates, not both." Choosing an active sheet deletes any draft of proposed rates the request had.

## Allocations replace the old set only when they differ

Lines 488-501

Asked as: allocations overwritten; allocations not saving; changed allocations replaced; empty allocation list

A save that sends allocations changes them only if they differ from the saved ones. Turbo then deletes the old allocations and saves the new ones. A save without allocations leaves them alone, and an empty list is allowed on save.

## Save with no changes does nothing

Lines 503-511

Asked as: saved but nothing changed; version didn't update; save no change; last edited unchanged

When an existing request is saved with no changes to its fields, allocations or proposed rates, Turbo saves nothing. The version and last edited time stay the same.

## Who can reject a request, and the reason

Lines 572-702

Asked as: who can reject request; approver can't reject; reject reason required; job holder reject

The selected approver, if they hold the project authorization claim, can reject in Awaiting approval or Awaiting creation. A job holder (procurement) can reject only in Awaiting creation. Others see "Only the selected approver, or a job holder once the request is approved, can reject this request." A reason is required: "Enter a rejection reason."

## Approval needs a completed checklist and a current review

Lines 594-619

Asked as: checklist answers wrong; checklist changed reload; approval checklist error; confirm checklist item

The approver must send every checklist answer as true for the current approval checklist version. Turbo says "No approval checklist is saved.", "Checklist version N is not supported.", or "The checklist changed. Reload the request and review it again." when the version is missing, wrong or old. A wrong answer count gives "Use only the N checklist answers for version N.", and an unconfirmed item gives "Confirm: " and its key. If the review token or document hash is out of date, Turbo says "The reviewed data changed. Reload and review it again."

## Request save and action use the request version

Lines 628-638

Asked as: request version mismatch; reload before continuing; action rejected stale copy; old version of request

Each action also carries the request version. If it does not match the saved version, Turbo says 'This request changed. Reload it before continuing.' This check runs before any action.

## Submitting a request

Lines 646-673

Asked as: send request for approval; submit rejected request; submit sends for approval; rejected needs editing first

Submitting needs the time claim and an editable request. A Rejected request must be edited first ("Edit the rejected request before submission."). Submitting sets the status to Awaiting approval, records the submitted time, clears earlier approval and issues a new submission token, even if nothing changed.

## Recalling a request

Lines 674-683

Asked as: take back request; cancel submitted request; withdraw request from approval; only creator recalls

Only the creator can recall ("Only the creator can recall this request."). Recall works only in Awaiting approval or Awaiting creation, otherwise "This action requires a submitted or approved request." It sets the status to Recalled, ends the submission token and clears approvals, including the PDF's and proposed rates'.

## Rejection reason is stored and approval is withdrawn

Lines 695-717

Asked as: rejection reason saved; pdf marked rejected; withdrawn approval notice; rejected after approval

A rejection saves the rejection time, the rejector and the trimmed reason on the request. It also sets the status of the authorization PDF to Rejected with the same details. An approver who rejects after approval counts as withdrawing the approval in the notice.

## Rejecting ends the submission

Lines 703-717

Asked as: reject ends submission; rejected request what next; request rejected by approver; reject clears approvals

Reject sets the status to Rejected and records who, when and why. It also clears the submission token and approvals, and marks the PDF record Rejected. A rejected request can then be edited by its creator, which returns it to Draft.

## Approving a request

Lines 718-732

Asked as: how to approve request; approve button error; only submitted can approve; approve no submission token

Only a submitted request (Awaiting approval, not rejected) can be approved, or Turbo says "Only a submitted request can be approved." Only the selected approver can approve ("Only the selected approver can approve this request."). The approver's project limit must cover the value. A request with no submission token gives "This request has no submission token. Recall and submit it again."

## Approval requires selected approver, limit and submission token

Lines 718-734

Asked as: who can approve; approval blocked limit; approver limit not enough; approve needs submission

Only a request in Awaiting approval with no rejection can be approved, or Turbo says 'Only a submitted request can be approved.' A request with no submission token gets 'This request has no submission token. Recall and submit it again.' Only the selected approver can approve. Their project limit must cover the project value (ProjectAuthorizationLimit).

## Approval fails if reviewed data changed

Lines 735-748

Asked as: data changed before approval; review token mismatch; reload and review again; document hash changed

On approval Turbo compares the review token and document hash that the approver reviewed with the current ones. If they differ, it says 'The reviewed data changed. Reload and review it again.' Turbo runs the full completeness check on approval, including the approver's authority.

## What approval records

Lines 749-766

Asked as: what approval records; approval saves what; approval moves to creation; approved rates after approval

Approval marks the PDF record Approved with the approver, time, project value, approver's limit, request version and checklist version. It also approves any proposed rates, saves a fresh review token and moves the request to Awaiting creation.

## Creating the job from a request

Lines 767-784

Asked as: create job from request; who can create job; not awaiting creation; creation checklist fails

Only a job holder (the job claim) can create, and only when the request is Awaiting creation ("This request is not awaiting creation."). A request with no submission token gives "This request has no submission token. Recall and submit it again." The creator must complete the creation checklist, with the same checklist checks as approval.

## Stale or incomplete approval blocks creation

Lines 787-803

Asked as: approval stale creating job; saved approval incomplete; can't create job approval; invalid saved checklist

Creation is refused with "The approval is stale. Recall the request for a new approval." when the PDF is not Approved, the version or value changed since approval, or the review token or document hash does not match. It is refused with "The saved approval is incomplete. Recall the request for a new approval." when the approver, approval time or limit is missing or below the value, or when a rejection exists. A bad saved approval checklist gives "The saved checklist is invalid." Creation does not recheck the approver's current limit.

## Proposed rates at creation

Lines 804-810

Asked as: proposed rates on job; rename rate sheet creation; requester rate sheet name; rates not approved for job

When creating, Turbo requires the proposed rates to be approved. The job creator may rename the rate sheet, and an empty name keeps the requester's name. Approved proposed rates become an active rate sheet that the job uses.

## Saving the attachment and request after an action

Lines 825-845

Asked as: action rolled back; notice failed after submit; transaction failed submit; saving pdf and request

After submit, recall, reject or approve, Turbo saves the authorization PDF record and then the request in the same transaction. If a notice fails, the whole action rolls back. Notices are sent only after the transaction succeeds.

## Notifications after submit, approve and reject

Lines 833-839

Asked as: no email after submit; notice after approval; reject sends notification; who gets notified

After submit, approve or reject, Turbo queues a notice. A reject of an Awaiting creation request by its own approver counts as a withdrawn approval. A job holder's reject does not.

## Approved proposed rates become an active rate sheet

Lines 852-857

Asked as: rates become rate sheet; job rate sheet activated; opening rate sheet period; approved rates job sheet

At job creation, Turbo activates the request's approved proposed rates as an active rate sheet. The job uses it and gets its opening rate sheet period on it.

## What job creation copies and sets

Lines 858-878

Asked as: what copies to job; job fields from request; job status after creation; fn agreement on job

The new job copies the request's business fields except the approver. Its status is the request's project status, and fn agreement is true only when the answer is "Yes". Its authorizing document is PA. Turbo records the creator, created and submitted times, who created the job, and the completion time.

## Created job gets creation details and PA document

Lines 862-878

Asked as: job created details; authorizing document PA; job first nation agreement; who created the job

The new job takes the request's business fields, except the approver. Its status is the request's project status. Its first-nation agreement is true only when the answer is 'Yes'. Its authorizing document is set to PA. Turbo records the request's creator, created and submitted dates, the person who created the job, and the completion time.

## Allocations, PDF and note move to the job

Lines 883-923

Asked as: allocations move to job; request deleted after job; client note from request; pdf moved to job

Request allocations become job time allocations, and the request's own are deleted. The authorization PDF moves to the job as its initial approval, with approval sequence 1. Non-blank notes become the job's first client note, written by the request creator. The request is then deleted.

## Authorization PDF becomes the job's first approval

Lines 903-923

Asked as: pdf first approval job; approval sequence one; attachment moves to job; request removed after creation

When the job is created, the authorization PDF moves from the request to the job. It is marked as the initial approval with approval sequence 1. Turbo then deletes the request.

## Merge references transfer to the job

Lines 917-919

Asked as: merge references move; merge references to job; references transferred job; merge reference missing

When the job is created, Turbo moves the request's merge references to the new job and its note (transferRequestMergeReferences). The rules for that step are in another file.

## Request notes become the job's first client note

Lines 926-947

Asked as: request notes become note; first client note job; notes copied to client; blank notes no note

If the request notes are not blank after trimming, Turbo makes them the job's first client note. The note belongs to the job's client and keeps the request creator as its author. Blank notes make no note.

## Deleting a request

Lines 948-994

Asked as: delete my request; can't delete request; only abandoned deleted; admin delete old request

The creator can delete a request they can still edit (Draft, Recalled or Rejected). Anyone else needs the admin claim, and the request must be abandoned, otherwise Turbo says "Only abandoned requests can be deleted." Deleting removes the PDF record and the allocations as well.
