<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/admin_profiles.md","sha256":"f5a8be6e9072e0efd041788dd570ee342ad4647524f5363045443618dd6ba67a","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/admin_profiles.md

## Nobody can list admin profiles through the API

Lines 7-9

Asked as: list all admin profiles; see every admin profile; get admin profiles from the API; pull the admin list

No access rule exists for listing admin profiles (admin_profiles). Nobody can list them through the API, except a PocketBase superuser. Turbo's own server code can still read them.

## Nobody can view one admin profile through the API

Lines 11-13

Asked as: open one admin profile; view a single admin record; show admin profile details; read someone's admin profile via API

No access rule exists for viewing a single admin profile (admin_profiles). Nobody can do this through the API, except a PocketBase superuser. Turbo's own code can still read it.

## Nobody can create admin profiles through the API

Lines 15-17

Asked as: add a new admin profile; make an admin profile through API; create admin record; new admin profile not saving

No access rule exists for creating an admin profile (admin_profiles). Nobody can create one through the API, except a PocketBase superuser. Turbo's own server code can still create them.

## Only admins can update admin profiles, not identity fields

Lines 19-27

Asked as: edit my admin profile; change admin profile details; update admin record refused; can't save admin profile changes; non admin editing profile

A signed-in person with the admin claim (user_claims name 'admin') can update an admin profile. A request must not change the record id (id), the user (uid) or the legacy user ID (legacy_uid). Anyone without the admin claim is refused.

## Nobody can delete admin profiles through the API

Lines 29-31

Asked as: remove an admin profile; delete admin record; get rid of admin profile; admin profile won't delete

No access rule exists for deleting an admin profile (admin_profiles). Nobody can delete one through the API, except a PocketBase superuser. Turbo's own code can still delete them.

## User and default branch are required

Lines 38-38

Asked as: admin profile missing user; no default branch saved; which branch to pick; profile needs user and branch

An admin profile requires a user (uid) and a default branch (default_branch). Each holds one selection.

## Work week hours must be 8 to 40

Lines 39-39

Asked as: work week hours error; weekly hours too high; hours per week limit; work week out of range

The work week hours (work_week_hours) on an admin profile are required and must be between 8 and 40.

## Default charge-out rate must be 50 to 1000

Lines 41-41

Asked as: charge out rate too low; default rate range; rate won't save; default rate limit

The default charge-out rate (default_charge_out_rate) is required and must be between 50 and 1000.

## Skip minimum time check has three allowed values

Lines 43-43

Asked as: skip min time options; minimum time check setting; skip time check choices; skip minimum time value rejected

The skip minimum time check setting (skip_min_time_check) is required. It accepts only "no", "on_next_bundle" or "yes".

## Opening balance limits and date format

Lines 44-46

Asked as: opening balance too high; opening OP maximum; opening OV limit; opening date wrong format

The opening date (opening_date) must look like YYYY-MM-DD when it is set. The opening OP value (opening_op) can be at most 332. The opening OV value (opening_ov) can be at most 200.

## Payroll ID format

Lines 47-47

Asked as: payroll number invalid; payroll id format error; CMS number format; payroll number starting with zero

The payroll ID (payroll_id) is required. It is either a whole number that does not start with 0, or CMS followed by 1 or 2 digits, such as CMS12.

## Mobile phone number format

Lines 51-51

Asked as: phone number rejected; mobile number format; mobile phone won't save; wrong phone format

The mobile phone (mobile_phone) must look like +1 (555) 555-5555 when it is set.

## Vehicle insurance expiry date format

Lines 53-53

Asked as: car insurance date format; vehicle insurance expiry error; insurance expiry won't save; insurance date wrong

The personal vehicle insurance expiry (personal_vehicle_insurance_expiry) must look like YYYY-MM-DD when it is set.

## User source anchor format

Lines 60-60

Asked as: user anchor too long; source anchor format; anchor not hex; user anchor invalid characters

The user source anchor (user_source_anchor) is at most 32 characters. It must be 32 lowercase hexadecimal characters when it is set.

## User and payroll ID must be unique

Lines 66-67

Asked as: duplicate user on profile; same payroll number twice; payroll id already used; user already has profile

Two admin profiles cannot have the same user (uid). Two admin profiles cannot have the same payroll ID (payroll_id).

## User source anchor must be unique when set

Lines 68-68

Asked as: duplicate source anchor; anchor already taken; same anchor on two profiles; anchor used before

Two admin profiles cannot have the same user source anchor (user_source_anchor). Empty anchors are not compared.
