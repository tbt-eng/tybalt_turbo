<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: vendors

`vendors` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/vendors/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): vendors

```
@request.auth.id != ""
```

## Who can view one record (view rule): vendors

```
@request.auth.id != ""
```

## Who can create a record (create rule): vendors

```
@request.auth.id != "" &&
(
  @request.auth.user_claims_via_uid.cid.name ?= 'time' ||
  @request.auth.user_claims_via_uid.cid.name ?= 'payables_admin'
)
```

## Who can update a record (update rule): vendors

```
@request.auth.user_claims_via_uid.cid.name ?= 'payables_admin'
```

## Who can delete a record (delete rule): vendors

```
@request.auth.id != "" &&
@request.auth.user_claims_via_uid.cid.name ?= 'payables_admin' &&

// prevent deletion of vendors if there are referencing purchase orders
@collection.purchase_orders.vendor != id &&

// prevent deletion of vendors if there are referencing expenses
@collection.expenses.vendor != id
```

## Fields: vendors

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `name` | text | required; min 3 |
| `alias` | text | min 3 |
| `status` | select | required; values ["Active","Inactive"]; maxSelect 1 |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |
| `_imported` | bool |  |

## Indexes: vendors

A unique index means that two records can't have the same values in its fields.

- `CREATE UNIQUE INDEX `idx_GCZxhiM` ON `vendors` (`name`)`
- `CREATE UNIQUE INDEX `idx_c8OTvkU` ON `vendors` (`alias`) WHERE `alias` != ''`
