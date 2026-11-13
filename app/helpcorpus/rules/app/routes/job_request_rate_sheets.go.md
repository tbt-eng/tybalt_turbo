<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_request_rate_sheets.go","sha256":"611d91e6da8dbbcf2dc24f861d8a2971666e65ec191c32a202bd014902a6d225","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_request_rate_sheets.go

## Proposed rates reason has a 500-character limit

Lines 39-82

Asked as: 500 character limit reason; reason too long; rate reason max length; proposed rates reason limit

The reason for proposed rates can have at most 500 characters ("the reason can have at most 500 characters"). Turbo does not require a reason when it saves a draft request. The code says that submission requires one, but that check is not in this file.

## Chosen rate sheet in a project request must be active

Lines 51-65

Asked as: draft rate sheet not allowed; can't pick draft sheet; select active rate sheet; rate sheet invalid for project

When a project request names a rate sheet, Turbo refuses a draft rate sheet unless it is the request's own draft. The message is "Select an active rate sheet." If the rate sheet cannot be found, the message is "Select a valid rate sheet." An empty choice passes this check.

## Proposed rates need a name and a starting rate sheet

Lines 71-89

Asked as: proposed rates name missing; save rates without name; no starting rate sheet; rate sheet required to save

A proposed-rates save is refused when the name is empty ("name is required") or no starting sheet is chosen ("choose the rate sheet to start from"). Turbo trims spaces from the name, reason and starting sheet before it checks them.

## Proposed rates must start from an active rate sheet

Lines 91-99

Asked as: start from inactive sheet; based on draft sheet error; proposed rates base sheet; can't base rates on draft

The starting rate sheet (based_on) of proposed rates must exist and must not be a draft, or Turbo says "Start the proposed rates from an active rate sheet." An inactive sheet is accepted only when the request's own draft already started from that sheet.

## Proposed rates must use known roles

Lines 100-104

Asked as: unknown role in rates; role not recognised rate; rate role doesn't exist; invalid role proposed rate

Every entry in proposed rates must name an existing role. Otherwise Turbo refuses with "A proposed rate has an unknown role."

## Proposed rate sheet name must be unique in any letter case

Lines 106-121

Asked as: rate sheet name taken; duplicate rate sheet name; name already used rates; same name rate sheet exists

Turbo compares the proposed name with every rate sheet, including all revisions, ignoring letter case. The request's own draft is excluded from the comparison. If another sheet has the name, Turbo says "another rate sheet already has this name".

## Proposed rates that match an active rate sheet are refused

Lines 123-132

Asked as: rates already exist as sheet; duplicate of active rate sheet; use existing rate sheet instead; proposal identical to sheet

Turbo compares the proposed rates with the active rate sheets. If an active sheet has exactly the same rates, Turbo refuses the proposal and names that sheet, so the requester can choose it instead. The comparison itself is in hooks.CompareWithActiveRateSheets.

## Proposal is rewritten only when something changed

Lines 134-166

Asked as: draft not updating; changes to proposed rates not saved; nothing saved proposal; proposal unchanged no update

A new proposal always counts as changed. For an existing draft, Turbo rewrites it only when the name, reason or starting sheet differs, or when the rates differ. If nothing differs, the request points to the draft and nothing is rewritten.

## Proposed rates are saved as an inactive draft rate sheet

Lines 151-211

Asked as: where is my draft rate sheet; new rate sheet inactive; draft rate sheet effective date; rate sheet revision zero

A new draft rate sheet gets revision 0, is inactive, and belongs to the request. Its effective date is set to today (hooks.RateSheetToday). Job creation sets the date when the sheet becomes active. The request then points to the draft as its rate sheet.

## Saving proposed rates updates entries by role

Lines 176-208

Asked as: rate entries not updating; role rate lost on save; decimal rate cut off; rate change dropped role

Turbo updates a draft's entries by role. It adds an entry for a new role, updates one whose rate or overtime rate changed, and leaves unchanged ones alone. It deletes entries for roles that the proposal no longer lists. Turbo stores the rate as a whole number (int), so the rate is truncated.

## Replaced request draft is deleted

Lines 213-220

Asked as: old draft still there; leftover draft rate sheet; draft remains after switching sheet; delete unused draft rates

When a request no longer uses its draft, Turbo deletes the draft after the request is saved with its new rate sheet.

## Approver sees proposed rates compared with the source sheet

Lines 226-258

Asked as: which rates changed from source; proposed rates comparison view; rate marked changed approver; highlighted rates differ from original

When a request has a draft, the editor and the approver see it marked as proposed, with the sheet it started from. Each rate shows the source's rate and overtime rate. A rate is marked changed when either value differs. A role that the source sheet lacks is marked changed. A request that uses an active sheet shows none of this.

## Approval checklist wording for proposed rates

Lines 260-267

Asked as: checklist wording proposed rates; rate sheet checklist text; approval checklist rates statement; proposed rates approval statement

For proposed rates, the rate sheet statement of the approval checklist reads "The proposed rates are approved for this client." The checklist key stays rate_sheet_correct and the stored version stays the same, so approvals need no new checklist version.
