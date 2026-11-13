<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/rate_sheets_api.go","sha256":"984a724da91dd305a8dfb4c0e45394cfc1ba301133950a8ed8a6c8965871f640","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/rate_sheets_api.go

## Only job claim holders can activate a rate sheet

Lines 29-36

Asked as: cannot activate rate sheet; no permission to activate; only procurement can activate; who can turn on rates

Activating a rate sheet requires the job claim (procurement). Anyone else gets "you do not have permission to activate rate sheets".

## Proposed rates of a project request cannot be activated directly

Lines 44-47

Asked as: activate proposed rates; draft rate sheet activate; project request rates activate; proposed rates not activating

Activating a rate sheet that is a draft (the proposed rates of a project request) is refused with the error from hooks.RateSheetDraftError("active"). The proposed rates are activated when the job is created.

## Activating an already active rate sheet changes nothing

Lines 49-56

Asked as: already active rate sheet; activate again; rate sheet already on; nothing happens activating

If the rate sheet is already active, Turbo answers "rate sheet is already active" and changes nothing.

## Older revision cannot be activated when a newer one exists

Lines 58-67

Asked as: older revision activate; newer revision exists; cannot activate old revision; previous revision won't activate

Turbo refuses to activate a rate sheet when a newer revision with the same name exists. The message is "cannot activate - a newer revision exists".

## Incomplete rate sheet cannot be activated

Lines 69-78

Asked as: rate sheet missing roles; incomplete rate sheet; activate fails missing entries; role missing from rates

Activation needs an entry for every role. Otherwise Turbo refuses with "rate sheet is incomplete - missing entries for some roles" and lists the missing roles.

## Activation runs the duplicate check and deactivates other revisions atomically

Lines 80-98

Asked as: activation failed nothing changed; other revisions deactivated; duplicate check activation; rate sheet activated message

Activation runs in one transaction. It checks for a duplicate active sheet, deactivates the other revisions and activates this sheet through hooks.ActivateRateSheet. If the check fails, Turbo returns that error and activates nothing. On success it answers "rate sheet activated".

## Only job claim holders can deactivate a rate sheet

Lines 108-115

Asked as: cannot deactivate rate sheet; no permission to deactivate; turn off rate sheet permission; only procurement can deactivate

Deactivating a rate sheet requires the job claim (procurement). Anyone else gets "you do not have permission to deactivate rate sheets".

## Deactivating a rate sheet

Lines 123-142

Asked as: turn off rate sheet; switch off rates; rate sheet already inactive; deactivate rate sheet button

Deactivation sets the rate sheet inactive and answers "rate sheet deactivated". If it is already inactive, Turbo answers "rate sheet is already inactive" and changes nothing. Deactivation does not check for drafts or revisions.

## Only admins can update a rate sheet entry

Lines 159-166

Asked as: cannot edit rate entry; no permission editing rates; only admin change rate; edit rate sheet entry denied

Updating a rate sheet entry requires the admin claim. Anyone else gets "you do not have permission to update rate sheet entries". Only rate and overtime_rate can be changed.

## Entry update needs rate or overtime rate

Lines 174-183

Asked as: entry update empty; rate entry not found; nothing to save in entry; which field to update rates

An entry update must include rate or overtime_rate. Otherwise Turbo says "at least one of rate or overtime_rate must be provided". A missing entry gives "rate sheet entry not found".

## Entry rate must be a whole number of at least 1

Lines 186-195

Asked as: rate not whole number; rate decimal not allowed; rate too low; rate has to be integer

When rate is given, it must be at least 1 ("rate must be at least 1") and a whole number ("rate must be a whole number").

## Entry overtime rate must be at least 1 but may have decimals

Lines 196-201

Asked as: overtime rate decimal; overtime rate too low; overtime with cents; overtime below one

When overtime_rate is given, it must be at least 1, or Turbo says "overtime_rate must be at least 1". Unlike rate, it need not be a whole number.

## Entry change cannot duplicate another active sheet

Lines 203-238

Asked as: entry matches another sheet; rates same as active sheet; duplicate rates after edit; edit makes rates identical

Changing an entry of an active sheet is refused if its rates become exactly equal to another active sheet's rates. Revisions with the same name are excluded from the comparison. An inactive sheet is not checked here, because activation checks it.

## Entries of proposed-rates drafts cannot be edited directly

Lines 211-214

Asked as: cannot edit proposed rates; edit draft rates; change project request rates; proposed rates locked

Editing an entry of a draft rate sheet (the proposed rates of a project request) is refused with the error from hooks.RateSheetDraftError("entry"). The requester changes proposed rates through the request.

## Entry update returns the saved rates

Lines 239-252

Asked as: rate entry saved response; what did entry update return; entry updated confirmation; saved rates after edit

When an admin updates a rate sheet entry, Turbo answers "rate sheet entry updated". The answer includes the entry's id, its rate and its overtime rate as saved. A failure from the duplicate or draft checks comes back with its own status and message.

## Only the first problem on a field is reported

Lines 284-291

Asked as: only first error shown; more errors on field; fixed one error another appears; multiple errors same field

When a rate sheet request has several problems on one field, Turbo reports only the first one for that field. The person sees the next problem only after fixing the first and sending the request again.

## Duplicate name and revision gives a conflict message

Lines 293-312

Asked as: name already taken; duplicate rate sheet name; same name and revision conflict; revision conflict message

A rate sheet error about uniqueness is reported as a conflict, not as a plain validation error. The message reads "name and revision must be unique" on the name field. The conflict covers the case where the same name and revision already exist.

## Entry errors are named by entry position

Lines 327-333

Asked as: which entry has the error; entry number error; entries position error; error row number rate sheet

When Turbo rejects an entry of a new rate sheet, it names the field as "entries.N.field", such as "entries.2.rate". N is the entry's position in the request, counting from 0.

## Rate sheet entries validation

Lines 335-362

Asked as: rate sheet needs entry; duplicate role in rate sheet; rate sheet no roles; entry minimum rules

A new rate sheet needs at least one entry. Each entry needs a role, and each role may appear once. Rate must be at least 1 and a whole number. Overtime rate must be at least 1.

## Creating a rate sheet needs the job claim for revision 0

Lines 379-388

Asked as: cannot create rate sheet; no permission new rate sheet; create copy rate sheet denied; who can make rate sheets

A new rate sheet or copy (revision 0) requires the job claim. Anyone else gets "you do not have permission to create rate sheets".

## Revising a rate sheet needs revise or admin claim

Lines 389-402

Asked as: cannot make revision; revise rate sheet denied; no permission revision; new revision not allowed

A revision (revision above 0) requires the rate_sheet_revise claim or the admin claim. Anyone else gets "you do not have permission to revise rate sheets".

## Rate sheet creation validates name, date and revision

Lines 404-432

Asked as: rate sheet name required; invalid effective date; validation error rate sheet; negative revision number

Name is required, after trimming spaces. The effective date is required and must be a real date in YYYY-MM-DD format. Revision cannot be negative ("revision must be 0 or greater"). Any failure gives "Validation error" with the fields.

## Revision effective date cannot precede the previous revision

Lines 434-445

Asked as: effective date too early; date before previous revision; revision date earlier than before; effective date conflicts previous

For a revision, the effective date must be on or after the previous revision's effective date. Otherwise Turbo says "effective date must be on or after the previous revision's effective date" and shows that previous date. The same date is allowed.

## New rate sheets start inactive and are created atomically

Lines 447-575

Asked as: new sheet not active; rate sheet created but inactive; partial rate sheet saved; how many entries created

Turbo creates the rate sheet and all its entries in one transaction, so a failure saves nothing. The new sheet is inactive and must be activated afterwards. Turbo answers with the id and the number of entries created.

## Rate sheet name and revision must be unique

Lines 458-475

Asked as: name and revision exists; rate sheet already exists; duplicate rate sheet; same name revision taken

Creating a rate sheet with a name and revision that already exist is refused with "a rate sheet with this name and revision already exists".

## Name used by a project request's proposed rates is blocked

Lines 477-485

Asked as: name reserved by project; name used by draft; project request name blocked; rate sheet name taken by request

A new rate sheet cannot use a name that a project request's proposed rates (a draft) already use. Turbo returns the error from hooks.RateSheetDraftError("name").

## New rate sheet cannot duplicate an active sheet's rates

Lines 487-496

Asked as: rates same as active sheet; new sheet duplicates rates; identical rates existing sheet; copy rates already used

Turbo refuses a new sheet whose rates exactly match another active sheet. Sheets with the same name are excluded, so a revision that changes only the effective date is allowed.

## Compare ignores active sheets with the same name

Lines 602-608

Asked as: compare ignores same name; compare own revisions; name compare skipped; comparison excludes sheet name

The compare request takes a name and a list of entries. Turbo trims the name. Active rate sheets with that name, which include the sheet's own revisions, are not compared with the entries.

## Compare rate sheets needs job, revise or admin claim

Lines 615-630

Asked as: cannot compare rates; no permission compare; compare button denied; compare access rights

Comparing unsaved entries with active rate sheets requires the job, rate_sheet_revise or admin claim. Anyone else gets "you do not have permission to compare rate sheets".

## Compare returns exact and near matches without changing anything

Lines 632-646

Asked as: check rates before saving; preview duplicate rates; near matches rates; compare without saving

The compare action returns the exact matches that would block a save or activation, and a few near matches as information. Active sheets with the given name are not compared. It changes no record.
