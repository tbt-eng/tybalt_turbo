<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: claims

`claims` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/claims/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): claims

```
@request.auth.id != ""
```

## Who can view one record (view rule): claims

```
@request.auth.id != ""
```

## Who can create a record (create rule): claims

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can update a record (update rule): claims

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can delete a record (delete rule): claims

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Fields: claims

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `name` | text | required; min 2 |
| `description` | text | required; min 5 |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |

## Indexes: claims

A unique index means that two records can't have the same values in its fields.

- `CREATE UNIQUE INDEX `idx_3KEX8wA` ON `claims` (`name`)`
