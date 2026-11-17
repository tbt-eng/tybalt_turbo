<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_request_rate_sheets.go","sha256":"725c2426c2f36198ed8e8bfa603100e9d72508f06b309f7c9cd70cff9f2e9e1a","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_request_rate_sheets.go

## Proposed rates are a private draft until job creation

Lines 16-21

Asked as: draft rates private; who can see proposed rates; edit proposed rates draft; rate sheet hidden until job

A project request that proposes new rates saves them as a draft rate sheet that the request owns. Only the request's creator changes the draft, and only by saving the request, under the same edit rules as the request. Turbo does not list or use the draft elsewhere until job creation activates it.

## Chosen rate sheet must not be another draft

Lines 53-67

Asked as: choose rate sheet draft; rate sheet not active; select a valid rate sheet; pick rate sheet for request

A project request that chooses a rate sheet must pick an existing one. If it does not exist Turbo says 'Select a valid rate sheet.' If it is a draft, Turbo says 'Select an active rate sheet.' The request's own draft is allowed.

## Proposed rates need a reason of at most 500 characters

Lines 75-85

Asked as: reason required proposed rates; reason too long; 500 characters reason; why proposing new rates

Every save of a project request with proposed rates needs a reason (proposed_rates.reason). Turbo trims it and refuses an empty one with 'reason is required'. A reason over 500 characters gets 'the reason can have at most 500 characters'.

## Proposed rates must name the rate sheet to start from

Lines 86-92

Asked as: choose rate sheet to start from; start from which rate sheet; based on sheet missing; no starting rate sheet

A project request save with proposed rates needs the rate sheet to start from (based_on). Without it Turbo says 'choose the rate sheet to start from'. Turbo also checks the entries with validateRateSheetEntryInputs, whose rules are in another file.

## Proposed rates must start from an active rate sheet

Lines 94-102

Asked as: start from inactive sheet; starting sheet is draft; base rates on draft sheet; start from active rate sheet

The starting sheet must exist and must not be a draft, or Turbo says 'Start the proposed rates from an active rate sheet.' An inactive sheet is allowed only when the request's own draft already started from it.

## Proposed rate with unknown role is refused

Lines 103-107

Asked as: unknown role on rates; role not found proposed rate; wrong role in rate entry; invalid role rates

Each proposed rate entry must use an existing role. Otherwise Turbo says 'A proposed rate has an unknown role.'

## Proposed rates sheet name must be unique

Lines 109-120

Asked as: rate sheet name taken; duplicate rate sheet name; name already used rates; same name as another sheet

The name of proposed rates is optional. If another rate sheet has the same name in any letter case, Turbo refuses with 'another rate sheet already has this name'. The request's own draft does not count as another sheet.

## Proposed rates identical to an active sheet are refused

Lines 122-132

Asked as: rates same as existing sheet; duplicate rate sheet rates; proposed rates already exist; identical rates active sheet

Turbo compares the proposed rates with the active rate sheets (CompareWithActiveRateSheets). If an active sheet has exactly the same rates, Turbo refuses the save and names that sheet, so the requester can choose it instead.

## Draft is rewritten only when the proposal changed

Lines 134-166

Asked as: no change to proposal; resaving same rates; keep draft approval; proposal unchanged draft

A proposal counts as changed when it is new, or when its name, reason, starting sheet or rates differ from the saved draft. If nothing changed, Turbo only points the request at the existing draft and keeps its approval.

## Saving proposed rates creates an inactive draft rate sheet

Lines 151-213

Asked as: where do new rates go; draft rate sheet created; rate sheet not active yet; request points to new sheet

The first save of proposed rates creates a draft rate sheet owned by the request, with revision 0 and active false. Turbo sets its effective date to today, and job creation sets the real date when the sheet becomes active. Turbo updates the draft's entries to match the proposal: it adds, changes and deletes entries by role. It points the request's rate sheet at the draft.

## Changed proposed rates lose their approval

Lines 170-171

Asked as: rates approval lost after edit; edited rates need approval again; approval cleared after change; changed rates unapproved

When a save changes the draft, Turbo clears the rates approver and approval time. The approver must approve the new rates again.

## Job creation can rename the proposed rates

Lines 224-267

Asked as: rename proposed rates; change rate sheet name at job; job rate sheet new name; rename sheet on job creation

When the job is created, the job holder can give the proposed rates a new name. Turbo ignores an empty name or the same name. If another rate sheet has the name, Turbo says 'Another rate sheet already has this name.' The old name is kept in rate_sheet_name_changes, and the approval stays valid. A request with no proposed rates gets 'This request has no proposed rates to rename.'

## Request view shows proposed rates beside the source rates

Lines 269-310

Asked as: compare new rates with source; show rate differences request; changed rate highlighted; see old and new rates

For a request with proposed rates, Turbo adds the draft, its rates approver and the source sheet to the request's related records. Each rate shows the source rate and overtime rate and is marked changed. A role missing from the source sheet is marked changed. Turbo adds nothing when the request uses an active sheet.

## Approver must state approval of the new rate sheet

Lines 322-337

Asked as: approve new rate sheet; approver statement rates; approval text missing rates; confirm creation of rate sheet

When the approver approves a request with proposed rates, they must also approve the creation of the new rate sheet and its rates. Otherwise Turbo says 'Approve the creation of the new rate sheet and its rates.' Turbo then records the approver and the time on the draft. A request with an active sheet needs no such statement.

## New submission clears the proposed rates approval

Lines 339-351

Asked as: resubmit request rates approval; submit again clears approval; new submission approval reset; approval after resubmitting

Turbo clears the draft's rates approver and approval time, so a new submission needs a new approval.

## Job creation needs approved proposed rates

Lines 353-364

Asked as: cannot create job rates; job creation blocked unapproved; recall request rates; proposed rates not approved

Turbo refuses to create a job from a request whose draft was not approved by the request's approver. The approval must have a time. The message is 'The proposed rates are not approved. Recall the request for a new approval.' A request with an active sheet is not checked.

## Compare rates route shows exact and close sheets

Lines 366-402

Asked as: compare my rates; find matching rate sheet; similar existing rate sheets; check rates against active sheets

Anyone with the time claim can compare unsaved proposed rates with the active rate sheets while project requests are enabled. Turbo returns exact matches and a few close sheets. It leaves the starting sheet out of the close ones and changes nothing.
