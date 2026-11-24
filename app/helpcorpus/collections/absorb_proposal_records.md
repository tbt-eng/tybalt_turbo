<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: absorb_proposal_records

`absorb_proposal_records` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/absorb_proposal_records/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): absorb_proposal_records

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can view one record (view rule): absorb_proposal_records

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can create a record (create rule): absorb_proposal_records

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can update a record (update rule): absorb_proposal_records

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can delete a record (delete rule): absorb_proposal_records

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Fields: absorb_proposal_records

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `proposal` | relation | required; maxSelect 1; to `absorb_proposals`; cascadeDelete |
| `collection_name` | select | required; values ["clients","client_contacts","vendors"]; maxSelect 1 |
| `record_id` | text | required |
| `record_name` | text |  |
| `open` | bool |  |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |

## Indexes: absorb_proposal_records

A unique index means that two records can't have the same values in its fields.

- `CREATE UNIQUE INDEX `idx_absorb_proposal_records_open` ON `absorb_proposal_records` (`collection_name`, `record_id`) WHERE `open` = 1`
- `CREATE INDEX `idx_absorb_proposal_records_proposal` ON `absorb_proposal_records` (`proposal`)`
