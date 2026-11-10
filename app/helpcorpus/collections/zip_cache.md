<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: zip_cache

`zip_cache` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/zip_cache/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): zip_cache

```
@request.auth.user_claims_via_uid.cid.name ?= 'report'
```

## Who can view one record (view rule): zip_cache

```
@request.auth.user_claims_via_uid.cid.name ?= 'report'
```

## Who can create a record (create rule): zip_cache

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can update a record (update rule): zip_cache

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can delete a record (delete rule): zip_cache

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Fields: zip_cache

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `key` | text | required |
| `class` | text | required |
| `zip` | file | required; maxSelect 1; maxSize 536870912 |
| `manifest` | json |  |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |
