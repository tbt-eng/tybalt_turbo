<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/cron/cron.go","sha256":"a72cda950845f07287a41fd3f4633660b46a85b75b10853a1531e75828e2af2d","model":"claude-sonnet-5-5"} -->

# Rules in app/cron/cron.go

## Replication health check runs every 5 minutes

Lines 11-15

Asked as: litestream backup stale; snapshot older than 35 minutes; replication health alert; S3 backup not updating

Every 5 minutes, Turbo checks litestream replication health. It verifies that the latest snapshot in S3 is not older than 35 minutes.

## Second approver reminder runs daily at 9pm UTC

Lines 17-26

Asked as: PO waiting for second approval; second approver reminder email; purchase order not approved after 24 hours; who needs to approve the PO

Turbo sends po_second_approval_required notifications every day at 9pm UTC. They go to the qualified second approvers of a purchase order (PO) that awaits second approval. Turbo sends one only when the priority second approver (priority_second_approver) has not approved the PO after 24 hours. Turbo sends them immediately.

## Timesheet submission reminders on Tuesday to Thursday at 8am UTC

Lines 28-32

Asked as: timesheet reminder email; haven't submitted my timesheet; why am I getting timesheet emails; last week's time sheet reminder

Turbo sends timesheet_submission_reminder notifications at 8am UTC on Tuesday, Wednesday and Thursday. They remind staff who have not submitted their time sheet for the previous week. Turbo sends them immediately.

## Expense approval reminders on Thursday and Friday at 9am UTC

Lines 34-38

Asked as: expense approval reminder; manager reminder to approve expenses; expense approval emails Thursday; staff expenses waiting for me

Turbo sends expense_approval_reminder notifications at 9am UTC on Thursday and Friday. They remind managers to approve the expenses that their staff submitted. Turbo sends them immediately.

## Timesheet approval reminders on Tuesday to Thursday at noon UTC

Lines 40-44

Asked as: timesheet approval reminder; approve staff time sheets reminder; timesheets waiting for my approval; time sheet approval emails

Turbo sends timesheet_approval_reminder notifications at 12pm UTC on Tuesday, Wednesday and Thursday. They remind managers to approve the time sheets that their staff submitted. Turbo sends them immediately.

## Expense uncommit request reminders on weekdays at 1pm UTC

Lines 46-50

Asked as: uncommit request reminder; admin reminder uncommit; uncommit requests waiting; expense uncommit emails

Turbo sends expense_uncommit_reminder notifications at 1pm UTC, Monday to Friday. They remind admin holders that uncommit requests are waiting. Turbo sends them immediately.

## Old computer logins are deleted daily at 7am UTC

Lines 52-57

Asked as: old computer logins removed; login history deleted; why are my logins gone; computer login cleanup

Every day at 7am UTC (2am or 3am Eastern), Turbo deletes computer logins and raw logins that are older than COMPUTER_LOGIN_RETENTION_DAYS. Turbo never deletes computers in this job.

## Currency rates refresh on weekday evenings at 10pm UTC

Lines 59-65

Asked as: exchange rates not updating; currency rate refresh failed; foreign exchange rate missing; Bank of Canada rates not loaded

Turbo refreshes foreign-exchange rates at 10pm UTC, Monday to Friday. The time is set for after the Bank of Canada business-day feed is expected to be published. When the refresh fails, Turbo logs an error and does not retry until the next scheduled run.
