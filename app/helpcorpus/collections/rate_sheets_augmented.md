<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: rate_sheets_augmented

`rate_sheets_augmented` is a view: a saved query that Turbo reads and never writes. Turbo's screens reach it through PocketBase's API (`/api/collections/rate_sheets_augmented/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): rate_sheets_augmented

```
@request.auth.id != ""
```

## Who can view one record (view rule): rate_sheets_augmented

```
@request.auth.id != ""
```

## The query of the view: rate_sheets_augmented

```sql
SELECT 
  rs.id AS id,
  rs.code AS code,
  rs.name AS name,
  rs.effective_date AS effective_date,
  rs.revision AS revision,
  rs.active AS active,
  rs.created AS created,
  rs.updated AS updated,
  COUNT(u.job) AS job_count,
  CAST(COALESCE(f.job_numbers, '') AS TEXT) AS job_numbers,
  CAST(COALESCE(f.client_names, '') AS TEXT) AS client_names
FROM rate_sheets rs
LEFT JOIN (
  SELECT j.rate_sheet AS rate_sheet, j.id AS job FROM jobs j WHERE j.rate_sheet != ''
  UNION
  SELECT h.rate_sheet AS rate_sheet, h.job AS job FROM job_rate_sheets h
) u ON u.rate_sheet = rs.id
LEFT JOIN (
  SELECT d.code AS code,
    GROUP_CONCAT(CASE WHEN d.kind = 'job' THEN d.value END, char(10)) AS job_numbers,
    GROUP_CONCAT(CASE WHEN d.kind = 'client' THEN d.value END, char(10)) AS client_names
  FROM (
    SELECT s.code AS code, 'job' AS kind, cj.number AS value
    FROM rate_sheets s
    JOIN (
      SELECT j.rate_sheet AS rate_sheet, j.id AS job FROM jobs j WHERE j.rate_sheet != ''
      UNION
      SELECT h.rate_sheet AS rate_sheet, h.job AS job FROM job_rate_sheets h
    ) cu ON cu.rate_sheet = s.id
    JOIN jobs cj ON cj.id = cu.job
    UNION
    SELECT s.code AS code, 'client' AS kind, c.name AS value
    FROM rate_sheets s
    JOIN (
      SELECT j.rate_sheet AS rate_sheet, j.id AS job FROM jobs j WHERE j.rate_sheet != ''
      UNION
      SELECT h.rate_sheet AS rate_sheet, h.job AS job FROM job_rate_sheets h
    ) cu ON cu.rate_sheet = s.id
    JOIN jobs cj ON cj.id = cu.job
    JOIN clients c ON c.id = cj.client
  ) d
  GROUP BY d.code
) f ON f.code = rs.code
WHERE COALESCE(rs.request, '') = ''
GROUP BY rs.id, rs.code, rs.name, rs.effective_date, rs.revision, rs.active, rs.created, rs.updated, f.job_numbers, f.client_names
```

## Fields: rate_sheets_augmented

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; pattern "^[a-z0-9]+$" |
| `code` | text | required; pattern "^\\d{4}-(0[1-9]\|1[0-2])-\\d{5}$" |
| `name` | text |  |
| `effective_date` | text | required; pattern "^\\d{4}-\\d{2}-\\d{2}$" |
| `revision` | number | onlyInt |
| `active` | bool |  |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |
| `job_count` | number | onlyInt |
| `job_numbers` | text |  |
| `client_names` | text |  |
