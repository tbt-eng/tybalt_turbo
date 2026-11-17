<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_request_views.go","sha256":"e1e62178273ce46039e4dd188298d18432ec4f89706ea96910cfa83542852b26","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_request_views.go

## Request detail shows linked records and people

Lines 27-42

Asked as: request page missing client; linked records on request; request details empty fields; request detail error record not found

A project request's detail page includes the records the request links to. These are the client, job owner, contact, manager, alternate manager, creator, approver, branch, rate sheet, parent project, proposal and invoicing profile (client_invoicing_information). A link that is empty is left out. If a linked record cannot be found, the detail fails.

## Billing contact comes from the invoicing profile

Lines 43-51

Asked as: billing contact missing on request; who gets billed on request; invoicing profile contact; billing contact not showing

When a project request has an invoicing profile, and that profile names a contact, the request detail adds that contact as the billing contact (billing_contact). A request with no invoicing profile, or a profile with no contact, has no billing contact.

## Bill-to name and rate order on request

Lines 52-80

Asked as: bill to name wrong; client alias on request; rates listed in wrong order; rate sheet order request

The request shows the client name that the invoicing profile selects (official or alias) as the bill-to name. Its rate sheet rates are listed highest rate first.

## Person names come from the profile, with a fallback

Lines 87-96

Asked as: manager name wrong on request; approver shows user name; phone number missing on request; profile name not showing

When a request page shows a person, such as a manager or approver, it uses the name on that person's profile: given name and surname joined and trimmed. It also shows the profile's phone number. A person with no profile shows the user account's name and an empty phone number.

## Review token ties approval to one submission

Lines 98-161

Asked as: approval gone after edit; approval stale after change; why did approval reset; changed request needs re-approval

The review token is a SHA-256 digest of the request's version, its submission token, its business fields, related records, division hours and the attachment's document hash. A change to any of them makes an earlier approval stale. The digest leaves out the approver, the creator, the rate sheet's code, name and updated time, and the rates approver and source sheet records, so renames or profile edits do not make it stale.

## Request detail opens only for requests the viewer may see

Lines 166-192

Asked as: cannot open request; can't see request; request access denied; who can view request

Turbo opens a project request's detail page only after loadVisibleRequest accepts the viewer and the request. The rules for who may see a request are in that function, which is not in this file. The detail then lists the request's attachment and its division hours, and shows each division by code and name.

## Request detail includes a review token

Lines 185-192

Asked as: review link not working; no review token; review button missing; request review code

The detail page for a project request returns a review token (review_token), which requestReviewToken builds from the request's current saved data and its submission token. The page uses this token for the review and approval steps.

## Request detail permission flags

Lines 193-220

Asked as: no edit button on request; buttons missing on request; cannot submit or approve button; request actions greyed out

The project request detail page returns which actions the viewer may take. Editing, submitting, recalling, rejecting, approving and creating need project request approval to be on (the IsJobRequestApprovalRequired setting). Deleting needs it too, unless an admin deletes an abandoned request.

## Who can create the job from a request

Lines 197-220

Asked as: cannot create job from request; create job button missing; job creation not allowed; request awaiting creation job

A person with the job claim (procurement) can create the job from a request only while its status is Awaiting creation.

## Who can delete a project request

Lines 201-220

Asked as: cannot delete request; delete request button missing; remove old abandoned request; who can delete request

The creator can delete an editable request when approval is on. An admin (the admin claim) can delete an abandoned request (abandonedRequest).

## Who can edit or submit a project request

Lines 209-220

Asked as: cannot edit request; submit button not working; rejected request resubmit; edit request greyed out

Only the creator of a project request can edit it, and only while its status is Draft, Recalled or Rejected. The creator can submit it in Draft or Recalled status, but not when it is Rejected (can_submit).

## Who can recall a project request

Lines 210-220

Asked as: withdraw request for changes; take back request; cannot recall request; undo request submission

Only the creator can recall a project request, and only while its status is Awaiting approval or Awaiting creation.

## Who can approve a project request

Lines 211-220

Asked as: cannot approve request; approve button missing; not my request to approve; approval limit too low

A person can approve a project request only if the request is Awaiting approval, that person is its selected approver, and their project authorization limit covers the request's project value (ProjectAuthorizationLimit returns no error).

## Request queues and who may open them

Lines 228-270

Asked as: request queue empty; unknown request queue error; approval list not showing; abandoned requests list access

The creation queue needs approval to be on and the job claim. The approval and approved queues need the project authorization claim and show only the viewer's own requests as approver, with status Awaiting approval or Awaiting creation. The abandoned queue needs the admin claim. The all queue needs the job claim. The mine queue (also the default) shows the viewer's own requests. Any other queue name gives "Unknown request queue."

## Abandoned queue cutoff

Lines 254-259

Asked as: how old is abandoned; abandoned request days; requests considered abandoned; stale request cutoff

The abandoned queue lists requests whose last edit is on or before JOB_REQUEST_ABANDONMENT_DAYS days ago.

## Request list paging and order

Lines 277-306

Asked as: request list order; requests per page; request list sorting; wrong order on request list

The request list shows newest submitted first, then newest created. The default page size is 50, and a page size below 1 or above 100 also gives 50. A page below 1 becomes page 1.

## Who can load request form options

Lines 317-330

Asked as: cannot open request form; you cannot use project requests; request form options not loading; no access to new request

Only active users with the time, job, project_authorization or admin claim can load the request form options. Others get "You cannot use project requests."

## Request form option lists

Lines 331-382

Asked as: missing proposal in dropdown; parent project not listed; rate sheet missing from form; which users appear in request form

The options leave out rate sheets that belong to a request (proposed rates). They list only active users. Proposals offered are jobs whose number starts with P and whose status is In Progress, Submitted or Awarded. Parent projects offered are jobs whose number matches two digits, a dash and three or four digits, such as 25-123.

## Authorization approvers lookup by value

Lines 391-419

Asked as: who can approve this value; approver list empty; project value approver lookup; enter positive whole dollar value

The approver lookup needs a positive whole-dollar project value, or it says "Enter a positive whole-dollar project value." It lists the people with the project authorization claim whose limit covers that value, each with their maximum value. People whose limit does not cover it are left out.

## Duplicate check for a request

Lines 421-496

Asked as: duplicate job warning; possible duplicate request; similar job already exists; same client PO warning

The duplicate check compares a request with existing jobs, but skips proposals (numbers starting with P). It flags the same client, proposal, parent, location, client PO or client reference, a similar description, and a parent project. The description is similar when either one contains the other, ignoring case and spaces; an empty description never matches.

## Pending requests in duplicate check need job claim

Lines 456-478

Asked as: cannot see possible duplicates; pending requests not showing; duplicate check missing requests; who sees duplicate requests

Only people with the job claim also see other project requests that look like duplicates. Others see matching jobs only.

## Project workflow badge counts

Lines 498-541

Asked as: badge count wrong; notification number on workflow; pending approvals count; creation count not showing

Badge counts apply only to active users. A person with the job claim sees a creation count when approval is on. A person with the project authorization claim sees an approval count when approval is on, and always sees a count of pending value approvals assigned to them for jobs.
