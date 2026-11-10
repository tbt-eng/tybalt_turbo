<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/cron/cron.go","sha256":"139c7e199cb26cf5db8da70b0ad2e75cbba11411d085c49e8762d04fb95b691c","model":"claude-sonnet-5-5"} -->

# Rules in app/cron/cron.go

## Litestream replication health check runs every 5 minutes

Lines 11-15

Asked as: litestream replication check; snapshot older than 35 minutes; backup snapshot stale alert; S3 replication health

Turbo checks litestream replication every 5 minutes (cron "litestream_replication_check"). The check verifies that the latest snapshot in S3 is not older than 35 minutes. The check itself is in checkLitestreamReplication, which is not in this file.

## Second approval reminder sent daily at 9pm UTC

Lines 17-26

Asked as: second approval reminder; purchase order waiting approval; PO approver nudge 24 hours; who gets notified for second approval

Every day at 9pm UTC, Turbo sends po_second_approval_required notifications to the qualified second approvers. It sends them when a purchase order awaits second approval and the priority second approver has not approved it after 24 hours. The notifications are sent immediately (true flag). The selection of recipients is in QueuePoSecondApproverNotifications, which is not in this file.

## Time sheet submission reminders Tuesday to Thursday 8am UTC

Lines 28-32

Asked as: timesheet not submitted reminder; when do timesheet reminders go out; reminder to submit last week's hours; missing timesheet email

At 8am UTC on Tuesday, Wednesday and Thursday, Turbo sends timesheet_submission_reminder notifications. They remind staff who have not submitted their time sheet for the previous week. The notifications are sent immediately.

## Expense approval reminders Thursday and Friday 9am UTC

Lines 34-38

Asked as: expense approval reminder email; manager reminded to approve expenses; when are expense reminders sent; pending expenses nudge

At 9am UTC on Thursday and Friday, Turbo sends expense_approval_reminder notifications. They remind managers to approve the expenses that their staff submitted. The notifications are sent immediately.

## Time sheet approval reminders Tuesday to Thursday noon UTC

Lines 40-44

Asked as: timesheet approval reminder; manager reminder to approve hours; noon reminder approve timesheets; staff timesheets waiting for me

At 12pm UTC on Tuesday, Wednesday and Thursday, Turbo sends timesheet_approval_reminder notifications. They remind managers to approve the time sheets that their staff submitted. The notifications are sent immediately.

## Expense uncommit request reminders on weekdays 1pm UTC

Lines 46-50

Asked as: uncommit request reminder; admin reminder expense uncommit; uncommit requests waiting; daily admin reminder weekdays

At 1pm UTC from Monday to Friday, Turbo sends expense_uncommit_reminder notifications. They remind admin holders that expense uncommit requests are waiting. The notifications are sent immediately.

## Currency rates refresh weekday evenings at 10pm UTC

Lines 52-58

Asked as: exchange rates not updating; currency sync failed; FX rates refresh time; why did currency rates fail

At 10pm UTC from Monday to Friday, Turbo refreshes foreign-exchange rates (SyncCurrencyRates). The time is chosen so that the Bank of Canada business-day feed is expected to be published. If the sync fails, Turbo only logs the error "currency rate sync failed". It does not retry before the next scheduled run, and it sends no notification.
