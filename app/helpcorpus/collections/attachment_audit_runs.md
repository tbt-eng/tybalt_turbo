<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: attachment_audit_runs

`attachment_audit_runs` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/attachment_audit_runs/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): attachment_audit_runs

```
@request.auth.user_claims_via_uid.cid.name ?= 'admin'
```

## Who can view one record (view rule): attachment_audit_runs

```
@request.auth.user_claims_via_uid.cid.name ?= 'admin'
```

## Who can create a record (create rule): attachment_audit_runs

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can update a record (update rule): attachment_audit_runs

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can delete a record (delete rule): attachment_audit_runs

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Fields: attachment_audit_runs

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `target_key` | text | required |
| `label` | text | required |
| `collection_name` | text | required |
| `field_name` | text | required |
| `status` | select | required; values ["running","completed","failed"]; maxSelect 1 |
| `requested_by` | relation | maxSelect 1; to `users` |
| `started_at` | date |  |
| `finished_at` | date |  |
| `total_records` | number | onlyInt |
| `referenced_records` | number | onlyInt |
| `matching_records` | number | onlyInt |
| `missing_records` | number | onlyInt |
| `orphaned_files` | number | onlyInt |
| `error` | text |  |
| `missing_report` | file | maxSelect 1; maxSize 536870912 |
| `orphaned_report` | file | maxSelect 1; maxSize 536870912 |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |

## Indexes: attachment_audit_runs

A unique index means that two records can't have the same values in its fields.

- `CREATE UNIQUE INDEX `idx_attachment_audit_runs_target_key` ON `attachment_audit_runs` (`target_key`)`
