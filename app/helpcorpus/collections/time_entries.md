<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: time_entries

`time_entries` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/time_entries/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): time_entries

```
@request.auth.id = uid ||
(tsid.submitted = true && @request.auth.id = tsid.approver) ||
(tsid.submitted = true && @request.auth.id ?= tsid.time_sheet_reviewers_via_time_sheet.reviewer) ||
(tsid.submitted = true && tsid.approved != '' && tsid.committed = '' && @request.auth.user_claims_via_uid.cid.name ?= 'commit') ||
(tsid.committed != '' && @request.auth.user_claims_via_uid.cid.name ?= 'report')
```

## Who can view one record (view rule): time_entries

```
@request.auth.id = uid ||
(tsid.submitted = true && @request.auth.id = tsid.approver) ||
(tsid.submitted = true && @request.auth.id ?= tsid.time_sheet_reviewers_via_time_sheet.reviewer) ||
(tsid.submitted = true && tsid.approved != '' && tsid.committed = '' && @request.auth.user_claims_via_uid.cid.name ?= 'commit') ||
(tsid.committed != '' && @request.auth.user_claims_via_uid.cid.name ?= 'report')
```

## Who can create a record (create rule): time_entries

```
@request.auth.id != "" &&
@request.auth.user_claims_via_uid.cid.name ?= 'time' &&
// if present, the category belongs to the job, otherwise is blank
(
  // compare the new category to the new job
  ( @request.body.job:isset = true && @request.body.category.job = @request.body.job ) ||
  @request.body.category = ""
)
```

## Who can update a record (update rule): time_entries

```
// the creating user can edit if the entry is not yet part of a timesheet
uid = @request.auth.id &&
@request.auth.user_claims_via_uid.cid.name ?= 'time' &&
tsid = "" &&

// uid must not change after create
@request.body.uid:changed = false &&

// if present, the category belongs to the job, otherwise is blank
(
  // the job is unchanged, compare the new category to job
  ( @request.body.job:isset = false && @request.body.category.job = job ) ||
  // the job has changed, compare the new category to the new job
  ( @request.body.job:isset = true && @request.body.category.job = @request.body.job ) ||
  @request.body.category = ""
)
```

## Who can delete a record (delete rule): time_entries

```
@request.auth.id = uid &&
@request.auth.user_claims_via_uid.cid.name ?= 'time' &&
tsid = ""
```

## Fields: time_entries

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `division` | relation | maxSelect 1; to `divisions` |
| `uid` | relation | required; maxSelect 1; to `users` |
| `hours` | number | max 18 |
| `description` | text |  |
| `time_type` | relation | required; maxSelect 1; to `time_types` |
| `meals_hours` | number | max 3 |
| `job` | relation | maxSelect 1; to `jobs` |
| `work_record` | text | pattern "^[FKQ][0-9]{2}-[0-9]{3,}(-[0-9]+)?$" |
| `work_record_subject_id` | relation | maxSelect 1; to `work_records_subjects` |
| `payout_request_amount` | number |  |
| `date` | text | required; pattern "^\\d{4}-\\d{2}-\\d{2}$" |
| `week_ending` | text | required; pattern "^\\d{4}-\\d{2}-\\d{2}$" |
| `tsid` | relation | maxSelect 1; to `time_sheets` |
| `category` | relation | maxSelect 1; to `categories` |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |
| `_imported` | bool |  |
| `branch` | relation | maxSelect 1; to `branches` |
| `role` | relation | maxSelect 1; to `rate_roles` |

## Indexes: time_entries

A unique index means that two records can't have the same values in its fields.

- `CREATE INDEX `idx_jgvQezNmMn` ON `time_entries` (`uid`)`
- `CREATE INDEX `idx_ljFGUYCrIB` ON `time_entries` ( `branch`, `job` )`
- `CREATE INDEX `idx_7JBLPOOySg` ON `time_entries` (`tsid`)`
- `CREATE INDEX `idx_0Htw2Bikbl` ON `time_entries` (`week_ending`)`
- `CREATE INDEX `idx_9QORseSHUO` ON `time_entries` ( `job`, `date DESC` )`
- `CREATE INDEX `idx_9fjDrGa7M7` ON `time_entries` (`division`)`
- `CREATE INDEX `idx_ZJij6mtUY4` ON `time_entries` (`time_type`)`
- `CREATE INDEX `idx_EXtVt4doe6` ON `time_entries` (`branch`)`
- `CREATE INDEX `idx_oUAVmzP0Gf` ON `time_entries` (`category`)`
- `CREATE INDEX `idx_a7yXF8hMf0` ON `time_entries` ( `tsid`, `time_type` )`
- `CREATE INDEX `idx_ZX8ABROs07` ON `time_entries` ( `uid`, `week_ending` )`
- `CREATE INDEX `idx_8f2UYSxtZu` ON `time_entries` (`work_record`)`
- `CREATE UNIQUE INDEX idx_time_entries_work_record_subject_id ON time_entries (work_record_subject_id) WHERE work_record_subject_id != ''`
- `CREATE INDEX `idx_time_entries_time_type_date` ON `time_entries` (`time_type`, `date`)`
- `CREATE INDEX `idx_time_entries_job_division_date_work` ON `time_entries` (`job`, `division`, `date`) WHERE `hours` > 0`
