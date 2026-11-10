<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: client_invoicing_information

`client_invoicing_information` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/client_invoicing_information/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): client_invoicing_information

```
@request.auth.id != '' && @request.auth.admin_profiles_via_uid.active = true
```

## Who can view one record (view rule): client_invoicing_information

```
@request.auth.id != '' && @request.auth.admin_profiles_via_uid.active = true
```

## Who can create a record (create rule): client_invoicing_information

```
@request.auth.id != '' && @request.auth.admin_profiles_via_uid.active = true && (@request.auth.user_claims_via_uid.cid.name ?= 'job' || @request.auth.user_claims_via_uid.cid.name ?= 'time')
```

## Who can update a record (update rule): client_invoicing_information

```
@request.auth.id != '' && @request.auth.admin_profiles_via_uid.active = true && @request.auth.user_claims_via_uid.cid.name ?= 'job'
```

## Who can delete a record (delete rule): client_invoicing_information

```
@request.auth.id != '' && @request.auth.admin_profiles_via_uid.active = true && @request.auth.user_claims_via_uid.cid.name ?= 'job' && @collection.jobs.invoicing_information != id && @collection.job_requests.invoicing_information != id
```

## Fields: client_invoicing_information

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |
| `client` | relation | required; maxSelect 1; to `clients` |
| `contact` | relation | required; maxSelect 1; to `client_contacts` |
| `creator` | relation | maxSelect 1; to `users` |
| `name` | text | max 120 |
| `invoicing_instructions` | text |  |
| `fax` | text |  |
| `billing_name` | select | values ["name","alias"]; maxSelect 1 |

## Indexes: client_invoicing_information

A unique index means that two records can't have the same values in its fields.

- `CREATE INDEX idx_invoice_client ON client_invoicing_information (client)`
- `CREATE INDEX idx_invoice_contact ON client_invoicing_information (contact)`
