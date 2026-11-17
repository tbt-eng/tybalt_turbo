<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_request_views.go","sha256":"35d6b9fd583b14b4f3dab28978e3abcdb82db426d8ca3ff348d3c9fd1c4da59e","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_request_views.go

## Request detail shows names of related records

Lines 27-51

Asked as: request shows user names; blank contact name on request; billing contact missing; names missing on project request

The project request detail page loads the client, job owner, contact, manager, alternate manager, creator, approver, branch, rate sheet, parent, proposal and invoicing profile of the request. People (users) appear with their profile given name, surname and phone, and the name is the profile's given name plus surname. The invoicing profile's contact is added as the billing contact.

## Bill-to name follows invoicing profile

Lines 43-65

Asked as: wrong bill to name; billing name on invoice; bill-to name incorrect; alias instead of client name

A request's bill-to name is chosen by utilities.BillingName from the client and the invoicing profile, as official name or alias. This file does not state the result. The billing contact comes from the invoicing profile's contact.

## Request rate sheet entries order

Lines 66-80

Asked as: rates listed in wrong order; highest rate first request; rate sheet order mixed up

The rate sheet entries shown for a request are sorted by highest rate first, and include each role's name, rate and overtime rate.

## Request rates include proposed rates

Lines 78-84

Asked as: proposed rates missing from request; rate list incomplete on request; proposed rates not showing

After the rate sheet entries, the request's proposed rates are added to the rate list (requestProposedRates), and the result is returned as rate_sheet_entries. Each entry gives the role, rate, overtime rate and updated time.

## Person details fall back to user record

Lines 87-96

Asked as: person name wrong on request; phone blank for person; name shows from user record; profile name not used

A person in a project request shows the name and email from their user record, with an empty phone. When the person has a profile, Turbo replaces the name with the profile's given name and surname and fills in the phone.

## Review token ignores approver and creator details

Lines 106-150

Asked as: review token changed; token different after edit; why review token changes; token changed after resubmit

The review token is a SHA-256 digest of the request's version, submission token, related records, business fields, division allocations (hours and whether each division is active) and the attachment's document hash. It leaves out the approver and creator details, the proposed-rates approver contact and the source rate sheet record. Changing any of the included data, or recalling and resubmitting (which gives a new submission token), gives a new token.

## Review token includes allocations and document hash

Lines 120-149

Asked as: allocation change changes token; document change review token; hours edit token different

The review token covers the request's version, submission token, business fields, related records, each division allocation (division, hours, whether the division is active), and the attachment's document hash. A change to any of these changes the token.

## Request detail needs visible request

Lines 154-160

Asked as: cannot open request; request not found; no access to request; request will not load

The request detail loads the request only if the person is allowed to see it (loadVisibleRequest). Otherwise Turbo returns that function's error.

## Request detail shows division code and name for allocations

Lines 169-172

Asked as: division number showing on request; allocation shows division id; division name missing on request; hours allocation division wrong

When a project request's detail is shown, Turbo expands each hours allocation's division. The page shows the division's code and name instead of its id. If the expansion fails, the detail request fails.

## Request detail includes review token, checklist and status flag

Lines 177-208

Asked as: checklist missing on request; all buttons greyed out; approval turned off buttons; no review token showing

The request detail returns a review token, and the current request checklist (version and questions) for the request's scope. It also returns whether project request approval is enabled (IsJobRequestApprovalRequired). When approval is not enabled, every permission (edit, submit, recall, reject, approve, create) is false, and delete is allowed only to an admin on an abandoned request.

## Project request detail: who can create the project

Lines 185-208

Asked as: create project button missing; make project from request; who can create project; cannot create project from request

A person with the job claim (procurement) can create the project (can_create) from a request whose status is Awaiting creation. Project request approval must be enabled.

## Project request detail: who can delete

Lines 189-208

Asked as: delete project request; cannot delete request; remove abandoned request; delete draft request button

The creator can delete a request they can edit (Draft, Recalled or Rejected) while approval is enabled. An admin (the admin claim) can also delete a request that is abandoned (abandonedRequest), even when approval is disabled.

## Project request detail: reject permission

Lines 193-208

Asked as: reject button missing; cannot reject request; send back request option

The reject button (can_reject) shows only if approval is enabled and canRejectRequest allows the current user. This file does not state who canRejectRequest allows.

## Project request detail: who can edit, submit, recall

Lines 197-208

Asked as: edit button missing on request; submit request not showing; recall request option missing; cannot edit project request

On a project request's detail, the creator can edit it (can_edit) while its status is Draft, Recalled or Rejected. The creator can submit it (can_submit) in the same cases, except Rejected. The creator can recall it (can_recall) while it is Awaiting approval or Awaiting creation. All of these need project request approval to be enabled (IsJobRequestApprovalRequired).

## Project request detail: who can approve

Lines 199-208

Asked as: approve button missing; cannot approve project request; approval limit too low; approver cannot approve

A person can approve a project request (can_approve) only if the request is Awaiting approval, they are its chosen approver, and their project authorization limit covers the request's project value (ProjectAuthorizationLimit returns no error). Project request approval must be enabled.

## Attachment link on request detail

Lines 204-207

Asked as: attachment link on request; download request attachment; request document link broken

When a request has an attachment, the detail page returns its link as /api/job_attachments/{id}/document.

## Request queues need claims

Lines 216-252

Asked as: cannot open queue; queue access denied; missing queue permission; request queue refused

The creation queue needs the job claim and approval to be enabled. The all queue needs the job claim. The approval and approved queues need the project authorization claim. The abandoned queue needs the admin claim. Without the claim the request is refused. The inactive-user check (workflowActive) applies to every queue.

## Request queue needs active workflow user

Lines 216-229

Asked as: queue not loading for inactive user; workflow access refused; queue blocked for user

Every request queue first checks that the user may use the workflow (workflowActive). The creation queue also requires project request approval to be enabled (requestEnabled).

## Request queue contents

Lines 230-254

Asked as: what is in approval queue; approved queue empty; my requests list; all requests queue contents

The approval queue lists requests assigned to you that are Awaiting approval. The approved queue lists requests you approved that are Awaiting creation. The mine queue (the default, also used for an empty queue) lists requests you created. The all queue lists every request. The creation queue uses creationQueueFilter.

## Abandoned requests queue cutoff

Lines 242-247

Asked as: abandoned requests list; old requests abandoned; how many days abandoned; requests not edited for days

The abandoned queue lists requests whose last_edited is on or before JOB_REQUEST_ABANDONMENT_DAYS days ago (the cutoff is that many days before now, UTC).

## Queue default is the person's own requests

Lines 248-254

Asked as: empty queue shows my requests; default request list; requests I created only

An empty queue name or 'mine' lists the requests the person created. The 'all' queue lists every request and needs the job claim (procurement).

## Unknown request queue refused

Lines 255-256

Asked as: unknown request queue error; queue name not recognised; invalid queue name

A queue name that is not creation, approval, approved, abandoned, all, mine or empty gets the message "Unknown request queue."

## Request list paging and order

Lines 265-294

Asked as: request list page size; requests sorted wrong; how many requests per page; client name missing in list

The request list is sorted by submitted (newest first), then created, then id. The page defaults to 1 if below 1. The page size (perPage) defaults to 50 when below 1 or above 100. Each row includes the client name when the client is found.

## Request options need a work claim

Lines 305-318

Asked as: cannot use project requests; project request form not opening; no access to request form

Only someone with the time, job, project_authorization or admin claim can load the project request form options. Anyone else gets "You cannot use project requests."

## Request form option lists

Lines 320-348

Asked as: dropdown options missing on request; clients missing from request form; branch list empty on request

The form options list all clients, contacts, invoicing profiles, branches and divisions. Rate sheets that belong to a request (proposed rates) are left out. Only active users appear as people.

## Request options list excludes proposed rate sheets

Lines 320-333

Asked as: proposed rate sheet not in dropdown; rate sheet missing from request form; rate sheet choices wrong

The request form's rate sheet choices leave out rate sheets that belong to a request (their proposed rates). All other lists (clients, contacts, invoicing profiles, branches, divisions) return every record.

## Only active users offered as people

Lines 334-348

Asked as: inactive user not in list; person missing from request dropdown; people list on request form

The request form's people list includes only active users, shown with their profile name and phone.

## Proposal and parent project choices

Lines 353-370

Asked as: proposal not in list; parent project not showing; no proposals offered; which jobs can be parent

Proposals offered are jobs whose number starts with P and whose status is In Progress, Submitted or Awarded. Parent projects offered are jobs whose number matches two digits, a dash, then 3 or 4 digits (such as 25-123), so subjobs are not offered as parents.

## Parent project number pattern

Lines 355-368

Asked as: subjob cannot be parent; why job is not parent; job number not accepted as parent

A job counts as a possible parent only if its number is two digits, a hyphen, then three or four digits (such as 25-123), and does not start with P. Subjobs do not qualify. Proposals are numbers starting with P with status In Progress, Submitted or Awarded.

## Request options report whether approval is required

Lines 371-375

Asked as: approval required flag; approval setting on request form; is approval turned on

The project request options response includes an enabled flag. It is set from the setting that says whether project request approval is required (IsJobRequestApprovalRequired).

## Authorization approvers list needs an active workflow user

Lines 379-383

Asked as: approver list empty; cannot see approvers; no approvers available

The list of authorization approvers for a project value is available only to users who pass the workflow active check (workflowActive). Other users get that check's refusal and no list.

## Authorization approvers list by project value

Lines 384-405

Asked as: who can approve this amount; approver for project value; whole dollar value error; approver limit too low for value

The approver list needs a positive whole-dollar value, otherwise it says "Enter a positive whole-dollar project value." It lists people with the project_authorization claim whose limit covers that value, and shows each person's max_value. People whose limit does not cover the value are left out.

## Duplicate project check reasons

Lines 415-483

Asked as: possible duplicate job warning; why duplicate warning; same description warning; duplicate project reasons

A request is compared with each non-proposal job (numbers starting with P are skipped). The reasons are: same client, proposal, parent, location, client PO or client reference (only when the value is not empty), a similar description, and parent project. Similar description means that, after trimming and ignoring case, one non-empty description contains the other. Parent project means the job is the request's parent.

## Duplicate check skips proposals

Lines 420-443

Asked as: duplicate warning shows proposal; duplicate matches list parent; duplicate check shows wrong jobs

The duplicate check compares the request to existing jobs but skips jobs whose number starts with P. Each match lists its parent number, client name and manager name.

## Pending request duplicates only for procurement

Lines 444-466

Asked as: pending duplicate requests not showing; other requests for same job; see duplicate pending requests

Only a person with the job claim (procurement) also sees other pending project requests that duplicate this one. The request itself is left out.

## Duplicate description match ignores case

Lines 469-483

Asked as: duplicate description case; similar description flagged; description match ignoring capitals

A description is similar when, after trimming and lowercasing, one contains the other. Both descriptions must be non-empty. Empty fields never count as the same client, proposal, parent, location, client PO or client reference. A job that is the request's parent gets the reason 'Parent project'.

## Workflow badge counts

Lines 486-529

Asked as: notification count requests; workflow badge number; count of approvals waiting; badge count wrong

Badge counts are skipped for inactive users. A person with the job claim sees the count of requests awaiting creation, only if approval is enabled. A person with the project_authorization claim sees the count of requests awaiting their approval (only if approval is enabled) and the count of pending value change approvals assigned to them.

## No workflow badges for inactive users

Lines 486-528

Asked as: badge missing for inactive user; no badge showing; workflow badge hidden

If the user cannot use the workflow, no project workflow badges are added and no error is returned. The creation badge needs approval enabled and the job claim. The approval badge needs approval enabled and the project authorization claim. The value approval badge counts pending job attachments assigned to the person, and appears without approval enabled.
