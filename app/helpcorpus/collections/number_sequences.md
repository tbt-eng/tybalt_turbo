<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: number_sequences

`number_sequences` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/number_sequences/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): number_sequences

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can view one record (view rule): number_sequences

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can create a record (create rule): number_sequences

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can update a record (update rule): number_sequences

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can delete a record (delete rule): number_sequences

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Fields: number_sequences

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `scope` | text | required |
| `value` | number | required; onlyInt |

## Indexes: number_sequences

A unique index means that two records can't have the same values in its fields.

- `CREATE UNIQUE INDEX idx_number_sequences_scope ON number_sequences (scope)`
