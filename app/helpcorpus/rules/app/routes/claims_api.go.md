<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/claims_api.go","sha256":"565ef6d8e71ff2803fbc2656cc870783bde6647ebc605f0fc0d871c172404ddd","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/claims_api.go

## Only admins can use claim detail and assignment actions

Lines 69-78

Turbo checks for the admin claim before it shows claim details, shows assignable users or bulk assigns a claim. A person without the admin claim gets "you do not have permission to " plus the action: "view claim details", "view claim assignment details" or "bulk assign claims".

## Viewing claim details needs an admin and an existing claim

Lines 80-109

An admin can view a claim's name, description and holders. Turbo shows "claim not found" when the claim id does not exist. The holders come from the claim details query (claimDetailsQuery).

## Viewing assignable users for a claim

Lines 111-137

An admin can list the users who can be assigned a claim. The list comes from the assignable users query (claimAssignableUsersQuery), and the result shows the claim's name and description. Turbo shows "claim not found" when the claim does not exist.

## Bulk assigning a claim: checks and result

Lines 139-172

An admin can assign one claim to many users at once. Turbo checks the claim first and shows "claim not found" if it is missing. A body that Turbo cannot read gives "invalid request body". On success Turbo returns the number of users assigned (assigned_count) and the number skipped (skipped_count).

## Bulk assign removes duplicate user ids

Lines 158-158

Turbo cleans the list of user ids with normalizeDistinctRecordIDs before it assigns. The rules for the cleaning are in another file.

## Bulk assign is all or nothing

Lines 159-168

Turbo runs the bulk assign in one transaction. If any user id does not exist or any save fails, Turbo shows "failed to bulk assign claim" and the transaction rolls back, so no user in that request gets the claim.

## Empty user list assigns nothing

Lines 175-177

When the list of user ids is empty, Turbo assigns no claims and skips none. It returns zero for both counts and shows no error.

## Bulk assign skips users who already hold the claim

Lines 186-214

Turbo skips a user who already holds the claim, counts them as skipped and creates no second copy. Turbo creates a new user claim record for every other user and counts them as assigned.

## Any signed-in request can list claims without an admin check

Lines 219-227

The claims list handler runs no admin check. It returns every claim with its id, name, description and holder count (holder_count). Whether a person must sign in to reach it depends on the route setup in another file.
