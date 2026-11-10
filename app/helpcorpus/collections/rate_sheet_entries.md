<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: rate_sheet_entries

`rate_sheet_entries` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/rate_sheet_entries/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): rate_sheet_entries

```
@request.auth.id != ""
```

## Who can view one record (view rule): rate_sheet_entries

```
@request.auth.id != ""
```

## Who can create a record (create rule): rate_sheet_entries

```
@request.auth.id != "" &&
@request.auth.user_claims_via_uid.cid.name ?= 'job'
```

## Who can update a record (update rule): rate_sheet_entries

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can delete a record (delete rule): rate_sheet_entries

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Fields: rate_sheet_entries

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `role` | relation | required; maxSelect 1; to `rate_roles` |
| `rate_sheet` | relation | required; maxSelect 1; to `rate_sheets`; cascadeDelete |
| `rate` | number | required; min 1; onlyInt |
| `overtime_rate` | number | required; min 1 |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |

## Indexes: rate_sheet_entries

A unique index means that two records can't have the same values in its fields.

- `CREATE UNIQUE INDEX `idx_MLXiy2bT4z` ON `rate_sheet_entries` ( `role`, `rate_sheet` )`
