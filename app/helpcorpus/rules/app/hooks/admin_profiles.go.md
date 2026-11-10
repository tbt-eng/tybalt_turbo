<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/admin_profiles.go","sha256":"22f8fb10585772d2e5012a7f291677885c97d8b82df3a20ee400df38e52f0afc","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/admin_profiles.go

## Only admins can create admin profiles directly

Lines 13-25

Asked as: cannot create admin profile; admin profile permission denied; edit admin profiles directly; no permission admin profile

Turbo refuses a direct create of an admin profile (admin_profiles collection) unless the person holds the admin claim. A person without it sees "you do not have permission to edit admin profiles directly". Limited editors update admin profiles through the custom route layer instead.

## Admin profile requires a user ID

Lines 30-42

Asked as: uid missing admin profile; admin profile user blank; save admin profile no user; uid is required

Turbo refuses to save an admin profile when the user ID (uid) is empty or only spaces. The error is "admin profile validation error", with the field message "uid is required" (code "required").

## Non-zero opening balances need an opening date

Lines 44-58

Asked as: opening balance no date; opening_date required; blank opening date error; opening balance without date

Turbo refuses to save an admin profile when the opening balance for OP (opening_op) or OV (opening_ov) is not zero and the opening date (opening_date) is empty. The error is "admin profile validation error", with the field message "opening_date is required when opening balances are non-zero". A blank opening date is allowed when both balances are zero.

## Opening date must be Sunday after pay period end

Lines 60-73

Asked as: opening date not sunday; pay period opening date wrong; opening date validation failed; sunday after pay period

When an admin profile has an opening date (opening_date), Turbo checks it with the time-off opening date validation (ValidateTimeOffOpeningDate). If the check fails, Turbo refuses with "admin profile validation error" and the field message "opening_date must be the Sunday after a pay period ending date". The check runs whenever a date is given, even if both balances are zero.

## Default branch must be usable by the profile's user

Lines 75-75

Asked as: default branch not allowed; user cannot use branch; admin profile branch error; wrong default branch saving

When an admin profile is saved, Turbo checks that the user (uid) can use the chosen default branch (default_branch), using the branch check EnsureUserCanUseBranch. This check runs after the uid and opening date checks. The rules for what counts as usable are in that function, not in this file.
