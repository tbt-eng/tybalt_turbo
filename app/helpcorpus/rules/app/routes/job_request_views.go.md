<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_request_views.go","sha256":"5eee8be78009d41038dfa36fd7c4eceebecc904e15ff2e65bfebc43e2ca3fc4f","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_request_views.go

## Request detail includes the records the request links to

Lines 27-42

A project request's detail loads the linked client, job owner, contact, manager, alternate manager, creator, approver, branch, rate sheet, parent, proposal and invoicing profile (client, job_owner, contact, manager, alternate_manager, creator, approver, branch, rate_sheet, parent, proposal, invoicing_information). Turbo loads only the links the request has filled in. If a linked record cannot be found, the detail does not load.

## Linked people show their profile name

Lines 37-40

Each linked person (manager, alternate manager, creator, approver) is shown by the workflow person view, not as a raw user record. The next rule describes the name and phone it uses.

## Request shows the invoicing profile's contact as billing contact

Lines 43-51

When a request has an invoicing profile with a contact, the related data includes that contact as the billing contact.

## Bill-to name follows the invoicing profile

Lines 52-65

The bill-to name on a project request comes from utilities.BillingName, using the client and the invoicing profile (if one is set). The name is the official name or the alias, as the profile chooses.

## Related rate sheet lists its rates by role

Lines 66-80

When a request has a rate sheet, the related data lists its entries sorted by role. Each entry gives the role name, the rate, the overtime rate and the updated time.

## Staff names come from the profile

Lines 83-92

A person's name in workflow data comes from their profile: given name plus surname, trimmed. The person's phone also comes from the profile. Without a profile, the user's own name is used and the phone is empty.

## Review token ties approval to one submission

Lines 102-139

The review token is a SHA-256 digest of the request's version, its submission token, its business fields, its related records without approver and creator, its allocations and the attachment's document hash. Any change to these changes the token. Each division allocation also contributes its hours and whether the division is active.

## Project request detail shows what the viewer can do

Lines 143-197

The detail page for a project request returns permission flags for the viewer. It also returns the review token, the request's checklist (version and questions), the allocations, the attachment and the related records. Every flag except can_delete is false when project request approval is not required (enabled).

## Who can create the job from a request

Lines 174-197

A person with the job claim can create the job from a project request, but only when its status is Awaiting creation.

## Only admins delete abandoned project requests

Lines 178-197

A person with the admin claim can delete a project request only when abandonedRequest says it is abandoned. Delete is available whether or not approval is required.

## Reject permission comes from canRejectRequest

Lines 182-197

The reject flag on a project request page is true when approval is required and canRejectRequest allows the viewer to reject it. This file does not say what canRejectRequest checks.

## Who can edit or submit a project request

Lines 186-197

Only the creator can edit a project request, and only when its status is Draft, Recalled or Rejected. The creator can submit it (can_submit) in Draft or Recalled, but not in Rejected.

## Who can recall a project request

Lines 187-197

Only the creator can recall a project request, and only while its status is Awaiting approval or Awaiting creation.

## Who can approve a project request

Lines 188-197

A person can approve a project request only when its status is Awaiting approval and they are the chosen approver. Their project authorization limit (ProjectAuthorizationLimit) must also cover the request's project value.

## Project request queues and who may open them

Lines 205-240

The creation queue needs the job claim and approval to be required. The approval queue needs the project authorization claim. The abandoned queue needs the admin claim. The all queue needs the job claim. The mine queue (the default) needs no claim. Any other queue gives the error "Unknown request queue." The user must also pass workflowActive.

## Approval queue shows requests awaiting the approver

Lines 219-223

The approval queue lists requests where the viewer is the approver and the status is Awaiting approval.

## Abandoned queue uses last edit cutoff

Lines 224-229

The abandoned queue lists requests whose last_edited time is on or before the cutoff. The cutoff is JOB_REQUEST_ABANDONMENT_DAYS days before now.

## My requests queue shows requests the viewer created

Lines 235-236

The mine queue, which is the default, lists only requests the viewer created.

## Request list paging and order

Lines 247-276

The list is sorted by submitted time, newest first, then by created time and then by id. The page size defaults to 50, and a perPage below 1 or above 100 also gives 50. A page below 1 becomes page 1. Each row includes the client name.

## Request options need a project-related claim

Lines 285-300

A user needs one of the time, job, project_authorization or admin claims to load the project request options. Otherwise Turbo says "You cannot use project requests."

## Request options list only active users

Lines 309-323

The people list in the options includes only active users.

## Proposal and parent choices in request options

Lines 330-345

Jobs whose number starts with P are offered as proposals only when their status is In Progress, Submitted or Awarded. Jobs whose number matches the pattern two digits, a dash, then three or four digits (for example 25-123) are offered as parents. Subjobs and other numbers are not offered as parents.

## Request options report whether project request approval is required

Lines 346-350

The project request options response includes an enabled flag. The flag comes from the job request approval setting (IsJobRequestApprovalRequired).

## Approver lookup is for active workflow users only

Lines 354-358

The approver lookup for a project value first checks that the signed-in user is allowed to use the project workflow (workflowActive). If that check fails, the lookup returns that error and lists no approvers. This file does not show the message.

## Approver lookup needs a positive whole-dollar value

Lines 359-362

The approver lookup refuses a value that is not a positive whole-dollar amount. The message is "Enter a positive whole-dollar project value."

## Approver lookup lists only approvers whose limit covers the value

Lines 363-381

The lookup goes through users with the project_authorization claim. It skips anyone for whom ProjectAuthorizationLimit returns an error for the value. It shows each remaining person with their max_value.

## Duplicate check skips proposals and shows pending requests only to job staff

Lines 390-441

The duplicate check lists existing jobs that match and skips jobs whose number starts with P. It lists other pending project requests only for a person with the job claim, and it leaves out the request itself.

## Duplicate check match reasons

Lines 444-459

A project request matches another record when client, proposal, parent, location, client PO or client reference are equal and not empty. Descriptions also match when, trimmed and ignoring case, one contains the other. Both must be non-empty. The match is also flagged "Parent project" when the other record is the request's parent.

## Badge counts for project workflow

Lines 461-504

Badges appear only for a user who passes workflowActive. A person with the job claim sees the creation count when approval is required. A person with the project_authorization claim sees the approval count when approval is required. That person also always sees the count of pending value approvals assigned to them.
