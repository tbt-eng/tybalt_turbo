<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/claim_assignable_users.sql","sha256":"1186672caee5907f0b12fc43355d80f7bc23d02e7348d1c040b237bd5e749b50","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/claim_assignable_users.sql

## Users without a profile still appear, with blank names

Lines 3-12

The assignable list keeps a user who has no profile (LEFT JOIN on profiles). Missing given name, surname, name, username, email and admin profile id show as empty text, not as null.

## Only users with an active admin profile can be given a claim

Lines 9-11

The list of users who can be assigned a claim includes only users who have an active admin profile (admin_profiles.active = 1). Users with no admin profile, or an inactive one, do not appear.

## Users who already hold the claim are left out of the assignable list

Lines 13-18

A user who already holds the chosen claim (a row in user_claims with the same claimId) does not appear in the list of assignable users. The list shows only users who do not hold that claim yet.

## Assignable users are sorted by surname, then given name

Lines 19-22

Turbo sorts the assignable users first by surname. If the surname is blank or missing, it uses the user's name, then email, then username. It then sorts by given name, using the username when the given name is blank. The user id breaks any remaining tie.
