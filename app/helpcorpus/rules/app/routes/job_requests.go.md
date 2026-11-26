<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_requests.go","sha256":"dedc8c801b15052a132d242c6800b3f9678de3892c3ac5624d1cb09f0d65d99b","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_requests.go

## Request routes need sign-in and name the project request endpoints

Lines 55-100

Asked as: check the request fields; request form error; request routes login; attachment upload error

All project request and job attachment routes require a signed-in user. Validation errors from these routes show "Check the request fields."

## Only active staff can use project request actions

Lines 106-118

Asked as: not active staff; inactive account request; sign in to continue; can't use project request

Turbo refuses a project request action with "An active staff member is required." when the person is not active staff. A person who is not signed in sees "Sign in to continue."

## Claim checks name the missing claim

Lines 119-131

Asked as: missing claim permission; claim is required; don't have permission; lacking access role

When an action needs a claim the person lacks, Turbo says "The <claim> claim is required." The person must also be active staff.

## Project requests can be disabled

Lines 132-141

Asked as: project requests disabled; can't save project request; requests turned off; request actions blocked

When Turbo does not require project request approval (IsJobRequestApprovalRequired), saving, submitting, approving, creating and deleting a project request are refused with "Project requests are disabled.".

## A request is abandoned after a set number of days without edits

Lines 142-144

Asked as: request abandoned; old request expired; unedited request days; stale request closed

A project request is abandoned when its last edit (last_edited) is the abandonment limit (JOB_REQUEST_ABANDONMENT_DAYS) or more days old. A request edited exactly that long ago counts as abandoned.

## Who can view a project request

Lines 145-176

Asked as: who can see request; cannot view this request; project request not found; why can't open request

The creator can view a request. A job holder (the job claim) can view any request. A project authorization holder (the project_authorization claim) can view a request only when selected as its approver. An admin (the admin claim) can view only an abandoned request. Anyone else sees "You cannot view this request." A missing request gives "Project request not found."

## Only the creator edits a request, and only in some states

Lines 177-192

Asked as: can't edit request; only creator edits; recall before editing; request locked for editing

Only the creator can edit a request: others see "Only the creator can edit this request." Editing works in status Draft, Recalled or Rejected. In other states Turbo says "Recall this request before editing it." Editing also needs project requests to be enabled.

## A stale version is refused

Lines 204-211

Asked as: request changed reload; version mismatch; someone else edited request; page out of date

If the version the page sends differs from the request's saved version, Turbo says "This request changed. Reload it before continuing."

## Editing a request resets its submission and approvals

Lines 212-249

Asked as: editing resets approval; edit clears submission; rejected back to draft; authorization PDF pending

Saving a business edit raises the request version, sets last_edited, and clears the submission token. A Rejected request returns to Draft and its rejection details clear. The linked authorization PDF record returns to Pending, loses its approval and gets the new project value.

## Project value must be a whole positive-or-zero CAD amount

Lines 250-265

Asked as: project value whole dollars; value not a number; project value cents; invalid project amount

Project value must be a whole number of at least 0, no larger than 9007199254740991. Otherwise Turbo says "Project value must be a whole CAD dollar amount."

## Only listed fields can be edited on a request

Lines 254-262

Asked as: field not editable; cannot change field; can't edit this field; read only field

Turbo refuses any field outside the allowed business fields with "Field is not editable: <field>".

## Request notes limits

Lines 266-277

Asked as: notes length limit; notes too short; notes character limit; request notes validation

Notes must be text. Turbo trims spaces and requires 10 to 1000 characters, or blank. Otherwise it says "Notes must be 10 to 1000 characters, or blank." The limits exist because job creation copies the notes into a client note.

## Fee basis must be set explicitly

Lines 278-283

Asked as: fee basis required; select time and materials; choose fee basis; fee basis not set

The fee basis (time_and_materials) must be true or false, or Turbo says "Select a fee basis." Setting it marks the basis as chosen. Submission without a chosen basis gives the same message.

## Draft validation of status, approver, parent and proposal

Lines 291-330

Asked as: project status active closed; invalid approver limit; parent project number; proposal status error

A request's project status must be Active or Closed, or Turbo says "Select Active or Closed for the project status." A selected approver's project limit must cover the project value. A parent must be a valid top-level project (number like 25-123) with the same client as the request. A proposal must be for the same client and have status In Progress, Submitted or Awarded, otherwise Turbo names its status and refuses.

## Job is created only from an Awarded proposal

Lines 332-357

Asked as: create job from proposal; proposal not awarded; why can't create job; award proposal first

A request with a proposal can be submitted and approved while the proposal is In Progress or Submitted. Creating the job needs the proposal to be Awarded, or Turbo says "Proposal <number> is <status>. Award it before you create the job." Turbo does not award the proposal itself.

## Required fields and PDF before submit or approve

Lines 358-406

Asked as: required field missing; cannot submit request; upload authorization PDF; project value must be positive

Description, location, client, contact, manager, branch, rate sheet, award date, approver and fn_agreement_answer must not be blank, or Turbo says "Required field: <field>". Project value must be above 0 ("Project value must be positive."). The request needs allocations, valid contact emails and an uploaded authorization PDF ("Upload an authorization PDF.").

## Allocation rules for requests

Lines 407-431

Asked as: allocation hours decimal; invalid division; inactive division; allocation whole hours

Allocation hours must be whole numbers ("Allocation hours must be whole numbers."). Each division must exist ("Select a valid division.") and be active ("Select an active division.").

## Editing an existing request needs it to be visible and editable

Lines 432-471

Asked as: cannot save existing request; edit request permission; save fails view check

Saving changes to an existing project request loads it through the view check, then the edit check, then the version check. A person who cannot view the request gets "You cannot view this request." Project requests being enabled is checked first.

## Creating a request needs the time claim

Lines 445-459

Asked as: create request time claim; who can create request; new request permission

A new project request needs the time claim. It starts as Draft with project status Active and version 1.

## A new request starts as an Active draft at version 1

Lines 446-459

Asked as: new request default status; request starts draft; created request version

A new project request is created with its creator set to the signed-in person, status Draft, project status Active, version 1 and the current time as last edited.

## Choose a rate sheet or propose rates, not both

Lines 483-533

Asked as: rate sheet and proposed rates; both rate sheet and rates; rate sheet conflict; propose rates with sheet

A request cannot both choose a rate sheet and propose new rates: Turbo says "Choose a rate sheet or propose new rates, not both." Choosing an active sheet drops the request's draft proposed rates. An unchanged save of an existing request changes nothing.

## Saving allocations replaces the old ones

Lines 510-575

Asked as: allocations not saving; changed allocations replaced; update allocation hours

When the allocations sent differ from the saved ones, Turbo deletes the old allocations and saves the new ones.

## Saving a request with no changes does nothing

Lines 524-533

Asked as: save with no changes; unchanged save version; nothing changed saving request

When an existing project request is saved and no business field, allocation or proposed rate changed, Turbo saves nothing. The version and last edited time stay the same, and the submission is not reset.

## Saving a new request saves its proposed rates after the request

Lines 534-541

Asked as: new request proposed rates; rates saved with new request; proposed rates new request

For a new request with proposed rates, Turbo saves the request first, then saves the proposed rate draft, then saves the request again, because the draft refers to the request.

## Checklist scope depends on request status

Lines 586-592

Asked as: which checklist shows; creation checklist vs approval; checklist for request status

A request that is Awaiting creation answers the creation checklist. Any other request answers the approval checklist.

## Who can reject a request

Lines 594-612

Asked as: who can reject request; reject button missing; approver reject permission

The selected approver, who must hold the project_authorization claim, can reject while the request is Awaiting approval or Awaiting creation. A job holder can reject only once the request is Awaiting creation. Others see "Only the selected approver, or a job holder once the request is approved, can reject this request."

## Checklist must be fully confirmed

Lines 614-641

Asked as: checklist not confirmed; confirm checklist questions; checklist changed reload; checklist answers false

The checklist version must match the current one (else "The checklist changed. Reload the request and review it again."). The answers must cover exactly the questions, and every answer must be true, or Turbo says "Confirm: <question>". With no saved checklist Turbo refuses.

## Reloading an out-of-date request page is refused on save and actions

Lines 642-660

Asked as: page out of date save; stale page action; reload before action

Save and every action on a project request compare the version the page sends with the saved version (version). When they differ, Turbo refuses with "This request changed. Reload it before continuing." A new request skips the save check because it has no saved version.

## Submitting a request

Lines 668-695

Asked as: submit project request; send request for approval; submit request fails

Submitting needs the time claim and an editable request. A Rejected request must be edited first ("Edit the rejected request before submission."). Submit validates the request fully, clears earlier approvals, sets status Awaiting approval, records the submit time and issues a new submission token.

## Submit needs the time claim and a non-rejected request

Lines 668-677

Asked as: cannot submit rejected; rejected request submit; submit needs edit first

Submitting needs the time claim and an editable request. A Rejected request gets "Edit the rejected request before submission."

## Submitting resets approvals and issues a new token

Lines 678-695

Asked as: submit resets approvals; new submission token; submit doesn't change version

Submit validates the request fully, clears the attachment approval and proposed rates approval, and sets status Awaiting approval and the submitted time. It issues a new 32-character submission token even when nothing changed. The version and last edited time stay the same.

## Recalling or rejecting a request

Lines 696-739

Asked as: recall or reject request; withdraw request; reject status requirements

Recall and reject work only on Awaiting approval or Awaiting creation requests ("This action requires a submitted or approved request."). Only the creator can recall (status Recalled). Reject needs a non-blank reason ("Enter a rejection reason."). Both end the submission token and clear approvals.

## Recall and reject need a submitted or approved request

Lines 696-705

Asked as: recall not allowed; only creator recall; recall button wrong status

Recall and reject work only on a request that is Awaiting approval or Awaiting creation. Otherwise Turbo says "This action requires a submitted or approved request." Only the creator can recall: others get "Only the creator can recall this request." Recall sets status Recalled.

## Rejecting needs a permitted person and a reason

Lines 706-724

Asked as: reject reason required; who can reject; reject without reason

A reject by someone who is not the selected approver, or a job holder once approved, gets "Only the selected approver, or a job holder once the request is approved, can reject this request." A blank reason gets "Enter a rejection reason." Reject sets status Rejected and records the time, rejector and trimmed reason.

## Recall and reject end the submission and clear approvals

Lines 725-739

Asked as: recall clears approvals; reject clears token; attachment marked rejected

Recall and reject clear the submission token and the proposed rates approval. They reset the attachment approval. On reject, the attachment is marked Rejected with the same rejector, time and reason.

## Approving a request

Lines 740-788

Asked as: approve project request; approve button; approval steps

Only a submitted, unrejected request in Awaiting approval can be approved. Only the selected approver can approve, and their project limit must cover the value. The request needs a submission token, a fully true checklist, and a matching review token and document hash, or Turbo says "The reviewed data changed. Reload and review it again." Approval records the limit and moves the request to Awaiting creation.

## Only the selected approver can approve

Lines 740-752

Asked as: who can approve; only approver approves; can't approve request

Only the approver selected on the request can approve it. Anyone else gets "Only the selected approver can approve this request." A request that is not Awaiting approval, or was rejected, gets "Only a submitted request can be approved."

## Approval requires a submission token

Lines 741-749

Asked as: no submission token; token missing approve; recall and submit again

Approving or creating a job fails with "This request has no submission token. Recall and submit it again." when the request has no token. Only submit issues a token, and edit, recall and reject clear it.

## Approver's limit and review data must match to approve

Lines 750-780

Asked as: reviewed data changed; approval document hash; limit does not cover

The approver's project limit (ProjectAuthorizationLimit) must cover the project value. The checklist and review token the page sends must match, and the document hash must match the attachment. Otherwise Turbo refuses with "The reviewed data changed. Reload and review it again." Approval records the limit as the approval maximum value.

## Creating the job from an approved request

Lines 789-868

Asked as: create job from request; job creation errors; approval stale create job

Creating needs the job claim and status Awaiting creation. It checks the creation checklist, the awarded proposal and that the approval is current. A stale approval gives "The approval is stale. Recall the request for a new approval." A changed page gives "The request changed after this page loaded. Reload the page, then create the job." An incomplete approval or approval limit below the project value is refused.

## Creating a job needs the job claim and an approved request

Lines 789-808

Asked as: create job needs claim; not awaiting creation; job claim creation

Creating a job from a request needs the job claim. The request must be Awaiting creation, or Turbo refuses with "This request is not awaiting creation." The creation checklist must be fully confirmed, and the request is validated without checking the approver's authority again.

## Stale approval blocks job creation

Lines 818-829

Asked as: approval is stale; recall for new approval; document changed after approval

Job creation is refused with "The approval is stale. Recall the request for a new approval." when the attachment is not Approved, its request version or project value differs from the request, its approved digest does not match, or the document hash differs.

## Out-of-date page blocks job creation but keeps the approval

Lines 830-837

Asked as: page changed after load; review token mismatch; reload then create job

When the review token the page sends is not the current one, Turbo refuses with "The request changed after this page loaded. Reload the page, then create the job." The saved approval stays valid, so no recall is needed.

## Incomplete saved approval blocks job creation

Lines 838-847

Asked as: saved approval incomplete; approval checklist invalid; rejection blocks job

Job creation is refused with "The saved approval is incomplete. Recall the request for a new approval." when the approver differs, there is no approval time, the approval maximum is below the project value, or a rejection is recorded. A saved approval checklist version that no longer exists gives "The saved checklist is invalid."

## Job creation records checklist and handles proposed rates

Lines 848-854

Asked as: rename rate sheet job; creation checklist version; proposed rate sheet name

Job creation records the creation checklist version on the attachment, requires proposed rates to be approved if the request has them, and lets the job creator rename the proposed rate sheet (rate_sheet_name). An empty name keeps the requester's name.

## Job creation returns the new job id and sends notices

Lines 859-867

Asked as: job id returned; created notice; new job number after create

After a successful create, Turbo answers with the new job's id as both id and job_id, and sends the queued notices. A "created" notice is queued for the request.

## Approved request is saved with its attachment and notice

Lines 869-883

Asked as: notice after approve; withdrawn approval flag; attachment saved on submit

After submit, recall, reject or approve, Turbo saves the authorization attachment if there is one, then the request. After submit, approve or reject, Turbo then queues a notice. A reject by the approver on a request that is Awaiting creation is flagged as a withdrawn approval. A job holder's reject is not.

## What job creation copies and removes

Lines 896-968

Asked as: what copies to job; request deleted after job; job fields from request

Job creation copies the request's business fields except the approver. It sets the job status from the project status, sets the First Nations agreement flag when the answer is Yes, and sets the authorizing document to PA. The request's allocations become job allocations, the PDF record moves to the job as its first approval, and the notes become a client note by the request creator. The request is then deleted.

## Approved proposed rates become an active rate sheet at creation

Lines 896-901

Asked as: proposed rates active sheet; rate sheet activated job; approved rates become sheet

Before the job is saved, Turbo activates the request's approved proposed rates as an active rate sheet (ActivateRequestRateSheet), so the job can use it.

## Job creation sets the job's origin fields

Lines 906-922

Asked as: job status from request; first nations agreement flag; created_by job; authorizing document PA

The new job takes its status from the request's project status. Its First Nations agreement (fn_agreement) is true only when the answer is "Yes". Its authorizing document is "PA". Turbo also records the request's creator, created and submitted times, the person who created the job (created_by) and the completion time.

## Job creation moves allocations and the authorization attachment

Lines 927-966

Asked as: allocations move to job; authorization PDF moves; request deleted after create

Job creation copies each request allocation to a job time allocation and deletes the request allocation. The authorization attachment moves from the request to the job, is marked as the initial approval and gets approval sequence 1. Turbo then deletes the request.

## Request notes become the job's first client note

Lines 970-991

Asked as: request notes client note; notes copied to job; job first note

When the request notes are not blank after trimming, job creation saves them as a client note on the new job and its client. The note's author is the request creator. Blank notes make no note.

## Deleting a request

Lines 992-1038

Asked as: delete project request; cannot delete request; only abandoned delete

The creator can delete a request they can still edit. Otherwise only an admin can delete it, and only when abandoned, or Turbo says "Only abandoned requests can be deleted." Deleting also removes its authorization record and allocations.
