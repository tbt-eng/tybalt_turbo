<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: clients

`clients` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/clients/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): clients

```
@request.auth.id != ""
```

## Who can view one record (view rule): clients

```
@request.auth.id != ""
```

## Who can create a record (create rule): clients

```
@request.auth.id != "" &&
@request.auth.user_claims_via_uid.cid.name ?= 'job' && @request.auth.admin_profiles_via_uid.active = true
```

## Who can update a record (update rule): clients

```
@request.auth.id != "" &&
@request.auth.user_claims_via_uid.cid.name ?= 'job' && @request.auth.admin_profiles_via_uid.active = true
```

## Who can delete a record (delete rule): clients

```
(@request.auth.id != "" &&
@request.auth.user_claims_via_uid.cid.name ?= 'job' &&
// prevent deletion of clients if there are referencing jobs
@collection.jobs.client != id) && @request.auth.admin_profiles_via_uid.active = true && @collection.job_requests.client != id && @collection.client_invoicing_information.client != id
```

## Fields: clients

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `name` | text | required; min 2 |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |
| `_imported` | bool |  |
| `business_development_lead` | relation | required; maxSelect 1; to `users` |
| `address` | text |  |
| `city` | text |  |
| `province_state` | text |  |
| `postal_code` | text |  |
| `country` | text |  |
| `phone` | text |  |
| `alias` | text | max 120 |

## Indexes: clients

A unique index means that two records can't have the same values in its fields.

- `CREATE UNIQUE INDEX `idx_aXJh3FO` ON `clients` (`name`)`
