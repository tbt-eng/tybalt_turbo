<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: invoices

`invoices` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/invoices/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): invoices

```
@request.auth.user_claims_via_uid.cid.name ?= "invoicing" || @request.auth.user_claims_via_uid.cid.name ?= "report"
```

## Who can view one record (view rule): invoices

```
@request.auth.user_claims_via_uid.cid.name ?= "invoicing" || @request.auth.user_claims_via_uid.cid.name ?= "report"
```

## Who can create a record (create rule): invoices

```
@request.auth.user_claims_via_uid.cid.name ?= "invoicing"
```

## Who can update a record (update rule): invoices

```
@request.auth.user_claims_via_uid.cid.name ?= "invoicing" && status = "Draft"
```

## Who can delete a record (delete rule): invoices

```
@request.auth.user_claims_via_uid.cid.name ?= "invoicing" && status = "Draft"
```

## Fields: invoices

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `client` | relation | required; maxSelect 1; to `clients` |
| `number` | text | pattern "^(\\d{2}(0[1-9]\|1[0-2])\\d{3})?$" |
| `date` | text | pattern "^(\\d{4}-\\d{2}-\\d{2})?$" |
| `total` | number |  |
| `status` | select | required; values ["Draft","Issued"]; maxSelect 1 |
| `creator` | relation | maxSelect 1; to `users` |
| `issued_by` | relation | maxSelect 1; to `users` |
| `issued_at` | text |  |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |

## Indexes: invoices

A unique index means that two records can't have the same values in its fields.

- `CREATE UNIQUE INDEX `idx_invoices_number` ON `invoices` (`number`) WHERE `number` != ''`
- `CREATE INDEX `idx_invoices_client` ON `invoices` (`client`)`
