<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/absorb_records.go","sha256":"6a33ee89add215581cf189fe579d3095c5e1bf9c4718c19005d38d0594924fb4","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/absorb_records.go

## Absorbing records with an unfound request body is refused

Lines 43-46

Asked as: invalid request body; can't read absorb data; merge form not read; absorb request body error

Turbo refuses an absorb when the request body cannot be read. It shows the message "Invalid request body".

## Absorb needs at least one record ID

Lines 48-50

Asked as: merge with no records; nothing selected to merge; no ids to merge; absorb button no selection

Turbo refuses an absorb request with no IDs to absorb. It shows "No IDs provided to absorb".

## A record cannot be absorbed into itself

Lines 52-57

Asked as: merge record into itself; target also selected; absorb same record twice; merge target in selection

If the target record's ID is among the IDs to absorb, Turbo refuses with "Cannot absorb a record into itself".

## Absorbing requires the collection's absorb permission

Lines 59-68

Asked as: can't merge records; no permission to merge; merge access denied; absorb not allowed for me

Only a person holding the claims that the collection requires (the absorb claim, or the claims that the collection's config asks for instead) can absorb records. Anyone else sees "User does not have permission to absorb records".

## Absorbs that change jobs are blocked when job editing is off

Lines 70-79

Asked as: merge jobs blocked; job editing turned off; can't merge jobs; absorb jobs disabled

For collections whose absorb changes jobs (ChangesJobs), Turbo refuses if job editing is disabled. It shows "Job editing is disabled; cannot absorb records that modify jobs".

## Vendor absorb and undo depend on expenses editing

Lines 81-84

Asked as: vendor merge locked; vendor absorb expenses blocked; expenses editing disabled vendor; vendor merge fails expenses

Before an absorb or an undo, Turbo calls requireExpensesEditing for the collection. This check applies to vendor absorbs, and its outcome comes from that function, which is defined elsewhere.

## Absorb target and records must exist

Lines 92-103

Asked as: target record missing; record to merge not found; merge can't find record; absorb missing record

Turbo refuses an absorb with 'Failed to find target record' when the target record does not exist. It refuses with 'Failed to find record to absorb' when any listed record does not exist.

## Records to absorb must share the same parent value

Lines 105-115

Asked as: merge different parent; records from different parents; absorb mismatched parent field; can't merge different parent values

When the collection has a parent constraint, every record to absorb must have the same value in the parent field as the target. Otherwise Turbo refuses with 'Cannot absorb records with different <field> values'.

## Absorbed client with alias-billing profiles needs kept client alias

Lines 118-131

Asked as: merge clients alias error; client alias needed to merge; invoicing profile alias merge; merge clients needs alias

When absorbing clients, if the kept client has no alias and any client being absorbed has an invoicing profile that bills its alias, Turbo refuses. The message is 'Invoicing profiles of a client being absorbed bill its alias. Give the kept client an alias first.'

## Collection rules can refuse an absorb

Lines 135-139

Asked as: absorb blocked by rule; collection rule stops merge; merge refused by config; merge blocked by setup rule

A rule in the collection's absorb config can refuse the absorb. Turbo then shows that rule's own message.

## Absorb refusal by collection rules shows its own message

Lines 135-142

Asked as: absorb failed generic message; why did merge fail; merge error message; absorb failed no reason

When a rule of the collection's config refuses the absorb, Turbo shows that rule's message. Any other failure shows "Failed to absorb records".

## Absorb success message

Lines 144-146

Asked as: absorb worked confirmation; merge success message; how many records merged; absorb completed message

After a successful absorb, Turbo replies 'Successfully absorbed N records into <target id>', with N the number of IDs given.

## Absorb is all or nothing

Lines 161-344

Asked as: merge half done; absorb partly applied; merge rolled back everything; absorb left changes behind

The absorb runs in one transaction. If any step fails, nothing changes.

## Collection preparation can add IDs or leave nothing to absorb

Lines 163-177

Asked as: preparation removed all ids; nothing left after prep; collection prep empties merge; merge prep leaves nothing

Some collections run their own preparation before the absorb. It sees the data inside the transaction and may add IDs. If no IDs remain, Turbo refuses with 'there is nothing to absorb'.

## Absorb with nothing left to absorb is refused

Lines 165-177

Asked as: no ids left to merge; absorb says nothing to absorb; merge nothing remaining; nothing to merge error

When the collection's own preparation leaves no IDs to absorb, Turbo refuses with "there is nothing to absorb".

## Only one pending absorb per collection

Lines 203-210

Asked as: already pending absorb; can't start second merge; existing merge must be undone; merge already in progress

Turbo refuses an absorb while the collection already has an absorb action. The existing action must be undone or deleted first.

## Absorb moves references to the target and flags them for writeback

Lines 218-276

Asked as: references moved to target; legacy system writeback flag; _imported set to false; merge flags for legacy sync

Turbo repoints every reference to an absorbed record at the target. In tables with an _imported column it also sets _imported to false and updates the 'updated' timestamp, so the change is written back to the legacy system.

## Absorb sets references to the target and updates timestamps

Lines 245-268

Asked as: timestamps updated on absorb; updated time after merge; legacy writeback after merge; absorb updates updated field

Absorb points every reference to an absorbed record at the target. On tables with an _imported column, it also sets _imported to false and updates the updated time, so the change is written back to the legacy system.

## Absorb checks client links stay consistent

Lines 278-280

Asked as: client relation mismatch merge; linked to another client error; merge rolled back client link; invoicing profile client mismatch

After references move, Turbo checks every changed job, project request, contact and invoicing profile. If its client relations no longer agree, the whole absorb is rolled back with 'merge would leave ... linked to another client'. An invoicing profile that would bill a missing client alias also rolls it back.

## Absorb saves an undo record and deletes absorbed records

Lines 282-331

Asked as: absorb saves undo data; what gets deleted on merge; absorb undo record created; merge keeps backup of records

Turbo saves the absorbed records, any rows they own, and the old reference values in an absorb action. It then deletes the owned rows and the absorbed records. An absorb that moves no references still records each table and column with no rows.

## Absorb with no moved references still records tables and columns

Lines 282-292

Asked as: merge moved no references; absorb with no references; empty absorb action saved; merge with no linked rows

When an absorb moves no reference, Turbo records each reference table and column with no rows. This lets the absorb action be saved and undone.

## Absorbed records' owned rows are deleted

Lines 304-331

Asked as: owned rows deleted on merge; child rows removed absorb; dependent rows deleted; merge deletes owned records

Absorb deletes the rows that the absorbed records own, then the records. Turbo saves those rows in the absorb action first so undo can restore them.

## Absorb completion notifies clients

Lines 346-349

Asked as: screen not refreshing after merge; absorb completed event; refresh after merge; clients not updated after absorb

When an absorb or undo finishes, Turbo sends a '<collection>/absorb_completed' event to subscribed clients so they can refresh. A failure to send does not fail the absorb.

## Undo absorb is blocked by expenses editing for vendors

Lines 365-368

Asked as: undo absorb vendor expenses; undo blocked by expenses editing; vendor undo locked; expenses editing undo merge

Undo absorb runs the same expenses-editing check as absorb, before the permission check.

## Undo absorb requires the absorb permission

Lines 370-381

Asked as: undo absorb permission denied; can't undo merge; undo needs absorb claims; no permission to undo merge

Undoing an absorb needs the same claims as absorbing. Otherwise Turbo refuses with 'User does not have permission to undo absorb'.

## Undo needs an existing absorb action

Lines 383-393

Asked as: no absorb to undo; undo nothing found; no absorb action found; undo merge missing

Turbo refuses an undo with 'No absorb action found for collection' when the collection has no absorb action. Each collection has at most one.

## Undo child absorbs first

Lines 395-411

Asked as: child absorb blocks undo; undo parent absorb first; child merge must be undone; undo order for absorbs

Turbo refuses to undo an absorb while a dependent child collection has its own absorb. The message is 'Cannot undo X absorb while Y absorb exists; undo Y absorb first'.

## Undo restores records, owned rows and references

Lines 413-497

Asked as: undo merge restores records; bring back absorbed records; revert absorb; restore merged records

Undo recreates the absorbed records with their original IDs and data, without schema validation. It recreates their owned rows and restores each reference to its old value. It then deletes the absorb action so a new absorb can run.

## Undo restores references in all tables and deletes action

Lines 416-503

Asked as: undo failed message; undo restore references; undo absorb error; undo merge failed

Undo recreates the absorbed records with their old IDs and without schema validation. It restores their owned rows and sets each tracked reference back to its old value. Then it deletes the absorb action so a new absorb can run. If any step fails, Turbo shows "Failed to undo absorb".

## Undo skips deleted records and rechecks client links

Lines 484-503

Asked as: undo skips deleted records; undo recheck client links; undo rolled back; undo client link error

After undo restores references, Turbo rechecks client links. It skips records deleted since the absorb. A failure rolls back the undo, which Turbo reports as 'Failed to undo absorb'.

## Undo success message

Lines 511-513

Asked as: undo success message; absorb undone confirmation; undo worked message; merge undo completed

After a successful undo, Turbo replies 'Successfully undid absorb operation'.

## Changed request relations bump request version

Lines 757-762

Asked as: project request version bump; request review stale; request relations changed; absorb changes request version

When an absorb or undo changes a project request's relations, Turbo raises its version by 1 and sets its last-edited and updated times to now. Any open review of that request becomes stale.

## Creating a job moves pending absorb-undo references

Lines 768-808

Asked as: job created from request undo; client note moved to job; request absorb undo data; new job from request references

When a job is created from a project request, undo data of pending absorbs moves from the request to the job. The reference moves only if the job still holds the absorb target in that field. The client note made from the request's notes moves with it. The request's entries are then removed from the undo data.
