<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: time_amendments

`time_amendments` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/time_amendments/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): time_amendments

```
@request.auth.user_claims_via_uid.cid.name ?= 'tame' ||
@request.auth.user_claims_via_uid.cid.name ?= 'report'
```

## Who can view one record (view rule): time_amendments

```
@request.auth.user_claims_via_uid.cid.name ?= 'tame' ||
@request.auth.user_claims_via_uid.cid.name ?= 'report'
```

## Who can create a record (create rule): time_amendments

```
@request.auth.user_claims_via_uid.cid.name ?= 'tame' &&
// no tsid is submitted, it's set in the hook
(@request.body.tsid:isset = false || @request.body.tsid = "") &&

// no committed properties are submitted
@request.body.committed:isset = false &&
@request.body.committer:isset = false &&
@request.body.committed_week_ending:isset = false &&

// if present, the category belongs to the job, otherwise is blank
(
  // compare the new category to the new job
  ( @request.body.job:isset = true && @request.body.category.job = @request.body.job ) ||
  @request.body.category = ""
)
```

## Who can update a record (update rule): time_amendments

```
@request.auth.id != "" &&
@request.auth.user_claims_via_uid.cid.name ?= 'tame' &&
committed = "" &&

// subject employee must remain stable after create
@request.body.uid:changed = false &&

// no tsid is submitted, it's set in the hook
(@request.body.tsid:isset = false || tsid = @request.body.tsid) &&

// no committed properties are submitted
(@request.body.committed:isset = false || committed = @request.body.committed) &&
(@request.body.committer:isset = false || committer = @request.body.committer) &&
(@request.body.committed_week_ending:isset = false || committed_week_ending = @request.body.committed_week_ending) &&

// if present, the category belongs to the job, otherwise is blank
(
  // the job is unchanged, compare the new category to job
  ( @request.body.job:isset = false && @request.body.category.job = job ) ||
  // the job has changed, compare the new category to the new job
  ( @request.body.job:isset = true && @request.body.category.job = @request.body.job ) ||
  @request.body.category = ""
)
```

## Who can delete a record (delete rule): time_amendments

```
@request.auth.user_claims_via_uid.cid.name ?= 'tame' &&
committed = ""
```

## Fields: time_amendments

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `division` | relation | maxSelect 1; to `divisions` |
| `uid` | relation | required; maxSelect 1; to `users` |
| `hours` | number | required; min -18; max 18 |
| `description` | text |  |
| `time_type` | relation | required; maxSelect 1; to `time_types` |
| `meals_hours` | number | max 3 |
| `job` | relation | maxSelect 1; to `jobs` |
| `work_record` | text | pattern "^[FKQ][0-9]{2}-[0-9]{3,}(-[0-9]+)?$" |
| `payout_request_amount` | number |  |
| `date` | text | required; pattern "^\\d{4}-\\d{2}-\\d{2}$" |
| `week_ending` | text | required; pattern "^\\d{4}-\\d{2}-\\d{2}$" |
| `tsid` | relation | maxSelect 1; to `time_sheets` |
| `category` | relation | maxSelect 1; to `categories` |
| `creator` | relation | required; maxSelect 1; to `users` |
| `committed` | date |  |
| `committer` | relation | maxSelect 1; to `users` |
| `committed_week_ending` | text | pattern "^\\d{4}-\\d{2}-\\d{2}$" |
| `skip_tsid_check` | bool |  |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |
| `_imported` | bool |  |
| `branch` | relation | required; maxSelect 1; to `branches` |

## Indexes: time_amendments

A unique index means that two records can't have the same values in its fields.

- `CREATE INDEX `idx_Yq0DnuvO5A` ON `time_amendments` ( `uid`, `committed_week_ending` ) WHERE committed != ''`
- `CREATE INDEX `idx_time_amendments_job` ON `time_amendments` (`job`)`
- `CREATE INDEX `idx_time_amendments_job_date` ON `time_amendments` ( `job`, `date` )`
