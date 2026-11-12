<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/rate_sheets_rename.go","sha256":"434f1e9355b3353e506e2ac6b357555c8a0152d4948493e7cdce6a90c0c37981","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/rate_sheets_rename.go

## Rename changes no rates and does not reprice jobs

Lines 27-36

Asked as: does renaming change prices; rename reprice jobs; name change alters rates; rename affects job cost

Renaming a rate sheet changes only the name on its revisions. No rates change, and jobs keep their rate sheet id, so Turbo reprices nothing.

## Only admins can rename a rate sheet

Lines 42-48

Asked as: can't rename rate sheet; no permission to rename; non admin rename sheet; who can change sheet name

Renaming a rate sheet requires the admin claim. Anyone else gets "you do not have permission to rename rate sheets". Revising a rate sheet does not allow a rename.

## New rate sheet name is required

Lines 54-59

Asked as: blank rate sheet name; empty name error; name only spaces; rename name missing

Turbo trims spaces from the new name. If nothing is left, it refuses with a validation error on name: "name is required".

## Rename is all or nothing

Lines 63-108

Asked as: rename partly failed; half renamed revisions; rename rolled back; some revisions not renamed

Turbo renames every revision and writes the name history row in one transaction. If any step fails, no revision changes.

## Rename of an unknown rate sheet fails

Lines 64-67

Asked as: rename nonexistent sheet; wrong rate sheet id; rename sheet not existing

If the rate sheet id is not found, Turbo answers "rate sheet not found".

## Rename must change the name

Lines 68-73

Asked as: same name error; rename to same name; name unchanged error; nothing changed rename

If the trimmed new name equals the current name, Turbo refuses with "the new name is the same as the current name". The comparison is case-sensitive, so changing only the capitals is a valid rename.

## Case-only rename of a sheet's own name is allowed

Lines 69-84

Asked as: change capitals of name; fix case of sheet name; uppercase rate sheet name; only letter case rename

The name check ignores capitals, so two sheets cannot differ only in capitals. A sheet can still change the capitals of its own name, because its own revisions do not count as another sheet. A same-name check is exact and case-sensitive, so a case-only change is not 'unchanged'.

## Another sheet's current name blocks a rename; case ignored

Lines 75-84

Asked as: name already used by other sheet; duplicate rate sheet name; name taken rename; sheet name exists elsewhere

Turbo refuses with "another rate sheet already has this name" if a sheet with a different name already uses the new name. The comparison ignores case, so two sheets cannot differ only in capitals. A sheet may change the capitals of its own name.

## Earlier names of other sheets can be reused

Lines 75-84

Asked as: reuse old name of another sheet; name previously used elsewhere; take name another sheet dropped

Turbo checks only current names. A sheet may take a name that another sheet used before.

## Rename changes every revision of the sheet

Lines 86-96

Asked as: revisions still show old name; all versions renamed; rename applies to every version; old name left on revisions

Revisions are linked by name. A rename sets the new name on every revision with the old name, in one transaction. If any step fails, no revision changes. Rates are not changed and jobs keep their rate sheet, so nothing is repriced.

## Rename is recorded in name history

Lines 98-107

Asked as: rename audit log; log of sheet name changes; sheet name change history

Each rename writes one rate sheet name change row. It holds the rate sheet in the URL, the old name, the new name and the admin who renamed it.

## Name history records who renamed and which revision

Lines 102-106

Asked as: who renamed rate sheet; renamed by which admin; record of rename user

Each rename row stores the revision in the URL (rate_sheet), the old name, the new name and the admin who renamed it (renamed_by).

## Rename fails with a conflict message when the name is taken

Lines 109-118

Asked as: name already taken rename; rate sheet name conflict; unique name error saving revision; rename clashes with existing name

When another rate sheet already has the name, Turbo refuses the rename with 'another rate sheet already has this name'. A unique-name error raised while saving a revision gives the same message.

## Other rename validation errors show as Validation error

Lines 115-120

Asked as: validation error rename; rename rejected field message; generic validation failure rename

A validation failure that is not a unique-name error, such as an unchanged name, is returned as 'Validation error' with the field message.

## Missing rate sheet on rename gives not found

Lines 121-123

Asked as: rate sheet not found; id in url missing; rename sheet id not found

When no rate sheet has the id in the URL, Turbo answers 'rate sheet not found'.

## Rename response contents

Lines 127-138

Asked as: what rename returns; revisions renamed count; rename response fields; old and new name returned

After a rename, Turbo returns the id, the new name, the old name, the number of revisions renamed (revisions_renamed) and the list of former names.

## Any signed-in user can read a sheet's former names

Lines 142-161

Asked as: view old names of sheet; previous names visible to everyone; see former sheet names; sheet name history lookup

The former-names lookup returns the sheet's current name and its earlier names, newest first. It does no claim check in this file. An unknown sheet gives "rate sheet not found".

## Unknown sheet on former names lookup gives not found

Lines 147-150

Asked as: former names not found; old names unknown sheet; previous names lookup error

When no rate sheet has the id, the former names request answers 'rate sheet not found'.

## Former names list skips current name and repeats

Lines 164-191

Asked as: old names missing current name; duplicate names in former list; renamed back shows only once; former names repeated entries

Former names come from the old names in the history rows of the sheet's revisions, newest first. The current name and repeated names are left out. A sheet renamed A to B and back to A shows only B.

## Former names are newest first

Lines 171-177

Asked as: order of previous names; former names sorted by date; list of old names order

Former names come from the history rows of revisions that now have the sheet's name. They are ordered newest first, with ties broken by row order.
