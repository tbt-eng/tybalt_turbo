package main

import (
	"encoding/csv"
	"net/http"
	"reflect"
	"testing"
	"tybalt/internal/testutils"

	"github.com/pocketbase/pocketbase/tests"
)

// The fixture user my.time.csv@example.com has one time entry on a sheet in
// each state (committed, approved, rejected, submitted, recalled) and one
// entry with no time sheet. See time_sheets.csv and time_entries.csv.
func TestMyTimeEntriesCSV(t *testing.T) {
	ownerToken, err := testutils.GenerateRecordToken("users", "my.time.csv@example.com")
	if err != nil {
		t.Fatal(err)
	}

	const url = "/api/time_entries/mine/csv"

	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthenticated request is refused",
			Method:          http.MethodGet,
			URL:             url,
			ExpectedStatus:  http.StatusUnauthorized,
			ExpectedContent: []string{`"data":{}`},
			TestAppFactory:  testutils.SetupTestApp,
		},
		{
			Name:           "caller gets only their own time entries, newest first",
			Method:         http.MethodGet,
			URL:            url,
			Headers:        map[string]string{"Authorization": ownerToken},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`date,week_ending,time_type,time_type_name,hours,job_number,job_description,division,division_name,category,work_record,payout_request_amount,description,time_sheet_status`,
			},
			NotExpectedContent: []string{
				// time entries of author@soup.com
				`2024-07-03`,
				`f2j5a8vk006baub`,
				// ids of the caller's own rows
				`temytimecsv`,
				`tsmytimecsv`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(tb testing.TB, _ *tests.TestApp, res *http.Response) {
				defer res.Body.Close()

				if got := res.Header.Get("Content-Type"); got != "text/csv; charset=utf-8" {
					tb.Fatalf("unexpected content type %q", got)
				}

				records, err := csv.NewReader(res.Body).ReadAll()
				if err != nil {
					tb.Fatalf("failed to read csv: %v", err)
				}

				want := [][]string{
					{"2023-04-03", "2023-04-08", "R", "Hours Worked", "4", "", "", "MD", "Drilling Services", "", "", "", "unbundled entry", "not submitted"},
					{"2023-03-27", "2023-04-01", "R", "Hours Worked", "6", "", "", "MD", "Drilling Services", "", "", "", "recalled week entry", "not submitted"},
					{"2023-03-20", "2023-03-25", "OTO", "Request Overtime Payout", "0", "", "", "", "", "", "", "250.50", "", "pending"},
					{"2023-03-13", "2023-03-18", "OV", "Vacation", "8", "", "", "", "", "", "", "", "", "rejected"},
					{"2023-03-06", "2023-03-11", "R", "Hours Worked", "8", "", "", "MD", "Drilling Services", "", "", "", "approved week entry", "approved"},
					{"2023-02-27", "2023-03-04", "R", "Hours Worked", "7.5", "24-321", "Clarke St Boiler Replacement", "MD", "Drilling Services", "paramecium", "", "", "Boiler survey, east wing", "committed"},
				}
				if len(records) != len(want)+1 {
					tb.Fatalf("expected header and %d rows, got %d records: %v", len(want), len(records), records)
				}
				if !reflect.DeepEqual(records[1:], want) {
					tb.Fatalf("unexpected rows:\n got: %v\nwant: %v", records[1:], want)
				}
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
