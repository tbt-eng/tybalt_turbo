<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_request_rate_sheets.go","sha256":"12a0bab430c146a204dbcebad2c24d9c1161819c484b9e907e2bced574b5671e","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_request_rate_sheets.go

## Proposed rates reason is limited to 500 characters

Lines 41-84

Asked as: reason too long; reason character limit; 500 characters; reason text too long

The reason for proposed rates can have at most 500 characters. A longer reason gives 'the reason can have at most 500 characters'. A draft request can be saved with no reason. The code says submission requires one, but this file does not check it.

## Chosen rate sheet in a project request must be active

Lines 53-67

Asked as: rate sheet not active; draft rate sheet chosen; select active rate sheet; invalid rate sheet on project

When a project request chooses a rate sheet, Turbo checks it. An unknown sheet gives 'Select a valid rate sheet.' A draft sheet gives 'Select an active rate sheet.' The request's own draft is allowed, and an empty choice passes.

## Proposed rates need a name and a starting sheet

Lines 73-87

Asked as: proposed rates name missing; blank rates name; no starting rate sheet; pick base rate sheet

Proposed rates in a project request must have a name (name) and a rate sheet to start from (based_on). Turbo trims spaces first. A blank name gives 'name is required'. A blank starting sheet gives 'choose the rate sheet to start from'.

## Proposed rates must start from an active rate sheet

Lines 93-101

Asked as: start from inactive sheet; based on deactivated sheet; proposed rates source inactive; starting sheet not active

The starting sheet of proposed rates must exist and not be a draft. It must also be active. Otherwise Turbo says 'Start the proposed rates from an active rate sheet.' One exception: a sheet that was deactivated is still allowed when it is the source the request's own draft already started from.

## Proposed rates must use known roles

Lines 102-106

Asked as: unknown role in rates; rate role not found; wrong role on rate; invalid role proposed rate

Every rate entry in proposed rates must name an existing role. An unknown role gives 'A proposed rate has an unknown role.'

## Proposed rate sheet name must be unique in any letter case

Lines 108-123

Asked as: rate sheet name taken; duplicate sheet name; name already used rates; same name as sheet

The name of proposed rates must differ from the name of every other rate sheet, in any revision and ignoring letter case. The request's own draft does not count. A clash gives 'another rate sheet already has this name'.

## Proposed rates that match an active sheet are refused

Lines 125-134

Asked as: rates identical to existing sheet; save refused same rates; rates match existing sheet; duplicate rate sheet error

If the proposed rates are exactly the rates of an active rate sheet, Turbo refuses the save. The error names that sheet (via DuplicateRateSheetError) so the requester can choose it instead.

## Saving unchanged proposed rates changes nothing

Lines 136-147

Asked as: save without changes; nothing changed resave; keep approval unchanged edit; resubmit same proposal

Turbo compares a request's existing draft with the new proposal on name, reason, starting sheet and rates. If none differ, it does not rewrite the draft or its entries and keeps the existing approval. A new draft always counts as changed.

## Saving changed proposed rates creates or updates a draft sheet

Lines 153-215

Asked as: draft rate sheet created; where do proposed rates go; rate entries updated draft; new draft revision zero

Changed proposed rates go into a draft rate sheet that the request owns. A new draft has revision 0 and is inactive. Turbo sets its effective date to today, though job creation sets the real date. It adds or updates entries by role, skips unchanged ones, and deletes entries for roles no longer listed. The request then points to the draft.

## Changed proposed rates lose their approval

Lines 172-173

Asked as: approval lost after edit; rates approval reset; need approve again rates; approver cleared on change

When proposed rates change on save, Turbo clears the draft's rates approver and approval time. The approver must approve again.

## Rate sheet draft is deleted when the request stops using it

Lines 217-224

Asked as: old draft sheet remains; leftover draft rate sheet; unused draft removed; stale rate sheet draft

A draft that the request no longer uses is deleted after the request is saved with its new rate sheet.

## Job creation can rename the proposed rates

Lines 226-272

Asked as: rename rates when creating job; change proposed rates name; job holder rename sheet; rename rate sheet on job

When a job is created from a request, the job holder (the job claim) can give the proposed rates a new name. A blank or unchanged name does nothing. The approval stays, because the approver approved the rates and not the name. Turbo records the old name, the new name and who renamed it. A request with no proposed rates gives 'This request has no proposed rates to rename.' A name that another sheet has, in any letter case, gives 'Another rate sheet already has this name.'

## Approver sees proposed rates beside the source sheet

Lines 274-315

Asked as: approver cannot see changes; compare proposed rates source; which rates changed; show changed rate highlights

For a request with proposed rates, Turbo adds the draft's reason, the approver of the rates, and the sheet it started from. Each rate is shown beside the source rate and overtime rate. A rate is marked changed if its role is not in the source or if its rate or overtime rate differs. Nothing is added when the request uses an active sheet.

## Approver must state approval of proposed rates

Lines 327-342

Asked as: approve rates statement; confirm rate sheet approval; approve new rate sheet; approval checkbox missing

When the approver approves a request with proposed rates, they must confirm the approval statement. Otherwise Turbo says 'Approve the creation of the new rate sheet and its rates.' Turbo then records the approver and the time on the draft. It does nothing if the request uses an active sheet.

## New submission clears the proposed rates approval

Lines 344-356

Asked as: resubmit request approval reset; submit again clears approval; new submission approval gone; reset rates approver

Turbo clears the approver and approval time on the request's draft, so a new submission needs a new approval. It does nothing if the request has no draft or no approval is set.

## Job creation requires approved proposed rates

Lines 358-369

Asked as: cannot create job unapproved; job creation blocked rates; rates not approved error; recall request approval

Turbo refuses to create a job from a request with proposed rates unless the draft has an approval time and a rates approver, and that approver is the request's approver. The error is 'The proposed rates are not approved. Recall the request for a new approval.'

## Compare proposed rates with active rate sheets

Lines 371-406

Asked as: compare rates with sheets; check rates before saving; find similar rate sheets; preview rate matches

Anyone with the time claim can compare unsaved proposed rates with the active sheets, but only while project requests are enabled. The result lists the exact matches that a save would refuse, plus the closest near matches. The sheet the requester started from is left out of the near matches. The comparison changes nothing.
