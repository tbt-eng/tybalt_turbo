<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: rate_roles

`rate_roles` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/rate_roles/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): rate_roles

```
@request.auth.id != ""
```

## Who can view one record (view rule): rate_roles

```
@request.auth.id != ""
```

## Who can create a record (create rule): rate_roles

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can update a record (update rule): rate_roles

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can delete a record (delete rule): rate_roles

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Fields: rate_roles

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `name` | text | required |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |

## Indexes: rate_roles

A unique index means that two records can't have the same values in its fields.

- `CREATE UNIQUE INDEX `idx_aNRplyVmkw` ON `rate_roles` (`name`)`
