<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/admin_profiles_save.go","sha256":"1ca476a6327881951e1919f61268c388b1da055c7c20835426ad28ccbdd19b6e","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/admin_profiles_save.go

## HR can edit only a fixed set of admin profile fields

Lines 22-38

A person with the HR claim (hr) can change only these admin profile fields through the limited save: active, allow_personal_reimbursement, default_branch, default_charge_out_rate, job_title, mobile_phone, off_rotation_permitted, opening_date, opening_op, opening_ov, payroll_id, personal_vehicle_insurance_expiry, salary, skip_min_time_check and time_sheet_expected.

## Time off manager can edit only four admin profile fields

Lines 40-45

A person with the time off manager claim (time_off_manager) can change only opening_date, opening_op, opening_ov and payroll_id through the limited save.

## Only admins can save an admin profile with claims

Lines 49-55

Saving an admin profile together with its claims requires the admin claim. Anyone else gets the message "you do not have permission to save admin profiles with claims".

## Invalid request body is refused on admin profile save

Lines 57-63

Both admin profile saves refuse a request body that cannot be read, with the message "invalid request body".

## Admin profile save is all or nothing

Lines 67-84

Turbo saves the admin profile and its claim changes in one transaction. If any step fails, nothing is saved. Validation failures show "Validation error".

## Saving without a record id creates a profile

Lines 86-90

When the admin save has no record id, Turbo creates a new admin profile and answers as created. With an id, it updates that profile.

## Limited admin profile save needs HR or time off manager claim

Lines 96-106

Turbo refuses a limited admin profile save unless the person has the HR claim (hr) or the time off manager claim (time_off_manager). The message is "you do not have permission to save limited admin profile fields".

## Limited admin profile save requires a record id

Lines 108-114

Turbo refuses a limited admin profile save with no record id. The message is "record id is required".

## Limited save needs at least one field change

Lines 123-131

Turbo refuses a limited admin profile save that has no field changes (missing or empty admin_profile). The message is "at least one admin profile field change is required".

## Limited save refuses fields outside the person's claims

Lines 133-135

If a limited save includes any field the person's claims do not allow, Turbo refuses the whole save. The message is "you do not have permission to edit one or more admin profile fields". A person with both claims may edit the combined set of fields.

## Limited admin profile save is all or nothing

Lines 138-155

Turbo saves a limited admin profile change in one transaction. When the save fails, it stores none of the changes. A failed save returns the hook error or a "Validation error".

## Limited save returns only the fields the person may edit

Lines 157-160

After a limited save, Turbo returns the record id and only the fields that the person's claims allow them to edit, not the full profile.

## Admin profile to update must exist

Lines 180-192

Turbo refuses to update an admin profile whose id is not found, with the message "admin profile not found". The limited save does the same.

## Admin profile uid cannot be changed

Lines 193-225

When Turbo updates an existing admin profile, it refuses a uid that differs from the stored one. The message is "uid cannot be changed".

## Admin profile requires a uid

Lines 201-213

Turbo refuses an admin profile save when uid is empty after trimming spaces. The message is "uid is required".

## Claim list is trimmed and de-duplicated

Lines 227-227

Turbo trims each claim id, drops empty ones and removes duplicates before it syncs the claims.

## Saving a profile replaces the user's claims

Lines 228-230

The admin save makes the user's claims match the submitted list. It deletes claims not in the list and adds missing ones, and it keeps claims already held. An empty list removes all of the user's claims.

## Claim sync runs before profile validation

Lines 228-244

The admin save changes the user's claims before it runs the profile processing and the save. Both are in one transaction, so a later failure undoes the claim changes.

## Admin profile rules run before saving

Lines 232-244

Both admin profile saves run the admin profile processing (hooks.ProcessAdminProfile) on the record before they save it. Any error from it stops the save.

## Admin profile save returns validation errors

Lines 239-244

When saving the admin profile record or a user claim fails validation, Turbo returns "Validation error" with the field errors.

## Limited save changes only the submitted fields

Lines 255-288

A limited save loads the existing profile, sets only the submitted fields and leaves claims untouched.

## Limited save refuses an unknown admin profile

Lines 255-267

When a person saves limited fields and no admin profile has the given record id, Turbo refuses with "admin profile not found".

## Limited save runs admin profile rules before saving

Lines 273-285

Turbo runs the admin profile processing rules (ProcessAdminProfile) on the changed record before it saves. When these rules refuse, the profile stays unchanged.

## Limited save checks every submitted field

Lines 298-308

Turbo refuses the whole limited save with "you do not have permission to edit one or more admin profile fields" when any submitted field is outside the person's allowed set. It checks this before it saves anything.

## Limited fields shown depend on the person's claims

Lines 310-323

The limited editable field set is empty without the HR claim or the time off manager claim. An HR person (the hr claim) gets the HR fields. A time off manager (the time_off_manager claim) gets the four time off fields. A person with both claims gets the union of the two sets.

## Claim list ignores blank and repeated ids

Lines 325-335

Turbo trims each claim id. It drops blank ids and repeated ids and keeps the first of each.

## Claims not in the new list are removed

Lines 337-392

When Turbo saves a profile with claims, it deletes each current claim of the user that is not in the new claim list (user_claims). It adds the claims that the user does not have yet. An empty list removes all claims of the user.
