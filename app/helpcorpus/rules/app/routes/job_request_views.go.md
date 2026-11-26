<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_request_views.go","sha256":"e154b6898057474d383a164065943f00474f8aaacf1f7b2cd0211e8d04580b2b","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_request_views.go

## Project request detail loads related records for display

Lines 26-50

Asked as: request shows no client; billing contact missing on request; job owner not showing; request detail missing details

The request detail includes the client, job owner, contact, manager, alternate manager, creator, approver, branch, rate sheet, parent job, proposal and invoicing profile. When the invoicing profile has a contact, the detail also shows it as the billing contact (billing_contact).

## Bill-to name on a project request

Lines 51-64

Asked as: wrong bill to name; client alias on invoice; billing name wrong on request; invoice name differs from client

A project request's detail shows a bill-to name (bill_to_name) for its client. Turbo picks the official name or the alias according to the request's invoicing profile, through utilities.BillingName. The profile is optional.

## Rates shown on a request, highest first

Lines 65-84

Asked as: rates listed in wrong order; highest rate not at top; rate sheet order on request; rates not sorted

A project request's rate sheet rates are listed with the highest rate first. Each rate shows its role name, rate and overtime rate. Proposed rates are added by requestProposedRates.

## Person names come from the profile when there is one

Lines 86-95

Asked as: wrong name for approver; phone number missing on request; user name shown instead of full name; profile name not showing

Where a request names a person (manager, approver, creator), Turbo shows the name from the person's profile, given name and surname. Without a profile it falls back to the user's name. The phone number comes from the profile and is empty without one.

## Review token ties approval to one submission

Lines 97-107

Asked as: approval lost after resubmit; old review page still accepted; recalled request approval rejected; approve token invalid after recall

The review token is a digest of the request's saved data and its submission token. The submission token is a random value issued only by submit. A recalled and resubmitted request gets a new token, so an old review page's checked checklist is not accepted. A draft may have an empty token for viewing, but approving or creating must reject an empty token.

## Awarding the proposal does not make approval stale

Lines 109-170

Asked as: awarding proposal clears approval; proposal award resets approval; approval stale after award; award makes request stale

The review token leaves out the linked proposal's status, updated time and _imported mark. Procurement can award the proposal after approval and before creation without making the approval stale. Any other edit to the proposal makes it stale. The legacy digest still holds the whole proposal record.

## Review token ignores approver, creator and rate-sheet display data

Lines 127-158

Asked as: renamed rate sheet clears approval; profile edit resets approval; rate sheet code change approval; approval stale after rename

The review token (digest) leaves out the approver's and creator's details, the rate sheet's code, name and updated time, and the source-sheet and missing-role display data. A rename, a code change or a profile edit does not make an approval stale. A change to the other rate sheet fields or to any rate does make it stale.

## Review token covers submission, fields, allocations and document

Lines 171-200

Asked as: changed hours reset approval; edit allocation makes approval stale; attachment change approval gone; field edit clears approval

The review token holds the version, the submission token, the business fields, each division allocation (division, hours, whether the division is active) and the attachment's document hash. A change to any of these makes the saved approval stale.

## Request detail opens only for requests the person may see

Lines 205-211

Asked as: cannot open request page; request not visible to me; access denied on request; request page won't load for me

Opening a project request's detail page first calls loadVisibleRequest with the signed-in person. Turbo builds the page from that record only when that call lets the person see it. Another file defines who may see a request.

## Request detail shows each allocation's division by code and name

Lines 216-223

Asked as: division shows id not name; allocation division blank; division code missing on request; hours division shows number

The request detail page expands each hours allocation's division (division), so it shows the division's code and name and not its id. If the expansion fails, Turbo does not return the detail page.

## Request detail returns a review token and the approval setting

Lines 224-246

Asked as: review token missing; approval setting not showing; no approval flag on request; request response missing approval info

The request detail response includes a review token (review_token) computed from the request as saved now. It also includes whether project request approval is required (enabled, from IsJobRequestApprovalRequired). Turbo checks the person's job and admin claims and whether they can reject the request for the permission flags on the page.

## Request detail permissions

Lines 248-259

Asked as: buttons missing on request; cannot edit my request; cannot recall my request; edit button greyed out

Every permission is false when project request approval is not enabled (IsJobRequestApprovalRequired), except that an admin (the admin claim) can still delete an abandoned request. The creator can edit a request in Draft, Recalled or Rejected status. Submit follows the edit rule except that a Rejected request cannot be submitted. The creator can recall a request that is Awaiting approval or Awaiting creation.

## Who can approve or create from request detail

Lines 249-259

Asked as: approve button not showing; cannot create job from request; approver cannot approve; create job button missing

A person can approve only when the request is Awaiting approval, they are the chosen approver, and their project authorization limit covers the project value (ProjectAuthorizationLimit). A person can create the job only with the job claim (procurement) and when the status is Awaiting creation. Rejection depends on canRejectRequest.

## Request detail includes the checklist and document link

Lines 251-259

Asked as: approval checklist missing; attachment link not showing; no questions on approval; document link missing

The request detail includes the current approval checklist (version and questions) for the request's scope. When the request has an attachment, it adds a link to its document.

## Request list queues and who can open them

Lines 267-308

Asked as: unknown request queue; queue not showing; cannot open approval queue; queue empty for me

Request queues have these rules. The creation queue needs approval to be enabled and the job claim (procurement). The approval queue and the approved queue need the project authorization claim. They show the person's own requests that are Awaiting approval or Awaiting creation. The abandoned queue needs the admin claim. The all queue needs the job claim. The default and mine queues show the person's own requests. Any other queue gives "Unknown request queue."

## Abandoned queue cutoff

Lines 293-298

Asked as: abandoned list wrong days; old requests in abandoned; abandoned cutoff wrong; stale requests not in abandoned

The abandoned queue lists requests whose last edit is on or before the cutoff of JOB_REQUEST_ABANDONMENT_DAYS days ago.

## Request list paging and order

Lines 316-345

Asked as: request list in wrong order; too many requests per page; list sorted oldest first; page size wrong on requests

The request list is newest submitted first, then newest created. The page size is 50 by default. A perPage below 1 or above 100 gives 50. A page below 1 gives page 1. Each row shows the client name.

## Request form options need a claim and list only some records

Lines 354-441

Asked as: cannot use project requests; proposed rates in rate sheet list; inactive staff in dropdown; proposal dropdown empty

Only people with the time, job, project authorization or admin claim can load the request form options. Others get "You cannot use project requests." The rate sheets listed leave out the proposed rates of any request. Staff are listed only if their admin profile is active. Proposals are those whose number starts with P and whose status is In Progress, Submitted or Awarded.

## Approver lookup by project value

Lines 445-472

Asked as: no approvers listed; enter positive dollar value; approver list empty for amount; approver missing for value

The approver list needs a positive whole-dollar value, or it says "Enter a positive whole-dollar project value." It lists people with the project authorization claim whose limit covers that value, each with their maximum value (max_value). People whose limit does not cover the value are left out.

## Project workflow badge counts

Lines 475-518

Asked as: count badge not showing; pending approvals number wrong; badge shows for inactive user; jobs badge count missing

People with the job claim see a count of requests Awaiting creation, when approval is enabled. People with the project authorization claim see a count of requests assigned to them that are Awaiting approval, when approval is enabled. They also see a count of pending value change approvals assigned to them. Inactive users get no badges.
