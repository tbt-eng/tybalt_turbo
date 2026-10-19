-- Appended to job_priced_time_entries.sql. Price entries before grouping
-- because an employee can use more than one role on the same job. A missing
-- employee profile gives blank names; the Go row fields cannot hold NULL.
, sheet_labels AS (
  -- The sheets that priced each row's entries. Time in different periods can
  -- use different sheets.
  SELECT uid, GROUP_CONCAT(label, '; ') AS rate_sheets
  FROM (SELECT DISTINCT uid, rate_sheet_name || ' (rev. ' || rate_sheet_revision || ')' AS label
    FROM priced_entries WHERE rate_sheet_name != '' ORDER BY label)
  GROUP BY uid
), base AS (
  SELECT te.number, COALESCE(p.given_name, '') AS given_name, COALESCE(p.surname, '') AS surname, te.uid,
    COALESCE(sl.rate_sheets, '') AS rate_sheets,
    SUM(te.hours) AS hours,
    COALESCE(SUM(te.value), 0) AS value,
    SUM(COALESCE(SUM(te.value), 0)) OVER () AS total,
    SUM(te.meals_hours) AS meals_hours,
    SUM(te.estimated_hours) AS estimated_hours,
    SUM(SUM(te.estimated_hours)) OVER () AS total_estimated_hours,
    SUM(te.unpriced_hours) AS unpriced_hours,
    SUM(SUM(te.unpriced_hours)) OVER () AS total_unpriced_hours
  FROM priced_entries te
  LEFT JOIN sheet_labels sl ON sl.uid IS te.uid
  LEFT JOIN profiles p ON te.uid = p.uid
  GROUP BY te.uid
)
SELECT *,
  CASE WHEN unpriced_hours > 0 THEN 'Incomplete'
       WHEN estimated_hours > 0 THEN 'Estimated' ELSE 'Rate sheet' END AS value_status,
  CASE WHEN total_unpriced_hours > 0 THEN 'Incomplete'
       WHEN total_estimated_hours > 0 THEN 'Estimated' ELSE 'Rate sheet' END AS total_status,
  CASE WHEN total_unpriced_hours = 0 THEN ROUND(value * 100 / NULLIF(total, 0), 1) END AS percent
FROM base
ORDER BY surname, given_name, uid;
