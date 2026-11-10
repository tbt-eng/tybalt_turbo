<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/claims_details.sql","sha256":"780cde7a724c79717a00f22cf29e921333f1b48212a9786aadc381342be952d1","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/claims_details.sql

## Claim holders list leaves out users with no profile

Lines 1-4

Asked as: user missing from claim holders; holder not showing in list; claim holder without profile; who holds claim not listed

The list of people who hold a claim (user_claims for claimId) shows only users who have a profile (profiles). A user who holds the claim but has no profile does not appear in the list.

## Claim holders without an admin profile still appear

Lines 5-7

Asked as: admin profile link empty; claim holder no admin profile; admin profile id blank; broken link to admin profile

A user who holds the claim but has no admin profile (admin_profiles) still appears in the list of claim holders. Their admin profile id (admin_profile_id) is empty, so the link to their admin profile does not work.

## Claim holders list is sorted by surname, then given name

Lines 8-18

Asked as: claim holders sort order; order of people on claim; claim holders alphabetical; names sorted by surname

The list of claim holders is sorted by surname, then by given name. It shows each person's given name and surname, and includes only users who hold the claim whose id is claimId (uc.cid).
