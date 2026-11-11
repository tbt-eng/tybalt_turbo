<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: job_rate_sheets

`job_rate_sheets` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/job_rate_sheets/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): job_rate_sheets

```
@request.auth.id != ""
```

## Who can view one record (view rule): job_rate_sheets

```
@request.auth.id != ""
```

## Who can create a record (create rule): job_rate_sheets

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can update a record (update rule): job_rate_sheets

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can delete a record (delete rule): job_rate_sheets

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Fields: job_rate_sheets

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `job` | relation | required; maxSelect 1; to `jobs`; cascadeDelete |
| `rate_sheet` | relation | required; maxSelect 1; to `rate_sheets` |
| `effective_date` | text | pattern "^(\\d{4}-\\d{2}-\\d{2})?$" |
| `reason` | text | max 500 |
| `created_by` | relation | maxSelect 1; to `users` |
| `updated_by` | relation | maxSelect 1; to `users` |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |

## Indexes: job_rate_sheets

A unique index means that two records can't have the same values in its fields.

- `CREATE UNIQUE INDEX `idx_job_rate_sheets_job_date` ON `job_rate_sheets` (`job`, `effective_date`)`
- `CREATE UNIQUE INDEX `idx_job_rate_sheets_opening` ON `job_rate_sheets` (`job`) WHERE `effective_date` = ''`
