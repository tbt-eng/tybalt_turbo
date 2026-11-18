<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_request_views.go","sha256":"0c562bbc5fb42aec6f406980cab88a70ac61494474df7a45d8e159ce4bf6a5be","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_request_views.go

## Project request detail shows related records with people's profile names

Lines 27-51

Asked as: name shows wrong on request; contact name missing; people shown as ids; billing contact on request; profile phone on request

The project request detail page loads the related records of the request: client, job owner, contact, manager, alternate manager, creator, approver, branch, rate sheet, parent, proposal and invoicing profile. For each person it shows the name from the profile (given name and surname), with the profile phone, when a profile exists. When the invoicing profile has a contact, the page also shows that contact as the billing contact (billing_contact).

## Bill-to name follows the invoicing profile

Lines 52-65

Asked as: bill to name wrong; invoice name alias; client name on invoice; official name or alias

A project request shows a bill-to name (bill_to_name) for its client. Turbo works it out with BillingName from the client and the invoicing profile, which chooses the official name or the alias.

## Rate sheet entries listed highest rate first

Lines 66-80

Asked as: rate sheet order; rates sorted wrong; highest rate first; rate list order

A project request shows its rate sheet's entries ordered by rate, highest first. Each entry has its role, rate and overtime rate.

## Project request detail lists proposed rates and the rate sheet's entries

Lines 66-85

Asked as: rates missing on request; proposed rates not showing; rate sheet empty on request; rate sheet blank

The detail page lists each entry of the request's rate sheet with its role, rate and overtime rate. The page then adds the request's proposed rates through requestProposedRates, which is defined in another file. The page always returns a rate_sheet_entries list, which is empty when the request has no rate sheet.

## Review token needs a submission token

Lines 98-105

Asked as: approve fails empty token; cannot approve draft; submission token missing; approve rejected no token

The review token includes a random submission token (submission_token) that only submit issues. A draft can have an empty one for viewing or printing. Approve and create must reject an empty one.

## Review token ignores display-only details

Lines 106-137

Asked as: approval stale after rename; rename makes approval stale; profile edit resets approval; approval reset after edit

The review token (a digest of the request) leaves out the approver, the creator, the rates approver's contact details, the source sheet's record and the missing-role lists. It also leaves out the rate sheet's code, name and updated time. A rename or a profile edit does not make an approval stale.

## Review token changes when approved content changes

Lines 138-167

Asked as: approval lost after edit; checklist no longer matches; changed request needs reapproval; approval cleared after change

The review token covers the version, the submission token, the business fields, the rate sheet's other fields and its rates, the division hours and active flags, and the attachment's document hash. A change to any of these makes the review token differ, so an earlier approval or checklist no longer matches.

## Project request detail returns the review token and checklist

Lines 172-202

Asked as: where is checklist; checklist not showing; attachment link missing; request detail shows nothing

The detail page returns the request, its allocations (with the division expanded), its attachment (with a document_url), the related records, the review token, and the current checklist for the request's scope with its version and questions. The page also returns whether project request approval is required (enabled). Turbo shows the detail only if the request is visible to the person (loadVisibleRequest, defined in another file).

## Project request actions are all off when approval is not required

Lines 199-226

Asked as: buttons greyed out; no buttons on request; approval not turned on; all actions disabled

On the detail page, the permissions to edit, submit, recall, reject, approve and create are true only when project request approval is required (enabled). The delete permission also holds for an admin when the request is abandoned (abandonedRequest, defined in another file), whatever the setting.

## Project request detail: who can create the job

Lines 203-226

Asked as: create job button missing; make job from request; turn request into job; cannot create job

A person with the job claim (procurement) can create the job from a project request, only when its status is Awaiting creation and approval is required.

## Project request detail: who can delete

Lines 207-226

Asked as: cannot delete request; remove project request; delete abandoned request; no delete button

The creator can delete a project request they can edit, when approval is required. An admin can delete an abandoned request (abandonedRequest) even if approval is not required.

## Project request detail: reject permission

Lines 211-226

Asked as: no reject button; reject option missing; cannot reject request

Whether the Reject button shows depends on canRejectRequest for the current user, and only when approval is required. This file does not state who passes that check.

## Project request detail: who can edit or submit

Lines 215-226

Asked as: edit button missing draft; cannot submit request; send for approval missing; cannot change request

On a project request's detail page, the creator can edit it when its status is Draft, Recalled or Rejected, and only when project request approval is required (enabled). The creator can submit it in the same cases, except a Rejected request.

## Project request detail: who can recall

Lines 216-226

Asked as: take back request; withdraw request; recall button missing; pull back submitted request

The creator can recall a project request while its status is Awaiting approval or Awaiting creation. This applies only when project request approval is required.

## Project request detail: who can approve

Lines 217-226

Asked as: approve button missing; approver cannot approve; over my limit approve; cannot approve amount

Only the selected approver can approve a project request, and only while its status is Awaiting approval. The approver's project authorization limit must also cover the request's project value (project_value), or Turbo does not offer approval.

## Request lists need an active user

Lines 234-236

Asked as: unknown queue error; request list empty no user; inactive user request list; queue not recognised

Every project request list needs an active user. Turbo rejects an unknown queue name with "Unknown request queue."

## Creation queue

Lines 240-247

Asked as: creation queue empty; requests waiting to be created; creation list missing

The creation queue (queue=creation) needs the job claim (procurement) and project request approval to be on. It lists requests that match creationQueueFilter.

## Approval queue

Lines 248-252

Asked as: requests to approve list; my approvals missing; approval list empty

The approval queue needs the project authorization claim. It lists requests where the user is the approver and the status is Awaiting approval.

## Approved queue

Lines 253-259

Asked as: requests I approved; approved list empty; approved requests missing

The approved queue needs the project authorization claim. It lists requests that this person approved and that have the status Awaiting creation.

## Abandoned queue

Lines 260-265

Asked as: old stale requests; abandoned requests list; requests not touched for days; inactive old requests

The abandoned queue needs the admin claim. It lists requests last edited on or before the cutoff, which is JOB_REQUEST_ABANDONMENT_DAYS days ago.

## All queue

Lines 266-270

Asked as: all requests list; see every request; full request list

The all queue needs the job claim and lists every project request.

## Mine queue is the default

Lines 271-272

Asked as: my requests; default list empty; requests I created; own requests list

With no queue given, or queue=mine, Turbo lists the project requests that the user created. No claim is needed.

## Request list paging

Lines 283-312

Asked as: requests page size; request list order; next page requests; newest requests first

The request list is sorted by newest submitted, then newest created. The page size is 50 unless perPage is between 1 and 100. A page below 1 becomes page 1.

## Request list shows client name

Lines 301-307

Asked as: client name missing in list; list shows no client; client column blank

Each item in the request list carries its client's name (client_name) when the client can be found.

## Project request options need a claim

Lines 323-336

Asked as: cannot use project requests; no access to project requests; project request options denied

Only a person with the time, job, project_authorization or admin claim can load the project request options, and the user must be active. Anyone else gets "You cannot use project requests."

## Options hide proposed-rate sheets and inactive users

Lines 338-366

Asked as: rate sheet missing from dropdown; proposed rate sheet not listed; inactive user not in list; user missing from options

The options list rate sheets that belong to no request, so a request's proposed rates are never listed. The user list holds only active users.

## Options list proposals and parent projects

Lines 373-388

Asked as: parent project dropdown; proposal not in list; subjob not showing as parent; proposals missing from list

Proposals are jobs whose number starts with P and whose status is In Progress, Submitted or Awarded. Parent projects are jobs whose number matches two digits, a hyphen, then 3 or 4 digits, such as 25-123. Subjobs are not listed as parents.

## Project request options list all clients, contacts, branches and divisions

Lines 389-393

Asked as: dropdown missing client; contact list empty; division missing from options; branch not listed

The options list returns every client, contact, invoicing profile, branch, division and rate sheet (rate sheets exclude proposed ones). It also returns whether approval is required (enabled).

## Approvers list needs an active user

Lines 397-401

Asked as: approvers list blank; inactive approvers list; no approvers shown

The list of approvers (authorization approvers) refuses an inactive user through workflowActive, which is defined in another file.

## Approvers list needs a whole-dollar value

Lines 402-405

Asked as: project value error; enter positive value; approver value decimal; whole dollar amount required

The list of authorization approvers needs a positive whole-dollar project value. Otherwise Turbo says "Enter a positive whole-dollar project value."

## Approvers list shows only those with enough limit

Lines 406-424

Asked as: who can approve amount; approver max value; approvers not showing limit; approver limit too low

The approvers list holds the people with the project authorization claim whose limit works for the value. It shows each person's maximum value (max_value). Turbo skips anyone for whom the limit check fails.

## Possible duplicate jobs for a request

Lines 433-502

Asked as: possible duplicate warning; same client PO job; duplicate job check; similar description warning

The duplicates check lists existing jobs, but not proposals (numbers starting with P), that share the request's client, proposal, parent, location, client PO or client reference. The match on each field is exact, and an empty field never matches. It also flags a similar description when either description contains the other, ignoring case and surrounding spaces. It adds "Parent project" when the job is the request's parent.

## Pending duplicate requests only for procurement

Lines 462-483

Asked as: duplicate requests warning; similar requests other people; duplicate request alert

Only a person with the job claim (procurement) sees other project requests that look like duplicates. Turbo leaves out the request itself.

## Menu badges for project workflow

Lines 504-547

Asked as: badge count missing; approval count menu; value change attachments count; menu count not showing

Turbo skips the project badges for an inactive user. A person with the job claim sees a creation count when approval is required. A person with the project authorization claim sees the count of requests awaiting their approval, only when approval is required. That person also sees the count of Pending value-change attachments assigned to them.
