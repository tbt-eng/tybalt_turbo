<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/profiles.md","sha256":"0d59206aaf689bbc429e8102a1ffcdd172b51f9e11382a5af1c820265089d1c1","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/profiles.md

## Who can list profiles

Lines 9-20

Asked as: see someone's profile list; who can list profiles; profiles missing from my list; manager can't see profile

A signed-in person can list a profile if it is their own (uid), or if they are its manager (manager). A holder of the tame claim or the job claim (procurement) can list all profiles. A holder of the legacy_po_create_update claim can list only profiles whose user has an active admin profile (admin_profiles_via_uid.active = true).

## Who can view one profile

Lines 24-26

Asked as: open another person's profile; view colleague profile; can't see profile page; look up staff profile

Any signed-in person can view any single profile record. The view rule checks only that the person is signed in, so it is wider than the list rule.

## Who can create a profile

Lines 30-33

Asked as: create my profile; add new profile for me; make profile for myself; can't add profile for colleague

A signed-in person can create a profile only for themselves. The uid of the new record must be their own user id.

## Who can update a profile

Lines 37-41

Asked as: edit my profile; change my profile details; manager can't edit profile; update someone else's profile

A signed-in person can update only their own profile (uid is their id). A request must not change uid. Managers and admins cannot update another person's profile through the API.

## Nobody can delete a profile through the API

Lines 45-45

Asked as: delete my profile; remove profile; can't delete profile; get rid of old profile

The profiles collection has no delete rule, so nobody can delete a profile through the API. Only a PocketBase superuser can delete one.

## Name format limits on a profile

Lines 52-53

Asked as: name rejected by profile; hyphen in surname; apostrophe in name; name too short or long

The surname and the given name are required and must have 2 to 48 characters. They allow letters only (a-z, A-Z). A surname can also have hyphens or apostrophes between letter groups (for example O'Neil). A given name allows hyphens but no apostrophes. Spaces and digits are refused.

## Manager is required; alternate manager is optional

Lines 54-55

Asked as: set profile manager; no manager on profile; optional alternate approver; backup manager field

A profile must have one manager (manager, a user). The alternate manager (alternate_manager) is optional and is limited to one user.

## Allowed values for notification type and default expense payment type

Lines 58-63

Asked as: notification type options; default payment type dropdown; expense payment type choices; email format for notifications

The notification type (notification_type) is either email_text or email_html. The default expense payment type (default_expense_payment_type) is one of OnAccount, Expense, CorporateCreditCard, Allowance, FuelCard, Mileage or PersonalReimbursement.

## Each user has at most one profile

Lines 70-70

Asked as: duplicate profile error; two profiles same person; already has a profile; profile already exists for user

A unique index on uid means that two profiles cannot belong to the same user.
