<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: jobs

`jobs` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/jobs/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): jobs

```
@request.auth.id != ""
```

## Who can view one record (view rule): jobs

```
@request.auth.id != ""
```

## Who can create a record (create rule): jobs

```
@request.auth.id != "" &&
@request.auth.user_claims_via_uid.cid.name ?= 'job'
```

## Who can update a record (update rule): jobs

```
@request.auth.id != '' && @request.auth.admin_profiles_via_uid.active = true && (@request.auth.user_claims_via_uid.cid.name ?= 'job' || (number ~ 'P%' && (manager = @request.auth.id || alternate_manager = @request.auth.id))) && @request.body.number:changed = false && (number ~ 'P%' || (@request.body.client:changed = false && @request.body.project_value:changed = false))
```

## Who can delete a record (delete rule): jobs

```
(@request.auth.id != "" &&
@request.auth.user_claims_via_uid.cid.name ?= 'admin' &&

// prevent deletion of jobs if there are referencing time_entries
@collection.time_entries.job != id &&

// prevent deletion of jobs if there are referencing purchase orders
@collection.purchase_orders.job != id &&

// prevent deletion of jobs if there are referencing expenses
@collection.expenses.job != id) && @collection.job_attachments.job != id
```

## Fields: jobs

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `number` | text | pattern "^(P)?[0-9]{2}-[0-9]{3,4}L?(-[0-9]{1,2})?(-[0-9])?$" |
| `description` | text | required; min 3 |
| `client` | relation | required; maxSelect 1; to `clients` |
| `contact` | relation | maxSelect 1; to `client_contacts` |
| `manager` | relation | maxSelect 1; to `users` |
| `alternate_manager` | relation | maxSelect 1; to `users` |
| `fn_agreement` | bool |  |
| `status` | select | values ["Active","Closed","Cancelled","Awarded","Not Awarded","Submitted","In Progress","No Bid"]; maxSelect 1 |
| `project_award_date` | text | pattern "^\\d{4}-\\d{2}-\\d{2}$" |
| `proposal_opening_date` | text | pattern "^\\d{4}-\\d{2}-\\d{2}$" |
| `proposal_submission_due_date` | text | pattern "^\\d{4}-\\d{2}-\\d{2}$" |
| `proposal` | relation | maxSelect 1; to `jobs` |
| `job_owner` | relation | maxSelect 1; to `clients` |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |
| `_imported` | bool |  |
| `branch` | relation | required; maxSelect 1; to `branches` |
| `location` | text | pattern "^[23456789CFGHJMPQRVWX]{8}\\+[23456789CFGHJMPQRVWX]{2,3}$" |
| `outstanding_balance` | number |  |
| `outstanding_balance_date` | text | pattern "^\\d{4}-\\d{2}-\\d{2}$" |
| `parent` | relation | maxSelect 1; to `jobs` |
| `authorizing_document` | select | values ["Unauthorized","PO","PA"]; maxSelect 1 |
| `client_po` | text | max 64 |
| `client_reference_number` | text |  |
| `proposal_value` | number | onlyInt |
| `time_and_materials` | bool |  |
| `rate_sheet` | relation | maxSelect 1; to `rate_sheets` |
| `project_value` | number | onlyInt |
| `project_completion_date` | text | pattern "^\\d{4}-\\d{2}-\\d{2}$" |
| `invoicing_information` | relation | maxSelect 1; to `client_invoicing_information` |
| `request_creator` | relation | maxSelect 1; to `users` |
| `created_by` | relation | maxSelect 1; to `users` |
| `request_created` | date |  |
| `request_submitted` | date |  |
| `creation_completed` | date |  |

## Indexes: jobs

A unique index means that two records can't have the same values in its fields.

- `CREATE UNIQUE INDEX `idx_V1RKd7H` ON `jobs` (`number`)`
- `CREATE INDEX `idx_d1R7JSCuuJ` ON `jobs` (`proposal`)`
- `CREATE INDEX `idx_SCUCvwyln3` ON `jobs` (`parent`)`
- `CREATE INDEX `idx_jobs_status` ON `jobs` (`status`)`
- `CREATE INDEX `idx_jobs_invoicing_information` ON `jobs` (invoicing_information)`
- `CREATE INDEX `idx_jobs_client` ON `jobs` (`client`)`
