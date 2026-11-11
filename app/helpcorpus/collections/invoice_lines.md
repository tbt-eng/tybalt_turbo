<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: invoice_lines

`invoice_lines` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/invoice_lines/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): invoice_lines

```
@request.auth.user_claims_via_uid.cid.name ?= "invoicing" || @request.auth.user_claims_via_uid.cid.name ?= "report"
```

## Who can view one record (view rule): invoice_lines

```
@request.auth.user_claims_via_uid.cid.name ?= "invoicing" || @request.auth.user_claims_via_uid.cid.name ?= "report"
```

## Who can create a record (create rule): invoice_lines

```
@request.auth.user_claims_via_uid.cid.name ?= "invoicing"
```

## Who can update a record (update rule): invoice_lines

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can delete a record (delete rule): invoice_lines

```
@request.auth.user_claims_via_uid.cid.name ?= "invoicing"
```

## Fields: invoice_lines

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `invoice` | relation | required; maxSelect 1; to `invoices` |
| `kind` | select | required; values ["expense","po"]; maxSelect 1 |
| `expense` | relation | maxSelect 1; to `expenses` |
| `purchase_order` | relation | maxSelect 1; to `purchase_orders` |
| `amount` | number | required; min 0.01 |
| `creator` | relation | maxSelect 1; to `users` |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |

## Indexes: invoice_lines

A unique index means that two records can't have the same values in its fields.

- `CREATE INDEX `idx_invoice_lines_invoice` ON `invoice_lines` (`invoice`)`
- `CREATE UNIQUE INDEX `idx_invoice_lines_expense` ON `invoice_lines` (`expense`) WHERE `expense` != ''`
- `CREATE UNIQUE INDEX `idx_invoice_lines_purchase_order` ON `invoice_lines` (`purchase_order`) WHERE `purchase_order` != ''`
