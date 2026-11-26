<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_request_views.go","sha256":"4d900e2a8be38a5443a9fe3b463c36a2fbbbdedb997573bb81965dca68d58d43","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_request_views.go

## Request page fails if a linked record is missing

Lines 27-42

Asked as: request page error; project request won't open; broken link on request; missing client record error

Turbo loads each record that a project request links to: client, job owner, contact, manager, alternate manager, creator, approver, branch, rate sheet, parent, proposal and invoicing profile. If any linked record cannot be found, the request page returns an error and does not load. Turbo leaves out a link that is empty.

## Request page shows the invoicing profile's billing contact

Lines 43-51

Asked as: billing contact missing; no billing contact shown; invoice contact not on request

When a project request has an invoicing profile (invoicing_information) with a contact, Turbo adds that contact to the request page as the billing contact (billing_contact). A request with no invoicing profile shows no billing contact.

## Request page shows the client's bill-to name

Lines 52-65

Asked as: wrong bill to name; alias billed instead of name; invoice shows wrong client name

Turbo sets the bill-to name (bill_to_name) of a project request from the client and the request's invoicing profile. The profile decides whether the official name or the alias is billed. The profile is optional. Without one, the name comes from the client alone (utilities.BillingName).

## Request rates list shows highest rate first

Lines 66-85

Asked as: rate list order wrong; rates sorted lowest first; rate sheet entries order

When a project request has a rate sheet, Turbo lists the sheet's entries from the highest rate to the lowest. Each entry shows its role name, rate and overtime rate. Turbo adds the request's proposed rates to this list.

## Manager and approver names come from the profile

Lines 87-96

Asked as: name shows email; approver name missing; person shows no phone

Turbo shows a person on a request by the given name and surname of their profile, and shows their phone. Without a profile, Turbo shows the user's own name and email, with an empty phone.

## Review token needs a submission token to approve or create

Lines 98-108

Asked as: cannot approve checklist; approve button rejected; empty submission token

A review digest binds the review page to the saved data and to one submission, using a random submission token that only submit issues. Recalling and resubmitting therefore gives a new digest, and an old checked checklist is not accepted. A draft can have an empty token for viewing or printing. Approval and creation reject an empty token.

## Two older forms of the review digest exist

Lines 109-138

Asked as: old review digest; legacy digest version; digest format outdated

Turbo can rebuild two older forms of a request's review digest. One holds the invoicing profile, its billing contact and the bill-to name (reviewDigestWithProfile). The other also holds the whole proposal record, with its status (reviewDigestWholeProposal). Turbo issues only the current form, and it uses the older forms only to recognise approvals saved before the change.

## Older approval digests still match

Lines 140-154

Asked as: approval still valid after change; old approval still counted; approval not stale

An approval's saved digest still counts as current if it equals the current digest or one of the two older forms (with the invoicing profile, or with the whole proposal record). Any other saved digest is stale.

## Invoicing profile excluded from approval review digest

Lines 177-186

Asked as: change invoicing profile after approval; billing profile edit approval; clear invoicing profile approval

A project request's current review digest (reviewDigestCurrent) leaves out the invoicing profile, its billing contact and the bill-to name. Procurement can set, change or clear the invoicing profile until the job is created without making an approval stale. A change to the client, including its name or alias, still makes the approval stale.

## Approver and creator changes do not stale an approval

Lines 187-201

Asked as: changing approver after approval; creator edit resets approval; profile rename approval

The review digest of a project request leaves out the approver and creator details. It also leaves out the proposed-rates approver's contact details, the source rate sheet's own record and the missing-role rows. A profile edit, rename or newer revision of those does not change what was approved.

## Rate sheet code, name and updated time ignored in digest

Lines 202-213

Asked as: rename rate sheet approval; rate sheet code changed; rate sheet name change

The review digest binds the rate sheet by its id and leaves out its code, name and updated time. A rename or code change does not make the approval stale. A change to the revision, active flag, effective date, reason, based_on, rates_approver, rates_approved or any rate does make it stale.

## Awarding the linked proposal does not stale approval

Lines 214-225

Asked as: award proposal after approval; proposal awarded approval lost; proposal edit resets approval

The review digest leaves out the linked proposal's status, updated time and _imported mark. Procurement can award the proposal after approval and before job creation without making the approval stale. Any other edit to the proposal makes the approval stale.

## Review digest covers submission token, allocations and document

Lines 226-258

Asked as: recall and resubmit new token; checklist reset after recall; allocation change review

The review digest includes the request's version, its submission token, its business fields, each division allocation (division, hours, whether the division is active) and the attachment's document hash. Recalling and resubmitting therefore issues a new token. A page that was reviewed earlier no longer matches.

## Request page shows divisions and the current checklist

Lines 263-289

Asked as: division code not shown; hour allocations show id; checklist not loading

The request page loads the request, its attachment and its hour allocations in one transaction. Each allocation shows its division code and name, not its id. The page also gets the review token and the current checklist version and questions for the request's scope (requestChecklistScope).

## Request detail tells the page if approval is on

Lines 290-335

Asked as: approval disabled buttons gone; no edit button; approval turned off

The request detail and the list carry an enabled flag, which shows whether project request approval is required (IsJobRequestApprovalRequired). When it is off, the detail page gives no edit, submit, recall, reject, approve or create permission. The delete permission for an abandoned request (admin) does not depend on it.

## Request detail permissions

Lines 306-335

Asked as: cannot edit draft; recall button missing; cannot submit request

All request actions need project request approval to be enabled (IsJobRequestApprovalRequired). Only the creator can edit a request in Draft, Recalled or Rejected status. The creator can submit it unless it is Rejected. The creator can recall it while it is Awaiting approval or Awaiting creation.

## Approval blocked message for selected approver

Lines 317-330

Asked as: approval limit message; cannot approve over limit; no project authority

When the selected approver's limit does not cover the value, the request detail shows only them "Your approval limit does not cover this project value." If they lack the project authorization claim it shows "You do not have project approval authority." Nobody else sees this message.

## Who can approve, create and delete a request

Lines 321-335

Asked as: who can approve request; delete abandoned request; cannot delete request

Only the selected approver can approve, and only while the request is Awaiting approval and the approver's project limit covers the project value. A job-claim holder can create a job from a request that is Awaiting creation. A creator can delete an editable request. An admin can delete an abandoned request (abandonedRequest).

## Request queues and who may open them

Lines 343-384

Asked as: unknown request queue; request queue names; where is my queue

The queues are: creation (job claim, needs approval enabled), approval and approved (project authorization claim, the requests where you are the approver, with status Awaiting approval or Awaiting creation), abandoned (admin claim), all (job claim), and mine (default, your own requests). Any other queue name gives "Unknown request queue."

## Request list for the queue "all" and "creation" needs the job claim

Lines 349-379

Asked as: cannot see all requests; creation queue missing; all tab not showing

The creation queue and the all queue both need procurement (the job claim). The creation queue also needs project approval to be turned on. The abandoned queue needs the admin claim.

## Abandoned queue cutoff

Lines 369-374

Asked as: abandoned requests list; stale requests days; old requests cutoff

The abandoned queue lists requests whose last_edited is on or before the cutoff of now minus JOB_REQUEST_ABANDONMENT_DAYS days. Only an admin can open it.

## Request list paging and order

Lines 392-402

Asked as: request list order; newest requests first; request list sorting

The request list shows the newest submitted first, then the newest created. It defaults to page 1 and 50 items per page. A perPage below 1 or above 100 falls back to 50. A page below 1 becomes 1.

## Request list page size defaults to 50

Lines 392-399

Asked as: list 50 per page; items per page request; default page size requests

The request list shows 50 requests a page when perPage is missing, below 1 or above 100. A page below 1 becomes page 1. The list reports the total items and total pages.

## Request list shows client name

Lines 402-416

Asked as: client name missing list; request list client column; client name in list

Each row in the request list carries the client's name when the client is found. The list is sorted by newest submission, then newest creation, then id.

## Request list reports paging totals and whether approval is on

Lines 406-421

Asked as: total pages requests; how many pages; empty page request list

The request list reports the total number of requests (totalItems), the page, the page size and the number of pages. The number of pages is the total divided by the page size, rounded up. A page past the last one returns an empty list. The list also reports whether project request approval is required (enabled).

## Project request options need a qualifying claim

Lines 432-445

Asked as: cannot use project requests; request form not opening; project request access denied

Loading the request form's options needs an active account and at least one of the time, job, project_authorization or admin claims. Otherwise Turbo says "You cannot use project requests."

## Request form option lists

Lines 449-512

Asked as: rate sheet dropdown; proposal dropdown missing; people list on form

The form's rate sheet list leaves out the proposed rates that belong to a request. Its people list holds only people with an active admin profile, named from their profile when they have one. Its proposals list holds only proposals (numbers starting with P) with status In Progress, Submitted or Awarded.

## Request form lists only active people with an admin profile

Lines 474-501

Asked as: missing person in request; inactive admin in dropdown; user not in list

The people in the request form are users whose admin profile is active. Turbo names each from the profile when there is one.

## Request form options say whether approval is required

Lines 513-517

Asked as: approval required option; form approval setting; approval flag on form

The project request options include a flag (enabled) that says whether project request approval is required. The form loads this flag with its lists.

## Authorization approvers lookup needs an active signed-in user

Lines 521-525

Asked as: approver lookup failed; cannot find approvers; workflow check error

Before Turbo lists the approvers who can cover a project value, it checks that the caller passes the workflowActive check. A caller who fails that check gets its error and no list.

## Authorization approvers lookup by project value

Lines 526-547

Asked as: who can approve amount; approver limit list; find approver for value

Looking up approvers needs a positive whole-dollar value, otherwise Turbo says "Enter a positive whole-dollar project value." It lists project authorization claim holders whose limit covers that value, each with their max_value. Holders whose limit does not cover it are skipped.

## Authorization approvers lookup needs a positive whole-dollar value

Lines 526-529

Asked as: enter positive whole dollar; project value invalid; cents in project value

The approvers lookup refuses a missing, non-numeric or non-positive value or a value that is not a whole number of dollars. Turbo answers "Enter a positive whole-dollar project value."

## Project workflow badge counts

Lines 551-593

Asked as: badge count requests; approval notification count; pending value change count

Job-claim holders see a badge count of requests Awaiting creation when approval is enabled. Project authorization holders see a count of requests Awaiting approval that name them as approver, only when approval is enabled. They also see a count of pending value-change attachments assigned to them, even when approval is disabled.
