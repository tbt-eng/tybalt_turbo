package main

import (
	"encoding/csv"
	"encoding/json"
	"math"
	"net/http"
	"slices"
	"testing"
	"tybalt/internal/testutils"

	"github.com/pocketbase/pocketbase/tests"
)

const (
	branchHoursThunderBayID = "80875lm27v8wgi4"
	branchHoursTorontoID    = "xeq9q81q5307f70"
	branchHoursMay2040      = "&start_date=2040-05-01&end_date=2040-05-31"
)

type branchHoursAPIRow struct {
	Direction      string  `json:"direction"`
	JobID          string  `json:"jobId"`
	JobNumber      string  `json:"jobNumber"`
	JobDescription string  `json:"jobDescription"`
	Client         string  `json:"client"`
	WorkBranch     string  `json:"workBranch"`
	UID            string  `json:"uid"`
	PayrollID      string  `json:"payrollId"`
	Surname        string  `json:"surname"`
	GivenName      string  `json:"givenName"`
	StaffBranch    string  `json:"staffBranch"`
	Hours          float64 `json:"hours"`
}

type branchHoursAPIResponse struct {
	Branch string              `json:"branch"`
	Rows   []branchHoursAPIRow `json:"rows"`
}

func TestBranchHoursReportAccess(t *testing.T) {
	url := "/api/branch_hours?branch=" + branchHoursThunderBayID + branchHoursMay2040
	cases := []struct {
		name   string
		email  string
		status int
	}{
		{"ordinary user cannot view branch hours", "valuation1@example.com", http.StatusForbidden},
		{"report claim does not grant branch hours", "fatt@mac.com", http.StatusForbidden},
		{"branch manager can view another branch", "u_no_claims@example.com", http.StatusOK},
		{"kpi can view branch hours", "time@test.com", http.StatusOK},
		{"admin can view branch hours", "admin.only@example.com", http.StatusOK},
	}
	(&tests.ApiScenario{
		Name:            "branch hours requires authentication",
		Method:          http.MethodGet,
		URL:             url,
		ExpectedStatus:  http.StatusUnauthorized,
		ExpectedContent: []string{"The request requires valid record authorization token."},
		TestAppFactory:  testutils.SetupTestApp,
	}).Test(t)
	for _, tc := range cases {
		expected := []string{`"rows"`}
		if tc.status == http.StatusForbidden {
			expected = []string{"You need to be a branch manager or hold the kpi or admin claim to view this report."}
		}
		(&tests.ApiScenario{
			Name:            tc.name,
			Method:          http.MethodGet,
			URL:             url,
			Headers:         map[string]string{"Authorization": mustGenerateEmployeeBranchHoursToken(t, tc.email)},
			ExpectedStatus:  tc.status,
			ExpectedContent: expected,
			TestAppFactory:  testutils.SetupTestApp,
		}).Test(t)
	}
}

func TestBranchHoursReportData(t *testing.T) {
	token := mustGenerateEmployeeBranchHoursToken(t, "time@test.com")
	cases := []struct {
		name   string
		branch string
		want   []branchHoursAPIRow
	}{
		{
			name:   "staff branch sees internal and supplied time by job and employee",
			branch: branchHoursThunderBayID,
			want: []branchHoursAPIRow{
				{Direction: "internal", JobNumber: "98-8101", WorkBranch: "Thunder Bay", PayrollID: "900001", StaffBranch: "Thunder Bay", Hours: 4},
				{Direction: "internal", JobNumber: "", WorkBranch: "Thunder Bay", PayrollID: "900001", StaffBranch: "Thunder Bay", Hours: 1.5},
				{Direction: "internal", JobNumber: "", WorkBranch: "Thunder Bay", PayrollID: "FIX_INACTIVE", StaffBranch: "Thunder Bay", Hours: 3},
				{Direction: "supplied", JobNumber: "98-8101", WorkBranch: "Toronto", PayrollID: "900001", StaffBranch: "Thunder Bay", Hours: 2},
				{Direction: "supplied", JobNumber: "98-8101", WorkBranch: "Toronto", PayrollID: "FIX_INACTIVE", StaffBranch: "Thunder Bay", Hours: 2.5},
			},
		},
		{
			name:   "work branch sees time received from other branches",
			branch: branchHoursTorontoID,
			want: []branchHoursAPIRow{
				{Direction: "received", JobNumber: "98-8101", WorkBranch: "Toronto", PayrollID: "900001", StaffBranch: "Thunder Bay", Hours: 2},
				{Direction: "received", JobNumber: "98-8101", WorkBranch: "Toronto", PayrollID: "FIX_INACTIVE", StaffBranch: "Thunder Bay", Hours: 2.5},
			},
		},
	}
	for _, tc := range cases {
		(&tests.ApiScenario{
			Name:            tc.name,
			Method:          http.MethodGet,
			URL:             "/api/branch_hours?branch=" + tc.branch + branchHoursMay2040,
			Headers:         map[string]string{"Authorization": token},
			ExpectedStatus:  http.StatusOK,
			ExpectedContent: []string{`"rows"`},
			TestAppFactory:  testutils.SetupTestApp,
			AfterTestFunc: func(tb testing.TB, _ *tests.TestApp, response *http.Response) {
				payload := decodeBranchHoursResponse(tb, response)
				if len(payload.Rows) != len(tc.want) {
					tb.Fatalf("row count = %d, want %d: %+v", len(payload.Rows), len(tc.want), payload.Rows)
				}
				for index, want := range tc.want {
					got := payload.Rows[index]
					if got.Direction != want.Direction || got.JobNumber != want.JobNumber ||
						got.WorkBranch != want.WorkBranch || got.PayrollID != want.PayrollID ||
						got.StaffBranch != want.StaffBranch || got.Hours != want.Hours {
						tb.Fatalf("row %d = %+v, want %+v", index, got, want)
					}
					if got.JobNumber != "" && (got.JobID != "jobbrmatch0001" ||
						got.JobDescription != "Branch mismatch fixture job" || got.Client != "Sepulchi & Co") {
						tb.Fatalf("row %d job details = %+v", index, got)
					}
				}
			},
		}).Test(t)
	}
}

// Every branch's Branch Hours totals must agree with the KPI Employee Branch
// Hours report: time worked for the branch equals its branch and NoJob columns,
// and time by the branch's staff equals their totals.
func TestBranchHoursReportReconcilesWithEmployeeBranchHours(t *testing.T) {
	token := mustGenerateEmployeeBranchHoursToken(t, "time@test.com")
	dates := "start_date=2000-01-01&end_date=2099-12-31"

	var kpi employeeBranchHoursAPIResponse
	(&tests.ApiScenario{
		Name:            "load KPI employee branch hours",
		Method:          http.MethodGet,
		URL:             "/api/kpi/reports/employee_branch_hours?" + dates,
		Headers:         map[string]string{"Authorization": token},
		ExpectedStatus:  http.StatusOK,
		ExpectedContent: []string{`"rows"`},
		TestAppFactory:  testutils.SetupTestApp,
		AfterTestFunc: func(tb testing.TB, _ *tests.TestApp, response *http.Response) {
			defer response.Body.Close()
			if err := json.NewDecoder(response.Body).Decode(&kpi); err != nil {
				tb.Fatalf("failed to decode KPI report: %v", err)
			}
		},
	}).Test(t)

	app := testutils.SetupTestApp(t)
	branches, err := app.FindAllRecords("branches")
	app.Cleanup()
	if err != nil {
		t.Fatalf("failed to load branches: %v", err)
	}

	checked := 0
	for _, branch := range branches {
		name := branch.GetString("name")
		var workedFor, staffTotal float64
		for _, row := range kpi.Rows {
			workedFor += row[name].(float64) + row[name+"NoJob"].(float64)
			if row["defaultBranch"] == name {
				staffTotal += row["total"].(float64)
			}
		}
		(&tests.ApiScenario{
			Name:            "branch hours reconcile for " + name,
			Method:          http.MethodGet,
			URL:             "/api/branch_hours?branch=" + branch.Id + "&" + dates,
			Headers:         map[string]string{"Authorization": token},
			ExpectedStatus:  http.StatusOK,
			ExpectedContent: []string{`"rows"`},
			TestAppFactory:  testutils.SetupTestApp,
			AfterTestFunc: func(tb testing.TB, _ *tests.TestApp, response *http.Response) {
				totals := map[string]float64{}
				for _, row := range decodeBranchHoursResponse(tb, response).Rows {
					totals[row.Direction] += row.Hours
				}
				if got := totals["internal"] + totals["received"]; math.Abs(got-workedFor) > 1e-9 {
					tb.Fatalf("%s worked-for hours = %v, KPI = %v", name, got, workedFor)
				}
				if got := totals["internal"] + totals["supplied"]; math.Abs(got-staffTotal) > 1e-9 {
					tb.Fatalf("%s staff hours = %v, KPI = %v", name, got, staffTotal)
				}
			},
		}).Test(t)
		if workedFor > 0 || staffTotal > 0 {
			checked++
		}
	}
	if checked < 2 {
		t.Fatalf("reconciliation covered %d branches with hours, want at least 2", checked)
	}
}

func TestBranchHoursReportCSV(t *testing.T) {
	(&tests.ApiScenario{
		Name:            "branch hours CSV has one row per job and employee",
		Method:          http.MethodGet,
		URL:             "/api/branch_hours?branch=" + branchHoursThunderBayID + branchHoursMay2040 + "&format=csv",
		Headers:         map[string]string{"Authorization": mustGenerateEmployeeBranchHoursToken(t, "time@test.com")},
		ExpectedStatus:  http.StatusOK,
		ExpectedContent: []string{"direction,jobNumber"},
		TestAppFactory:  testutils.SetupTestApp,
		AfterTestFunc: func(tb testing.TB, _ *tests.TestApp, response *http.Response) {
			defer response.Body.Close()
			if got := response.Header.Get("Content-Disposition"); got != `attachment; filename="branch_hours_Thunder_Bay_2040-05-01_2040-05-31.csv"` {
				tb.Fatalf("Content-Disposition = %q", got)
			}
			records, err := csv.NewReader(response.Body).ReadAll()
			if err != nil {
				tb.Fatalf("failed to parse branch hours CSV: %v", err)
			}
			wantColumns := []string{"direction", "jobNumber", "jobDescription", "client", "workBranch", "payrollId", "surname", "givenName", "staffBranch", "hours"}
			if !slices.Equal(records[0], wantColumns) {
				tb.Fatalf("CSV columns = %v, want %v", records[0], wantColumns)
			}
			if len(records) != 6 {
				tb.Fatalf("CSV row count = %d, want 6", len(records))
			}
			want := []string{"supplied", "98-8101", "Branch mismatch fixture job", "Sepulchi & Co", "Toronto", "FIX_INACTIVE", "User", "Inactive", "Thunder Bay", "2.5"}
			if !slices.Equal(records[5], want) {
				tb.Fatalf("last CSV row = %v, want %v", records[5], want)
			}
		},
	}).Test(t)
}

func TestBranchHoursReportValidation(t *testing.T) {
	token := mustGenerateEmployeeBranchHoursToken(t, "time@test.com")
	cases := []struct {
		name  string
		query string
		want  string
	}{
		{"branch is required", "start_date=2040-05-01&end_date=2040-05-31", "Branch is required"},
		{"unknown branch is rejected", "branch=missingbranch01" + branchHoursMay2040, "Unknown branch"},
		{"invalid date is rejected", "branch=" + branchHoursThunderBayID + "&start_date=2040-02-30&end_date=2040-05-31", "Start_date must be a valid date in YYYY-MM-DD format"},
		{"reversed range is rejected", "branch=" + branchHoursThunderBayID + "&start_date=2040-06-01&end_date=2040-05-31", "Start_date must be on or before end_date"},
		{"unknown format is rejected", "branch=" + branchHoursThunderBayID + branchHoursMay2040 + "&format=pdf", "Format must be json or csv"},
	}
	for _, tc := range cases {
		(&tests.ApiScenario{
			Name:            tc.name,
			Method:          http.MethodGet,
			URL:             "/api/branch_hours?" + tc.query,
			Headers:         map[string]string{"Authorization": token},
			ExpectedStatus:  http.StatusBadRequest,
			ExpectedContent: []string{tc.want},
			TestAppFactory:  testutils.SetupTestApp,
		}).Test(t)
	}
	(&tests.ApiScenario{
		Name:            "empty range returns no rows",
		Method:          http.MethodGet,
		URL:             "/api/branch_hours?branch=" + branchHoursThunderBayID + "&start_date=1900-01-01&end_date=1900-01-31",
		Headers:         map[string]string{"Authorization": token},
		ExpectedStatus:  http.StatusOK,
		ExpectedContent: []string{`"branch":"Thunder Bay"`, `"rows":[]`},
		TestAppFactory:  testutils.SetupTestApp,
	}).Test(t)
}

func decodeBranchHoursResponse(tb testing.TB, response *http.Response) branchHoursAPIResponse {
	tb.Helper()
	defer response.Body.Close()
	var payload branchHoursAPIResponse
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		tb.Fatalf("failed to decode branch hours response: %v", err)
	}
	return payload
}
