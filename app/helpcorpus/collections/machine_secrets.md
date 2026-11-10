<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: machine_secrets

`machine_secrets` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/machine_secrets/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): machine_secrets

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can view one record (view rule): machine_secrets

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can create a record (create rule): machine_secrets

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can update a record (update rule): machine_secrets

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can delete a record (delete rule): machine_secrets

```
@request.auth.user_claims_via_uid.cid.name ?= 'admin'
```

## Fields: machine_secrets

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `sha256_hash` | text | required; pattern "^[a-fA-F0-9]{64}$" |
| `salt` | text | required; pattern "^[a-zA-Z0-9]{16,128}$" |
| `role` | select | required; values ["legacy_writeback"]; maxSelect 1 |
| `expiry` | date | required |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |
