<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_request_rate_sheets.go","sha256":"1de5f528b0bf70d578bd4ab09eb851cc2f8a698dd5e74eeef567057f280263dd","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_request_rate_sheets.go

## Proposed rates reason is limited to 500 characters

Lines 41-87

Asked as: reason too long; reason character limit; 500 characters; rates reason too long

The reason for proposed rates can have at most 500 characters. Longer text is refused with "the reason can have at most 500 characters".

## Chosen rate sheet in a project request must be active

Lines 53-67

Asked as: rate sheet is draft; draft rate sheet not allowed; select active rate sheet; rate sheet not active

When a project request names a rate sheet, Turbo refuses a sheet that does not exist with "Select a valid rate sheet.". It refuses a draft sheet with "Select an active rate sheet.". The request's own draft of proposed rates is allowed.

## Proposed rates need a name, reason and source sheet

Lines 73-94

Asked as: proposed rates missing name; reason is required; no source rate sheet; proposed rates empty fields

Every save of proposed rates needs a name (proposed_rates.name), a reason and a rate sheet to start from (based_on). Turbo trims spaces first. A missing reason shows "reason is required", and a missing source shows "choose the rate sheet to start from". The rate entries are also checked against the rate sheet entry rules.

## Proposed rates must start from an active rate sheet

Lines 96-104

Asked as: starting from inactive sheet; source rate sheet deactivated; base sheet not active; proposed rates source draft

The source sheet (based_on) must exist and must not be a draft. Otherwise Turbo says "Start the proposed rates from an active rate sheet.". An inactive source is allowed only when it is already the source of the request's own draft, because a draft keeps its source after that sheet is deactivated.

## Proposed rates must use known roles

Lines 105-109

Asked as: unknown role in rates; role not found; invalid role on proposed rate; wrong role in proposed rates

Every entry in proposed rates must use an existing role. Otherwise Turbo refuses with "A proposed rate has an unknown role.".

## Proposed rate sheet name must be unique in any letter case

Lines 111-126

Asked as: rate sheet name taken; duplicate rate sheet name; name already used; same rate sheet name

Turbo refuses a proposed name that another rate sheet already uses, ignoring letter case. Every revision of every sheet counts, but the request's own draft does not. The message is "another rate sheet already has this name".

## Proposed rates that match an active sheet are refused

Lines 128-137

Asked as: rates same as existing sheet; duplicate of active sheet; identical rate sheet exists; proposed rates already exist

If the proposed rates exactly match an active rate sheet, Turbo refuses the save and names that sheet, so the requester can choose it instead. The comparison is done by hooks.CompareWithActiveRateSheets and hooks.DuplicateRateSheetError, which hold the details.

## Saving proposed rates only changes the draft if something changed

Lines 139-171

Asked as: save with no changes; nothing changed on rates; draft not updated; rates saved unchanged

A request's proposed rates are stored as a draft rate sheet (revision 0, not active) that the request owns. A new draft counts as changed. An existing draft counts as changed only if its name, reason, source sheet or rates differ. When nothing changed, Turbo leaves the draft and its approval alone and points the request at it.

## Changed proposed rates lose their approval

Lines 172-218

Asked as: approval lost after edit; rates need reapproval; editing rates clears approval; approval reset on change

When a draft is changed and saved, Turbo clears its rates approver and approval time, so the rates need a new approval. It sets the draft's effective date to today's date, and job creation sets the real date later. Turbo updates entries by role: it adds new roles, updates entries whose rate or overtime rate differ, and deletes roles that are no longer proposed.

## Draft is deleted when the request stops using it

Lines 220-227

Asked as: old draft still there; leftover draft rate sheet; remove unused draft; draft not deleted after switch

Turbo deletes the request's draft rate sheet when the request no longer uses it, after the request is saved with its new rate sheet. It does nothing if there is no draft.

## Job creator can rename proposed rates

Lines 229-275

Asked as: rename proposed rates; change proposed rates name; cannot rename rates; rename rate sheet on job

When a job is created from a request, the person with the job claim can give the proposed rates a new name. An empty or unchanged name does nothing. If the request has no proposed rates, Turbo says "This request has no proposed rates to rename.". If another sheet already has the name, ignoring case, Turbo says "Another rate sheet already has this name.". The approval stands, and Turbo records the old and new names and who renamed it.

## Editor and approver see the proposal beside the source rates

Lines 277-318

Asked as: compare proposed with source; see changed rates; highlight rate differences; which rates differ from source

For a request with proposed rates, Turbo adds the draft's approver, the source sheet and the source rates next to each proposed rate. It marks a rate as changed when the rate or overtime rate differs from the source, or when the role is not on the source sheet. It does nothing if the request uses an active sheet.

## Approver must confirm the new rate sheet

Lines 330-345

Asked as: approve new rate sheet; approval statement missing; confirm rate sheet creation; approver must confirm rates

When the approver approves a request with proposed rates, they must make a separate statement that approves creating the new rate sheet and its rates. Without it Turbo says "Approve the creation of the new rate sheet and its rates.". With it Turbo records the approver and the time on the draft. It does nothing if the request uses an active sheet.

## New submission clears the proposed rates approval

Lines 347-359

Asked as: resubmit clears rates approval; new submission approval reset; request resubmitted rates unapproved; approval removed on submit

Turbo clears the draft's rates approver and approval time, so a new submission needs a new approval. It saves only if an approval was set.

## Job creation needs the approver's approval of proposed rates

Lines 361-372

Asked as: cannot create job rates; proposed rates not approved; job creation blocked; recall request for approval

Turbo refuses to create the job when the draft's rates were not approved. The approval must be by the request's own approver and have an approval time. The message is "The proposed rates are not approved. Recall the request for a new approval.".

## Compare proposed rates with active sheets before saving

Lines 374-409

Asked as: check rates against sheets; compare before saving; find similar rate sheets; preview rate matches

Anyone who can start a project request (the time claim) can compare unsaved proposed rates with the active rate sheets, if project requests are enabled. The result lists exact matches, which a save would refuse, and a few close sheets. The sheet the requester started from is left out of the close sheets. The comparison changes nothing.
