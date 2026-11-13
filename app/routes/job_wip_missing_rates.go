package routes

import (
	"net/http"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type WIPMissingRateRole struct {
	ID    string  `db:"role_id" json:"id"`
	Name  string  `db:"role_name" json:"name"`
	Hours float64 `db:"hours" json:"hours"`
}

type WIPMissingRateSheet struct {
	ID       string               `json:"id"`
	Name     string               `json:"name"`
	Revision int                  `json:"revision"`
	Roles    []WIPMissingRateRole `json:"roles"`
}

// WIPMissingRates lists, for each rate sheet that priced the job's time, the
// roles it has no rate for. A role can be priced in one period of the job's
// rate sheet history and not in the next.
type WIPMissingRates struct {
	AsOf       string                `json:"as_of"`
	RateSheets []WIPMissingRateSheet `json:"rate_sheets"`
}

func createGetJobWIPMissingRatesHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		today := time.Now().UTC().Format("2006-01-02")
		asOf := e.Request.URL.Query().Get("as_of")
		if asOf == "" {
			asOf = today
		}
		if _, err := time.Parse("2006-01-02", asOf); err != nil || asOf > today {
			return e.BadRequestError("as_of must be a valid date no later than today", nil)
		}
		jobID := e.Request.PathValue("id")
		if _, err := app.FindRecordById("jobs", jobID); err != nil {
			return e.NotFoundError("Job not found", nil)
		}
		var rows []struct {
			RateSheetID       string  `db:"rate_sheet_id"`
			RateSheetName     string  `db:"rate_sheet_name"`
			RateSheetRevision int     `db:"rate_sheet_revision"`
			RoleID            string  `db:"role_id"`
			RoleName          string  `db:"role_name"`
			Hours             float64 `db:"hours"`
		}
		// Use the same pricing rules and date limit as WIP. Each entry carries
		// the sheet of the period its date falls in.
		err := app.DB().NewQuery(singleJobScopeQuery + jobPricedTimeEntriesQuery + `
			SELECT rs.id AS rate_sheet_id, rs.name AS rate_sheet_name,
				rs.revision AS rate_sheet_revision,
				p.role AS role_id,
				COALESCE(NULLIF(r.name, ''), p.role) AS role_name,
				SUM(p.hours) AS hours
			FROM priced_entries p
			JOIN rate_sheets rs ON rs.id = p.rate_sheet
			LEFT JOIN rate_roles r ON r.id = p.role
			WHERE COALESCE(p.role, '') != '' AND p.sheet_rate IS NULL
			GROUP BY rs.id, p.role
			HAVING SUM(p.hours) != 0
			ORDER BY rs.name, rs.revision, role_name, role_id
		`).Bind(dbx.Params{
			"job_id": jobID, "start_date": "", "end_date": asOf,
		}).All(&rows)
		if err != nil {
			return e.InternalServerError("Failed to load missing rates", err)
		}
		result := WIPMissingRates{AsOf: asOf, RateSheets: []WIPMissingRateSheet{}}
		for _, row := range rows {
			n := len(result.RateSheets)
			if n == 0 || result.RateSheets[n-1].ID != row.RateSheetID {
				result.RateSheets = append(result.RateSheets, WIPMissingRateSheet{
					ID: row.RateSheetID, Name: row.RateSheetName, Revision: row.RateSheetRevision,
					Roles: []WIPMissingRateRole{},
				})
				n++
			}
			result.RateSheets[n-1].Roles = append(result.RateSheets[n-1].Roles,
				WIPMissingRateRole{ID: row.RoleID, Name: row.RoleName, Hours: row.Hours})
		}
		return e.JSON(http.StatusOK, result)
	}
}
