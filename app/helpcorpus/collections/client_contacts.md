<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: client_contacts

`client_contacts` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/client_contacts/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): client_contacts

```
@request.auth.id != ""
```

## Who can view one record (view rule): client_contacts

```
@request.auth.id != ""
```

## Who can create a record (create rule): client_contacts

```
@request.auth.id != "" &&
@request.auth.user_claims_via_uid.cid.name ?= 'job' && @request.auth.admin_profiles_via_uid.active = true
```

## Who can update a record (update rule): client_contacts

```
@request.auth.id != "" &&
@request.auth.user_claims_via_uid.cid.name ?= 'job' && @request.auth.admin_profiles_via_uid.active = true
```

## Who can delete a record (delete rule): client_contacts

```
(@request.auth.id != "" &&
@request.auth.user_claims_via_uid.cid.name ?= 'job' &&
// prevent deletion of contacts if there are referencing jobs
@collection.jobs.contact != id) && @request.auth.admin_profiles_via_uid.active = true && @collection.job_requests.contact != id && @collection.client_invoicing_information.contact != id
```

## Fields: client_contacts

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `surname` | text | required |
| `given_name` | text | required |
| `email` | email |  |
| `client` | relation | required; maxSelect 1; to `clients`; cascadeDelete |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |
| `_imported` | bool |  |
| `address` | text |  |
| `city` | text |  |
| `province_state` | text |  |
| `postal_code` | text |  |
| `country` | text |  |
| `phone` | text |  |
| `no_email` | bool |  |

## Indexes: client_contacts

A unique index means that two records can't have the same values in its fields.

- `CREATE INDEX `idx_KxKk01Y` ON `client_contacts` ( `surname`, `given_name` )`
- `CREATE UNIQUE INDEX `idx_0tGweKG7kJ` ON `client_contacts` (`email`) WHERE email != ''`
