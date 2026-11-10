<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: mileage_reset_dates

`mileage_reset_dates` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/mileage_reset_dates/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): mileage_reset_dates

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can view one record (view rule): mileage_reset_dates

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can create a record (create rule): mileage_reset_dates

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can update a record (update rule): mileage_reset_dates

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can delete a record (delete rule): mileage_reset_dates

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Fields: mileage_reset_dates

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `date` | text | required; pattern "^\\d{4}-\\d{2}-\\d{2}$" |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |
| `_imported` | bool |  |

## Indexes: mileage_reset_dates

A unique index means that two records can't have the same values in its fields.

- `CREATE UNIQUE INDEX `idx_0o9BK5jvdD` ON `mileage_reset_dates` (`date`)`
