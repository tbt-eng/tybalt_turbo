<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/notifications/queue_reminders.go","sha256":"7065174dcb5ee5dff2415b2aa9ac1218419f7d7710036cbe1edca1d3d2b89d2e","model":"claude-sonnet-5-5"} -->

# Rules in app/notifications/queue_reminders.go

## Second approver reminders go to qualified second approvers

Lines 27-66

Turbo queues a 'po_second_approval_required' notification for each user in the view pending_items_for_qualified_po_second_approvers whose num_pos_qualified is above 0. The notification links to /pos/pending. Turbo does not check for an earlier pending reminder, so a user can get a new one each time the job runs.

## Second approver reminder job stops at first failure

Lines 35-62

If saving one second approval notification fails, Turbo stops the job and returns the error 'error saving second approval notification'. Recipients later in the list get no reminder in that run. If the lookup of the view fails, Turbo returns 'error querying pending_items_for_qualified_po_second_approvers'.

## Queued notifications are sent only when requested

Lines 68-74

Reminders are created as deferred, system notifications. They are sent right away only when the send option is true. Otherwise they stay pending for later delivery.

## Time sheet submission reminder covers the previous week

Lines 82-97

The default timesheet submission reminder run works out the week ending of today's date and subtracts 7 days. It then queues reminders for that previous week ending.

## Submission reminder targets staff expected to submit without a submitted sheet

Lines 105-119

Turbo queues a 'timesheet_submission_reminder' for each user who has no submitted time sheet (submitted = 1) for the week ending and whose admin profile has time_sheet_expected = 1. A user with no admin profile is treated as not expected and gets no reminder. A time sheet that is saved but not submitted does not stop the reminder.

## Submission reminder is deduplicated per week ending

Lines 120-131

Turbo skips a user who already has a pending or inflight submission reminder for the same WeekEnding. A pending reminder for a different week does not block a new one. The reminder links to /time/entries/list.

## Time sheet approval reminders go to managers with waiting sheets

Lines 148-170

Turbo queues a 'timesheet_approval_reminder' for each approver who has a time sheet that is submitted and not approved, committed or rejected. Sheets with no approver are ignored. Each manager gets one reminder however many sheets wait. The reminder links to /time/sheets/pending.

## Approval reminders are limited to one per 24 hours

Lines 163-165

Turbo skips a recipient who already has a pending or inflight reminder of the same template created in the last 1 day. A reminder that was already sent does not block a new one. This applies to time sheet approval, expense approval and expense uncommit reminders.

## Expense approval reminders go to managers with waiting expenses

Lines 184-206

Turbo queues an 'expense_approval_reminder' for each approver who has an expense that is submitted and not approved, committed or rejected. Expenses with no approver are ignored. Each manager gets one reminder. It links to /expenses/pending.

## Expense approval reminder 24-hour limit

Lines 199-201

An approver with a pending or inflight expense approval reminder created in the last 1 day gets no new one.

## Uncommit request reminders go to admins while requests are pending

Lines 218-243

Turbo queues an 'expense_uncommit_reminder' for every user with the admin claim, but only when at least one expense uncommit request has status 'pending'. The reminder carries the count of pending requests and links to /expenses/uncommit-requests. An admin with a pending or inflight reminder created in the last 1 day gets no new one.

## Reminder job needs its template and skips failing recipients

Lines 246-327

A reminder job fails with 'error finding notification template' if the notification template for its code is missing. A recipient with no ID, a failed duplicate check or a failed notification save is skipped and logged. The job continues with the other recipients.

## Duplicate check counts only pending and inflight notifications

Lines 281-367

Before it creates a reminder, Turbo looks for a notification with the same recipient and template with status 'pending' or 'inflight' that also matches the job's own condition. Notifications with any other status do not count as duplicates. A job with no condition does no duplicate check.
