### Notification Templates for Timesheets and Expenses

Notification templates for the notifications described in [Issue #86](https://github.com/tbt-eng/tybalt_turbo/issues/86). All six are built: `app/notifications/queue_reminders.go` queues the three reminders from cron jobs, and `app/notifications/queue_events.go` sends the other three when the event happens. Turbo sends a template only while its code is `true` in the `notifications` row of `app_config`.

Each link is the `{{.ActionURL}}` placeholder. The code fills it with the app URL and the path of the page: `/time/entries/list`, `/expenses/pending`, `/time/sheets/pending`, `/time/sheets/<id>/details` or `/expenses/<id>/details`. Turbo refuses to send an email whose text still has the old `{APP_URL}` or `{:RECORD_ID}` placeholder.

---

### `timesheet_submission_reminder`

- **Code**: `timesheet_submission_reminder`
- **Description**: Sent to users who have not submitted their timesheet for the previous week. Only users whose admin profile has `time_sheet_expected` set get it.
- **Subject**: `Please submit a timesheet for last week`
- **Text email**:

```text
Hello {{.RecipientName}},

You have not submitted your timesheet for the week ending {{.WeekEnding}}.

Please review your time entries then submit here:

{{.ActionURL}}
```

---

### `expense_approval_reminder`

- **Code**: `expense_approval_reminder`
- **Description**: Sent to managers with one or more expenses awaiting their approval.
- **Subject**: `Expenses await your approval`
- **Text email**:

```text
Hello {{.RecipientName}},

One or more expenses are awaiting your approval.

Please review and approve or reject them here:

{{.ActionURL}}
```

---

### `timesheet_approval_reminder`

- **Code**: `timesheet_approval_reminder`
- **Description**: Sent to managers with one or more timesheets awaiting their approval.
- **Subject**: `Time sheets await your approval`
- **Text email**:

```text
Hello {{.RecipientName}},

One or more timesheets are awaiting your approval.

Please review and approve or reject them here:

{{.ActionURL}}
```

---

### `timesheet_rejected`

- **Code**: `timesheet_rejected`
- **Description**: Sent when a timesheet is rejected to the employee, the rejector, and the employee's manager (if different from the rejector).
- **Subject**: `A time sheet was rejected`
- **Text email**:

```text
Hello {{.RecipientName}},

The timesheet for {{.EmployeeName}} for the week ending {{.WeekEnding}} was rejected by {{.RejectorName}} for the following reason:

{{.RejectionReason}}

You can review the timesheet and make any required changes here:

{{.ActionURL}}
```

---

### `expense_rejected`

- **Code**: `expense_rejected`
- **Description**: Sent when an expense is rejected to the employee, the rejector, and the employee's manager (if different from the rejector).
- **Subject**: `An expense was rejected`
- **Text email**:

```text
Hello {{.RecipientName}},

The expense submitted by {{.EmployeeName}} on {{.ExpenseDate}} for {{.ExpenseAmount}} was rejected by {{.RejectorName}} for the following reason:

{{.RejectionReason}}

You can review the expense and make any required changes here:

{{.ActionURL}}
```

---

### `timesheet_shared`

- **Code**: `timesheet_shared`
- **Description**: Sent to newly added viewers when a timesheet is shared with them for review.
- **Subject**: `A time sheet has been shared with you`
- **Text email**:

```text
Hello {{.RecipientName}},

{{.UserName}} has shared a timesheet for {{.EmployeeName}} for the week ending {{.WeekEnding}} with you for review.

You can view the shared timesheet here:

{{.ActionURL}}
```
