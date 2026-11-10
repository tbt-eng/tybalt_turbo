<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_attachments.go","sha256":"5e7141aba0e9b6c8fd4d498bb8c9ab804917c0b198d73507fd3c35ad9de2ab9e","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_attachments.go

## Authorization PDF upload size and format limits

Lines 24-56

Turbo accepts at most one authorization PDF per upload, and the file must be at most 20 MB. Otherwise it says "The PDF must be at most 20 MB." The file must start with "%PDF-"; otherwise Turbo says "Upload a valid PDF file." A missing or extra file gives "Upload one authorization PDF."

## Optional PDF on corrections

Lines 27-34

When an upload is not required (editing a value change), Turbo accepts a form with no PDF and keeps the existing one. Creating a value change requires a PDF.

## Duplicate authorization PDF is refused

Lines 57-75

Turbo identifies each PDF by its SHA-256 hash and refuses a PDF already attached to another project request or job approval: "This PDF already belongs to a project request or job approval." If the other copy belongs to an abandoned project request, the message names that request ID and tells the user to give it to an admin to resolve the duplicate.

## Add, replace or remove PDF on a project request

Lines 76-151

Only a project request that the user can see and edit, at the current version, accepts a PDF change. Removing a PDF when none exists does nothing but still marks the request edited. Uploading the same PDF again only checks that the stored file exists. A new PDF replaces the old one, clears any approval on the attachment, records the uploader, upload time and the request's project value, and increases the attachment version by 1.

## Who can view an authorization PDF or attachment

Lines 152-199

A user must be active in the workflow. For a project request attachment, anyone who can view that request can view it. For a job value change, staff with the job claim or admin claim can view it. The job's manager and alternate manager can view it. A holder of the project authorization claim can view it only if they are the assigned approver or the approver. Otherwise Turbo says "You cannot view this PDF."

## PDF download is not cached

Lines 200-207

Turbo serves an authorization PDF with Cache-Control "private, no-store". A missing attachment gives "Attachment not found."

## Value change list: queue or job required

Lines 209-247

The list of value changes needs an active workflow user. With queue=approval it needs the project authorization claim and shows only Pending changes assigned to the user. Otherwise a job must be given, or Turbo says "Select a job." Only changes the user can view appear. They are sorted by newest approval sequence, then newest creation.

## Value change list row details

Lines 248-290

Each row gives the PDF link, the job number, description, current project value, related records and names of the uploader, approver, assigned approver and rejector. A row is marked current only if it is Approved and its approval sequence matches the job's latest approval. Pending rows also carry a review token.

## Value change checks use the job's stored time allocations

Lines 295-305

When Turbo checks a value change, it reads the job's saved time allocations (division and hours), sorted by division. It does not use allocations sent with the value change.

## Review token ties approval to job data

Lines 307-368

The review token is a hash of the job (without created, updated and import flag), the latest approval sequence, related records, and time allocations (division and hours only). Any change to these changes the token. A save records it, a rejection compares it with the page shown, and an approval compares it with the saved token.

## Approved value sets authorizing document to PA

Lines 370-376

An approved value change sets the job's project value and sets its authorizing document to "PA", even on older projects with no authorizing document.

## Value change saved only by job editors

Lines 389-413

Creating or editing a value change needs the job claim (procurement). The value must be a positive whole-dollar amount, or Turbo says "Enter a positive whole-dollar project value." The assigned approver's project limit must cover the value, checked by ProjectAuthorizationLimit.

## Approved value change cannot be edited

Lines 423-436

Editing a value change that is Approved, or an attachment without a job, is refused: "Approved evidence cannot be edited." A stale form version gives "This value change changed. Reload it."

## Proposals have no value approval

Lines 437-443

A value change on a job whose number starts with P is refused: "Proposals do not use project-value approval."

## Stored job details must be valid

Lines 444-456

Saving a value change checks the job's stored time allocations and the job itself as if the new value were approved. It ignores any job fields or allocations sent with the request, which are edited in the job editor.

## Only one open value change per job

Lines 457-463

A job cannot have a second value change that is not Approved (Pending or Rejected). Turbo says "Resolve the current pending value change first."

## Unchanged save keeps state

Lines 468-481

Saving a value change with the same job data, value, assigned approver and PDF (or no new PDF) changes nothing. The version stays and a rejected change stays rejected. Any change to the job data or related records needs fresh approval.

## Saving a change resubmits for approval

Lines 482-508

A real change clears the earlier approval, increases the version by 1, records the review token, and saves the new value, approver and PDF (uploader and time). It notifies the approver as "submitted". The notice says it is a revision only if the assigned approver is the same.

## Saving a value change sends notices only after it succeeds

Lines 510-515

Turbo sends the notice about a submitted or revised value change only after the save has fully succeeded. If any check refuses the save, no notice goes out. On success it returns the saved value change.

## Approve or reject refuses unreadable request data

Lines 519-522

When Turbo cannot read the data sent with an approval or rejection of a value change, it refuses the action. It shows the message "Invalid action data." Nothing is reviewed.

## Only the selected approver reviews a value change

Lines 525-538

Approving or rejecting needs the project authorization claim. Only a Pending value change can be reviewed: "Only a pending value change can be reviewed." Only the assigned approver may do it: "Only the selected approver can review this value change."

## Review needs the version and PDF seen

Lines 539-542

A review fails if the version or PDF hash differs from what the reviewer saw: "The reviewed value change changed. Reload it."

## Rejecting a value change

Lines 544-565

Rejection needs a non-blank reason, or Turbo says "Enter a rejection reason." The job data must match the page shown, or Turbo says "The job or its related data changed. Reload it before rejecting." Then the status becomes Rejected, and the rejector, time, trimmed reason and token are recorded.

## Approving a value change

Lines 566-616

Approval needs the PDF file to exist, otherwise "The authorization PDF is unavailable. Restore or replace it before approval." The approver's limit must cover the value. The job data must match the saved token, otherwise "The job changed. Ask a job editor to save and review this value change again." The job is validated, then updated with the new value. The change records the previous value, the next approval sequence (latest plus 1, or 1), the approver, time and limit, and becomes Approved.

## Decision bumps version and notifies

Lines 618-630

Approval or rejection increases the attachment version by 1, so open delete confirmations become stale. It queues an "approved" or "rejected" notice.

## Approval or rejection notices are sent only after success

Lines 634-638

Turbo sends the notice about an approved or rejected value change only after the decision is saved. If any check refuses the decision, no notice goes out. On success it returns the updated value change.

## Deleting a value change refuses unreadable request data

Lines 643-646

When Turbo cannot read the data sent with a delete of a value change (job attachment), it refuses the action. It shows the message "Invalid action data." Nothing is deleted.

## Only admins delete value changes

Lines 647-665

Deleting a value change needs the admin claim. An attachment without a job is refused: "Use abandoned-request deletion for this attachment." The initial project approval cannot be deleted. The version and PDF hash must match what the admin saw, or Turbo says "This value change changed. Reload it and confirm deletion again."

## Deleting a pending or rejected value change

Lines 666-666

An admin can delete a Pending or Rejected value change. It is just removed, and the job's value does not change.

## Deleting an approved value change restores the old value

Lines 667-713

Deleting an Approved value change first needs every Pending or Rejected change on the job to be deleted: "Delete the pending or rejected value change first." Only the latest approved change can be deleted: "Only the latest approved value change can be deleted." The job's current value must equal the approval's value, otherwise "The current project value does not match this approval." The job is restored to the previous project value.

## Negative restored value is kept

Lines 696-705

If the restored previous value is negative, Turbo validates the other job fields with a value of 0, then saves the exact negative value without validating it.

## Missing PDF file can be recovered by identical upload

Lines 722-741

Turbo checks that the stored PDF exists. If it is missing and the user uploaded an identical PDF (same hash), Turbo re-uploads it. Otherwise it says "The authorization PDF is unavailable. Restore or replace it before approval."
