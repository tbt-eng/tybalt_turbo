<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/rate_sheets_rename.go","sha256":"66273b71d23b05e0171cbf2729786e8a8ff7f38baba50a38997534eb76ea51cf","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/rate_sheets_rename.go

## Rate sheet name is optional; empty name removes it

Lines 30-32

Asked as: remove rate sheet name; clear the name; blank rate sheet name; optional name empty

Turbo trims spaces from the new name before it uses it. The name is an optional description, so renaming to an empty name removes the name.

## Rename applies to every revision of the sheet

Lines 34-38

Asked as: rename all revisions; revisions renamed count; rename changes every version; does renaming reprice jobs

Revisions of a sheet share its code and name. A rename changes the name on all revisions with that code in one transaction. The response gives the number as revisions_renamed. The code does not change, no rates change, and jobs keep their rate sheet, so nothing is repriced.

## Name already in use by another sheet is refused

Lines 40-42

Asked as: name already taken; duplicate rate sheet name; name used by another sheet; rename refused name exists

Turbo refuses a rename when another rate sheet (with a different code) already has the new name. The message names the kind of holder, and it comes from RateSheetNameHeld. The check ignores case, so two sheets cannot differ only in capitals. A sheet can change the capitals of its own name.

## Only current names block a rename

Lines 40-41

Asked as: name used before; old name taken by other sheet; reuse a former name

Turbo checks only the current names of other sheets. A sheet can take a name that another sheet used before.

## Who can rename a rate sheet

Lines 47-57

Asked as: who can rename rate sheet; need admin to rename; rename permission rate sheet; cannot change rate sheet name

Only staff with the admin claim or the rate sheet revise claim (rate_sheet_revise) can rename a rate sheet. Anyone else gets the message "you do not have permission to rename rate sheets".

## Rename refusal message for missing permission

Lines 55-57

Asked as: no permission to rename; rename button not allowed; missing rate sheet revise claim; cannot rename, access denied

A person who holds neither the admin claim nor the rate sheet revise claim (rate_sheet_revise) gets the message "you do not have permission to rename rate sheets".

## Rename trims spaces around the new rate sheet name

Lines 63-63

Asked as: extra spaces in name; whitespace rate sheet name; leading trailing spaces; spaces only name

When a person renames a rate sheet, Turbo trims leading and trailing spaces from the new name before it compares or saves it. A name of only spaces counts as an empty name.

## Rename compares the new name to the current name after trimming

Lines 63-83

Asked as: only capitals changed; change case of name; fix capital letters in name; rename same name different case

Turbo compares the trimmed new name with the sheet's current name (name). The comparison is exact, so a change of capitals alone counts as a change and is allowed.

## Rename is all-or-nothing

Lines 69-117

Asked as: rename half done; rename failed partway; name not saved everywhere; rollback rename rate sheet

Turbo renames every revision and writes the name history row in one transaction. If any step fails, no revision keeps the new name and no history row is written.

## Renaming an unknown rate sheet is refused

Lines 70-73

Asked as: cannot find rate sheet; wrong rate sheet id; missing rate sheet rename; sheet does not exist

When no rate sheet has the id in the rename request, Turbo refuses with the message "rate sheet not found".

## Draft rate sheets cannot be renamed

Lines 74-77

Asked as: cannot rename draft; rename proposed rate sheet; draft rates in request; project request sheet rename

Turbo refuses to rename a rate sheet that is a draft (proposed rates in a project request). It gives the draft error from RateSheetDraftError for "name". The requester names proposed rates in the request.

## New name must differ from current name

Lines 78-83

Asked as: name not changing; same name error; nothing changed on rename; rename to same name

If the trimmed new name equals the sheet's current name, Turbo refuses with the message "the new name is the same as the current name" (name_unchanged). The check is on the name after trimming and is case-sensitive, so a change of capitals alone is allowed.

## Rename checks the name against the code's other sheets

Lines 85-93

Asked as: how name clash is checked; name holder lookup; check other sheets for name

Before it renames, Turbo looks for a current holder of the new name, passing the sheet's code (hooks.RateSheetNameHolder). Which sheets count as another holder depends on that function in another file.

## Rename is recorded in name history

Lines 107-116

Asked as: rename audit log; who renamed rate sheet; name change history; track rate sheet renames

Each rename writes one history row (rate_sheet_name_changes) with the old name, the new name and the person who renamed it. The row is anchored to the revision in the URL.

## Rename error responses

Lines 118-139

Asked as: rename error message; another sheet has this name; failed to rename rate sheet; validation error rename

A rename of an unknown rate sheet gives "rate sheet not found". A name that clashes on a unique rule gives "another rate sheet already has this name". Other validation problems give "Validation error". Anything else gives "failed to rename rate sheet".

## Rename response lists former names

Lines 141-153

Asked as: rename response fields; old name returned after rename; list of previous names after rename

A successful rename returns the id, code, new name, old name, number of revisions renamed and the former names of the sheet.

## Any signed-in user can read former names and codes

Lines 157-184

Asked as: view old names of sheet; see previous codes; rate sheet history lookup; former names endpoint

The former-names endpoint returns a rate sheet's current name and code, its earlier names and its earlier codes. Any signed-in user can read it. An unknown sheet gives "rate sheet not found".

## Former names list rules

Lines 187-215

Asked as: previous names list wrong; old names repeated; why name missing from history; former names order

Former names come from the old names in the rename history of all revisions with the sheet's code, newest first. The list leaves out the current name, empty names and repeats. A sheet renamed A to B and back to A shows only B.

## Former codes list rules

Lines 218-242

Asked as: old rate sheet code; previous code of sheet; legacy code lookup; printed document old code

Former codes are the old codes from the code-change history of the sheet, newest first. The list leaves out the current code and repeats. Printed documents and the legacy system keep old codes, so this list leads back to the sheet.
