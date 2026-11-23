<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_request_views.go","sha256":"0a2f75360b1065b5ac752ffa7d9fafc404626525a13e081011f06ee5f000e09a","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_request_views.go

## Request detail lists the related records and fails if one is missing

Lines 26-41

Asked as: request detail fails; missing related record; detail page error; record not found on request

The request detail includes the client, job owner, contact, manager, alternate manager, creator, approver, branch, rate sheet, parent job, proposal and invoicing profile, for each one the request names. If a named record cannot be found, the whole detail fails instead of leaving that record out.

## Request detail shows the billing contact of the invoicing profile

Lines 42-50

Asked as: billing contact missing; invoicing contact not showing; who to bill contact; billing contact on request

When a project request names an invoicing profile (invoicing_information) that has a contact, the request detail includes that contact as the billing contact (billing_contact). A request without an invoicing profile, or a profile without a contact, has no billing contact.

## Request detail bills the client name chosen by invoicing profile

Lines 51-64

Asked as: bill to name; billing name wrong; client alias on invoice; official name vs alias

The project request detail shows a bill-to name (bill_to_name) from the client and the invoicing profile. The invoicing profile chooses the official name or the alias, through utilities.BillingName.

## Request rate sheet rates list highest rate first

Lines 65-79

Asked as: rate sheet order; highest rate first; rates sorted; sheet rates list

The request detail lists the rate sheet's entries with the highest rate first. Each entry shows role, rate and overtime rate.

## Request detail names people from their profile

Lines 86-95

Asked as: manager name shows username; person name blank; profile name missing; approver name wrong

The request detail shows the manager, alternate manager, creator and approver by the given name and surname on their profile. When a person has no profile, Turbo uses the user name and an empty phone.

## Project request review token changes when request is resubmitted

Lines 97-167

Asked as: review token changed; old review page checklist; resubmit checklist stale; approval page out of date

Turbo builds a review token (a SHA-256 digest) for a project request from its version, its submission token (submission_token), its business fields, related records, division hours and the attachment's document hash. Recalling and resubmitting gives a new submission token, so an old review page's checked checklist no longer matches.

## Review token ignores approver, creator and display-only rate sheet details

Lines 110-136

Asked as: rename rate sheet stale approval; profile edit approval stale; what makes approval stale; rate change review token

The review token leaves out the approver, the creator, the rates approver's contact details, the source rate sheet's record and the missing-role rows. It also leaves out the rate sheet's code, name and updated time. Renaming a sheet or editing a profile does not make an approval stale. A changed rate, revision, active flag or effective date does.

## Request detail opens only requests the person may see

Lines 171-185

Asked as: cannot open request; can't see request; request not visible; request access denied

Turbo loads a request's detail with a check of what the signed-in person may see (loadVisibleRequest). The detail also returns the attachment and the allocations. Everything is read in one transaction, so the page shows one consistent state.

## Request detail shows allocation divisions by code and name

Lines 186-189

Asked as: division shows id; allocation division code; division name missing; allocation shows number

The request detail expands each allocation's division, so the page shows the division's code and name and not its id. If the expansion fails, the detail fails.

## Request detail returns the review token and current checklist

Lines 190-220

Asked as: review token missing; current checklist version; checklist questions request; request checklist not showing

The request detail returns a review token (review_token) for the request's current saved data and submission. It also returns the current checklist version and questions for the request's scope.

## Approval setting switches off project request actions in the detail

Lines 198-225

Asked as: no approval buttons; approval turned off; cannot edit request approval off; actions greyed out

When approval of project requests is not required (IsJobRequestApprovalRequired is false), the detail allows no edit, submit, recall, reject, approve or create. An admin (the admin claim) can still delete an abandoned request.

## Who can edit, submit or recall a project request

Lines 214-225

Asked as: edit my request; submit button missing; recall request; rejected request resubmit

Only the creator of a project request can edit it, and only in status Draft, Recalled or Rejected. The creator can submit it in Draft or Recalled, but not Rejected, so a rejected request cannot be submitted from the detail page. The creator can recall it in Awaiting approval or Awaiting creation. All of these actions are off when project request approval is not required (IsJobRequestApprovalRequired).

## Who can approve, create or reject a project request

Lines 214-225

Asked as: approve button missing; create job button; reject request option; why cannot approve

A person can approve a request only if its status is Awaiting approval, they are its chosen approver, and their project authorization limit covers the project value. Only a person with the job claim (procurement) can create the job, and only in status Awaiting creation. Reject depends on canRejectRequest. These actions are off when project request approval is not required.

## Who can delete a project request

Lines 214-225

Asked as: delete request; remove draft request; abandoned request delete; cannot delete request

The creator can delete a request they can edit (Draft, Recalled or Rejected) while approval is required. An administrator (the admin claim) can also delete an abandoned request (abandonedRequest).

## Project request queues and who can open them

Lines 233-275

Asked as: request queue empty; approval queue; queue unknown; which requests I see

The creation queue needs approval to be required and the job claim. It lists requests Awaiting creation. The approval queue needs the project authorization claim and lists requests assigned to you with status Awaiting approval. The approved queue lists requests you approved that are still Awaiting creation. The all queue needs the job claim. The abandoned queue needs the admin claim. The mine queue, which is also the default, lists your own requests. Any other queue gives "Unknown request queue."

## Abandoned requests are those not edited for a set number of days

Lines 259-264

Asked as: stale requests; old unedited requests; inactive requests list; requests not touched

The abandoned queue lists requests whose last edit (last_edited) is at or before now minus JOB_REQUEST_ABANDONMENT_DAYS days. A request edited exactly at the cutoff counts as abandoned.

## Project request list paging and order

Lines 282-311

Asked as: request list order; requests per page; paging requests; how many requests shown

The list sorts by newest submitted, then newest created. Page size is 50 by default. A perPage below 1 or above 100 falls back to 50, and a page below 1 becomes 1.

## Request options need a qualifying claim

Lines 320-335

Asked as: cannot use project requests; no access project requests; request form locked; need claim for requests

A person needs the time, job, project_authorization or admin claim to load the project request form options. Otherwise Turbo says "You cannot use project requests."

## Request form options leave out proposed rate sheets and inactive people

Lines 339-402

Asked as: rate sheet missing from form; inactive staff in list; proposal dropdown empty; why proposal not listed

The form's rate sheet list excludes sheets that belong to a request, so proposed rates are not listed. The people list includes only staff whose admin profile is active. The proposals list has only proposals (number starts with P) with status In Progress, Submitted or Awarded.

## Authorization approver list by project value

Lines 411-440

Asked as: approver list empty; authorization limit approver; project value approver; who can approve amount

The approver list needs a positive whole-dollar value. Otherwise Turbo says "Enter a positive whole-dollar project value." It lists people with the project_authorization claim whose limit covers the value, and shows each person's maximum value. People whose limit does not cover it are left out.

## Project workflow badge counts

Lines 441-485

Asked as: badge count wrong; workflow count; pending requests number; value change count

A person with the job claim sees a count of requests Awaiting creation, only when approval is required. A person with the project_authorization claim sees a count of requests Awaiting approval assigned to them, only when approval is required. They also see a count of pending value-change attachments assigned to them, even when approval is not required.
