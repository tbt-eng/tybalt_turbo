<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_attachments.go","sha256":"75b646c98fcbb63002093e7fb99b039e2c18072a6cd0b3cc358cf4d0f9128c55","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_attachments.go

## Authorization PDF upload size limit

Lines 25-51

Asked as: pdf too large; attachment over 20 MB; upload size limit; can't upload big pdf

Turbo rejects an authorization PDF larger than 20 MB with 'The PDF must be at most 20 MB.' The whole request body is capped at 21 MB.

## Exactly one authorization PDF per upload

Lines 25-38

Asked as: upload two pdfs; multiple files attached; no file selected; file missing from upload

When a PDF is required, the upload must contain exactly one file in the 'document' field. Otherwise Turbo says 'Upload one authorization PDF.' When the PDF is optional and none is sent, Turbo accepts the request and keeps the existing PDF. A file that cannot be read gives 'Cannot read the authorization PDF.'

## Authorization PDF must be a real PDF

Lines 52-54

Asked as: wrong file type; image instead of pdf; renamed file won't upload; not a real pdf

Turbo checks that the uploaded file starts with '%PDF-'. Otherwise it refuses with 'Upload a valid PDF file.'

## Same PDF cannot be used twice

Lines 58-76

Asked as: same pdf uploaded twice; pdf already used elsewhere; duplicate attachment; reuse old approval pdf

Turbo hashes each PDF (SHA-256) and refuses one that already belongs to another project request or job approval: 'This PDF already belongs to a project request or job approval.' If the other copy is on an abandoned project request, the message names that request ID and tells the person to give it to an admin to resolve the duplicate.

## Removing a project request PDF needs valid action data

Lines 83-88

Asked as: can't remove pdf; remove attachment button broken; invalid action data removing

When someone removes the authorization PDF from a project request, Turbo reads the request's version from the submitted action. If the action data cannot be read, Turbo refuses with "Invalid action data." and removes nothing.

## Uploading a PDF to a project request requires a file

Lines 89-94

Asked as: upload pdf to project request; project request pdf missing; version number wrong upload

Uploading an authorization PDF to a project request always needs one PDF. A missing or extra file gets the PDF upload error. Turbo takes the request version from the form's "version" value. A version that is not a whole number counts as 0.

## Upload or remove PDF on a project request

Lines 96-145

Asked as: can't change pdf on request; replace pdf project request; delete pdf request not allowed

Uploading or removing the authorization PDF on a project request works only when the person can see the request, the request is editable by them, and the version they sent matches the current one. Each upload or removal marks the request as edited.

## Project request PDF actions need an editable request at the current version

Lines 96-106

Asked as: request already changed; pdf version out of date; stale project request; locked project request

To upload or remove the PDF on a project request, the request must be visible to the person, be editable by them, and match the version they sent. Turbo checks all three before it changes anything. A successful upload or removal marks the request as edited.

## A project request's PDF upload is saved and removal is cleaned up

Lines 111-145

Asked as: pdf attachment saved; remove attachment record; approval cleared after upload; attachment version goes up

Removing a project request's PDF deletes its attachment record. Uploading the PDF already attached, with the same hash, changes nothing except re-checking the file. A new upload records the uploader, the upload time and the request's project value. It also raises the attachment version by 1 and clears any earlier approval.

## Removing a PDF with none attached

Lines 112-118

Asked as: remove pdf nothing attached; no pdf to delete; delete empty attachment

Removing the PDF from a project request that has no PDF does nothing and gives no error. Otherwise Turbo deletes the attachment.

## Uploading a PDF to a project request

Lines 119-144

Asked as: replace project request pdf; new pdf overwrites old; restore missing stored file

Uploading a different PDF to a project request replaces the earlier one. Turbo records the uploader, the upload time and the request's project value, clears any earlier approval, and raises the attachment version by 1. Uploading the identical PDF changes nothing except restoring a missing stored file.

## Who can view an authorization PDF

Lines 153-209

Asked as: can't open pdf; view authorization pdf access; pdf download forbidden; who sees the pdf

For a PDF on a project request, anyone who can view that request can view it. For a PDF on a job's value change, these people can view it: a job editor (the job claim), an admin, the project's manager or alternate manager, and an approver (the project_authorization claim) who is the assigned approver or the approver of that change. Anyone else gets 'You cannot view this PDF.' The download is not cached by browsers.

## Value change list needs a job or a queue

Lines 210-233

Asked as: select a job; value changes list blank; no job chosen changes

The list of value changes shows either the approval queue or the changes for one job. Without the 'approval' or 'approved_by_me' queue, Turbo requires a job and otherwise says 'Select a job.'

## Approval queue lists only my pending changes

Lines 221-226

Asked as: my approval queue; pending approvals to review; approvals waiting for me

The approval queue needs the project authorization claim. It lists only pending value changes assigned to the signed-in approver. Changes the person cannot view are left out.

## Value change list contents and order

Lines 236-291

Asked as: value change list order; newest approvals first; list shows approver names

Each listed change shows the job number, description, current project value and related records. Newest approvals come first (by approval sequence, then creation time). Names of the uploader, approver, assigned approver and rejector are added. Turbo flags as current the approved change that set the job's present value. A pending change also gets a review token.

## Value change list is limited by what the viewer may see

Lines 244-251

Asked as: missing changes in list; why can't I see change; hidden value changes

The value change list skips any change the signed-in person cannot view. A job's manager and alternate manager, procurement (the job claim) and admins can view all of the job's changes. An approver (the project_authorization claim) sees only changes where they are the assigned or actual approver.

## Value change list marks the current approval and gives a review token

Lines 262-287

Asked as: which approval is current; current approved change flag; review token missing

Each row in the value change list carries a flag (current) that is true only for the approved change whose approval sequence matches the job's latest approval sequence. Each pending row also carries a current_review_token that the reviewer sends back when rejecting. Each row shows the job number, description, effective project value, and the names of the uploader, approver, assigned approver and rejector.

## Approved-by-me list covers 30 days

Lines 299-332

Asked as: my approvals last month; approved by me list; recent approvals history

An approver (the project authorization claim) can list their own approvals of new jobs and value changes from the last 30 days (APPROVED_BY_ME_DAYS). Without that claim Turbo refuses.

## Value change review data token changes when job or related data change

Lines 346-397

Asked as: review token changed; data changed token mismatch; timestamps change token

Turbo identifies the data a value change is reviewed against with a token (valueReview.token). The token covers the job's fields (not created, updated or _imported), the latest approval sequence, the related records (requestRelated) and the division and hours of the job's time allocations. A change to any of these changes the token. Timestamps and allocation row identities do not change it.

## Latest approval sequence is read from approved changes only

Lines 366-374

Asked as: approval sequence number; delete approval invalidates review; latest approval number

Turbo finds a job's latest approval sequence from its approved value changes (status Approved), highest approval_sequence first. A job with no approved change has sequence 0. The review token includes this sequence, so a new approval, or the deletion of one, invalidates earlier review tokens.

## Approving a value change also sets the job's authorizing document to PA

Lines 409-426

Asked as: authorizing document PA; job authorizing document changed; approval sets project value

When a value change is approved, Turbo sets the job's project value to the approved value and sets its authorizing document (authorizing_document) to "PA". It does this even on older projects that have no authorizing document. Turbo applies the same two changes when it checks that the job would pass validation with the new value.

## Only job editors save a value change

Lines 428-432

Asked as: can't save value change; save value change permission; edit project value not allowed

Creating or editing a project value change needs a job editor (the job claim).

## New value change requires a PDF; edit keeps it

Lines 433-437

Asked as: value change needs pdf; edit without pdf; attach pdf to change

Creating a value change requires an authorization PDF. When editing, the PDF is optional and the existing one stays if none is sent.

## Project value must be positive whole dollars

Lines 438-441

Asked as: project value decimal; negative project value; cents not allowed value

A value change must have a positive whole-dollar project value. Otherwise Turbo says 'Enter a positive whole-dollar project value.'

## Assigned approver's limit must cover the value

Lines 450-452

Asked as: approver limit too low; assigned approver can't authorize; authorization limit exceeded

The assigned approver must be able to authorize the value. Turbo checks this with ProjectAuthorizationLimit and refuses the save if the check fails.

## Approved value changes cannot be edited

Lines 462-469

Asked as: edit approved value change; change approved evidence; approved change locked

Editing an attachment that is approved, or one that has no job, is refused with 'Approved evidence cannot be edited.'

## Stale value change edit is refused

Lines 470-474

Asked as: someone else edited change; value change changed reload; version conflict save

If the version sent differs from the saved version, Turbo refuses with 'This value change changed. Reload it.'

## Proposals have no value approval

Lines 480-482

Asked as: proposal project value approval; P job value change; proposal approval not needed

A job whose number starts with P is a proposal. Turbo refuses a value change on it: 'Proposals do not use project-value approval.'

## Job must already be valid for a value change

Lines 483-495

Asked as: job validation failed value change; allocations check value; job fields ignored

A value change does not edit the job's details or time allocations. Turbo checks the job's stored allocations and checks that the job would pass validation with the new value. Any job fields sent with the change are ignored.

## One open value change per job

Lines 496-502

Asked as: another pending change; already pending value change; rejected change blocks new

A job can have only one pending or rejected value change at a time. Turbo refuses another with 'Resolve the current pending value change first.'

## Unchanged save keeps state

Lines 507-520

Asked as: saved without changes; save same value change; rejected stays rejected

If the job data, value, assigned approver and PDF all match the last save, Turbo keeps the change as it is. A rejected change stays rejected, and the version does not change. A correction in the job editor also needs fresh approval.

## A real change resubmits for approval

Lines 521-547

Asked as: resubmit for approval after edit; submitted notice sent; approval cleared after change

Saving a changed value, approver, PDF or job data clears earlier approval details, raises the version by 1 and records the review data. It sends a 'submitted' notice. The notice says it replaces an earlier version only when the assigned approver is the same.

## Value change notices go out after the save or decision is stored

Lines 548-554

Asked as: email not sent; notice after save failed; notification timing

Turbo sends the notices for a submitted, approved or rejected value change only after the database transaction succeeds. If the transaction fails, Turbo sends no notice and returns the error. On success it returns the saved value change record.

## Approving or rejecting a value change needs valid action data

Lines 556-561

Asked as: approve button not working; reject invalid action; decision refused action data

When an approver approves or rejects a value change, Turbo first reads the submitted action. If the action data cannot be read, Turbo refuses with "Invalid action data." and makes no decision.

## Review needs the project authorization claim

Lines 564-574

Asked as: can't approve value change; reject value change permission; only pending reviewed

Approving or rejecting a value change needs the project authorization claim. Only a pending change with a job can be reviewed: 'Only a pending value change can be reviewed.'

## Only the selected approver reviews

Lines 575-577

Asked as: not my approval; wrong approver reviewing; assigned approver only

Only the assigned approver can approve or reject a value change. Others get 'Only the selected approver can review this value change.'

## Review must match the version and PDF seen

Lines 578-581

Asked as: reviewed changed reload; pdf changed before approve; stale approval page

A review applies only if the version and PDF hash match what the reviewer saw. Otherwise Turbo says 'The reviewed value change changed. Reload it.'

## Rejection needs a reason and fresh data

Lines 583-604

Asked as: reject without reason; rejection needs comment; job changed before reject

Rejecting a value change needs a non-blank reason: 'Enter a rejection reason.' If the job or related data changed since the page was shown, Turbo refuses with 'The job or its related data changed. Reload it before rejecting.' A rejection records the rejector, time, reason and reviewed data, and marks the change Rejected.

## Rejected value change keeps its token and stays rejected until corrected

Lines 597-604

Asked as: rejected stays rejected; rejection record reason; correct rejected change

When a change is rejected, Turbo records the review token of the data the approver saw, together with the rejector, the time and the trimmed reason. Only a later correction, in the job editor or in a save of the change, makes the change different from this recorded data. Saving the change unchanged leaves it rejected.

## Approval requirements

Lines 605-634

Asked as: can't approve pdf unavailable; job changed ask editor; approval validation failed

Approving checks that the stored PDF exists. If it is missing, Turbo says 'The authorization PDF is unavailable. Restore or replace it before approval.' The approver's limit must cover the value. The job and related data must be unchanged since the last save, or Turbo says 'The job changed. Ask a job editor to save and review this value change again.' The job's stored allocations and validation must also pass.

## Approval refuses when the approver's limit does not cover the value

Lines 606-612

Asked as: approval limit exceeded; approval max value; pdf missing at approve

When approving a value change, Turbo checks the approver's project authorization limit (ProjectAuthorizationLimit) against the change's project value. It stores that limit as approval_max_value on the change. Turbo checks the stored PDF exists at approval, and refuses with "The authorization PDF is unavailable. Restore or replace it before approval." if it is missing.

## Approval refuses a stale review token

Lines 617-623

Asked as: job changed cannot approve; approve after job edit; stale review approve

Turbo refuses to approve a value change when the job or its related data changed since the change was last saved. The message is "The job changed. Ask a job editor to save and review this value change again."

## What approval does

Lines 636-669

Asked as: approve sets value; previous value recorded; what approval does

Approval sets the job's project value to the approved value. It sets the job's authorizing document to 'PA'. It records the previous value, the next approval sequence number, the approver, the time and the approver's limit, and marks the change Approved. Each decision raises the version by 1 and sends an approved or rejected notice.

## Each decision on a value change raises its version and notifies

Lines 657-669

Asked as: version up after approve; delete confirmation stale; decision notice queued

After an approval or a rejection, Turbo raises the change's version by 1, even when the PDF and value stay the same. A delete confirmation opened while the change was pending then becomes stale. Turbo also queues a notice of the decision (approved or rejected).

## Decision notices follow a successful approval or rejection

Lines 670-677

Asked as: notice sent after approve; updated change returned; notices after rejection

After an approval or rejection commits, Turbo sends the queued notices and returns the updated value change. A refused decision sends no notices.

## Deleting a value change needs valid action data

Lines 680-685

Asked as: delete button error; delete invalid action; admin delete fails

When an admin deletes a value change, Turbo first reads the submitted action. If the action data cannot be read, Turbo refuses with "Invalid action data." and deletes nothing.

## Only admins delete value changes

Lines 686-699

Asked as: delete value change permission; non admin delete; delete needs admin

Deleting a value change needs an admin. An attachment with no job gets 'Use abandoned-request deletion for this attachment.' The initial project approval cannot be deleted.

## Delete value change refuses attachments of project requests and the initial approval

Lines 694-699

Asked as: delete initial approval; delete project request attachment; abandoned request deletion

An admin cannot delete an attachment that has no job. Turbo refuses with "Use abandoned-request deletion for this attachment." An admin also cannot delete the initial project approval. Turbo refuses with "The initial project approval cannot be deleted."

## Delete needs the reviewed version

Lines 700-704

Asked as: delete says changed reload; confirm deletion again; delete stale version

Deleting a value change requires the version and PDF hash the admin saw. Otherwise Turbo says 'This value change changed. Reload it and confirm deletion again.'

## Deleting an approved value change

Lines 706-752

Asked as: delete approved value change; latest approval delete; pending blocks delete

An approved change can be deleted only if no pending or rejected change exists on the job ('Delete the pending or rejected value change first.'). It must be the latest approval ('Only the latest approved value change can be deleted.'). The job's value must still equal the approved value, or Turbo says 'The current project value does not match this approval.' Deleting then restores the job's previous project value.

## Approved value change deletion restores the previous project value

Lines 706-733

Asked as: project value restore after delete; previous value reverted; value does not match approval

Deleting an approved value change sets the job's project value back to the value it had before (previous_project_value). Turbo refuses with "The current project value does not match this approval." when the job's current value differs from the approved value. The restore also runs the job's workflow checks.

## Restoring a negative earlier value

Lines 735-745

Asked as: negative previous value; restore negative value; job value negative delete

If the restored previous value is negative, Turbo validates the rest of the job with value 0, then saves the exact negative value without further validation.

## Deleting a pending or rejected change

Lines 752-752

Asked as: delete pending change; remove rejected change; delete does not change value

Deleting a pending or rejected value change just removes it. The job's value is not changed.

## Missing PDF file is re-uploaded only from an identical upload

Lines 761-781

Asked as: pdf file missing; restore missing pdf; stored pdf gone; reupload identical pdf

Before a change is approved, or when a save keeps the same PDF, Turbo checks that the stored PDF file exists. If the file is missing and the upload has the same hash as the record, Turbo restores the file from the upload. Otherwise it refuses with "The authorization PDF is unavailable. Restore or replace it before approval."
