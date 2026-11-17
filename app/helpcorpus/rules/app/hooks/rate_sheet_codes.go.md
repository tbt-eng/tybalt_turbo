<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/rate_sheet_codes.go","sha256":"4c98965fe96983075b8f9bf7eb337946e111f8811b6d0a7e23b8ecce7c1a446f","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/rate_sheet_codes.go

## Rate sheet code format and serial

Lines 41-58

Asked as: code looks wrong; rate sheet code format; code won't accept my input; what does the code look like

A rate sheet code (code) has the form YYYY-MM-NNNNN, such as 2026-11-00004. The month must be 01 to 12 and the five-digit serial must be at least 00001. A code in any other form is not valid.

## Month has at most 99999 rate sheet codes

Lines 43-44

Asked as: can't create more codes this month; out of codes for the month; month code limit reached; no codes left this month

When a month has used serial 99999, Turbo cannot issue another code. It fails with the error "no rate sheet code is left in <month>".

## Rate sheet label shows code, revision and name

Lines 65-79

Asked as: how is a rate sheet shown; sheet name in brackets; rev number shown on sheet; how do labels show rate sheets

Turbo names a rate sheet in labels and messages by its code and revision, such as 2026-11-00004 (rev. 2). When the sheet has a name, the name follows in quotes. A blank name is left out.

## Next rate sheet code is issued per month

Lines 81-109

Asked as: next code for new sheet; why did my new code skip a number; which number does the next sheet get; code numbering restarts monthly

Turbo issues a new code for the current UTC month. The serial is one more than the highest serial already issued or already on a rate sheet in that month. The serial starts at 00001 in each month.

## Rate sheet code month uses UTC

Lines 84-85

Asked as: which month is my code in; wrong month on rate sheet code; timezone for rate sheet code; code month doesn't match my date

Turbo takes the month of a new rate sheet code from the current time in UTC. The serial counts up from 00001 in each month.

## Issued codes include sheets saved outside Turbo

Lines 92-100

Asked as: code already used by imported sheet; duplicate code from outside import; sheet made outside app has same code; new code jumps past old sheets

When Turbo issues the next code, it takes the highest serial of that month from both the registry (rate_sheet_codes) and the codes on existing rate sheets. It adds 1 to the higher of the two. So a code on a sheet written outside Turbo's hooks is never issued again.

## Issued rate sheet codes are never reused

Lines 111-138

Asked as: reuse a deleted sheet's code; old code came back; get a code back after delete; why isn't my deleted code free

Turbo records every code it issues in rate_sheet_codes and never clears it. A code is not issued again, even after its sheet is deleted or absorbed.

## Rate sheet names must be unique

Lines 155-184

Asked as: sheet name already in use; two sheets same title; duplicate rate sheet title; name clash with another sheet

A rate sheet name cannot equal the name of a sheet with another code. The comparison ignores capitals and surrounding spaces. An empty name is never in conflict, so several sheets may have no name. Revisions of one sheet share its name.

## Message when a rate sheet name is taken

Lines 186-207

Asked as: name already taken error; project request has the same name; can't use this sheet title; name belongs to a proposal

When another rate sheet has the name, Turbo refuses with "another rate sheet already has this name". When the name belongs to a project request's proposed rates, the message is "a project request's proposed rates already have this name".

## Rate sheet code error message and fields

Lines 209-217

Asked as: red error under code box; invalid code error message; where do code errors appear; code field shows an error

When Turbo refuses a rate sheet code, it shows the error under the code field (invalid_rate_sheet_code) with the message of the specific refusal.

## Undo of a rate sheet absorb restores the original code

Lines 219-231

Asked as: bring back absorbed sheet; undo merge keeps old code; restore sheet after absorb undo; undo absorb code rejected

A pending rate sheet absorb (absorb_actions on rate_sheets) records the code of each absorbed sheet. When undo recreates that sheet, Turbo accepts the code it carries. This holds even when no sheet with that code exists any more, and when the code's other revisions exist. Turbo then skips the draft and name checks for that sheet.

## Only a superuser reaches the collection API update

Lines 233-243

Asked as: edit sheet through API as non-admin; who can update rate sheets via API; API update refused for sheet; only admins can change sheet data

The collection API has no update rule for rate sheets, so only a superuser can update a sheet through it. Turbo marks that save as a superuser code edit, and the marker cannot come from a request body.

## New rate sheet gets its code automatically

Lines 251-271

Asked as: how do I get a code for sheet; code missing on new sheet; sheet saved without a code; assign code when saving

Nobody chooses a code. A new sheet with no code and revision 0 gets the next code when it is saved. Turbo trims spaces from its name, and refuses the name if another sheet holds it.

## New revision needs the code of the sheet it revises

Lines 256-259

Asked as: new revision missing code; revision without original code; add revision to existing sheet; why won't my revision save

A new sheet with no code and a revision other than 0 is refused with "a new revision takes the code of the sheet it revises".

## New rate sheet cannot carry its own code

Lines 273-294

Asked as: I typed a code for new sheet; can't enter my own code; new sheet code rejected; preset code on new sheet

A new sheet that carries a code in the wrong form is refused with "a rate sheet code has the form YYYY-MM-00001". A code in the right form that no sheet has is refused with "a new rate sheet gets its code when it is saved". The one exception is the undo of a rate sheet absorb, which recreates the sheet with its code.

## Cannot revise a draft rate sheet

Lines 295-299

Asked as: revise a draft sheet; draft blocks new revision; can't make revision of draft; draft error on code field

A new revision is refused when any revision of that sheet is a draft. Turbo returns the draft error on the code field. This check is skipped when the undo of an absorb recreates the sheet.

## A revision keeps its sheet's name

Lines 300-309

Asked as: revision name changed error; rename revision not allowed; why can't I rename this revision; revision name must match sheet

A new revision must have the same name as the latest revision of its sheet. Otherwise Turbo refuses with "a revision keeps the name of its rate sheet; rename the sheet instead".

## Only a superuser can change a rate sheet code

Lines 320-323

Asked as: change an existing sheet code; edit code of saved sheet; code is locked after saving; fix wrong code on sheet

Changing the code of a saved rate sheet is refused for everyone but a PocketBase superuser who saves through the collection API. The message is "only a superuser can change a rate sheet code".

## Rules for a superuser's code change

Lines 324-341

Asked as: superuser changes code rules; code used before refused; change code of sheet with revisions; admin code change restrictions

A new code must be in the form YYYY-MM-00001, or Turbo refuses it. A code that was issued before is refused with "this code was issued before; a code is never used twice". A sheet that has other revisions is refused with "the code of a rate sheet with revisions can't change".

## Renaming a saved rate sheet

Lines 345-355

Asked as: rename saved sheet; change sheet title after saving; saved sheet name won't change; edit name of existing sheet

When a saved sheet's name changes, Turbo trims it and refuses it if another sheet holds the name. The same refusal messages as for a new sheet apply.

## Code changes are logged and the new code is reserved

Lines 359-375

Asked as: log of code changes; record of old and new code; where are code changes kept; new code reserved after change

After a superuser changes a code, Turbo records the new code as issued. It also adds an entry in rate_sheet_code_changes with the sheet, the old code and the new code.

## Hooks apply on every save path

Lines 377-409

Asked as: checks skipped on API save; do name checks apply everywhere; save from another route ignores checks; transaction checks on sheet save

Code and name checks run inside a transaction on every save of a rate sheet, whether it comes from a route, the collection API or a superuser's edit.

## Collection API update of rate sheets is for superusers only

Lines 410-418

Asked as: API update fails for normal users; sheet edit via collection API; only admin hits update hook; non-admin editing sheet over API

The collection API has no update rule for rate sheets, so only a superuser reaches the update request hook. Turbo then marks the save as a superuser code edit.

## Hook errors on rate sheets reach the collection API as field errors

Lines 421-444

Asked as: generic failed to save message; API save error hides reason; field errors on sheet save; why no clear error from API

When a rate sheet save through the collection API fails with one of Turbo's hook errors, Turbo returns that error's own status, message and per-field codes. It does not show PocketBase's generic "Failed to ..." message. Other errors pass through unchanged.
