<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/cron/absorb_proposal_retention.go","sha256":"e7490868d09e3b6eed592ae89f4bebf62a6ae2c6326c947028e5a0fe3ec40146","model":"claude-sonnet-5-5"} -->

# Rules in app/cron/absorb_proposal_retention.go

## Decided absorb proposals are deleted after the retention period

Lines 13-17

Asked as: old absorb proposals deleted; absorb proposal auto delete; retention period absorb; proposals removed after days

A daily job deletes each absorb proposal that is not pending and was decided more than ABSORB_PROPOSAL_RETENTION_DAYS days ago. The statuses are accepted, declined, withdrawn or obsolete. The cutoff is the current time in UTC minus that many days (the number itself is set in the constants file). A proposal decided before the cutoff is deleted, and one decided exactly at the cutoff stays.

## Pending absorb proposals are never deleted by retention

Lines 30-32

Asked as: pending absorb proposal kept; old pending proposal not deleted; why is my proposal still there; undecided proposal never removed

The retention job never deletes an absorb proposal with status pending, however old it is. It also skips any proposal that has an empty decided_at (decided_at).

## Deleting an absorb proposal also deletes its record rows

Lines 30-54

Asked as: absorb proposal records deleted; proposal delete transaction rollback; related rows removed with proposal; absorb record rows cleanup

When the retention job deletes an expired absorb proposal, it first deletes that proposal's rows in absorb_proposal_records, then the proposal. Both deletes run in one transaction. If either fails, the transaction rolls back and nothing is deleted. A closed proposal is not kept as history, because a committed absorb keeps no record of its own.
