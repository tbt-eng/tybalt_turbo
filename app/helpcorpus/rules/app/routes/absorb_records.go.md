<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/absorb_records.go","sha256":"ef7db93ae9b1e8a82949ae5d29d9513c87824f3198ec655c510bb7387159ed12","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/absorb_records.go

## Absorb request with an unreadable body is refused

Lines 43-46

When Turbo cannot read the absorb request body, it refuses the absorb with the message "Invalid request body". Turbo makes this check before it checks permissions or looks at any record.

## Absorb needs a list of records to absorb

Lines 48-50

Turbo refuses an absorb request whose list of IDs to absorb (ids_to_absorb) is empty. It shows the message "No IDs provided to absorb".

## A record cannot be absorbed into itself

Lines 52-57

Turbo refuses an absorb when the kept record's ID is also in the list to absorb. It shows the message "Cannot absorb a record into itself".

## Only staff with the absorb claim can absorb

Lines 59-67

Only a person with the absorb claim can absorb records. Anyone else sees "User does not have permission to absorb records".

## Absorbing clients or contacts needs job editing on

Lines 69-78

For clients and client contacts, Turbo refuses an absorb while job editing is disabled. It shows "Job editing is disabled; cannot absorb records that modify jobs", because these absorbs change jobs.

## Vendor absorb and undo need expense editing on

Lines 80-83

Both absorb and undo-absorb call requireExpensesEditing for the collection, which applies to vendor absorbs. This file does not show the exact condition or message.

## Absorb target record must exist

Lines 91-95

Turbo refuses an absorb when the target record (the id in the path) is not found in the collection. The message is "Failed to find target record".

## Every record to absorb must exist

Lines 97-102

Turbo refuses an absorb when any record in the list to absorb is not found in the collection. The message is "Failed to find record to absorb".

## Absorbed records must share the target's parent value

Lines 104-114

When the collection has a parent constraint, every record to absorb must have the same value in the parent field as the target. Otherwise Turbo refuses with "Cannot absorb records with different <field> values". The field name comes from the constraint (absorb.GetRefConfigs).

## Kept client needs an alias if absorbed clients' profiles bill an alias

Lines 117-130

When absorbing clients, if the kept client has no alias, Turbo checks each client being absorbed for invoicing profiles that bill its alias. If any exist, Turbo refuses with "Invoicing profiles of a client being absorbed bill its alias. Give the kept client an alias first."

## Absorb success message

Lines 138-140

A successful absorb replies "Successfully absorbed N records into <target id>", where N is the number of ids in the request.

## Only one pending absorb per collection

Lines 180-187

Turbo refuses an absorb when the collection already has a recorded absorb action. The existing absorb must be undone or deleted first. The error says "there is an existing absorb action for collection <name> that must be undone or deleted first".

## Absorb moves all references to the target and deletes the absorbed records

Lines 195-289

An absorb points every reference to an absorbed record (in the tables from absorb.GetRefConfigs) at the target record, then deletes the absorbed records. It saves the old records and old reference values as an absorb action so the absorb can be undone. All of this happens in one transaction, so any failure rolls back the whole absorb.

## Absorb marks changed legacy-synced rows for writeback

Lines 222-245

For each updated reference row in a table with an _imported column, Turbo sets _imported to false and sets updated to the current time. This makes the change write back to the legacy system, because the direct update skips the normal hooks.

## Merge refuses to leave records linked to another client

Lines 255-257

After the references move, Turbo re-validates every changed job, job request, client contact and invoicing profile. If its client relations no longer agree, the absorb fails with "merge would leave <table> <id> linked to another client". An invoicing profile that would bill a client alias that does not exist also fails the absorb.

## Absorb completion notifies other users

Lines 295-298

After an absorb or an undo completes, Turbo sends an "<collection>/absorb_completed" realtime event to subscribed clients so they can refresh. A failure to send is logged and does not fail the absorb.

## Only staff with the absorb claim can undo an absorb

Lines 319-329

Undoing an absorb requires the absorb claim. Without it, Turbo refuses with "User does not have permission to undo absorb".

## Undo needs an existing absorb action

Lines 331-341

Turbo refuses an undo when the collection has no recorded absorb action. The message is "No absorb action found for collection".

## Child absorbs must be undone before parent absorbs

Lines 343-359

Turbo refuses to undo an absorb while a dependent child collection (from absorb.GetUndoBlockers) still has an absorb action. The message is "Cannot undo <collection> absorb while <child> absorb exists; undo <child> absorb first".

## Undo restores records and references

Lines 364-454

An undo recreates the absorbed records with their saved data, without schema validation, because old records may not meet current required fields. It then sets each changed reference back to its old value and deletes the absorb action, so a new absorb is allowed. It runs in one transaction and rolls back on any failure. It replies "Successfully undid absorb operation".

## Undo uses latest saved undo data

Lines 365-378

The undo reads the absorb action inside its own transaction, so references transferred to a job by project creation are not lost.

## Undo skips deleted records when validating clients

Lines 425-427

After the restore, an undo re-validates client relations on jobs, requests, contacts and invoicing profiles. It skips records that were deleted since the merge. A record that fails validation fails the undo.

## Absorb saves full copies of absorbed records for undo

Lines 498-530

Before an absorb deletes records, Turbo saves each absorbed record's original id and the value of every field in the collection's schema. Turbo stores these copies in the absorb action (absorbed_records). Undo uses them to recreate the records.

## Absorb completion event goes only to subscribed users

Lines 532-569

After an absorb or an undo, Turbo sends the event "<collection>/absorb_completed" to connected users. Only users who subscribe to that topic for the collection receive it, and they can then refresh their data.

## Absorb remembers each changed reference's old value

Lines 571-595

For every reference that an absorb moves, Turbo remembers the table, the column, the record id and the old value. Turbo saves this in the absorb action (updated_references). Undo needs it to point each reference back at the absorbed record.

## Merge checks four record types for client links that disagree

Lines 597-619

After a merge or an undo, Turbo checks every job, project request (job_requests), client contact and invoicing profile that the operation changed. The client relations on each record must still agree. If they do not, Turbo refuses with "merge would leave <table> <id> linked to another client" and the reason.

## Merge refuses an invoicing profile billing a missing alias

Lines 620-624

After a merge or an undo, Turbo checks each changed invoicing profile (client_invoicing_information) for its billing name. If the profile would bill a client alias that does not exist, Turbo refuses with "merge would leave invoicing profile <id> billing a client alias that does not exist".

## Merge or undo makes open request reviews stale

Lines 625-630

When a merge or undo changes a job request, Turbo adds 1 to its version and sets last_edited and updated to now. A review of that request that is open then becomes stale.

## Creating a job from a request moves pending undo references

Lines 636-675

When a job is created from a request, Turbo moves any pending client or contact absorb-undo references for that request to the new job. It moves them only if the job's field still holds the absorb target, so a relation the creator changed after the merge is not restored. If the client note made from the request's notes is under the target client, it moves too.
