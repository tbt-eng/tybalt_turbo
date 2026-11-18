-- Every rate sheet revision that a job of this client uses now (jobs.rate_sheet)
-- or that a period of a job's rate sheet history (job_rate_sheets) names. Jobs
-- where the client is only the job owner do not count. current_jobs counts the
-- client's jobs on the revision now; past_jobs counts the ones that only a
-- period of their history puts on it. A request's proposed rates are left out,
-- as the rate_sheets read rule leaves them out. Sheets with current jobs come
-- first, then the newest code and revision.
SELECT
  rs.id AS id,
  rs.code AS code,
  rs.name AS name,
  rs.revision AS revision,
  rs.effective_date AS effective_date,
  rs.active AS active,
  COUNT(DISTINCT CASE WHEN u.current = 1 THEN u.job END) AS current_jobs,
  COUNT(DISTINCT u.job) - COUNT(DISTINCT CASE WHEN u.current = 1 THEN u.job END) AS past_jobs
FROM (
  SELECT j.id AS job, j.rate_sheet AS rate_sheet, 1 AS current
  FROM jobs j
  WHERE j.client = {:id} AND j.rate_sheet != ''
  UNION
  SELECT h.job AS job, h.rate_sheet AS rate_sheet, 0 AS current
  FROM jobs j
  JOIN job_rate_sheets h ON h.job = j.id
  WHERE j.client = {:id}
) u
JOIN rate_sheets rs ON rs.id = u.rate_sheet
WHERE COALESCE(rs.request, '') = ''
GROUP BY rs.id
ORDER BY current_jobs > 0 DESC, rs.code DESC, rs.revision DESC
