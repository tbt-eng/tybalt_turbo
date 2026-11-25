<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/client_profile_replace.go","sha256":"cb8b9210c1c156188122f677067867e813b024a2d1b53dea3b04f74065e70ede","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/client_profile_replace.go

## Replacing an invoicing profile is all or nothing

Lines 33-51

Asked as: replace invoicing profile all or nothing; replacement fails partway; half done profile change; rollback profile replace

Replacing an invoicing profile moves every job and job request that uses it to another profile of the same client. Turbo then deletes the old profile. It does all of this in one transaction, so a failure at any step changes nothing.

## Replacing an invoicing profile requires sign-in

Lines 39-41

Asked as: replace profile not signed in; unauthorized replace profile; log in to replace profile; session expired profile replace

Turbo refuses the replace-invoicing-profile action for a person who is not signed in. The message is "unauthorized".

## Replacement moves jobs and requests of every status

Lines 67-91

Asked as: closed job blocks profile replace; old closed jobs stop replacement; legacy job prevents profile change; closed request blocking profile swap

Replacing an invoicing profile moves Closed jobs and job requests in any status as well as the others. Turbo updates them directly and does not save each record. The job and request edit checks therefore do not run, so a legacy or older Closed job cannot block the replacement. Turbo checks only the two profiles.

## Who may replace an invoicing profile

Lines 94-104

Asked as: who can replace invoicing profile; need job claim to replace; procurement permission profile change; why can't I replace profile

Only an active user with the job claim (procurement) may replace a client's invoicing profile. Anyone else sees "An active user with the job claim is required to replace an invoicing profile."

## Replacing a profile follows the jobs editing switch

Lines 105-112

Asked as: replace profile jobs editing off; jobs editing disabled profile; cannot change profile editing switch; jobs editing turned off error

Replacing an invoicing profile edits jobs, so it works only while jobs editing is enabled. When editing is off, Turbo refuses with the jobs-editing-disabled error (ErrJobsEditingDisabled).

## Profile must belong to the client

Lines 114-117

Asked as: invoicing profile not found for client; wrong client profile replace; profile belongs to other client; profile missing from client

Turbo looks up the profile to replace under the client in the request. If the profile does not exist or belongs to another client, it says "Invoicing profile not found for this client."

## Replacement profile is required

Lines 118-120

Asked as: no replacement profile chosen; which profile to use instead; pick replacement profile; missing profile to switch to

The request must name a replacement profile. Without one Turbo says "Select the invoicing profile to use instead."

## Replacement must differ from the profile

Lines 121-123

Asked as: same profile selected; replace profile with itself; choose different profile; replacement same as old

The replacement cannot be the profile being replaced. Turbo says "Select a different invoicing profile."

## Replacement must exist and share the client

Lines 124-130

Asked as: replacement profile not found; new profile belongs another client; selected profile wrong client; replacement profile missing

If the replacement profile does not exist, Turbo says "The selected invoicing profile was not found." If it belongs to another client, Turbo says "The selected invoicing profile must belong to this client."

## Pending client or contact merge blocks replacement

Lines 140-152

Asked as: merge blocking profile replace; client merge pending; contact merge waiting commit; undo merge before replacing

Turbo refuses the replacement while a client or contact merge that involves these profiles, or their jobs or job requests, is waiting to be committed or undone. The message is "A client or contact merge that involves these invoicing profiles or their jobs or requests is waiting to be committed or undone. Finish it before replacing the profile."

## Jobs move to the replacement profile

Lines 154-163

Asked as: jobs moved to new profile; job profile changed automatically; jobs switch to replacement profile; jobs not updating after replace

Every job that uses the old profile moves to the replacement, whatever its status, Closed included. Turbo sets the job's _imported to false and updated to now, so the jobs writeback picks it up. It does not run the job save hooks or validation.

## Job requests move and become stale

Lines 164-173

Asked as: job requests stale after replace; request needs review again; request version changed profile; approval outdated after profile change

Every job request that uses the old profile moves to the replacement, in any status. Turbo adds 1 to the request's version and sets last_edited to now, so a review in progress or an approval waiting for creation is stale and the request needs review again. It does not run the request save hooks.

## Old profile is deleted and counts are returned

Lines 174-184

Asked as: old profile deleted after replace; how many jobs moved; profile replace counts; old invoicing profile removed

After the move, Turbo checks that every record that uses the replacement shares its client, then deletes the old profile. It reports how many jobs and job requests moved. The whole action runs in one transaction, so a failure changes nothing.
