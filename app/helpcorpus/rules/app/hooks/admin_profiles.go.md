<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/admin_profiles.go","sha256":"f0e7aae320fd2cc4e64eee54bc4e6373e7f0c282c61d9a53ac10f0e9645e9b76","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/admin_profiles.go

## Only admins can edit admin profiles directly

Lines 15-27

Asked as: can't edit admin profile; permission denied admin profile; limited editor cannot change profile; why can't I change admin settings

Turbo refuses a direct create or update of an admin profile by anyone without the admin claim. The message is "you do not have permission to edit admin profiles directly". Limited editors change profiles through Turbo's custom routes instead.

## Admin profile requires a user

Lines 31-44

Asked as: admin profile won't save; uid missing error; no user selected profile; profile blank user field

Turbo refuses to save an admin profile with an empty uid (the user). The error is "uid is required" (code required). A uid of only spaces counts as empty.

## Opening balances need an opening date

Lines 46-60

Asked as: opening balance error; opening date missing; balance won't save without date; opening date needed

Turbo refuses to save an admin profile when opening_op or opening_ov is not zero and opening_date is empty. The error is "opening_date is required when opening balances are non-zero".

## Opening date must be Sunday after pay period end

Lines 62-75

Asked as: opening date rejected; wrong opening date; date not sunday; pay period date error

When an admin profile has an opening_date, Turbo checks it with ValidateTimeOffOpeningDate. If the check fails, Turbo refuses the save with "opening_date must be the Sunday after a pay period ending date".

## Default branch must be usable by the person

Lines 77-77

Asked as: default branch not allowed; branch can't be used; wrong default branch; profile branch error

When Turbo saves an admin profile, it checks the default_branch against the profile's user with EnsureUserCanUseBranch. The rule for what that check allows is in another file.

## Deactivating a profile signs the person out

Lines 80-106

Asked as: user logged out after deactivating; deactivated user kicked out; sessions ended when profile changed; changing uid signs user out

When an admin profile goes from active to inactive, Turbo gives the user a new token key after the save. Every token the person holds becomes invalid at once. The same happens if an active profile's uid changes, and the old user is signed out. This covers every save, including the editor and the active route.

## Deactivation save can leave profile inactive if sign-out fails

Lines 91-106

Asked as: deactivate error but still inactive; save failed profile inactive; profile stuck inactive; sign out failed on save

The token key changes after the profile is written. If the key change fails outside a transaction, as in the collection API, the save returns an error but the profile stays inactive. Inside a transaction, both writes roll back together.

## Creating an inactive profile signs the person out

Lines 108-121

Asked as: new inactive profile logs user out; impersonation session ended; creating inactive profile logout; superuser session lost

When an admin profile is created as inactive, Turbo gives the user a new token key after the save. This invalidates tokens obtained another way, such as a superuser's impersonation. A profile created as active changes nothing.

## Deleting an admin profile signs the person out

Lines 123-132

Asked as: deleted profile user logged out; delete admin profile sign out; removing profile kicks user out; who can delete admin profile

When an admin profile is deleted, Turbo gives its user a new token key, because a user with no profile counts as inactive. Only a superuser can delete a profile.

## Missing user is not an error when signing out

Lines 134-149

Asked as: sign out error missing user; no user to log out; logout fails user not found; token revoke no user error

When Turbo revokes tokens and the uid is empty or the user does not exist, it does nothing and raises no error. Other errors are returned.
