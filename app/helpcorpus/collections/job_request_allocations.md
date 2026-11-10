<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: job_request_allocations

`job_request_allocations` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/job_request_allocations/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): job_request_allocations

```
@request.auth.id != '' && @request.auth.admin_profiles_via_uid.active = true && (request.creator = @request.auth.id || @request.auth.user_claims_via_uid.cid.name ?= 'job' || (request.approver = @request.auth.id && @request.auth.user_claims_via_uid.cid.name ?= 'project_authorization'))
```

## Who can view one record (view rule): job_request_allocations

```
@request.auth.id != '' && @request.auth.admin_profiles_via_uid.active = true && (request.creator = @request.auth.id || @request.auth.user_claims_via_uid.cid.name ?= 'job' || (request.approver = @request.auth.id && @request.auth.user_claims_via_uid.cid.name ?= 'project_authorization'))
```

## Who can create a record (create rule): job_request_allocations

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can update a record (update rule): job_request_allocations

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can delete a record (delete rule): job_request_allocations

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Fields: job_request_allocations

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |
| `request` | relation | required; maxSelect 1; to `job_requests`; cascadeDelete |
| `division` | relation | required; maxSelect 1; to `divisions` |
| `hours` | number | onlyInt |

## Indexes: job_request_allocations

A unique index means that two records can't have the same values in its fields.

- `CREATE UNIQUE INDEX idx_request_division ON job_request_allocations (request,division)`
