<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: request_checklists

`request_checklists` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/request_checklists/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): request_checklists

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can view one record (view rule): request_checklists

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can create a record (create rule): request_checklists

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can update a record (update rule): request_checklists

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can delete a record (delete rule): request_checklists

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Fields: request_checklists

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `version` | number | required; min 1; onlyInt |
| `questions` | json | required |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |
| `scope` | select | required; values ["approval","creation"]; maxSelect 1 |

## Indexes: request_checklists

A unique index means that two records can't have the same values in its fields.

- `CREATE UNIQUE INDEX `idx_request_checklists_scope_version` ON `request_checklists` (`scope`, `version`)`
