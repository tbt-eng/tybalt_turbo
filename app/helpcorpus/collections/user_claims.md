<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: user_claims

`user_claims` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/user_claims/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): user_claims

```
@request.auth.id = uid ||
(
  @request.auth.id != "" &&
  @request.auth.user_claims_via_uid.cid.name ?= 'admin'
)
```

## Who can view one record (view rule): user_claims

```
@request.auth.id = uid ||
(
  @request.auth.id != "" &&
  @request.auth.user_claims_via_uid.cid.name ?= 'admin'
)
```

## Who can create a record (create rule): user_claims

```
@request.auth.id != "" &&
@request.auth.user_claims_via_uid.cid.name ?= 'admin'
```

## Who can update a record (update rule): user_claims

```
@request.auth.id != "" &&
@request.auth.user_claims_via_uid.cid.name ?= 'admin'
```

## Who can delete a record (delete rule): user_claims

```
@request.auth.id != "" &&
@request.auth.user_claims_via_uid.cid.name ?= 'admin'
```

## Fields: user_claims

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `uid` | relation | required; maxSelect 1; to `users` |
| `cid` | relation | required; maxSelect 1; to `claims` |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |
| `_imported` | bool |  |

## Indexes: user_claims

A unique index means that two records can't have the same values in its fields.

- `CREATE UNIQUE INDEX `idx_6dSZCrb` ON `user_claims` ( `uid`, `cid` )`
