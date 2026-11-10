<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/admin_profiles_augmented.md","sha256":"b59bac4655d726f7497934f5508e94883d4deafadd174a7e18ce894f6ed01422","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/admin_profiles_augmented.md

## Admin profiles view cannot be created, updated or deleted through the API

Lines 5-5

Asked as: edit admin profile via api; cannot update admin profile; delete admin profile api; change profile through view

The admin_profiles_augmented view has no create, update or delete rules in this file. It is a saved query that Turbo reads and never writes, so nobody can change records through it. Admin profiles are changed in the admin_profiles collection.

## Who can list admin profiles with PO approver details

Lines 7-16

Asked as: who can list admin profiles; list admin profiles permission; see all admin profiles; why can't I see admin profiles

Only a signed-in person who holds the admin, hr or time_off_manager claim can list records in the admin_profiles_augmented view. Everyone else, including signed-out visitors, fails the list rule.

## Who can view one admin profile with PO approver details

Lines 18-27

Asked as: view single admin profile; who can open admin profile; admin profile access denied; cannot open a profile

Only a signed-in person who holds the admin, hr or time_off_manager claim can view a single record in the admin_profiles_augmented view. Everyone else fails the view rule.

## View adds name and PO approver limits to each admin profile

Lines 32-80

Asked as: po approver limits in profile; show approver max amount; profile name and approver limits; approver spending limit columns

Each admin profile row also shows the person's given name and surname from their profile. The row also shows the person's PO approver limits (po_approver_max_amount, po_approver_project_max, po_approver_sponsorship_max, po_approver_staff_and_social_max, po_approver_media_and_event_max, po_approver_computer_max) and divisions. These come from the po_approver_props record linked to the person's po_approver claim.

## PO approver divisions default to an empty list

Lines 61-61

Asked as: empty divisions shows brackets; po approver divisions blank; divisions show []; no divisions listed

When a person has no PO approver divisions, the view shows an empty list ('[]') in po_approver_divisions instead of an empty value.

## PO approver details are empty without the po_approver claim

Lines 65-80

Asked as: approver columns are null; missing po approver limits; no approver details shown; profile shows no limits

The PO approver columns are filled only when the person holds the po_approver claim. Without that claim, the columns are empty (null) while the admin profile row still appears. A person with the claim but no po_approver_props record also shows empty limits.

## Work week hours must be between 8 and 40

Lines 90-90

Asked as: weekly hours limit; work week hours too high; hours per week range; work week too few hours

The work_week_hours field is required and must be between 8 and 40.

## Default charge-out rate must be between 50 and 1000

Lines 92-92

Asked as: charge out rate limit; default hourly rate range; rate too low or high; charge rate validation

The default_charge_out_rate field is required and must be between 50 and 1000.

## Skip minimum time check has three allowed values

Lines 94-94

Asked as: minimum time check setting; skip min time options; min time check value; skip time check choices

The skip_min_time_check field is required and must be one of no, on_next_bundle or yes.

## Opening date and vehicle insurance expiry use YYYY-MM-DD format

Lines 95-95

Asked as: date format error; opening date wrong format; insurance expiry date format; date must be year-month-day

The opening_date and personal_vehicle_insurance_expiry fields must have the format YYYY-MM-DD (four digits, two digits, two digits).

## Opening overtime payout and overtime vacation hours have maximums

Lines 96-97

Asked as: overtime payout maximum; opening overtime too large; overtime vacation hours cap; opening op limit

The opening_op field has a maximum of 332. The opening_ov field has a maximum of 200.

## Payroll ID format

Lines 98-98

Asked as: payroll number invalid; payroll id format; employee payroll number rules; cms payroll code

The payroll_id field is required. It must be a whole number that does not start with 0, or the letters CMS followed by 1 or 2 digits (such as CMS7).

## Mobile phone number format

Lines 102-102

Asked as: phone number format; mobile number invalid; cell phone format; how to enter mobile

The mobile_phone field must look like +1 (123) 456-7890. The field is optional.

## Default branch is required

Lines 105-105

Asked as: branch missing on profile; no default office; choose home branch; default office required

The default_branch field is required and must point to one branch.

## Given name and surname rules

Lines 106-107

Asked as: name rejected; first name characters allowed; last name apostrophe hyphen; name too long or short

The given_name is required, 2 to 48 characters, letters only, with hyphens allowed between letter groups. The surname is required, 2 to 48 characters, letters only, with hyphens or apostrophes allowed between letter groups.
