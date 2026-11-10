<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/profiles.go","sha256":"8892fc9feef12da5a6fcb1450e0c9d3524d4563785ca4032b0160302be1f9e60","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/profiles.go

## Profile default division must be an active division

Lines 12-18

When a profile is created or updated, Turbo checks the default division (default_division) with EnsureActiveDivision. If the check fails, Turbo rejects the save with the error from that check. The exact rule and message live in EnsureActiveDivision, which is not in this file.

## Profile manager must be an active user

Lines 23-34

When a profile has a manager, Turbo checks that the manager is an active user. If not, it refuses the save with "manager must be an active user" (manager_not_active). A profile with no manager skips this check.

## Profile manager must hold the tapr claim

Lines 36-45

When a profile has a manager, the manager must have the tapr claim. Otherwise Turbo refuses the save with "manager must have tapr claim" (invalid_manager). The active check runs first, so an inactive manager gets the not-active message.

## Profile alternate manager must be an active user

Lines 47-58

When a profile has an alternate manager (alternate_manager), that person must be an active user. Otherwise Turbo refuses the save with "alternate manager must be an active user" (alternate_manager_not_active). The alternate manager is optional, and an empty value skips the checks.

## Profile alternate manager must hold the tapr claim

Lines 60-69

When a profile has an alternate manager, that person must have the tapr claim. Otherwise Turbo refuses the save with "alternate manager must have tapr claim" (invalid_alternate_manager). This check runs after the active check.
