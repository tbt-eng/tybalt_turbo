<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/rate_sheets_api.go","sha256":"db3cc68c3215223fa79faa8c57b15bc6ef54a587d820383d4bcc5d426185ef9c","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/rate_sheets_api.go

## Only procurement can activate a rate sheet

Lines 29-42

Asked as: cannot activate rate sheet; no permission to activate; who can activate sheets; activate button not working

Activating a rate sheet requires a procurement user (the job claim). Anyone else gets "you do not have permission to activate rate sheets". An unknown rate sheet gets "rate sheet not found".

## Activating an already active rate sheet changes nothing

Lines 44-51

Asked as: rate sheet already active; activate does nothing; nothing happens when activating

If the rate sheet is already active (active), Turbo replies "rate sheet is already active" and saves nothing. This is not an error.

## An older revision cannot be activated

Lines 53-62

Asked as: old revision won't activate; activate older version; newer version exists error

Turbo refuses to activate a rate sheet when a newer revision with the same name exists. The message is "cannot activate - a newer revision exists".

## A rate sheet must cover all roles to activate

Lines 64-73

Asked as: rate sheet incomplete; missing roles on sheet; role missing from rate sheet

Turbo checks completeness with hooks.ValidateRateSheetComplete. If any role has no entry, it refuses with "rate sheet is incomplete - missing entries for some roles" and lists the missing roles.

## Activation blocks a duplicate of another active sheet

Lines 78-92

Asked as: rates same as active sheet; activation blocked by duplicate; identical to another active sheet

On activation Turbo compares the stored entries with the other active rate sheets. Revisions with the same name are excluded. If the rates match exactly, activation fails with the duplicate rate sheet error and the sheet stays inactive.

## Activating a sheet deactivates its other revisions

Lines 94-115

Asked as: other revision switched off; previous version deactivated; only one revision active

Activation deactivates the other revisions with the same name and then activates this sheet, in one transaction. So only one revision of a name is active at a time.

## Only procurement can deactivate a rate sheet

Lines 125-138

Asked as: cannot deactivate sheet; turn off rate sheet denied; deactivate permission error

Deactivating a rate sheet requires a procurement user (the job claim). Anyone else gets "you do not have permission to deactivate rate sheets". An unknown sheet gets "rate sheet not found".

## Deactivating an inactive rate sheet changes nothing

Lines 140-159

Asked as: sheet already inactive; deactivate does nothing; already switched off

If the rate sheet is already inactive, Turbo replies "rate sheet is already inactive" and saves nothing. Otherwise it sets active to false and replies "rate sheet deactivated". Deactivation runs no completeness or duplicate check.

## Only admins can edit a rate sheet entry

Lines 176-199

Asked as: cannot edit rate entry; not allowed to change rates; update entry permission denied

Updating a rate sheet entry requires the admin claim. Anyone else gets "you do not have permission to update rate sheet entries". An unknown entry gets "rate sheet entry not found".

## An entry update must change rate or overtime rate

Lines 191-194

Asked as: empty entry update; update entry with nothing; which fields can I change

An entry update must include rate or overtime_rate. If it has neither, Turbo says "at least one of rate or overtime_rate must be provided". Only these two fields can be changed, and a field that is left out stays as it was.

## Entry rate must be a whole number of at least 1

Lines 203-212

Asked as: rate has decimals; rate below one; fraction not allowed in rate

When updating an entry, a rate below 1 is refused with "rate must be at least 1". A rate with a fraction is refused with "rate must be a whole number".

## Entry overtime rate must be at least 1 but may have decimals

Lines 213-218

Asked as: overtime rate decimal; overtime rate below one; overtime fraction allowed

When updating an entry, an overtime_rate below 1 is refused with "overtime_rate must be at least 1". Unlike rate, the overtime rate may have a fraction.

## Editing an active sheet's entry cannot duplicate another active sheet

Lines 220-251

Asked as: edit makes duplicate sheet; changed rates match other sheet; entry change blocked duplicate

If the entry belongs to an active rate sheet, Turbo applies the change and compares the sheet with the other active sheets. Revisions with the same name are excluded. An exact match is refused as a duplicate. Entries of an inactive sheet are not checked, because activation checks them.

## Entry update failure on a duplicate returns the duplicate error

Lines 252-258

Asked as: DuplicateRateSheetError; save fails same rates; duplicate error on save

When an active sheet's entry change would make its rates equal to another active sheet's, Turbo refuses the save and returns the duplicate-sheet error (DuplicateRateSheetError) instead of a generic failure. The same applies to a duplicate found at activation.

## Entry update responds with the saved rates

Lines 260-265

Asked as: what entry update returns; confirmation after rate change; updated rate response

After a successful entry update, Turbo replies "rate sheet entry updated" with the entry's id, its rate and its overtime rate.

## Only the first error for each field is reported

Lines 297-304

Asked as: only one error shown; hidden validation errors; more errors after fixing

When a rate sheet create request breaks more than one check on the same field, Turbo reports only the first error for that field. The other errors for that field stay hidden until the first is fixed (setValidationError).

## Duplicate name and revision error names the name field

Lines 306-310

Asked as: name and revision must be unique; duplicate revision name error; same name and revision

When a new rate sheet has the same name and revision as an existing one, Turbo attaches the error to the name field. The error reads "name and revision must be unique" (validation_not_unique).

## Creating a new sheet or copy needs procurement; a revision needs more

Lines 363-386

Asked as: cannot create rate sheet; cannot revise rate sheet; copy sheet permission denied

A rate sheet with revision 0 (new or copy) needs a procurement user (the job claim), or Turbo says "you do not have permission to create rate sheets". A revision above 0 needs the rate_sheet_revise or admin claim, or Turbo says "you do not have permission to revise rate sheets".

## Rate sheet create fields are validated

Lines 388-415

Asked as: validation error creating sheet; name required error; rate sheet form errors

Name is trimmed and required. The effective date is required and must look like YYYY-MM-DD and be a valid date. Revision cannot be negative ("revision must be 0 or greater"). At least one entry is required. Turbo returns "Validation error" with the field messages.

## Effective date must be a real YYYY-MM-DD date

Lines 390-405

Asked as: effective date invalid; date format error; effective date required

A new rate sheet needs an effective date. A missing date gives "effective date is required". A date not in YYYY-MM-DD form gives "effective date must be in YYYY-MM-DD format". A date in that form is also checked as a real date (IsValidDate). The name is trimmed first, and a blank name gives "name is required".

## Revision cannot be negative

Lines 407-410

Asked as: negative revision number; revision below zero; revision number error

A new rate sheet with a revision below 0 is refused with "revision must be 0 or greater".

## At least one entry is required to create a sheet

Lines 412-415

Asked as: no entries on sheet; empty rate sheet; create sheet without entries

Creating a rate sheet with no entries is refused with "at least one entry is required".

## Each new entry needs a unique role and valid rates

Lines 417-436

Asked as: duplicate role on sheet; role repeated in sheet; entry rate invalid

Each entry needs a role, and a role cannot repeat ("entry role must be unique"). The rate must be at least 1 and a whole number. The overtime rate must be at least 1 and may have decimals.

## A revision's date cannot precede the previous revision's

Lines 442-453

Asked as: revision date too early; date earlier than previous revision; effective date before last revision

For a revision above 0, Turbo checks the effective date with hooks.CheckRevisionEffectiveDate. If it is earlier than the previous revision's date, it refuses with "effective date must be on or after the previous revision's effective date" and gives that date. A date that equals the previous one is allowed.

## New rate sheets start inactive and are created whole

Lines 455-541

Asked as: new sheet not active; why is new sheet inactive; sheet saved partially

Turbo creates the sheet with active set to false, together with all its entries, in one transaction. If any entry fails, nothing is saved. An entry error is reported with the entry's position, such as entries.0.role.

## A rate sheet name and revision must be unique

Lines 466-483

Asked as: rate sheet already exists; same name and revision exists; duplicate rate sheet name

Creating a sheet whose name and revision already exist is refused with "a rate sheet with this name and revision already exists".

## A new sheet cannot duplicate another active sheet's rates

Lines 485-494

Asked as: new sheet same rates as active; create blocked duplicate rates; copy of active sheet rejected

Turbo refuses to create a sheet whose entries match an active sheet exactly. Sheets with the same name are excluded, so a revision that changes only the effective date is allowed.

## Entry errors name the entry position

Lines 515-565

Asked as: which entry has the error; entries.0 error; error on entry row

If one entry fails to save during creation, the error names the entry by its position (entries.N.field). A role that violates the unique index gives "entry role must be unique". Nothing is saved.

## Successful creation reports the entries created

Lines 567-572

Asked as: how many entries created; create confirmation message; entries_created count

A successful create returns the new sheet's id, name, revision and the number of entries created (entries_created).

## Compare rate sheets needs a job, revise or admin claim

Lines 615-628

Asked as: cannot compare rate sheets; compare tool permission denied; compare access denied

The compare tool is open to procurement (the job claim) and to users with the rate_sheet_revise or admin claim. Anyone else gets "you do not have permission to compare rate sheets".

## Compare returns exact and near matches without saving

Lines 630-643

Asked as: compare without saving; check rates before saving; preview duplicate check

Compare takes entries and a name. It returns the exact matches with active sheets, which a save or an activation would block. It also returns a few near matches for information. Active sheets with the same name are excluded, and nothing is changed.

## Compare shows exact matches and a few near matches

Lines 635-643

Asked as: similar rate sheets; near match rates; compare shows similar sheets

Compare returns the exact matches, which would block a save or activation, and the closest near matches as information. The near matches are limited by RateSheetNearMatchMaxRoles and RateSheetNearMatchLimit (defined in another file). The name is trimmed and sheets with that name are left out.
