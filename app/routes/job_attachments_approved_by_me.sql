-- Project values the signed-in person approved that are now on jobs, from the
-- last 30 days, newest first. A new job counts from when the job was created,
-- since its request can be approved long before. A later value change counts
-- from its approval.
WITH approvals AS (
  SELECT
    a.id,
    a.job,
    j.number AS job_number,
    j.description AS job_description,
    a.project_value,
    CASE WHEN a.initial_approval THEN 'new_job' ELSE 'value_change' END AS kind,
    CASE WHEN a.initial_approval THEN j.created ELSE a.approved END AS date
  FROM job_attachments a
  JOIN jobs j ON j.id = a.job
  WHERE a.job != ''
    AND a.status = 'Approved'
    AND a.approver = {:uid}
)
SELECT id, job, job_number, job_description, project_value, kind, date
FROM approvals
WHERE date >= {:cutoff}
ORDER BY date DESC, id
