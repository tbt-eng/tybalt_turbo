<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_request_views.go","sha256":"251e3696eb9acfc8f25dda3026f9e3fb6cb533659e7493c21da9fe2fc4083e55","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_request_views.go

## Project request detail shows contact details of related people

Lines 27-96

Asked as: contact details on request page; manager phone number missing; who is the approver contact; related person name and email

The project request page shows each related person (manager, alternate manager, creator, approver) by name, email and phone. Turbo takes the given name, surname and phone from the person's profile when one exists. Without a profile, the phone is empty and the user's own name is used.

## Billing contact comes from the invoicing profile

Lines 43-51

Asked as: billing contact on request; who gets the invoice contact; invoicing profile contact name; billing person not showing

When a project request has an invoicing profile (invoicing_information) that names a contact, the request page shows that contact as the billing contact (billing_contact).

## Bill-to name follows the invoicing profile

Lines 52-65

Asked as: bill to name wrong; invoice name alias; client name on invoice; which name goes on the bill

A request shows a bill-to name chosen by the client's invoicing profile (BillingName): the official name or an alias.

## Rate sheet rates listed highest first

Lines 66-80

Asked as: rate sheet order; rates sorted wrong; highest rate at top; rate list ordering

The rates shown for a request's rate sheet sort by highest rate first.

## Review token ties approval to one submission

Lines 98-156

Asked as: approval lost after resubmit; approval invalid after recall; why did my approval reset; resubmitted request needs approval again

The review token is a SHA-256 digest of the request data, its allocations, its document hash and a random submission token. The approver and creator are left out, and so is the rate sheet code. A recalled and resubmitted request gets a new token. Approval and creation must reject an empty token.

## Review token ignores approver, creator and display-only details

Lines 106-156

Asked as: approval reset after editing contact; changing approver resets approval; approval cleared after small edit; does editing approver undo approval

The review token (requestReviewToken) leaves out the approver and creator. It also leaves out the proposed rates approver's contact details, the source rate sheet's own record and the rate sheet's code. A change to those does not invalidate an earlier approval. A change to the request's business fields, allocations, attachment hash, version or submission token does.

## Project request detail only opens requests you may see

Lines 160-166

Asked as: cannot open project request; request not found for me; access denied on request; can't view request details

The project request detail page loads the request with loadVisibleRequest for the signed-in user. A request that this check does not allow is refused and no detail is returned. The rule for who may see a request is in another file.

## Project request detail shows divisions by code and name

Lines 171-178

Asked as: division name not showing; division id instead of name; allocation division code missing; request page fails loading divisions

The detail page returns the request's hour allocations with each division expanded. The page can show the division's code and name instead of its id. If the expansion fails, the whole page fails.

## Project request page shows what you may do

Lines 187-214

Asked as: why no edit button; buttons missing on request; greyed out approve button; cannot see submit option

The project request detail page returns permission flags. Editing, submitting, recalling, rejecting, approving and creating all require that project request approval is switched on (IsJobRequestApprovalRequired). Deleting a draft the creator can edit also requires it.

## Who can create the job from a project request

Lines 191-214

Asked as: who can make the job; create job from request button; job creation not available; procurement create job

A person with the job claim (procurement) can create the job only while the request status is Awaiting creation.

## Who can delete a project request

Lines 195-214

Asked as: remove project request; who can delete request; delete abandoned request; admin remove draft request

The creator can delete a request they can edit. An admin (the admin claim) can also delete an abandoned request (abandonedRequest). The admin path does not need project request approval to be switched on.

## Who can edit or submit a project request

Lines 203-214

Asked as: cannot edit request; why can't I submit; rejected request cannot submit; change my request after draft

Only the creator can edit a project request, and only while its status is Draft, Recalled or Rejected. The creator can submit a Draft or Recalled request but not a Rejected one (can_submit is false for Rejected).

## Who can recall a project request

Lines 204-214

Asked as: take back request; withdraw my request; pull request back; undo submission

Only the creator can recall a project request, and only while its status is Awaiting approval or Awaiting creation.

## Who can approve a project request from the page

Lines 205-214

Asked as: approve button missing; can I approve this request; approver cannot sign off; approval limit too low

A person can approve only when the status is Awaiting approval, the person is the selected approver, and their project authorization limit covers the project value (ProjectAuthorizationLimit returns no error).

## Project request detail includes current checklist and attachment link

Lines 206-214

Asked as: checklist not showing; attachment link missing; where is the request document; request questions missing

The detail page returns the current request checklist (version and questions) for the request's scope (requestChecklistScope). If the request has an attachment, the page adds a document link of the form /api/job_attachments/{id}/document.

## Project request queues check workflow access first

Lines 222-226

Asked as: request queue refused; cannot load request list; project requests blocked; workflow access list error

Every project request queue first runs the workflow-active check (workflowActive) for the signed-in user. If the check fails, the list is refused before any queue claim is checked. This file does not show the message or who passes the check.

## Project request queues need claims

Lines 228-258

Asked as: empty request queue; no creation queue; cannot see approved list; queue not showing for me

The creation queue needs approval to be switched on and the job claim. The approval and approved queues need the project authorization claim. The abandoned queue needs the admin claim. The all queue needs the job claim. Without the claim the list is refused.

## Approval queue contents

Lines 236-247

Asked as: what is in my approval queue; requests waiting for my approval; approved requests list; which requests do I approve

The approval queue lists requests where the person is the approver and the status is Awaiting approval. The approved queue lists requests where the person is the approver and the status is Awaiting creation.

## Abandoned requests queue cutoff

Lines 248-253

Asked as: old abandoned requests; stale requests list; how old before abandoned; requests not edited for weeks

The abandoned queue lists requests whose last edit (last_edited) is on or before the cutoff of JOB_REQUEST_ABANDONMENT_DAYS days ago.

## Default project request queue is your own

Lines 259-262

Asked as: my requests list; requests I created; unknown queue error; default request list

With no queue given, or queue mine, Turbo lists the requests the person created. This needs no claim. Any other queue name gives "Unknown request queue."

## Project request list paging and order

Lines 271-300

Asked as: request list order; requests per page; how many requests shown; client name missing in list

The list is sorted by newest submitted, then newest created. Pages start at 1. The page size is 50 by default, and 50 is used when perPage is below 1 or above 100. Each row shows the client name.

## Request form options need workflow access

Lines 311-313

Asked as: request form empty; no options on request form; form dropdowns blank; workflow blocks request form

The project request form options run the workflow-active check (workflowActive) before anything else. If the check fails, no options are returned. This file does not show the message.

## Project request options need a claim

Lines 314-324

Asked as: cannot use project requests; request form no access; why no form options; missing permission for request form

Loading the options for the request form needs one of the time, job, project_authorization or admin claims. Otherwise Turbo says "You cannot use project requests."

## Request form options leave out proposed rate sheets

Lines 329-332

Asked as: rate sheet missing from form; proposed rate sheet not listed; which rate sheets to pick; rate sheet choices short

The rate sheet choices exclude the sheets that belong to a project request (proposed rates). Only sheets with no request are listed.

## Request form lists only active people

Lines 344-354

Asked as: inactive person in dropdown; staff not listed; missing user in form; former employee shown

The people choices include only active users.

## Proposal and parent project choices

Lines 361-376

Asked as: proposal not in list; parent project missing; subjob not offered; which job numbers appear

Proposals (numbers starting with P) are offered only when their status is In Progress, Submitted or Awarded. Parent projects are jobs whose number looks like two digits, a dash, then 3 or 4 digits, such as 25-123. Subjobs are not offered.

## Request form options say if approval is required

Lines 377-382

Asked as: is approval needed; approval switched off; form says not enabled; request approval required flag

The form options include an "enabled" flag. It is true when job request approval is required (IsJobRequestApprovalRequired).

## Approver lookup needs workflow access

Lines 387-389

Asked as: no approvers listed; approver lookup empty; who can approve this value; approver search blank

The approver lookup for a project value runs the workflow-active check (workflowActive) before it reads the value. If the check fails, no approvers are listed.

## Approver lookup needs a whole-dollar value

Lines 390-393

Asked as: approver lookup error; value must be whole dollars; project value invalid; cents in project value

Looking up approvers needs a positive whole-dollar project value. Otherwise Turbo says "Enter a positive whole-dollar project value."

## Approver lookup lists those whose limit covers the value

Lines 394-411

Asked as: approver not showing; limit too small for approver; max value of approvers; who has enough authority

Only holders of the project authorization claim whose limit covers the value are listed. Each shows their max_value. Others are left out.

## Duplicate check on a project request

Lines 421-490

Asked as: possible duplicate request; request already exists; duplicate warning; same job warning

Turbo compares the request with existing jobs and skips proposals. Reasons are same client, proposal, parent, location, client PO or client reference (only when set), similar description, and parent project.

## Pending request duplicates only for procurement

Lines 450-471

Asked as: other requests like this; duplicate requests for procurement; see similar pending requests; who sees duplicate list

Other project requests that look like duplicates are listed only for people with the job claim (procurement). The request itself is left out.

## Similar description matching

Lines 482-485

Asked as: similar description warning; description matches another; same wording request; why description flagged

Descriptions match when, ignoring case and outer spaces, one contains the other. Empty descriptions never match.

## Menu badges for project work

Lines 492-535

Asked as: menu count missing; badge number wrong; how many awaiting creation; pending approvals count

Procurement sees a count of requests awaiting creation, only when approval is on. Holders of the project authorization claim see their requests awaiting approval (when approval is on) and their pending value approvals. Inactive users get no badges.
