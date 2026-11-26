<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/client_profile_replace.go","sha256":"4fb92adc76a8565b394f2277e41c2e46c5d0b2cee3c969d937a8dcba60c86bcb","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/client_profile_replace.go

## Replace result reports how many jobs and requests moved

Lines 18-23

Asked as: replace result count; how many jobs moved; how many requests moved; replace summary numbers

After a successful replace, Turbo reports the number of jobs and the number of project requests (job requests) that moved to the replacement invoicing profile.

## Replacing an invoicing profile is all or nothing

Lines 33-51

Asked as: partial replace failure; replace rolled back; replace half done; invoicing profile changed partly

Turbo moves the jobs and project requests and deletes the old invoicing profile in one transaction. If any step fails, nothing changes.

## Replacing an invoicing profile requires being signed in

Lines 39-41

Asked as: replace profile not logged in; replace unauthorized; not signed in replace; session expired replace

Turbo refuses a replace-invoicing-profile request from someone who is not signed in. The reply is "unauthorized". Turbo checks this before anything else.

## Replace changes only the invoicing profile, bypassing edit hooks

Lines 67-80

Asked as: replace blocked by closed job; replace ignores edit rules; replace fails validation old job; replace skips job validation

Turbo updates the jobs and requests with direct SQL, not by saving each record. The job and request edit rules do not run, so a creator-only edit rule or a full job validation of an older Closed job cannot block the replace. Only the invoicing profile field changes.

## Requests awaiting approval keep version and approval after replace

Lines 85-90

Asked as: request awaiting approval after replace; approval lost after profile change; request version after replace; approved request profile swap

The invoicing profile is not part of what an approval covers. A project request that awaits approval or job creation keeps its version and its approval when its profile is replaced.

## Replace moves records in every status, including Closed

Lines 92-94

Asked as: replace closed jobs; move closed job profile; replace profile on closed requests; closed job wrong invoicing profile

Replacing an invoicing profile works the same for every status. Closed jobs and project requests in any status move to the replacement profile too.

## Only active job-claim holders can replace an invoicing profile

Lines 96-106

Asked as: who can replace invoicing profile; replace button missing; need procurement to replace profile; cannot replace invoicing profile permission

Replacing a client's invoicing profile requires an active user who holds the job claim (procurement). Anyone else gets: "An active user with the job claim is required to replace an invoicing profile."

## Replacing an invoicing profile follows the jobs editing switch

Lines 107-114

Asked as: replace profile jobs editing off; jobs editing disabled replace; cannot change profile jobs locked; editing switch replace profile

When jobs editing is turned off, Turbo refuses to replace an invoicing profile. It returns the jobs-editing-disabled error (ErrJobsEditingDisabled), because the action edits jobs.

## The profile to replace must belong to the client

Lines 116-119

Asked as: invoicing profile not found; old profile wrong client; profile belongs to other client; profile not for this client

If the invoicing profile does not exist, or belongs to a different client than the one in the request, Turbo refuses with "Invoicing profile not found for this client."

## A replacement profile is required

Lines 120-122

Asked as: no replacement profile chosen; blank replacement profile; pick new invoicing profile; replacement missing

If no replacement is chosen, Turbo refuses with "Select the invoicing profile to use instead."

## The replacement must differ from the profile being replaced

Lines 123-125

Asked as: same profile replacement; replace profile with itself; pick different profile; both profiles same

If the replacement is the same profile as the one being replaced, Turbo refuses with "Select a different invoicing profile."

## The replacement must exist and belong to the same client

Lines 126-132

Asked as: replacement profile not found; new profile other client; replacement wrong client; chosen profile unavailable

If the replacement profile is not found, Turbo says "The selected invoicing profile was not found." If it belongs to another client, Turbo says "The selected invoicing profile must belong to this client."

## Replacing a profile moves all its jobs and requests, then deletes it

Lines 134-198

Asked as: old profile deleted after replace; move all jobs to new profile; profile replacement count; delete invoicing profile moving jobs

Turbo moves every job and every project request (job request) that uses the profile to the replacement profile, then deletes the old profile. The jobs and requests move in any status, including Closed jobs. It returns the number of jobs and requests moved. It all runs in one transaction, so a failure changes nothing.

## A pending client or contact merge blocks the replacement

Lines 142-154

Asked as: merge pending replace profile; client merge blocks replace; contact merge waiting commit; undo merge before replacing

If a client or contact merge that touches these profiles, or their jobs or job requests, is waiting to be committed or undone, Turbo refuses. The message is: "A client or contact merge that involves these invoicing profiles or their jobs or requests is waiting to be committed or undone. Finish it before replacing the profile."

## Old approvals are upgraded before the move

Lines 156-166

Asked as: approval stale after replace; old digest approval; approval upgrade profile change; legacy approval replace

Before the move, Turbo upgrades any request approval saved under the old digest, which included the profile, to the current digest (upgradeLegacyRequestApproval). The move then does not make that approval stale.

## Moved jobs are flagged for writeback to the legacy system

Lines 169-177

Asked as: jobs not written back; legacy system sync profile change; imported flag job profile; job writeback after profile move

Each moved job gets _imported set to false and its updated time set to now, as a normal profile edit does, so the jobs writeback picks it up.

## Moved requests in Draft, Recalled or Rejected get a new version

Lines 178-187

Asked as: request version changed after replace; form must reload replace; draft request version bump; rejected request edit form stale

A request in Draft, Recalled or Rejected status gets its version raised by 1 and last_edited set to now. An open edit form must then reload. Requests in other statuses keep their version and last_edited, so an approval stays valid because the profile is not part of what an approval covers. Every moved request still gets a new updated time.

## Replacement is checked for client consistency and the old profile's delete is guarded

Lines 188-195

Asked as: client mismatch after replace; old profile delete refused; profile still in use delete; replace consistency check failed

After the move, Turbo checks that every record using the replacement shares its client (ValidateClientReferences). It then deletes the old profile, and the delete hook refuses if anything still uses it.
