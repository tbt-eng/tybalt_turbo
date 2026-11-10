<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/job_attachments.md","sha256":"d96c35f93fbb52cdde59d2a3376e85d0796befa16d5b1025162731bc9d2ac499","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/job_attachments.md

## Who can list or view job attachments

Lines 10-16

Asked as: who can see attachments; view job attachments permission; attachment access rules; why can't I see the attachment

A person must be signed in and have an active admin profile (admin_profiles_via_uid.active). For an attachment on a project request, the request's creator passes. The request's approver passes only with the project authorization claim. Anyone with the job claim passes. For an attachment on a job, anyone with the job claim or the admin claim passes. The attachment's assigned approver (assigned_approver) passes only with the project authorization claim. The view rule is the same filter as the list rule.

## Nobody can create, update or delete attachments through the API

Lines 21-29

Asked as: cannot edit attachment; delete attachment via API; api update attachment forbidden; modify attachments not allowed

The job_attachments collection has no create, update or delete rule. Nobody can do these actions through the API, except a PocketBase superuser. Turbo's own server code can still do them.

## Attachment document must be a PDF up to 20 MB

Lines 40-40

Asked as: upload pdf size limit; attachment file too large; only pdf files allowed; 20 MB upload limit

The document field is required and holds one file. The file must be a PDF (application/pdf) and at most 20971520 bytes (20 MB).

## Document hash format and uniqueness

Lines 41-41

Asked as: duplicate document upload; document already attached; same file uploaded twice; invalid document hash

The document_hash field is required and must be 64 lowercase hexadecimal characters. A unique index means two attachments cannot have the same document hash, so Turbo refuses a document that is already attached.

## Attachment status values

Lines 45-45

Asked as: attachment status options; allowed attachment states; what statuses can attachment have; attachment status missing

The status of an attachment is required and is one of Pending, Approved or Rejected.

## One attachment per project request

Lines 67-67

Asked as: second attachment on request; request already has attachment; only one file per request; attach another file to request

A unique index allows a project request to have only one attachment (idx_attachment_request). The index does not apply to attachments with no request.

## One unapproved attachment per job at a time

Lines 68-68

Asked as: cannot add attachment to job; job already has pending attachment; job attachment blocked; rejected attachment still open on job

A job can have only one attachment whose status is not Approved (idx_attachment_pending_job). A job needs its Pending or Rejected attachment to be resolved before another one is added.

## Approval sequence is unique per job for approved attachments

Lines 69-69

Asked as: duplicate approval sequence; approval number already used; same sequence on two approved; approval sequence conflict

Two approved attachments on the same job cannot have the same approval_sequence (idx_attachment_sequence). The index applies only to attachments with a job and status Approved.
