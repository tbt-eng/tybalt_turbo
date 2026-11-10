<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/absorb_records.go","sha256":"ef7db93ae9b1e8a82949ae5d29d9513c87824f3198ec655c510bb7387159ed12","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/absorb_records.go

## Absorb request with an unreadable body is refused

Lines 43-46

Asked as: invalid request body; absorb body unreadable; absorb request won't parse; malformed absorb request

When Turbo cannot read the absorb request body, it refuses the absorb with the message "Invalid request body". Turbo makes this check before it checks permissions or looks at any record.

## Absorb needs a list of records to absorb

Lines 48-50

Asked as: no records selected to absorb; absorb with nothing selected; empty absorb list; merge without any records

Turbo refuses an absorb request whose list of IDs to absorb (ids_to_absorb) is empty. It shows the message "No IDs provided to absorb".

## A record cannot be absorbed into itself

Lines 52-57

Asked as: merge a record into itself; absorb same record twice; keep and absorb same client; merge record with itself

Turbo refuses an absorb when the kept record's ID is also in the list to absorb. It shows the message "Cannot absorb a record into itself".

## Only staff with the absorb claim can absorb

Lines 59-67

Asked as: can't merge clients; no permission to merge; cannot absorb duplicates; absorb button not allowed

Only a person with the absorb claim can absorb records. Anyone else sees "User does not have permission to absorb records".

## Absorbing clients or contacts needs job editing on

Lines 69-78

Asked as: merge clients job editing off; can't merge clients with jobs disabled; absorb contacts job editing disabled; enable job editing to merge

For clients and client contacts, Turbo refuses an absorb while job editing is disabled. It shows "Job editing is disabled; cannot absorb records that modify jobs", because these absorbs change jobs.

## Vendor absorb and undo need expense editing on

Lines 80-83

Asked as: merge vendors blocked; vendor absorb not working; merge vendors expenses editing; can't merge vendors

Both absorb and undo-absorb call requireExpensesEditing for the collection, which applies to vendor absorbs. This file does not show the exact condition or message.

## Absorb target record must exist

Lines 91-95

Asked as: target record missing; merge target not found; kept record doesn't exist; absorb destination not found

Turbo refuses an absorb when the target record (the id in the path) is not found in the collection. The message is "Failed to find target record".

## Every record to absorb must exist

Lines 97-102

Asked as: duplicate record missing; record to merge not found; one of the merged records gone; absorb source not found

Turbo refuses an absorb when any record in the list to absorb is not found in the collection. The message is "Failed to find record to absorb".

## Absorbed records must share the target's parent value

Lines 104-114

Asked as: merge records different parent; absorb different owner; merge across different parents; parent field mismatch merge

When the collection has a parent constraint, every record to absorb must have the same value in the parent field as the target. Otherwise Turbo refuses with "Cannot absorb records with different <field> values". The field name comes from the constraint (absorb.GetRefConfigs).

## Kept client needs an alias if absorbed clients' profiles bill an alias

Lines 117-130

Asked as: merge clients alias error; client needs alias to merge; invoicing profile alias merge; alias blocks client merge

When absorbing clients, if the kept client has no alias, Turbo checks each client being absorbed for invoicing profiles that bill its alias. If any exist, Turbo refuses with "Invoicing profiles of a client being absorbed bill its alias. Give the kept client an alias first."

## Absorb success message

Lines 138-140

Asked as: merge done confirmation; how many records merged; absorb success text; merge completed message

A successful absorb replies "Successfully absorbed N records into <target id>", where N is the number of ids in the request.

## Only one pending absorb per collection

Lines 180-187

Asked as: already merged once; second merge blocked; merge already pending; undo previous merge first

Turbo refuses an absorb when the collection already has a recorded absorb action. The existing absorb must be undone or deleted first. The error says "there is an existing absorb action for collection <name> that must be undone or deleted first".

## Absorb moves all references to the target and deletes the absorbed records

Lines 195-289

Asked as: merge moves everything; references repointed on merge; what happens when merging; merged records deleted

An absorb points every reference to an absorbed record (in the tables from absorb.GetRefConfigs) at the target record, then deletes the absorbed records. It saves the old records and old reference values as an absorb action so the absorb can be undone. All of this happens in one transaction, so any failure rolls back the whole absorb.

## Absorb marks changed legacy-synced rows for writeback

Lines 222-245

Asked as: merge not syncing legacy; legacy writeback after merge; imported rows merge; merge not reaching old system

For each updated reference row in a table with an _imported column, Turbo sets _imported to false and sets updated to the current time. This makes the change write back to the legacy system, because the direct update skips the normal hooks.

## Merge refuses to leave records linked to another client

Lines 255-257

Asked as: merge client link conflict; merge leaves job on other client; merge fails linked to client; invoicing profile alias missing on merge

After the references move, Turbo re-validates every changed job, job request, client contact and invoicing profile. If its client relations no longer agree, the absorb fails with "merge would leave <table> <id> linked to another client". An invoicing profile that would bill a client alias that does not exist also fails the absorb.

## Absorb completion notifies other users

Lines 295-298

Asked as: other users not refreshing after merge; merge notification missing; screen not updating after merge; refresh after merge event

After an absorb or an undo completes, Turbo sends an "<collection>/absorb_completed" realtime event to subscribed clients so they can refresh. A failure to send is logged and does not fail the absorb.

## Only staff with the absorb claim can undo an absorb

Lines 319-329

Asked as: can't undo merge; no permission to undo merge; undo merge not allowed; revert merge access denied

Undoing an absorb requires the absorb claim. Without it, Turbo refuses with "User does not have permission to undo absorb".

## Undo needs an existing absorb action

Lines 331-341

Asked as: nothing to undo; no merge to revert; undo merge not found; no previous merge

Turbo refuses an undo when the collection has no recorded absorb action. The message is "No absorb action found for collection".

## Child absorbs must be undone before parent absorbs

Lines 343-359

Asked as: undo merge blocked by child; undo parent merge child exists; child merge must be undone first; undo order merges

Turbo refuses to undo an absorb while a dependent child collection (from absorb.GetUndoBlockers) still has an absorb action. The message is "Cannot undo <collection> absorb while <child> absorb exists; undo <child> absorb first".

## Undo restores records and references

Lines 364-454

Asked as: undo merge restore records; revert merge success message; bring back merged records; undo completed text

An undo recreates the absorbed records with their saved data, without schema validation, because old records may not meet current required fields. It then sets each changed reference back to its old value and deletes the absorb action, so a new absorb is allowed. It runs in one transaction and rolls back on any failure. It replies "Successfully undid absorb operation".

## Undo uses latest saved undo data

Lines 365-378

Asked as: undo after project creation; undo lost job references; undo missing transferred data; project creation undo issue

The undo reads the absorb action inside its own transaction, so references transferred to a job by project creation are not lost.

## Undo skips deleted records when validating clients

Lines 425-427

Asked as: undo client validation error; undo fails deleted record; undo check clients error; undo skips deleted

After the restore, an undo re-validates client relations on jobs, requests, contacts and invoicing profiles. It skips records that were deleted since the merge. A record that fails validation fails the undo.

## Absorb saves full copies of absorbed records for undo

Lines 498-530

Asked as: what gets saved for undo; undo copy of records; backup of merged records; merged records saved copy

Before an absorb deletes records, Turbo saves each absorbed record's original id and the value of every field in the collection's schema. Turbo stores these copies in the absorb action (absorbed_records). Undo uses them to recreate the records.

## Absorb completion event goes only to subscribed users

Lines 532-569

Asked as: who sees merge update; merge event subscribers; realtime merge notice; merge update to users

After an absorb or an undo, Turbo sends the event "<collection>/absorb_completed" to connected users. Only users who subscribe to that topic for the collection receive it, and they can then refresh their data.

## Absorb remembers each changed reference's old value

Lines 571-595

Asked as: old references saved; track old values merge; what references were changed; merge remembers links

For every reference that an absorb moves, Turbo remembers the table, the column, the record id and the old value. Turbo saves this in the absorb action (updated_references). Undo needs it to point each reference back at the absorbed record.

## Merge checks four record types for client links that disagree

Lines 597-619

Asked as: client mismatch after merge; merge linked to wrong client; job client conflict merge; undo client relation error

After a merge or an undo, Turbo checks every job, project request (job_requests), client contact and invoicing profile that the operation changed. The client relations on each record must still agree. If they do not, Turbo refuses with "merge would leave <table> <id> linked to another client" and the reason.

## Merge refuses an invoicing profile billing a missing alias

Lines 620-624

Asked as: invoicing profile alias error; billing name missing alias; profile bills missing alias; merge invoicing profile error

After a merge or an undo, Turbo checks each changed invoicing profile (client_invoicing_information) for its billing name. If the profile would bill a client alias that does not exist, Turbo refuses with "merge would leave invoicing profile <id> billing a client alias that does not exist".

## Merge or undo makes open request reviews stale

Lines 625-630

Asked as: review out of date after merge; request review stale; merge changes request version; open review outdated

When a merge or undo changes a job request, Turbo adds 1 to its version and sets last_edited and updated to now. A review of that request that is open then becomes stale.

## Creating a job from a request moves pending undo references

Lines 636-675

Asked as: job from request undo; new job lost merge link; request job client restore; client note after job creation

When a job is created from a request, Turbo moves any pending client or contact absorb-undo references for that request to the new job. It moves them only if the job's field still holds the absorb target, so a relation the creator changed after the merge is not restored. If the client note made from the request's notes is under the target client, it moves too.
