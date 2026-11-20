<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: computer_raw_logins

`computer_raw_logins` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/computer_raw_logins/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): computer_raw_logins

```
@request.auth.id != '' && (
  @request.auth.user_claims_via_uid.cid.name ?= 'computers' ||
  @request.auth.user_claims_via_uid.cid.name ?= 'admin'
)
```

## Who can view one record (view rule): computer_raw_logins

```
@request.auth.id != '' && (
  @request.auth.user_claims_via_uid.cid.name ?= 'computers' ||
  @request.auth.user_claims_via_uid.cid.name ?= 'admin'
)
```

## Who can create a record (create rule): computer_raw_logins

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can update a record (update rule): computer_raw_logins

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can delete a record (delete rule): computer_raw_logins

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Fields: computer_raw_logins

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `report` | json |  |
| `computer_name` | text |  |
| `error` | text |  |
| `created` | autodate | onCreate |

## Indexes: computer_raw_logins

A unique index means that two records can't have the same values in its fields.

- `CREATE INDEX `idx_computer_raw_logins_computer_name` ON `computer_raw_logins` (`computer_name`)`
- `CREATE INDEX `idx_computer_raw_logins_created` ON `computer_raw_logins` (`created`)`
