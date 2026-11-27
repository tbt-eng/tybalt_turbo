<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/absorb_records.go","sha256":"7499cd7b5cc9aa163d460a6cb0742abaa5e73937114588c66030e96175c0a69f","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/absorb_records.go

## Absorb needs at least one record to absorb

Lines 51-53

Asked as: nothing selected to merge; no records chosen; no ids given; merge button does nothing

Turbo refuses an absorb request that lists no records to absorb. It shows "No IDs provided to absorb".

## A record cannot be absorbed into itself

Lines 55-60

Asked as: merge record into itself; target also selected to merge; same record twice merge

Turbo refuses an absorb when the target record is also in the list of records to absorb. It shows "Cannot absorb a record into itself".

## Absorb requires the collection's absorb claim

Lines 62-71

Asked as: can't merge records; no permission to merge; merge button not allowed; who can absorb

Only staff with the claims that the collection needs for absorbing (the absorb claim, or the claims the collection's config asks for instead) may absorb records. Anyone else sees "User does not have permission to absorb records".

## Absorbs that change jobs need job editing enabled

Lines 73-82

Asked as: merge blocked job editing off; job editing disabled merge; can't merge jobs

For collections whose absorb changes jobs (ChangesJobs), Turbo refuses the absorb while job editing is disabled. It shows "Job editing is disabled; cannot absorb records that modify jobs".

## Vendor absorb and undo depend on expenses editing

Lines 84-87

Asked as: merge vendors expenses editing; vendor merge blocked; undo vendor merge error

Both the absorb and the undo of a vendor absorb first run requireExpensesEditing. The result of that check is defined in another file.

## Target and absorbed records must exist

Lines 95-107

Asked as: merge missing record; record not found merge; target deleted merge

Turbo refuses the absorb if the target record does not exist ("Failed to find target record"). It also refuses if any record to absorb does not exist ("Failed to find record to absorb").

## Records must share the same parent value

Lines 108-119

Asked as: merge different parent values; can't merge different parents; merge records different project

When the collection has a parent constraint, every record to absorb must have the same value in that field as the target. Otherwise Turbo shows "Cannot absorb records with different <field> values".

## Kept client needs an alias if absorbed clients' profiles bill an alias

Lines 121-134

Asked as: merge clients alias invoicing; client alias needed merge; invoicing profile alias merge

When absorbing clients, if the kept client has no alias and any absorbed client has an invoicing profile that bills its alias, Turbo refuses. It shows "Invoicing profiles of a client being absorbed bill its alias. Give the kept client an alias first."

## Absorb that carries out a proposal is checked against it

Lines 136-144

Asked as: merge proposal not found; proposal already decided; approve merge proposal; proposal wrong record type

When the request names an absorb proposal, Turbo checks it first. The proposal must exist ("absorb proposal not found"). It must be for the same kind of record ("the absorb proposal is for another kind of record"). It must still be pending ("this absorb proposal has already been decided").

## Vendors with the same invoice number cannot be absorbed together

Lines 146-156

Asked as: merge vendors same invoice number; duplicate invoice number vendors; vendor invoice clash merge

When absorbing vendors, Turbo refuses if two of the vendors (target included) each have a vendor invoice with the same number. It shows "These vendors each have an invoice numbered <numbers>. Correct the invoice number on one of them before merging." Blank invoice number keys are ignored.

## Collection rules can refuse an absorb

Lines 158-167

Asked as: merge refused by collection rule; nothing to merge message; collection rule blocks merge

If a rule of the collection's config refuses the absorb, Turbo shows that rule's message. If the collection's prepare step leaves nothing to absorb, the message is "there is nothing to absorb".

## Absorb success message

Lines 169-171

Asked as: merge success message; records merged confirmation; how many merged

After a successful absorb Turbo says "Successfully absorbed N records into <target id>" and notifies open clients of the collection to refresh.

## Proposal absorb may drop or swap records but not add

Lines 175-208

Asked as: proposal added record merge; reviewer add record proposal; records not from proposal

The target and the records to absorb must all come from the proposal. A reviewer may leave records out or keep another one as target, but not add a record. Otherwise Turbo says "every record must come from the absorb proposal". Turbo then validates the proposal again, because the records may have changed.

## Absorb records who did it and the proposal

Lines 221-225

Asked as: who merged records; record merge user audit; merge history user

The absorb action stores the target, the proposal (if any) and the user who absorbed. The absorb route always passes the user.

## Absorb and undo run as one transaction

Lines 236-427

Asked as: merge failed nothing changed; merge rolled back; partial merge failed

The whole absorb runs in one transaction. If any step fails, nothing changes and Turbo reports "Failed to absorb records". After success, Turbo notifies open screens (the absorb_completed event) so they refresh. A failure of that notice does not undo the absorb.

## Absorb stops when the collection's rules leave nothing to absorb

Lines 238-252

Asked as: merge empty after preparation; nothing left to merge; merge prep removed all

Before absorbing, some collections run their own preparation, which can change the list of records to absorb. If the list is then empty, Turbo refuses with "there is nothing to absorb". The preparation runs inside the same transaction as the absorb.

## Only one pending absorb per collection

Lines 278-285

Asked as: can't merge another pending merge; earlier merge still pending; merge blocked by previous merge

Turbo cannot absorb records in a collection while an earlier absorb action for it exists. That action must be undone or deleted first.

## Absorb moves references to the target and deletes the records

Lines 287-406

Asked as: merge references repointed; merge deletes records; what happens when merging

Turbo repoints every reference to an absorbed record to the target, in one transaction. It saves the absorbed records, the rows they own and the old reference values in an absorb action so it can undo. Then it deletes the owned rows and the absorbed records.

## Absorb flags rewritten rows for legacy writeback

Lines 321-343

Asked as: merge imported flag legacy; legacy system not updated merge; _imported merge

For rewritten references in tables with an _imported column, Turbo sets _imported to false and updates the updated time. That sends the change to the legacy system.

## Absorb records the rewritten references for undo

Lines 357-377

Asked as: merge reference record undo; unused rate sheet merge; what references moved

Turbo records each reference it moves to the target (table, column, record and old value) in the absorb action. When no reference moves, such as a rate sheet name that nothing uses, it records each table and column with no rows instead. Undo uses this record.

## Absorb deletes rows owned by absorbed records

Lines 379-393

Asked as: merged record child rows deleted; owned rows merge deleted; merge removes child items

Turbo saves the rows that an absorbed record owns, in the absorb action, and then deletes them with the record. Undo can recreate them.

## Undo needs the absorb claim

Lines 445-456

Asked as: can't undo merge; no permission undo merge; undo button not allowed

Only staff with the collection's absorb claims may undo an absorb. Anyone else sees "User does not have permission to undo absorb".

## Undo needs an existing absorb action

Lines 458-468

Asked as: no merge to undo; undo nothing found; no absorb action

Turbo undoes the single pending absorb of a collection. If there is none, it shows "No absorb action found for collection".

## Child absorbs must be undone first

Lines 470-486

Asked as: can't undo merge dependent; undo blocked other merge; undo order merges

Turbo refuses to undo an absorb while an absorb exists in a collection that depends on it. It shows "Cannot undo X absorb while Y absorb exists; undo Y absorb first".

## Undo restores records, owned rows and references

Lines 491-601

Asked as: undo merge restore records; bring back merged records; reverse merge

Undo recreates the absorbed records with their original data and owned rows, without schema validation. It restores each saved reference to its old value and then deletes the absorb action, so a new absorb can follow. Undo also flags tables with an _imported column for legacy writeback. Turbo then says "Successfully undid absorb operation".

## Undo restores references and flags rows for legacy writeback

Lines 539-570

Asked as: undo merge legacy system; undo references legacy writeback; undo imported flag

Undo sets each rewritten reference back to its old value. For a table with an _imported column, it also sets _imported to false and updates the updated time, so the change is sent to the legacy system again.

## Undo notifies screens and reports success

Lines 590-601

Asked as: undo merge success message; undo failed nothing changed; undo notifies screens

After undo restores the records, Turbo notifies open screens that the absorb has completed. It then shows "Successfully undid absorb operation". If the undo fails, Turbo shows "Failed to undo absorb" and changes nothing.

## Absorb action stores target, records, references, proposal and user

Lines 625-645

Asked as: merge action saved record; merge log entry; absorb action save failed

Turbo saves one absorb action (absorb_actions) for each absorb. It holds the collection name, the target id, the absorbed records, the updated references, the proposal id and the id of the user who absorbed (absorbed_by). A save that fails stops the absorb.

## Absorb history saves each record's fields and owned rows

Lines 647-688

Asked as: merged records copy saved; restore deleted records snapshot; merge backup of fields

Before an absorb deletes records, Turbo saves the id and every schema field of each absorbed record in the absorb action (absorbed_records). If the collection's config lists owned rows, Turbo also saves the rows each record owns, grouped by table. Undo uses this copy to bring the records back.

## Undo recreates owned rows without validation

Lines 711-746

Asked as: undo owned rows recreate; restore child rows undo; undo unexpected format

When undo brings an absorbed record back, Turbo recreates each row that the record owned, table by table, after the record itself. Turbo saves these rows without schema validation, because they may predate the current rules. Turbo stops the undo if the saved owned rows have an unexpected format.

## Screens get an absorb_completed event after absorb or undo

Lines 748-785

Asked as: screens not refreshing after merge; live refresh merge; merge notification clients

After an absorb or an undo finishes, Turbo sends the event "<collection>/absorb_completed" to every connected client that subscribes to that topic. Those screens can then refresh their data. The event carries the action "absorb_completed" and the collection name.

## Absorb tracks old reference values per table, column and record

Lines 787-811

Asked as: old reference values merge; tracking references json; undo reference values saved

Turbo tracks every rewritten reference by table, column and record id, and keeps the old value. It saves this tracking in the absorb action as JSON (updated_references), so undo can put each old value back.

## Client links must stay consistent after absorb or undo

Lines 813-846

Asked as: client links broken after merge; job client mismatch merge; undo merge client link error

After an absorb or undo, Turbo checks every changed job, project request, contact and invoicing profile whose client relation was rewritten. It refuses the merge if one would link to another client. A changed invoicing profile must not bill a client alias that does not exist. Undo skips records deleted since the absorb.

## Absorb touches only requests whose client links changed

Lines 813-841

Asked as: merge validates client links only; rate sheet move merge check; merge checks jobs

After an absorb, Turbo validates client links only for records where a client, job owner or contact column was rewritten. A job moved to another rate sheet keeps its client links and is not checked again.

## Undo skips records deleted since the absorb

Lines 819-836

Asked as: undo skips deleted records; deleted job undo merge; undo missing records

When undo checks client links, it skips jobs, requests, contacts and invoicing profiles that were deleted since the absorb. Absorb itself does not skip them: a missing record is an error.

## Absorb checks the billing name of invoicing profiles

Lines 842-846

Asked as: invoicing profile alias error merge; billing name nonexistent client merge; profile billing alias merge fails

When an absorb or undo changes the client of an invoicing profile, Turbo checks that the profile does not bill a client alias that does not exist. If it does, the absorb fails with "merge would leave invoicing profile ... billing a client alias that does not exist".

## Absorb makes open project request reviews stale

Lines 847-852

Asked as: project request review stale merge; request version changed merge; open review outdated

When an absorb or undo rewrites a project request's references, Turbo raises its version by 1 and updates its last edited time. Any open review of the request becomes stale.

## Creating a job from a request keeps undo data

Lines 858-912

Asked as: request to job undo data; create job from request merge undo; job from request keeps merge

When a request becomes a job, Turbo moves the request's pending absorb-undo references to the job. It does this only if the job's field still holds the absorb target. The client note and the opening rate sheet period move with them.

## Vendor invoice clash check ignores blank numbers

Lines 914-937

Asked as: vendors invoice number clash ignored blank; duplicate invoice number check; blank invoice number clash

Two vendors clash when their expense documents have the same invoice number key and different vendors. Documents with a blank invoice number key are ignored. Turbo lists each clashing number once, in order.
