<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/absorb_records.go","sha256":"99dcd2b0c993d107e23058baa0d0b91ba23259f1a4016b08440e7321def6a811","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/absorb_records.go

## Absorb request body must be valid

Lines 44-49

Asked as: absorb error invalid body; request body unreadable; can't read absorb request; absorb fails bad data

Turbo refuses an absorb when the request body cannot be read. The message is "Invalid request body".

## Absorb needs at least one record to absorb

Lines 51-53

Asked as: nothing selected to merge; no records picked to absorb; merge with no ids; merge button no records

An absorb request with an empty list of records to absorb (ids_to_absorb) is refused with 'No IDs provided to absorb'.

## A record cannot be absorbed into itself

Lines 55-60

Asked as: merge record into itself; same record selected twice; target also in selected list

If the target record is also in the list of records to absorb, Turbo refuses with 'Cannot absorb a record into itself'.

## Absorbing needs the collection's absorb permission

Lines 62-71

Asked as: can't merge records; no rights to merge; absorb button not allowed; missing permission to combine

Only a person with the absorb claim, or the claims that the collection's config asks for instead, can absorb. Anyone else gets 'User does not have permission to absorb records'.

## Absorb is refused when a claim check fails

Lines 62-68

Asked as: can't verify my permissions; claims check error absorb; permission lookup failed

If Turbo cannot check the user's claims, the absorb is refused with "Failed to check user claims". This check is separate from the refusal "User does not have permission to absorb records".

## Absorbs that change jobs need job editing enabled

Lines 73-82

Asked as: merge jobs turned off; job editing disabled absorb; can't merge jobs; absorb changes jobs blocked

For collections whose absorb changes jobs, Turbo refuses with 'Job editing is disabled; cannot absorb records that modify jobs' when job editing is turned off.

## Vendor absorb and undo need expenses editing enabled

Lines 84-87

Asked as: vendor merge expenses editing off; can't merge vendors expenses; undo vendor merge blocked; expenses editing disabled merge

Turbo checks the expenses-editing setting (requireExpensesEditing) before an absorb and before an undo. This applies to vendor absorbs, and the helper's error is returned when editing is off.

## Absorb is refused when the collection is not supported

Lines 89-93

Asked as: merge not available for this; absorb unsupported collection; can't merge this type; absorb fails for collection

Turbo looks up the reference configuration of the collection before it absorbs. If the collection has none, the absorb fails with "Failed to absorb records".

## Target and absorbed records must exist

Lines 95-105

Asked as: record not found merge; merge deleted client; target record missing; absorbed record gone

Turbo refuses with 'Failed to find target record' if the target is missing. It refuses with 'Failed to find record to absorb' if any record in the list is missing.

## Records must share the same parent value

Lines 108-118

Asked as: merge records different parent; absorb different client; different field values merge; records belong to different

If the collection has a parent constraint, every record to absorb must have the same value in that field as the target. Otherwise Turbo refuses with 'Cannot absorb records with different <field> values'.

## Absorbed client's invoicing profiles billing an alias need a kept alias

Lines 121-134

Asked as: client alias needed merge; invoicing profile alias merge; merge client alias error; kept client no alias

When absorbing clients, if the kept client has no alias and an absorbed client has an invoicing profile that bills its alias, Turbo refuses. The message is 'Invoicing profiles of a client being absorbed bill its alias. Give the kept client an alias first.'

## Absorb that carries out a proposal is checked again

Lines 136-144

Asked as: proposal merge already decided; absorb proposal not pending; proposal wrong record type; proposal absorb failed check

When an absorb names an absorb proposal, Turbo checks the proposal first. It refuses with 'absorb proposal not found' if it does not exist. It refuses with 'the absorb proposal is for another kind of record' if the collection differs. It refuses with 'this absorb proposal has already been decided' unless the proposal is pending.

## Collection rules can refuse an absorb

Lines 146-155

Asked as: merge blocked by rule; nothing to absorb message; collection rule stopped merge; preparation step stopped absorb

If a rule of the collection's config refuses the absorb, Turbo shows that rule's message. If the collection's preparation step leaves nothing to absorb, the message is 'there is nothing to absorb'.

## Absorb success message

Lines 157-159

Asked as: merge succeeded message; how many records merged; absorb confirmation text; merge done message

A successful absorb replies 'Successfully absorbed N records into <target id>'. N is the number of records in the request.

## Absorb proposal must exist, match the collection and be pending

Lines 167-180

Asked as: proposal not found absorb; proposal already decided; proposal wrong collection; proposal no longer open

When an absorb carries out a proposal, Turbo checks the proposal itself. It refuses with "absorb proposal not found" if the proposal does not exist. It refuses with "the absorb proposal is for another kind of record" if the proposal is for another collection. It refuses with "this absorb proposal has already been decided" if the proposal is not pending (absorb.ProposalPending).

## Proposal absorb may drop records but not add any

Lines 185-195

Asked as: records not in proposal; added extra record to merge; every record must come from proposal; merge picked other record

The target and every absorbed record must come from the proposal, either as its target or as one of its proposed records. The reviewer may leave records out or keep another proposed record as target. Any other record is refused with 'every record must come from the absorb proposal'. Turbo then runs the proposal validation again, since the records may have changed.

## Proposal absorb refuses records outside the proposal

Lines 185-193

Asked as: record outside proposal; merge uses unlisted record; not part of proposal; proposal merge wrong records

In an absorb that carries out a proposal, the target and every record to absorb must be the proposal's target or one of its proposed records. Otherwise Turbo refuses with "every record must come from the absorb proposal". The reviewer may leave records out or keep another of them as target.

## Proposal checks run again at absorb time

Lines 194-195

Asked as: proposal check again merge; records changed since proposal; revalidate proposal on absorb; proposal validation at merge

When an absorb carries out a proposal, Turbo runs the proposal validation again (absorb.ValidateProposal) on the chosen target and records, because the records may have changed since the proposal was made.

## Absorb is all-or-nothing

Lines 224-407

Asked as: merge half done; partial absorb rollback; merge failed nothing changed; all or nothing merge

Turbo runs the whole absorb in one transaction. If any step fails, the absorb fails with "Failed to absorb records" and nothing changes.

## Absorb with no records left is refused

Lines 226-240

Asked as: rules removed all records; nothing left to merge; preparation cleared records; merge list empty after rules

When a collection's own preparation rules (config.Prepare) leave no ids to absorb, Turbo refuses with "there is nothing to absorb". The rules run first, inside the transaction, and may add ids to the list.

## Only one pending absorb per collection

Lines 266-273

Asked as: already a merge pending; second merge blocked; undo existing merge first; merge already in progress

If the collection already has an absorb action, a new absorb fails. The existing one must be undone or deleted first.

## Absorb moves references to the target and deletes the absorbed records

Lines 275-394

Asked as: merged records deleted references moved; what happens to merged client; references repointed on merge; merge deletes absorbed record

Turbo changes every reference to an absorbed record so it points to the target. It then deletes the rows that the absorbed records own, then the absorbed records. It saves them all in an absorb action so the absorb can be undone. All of it happens in one transaction, so a failure rolls everything back.

## Moved references are marked for legacy writeback

Lines 309-322

Asked as: legacy system sync merge; imported flag merge; writeback after merge; merge not sent to legacy

For tables with an _imported column, an absorb sets _imported to false and updates the updated time on each moved reference, so the change is written back to the legacy system.

## Absorb must not leave records linked to another client

Lines 341-343

Asked as: job linked wrong client merge; merge leaves other client; invoicing profile client error; merge breaks client link

After the references move, Turbo checks every job, project request, contact and invoicing profile whose client relation changed. If it would leave a record linked to another client, the absorb fails. An invoicing profile that would bill a missing client alias also fails.

## Absorb that moves no references still records them

Lines 345-355

Asked as: merge moved nothing; rate sheet not used merge; empty merge action; merge no references changed

If no reference moves, such as a rate sheet name that nothing uses, Turbo records each table and column with no rows in the action.

## Absorbed records are saved with the rows they own

Lines 367-381

Asked as: child rows merge saved; owned records undo merge; merge deletes child rows; saved rows before merge delete

Before it deletes absorbed records, Turbo saves every field of each record and the rows it owns, such as child rows in other tables. It then deletes those owned rows with plain SQL, so undo can restore them.

## Absorb completion notifies connected users

Lines 409-412

Asked as: screen not refreshing after merge; live update absorb; notification merge completed; merge event not sent

After an absorb or an undo, Turbo sends a realtime event "<collection>/absorb_completed" to clients that subscribe to it, so their screens refresh. A failure to send does not fail the absorb.

## Undo needs the same permission as absorb

Lines 433-444

Asked as: can't undo merge; undo permission denied; no rights undo absorb; who can undo merge

Undoing an absorb needs the same claims as absorbing. Anyone else gets 'User does not have permission to undo absorb'.

## Undo order and permission messages

Lines 437-456

Asked as: undo no permission message; undo nothing pending; no absorb to undo; undo error message

Undo is refused with "User does not have permission to undo absorb" without the required claims. It is refused with "No absorb action found for collection" when no absorb is pending.

## Undo needs a pending absorb action

Lines 446-456

Asked as: undo without merge; nothing to undo; no absorb action found; undo when no merge

If the collection has no absorb action, undo is refused with 'No absorb action found for collection'.

## Child absorbs must be undone first

Lines 458-474

Asked as: undo blocked by other merge; dependent merge undo first; child collection merge exists; undo parent merge error

Turbo refuses to undo an absorb while a dependent collection still has its own absorb. The message is 'Cannot undo X absorb while Y absorb exists; undo Y absorb first'.

## Undo success message

Lines 478-589

Asked as: undo worked message; undo confirmation; absorb undone text; undo success text

When an undo succeeds, Turbo answers "Successfully undid absorb operation". An undo that fails inside the transaction is rolled back entirely.

## Undo restores records, owned rows and references

Lines 479-576

Asked as: bring back merged records; restore merged client; undo restores references; undo merge restore data

Undo recreates the absorbed records with their original ids and data, and the rows they owned. It skips schema validation, because old records may predate current rules. It sets every moved reference back to its old value and deletes the absorb action, so a new absorb can follow. It reads the action again inside the transaction, so references moved since the absorb are kept.

## Undo reads latest undo data

Lines 479-493

Asked as: undo after job created; job creation references undo; latest absorb data undo; undo keeps moved references

Undo reads the latest saved absorb action inside its transaction. This keeps the references that job creation moved to a job.

## Undo recreates records without validation

Lines 510-515

Asked as: restore record validation error; undo fails required field; old records undo validation; undo skips field checks

Undo recreates absorbed records and their owned rows without schema validation, because they may predate current rules such as newly required fields.

## Undo marks restored references for legacy writeback

Lines 527-548

Asked as: restored references legacy sync; undo imported flag; undo writeback; undo sent to legacy

When undo restores references in a table with an _imported column, Turbo sets _imported to false and updates the updated time, so the change is written back to the legacy system again.

## Undo skips records deleted since the absorb

Lines 560-562

Asked as: undo deleted records error; undo client relation deleted; undo wrong client link; undo skips removed records

When undo re-checks client relations, it skips records that were deleted since the absorb. It still fails if a restored relation would link to another client.

## Undo deletes the absorb action

Lines 564-589

Asked as: can merge again after undo; clear merge after undo; remove merge record; reset merge after undo

After a successful undo Turbo deletes the absorb action, so the collection can absorb again. It then shows "Successfully undid absorb operation".

## Absorb history records who absorbed and the proposal

Lines 613-633

Asked as: who merged records; merge history user; audit merge who; which user absorbed

Turbo saves the absorb action with the collection, target, absorbed records, moved references, the proposal (if any) and the user who absorbed (absorbed_by). The route always passes the user, so the history of an absorb names them.

## Client absorb checks invoicing profile alias after the move

Lines 801-834

Asked as: invoicing profile alias after merge; alias missing after merge; merge alias check; profile alias error merge

After an absorb moves references, Turbo checks that the changed invoicing profiles still bill an alias that exists, and fails the absorb otherwise.

## Changed project requests become stale on absorb or undo

Lines 835-840

Asked as: project request review stale; version changed merge request; merge changes project request; request edited after merge

When an absorb or undo changes a client relation of a project request, Turbo adds 1 to its version and updates its last-edited and updated times. Any open review of that request becomes stale.

## Creating a job moves a request's undo references to the job

Lines 846-898

Asked as: job created from request undo; request merge moved to job; client note job creation; undo references job

When a job is created from a project request, the request's pending undo references move to the job. The same applies to the client note and the opening rate sheet period. Turbo skips a relation that the creator changed after the absorb, and removes the request's entries from the action.

## Moved job request relations are kept only if unchanged

Lines 861-891

Asked as: creator changed client after merge; rate sheet period job merge; undo skips changed job; job merge keeps changes

When a job is created from a request, Turbo moves the request's pending undo references to the job. It skips a relation that the creator changed after the absorb, so undo does not restore it. For a client change, it also moves the client note made from the request's notes, and for a rate sheet change the job's opening rate sheet period.
