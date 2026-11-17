<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: rate_sheet_codes

`rate_sheet_codes` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/rate_sheet_codes/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): rate_sheet_codes

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can view one record (view rule): rate_sheet_codes

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can create a record (create rule): rate_sheet_codes

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can update a record (update rule): rate_sheet_codes

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can delete a record (delete rule): rate_sheet_codes

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Fields: rate_sheet_codes

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `code` | text | required; pattern "^\\d{4}-(0[1-9]\|1[0-2])-\\d{5}$" |
| `month` | text | required; pattern "^\\d{4}-(0[1-9]\|1[0-2])$" |
| `serial` | number | required; min 1; max 99999; onlyInt |
| `created` | autodate | onCreate |

## Indexes: rate_sheet_codes

A unique index means that two records can't have the same values in its fields.

- `CREATE UNIQUE INDEX `idx_rate_sheet_codes_code` ON `rate_sheet_codes` (`code`)`
- `CREATE UNIQUE INDEX `idx_rate_sheet_codes_month_serial` ON `rate_sheet_codes` (`month`, `serial`)`
