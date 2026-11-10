<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_request_views.go","sha256":"68c8e64aa7c8b23b488c6b0d782536584b8fa6ec79aa31b9e077731dc2fcbdbe","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_request_views.go

## Project request related records shown with the request

Lines 27-42

Asked as: request details won't load; missing record on request; project request blank error; client or contact not found

Turbo loads the related records for a project request: client, job owner, contact, manager, alternate manager, creator, approver, branch, rate sheet, parent, proposal and invoicing profile. A person shows the name and phone from their profile. If a related record cannot be found, the request detail cannot be loaded.

## Billing contact comes from the invoicing profile

Lines 43-51

Asked as: billing contact missing; wrong invoice contact; who gets the invoice contact

When a request has an invoicing profile with a contact, the related data adds that contact as the billing contact.

## Invoice name comes from the invoicing profile

Lines 52-65

Asked as: wrong name on invoice; client alias on invoice; bill to name wrong

The related data for a request gives a bill_to_name. The invoicing profile chooses whether it is the client's official name or an alias (utilities.BillingName). Without a profile, BillingName gets no profile.

## Rate sheet entries listed highest rate first

Lines 66-81

Asked as: rates listed in wrong order; rate sheet sorting; role rates not showing

A request's related data lists the rate sheet's entries, each with role name, rate and overtime rate, with the highest rate first. A request without a rate sheet gets an empty list.

## Request people show profile name and phone

Lines 84-93

Asked as: name missing on request; no profile name shown; phone number blank

Turbo shows a person's name as the given name and surname from their profile, trimmed. If the person has no profile, it shows the name on the user record and an empty phone.

## Review token ties the review page to one submission

Lines 95-140

Asked as: old checklist still accepted; checklist ticked after change; review stale after recall; submission token reset

The review token (a SHA-256 digest) covers the request's version, its submission_token, its business fields, its related records without approver and creator, its division allocations (division, hours, active) and the attachment's document hash. A change to any of these gives a different token. Recalling and resubmitting gives a new submission token, so an old checked checklist is not accepted. A draft can have an empty submission token for viewing, but approve and create must reject an empty one.

## Review token leaves out approver and creator

Lines 103-133

Asked as: changing approver resets review; approver change breaks checklist; creator change review token

The review token (requestReviewToken) ignores the approver and the creator of the request. Changing either of them does not change the token. The token does change with the request version, the submission token, the business fields, the allocations (division, hours, whether the division is active) and the attachment's document hash.

## Project request detail shows checklist and attachment link

Lines 151-198

Asked as: attachment link missing; where is the request document; checklist not showing on request; division code missing on request

The request detail page shows the current checklist version and questions for the request's scope. It also shows each allocation's division with its code and name. If the request has an attachment, Turbo adds a link to its document (/api/job_attachments/<id>/document). All of this is visible only to people allowed to view that request.

## Project request detail permissions need approval switch on

Lines 171-198

Asked as: request buttons missing; cannot edit request; actions gone approval turned off

On the request detail page, the edit, submit, recall, reject, approve and create permissions all need the job request approval requirement (IsJobRequestApprovalRequired) to be on. Delete of an editable request needs it too. An admin can still delete an abandoned request when it is off.

## Who can create the job from a project request

Lines 175-198

Asked as: create job button missing; cannot create job from request; procurement create job option

The detail page shows can_create only when approval is on, the person has the job claim (procurement) and the request is Awaiting creation.

## Who can delete a project request

Lines 179-198

Asked as: cannot delete request; delete abandoned request; remove my draft request

The creator can delete an editable request (Draft, Recalled or Rejected) while approval is on. An admin (the admin claim) can also delete a request that is abandoned (abandonedRequest), even when approval is off.

## Reject permission is decided by canRejectRequest

Lines 183-198

Asked as: reject button missing; cannot reject request; who can reject

The detail page shows can_reject only when approval is on and canRejectRequest allows this person. This file does not show that function's rules.

## Project request detail: who can edit, submit, recall, delete

Lines 187-198

Asked as: can't edit my draft; submit button missing; recall my request; cannot submit rejected request

On the project request detail page, only the creator can edit a request in Draft, Recalled or Rejected status (can_edit). A creator can submit it (can_submit) in Draft or Recalled status, but not in Rejected status. A creator can recall it (can_recall) while it is Awaiting approval or Awaiting creation. All of these need project request approval to be on (enabled).

## Who can approve a project request on the detail page

Lines 189-198

Asked as: approve button missing; cannot approve request; approver can't approve

The detail page lets a person approve (can_approve) only when approval is on, the request is Awaiting approval, the person is the chosen approver, and their project authorization limit covers the request's project value (ProjectAuthorizationLimit returns no error).

## Request queues and who can open them

Lines 206-241

Asked as: unknown request queue; approval queue empty; creation tab missing; requests tab not showing

The 'mine' queue (the default) shows the person's own requests. The 'creation' queue needs project request approval to be on and the job claim (procurement). The 'approval' queue needs the project authorization claim and shows requests Awaiting approval where the person is the approver. The 'all' queue needs the job claim. An unknown queue gives 'Unknown request queue.'. Every queue first requires an active user.

## Abandoned requests queue

Lines 225-230

Asked as: abandoned requests list; old unused requests; stale requests view

The 'abandoned' queue needs the admin claim. It shows requests last edited on or before the cutoff, which is JOB_REQUEST_ABANDONMENT_DAYS days ago.

## Request list paging and sort order

Lines 248-277

Asked as: request list order; how many requests per page; requests sorted wrong

The request list is sorted by newest submitted, then newest created. The default page is 1 and a page below 1 becomes 1. The default page size is 50, and a perPage below 1 or above 100 becomes 50. The list shows the client name for each request.

## Request form options need a project claim

Lines 286-301

Asked as: cannot use project requests; no access to request form; project request access denied

The request form options page requires an active user. Turbo refuses with "You cannot use project requests." unless the user holds the time, job, project authorization or admin claim.

## Request form options need a work claim

Lines 288-301

Asked as: project requests forbidden; missing time claim; request options access denied

Only a person with the time, job, project_authorization or admin claim can load the project request options. Anyone else gets 'You cannot use project requests.'. The person must also be an active user.

## Request form options list reference data

Lines 302-309

Asked as: dropdown lists empty; no divisions in request form; request form choices

The request form options return all clients, contacts, invoicing profiles, branches, rate sheets and divisions, plus the approval-required setting (enabled).

## Options list only active people

Lines 314-324

Asked as: inactive user in people list; former staff in dropdown; leavers showing as people

The options list of people includes only active users.

## Which proposals and parent projects are offered

Lines 331-346

Asked as: proposal missing from dropdown; parent job not offered; subjob as parent option

The options offer a proposal (number starts with P) only when its status is In Progress, Submitted or Awarded. They offer a job as a parent only when its number looks like 25-123 or 25-1234 (two digits, a dash, three or four digits). Subjobs are not offered as parents.

## Approver list needs a positive whole-dollar value

Lines 357-363

Asked as: invalid project value; positive whole-dollar project value; approver list error

The authorization approvers list needs an active user. It refuses with "Enter a positive whole-dollar project value." when the value is missing, not a number, or not a positive whole amount.

## Authorization approvers list by project value

Lines 360-382

Asked as: approvers for project value; who can approve this amount; approver list empty

The approver list needs a whole-dollar project value above zero, or it gives 'Enter a positive whole-dollar project value.'. It lists only people with the project_authorization claim whose limit covers that value, with their max_value. People whose limit does not cover it are left out.

## Approvers list only those whose limit covers the value

Lines 364-382

Asked as: approver not listed; approver limit too low; approver max value

Turbo lists people with the project authorization claim, ordered by user id. It leaves out anyone whose project limit (ProjectAuthorizationLimit) does not accept the value. Each person listed shows their max_value.

## Duplicate check on a request that you can view

Lines 385-419

Asked as: possible duplicate jobs; duplicate job check; same job twice

The duplicates list works only for a request the user may view (loadVisibleRequest). It checks all jobs except proposals (numbers starting with P) and shows the matching reasons, the parent number, the client name and the manager name.

## Possible duplicate jobs for a request

Lines 396-460

Asked as: why flagged as duplicate; similar job description match; duplicate warning reasons

The duplicate check compares a request with existing jobs, skipping proposals (numbers starting with P). It reports the reasons: same client, proposal, parent, location, client PO or client reference (only when the request's value is not empty), a similar description, and parent project. A similar description means that one description contains the other, ignoring case and spaces, and neither is empty.

## Only procurement sees duplicate pending requests

Lines 420-441

Asked as: pending duplicate requests hidden; duplicate request not shown; other requests not listed

The duplicate check lists other project requests that match only for a person with the job claim (procurement). Others see only jobs. The request itself is left out.

## Duplicate match reasons

Lines 445-460

Asked as: duplicate match reason; matched on client PO; similar description flagged

A request matches another record on same client, proposal, parent, location, client PO or client reference when the value is not empty and is equal. A description gives Similar description when both are not empty and one contains the other, ignoring case and spaces. A request whose parent is that job gives Parent project.

## Menu badge counts for project workflow

Lines 462-505

Asked as: menu count missing; approval count badge; creation badge number wrong

Procurement (job claim) sees a creation count only when project request approval is on. A person with the project_authorization claim sees an approval count (only when approval is on) and a count of pending value approvals assigned to them. Inactive users get no badges.

## Badges skipped for inactive users

Lines 462-491

Asked as: no menu badges showing; badges missing for me; inactive user badges

Project workflow menu badges are not added when the user is not active. The counts also depend on claims: creation needs the job claim, approval needs the project authorization claim. The creation and approval request counts show only when the approval requirement is on.

## Value approvals badge counts pending attachments

Lines 499-503

Asked as: value approvals badge; pending attachment count; attachment approval menu count

A person with the project authorization claim sees a badge /jobs/value-approvals. It counts pending job attachments assigned to them as approver that belong to a job. This badge does not depend on the approval requirement.
