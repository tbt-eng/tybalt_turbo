<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: time_off

`time_off` is a view: a saved query that Turbo reads and never writes. Turbo's screens reach it through PocketBase's API (`/api/collections/time_off/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): time_off

```
@request.auth.id = id ||
@request.auth.id = manager_uid ||
@request.auth.user_claims_via_uid.cid.name ?= 'report'
```

## Who can view one record (view rule): time_off

```
@request.auth.id = id ||
@request.auth.id = manager_uid ||
@request.auth.user_claims_via_uid.cid.name ?= 'report'
```

## The query of the view: time_off

```sql
SELECT 
  p.uid as id,
  CAST(CONCAT(p.surname, ', ', p.given_name) AS TEXT) AS name,
  p.manager AS manager_uid,
  CAST(CONCAT(mp.surname, ', ', mp.given_name) AS TEXT) AS manager,
  ap.opening_date,
  ap.opening_ov,
  ap.opening_op,
  ap.salary,
  CAST(COALESCE(SUM(CASE WHEN src.code = 'OV' THEN src.hours ELSE 0 END), 0) AS REAL) AS used_ov,
  CAST(COALESCE(SUM(CASE WHEN src.code = 'OP' THEN src.hours ELSE 0 END), 0) AS REAL) AS used_op,
  CAST(COALESCE(SUM(CASE WHEN src.code = 'OV' AND src.tsid != '' THEN src.hours ELSE 0 END), 0) AS REAL) AS timesheet_ov,
  CAST(COALESCE(SUM(CASE WHEN src.code = 'OP' AND src.tsid != '' THEN src.hours ELSE 0 END), 0) AS REAL) AS timesheet_op,
  CAST(MAX(CASE WHEN src.code = 'OV' THEN src.date END) AS TEXT) AS last_ov,
  CAST(MAX(CASE WHEN src.code = 'OP' THEN src.date END) AS TEXT) AS last_op
FROM 
  admin_profiles ap
JOIN
    profiles p ON p.uid = ap.uid AND ap.untracked_time_off = false AND ap.time_sheet_expected = true
LEFT JOIN
    profiles mp ON p.manager = mp.uid
LEFT JOIN (
  -- time_entries (committed only)
  SELECT 
  te.uid,
  te.hours,
  te.tsid,
  te.date,
  te.week_ending,
  tt.code
  FROM time_entries te
  JOIN time_types tt ON te.time_type = tt.id
  JOIN time_sheets ts ON te.tsid = ts.id
  WHERE te.tsid != '' AND ts.committed != '' AND tt.code IN ('OV','OP')
  UNION ALL
  -- time_amendments (committed only)
  SELECT 
  ta.uid,
  ta.hours,
  IFNULL(ta.tsid, '') AS tsid,
  ta.date,
  ta.committed_week_ending AS week_ending,
  tt2.code
  FROM time_amendments ta
  JOIN time_types tt2 ON ta.time_type = tt2.id
  WHERE ta.committed != '' AND tt2.code IN ('OV','OP')
) AS src 
  ON p.uid = src.uid
  AND src.week_ending > ap.opening_date
GROUP BY 
  p.uid, p.surname, p.given_name, p.manager, mp.surname, mp.given_name, ap.opening_date, ap.opening_ov, ap.opening_op, ap.salary
```

## Fields: time_off

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; pattern "^[a-z0-9]+$" |
| `name` | text |  |
| `manager_uid` | relation | required; maxSelect 1; to `users` |
| `manager` | text |  |
| `opening_date` | text | pattern "^\\d{4}-\\d{2}-\\d{2}$" |
| `opening_ov` | number | max 200 |
| `opening_op` | number | max 332 |
| `salary` | bool |  |
| `used_ov` | number |  |
| `used_op` | number |  |
| `timesheet_ov` | number |  |
| `timesheet_op` | number |  |
| `last_ov` | text |  |
| `last_op` | text |  |
