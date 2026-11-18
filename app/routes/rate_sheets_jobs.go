package routes

import (
	_ "embed" // for go:embed
	"net/http"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

//go:embed rate_sheets_jobs.sql
var rateSheetJobsQuery string

// rateSheetJob is a job that uses or used one revision of a rate sheet.
type rateSheetJob struct {
	ID          string `db:"id" json:"id"`
	Number      string `db:"number" json:"number"`
	Description string `db:"description" json:"description"`
	ClientID    string `db:"client_id" json:"client_id"`
	ClientName  string `db:"client_name" json:"client_name"`
	// Current is true when the job uses the revision now, false when only a
	// period of its rate sheet history did.
	Current bool `db:"current" json:"current"`
}

// createGetRateSheetJobsHandler returns the jobs of this exact revision of a
// rate sheet, by client name and job number. Other revisions of the code
// have their own. Any signed-in user can read them, as they can list jobs
// and clients.
func createGetRateSheetJobsHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		id := e.Request.PathValue("id")
		if _, err := app.FindRecordById("rate_sheets", id); err != nil {
			return e.Error(http.StatusNotFound, "rate sheet not found", err)
		}
		jobs := []rateSheetJob{}
		if err := app.DB().NewQuery(rateSheetJobsQuery).Bind(dbx.Params{"id": id}).All(&jobs); err != nil {
			return e.Error(http.StatusInternalServerError, "failed to load jobs", err)
		}
		return e.JSON(http.StatusOK, jobs)
	}
}
