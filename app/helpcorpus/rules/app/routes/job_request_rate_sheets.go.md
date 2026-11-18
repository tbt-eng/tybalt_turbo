<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_request_rate_sheets.go","sha256":"e7d2f7599c5510137470dceae886d274a4515e28679f606bd042a86c7e0ad765","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_request_rate_sheets.go

## Proposed rates reason is limited to 500 characters

Lines 41-42

Asked as: reason too long; reason character limit; 500 characters reason; rate reason too long

A proposed rates reason (rate_sheets.reason) can have at most 500 characters. This is the length of the rate_sheets.reason field.

## Only the project request's own draft is chosen among drafts

Lines 44-51

Asked as: my draft rate sheet; which draft belongs to request; request draft sheet; only one draft per request

A project request owns at most one draft rate sheet from its proposed rates. Turbo uses the first draft that hooks.RequestRateSheetDrafts returns for the request. If the request has no draft, it has no own draft.

## Chosen rate sheet in a project request must be active

Lines 53-67

Asked as: select an active rate sheet; rate sheet not active; choose inactive sheet; draft sheet rejected

When a requester chooses a rate sheet in a project request, Turbo refuses a draft sheet unless it is the request's own draft. The message is "Select an active rate sheet." An unknown sheet gives "Select a valid rate sheet." An empty choice passes this check.

## A request may choose its own draft rate sheet

Lines 55-62

Asked as: choose my own draft; use request's draft sheet; pick draft as rate sheet; invalid rate sheet choice

When a project request chooses a rate sheet, an empty choice or the request's own draft passes the check without a lookup. Any other sheet must exist, or Turbo says "Select a valid rate sheet."

## Proposed rates need a reason of at most 500 characters

Lines 75-85

Asked as: reason is required; missing reason proposed rates; empty reason rates; save rates without reason

Every save of proposed rates in a project request needs a reason (reason). The reason is trimmed. An empty reason gives "reason is required". A reason over 500 characters gives "the reason can have at most 500 characters".

## Proposed rates need a rate sheet to start from

Lines 86-92

Asked as: choose rate sheet to start from; based on sheet missing; no starting rate sheet; which sheet to copy rates from

Every save of proposed rates needs the rate sheet to start from (based_on). If it is empty, Turbo says "choose the rate sheet to start from". Turbo also checks the rate entries with the rate sheet entry rules and returns all errors together.

## Proposed rates must start from an active sheet

Lines 96-102

Asked as: start from inactive sheet; based on sheet deactivated; start from draft sheet; starting sheet not active

The sheet to start from must exist and must not be a draft. It must be active. Otherwise Turbo says "Start the proposed rates from an active rate sheet." One exception: if the request's own draft already starts from that sheet, the sheet may since have been deactivated and still passes.

## Proposed rate with unknown role is refused

Lines 103-107

Asked as: unknown role on rate; role does not exist; invalid role in proposed rates; rate for missing role

If any proposed rate names a role that does not exist, Turbo refuses the save with "A proposed rate has an unknown role."

## Proposed rate sheet name must not be held by another sheet

Lines 109-121

Asked as: rate sheet name taken; duplicate rate sheet name; name already used by sheet; proposed rates name exists

The name of proposed rates is optional. If another rate sheet holds the name, in any letter case, Turbo refuses the save with an error on the name. The request's own draft is not counted as another sheet. The message comes from RateSheetNameHeld.

## Proposed rates that match an active sheet are refused

Lines 123-133

Asked as: rates same as existing sheet; duplicate rates refused; identical to active rate sheet; rates already exist as sheet

Turbo refuses proposed rates that exactly match the rates of an active rate sheet. It names that sheet (DuplicateRateSheetError), so the requester can choose it instead.

## Draft changes only when name, reason, source or rates change

Lines 135-167

Asked as: save draft without changes; resave proposed rates; nothing changed in draft; draft not updated

When a request already has a draft, saving proposed rates counts as a change only if the name, reason, based_on sheet or the rates differ from what is stored. A request with no draft always creates one. If nothing changed, Turbo only points the request at the existing draft and leaves its approval alone.

## Saving changed proposed rates creates an inactive draft and clears approval

Lines 153-177

Asked as: draft lost approval after save; new draft needs approval; approval cleared after edit; draft rate sheet inactive

A new draft rate sheet starts at revision 0, inactive, and owned by the request. Saving changed rates clears the rates approver and approval time, so the rates need a new approval. Turbo sets the effective date to today, and job creation sets the real date.

## Draft entries are synced to the proposed rates

Lines 179-214

Asked as: draft rates not updating; role missing from draft; old roles still on draft; draft entries out of sync

When proposed rates are saved, Turbo adds an entry for each new role. It updates an entry only if its rate or overtime rate differs. It deletes entries for roles that are no longer proposed.

## Renaming proposed rates at job creation

Lines 225-269

Asked as: rename proposed rates at job; no proposed rates to rename; change rate sheet name when creating job; rename draft sheet

When procurement creates the job, they may rename the proposed rates. An empty name or the same name does nothing. If the request has no draft, Turbo says "This request has no proposed rates to rename." If another sheet holds the name, Turbo refuses with that message. A rename keeps the approval and records the old name, the new name and who renamed it.

## Request shows proposed rates compared with the source sheet

Lines 277-322

Asked as: compare proposed rates to source; rate changed flag; missing roles on request; request rates versus source

For a request with proposed rates, Turbo adds to the request's related records the flag rate_sheet_proposed, the source sheet and the rates approver. It also adds rate_sheet_missing_roles, the roles without a rate (empty when complete). Each rate shows the source's rate and overtime rate, and a changed flag. A role not on the source counts as changed. A request with an active sheet gets none of this.

## Approving proposed rates needs an explicit statement

Lines 338-349

Asked as: approve proposed rates; approval statement missing; approve new rate sheet; approver cannot approve rates

When the approver approves a request with proposed rates, they must make the statement of approval. Otherwise Turbo says "Approve the creation of the new rate sheet and its rates." If the statement is made, Turbo records the approver and the time on the draft. A request with an active sheet needs nothing.

## New submission clears the proposed rates approval

Lines 353-363

Asked as: resubmit request approval cleared; submitting clears rates approval; approval gone after resubmit; need new rates approval

Turbo clears the rates approver and approval time on the request's draft, so a new submission needs a new approval. It does nothing if there is no approval to clear.

## Job creation needs the request approver's approval of proposed rates

Lines 367-376

Asked as: cannot create job unapproved rates; job creation blocked rates; proposed rates not approved; recall request for approval

Turbo refuses to create a job from a request with proposed rates unless the draft has an approval time and a rates approver who is the request's approver. The message is "The proposed rates are not approved. Recall the request for a new approval."

## Compare proposed rates with active sheets

Lines 378-414

Asked as: check rates against active sheets; preview duplicate rates; compare unsaved rates; see similar rate sheets

Anyone with the time claim can compare unsaved proposed rates with the active rate sheets, if project requests are enabled. Turbo returns the exact matches, which a save would refuse. It also returns a few close sheets, leaving out the sheet in based_on. Nothing is saved.
