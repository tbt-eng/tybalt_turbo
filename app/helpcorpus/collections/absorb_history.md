<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: absorb_history

`absorb_history` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/absorb_history/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): absorb_history

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can view one record (view rule): absorb_history

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can create a record (create rule): absorb_history

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can update a record (update rule): absorb_history

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can delete a record (delete rule): absorb_history

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Fields: absorb_history

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `collection_name` | text | required |
| `target_id` | text | required |
| `target_name` | text |  |
| `absorbed_records` | json | required; maxSize 2000000 |
| `updated_references` | json | required; maxSize 2000000 |
| `proposal` | text |  |
| `absorbed_by` | relation | maxSelect 1; to `users` |
| `absorbed_at` | date |  |
| `committed_by` | relation | maxSelect 1; to `users` |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |

## Indexes: absorb_history

A unique index means that two records can't have the same values in its fields.

- `CREATE INDEX `idx_absorb_history_created` ON `absorb_history` (`created`)`
