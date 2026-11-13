package main

import (
	"net/http"
	"strings"
	"testing"

	"tybalt/internal/testutils"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

// Fixtures:
//   - test_job_w_rs: project on c41ofep525bcacj (active, effective 2025-01-07),
//     with an opening row from the migration
//   - u09fwwcg07y03m7: project with no rate sheet
//   - test_rs_alt_actv: active sheet, effective 2025-01-01
//   - author@soup.com (f2j5a8vk006baub): active, job claim, no rate_sheet_revise
//   - m3kzmuowqlzuic0: the rate_sheet_revise claim

const rateSheetReviseClaim = "m3kzmuowqlzuic0"

func grantRateSheetRevise(tb testing.TB, app *tests.TestApp) {
	tb.Helper()
	_, err := app.DB().NewQuery(`INSERT INTO user_claims (id, uid, cid, created, updated, _imported)
		VALUES ('ucrevisetest001', 'f2j5a8vk006baub', {:cid}, '2026-10-19 00:00:00.000Z', '2026-10-19 00:00:00.000Z', 0)`).
		Bind(dbx.Params{"cid": rateSheetReviseClaim}).Execute()
	if err != nil {
		tb.Fatal(err)
	}
}

func insertRateSheetRow(tb testing.TB, app *tests.TestApp, id, job, sheet, date string) {
	tb.Helper()
	_, err := app.DB().NewQuery(`INSERT INTO job_rate_sheets (id, job, rate_sheet, effective_date, reason, created_by, updated_by, created, updated)
		VALUES ({:id}, {:job}, {:sheet}, {:date}, 'fixture', '', '', '2026-10-19 00:00:00.000Z', '2026-10-19 00:00:00.000Z')`).
		Bind(dbx.Params{"id": id, "job": job, "sheet": sheet, "date": date}).Execute()
	if err != nil {
		tb.Fatal(err)
	}
}

func jobRateSheet(tb testing.TB, app *tests.TestApp, job string) string {
	tb.Helper()
	record, err := app.FindRecordById("jobs", job)
	if err != nil {
		tb.Fatal(err)
	}
	return record.GetString("rate_sheet")
}

func rateSheetRowCount(tb testing.TB, app *tests.TestApp, job string) int64 {
	tb.Helper()
	n, err := app.CountRecords("job_rate_sheets", dbx.HashExp{"job": job})
	if err != nil {
		tb.Fatal(err)
	}
	return n
}

func TestJobRateSheetHistory(t *testing.T) {
	token, err := testutils.GenerateRecordToken("users", "author@soup.com")
	if err != nil {
		t.Fatal(err)
	}
	headers := map[string]string{"Authorization": token}
	future := "2999-01-01"

	scenarios := []tests.ApiScenario{
		{
			Name:            "history lists the opening period and the user's permissions",
			Method:          http.MethodGet,
			URL:             "/api/jobs/test_job_w_rs/rate_sheets",
			Headers:         headers,
			ExpectedStatus:  http.StatusOK,
			ExpectedContent: []string{`"can_revise":false`, `"can_set_first":false`, `"starts":""`, `"id":"c41ofep525bcacj"`},
			TestAppFactory:  testutils.SetupTestApp,
		},
		{
			Name:            "a dated period needs rate_sheet_revise",
			Method:          http.MethodPost,
			URL:             "/api/jobs/test_job_w_rs/rate_sheets",
			Body:            strings.NewReader(`{"rate_sheet":"test_rs_alt_actv","effective_date":"2026-01-02","reason":"Client accepted new rates"}`),
			Headers:         headers,
			ExpectedStatus:  http.StatusForbidden,
			ExpectedContent: []string{"rate_sheet_revise"},
			TestAppFactory:  testutils.SetupTestApp,
		},
		{
			Name:            "a reviser adds a dated period and the job follows it",
			Method:          http.MethodPost,
			URL:             "/api/jobs/test_job_w_rs/rate_sheets",
			Body:            strings.NewReader(`{"rate_sheet":"test_rs_alt_actv","effective_date":"2026-01-02","reason":"Client accepted new rates"}`),
			Headers:         headers,
			BeforeTestFunc:  func(tb testing.TB, app *tests.TestApp, _ *core.ServeEvent) { grantRateSheetRevise(tb, app) },
			ExpectedStatus:  http.StatusOK,
			ExpectedContent: []string{`"starts":"2026-01-02"`, `"through":"2026-01-01"`},
			TestAppFactory:  testutils.SetupTestApp,
			AfterTestFunc: func(tb testing.TB, app *tests.TestApp, _ *http.Response) {
				if got := jobRateSheet(tb, app, "test_job_w_rs"); got != "test_rs_alt_actv" {
					tb.Errorf("jobs.rate_sheet = %q; want test_rs_alt_actv", got)
				}
				if n := rateSheetRowCount(tb, app, "test_job_w_rs"); n != 2 {
					tb.Errorf("rows = %d; want 2", n)
				}
			},
		},
		{
			Name:            "a period on the same sheet as the one before is refused",
			Method:          http.MethodPost,
			URL:             "/api/jobs/test_job_w_rs/rate_sheets",
			Body:            strings.NewReader(`{"rate_sheet":"c41ofep525bcacj","effective_date":"2026-01-02","reason":"No change"}`),
			Headers:         headers,
			BeforeTestFunc:  func(tb testing.TB, app *tests.TestApp, _ *core.ServeEvent) { grantRateSheetRevise(tb, app) },
			ExpectedStatus:  http.StatusBadRequest,
			ExpectedContent: []string{"same_as_previous"},
			TestAppFactory:  testutils.SetupTestApp,
		},
		{
			Name:            "a future start date is refused",
			Method:          http.MethodPost,
			URL:             "/api/jobs/test_job_w_rs/rate_sheets",
			Body:            strings.NewReader(`{"rate_sheet":"test_rs_alt_actv","effective_date":"` + future + `","reason":"Later rates"}`),
			Headers:         headers,
			BeforeTestFunc:  func(tb testing.TB, app *tests.TestApp, _ *core.ServeEvent) { grantRateSheetRevise(tb, app) },
			ExpectedStatus:  http.StatusBadRequest,
			ExpectedContent: []string{`"future"`},
			TestAppFactory:  testutils.SetupTestApp,
		},
		{
			Name:            "a start date before the revision took effect is refused",
			Method:          http.MethodPost,
			URL:             "/api/jobs/test_job_w_rs/rate_sheets",
			Body:            strings.NewReader(`{"rate_sheet":"test_rs_alt_actv","effective_date":"2024-12-31","reason":"Backdated"}`),
			Headers:         headers,
			BeforeTestFunc:  func(tb testing.TB, app *tests.TestApp, _ *core.ServeEvent) { grantRateSheetRevise(tb, app) },
			ExpectedStatus:  http.StatusBadRequest,
			ExpectedContent: []string{"before_revision"},
			TestAppFactory:  testutils.SetupTestApp,
		},
		{
			Name:            "a dated period needs a reason",
			Method:          http.MethodPost,
			URL:             "/api/jobs/test_job_w_rs/rate_sheets",
			Body:            strings.NewReader(`{"rate_sheet":"test_rs_alt_actv","effective_date":"2026-01-02","reason":"  "}`),
			Headers:         headers,
			BeforeTestFunc:  func(tb testing.TB, app *tests.TestApp, _ *core.ServeEvent) { grantRateSheetRevise(tb, app) },
			ExpectedStatus:  http.StatusBadRequest,
			ExpectedContent: []string{`"reason"`},
			TestAppFactory:  testutils.SetupTestApp,
		},
		{
			Name:            "the opening period cannot be deleted",
			Method:          http.MethodDelete,
			URL:             "/api/jobs/test_job_w_rs/rate_sheets/opening",
			Headers:         headers,
			BeforeTestFunc:  func(tb testing.TB, app *tests.TestApp, _ *core.ServeEvent) { grantRateSheetRevise(tb, app) },
			ExpectedStatus:  http.StatusBadRequest,
			ExpectedContent: []string{"opening period cannot be deleted"},
			TestAppFactory:  testutils.SetupTestApp,
		},
		{
			Name:            "changing the opening period's sheet changes the job's sheet",
			Method:          http.MethodPatch,
			URL:             "/api/jobs/test_job_w_rs/rate_sheets/opening",
			Body:            strings.NewReader(`{"rate_sheet":"test_rs_alt_actv","reason":"Wrong sheet chosen at setup"}`),
			Headers:         headers,
			BeforeTestFunc:  func(tb testing.TB, app *tests.TestApp, _ *core.ServeEvent) { grantRateSheetRevise(tb, app) },
			ExpectedStatus:  http.StatusOK,
			ExpectedContent: []string{`"id":"test_rs_alt_actv"`},
			TestAppFactory:  testutils.SetupTestApp,
			AfterTestFunc: func(tb testing.TB, app *tests.TestApp, _ *http.Response) {
				if got := jobRateSheet(tb, app, "test_job_w_rs"); got != "test_rs_alt_actv" {
					tb.Errorf("jobs.rate_sheet = %q; want test_rs_alt_actv", got)
				}
			},
		},
		{
			Name:    "deleting a period that would leave neighbours on one sheet is refused",
			Method:  http.MethodDelete,
			URL:     "/api/jobs/test_job_w_rs/rate_sheets/rsperiodmid0001",
			Headers: headers,
			BeforeTestFunc: func(tb testing.TB, app *tests.TestApp, _ *core.ServeEvent) {
				grantRateSheetRevise(tb, app)
				insertRateSheetRow(tb, app, "rsperiodmid0001", "test_job_w_rs", "test_rs_alt_actv", "2026-01-01")
				insertRateSheetRow(tb, app, "rsperiodlate001", "test_job_w_rs", "c41ofep525bcacj", "2026-02-01")
			},
			ExpectedStatus:  http.StatusBadRequest,
			ExpectedContent: []string{"same_as_previous"},
			TestAppFactory:  testutils.SetupTestApp,
			AfterTestFunc: func(tb testing.TB, app *tests.TestApp, _ *http.Response) {
				if n := rateSheetRowCount(tb, app, "test_job_w_rs"); n != 3 {
					tb.Errorf("rows = %d; want 3 after the refused delete", n)
				}
			},
		},
		{
			Name:            "a job holder sets a first rate sheet without rate_sheet_revise",
			Method:          http.MethodPost,
			URL:             "/api/jobs/u09fwwcg07y03m7/rate_sheets",
			Body:            strings.NewReader(`{"rate_sheet":"c41ofep525bcacj"}`),
			Headers:         headers,
			ExpectedStatus:  http.StatusOK,
			ExpectedContent: []string{`"starts":""`, `"id":"c41ofep525bcacj"`},
			TestAppFactory:  testutils.SetupTestApp,
			AfterTestFunc: func(tb testing.TB, app *tests.TestApp, _ *http.Response) {
				if got := jobRateSheet(tb, app, "u09fwwcg07y03m7"); got != "c41ofep525bcacj" {
					tb.Errorf("jobs.rate_sheet = %q; want c41ofep525bcacj", got)
				}
				row, err := app.FindFirstRecordByFilter("job_rate_sheets", "job = 'u09fwwcg07y03m7'")
				if err != nil {
					tb.Fatal(err)
				}
				if row.GetString("created_by") != "f2j5a8vk006baub" {
					tb.Errorf("created_by = %q; want the job holder", row.GetString("created_by"))
				}
			},
		},
		{
			Name:            "a first rate sheet cannot replace an existing one",
			Method:          http.MethodPost,
			URL:             "/api/jobs/test_job_w_rs/rate_sheets",
			Body:            strings.NewReader(`{"rate_sheet":"test_rs_alt_actv"}`),
			Headers:         headers,
			ExpectedStatus:  http.StatusBadRequest,
			ExpectedContent: []string{"Enter the date the new rate sheet starts"},
			TestAppFactory:  testutils.SetupTestApp,
		},
		{
			Name:            "proposals have no rate sheet history",
			Method:          http.MethodGet,
			URL:             "/api/jobs/fcpropimpinprog/rate_sheets",
			Headers:         headers,
			ExpectedStatus:  http.StatusBadRequest,
			ExpectedContent: []string{"Proposals do not have rate sheets"},
			TestAppFactory:  testutils.SetupTestApp,
		},
		{
			Name:            "authentication is required",
			Method:          http.MethodGet,
			URL:             "/api/jobs/test_job_w_rs/rate_sheets",
			ExpectedStatus:  http.StatusUnauthorized,
			ExpectedContent: []string{"{"},
			TestAppFactory:  testutils.SetupTestApp,
		},
	}
	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
