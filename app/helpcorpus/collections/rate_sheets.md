<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: rate_sheets

`rate_sheets` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/rate_sheets/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): rate_sheets

```
@request.auth.id != "" && request = ""
```

## Who can view one record (view rule): rate_sheets

```
@request.auth.id != "" && request = ""
```

## Who can create a record (create rule): rate_sheets

```
@request.auth.id != "" &&
// A draft is made only through a project request.
@request.body.request:isset = false &&
@request.body.based_on:isset = false &&
@request.body.reason:isset = false && (
  @request.auth.user_claims_via_uid.cid.name ?= 'rate_sheet_revise' ||
  @request.auth.user_claims_via_uid.cid.name ?= 'admin' ||
  (@request.auth.user_claims_via_uid.cid.name ?= 'job' && @request.body.revision = 0)
)
```

## Who can update a record (update rule): rate_sheets

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can delete a record (delete rule): rate_sheets

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Fields: rate_sheets

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `effective_date` | text | required; pattern "^\\d{4}-\\d{2}-\\d{2}$" |
| `name` | text | required |
| `revision` | number | onlyInt |
| `active` | bool |  |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |
| `request` | relation | maxSelect 1; to `job_requests` |
| `reason` | text | max 500 |
| `based_on` | relation | maxSelect 1; to `rate_sheets` |

## Indexes: rate_sheets

A unique index means that two records can't have the same values in its fields.

- `CREATE UNIQUE INDEX `idx_dHsmkOUme0` ON `rate_sheets` ( `name`, `revision` )`
- `CREATE INDEX `idx_rate_sheets_request` ON `rate_sheets` (`request`)`
