<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: expense_rates

`expense_rates` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/expense_rates/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): expense_rates

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can view one record (view rule): expense_rates

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can create a record (create rule): expense_rates

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can update a record (update rule): expense_rates

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can delete a record (delete rule): expense_rates

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Fields: expense_rates

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `effective_date` | text | required; pattern "^\\d{4}-\\d{2}-\\d{2}$" |
| `breakfast` | number | required |
| `lunch` | number | required |
| `dinner` | number | required |
| `lodging` | number | required |
| `mileage` | json | required; maxSize 2000000 |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |

## Indexes: expense_rates

A unique index means that two records can't have the same values in its fields.

- `CREATE INDEX `idx_2mgZH2F08o` ON `expense_rates` (`effective_date`)`
