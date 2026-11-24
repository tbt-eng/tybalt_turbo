<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/absorb_records.go","sha256":"daf2b031664fd5f5a0abd95e6ab4c48ea264d23e41065eba8621a9fb5ac5bf6a","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/absorb_records.go

## Absorb needs at least one record to absorb

Lines 51-53

Asked as: absorb nothing; no records selected to merge; absorb with empty list; nothing to absorb in request

Turbo refuses an absorb when the request lists no records to absorb. It shows "No IDs provided to absorb".

## A record cannot be absorbed into itself

Lines 55-60

Asked as: merge a record into itself; absorb same record twice; target also in the list; choose the same client twice

Turbo refuses an absorb when the target record is also in the list of records to absorb. It shows "Cannot absorb a record into itself".

## Only users with absorb access can absorb records

Lines 62-71

Asked as: cannot absorb, no permission; can't merge records access; absorb button missing for me; not allowed to merge clients

Turbo refuses an absorb when the user lacks the absorb claim, or the claims that the collection's absorb config requires instead. It shows "User does not have permission to absorb records".

## Absorbs that change jobs need job editing enabled

Lines 73-82

Asked as: absorb blocked job editing off; merge jobs not allowed; job editing turned off; can't merge jobs

When the collection's absorb changes jobs (ChangesJobs), Turbo refuses the absorb while job editing is disabled. It shows "Job editing is disabled; cannot absorb records that modify jobs".

## Vendor absorbs and undos need expenses editing enabled

Lines 84-87

Asked as: vendor merge blocked; expenses editing disabled absorb; can't merge vendors; undo vendor merge not working

Turbo checks that expenses editing is enabled (requireExpensesEditing) before it absorbs or undoes an absorb for a collection. The check is meant for vendor absorbs. The function's own rules are in another file.

## Absorb target record must exist

Lines 95-99

Asked as: target record missing; merge into deleted record; can't find the kept record; target does not exist

When a person absorbs records, Turbo looks up the target record first. If it does not exist, Turbo refuses with 'Failed to find target record'.

## Every record to absorb must exist

Lines 101-106

Asked as: record to merge not found; one of the records is gone; missing record in merge; duplicate record already deleted

Turbo checks that each record to absorb exists in the collection. If one is missing, Turbo refuses with 'Failed to find record to absorb'.

## Absorbed records must share the parent value of the target

Lines 108-118

Asked as: merge records different parent; different client on records; can't merge across parents; records belong to different owners

For a collection with a parent constraint, each record to absorb must have the same parent field value as the target. Otherwise Turbo refuses with 'Cannot absorb records with different <field> values'.

## Absorbing clients needs an alias on the kept client

Lines 121-134

Asked as: merging clients alias error; client alias needed for merge; invoice bills alias when merging; kept client has no alias

If the kept client has no alias, Turbo checks each client being absorbed. If any invoicing profile of an absorbed client bills its alias, Turbo refuses with 'Invoicing profiles of a client being absorbed bill its alias. Give the kept client an alias first.'

## Absorb carrying out a proposal must use a pending proposal

Lines 136-144

Asked as: approve merge proposal fails; proposal not found on merge; merge proposal already decided; proposal wrong type for merge

When an absorb names an absorb proposal, Turbo refuses it if the proposal is not found ("absorb absorb proposal not found" is shown as "absorb proposal not found"). It refuses it with "the absorb proposal is for another kind of record" if the proposal is for another collection. It refuses it with "this absorb proposal has already been decided" unless the proposal is pending.

## Collection rules can add records or refuse the absorb

Lines 146-155

Asked as: merge rule stopped it; nothing left to merge; config blocks the merge; merge refused by collection rule

If the collection's config has a prepare step, it checks the request first inside the transaction and may add ids to absorb. A rule of the config that refuses the absorb shows its own message. If nothing is left to absorb, Turbo refuses with "there is nothing to absorb".

## Absorb success message and refresh broadcast

Lines 157-159

Asked as: merge done message; records merged successfully; screen not refreshed after merge; merge success confirmation

A successful absorb replies "Successfully absorbed N records into <target id>". Turbo then broadcasts a "<collection>/absorb_completed" event so open screens refresh. If the broadcast fails, the absorb still succeeds.

## Proposal absorb may only use records from the proposal

Lines 167-193

Asked as: merge adds record not in proposal; record not on the proposal; can't add records to proposal merge; proposal merge wrong records

When an absorb carries out a proposal, the target and every absorbed record must come from the proposal: the proposed target or a proposed record. The reviewer may leave records out or keep another proposed record as target, but cannot add one. Otherwise Turbo refuses with "every record must come from the absorb proposal".

## Proposal checks run again at absorb time

Lines 194-195

Asked as: proposal no longer valid; merge proposal check failed; records changed since proposal; proposal validation error on merge

When an absorb carries out a proposal, Turbo validates the proposal again with absorb.ValidateProposal, because the records may have changed since the proposal was made. A refusal from that check is shown with its own code and message.

## Only one pending absorb per collection

Lines 258-265

Asked as: merge already in progress; another merge pending; can't start second merge; existing merge must be undone

Turbo refuses an absorb if the collection already has an absorb action. The existing absorb must be undone or deleted first, with the error "there is an existing absorb action for collection ... that must be undone or deleted first".

## Absorb moves references to the target and deletes the absorbed records

Lines 267-395

Asked as: what merge changes; references moved after merge; records deleted by merge; merge rollback on failure

Absorb changes every reference to an absorbed record in the configured tables and columns so it points to the target record. It deletes the rows the absorbed records own, then deletes the absorbed records. Turbo saves the records, owned rows and old reference values in an absorb action so the absorb can be undone. All of this happens in one transaction, so a failure rolls back everything.

## Absorb marks changed records for legacy writeback

Lines 300-314

Asked as: merge not sent to legacy system; legacy system not updated; imported flag after merge; changed rows resent to legacy

When a table that absorb updates has an _imported column, Turbo sets _imported to false and refreshes the updated time on the changed rows. This makes the legacy system receive the change.

## Absorb with no moved references still records its tables

Lines 337-347

Asked as: merge with no references; empty merge history; merge shows no tables; rate sheet merge nothing moved

If an absorb moves no reference, such as a rate sheet name that nothing uses, Turbo records each configured table and column with no rows. This keeps the absorb action's updated references from being blank.

## Undo absorb requires absorb access

Lines 425-436

Asked as: cannot undo merge no access; undo permission denied; who can undo a merge; undo absorb not allowed

Undoing an absorb needs the same access as absorbing, the absorb claim or the claims the collection's config asks for instead. Without it Turbo refuses with "User does not have permission to undo absorb".

## Undo needs an absorb to undo

Lines 438-448

Asked as: nothing to undo; no merge to reverse; undo button but no merge; no absorb found to undo

Turbo refuses an undo with "No absorb action found for collection" when the collection has no absorb action. A collection has at most one at a time.

## Undo child collection absorbs first

Lines 450-466

Asked as: undo blocked by child merge; child collection merged first; parent merge can't be undone; undo order for merges

A parent collection's absorb cannot be undone while an absorb exists for a collection that depends on it. Turbo refuses with "Cannot undo X absorb while Y absorb exists; undo Y absorb first".

## Undo recreates records and restores references

Lines 471-568

Asked as: restore merged records; bring back deleted clients; reverse merge fully; undo keeps later project references

Undo recreates the absorbed records with their original data and ids, plus the rows they owned, without running current validation. It then sets each saved reference back to its old value, and deletes the absorb action so a new absorb is allowed. Everything happens in one transaction. Undo uses the latest saved data, so references moved by project creation are kept.

## Undo marks restored records for legacy writeback

Lines 519-532

Asked as: undo not sent to legacy; restored records legacy sync; imported flag after undo; legacy not updated on undo

When undo restores references in a table with an _imported column, Turbo sets _imported to false and refreshes the updated time. The legacy system then receives the record again.

## Undo skips records deleted since the absorb

Lines 552-554

Asked as: deleted record after merge undo; undo error missing client; records removed since merge; client check after undo

When Turbo checks client references after an undo, it skips records that were deleted since the absorb.

## Undo success message

Lines 574-581

Asked as: undo done message; undo success confirmation; screen not refreshed after undo; undo completed notice

A successful undo replies "Successfully undid absorb operation". Turbo broadcasts a "<collection>/absorb_completed" event so open screens refresh.

## Merged client references must still agree

Lines 792-825

Asked as: merged client links wrong; job linked to other client; merge leaves contact mismatched; invoicing profile alias missing after merge

After an absorb or an undo, Turbo checks every rewritten job, project request, client contact and invoicing profile whose client relation changed. If they would link to another client, the absorb fails with "merge would leave ... linked to another client". For invoicing profiles whose client changed, it also fails when the profile bills a client alias that does not exist.

## Absorb makes open reviews of project requests stale

Lines 826-831

Asked as: review stale after merge; request changed during review; project request version conflict; merge makes review out of date

When an absorb or undo rewrites a project request, Turbo adds 1 to its version and sets its last edited and updated times to now. This makes any open review of that request stale.

## Creating a job from a request keeps absorb undo data

Lines 837-889

Asked as: request to job merge undo; job created from request undo; merge data after job creation; rate sheet moved to job

When a project request is turned into a job, absorb undo data that pointed to the request moves to the job. This covers clients, client contacts and rate sheets. The job's client note and its opening rate sheet period move too. A relation that the creator changed after the merge is not restored. The request's own entry is then removed from the undo data.
