<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/rate_sheets_rename.go","sha256":"ea343c4333b12bc4fd0d0733dd0b42b7f70c5f470bc9e133eb35113a26e58bf5","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/rate_sheets_rename.go

## Only admins can rename a rate sheet

Lines 29-50

Asked as: non-admin rename rate sheet; who can rename a rate sheet; rename permission denied; only admins rename rate sheet

Only an admin (the admin claim) can rename a rate sheet. Anyone else gets "you do not have permission to rename rate sheets". Revising a rate sheet does not allow a rename.

## Rename changes every revision at once

Lines 32-35

Asked as: rename rate sheet all revisions; rename affects every revision; rename failed nothing saved; rename rate sheet repricing jobs

Revisions of a sheet are linked by name. A rename sets the new name on every revision that has the old name, in one transaction. If any save fails, none of the changes are kept. No rates change, and jobs keep their rate sheet, so nothing is repriced.

## Another sheet's current name blocks the rename

Lines 37-39

Asked as: rate sheet name already taken; duplicate rate sheet name; another sheet has that name; rename clashes with existing sheet

Turbo refuses the rename with "another rate sheet already has this name" when a sheet with a different name already has the new name. The check ignores case, so two sheets cannot differ only in capitals. A sheet can still change the capitals of its own name. Only current names count, so a sheet can take a name that another sheet used before.

## Rename needs a non-blank name

Lines 56-61

Asked as: blank rate sheet name; rate sheet name is required; empty name rate sheet; rename with spaces only

Turbo trims spaces from the new name of a rate sheet. If nothing is left, it refuses with a validation error on name: "name is required".

## Draft rate sheets cannot be renamed

Lines 70-73

Asked as: cannot rename draft rate sheet; rename draft sheet; draft rate sheet proposed rates; rate sheet still in draft

Turbo refuses to rename a rate sheet that is a draft (hooks.IsRateSheetDraft). It returns the draft error that hooks.RateSheetDraftError gives for the name field. The requester names proposed rates in the project request.

## New name must differ from the current name

Lines 74-79

Asked as: new name same as old; rename to same name; rate sheet name unchanged; rename identical name

If the trimmed new name equals the sheet's current name exactly, Turbo refuses with "the new name is the same as the current name". This comparison is case-sensitive, so a change of capitals only is allowed.

## Name check ignores case but exempts own revisions

Lines 81-90

Asked as: case only rename rate sheet; rename different capitals; rate sheet name differs by case; revisions counted as duplicate name

The in-use check counts rate sheets whose lowercase name equals the lowercase new name and whose name differs from the old name. Revisions of the same sheet are therefore not counted as a clash.

## Rename writes a name history row

Lines 104-113

Asked as: rate sheet rename history; who renamed rate sheet; audit log rate sheet name; track rate sheet name changes

Each rename adds one history row (rate_sheet_name_changes). It holds the revision in the URL as the anchor, the old name, the new name and the admin who renamed it (renamed_by).

## Rename response contents and errors

Lines 115-148

Asked as: rename response details; rename reply revisions count; rate sheet not found rename; validation_not_unique conflict rename

A successful rename returns the id, the new name, the old name, the number of revisions renamed (revisions_renamed) and the former names. If the rate sheet does not exist, Turbo says "rate sheet not found". A name conflict is reported as a conflict error with the code validation_not_unique.

## Any signed-in user can read former names

Lines 152-172

Asked as: old names of rate sheet; previous names rate sheet; who can see former names; look up earlier rate sheet names

The former-names lookup for a rate sheet returns its id, its current name and its earlier names. Any signed-in user can use it. If the sheet does not exist, Turbo says "rate sheet not found".

## Former names: newest first, no repeats

Lines 174-202

Asked as: former names list order; why old name missing; renamed back to original name; repeated names rate sheet history

Former names come from the old names in the history rows whose anchor revision has the current name. They are listed newest first. The current name and repeated names are left out, so a sheet renamed A to B and back to A shows only B.
