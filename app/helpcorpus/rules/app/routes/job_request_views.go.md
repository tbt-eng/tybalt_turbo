<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_request_views.go","sha256":"a2e8df4351827b6afa0cfd4052bfbf38812f6d6d2d2f0d43b0a11561d6ad24eb","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_request_views.go

## Request page shows related records for each linked field

Lines 27-42

Asked as: request page not loading; linked record missing error; project request won't open; page fails to load

The project request page loads the linked records for client, job owner, contact, manager, alternate manager, creator, approver, branch, rate sheet, parent job, proposal and invoicing profile (invoicing_information). It shows only the links that are set. The page fails to load if a linked record is missing.

## Linked people on a request show profile names and phone

Lines 37-39

Asked as: approver shows email not name; manager phone number blank; person name wrong on request; contact phone missing

The request page shows each linked person (manager, alternate manager, creator, approver) with the name from their profile (given name and surname), plus their phone. Without a profile, Turbo shows the user's name and email and an empty phone.

## Request page shows the invoicing profile's billing contact

Lines 43-51

Asked as: billing contact not showing; invoice contact on request; who is the billing contact; billing person missing

When a project request has an invoicing profile (invoicing_information) that names a contact, the request page shows that contact as the billing contact (billing_contact). A request with no invoicing profile, or a profile with no contact, shows no billing contact.

## Bill-to name on a project request follows the invoicing profile

Lines 52-65

Asked as: bill to name wrong; billing name on project request; client name on invoice; wrong name for billing

A project request's bill-to name (bill_to_name) comes from utilities.BillingName for the client and the chosen invoicing profile. The profile is optional. Without one, the function gets the client alone.

## Rate sheet rates listed highest first on a request

Lines 66-80

Asked as: rate sheet order; highest rate first; rates listed in wrong order; rate list sorting

A project request shows the rates of its rate sheet with the highest rate first. Each entry shows the role name, the rate and the overtime rate.

## Request page rate list adds proposed rates from another function

Lines 81-84

Asked as: proposed rates not showing; extra rates on request page; rates missing from request; where do proposed rates come from

The request page lists the rate sheet's entries, then calls requestProposedRates to add the request's proposed rates to that list. This file does not show how those rates are chosen or combined.

## Approver and manager names come from the profile first

Lines 87-96

Asked as: name shows email instead of profile; approver shown with wrong name; phone number empty for approver; profile name not used

Where a request page names a person, such as the manager, approver or creator, Turbo uses the given name and surname from that person's profile. If the person has no profile, it uses the user record's name. The phone number comes from the profile and is empty when there is no profile.

## Review token ties approval to one submission

Lines 98-108

Asked as: approval link stopped working; old review page invalid; resubmit broke approval; approval token expired

The review token covers the saved request data, its allocations, the attachment's document hash and a random submission token (submission_token) issued on submit. A recall and resubmit of unchanged data gets a new token, so an old review page does not stay valid. A draft can have an empty submission token for viewing, but approve and create reject an empty one.

## Older approval digests still count as matching

Lines 126-154

Asked as: old approval still counts; approval still valid after edit; digest mismatch on approval; older approval matching

An approval still matches if its saved digest equals the current digest or one of two older forms. The older forms hold the invoicing profile, or the whole proposal. They hold more data, so a match proves that the data the current digest covers did not change.

## Test helper can build older review digest forms

Lines 156-170

Asked as: test digest helper; unknown digest form error; build old digest test; digest form in tests

A test-only function builds an older form of the review digest, "profile" or "proposal", for a request. It rejects any other form name with the error "unknown digest form".

## Review token leaves out invoicing profile and approver

Lines 177-188

Asked as: billing change makes approval stale; invoicing profile edit approval; clear invoice profile approval; approver change approval invalid

The current review digest leaves out the invoicing profile, its billing contact and the bill-to name. Procurement can set, change or clear the profile until the job is created without making an approval stale. A change to the client's name or alias still makes the approval stale. The digest also leaves out the approver and the creator.

## Proposed-rates display details do not stale an approval

Lines 189-213

Asked as: rates approver profile edit approval; rate sheet rename approval; rate sheet code change stale; proposed rates display change

The review digest leaves out the rates approver's contact details, the source sheet's own record and the missing-role rows. It also leaves out the rate sheet's code, name and updated time. A rename, a code change or a profile edit does not make an approval stale. A change to the revision, active, effective date, reason, based_on, rates approver, rates approved or any rate does.

## Review token ignores rates approver details and missing roles

Lines 189-201

Asked as: approver contact change stale; missing roles approval; source sheet rename approval; rates approver details approval

The review token leaves out the proposed rates approver's contact details, the source rate sheet's own record and the missing-roles rows. Editing the approver's profile or renaming the source sheet does not make an approval stale.

## Review token leaves out rate sheet code, name and updated time

Lines 202-213

Asked as: rate sheet renamed approval; rate sheet code change; rate sheet updated time approval; rate sheet name change stale

The review token leaves out the rate sheet's code, name and updated time, because the request binds the sheet by its id. A rename or code change does not make an approval stale. A change to the revision, active flag, effective date or rates does.

## Awarding the proposal does not stale an approval

Lines 214-225

Asked as: award proposal makes approval stale; proposal status change approval; proposal marked awarded after approval; proposal edit approval invalid

The review digest leaves out the linked proposal's status, updated time and _imported mark. Procurement can award the proposal after approval and before creation without making the approval stale. Any other edit to the proposal makes the approval stale. Job creation needs an Awarded proposal. Approval does not.

## Review token includes version, submission and business fields

Lines 226-258

Asked as: approval no longer valid; review token mismatch; submission token changed; why approval went stale

The review token (digest) covers the request's version, its submission token (submission_token), its business fields, its related records, its division hours and the PDF's document hash. A change to any of them makes the old token no longer match.

## Division active flag is part of the digest

Lines 233-245

Asked as: division inactive approval stale; hours allocation change approval; division active flag stale; allocation edit breaks approval

The review digest holds each hours allocation with its division and the division's active flag. A change to the allocations or a division's active flag makes an approval stale.

## Replacing the authorization PDF makes an approval stale

Lines 246-252

Asked as: replace PDF approval invalid; new attachment approval stale; document changed after approval; swap authorization document approval

The review token includes the hash of the request's attachment document (document_hash). A different document therefore gives a different token, and an approval given on the old document no longer matches. A request with no attachment adds no hash.

## Request page permissions

Lines 306-335

Asked as: cannot edit project request; submit button missing; recall button not showing; why request not editable

A project request is editable by its creator while its status is Draft, Recalled or Rejected. Submit has the same rule, except that a Rejected request cannot be submitted. Recall is available to the creator while the status is Awaiting approval or Awaiting creation. Every one of these actions needs the approval workflow to be on (enabled).

## Who can delete a project request

Lines 306-335

Asked as: delete project request; remove abandoned request; admin delete request; delete draft request

The creator can delete an editable request (Draft, Recalled or Rejected) when the workflow is on. An admin can delete an abandoned request (abandonedRequest) even when the workflow is off.

## Only the selected approver can approve, within their limit

Lines 308-335

Asked as: approve button missing; cannot approve request; no approve option; approver limit approve

The page offers Approve only when the workflow is on, the status is Awaiting approval, the person is the selected approver and their project authorization limit covers the project value.

## Request page shows the current checklist

Lines 309-312

Asked as: checklist not showing; approval checklist questions; request checklist version; checklist missing on request

The request page loads the current approval checklist version and its questions for the request's scope (requestChecklistScope). The caller of the page receives them with the review token.

## Approver told why Approve is blocked

Lines 317-330

Asked as: why cannot approve; approval blocked message; no project authority message; approval limit reason

The selected approver of an Awaiting approval request, whose limit does not cover it, sees a reason. It is "You do not have project approval authority." if they lack the project authorization claim. Otherwise it is "Your approval limit does not cover this project value." Nobody else sees the reason.

## Create job needs the job claim and Awaiting creation

Lines 335-335

Asked as: create job button missing; no create job option; job claim create; cannot create job from request

The page offers creating a job from a request only when the workflow is on, the person has the job claim (procurement) and the status is Awaiting creation.

## Request queues and their required claims

Lines 343-382

Asked as: unknown request queue; queue access denied; approval queue not visible; which queues can I see

The creation queue needs the workflow on and the job claim. The approval and approved queues need the project authorization claim. The abandoned queue needs the admin claim. The all queue needs the job claim. The mine queue, which is the default, lists the person's own requests and needs no claim. An unknown queue gives "Unknown request queue."

## Queue contents

Lines 354-371

Asked as: approval queue empty; what shows in approved queue; abandoned requests list; queue contents explained

The approval queue lists requests where the person is the approver and the status is Awaiting approval. The approved queue lists those with the status Awaiting creation. The abandoned queue lists requests last edited on or before JOB_REQUEST_ABANDONMENT_DAYS days ago.

## Request list paging and order

Lines 383-426

Asked as: request list paging; how many requests per page; request list sort order; newest requests first

The request list shows 50 per page by default. A perPage below 1 or above 100 also gives 50, and a page below 1 gives page 1. Requests are sorted by newest submitted, then newest created.

## Project request options need one of four claims

Lines 427-439

Asked as: cannot use project requests; request form options error; no access to request form; request form won't load options

Loading the request form's options needs the time, job, project authorization or admin claim. Without one Turbo refuses with "You cannot use project requests."

## Option lists on the request form

Lines 440-506

Asked as: rate sheet missing from dropdown; proposal missing from list; people dropdown missing someone; inactive person in list

The options leave out rate sheets that belong to a request's proposed rates. The people list holds only people whose admin profile is active. A person's name comes from their profile when one exists. The proposals list holds only proposals (number starts with P) with status In Progress, Submitted or Awarded.

## Approver search needs a positive whole-dollar value

Lines 515-520

Asked as: approver search needs value; positive whole dollar error; enter project value; approver search error

The list of possible approvers needs a positive whole-dollar project value. Otherwise Turbo says "Enter a positive whole-dollar project value."

## Approver search lists people whose limit covers the value

Lines 521-539

Asked as: approver list missing person; who can approve this amount; approver max value; people over limit approver list

The list holds people with the project authorization claim whose limit covers the value, each with their maximum value (max_value). People whose limit does not cover it are skipped.

## Menu badges for project workflow

Lines 542-583

Asked as: menu badge count; notification badge requests; pending approvals badge; job badge number

A person with the job claim sees a badge counting requests Awaiting creation, when the workflow is on. A person with the project authorization claim sees a badge of requests Awaiting approval assigned to them, when the workflow is on. They also see a badge of pending value approvals assigned to them, when the workflow is on or off.
