<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_request_views.go","sha256":"a40f61ddaaf69f01dfab9aa116cddf4a59f1134552a27ab6ae4eea05f6cc192a","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_request_views.go

## Project request detail loads related records for display

Lines 26-50

Asked as: project request details page missing client; job owner not showing on request; billing contact missing on request page; request page related records empty

A project request's detail page includes the related client, job owner, contact, manager, alternate manager, creator, approver, branch, rate sheet, parent job, proposal and invoicing profile. People (users) show their profile name, email and phone when a profile exists. If an invoicing profile has a contact, that contact appears as the billing contact (billing_contact).

## Bill-to name follows the invoicing profile

Lines 42-64

Asked as: bill-to name wrong on request; invoice name shows alias; billing name on project request; which name goes on invoice

A project request's related data gives a bill-to name (bill_to_name). Turbo takes it from the client and the invoicing profile (utilities.BillingName), which chooses the official name or the alias. The invoicing profile's contact is added as the billing contact.

## Request rate sheet entries are listed highest rate first

Lines 65-83

Asked as: rate sheet order on request; highest rate listed first; request rates in wrong order; overtime rate not showing

The rates of a request's rate sheet are listed with the highest rate first. Each entry shows the role name, rate and overtime rate. Proposed rates are added by requestProposedRates.

## Person details use profile name and phone when available

Lines 86-95

Asked as: person shows no phone number; wrong name on request person; profile name missing on request; user name shown instead of profile

Turbo shows a person on a project request with the name, email and phone from their profile (given name and surname, phone). If the person has no profile, it shows the user's name and email with an empty phone.

## Review token binds one submission and the attached document

Lines 97-201

Asked as: approval stale after recall; old review page checklist; resubmitted request still shows old checks; changes after approval not counted

The review token (requestReviewToken) hashes the request version, the submission token, the business fields, the related records, the division hours, and the attachment's document hash. Recall and resubmit issue a new submission token, so an old review page's checked checklist no longer matches. A draft can have an empty submission token for viewing. Approve and create reject an empty token.

## Review token excludes approver, creator and rate sheet display data

Lines 132-158

Asked as: rename rate sheet approval stale; rate sheet change makes approval stale; rate sheet code changed approval; approver change affects approval

The review token of a project request leaves out the approver, the creator, the proposed-rates approver's contact details, the source rate sheet, and the missing-role rows. It also leaves out the rate sheet's code, name and updated time. Renaming a sheet or changing its code does not make the approval stale. A change to the revision, the effective date, the active flag or any rate does.

## Project request review token ignores proposal award

Lines 159-170

Asked as: proposal awarded after approval; awarding proposal resets approval; proposal status changed approval stale; approval lost when proposal awarded

The review token (digest) of a project request leaves out the linked proposal's status, updated time and _imported mark. Awarding the proposal after approval does not make the approval stale. Any other edit to the proposal does. The legacy digest keeps the whole proposal record.

## Project request detail visible only to permitted viewers

Lines 205-235

Asked as: cannot open project request; can't see request details; request detail access denied; attachment link missing on request

Opening a project request's detail loads it through a visibility check for the signed-in person (loadVisibleRequest). The check's rules are in another file. The page also returns the attachment with a document link, the allocations with expanded divisions, the review token, the current checklist, and whether approval is required (enabled).

## Request detail permissions

Lines 244-273

Asked as: what can I do with request; buttons missing on project request; who can approve request; why can't I submit request

All the permissions need project request approval to be enabled (IsJobRequestApprovalRequired). The creator can edit a request in Draft, Recalled or Rejected status. The creator can submit it unless it is Rejected. The creator can recall it while it is Awaiting approval or Awaiting creation. Only the selected approver can approve it, in Awaiting approval, and only when their project authorization limit covers the project value. A job claim holder (procurement) can create the job in Awaiting creation status. Rejecting depends on canRejectRequest.

## Who can delete a project request

Lines 248-273

Asked as: delete project request; remove abandoned request; cannot delete request; admin delete old request

The creator can delete an editable request (Draft, Recalled or Rejected) while approval is enabled. An admin (the admin claim) can delete an abandoned request (abandonedRequest), even when approval is not enabled.

## Request edit and submit allowed for creator in draft states

Lines 248-248

Asked as: edit rejected request; submit button missing on request; cannot submit rejected request; edit request after rejection

A project request is editable only by its creator, and only while its status is Draft, Recalled or Rejected. Editing and submitting both need approval to be required (enabled). A Rejected request can be edited but not submitted (can_submit is false).

## Request recall, reject, approve and create permissions

Lines 249-273

Asked as: recall request to edit; approve button not working; who can create job from request; reject request permission

The creator can recall a request that is Awaiting approval or Awaiting creation. Only the selected approver whose project limit covers the project value can approve, and only while the status is Awaiting approval. A committer of jobs (the job claim) can create the job only while the status is Awaiting creation. All of these need approval to be required.

## Approver told why Approve is unavailable

Lines 259-272

Asked as: approve button disabled; why can't I approve; approval limit too low; no project approval authority message

The selected approver of a request in Awaiting approval sees a reason when their authority does not cover it. The reason is "Your approval limit does not cover this project value." If they lack the project authorization claim, it is "You do not have project approval authority." Nobody else receives this reason.

## Project request queues and who can open them

Lines 281-323

Asked as: request queue empty; unknown request queue error; approval queue missing requests; abandoned requests list not showing

The creation queue needs approval enabled and the job claim (procurement). The approval and approved queues need the project authorization claim and show only requests where the person is the approver: Awaiting approval, or Awaiting creation. The abandoned queue needs the admin claim. The all queue needs the job claim. The mine queue, which is the default, shows the person's own requests. Any other name gives "Unknown request queue." An inactive person is refused first (workflowActive).

## Abandoned queue cutoff

Lines 307-312

Asked as: when request counts as abandoned; abandonment days cutoff; old requests in abandoned list; edited exactly on cutoff

The abandoned queue lists requests whose last edit is at or before the cutoff of JOB_REQUEST_ABANDONMENT_DAYS days ago. A request edited exactly at the cutoff counts as abandoned.

## Request list paging and order

Lines 330-359

Asked as: request list order; requests per page wrong; how many requests per page; request list paging problem

The request list shows the newest submitted first, then the newest created. The default page size is 50. A perPage below 1 or above 100 falls back to 50. A page below 1 becomes page 1. Each row shows the client name.

## Request form options need a work claim

Lines 370-383

Asked as: cannot use project requests; project request form not loading; no access to request form; request form permission denied

Only a person with the time, job, project_authorization or admin claim can load the project request form options. Anyone else gets "You cannot use project requests."

## Request form option lists

Lines 386-450

Asked as: missing rate sheet in request form; which users appear on request; proposal missing from request dropdown; branch list in request form

The form lists the invoicing profiles, branches, divisions and rate sheets, but not rate sheets that hold a request's proposed rates. Users are the people whose admin profile is active. A user is named from their profile when one exists. Proposals are jobs whose number starts with P and whose status is In Progress, Submitted or Awarded.

## Request form options report whether approval is required

Lines 451-456

Asked as: is request approval turned on; approval setting on request form; approval required flag; request approval switched off

The project request form options response includes a flag (enabled) that says whether project request approval is required. Turbo reads it from the job request approval setting (IsJobRequestApprovalRequired) when the form loads.

## Authorization approvers lookup needs an active user

Lines 459-463

Asked as: approver lookup refused; cannot find approvers; inactive user approver lookup; approver list not loading

Before Turbo lists the approvers who can authorize a project value, it checks that the caller is an active user (workflowActive). Turbo refuses the lookup when that check fails. The check itself is in another file.

## Authorization approvers lookup

Lines 464-486

Asked as: whole dollar project value error; who can authorize this amount; approver limit amount too high; list of approvers for value

The approver lookup needs a positive whole-dollar value. Otherwise it gives "Enter a positive whole-dollar project value." It lists the people with the project authorization claim whose limit covers that value (utilities.ProjectAuthorizationLimit). It shows each person's maximum value. It skips the people whose limit does not cover it.

## Project workflow badge counts

Lines 489-532

Asked as: project badge count wrong; pending requests count missing; no project badges showing; approval count not showing

Inactive people get no project badges. With approval enabled, a person with the job claim sees a count of requests awaiting creation. A person with the project authorization claim sees a count of requests awaiting their approval, shown only when approval is enabled. That person also sees a count of pending project value approvals assigned to them. That count shows even when request approval is disabled.
