<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_request_views.go","sha256":"3487fbe5277f236410e023c49a49ff8d72e2d407e589e3c142bf9bb255a3459a","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_request_views.go

## Request detail shows the full record of each linked person or record

Lines 27-42

Asked as: linked person missing on request; request page error; related job not showing; linked client not appearing

A project request's detail lists the linked client, job owner, contact, manager, alternate manager, creator, approver, branch, rate sheet, parent job, proposal and invoicing profile. A link that is empty is left out. If a linked record is missing, the detail fails.

## Request detail bills the client name the profile chooses

Lines 43-80

Asked as: wrong client name on invoice; alias instead of official name; bill to name wrong; which client name billed

A request's details show a bill-to name that is the official name or the alias, chosen by the invoicing profile (utilities.BillingName). They show the billing contact from that profile and the rate sheet entries, highest rate first.

## Request detail includes the invoicing contact

Lines 43-51

Asked as: invoice contact not showing; who receives the invoice; billing contact missing on request; invoicing contact on request

When the request has an invoicing profile that names a contact, the detail also shows that contact as the billing contact (billing_contact).

## Request detail lists rate sheet rates, highest rate first

Lines 66-84

Asked as: rates listed in wrong order; rate sheet rates on request; overtime rate not showing; rates sorted highest first

When the request has a rate sheet, its detail lists each role's rate and overtime rate, ordered with the highest rate first. Proposed rates are added by requestProposedRates, whose result this file does not show.

## Request detail shows people with profile name and phone

Lines 87-96

Asked as: approver name shown wrong; manager phone number missing; name shown for approver; person has no profile name

Each person in a request's detail (manager, alternate manager, creator, approver) shows the name from their profile (given name and surname) and their phone. If the person has no profile, the detail shows the user's name and email and an empty phone.

## Project request review token needs a fresh submission

Lines 98-105

Asked as: old review page after resubmit; approval stale after recall; checklist out of date after resubmitting; resubmitted request still needs review

The review token includes the submission token that submit issues. A recalled and resubmitted request gets a new token even when its data and PDF are unchanged. An old review page and its checked checklist no longer match.

## Project request review token excludes display-only details

Lines 106-168

Asked as: rate sheet name change approval; editing approver makes approval stale; approval still valid after edit; changing rate sheet name affects approval

Turbo builds a review token (a SHA-256 digest) for a project request from its saved data, its submission token (submission_token), its allocations and its attachment hash. It leaves out the approver, the creator, the rate sheet's code, name and updated time, and display-only rate sheet details. A change to those details does not make an approval stale. A change to the rates, the revision or the allocations does.

## Request detail includes checklist and attachment link

Lines 179-226

Asked as: checklist missing on request; attachment link not showing; where is the attachment; allocations not showing on request

The request detail response includes the current checklist version and questions for the request's scope. It also includes the attachment with a document link, the allocations with expanded divisions, and the review token.

## Who can approve or create from a request

Lines 211-226

Asked as: cannot approve request; approve button not showing; create job button missing; approval limit too low to approve

The chosen approver can approve a request only when its status is Awaiting approval and the approver's project authorization limit covers the project value. A person with the job claim can create the job only when the status is Awaiting creation. Rejection depends on canRejectRequest, which is defined in another file.

## Request detail page permissions

Lines 215-226

Asked as: no edit button on request; cannot submit my request; recall button missing; request buttons greyed out

The request detail page reports what the viewer can do, and every action except deleting an abandoned request needs project request approval to be enabled (IsJobRequestApprovalRequired). The creator can edit when the status is Draft, Recalled or Rejected. The creator can submit in Draft or Recalled, but not in Rejected. The creator can recall when the status is Awaiting approval or Awaiting creation.

## Who can delete a project request

Lines 215-226

Asked as: delete request button missing; cannot delete my request; remove abandoned request; delete old request

The creator can delete an editable request (Draft, Recalled or Rejected) while approval is enabled. An admin (the admin claim) can delete a request that abandonedRequest says is abandoned, even when approval is disabled.

## Request queues and who may open them

Lines 234-275

Asked as: unknown request queue error; which queue shows my requests; queue will not open; my requests list empty

The queue parameter picks the list of requests. The creation queue needs approval enabled and the job claim. The approval and approved queues need the project authorization claim and show requests assigned to the viewer. The abandoned queue needs the admin claim. The all queue needs the job claim. An empty or mine queue shows the viewer's own requests. Any other value gives "Unknown request queue."

## Queue contents by status

Lines 248-259

Asked as: request missing from approval list; approved list is empty; why request not in queue; request awaiting creation not showing

The approval queue shows requests assigned to the viewer with status Awaiting approval. The approved queue shows requests the viewer is approver of with status Awaiting creation. The creation queue uses creationQueueFilter, which is defined in another file.

## Abandoned requests cutoff

Lines 260-265

Asked as: when request becomes abandoned; how many days abandoned; abandoned cutoff date; old requests marked abandoned

The abandoned queue lists requests whose last edit (last_edited) is on or before JOB_REQUEST_ABANDONMENT_DAYS days ago. The comparison includes the cutoff time itself. The value of the constant is set in another file.

## Request list paging and order

Lines 283-312

Asked as: requests per page; request list sort order; client name in request list; how many requests shown

The request list is sorted by newest submitted, then newest created. Pages hold 50 items by default. A perPage value below 1 or above 100 becomes 50. A page below 1 becomes 1. Each row shows the client's name.

## Request options need a project-request role

Lines 323-336

Asked as: cannot use project requests; no access to request form; project request access denied; why can't open request form

Only active staff with the time, job, project authorization or admin claim can load the request form options. Others see "You cannot use project requests."

## Request form options contents

Lines 338-393

Asked as: rate sheet missing from dropdown; proposal not in list; parent project dropdown empty; subjob not in parent list

The form options list only active users. Rate sheets that belong to a request (proposed rates) are left out. Proposals are jobs whose number starts with P and whose status is In Progress, Submitted or Awarded. Parent projects are jobs whose number matches two digits, a hyphen and 3 or 4 digits, such as 25-123. Subjobs are not listed as parents.

## Authorization approvers lookup by project value

Lines 402-424

Asked as: who can approve this amount; approver dropdown empty; whole dollar value error; approver list for project value

The approver lookup needs a positive whole-dollar value. Otherwise it gives "Enter a positive whole-dollar project value." It lists only holders of the project authorization claim whose limit covers that value (ProjectAuthorizationLimit), with each person's maximum value (max_value).

## Menu badge counts for project workflow

Lines 427-470

Asked as: menu badge count wrong; pending approvals number on menu; notification count not showing; menu count for requests

Project workflow badge counts show only for active staff. Job claim holders see the count of requests awaiting creation when approval is enabled. Project authorization holders see the count of requests awaiting their approval when approval is enabled. They also see the count of pending value-change attachments assigned to them, even when approval is disabled.
