<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: time_sheets

`time_sheets` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/time_sheets/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): time_sheets

```
@request.auth.id = uid ||
(submitted = true && @request.auth.id = approver) ||
(submitted = true && @request.auth.id ?= time_sheet_reviewers_via_time_sheet.reviewer) ||
(submitted = true && approved != '' && @request.auth.user_claims_via_uid.cid.name ?= 'commit') ||
(committed != '' && @request.auth.user_claims_via_uid.cid.name ?= 'report')
```

## Who can view one record (view rule): time_sheets

```
@request.auth.id = uid ||
(submitted = true && @request.auth.id = approver) ||
(submitted = true && @request.auth.id ?= time_sheet_reviewers_via_time_sheet.reviewer) ||
(submitted = true && approved != '' && @request.auth.user_claims_via_uid.cid.name ?= 'commit') ||
(committed != '' && @request.auth.user_claims_via_uid.cid.name ?= 'report')
```

## Who can create a record (create rule): time_sheets

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can update a record (update rule): time_sheets

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can delete a record (delete rule): time_sheets

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Fields: time_sheets

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `uid` | relation | required; maxSelect 1; to `users` |
| `work_week_hours` | number | required; min 8; max 40 |
| `salary` | bool |  |
| `week_ending` | text | required; pattern "^\\d{4}-\\d{2}-\\d{2}$" |
| `submitted` | bool |  |
| `approver` | relation | maxSelect 1; to `users` |
| `rejector` | relation | maxSelect 1; to `users` |
| `rejection_reason` | text |  |
| `approved` | date |  |
| `rejected` | date |  |
| `committed` | date |  |
| `committer` | relation | maxSelect 1; to `users` |
| `payroll_id` | text | required; pattern "^(?:[1-9]\\d*\|CMS[0-9]{1,2})$" |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |
| `_imported` | bool |  |

## Indexes: time_sheets

A unique index means that two records can't have the same values in its fields.

- `CREATE UNIQUE INDEX `idx_NSP4DAc` ON `time_sheets` ( `uid`, `week_ending` )`
