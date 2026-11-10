<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: payroll_report_week_endings

`payroll_report_week_endings` is a view: a saved query that Turbo reads and never writes. Turbo's screens reach it through PocketBase's API (`/api/collections/payroll_report_week_endings/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): payroll_report_week_endings

```
@request.auth.user_claims_via_uid.cid.name ?= 'report'
```

## Who can view one record (view rule): payroll_report_week_endings

```
@request.auth.user_claims_via_uid.cid.name ?= 'report'
```

## The query of the view: payroll_report_week_endings

```sql
WITH raw_payroll_periods AS (
  SELECT
    week_ending
  FROM time_sheets
  WHERE committed != ''
    AND (CAST(JULIANDAY(week_ending) - JULIANDAY('2025-03-01') AS INTEGER)) % 14 = 0
  GROUP BY week_ending

  UNION

  SELECT
    pay_period_ending AS week_ending
  FROM expenses
  WHERE committed != ''
    AND pay_period_ending != ''
    AND (CAST(JULIANDAY(pay_period_ending) - JULIANDAY('2025-03-01') AS INTEGER)) % 14 = 0
  GROUP BY pay_period_ending
),
payroll_periods AS (
  SELECT
    'p' || REPLACE(raw.week_ending, '-', '') AS id,
    raw.week_ending
  FROM raw_payroll_periods raw
),
timesheet_counts AS (
  SELECT
    CASE
      WHEN (CAST(JULIANDAY(week_ending) - JULIANDAY('2025-03-01') AS INTEGER)) % 14 = 0
        THEN week_ending
      ELSE date(week_ending, '+7 days')
    END AS payroll_week_ending,
    COUNT(*) AS committed_timesheet_count
  FROM time_sheets
  WHERE committed != ''
  GROUP BY 1
),
expense_counts AS (
  SELECT
    pay_period_ending AS payroll_week_ending,
    COUNT(*) AS committed_expense_count
  FROM expenses
  WHERE committed != ''
    AND pay_period_ending != ''
    AND (CAST(JULIANDAY(pay_period_ending) - JULIANDAY('2025-03-01') AS INTEGER)) % 14 = 0
  GROUP BY pay_period_ending
),
placeholder_time_rows AS (
  SELECT
    ts.week_ending AS source_week_ending
  FROM time_sheets ts
  LEFT JOIN admin_profiles ap ON ap.uid = ts.uid
  WHERE ts.committed != ''
    AND (LENGTH(COALESCE(ap.payroll_id, '')) = 9 AND COALESCE(ap.payroll_id, '') GLOB '9[0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9]')

  UNION ALL

  SELECT
    ta.committed_week_ending AS source_week_ending
  FROM time_amendments ta
  LEFT JOIN admin_profiles ap ON ap.uid = ta.uid
  WHERE ta.committed != ''
    AND ta.committed_week_ending != ''
    AND (LENGTH(COALESCE(ap.payroll_id, '')) = 9 AND COALESCE(ap.payroll_id, '') GLOB '9[0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9]')
),
placeholder_time_counts AS (
  SELECT
    CASE
      WHEN (CAST(JULIANDAY(source_week_ending) - JULIANDAY('2025-03-01') AS INTEGER)) % 14 = 0
        THEN source_week_ending
      ELSE date(source_week_ending, '+7 days')
    END AS payroll_week_ending,
    SUM(
      CASE
        WHEN (CAST(JULIANDAY(source_week_ending) - JULIANDAY('2025-03-01') AS INTEGER)) % 14 = 0
          THEN 0
        ELSE 1
      END
    ) AS placeholder_payroll_id_week1_time_count,
    SUM(
      CASE
        WHEN (CAST(JULIANDAY(source_week_ending) - JULIANDAY('2025-03-01') AS INTEGER)) % 14 = 0
          THEN 1
        ELSE 0
      END
    ) AS placeholder_payroll_id_week2_time_count
  FROM placeholder_time_rows
  GROUP BY 1
),
placeholder_expense_counts AS (
  SELECT
    e.pay_period_ending AS payroll_week_ending,
    COUNT(*) AS placeholder_payroll_id_expense_count
  FROM expenses e
  LEFT JOIN admin_profiles ap ON ap.uid = e.uid
  WHERE e.committed != ''
    AND e.pay_period_ending != ''
    AND (CAST(JULIANDAY(e.pay_period_ending) - JULIANDAY('2025-03-01') AS INTEGER)) % 14 = 0
    AND (LENGTH(COALESCE(ap.payroll_id, '')) = 9 AND COALESCE(ap.payroll_id, '') GLOB '9[0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9]')
  GROUP BY e.pay_period_ending
)
SELECT
  p.id,
  p.week_ending,
  COALESCE(t.committed_timesheet_count, 0) AS committed_timesheet_count,
  COALESCE(e.committed_expense_count, 0) AS committed_expense_count,
  COALESCE(pt.placeholder_payroll_id_week1_time_count, 0) AS placeholder_payroll_id_week1_time_count,
  COALESCE(pt.placeholder_payroll_id_week2_time_count, 0) AS placeholder_payroll_id_week2_time_count,
  COALESCE(pe.placeholder_payroll_id_expense_count, 0) AS placeholder_payroll_id_expense_count
FROM payroll_periods p
LEFT JOIN timesheet_counts t
  ON t.payroll_week_ending = p.week_ending
LEFT JOIN expense_counts e
  ON e.payroll_week_ending = p.week_ending
LEFT JOIN placeholder_time_counts pt
  ON pt.payroll_week_ending = p.week_ending
LEFT JOIN placeholder_expense_counts pe
  ON pe.payroll_week_ending = p.week_ending
ORDER BY p.week_ending DESC;
```

## Fields: payroll_report_week_endings

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; pattern "^[a-z0-9]+$" |
| `week_ending` | json | maxSize 1 |
| `committed_timesheet_count` | json | maxSize 1 |
| `committed_expense_count` | json | maxSize 1 |
| `placeholder_payroll_id_week1_time_count` | json | maxSize 1 |
| `placeholder_payroll_id_week2_time_count` | json | maxSize 1 |
| `placeholder_payroll_id_expense_count` | json | maxSize 1 |
