<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/cron/cron.go","sha256":"47176de920625885fda10257862da530d4d6ea2d00978187275f470b07960b6f","model":"claude-sonnet-5-5"} -->

# Rules in app/cron/cron.go

## Litestream replication check runs every 5 minutes

Lines 11-15

Asked as: litestream snapshot old; replication health check; S3 backup stale; backup not replicating

Every 5 minutes, Turbo checks litestream replication health. It verifies that the latest snapshot in S3 is not older than 35 minutes.

## Second approver reminder runs daily at 9pm UTC

Lines 17-26

Asked as: second approver not approving; PO waiting too long; reminder for second approval; purchase order stuck awaiting approval

Every day at 9pm UTC, Turbo sends po_second_approval_required notifications to the qualified second approvers. It sends them when a purchase order awaits second approval and the priority second approver (priority_second_approver) has not approved it after 24 hours. Turbo sends these notifications immediately.

## Time sheet submission reminders on Tuesday to Thursday

Lines 28-32

Asked as: forgot to submit timesheet; timesheet reminder email; haven't filled in my hours; weekly timesheet nag

At 8am UTC on Tuesday, Wednesday and Thursday, Turbo sends timesheet_submission_reminder notifications. They remind staff who have not submitted their time sheet for the previous week. Turbo sends them immediately.

## Expense approval reminders on Thursday and Friday

Lines 34-38

Asked as: manager expense reminder; expenses waiting for my approval; why did I get expense email; pending expense approvals notice

At 9am UTC on Thursday and Friday, Turbo sends expense_approval_reminder notifications. They remind managers to approve the expenses that their staff submitted. Turbo sends them immediately.

## Time sheet approval reminders on Tuesday to Thursday

Lines 40-44

Asked as: timesheets waiting for approval; manager timesheet reminder; staff hours need approving; approval reminder for hours

At 12pm UTC on Tuesday, Wednesday and Thursday, Turbo sends timesheet_approval_reminder notifications. They remind managers to approve the time sheets that their staff submitted. Turbo sends them immediately.

## Expense uncommit request reminders on weekdays

Lines 46-50

Asked as: uncommit request waiting; admin reminder expenses; expense uncommit notification; uncommit requests pending

At 1pm UTC on Monday to Friday, Turbo sends expense_uncommit_reminder notifications. They remind admin holders that expense uncommit requests are waiting. Turbo sends them immediately.

## Old computer logins and raw logins are deleted daily

Lines 52-57

Asked as: old logins removed; login history cleanup; computer logins deleted; raw logins purge

Every day at 7am UTC, Turbo deletes computer logins and raw logins older than COMPUTER_LOGIN_RETENTION_DAYS. Turbo never deletes the computers themselves.

## Decided absorb proposals are deleted after retention period

Lines 59-64

Asked as: absorb proposals disappeared; old absorb proposals cleared; decided proposals removed; proposal cleanup

Every day at 7:10am UTC, Turbo deletes absorb proposals that were decided more than ABSORB_PROPOSAL_RETENTION_DAYS ago. Turbo never deletes pending proposals.

## Absorb history is deleted after retention period

Lines 66-70

Asked as: absorb history gone; absorb records deleted; old absorb rows purged; absorb history cleanup

Every day at 7:20am UTC, Turbo deletes the history rows of absorbs that were committed more than ABSORB_HISTORY_RETENTION_DAYS ago. This runs after the absorb proposal retention job.

## Currency rates refresh on weekday evenings

Lines 72-78

Asked as: exchange rates not updated; currency rates weekend; FX rates refresh time; no new exchange rate today

At 10pm UTC on Monday to Friday, Turbo refreshes the foreign-exchange rates (SyncCurrencyRates). This runs after the Bank of Canada business-day feed is expected to be published. Turbo does not refresh rates on weekends.
