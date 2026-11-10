<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: job_requests

`job_requests` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/job_requests/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): job_requests

```
@request.auth.id != '' && @request.auth.admin_profiles_via_uid.active = true && (creator = @request.auth.id || @request.auth.user_claims_via_uid.cid.name ?= 'job' || (approver = @request.auth.id && @request.auth.user_claims_via_uid.cid.name ?= 'project_authorization'))
```

## Who can view one record (view rule): job_requests

```
@request.auth.id != '' && @request.auth.admin_profiles_via_uid.active = true && (creator = @request.auth.id || @request.auth.user_claims_via_uid.cid.name ?= 'job' || (approver = @request.auth.id && @request.auth.user_claims_via_uid.cid.name ?= 'project_authorization'))
```

## Who can create a record (create rule): job_requests

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can update a record (update rule): job_requests

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can delete a record (delete rule): job_requests

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Fields: job_requests

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |
| `creator` | relation | required; maxSelect 1; to `users` |
| `status` | select | required; values ["Draft","Awaiting approval","Awaiting creation","Rejected","Recalled"]; maxSelect 1 |
| `project_status` | select | values ["Active","Closed"]; maxSelect 1 |
| `version` | number | onlyInt |
| `submission_token` | text | max 32; pattern "^[a-z0-9]{32}$" |
| `approver` | relation | maxSelect 1; to `users` |
| `submitted` | date |  |
| `rejected` | date |  |
| `last_edited` | date |  |
| `rejector` | relation | maxSelect 1; to `users` |
| `rejection_reason` | text |  |
| `client` | relation | maxSelect 1; to `clients` |
| `job_owner` | relation | maxSelect 1; to `clients` |
| `contact` | relation | maxSelect 1; to `client_contacts` |
| `manager` | relation | maxSelect 1; to `users` |
| `alternate_manager` | relation | maxSelect 1; to `users` |
| `branch` | relation | maxSelect 1; to `branches` |
| `rate_sheet` | relation | maxSelect 1; to `rate_sheets` |
| `parent` | relation | maxSelect 1; to `jobs` |
| `proposal` | relation | maxSelect 1; to `jobs` |
| `invoicing_information` | relation | maxSelect 1; to `client_invoicing_information` |
| `description` | text |  |
| `location` | text |  |
| `client_po` | text |  |
| `client_reference_number` | text |  |
| `project_award_date` | text |  |
| `project_completion_date` | text |  |
| `project_value` | number | onlyInt |
| `time_and_materials` | bool |  |
| `time_and_materials_set` | bool |  |
| `fn_agreement_answer` | select | values ["Yes","No"]; maxSelect 1 |
| `notes` | text | min 10; max 1000 |

## Indexes: job_requests

A unique index means that two records can't have the same values in its fields.

- `CREATE INDEX idx_requests_creator_status ON job_requests (creator,status)`
- `CREATE INDEX idx_requests_approver_status ON job_requests (approver,status)`
- `CREATE INDEX idx_requests_status_submitted ON job_requests (status,submitted)`
- `CREATE INDEX idx_requests_proposal ON job_requests (proposal)`
- `CREATE INDEX idx_requests_parent ON job_requests (parent)`
- `CREATE INDEX idx_requests_client ON job_requests (client)`
- `CREATE INDEX idx_requests_last_edited ON job_requests (last_edited)`
