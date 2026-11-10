<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: time_sheet_reviewers

`time_sheet_reviewers` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/time_sheet_reviewers/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): time_sheet_reviewers

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can view one record (view rule): time_sheet_reviewers

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can create a record (create rule): time_sheet_reviewers

```
// The caller is authenticated
@request.auth.id != "" &&
// The caller is the approver of the specified time_sheet
@request.auth.time_sheets_via_approver.id ?= time_sheet
```

## Who can update a record (update rule): time_sheet_reviewers

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can delete a record (delete rule): time_sheet_reviewers

```
// The caller is authenticated
@request.auth.id != "" &&
// The caller is the approver of the specified time_sheet
@request.auth.time_sheets_via_approver.id ?= time_sheet
```

## Fields: time_sheet_reviewers

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `time_sheet` | relation | required; maxSelect 1; to `time_sheets`; cascadeDelete |
| `reviewer` | relation | required; maxSelect 1; to `users` |
| `reviewed` | date |  |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |

## Indexes: time_sheet_reviewers

A unique index means that two records can't have the same values in its fields.

- `CREATE UNIQUE INDEX `idx_MVTW8sD` ON `time_sheet_reviewers` ( `time_sheet`, `reviewer` )`
