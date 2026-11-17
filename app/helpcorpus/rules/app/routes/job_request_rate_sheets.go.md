<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_request_rate_sheets.go","sha256":"21a6f6758e019bd3ee740664f3065cfbc835487055d3a42d68c421a483f145c4","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_request_rate_sheets.go

## A project request owns at most one proposed-rates draft

Lines 16-51

Asked as: two drafts on one request; which draft rate sheet; request has no draft; draft rates missing

Proposed rates are saved as a draft rate sheet that the project request owns. Turbo looks up the draft by request and uses the first one it finds. When the request has none, it treats the request as having no draft.

## Proposed rates reason is limited to 500 characters

Lines 41-42

Asked as: reason too long; rates reason character limit; 500 characters reason; reason text too big

The reason (rate_sheets.reason) for proposed rates in a project request can have at most 500 characters (proposedRatesReasonMax). The 500 limit matches the length of the reason field on a rate sheet.

## Chosen rate sheet in a project request must be active

Lines 53-67

Asked as: cannot choose rate sheet; draft sheet not allowed on request; select active rate sheet error; wrong rate sheet chosen

When a project request chooses a rate sheet, Turbo refuses a draft sheet unless it is the request's own draft. The message is "Select an active rate sheet." A sheet that does not exist gives "Select a valid rate sheet." An empty choice passes this check.

## Proposed rates need a reason and a starting sheet

Lines 74-92

Asked as: reason required for rates; must pick starting rate sheet; proposed rates won't save; blank reason rejected

Every save of proposed rates in a project request needs a reason (proposed_rates.reason) and a rate sheet to start from (based_on). The messages are "reason is required" and "choose the rate sheet to start from". The reason can have at most 500 characters. Turbo trims spaces from the name, reason and based_on first.

## Proposed rates must start from an active rate sheet

Lines 94-102

Asked as: based on inactive sheet; start rates from old sheet; starting sheet is archived; base sheet not active

Turbo refuses proposed rates whose starting sheet is missing, is a draft, or is inactive. The message is "Start the proposed rates from an active rate sheet." An inactive starting sheet is allowed when the request's own draft already started from it.

## Proposed rate with unknown role is refused

Lines 103-107

Asked as: unknown role on rates; role not recognised in proposal; invalid role in rates; role does not exist rates

Turbo refuses proposed rates when any entry names a role that does not exist. The message is "A proposed rate has an unknown role."

## Proposed rate sheet name must be unique

Lines 109-121

Asked as: rate sheet name taken; duplicate rate sheet name; name already used for rates; sheet name exists already

Turbo refuses a proposed rate sheet name that another sheet already holds, in any letter case. The request's own draft is not counted as another sheet. The name itself is optional.

## Proposed rates that match an active sheet are refused

Lines 123-133

Asked as: rates same as existing sheet; duplicate of active sheet rates; proposal identical to sheet; rates already exist elsewhere

Turbo refuses proposed rates that exactly match the rates of an active rate sheet. The error names that sheet so the requester can choose it instead. The comparison skips the request's own draft code.

## Draft changes only when name, reason, source or rates change

Lines 135-167

Asked as: save did nothing to draft; draft not updating; no change saved proposal; edits not being kept

A new proposal always counts as changed. For an existing draft, Turbo treats the save as changed only if the name, reason, based_on sheet or the rates differ from the stored ones. An unchanged save leaves the draft alone and only points the request at it.

## Saving changed proposed rates creates or updates a draft sheet

Lines 152-214

Asked as: proposed rates not saved; where is my draft sheet; save proposal as draft; rates draft created

Turbo saves proposed rates as an inactive draft rate sheet (revision 0) that the request owns. It sets the effective date to today, which job creation later replaces. Entries are added, updated when the rate or overtime rate differs, or deleted when the role is no longer in the proposal. The request then points to the draft.

## Changed proposed rates clear the rates approval

Lines 171-172

Asked as: approval disappeared after edit; edited rates lost approval; rates need reapproval; approver removed after change

When proposed rates change, Turbo clears the draft's rates approver and approval time. The rates then need a new approval.

## Request's draft is deleted when no longer used

Lines 216-223

Asked as: old draft still there; unused draft remains; draft left after changing sheet; remove leftover draft

After the request is saved with a different rate sheet, Turbo deletes the draft that the request no longer uses. It does nothing when there is no draft.

## Renaming proposed rates when the job is created

Lines 225-269

Asked as: rename rates when creating job; change proposed sheet name; rates name change at job; no proposed rates to rename

The job holder creating the job can rename the request's proposed rate sheet. An empty or unchanged name does nothing. A request with no proposed rates gives "This request has no proposed rates to rename." A name held by another sheet is refused with that sheet's message. Turbo keeps the old name in the name-change history. The earlier approval of the rates stands.

## Request shows proposed rates beside the source sheet

Lines 271-312

Asked as: compare rates to source sheet; why rates show changed; rates beside original sheet; approver shown on request

For a request with proposed rates, Turbo shows the reason, the source sheet and the rates approver. Each rate is shown beside the source's rate and overtime rate. A rate is marked changed when the rate or overtime rate differs, or when the source has no rate for that role. A request that uses an active sheet shows none of this.

## Approver must state approval of the new rate sheet

Lines 324-339

Asked as: cannot approve request rates; approval refused new sheet; confirm new rate sheet approval; approve rates checkbox

When the approver approves a request with proposed rates, they must confirm the creation of the new rate sheet and its rates. Otherwise Turbo refuses with "Approve the creation of the new rate sheet and its rates." Turbo then records the approver and the time on the draft. A request that uses an active sheet skips this.

## New submission clears proposed rates approval

Lines 341-353

Asked as: resubmit cleared approval; approval gone after resubmit; submit again rates approval; reset rates approver

Turbo clears the rates approver and approval time on the request's draft, so a new submission needs a new approval. It does nothing when no approval is stored.

## Job creation needs approved proposed rates

Lines 355-366

Asked as: cannot create job from request; job creation blocked rates; rates not approved for job; recall request approval

Turbo refuses to create a job from a request with proposed rates unless the draft is approved by the request's own approver and has an approval time. The message is "The proposed rates are not approved. Recall the request for a new approval."

## Compare proposed rates with active rate sheets

Lines 368-404

Asked as: check rates before saving; find matching rate sheets; preview rates against sheets; similar rate sheets lookup

Anyone with the time claim can compare unsaved proposed rates, if project requests are enabled. Turbo returns the active sheets that match exactly, and a few close ones. The sheet in based_on is left out of the close ones. Nothing is saved.
