<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: time_amendments_augmented

`time_amendments_augmented` is a view: a saved query that Turbo reads and never writes. Turbo's screens reach it through PocketBase's API (`/api/collections/time_amendments_augmented/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): time_amendments_augmented

```
// copy the listRule and viewRule from the time_amendments collection api rules
@request.auth.user_claims_via_uid.cid.name ?= 'tame' ||
@request.auth.user_claims_via_uid.cid.name ?= 'report'
```

## Who can view one record (view rule): time_amendments_augmented

```
// copy the listRule and viewRule from the time_amendments collection api rules
@request.auth.user_claims_via_uid.cid.name ?= 'tame' ||
@request.auth.user_claims_via_uid.cid.name ?= 'report'
```

## The query of the view: time_amendments_augmented

```sql
SELECT 
  ta.id,
  ta.division,
  ta.uid,
  ta.hours,
  ta.description,
  ta.time_type,
  ta.meals_hours,
  ta.job,
  ta.work_record,
  ta.payout_request_amount,
  ta.date,
  ta.week_ending,
  ta.tsid,
  ta.category,
  ta.creator,
  ta.committed,
  ta.committer,
  ta.committed_week_ending,
  ta.skip_tsid_check,
  (p0.given_name || ' ' || p0.surname) as uid_name,
  (p1.given_name || ' ' || p1.surname) as creator_name,
  (p2.given_name || ' ' || p2.surname) as committer_name,
  tt.code as time_type_code,
  tt.name as time_type_name,
  j.number as job_number,
  j.description as job_description,
  c.name as category_name,
  d.code as division_code,
  d.name as division_name
FROM time_amendments ta
LEFT JOIN time_types tt ON ta.time_type = tt.id
LEFT JOIN jobs j ON ta.job = j.id
LEFT JOIN categories c ON ta.category = c.id
LEFT JOIN divisions d ON ta.division = d.id
LEFT JOIN profiles p0 ON ta.uid = p0.uid
LEFT JOIN profiles p1 ON ta.creator = p1.uid
LEFT JOIN profiles p2 ON ta.committer = p2.uid;
```

## Fields: time_amendments_augmented

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; pattern "^[a-z0-9]+$" |
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
| `uid_name` | json | maxSize 1 |
| `creator_name` | json | maxSize 1 |
| `committer_name` | json | maxSize 1 |
| `time_type_code` | text | required; min 1 |
| `time_type_name` | text | required; min 2 |
| `job_number` | text | pattern "^(P)?[0-9]{2}-[0-9]{3,4}L?(-[0-9]{1,2})?(-[0-9])?$" |
| `job_description` | text | required; min 3 |
| `category_name` | text | required; min 3 |
| `division_code` | text | required; min 1 |
| `division_name` | text | required; min 2 |
