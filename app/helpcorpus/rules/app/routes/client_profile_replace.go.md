<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/client_profile_replace.go","sha256":"a1f7e18089a13b7c07df62e1110c862ecb837cbdbf812d396f5f02b07d0bc0ae","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/client_profile_replace.go

## Replace action returns counts of moved jobs and requests

Lines 18-23

Asked as: how many jobs moved; replace result counts; how many requests were moved; replace summary numbers

When the replace succeeds, Turbo reports how many jobs and how many job requests it moved to the replacement profile (jobs and requests).

## Replace invoicing profile runs as one all-or-nothing transaction

Lines 33-64

Asked as: replace failed partway; partial replace changes; error during profile replace; rollback invoicing profile

Turbo moves the jobs and requests and deletes the profile in one transaction. If any step fails, nothing changes. Turbo shows the error message from the failing check.

## Replacing an invoicing profile requires a signed-in user

Lines 39-41

Asked as: unauthorized replace profile; not logged in replace; replace without sign in; session missing replace

The replace-invoicing-profile action refuses a request that has no signed-in user. Turbo answers 'unauthorized'.

## Replacing an invoicing profile requires the job claim

Lines 96-102

Asked as: who can replace profile; no permission to replace invoicing; procurement role needed; can't replace profile as user

Only a person with the job claim (procurement) can replace a client's invoicing profile. Anyone else sees: "A user with the job claim is required to replace an invoicing profile."

## Replacing an invoicing profile needs job editing switched on

Lines 103-110

Asked as: replace disabled jobs editing; job editing turned off; can't replace profile now; jobs locked replace

Replacing an invoicing profile edits jobs, so it follows the jobs editing switch. When job editing is disabled, Turbo refuses the action and returns utilities.ErrJobsEditingDisabled (its message is defined in another file).

## Profile must exist and belong to the client

Lines 112-115

Asked as: profile not found for client; invoicing profile missing; wrong client profile; can't find invoicing profile

Turbo refuses with "Invoicing profile not found for this client." when the profile does not exist or belongs to a different client than the one in the request.

## Replacement profile is required

Lines 116-118

Asked as: no replacement chosen; pick new invoicing profile; replacement field empty; which profile to use instead

If no replacement is given, Turbo refuses with "Select the invoicing profile to use instead."

## Replacement must differ from the profile

Lines 119-121

Asked as: same profile replacement; replace profile with itself; choose same invoicing profile; can't pick current profile

Choosing the profile being replaced as its own replacement is refused with "Select a different invoicing profile."

## Replacement must exist and share the client

Lines 122-128

Asked as: replacement profile not found; replacement belongs other client; wrong client replacement profile; new profile different client

If the replacement profile is not found, Turbo says "The selected invoicing profile was not found." If it belongs to another client, Turbo says "The selected invoicing profile must belong to this client."

## Replace moves all jobs and requests, then deletes the profile

Lines 130-194

Asked as: move jobs to another profile; merge invoicing profiles; closed jobs moved too; delete old invoicing profile

Turbo moves every job and project request that uses the profile to the replacement profile, in any status (Closed jobs included). It then deletes the old profile and reports how many jobs and requests moved. It runs in one transaction, so a failure changes nothing.

## Pending client or contact merge blocks replacement

Lines 138-150

Asked as: merge waiting to commit; pending merge blocks replace; client merge in progress; contact merge stuck undo

Turbo refuses the replacement while a client or contact merge that involves either profile, or their jobs or project requests, is waiting to be committed or undone. The message is "A client or contact merge that involves these invoicing profiles or their jobs or requests is waiting to be committed or undone. Finish it before replacing the profile."

## Draft, Recalled and Rejected requests get a new version

Lines 152-183

Asked as: request edit form stale; draft request version changed; approval lost after move; moved request reload form

A moved project request in Draft, Recalled or Rejected status gets version + 1 and a new last_edited, so an edit form opened earlier must reload. Requests in other statuses keep their version and last_edited, and so keep their approval. All moved requests get a new updated time.

## Moved jobs are marked for the jobs writeback

Lines 164-173

Asked as: jobs writeback not syncing; moved jobs not exported; job sync flag after replace; jobs not sent to accounting

Each moved job gets _imported set to false and updated set to now, so the jobs writeback picks it up. The job's other hooks and validation do not run.

## Replacement is checked against client after the move

Lines 184-191

Asked as: profile still in use; delete old profile refused; record uses other client; replacement client mismatch after move

After the move, Turbo checks that every record that uses the replacement shares its client. If the old profile is still in use, the delete is refused.
