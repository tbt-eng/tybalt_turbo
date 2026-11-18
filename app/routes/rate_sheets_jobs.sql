-- The jobs that use or used one revision of a rate sheet, with their clients.
-- current is true when the job uses the revision now (jobs.rate_sheet), and
-- false when only a period of its rate sheet history (job_rate_sheets) does.
SELECT
  j.id AS id,
  j.number AS number,
  j.description AS description,
  COALESCE(c.id, '') AS client_id,
  COALESCE(c.name, '') AS client_name,
  j.rate_sheet = {:id} AS current
FROM jobs j
LEFT JOIN clients c ON c.id = j.client
WHERE j.rate_sheet = {:id}
  OR j.id IN (SELECT h.job FROM job_rate_sheets h WHERE h.rate_sheet = {:id})
ORDER BY client_name COLLATE NOCASE, j.number
