<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: time_types

`time_types` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/time_types/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): time_types

```
@request.auth.id != ""
```

## Who can view one record (view rule): time_types

```
@request.auth.id != ""
```

## Who can create a record (create rule): time_types

```
@request.auth.id != "" &&
@request.auth.user_claims_via_uid.cid.name ?= 'tt'
```

## Who can update a record (update rule): time_types

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can delete a record (delete rule): time_types

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Fields: time_types

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `code` | text | required; min 1 |
| `name` | text | required; min 2 |
| `description` | text |  |
| `allowed_fields` | json | required; maxSize 2000000 |
| `required_fields` | json | maxSize 2000000 |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |

## Indexes: time_types

A unique index means that two records can't have the same values in its fields.

- `CREATE UNIQUE INDEX `idx_fQtszvd` ON `time_types` (`code`)`
