<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/rate_sheets_api.go","sha256":"d8d9fa104737edec05b65f3a39e9cdf3df592e7a59772fc97627f507369636b7","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/rate_sheets_api.go

## Only users with the job claim can activate a rate sheet

Lines 28-34

Asked as: cannot activate rate sheet; no permission to activate; activate button refused; why can't I turn on sheet

Activating a rate sheet requires the job claim. Without it, Turbo refuses with "you do not have permission to activate rate sheets".

## Activating an already active rate sheet changes nothing

Lines 43-49

Asked as: rate sheet already active; activate does nothing; already switched on; activating shows already active

If the rate sheet is already active, Turbo answers "rate sheet is already active". It makes no change and skips the revision and completeness checks.

## Older revision cannot be activated when a newer one exists

Lines 51-60

Asked as: old version won't activate; newer version exists; can't activate previous version; older rate sheet blocked

When a rate sheet with the same name has a newer revision, Turbo refuses to activate the older one. The message is "cannot activate - a newer revision exists".

## Incomplete rate sheet cannot be activated

Lines 62-71

Asked as: rate sheet missing roles; activation fails missing roles; sheet has gaps; rate sheet not complete

A rate sheet must have entries for all roles before activation. Otherwise Turbo refuses with "rate sheet is incomplete - missing entries for some roles" and lists the missing roles (missing_roles).

## Activation deactivates other revisions with the same name first

Lines 73-89

Asked as: other versions switched off; previous version deactivated; two active versions at once; activating turns off others

Turbo deactivates the other revisions with the same name before it activates this one. It uses no transaction, so this order keeps two revisions from being active at once. The success message is "rate sheet activated".

## Only users with the job claim can deactivate a rate sheet

Lines 99-106

Asked as: cannot deactivate rate sheet; no permission to deactivate; can't turn off sheet; deactivate refused

Deactivating a rate sheet requires the job claim. Without it, Turbo refuses with "you do not have permission to deactivate rate sheets".

## Deactivating an inactive rate sheet changes nothing

Lines 114-133

Asked as: rate sheet already inactive; deactivate does nothing; already switched off; turning off changes nothing

If the rate sheet is already inactive, Turbo answers "rate sheet is already inactive" and makes no change. Otherwise it sets the sheet inactive and answers "rate sheet deactivated".

## Only admins can update a rate sheet entry

Lines 150-157

Asked as: can't edit rate entry; only admins edit rates; no permission update entry; job claim not enough editing

Updating a rate sheet entry requires the admin claim. Without it, Turbo refuses with "you do not have permission to update rate sheet entries". The job claim does not suffice.

## Entry update needs a rate or an overtime rate

Lines 165-168

Asked as: rate update empty; must provide rate or overtime; update with no values; entry update needs a value

An entry update must contain a rate, an overtime rate or both. Otherwise Turbo refuses with "at least one of rate or overtime_rate must be provided". Only these two fields can change.

## Entry update checks that the entry exists

Lines 170-174

Asked as: rate entry missing; entry not found; can't find rate line; rate sheet line gone

If the rate sheet entry does not exist, Turbo refuses with "rate sheet entry not found".

## Entry update changes only the fields provided

Lines 176-204

Asked as: update overwrote other rate; changing rate clears overtime; partial rate update; only one rate changed

Turbo changes only the rate or overtime rate that the request contains and leaves the other as it was. It then answers "rate sheet entry updated" with the saved values.

## Entry rate must be a whole number of at least 1

Lines 177-186

Asked as: rate decimal not allowed; rate below 1; rate must be whole; rate has cents

When an entry update includes a rate, a rate below 1 is refused with "rate must be at least 1". A rate that is not a whole number is refused with "rate must be a whole number".

## Entry overtime rate must be at least 1 but may have decimals

Lines 187-192

Asked as: overtime rate decimal; overtime under 1; overtime with cents; overtime fraction allowed

When an entry update includes an overtime rate, a value below 1 is refused with "overtime_rate must be at least 1". Unlike the rate, the overtime rate need not be a whole number.

## Only the first error for each field is reported

Lines 236-243

Asked as: multiple errors only one shown; why only one error; hidden errors per field; other field errors missing

When the rate sheet creation checks find several errors for one field, Turbo keeps only the first and ignores the rest. Other fields still get their own errors.

## Duplicate rate sheet name and revision error names the name field

Lines 245-249

Asked as: name already taken; duplicate rate sheet name; same name and revision error; name field error sheet

When a rate sheet with the same name and revision already exists, the conflict error data has a name field error. Its message is "name and revision must be unique" (code validation_not_unique).

## Any unique-style validation error is treated as a name and revision conflict

Lines 251-264

Asked as: unique error message; validation not unique; duplicate detection text; unique constraint wording

When a validation error has the code validation_not_unique, or its text contains "unique" in any letter case, Turbo treats it as a duplicate. Turbo checks this when a rate sheet save fails with validation errors.

## Creating a new rate sheet or copy needs the job claim

Lines 303-311

Asked as: can't create rate sheet; no permission create copy; copy rate sheet denied; new sheet not allowed

A rate sheet with revision 0 (new or copy) requires the job claim. Without it, Turbo refuses with "you do not have permission to create rate sheets".

## Revising a rate sheet needs the revise or admin claim

Lines 312-325

Asked as: can't revise rate sheet; revision permission denied; job claim cannot revise; revise refused no access

A rate sheet with revision above 0 requires the rate_sheet_revise or the admin claim. Without either, Turbo refuses with "you do not have permission to revise rate sheets". The job claim alone does not suffice.

## Validation errors are returned together

Lines 327-379

Asked as: validation error list; several errors at once; entries.0.rate error; all field errors together

Turbo collects all field errors and returns them as "Validation error". Entry errors use keys like entries.0.rate. Only the first error per field is kept. Checks that need the database run only after these pass.

## Rate sheet name is required

Lines 329-333

Asked as: rate sheet name missing; blank name rejected; name empty error; name spaces only

Turbo trims the name of spaces. An empty name fails validation with "name is required".

## Effective date is required and must be a valid YYYY-MM-DD date

Lines 335-344

Asked as: effective date error; invalid date rate sheet; date format YYYY-MM-DD; effective date blank

The effective date is trimmed. An empty date fails with "effective date is required". A date not in YYYY-MM-DD format fails with "effective date must be in YYYY-MM-DD format". A date in the right format must also be a real date (IsValidDate).

## Revision cannot be negative

Lines 346-349

Asked as: negative revision; revision below zero; minus revision number; revision number invalid

A negative revision fails validation with "revision must be 0 or greater".

## Rate sheet needs at least one entry

Lines 351-354

Asked as: no entries rate sheet; empty rate sheet create; rate sheet without lines; entries required creating

A rate sheet creation with no entries fails with "at least one entry is required".

## Each entry needs a unique role

Lines 357-366

Asked as: duplicate role entry; role required entry; same role twice; entry role repeated

Each entry needs a role. A missing role fails with "entry role is required". A role repeated in the same request fails with "entry role must be unique".

## Entry rate must be a whole number of at least 1 on creation

Lines 367-371

Asked as: entry rate decimal create; entry rate under 1; new entry rate cents; creation rate too low

Each entry rate below 1 fails with "entry rate must be at least 1". A rate that is not a whole number fails with "entry rate must be a whole number".

## Entry overtime rate must be at least 1 on creation

Lines 372-374

Asked as: overtime under 1 new sheet; entry overtime too low; overtime minimum creating; overtime cents new entries

Each entry overtime rate below 1 fails with "entry overtime_rate must be at least 1". Decimals are allowed.

## Revision effective date cannot precede the previous revision

Lines 381-392

Asked as: effective date before previous revision; date earlier than last revision; revision date too early; previous date later

For a revision above 0, the effective date must be on or after the previous revision's date. Otherwise Turbo refuses with "effective date must be on or after the previous revision's effective date" and gives the previous date (previous_effective_date). The comparison is made by CheckRevisionEffectiveDate in another file.

## New rate sheets start inactive and are created with all entries at once

Lines 394-497

Asked as: new rate sheet inactive; why is new sheet not active; partial sheet saved; entries all created together

Turbo creates the rate sheet and all its entries in one transaction, so no partial sheet remains if one entry fails. The new sheet is inactive (active=false) and must be activated separately. The response gives the id, name, revision and the count of entries created (entries_created).

## Rate sheet name and revision must be unique

Lines 405-422

Asked as: rate sheet already exists; duplicate name and revision; same sheet already there; name revision conflict

If a rate sheet with the same name and revision exists, Turbo refuses with "a rate sheet with this name and revision already exists". The same message appears when the database reports a unique conflict.

## Entry save errors report the entry position

Lines 451-478

Asked as: entries.N error; which entry failed; entry role duplicate database; entry save failed position

If saving an entry fails validation, Turbo reports "Validation error" with the field prefixed by entries.N. A database unique failure on an entry is reported as "entry role must be unique" for that entry.
