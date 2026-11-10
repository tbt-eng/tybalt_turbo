<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/admin_profiles.md","sha256":"cc27e1455b8e074f7d0a79577d7e46d1ac5e057fbc52c362d59270980ffc84be","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/admin_profiles.md

## Nobody can list admin profiles through the API

Lines 7-9

Asked as: list admin profiles; see all admin profiles; cannot list profiles via API; fetch all admin records

The list action on admin profiles (admin_profiles) has no rule. Nobody can list the records through the API, except a PocketBase superuser. Turbo's own server code can still read them.

## Nobody can view one admin profile through the API

Lines 11-13

Asked as: view one admin profile; open a profile record; cannot read admin profile; get single profile api

The view action on an admin profile (admin_profiles) has no rule. Nobody can view a single record through the API, except a PocketBase superuser. Turbo's own server code can still read it.

## Nobody can create an admin profile through the API

Lines 15-17

Asked as: create admin profile api; add new admin profile; cannot make profile record; post admin profile

The create action on admin profiles (admin_profiles) has no rule. Nobody can create a record through the API, except a PocketBase superuser. Turbo's own server code can still create one.

## Only an admin can update an admin profile

Lines 19-27

Asked as: who can edit admin profile; update profile not allowed; change admin profile details; edit profile refused

A signed-in person who holds the admin claim (user claim named 'admin') can update an admin profile through the API. A request must not change the record id (id), the user it belongs to (uid) or the legacy user (legacy_uid). Anyone else is refused.

## Nobody can delete an admin profile through the API

Lines 29-31

Asked as: delete admin profile; remove a profile record; cannot delete admin profile; erase admin profile

The delete action on admin profiles (admin_profiles) has no rule. Nobody can delete a record through the API, except a PocketBase superuser. Turbo's own server code can still delete one.

## Work week hours must be 8 to 40

Lines 39-39

Asked as: work week hours limit; hours per week range; weekly hours too high; work week hours wrong

The work week hours (work_week_hours) on an admin profile are required and must be between 8 and 40.

## Default charge-out rate must be 50 to 1000

Lines 41-41

Asked as: charge out rate limit; default rate range; hourly rate too low; default charge rate invalid

The default charge-out rate (default_charge_out_rate) on an admin profile is required and must be between 50 and 1000.

## Skip minimum time check has three allowed values

Lines 43-43

Asked as: skip minimum time check; skip min time options; minimum time check setting; skip time check values

The skip minimum time check setting (skip_min_time_check) is required. It accepts only 'no', 'on_next_bundle' or 'yes'.

## Opening date and insurance expiry must be YYYY-MM-DD

Lines 44-44

Asked as: date format wrong; opening date format; insurance expiry date invalid; date must be dd mm yyyy

The opening date (opening_date) and the personal vehicle insurance expiry (personal_vehicle_insurance_expiry) must match the pattern four digits, two digits, two digits, separated by hyphens (YYYY-MM-DD). Both are optional.

## Opening balances have maximums

Lines 45-46

Asked as: opening overtime payout max; opening balance limit; opening overtime value too high; opening op cap

The opening overtime payout (opening_op) is at most 332. The opening overtime value (opening_ov) is at most 200.

## Payroll ID format

Lines 47-47

Asked as: payroll number format; payroll id invalid; employee number rules; CMS payroll code

The payroll ID (payroll_id) is required. It is either a whole number that does not start with 0, or 'CMS' followed by 1 or 2 digits.

## Mobile phone format

Lines 51-51

Asked as: phone number format; mobile number invalid; cell phone format; phone wont save

The mobile phone (mobile_phone) is optional. When given, it must look like +1 (123) 456-7890.

## Default branch is required

Lines 57-57

Asked as: home branch required; no default office; profile missing branch; user missing on profile

Every admin profile must have one default branch (default_branch), chosen from the branches collection. A user (uid) is also required, and a profile belongs to one user.

## One admin profile per user, and unique payroll IDs

Lines 65-66

Asked as: duplicate admin profile; user already has profile; payroll id already used; same payroll number twice

Two admin profiles can't have the same user (uid). Two admin profiles can't have the same payroll ID (payroll_id).
