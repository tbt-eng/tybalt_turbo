<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: time_tracking

`time_tracking` is a view: a saved query that Turbo reads and never writes. Turbo's screens reach it through PocketBase's API (`/api/collections/time_tracking/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): time_tracking

```
@request.auth.user_claims_via_uid.cid.name ?= 'report'
```

## Who can view one record (view rule): time_tracking

```
@request.auth.user_claims_via_uid.cid.name ?= 'report'
```

## The query of the view: time_tracking

```sql
WITH placeholder_expense_counts AS (
  SELECT
    e.committed_week_ending AS week_ending,
    COUNT(*) AS placeholder_payroll_id_expense_count
  FROM expenses e
  LEFT JOIN admin_profiles ap ON ap.uid = e.uid
  WHERE e.committed != ''
    AND e.committed_week_ending != ''
    AND (LENGTH(COALESCE(ap.payroll_id, '')) = 9 AND COALESCE(ap.payroll_id, '') GLOB '9[0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9]')
  GROUP BY e.committed_week_ending
)
SELECT
  MIN(ts.id) id,
  ts.week_ending,
  SUM(CASE WHEN ts.committed != '' THEN 1 ELSE 0 END) committed_count,
  SUM(CASE WHEN ts.approved != '' AND ts.committed = '' AND ts.rejected = '' THEN 1 ELSE 0 END) approved_count,
  SUM(CASE WHEN ts.submitted = 1 AND ts.approved = '' AND ts.committed = '' AND ts.rejected = '' THEN 1 ELSE 0 END) submitted_count,
  COALESCE(MAX(pe.placeholder_payroll_id_expense_count), 0) AS placeholder_payroll_id_expense_count
FROM time_sheets ts
LEFT JOIN placeholder_expense_counts pe
  ON pe.week_ending = ts.week_ending
GROUP BY ts.week_ending
ORDER BY ts.week_ending DESC;
```

## Fields: time_tracking

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; pattern "^[a-z0-9]+$" |
| `week_ending` | text | required; pattern "^\\d{4}-\\d{2}-\\d{2}$" |
| `committed_count` | json | maxSize 1 |
| `approved_count` | json | maxSize 1 |
| `submitted_count` | json | maxSize 1 |
| `placeholder_payroll_id_expense_count` | json | maxSize 1 |
