<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/rate_sheets.go","sha256":"03853d0134437a808ba20d61d260e0434c1143ad1ef5d706164d2e1f3938e825","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/rate_sheets.go

## Rate sheet checks run only when it becomes active

Lines 17-22

Asked as: save rate sheet inactive; rate sheet check skipped; completeness check not running; rate sheet saved as draft

Turbo checks a rate sheet for completeness only when the sheet is saved with active set to true. A save that leaves the sheet inactive skips the check.

## Already-active rate sheet is not re-checked on update

Lines 24-31

Asked as: edit active rate sheet; update live rate sheet; no check on active sheet edit; changed active rate sheet

When a rate sheet that is already active is updated, Turbo does not check again that it has an entry for every role. The check runs on a new sheet created as active, or on an inactive sheet that becomes active.

## Activating a rate sheet needs an entry for every role

Lines 33-56

Asked as: cannot activate rate sheet; rate sheet missing role; activate rate sheet fails; add entries for all roles

Turbo refuses to set a rate sheet to active (active=true) when any role (rate_roles) has no entry (rate_sheet_entries) on that sheet. It shows the message "rate sheet must have entries for all roles before activation" and lists the missing roles (missing_roles). The overall error is "rate sheet is missing entries for some roles".

## Rate sheet completeness check failure message

Lines 33-42

Asked as: unable to verify rate sheet completeness; failed to validate rate sheet entries; rate sheet check error; cannot validate rate sheet

If Turbo cannot check the rate sheet's entries, it refuses the activation. It says 'failed to validate rate sheet entries' and, on the active field, 'unable to verify rate sheet completeness'.

## Incomplete rate sheet lists the missing roles

Lines 44-84

Asked as: which roles are missing; rate sheet missing roles list; roles without rate entry; rate sheet incomplete activation

When a rate sheet lacks an entry for some role, Turbo refuses to activate it. It says 'rate sheet is missing entries for some roles', and the active field says 'rate sheet must have entries for all roles before activation'. The error lists the missing roles (missing_roles) by name, in alphabetical order.

## Revised rate sheet effective date cannot precede previous revision

Lines 89-119

Asked as: effective date too early; effective date before previous revision; revision date earlier than last; rate sheet date rejected

When a person creates a rate sheet revision above 0, its effective date must be on or after the previous revision's effective date. An earlier date is refused with 'effective date must be on or after the previous revision's effective date'. The field message names the previous date. The same date is allowed.

## Effective date check skips updates and first revisions

Lines 90-99

Asked as: why date not checked on edit; first revision date allowed; effective date not validated on update; revision zero effective date

Turbo checks the effective date only when a rate sheet is created, not when it is updated. It also skips revision 0, the first revision.

## Previous revision is found by shared code

Lines 128-156

Asked as: no previous revision found; same code revision missing; effective date with no earlier revision; revision lookup by code

Turbo compares the new effective date with the revision numbered one lower that has the same code. If it finds no such revision, it allows the new date without error.

## Newer revision blocks activating an older one

Lines 158-175

Asked as: cannot activate older revision; newer revision exists; outdated rate sheet activation; old revision blocked

Turbo treats a rate sheet as outdated when another sheet with the same code has a higher revision number. This is used to prevent activating older revisions. The function only reports the fact. The refusal itself is not in this file.

## Activating a revision deactivates the other revisions

Lines 177-191

Asked as: other revisions deactivated; only one active revision; activating turns off old sheet; previous rate sheet switched off

When a revision is activated, Turbo can set every other active rate sheet with the same code to inactive. This keeps one revision active at a time.
