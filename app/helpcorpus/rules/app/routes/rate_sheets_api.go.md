<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/rate_sheets_api.go","sha256":"af8ac84c6ad990cff579974aa12ae0803ab24e0ebb0639b16a828f2d59d4db9c","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/rate_sheets_api.go

## Only the job claim can activate a rate sheet

Lines 29-36

Asked as: can't activate rate sheet; no permission to activate; who can activate rates; activate rate sheet denied

Activating a rate sheet requires a person with the job claim (procurement). Anyone else gets "you do not have permission to activate rate sheets".

## Activating a missing rate sheet is refused

Lines 38-42

Asked as: rate sheet missing; activate sheet doesn't exist; rate sheet not found; can't find rate sheet

When a committer-style user with the job claim activates a rate sheet that does not exist, Turbo refuses with "rate sheet not found". The same message appears when deactivating a missing rate sheet.

## Activating an already active rate sheet changes nothing

Lines 44-51

Asked as: rate sheet already on; activate again does nothing; sheet already active; nothing happens activating

If the rate sheet is already active, Turbo answers "rate sheet is already active" and does nothing more. It skips the revision, completeness and duplicate checks.

## Older revisions cannot be activated

Lines 53-62

Asked as: old revision won't activate; activate older version; newer version exists; can't activate old revision

Turbo refuses to activate a rate sheet when a newer revision with the same name exists. The message is "cannot activate - a newer revision exists".

## Incomplete rate sheet cannot be activated

Lines 64-73

Asked as: rate sheet missing roles; activation fails missing roles; incomplete rates; roles left out of sheet

Activation requires entries for all roles (checked by hooks.ValidateRateSheetComplete). If any role is missing, Turbo refuses with "rate sheet is incomplete - missing entries for some roles" and lists the missing roles.

## Activation blocked by identical active rate sheet

Lines 78-92

Asked as: same rates as another sheet; duplicate rates active; identical rate sheet exists; activation duplicate error

On activation Turbo compares the stored entries with other active rate sheets. Revisions with the same name are excluded. If any sheet has exactly the same rates, activation is refused with a duplicate rate sheet error.

## Activating deactivates other revisions of the same name

Lines 94-115

Asked as: activating turns off other revisions; only one revision active; previous revision switched off; new revision replaces old

Activating a rate sheet first deactivates the other revisions with the same name, then activates this one, in one transaction. So only one revision of a name is active at a time. Then Turbo answers "rate sheet activated".

## Only the job claim can deactivate a rate sheet

Lines 125-132

Asked as: who can deactivate rates; can't deactivate rate sheet; no permission to deactivate; deactivate permission denied

Deactivating a rate sheet requires a person with the job claim (procurement). Anyone else gets "you do not have permission to deactivate rate sheets".

## Deactivating a rate sheet

Lines 140-159

Asked as: turn off rate sheet; deactivate rate sheet; sheet already inactive; switch off rates

Deactivation sets the rate sheet to inactive and answers "rate sheet deactivated". If it is already inactive, Turbo answers "rate sheet is already inactive" and changes nothing. It runs no other check.

## Only admins can edit rate sheet entries

Lines 163-183

Asked as: who can change rates; can't edit rate entry; update rate permission denied; only admins change rates

Updating a rate sheet entry requires the admin claim. Anyone else gets "you do not have permission to update rate sheet entries". Only rate and overtime rate can be changed.

## Entry update needs a rate or overtime rate

Lines 191-200

Asked as: entry update no rate; update needs a rate; overtime rate only update; nothing to update rate

An entry update must provide rate, overtime_rate or both. Otherwise Turbo says "at least one of rate or overtime_rate must be provided". If the entry does not exist, it says "rate sheet entry not found".

## Rate sheet entry update: missing entry is refused

Lines 196-200

Asked as: rate entry missing; can't find rate entry; update entry not found; entry doesn't exist

When an admin updates a rate sheet entry that does not exist, Turbo refuses with "rate sheet entry not found". The admin claim and a rate or overtime rate in the request are checked before the entry is looked up.

## Entry rate must be a whole number of at least 1

Lines 203-212

Asked as: rate has decimals; rate too low; rate fraction not allowed; rate less than one

When updating an entry, a rate below 1 is refused with "rate must be at least 1". A rate with a fraction is refused with "rate must be a whole number".

## Rate sheet entry update rejects rate below 1 with message

Lines 203-218

Asked as: rate below one; rate decimal refused; update rate rejected; refused rate not saved

When an admin sets a rate below 1, Turbo refuses with "rate must be at least 1". When the rate is not a whole number, Turbo refuses with "rate must be a whole number". An overtime rate below 1 gets "overtime_rate must be at least 1". A refused update saves nothing, because the checks run before the save.

## Entry overtime rate minimum is 1 but may have decimals

Lines 213-230

Asked as: overtime rate too low; overtime decimals allowed; overtime under one; overtime minimum value

When updating an entry, an overtime rate below 1 is refused with "overtime_rate must be at least 1". Unlike the rate, the overtime rate may have decimals. Updating one field leaves the other unchanged.

## Entry update returns the stored rates

Lines 220-230

Asked as: rate update result; what rates were saved; confirm updated rates; rate update reply

After a successful update, Turbo replies "rate sheet entry updated" with the entry's saved rate and overtime rate. Only the fields sent in the request change.

## Omitted revision on create counts as revision 0

Lines 241-247

Asked as: revision left blank; no revision number; create sheet without revision; revision empty on create

A create-rate-sheet request has the fields name, effective_date, revision and entries. If the revision is missing, it is 0. Turbo then treats the request as a new sheet or a copy, and it needs the job claim.

## Only the first error for each field is reported

Lines 262-269

Asked as: multiple errors same field; only one error shown; repeated field error; why one error only

When a rate sheet request fails validation on one field more than once, Turbo keeps only the first error for that field. For example, a blank role reports only that it is required. Later errors on the same field are dropped.

## Name and revision conflict reports a unique error on name

Lines 271-275

Asked as: name and revision taken; sheet name already used; duplicate name revision; name conflict error

When a rate sheet with the same name and revision already exists, the error is attached to the name field. Its message is "name and revision must be unique".

## Unique-type validation errors count as a conflict

Lines 277-290

Asked as: unique error on create; create fails as duplicate; not unique validation; conflict error creating sheet

When creating a rate sheet fails validation, Turbo treats it as a name and revision conflict if any field error has the code validation_not_unique or its text contains "unique". Otherwise it reports an ordinary validation error.

## Entry errors are labelled with the entry position

Lines 305-311

Asked as: which entry has error; entry error position; error row number; entries index error

When an entry fails validation while a rate sheet is created, its error is labelled "entries.N.field". N is the entry's position in the list, starting at 0. For example, "entries.2.rate" is the rate of the third entry.

## Creating a rate sheet needs job claim for revision 0

Lines 328-337

Asked as: who can create rate sheet; can't create new sheet; create copy permission; create rate sheet denied

A new rate sheet or a copy (revision 0) requires the job claim. Anyone else gets "you do not have permission to create rate sheets".

## Revising a rate sheet needs revise or admin claim

Lines 338-351

Asked as: can't revise rate sheet; revision permission denied; who can make revision; new revision not allowed

A revision (revision above 0) requires the rate_sheet_revise claim or the admin claim. Without either, Turbo says "you do not have permission to revise rate sheets".

## Create validation errors are returned together

Lines 353-405

Asked as: several errors at once; all validation errors; create sheet errors list; empty entries error

Turbo checks the name, effective date, revision and entries together. It returns every failure in one "Validation error" response, and creates nothing when any check fails. An empty entry list gives "at least one entry is required".

## Rate sheet name and effective date are required

Lines 355-370

Asked as: missing sheet name; blank effective date; name required rate sheet; effective date missing

The name is trimmed and must not be empty ("name is required"). The effective date is trimmed, is required ("effective date is required") and must look like YYYY-MM-DD ("effective date must be in YYYY-MM-DD format"). It must also be a real date (utilities.IsValidDate).

## Effective date must be YYYY-MM-DD and a real date

Lines 361-370

Asked as: date format wrong; invalid effective date; date not accepted; wrong date format

When a rate sheet is created, the effective date must match YYYY-MM-DD. Otherwise Turbo says "effective date must be in YYYY-MM-DD format". A date in that format must also pass the date validity check (IsValidDate).

## Revision cannot be negative

Lines 372-375

Asked as: negative revision number; revision minus one; revision below zero; revision number invalid

A rate sheet revision below 0 is refused with "revision must be 0 or greater".

## New rate sheet needs entries with valid roles and rates

Lines 377-405

Asked as: new sheet entries invalid; entry rate invalid create; no entries in sheet; sheet entries validation

A rate sheet needs at least one entry ("at least one entry is required"). Each entry needs a role ("entry role is required"), and a role cannot repeat ("entry role must be unique"). Each rate must be at least 1 ("entry rate must be at least 1") and a whole number ("entry rate must be a whole number"). Each overtime rate must be at least 1 ("entry overtime_rate must be at least 1"), and may have decimals. Turbo returns all errors as "Validation error".

## Each rate sheet entry needs a unique role

Lines 383-392

Asked as: same role twice; duplicate role entries; role repeated in sheet; role listed twice

When a rate sheet is created, each entry needs a role ("entry role is required"). A role that appears twice in the request is refused with "entry role must be unique".

## Create entry rate must be at least 1 and whole

Lines 393-400

Asked as: create rate has decimals; create rate too low; new sheet rate invalid; overtime create minimum

When a rate sheet is created, an entry rate below 1 is refused with "entry rate must be at least 1". A rate that is not a whole number is refused with "entry rate must be a whole number". An overtime rate below 1 is refused with "entry overtime_rate must be at least 1"; the overtime rate may have decimals.

## Revision effective date cannot precede previous revision

Lines 407-418

Asked as: revision date too early; effective date before previous; date earlier than last revision; revision date older

For a revision above 0, the effective date must be on or after the previous revision's date. The same date is allowed. Otherwise Turbo says "effective date must be on or after the previous revision's effective date" and gives the previous date.

## Created rate sheet starts inactive and is saved atomically

Lines 420-537

Asked as: new sheet status; created sheet inactive; partial save entries; sheet saved all or nothing

A new rate sheet is created inactive (active=false) with all its entries in one transaction. If any entry fails, nothing is saved. Turbo answers with the sheet id and the number of entries created. A failing entry is reported as entries.N.field.

## Rate sheet name and revision must be unique

Lines 431-448

Asked as: rate sheet already exists; same name and revision; sheet duplicate exists; already have this revision

Turbo refuses to create a rate sheet when one with the same name and revision exists. The message is "a rate sheet with this name and revision already exists".

## New rate sheet blocked if rates match an active sheet

Lines 450-459

Asked as: new sheet same rates; rates match active sheet; create blocked duplicate rates; same as active rates

Creating a rate sheet is refused when its entries exactly match another active rate sheet. Sheets with the same name are excluded, so a revision that changes only the effective date is allowed.

## Creating a rate sheet fails whole if any entry fails

Lines 480-530

Asked as: one entry failed create; partial sheet not saved; entry failure cancels sheet; create rolled back

Turbo creates the rate sheet and all its entries in one transaction. If an entry fails validation, its error is shown with the entry's position (entries.N.field) and nothing is saved.

## Compare ignores active sheets with the same name

Lines 565-571

Asked as: compare same name sheets; skipped sheets same name; compare ignores own name; comparison leaves out name

A comparison of unsaved rate sheet entries leaves out active sheets that have the name of the sheet being made. Its name is the one in the request.

## Only the job claim can compare rate sheets

Lines 573-585

Asked as: who can compare rates; compare permission denied; can't compare rate sheets; compare access needed

Comparing entries with active rate sheets requires the job claim. Anyone else gets "you do not have permission to compare rate sheets".

## Compare returns exact and near matches without changing anything

Lines 587-600

Asked as: check rates against active; compare rate matches; similar rate sheets; preview duplicate rates

The compare action checks unsaved entries against active rate sheets. Active sheets with the same (trimmed) name are skipped. It returns the exact matches, which would block a save or activation, and a few close sheets as information. It changes nothing.
