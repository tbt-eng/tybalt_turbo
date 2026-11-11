<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: work_records_notes

`work_records_notes` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/work_records_notes/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): work_records_notes

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can view one record (view rule): work_records_notes

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can create a record (create rule): work_records_notes

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can update a record (update rule): work_records_notes

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can delete a record (delete rule): work_records_notes

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Fields: work_records_notes

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `work_record` | relation | required; maxSelect 1; to `work_records`; cascadeDelete |
| `uid` | relation | required; maxSelect 1; to `users` |
| `note` | text | required |
| `created` | autodate | onCreate |

## Indexes: work_records_notes

A unique index means that two records can't have the same values in its fields.

- `CREATE INDEX idx_work_records_notes_work_record_created ON work_records_notes (work_record, created)`
