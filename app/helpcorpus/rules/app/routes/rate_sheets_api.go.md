<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/rate_sheets_api.go","sha256":"cc4ae278fd2f84f86a843bffa739fd27908851aee82427feda6592841116c110","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/rate_sheets_api.go

## Only job claim holders can activate a rate sheet

Lines 29-42

Asked as: activate rate sheet no permission; can't turn on rate sheet; enable rates access denied; who can activate rates

Activating a rate sheet requires the procurement claim (the job claim). Anyone else sees "you do not have permission to activate rate sheets". If the sheet id does not exist, Turbo says "rate sheet not found".

## Proposed rates of a project request cannot be activated

Lines 44-47

Asked as: activate request proposed rates; turn on draft rates; proposed rates won't activate; draft rate sheet activation

If the rate sheet is a draft (hooks.IsRateSheetDraft), activation is refused with the error from hooks.RateSheetDraftError("active"). A request's proposed rates are activated when the job is created.

## Activating an already active rate sheet changes nothing

Lines 49-56

Asked as: rate sheet already on; turn on again nothing happens; activate sheet already active

If the rate sheet is already active, Turbo answers "rate sheet is already active" and does nothing else. It is not an error.

## An older revision cannot be activated

Lines 58-67

Asked as: old revision won't activate; activate previous version of rates; newer version exists error

If a newer revision of the same code exists, activation is refused with "cannot activate - a newer revision exists".

## Incomplete rate sheet cannot be activated

Lines 69-78

Asked as: rate sheet missing roles; sheet won't activate missing role; roles missing before activating

A rate sheet must have an entry for every role before activation. Otherwise Turbo says "rate sheet is incomplete - missing entries for some roles" and lists the missing roles (missing_roles).

## Activation checks duplicates and deactivates other revisions atomically

Lines 80-98

Asked as: two identical rate sheets; same rates activated twice; duplicate rate sheet activation

Activation runs in one transaction (hooks.ActivateRateSheet). The duplicate check, the deactivation of other revisions and the activation of this sheet succeed or fail together, so two identical sheets cannot both be activated. Errors from that step are returned with their own status and message. Turbo does not state here what counts as a duplicate.

## Error responses when activating or editing a rate sheet fail

Lines 83-92

Asked as: activation failed generic error; rate sheet fails to save; why did activation fail

When a rate sheet activation fails a duplicate or other hook check, Turbo returns the error that the check gives, with its own status. Other failures return "failed to activate rate sheet". Entry edits likewise return "failed to update rate sheet entry" for unexpected failures.

## Only job claim holders can deactivate a rate sheet

Lines 108-121

Asked as: turn off rate sheet no permission; can't deactivate rate sheet; disable rate sheet access

Deactivating a rate sheet requires the procurement claim (the job claim). Anyone else sees "you do not have permission to deactivate rate sheets". An unknown id gives "rate sheet not found".

## Deactivating an inactive rate sheet changes nothing

Lines 123-142

Asked as: rate sheet already off; deactivate draft sheet; turn off inactive rates

If the rate sheet is already inactive, Turbo answers "rate sheet is already inactive". Otherwise it sets active to false and answers "rate sheet deactivated". Unlike activation, deactivation does not refuse draft sheets.

## Only admins can edit rate sheet entries

Lines 152-166

Asked as: can't change entry rates; edit rate sheet entry permission; only admin change rates

Updating a rate sheet entry requires the admin claim. Others see "you do not have permission to update rate sheet entries". Only rate and overtime_rate can be changed.

## Entry update needs rate or overtime rate

Lines 174-183

Asked as: empty entry update; change entry nothing sent; entry not found rate

An entry update must include at least one of rate or overtime_rate. Otherwise Turbo says "at least one of rate or overtime_rate must be provided". An unknown entry gives "rate sheet entry not found".

## Entry rate limits on update

Lines 185-201

Asked as: rate below one; rate decimal not allowed; change only overtime rate

A new rate must be at least 1 ("rate must be at least 1") and a whole number ("rate must be a whole number"). A new overtime rate must be at least 1 ("overtime_rate must be at least 1"), and it may have decimals. Only the fields sent are changed.

## Entry update stores rate as whole number, overtime rate may be fractional

Lines 186-201

Asked as: rate with cents; overtime decimals allowed; fractional overtime rate

On an entry edit, Turbo refuses a rate below 1 with "rate must be at least 1" and a rate with decimals with "rate must be a whole number". Overtime rate below 1 is refused with "overtime_rate must be at least 1", but a decimal overtime rate is allowed.

## Entry edit on active sheet cannot duplicate another active sheet

Lines 203-238

Asked as: edit rate makes duplicate; rates match another active sheet; entry change duplicate error

When the entry's sheet is active, Turbo applies the change and compares the sheet's rates with other active sheets, excluding other revisions of the same code. If they would match exactly, the update is refused with the duplicate error. An inactive sheet is not checked here, because activation checks it.

## Entries of proposed rates cannot be edited directly

Lines 211-214

Asked as: edit request rates directly; change proposed rate entry; project request rates locked

If the entry belongs to a project request's proposed rates (a draft sheet), the update is refused with the error from hooks.RateSheetDraftError("entry"). The requester changes proposed rates through the request.

## Only the first error per field is reported

Lines 288-295

Asked as: only one error shown; hidden validation errors; fix one error another appears

When a rate sheet request breaks several rules on the same field, Turbo keeps only the first error for that field. The other errors on that field are dropped (setValidationError). Fixing it can reveal the next error on that field.

## Unique-error conflicts on rate sheet creation

Lines 297-316

Asked as: code revision not unique; duplicate code error on save; unique validation error create

A failed save is treated as a duplicate (code and revision) when any validation error is coded validation_not_unique or its text contains 'unique'. The error data then uses the key 'revision' with the message 'code and revision must be unique' (rateSheetUniqueValidationErrors).

## Rate sheet entry validation

Lines 339-366

Asked as: need at least one entry; duplicate role in sheet; entry rate must be whole

A new sheet or a request's proposed rates need at least one entry ("at least one entry is required"). Each entry needs a role ("entry role is required"), and a role may appear only once ("entry role must be unique"). Each rate must be at least 1 ("entry rate must be at least 1") and a whole number ("entry rate must be a whole number"). Each overtime rate must be at least 1 ("entry overtime_rate must be at least 1"), and it may have decimals.

## Creating a new rate sheet or copy needs job claim

Lines 383-392

Asked as: create rate sheet no permission; copy rate sheet permission; new rate sheet denied

Creating a rate sheet with revision 0 (new or copy) requires the procurement claim (the job claim). Otherwise Turbo says "you do not have permission to create rate sheets".

## Revising a rate sheet needs revise or admin claim

Lines 393-406

Asked as: revise rate sheet no permission; make revision denied; can't create new revision

Creating a revision (revision above 0) requires the rate_sheet_revise claim or the admin claim. Otherwise Turbo says "you do not have permission to revise rate sheets".

## Rate sheet creation validation reports all errors together

Lines 408-434

Asked as: multiple validation errors at once; all field errors listed; validation error list create

Turbo checks the effective date, the revision number and the entries together. It refuses with "Validation error" listing every field that failed, before it checks the code or looks for the sheet to revise.

## New rate sheet cannot send a code

Lines 410-442

Asked as: code not allowed on new sheet; why is code rejected; code assigned automatically

A sheet with revision 0 may not send a code. Turbo refuses with "a new rate sheet gets its code when it is saved" and assigns the code itself on save. Name is optional and trimmed.

## Effective date is required and must be valid

Lines 414-423

Asked as: effective date blank error; wrong date format rate sheet; YYYY-MM-DD date error

A new rate sheet needs an effective date in YYYY-MM-DD format. A blank date gives "effective date is required". A wrong format gives "effective date must be in YYYY-MM-DD format". A date that fails utilities.IsValidDate also gives an error. Spaces around the date are trimmed.

## Revision cannot be negative

Lines 425-428

Asked as: negative revision number; revision below zero; minus revision error

A revision below 0 is refused with "revision must be 0 or greater".

## Revision finds its sheet by code, or by name if no code

Lines 443-454

Asked as: no rate sheet to revise; revise by name; which sheet gets revised

A revision takes the code of the latest revision of the sheet it revises. If no code is sent, Turbo finds the sheet by name instead. If no sheet is found, Turbo says "there is no rate sheet to revise". If the name is blank too, nothing is found.

## Proposed rates cannot be revised

Lines 455-459

Asked as: revise request rates; draft rates revision denied; new revision of proposed rates

A project request's proposed rates (a draft sheet) have no revisions. Revising one is refused with the error from hooks.RateSheetDraftError("name").

## Revision keeps the sheet's name

Lines 460-465

Asked as: revision name missing; revision loses its name; name carried over revision

If a revision sends no name, it takes the name of the sheet it revises.

## Revision effective date cannot be before the previous revision

Lines 467-478

Asked as: effective date earlier than previous; revision date too early; backdated revision date

A revision's effective date must be on or after the previous revision's date. Otherwise Turbo says "effective date must be on or after the previous revision's effective date" and gives previous_effective_date. The same date is allowed.

## Revision effective date check error includes previous date

Lines 468-478

Asked as: previous effective date shown; date before last revision; same date revision allowed

When a revision's effective date falls before the previous revision's, Turbo refuses with "effective date must be on or after the previous revision's effective date" and shows the previous effective date. A date equal to the previous one is allowed.

## Duplicate revision number is refused

Lines 492-503

Asked as: revision number already exists; revision conflict error; this revision exists already

If the sheet already has a revision with the requested number, Turbo answers with a conflict: "this rate sheet already has this revision".

## New rate sheet name must be unique, ignoring letter case

Lines 504-523

Asked as: rate sheet name taken; duplicate sheet name; same name different case

A new sheet (revision 0) cannot share a name with another rate sheet, in any letter case. Turbo says "another rate sheet already has this name". Several sheets may have no name. If the name is held by a request's proposed rates, Turbo refuses with the draft error from hooks.RateSheetDraftError("name").

## New sheet cannot match an active sheet's rates

Lines 525-534

Asked as: same rates as active sheet; revision only changes date; duplicate rates blocked save

A save is blocked when its rates are identical to another active sheet's, with the duplicate error. The other revisions of the revised sheet are excluded, so a revision that changes only the effective date is allowed.

## New rate sheets start inactive and are created with entries atomically

Lines 536-609

Asked as: new rate sheet inactive; why is new sheet off; sheet saved partly

Turbo creates the rate sheet with active set to false, together with all its entries, in one transaction. If any step fails, nothing is saved. A new sheet gets its code on save. Entry errors are reported with keys like entries.N.field, and a repeated role in the database gives "entry role must be unique".

## Entry validation errors are tied to the entry position

Lines 559-582

Asked as: error on entry row; which entry has the error; duplicate role row

When an entry fails to save, Turbo reports the error against that entry's position in the list. A role repeated within the sheet gives "entry role must be unique".

## Duplicate revision or name returns conflict message

Lines 587-609

Asked as: already has this revision message; name already held conflict; code revision exists message

A revision number that already exists on a sheet is refused with "this rate sheet already has this revision". A new sheet whose name another sheet holds is refused with "another rate sheet already has this name". A failed save on code and revision gives "a rate sheet with this code and revision already exists".

## Duplicate code and revision conflict

Lines 599-607

Asked as: code and revision already exists; unique constraint save failed; conflict on save code

If saving fails on a unique constraint, Turbo says "a rate sheet with this code and revision already exists".

## Successful creation returns id, code and entry count

Lines 611-617

Asked as: what does create return; new sheet id response; entries created count

A successful create returns the new sheet's id, code, name, revision and entries_created.

## Comparing rate sheets needs a job, revise or admin claim

Lines 655-675

Asked as: compare rates permission denied; can't compare rate sheets; comparison access needed

Comparing unsaved entries with active rate sheets needs the job claim, the rate_sheet_revise claim or the admin claim. Otherwise Turbo says "you do not have permission to compare rate sheets". The comparison changes nothing.

## Compare excludes the revised sheet and returns exact and near matches

Lines 682-698

Asked as: near matches shown; compare results exact matches; similar rate sheets list

A comparison excludes the revised sheet's own code, found by code or else by name. It returns exact matches, which a save would block, and a few near matches as information (limited by hooks.RateSheetNearMatchMaxRoles and RateSheetNearMatchLimit).

## Compare and save refuse rate sheets by exact rates only

Lines 690-698

Asked as: near match blocks save; similar rates allowed; only identical rates blocked

Turbo treats a rate sheet as a duplicate only when its rates match another active sheet's exactly. Near matches are shown for information only and never block a save.

## Finding the sheet a revision or comparison refers to

Lines 702-722

Asked as: latest revision by code; find sheet by name; empty code and name

Turbo looks up the sheet by code and takes the latest revision of that code's sheets (revisedRateSheet). When the code is empty, it finds the sheet by exact name, takes the highest revision with that name, and uses that sheet's code. When code and name are both empty, no sheet is found. A revision then fails with 'there is no rate sheet to revise', and a comparison excludes no sheet.
