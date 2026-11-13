<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/absorb_records.go","sha256":"d13db50e03c0b4fd241bbade65de87e0dcb11e18dc81d014a5e0d42cc1f68736","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/absorb_records.go

## Absorb needs at least one record to absorb

Lines 48-50

Asked as: absorb nothing selected; no records selected; empty absorb list; absorb button does nothing

Turbo refuses an absorb request that lists no records to absorb. The message is "No IDs provided to absorb".

## A record cannot be absorbed into itself

Lines 52-57

Asked as: merge record into itself; target also selected; absorb same record; merge with itself

Turbo refuses an absorb when the target record is also in the list of records to absorb. The message is "Cannot absorb a record into itself".

## Absorb and undo absorb need the collection's absorb permission

Lines 59-68

Asked as: can't absorb records; no permission to merge; cannot undo absorb access; absorb access denied

Only a person with the absorb claim, or the claims that the collection's config asks for instead (absorb.HasRequiredClaims), can absorb records. Anyone else gets "User does not have permission to absorb records". The same claims allow undoing an absorb, and others get "User does not have permission to undo absorb".

## Absorbs that change jobs need job editing enabled

Lines 70-79

Asked as: absorb fails job editing off; jobs locked cannot merge; job editing turned off; merge blocked jobs disabled

When the collection's config marks that absorbing changes jobs (ChangesJobs), Turbo checks that job editing is enabled. If it is off, Turbo refuses with "Job editing is disabled; cannot absorb records that modify jobs".

## Vendor absorb and undo need expense editing enabled

Lines 81-84

Asked as: vendor merge not working; expense editing off merge; cannot merge vendors; vendor absorb blocked

Before an absorb or an undo of an absorb, Turbo calls requireExpensesEditing for the collection. The comment says this applies to vendor absorbs. This file does not show the refusal message. On undo, this check runs before the permission check.

## Absorb target record must exist

Lines 92-96

Asked as: merge target not found; kept record missing; target does not exist; absorb target deleted

When an absorb request names a target record that does not exist in the collection, Turbo refuses with 'Failed to find target record'.

## Every record to absorb must exist

Lines 98-103

Asked as: selected record gone; merge record not found; duplicate already deleted; record to merge missing

If any record in the list to absorb does not exist in the collection, Turbo refuses with 'Failed to find record to absorb'.

## Absorbed records must share the target's parent value

Lines 105-115

Asked as: merge different parent field; parent mismatch absorb; cannot merge different values; records have different parent

Some collections have a parent constraint field. Turbo refuses an absorb when a record's value in that field differs from the target's, with 'Cannot absorb records with different <field> values'.

## Absorbing clients needs an alias on the kept client

Lines 118-131

Asked as: merge clients alias error; client alias required; invoicing profile alias merge; merge client needs alias

When the kept client has no alias, Turbo refuses if any client being absorbed has invoicing profiles that bill its alias. The message is 'Invoicing profiles of a client being absorbed bill its alias. Give the kept client an alias first.'

## Collection rules can refuse an absorb

Lines 134-139

Asked as: absorb blocked by rule; merge refused custom rule; collection rule stopped merge; absorb refused with message

If a rule of the collection's config refuses the absorb, Turbo shows that rule's own message. Such rules can also add ids to absorb. If nothing is left to absorb, the message is 'there is nothing to absorb'.

## Successful absorb confirmation message

Lines 144-146

Asked as: absorb worked message; how many records merged; absorb confirmation; merge success message

After a successful absorb, Turbo replies 'Successfully absorbed N records into <target id>', where N is the number of ids requested.

## Absorb moves references to the target and deletes the absorbed records

Lines 160-331

Asked as: what happens to references; merged records deleted; absorb moves links; absorb transaction details

Turbo changes every reference to an absorbed record so it points to the target. It then deletes the absorbed records and the rows they own. All of this runs in one transaction, so a failure changes nothing.

## Absorb with no usable target or ids to absorb is refused

Lines 163-177

Asked as: says nothing to absorb; no ids after prepare; merge does nothing; absorb nothing left

When a collection's own rules (its Prepare step) leave nothing to absorb, Turbo refuses the absorb with the message 'there is nothing to absorb'. The collection's rules run first and can add ids. They see the data inside the same transaction that absorbs.

## Only one pending absorb per collection

Lines 203-210

Asked as: cannot absorb second time; already pending merge; previous absorb still pending; second merge blocked

Turbo refuses a new absorb in a collection while an earlier absorb action for it exists. The earlier one must be undone or deleted first.

## Absorb saves data so it can be undone

Lines 212-302

Asked as: absorb backup data; where is undo data; undo information stored; saved before merge

Before it deletes anything, Turbo saves the absorbed records, the rows they own and every reference it changed in an absorb action. Undo uses this record. When no reference moves, Turbo records each table and column with no rows.

## Absorb marks changed records for legacy writeback

Lines 246-259

Asked as: sync to legacy system; imported flag absorb; legacy not updated after merge; writeback after absorb

For tables with an _imported column, absorb sets _imported to false and updates the 'updated' time on each changed record. The legacy system then receives the change.

## Absorb with no moved references still records every table and column

Lines 282-292

Asked as: absorb with no references; unused rate sheet merge; merge nothing linked; no references moved

When an absorb moves no reference, Turbo still saves the absorb action. It records each table and column with no rows, because the saved references cannot be blank. An example is a rate sheet name that nothing uses.

## Absorb deletes rows owned by absorbed records

Lines 304-318

Asked as: child rows removed on merge; owned rows deleted absorb; children deleted when merging; merge deletes child rows

Turbo deletes the rows that the absorbed records own, such as child rows in an owned table, before it deletes the records. It saves those rows in the absorb action first, so undo can recreate them.

## Absorb and undo notify open screens

Lines 346-349

Asked as: screen not refreshing after merge; absorb not showing; need reload after absorb; open screens not updated

When an absorb or undo finishes, Turbo sends an absorb_completed event to clients subscribed to that collection so they refresh. A failure to send it does not fail the operation.

## Undo absorb needs an existing absorb action

Lines 383-393

Asked as: undo absorb not found; no absorb to undo; undo does nothing; nothing to undo merge

Undo refuses with 'No absorb action found for collection' when the collection has no pending absorb.

## Undo must happen in reverse order of dependent absorbs

Lines 395-411

Asked as: cannot undo child absorb exists; undo parent absorb blocked; undo order merge; dependent absorb undo

Turbo refuses to undo an absorb while an absorb exists in a dependent child collection. The message is 'Cannot undo X absorb while Y absorb exists; undo Y absorb first'.

## Failed absorb or undo changes nothing

Lines 413-517

Asked as: absorb failed partway; undo failed nothing changed; rollback absorb; merge half done

Absorb and undo each run in one transaction. If any step fails, Turbo rolls back everything. Undo then shows the message 'Failed to undo absorb'.

## Undo recreates absorbed records and restores references

Lines 416-513

Asked as: restore merged records; bring back absorbed clients; undo restores deleted data; recreate absorbed records

Undo recreates each absorbed record with its original id and data, and the rows it owned. It sets the changed references back to their old values. Then it deletes the absorb action so a new absorb is possible. Recreated records skip validation because they may predate current rules.

## Undo restores references from the latest saved undo data

Lines 416-430

Asked as: undo after job created; undo lost references; project request undo references; undo uses latest data

Undo reads the absorb action again inside its transaction. It uses the latest saved records and references, so references that job creation moved from a project request are not lost.

## Undo recreates owned rows after the absorbed record

Lines 440-458

Asked as: undo missing child rows; restore children after undo; owned rows back on undo; undo order of rows

Undo recreates each absorbed record first and then the rows it owned. Turbo skips validation when it saves them, so records and rows that predate current rules come back as they were.

## Undo marks changed records for legacy writeback

Lines 465-477

Asked as: undo not synced legacy; restored references imported flag; legacy writeback undo; undo sets imported

For tables with an _imported column, undo sets _imported to false and updates the 'updated' time on restored references.

## Undo deletes the absorb action

Lines 501-510

Asked as: can absorb again after undo; clear undo record; undo removes absorb record; reset absorb after undo

After a successful undo, Turbo deletes the absorb action for the collection. This allows a new absorb on that collection.

## Undo confirmation message

Lines 524-526

Asked as: undo worked message; undid absorb confirmation; undo success message; absorb reversed message

After a successful undo, Turbo replies 'Successfully undid absorb operation'.

## Absorb notifies open screens after it completes

Lines 671-708

Asked as: screen refresh after merge; absorb event send failed; merge notification not sent; screens not updating absorb

After an absorb or an undo, Turbo sends an absorb_completed event for the collection to connected users who subscribe to it, so their screens refresh. A failure to send is logged and does not fail the absorb.

## Absorb is refused if it leaves records linked to different clients

Lines 736-769

Asked as: jobs linked to wrong client; absorb moved contacts wrong client; invoicing profile client mismatch; merge links another client

After references move, Turbo checks the jobs, project requests, client contacts and invoicing profiles it changed. If a client relation would link to another client, the absorb fails. An invoicing profile that would bill a client alias that does not exist also fails it.

## Undo skips records deleted since the absorb

Lines 742-759

Asked as: undo error deleted record; client relation check undo; undo fails deleted records; undo skips missing records

When Turbo checks client relations after an undo, it skips records that were deleted after the absorb.

## Absorb makes open project request reviews stale

Lines 770-775

Asked as: request review out of date; request version changed merge; reviewer conflict after absorb; project request stale

When an absorb or undo changes a project request's relations, Turbo raises that request's version by 1. It also sets its last edited and updated times to now, so any open review of it is stale.

## Creating a job from a request keeps pending undo data

Lines 781-826

Asked as: create job from request undo; project request job undo data; job creation undo references lost; request entries moved to job

When a job is created from a project request, Turbo moves the request's entries in pending absorb undo data to the new job. This covers client, contact and rate sheet absorbs. Turbo does not move an entry if the creator changed that relation after the merge.

## Request notes and opening rate sheet follow the job on undo

Lines 806-823

Asked as: client note undo; opening rate sheet undo; notes moved back on undo; request notes follow job

For a moved client reference, Turbo also tracks the client note made from the request's notes. For a moved rate sheet, it tracks the job's opening rate sheet period rows. Undo then moves them back with the job.
