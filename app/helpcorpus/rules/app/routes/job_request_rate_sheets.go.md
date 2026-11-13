<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_request_rate_sheets.go","sha256":"8916d60a875ba77ddae61030aba3f817da72c27d5adddb1fdc0e84acf7f177cd","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_request_rate_sheets.go

## Proposed rates are saved as a draft rate sheet

Lines 14-19

Asked as: save new rates as draft; proposed rates draft sheet; rates saved but not active

A project request that proposes new rates saves them as a draft rate sheet that the request owns (revision 0, not active). The draft is not listed or used elsewhere until job creation activates it. The request points to the draft as its rate sheet.

## Reason limited to 500 characters; optional in a draft

Lines 39-40

Asked as: reason too long; reason character limit; reason 500 characters

The reason for proposed rates can have at most 500 characters (error 'the reason can have at most 500 characters'). A draft request can be saved without a reason, because this check does not require one. Submission requires a reason.

## Chosen rate sheet must be valid and not a draft

Lines 51-65

Asked as: select a valid rate sheet; choose another draft sheet; rate sheet choice refused

When a request chooses a rate sheet, Turbo refuses an unknown sheet with 'Select a valid rate sheet.' It refuses another draft with 'Select an active rate sheet.' The request's own draft is allowed, and an empty choice passes this check.

## Proposed rates need a name and a starting sheet

Lines 71-89

Asked as: proposed rates missing name; rate sheet to start from; can't submit rate request

Proposed rates must have a name (error 'name is required') and a rate sheet to start from (error 'choose the rate sheet to start from'). Turbo trims spaces from the name, reason and starting sheet before it checks them. The rate entries are also validated.

## Proposed rates must start from an active rate sheet

Lines 91-99

Asked as: start from inactive sheet; starting sheet is draft; source sheet deactivated

Turbo refuses proposed rates that start from a missing sheet or a draft sheet with 'Start the proposed rates from an active rate sheet.' It also refuses an inactive sheet. The exception is when the request's own draft already started from that sheet, because a draft keeps its source after the source is deactivated.

## Proposed rate with unknown role is refused

Lines 100-104

Asked as: unknown role in rates; role does not exist; rate role not found

If any proposed rate names a role that does not exist, Turbo refuses the request with 'A proposed rate has an unknown role.'

## Proposed rate sheet name must be unique

Lines 106-121

Asked as: rate sheet name taken; duplicate rate sheet name; name already used by sheet

Turbo refuses a proposed rate sheet name that any other rate sheet revision already has, ignoring letter case, with 'another rate sheet already has this name'. The request's own draft does not count as a conflict.

## Rates that match an active sheet are refused

Lines 123-132

Asked as: rates same as active sheet; duplicate of existing rate sheet; rates already exist as sheet

Turbo refuses proposed rates that exactly match an active rate sheet. The error names that sheet so the requester can choose it instead.

## Draft is rewritten only when something changed

Lines 134-166

Asked as: draft not updated; nothing changed in draft; draft saved again

Turbo saves the draft again only when the name, reason, starting sheet or rates differ from what is stored. A new draft always counts as changed. If nothing changed, the request just points to the existing draft.

## Saving proposed rates sets the draft's effective date and entries

Lines 167-210

Asked as: draft effective date; draft rates not updating; draft entries changed

When it saves a draft, Turbo sets its effective date to today's date (RateSheetToday). Job creation sets the real activation date later. Turbo updates the draft's entries to match the proposal: it adds new roles, updates entries whose rate or overtime rate changed, and deletes roles no longer proposed.

## Unused request draft is deleted

Lines 213-220

Asked as: old draft still there; leftover draft sheet; draft removed after change

After the request is saved with a different rate sheet, Turbo deletes the draft the request no longer uses. It does nothing when there is no draft.

## Job creator can rename proposed rates

Lines 222-268

Asked as: rename proposed rates; change rate sheet name; rename job rates

When a job is created from a request, the job holder (job claim) can give the proposed rates a new name. An empty or unchanged name does nothing. The approval stands, because the approver approved the rates and not the name. Turbo records the old name, the new name and who renamed it in a name-change record.

## Rename refused without a draft or with a taken name

Lines 237-253

Asked as: cannot rename rates; rename fails no draft; rename name already taken

Renaming proposed rates fails with 'This request has no proposed rates to rename.' when the request has no draft. It fails with 'Another rate sheet already has this name.' when another sheet has the name, ignoring letter case.

## Approver sees proposed rates beside source rates

Lines 270-306

Asked as: compare proposed rates source; rates changed highlight; proposal beside original rates

For a request with proposed rates, Turbo shows the proposal and the sheet it started from. Each rate is shown beside the source's rate and overtime rate. A rate is marked changed when its role is new or its rate or overtime rate differs. Nothing is added when the request uses an active sheet.

## Approval checklist wording for proposed rates

Lines 308-315

Asked as: rate sheet checklist wording; approval checklist proposed rates; checklist text for rates

For proposed rates, the rate sheet statement of the approval checklist reads 'The proposed rates are approved for this client.' The key rate_sheet_correct and the checklist version stay the same, so approvals need no new checklist version.
