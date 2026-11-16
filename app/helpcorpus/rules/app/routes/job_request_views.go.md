<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_request_views.go","sha256":"e723e970f6bbf2901027b3a9a14693605aa9248dfd95e25689af2426dc3b3c62","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_request_views.go

## Project request page shows the people and records it relates to

Lines 27-42

Asked as: project request page missing data; request fails to load; related client or contact not showing; person details on request

For a project request, Turbo loads the related client, job owner, contact, manager, alternate manager, creator, approver, branch, rate sheet, parent job, proposal and invoicing profile (invoicing_information). A person field (manager, alternate manager, creator, approver) shows the person's workflow details instead of the raw user record. If any of these records is missing, the whole request fails.

## Related records show the billing name from the invoicing profile

Lines 43-65

Asked as: bill to name wrong; billing name alias; invoicing contact missing; billing contact on request

The request's related data includes bill_to_name, which comes from utilities.BillingName using the client and the invoicing profile. The profile chooses the official name or an alias. The invoicing profile's contact is added as billing_contact.

## Rate sheet entries list highest rate first

Lines 66-85

Asked as: rate sheet order; rates listed in wrong order; proposed rates not showing; overtime rate on request

A request's rate sheet entries are listed with the highest rate first, with the role name, rate and overtime rate. Turbo then adds the request's proposed rates through requestProposedRates.

## Review token ties the review page to one submission

Lines 98-143

Asked as: review checklist no longer matches; checked checklist invalid after recall; review token mismatch; resubmit breaks approval checklist

Turbo computes a review token (SHA-256 digest) from the request's version, its submission token, its business fields, related records, division hours and the attached document hash. Turbo leaves out the approver and creator. A recall and resubmit gives a new submission token, so an old checked checklist no longer matches. An empty submission token is fine for viewing or printing a draft, but approve and create must reject it.

## Project request detail shows the permissions of the viewer

Lines 147-208

Asked as: which buttons can I see; why no approve button; what actions can I take; request detail permissions

The project request detail page tells the viewer which actions they can take: edit, submit, recall, reject, approve, create and delete. Every action except delete by an admin needs project request approval to be switched on (IsJobRequestApprovalRequired). Turbo loads the request with loadVisibleRequest, so only a person who can see it gets the detail.

## Only procurement can create a job from an approved request

Lines 178-201

Asked as: create job from request; cannot make job; turn approved request into job; who can create job

A person with the job claim (procurement) can create the job from a project request whose status is Awaiting creation (can_create).

## Who can delete a project request

Lines 182-201

Asked as: delete project request; cannot delete request; remove abandoned request; admin delete request

The creator can delete a request they can edit (Draft, Recalled or Rejected) while approval is switched on. An admin (the admin claim) can delete an abandoned request (abandonedRequest) even when approval is off.

## Only the creator can edit a draft, recalled or rejected request

Lines 190-201

Asked as: edit someone else's request; cannot edit request; change draft request; cannot submit rejected request

A person can edit a project request only when they created it and its status is Draft, Recalled or Rejected. They can submit it from those states too, except a Rejected request: it can be edited but not submitted from the detail page (can_submit).

## Only the creator can recall a submitted request

Lines 191-201

Asked as: withdraw submitted request; take back request; cancel my submission; no recall button

The creator of a project request can recall it while its status is Awaiting approval or Awaiting creation (can_recall). Nobody else gets the recall action.

## Approve needs the chosen approver whose project limit covers the value

Lines 191-201

Asked as: cannot approve request; approve button missing; approval limit too low; not the approver

A person can approve a project request only when its status is Awaiting approval, they are the selected approver, and their project authorization limit covers the project value (ProjectAuthorizationLimit returns no error). If the limit does not cover the value, Turbo does not offer approval.

## Request queues and who can open them

Lines 209-251

Asked as: request queue empty; missing approval queue; which list shows my requests; cannot open queue

Every queue requires an active user (workflowActive). The creation queue needs approval to be switched on and the job claim (procurement). The approval and approved queues need the project authorization claim and show only requests where the person is the approver, with status Awaiting approval or Awaiting creation. The abandoned queue needs the admin claim. The all queue needs the job claim. The mine queue, which is the default, shows the requests the person created.

## Abandoned queue cutoff

Lines 235-240

Asked as: abandoned requests cutoff; old requests not listed; how long before abandoned; stale requests list

The abandoned queue lists requests whose last_edited is on or before JOB_REQUEST_ABANDONMENT_DAYS days ago. A request edited exactly at the cutoff is included.

## Unknown queue is refused

Lines 248-249

Asked as: invalid queue name; queue not found error; wrong list name; unrecognised request list

A request list with any other queue name is refused with the message "Unknown request queue."

## Request list paging and order

Lines 258-287

Asked as: request list order; requests per page; page size wrong; client name missing in list

The request list is sorted by newest submitted first, then newest created. The page size is 50 by default. A perPage below 1 or above 100 falls back to 50, and a page below 1 becomes page 1. Each row shows the client name.

## Request options need a project-related claim

Lines 296-311

Asked as: cannot use project requests; no access to project request options; project request dropdowns empty; project options blocked

Staff can load the options for a project request (clients, contacts, invoicing profiles, branches, rate sheets, divisions, active users, proposals, parent projects) only with the time, job, project authorization or admin claim. Anyone else gets "You cannot use project requests."

## Options leave out request rate sheets and inactive users

Lines 313-341

Asked as: rate sheet dropdown missing; inactive user not in list; user picker empty; old rate sheets showing

The rate sheet options leave out the proposed rate sheets that belong to a request. The user options list only active users.

## Proposal and parent options

Lines 346-363

Asked as: parent project not listed; proposal not in dropdown; subjob as parent; which proposals appear

Proposals in the options start with P and have status In Progress, Submitted or Awarded. Parent project options are jobs whose number matches the pattern of two digits, a hyphen, then three or four digits, such as 25-123. Subjobs and other numbers are not offered as parents.

## Approver lookup by project value

Lines 372-401

Asked as: who can approve this amount; approver list empty; find approver for value; approver limit lookup

The approver lookup needs a positive whole-dollar value. Otherwise Turbo says "Enter a positive whole-dollar project value." It lists the holders of the project authorization claim whose limit covers that value, each with the max_value limit. Holders whose limit does not cover it are left out.

## Duplicate check on a request

Lines 402-461

Asked as: possible duplicate job; duplicate warning; same project already exists; check for duplicates

The duplicate check compares a request with existing jobs. Proposals are skipped. It also compares the request with other pending project requests, but only for staff with the job claim (procurement).

## Duplicate reasons

Lines 462-477

Asked as: why flagged as duplicate; duplicate match reason; parent project warning; same description flagged

A job or request counts as a duplicate when it has the same client, proposal, parent, location, client PO or client reference (an empty value never matches). It also counts when one description, trimmed and case-insensitive, contains the other. It is flagged "Parent project" when it is the request's parent.

## Menu badges for project workflow

Lines 479-522

Asked as: menu count wrong; badge number missing; pending approvals count; no notification count

Badge counts appear only for an active user. Procurement see the count of requests in the creation queue when approval is on. Holders of the project authorization claim see the count of requests awaiting their approval when approval is on. They also see the count of pending value-change approvals assigned to them, even when approval is off.
