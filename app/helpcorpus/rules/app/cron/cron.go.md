<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/cron/cron.go","sha256":"e18fefb0df16aa40b87d3e0c9ead700860071fcd4193218202dd41bcae064724","model":"claude-sonnet-5-5"} -->

# Rules in app/cron/cron.go

## Litestream replication check every 5 minutes

Lines 11-15

Asked as: litestream not replicating; snapshot older than 35 minutes; backup replication health check; S3 backup stale

Turbo checks the health of litestream replication every 5 minutes (litestream_replication_check). It verifies that the latest snapshot in S3 is not older than 35 minutes.

## Second approval notification for purchase orders runs daily at 9pm UTC

Lines 17-26

Asked as: purchase order second approval reminder; PO waiting for second approver; second approver not notified; PO approval after 24 hours

Turbo runs the job po_second_approval_notifications every day at 9pm UTC. It notifies the qualified second approvers (po_second_approval_required) of a PO that awaits second approval when the priority second approver has not approved it after 24 hours. Turbo sends the notifications at once.

## Timesheet submission reminders on Tuesday to Thursday at 8am UTC

Lines 28-32

Asked as: timesheet not submitted reminder; reminder to submit timesheet; missing timesheet last week; timesheet nag email Tuesday

Turbo runs the job timesheet_submission_reminders at 8am UTC on Tuesday, Wednesday and Thursday. It reminds staff who have not submitted their time sheet for the previous week (timesheet_submission_reminder). Turbo sends the reminders at once.

## Expense approval reminders on Thursday and Friday at 9am UTC

Lines 34-38

Asked as: manager expense approval reminder; pending expenses to approve; expense reminder emails Thursday; staff expenses waiting for me

Turbo runs the job expense_approval_reminders at 9am UTC on Thursday and Friday. It reminds managers to approve the submitted expenses of their staff (expense_approval_reminder).

## Timesheet approval reminders on Tuesday to Thursday at 12pm UTC

Lines 40-44

Asked as: timesheets waiting for manager approval; timesheet approval email reminder; manager reminder approve time sheets; staff timesheets to sign off

Turbo runs the job timesheet_approval_reminders at 12pm UTC on Tuesday, Wednesday and Thursday. It reminds managers to approve the submitted time sheets of their staff (timesheet_approval_reminder).

## Expense uncommit request reminders on weekdays at 1pm UTC

Lines 46-50

Asked as: expense uncommit request waiting; admin reminder uncommit; uncommit requests pending alert; weekday expense uncommit email

Turbo runs the job expense_uncommit_request_reminders at 1pm UTC, Monday to Friday. It reminds admin holders that expense uncommit requests are waiting (expense_uncommit_reminder).

## Daily deletion of old computer logins

Lines 52-58

Asked as: old computer logins deleted; login history cleanup; computer login retention; raw logins disappearing

Turbo runs the job computer_login_retention at 7am UTC every day. It deletes computer logins, and raw logins without a computer name, that are older than COMPUTER_LOGIN_RETENTION_DAYS. A raw login with a computer name never expires, and Turbo never deletes computers.

## Daily deletion of old decided absorb proposals

Lines 60-65

Asked as: old absorb proposals removed; decided absorb proposals cleanup; absorb proposal deleted automatically; why did my absorb proposal vanish

Turbo runs the job absorb_proposal_retention at 7:10am UTC every day. It deletes absorb proposals that were decided more than ABSORB_PROPOSAL_RETENTION_DAYS ago. Turbo never deletes pending proposals.

## Daily deletion of old absorb history rows

Lines 67-71

Asked as: absorb history rows deleted; old absorb history cleanup; absorb history purged; absorb records disappearing

Turbo runs the job absorb_history_retention at 7:20am UTC every day, after the proposal retention job. It deletes the history rows of absorbs committed more than ABSORB_HISTORY_RETENTION_DAYS ago.

## Currency rates refresh on weekday evenings

Lines 73-79

Asked as: exchange rates not updated; currency rate refresh failed; FX rates missing evening; Bank of Canada rates not loading

Turbo refreshes foreign-exchange rates at 10pm UTC, Monday to Friday, after the Bank of Canada business-day feed is expected. If the refresh fails, Turbo only logs the error and does not retry until the next run.
