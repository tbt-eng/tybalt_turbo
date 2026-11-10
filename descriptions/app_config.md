# app_config Reference

The `app_config` collection is a key-value store for runtime configuration. Each record has a `key` (the domain), a JSON `value` object, and an optional `description`. Authenticated users can read configs. The collection has no create, update or delete rule, so only a PocketBase superuser can change a config.

Configuration is read via helpers in `app/utilities/config.go`. Missing or invalid values fall back to defaults.

---

## Domain: `jobs`

Controls job editing and the Project Request creation workflow.

| Property | Type | Default | Description |
| --- | --- | --- | --- |
| `create_edit_absorb` | bool | `true` | Enables job creation, updates, admin-only manual renumbering, and client/contact merges. When false, these operations return HTTP 403. |
| `require_job_request_approval` | bool | `false` | Requires an approved Project Request for each new project and subjob, including imports and admin creation. When false, holders of `job` can create projects directly. |

Missing settings use the defaults. An invalid request-approval Boolean or a database read error blocks creation. Proposals remain outside the
request approval workflow. When request approval is disabled, retained requests
keep their state and age. Turbo hides the Job Create Queue and its badge, and blocks
request creation and job creation from a request. Enabling the setting again
restores the queue from retained request state.

This flag does not restrict work on existing jobs. Project-value changes always
require approval, including when this flag is false. See the
[Project Authorization specification](project_authorizations.md).

---

## Domain: `expenses`

Controls expense, PO, and vendor editing, plus expense validation thresholds.

| Property                    | Type   | Default   | Description                                                                                                                            |
|-----------------------------|--------|-----------|----------------------------------------------------------------------------------------------------------------------------------------|
| `create_edit_absorb`        | bool   | `true`    | Enables expense, purchase order, and vendor creation/updating/absorb. When `false`, these operations return HTTP 403.                  |
| `no_po_expense_limit`       | number | `100.0`   | A non-exempt expense of this many dollars or more requires a PO. Set to `0` to require a PO for all non-exempt expenses. Must be >= 0. |
| `uncommit_payroll_lock_days` | number | `7` | Days after `pay_period_ending` that a payroll-paid expense can still be requested for uncommit. See [uncommit_requests.md](uncommit_requests.md). Must be >= 0. |
| `uncommit_non_payroll_lock_days` | number | `90` | Days after commit that any other expense can still be requested for uncommit. Must be >= 0. |
| `po_expense_allowed_excess` | object | see below | Controls how much total expenses on a PO can exceed the PO total.                                                                      |

While `create_edit_absorb` is `false`, Turbo also refuses to delete an expense, PO or vendor. It refuses the expense actions (submit, recall, copy, approve, reject, commit, uncommit, and each step of an uncommit request) and the PO actions (approve, reject, cancel, close, make cumulative, revise, and creating a PO group).

### `po_expense_allowed_excess` sub-object

| Property  | Type   | Default       | Description                                                                                                                                            |
|-----------|--------|---------------|--------------------------------------------------------------------------------------------------------------------------------------------------------|
| `percent` | number | `5`           | Percentage overage allowed (0-100). Stored as human-readable percent, normalized internally to a fraction.                                             |
| `value`   | number | `100.0`       | Absolute dollar overage allowed. Must be >= 0.                                                                                                         |
| `mode`    | string | `"lesser_of"` | How to combine percent and value. `"lesser_of"` picks the smaller excess (more restrictive). `"greater_of"` picks the larger excess (more permissive). |

**Mode examples** for a $2,000 PO with percent=5, value=100:

- 5% of $2,000 = $100, absolute = $100 — both equal, limit is $2,100 either way
- For a $5,000 PO: 5% = $250, absolute = $100
  - `"lesser_of"` → uses $100 → limit is $5,100
  - `"greater_of"` → uses $250 → limit is $5,250

**Fail mode:** open (editing defaults to enabled)

---

## Domain: `time`

Controls time-entry and time-amendment editing, plus selected timesheet workflow mutations.

| Property      | Type | Default | Description                                                                                     |
|---------------|------|---------|-------------------------------------------------------------------------------------------------|
| `create_edit` | bool | `true`  | Enables time entry and time amendment creation/editing/deletion, time entry copy (one entry, or a week's entries to the next week), plus timesheet bundle, approve and uncommit. When `false`, these operations return HTTP 403. |

**Fail mode:** open (defaults to enabled)

---

## Domain: `purchase_orders`

Controls purchase order workflow behavior.

| Property                         | Type   | Default | Description                                                                                         |
|----------------------------------|--------|---------|-----------------------------------------------------------------------------------------------------|
| `second_stage_timeout_hours`     | number | `24.0`  | Hours a PO waits in "pending second approver" status before timing out. Must be > 0.                |
| `enable_legacy_po_create_update` | bool   | `false` | Enables the hidden legacy PO create/update flow for holders of the `legacy_po_create_update` claim. |

---

## Domain: `notifications`

Toggles individual notification templates on/off. Each property key is a template code. All default to `false` (fail-closed — notifications must be explicitly enabled).

| Template Code                          | Description                                 |
|----------------------------------------|---------------------------------------------|
| `po_approval_required`                 | PO requires first-approver action           |
| `po_revision_approval_required`        | Revised PO requires first-approver action   |
| `po_second_approval_required`          | PO requires second-approver action          |
| `po_priority_second_approval_required` | PO requires priority second-approver action |
| `po_active`                            | PO has been approved and is now active      |
| `po_rejected`                          | PO has been rejected                        |
| `expense_rejected`                     | Expense has been rejected                   |
| `expense_approval_reminder`            | Reminder to approve pending expenses        |
| `expense_uncommit_approved`            | Uncommit request was approved               |
| `expense_uncommit_declined`            | Uncommit request was declined               |
| `expense_uncommit_reminder`            | Reminder that uncommit requests are waiting |
| `timesheet_submission_reminder`        | Reminder to submit timesheet                |
| `timesheet_approval_reminder`          | Reminder to approve pending timesheets      |
| `timesheet_rejected`                   | Timesheet has been rejected                 |
| `timesheet_shared`                     | Timesheet has been shared with a viewer     |

The same row also switches the template codes that are not in this table: `job_request_submitted`, `job_request_approved`, `job_request_rejected`, `job_request_created`, `job_value_submitted`, `job_value_approved`, `job_value_rejected`, `job_client_changed`, `po_group_approval_required` and `po_group_second_approval_required`. The migrations that added those codes wrote `true` for them in this row, and did the same for `po_revision_approval_required` and the three `expense_uncommit_*` codes.

**Fail mode:** closed (defaults to disabled)

---

## Domain: `litestream`

Controls Litestream replica health checking and alerting.

| Property                  | Type   | Default              | Description                                                                                                                     |
|---------------------------|--------|----------------------|---------------------------------------------------------------------------------------------------------------------------------|
| `snapshot_prefix`         | string | `"litestream/0009/"` | S3 prefix checked by the cron health check when looking for the newest snapshot.                                                |
| `stale_threshold_minutes` | number | `35.0`               | Maximum allowed age of the newest snapshot before the replica is considered stale. Must be > 0.                                 |
| `alert_email`             | string | unset                | Email address notified when Litestream health is unhealthy: missing creds, S3 listing errors, no snapshots, or stale snapshots. |
| `restart_on_stale`        | bool   | `false`              | When `true`, exit the process after a stale-snapshot alert so the platform restarts the app. Stale only.                        |

**Fail mode:** alerting off unless `alert_email` is configured

---

## Domain: `help_chat`

Controls Help chat, which answers questions about Turbo in Help. [help_chat.md](help_chat.md) has the steps to set it up and to change the model.

| Property | Type | Default | Description |
| --- | --- | --- | --- |
| `enabled` | bool | `false` | Turns the chat on. When false, `/api/help/ask` returns HTTP 403. |
| `provider` | string | `"anthropic"` | The company whose API runs the model: `"anthropic"`, `"openai"` or `"openrouter"`. The server needs that provider's key. |
| `model` | string | see below | A model id of the provider. Each turn records the provider and the model that answered it. |
| `effort` | string | `"low"` | How much the model thinks before it answers. |
| `max_output_tokens` | number | `8000` | The most tokens in one model reply, with its thinking. 500 to 64000. |
| `max_question_chars` | number | `2000` | The longest question. 20 to 4000. |
| `questions_per_hour` | number | `30` | The most questions from one user in an hour. 1 to 1000. |
| `max_tool_calls` | number | `6` | The most searches and reads in one turn. 0 to 20. |
| `history_turns` | number | `10` | How many earlier turns of a chat go to the model. 0 to 50. |
| `turn_timeout_seconds` | number | `60` | The longest time for one turn. 5 to 300. |

Each provider has its own key, models and effort values:

| `provider` | Secret on the server | `model` | `effort` |
| --- | --- | --- | --- |
| `anthropic` | `ANTHROPIC_API_KEY` | `"claude-haiku-5-5"` (default) or `"claude-sonnet-5-5"` | `"low"`, `"medium"`, `"high"` |
| `openai` | `OPENAI_API_KEY` | Required. Any model id from OpenAI. | `"minimal"`, `"low"`, `"medium"`, `"high"`, or `""` to send none |
| `openrouter` | `OPENROUTER_API_KEY` | Required. Any model id from OpenRouter, such as `"openai/gpt-5-mini"`. | The same as `openai` |

A key is never a value in `app_config`. Every signed-in user can read this collection, and its rows are in the database backups.

**Fail mode:** closed. With no row the chat is off. A value that is out of range, a provider that is not in the list, a model that the provider does not allow, or no key on the server for the provider also leaves the chat off, and the routes return HTTP 503. The server log says which.

---

## Example values

```json
// key: "jobs"
{ "create_edit_absorb": true, "require_job_request_approval": false }

// key: "expenses"
{
  "create_edit_absorb": true,
  "no_po_expense_limit": 100.0,
  "po_expense_allowed_excess": {
    "percent": 5,
    "value": 100.0,
    "mode": "lesser_of"
  }
}

// key: "time"
{ "create_edit": true }

// key: "purchase_orders"
{
  "second_stage_timeout_hours": 24,
  "enable_legacy_po_create_update": false
}

// key: "notifications"
{
  "po_approval_required": true,
  "po_active": true,
  "expense_rejected": true,
  "timesheet_submission_reminder": false
}

// key: "help_chat"
{ "enabled": true, "provider": "anthropic", "model": "claude-haiku-5-5" }

// key: "litestream"
{
  "snapshot_prefix": "litestream/0009/",
  "stale_threshold_minutes": 35,
  "alert_email": "ops@example.com",
  "restart_on_stale": false
}
```
