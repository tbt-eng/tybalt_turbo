<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: expense_uncommit_requests

`expense_uncommit_requests` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/expense_uncommit_requests/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): expense_uncommit_requests

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can view one record (view rule): expense_uncommit_requests

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can create a record (create rule): expense_uncommit_requests

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can update a record (update rule): expense_uncommit_requests

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can delete a record (delete rule): expense_uncommit_requests

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Fields: expense_uncommit_requests

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `expense` | relation | maxSelect 1; to `expenses` |
| `requester` | relation | required; maxSelect 1; to `users` |
| `reason` | relation | maxSelect 1; to `expense_uncommit_reasons` |
| `note` | text | max 1000 |
| `outcome` | select | required; values ["uncommit_and_reject","uncommit_only"]; maxSelect 1 |
| `status` | select | required; values ["pending","approved","declined","cancelled"]; maxSelect 1 |
| `decided_by` | relation | maxSelect 1; to `users` |
| `decided_at` | date |  |
| `decision_comment` | text | max 1000 |
| `employee_name` | text |  |
| `po_number` | text |  |
| `vendor_name` | text |  |
| `description` | text |  |
| `currency_code` | text |  |
| `payment_type` | text |  |
| `total` | number |  |
| `committed` | date |  |
| `committer` | relation | maxSelect 1; to `users` |
| `committed_week_ending` | text |  |
| `pay_period_ending` | text |  |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |

## Indexes: expense_uncommit_requests

A unique index means that two records can't have the same values in its fields.

- `CREATE UNIQUE INDEX `idx_expense_uncommit_requests_pending` ON `expense_uncommit_requests` (`expense`) WHERE `status` = 'pending'`
