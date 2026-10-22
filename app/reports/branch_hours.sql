-- Branch Hours: one row per job, employee and work branch for the time that
-- involves the selected branch. Time qualifies by the same rules as
-- employee_branch_hours.sql, so a branch's totals reconcile with that report.
-- The staff branch is the employee's current default branch. The work branch
-- is the branch on the time entry.
WITH qualifying_time AS (
  SELECT
    te.uid,
    te.branch AS work_branch,
    ap.default_branch AS staff_branch,
    COALESCE(te.job, '') AS job,
    CAST(te.hours AS REAL) AS hours
  FROM time_entries te
  INNER JOIN time_types tt
    ON tt.id = te.time_type
  INNER JOIN branches entry_branch
    ON entry_branch.id = te.branch
  INNER JOIN admin_profiles ap
    ON ap.uid = te.uid
  WHERE tt.code IN ('R', 'RT')
    AND te.date >= {:start_date}
    AND te.date <= {:end_date}
    AND (te.branch = {:branch} OR ap.default_branch = {:branch})
)
SELECT
  CASE
    WHEN qt.staff_branch = {:branch} AND qt.work_branch = {:branch} THEN 'internal'
    WHEN qt.staff_branch = {:branch} THEN 'supplied'
    ELSE 'received'
  END AS direction,
  qt.job AS jobId,
  COALESCE(j.number, '') AS jobNumber,
  COALESCE(j.description, '') AS jobDescription,
  COALESCE(c.name, '') AS client,
  work_branch.name AS workBranch,
  qt.uid AS uid,
  COALESCE(ap.payroll_id, '') AS payrollId,
  COALESCE(p.surname, '') AS surname,
  COALESCE(p.given_name, '') AS givenName,
  COALESCE(staff_branch.name, '') AS staffBranch,
  SUM(qt.hours) AS hours
FROM qualifying_time qt
INNER JOIN branches work_branch
  ON work_branch.id = qt.work_branch
INNER JOIN admin_profiles ap
  ON ap.uid = qt.uid
INNER JOIN profiles p
  ON p.uid = qt.uid
LEFT JOIN branches staff_branch
  ON staff_branch.id = qt.staff_branch
LEFT JOIN jobs j
  ON j.id = qt.job
LEFT JOIN clients c
  ON c.id = j.client
GROUP BY
  qt.job,
  qt.uid,
  qt.work_branch,
  qt.staff_branch
ORDER BY
  direction,
  jobNumber = '',
  jobNumber,
  work_branch.name,
  p.surname,
  p.given_name,
  qt.uid
