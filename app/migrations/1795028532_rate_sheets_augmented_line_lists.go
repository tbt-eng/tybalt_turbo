package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// rate_sheets_augmented joins job_numbers and client_names with a line break,
// not a comma (issue 207). A client name can contain a comma, so the UI could
// not split the comma separated list to show which entry a search matched.
//
// SQLite's GROUP_CONCAT takes no separator with DISTINCT, so the UNION in d
// makes the values distinct: one row for each code, kind and value. Then
// GROUP_CONCAT joins each kind with char(10). Everything else is as in
// 1795025548_rate_sheets_augmented_find_by_job.go.
//
// See descriptions/job_rate_sheet_history.md, Find a sheet by job or client.

const rateSheetsAugmentedLineListsQuery = "SELECT \n  rs.id AS id,\n  rs.code AS code,\n  rs.name AS name,\n  rs.effective_date AS effective_date,\n  rs.revision AS revision,\n  rs.active AS active,\n  rs.created AS created,\n  rs.updated AS updated,\n  COUNT(u.job) AS job_count,\n  CAST(COALESCE(f.job_numbers, '') AS TEXT) AS job_numbers,\n  CAST(COALESCE(f.client_names, '') AS TEXT) AS client_names\nFROM rate_sheets rs\nLEFT JOIN (\n  SELECT j.rate_sheet AS rate_sheet, j.id AS job FROM jobs j WHERE j.rate_sheet != ''\n  UNION\n  SELECT h.rate_sheet AS rate_sheet, h.job AS job FROM job_rate_sheets h\n) u ON u.rate_sheet = rs.id\nLEFT JOIN (\n  SELECT d.code AS code,\n    GROUP_CONCAT(CASE WHEN d.kind = 'job' THEN d.value END, char(10)) AS job_numbers,\n    GROUP_CONCAT(CASE WHEN d.kind = 'client' THEN d.value END, char(10)) AS client_names\n  FROM (\n    SELECT s.code AS code, 'job' AS kind, cj.number AS value\n    FROM rate_sheets s\n    JOIN (\n      SELECT j.rate_sheet AS rate_sheet, j.id AS job FROM jobs j WHERE j.rate_sheet != ''\n      UNION\n      SELECT h.rate_sheet AS rate_sheet, h.job AS job FROM job_rate_sheets h\n    ) cu ON cu.rate_sheet = s.id\n    JOIN jobs cj ON cj.id = cu.job\n    UNION\n    SELECT s.code AS code, 'client' AS kind, c.name AS value\n    FROM rate_sheets s\n    JOIN (\n      SELECT j.rate_sheet AS rate_sheet, j.id AS job FROM jobs j WHERE j.rate_sheet != ''\n      UNION\n      SELECT h.rate_sheet AS rate_sheet, h.job AS job FROM job_rate_sheets h\n    ) cu ON cu.rate_sheet = s.id\n    JOIN jobs cj ON cj.id = cu.job\n    JOIN clients c ON c.id = cj.client\n  ) d\n  GROUP BY d.code\n) f ON f.code = rs.code\nWHERE COALESCE(rs.request, '') = ''\nGROUP BY rs.id, rs.code, rs.name, rs.effective_date, rs.revision, rs.active, rs.created, rs.updated, f.job_numbers, f.client_names"

func init() {
	m.Register(func(app core.App) error {
		return setRateSheetsAugmentedQuery(app, rateSheetsAugmentedLineListsQuery)
	}, func(app core.App) error {
		return setRateSheetsAugmentedQuery(app, rateSheetsAugmentedFindByJobQuery)
	})
}
