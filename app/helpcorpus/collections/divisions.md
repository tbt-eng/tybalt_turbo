<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: divisions

`divisions` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/divisions/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): divisions

```
@request.auth.id != ""
```

## Who can view one record (view rule): divisions

```
@request.auth.id != ""
```

## Who can create a record (create rule): divisions

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can update a record (update rule): divisions

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can delete a record (delete rule): divisions

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Fields: divisions

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `code` | text | required; min 1 |
| `name` | text | required; min 2 |
| `active` | bool |  |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |

## Indexes: divisions

A unique index means that two records can't have the same values in its fields.

- `CREATE UNIQUE INDEX `idx_rbNPJNF` ON `divisions` (`code`)`
