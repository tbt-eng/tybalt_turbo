<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/rate_sheets_rename.go","sha256":"d8c7120e124188122a57b614f45fd5c300f031171ec0c5bb6b17ccce2dab9547","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/rate_sheets_rename.go

## An empty name removes the rate sheet name

Lines 29-31

Asked as: clear the sheet name; remove name from rate sheet; blank rate sheet name; delete the name of a sheet

The name is an optional description. A rename to an empty name is allowed and removes the name.

## Rename applies to every revision of the sheet

Lines 33-35

Asked as: rename all revisions; rename changes every version; name not updating on other revisions; same name on all revisions

Revisions of a sheet share its code and its name. A rename changes the name on every revision with that code in one transaction. If any save fails, no revision changes.

## Rename records history and changes nothing else

Lines 33-37

Asked as: rename history log; does renaming reprice jobs; rename audit trail; rename keeps rates and codes

Each rename writes one name change record (rate_sheet_name_changes) with the old name, the new name, the person who renamed, and the revision in the URL. The code, the rates and the jobs' rate sheet links do not change, so nothing is repriced.

## Current names must be unique, ignoring case

Lines 39-41

Asked as: name already taken; duplicate rate sheet name; rate sheet name exists; same name different capitals

Turbo refuses a rename when another rate sheet already has the new name. The message is "another rate sheet already has this name". The check ignores case, so two sheets cannot differ only in capitals, but a sheet can change the capitals of its own name. Only current names count, so a sheet may take a name that another sheet used before.

## Who can rename a rate sheet

Lines 46-56

Asked as: cannot rename rate sheet; no permission to rename; rename button missing; need revise claim

Only holders of the admin claim or the rate sheet revise claim (rate_sheet_revise) can rename a rate sheet. Anyone else gets the message "you do not have permission to rename rate sheets".

## Rename removes spaces around the new name

Lines 62-62

Asked as: extra spaces in name; trailing spaces rate sheet; name with spaces trimmed; whitespace only name

Turbo trims spaces from both ends of the new name before it compares or saves it. A name of only spaces becomes an empty name.

## Proposed rates in a request cannot be renamed

Lines 72-75

Asked as: rename blocked draft sheet; proposed rates cannot rename; project request rate sheet name; can't rename proposed sheet

Turbo refuses to rename a rate sheet that is a draft (proposed rates in a project request). The error is the draft error for the field "name" (hooks.RateSheetDraftError). The requester names proposed rates in the request.

## New name must differ from current name

Lines 76-81

Asked as: name same as before; rename no change; just changing capitals rejected; name unchanged error

Turbo refuses a rename when the trimmed new name equals the sheet's current name. The message is "the new name is the same as the current name" (name_unchanged). This comparison is exact, so a change of capitals alone passes.

## Name-in-use check skips the sheet's own code

Lines 83-90

Asked as: name used by another sheet; sheet with that name exists; own name conflict rename

The name-in-use check takes the sheet's code (hooks.RateSheetNameHolder with the code). A different sheet that holds the name blocks the rename.

## Rename saves every revision and the history row together

Lines 92-114

Asked as: partial rename failed; rename half applied; rename rolled back; some revisions not renamed

Turbo renames every revision of the rate sheet and writes the history row (rate_sheet_name_changes) in one transaction. If any save fails, no revision keeps the new name and no history row is written.

## Rename history row records old name, new name and who renamed

Lines 104-113

Asked as: who renamed sheet; rename record user; log of old and new name

Each rename writes one history row (rate_sheet_name_changes). It holds the revision in the URL (rate_sheet), the old name, the new name and the person who renamed (renamed_by, the signed-in user).

## Rename error responses

Lines 115-135

Asked as: rate sheet missing error; rename conflict message; validation error on rename; sheet not found when renaming

If the rate sheet does not exist, Turbo says "rate sheet not found". A unique-name validation error gives a conflict. Other validation errors give "Validation error".

## Rename result shows what changed

Lines 137-149

Asked as: rename response details; how many revisions changed; what did rename return

A successful rename returns the id, the code, the new name, the old name, the number of revisions renamed (revisions_renamed) and the former names.

## Any signed-in user can read former names and codes

Lines 153-181

Asked as: previous names of sheet; old names of rate sheet; see earlier codes; rate sheet name history

Any signed-in user can read a rate sheet's current name, former names and former codes, newest first. An unknown sheet gives "rate sheet not found".

## Former names list leaves out some names

Lines 183-212

Asked as: missing old names; previous name not listed; name history gaps; why is my old name gone

The former names list comes from the old names in the rename history of all revisions with the sheet's code, newest first. It leaves out the current name, empty names and repeats. A sheet renamed A to B and back to A shows only B.

## Former codes list

Lines 214-239

Asked as: old code of sheet; previous code lookup; legacy code find sheet; code from printed document

The former codes list gives the old codes from code changes (rate_sheet_code_changes) of sheets that have this code now, newest first. It leaves out the current code and repeats. Printed documents and the legacy system keep old codes, and this list leads back to the sheet.
