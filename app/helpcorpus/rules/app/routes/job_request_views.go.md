<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_request_views.go","sha256":"1f070449230d2c12cecca62c21617649e1e261e369f6593cb1975f4e35393649","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_request_views.go

## Project request page shows linked records

Lines 29-42

Asked as: linked records missing on request; blank fields on request page; request page error deleted record; related client not showing; job owner missing from request

The request detail shows each record the project request links to: client, job owner, contact, manager, alternate manager, creator, approver, branch, rate sheet, parent, proposal and invoicing profile. A link that is blank is left out. The page fails with an error if a linked record no longer exists.

## Billing contact comes from the invoicing profile

Lines 43-51

Asked as: billing contact not showing; who gets the invoice; invoice contact on request; billing person missing; no contact on request

When the request's invoicing profile (invoicing_information) has a contact, the request page also shows that contact as the billing contact (billing_contact). A request with no invoicing profile, or a profile with no contact, shows no billing contact.

## Bill-to name on a project request

Lines 52-65

Asked as: invoice name wrong; bill to name missing; client name on invoice; billing name blank; name on invoice incorrect

The request detail shows a bill-to name (bill_to_name) for the client. The name comes from utilities.BillingName, which uses the client and the invoicing profile (official name or alias). A request with a client but no invoicing profile still gets a bill-to name.

## Rate sheet entries on a project request, highest rate first

Lines 66-81

Asked as: rates listed in wrong order; rate sheet entries missing; role rates on request; overtime rate not showing; rate list empty

A request with a rate sheet lists the sheet's entries with role name, rate and overtime rate. Turbo sorts them with the highest rate first. A request with no rate sheet shows an empty list.

## People on a project request show profile name and phone

Lines 84-93

Asked as: manager name shows wrong; phone number missing for approver; user name instead of profile; approver contact details missing; person shows no phone

Manager, approver and other people on a request show their id, name, email and phone. When the person has a profile, the name is the profile's given name and surname joined, and the phone comes from the profile. When the person has no profile, the name is the user account's name and the phone is empty.

## Review token covers saved data and the submission

Lines 95-140

Asked as: approval code changed; submission token empty; why new token on resubmit; review code differs after recall; token different after resubmit

The review token is a SHA-256 digest of the request's version, its submission token (submission_token), its business fields, the related records (not the approver or creator), the division allocations (division, hours, whether the division is active) and the attachment's document hash. A recall and resubmit of unchanged data gives a new token, because each submit issues a new submission token. The token does not replace the permission and state checks. Approve and create must reject an empty submission token, while a draft can have one for viewing or printing.

## Request detail includes checklist, attachment link and allocations

Lines 155-198

Asked as: checklist missing on request; attachment link not showing; division split not shown; request detail missing attachment; allocations not displayed

The request detail returns the current checklist version and questions for the request's scope. It also returns the attachment with a document link, and the division allocations with each division expanded.

## Who can create the job from a request, and who can reject

Lines 183-198

Asked as: create job button missing; procurement make job from request; reject button missing; turn request into job; who can reject request

Procurement (the job claim) can create the job from a request that is Awaiting creation, when approval is enabled. The reject button depends on canRejectRequest, which is defined in another file. Approval being enabled is also required.

## Who can edit, submit, recall or delete a project request

Lines 187-198

Asked as: cannot edit my request; take back submitted request; delete my draft request; edit button greyed out; submit button missing

Only the creator can edit a request, and only while it is Draft, Recalled or Rejected. Submit follows the same rule but is not offered for Rejected requests. The creator can recall a request that is Awaiting approval or Awaiting creation. All of these require that project request approval is enabled (IsJobRequestApprovalRequired). Delete is also offered to the creator under the edit rule. An admin (the admin claim) can delete an abandoned request (abandonedRequest) even when approval is disabled.

## Who can approve a project request

Lines 189-198

Asked as: approve button missing; cannot approve request; approval limit too low; approve someone else's request; approver cannot sign off

The assigned approver can approve a request that is Awaiting approval, and only if their project authorization limit covers the request's project value (ProjectAuthorizationLimit). Approval must be enabled. Another person cannot approve.

## Project request queues and who can open them

Lines 206-247

Asked as: request queue empty; where do I find approvals; awaiting approval list; unknown request queue; abandoned requests tab

The queues need an active user. 'mine' (the default) shows requests the person created. 'creation' needs procurement (the job claim) and approval enabled. 'approval' needs the project authorization claim and shows requests assigned to that approver that are Awaiting approval. 'approved' shows the same approver's requests that are Awaiting creation. 'all' needs the job claim. 'abandoned' needs the admin claim. An unknown queue gives the message 'Unknown request queue.'.

## Abandoned requests queue cutoff

Lines 232-237

Asked as: when request becomes abandoned; abandoned after how many days; request marked abandoned too soon; last edited cutoff; abandonment days setting

The abandoned queue lists requests whose last edit (last_edited) is at or before the cutoff of JOB_REQUEST_ABANDONMENT_DAYS days ago. A request edited exactly at the cutoff counts as abandoned.

## Request list paging and sort order

Lines 249-292

Asked as: request list order; requests per page; newest requests first; page size requests; list paging wrong

The list shows newest submitted first, then newest created. The default page is 1, and a page below 1 becomes 1. The page size is 50 by default, and a size below 1 or above 100 becomes 50. Each row includes the client name.

## Who can load the project request form options

Lines 293-331

Asked as: cannot use project requests; no access to request form; project request permission denied; why can't open request form; request form locked

The options need an active user with at least one of these claims: time, job, project authorization or admin. Anyone else gets 'You cannot use project requests.'. Only active users appear in the list of people.

## Which jobs appear as proposals and parents in form options

Lines 332-353

Asked as: parent project dropdown empty; proposal jobs not showing; which jobs can be parent; subjob missing from parent list; proposal list empty

Proposals are jobs whose number starts with P and whose status is In Progress, Submitted or Awarded. Parent projects are jobs whose number matches two digits, a dash, and three or four digits (such as 25-123). Subjobs and other jobs are not offered as parents.

## Authorization approvers list needs a whole-dollar value

Lines 362-391

Asked as: who can approve amount; approver list empty; whole dollar value; approver limit covers value; approvers for project value

The approvers list needs a positive whole-dollar project value. Otherwise Turbo says 'Enter a positive whole-dollar project value.'. It lists only the people with the project authorization claim whose limit covers that value, and shows each person's maximum value. People whose limit does not cover it are left out.

## Duplicate check lists similar jobs and requests

Lines 392-451

Asked as: duplicate job warning; similar jobs check; possible duplicate requests; proposals ignored in duplicates; check for similar requests

The duplicate check compares a request with existing jobs, skipping proposals (numbers starting with P). Only procurement (the job claim) also sees other pending requests that look like duplicates. The request itself is left out.

## Reasons a job or request counts as a duplicate

Lines 452-467

Asked as: why flagged as duplicate; similar description warning; same client PO duplicate; parent project duplicate reason; duplicate match reason

Turbo gives a reason when a non-empty client, proposal, parent, location, client PO or client reference number matches. It also gives 'Similar description' when one description contains the other, ignoring case and spaces; an empty description never matches. When the request's parent is the other record, it adds 'Parent project'.

## Project workflow badge counts

Lines 469-513

Asked as: badge count wrong; number on approvals icon; pending requests count; no count showing; value change attachments count

Procurement (the job claim) sees a count of requests Awaiting creation, only when approval is enabled. Holders of the project authorization claim see a count of requests Awaiting approval assigned to them, only when approval is enabled. They also see a count of pending value-change attachments assigned to them, whether or not approval is enabled. An inactive user gets no counts.
