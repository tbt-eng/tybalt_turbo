package reports

import (
	_ "embed"
	"encoding/csv"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

//go:embed branch_hours.sql
var branchHoursQuery string

type BranchHoursRow struct {
	Direction      string  `db:"direction" json:"direction"`
	JobID          string  `db:"jobId" json:"jobId"`
	JobNumber      string  `db:"jobNumber" json:"jobNumber"`
	JobDescription string  `db:"jobDescription" json:"jobDescription"`
	Client         string  `db:"client" json:"client"`
	WorkBranch     string  `db:"workBranch" json:"workBranch"`
	UID            string  `db:"uid" json:"uid"`
	Surname        string  `db:"surname" json:"surname"`
	GivenName      string  `db:"givenName" json:"givenName"`
	StaffBranch    string  `db:"staffBranch" json:"staffBranch"`
	Hours          float64 `db:"hours" json:"hours"`
}

type BranchHoursReport struct {
	Branch string           `json:"branch"`
	Rows   []BranchHoursRow `json:"rows"`
}

var branchHoursCSVColumns = []string{
	"direction",
	"jobNumber",
	"jobDescription",
	"client",
	"workBranch",
	"surname",
	"givenName",
	"staffBranch",
	"hours",
}

var branchHoursFilenameUnsafe = regexp.MustCompile(`[^A-Za-z0-9]+`)

// CreateBranchHoursReportHandler returns hours by job and employee for the
// time that involves one branch, either because the employee belongs to the
// branch or because the time was worked for it. The route is responsible for
// requiring branch report access.
func CreateBranchHoursReportHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		queryValues := e.Request.URL.Query()
		branchID := strings.TrimSpace(queryValues.Get("branch"))
		if branchID == "" {
			return e.Error(http.StatusBadRequest, "branch is required", nil)
		}
		startDate, err := parseEmployeeBranchHoursDate(e, queryValues.Get("start_date"), "start_date")
		if err != nil {
			return err
		}
		endDate, err := parseEmployeeBranchHoursDate(e, queryValues.Get("end_date"), "end_date")
		if err != nil {
			return err
		}
		if startDate.After(endDate) {
			return e.Error(http.StatusBadRequest, "start_date must be on or before end_date", nil)
		}

		format := strings.ToLower(strings.TrimSpace(queryValues.Get("format")))
		if format == "" {
			format = "json"
		}
		if format != "json" && format != "csv" {
			return e.Error(http.StatusBadRequest, "format must be json or csv", nil)
		}

		branch, err := app.FindRecordById("branches", branchID)
		if err != nil {
			return e.Error(http.StatusBadRequest, "unknown branch", nil)
		}

		report := BranchHoursReport{Branch: branch.GetString("name"), Rows: []BranchHoursRow{}}
		if err := app.DB().NewQuery(branchHoursQuery).Bind(dbx.Params{
			"branch":     branch.Id,
			"start_date": startDate.Format(time.DateOnly),
			"end_date":   endDate.Format(time.DateOnly),
		}).All(&report.Rows); err != nil {
			return e.Error(http.StatusInternalServerError, "failed to create branch hours report", err)
		}

		if format == "json" {
			return e.JSON(http.StatusOK, report)
		}

		csvText, err := branchHoursCSV(report.Rows)
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to create branch hours CSV", err)
		}
		filename := fmt.Sprintf(
			"branch_hours_%s_%s_%s.csv",
			strings.Trim(branchHoursFilenameUnsafe.ReplaceAllString(report.Branch, "_"), "_"),
			startDate.Format(time.DateOnly),
			endDate.Format(time.DateOnly),
		)
		e.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
		e.Response.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
		return e.String(http.StatusOK, csvText)
	}
}

func branchHoursCSV(rows []BranchHoursRow) (string, error) {
	var builder strings.Builder
	writer := csv.NewWriter(&builder)
	if err := writer.Write(branchHoursCSVColumns); err != nil {
		return "", err
	}
	for _, row := range rows {
		if err := writer.Write([]string{
			row.Direction,
			row.JobNumber,
			row.JobDescription,
			row.Client,
			row.WorkBranch,
			row.Surname,
			row.GivenName,
			row.StaffBranch,
			strconv.FormatFloat(row.Hours, 'f', -1, 64),
		}); err != nil {
			return "", err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", err
	}
	return builder.String(), nil
}
