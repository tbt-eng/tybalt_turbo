-- Shared pricing rules for WIP and the staff and division summaries.
-- Each entry uses the rate sheet in effect on its date, from the project's
-- rate sheet history. A project with no history rows uses jobs.rate_sheet for
-- all its time: imported projects have none until their history is written.
-- Use the assigned revision even when it is inactive, so a change to a
-- sheet's rates recalculates past values. Keep the existing selection of
-- committed and uncommitted entries; amendments and meals are not part of the
-- value.
, rate_history AS (
  SELECT jrs.job, jrs.rate_sheet, jrs.effective_date AS starts
  FROM selected_jobs j
  JOIN job_rate_sheets jrs ON jrs.job = j.id
  UNION ALL
  SELECT j.id, j.rate_sheet, ''
  FROM selected_jobs j
  WHERE COALESCE(j.rate_sheet, '') != ''
    AND NOT EXISTS (SELECT 1 FROM job_rate_sheets jrs WHERE jrs.job = j.id)
), rate_periods AS (
  -- The opening period's empty start sorts first and is before every date.
  SELECT job, rate_sheet, starts,
    LEAD(starts) OVER (PARTITION BY job ORDER BY starts) AS ends
  FROM rate_history
), entry_rates AS (
  SELECT te.job, te.uid, te.division, te.date, te.hours, te.meals_hours, te.role, j.number, p.rate_sheet,
    COALESCE(rs.name, '') AS rate_sheet_name, rs.revision AS rate_sheet_revision,
    rse.rate AS sheet_rate,
    -- A zero or missing employee default cannot supply a dollar value.
    CASE WHEN ap.default_charge_out_rate > 0 THEN ap.default_charge_out_rate END AS default_rate
  -- Start with the selected jobs so a materialized report scope does not cause
  -- SQLite to scan every time entry before it applies the job filter.
  FROM selected_jobs j
  CROSS JOIN time_entries te ON te.job = j.id
  LEFT JOIN rate_periods p ON p.job = j.id
    AND te.date >= p.starts AND (p.ends IS NULL OR te.date < p.ends)
  LEFT JOIN rate_sheets rs ON rs.id = p.rate_sheet
  LEFT JOIN rate_sheet_entries rse
    ON rse.rate_sheet = p.rate_sheet AND rse.role = te.role
  LEFT JOIN admin_profiles ap ON ap.uid = te.uid
  WHERE te.date >= {:start_date}
    AND te.date <= {:end_date}
), priced_entries AS (
  SELECT *,
    hours * COALESCE(sheet_rate, default_rate) AS value,
    CASE WHEN sheet_rate IS NULL AND default_rate IS NOT NULL THEN hours ELSE 0 END AS estimated_hours,
    CASE WHEN COALESCE(sheet_rate, default_rate) IS NULL THEN hours ELSE 0 END AS unpriced_hours
  FROM entry_rates
)
