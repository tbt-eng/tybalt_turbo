<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_requests.go","sha256":"720ad8268a35e865e5176ffc01138eda4b51b31fffd37e1b6e6e180b6d0cc2f1","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_requests.go

## Action routes available on a project request

Lines 62-65

Asked as: submit button missing; can't send my project request; upload pdf to request; delete my request

Turbo offers submit, recall, reject, approve and create actions on a project request, plus delete and document upload. They all require sign-in.

## Active staff are required for request actions

Lines 96-121

Asked as: not an active staff error; can't act on request; missing claim message; why am I blocked

A person who is not an active staff member gets "An active staff member is required." Turbo also needs a signed-in user. A claim check gives "The <claim> claim is required." when the person lacks the claim.

## Project requests can be switched off

Lines 122-131

Asked as: cannot save project request; project requests turned off; can't submit anything right now; requests greyed out

When project requests are disabled (IsJobRequestApprovalRequired is false), Turbo refuses to save, submit, recall, reject, approve or create a request. The message is "Project requests are disabled."

## Request is abandoned after a set number of days

Lines 132-134

Asked as: request went stale; old draft expired; why was my request abandoned; request too old to use

A request is abandoned when its last_edited date is not later than JOB_REQUEST_ABANDONMENT_DAYS days ago. A request edited exactly that many days ago counts as abandoned.

## Who can view a project request

Lines 135-166

Asked as: can't open someone's request; cannot see request; who can view request; request hidden from me

The creator can view their own request. A job claim holder (procurement) can view any request. A project authorization claim holder can view a request only when selected as its approver. An admin can view a request only when it is abandoned. Anyone else gets "You cannot view this request."

## Only the creator edits, and only in some states

Lines 167-182

Asked as: can't edit my request; edit button greyed out; someone else's request locked; recall before editing

Only the creator can edit a request. Others get "Only the creator can edit this request." A request can be edited when it is Draft, Recalled or Rejected. In any other state the message is "Recall this request before editing it."

## Stale version is refused

Lines 194-201

Asked as: request changed under me; someone else edited it; version mismatch error; reload request message

Save and workflow actions compare the version the person saw with the request's current version. If they differ, Turbo says "This request changed. Reload it before continuing."

## Editing a request resets its approval state

Lines 202-239

Asked as: editing resets approval; changing request clears approval; rejected request back to draft; pdf approval went pending

Saving a changed request raises its version by 1 and sets last edited to now. It clears the submission token and any rejection. A Rejected request goes back to Draft. The attached authorization PDF's approval is cleared to Pending and takes the new project value.

## Project value must be a whole non-negative CAD amount

Lines 240-255

Asked as: project value not accepted; value has cents; negative project amount; dollar amount error

When a person saves a project request, the project value must be a whole number of at least 0 (at most 9007199254740991). Otherwise Turbo says 'Project value must be a whole CAD dollar amount.'

## Request fields that cannot be edited are refused

Lines 244-251

Asked as: field won't save; can't change that field; field not editable message; extra field refused

When a person saves a project request, Turbo accepts only the listed business fields (requestBusinessFields). Any other field is refused with 'Field is not editable: <field>'.

## Request notes must be 10 to 1000 characters or blank

Lines 256-267

Asked as: notes too short; notes too long; notes length error; notes won't save

Turbo trims the notes on a project request. Notes must be text, and after trimming they must have 10 to 1000 characters or be blank. Otherwise Turbo says 'Notes must be 10 to 1000 characters, or blank.' or 'Notes must be text.'

## Fee basis must be chosen as true or false

Lines 268-273

Asked as: time and materials missing; fee basis error; pick billing basis; t&m not saving

When a request sets time and materials (time_and_materials), the value must be true or false, or Turbo says 'Select a fee basis.' A valid value also marks the fee basis as chosen (time_and_materials_set).

## Project status must be Active or Closed

Lines 282-284

Asked as: project status error; status must be active; can't pick project status; closed or active choice

Saving a request requires a project status (project_status) of Active or Closed. Otherwise Turbo says 'Select Active or Closed for the project status.'

## Selected approver's project limit is checked when saving

Lines 288-292

Asked as: approver limit too low; approver cannot cover value; authority limit error; approver's spending limit

When a request has an approver, saving checks the approver's project authorization limit against the project value (ProjectAuthorizationLimit). The check is skipped when the caller turns authority checks off, as at job creation.

## Parent project must be a valid top-level project with same client

Lines 293-304

Asked as: parent project invalid; subjob client mismatch; top level project only; wrong parent project number

If a request names a parent, the parent must exist ('Select a valid parent project.'). Its number must look like 25-123 or 25-1234 ('Only a top-level project can be a parent.'). The request's client must match the parent's client ('The subjob must use the parent project's client.').

## Linked proposal must fit the client and status

Lines 305-318

Asked as: proposal not accepted; proposal wrong client; proposal number format; proposal not awarded status

A linked proposal must exist ('Select a valid proposal.'). Its number must start with P and it must belong to the request's client ('Select a proposal for this client.'). Its status must be In Progress, Submitted or Awarded ('Select an In Progress, Submitted, or Awarded proposal.').

## Required fields before a request is submitted

Lines 325-335

Asked as: what's needed before submit; required field blocking submit; project value zero error; pick fee basis first

Submitting or approving needs description, location, client, contact, manager, branch, rate sheet, project award date, invoicing information, approver and fn agreement answer. A blank one gives 'Required field: <field>'. A fee basis must be chosen ('Select a fee basis.') and the project value must be above 0 ('Project value must be positive.').

## Linked proposal must be Awarded to submit or approve

Lines 336-344

Asked as: proposal not awarded yet; linked proposal blocks submit; set proposal to awarded; can't submit with proposal

When a request links a proposal, the proposal must have status Awarded before submit, approve or create. Otherwise Turbo says 'Set the linked proposal to Awarded before continuing.'

## Request save validates client, contact and invoicing references

Lines 345-347

Asked as: client reference error; contact invoicing check failed; client details invalid on save; reference validation fails

Turbo checks the client references of a request (ValidateClientReferences) on every save, and again on submit, approve and create. The rules for those references are in another file.

## Contact emails and project data are validated before submit

Lines 348-364

Asked as: contact email rejected; project data check failed; can't submit bad emails; validation fails before submit

Before a request is submitted or approved, Turbo validates the contact emails (ValidateRequestContactEmails) and the project data (hooks.ValidateProjectRequestData). The details of both checks are in other files.

## Request allocations are required and checked when submitting

Lines 351-361

Asked as: allocations missing on submit; division hours required; no time split entered; allocation check failed

When a request is validated as complete, Turbo loads its saved division allocations and requires them, with whole hours and active divisions. Saving a draft with an empty allocation list skips this check.

## Authorization PDF is required to submit

Lines 365-374

Asked as: no authorization document; need signed pdf; can't submit without pdf; upload authorization file

A request needs an uploaded authorization PDF with a document hash before it can be submitted, approved or created. Otherwise Turbo says 'Upload an authorization PDF.'

## Allocations must be whole hours in active divisions

Lines 377-401

Asked as: allocation hours decimal; invalid division chosen; inactive division error; hours must be whole

Each division allocation on a request must have whole-number hours, a valid division ('Select a valid division.') and an active division ('Select an active division.'). A complete request needs allocations. A save may omit them.

## Creating a request needs the time claim

Lines 415-429

Asked as: can't create request; no time claim access; who can start request; missing time permission

A person needs the time claim to create a request. A new request starts as Draft, with project status Active and version 1.

## New requests start as Active drafts

Lines 420-429

Asked as: new request defaults; what status is new; how a request starts; owner of new request

A new project request starts with status Draft, project status Active, version 1, and the creator as its owner. The last edited time is set to now.

## Request saves reject changes made on an outdated version

Lines 430-441

Asked as: save rejected outdated copy; someone saved before me; old version refused; reload before saving

When a creator saves a project request, the version sent must match the request's current version. Turbo refuses a mismatch with 'This request changed. Reload it before continuing.' Submit, recall, reject, approve and create actions use the same check.

## Choose a rate sheet or propose rates, not both

Lines 453-479

Asked as: rate sheet and new rates; both rate options chosen; proposed rates conflict; choosing sheet deletes draft

A request cannot both choose a rate sheet and propose new rates. Turbo says 'Choose a rate sheet or propose new rates, not both.' Proposed rates set the request's rate sheet to its own draft. Choosing another sheet drops that draft.

## Saving a draft replaces its allocations only when they differ

Lines 480-493

Asked as: allocations overwritten on save; draft allocations replaced; saved allocations changed; empty allocations on draft

When a draft is saved, Turbo compares the sent allocations with the saved ones. It deletes the old allocations and saves the new ones only if they differ. Omitting allocations leaves the saved ones unchanged.

## Saving with no changes does nothing

Lines 494-503

Asked as: save did nothing; no change saved; version not bumped; nothing changed when saving

When a person saves an existing request without changing any business field, allocation or proposed rate, Turbo makes no change. The version and last edited time stay the same.

## Saving a new request also saves proposed rates

Lines 504-525

Asked as: proposed rates not saved; new request rates lost; rate draft on create; rate sheet deletes draft

When a new request has proposed rates, Turbo saves the request first, then saves the rate draft, then saves the request again. A rate sheet chosen in place of an existing draft deletes that draft.

## Request status decides which checklist applies

Lines 556-562

Asked as: which checklist shows; creation checklist vs approval; why different checklist; checklist depends on status

A request awaiting creation uses the creation checklist. Any other request uses the approval checklist.

## Who can reject a request, and the reason

Lines 564-582

Asked as: who can reject; reject button missing; reason needed to reject; job holder reject rights

The selected approver (project_authorization claim) can reject while a request is Awaiting approval or Awaiting creation. A job holder (job claim) can reject only once it is Awaiting creation. Others get 'Only the selected approver, or a job holder once the request is approved, can reject this request.' A rejection reason is required ('Enter a rejection reason.'). It sets status Rejected.

## Approval needs a current reviewed checklist and data

Lines 586-611

Asked as: checklist out of date; answer all questions; confirm checklist error; reviewed data mismatch

Approving or creating requires a saved checklist whose version is current. A page with an older checklist is refused ('The checklist changed. Reload the request and review it again.'). Every question must be answered true, with no extra answers ('Confirm: <key>'). The review token and document hash must match, or Turbo says 'The reviewed data changed. Reload and review it again.'

## Every request action requires project requests to be enabled

Lines 621-623

Asked as: requests disabled error; can't approve right now; all request actions blocked; project requests switched off

Submit, recall, reject, approve and create each check that project requests are enabled. If not, Turbo refuses with 'Project requests are disabled.'

## Submitting a request

Lines 638-662

Asked as: how to submit request; submit makes awaiting approval; submission token issued; sending request for approval

Submit needs the time claim, the creator, and a Draft or Recalled request. A Rejected request gives 'Edit the rejected request before submission.' A valid submit sets status Awaiting approval, records the submitted time, issues a new submission token and clears earlier approval on the PDF.

## Submit requires the time claim and an editable request

Lines 638-650

Asked as: who can submit; rejected request submit error; submit needs time claim; submit fails validation

Only a user with the time claim who created the request can submit it. The request must be Draft or Recalled; a Rejected request gets 'Edit the rejected request before submission.' Turbo then runs full validation and checks the approver's project limit.

## Recalling a request

Lines 663-703

Asked as: take back request; withdraw submitted request; only creator recalls; recall button failing

Only the creator can recall a request ('Only the creator can recall this request.'). The request must be Awaiting approval or Awaiting creation, or Turbo says 'This action requires a submitted or approved request.' Recall sets status Recalled and clears the submission token and the PDF's approval.

## Recall is limited to the creator while submitted or approved

Lines 663-672

Asked as: recall not allowed; can't recall approved request; recall wrong status; pull back my request

Only a request in Awaiting approval or Awaiting creation can be recalled: otherwise 'This action requires a submitted or approved request.' Only the creator can recall: 'Only the creator can recall this request.' Recall sets the status to Recalled.

## Rejecting records who, when and why

Lines 673-703

Asked as: who rejected my request; rejection details recorded; reject withdraws approval; why was it rejected

Rejecting a request sets status Rejected and stores the rejector, time and trimmed reason. It clears the submission token. The PDF is also marked Rejected with the same details. An approver who rejects an Awaiting creation request withdraws their approval.

## Reject needs permission and a reason

Lines 673-691

Asked as: can't reject request; reject reason blank; not the selected approver; reject permission denied

The reject action needs an active user who is the selected approver, or a job holder once approved. Others get 'Only the selected approver, or a job holder once the request is approved, can reject this request.' A blank reason gives 'Enter a rejection reason.'

## Recall or reject cancels the submission and approval

Lines 692-703

Asked as: recall clears approval; submission token cleared; pdf marked rejected; resubmit after reject

Recall and reject clear the submission token and the saved approval on the authorization document. A reject marks that document Rejected with the reason. Resubmission issues a new token.

## Approving a request

Lines 704-720

Asked as: approve button not working; only approver can approve; no submission token error; approver limit blocks approval

Only a request that is Awaiting approval and not rejected can be approved ('Only a submitted request can be approved.'). It needs a submission token, or Turbo says 'This request has no submission token. Recall and submit it again.' Only the selected approver can approve ('Only the selected approver can approve this request.'). The approver's project limit must cover the project value.

## Approve needs a submission token and the selected approver

Lines 705-720

Asked as: cannot approve request; approve says no token; not the approver; approval status wrong

Only a request in Awaiting approval with no rejection can be approved: otherwise 'Only a submitted request can be approved.' A request with no submission token gets 'This request has no submission token. Recall and submit it again.' Only the selected approver can approve: 'Only the selected approver can approve this request.' The approver's project limit must cover the project value.

## Approval is refused if reviewed data changed

Lines 721-742

Asked as: approval refused changed data; reviewed data mismatch on approve; hash mismatch approval; approver details recorded

When approving, the review token and document hash sent must match the current request and PDF. Otherwise Turbo refuses with 'The reviewed data changed. Reload and review it again.' Turbo records the approver, time, project limit, request version, checklist version and token on the authorization document.

## Approval sets status Awaiting creation

Lines 735-743

Asked as: what happens after approval; status awaiting creation; pdf marked approved; approval recorded

A valid approval marks the PDF Approved with the approver, time, project value, approver's limit and checklist version. The request status becomes Awaiting creation.

## Creating the job from an approved request

Lines 744-780

Asked as: create job button fails; approval is stale; saved approval incomplete; checklist invalid on create

Create needs the job claim and a request in Awaiting creation. It needs the submission token, the current creation checklist and a matching approval. Otherwise Turbo says 'The approval is stale. Recall the request for a new approval.' or 'The saved approval is incomplete...'. The saved approval checklist must still exist ('The saved checklist is invalid.').

## Create requires the job claim and an approved request

Lines 745-756

Asked as: who can create job; job claim required; not awaiting creation; create job no token

Only a user with the job claim can create the job, and only when the request is Awaiting creation: 'This request is not awaiting creation.' A request with no submission token gets 'This request has no submission token. Recall and submit it again.'

## Create refuses stale or incomplete approvals

Lines 760-780

Asked as: stale approval on create; incomplete approval message; unknown checklist saved; create job approval error

Creation checks the authorization document is Approved, matches the request version, project value and review token, and has the same approver, an approval time, a limit covering the value and no rejection. A stale approval gives 'The approval is stale. Recall the request for a new approval.' An incomplete one gives 'The saved approval is incomplete. Recall the request for a new approval.' An unknown saved approval checklist gives 'The saved checklist is invalid.' Creation does not re-check the approver's limit.

## Created job is recorded and requester notified

Lines 781-794

Asked as: job created notification; requester told about job; new job id returned; creation recorded

After creation Turbo records the creation checklist version on the authorization document, creates the job, and queues a 'created' notice. Turbo returns the new job's id.

## Notices queue on submit, approve and reject

Lines 804-810

Asked as: no notice sent; who gets emails; notification on reject; withdrawn approval email

After submit, approve or reject, Turbo queues a notice. A reject by the selected approver after approval is flagged as a withdrawn approval; a job holder's reject is not.

## What the new job takes from the request

Lines 823-849

Asked as: what carries to job; job fields from request; fn agreement on job; job authorizing document

Create copies the request's business fields except the approver into the new job. The job status is the request's project status. The job's fn agreement is true when the answer is Yes. Its authorizing document is PA. Turbo records the request creator, created and submitted times, who created it and when.

## Approved rates are activated before job creation

Lines 823-828

Asked as: rates not on new job; rate sheet activation; proposed rates go live; job rates missing

When the job is created, Turbo activates the request's proposed rate sheet in the same transaction (hooks.ActivateRequestRateSheet), so the job can use it.

## New job copies request fields except approver

Lines 833-841

Asked as: job fields copied; which fields move to job; job status from request; approver not copied

The new job copies the request's business fields that exist on jobs, except the approver. Its status is the request's project status, its First Nations agreement flag is true only when the answer is 'Yes', and its authorizing document is 'PA'.

## Request data moves to the job and the request is deleted

Lines 850-895

Asked as: request gone after job; what happens to request; allocations moved to job; request deleted on create

Creating a job moves the request's allocations to job time allocations and attaches the authorization PDF to the job as its first approval (sequence 1). Non-blank notes become the job's first client note, authored by the request creator. Then Turbo deletes the request.

## Request allocations become the job's time allocations

Lines 854-880

Asked as: hours moved to job; allocations on new job; pdf attached to job; first approval on job

At job creation, each request allocation becomes a job time allocation of the same division and hours, and the request allocation is deleted. The authorization document moves to the job as its initial approval with approval sequence 1.

## Request notes become the first client note

Lines 897-918

Asked as: notes missing on job; first client note; notes not carried over; job note from request

When the job is created, non-blank request notes become the job's first client note, authored by the request creator. Blank notes make no note.

## Deleting a project request

Lines 919-966

Asked as: delete my request; remove request; can't delete request; only admin delete

A creator can delete a request they can still edit. Other requests can be deleted only by an admin (admin claim), and only if abandoned ('Only abandoned requests can be deleted.'). Deleting also removes the request's PDF and allocations.

## Only the creator may delete an editable request

Lines 929-960

Asked as: delete button missing; delete abandoned request; admin delete old request; remove editable request

A creator can delete a request they can still edit. For any other request, the user needs the admin claim and the request must be abandoned: 'Only abandoned requests can be deleted.' Deleting also removes its authorization document and allocations.
