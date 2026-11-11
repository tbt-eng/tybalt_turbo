<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/absorb_records.go","sha256":"1f62115e908ba9d8bf38660c31f25687e4a40eacb8154461f0c4b8533e50dbc7","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/absorb_records.go

## Absorb request needs at least one record ID

Lines 43-50

Asked as: no IDs to absorb; absorb nothing selected; empty absorb list; absorb request empty; invalid request body absorb

Turbo refuses an absorb request with an empty list of IDs to absorb (ids_to_absorb). The message is "No IDs provided to absorb". A body that cannot be read gets "Invalid request body".

## A record cannot be absorbed into itself

Lines 52-57

Asked as: absorb record into itself; merge same record; target also in absorb list; cannot merge with itself; absorb yourself error

Turbo refuses an absorb when the target record's ID is also in the list of IDs to absorb. The message is "Cannot absorb a record into itself". This check runs before the permission check.

## Only holders of the absorb claim can absorb records

Lines 59-67

Asked as: no permission to absorb; can't merge records; absorb button missing; merge access denied; need absorb permission

Turbo refuses an absorb unless the person has the absorb claim. The message is "User does not have permission to absorb records".

## Absorbing clients or client contacts needs job editing on

Lines 69-78

Asked as: absorb clients job editing off; merge clients disabled; cannot merge contacts; job editing disabled merge; client merge blocked

When job editing is disabled, Turbo refuses to absorb clients or client contacts. The message is "Job editing is disabled; cannot absorb records that modify jobs".

## Vendor absorbs depend on the expenses editing setting

Lines 80-83

Asked as: merge vendors expenses setting; vendor absorb blocked; expenses editing disabled vendor; cannot merge vendors; vendor merge not allowed

Absorb and undo-absorb both call requireExpensesEditing for the collection, which guards vendor absorbs. This file does not show the rule it applies. For undo, this check runs before the absorb claim check.

## Absorb target and records must exist

Lines 91-102

Asked as: failed to find target record; record to merge missing; merge record not found; deleted record absorb; target record gone

Turbo refuses an absorb if the target record is missing, with the message "Failed to find target record". It refuses if any record to absorb is missing, with the message "Failed to find record to absorb".

## Absorbed records must share the same parent

Lines 104-114

Asked as: merge different parent; records different job; absorb different client field; merge mismatched parent; different values merge error

For collections with a parent constraint (from absorb.GetRefConfigs), the target and every record to absorb must have the same value in the parent field. Otherwise Turbo refuses with "Cannot absorb records with different <field> values".

## Absorbing clients needs the kept client to have an alias

Lines 117-130

Asked as: client alias required merge; kept client no alias; invoicing profile alias merge; merge clients alias error; give client alias first

When the kept client has no alias, Turbo checks each client being absorbed for invoicing profiles that bill its alias. If any exist, Turbo refuses with "Invoicing profiles of a client being absorbed bill its alias. Give the kept client an alias first."

## Vendors with the same invoice number cannot merge

Lines 132-142

Asked as: vendors same invoice number; duplicate invoice number merge; merge vendors invoice clash; invoice number conflict vendors; same invoice numbered vendors

Turbo refuses to absorb vendors when two or more of the vendors involved (target and absorbed) each have an invoice with the same number. The message is "These vendors each have an invoice numbered <numbers>. Correct the invoice number on one of them before merging." Invoices with an empty number key are ignored.

## Absorb success message and refresh

Lines 150-152

Asked as: merge success message; absorbed records successfully; merge completed refresh; how many records merged; merge done confirmation

A successful absorb returns "Successfully absorbed N records into <target id>". Turbo also tells connected clients to refresh the collection.

## Only one absorb per collection until undone

Lines 192-199

Asked as: already merged collection; second merge blocked; existing absorb must undo; can't merge again; merge already exists

Turbo refuses a new absorb in a collection while an earlier absorb for that collection still exists. It must be undone or deleted first. The error says there is an existing absorb action that must be undone or deleted first.

## Absorb repoints references and deletes absorbed records

Lines 201-301

Asked as: merge references repointed; what happens when merging; merged records deleted; merge transaction rollback; references moved to target

An absorb changes every reference to an absorbed record so it points to the target. Then Turbo deletes the absorbed records. All of this happens in one transaction, so a failure leaves nothing changed. Turbo saves the old records and old references so the absorb can be undone.

## Absorb marks imported rows for legacy writeback

Lines 235-248

Asked as: merge legacy writeback; imported rows merge; merge not syncing legacy; _imported reset merge; merge sync old system

When an absorb changes references in a table that has an _imported column, Turbo sets _imported to false and sets updated to the current time. This makes the change write back to the legacy system.

## Absorb is refused if it links a record to another client

Lines 267-269

Asked as: merge links client another; merge leaves linked to another client; merge client conflict; job linked to other client; merge would break client link

After repointing references, Turbo checks every changed job, project request, client contact and invoicing profile. If its client relations no longer agree, the absorb fails with "merge would leave <table> <id> linked to another client". An invoicing profile that would bill a client alias that does not exist also fails the absorb.

## Undo absorb needs the absorb claim

Lines 326-341

Asked as: undo absorb permission; cannot undo merge; undo merge access denied; need permission undo merge; who can undo absorb

Only holders of the absorb claim can undo an absorb. Others get "User does not have permission to undo absorb". Vendor undo also depends on the expenses editing setting.

## Undo needs an existing absorb

Lines 343-353

Asked as: no absorb to undo; undo merge nothing found; no merge action found; undo without merge; nothing to undo collection

Undo refuses with "No absorb action found for collection" if the collection has no absorb to undo. Each collection can have only one absorb at a time.

## Undo child absorbs before parent absorbs

Lines 355-371

Asked as: undo order child parent; cannot undo parent merge; undo child merge first; undo merge dependent error; reverse undo order merge

Turbo refuses to undo an absorb while an absorb in a dependent child collection exists. The message is "Cannot undo X absorb while Y absorb exists; undo Y absorb first". Undos must run in reverse order.

## Undo restores records and references

Lines 373-453

Asked as: undo merge restores records; bring back merged records; revert merge; restore deleted merged records; undo merge references

Undo recreates the absorbed records with their saved data and original IDs, skipping schema validation. It puts every changed reference back to its old value, then deletes the absorb record so a new absorb is possible. It reads the latest saved undo data first, so references moved by job creation are kept. If any step fails, Turbo rolls back the whole undo.

## Undo skips deleted records when checking clients

Lines 437-439

Asked as: undo merge client check; undo merge client error; undo leaves client linked; undo skipped deleted records; undo merge client conflict

After an undo, Turbo checks client relations on the restored records. It skips records deleted since the absorb. An undo that would leave a record linked to another client fails.

## Undo success message

Lines 459-466

Asked as: undo success message; merge undone confirmation; undo completed message; undid merge; undo finished refresh

A successful undo returns "Successfully undid absorb operation". Turbo tells connected clients to refresh.

## Absorb action is saved with its target and old references

Lines 491-508

Asked as: merge action saved; what is saved for merge; merge history record; absorb action stored; old references saved merge

Each absorb saves one absorb action. It holds the collection name, the kept record's id (target_id), the absorbed records and the old reference values (updated_references).

## Absorb stores full copies of absorbed records for undo

Lines 511-542

Asked as: merge copies records backup; saved copy before merge; backup of merged records; merged record snapshot; keep deleted records undo

Before an absorb, Turbo saves every absorbed record with its original id and all of its schema fields in the absorb action (absorbed_records). Undo uses this copy to bring the records back.

## Absorb and undo refresh open screens

Lines 548-581

Asked as: screens not refreshing merge; merge event update; other users see merge; open screen refresh after merge; absorb completed event

When an absorb or an undo finishes, Turbo sends an "<collection>/absorb_completed" event. Only clients subscribed to that topic receive it. These clients can then refresh their data.

## Merge checks that jobs, requests, contacts and profiles keep one client

Lines 613-631

Asked as: merge client relations mismatch; merge job client mismatch; merge contact client conflict; after merge client check; merge would link another client

After a merge or an undo, Turbo checks every job, project request (job_requests), client contact and invoicing profile that it changed. Their client relations must still agree. If one does not, Turbo refuses with "merge would leave <table> <id> linked to another client" plus the reason.

## Merge refused if an invoicing profile would bill a missing alias

Lines 632-636

Asked as: invoicing profile missing alias merge; profile bills missing alias; merge invoicing alias error; alias does not exist merge; billing profile alias merge fails

After a merge or an undo, Turbo checks each changed invoicing profile (client_invoicing_information). If it would bill a client alias that does not exist, Turbo refuses with "merge would leave invoicing profile <id> billing a client alias that does not exist".

## Absorb invalidates open reviews of changed project requests

Lines 637-642

Asked as: merge stale review; project request review outdated; request version merge; merge changes open review; review stale after merge

When an absorb or undo changes a project request's relations, Turbo increases its version by 1. It also sets last_edited and updated to now. This makes any open review of that request stale.

## Creating a job from a request keeps merge undo working

Lines 648-688

Asked as: job from request undo merge; create job merge undo lost; request notes merge undo; undo merge after job created; job creation breaks undo

When a job is created from a request, Turbo moves the request's saved merge-undo references for clients and contacts to the new job. It does this only when the job's field still equals the merge target. If the request's notes became a client note under the merged client, undo moves that note back too.

## Vendor invoice number clash check across merged vendors

Lines 690-713

Asked as: vendor invoice clash merge; duplicate vendor invoice merge; invoice number used twice vendors; which invoice numbers clash; merge invoice number check

Turbo looks at the kept vendor and all vendors being absorbed. It finds invoice numbers (vendor_invoice_number_key) that belong to more than one of them. Expense documents with a blank invoice number key are ignored. A number used twice under the same vendor is not a clash. The clashing numbers come back sorted.
