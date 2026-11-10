<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/claims_api.go","sha256":"565ef6d8e71ff2803fbc2656cc870783bde6647ebc605f0fc0d871c172404ddd","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/claims_api.go

## Only admins can use claim detail and assignment actions

Lines 69-78

Asked as: you do not have permission; admin claim required; not allowed to view claim details; no access to bulk assign; non-admin cannot assign claims

Turbo checks for the admin claim before it shows claim details, shows assignable users or bulk assigns a claim. A person without the admin claim gets "you do not have permission to " plus the action: "view claim details", "view claim assignment details" or "bulk assign claims".

## Viewing claim details needs an admin and an existing claim

Lines 80-109

Asked as: claim not found; view claim details; who holds this claim; claim holders list; claim description missing

An admin can view a claim's name, description and holders. Turbo shows "claim not found" when the claim id does not exist. The holders come from the claim details query (claimDetailsQuery).

## Viewing assignable users for a claim

Lines 111-137

Asked as: who can be assigned a claim; list assignable users; users available for claim; claim assignment candidates; assignable people for claim

An admin can list the users who can be assigned a claim. The list comes from the assignable users query (claimAssignableUsersQuery), and the result shows the claim's name and description. Turbo shows "claim not found" when the claim does not exist.

## Bulk assigning a claim: checks and result

Lines 139-172

Asked as: assign claim to many users; bulk assign result; assigned_count skipped_count; invalid request body; claim assignment response counts

An admin can assign one claim to many users at once. Turbo checks the claim first and shows "claim not found" if it is missing. A body that Turbo cannot read gives "invalid request body". On success Turbo returns the number of users assigned (assigned_count) and the number skipped (skipped_count).

## Bulk assign removes duplicate user ids

Lines 158-158

Asked as: duplicate user ids; remove duplicate users; same user twice in list; repeated user ids in assign; dedupe users when assigning

Turbo cleans the list of user ids with normalizeDistinctRecordIDs before it assigns. The rules for the cleaning are in another file.

## Bulk assign is all or nothing

Lines 159-168

Asked as: failed to bulk assign claim; bulk assign rolled back; one user id missing assign; partial claim assignment; nobody got the claim

Turbo runs the bulk assign in one transaction. If any user id does not exist or any save fails, Turbo shows "failed to bulk assign claim" and the transaction rolls back, so no user in that request gets the claim.

## Empty user list assigns nothing

Lines 175-177

Asked as: empty user list; assign zero users; no users selected for claim; bulk assign nothing selected; zero assigned zero skipped

When the list of user ids is empty, Turbo assigns no claims and skips none. It returns zero for both counts and shows no error.

## Bulk assign skips users who already hold the claim

Lines 186-214

Asked as: user already has claim; already assigned claim; skipped users claim; duplicate claim record; claim assigned twice

Turbo skips a user who already holds the claim, counts them as skipped and creates no second copy. Turbo creates a new user claim record for every other user and counts them as assigned.

## Any signed-in request can list claims without an admin check

Lines 219-227

Asked as: list all claims; claims without admin check; holder count of claims; show every claim; claims list no permission

The claims list handler runs no admin check. It returns every claim with its id, name, description and holder count (holder_count). Whether a person must sign in to reach it depends on the route setup in another file.
