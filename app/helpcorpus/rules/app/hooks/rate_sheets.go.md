<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/rate_sheets.go","sha256":"ed880ae299c00db705a0f0ee7e057f9b4cdd51d80933252c47b1dff2d8098cd8","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/rate_sheets.go

## Completeness check runs only when activating

Lines 19-31

Turbo checks that a rate sheet has entries for all roles only when the sheet is being saved as active. A rate sheet saved as inactive is not checked. A sheet that is already active is not checked again on update.

## Rate sheet completeness check failure message

Lines 33-42

If Turbo cannot check the rate sheet's entries because of a database error, it refuses the activation. The message is "unable to verify rate sheet completeness".

## Activating a rate sheet requires entries for every role

Lines 44-56

Turbo refuses to set a rate sheet active when any role in rate roles (rate_roles) has no entry on that rate sheet (rate_sheet_entries). The message is "rate sheet must have entries for all roles before activation". The error lists the missing role names.

## Missing roles are listed alphabetically

Lines 63-84

The list of roles missing from a rate sheet is sorted by role name. A role counts as missing when the sheet has no entry for it.

## Revised rate sheet effective date cannot be earlier than previous revision

Lines 89-119

When a person creates a rate sheet with a revision above 0, Turbo compares its effective date (effective_date) with that of the previous revision of the same name. If the new date is earlier, Turbo refuses it. The message is "effective date must be on or after <previous date> (previous revision's date)". The same date is allowed.

## Effective date check skips first revision and edits

Lines 90-99

Turbo checks the effective date only when a rate sheet is created, and not for revision 0. Edits to an existing rate sheet are not checked.

## Missing previous revision allows any effective date

Lines 144-147

If Turbo finds no previous revision (revision minus 1, same name), it accepts the new effective date without a check.

## Effective dates are compared as text

Lines 149-154

Turbo compares effective dates as YYYY-MM-DD text. A new date equal to the previous date is valid, and only a strictly earlier date is refused.

## Newer revision check blocks activating older revisions

Lines 157-174

Turbo can tell whether a rate sheet with the same name and a higher revision number exists. It is used to prevent activating an older revision. Only a strictly higher revision counts.

## Activating a revision deactivates other revisions

Lines 176-190

When a revision is activated, Turbo sets active to false on all other rate sheets with the same name that are active. This leaves one active revision at a time.
