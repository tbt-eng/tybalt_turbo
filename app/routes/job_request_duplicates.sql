-- Candidate jobs for the duplicate-work review of one project request.
-- A job is a candidate when it has the request's client, proposal, parent,
-- client PO or client reference, or when it is the request's parent.
-- Proposals (numbers that start with P) are left out. The Go code decides
-- which candidates are listed. Each row carries the names the page shows and
-- the date of the job's latest time entry, from one grouped query that uses
-- the time_entries (job, date) index.
WITH candidates AS (
  SELECT j.id
  FROM jobs j
  WHERE substr(j.number, 1, 1) != 'P'
    AND (
      ({:client} != '' AND j.client = {:client})
      OR ({:proposal} != '' AND j.proposal = {:proposal})
      OR ({:parent} != '' AND (j.id = {:parent} OR j.parent = {:parent}))
      OR ({:client_po} != '' AND lower(trim(j.client_po)) = {:client_po})
      OR ({:client_reference} != '' AND lower(trim(j.client_reference_number)) = {:client_reference})
    )
),
latest_time AS (
  SELECT te.job, MAX(te.date) AS last_time
  FROM time_entries te
  WHERE te.job IN (SELECT id FROM candidates)
  GROUP BY te.job
)
SELECT
  j.id,
  j.number,
  j.description,
  j.client,
  j.status,
  j.manager,
  j.parent,
  j.proposal,
  j.location,
  j.client_po,
  j.client_reference_number,
  j.project_award_date,
  COALESCE(c.name, '') AS client_name,
  COALESCE(p.number, '') AS parent_number,
  COALESCE(NULLIF(TRIM(COALESCE(pr.given_name, '') || ' ' || COALESCE(pr.surname, '')), ''), u.name, '') AS manager_name,
  COALESCE(lt.last_time, '') AS last_time
FROM candidates cand
JOIN jobs j ON j.id = cand.id
LEFT JOIN clients c ON c.id = j.client
LEFT JOIN jobs p ON j.parent != '' AND p.id = j.parent
LEFT JOIN users u ON j.manager != '' AND u.id = j.manager
LEFT JOIN profiles pr ON j.manager != '' AND pr.uid = j.manager
LEFT JOIN latest_time lt ON lt.job = j.id
