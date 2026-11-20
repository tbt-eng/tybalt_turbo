<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/rate_sheets_api.go","sha256":"10765eb8defb7906125ae5fcd5d08a1418fc6e9735e1a26bcd300b8ad4ad5b93","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/rate_sheets_api.go

## Only job claim holders can activate a rate sheet

Lines 29-36

Asked as: can't activate rate sheet; no permission turn on rates; only procurement activate; activate rates not allowed

A person without the job claim (procurement) cannot activate a rate sheet. Turbo answers "you do not have permission to activate rate sheets".

## Proposed rates of a project request cannot be activated

Lines 44-47

Asked as: activate project request rates; can't activate proposed rates; request rates draft activate; activate rates from request

A rate sheet that holds the proposed rates of a project request (a draft) cannot be activated. Turbo returns the draft error from hooks.RateSheetDraftError("active"). Its proposed rates are activated when the job is created.

## Activating an already active rate sheet changes nothing

Lines 49-56

Asked as: rate sheet already active; activate twice; sheet says active already; activate active sheet

If the rate sheet is already active, Turbo answers "rate sheet is already active" and makes no change. This is not an error.

## Older revision cannot be activated

Lines 58-67

Asked as: activate old revision; why can't I activate old version; newer revision exists; activate previous rate version

Turbo refuses to activate a rate sheet revision when a newer revision of the same code exists. The message is "cannot activate - a newer revision exists".

## Incomplete rate sheet cannot be activated

Lines 69-78

Asked as: rate sheet missing roles; rate sheet incomplete; roles missing activation; sheet has no entry for role

A rate sheet must have an entry for every role before it is activated. If roles are missing, Turbo refuses with "rate sheet is incomplete - missing entries for some roles" and lists the missing roles (missing_roles).

## Activation checks duplicates and deactivates other revisions

Lines 80-98

Asked as: activation duplicate sheet error; activate sheet already exists; rate sheet activation failed duplicate; activation other revisions deactivated

Activation runs in one transaction through hooks.ActivateRateSheet. That function checks for a duplicate sheet and deactivates the other revisions. If it finds a duplicate, Turbo returns its error. On success Turbo answers "rate sheet activated".

## Only job claim holders can deactivate a rate sheet

Lines 108-115

Asked as: can't deactivate rate sheet; turn off rate sheet no permission; procurement deactivate rates; deactivate not allowed

A person without the job claim (procurement) cannot deactivate a rate sheet. Turbo answers "you do not have permission to deactivate rate sheets".

## Deactivating a rate sheet

Lines 117-143

Asked as: turn off rate sheet; switch off rates; already inactive sheet; deactivated rate sheet

Deactivation sets the rate sheet to inactive and answers "rate sheet deactivated". If the sheet is already inactive, Turbo answers "rate sheet is already inactive" and changes nothing. Unlike activation, deactivation does not check for drafts, newer revisions or completeness.

## Active rate sheet cannot be deleted

Lines 163-166

Asked as: delete active rate sheet; remove active rates; can't delete live sheet; delete sheet still in use active

Turbo refuses to delete an active rate sheet. The message is "Deactivate the rate sheet before you delete it."

## Proposed rates of a project request cannot be deleted directly

Lines 167-171

Asked as: delete project request rates; remove proposed rates sheet; request rates can't delete; delete sheet from request

A rate sheet that holds the proposed rates of a project request cannot be deleted by itself. The message is "This rate sheet holds the proposed rates of a project request. Delete the request to delete it."

## Rate sheet used by a job cannot be deleted

Lines 173-198

Asked as: delete rate sheet job used; job used rate sheet remove; rate sheet job history delete; job once used these rates

Turbo refuses to delete a rate sheet that a job uses now, as its current rate sheet or in any period of its rate sheet history, or has used. The message is "A job uses or has used this rate sheet, so it can't be deleted."

## Rate sheet used by a project request cannot be deleted

Lines 199-201

Asked as: delete rate sheet request uses; project request using rates delete; remove rates a request uses; request still uses sheet

Turbo refuses to delete a rate sheet that a project request uses. The message is "A project request uses this rate sheet, so it can't be deleted."

## Rate sheet that others are based on cannot be deleted

Lines 202-204

Asked as: rate sheet based on delete; copied rates can't delete; sheet started from delete; other sheet depends on rates

Turbo refuses to delete a rate sheet when another sheet is based on it (based_on), such as the rates of a project request that started from it. The message is "The rates of a project request started from this rate sheet, so it can't be deleted."

## Rate sheet with name or code history cannot be deleted

Lines 205-210

Asked as: rate sheet old names delete; former code can't delete; rename history rate sheet; sheet previously renamed delete

Turbo refuses to delete a rate sheet that has former names or former codes. The messages are "This rate sheet has former names, so it can't be deleted." and "This rate sheet has former codes, so it can't be deleted."

## Pending rate sheet absorb blocks all rate sheet deletes

Lines 211-213

Asked as: absorb waiting delete rate sheet; can't delete any rate sheet; pending absorb commit undo; finish absorb before delete

Turbo refuses to delete any rate sheet while a rate sheet absorb waits to be committed or undone. This applies whatever sheets the absorb involves. The message is "A rate sheet absorb is waiting to be committed or undone. Finish it before you delete a rate sheet."

## Deleting a rate sheet removes entries but keeps its code

Lines 218-257

Asked as: code reserved after delete; deleted rate sheet code reuse; delete removes rate entries; code still reserved deleted

The check and the delete run in one transaction. Any revision can be deleted if it passes the checks. Its entries are deleted with it. Its code stays reserved and Turbo never issues it again. Turbo answers "rate sheet deleted".

## Only admins can delete a rate sheet

Lines 227-242

Asked as: admin delete rate sheet; no permission delete rates; delete rate sheet not admin; rate sheet not found delete

Only a person with the admin claim can delete a rate sheet. Anyone else gets "you do not have permission to delete rate sheets". If the sheet does not exist, Turbo answers "Rate sheet not found."

## Only admins can edit rate sheet entries

Lines 268-279

Asked as: edit rate entries no permission; change role rate admin; can't update rates; only admin change rates

Only a person with the admin claim can update a rate sheet entry. Anyone else gets "you do not have permission to update rate sheet entries". Only rate and overtime rate can change.

## Entry update needs rate or overtime rate

Lines 287-296

Asked as: update entry without rate; entry not found rate; overtime rate missing update; rate entry needs value

An entry update must include rate or overtime_rate. Otherwise Turbo answers "at least one of rate or overtime_rate must be provided". If the entry does not exist, Turbo answers "rate sheet entry not found".

## Entry rate limits on update

Lines 299-308

Asked as: rate must be whole number; rate below 1; rate decimal not allowed; rate zero error

A rate must be at least 1, or Turbo answers "rate must be at least 1". A rate must be a whole number, or Turbo answers "rate must be a whole number".

## Entry overtime rate limits on update

Lines 309-317

Asked as: overtime rate decimals; overtime rate too many decimals; overtime rate minimum; overtime rate below 1

An overtime rate must be at least 1, or Turbo answers "overtime_rate must be at least 1". It may have at most two decimal places, or Turbo answers with hooks.OvertimeRateDecimalsMessage.

## Entry edit cannot duplicate another active sheet

Lines 319-353

Asked as: same rates as another sheet; rates match active sheet edit; entry edit duplicate active; change makes sheet identical

If the entry's sheet is active, Turbo compares the sheet's rates after the change with the other active sheets. Revisions of the same code are excluded. If the rates equal another active sheet exactly, Turbo refuses the change. An inactive sheet is not checked here, because activation checks it.

## Entry edit blocked for proposed rates of a request

Lines 327-330

Asked as: can't edit request rate entry; project request rate edit blocked; proposed rate entry locked; change rates through request

Turbo refuses to edit an entry of a sheet that holds the proposed rates of a project request. The requester changes proposed rates through the request. Turbo returns hooks.RateSheetDraftError("entry").

## Rate sheet entry edit failure responses and success result

Lines 354-368

Asked as: saved rate after update; entry change not saved; overtime rate returned after edit; update entry result rate

After an admin updates an entry, Turbo returns the saved rate (whole number) and overtime rate. If a duplicate or draft check refuses the change, Turbo returns that error and does not save the entry.

## Only the first error on a field is reported

Lines 404-411

Asked as: only one error shown; fixing error reveals another; same field error hidden; another error after fix

When a rate sheet save finds several errors on the same field, Turbo reports only the first one for that field. Fixing it can then reveal a later error on the same field. This applies to the effective date, revision, code and entry fields.

## Rate sheet entry validation on create

Lines 455-484

Asked as: rate sheet needs entry; duplicate role on sheet; entry role required; create sheet entry validation

A new rate sheet needs at least one entry ("at least one entry is required"). Each entry needs a role ("entry role is required"), and each role may appear once ("entry role must be unique"). The rate must be a whole number of at least 1, and the overtime rate at least 1 with at most two decimal places. Turbo reports these errors with the entry's index.

## Creating a new rate sheet needs the job claim

Lines 501-510

Asked as: can't create rate sheet; copy rate sheet no permission; new rate sheet access denied; create sheets needs procurement

A new rate sheet or a copy (revision 0) needs the job claim (procurement). Without it Turbo answers "you do not have permission to create rate sheets".

## Revising a rate sheet needs revise or admin claim

Lines 511-524

Asked as: can't revise rate sheet; revision permission denied; make new revision no access; revise sheet not allowed

A revision (revision above 0) needs the rate_sheet_revise claim or the admin claim. Without either Turbo answers "you do not have permission to revise rate sheets".

## Rate sheet effective date rules

Lines 532-541

Asked as: effective date invalid; date format rate sheet; effective date required; rate sheet date wrong format

The effective date is required and must be written YYYY-MM-DD and be a valid date. The messages are "effective date is required" and "effective date must be in YYYY-MM-DD format". Turbo trims spaces first.

## Revision number cannot be negative

Lines 543-546

Asked as: negative revision number; revision below zero; revision number minus

A rate sheet revision must be 0 or greater. Otherwise Turbo answers "revision must be 0 or greater".

## New rate sheet cannot send a code

Lines 554-560

Asked as: code on new rate sheet; set code when creating; new sheet code error; can't type code new sheet

A new rate sheet (revision 0) may not send a code. Turbo answers "a new rate sheet gets its code when it is saved" and assigns the next code on save.

## Revision finds its sheet by code or by name

Lines 561-583

Asked as: which sheet to revise; revise rate sheet by name; no sheet to revise; find sheet for revision

A revision takes the code of the sheet it revises. If no code is sent, Turbo finds the latest revision of the sheet with that name. If there is no such sheet, Turbo answers "there is no rate sheet to revise". A revision with no name keeps its sheet's name.

## Revision with no sheet to revise is refused

Lines 568-572

Asked as: revision matches no sheet; revise code field error; no matching rate sheet revision; revision sheet not found

When a rate sheet revision matches no existing sheet by code or by name, Turbo refuses it with the validation error 'there is no rate sheet to revise' on the code field.

## Proposed rates cannot be revised

Lines 573-577

Asked as: revise project request rates; draft rates revision blocked; can't revise proposed rates; request rates no revisions

A sheet that holds the proposed rates of a project request has no revisions. Turbo returns hooks.RateSheetDraftError("name") when someone tries to revise it. It becomes an ordinary sheet when the job is created.

## Revision effective date cannot precede previous revision

Lines 585-596

Asked as: effective date earlier than previous; revision date before previous; previous effective date error; date older than last revision

A revision's effective date may equal the previous revision's date but not come before it. Otherwise Turbo answers "effective date must be on or after the previous revision's effective date" and gives previous_effective_date.

## New rate sheet starts inactive and saves atomically

Lines 598-733

Asked as: new sheet inactive after save; entries not saved revision; saved but not active; revision saved inactive

A new rate sheet or revision is saved inactive, with all its entries, in one transaction. If any entry fails, nothing is saved. Turbo answers with the id, code, name, revision and entries_created.

## Duplicate revision number refused

Lines 610-621

Asked as: revision number already exists; same revision twice; sheet already has revision

Turbo refuses to create a revision when the sheet already has that revision number. The message is "this rate sheet already has this revision".

## New rate sheet name must be unique

Lines 622-637

Asked as: rate sheet name taken; duplicate sheet name; name already used rates; sheet name exists

A new rate sheet cannot take a name that another sheet holds, in any letter case. Several sheets may have no name. The proposed rates of a project request hold their names too, and the message says so.

## New rate sheet cannot duplicate an active sheet

Lines 639-648

Asked as: same rates as active sheet; identical rates new sheet; save rates already active; new sheet copies active rates

Turbo refuses to save a rate sheet whose rates exactly equal another active sheet's rates. The revisions of the revised sheet are excluded, so a revision that changes only the effective date is allowed.

## Entry validation errors on create name the entry

Lines 682-696

Asked as: entry error which row; entry index error; role repeated error; entry error position

If an entry fails when Turbo saves it, the validation error is shown under entries.N.field, N being the entry's position. A role that repeats gives 'entry role must be unique'.

## Create conflicts return a conflict error

Lines 701-723

Asked as: code and revision exists; conflict saving rate sheet; name conflict save; rate sheet already exists conflict

When a rate sheet with the same code and revision already exists, Turbo refuses the save with 'a rate sheet with this code and revision already exists'. A duplicate revision found beforehand gives 'this rate sheet already has this revision'. A name already held gives a conflict with the name message.

## Compare finds the revised sheet by code, then by name

Lines 759-767

Asked as: compare rate sheet by code; compare by sheet name; compare finds revised sheet; match sheet to compare

A rate sheet compare request can send the code of the sheet being revised. Without a code, Turbo matches the sheet by the request's name, as it did before codes existed. A new sheet sends neither a code nor a name. The active rate sheets of the matched sheet's code are left out of the comparison.

## Compare rate sheets permission

Lines 776-789

Asked as: compare rates no permission; can't compare entries; compare access denied; compare needs job claim

Only a person with the job claim, the rate_sheet_revise claim or the admin claim can compare entries with the active rate sheets. Anyone else gets "you do not have permission to compare rate sheets".

## Compare returns exact and near matches without saving

Lines 791-813

Asked as: compare rates what matches; near matches rate sheet; check rates before save; compare without saving

Compare changes nothing. It returns the exact matches, which a save or activation would block, and a few near matches as information. It leaves out the revisions of the sheet being revised, found by code or by name.

## Compare excludes the revised sheet's own family

Lines 796-804

Asked as: compare shows own revisions; revisions counted as match; compare includes same sheet; new sheet compare all active

When a compare request names a sheet (by code or name) that exists, Turbo leaves that sheet's code out of the comparison, so its own revisions never count as matches. A new sheet sends neither and is compared with all active sheets.

## Revision lookup uses the latest revision of the sheet

Lines 816-836

Asked as: latest revision lookup; which revision used code; no code no name revise; sheet found by name

When Turbo looks for the sheet that a revision or a comparison refers to, it takes the latest revision of the sheet with that code. If no code is sent, it finds the sheet by exact name (the newest revision with that name) and uses that sheet's code. With neither code nor name, or no match, there is no sheet to revise.
