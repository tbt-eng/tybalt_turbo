<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_request_views.go","sha256":"9cc6b1549c969271ae025a840d3d1056bf6e6fbf318bdf2943db0f73c4a4dafb","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_request_views.go

## Request detail shows the linked records by role

Lines 27-42

Asked as: linked records on request; which client is linked; related contacts on request; request links to people

A project request's detail shows the records it links to: client, job owner, contact, manager, alternate manager, creator, approver, branch, rate sheet, parent project, proposal and invoicing profile. Turbo shows a linked record only when the request has one for that field (client, job_owner, contact, manager, alternate_manager, creator, approver, branch, rate_sheet, parent, proposal, invoicing_information).

## Bill-to name follows the invoicing profile

Lines 43-65

Asked as: bill to name wrong; billing name on request; alias instead of official name; client name on invoice

The request shows a bill-to name for the client. The name is the official name or the alias, whichever the invoicing profile chooses (utilities.BillingName). The billing contact is the contact of the invoicing profile.

## Billing contact comes from the invoicing profile

Lines 43-51

Asked as: billing contact missing; who gets the invoice; invoice contact not showing; no billing contact shown

When a project request has an invoicing profile, Turbo shows that profile's contact as the request's billing contact (billing_contact). When the profile has no contact, Turbo shows no billing contact.

## Request rate sheet entries list highest rate first

Lines 66-84

Asked as: rates in wrong order; highest rate at top; rate sheet sorting; overtime rate not showing

The rate sheet entries for a request are sorted from the highest rate first, and show the role, rate and overtime rate. Proposed rates are added by requestProposedRates.

## People shown on requests use profile name and phone

Lines 87-96

Asked as: manager phone number blank; person shows no phone; name missing on request; wrong name for approver

Turbo shows each person on a project request (manager, creator, approver and others) with name, email and phone. When the person has a profile, Turbo builds the name from the profile's given name and surname and takes the phone from the profile. When the person has no profile, Turbo shows the user record's name and an empty phone.

## Review token changes with each submission

Lines 98-112

Asked as: old checklist rejected; review page expired; resubmit keeps old checks; checked boxes not accepted

The review token includes a random submission token that only submit issues. Recall and resubmit of unchanged data gives a new token, so an old review page's checked checklist is not accepted. A draft can have an empty submission token for viewing or printing. Approve and create must reject an empty token, and permission and status checks remain separate.

## Project request review token excludes display-only details

Lines 106-165

Asked as: approval stale after rename; approval outdated after edit; rename made approval invalid; approval reset after profile edit

Turbo builds a review token (a SHA-256 digest) from the saved project request, its submission token (submission_token), its division allocations and the attachment's document hash. The token leaves out the approver, the creator, the rate sheet's code, name and updated time, and the approver's contact details. A rename or profile edit does not make an approval stale, but a change to the rates, revision, active flag, effective date or allocations does.

## Request detail page contents

Lines 169-223

Asked as: request page missing attachment; no approval checklist shown; where is the document link; detail page sections blank

The detail page of a project request shows the request, its division allocations (each with the division's code and name), its attachment with a document link, the linked records, a review token and the current approval checklist for the request's scope. It also shows whether project request approval is enabled. The page loads all of this in one transaction.

## Project request detail page permissions

Lines 208-223

Asked as: why can't I edit; submit button missing; no approve button; cannot recall request

The detail response gives flags for what the viewer can do. All flags except delete need project request approval to be enabled (IsJobRequestApprovalRequired). Only the creator can edit a request in Draft, Recalled or Rejected status. The creator can submit in Draft or Recalled, but not Rejected. The creator can recall in Awaiting approval or Awaiting creation. The assigned approver can approve in Awaiting approval only when their project authorization limit covers the project value. A holder of the job claim can create the job in Awaiting creation. Reject follows canRejectRequest.

## Who can delete a project request

Lines 212-223

Asked as: delete request button; cannot delete my request; remove abandoned request; admin delete request

The creator can delete an editable request (Draft, Recalled or Rejected) while approval is enabled. An admin (the admin claim) can delete an abandoned request (abandonedRequest) even when approval is disabled.

## Project request list queues and required claims

Lines 231-272

Asked as: queue not showing; unknown request queue; no access to approvals list; creation queue missing

Every queue needs an active user. The creation queue needs approval enabled and the job claim. The approval and approved queues need the project authorization claim. The abandoned queue needs the admin claim. The all queue needs the job claim. The default and mine queues need no claim and show the person's own requests. Any other queue gives "Unknown request queue."

## Queue contents for each project request view

Lines 245-269

Asked as: why request not in list; approved requests missing; abandoned requests list empty; my requests view

Creation shows requests in the creation queue (creationQueueFilter). Approval shows requests assigned to the approver in Awaiting approval. Approved shows requests the approver approved that are still in Awaiting creation. Abandoned shows requests last edited on or before the cutoff of JOB_REQUEST_ABANDONMENT_DAYS days ago. Mine shows requests the person created.

## Project request list paging and sorting

Lines 280-309

Asked as: list order wrong; page size limit; rows per page; client name missing in list

The list sorts by newest submitted, then newest created. The default page is 1, and a page below 1 becomes 1. The default page size is 50, and a size below 1 or above 100 becomes 50. Each row carries the client name.

## Who can load project request form options

Lines 320-333

Asked as: You cannot use project requests; can't open request form; no options for new request; access denied project request

Only an active user with the time, job, project authorization or admin claim can load the options. Others see "You cannot use project requests."

## Project request form option lists

Lines 335-363

Asked as: missing client in dropdown; rate sheet not in list; inactive person not shown; branch missing from options

The options list all clients, contacts, invoicing profiles, branches and divisions. The rate sheets list leaves out any sheet that belongs to a request (a request's proposed rates are listed nowhere). Only active users appear as people.

## Which jobs appear as proposals or parent projects

Lines 370-383

Asked as: job not offered as parent; proposal number not listed; subjob not available as parent; P job missing

A job whose number starts with P is a proposal, and it is offered only when its status is In Progress, Submitted or Awarded. Other jobs are offered as parents only when the number looks like 25-123 or 25-1234 (two digits, a hyphen, then three or four digits). Subjob numbers with a suffix are not offered as parents.

## Authorization approver lookup needs a whole-dollar value

Lines 394-423

Asked as: enter positive whole-dollar; approver list empty; who can approve this value; authorization limit approver search

The approver lookup needs a positive whole-dollar project value. Otherwise Turbo says "Enter a positive whole-dollar project value." It lists only holders of the project authorization claim whose limit covers that value, and shows each person's max value. People whose limit does not cover the value are left out.

## Duplicate check lists jobs and pending requests

Lines 430-481

Asked as: duplicate job warning; similar jobs found; pending request duplicate; check for existing jobs

The duplicate check lists existing jobs that match, and skips proposals (numbers starting with P). Only holders of the job claim also see other matching project requests. The request itself is never listed.

## Duplicate check reasons

Lines 484-499

Asked as: why flagged as duplicate; same client PO warning; similar description match; parent project duplicate reason

Turbo compares a request with other records and gives a reason for each match: Same client, Same proposal, Same parent, Same location, Same client PO, Same client reference. An empty field never matches. A description gives Similar description when one contains the other, ignoring case and spaces, and neither is empty. A request whose parent is the other job also gets Parent project.

## Menu badge counts for project workflow

Lines 501-544

Asked as: badge number on menu; approval count not showing; notification count wrong; pending value approvals count

Badges show only for active users. A job claim holder sees the count of requests in the creation queue when approval is enabled. A project authorization holder sees the count of requests awaiting their approval when approval is enabled. They also always see the count of pending value approvals assigned to them (attachments of jobs with status Pending).
