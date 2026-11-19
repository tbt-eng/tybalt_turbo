<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/absorb_records.go","sha256":"efde57f9ddbb62b7143a66cc49aea64006c96674adb0babe88841b3c6176ea1a","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/absorb_records.go

## Absorb needs a list of records to absorb

Lines 43-50

Asked as: no ids provided; absorb with nothing selected; invalid request body; empty absorb request

Turbo refuses an absorb request with no IDs and shows 'No IDs provided to absorb'. It also refuses a request body it cannot read ('Invalid request body').

## A record cannot be absorbed into itself

Lines 52-57

Asked as: absorb record into itself; merge with itself; same record twice; cannot absorb into itself

If the target record's ID is in the list of records to absorb, Turbo refuses and shows 'Cannot absorb a record into itself'.

## Absorbing requires the collection's absorb permission

Lines 59-68

Asked as: no permission to absorb; cannot merge records; absorb access denied; not allowed to merge

Turbo checks the user's claims with absorb.HasRequiredClaims for the collection. This is the absorb claim, or the claims the collection's config asks for instead. Without them, Turbo shows 'User does not have permission to absorb records'.

## Absorbs that change jobs are blocked when job editing is off

Lines 70-79

Asked as: job editing disabled; merge blocked jobs off; cannot absorb jobs; job editing turned off

For collections whose absorb config changes jobs (ChangesJobs), Turbo refuses while job editing is disabled. It shows 'Job editing is disabled; cannot absorb records that modify jobs'.

## Vendor absorb depends on expenses editing being enabled

Lines 81-84

Asked as: vendor merge expenses editing; expenses editing off vendors; merge vendors not working; vendor absorb blocked

Both absorb and undo of absorb call requireExpensesEditing for the collection. The check is meant for vendor absorbs. Its result and message come from another file.

## Target and absorbed records must exist

Lines 92-103

Asked as: failed to find target; record to absorb missing; merge target not found; deleted record merge

Turbo shows 'Failed to find target record' if the target is missing. It shows 'Failed to find record to absorb' if any listed record is missing.

## Absorbed records must share the target's parent value

Lines 105-115

Asked as: different parent values; merge different clients; records belong to different; parent field mismatch

Some collections have a parent constraint field. Every record to absorb must have the same value in that field as the target. Otherwise Turbo shows 'Cannot absorb records with different <field> values'.

## Absorbed clients' alias-billing profiles need the kept client to have an alias

Lines 118-131

Asked as: invoicing profile bills alias; client alias needed merge; merging clients alias error; give client alias first

When absorbing clients, if the kept client has no alias, Turbo checks each absorbed client. If any of their invoicing profiles bill the alias, it refuses. It shows 'Invoicing profiles of a client being absorbed bill its alias. Give the kept client an alias first.'

## Vendors with the same invoice number cannot be absorbed together

Lines 133-143

Asked as: duplicate invoice number vendors; same vendor invoice number; vendors invoice clash; merge vendors invoice numbers

Absorbing vendors is refused if two or more of the vendors involved, the target included, each have an expense document with the same vendor invoice number. Blank numbers do not count. Turbo shows 'These vendors each have an invoice numbered <numbers>. Correct the invoice number on one of them before merging.'

## Collection rules can refuse an absorb

Lines 145-154

Asked as: merge refused by rule; absorb failed message; collection rule blocked merge; merge failed unknown reason

If a rule in the collection's config refuses the absorb, Turbo shows that rule's message. Other failures show 'Failed to absorb records'.

## Successful absorb reports the count

Lines 156-158

Asked as: absorb success count; how many records merged; merged records confirmation; successfully absorbed records

A successful absorb replies 'Successfully absorbed N records into <target id>', where N is the number of IDs in the request.

## Absorb stores no data in a temp table for refused absorbs

Lines 173-175

Asked as: merge rolled back; absorb partly done; nothing changed after failed merge; failed merge leftovers

Absorb runs in one transaction. If any step fails, Turbo rolls back everything, and no record is deleted and no reference moves.

## Collection's prepare step can add IDs or find nothing to absorb

Lines 175-189

Asked as: nothing to absorb; merge list empty; prepare step empty list; no records left to merge

If the collection's config has a prepare step, it checks the request first inside the transaction. It may change the list of IDs. If the list ends up empty, Turbo refuses with 'there is nothing to absorb'.

## Only one pending absorb per collection

Lines 215-222

Asked as: already a pending merge; second merge blocked; merge already in progress; undo previous merge first

Turbo refuses a new absorb while the collection already has an absorb action. The earlier absorb must be undone or deleted first.

## Absorb moves references to the target and deletes the absorbed records

Lines 224-343

Asked as: what merge does; references repointed; merge deletes absorbed records; merge saves undo data

Turbo repoints every reference to the absorbed records so it names the target. It saves the absorbed records and the old reference values in an absorb action so the absorb can be undone. It then deletes the rows that the absorbed records own, and the absorbed records. All of this runs in one transaction.

## Absorb marks changed rows for legacy writeback

Lines 257-271

Asked as: imported flag merge; legacy writeback merge; merge sent to legacy system; imported column reset

For a table with an _imported column, the reference update sets _imported to false and refreshes updated. This makes the change get written back to the legacy system.

## Absorb is refused if it leaves linked records on another client

Lines 290-292

Asked as: linked records other client; merge leaves mismatched client; job client mismatch after merge; merge fails client link

After the references move, Turbo checks every changed job, job request, client contact and invoicing profile. If a client relation no longer agrees, the absorb fails. Invoicing profiles whose client changed must still bill an alias that exists.

## Absorb that moves no references still records every table and column

Lines 294-304

Asked as: merge moved no references; empty absorb action; nothing repointed merge; blank merge record

If no reference moved, Turbo saves the absorb action with each configured table and column set to empty. The saved action cannot be blank.

## Failed refresh notice does not undo a completed absorb

Lines 358-363

Asked as: refresh event failed; realtime update error after merge; screen not refreshed after merge; absorb completed event error

After an absorb commits, Turbo sends a realtime absorb_completed event so open screens can refresh. If sending that event fails, Turbo only logs the error. The absorb stays done and the person still gets the success result.

## Undo absorb is blocked when vendor expenses editing is off

Lines 377-380

Asked as: undo merge vendors blocked; expenses editing undo; cannot undo vendor merge; vendor undo expenses off

Undoing an absorb on the vendors collection checks first that expenses editing is enabled. Turbo refuses the undo when it is not. The check (requireExpensesEditing) runs before the permission check.

## Undo of absorb needs the same permission

Lines 382-393

Asked as: no permission undo merge; cannot undo absorb; undo access denied; not allowed to undo merge

Undoing an absorb needs the same claims as absorbing. Without them, Turbo shows 'User does not have permission to undo absorb'.

## Undo needs an existing absorb action

Lines 395-405

Asked as: no merge to undo; nothing to undo; no absorb found; undo missing action

Undo looks for the collection's single absorb action. If none exists, Turbo shows 'No absorb action found for collection'.

## Child collection absorbs must be undone first

Lines 407-423

Asked as: undo order merge; later merge must be undone; dependent merge blocks undo; undo previous merge first order

Turbo refuses to undo an absorb while an absorb exists in a collection that depends on it. It shows 'Cannot undo X absorb while Y absorb exists; undo Y absorb first'. Absorbs must be undone in last-in, first-out order.

## Undo restores records, owned rows and references

Lines 428-525

Asked as: restore merged records; bring back absorbed records; undo merge restores; records come back

Undo recreates the absorbed records with their original data, without schema validation. It also recreates their owned rows, then points the references back to their old values, and deletes the absorb action. It all runs in one transaction. It reads the latest saved undo data, so references moved by job creation are kept.

## Undo marks restored rows for legacy writeback

Lines 476-489

Asked as: undo imported flag; legacy system undo merge; restored rows legacy sync; undo writeback legacy

On undo, a table with an _imported column gets _imported set to false and updated refreshed. The legacy system receives the change again.

## Undo skips deleted records but keeps client checks

Lines 509-511

Asked as: undo deleted records skipped; undo client check; undo alias check; deleted since merge

On undo, Turbo skips records deleted since the absorb. For the rest, it re-checks client links and invoicing profile alias names.

## Undo success message

Lines 531-538

Asked as: undo success message; merge undone confirmation; undo worked; undo completed reply

A successful undo replies 'Successfully undid absorb operation'. Turbo then notifies subscribed clients so they refresh.

## Absorb action is recorded with the target and the changes

Lines 562-580

Asked as: absorb history record; what changed in merge; merge log; merge audit details

Each absorb creates one absorb action (absorb_actions) that holds the collection name, the target ID, the absorbed records and the references that changed.

## Absorb saves the absorbed records and their owned rows for undo

Lines 582-623

Asked as: what gets saved before merge; saved for undo; merge backup data; owned rows saved

Before it deletes anything, Turbo saves each absorbed record and the rows it owns (OwnedRows) in the absorb action. Undo uses this data to recreate them.

## Undo recreates records without validation

Lines 646-681

Asked as: undo fails validation; old records restore error; restore missing required field; undo old records invalid

Undo recreates absorbed records and their owned rows with their original data and skips current validation. Records that predate newer required fields can therefore be restored.

## Absorb finished sends a realtime refresh event

Lines 683-720

Asked as: screens not refreshing merge; realtime merge event; merge broadcast failed; other users see merge

When an absorb or an undo succeeds, Turbo sends a realtime event "<collection>/absorb_completed" to the clients that subscribe to it, so they can refresh. If the broadcast fails, the absorb still succeeds.

## Only clients, jobs, requests, contacts and invoicing profiles get client checks

Lines 754-776

Asked as: client check after merge; rate sheet change check; why merge checked clients; contacts client mismatch

After an absorb or an undo, Turbo checks the client relations of every job, project request, client contact and invoicing profile that changed. It checks only if a client, job owner or contact column was rewritten. A rewrite of another column, such as a rate sheet, does not trigger the check.

## Absorb or undo is refused if an invoicing profile bills a missing alias

Lines 777-781

Asked as: invoicing profile alias missing; billing name alias error; merge alias does not exist; profile billing wrong alias

When a client relation of an invoicing profile changes, Turbo checks that the profile's billing name still matches an alias of its client. If not, the absorb or undo fails with "merge would leave invoicing profile ... billing a client alias that does not exist".

## Changed job requests become new versions

Lines 782-787

Asked as: job request version changed; request edited after merge; review stale after merge; request version bump

When an absorb or undo changes the client-related references of a job request, Turbo raises that request's version by 1 and updates its last edited time. Any open review of the request becomes stale.

## Creating a job from a request keeps absorb undo data correct

Lines 793-847

Asked as: job from request undo data; creating job breaks undo; rate sheet periods moved; client note from request

When a job is created from a request, pending absorb actions for clients, contacts and rate sheets are updated. Undo data for the request's references moves to the new job, but only if the job still holds the absorb target. A client note made from the request notes and the job's rate sheet periods move back with the job on undo.

## Clashing vendor invoice numbers are found by number key

Lines 849-872

Asked as: which vendor invoice numbers clash; duplicate invoice numbers listed; blank invoice number merge; invoice number check

The vendor absorb check looks at the vendor invoice numbers of the target and the absorbed vendors. It reports a number only when more than one vendor has it. It ignores blank numbers and lists the numbers in order.
