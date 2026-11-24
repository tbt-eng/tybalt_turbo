package routes

import (
	_ "embed" // for go:embed
	"encoding/csv"
	"net/http"
	"strconv"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

//go:embed my_time_entries.sql
var myTimeEntriesQuery string

type myTimeEntryRow struct {
	Date                string  `db:"date"`
	WeekEnding          string  `db:"week_ending"`
	TimeType            string  `db:"time_type"`
	TimeTypeName        string  `db:"time_type_name"`
	Hours               float64 `db:"hours"`
	MealsHours          float64 `db:"meals_hours"`
	JobNumber           string  `db:"job_number"`
	JobDescription      string  `db:"job_description"`
	Division            string  `db:"division"`
	DivisionName        string  `db:"division_name"`
	Category            string  `db:"category"`
	Branch              string  `db:"branch"`
	Role                string  `db:"role"`
	WorkRecord          string  `db:"work_record"`
	PayoutRequestAmount float64 `db:"payout_request_amount"`
	Description         string  `db:"description"`
	TimeSheetStatus     string  `db:"time_sheet_status"`
}

var myTimeEntriesHeaders = []string{
	"date",
	"week_ending",
	"time_type",
	"time_type_name",
	"hours",
	"meals_hours",
	"job_number",
	"job_description",
	"division",
	"division_name",
	"category",
	"branch",
	"role",
	"work_record",
	"payout_request_amount",
	"description",
	"time_sheet_status",
}

// createMyTimeEntriesCSVHandler returns every time entry of the caller as a
// CSV file, newest first, whatever the state of its time sheet. It reads only
// rows whose uid is the caller, so it needs no claim.
func createMyTimeEntriesCSVHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		var rows []myTimeEntryRow
		if err := app.DB().NewQuery(myTimeEntriesQuery).Bind(dbx.Params{
			"uid": e.Auth.Id,
		}).All(&rows); err != nil {
			return e.Error(http.StatusInternalServerError, "failed to execute query", err)
		}

		var csvBuilder strings.Builder
		writer := csv.NewWriter(&csvBuilder)
		if err := writer.Write(myTimeEntriesHeaders); err != nil {
			return e.Error(http.StatusInternalServerError, "failed to write csv header", err)
		}
		for _, row := range rows {
			payout := ""
			if row.PayoutRequestAmount != 0 {
				payout = strconv.FormatFloat(row.PayoutRequestAmount, 'f', 2, 64)
			}
			if err := writer.Write([]string{
				row.Date,
				row.WeekEnding,
				row.TimeType,
				row.TimeTypeName,
				strconv.FormatFloat(row.Hours, 'f', -1, 64),
				strconv.FormatFloat(row.MealsHours, 'f', -1, 64),
				row.JobNumber,
				row.JobDescription,
				row.Division,
				row.DivisionName,
				row.Category,
				row.Branch,
				row.Role,
				row.WorkRecord,
				payout,
				row.Description,
				row.TimeSheetStatus,
			}); err != nil {
				return e.Error(http.StatusInternalServerError, "failed to write csv row", err)
			}
		}
		writer.Flush()
		if err := writer.Error(); err != nil {
			return e.Error(http.StatusInternalServerError, "failed to finalize csv", err)
		}

		e.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
		e.Response.Header().Set("Content-Disposition", `attachment; filename="my_time_entries.csv"`)
		return e.String(http.StatusOK, csvBuilder.String())
	}
}
