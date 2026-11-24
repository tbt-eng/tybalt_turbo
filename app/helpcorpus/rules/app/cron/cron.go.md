<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/cron/cron.go","sha256":"aabf2fff634aba1b18e482b631ca2b1412e774789c8308a71782507dd26920b8","model":"claude-sonnet-5-5"} -->

# Rules in app/cron/cron.go

## Replication health check runs every 5 minutes

Lines 11-15

Asked as: litestream snapshot old; replication stale alert; backup snapshot not updated; S3 backup check failing

Every 5 minutes (litestream_replication_check), Turbo checks that the latest litestream snapshot in S3 is not older than 35 minutes.

## Second approver PO reminder runs daily at 9pm UTC

Lines 17-26

Asked as: PO waiting for second approval; reminder to second approver; priority approver hasn't approved PO; purchase order second sign off

Every day at 9pm UTC (po_second_approval_notifications), Turbo sends po_second_approval_required notifications to the qualified second approvers. It sends them when a purchase order awaits second approval and the priority second approver (priority_second_approver) has not approved it after 24 hours. Turbo sends them immediately.

## Time sheet submission reminders on Tuesday to Thursday

Lines 28-32

Asked as: timesheet not submitted reminder; why did I get timesheet nag; remind me to submit my hours; last week's timesheet missing email

At 8am UTC on Tuesday, Wednesday and Thursday (timesheet_submission_reminders), Turbo sends timesheet_submission_reminder notifications. They go to staff who have not submitted their time sheet for the previous week.

## Expense approval reminders on Thursday and Friday

Lines 34-38

Asked as: manager expense approval nudge; expenses waiting for my approval email; reminder to approve staff expenses; Thursday expense reminder

At 9am UTC on Thursday and Friday (expense_approval_reminders), Turbo sends expense_approval_reminder notifications. They remind managers to approve the expenses that their staff submitted.

## Time sheet approval reminders on Tuesday to Thursday

Lines 40-44

Asked as: timesheets waiting for approval email; manager timesheet approval nudge; reminder to approve team hours; noon timesheet reminder

At 12pm UTC on Tuesday, Wednesday and Thursday (timesheet_approval_reminders), Turbo sends timesheet_approval_reminder notifications. They remind managers to approve the time sheets that their staff submitted.

## Expense uncommit request reminders on weekdays

Lines 46-50

Asked as: uncommit requests waiting; admin reminder uncommit; expense uncommit email; weekday uncommit alert

At 1pm UTC on Monday to Friday (expense_uncommit_request_reminders), Turbo sends expense_uncommit_reminder notifications. They remind admin holders that uncommit requests are waiting.

## Old computer logins and raw logins are deleted daily

Lines 52-57

Asked as: old logins deleted; login history cleanup; raw login logs purge; computer logins removed automatically

Every day at 7am UTC (computer_login_retention), Turbo deletes computer logins and raw logins older than COMPUTER_LOGIN_RETENTION_DAYS. Turbo never deletes the computers themselves.

## Decided absorb proposals are deleted after retention period

Lines 59-64

Asked as: old absorb proposals removed; decided proposals cleanup; absorb proposal history deleted; proposals disappearing

Every day at 7:10am UTC (absorb_proposal_retention), Turbo deletes absorb proposals that were decided more than ABSORB_PROPOSAL_RETENTION_DAYS ago. Turbo never deletes pending proposals.

## Currency rates refresh on weekday evenings

Lines 66-72

Asked as: exchange rates not updating; currency rate refresh failed; FX rates missing in evening; Bank of Canada rate sync

At 10pm UTC on Monday to Friday (currency_rate_sync), Turbo refreshes foreign-exchange rates, after the Bank of Canada business-day feed is expected to be published. If the refresh fails, Turbo only logs the error and does not retry until the next run.
