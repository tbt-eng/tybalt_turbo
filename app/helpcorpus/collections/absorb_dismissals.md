<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: absorb_dismissals

`absorb_dismissals` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/absorb_dismissals/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): absorb_dismissals

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can view one record (view rule): absorb_dismissals

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can create a record (create rule): absorb_dismissals

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can update a record (update rule): absorb_dismissals

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can delete a record (delete rule): absorb_dismissals

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Fields: absorb_dismissals

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `collection_name` | select | required; values ["clients","vendors"]; maxSelect 1 |
| `record_ids` | json | required; maxSize 100000 |
| `dismissed_by` | relation | maxSelect 1; to `users` |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |

## Indexes: absorb_dismissals

A unique index means that two records can't have the same values in its fields.

- `CREATE INDEX `idx_absorb_dismissals_collection` ON `absorb_dismissals` (`collection_name`)`
