<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/rate_sheets_api.go","sha256":"5ea763e7f468c51c30bb9ea781d623de57cca753ad5b0b8653851311d819609d","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/rate_sheets_api.go

## Only procurement can activate a rate sheet

Lines 29-36

Asked as: can't activate rate sheet; no permission to turn on sheet; who can activate rates; activate button denied

Activating a rate sheet requires the job claim (procurement). Anyone else gets "you do not have permission to activate rate sheets".

## Proposed rates in a project request cannot be activated directly

Lines 44-47

Asked as: activate proposed rates; project request rates won't activate; turn on draft rates; proposed rates activation

If the rate sheet is a draft of a project request's proposed rates (hooks.IsRateSheetDraft), activation is refused with the draft error from hooks.RateSheetDraftError("active"). Those rates are activated when the job is created.

## Activating an already active rate sheet does nothing

Lines 49-56

Asked as: rate sheet already on; activate sheet twice; sheet is already live; activating active sheet

If the rate sheet is already active, Turbo answers "rate sheet is already active" and changes nothing.

## An older revision cannot be activated

Lines 58-67

Asked as: old version won't activate; activate previous revision; newer version exists error; activate an old sheet

Activation is refused with "cannot activate - a newer revision exists" when a newer revision of the same code exists.

## An incomplete rate sheet cannot be activated

Lines 69-78

Asked as: missing roles on rate sheet; rate sheet won't activate; sheet missing a role; incomplete rate sheet error

Activation is refused with "rate sheet is incomplete - missing entries for some roles" when the sheet lacks an entry for some role. The response lists the missing roles.

## Activation checks duplicates and deactivates other revisions atomically

Lines 80-98

Asked as: activation failed duplicate; other revisions turned off; activate gives error; activation transaction error

Activation runs in one transaction through hooks.ActivateRateSheet. That function checks for a duplicate sheet and deactivates the other revisions. If it fails with a hook error, Turbo returns that error. Otherwise the answer is "rate sheet activated".

## Only procurement can deactivate a rate sheet

Lines 108-115

Asked as: can't turn off rate sheet; deactivate button denied; who can disable rates; no permission deactivate sheet

Deactivating a rate sheet requires the job claim (procurement). Anyone else gets "you do not have permission to deactivate rate sheets".

## Deactivating a rate sheet

Lines 123-142

Asked as: turn off rate sheet; disable rate sheet; sheet already inactive; switch off rates

Deactivation sets active to false and answers "rate sheet deactivated". If the sheet is already inactive, Turbo answers "rate sheet is already inactive" and changes nothing. Deactivation has no draft check and no check of jobs that use the sheet.

## Only admins can edit rate sheet entries

Lines 159-166

Asked as: can't change rates; edit rate entry denied; only admins edit rates; change rate permission

Updating a rate sheet entry requires the admin claim. Anyone else gets "you do not have permission to update rate sheet entries". Only rate and overtime_rate can change.

## Entry update needs rate or overtime rate

Lines 174-183

Asked as: nothing to update rate; rate entry not found; update rate blank; no rate sent

An entry update must send at least one of rate or overtime_rate. Otherwise Turbo says "at least one of rate or overtime_rate must be provided". An unknown entry gives "rate sheet entry not found".

## Entry rate limits

Lines 185-201

Asked as: rate can't be zero; rate decimal not allowed; overtime rate decimals; rate too low

A rate must be at least 1 ("rate must be at least 1") and a whole number ("rate must be a whole number"). An overtime rate must be at least 1 ("overtime_rate must be at least 1") but can have decimals. Only the fields sent are changed.

## Entry edit cannot duplicate another active sheet

Lines 203-238

Asked as: rates same as another sheet; change rate duplicates sheet; rates match active sheet; edit makes duplicate rates

When the sheet is active, an entry change is refused if the resulting rates would exactly equal another active sheet's rates. The sheet's own revisions are excluded. Turbo returns the duplicate error from hooks.DuplicateRateSheetError. An inactive sheet is not checked, because activation checks it.

## Entry edit blocked on draft sheets

Lines 211-214

Asked as: can't edit proposed rates; change project request rates; edit draft rate entry; draft rates locked

An entry of a project request's proposed rates (a draft sheet) cannot be changed here. Turbo returns the draft error from hooks.RateSheetDraftError("entry"). The requester changes proposed rates through the request.

## Entry update errors from hooks are shown as given

Lines 239-252

Asked as: update rate entry failed; rate entry update error; saved rate entry message; entry update message

When updating a rate sheet entry fails with a rule error (a draft sheet or a duplicate of another active sheet), Turbo shows that error with its own status. Any other failure shows 'failed to update rate sheet entry'. On success Turbo returns the message 'rate sheet entry updated' with the saved rate and overtime rate.

## Only the first error for each field is reported

Lines 288-295

Asked as: only one error shown; missing other error messages; why only first error; hidden field errors

When a rate sheet request has several problems with the same field, Turbo keeps the first error for that field and drops the later ones. The person fixes one problem at a time for each field.

## Entry errors name the entry by position

Lines 331-337

Asked as: which row has error; entries error position; role must be unique error; entries field error

When saving an entry fails validation, Turbo reports the error under 'entries.N.field', where N is the entry's position in the list. A database unique-constraint failure on an entry is reported as 'entry role must be unique' for that entry.

## Rate sheet entry input rules

Lines 339-366

Asked as: rate sheet needs entry; entry without role; rate sheet save validation; entry input rules

A new sheet needs at least one entry ("at least one entry is required"). Each entry needs a role ("entry role is required"), and each role may appear once ("entry role must be unique"). Each rate must be at least 1 and a whole number, and each overtime rate must be at least 1. Errors name the entry by position, such as entries.0.rate. Overtime rate may have decimals.

## Entry role must be unique within a rate sheet

Lines 343-356

Asked as: duplicate role on sheet; same role twice; role repeated in rate sheet; role listed again

A rate sheet's entries may use each role once. When an entry repeats a role, Turbo refuses with 'entry role must be unique'. The error names the entry by its position, such as entries.2.role.

## Entry rate must be whole and at least 1; overtime rate at least 1

Lines 343-366

Asked as: new sheet rate zero; overtime rate too low; sheet with no entries; new sheet rate decimal

When creating a rate sheet, each entry's rate must be at least 1 ('entry rate must be at least 1') and a whole number ('entry rate must be a whole number'). The overtime rate must be at least 1 ('entry overtime_rate must be at least 1'), and it may have decimals. A sheet with no entries is refused with 'at least one entry is required'. An entry without a role is refused with 'entry role is required'.

## Creating a new rate sheet needs the job claim

Lines 384-392

Asked as: can't create rate sheet; new rate sheet permission; copy rate sheet denied; create sheet not allowed

Creating a rate sheet with revision 0 (new or copy) requires the job claim (procurement). Otherwise: "you do not have permission to create rate sheets".

## Revising a rate sheet needs revise or admin claim

Lines 393-406

Asked as: can't revise rate sheet; revision permission denied; make new revision not allowed; revise permission error

Creating a revision (revision above 0) requires the rate_sheet_revise claim or the admin claim. Otherwise: "you do not have permission to revise rate sheets".

## Rate sheet validation errors returned together

Lines 408-434

Asked as: validation error rate sheet; several errors at once; why validation error appears; rate sheet errors together

Turbo collects the effective date, revision and entry errors and returns them together as 'Validation error' before it looks for the sheet to revise or saves anything.

## Rate sheet name is trimmed and optional

Lines 410-413

Asked as: rate sheet name blank; revision name missing; trim sheet name; sheet name optional

Turbo trims spaces from a new rate sheet's name. The name is an optional description. A revision with no name takes the name of the sheet it revises.

## Effective date required and checked

Lines 414-423

Asked as: effective date error; date format rate sheet; effective date required; wrong date format sheet

A new rate sheet or revision needs an effective date in YYYY-MM-DD format and it must be a valid date. Messages: "effective date is required" and "effective date must be in YYYY-MM-DD format". The date is trimmed of spaces first.

## Revision number cannot be negative

Lines 425-428

Asked as: revision negative; revision minus number; revision below zero; negative revision number

A revision below 0 is refused with "revision must be 0 or greater".

## New rate sheet cannot send a code

Lines 436-442

Asked as: code not allowed new sheet; enter code for new rate sheet; new sheet code error; why no code field

A new sheet (revision 0) cannot send a code. Turbo says "a new rate sheet gets its code when it is saved" and assigns the code on save.

## Revision must find an existing sheet

Lines 443-465

Asked as: no rate sheet to revise; revision can't find sheet; which sheet to revise; sheet not found for revision

A revision finds its sheet by code, or by name if no code is sent (using the latest revision). If none is found, Turbo says "there is no rate sheet to revise". The revision takes that sheet's code and, if no name is sent, its name.

## Proposed rates cannot be revised

Lines 455-459

Asked as: revise proposed rates; change project request rates revision; draft rates can't revise; revise draft sheet

A project request's proposed rates (a draft sheet) cannot be revised. Turbo refuses with the draft error from hooks.RateSheetDraftError("name").

## Revision effective date cannot precede previous revision

Lines 467-478

Asked as: date earlier than last revision; revision date too early; effective date before previous; date before last revision

A revision's effective date must be on or after the previous revision's date. The same date is allowed. Otherwise Turbo says "effective date must be on or after the previous revision's effective date" and gives the previous date.

## Duplicate revision number refused

Lines 492-503

Asked as: revision number already exists; same revision twice; revision already there; duplicate revision number

Creating a revision whose number already exists for that code is refused with "this rate sheet already has this revision". A unique-constraint failure on code and revision gives "a rate sheet with this code and revision already exists".

## Rate sheet names must be unique

Lines 504-519

Asked as: rate sheet name taken; name already used sheet; duplicate sheet name; name in use

A new sheet (revision 0) cannot take a name another sheet holds, in any letter case. Several sheets may have no name. Proposed rates in a project request also hold their name. Turbo returns the message from hooks.RateSheetNameHeld.

## New sheet cannot match an active sheet's rates

Lines 521-530

Asked as: rates same as active sheet; sheet duplicates existing rates; save blocked duplicate rates; rates already exist

Saving a sheet is refused if its rates exactly match another active sheet. The sheet's own revisions are excluded, so a revision that changes only the effective date is allowed. Turbo returns the duplicate error from hooks.DuplicateRateSheetError.

## Created sheet starts inactive and saves with entries atomically

Lines 532-613

Asked as: new sheet inactive; sheet saved not active; entries failed nothing saved; sheet saved without entries

A new sheet or revision is saved inactive (active false) with all its entries in one transaction. If any entry fails, nothing is saved. The answer gives the id, code and number of entries created. Entry errors name the entry position.

## Revision conflict and name conflict return a conflict error

Lines 583-605

Asked as: revision conflict error; code revision already exists; name conflict sheet; sheet conflict message

When a revision number already exists for the code, Turbo refuses with 'this rate sheet already has this revision'. A name already held by another sheet is refused with the message from the name check. A save that hits a database unique constraint on code and revision returns 'a rate sheet with this code and revision already exists'.

## Successful creation returns new sheet details

Lines 607-613

Asked as: what does create return; sheet id after save; new sheet details; details after creating sheet

After a rate sheet is created, Turbo returns its id, its code, its name, its revision and the number of entries created. A new sheet gets its code when it is saved.

## Compare treats the revised sheet's own revisions as not duplicates

Lines 641-686

Asked as: compare includes own revisions; compare duplicates same sheet; which sheets compared; compare excludes revisions

When comparing rates, Turbo finds the sheet being revised by code, or by name if there is no code. It leaves that sheet's revisions out of the comparison. A new sheet sends neither code nor name, so it is compared with all active sheets.

## Compare rate sheets permission

Lines 658-671

Asked as: can't compare rates; compare rate sheets denied; compare button not allowed; compare permission error

Comparing unsaved rates with active sheets requires the job, rate_sheet_revise or admin claim. Otherwise: "you do not have permission to compare rate sheets".

## Compare returns exact and near matches and changes nothing

Lines 678-695

Asked as: similar rate sheets; matching rates check; compare shows matches; rates close to other sheet

Compare skips the active revision of the sheet being revised, found by code or by name. It returns the exact matches, which would block a save or activation, and a few close sheets as information (the limits come from hooks.RateSheetNearMatchMaxRoles and hooks.RateSheetNearMatchLimit).

## Compare is limited to near matches count

Lines 691-694

Asked as: how many similar sheets; compare list limit; near match count; how many matches shown

Compare returns the exact matches and a short list of close sheets. The list size is set by hooks.RateSheetNearMatchMaxRoles and hooks.RateSheetNearMatchLimit.

## Revision finds its sheet by code, or by name if no code

Lines 698-718

Asked as: revise sheet by name; which code for revision; revise sheet without code; revision lookup by name

A revision looks up the latest revision of the sheet with the given code. If no code is sent, Turbo uses the sheet with the given name. With neither, it finds nothing and the revision is refused with 'there is no rate sheet to revise'. The new revision takes the code of the sheet it revises.
