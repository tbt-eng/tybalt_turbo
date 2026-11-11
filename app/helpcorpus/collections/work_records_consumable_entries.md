<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: work_records_consumable_entries

`work_records_consumable_entries` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/work_records_consumable_entries/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): work_records_consumable_entries

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can view one record (view rule): work_records_consumable_entries

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can create a record (create rule): work_records_consumable_entries

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can update a record (update rule): work_records_consumable_entries

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can delete a record (delete rule): work_records_consumable_entries

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Fields: work_records_consumable_entries

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `work_record` | relation | required; maxSelect 1; to `work_records`; cascadeDelete |
| `consumable` | relation | required; maxSelect 1; to `work_records_consumables` |
| `quantity_number` | number |  |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |

## Indexes: work_records_consumable_entries

A unique index means that two records can't have the same values in its fields.

- `CREATE UNIQUE INDEX idx_work_records_consumable_entries_work_record_consumable ON work_records_consumable_entries (work_record, consumable)`
