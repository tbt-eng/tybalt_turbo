<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/absorb/history.go","sha256":"af6dfdca37f09efaf63a0d5e7b4f6a82ded88a11e81f990e07bf511005c5454b","model":"claude-sonnet-5-5"} -->

# Rules in app/absorb/history.go

## Committed absorb is recorded in absorb history

Lines 9-16

Asked as: absorb committed history; where do committed absorbs go; undone absorb not in history; absorb history row missing

When an absorb is committed, Turbo copies the absorb action into the absorb history (absorb_history). It does this in the same transaction that deletes the action. A committed absorb always has a history row, and an absorb that was undone never has one.

## Absorb history contents copied from the action

Lines 21-35

Asked as: absorb history absorbed time; absorbed_at wrong time; absorb date shows creation; absorb history fields copied

The history row keeps the collection name, the target, the absorbed records, the updated references, the proposal and who absorbed. It sets the absorbed time (absorbed_at) to the time the absorb action was created, not the time of the commit.

## Absorb history keeps the target's name at commit time

Lines 27-30

Asked as: target name missing absorb; absorb history target name empty; kept record renamed absorb; target record not found absorb

Turbo saves the target's name (target_name) in the history row, because the kept record can change or be absorbed later. If Turbo cannot find the target record, it leaves the name empty and still saves the row.

## Committer is recorded only when a user committed

Lines 36-38

Asked as: who committed absorb; committed_by empty; absorb committed by system; absorb history committer blank

The history row records who committed (committed_by) only when a user made the request. If no user made it, the committer stays empty.
