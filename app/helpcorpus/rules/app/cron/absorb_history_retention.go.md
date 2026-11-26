<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/cron/absorb_history_retention.go","sha256":"a287ed03c3871c65487b52036b825c7fdefc11b7a0240d5b060dba9b4ee7a0a2","model":"claude-sonnet-5-5"} -->

# Rules in app/cron/absorb_history_retention.go

## Absorb history rows are deleted after the retention period

Lines 14-18

Asked as: absorb history cleanup; old absorb history deleted; absorb history retention period; why are absorb records disappearing

A daily job deletes the absorb history rows (absorb_history) of absorbs that were committed more than ABSORB_HISTORY_RETENTION_DAYS days ago. The number of days is a constant defined in another file, so this file does not show its value. The cutoff is the current time in UTC minus that many days.

## Absorb history is deleted only when strictly older than the cutoff

Lines 20-35

Asked as: row created exactly at cutoff; absorb history older than cutoff; which absorb history rows get deleted; commit time absorb history delete

Turbo deletes an absorb history row only when its created time (the commit time) is earlier than the cutoff. A row created exactly at the cutoff stays. The comparison is made as text in the PocketBase date layout. Turbo deletes the rows with plain SQL, because no hook watches the collection and no row depends on it.

## Daily absorb history cleanup runs on its own and logs the result

Lines 37-47

Asked as: daily cron absorb history; cleanup job logs deleted rows; absorb history delete failed; nightly absorb history job

The daily cron job deletes the expired absorb history rows and logs the cutoff and the number of rows deleted. If the delete fails, Turbo logs an error with the cutoff and does not retry in that run.
