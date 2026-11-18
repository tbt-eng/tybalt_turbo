<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/rate_sheets_api.go","sha256":"c1e5d7060cba75624f6e07b59025e033704544ae515f9ecca5cd6827f53e8599","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/rate_sheets_api.go

## Only job claim holders can activate a rate sheet

Lines 29-36

Asked as: cannot activate rate sheet; no permission to activate; turn on rate sheet; rate sheet won't go live

Activating a rate sheet requires the job claim (procurement). Without it Turbo refuses with "you do not have permission to activate rate sheets".

## Proposed rates of a project request cannot be activated

Lines 44-47

Asked as: activate project request rates; proposed rates won't activate; draft rates activation; request rates can't be turned on

A rate sheet that is a draft (the proposed rates of a project request) cannot be activated by hand. Turbo refuses with the error from hooks.RateSheetDraftError("active"). Those rates are activated when the job is created.

## Activating an already active rate sheet changes nothing

Lines 49-56

Asked as: rate sheet already active; activate twice; sheet is active already; nothing happened activating

If the rate sheet is already active, Turbo does not change it. It replies "rate sheet is already active" and treats the request as a success.

## Older revisions cannot be activated

Lines 58-67

Asked as: old revision won't activate; activate previous version; newer version exists; activating an old sheet

Turbo refuses to activate a rate sheet when a newer revision with the same code exists. The message is "cannot activate - a newer revision exists".

## Incomplete rate sheet cannot be activated

Lines 69-78

Asked as: rate sheet missing roles; can't activate missing entries; sheet has gaps; missing role rates

A rate sheet must have an entry for every role before it is activated. Otherwise Turbo refuses with "rate sheet is incomplete - missing entries for some roles" and lists the missing roles (missing_roles).

## Activation checks duplicates and deactivates other revisions atomically

Lines 80-98

Asked as: two identical sheets active; duplicate sheet activation; activation blocked duplicate; other revisions switched off

Activation runs in one transaction: hooks.ActivateRateSheet checks for a duplicate sheet, deactivates the other revisions and activates this one. Two identical sheets cannot both be activated at the same time. A refusal from that check is returned with its own message and status. This file does not state which checks that function makes.

## Only job claim holders can deactivate a rate sheet

Lines 108-115

Asked as: cannot turn off rate sheet; no permission to deactivate; switch off rate sheet; rate sheet won't stop

Deactivating a rate sheet requires the job claim (procurement). Without it Turbo refuses with "you do not have permission to deactivate rate sheets".

## Deactivating an inactive rate sheet changes nothing

Lines 123-142

Asked as: rate sheet already inactive; deactivate twice; sheet is off already; turn off nothing happened

If the rate sheet is already inactive, Turbo replies "rate sheet is already inactive" and changes nothing. Otherwise it sets the sheet inactive. Deactivation runs no completeness or duplicate check.

## Only admins can edit rate sheet entries

Lines 159-166

Asked as: cannot change rate entry; edit rate needs admin; rate changes not allowed; only admins change rates

Updating a rate sheet entry requires the admin claim. Without it Turbo refuses with "you do not have permission to update rate sheet entries".

## Entry update needs a rate or an overtime rate

Lines 174-204

Asked as: empty rate update; update rate with nothing; overtime rate only change; which rate fields to change

An entry update must include rate, overtime_rate or both. Otherwise Turbo refuses with "at least one of rate or overtime_rate must be provided". Only these two fields can be changed, and only the ones sent are updated.

## Entry rate must be a whole number of at least 1

Lines 186-195

Asked as: rate has decimals; rate below one; whole number rate; rate value too low

When an entry update sets a rate, a value below 1 is refused with "rate must be at least 1". A value with decimals is refused with "rate must be a whole number".

## Entry overtime rate minimum and decimals

Lines 196-204

Asked as: overtime rate decimals; overtime rate too low; overtime more than two decimals; overtime value rejected

When an entry update sets an overtime rate, a value below 1 is refused with "overtime_rate must be at least 1". More than two decimal places is refused with hooks.OvertimeRateDecimalsMessage.

## Entry edit on active sheet may not duplicate another active sheet

Lines 206-241

Asked as: rates same as active sheet; edit makes duplicate active; change matches other sheet; entry change blocked duplicate

When the entry's sheet is active, Turbo applies the change and compares the sheet's rates with the other active sheets. Revisions of the same code are excluded. If the rates would equal another active sheet's, Turbo refuses with the duplicate error and does not save. An inactive sheet is not checked here, because activation checks it.

## Entries of proposed rates cannot be edited directly

Lines 214-217

Asked as: cannot edit proposed rates; change project request rates; draft rate entry locked; edit requester rates

An entry that belongs to a draft rate sheet (a project request's proposed rates) cannot be edited here. Turbo refuses with hooks.RateSheetDraftError("entry"). The requester changes proposed rates through the request.

## Rate sheet entry edit failure responses

Lines 242-248

Asked as: entry save failed; rate update error; could not save entry; saving rate entry failing

When saving an edited rate sheet entry fails, Turbo returns the hook's own error message and status if the failure is a rule such as a draft or a duplicate. Any other failure gives "failed to update rate sheet entry".

## Successful entry edit returns new rates

Lines 250-255

Asked as: rate update confirmation; what rates were saved; entry updated reply; saved overtime rate shown

After a rate sheet entry is updated, Turbo replies "rate sheet entry updated" with the entry's saved rate and overtime rate.

## Rate sheet entries must be valid

Lines 342-371

Asked as: rate sheet needs entries; duplicate role in sheet; sheet entry validation; role missing from entry

A new rate sheet or revision needs at least one entry ("at least one entry is required"). Each entry needs a role ("entry role is required"), and a role may appear once ("entry role must be unique"). The rate must be a whole number of at least 1. The overtime rate must be at least 1 with at most two decimal places. Errors name the entry by its position.

## New rate sheet needs the job claim

Lines 389-397

Asked as: cannot create rate sheet; make new sheet no permission; need job claim new sheet; copy rate sheet refused

Creating a new rate sheet or a copy (revision 0) requires the job claim. Without it Turbo refuses with "you do not have permission to create rate sheets".

## Revising a rate sheet needs revise or admin claim

Lines 398-411

Asked as: cannot revise rate sheet; make new revision no permission; revision permission denied; update sheet version refused

Creating a revision (revision above 0) requires the rate_sheet_revise claim or the admin claim. Without either Turbo refuses with "you do not have permission to revise rate sheets".

## Revision keeps the name of its sheet

Lines 415-418

Asked as: revision without name; revision name missing; revision takes sheet name; name left blank revision

The name is optional. A revision with no name takes the name of the sheet it revises. A name is trimmed of spaces.

## Effective date is required and must be a valid YYYY-MM-DD date

Lines 419-428

Asked as: effective date missing; wrong date format; invalid effective date; date not accepted rate sheet

A rate sheet needs an effective date. A missing one gives "effective date is required". A date not in the form YYYY-MM-DD gives "effective date must be in YYYY-MM-DD format". A date in that form must also be a real date, checked by utilities.IsValidDate.

## Revision number cannot be negative

Lines 430-433

Asked as: negative revision number; revision below zero; revision number invalid

A rate sheet request with a revision below 0 is refused with "revision must be 0 or greater".

## New rate sheet cannot send a code

Lines 441-447

Asked as: code on new sheet; set code for new sheet; new sheet code refused; choose code myself

A new rate sheet (revision 0) cannot send a code. Turbo refuses with "a new rate sheet gets its code when it is saved". It assigns the code when it saves the sheet.

## Revision needs an existing rate sheet to revise

Lines 448-459

Asked as: no sheet to revise; revision cannot find sheet; wrong code for revision; revise by name not found

A revision finds the sheet it revises by code. If no code is sent, it finds the sheet by name. When no sheet is found, Turbo refuses with "there is no rate sheet to revise". The revision takes the code of the latest revision of that sheet.

## Proposed rates cannot be revised

Lines 460-464

Asked as: cannot revise proposed rates; revise project request rates; draft sheet revision blocked; revising request rates

A draft rate sheet (a project request's proposed rates) has no revisions. Trying to revise it gives the error from hooks.RateSheetDraftError("name"). The sheet becomes an ordinary sheet when the job is created.

## Revision effective date cannot precede previous revision

Lines 472-483

Asked as: effective date earlier than before; revision date too early; date before previous revision; effective date conflicts previous

A revision's effective date may equal the previous revision's date but not be earlier. Otherwise Turbo refuses with "effective date must be on or after the previous revision's effective date" and gives the previous date (previous_effective_date).

## New rate sheet is saved inactive with all its entries at once

Lines 485-586

Asked as: new sheet saved inactive; partial sheet saved; sheet not active after creating; why is new sheet off

Turbo creates the rate sheet and all its entries in one transaction, so no partial sheet is saved. The new sheet is inactive (active false) and must be activated later.

## Duplicate revision number is refused

Lines 497-508

Asked as: revision already exists; same revision number twice; duplicate revision number; revision number taken

A revision whose code and revision number already exist is refused with "this rate sheet already has this revision".

## New rate sheet name must be unique in any letter case

Lines 509-524

Asked as: rate sheet name taken; duplicate sheet name; name already used; same name different case

A new rate sheet (revision 0) cannot use a name another sheet holds, whatever the letter case. Several sheets may have no name. Proposed rates of a project request hold their name too, and the message says so. The message comes from hooks.RateSheetNameHeld.

## New sheet cannot match an active sheet's rates

Lines 526-535

Asked as: rates match active sheet; new sheet same rates; identical rates refused; revision same rates as active

Turbo refuses to save a new sheet or revision whose rates are identical to another active sheet. Revisions of the revised sheet are excluded, so a revision that changes only the effective date is allowed.

## Rate sheet creation failure responses

Lines 588-610

Asked as: rate sheet conflict error; code and revision clash; sheet creation rejected; bad entry position error

A conflict from a duplicate name or revision is answered as a conflict with its message. A code-and-revision clash gives "a rate sheet with this code and revision already exists". A bad entry is reported as a validation error with the entry's position.

## Created rate sheet response includes its code and entry count

Lines 612-618

Asked as: new sheet code shown; sheet response entry count; which code was assigned; how many entries created

When Turbo creates a rate sheet, it replies with the new sheet's id, its code, its name, its revision number and the number of entries it created (entries_created). A new sheet's code is the one Turbo assigns when it saves the sheet. A revision shows the code of the sheet it revises.

## Compare rate sheets requires job, revise or admin claim

Lines 663-676

Asked as: cannot compare rate sheets; compare permission denied; comparison not allowed; check rates no access

Comparing unsaved entries with active rate sheets requires the job, rate_sheet_revise or admin claim. Otherwise Turbo refuses with "you do not have permission to compare rate sheets".

## Compare shows exact and near matches without saving

Lines 683-699

Asked as: compare rates before saving; check similar rate sheets; find matching sheets; preview rate matches

The comparison changes nothing. It returns exact matches, which a save or activation would block, and a few close sheets as information. The revised sheet's own code is excluded when it is found by code or by name.

## Finding the sheet a revision revises

Lines 703-722

Asked as: which sheet gets revised; latest revision lookup; find sheet by code; no code or name

Turbo finds the latest revision of the sheet with the given code. With no code, it uses the sheet with the given name (the highest revision). With neither a code nor a name, it finds no sheet.
