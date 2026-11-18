package routes

import (
	_ "embed" // for go:embed
	"net/http"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

//go:embed client_rate_sheets.sql
var clientRateSheetsQuery string

// clientRateSheet is a revision of a rate sheet that the client's jobs use
// now or used in the past, with how many of them.
type clientRateSheet struct {
	ID            string `db:"id" json:"id"`
	Code          string `db:"code" json:"code"`
	Name          string `db:"name" json:"name"`
	Revision      int    `db:"revision" json:"revision"`
	EffectiveDate string `db:"effective_date" json:"effective_date"`
	Active        bool   `db:"active" json:"active"`
	// CurrentJobs counts the client's jobs on the revision now. PastJobs
	// counts the ones that only a period of their rate sheet history put on
	// it.
	CurrentJobs int `db:"current_jobs" json:"current_jobs"`
	PastJobs    int `db:"past_jobs" json:"past_jobs"`
}

// createGetClientRateSheetsHandler returns the rate sheet revisions of the
// client's jobs, current and past, for the client's details tab. It is
// separate from the client details so that query stays as it is. Any
// signed-in user can read it, as they can list rate sheets, jobs and clients.
func createGetClientRateSheetsHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		id := e.Request.PathValue("id")
		if _, err := app.FindRecordById("clients", id); err != nil {
			return e.Error(http.StatusNotFound, "client not found", err)
		}
		sheets := []clientRateSheet{}
		if err := app.DB().NewQuery(clientRateSheetsQuery).Bind(dbx.Params{"id": id}).All(&sheets); err != nil {
			return e.Error(http.StatusInternalServerError, "failed to load client rate sheets", err)
		}
		return e.JSON(http.StatusOK, sheets)
	}
}
