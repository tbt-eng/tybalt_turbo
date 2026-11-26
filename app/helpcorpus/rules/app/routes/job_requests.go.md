<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_requests.go","sha256":"2307b7a506d7345b6a39e6fc372fef09b16d3806811c6b03051350e6add7bd13","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_requests.go

## Abandoned requests are those not edited for the abandonment period

Lines 27-27

Asked as: request went stale; old job request; request not touched for days; abandoned job request

A project request counts as abandoned when its last edit (last_edited) is at or before now minus JOB_REQUEST_ABANDONMENT_DAYS days. A request edited exactly that many days ago is already abandoned. The number of days comes from utilities.JOB_REQUEST_ABANDONMENT_DAYS, defined in another file.

## Project request rules apply only to listed business fields

Lines 30-30

Asked as: which fields can I edit; fields copied to new job; can't change this field; editable request fields

A project request's editable fields are: project_status, approver, client, job_owner, contact, manager, alternate_manager, branch, rate_sheet, parent, proposal, invoicing_information, description, notes, location, client_po, client_reference_number, project_award_date, project_completion_date, project_value, time_and_materials and fn_agreement_answer. These fields are also copied to the new job at creation, except the approver.

## Request save bodies carry rates, allocations and a version

Lines 32-39

Asked as: propose new rates; save rates with request; allocations not saving; rates on job request

A save of a project request can carry the request fields, a list of division allocations, the version the person saw, and proposed rates. Proposed rates replace the choice of an active rate sheet.

## Approver can send a reason, checklist answers and review data

Lines 40-53

Asked as: approver rejection reason; approver checklist answers; approve proposed rates; new rate sheet name

An action on a project request carries the version, the review token, the document hash, the checklist version and answers, and a rejection reason. It also carries the approver's statement that they approve proposed rates (rates_approved) and an optional new rate sheet name for the job creator.

## Project request routes require sign-in

Lines 55-85

Asked as: asked to sign in; signed out request; log in to request; need to log in

All project request and job attachment routes require a signed-in user. Turbo says 'Sign in to continue.' when there is no signed-in user.

## Unknown actions on an attachment or its document

Lines 70-72

Asked as: attachment action not working; authorization document action; unknown request action; wrong action on PDF

Only the listed actions submit, recall, reject, approve and create exist for a project request. Any of them can also be sent as an action on a request's authorization document route.

## Workflow errors from hooks and validation are returned as messages

Lines 87-98

Asked as: error message on request; request validation error; why request failed; hook error message

When a project request route fails with a hook error, Turbo shows that error's message and details. When field validation fails, Turbo shows 'Check the request fields.' with the field errors.

## Active staff only for project request actions

Lines 104-129

Asked as: not active staff; can't act on request; missing claim error; staff access required

Turbo refuses a user who is not an active staff member with "An active staff member is required." A missing claim is refused with "The <claim> claim is required."

## Project requests can be switched off

Lines 130-139

Asked as: project requests turned off; can't save job request; requests disabled setting; job requests switched off

When the project request setting is off, Turbo refuses to save, act on or delete project requests with "Project requests are disabled." The check is IsJobRequestApprovalRequired. Delete is the exception, because it checks the setting only through the editable check.

## Abandoned request definition

Lines 140-142

Asked as: when is request abandoned; how old before abandoned; request counts as abandoned; abandonment period

A request is abandoned when its last edit (last_edited) is JOB_REQUEST_ABANDONMENT_DAYS days old or older. A request edited exactly that many days ago counts as abandoned.

## Who can view a project request

Lines 143-174

Asked as: who can see request; can't view request; request hidden from me; access to job request

The creator can always view a request. A procurement user (the job claim) can view every request. A project authorization approver (the project_authorization claim) can view only requests that name them as approver. An admin can view only abandoned requests. Anyone else gets "You cannot view this request."

## Save loads a request the person can view

Lines 161-174

Asked as: request not found; can't open request to save; cannot view this request; save existing request fails

Saving or acting on an existing project request first requires that the person can view it. Otherwise Turbo says 'Project request not found.' or 'You cannot view this request.'

## Only the creator edits, and only in some states

Lines 175-190

Asked as: can't edit request; only creator can edit; recall before editing; locked request edits

Only the creator can edit a request, otherwise Turbo says "Only the creator can edit this request." A request can be edited only when it is Draft, Recalled or Rejected. In any other status Turbo says "Recall this request before editing it."

## Duplicate lookups on the attachment

Lines 191-197

Asked as: missing request PDF; no attachment on request; authorization attachment missing; where is the attachment

Each project request has at most one authorization attachment, found by its request. A request without one has none, which the submit check treats as missing the PDF.

## Allocations sort by division when read

Lines 198-200

Asked as: allocations order; division order list; allocations sorted wrong

A request's allocations are listed ordered by division.

## Stale version is refused

Lines 202-209

Asked as: request changed reload; someone else edited request; out of date version; version mismatch error

If the version sent does not match the request's current version, Turbo says "This request changed. Reload it before continuing."

## Editing a request resets its approval

Lines 210-247

Asked as: edit resets approval; approval gone after edit; editing clears approval; request version increase

A saved edit raises the request version, updates last_edited and clears the submission token. A Rejected request returns to Draft, and its rejection details are cleared. The attachment's approval is reset to Pending, and it takes the new project value.

## Rejected request becomes a draft after editing

Lines 230-233

Asked as: rejected request after edit; rejected back to draft; clear rejection on edit

Editing a Rejected request sets its status back to Draft and clears the rejector, time and reason.

## Project value must be a whole amount

Lines 248-262

Asked as: project value whole dollars; value must be integer; cents in project value; negative project value

Project value must be a whole CAD dollar amount of 0 or more, otherwise Turbo says "Project value must be a whole CAD dollar amount." A request must have a value above 0 to submit.

## Only listed fields can be edited

Lines 252-263

Asked as: field not editable; can't change that field; save refused field

A save that includes any field outside the request's editable business fields is refused with "Field is not editable: <field>".

## Notes length limit

Lines 264-275

Asked as: notes too long; notes length error; notes character limit

Notes are trimmed. They must be text, and 10 to 1000 characters, or blank. Otherwise Turbo says "Notes must be 10 to 1000 characters, or blank."

## Fee basis must be chosen

Lines 276-281

Asked as: fee basis missing; choose time and materials; select fee basis

The fee basis (time_and_materials) must be true or false, otherwise Turbo says "Select a fee basis." Submit and approve also say this until a fee basis is set.

## Draft checks: project status

Lines 289-292

Asked as: project status wrong; active or closed status; status must be active

The project status must be Active or Closed, otherwise Turbo says "Select Active or Closed for the project status." A new request starts as Draft with project status Active.

## Approver's project limit is checked

Lines 296-300

Asked as: approver project limit; approver limit too low; approver can't cover value

When a request names an approver, Turbo checks that approver's project authorization limit against the project value, through ProjectAuthorizationLimit. Approval and submit run this check. Creating the job skips it.

## Parent project rules

Lines 301-312

Asked as: parent project invalid; subjob client mismatch; top-level parent project

A parent must be an existing job, otherwise Turbo says "Select a valid parent project." Its number must be top-level, like 25-123, otherwise "Only a top-level project can be a parent." The request's client must equal the parent's client, otherwise "The subjob must use the parent project's client."

## Proposal rules for a request

Lines 313-326

Asked as: proposal not valid; wrong proposal status; select proposal for client

A linked proposal must exist, start with P and belong to the request's client, otherwise Turbo says "Select a proposal for this client." Its status must be In Progress, Submitted or Awarded, otherwise Turbo says "Proposal <number> is <status>. Select an In Progress, Submitted, or Awarded proposal."

## Job is created only from an awarded proposal

Lines 330-355

Asked as: can't create job yet; proposal not awarded; award proposal before job

A request can be submitted and approved while its proposal is In Progress or Submitted. Creating the job needs the proposal Awarded, otherwise Turbo says "Proposal <number> is <status>. Award it before you create the job." Turbo does not award the proposal itself. This check runs before the approval checks.

## Required fields before submit

Lines 356-370

Asked as: required field before submit; can't submit request; missing field to submit

Submit and approve require these fields: description, location, client, contact, manager, branch, rate sheet, award date, invoicing information, approver and fund agreement answer. Otherwise Turbo says "Required field: <field>". They also require a fee basis and a positive project value.

## Request validation checks contact emails and hooks

Lines 373-392

Asked as: contact email error; client reference check; submit validation fails

Submit and approve also check the client references, the contact emails and the project request data rules (ValidateProjectRequestData in hooks). Their messages come from other files.

## Allocations are required and must be whole hours

Lines 379-389

Asked as: need division allocations; allocation hours decimal; inactive division allocation

A complete request needs division allocations. Hours must be whole numbers, otherwise "Allocation hours must be whole numbers." Each division must exist ("Select a valid division.") and be active ("Select an active division.").

## Authorization PDF required

Lines 393-403

Asked as: upload PDF to request; no authorization file; PDF not uploaded

A request needs an uploaded authorization PDF with a stored hash, otherwise Turbo says "Upload an authorization PDF."

## Allocation divisions must be valid, active and whole hours

Lines 405-429

Asked as: invalid allocation hours; allocation division inactive; allocation hours not whole

Each allocation hours value must be a finite whole number, otherwise 'Invalid allocation hours.' or 'Allocation hours must be whole numbers.' The division must exist ('Select a valid division.') and be active ('Select an active division.'). A save with no allocations skips these checks.

## Time claim holders can create requests

Lines 437-447

Asked as: who can create request; no time claim; can't start new request

Creating a new project request requires the time claim. Project requests must also be enabled.

## Creating a request needs the time claim

Lines 444-457

Asked as: can't create request; need time claim; new request not allowed

Only a user with the time claim can create a new project request. It starts as Draft, version 1.

## New project request starts as an Active draft

Lines 448-457

Asked as: new request starts draft; what status when created; first version of request

A new project request starts with status Draft, project status Active and version 1. Its creator is the person who saved it, and its last edited time is now.

## Editing an existing request needs the same checks as other edits

Lines 458-469

Asked as: save changes to request; can't update request; edit existing request checks

When a person saves changes to an existing project request, Turbo loads it only if the person can view it. Turbo then requires that the person is the creator and that the request is Draft, Recalled or Rejected. Turbo also requires that the version the person saw is the current one.

## Save compares the request before and after the change

Lines 470-480

Asked as: save changes nothing; did my edit save; detect request change

Before it applies the edits, Turbo records the request's business fields and whether the fee basis was chosen (time_and_materials_set). After the edit, any difference counts as a change. A save that changes nothing does not update the request.

## Rate sheet or proposed rates, not both

Lines 481-507

Asked as: rate sheet and proposed rates; both rate options error; replace draft rates

A save cannot both choose a rate sheet and propose new rates, otherwise Turbo says "Choose a rate sheet or propose new rates, not both." Choosing a different sheet deletes the request's draft proposed rates.

## Saving allocations replaces the old ones

Lines 508-573

Asked as: allocations replaced on save; old allocations deleted; update division hours

If the allocations differ, Turbo deletes the old rows and saves the new ones. Allocations are checked only when sent.

## Unchanged save does nothing

Lines 522-531

Asked as: save with no changes; unchanged request version; saving does nothing

Saving an existing request with no change to its fields, allocations or proposed rates does not raise the version or change anything.

## Which checklist applies depends on the request status

Lines 585-590

Asked as: which checklist to answer; creation checklist; approval checklist why

A request that is Awaiting creation uses the creation checklist. In any other status it uses the approval checklist.

## Who can reject and the reason

Lines 592-610

Asked as: who can reject request; reject reason required; can't reject request

The selected approver (with the project_authorization claim) can reject while the request is Awaiting approval or Awaiting creation. A job holder (the job claim) can reject only once it is Awaiting creation. Anyone else gets "Only the selected approver, or a job holder once the request is approved, can reject this request." A rejection needs a reason, otherwise "Enter a rejection reason."

## Checklist must be answered fully

Lines 612-639

Asked as: checklist not confirmed; confirm checklist question; checklist not saved

Approval and creation each need the current checklist for their step. Every answer must be true, otherwise Turbo says "Confirm: <question>". If no checklist is saved, Turbo says "No <scope> checklist is saved." An unsupported version is refused. An old version gets "The checklist changed. Reload the request and review it again." The number of answers must match the questions.

## Checklist version and answers must match exactly

Lines 614-638

Asked as: checklist version error; checklist changed reload; checklist answers mismatch

If no checklist is saved, Turbo says 'No <scope> checklist is saved.' A version below 1 or above the current one is 'Checklist version N is not supported.' An older version gets 'The checklist changed. Reload the request and review it again.' The answers must be exactly as many as the questions, and each must be true, otherwise Turbo says 'Confirm: <question key>'.

## Actions run only while project requests are enabled

Lines 648-651

Asked as: actions disabled; project requests off action; can't submit when disabled

Each action on a project request first checks that project requests are enabled; if not, Turbo says 'Project requests are disabled.'

## Every action checks the version the person saw

Lines 656-658

Asked as: version changed before action; stale request action; reload before submit

Submit, recall, reject, approve and create refuse when the request's version differs from the version sent: 'This request changed. Reload it before continuing.'

## Submit rules

Lines 666-693

Asked as: how to submit request; submit from draft; submitted request steps

Submitting needs the time claim, and only the creator can do it from Draft or Recalled. A Rejected request cannot be submitted: "Edit the rejected request before submission." Submit moves the request to Awaiting approval, records the submitted time and issues a new submission token, even if the data is unchanged. It clears any earlier approval.

## Submit requires the time claim and an editable request

Lines 666-675

Asked as: submit rejected request; can't submit rejected; who can submit

Submitting a request requires the time claim and that the person is the creator with status Draft, Recalled or Rejected. A Rejected request must be edited first: 'Edit the rejected request before submission.'

## Submit sets awaiting approval and issues a new token

Lines 676-693

Asked as: submit request token; awaiting approval status; new submission token

A successful submit sets status Awaiting approval, records the submitted time, issues a new random 32-character submission token and clears any earlier attachment approval and proposed-rate approval. It does not change the version or last edit.

## Recall and reject need a submitted request

Lines 694-737

Asked as: recall or reject not allowed; action needs submitted request; reject before submit

Recall and reject work only on a request that is Awaiting approval or Awaiting creation. Otherwise Turbo says "This action requires a submitted or approved request." Either action ends the submission and clears the approval.

## Reject and recall need a submitted request

Lines 694-697

Asked as: reject request not submitted; recall request status; submitted request required

Recall and reject work only when the request is Awaiting approval or Awaiting creation. Otherwise Turbo says 'This action requires a submitted or approved request.'

## Only the creator can recall

Lines 698-703

Asked as: who can recall; recall someone else's request; only creator recalls

Only the creator can recall a request, otherwise Turbo says "Only the creator can recall this request." Recall sets the status to Recalled.

## Recall clears the rejection and submission token

Lines 698-730

Asked as: recall clears rejection; take back submission; withdraw request submission

A recall sets status Recalled, clears any rejection details, ends the submission by clearing its token and clears the proposed-rate and attachment approval.

## Rejection records who, when and why

Lines 714-737

Asked as: rejection recorded details; who rejected request; rejection reason saved

A rejection sets the status to Rejected and stores the time, the rejector and the trimmed reason. The attachment is also marked Rejected.

## Approve: state and approver

Lines 738-750

Asked as: approve request not allowed; only selected approver; no submission token approval

Only a request that is Awaiting approval and not rejected can be approved, otherwise "Only a submitted request can be approved." Only the selected approver can approve, otherwise "Only the selected approver can approve this request." A request without a submission token is refused with "This request has no submission token. Recall and submit it again."

## Approval needs an Awaiting approval request and the selected approver

Lines 739-757

Asked as: approve needs awaiting status; why can't I approve; approver limit approval

Only a request in status Awaiting approval with no rejection can be approved: 'Only a submitted request can be approved.' Only the selected approver can approve: 'Only the selected approver can approve this request.' The approver's project limit must also cover the project value.

## Approval needs a submission token

Lines 745-747

Asked as: submission token missing; recall and submit again; no token error

Approval and creation refuse a request with no submission token: 'This request has no submission token. Recall and submit it again.'

## Approve: limit, checklist and review data

Lines 751-768

Asked as: reviewed data changed; approval checklist complete; approve limit checklist

Approval needs the approver's limit to cover the project value, a fully confirmed current checklist, and a complete request. The review token and document hash sent must match the current ones, otherwise Turbo says "The reviewed data changed. Reload and review it again."

## Review token and document hash must match at approval

Lines 762-768

Asked as: review token mismatch; document hash changed; reload and review again

Approval requires the review token and document hash the approver reviewed. If either differs from the current ones, Turbo says 'The reviewed data changed. Reload and review it again.'

## Approval result

Lines 769-786

Asked as: what approval records; approval result saved; approved request next step

Approval marks the attachment Approved and records the approver, the time, the value, the limit (approval_max_value), the request version and the checklist version. The request moves to Awaiting creation. Turbo then saves a new review token for the approved state. Approving also calls approveProposedRates with the approver's rates_approved statement.

## Create job: who and when

Lines 787-806

Asked as: who can create job; not awaiting creation; create job without token

Only a user with the job claim can create the job, and only from Awaiting creation, otherwise Turbo says "This request is not awaiting creation." A request without a submission token is refused. Creation checks the request without the approver's limit.

## Creation needs an approved request and stale approval is refused

Lines 791-841

Asked as: approval is stale; saved approval incomplete; creation refused approval

Creation needs the job claim and status Awaiting creation ('This request is not awaiting creation.'). It refuses with 'The approval is stale. Recall the request for a new approval.' when the attachment is not Approved, its request version or project value differs from the request's, or the review token or document hash differs. If the approver, approval time or limit does not cover the value, or either record is rejected, Turbo says 'The saved approval is incomplete. Recall the request for a new approval.' A saved checklist version that no longer exists gives 'The saved checklist is invalid.'

## Job creation needs the creation checklist and a token

Lines 797-806

Asked as: create job checklist; creation needs token; job creation requirements

Creating the job requires a submission token and that the creation checklist is fully confirmed. The request is then validated without checking the approver's limit again.

## Create job: stale or incomplete approval

Lines 807-841

Asked as: stale approval creating job; approval missing limit; saved checklist invalid

Creation is refused with "The approval is stale. Recall the request for a new approval." if the attachment is not Approved, its version or value differs from the request, or the tokens or document hash do not match. It is refused with "The saved approval is incomplete. Recall the request for a new approval." if the approver, approval time or limit is missing or too low, or if anything is marked rejected. A missing approval checklist gives "The saved checklist is invalid."

## Creation stores which checklist the creator answered

Lines 842-842

Asked as: creation checklist version saved; which checklist creator answered; record creation checklist

After the approval checks pass, the creation checklist version the creator confirmed is saved on the attachment (creation_checklist_version).

## Rate sheet name can change at creation

Lines 843-895

Asked as: rename rate sheet; change rate sheet name; proposed rates at creation

Proposed rates need approval before creation (requireProposedRatesApproval). The job creator can rename them (rate_sheet_name). Approved proposed rates become an active rate sheet when the job is created.

## Reject and recall notifications

Lines 871-877

Asked as: withdrawn approval notice; notification after reject; approval withdrawn email

After submit, approve or reject, Turbo queues a notice. A rejection by the approver once the request is Awaiting creation is marked as a withdrawn approval; a job holder's rejection is not.

## Job created from the request

Lines 890-917

Asked as: job from request fields; fund agreement flag; authorizing document PA

Creating the job copies the request's business fields except the approver. The job status comes from the project status. The fund agreement flag is true when the answer is "Yes". The authorizing document is "PA". Turbo records the request creator and the dates, the creating user and the completion time.

## Job is created from request fields

Lines 890-916

Asked as: fields copied to job; job status from project; job created details

Creating the job copies the request's business fields, except the approver, to a new job. The job's status is the request's project status, its fn_agreement is true only when the answer is 'Yes', and its authorizing document is 'PA'. It records the request's creator, created time and submitted time, who created it and when creation completed. Approved proposed rates become an active rate sheet first.

## Allocations, attachment and notes move to the job

Lines 921-962

Asked as: allocations move to job; attachment moves to job; request deleted after job

Creation turns request allocations into job time allocations and deletes the request rows. The attachment moves to the job as its initial approval, with approval sequence 1. Notes become the job's first client note, written by the request creator. The request is then deleted.

## Request notes become the job's first client note

Lines 964-985

Asked as: notes become client note; request notes on job; first client note

When the request's notes are not blank, the created job gets a client note with those notes. The note's author is the request creator. Blank notes make no note.

## Who can delete a request

Lines 986-1031

Asked as: delete job request; can't delete request; who can delete request

The creator can delete a request they can still edit (Draft, Recalled or Rejected). Anyone else needs the admin claim, and the request must be abandoned, otherwise Turbo says "Only abandoned requests can be deleted." Deleting also removes its attachment and allocations.
