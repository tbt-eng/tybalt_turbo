package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// A search finds a rate sheet by a job or a client (issue 207, step 2).
// rate_sheets_augmented gets two text fields, job_numbers and client_names:
// the numbers of the jobs, and the names of their clients, that use or used
// any revision of the sheet's code, now (jobs.rate_sheet) or in a period of
// the rate sheet history (job_rate_sheets). Each is a comma separated list,
// or empty. So a job on revision 0 finds revision 1. job_count still counts
// the jobs of this revision only.
//
// The CAST makes PocketBase read each field as text.
//
// See descriptions/job_rate_sheet_history.md, Find a sheet by job or client.

const rateSheetsAugmentedFindByJobQuery = "SELECT \n  rs.id AS id,\n  rs.code AS code,\n  rs.name AS name,\n  rs.effective_date AS effective_date,\n  rs.revision AS revision,\n  rs.active AS active,\n  rs.created AS created,\n  rs.updated AS updated,\n  COUNT(u.job) AS job_count,\n  CAST(COALESCE(f.job_numbers, '') AS TEXT) AS job_numbers,\n  CAST(COALESCE(f.client_names, '') AS TEXT) AS client_names\nFROM rate_sheets rs\nLEFT JOIN (\n  SELECT j.rate_sheet AS rate_sheet, j.id AS job FROM jobs j WHERE j.rate_sheet != ''\n  UNION\n  SELECT h.rate_sheet AS rate_sheet, h.job AS job FROM job_rate_sheets h\n) u ON u.rate_sheet = rs.id\nLEFT JOIN (\n  SELECT s.code AS code,\n    GROUP_CONCAT(DISTINCT cj.number) AS job_numbers,\n    GROUP_CONCAT(DISTINCT c.name) AS client_names\n  FROM rate_sheets s\n  JOIN (\n    SELECT j.rate_sheet AS rate_sheet, j.id AS job FROM jobs j WHERE j.rate_sheet != ''\n    UNION\n    SELECT h.rate_sheet AS rate_sheet, h.job AS job FROM job_rate_sheets h\n  ) cu ON cu.rate_sheet = s.id\n  JOIN jobs cj ON cj.id = cu.job\n  LEFT JOIN clients c ON c.id = cj.client\n  GROUP BY s.code\n) f ON f.code = rs.code\nWHERE COALESCE(rs.request, '') = ''\nGROUP BY rs.id, rs.code, rs.name, rs.effective_date, rs.revision, rs.active, rs.created, rs.updated, f.job_numbers, f.client_names"

func init() {
	m.Register(func(app core.App) error {
		return setRateSheetsAugmentedQuery(app, rateSheetsAugmentedFindByJobQuery)
	}, func(app core.App) error {
		return setRateSheetsAugmentedQuery(app, rateSheetsAugmentedCodeQuery)
	})
}
