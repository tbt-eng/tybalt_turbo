<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/admin_profiles_active.go","sha256":"c30b228868a56363c4b1ef958421f89e5d0ea2e3c4a0286dffab7698ede92a6f","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/admin_profiles_active.go

## Toggle targets list excludes admin users

Lines 36-61

Asked as: admin users missing from toggle list; why isn't this admin in the list; which admins can be switched off; toggle list hides admins

The list of admin profiles whose active state can be toggled contains every admin profile whose user does not hold the admin claim. The list is sorted by profile id.

## Record id is required to set active state

Lines 71-77

Asked as: record id error when deactivating; blank id on admin profile; can't turn off profile no id

Setting an admin profile's active state fails with "record id is required" when the id in the path is empty or only spaces.

## Active value is required

Lines 79-91

Asked as: active missing error; invalid request body on toggle; toggle request rejected; no active value sent

Setting an admin profile's active state fails with "active is required" when the request leaves out the active value. An invalid request body fails with "invalid request body".

## Active change runs as one transaction

Lines 93-101

Asked as: half saved profile after error; change partly applied; rollback when toggle fails; profile unchanged after failure

The active-state change runs in a single transaction, so a failure leaves the profile unchanged.

## Setting active state saves without validation

Lines 113-117

Asked as: toggle skips validation; save admin profile without checks; profile saved despite validation errors; what does toggle send back

For a permitted target, Turbo sets the profile's active flag to the requested value and saves it without running record validation. The response returns the profile id, uid and the new active value.

## Who can list or change admin profile active state

Lines 121-133

Asked as: no permission to toggle admin profile; can't deactivate admin profile; who may change profile status; access denied on admin profiles

Only staff holding the admin, it or hr claim may list the toggle targets or set an admin profile's active state. Anyone else gets the message "you do not have permission to set admin profile active state".

## Admin profile must exist

Lines 136-148

Asked as: admin profile not found; profile id doesn't exist; wrong id for admin profile; can't find admin profile to toggle

Turbo refuses to set the active state and says "admin profile not found" when no admin profile has the given id.

## Admin profile needs a user

Lines 150-162

Asked as: profile has no user; admin profile missing uid; uid error when toggling; profile not linked to user

Turbo refuses to set the active state when the admin profile has no uid (blank or spaces only). The message is "uid is required".

## Admin users cannot be activated or deactivated here

Lines 164-179

Asked as: can't deactivate an admin; why can't I switch off this admin; admin account locked from toggle; cannot disable administrator here

Turbo refuses to change the active state of a profile whose user holds the admin claim, even for an admin. The message is "admin users cannot be activated or deactivated through this endpoint".

## Setting active flag on a non-admin profile

Lines 181-186

Asked as: turn admin profile on or off; set profile active flag; enable or disable admin profile; activate admin profile

When the checks pass, Turbo sets the admin profile's active field (active) to the requested true or false value. It saves the profile without running validation. Turbo then returns the profile id, the user id (uid) and the new active value.
