<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_request_rate_sheets.go","sha256":"ed20a055f3a3be6a8afbeda4f13ddf52235b7b5b025da12c51264c053e4b1d09","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_request_rate_sheets.go

## A draft rate sheet belongs to one project request

Lines 16-51

Asked as: draft rate sheet owner; proposed rates belong to request; who can edit draft rates; rate sheet draft for request

Proposed rates are kept as a draft rate sheet that the request owns. Turbo takes the request's first draft as its own. Only the request's creator changes the draft, and only by saving the request under the same edit rules as the rest of the request.

## A draft rate sheet is not listed or used elsewhere

Lines 16-22

Asked as: draft rates showing in lists; use draft sheet elsewhere; draft rate sheet not visible; draft sheet in other requests

A draft rate sheet is never listed or used elsewhere. Job creation activates it.

## Proposed rates reason is limited to 500 characters

Lines 41-42

Asked as: reason too long; proposed rates reason limit; reason character limit; reason over 500 characters

The reason for proposed rates in a project request can have at most 500 characters (proposedRatesReasonMax, the length of rate_sheets.reason). A longer reason is refused with 'the reason can have at most 500 characters'.

## Chosen rate sheet in a project request must be active

Lines 53-67

Asked as: cannot pick rate sheet; rate sheet not accepted on request; select active rate sheet error; draft sheet rejected in request

When a project request names a rate sheet, Turbo refuses a draft rate sheet unless it is the request's own draft. The message is "Select an active rate sheet." If the sheet cannot be found, the message is "Select a valid rate sheet."

## Proposed rates need a reason and a starting sheet

Lines 75-92

Asked as: reason missing on rates; starting sheet required; proposed rates won't save; must pick sheet to copy

Every save of proposed rates needs a reason (the error is "reason is required") and a rate sheet to start from ("choose the rate sheet to start from"). The reason has at most 500 characters. Turbo trims the name, the reason and the starting sheet before it checks them.

## Proposed rates must start from an active rate sheet

Lines 94-102

Asked as: starting sheet inactive error; copy from draft or inactive sheet; base sheet not active; proposed rates source sheet refused

Turbo refuses proposed rates that start from a missing sheet or from a draft sheet. The message is "Start the proposed rates from an active rate sheet." An inactive sheet is allowed only when it is the sheet the request's own draft already started from.

## Proposed rate with unknown role is refused

Lines 103-107

Asked as: role not recognised in rates; rate for missing role; invalid role on proposal; role does not exist error

If a proposed rate uses a role that does not exist, Turbo refuses the save with "A proposed rate has an unknown role."

## Proposed rates name must not be held by another sheet

Lines 109-121

Asked as: rate sheet name already taken; duplicate proposed rates name; name used by another sheet; sheet name exists already

Turbo refuses a proposed rates name that another rate sheet already holds, in any letter case. The request's own draft is not counted as another sheet. The name is optional.

## Proposed rates that match an active sheet are refused

Lines 123-133

Asked as: rates same as existing sheet; duplicate rate sheet refused; identical rates already exist; pick existing sheet instead

Turbo compares the proposed rates with the active rate sheets. If an active sheet has exactly the same rates, Turbo refuses the save and names that sheet, so the requester can choose it instead.

## Saving proposed rates only changes the draft when something differs

Lines 135-145

Asked as: save rates with no changes; nothing changed on proposal; draft not updated; proposal unchanged save

Turbo counts a proposal as changed when the request has no draft yet, or when the name, the reason, the starting sheet or the rates differ from the draft. If nothing changed, the request points at the existing draft and Turbo writes nothing.

## Proposed rates are saved as an inactive draft rate sheet

Lines 152-214

Asked as: where do new rates go; rates saved as draft; draft rates not active; new rate sheet created on save

A new draft is a rate sheet with revision 0 that is not active and belongs to the request. Its effective date is set to today, and job creation sets the real date later. Turbo saves only the rate entries that are new or whose rate or overtime rate changed. It deletes entries for roles that the proposal no longer lists. The request then points at the draft.

## Changed proposed rates lose their approval

Lines 171-172

Asked as: approval lost after edit; rates approval reset; need approval again after change; edited rates unapproved

When a save changes the proposed rates, Turbo clears the rates approver and the approval time on the draft. The request needs a new approval.

## Renaming proposed rates at job creation

Lines 225-269

Asked as: change proposed rates name; rename rate sheet at job; job creation rename sheet; rename draft before job

When a job is created, an empty name or an unchanged name does nothing. Turbo refuses a name that another sheet holds. Otherwise it renames the draft and records the old name, the new name and who renamed it in the rate sheet name changes. The approval stands. If the request has no proposed rates, the message is "This request has no proposed rates to rename."

## Request details show the proposal beside its source rates

Lines 280-330

Asked as: request shows source rates; rates compared to source; why rates marked changed; request rate details beside source

For a request with proposed rates, Turbo shows the draft's reason, the starting sheet and the rates approver. It lists the roles that have no rate and gives the source sheet's rates for them when they exist. Each rate shows the source rate and source overtime rate and is marked changed when either differs. A role the source sheet lacks is marked changed. A request that uses an active sheet shows none of this.

## Missing rates are listed by role name

Lines 336-361

Asked as: rates missing for roles; which roles have no rate; roles without rates list; missing role rate rows

Turbo lists each rate role that the draft has no rate for, ordered by role name. A row gives the role and, when the source sheet has a rate for it, the source rate and overtime rate.

## Approving proposed rates is a separate statement

Lines 373-388

Asked as: approve new rates statement; rates approval wording missing; approve rate sheet creation; approve proposed rates text

The approver must make a statement that approves the creation of the new rate sheet and its rates. Without it Turbo refuses with "Approve the creation of the new rate sheet and its rates." Turbo records the approver and the time on the draft. A request that uses an active sheet skips this.

## A new submission clears the proposed rates approval

Lines 390-402

Asked as: resubmit clears approval; approval gone after submit; submitting request removes rates approval; reset rates approval on resubmission

Turbo clears the draft's approver and approval time on a new submission, so the request needs a new approval. It does nothing when the draft has no approval.

## Job creation needs approved proposed rates

Lines 404-415

Asked as: cannot create job from request; job creation blocked; rates not approved for job; approver mismatch job creation

Turbo refuses to create a job from a request whose draft has no approval, whose approval was not made by the request's approver, or whose approval time is empty. The message is "The proposed rates are not approved. Recall the request for a new approval."

## Comparing unsaved proposed rates with active sheets

Lines 417-453

Asked as: check rates before saving; preview matching rate sheets; compare rates to existing sheets; see close sheets for rates

Anyone who holds the time claim can compare unsaved proposed rates, if project requests are enabled. Turbo returns the exact matches that a save would refuse and a few close sheets. It leaves the starting sheet out of the close sheets. It changes nothing.
