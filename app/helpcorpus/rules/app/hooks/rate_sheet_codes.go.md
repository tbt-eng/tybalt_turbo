<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/rate_sheet_codes.go","sha256":"1e02ac651a556c165aab0a3d1631826a35461709c3a661518fc9686c610f19a1","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/rate_sheet_codes.go

## Rate sheet code format and serial

Lines 41-58

Asked as: rate sheet code format; what does a rate sheet code look like; invalid rate sheet code number; serial 00000 code

A rate sheet code looks like 2026-11-00004: the UTC creation month of its first revision (YYYY-MM) and a five-digit serial. A serial of 00000 is not a valid code, and the month must be 01 to 12.

## Rate sheet code limit per month

Lines 44-103

Asked as: no rate sheet code left; too many rate sheet codes this month; run out of codes in month; monthly rate sheet code limit

A month has at most 99999 rate sheet codes. When the serial would pass 99999, Turbo refuses with 'no rate sheet code is left in <month>'.

## Rate sheet label format

Lines 65-74

Asked as: how rate sheet is labelled; rate sheet name shown in quotes; rate sheet label with revision; sheet name display format

Turbo names a rate sheet in messages and labels as its code and revision, such as 2026-11-00004 (rev. 2). When the sheet has a name (after trimming spaces), the name follows in quotes.

## How Turbo issues the next rate sheet code

Lines 84-109

Asked as: next rate sheet number; which code will my new sheet get; code numbering restarts monthly; how new rate sheet code is picked

Turbo issues the next code of the current UTC month. The serial is the highest serial already issued or already used on a sheet in that month, plus 1. Serials restart at 00001 each month, and a code is recorded in rate_sheet_codes when it is issued.

## A rate sheet code is never reused

Lines 111-138

Asked as: reuse deleted rate sheet code; old rate sheet number again; code of deleted sheet came back; rate sheet code recycled

Turbo keeps every issued code in rate_sheet_codes and never clears it. A code is not issued again, even after its sheet is deleted or absorbed.

## Rate sheet names must be unique

Lines 155-197

Asked as: rate sheet name already taken; duplicate rate sheet name; same name as another sheet; two sheets same title

Two rate sheets of different codes cannot have the same name. The comparison ignores capitals and surrounding spaces. Turbo refuses with 'another rate sheet already has this name'. An empty name is never in use, so several sheets can have no name. Revisions of one sheet share its name.

## Undo of a rate sheet absorb restores the code

Lines 209-284

Asked as: undo absorb rate sheet; restore deleted rate sheet; get back absorbed sheet old code; undo merge of rate sheet

When the undo of a pending rate sheet absorb recreates a deleted sheet, Turbo accepts the sheet with its old code. The draft and name checks do not run for it.

## Only a superuser can change a rate sheet code

Lines 223-233

Asked as: cannot edit rate sheet code; change code of rate sheet; who can change sheet number; rate sheet code locked

A change of the code of a saved rate sheet needs a PocketBase superuser saving through the collection API. Anyone else gets 'only a superuser can change a rate sheet code'.

## New rate sheet without a code gets one

Lines 246-261

Asked as: create rate sheet without code; new rate sheet assigned number; revision above zero new sheet; new sheet name trimmed

A new rate sheet with no code must be revision 0 and gets the next code automatically. If it has a revision above 0, Turbo refuses with 'a new revision takes the code of the sheet it revises'. Turbo trims the name and refuses a name that another sheet holds.

## Nobody can choose a code for a new sheet

Lines 263-280

Asked as: pick my own rate sheet code; choose code for new sheet; wrong rate sheet code format; new sheet with existing code

A new sheet that carries a code must be a new revision of an existing sheet with that code. If no sheet has the code, Turbo refuses with 'a new rate sheet gets its code when it is saved'. A code that is not in the form YYYY-MM-00001 is refused with 'a rate sheet code has the form YYYY-MM-00001'.

## Cannot revise a draft rate sheet

Lines 285-289

Asked as: revise a draft rate sheet; new revision of draft; draft sheet cannot be revised; draft error on code field

When a sheet with a code is saved as a new revision, Turbo refuses if any revision of that sheet is a draft. It gives the draft error from RateSheetDraftError on the code field.

## A revision keeps its sheet's name

Lines 290-299

Asked as: revision different name; rename revision of rate sheet; revision name must match sheet; new revision name mismatch

A new revision must have the same name as the latest revision of its sheet. Otherwise Turbo refuses with 'a revision keeps the name of its rate sheet; rename the sheet instead'.

## Rules for a superuser's code change

Lines 314-334

Asked as: change rate sheet number to old one; reuse code after superuser change; sheet with revisions code change; superuser code change rules

The new code must be in the form YYYY-MM-00001. It must never have been issued, or Turbo refuses with 'this code was issued before; a code is never used twice'. The sheet must have no other revisions, or Turbo refuses with 'the code of a rate sheet with revisions can't change'.

## Renaming a saved rate sheet

Lines 335-345

Asked as: rename rate sheet; change rate sheet name; new name already used by sheet; edit sheet title

When the name of a saved sheet changes, Turbo trims it and refuses it if another sheet has it, with 'another rate sheet already has this name'.

## Code changes are logged and the new code is issued

Lines 349-365

Asked as: log of code change; record of rate sheet number change; where code change is logged; rate_sheet_code_changes history

After a superuser changes a code, Turbo records the new code as issued so it is never used again. It logs the sheet, old code and new code in rate_sheet_code_changes.

## Code and name checks run on every save path

Lines 367-409

Asked as: checks on every save; save rate sheet via api bypass; collection api rate sheet save; superuser edit check bypass

Turbo checks a new sheet's code and name and a changed code or name on every path: the routes, the collection API and a superuser's edit. Each check and save runs in one transaction. Only a superuser reaches the collection API update.

## Collection API errors keep their message and fields

Lines 411-434

Asked as: generic error on rate sheet save; api error message unclear; which field rate sheet error; rate sheet save error details

When a rate sheet save through the collection API fails a code or name check, the response carries the check's own message and the field it applies to, not a generic failure message.
