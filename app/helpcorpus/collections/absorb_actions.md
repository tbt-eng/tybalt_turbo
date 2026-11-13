<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: absorb_actions

`absorb_actions` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/absorb_actions/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): absorb_actions

```
@request.auth.user_claims_via_uid.cid.name ?= 'absorb' || (collection_name = 'rate_sheets' && @request.auth.user_claims_via_uid.cid.name ?= 'admin')
```

## Who can view one record (view rule): absorb_actions

```
@request.auth.user_claims_via_uid.cid.name ?= 'absorb' || (collection_name = 'rate_sheets' && @request.auth.user_claims_via_uid.cid.name ?= 'admin')
```

## Who can create a record (create rule): absorb_actions

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can update a record (update rule): absorb_actions

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can delete a record (delete rule): absorb_actions

```
@request.auth.user_claims_via_uid.cid.name ?= 'absorb' || (collection_name = 'rate_sheets' && @request.auth.user_claims_via_uid.cid.name ?= 'admin')
```

## Fields: absorb_actions

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `collection_name` | text | required |
| `target_id` | text | required |
| `absorbed_records` | json | required; maxSize 2000000 |
| `updated_references` | json | required; maxSize 2000000 |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |

## Indexes: absorb_actions

A unique index means that two records can't have the same values in its fields.

- `CREATE UNIQUE INDEX `idx_T0t8iRR` ON `absorb_actions` (`collection_name`)`
