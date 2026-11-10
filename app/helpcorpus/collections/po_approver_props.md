<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: po_approver_props

`po_approver_props` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/po_approver_props/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): po_approver_props

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can view one record (view rule): po_approver_props

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can create a record (create rule): po_approver_props

```
@request.auth.id != "" &&
@request.auth.user_claims_via_uid.cid.name ?= 'admin'
```

## Who can update a record (update rule): po_approver_props

```
@request.auth.id != "" &&
@request.auth.user_claims_via_uid.cid.name ?= 'admin'
```

## Who can delete a record (delete rule): po_approver_props

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Fields: po_approver_props

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `max_amount` | number | required; onlyInt |
| `divisions` | relation | maxSelect 999; to `divisions` |
| `user_claim` | relation | required; maxSelect 1; to `user_claims`; cascadeDelete |
| `project_max` | number | onlyInt |
| `sponsorship_max` | number | onlyInt |
| `staff_and_social_max` | number | onlyInt |
| `media_and_event_max` | number | onlyInt |
| `computer_max` | number | onlyInt |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |

## Indexes: po_approver_props

A unique index means that two records can't have the same values in its fields.

- `CREATE UNIQUE INDEX `idx_KQah2XAlqx` ON `po_approver_props` (`user_claim`)`
