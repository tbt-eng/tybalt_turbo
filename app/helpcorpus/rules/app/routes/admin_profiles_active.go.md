<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/admin_profiles_active.go","sha256":"c30b228868a56363c4b1ef958421f89e5d0ea2e3c4a0286dffab7698ede92a6f","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/admin_profiles_active.go

## Toggle targets list excludes admin users

Lines 36-61

The list of admin profiles whose active state can be toggled contains every admin profile whose user does not hold the admin claim. The list is sorted by profile id.

## Record id is required to set active state

Lines 71-77

Setting an admin profile's active state fails with "record id is required" when the id in the path is empty or only spaces.

## Active value is required

Lines 79-91

Setting an admin profile's active state fails with "active is required" when the request leaves out the active value. An invalid request body fails with "invalid request body".

## Active change runs as one transaction

Lines 93-101

The active-state change runs in a single transaction, so a failure leaves the profile unchanged.

## Setting active state saves without validation

Lines 113-117

For a permitted target, Turbo sets the profile's active flag to the requested value and saves it without running record validation. The response returns the profile id, uid and the new active value.

## Who can list or change admin profile active state

Lines 121-133

Only staff holding the admin, it or hr claim may list the toggle targets or set an admin profile's active state. Anyone else gets the message "you do not have permission to set admin profile active state".

## Admin profile must exist

Lines 136-148

Turbo refuses to set the active state and says "admin profile not found" when no admin profile has the given id.

## Admin profile needs a user

Lines 150-162

Turbo refuses to set the active state when the admin profile has no uid (blank or spaces only). The message is "uid is required".

## Admin users cannot be activated or deactivated here

Lines 164-179

Turbo refuses to change the active state of a profile whose user holds the admin claim, even for an admin. The message is "admin users cannot be activated or deactivated through this endpoint".

## Setting active flag on a non-admin profile

Lines 181-186

When the checks pass, Turbo sets the admin profile's active field (active) to the requested true or false value. It saves the profile without running validation. Turbo then returns the profile id, the user id (uid) and the new active value.
