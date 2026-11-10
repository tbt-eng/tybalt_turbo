<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: categories

`categories` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/categories/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): categories

```
@request.auth.id != ""
```

## Who can view one record (view rule): categories

```
@request.auth.id != ""
```

## Who can create a record (create rule): categories

```
@request.auth.id != "" &&
@request.auth.user_claims_via_uid.cid.name ?= 'job'
```

## Who can update a record (update rule): categories

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can delete a record (delete rule): categories

```
@request.auth.id != "" &&
@request.auth.user_claims_via_uid.cid.name ?= 'job' &&

// prevent deletion of categories if there are referencing time_entries
@collection.time_entries.category != id &&

// prevent deletion of categories if there are referencing purchase orders
@collection.purchase_orders.category != id &&

// prevent deletion of categories if there are referencing expenses
@collection.expenses.category != id
```

## Fields: categories

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `name` | text | required; min 3 |
| `job` | relation | required; maxSelect 1; to `jobs`; cascadeDelete |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |
| `_imported` | bool |  |

## Indexes: categories

A unique index means that two records can't have the same values in its fields.

- `CREATE UNIQUE INDEX `idx_SF6A76x` ON `categories` ( `job`, `name` )`
- `CREATE INDEX `idx_XMSeZCfYU0` ON `categories` (`job`)`
