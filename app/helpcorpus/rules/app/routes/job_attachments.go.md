<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_attachments.go","sha256":"ff076ae5fe76006ef1d732ddd2fac64a755673a337e868f8407e094038c20598","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_attachments.go

## Authorization PDF upload size and format limits

Lines 25-57

Asked as: pdf too big; file over 20 MB; upload not a pdf; wrong file type upload; only one pdf allowed

Turbo accepts one authorization PDF per upload. The file must be at most 20 MB, otherwise it says "The PDF must be at most 20 MB." The file must start with %PDF-, otherwise it says "Upload a valid PDF file." More than one file, or none where one is required, gives "Upload one authorization PDF."

## Upload is optional when a correction keeps its PDF

Lines 28-35

Asked as: no file needed on edit; edit without re-uploading pdf; keep old pdf on change; save correction without file

When the PDF is not required (an edit to a value change) and no file is sent, Turbo accepts the request without a file and keeps the existing PDF. A missing file is not an error then.

## Same PDF cannot be used on two approvals

Lines 58-76

Asked as: same pdf already used; pdf already attached elsewhere; duplicate authorization file; pdf used on another approval

Turbo rejects a PDF whose hash matches another attachment, with "This PDF already belongs to a project request or job approval." If the other attachment belongs to an abandoned project request, the message names that request ID and tells the user to give it to an admin to resolve the duplicate.

## Attaching or removing a project request's authorization PDF

Lines 77-152

Asked as: remove pdf from project request; replace project request pdf; upload pdf to request; pdf not saving on request

A person who can edit the project request can upload or remove its authorization PDF, only while the request is editable and the version they hold is current. Removing when there is no attachment does nothing. Uploading the same PDF again only restores a missing file and changes nothing else. A new PDF clears any earlier approval on the attachment, records the uploader, the upload time and the request's project value, and raises the version by 1. Either action marks the request as edited.

## Who can view an authorization PDF

Lines 153-206

Asked as: cannot open authorization pdf; can't view pdf attachment; attachment not found; who can see the pdf

For a PDF on a project request, whoever can view the request can view it. For a PDF on a job's value change, these people can view it: procurement (the job claim), admins, the job's manager, the job's alternate manager, and an approver (the project_authorization claim) who is the assigned approver or the approver of that change. Anyone else gets "You cannot view this PDF." A missing attachment gives "Attachment not found." The PDF is served with no-store caching.

## Value change list needs a job or the approval queue

Lines 207-227

Asked as: select a job error; value change list empty; approval queue shows nothing; missing rows in queue

The value change list shows only attachments tied to a job. Without the approval queue, the caller must select a job, otherwise Turbo says "Select a job." The approval queue requires the project authorization claim and shows only Pending changes assigned to the caller. Rows the caller cannot view are left out.

## Value change list contents and ordering

Lines 228-291

Asked as: value change list order; which change is current; current row value change; why rows sorted this way

Turbo orders changes by approval sequence, newest first, then by creation time. Each row carries the job number, description, current project value, related records and the names of the uploader, approver, assigned approver and rejector. A row is marked current when it is Approved and its sequence equals the job's latest approval sequence. A Pending row also carries a review token.

## Approved-by-me list covers the last 30 days

Lines 293-326

Asked as: my approvals last 30 days; list of approvals I made; approved by me report; my past approvals

An approver (the project_authorization claim) can list their own approvals of new jobs and value changes from the last 30 days (APPROVED_BY_ME_DAYS). Without the claim, Turbo refuses.

## Review token ties a value change to the job data seen

Lines 340-401

Asked as: job changed since review; stale value change token; reload before approving; review data out of date

Turbo computes a review token from the job (leaving out created, updated and the import flag), its latest approval sequence, its related records and its time allocations (division and hours). A change in any of these changes the token. Saves record the token, rejection compares it with the token of the page shown, and approval compares it with the saved token.

## Approving a value change sets the project authorization to PA

Lines 403-420

Asked as: authorizing document PA; approval sets PA; project authorization becomes PA; old project no authorizing document

An approved value change sets the job's project value to the approved value and sets its authorizing document to "PA", including on older projects with no authorizing document.

## Only procurement can save a value change

Lines 422-446

Asked as: who can create value change; can't save project value; whole dollar project value error; approver limit too low

Only procurement (the job claim) can create or edit a project value change. The project value must be a positive whole-dollar amount, otherwise Turbo says "Enter a positive whole-dollar project value." The assigned approver's project limit must cover the value, as ProjectAuthorizationLimit checks. A new value change requires a PDF.

## Value change save refusals

Lines 449-496

Asked as: cannot edit approved change; proposal job P number; second pending change error; form out of date save

Turbo refuses to edit an Approved change with "Approved evidence cannot be edited." It refuses a stale form with "This value change changed. Reload it." It refuses proposals (job numbers starting with P) with "Proposals do not use project-value approval." It refuses a second open change on a job with "Resolve the current pending value change first." The job's stored details and allocations must be valid, and the job must pass validation with the new value. Fields sent with the value change are ignored.

## Unchanged value change save keeps its state

Lines 501-514

Asked as: saved without changes; resave value change no change; rejected stays rejected on save; save did nothing

If the job data matches the saved review token, and the value, the approver and the PDF are unchanged, a save changes nothing. A rejected change stays rejected. Any real correction, here or in the job editor, needs fresh approval.

## Saving a changed value change resubmits it

Lines 515-546

Asked as: edit value change again; resubmit changed value change; notify approver after edit; revision notice to approver

A real change clears any earlier approval or rejection, raises the version by 1 and stores the new value, approver and review token. A new PDF is checked for duplicates and replaces the old one. The assigned approver is notified, and the notice says the change is a revision only if the same person was already the approver.

## Only the selected approver can approve or reject

Lines 558-575

Asked as: not my approval to review; cannot approve value change; only assigned approver reviews; reviewed value change changed

Approving or rejecting a value change needs the project authorization claim. The change must be Pending, otherwise Turbo says "Only a pending value change can be reviewed." Only the assigned approver can review it, otherwise Turbo says "Only the selected approver can review this value change." The reviewer must hold the current version and PDF hash, otherwise Turbo says "The reviewed value change changed. Reload it."

## Rejecting a value change

Lines 577-598

Asked as: reject value change reason; rejection reason required; reject needs reason text; reject job changed reload

A rejection needs a reason, otherwise Turbo says "Enter a rejection reason." The review token of the page shown must match the job's current data, otherwise Turbo says "The job or its related data changed. Reload it before rejecting." A rejection sets the status to Rejected and records the rejector, the time and the trimmed reason. The version rises by 1 and the people concerned are notified.

## Approving a value change

Lines 599-665

Asked as: approve value change fails; pdf unavailable approval; job changed ask editor; approval limit blocks approve

Approval needs the PDF file to exist, otherwise Turbo says "The authorization PDF is unavailable. Restore or replace it before approval." The approver's project limit must cover the value. If the job changed since the save, Turbo says "The job changed. Ask a job editor to save and review this value change again." Allocations and the job must pass validation. Approval sets the job's new value, records the previous value, the approver, the time, the approver's limit (approval_max_value) and the next approval sequence, and sets the status to Approved. The version rises by 1 and the people concerned are notified.

## Approval or rejection notices go out only after success

Lines 666-671

Asked as: emails sent when approval fails; notification after failed approval; no notice if save fails; alerts on failed rejection

When an approver approves or rejects a value change, Turbo sends the notices only after the whole action has saved. If any check or save fails, the action returns the error and sends no notices.

## Deleting a value change needs valid action data

Lines 674-679

Asked as: delete attachment invalid data; delete returns invalid action data; delete request body error; attachment delete body fails

When an admin deletes a value change (job attachment) and the request body cannot be read, Turbo refuses with "Invalid action data." The delete handler reads the body before it checks the admin claim.

## Only admins can delete value change evidence

Lines 681-698

Asked as: delete value change attachment; delete attachment on abandoned request; cannot delete initial approval; admin delete pdf attachment

Only an admin can delete a value change attachment. Attachments on abandoned requests need a different deletion: "Use abandoned-request deletion for this attachment." The initial project approval cannot be deleted. The admin must hold the current version and PDF hash, otherwise Turbo says "This value change changed. Reload it and confirm deletion again."

## Deleting an approved value change restores the old value

Lines 700-746

Asked as: undo approved value change; delete approval restore previous value; project value does not match approval; old project value after delete

An admin can delete an Approved value change only if the job has no pending or rejected change ("Delete the pending or rejected value change first.") and it is the latest approved one ("Only the latest approved value change can be deleted."). The job's current value must equal the approval's value, otherwise Turbo says "The current project value does not match this approval." Deleting restores the previous project value. A negative previous value is restored exactly, after the other fields are validated. Deleting a pending or rejected change only removes it.

## Missing PDF file can be recovered by identical upload

Lines 755-775

Asked as: pdf file missing; restore missing pdf; authorization pdf unavailable; reupload same pdf to fix

Turbo checks that the stored PDF exists. If it is missing, an editable upload with the same hash restores it. Otherwise Turbo says "The authorization PDF is unavailable. Restore or replace it before approval."
