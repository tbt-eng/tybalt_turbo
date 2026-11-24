-- Every time entry of one user, newest first, for the CSV download on the
-- time sheets list page. The only parameter is the caller's uid.
--
-- The status comes from the entry's time sheet and uses the words of the
-- user docs (ui/src/docs/time.md). An entry with no time sheet, or on a
-- recalled sheet, is not submitted. A rejected sheet stays submitted until it
-- is recalled, so rejected is checked before approved and pending.
SELECT
  te.date AS date,
  te.week_ending AS week_ending,
  COALESCE(tt.code, '') AS time_type,
  COALESCE(tt.name, '') AS time_type_name,
  COALESCE(te.hours, 0) AS hours,
  COALESCE(te.meals_hours, 0) AS meals_hours,
  COALESCE(j.number, '') AS job_number,
  COALESCE(j.description, '') AS job_description,
  COALESCE(d.code, '') AS division,
  COALESCE(d.name, '') AS division_name,
  COALESCE(c.name, '') AS category,
  COALESCE(b.code, '') AS branch,
  COALESCE(r.name, '') AS role,
  COALESCE(te.work_record, '') AS work_record,
  COALESCE(te.payout_request_amount, 0) AS payout_request_amount,
  COALESCE(te.description, '') AS description,
  CASE
    WHEN ts.id IS NULL THEN 'not submitted'
    WHEN COALESCE(ts.committed, '') != '' THEN 'committed'
    WHEN COALESCE(ts.rejected, '') != '' THEN 'rejected'
    WHEN COALESCE(ts.approved, '') != '' THEN 'approved'
    WHEN COALESCE(ts.submitted, 0) = 1 THEN 'pending'
    ELSE 'not submitted'
  END AS time_sheet_status
FROM time_entries te
LEFT JOIN time_sheets ts ON ts.id = te.tsid
LEFT JOIN time_types tt ON tt.id = te.time_type
LEFT JOIN jobs j ON j.id = te.job
LEFT JOIN divisions d ON d.id = te.division
LEFT JOIN categories c ON c.id = te.category
LEFT JOIN branches b ON b.id = te.branch
LEFT JOIN rate_roles r ON r.id = te.role
WHERE te.uid = {:uid}
ORDER BY te.date DESC, te.created DESC, te.id DESC
