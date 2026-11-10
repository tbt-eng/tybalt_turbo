<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_requests.go","sha256":"db46724b2828255608558578e159fbcbbaa11295c60af576bae2e39c2da3f5f8","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_requests.go

## Abandonment period constant

Lines 27-27

Asked as: when request is abandoned; days until abandoned; request goes stale; abandonment period days

A request counts as abandoned after JOB_REQUEST_ABANDONMENT_DAYS days (from utilities) without edits.

## Sign in required for project request routes

Lines 46-74

Asked as: not signed in error; sign in message request; login needed for requests; attachment needs login

All project request and attachment routes require a signed-in user. An unauthenticated caller gets "Sign in to continue."

## Validation errors show a generic message

Lines 76-87

Asked as: check the request fields; form validation error; request field errors; generic request error message

Field validation failures on request routes return "Check the request fields." with the field details.

## Active staff member required for request actions

Lines 93-118

Asked as: not an active staff; inactive staff cannot request; missing claim error; staff access request denied

A person who is not signed in sees 'Sign in to continue.' A person who is not an active staff member sees 'An active staff member is required.' A missing claim gives 'The <claim> claim is required.'

## Project requests can be switched off

Lines 119-128

Asked as: project requests turned off; cannot create project request; project requests disabled; feature switched off

When project requests are not enabled (utilities.IsJobRequestApprovalRequired is false), Turbo refuses to create, edit, submit, recall, reject, approve or create a job from a request. It shows 'Project requests are disabled.'

## Abandoned request after the abandonment period

Lines 129-131

Asked as: request edited exactly at cutoff; old draft counts abandoned; last edit cutoff date; request abandoned by age

A project request counts as abandoned when its last edit is on or before the cutoff of JOB_REQUEST_ABANDONMENT_DAYS days ago. A request edited exactly at the cutoff is already abandoned.

## Who can view a project request

Lines 132-163

Asked as: cannot see my request; who can view request; request not found; cannot view this request

The creator can view their own request. Anyone with the job claim (procurement) can view any request. The selected approver with the project_authorization claim can view it. An admin can view only abandoned requests. Anyone else sees 'You cannot view this request.' A missing request gives 'Project request not found.'

## Only the creator can edit, and only in editable states

Lines 164-179

Asked as: cannot edit request; someone else's request edit; recall before editing; only creator edits

Only the creator can edit a project request, and only when it is Draft, Recalled or Rejected. Others see 'Only the creator can edit this request.' In any other state the creator sees 'Recall this request before editing it.'

## Stale version is refused

Lines 191-198

Asked as: request changed reload; version mismatch save; someone else edited request; outdated request version

If the version the caller sends differs from the request's current version, Turbo refuses with 'This request changed. Reload it before continuing.'

## Editing a request resets approval and token

Lines 199-236

Asked as: edit resets approval; approval cleared after edit; authorization PDF pending again; why rejected request returns draft

A saved business edit adds 1 to the version, sets last_edited, and clears the submission token. A Rejected request returns to Draft and its rejection details clear. A linked authorization PDF record returns to Pending with its approval cleared, copies the new project value and gets a new version.

## Only whole dollar project values

Lines 237-252

Asked as: project value cents; decimal project value; project value not whole; whole dollar amount

Project value must be a finite whole number, zero or more, up to 9007199254740991. Otherwise Turbo says 'Project value must be a whole CAD dollar amount.'

## Only listed fields can be edited

Lines 241-250

Asked as: field is not editable; cannot change this field; field locked on request; non editable field error

A save with a field outside the editable list gives 'Field is not editable: <field>'.

## Request notes limits

Lines 253-264

Asked as: notes too short; notes character limit; notes must be text; request notes length

Notes must be text. Turbo trims spaces. Notes must be 10 to 1000 characters, or blank. Otherwise it says 'Notes must be 10 to 1000 characters, or blank.' Non-text gives 'Notes must be text.'

## Fee basis must be a true or false choice

Lines 265-270

Asked as: select fee basis; time and materials error; fee basis missing; fee basis not chosen

Setting time_and_materials needs a true or false value, else 'Select a fee basis.' Setting it marks the fee basis as chosen.

## Draft checks on every save

Lines 278-289

Asked as: project status active closed; select project status; client reference invalid; approver limit exceeded

Project status must be Active or Closed, or Turbo says 'Select Active or Closed for the project status.' Client references are validated. A selected approver's project limit is checked against the project value (ProjectAuthorizationLimit).

## Parent project rules

Lines 290-301

Asked as: parent project invalid; subjob client mismatch; parent must be top level; parent project not found

A parent must exist ('Select a valid parent project.') and be top-level, with a number like 25-123 or 25-1234 ('Only a top-level project can be a parent.'). The request's client must match the parent's client ('The subjob must use the parent project's client.').

## Linked proposal rules

Lines 302-315

Asked as: linked proposal invalid; proposal wrong client; proposal status not awarded; proposal number starts P

A linked proposal must exist ('Select a valid proposal.'), have a number starting with P and the same client ('Select a proposal for this client.'). Its status must be In Progress, Submitted or Awarded ('Select an In Progress, Submitted, or Awarded proposal.').

## Completeness checks before submit, approve and create

Lines 318-373

Asked as: cannot submit missing fields; required field missing; upload authorization pdf; set proposal awarded

These fields must be non-blank: description, location, client, contact, manager, branch, rate_sheet, project_award_date, invoicing_information, approver, fn_agreement_answer ('Required field: <field>'). A fee basis must be chosen. Project value must be above zero. A linked proposal must be Awarded ('Set the linked proposal to Awarded before continuing.'). Contact emails, allocations and an uploaded authorization PDF are required ('Upload an authorization PDF.').

## Allocation rules

Lines 374-398

Asked as: invalid allocation hours; division not active; allocation hours decimal; select valid division

Each division allocation needs hours that are whole numbers and a valid, active division. Messages: 'Invalid allocation hours.', 'Allocation hours must be whole numbers.', 'Select a valid division.', 'Select an active division.' Allocations are optional on a draft save but required to submit.

## Request save runs only while project requests are enabled

Lines 406-409

Asked as: save fails project requests off; nothing saved disabled; cannot save while disabled

Creating or editing a project request first checks that project requests are enabled. If they are not, Turbo refuses with "Project requests are disabled." Nothing is saved.

## Creating a request

Lines 412-426

Asked as: create new project request; start project request; new request draft

Creating a project request needs the time claim. The new request starts as Draft with project status Active and version 1.

## Only staff with the time claim can start a project request

Lines 412-415

Asked as: time claim required; cannot start project request; no time claim access

Creating a new project request requires an active staff member who holds the time claim. Without it Turbo refuses with "The time claim is required." Editing or submitting later also needs the creator to be the one acting.

## New project request starts as an Active draft at version 1

Lines 417-426

Asked as: new request defaults; initial request status; request version starts 1

A new project request gets its creator from the signed-in user, status Draft, project status Active, version 1 and the current time as last edited.

## Editing an existing request needs visibility, editability and current version

Lines 427-438

Asked as: cannot save request changes; save existing request refused; editing request not allowed

Saving changes to an existing request requires that the person can view it, is the creator in an editable state (Draft, Recalled or Rejected), and sends the current version. Otherwise Turbo refuses with the matching message.

## Allocation changes are validated only when allocations are sent

Lines 450-463

Asked as: save without allocations; allocations not required draft; allocations optional saving

If the save request includes allocations, they are checked (whole hours, active divisions). If the save omits allocations, the existing ones stay unchanged. Allocations are not required to save a draft.

## Allocations are replaced when they change

Lines 458-501

Asked as: allocations deleted on save; old allocations replaced; allocations changed overwritten

If the allocations sent differ from the saved ones, Turbo deletes the old allocations and saves the new ones. If none are sent, the old ones stay.

## Save with no changes does nothing

Lines 464-473

Asked as: save no changes; saving unchanged request; version not bumped

Saving an existing request with no changes to the business fields or allocations does not change the request or bump its version.

## Request ID is returned and detail is shown after a save

Lines 502-508

Asked as: screen after save; save returns request detail; show request after saving

After a successful save, including a save that changed nothing, Turbo returns the request's detail.

## Admin profile and checklist helper scope

Lines 512-518

Asked as: which checklist shown; creation checklist awaiting; checklist depends on status

The checklist for a request depends on its status: a request Awaiting creation uses the creation checklist; any other status uses the approval checklist.

## Who can reject a request

Lines 523-538

Asked as: who can reject request; reject not allowed; approver reject request

The selected approver with the project_authorization claim can reject while the request is Awaiting approval or Awaiting creation. A job holder (job claim) can reject only when it is Awaiting creation. Others see 'Only the selected approver, or a job holder once the request is approved, can reject this request.'

## Checklist must be fully confirmed

Lines 542-567

Asked as: checklist changed reload; confirm checklist questions; no checklist saved

Approve and create each need the current checklist: a saved checklist ('No <scope> checklist is saved.'), a supported version, and the latest version, or Turbo says 'The checklist changed. Reload the request and review it again.' The answers must match the question count and all be true ('Confirm: <question>').

## Request actions need the request enabled, visibility and current version

Lines 568-592

Asked as: request action refused version; submit recall approve refused; action needs visibility

Submit, recall, reject, approve and create each require project requests to be enabled, that the person can view the request, and that the version sent matches the request's version. A stale version gives "This request changed. Reload it before continuing."

## Submit a request

Lines 594-618

Asked as: submit project request; request awaiting approval; submit rejected request

Submitting needs the time claim and an editable request. A Rejected request gives 'Edit the rejected request before submission.' The request must pass the completeness checks. Then it becomes Awaiting approval, records the submission time, gets a new submission token and the authorization PDF record resets to Pending.

## Recall a request

Lines 619-659

Asked as: withdraw submitted request; recall submitted request; only creator recall

Only the creator can recall ('Only the creator can recall this request.'), and only when it is Awaiting approval or Awaiting creation. Otherwise Turbo says 'This action requires a submitted or approved request.' The request becomes Recalled, the submission token clears and the authorization PDF approval clears.

## Reject a request

Lines 629-659

Asked as: reject reason required; enter rejection reason; rejected request marked

A rejection needs a non-blank reason ('Enter a rejection reason.'). The request becomes Rejected with the time, the rejector and the trimmed reason. The submission token clears and the authorization PDF record is marked Rejected.

## Approve a request

Lines 660-699

Asked as: approve project request; no submission token; reviewed data changed

Only a submitted request in Awaiting approval can be approved. It needs a submission token or Turbo says 'This request has no submission token. Recall and submit it again.' Only the selected approver can approve, and the project value must be within their limit. The checklist, the review token and the document hash must match what the approver reviewed ('The reviewed data changed. Reload and review it again.'). The request becomes Awaiting creation.

## Create a job from an approved request

Lines 700-749

Asked as: create job from request; approval stale recall; saved approval incomplete

Creating a job needs the job claim and a request in Awaiting creation. The submission token must exist, the creation checklist must be confirmed, and the approval must match the current version, value, token and document, or Turbo says 'The approval is stale. Recall the request for a new approval.' An incomplete saved approval gives 'The saved approval is incomplete. Recall the request for a new approval.' and an invalid approval checklist gives 'The saved checklist is invalid.' The approver's limit is not checked again at creation.

## Notifications after submit, approve and reject

Lines 760-766

Asked as: email after submit; notification after reject; withdrawn approval notice

Submit, approve and reject send a notification after saving. A reject by the approver when the request is Awaiting creation is treated as a withdrawn approval. A job holder's reject is not.

## Request action results and notices

Lines 767-777

Asked as: what create returns; notices sent after action; action result id

After create, Turbo returns the new job's id. After other actions it returns the updated request. Notices queued during the action are sent after the transaction succeeds.

## What job creation copies and deletes

Lines 779-846

Asked as: what copies to job; request converted to job; FN agreement job created

The new job copies the request's fields except the approver. Its status is the request's project status. The job is marked as an FN agreement when the answer is Yes, and its authorizing document is PA. Division allocations become job time allocations. The authorization PDF moves to the job as its first approval (sequence 1). The request and its allocations are deleted.

## Request notes become the job's first client note

Lines 848-869

Asked as: request notes client note; notes become job note; notes first job note

Non-blank request notes become a client note on the new job, authored by the request creator. Blank notes make no note.

## Deleting a project request

Lines 870-917

Asked as: delete project request; cannot delete request; only abandoned delete

The creator can delete a request that they can still edit (Draft, Recalled or Rejected). Anyone else needs the admin claim, and the request must be abandoned, or Turbo says 'Only abandoned requests can be deleted.' Deleting also removes its authorization PDF record and allocations.
