<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_request_views.go","sha256":"a8c6a8833e962b8d043e8c7f50c0ea15288a39703c0336e99a01026d0247181a","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_request_views.go

## Request detail lists related records and their people

Lines 26-41

The project request detail page (requestRelated) loads the related client, job owner, contact, manager, alternate manager, creator, approver, branch, rate sheet, parent, proposal and invoicing profile. A related person shows the given name and surname from their profile, plus their email and phone. When a person has no profile, the user name is shown and the phone is empty.

## Request detail shows the invoicing contact

Lines 42-50

When the request has an invoicing profile with a contact, the detail adds that contact as the billing contact (billing_contact).

## Bill-to name follows the invoicing profile

Lines 51-64

The bill-to name on a request comes from the client and the invoicing profile (BillingName), which chooses the official name or the alias.

## Request detail lists rate sheet rates by role

Lines 65-79

When the request has a rate sheet, the detail lists each rate sheet entry sorted by role. Each entry shows the role name, the rate and the overtime rate. A request with no rate sheet gets an empty list.

## Person details fall back to the user record without a profile

Lines 82-91

Wherever a project request shows a person (manager, approver, creator, and similar), Turbo builds the name from the person's profile as given name plus surname, trimmed (profiles). Phone comes from the profile too. When the person has no profile, Turbo shows the user record's name and an empty phone. Email always comes from the user record.

## Review token ties approval to one submission

Lines 93-138

The review token is a SHA-256 digest of the request version, submission token, business fields, related records (without approver and creator), division hours and active flags, and the attachment's document hash. A new submission changes the token, so an old review page's checked checklist is not accepted. A draft may have an empty submission token for viewing, but approve and create reject an empty one.

## Review token excludes approver and creator

Lines 101-137

The review token digest leaves out the approver and the creator, but it includes the version, the submission token, the business fields, the related records, the hours and active state of each division allocation, and the attachment document hash. A change to any of these changes the token.

## Request detail permissions need approval to be enabled

Lines 165-192

The request detail page gives every permission (edit, submit, recall, reject, approve, create) only when job request approval is required (IsJobRequestApprovalRequired). Delete is the exception: it needs only the admin claim and an abandoned request.

## Who can create a job from a request

Lines 169-192

A person with the job claim (procurement) can create the job from a request whose status is Awaiting creation. Approval must be enabled.

## Only admins can delete abandoned requests

Lines 173-192

The detail view allows delete only for an admin and only when the request counts as abandoned (abandonedRequest). Reject permission comes from canRejectRequest and also needs approval enabled.

## Who can edit or submit a project request

Lines 181-192

Only the creator of a project request can edit it, and only while its status is Draft, Recalled or Rejected. Submitting follows the same rule but is not allowed for a Rejected request (so a rejected request must be edited first). Both need project request approval to be enabled (IsJobRequestApprovalRequired).

## Who can recall a project request

Lines 182-192

Only the creator can recall a project request, and only while its status is Awaiting approval or Awaiting creation. Approval must be enabled for the setting.

## Who can approve a project request

Lines 183-192

A person can approve only when the status is Awaiting approval, they are the selected approver, and their project authorization limit covers the request's project value (ProjectAuthorizationLimit returns no error). Approval must be enabled.

## Request detail includes the current checklist

Lines 184-192

The request detail returns the current checklist version and questions for the request's scope, plus the attachment's document link.

## Request queues and who may open them

Lines 200-234

The creation queue needs approval enabled and the job claim. The approval queue needs the project authorization claim and shows requests assigned to you with status Awaiting approval. The abandoned queue needs the admin claim. The all queue needs the job claim. The mine queue (the default) shows requests you created. Any other queue name gives "Unknown request queue."

## Mine queue is the default

Lines 214-231

An empty queue or 'mine' shows the requests the user created, with no claim needed. The approval queue shows only the requests where the user is approver and the status is Awaiting approval.

## Abandoned queue cutoff

Lines 219-224

The abandoned queue lists requests whose last edit is on or before the cutoff of JOB_REQUEST_ABANDONMENT_DAYS days ago.

## Request list paging and sort order

Lines 242-258

The request list sorts by newest submitted, then newest created. The page defaults to 1 when it is below 1. The page size (perPage) defaults to 50 when it is below 1 or above 100.

## Request list reply fields and client name

Lines 259-271

Each listed request carries the client name when the client is found. The reply gives totalItems, page, perPage, totalPages and whether approval is enabled.

## Project request options need a qualifying claim

Lines 280-318

Only staff with the time, job, project_authorization or admin claim can load the project request options. Others get "You cannot use project requests." The list includes only active users.

## Options list shows only active users

Lines 296-318

The project request options include only active users. Each person has the name, email and phone from their profile. The options also list all clients, contacts, invoicing profiles, branches, rate sheets and divisions.

## Which proposals and parent projects are offered

Lines 325-340

Proposals (numbers starting with P) are offered only when their status is In Progress, Submitted or Awarded. Parent projects are jobs whose number matches the pattern of two digits, a dash, then 3 or 4 digits (such as 25-123). Subjobs and other numbers are not offered as parents.

## Authorization approvers must cover the project value

Lines 349-376

The approver list needs a positive whole-dollar value, otherwise "Enter a positive whole-dollar project value." It lists only staff with the project authorization claim whose limit covers the value, and shows each person's max_value.

## Approver lookup needs a positive whole-dollar value

Lines 354-376

The list of authorization approvers needs a positive whole-dollar value. Otherwise Turbo says 'Enter a positive whole-dollar project value.' It lists only people with the project_authorization claim whose limit covers the value, each with their max value. People whose limit does not cover it are left out.

## Duplicate check scope

Lines 390-435

The duplicate check skips proposals (numbers starting with P). Only staff with the job claim also see other pending project requests, and the request itself is excluded.

## Duplicate check skips proposals and ignores empty fields

Lines 390-413

The duplicate check against jobs skips proposals (numbers starting with P). Fields match only when they are non-empty and equal. Matching jobs appear with the reasons, the parent number, the client name and the manager name.

## Pending request duplicates need the job claim

Lines 414-435

The duplicate check lists other pending requests that look like this one only for a user with the job claim. It leaves out the request itself.

## Duplicate check reasons

Lines 439-454

A request matches a job or another request for the same client, proposal, parent, location, client PO or client reference (only when set). It also matches a similar description (case-insensitive, either contains the other) and a Parent project when the request's parent is that job.

## Description similarity ignores case and spaces

Lines 446-452

Two descriptions count as similar when both are non-empty and one contains the other. The check ignores case and surrounding spaces. A job whose id is this request's parent also gets the reason 'Parent project'.

## Menu badge counts for project workflow

Lines 456-498

Staff with the job claim see a count of requests awaiting creation when approval is enabled. Staff with the project authorization claim see a count of requests awaiting their approval (when enabled) and pending value change approvals assigned to them. Inactive users get no badges.

## Menu badges skip inactive users

Lines 456-459

The project workflow menu badges are not added when the user is not an active workflow user (workflowActive fails). No error is shown.

## Creation and approval badges need approval enabled

Lines 460-498

The creation badge count shows for users with the job claim only when approval is required. The approval badge count shows for users with the project authorization claim only when approval is required. The value approvals badge counts pending attachments assigned to the user with a job, and it does not depend on that setting.
