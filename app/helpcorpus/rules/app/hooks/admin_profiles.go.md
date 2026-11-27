<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/admin_profiles.go","sha256":"ee463856999e791caa8bfee2ce3aa34a969fa4802d98b813ed0a2b62c076d4ec","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/admin_profiles.go

## Only admins can create or edit admin profiles directly

Lines 15-27

Asked as: cannot edit admin profile; permission denied admin profile; limited editor admin profiles; who can change admin profiles

Only a person with the admin claim can create or edit an admin profile directly through the collection. Anyone else gets the message "you do not have permission to edit admin profiles directly". Limited editors update admin profiles through the custom route layer instead.

## Admin profile requires a user (uid)

Lines 32-44

Asked as: admin profile missing user; uid is required; blank uid error; admin profile validation error

Turbo refuses to save an admin profile when its user (uid) is empty or only spaces. The message is "admin profile validation error", and the uid field says "uid is required".

## Opening balances need an opening date

Lines 46-60

Asked as: opening balance needs date; opening_date required; save fails opening balance; overtime balance without date

Turbo refuses to save an admin profile when the opening overtime-payout balance (opening_op) or the opening overtime-vacation balance (opening_ov) is not zero and the opening date (opening_date) is empty. The opening_date field says "opening_date is required when opening balances are non-zero".

## Opening date must be Sunday after a pay period end

Lines 62-75

Asked as: opening date not sunday; pay period end date error; invalid opening date; sunday after pay period

When an admin profile has an opening date (opening_date), Turbo checks it with ValidateTimeOffOpeningDate. If the check fails, Turbo refuses the save with "opening_date must be the Sunday after a pay period ending date". The check runs even when both opening balances are zero.

## Default branch must be usable by the user

Lines 77-77

Asked as: default branch not allowed; user cannot use branch; default branch error; branch access denied profile

After the other checks pass, Turbo checks that the profile's user may use the default branch (default_branch), through EnsureUserCanUseBranch. The result of that check comes from another file.

## Deactivating an admin profile signs the person out

Lines 88-106

Asked as: deactivated user still logged in; admin profile inactive logout; disable admin signs out; tokens revoked deactivation

When an admin profile changes from active to inactive, Turbo gives the user a new token key. This makes every token the person holds invalid at once. It applies to every save of the profile, including the active route, the editor and saves made in code. Turbo does nothing when the profile was already inactive or stays active. It also does nothing when the user record is not found.
