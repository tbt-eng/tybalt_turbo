<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: job_attachments

`job_attachments` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/job_attachments/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): job_attachments

```
@request.auth.id != '' && @request.auth.admin_profiles_via_uid.active = true && ((request != '' && (request.creator = @request.auth.id || request.approver = @request.auth.id && @request.auth.user_claims_via_uid.cid.name ?= 'project_authorization' || @request.auth.user_claims_via_uid.cid.name ?= 'job')) || (job != '' && (@request.auth.user_claims_via_uid.cid.name ?= 'job' || @request.auth.user_claims_via_uid.cid.name ?= 'admin' || assigned_approver = @request.auth.id && @request.auth.user_claims_via_uid.cid.name ?= 'project_authorization')))
```

## Who can view one record (view rule): job_attachments

```
@request.auth.id != '' && @request.auth.admin_profiles_via_uid.active = true && ((request != '' && (request.creator = @request.auth.id || request.approver = @request.auth.id && @request.auth.user_claims_via_uid.cid.name ?= 'project_authorization' || @request.auth.user_claims_via_uid.cid.name ?= 'job')) || (job != '' && (@request.auth.user_claims_via_uid.cid.name ?= 'job' || @request.auth.user_claims_via_uid.cid.name ?= 'admin' || assigned_approver = @request.auth.id && @request.auth.user_claims_via_uid.cid.name ?= 'project_authorization')))
```

## Who can create a record (create rule): job_attachments

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can update a record (update rule): job_attachments

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can delete a record (delete rule): job_attachments

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Fields: job_attachments

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |
| `request` | relation | maxSelect 1; to `job_requests` |
| `job` | relation | maxSelect 1; to `jobs` |
| `document` | file | required; maxSelect 1; maxSize 20971520; mimeTypes ["application/pdf"] |
| `document_hash` | text | required; min 64; max 64; pattern "^[a-f0-9]{64}$" |
| `uploader` | relation | required; maxSelect 1; to `users` |
| `uploaded` | date |  |
| `project_value` | number | onlyInt |
| `status` | select | required; values ["Pending","Approved","Rejected"]; maxSelect 1 |
| `assigned_approver` | relation | maxSelect 1; to `users` |
| `approver` | relation | maxSelect 1; to `users` |
| `approved` | date |  |
| `previous_project_value` | number | onlyInt |
| `approval_max_value` | number | onlyInt |
| `request_version` | number | onlyInt |
| `approval_sequence` | number | onlyInt |
| `version` | number | onlyInt |
| `initial_approval` | bool |  |
| `checklist_version` | number | onlyInt |
| `rejector` | relation | maxSelect 1; to `users` |
| `rejected` | date |  |
| `rejection_reason` | text |  |
| `review_token` | text |  |
| `creation_checklist_version` | number | onlyInt |

## Indexes: job_attachments

A unique index means that two records can't have the same values in its fields.

- `CREATE UNIQUE INDEX idx_attachment_hash ON job_attachments (document_hash)`
- `CREATE UNIQUE INDEX idx_attachment_request ON job_attachments (request) WHERE request != ''`
- `CREATE UNIQUE INDEX idx_attachment_pending_job ON job_attachments (job) WHERE job != '' AND status != 'Approved'`
- `CREATE UNIQUE INDEX idx_attachment_sequence ON job_attachments (job,approval_sequence) WHERE job != '' AND status = 'Approved'`
- `CREATE INDEX idx_attachment_queue ON job_attachments (assigned_approver,status)`
