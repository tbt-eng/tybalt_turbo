package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// rate_sheets_augmented feeds the Rate Sheets list. Its job_count counted the
// jobs whose jobs.rate_sheet is the sheet. Time is priced from the rate sheet
// history (job_rate_sheets), so a job also uses a sheet when one of its
// periods is on it. The view now counts each such job once, as the rate sheet
// absorb rules do (app/hooks/rate_sheet_absorb.go).

const rateSheetsAugmentedCurrentOnlyQuery = "SELECT \n  rs.id AS id,\n  rs.name AS name,\n  rs.effective_date AS effective_date,\n  rs.revision AS revision,\n  rs.active AS active,\n  rs.created AS created,\n  rs.updated AS updated,\n  COUNT(j.id) AS job_count\nFROM rate_sheets rs\nLEFT JOIN jobs j ON j.rate_sheet = rs.id\nGROUP BY rs.id, rs.name, rs.effective_date, rs.revision, rs.active, rs.created, rs.updated"

// The UNION keeps one row for each sheet and job, so COUNT counts a job once
// even when it uses the sheet both now and in an earlier period.
const rateSheetsAugmentedHistoryQuery = "SELECT \n  rs.id AS id,\n  rs.name AS name,\n  rs.effective_date AS effective_date,\n  rs.revision AS revision,\n  rs.active AS active,\n  rs.created AS created,\n  rs.updated AS updated,\n  COUNT(u.job) AS job_count\nFROM rate_sheets rs\nLEFT JOIN (\n  SELECT j.rate_sheet AS rate_sheet, j.id AS job FROM jobs j WHERE j.rate_sheet != ''\n  UNION\n  SELECT h.rate_sheet AS rate_sheet, h.job AS job FROM job_rate_sheets h\n) u ON u.rate_sheet = rs.id\nGROUP BY rs.id, rs.name, rs.effective_date, rs.revision, rs.active, rs.created, rs.updated"

func setRateSheetsAugmentedQuery(app core.App, query string) error {
	collection, err := app.FindCollectionByNameOrId("rate_sheets_augmented")
	if err != nil {
		return err
	}
	collection.ViewQuery = query
	return app.Save(collection)
}

func init() {
	m.Register(func(app core.App) error {
		return setRateSheetsAugmentedQuery(app, rateSheetsAugmentedHistoryQuery)
	}, func(app core.App) error {
		return setRateSheetsAugmentedQuery(app, rateSheetsAugmentedCurrentOnlyQuery)
	})
}
