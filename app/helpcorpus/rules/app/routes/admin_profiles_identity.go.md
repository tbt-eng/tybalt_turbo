<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/admin_profiles_identity.go","sha256":"1ae091f6c9c338c1afacdc809152180cac43fe9a5d5ca2f00d458f94d8e7befd","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/admin_profiles_identity.go

## IT claim is the second claim that opens admin profile identity

Lines 16-16

Asked as: it claim access; IT team can log in; IT role admin profiles; why can't IT see screen

Turbo names the IT claim "it" (itClaimName). Turbo checks this claim, together with the admin claim, before it lets someone use the admin profile identity screens. A person with either claim gets in.

## Admin profile identity list contents and order

Lines 53-91

Asked as: admin profile list columns; sort order of admin profiles; what columns in admin list; provider count column

The list shows every admin profile with its active flag, legacy ID, the user's email and name, the profile's given name and surname, and the count of linked external sign-in providers (provider_count). Missing values show as empty text. Rows are sorted by surname, given name, email, then user ID.

## Admin profile identity list is open to admin or IT claim only

Lines 53-57

Asked as: can't see admin profile list; no permission admin list; access denied admin profiles page

Turbo checks for the admin or IT claim before it lists admin profile identities. A person with neither claim is refused with "you do not have permission to manage admin profile identity fields".

## Viewing one admin profile identity requires admin or IT claim

Lines 93-97

Asked as: can't open admin profile identity; blocked from viewing admin profile; permission denied view profile

Turbo checks for the admin or IT claim before it shows one admin profile's identity. A person with neither claim is refused with "you do not have permission to manage admin profile identity fields" (requireAdminProfileIdentityClaim).

## Saving a legacy ID requires admin or IT claim

Lines 109-113

Asked as: can't save legacy id; edit legacy uid denied; cannot change legacy id

Turbo checks for the admin or IT claim before it saves an admin profile's legacy ID (legacy_uid). A person with neither claim is refused with "you do not have permission to manage admin profile identity fields".

## Save with unreadable request body is refused

Lines 115-121

Asked as: invalid request body error; save fails bad request; request body unreadable

When Turbo cannot read the save request body, it refuses the save with the message "invalid request body" (code invalid_request_body).

## Blank legacy ID clears it without duplicate check

Lines 123-124

Asked as: remove legacy id; delete legacy id; empty legacy id; clear legacy uid

Surrounding spaces are trimmed from the legacy ID before saving. A blank value skips the duplicate check and sets the admin profile's legacy ID to empty.

## Saving legacy ID returns updated identity

Lines 125-140

Asked as: legacy id saved response; updated profile after save; validation error saving legacy id

The legacy ID save and the identity reload happen in one transaction. A successful save returns the updated admin profile identity. A save that fails validation is reported as a validation error.

## Legacy ID save and the identity response are one transaction

Lines 126-134

Asked as: partial save rolled back; change not kept if fails; transaction legacy id provider

Turbo saves the legacy ID and builds the identity response in one transaction. The same applies when it clears a provider. If either step fails, the change is not kept.

## Clearing a provider returns the refreshed identity

Lines 143-167

Asked as: provider removed from profile view; delete sign in provider refresh; unlink provider reload

After a provider is deleted, Turbo returns the admin profile identity without that provider. The deletion and reload happen in one transaction.

## Clearing an authorized provider requires admin or IT claim

Lines 143-147

Asked as: can't remove provider; unlink google login denied; delete external login permission

Turbo checks for the admin or IT claim before it clears an authorized provider. A person with neither claim is refused with "you do not have permission to manage admin profile identity fields".

## Only admin or IT claim holders manage admin profile identity

Lines 169-183

Asked as: who can manage admin profiles; admin or IT only; permission admin profile identity

Listing, viewing, saving a legacy ID and clearing an authorized provider on admin profiles all require the admin claim or the IT claim (it). Anyone else gets "you do not have permission to manage admin profile identity fields".

## Legacy ID must be unique across admin profiles

Lines 191-212

Asked as: legacy id already used; duplicate legacy id; legacy uid taken

Saving a non-empty legacy ID (legacy_uid) that another admin profile already uses is refused with "legacy_uid is already assigned to another admin profile". The profile's own current value does not count as a duplicate.

## Duplicate legacy ID error message

Lines 200-211

Asked as: admin profile identity validation error; duplicate legacy id message; legacy uid already assigned

When another admin profile already has the legacy ID, Turbo refuses the save with "admin profile identity validation error". The error on legacy_uid is duplicate_legacy_uid: "legacy_uid is already assigned to another admin profile".

## Clearing an authorized provider requires its ID

Lines 225-241

Asked as: remove provider no id; external auth id missing; provider id required

Clearing an authorized provider without an external auth ID is refused with "authorized provider id is required". Turbo checks that the admin profile exists first.

## Missing authorized provider ID error message

Lines 230-241

Asked as: authorized provider validation error; external auth required error; provider id missing message

When no authorized provider ID is given, Turbo refuses with "authorized provider validation error". The error on external_auth is code required: "authorized provider id is required".

## Provider must belong to this admin profile's user

Lines 243-271

Asked as: provider not linked to profile; remove another user's login; wrong provider for profile

Turbo deletes an authorized provider only if it is linked to the user of this admin profile. Otherwise it refuses with "authorized provider not found for this admin profile".

## Provider not found error message

Lines 254-265

Asked as: authorized provider not found; provider not found on profile; external auth not found

When the provider is not linked to this admin profile's user, Turbo refuses with "authorized provider not found". The error on external_auth is code not_found: "authorized provider not found for this admin profile".

## Clearing a provider deletes only that link

Lines 269-273

Asked as: does user account get deleted; unlink provider keeps user; remove login keeps profile

Turbo deletes the authorized provider (external auth) link and leaves the admin profile and user in place. The checks on the admin profile ID and the provider ID run before the delete.

## Admin profile ID is required and must exist

Lines 276-305

Asked as: profile id missing; no admin profile id; admin profile not found

Viewing, saving or clearing a provider needs an admin profile ID. A blank ID is refused with "admin profile id is required". An ID that matches no admin profile is refused with "admin profile not found".

## Missing admin profile ID error message

Lines 277-288

Asked as: admin profile id required error; id required validation; blank profile id

When the admin profile ID is blank, Turbo refuses with "admin profile identity validation error". The error on id is code required: "admin profile id is required".

## Unknown admin profile error message

Lines 290-302

Asked as: admin profile not found message; profile does not exist; wrong profile id

When no admin profile has the ID, Turbo refuses with "admin profile not found". The error on id has code not_found.

## Admin profile identity detail shows user, profile and providers

Lines 316-362

Asked as: what does profile view show; identity details contents; profile email name providers

Viewing one admin profile returns its legacy ID, the user's email and name, the profile's given name and surname, and each authorized provider with its ID, provider, provider ID, created and updated times. If the profile has no user ID, only the ID, user ID and legacy ID are returned. If no user or profile record is found, those fields stay empty.

## Admin profile with no user shows only ID and legacy ID

Lines 317-326

Asked as: profile without user; admin profile no email; missing uid profile

When an admin profile has no user ID (uid), its detail shows only the ID and legacy ID. Email, name, given name, surname and providers stay empty.

## Missing user or profile record leaves fields blank

Lines 328-359

Asked as: profile name blank; user record missing; empty email on profile

When the user record or the profile record is missing, Turbo does not fail. It leaves email, name and providers, or given name and surname, blank.

## Identity detail lists each authorized provider

Lines 336-349

Asked as: linked providers list; provider created updated times; which logins are linked

For each authorized provider, the detail shows its ID, provider, provider ID, created time and updated time. Without a user record, the providers list is not filled.

## Admin profile identity errors map to messages

Lines 364-373

Asked as: failed to save identity; failed to clear provider; unexpected error admin identity

A validation error on a save or a clear is returned with its own status and its data. An unexpected failure shows the message "failed to save admin profile identity" or "failed to clear authorized provider", or the load message for a view.
