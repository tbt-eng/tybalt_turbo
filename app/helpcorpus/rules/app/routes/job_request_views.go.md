<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_request_views.go","sha256":"0dd10f2e770e5dc38f55718be9e609d1bb5ac7b4ffc8b47369a86a19780e45e3","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_request_views.go

## Request detail loads linked records for each field that is set

Lines 27-42

Asked as: request detail missing client; details page won't load; linked record not found; job owner not showing; rate sheet entries empty

A project request's detail shows the full record behind each of these fields when the field is set: client, job owner (job_owner), contact, manager, alternate manager, creator, approver, branch, rate sheet, parent job, proposal and invoicing profile (invoicing_information). Turbo skips a field that is empty. The detail also always carries a rate sheet entries list (rate_sheet_entries), which is empty when the request has no rate sheet. If a linked record cannot be found, the detail fails to load.

## Bill-to name and billing contact

Lines 43-65

Asked as: bill to name wrong; who gets billed; billing contact on request; invoice name on request

The request page shows the bill-to name chosen by the client's invoicing profile (official name or alias, via BillingName). The billing contact is the contact on the invoicing profile.

## Rate sheet entries listed highest rate first

Lines 66-84

Asked as: rates listed in wrong order; rate sheet sort order; highest rate first; rate sheet entries order

A request's rate sheet entries are listed from the highest rate down, with role, rate and overtime rate. Proposed rates are added by requestProposedRates, a function in another file.

## Person name and phone come from the profile, with fallbacks

Lines 87-96

Asked as: wrong name for approver; manager phone number missing; name shows as username; person phone blank

Turbo shows each person on a project request (manager, approver, creator and others) with their name, email and phone. When the person has a profile, the name is the profile's given name and surname joined and trimmed, and the phone is the profile's phone. When the person has no profile, the name is the user account's name and the phone is empty.

## Review token ties approval to one submission

Lines 98-143

Asked as: resubmitted request same approval; recalled request approval token; approve fails empty token; why approval rejected after recall; submission token error

The review token (a SHA-256 digest) covers the request's version, submission token, business fields, related records without approver and creator, division hours, whether each division is active, and the attachment's document hash. A recalled and resubmitted request therefore gets a new token even if the data and PDF are unchanged. A draft can have an empty submission token for viewing, but approve and create must reject an empty one.

## Request detail shows what the viewer may do

Lines 144-216

Asked as: no buttons on request; can't see actions; approval turned off buttons; request page no permissions

The project request detail page returns permission flags with the request. Every flag except the admin-delete case is false when project request approval is switched off (IsJobRequestApprovalRequired).

## Who can create the job from a request

Lines 178-209

Asked as: create job button missing; can't create job; procurement cannot create job; job creation not allowed

A person with the procurement claim (the job claim) can create the job (can_create) only when the request status is Awaiting creation.

## Who can delete a project request

Lines 182-209

Asked as: delete draft request; remove project request; delete abandoned request admin; can't delete request

The creator can delete an editable request (Draft, Recalled or Rejected) while approval is enabled. An admin (the admin claim) can also delete an abandoned request (abandonedRequest), even when approval is switched off. can_reject comes from canRejectRequest, a function in another file.

## Who can edit, submit or recall a project request

Lines 190-209

Asked as: edit request button; submit project request; recall request; can't recall request; rejected request submit

Only the creator can edit a request in Draft, Recalled or Rejected status (can_edit). The creator can submit it (can_submit) except when it is Rejected. The creator can recall it (can_recall) only when its status is Awaiting approval or Awaiting creation.

## Who can approve a project request

Lines 192-209

Asked as: approve button missing; can't approve request; approve over my limit; approver authorization limit

A person can approve (can_approve) only when the request is Awaiting approval, they are its selected approver, and their project authorization limit covers the request's project value (ProjectAuthorizationLimit returns no error).

## Proposed rates checklist wording for approvers

Lines 197-204

Asked as: proposed rates question wording; rate sheet checklist text; approve proposed rates; proposed rates checklist

When a request has proposed rates and its checklist is the approval checklist, the rate sheet question shows the proposedRatesChecklistText wording. The approver then approves the rates themselves.

## Request list queues and required claims

Lines 217-258

Asked as: request queue unknown; my requests list; creation queue empty; approval queue missing; abandoned tab not showing; all requests tab

The project request list has these queues. "mine" (the default) shows the person's own requests. "creation" needs the job claim and shows requests awaiting creation, and it also needs approval to be enabled. "approval" needs the project authorization claim and shows requests assigned to the person with status Awaiting approval. "approved" needs the same claim and shows requests the person approved that are Awaiting creation. "abandoned" needs the admin claim. "all" needs the job claim. Any other queue gives "Unknown request queue."

## Request queues need an active user

Lines 217-222

Asked as: inactive user cannot see requests; queues error not active; user not active error; request list inactive account

Every queue and the options list check that the person is active (workflowActive) first. The error comes from that function in another file.

## Abandoned queue cutoff

Lines 243-248

Asked as: abandoned requests days; how old is abandoned; JOB_REQUEST_ABANDONMENT_DAYS; requests considered abandoned

The abandoned queue lists requests whose last edit is on or before JOB_REQUEST_ABANDONMENT_DAYS days ago. The comparison includes the cutoff moment.

## Request list paging and order

Lines 260-303

Asked as: requests per page; list page size; request list sort order; newest requests first; page 2 of requests

The list defaults to page 1 and 50 per page. A page below 1 becomes 1. A perPage below 1 or above 100 becomes 50. Requests sort by newest submitted, then newest created, then id. Each row shows the client name.

## Who can load project request form options

Lines 304-319

Asked as: cannot use project requests; no access to request form; request form options denied

Only a person with the time, job, project authorization or admin claim can load the request form options. Others get "You cannot use project requests."

## Request form options: rate sheets, people, proposals, parents

Lines 320-371

Asked as: proposal missing from dropdown; parent project list; rate sheet not in dropdown; users missing from request form; subjob parent dropdown

The options leave out rate sheets that belong to a request's proposed rates. They list only active users. Proposals are jobs whose number starts with P and whose status is In Progress, Submitted or Awarded. Parent projects are jobs whose number matches the pattern two digits, a hyphen, then three or four digits, such as 25-123. Subjobs and other numbers are not listed as parents.

## Authorization approvers list by project value

Lines 380-408

Asked as: who can approve dollar value; approver list empty; whole-dollar project value; authorization limit approvers

The approvers list needs a positive whole-dollar value, or Turbo says "Enter a positive whole-dollar project value." It lists holders of the project authorization claim whose limit covers that value (ProjectAuthorizationLimit), each with the maximum value. Holders whose limit does not cover the value are left out.

## Possible duplicate jobs for a request

Lines 410-485

Asked as: possible duplicate job; duplicate check; same job twice; matching jobs found

The duplicates check lists existing jobs, except proposals (numbers starting with P), that match the request. Matches are the same client, proposal, parent, location, client PO or client reference, or a description that contains the other or is contained by it, ignoring case and spaces. It also notes when the job is the request's parent. Empty values never match.

## Pending requests in duplicates check need job claim

Lines 445-466

Asked as: duplicate requests hidden; can't see pending duplicates; other requests match

Other project requests that look like duplicates are listed only to a person with the job claim. The request itself is excluded.

## Menu badge counts for project workflow

Lines 487-531

Asked as: menu count missing; badge number not showing; awaiting approval count; project notification count

When approval is enabled, a person with the job claim sees a count of requests awaiting creation. A person with the project authorization claim sees a count of requests awaiting their approval, and a count of pending project value approvals assigned to them. Without the claim, or for an inactive user, no counts are added.
