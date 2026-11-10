<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: currencies

`currencies` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/currencies/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): currencies

```
@request.auth.id != ""
```

## Who can view one record (view rule): currencies

```
@request.auth.id != ""
```

## Who can create a record (create rule): currencies

```
@request.auth.user_claims_via_uid.cid.name ?= "admin"
```

## Who can update a record (update rule): currencies

```
@request.auth.user_claims_via_uid.cid.name ?= "admin"
```

## Who can delete a record (delete rule): currencies

```
@request.auth.user_claims_via_uid.cid.name ?= "admin"
```

## Fields: currencies

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `code` | text | required; pattern "^[A-Z]{3}$" |
| `symbol` | text | required; pattern "^(\\p{Sc}\|[A-Z]{3})$" |
| `icon` | file | required; maxSelect 1; mimeTypes ["image/svg+xml"] |
| `rate` | number | required |
| `rate_date` | date |  |
| `ui_sort` | number |  |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |

## Indexes: currencies

A unique index means that two records can't have the same values in its fields.

- `CREATE UNIQUE INDEX `idx_xVaAmGcXs3` ON `currencies` (`code`)`
