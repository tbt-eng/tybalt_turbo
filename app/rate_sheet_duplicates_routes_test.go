package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"tybalt/internal/testutils"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

// Fixtures used here (app/test_seed_data):
//   - c41ofep525bcacj: "2025 Standard Rates", active, an entry for every role
//   - rs_dupact_001: "Duplicate Activation Sheet", inactive, the same entries
//     as c41ofep525bcacj
//   - rs_complete_all_001: "Complete Hook Sheet", inactive, distinct rates
//   - rs_revact_r0 (active) and rs_revact_r1 (inactive): "Revision Activation
//     Sheet" revisions 0 and 1, with identical entries

type rateSheetBodyEntry struct {
	Role         string  `json:"role"`
	Rate         float64 `json:"rate"`
	OvertimeRate float64 `json:"overtime_rate"`
}

// standardSheetEntries returns the entries of the active 2025 Standard Rates
// sheet, in the form of a create request.
func standardSheetEntries(t *testing.T) []rateSheetBodyEntry {
	t.Helper()
	app := testutils.SetupTestApp(t)
	records, err := app.FindRecordsByFilter("rate_sheet_entries", "rate_sheet = 'c41ofep525bcacj'", "", 0, 0, nil)
	if err != nil {
		t.Fatalf("failed to load entries: %v", err)
	}
	if len(records) == 0 {
		t.Fatal("expected entries on the standard sheet")
	}
	entries := make([]rateSheetBodyEntry, 0, len(records))
	for _, r := range records {
		entries = append(entries, rateSheetBodyEntry{
			Role:         r.GetString("role"),
			Rate:         r.GetFloat("rate"),
			OvertimeRate: r.GetFloat("overtime_rate"),
		})
	}
	return entries
}

func rateSheetBody(t *testing.T, name string, effectiveDate string, revision int, entries []rateSheetBodyEntry) string {
	t.Helper()
	body, err := json.Marshal(map[string]any{
		"name":           name,
		"effective_date": effectiveDate,
		"revision":       revision,
		"entries":        entries,
	})
	if err != nil {
		t.Fatal(err)
	}
	return string(body)
}

func countRateSheetsNamed(t testing.TB, app core.App, name string) int {
	t.Helper()
	records, err := app.FindRecordsByFilter("rate_sheets", "name = {:name}", "", 0, 0, map[string]any{"name": name})
	if err != nil {
		t.Fatalf("failed to count rate sheets: %v", err)
	}
	return len(records)
}

func TestCreateRateSheet_DuplicateRates(t *testing.T) {
	jobToken, err := testutils.GenerateRecordToken("users", "author@soup.com")
	if err != nil {
		t.Fatal(err)
	}
	reviseToken, err := testutils.GenerateRecordToken("users", "time@test.com")
	if err != nil {
		t.Fatal(err)
	}

	identical := standardSheetEntries(t)
	changed := standardSheetEntries(t)
	changed[0].Rate++

	scenarios := []tests.ApiScenario{
		{
			Name:           "a new sheet identical to an active sheet is blocked and names it",
			Method:         http.MethodPost,
			URL:            "/api/rate_sheets",
			Body:           strings.NewReader(rateSheetBody(t, "Copy of 2025 Standard Rates", "2026-03-01", 0, identical)),
			Headers:        map[string]string{"Authorization": jobToken},
			ExpectedStatus: http.StatusConflict,
			ExpectedContent: []string{
				`"status":409`,
				`"code":"duplicate_rate_sheet"`,
				`"id":"c41ofep525bcacj"`,
				`"name":"2025 Standard Rates"`,
				`"revision":0`,
			},
			ExpectedEvents: map[string]int{"OnRecordCreate": 0},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				if n := countRateSheetsNamed(t, app, "Copy of 2025 Standard Rates"); n != 0 {
					t.Errorf("expected no sheet to be saved, got %d", n)
				}
			},
		},
		{
			Name:            "a new sheet that differs in one rate is saved",
			Method:          http.MethodPost,
			URL:             "/api/rate_sheets",
			Body:            strings.NewReader(rateSheetBody(t, "Client X Rates", "2026-03-01", 0, changed)),
			Headers:         map[string]string{"Authorization": jobToken},
			ExpectedStatus:  http.StatusCreated,
			ExpectedContent: []string{`"name":"Client X Rates"`},
			TestAppFactory:  testutils.SetupTestApp,
		},
		{
			Name:            "a revision with the same rates as its earlier revision is saved",
			Method:          http.MethodPost,
			URL:             "/api/rate_sheets",
			Body:            strings.NewReader(rateSheetBody(t, "2025 Standard Rates", "2026-01-01", 1, identical)),
			Headers:         map[string]string{"Authorization": reviseToken},
			ExpectedStatus:  http.StatusCreated,
			ExpectedContent: []string{`"name":"2025 Standard Rates"`, `"revision":1`},
			TestAppFactory:  testutils.SetupTestApp,
		},
		{
			Name:            "a revision identical to another active sheet is blocked",
			Method:          http.MethodPost,
			URL:             "/api/rate_sheets",
			Body:            strings.NewReader(rateSheetBody(t, "Revise Test Sheet", "2026-02-01", 1, identical)),
			Headers:         map[string]string{"Authorization": reviseToken},
			ExpectedStatus:  http.StatusConflict,
			ExpectedContent: []string{`"code":"duplicate_rate_sheet"`, `"id":"c41ofep525bcacj"`},
			ExpectedEvents:  map[string]int{"OnRecordCreate": 0},
			TestAppFactory:  testutils.SetupTestApp,
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestActivateRateSheet_DuplicateRates(t *testing.T) {
	jobToken, err := testutils.GenerateRecordToken("users", "author@soup.com")
	if err != nil {
		t.Fatal(err)
	}

	assertActive := func(id string, want bool) func(testing.TB, *tests.TestApp, *http.Response) {
		return func(t testing.TB, app *tests.TestApp, res *http.Response) {
			record, err := app.FindRecordById("rate_sheets", id)
			if err != nil {
				t.Fatalf("failed to load %s: %v", id, err)
			}
			if record.GetBool("active") != want {
				t.Errorf("expected %s active=%v, got %v", id, want, record.GetBool("active"))
			}
		}
	}

	scenarios := []tests.ApiScenario{
		{
			Name:           "a sheet identical to another active sheet cannot be activated",
			Method:         http.MethodPost,
			URL:            "/api/rate_sheets/rs_dupact_001/activate",
			Headers:        map[string]string{"Authorization": jobToken},
			ExpectedStatus: http.StatusConflict,
			ExpectedContent: []string{
				`"code":"duplicate_rate_sheet"`,
				`"active":{`,
				`"id":"c41ofep525bcacj"`,
				`"name":"2025 Standard Rates"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc:  assertActive("rs_dupact_001", false),
		},
		{
			Name:            "a distinct complete sheet is activated",
			Method:          http.MethodPost,
			URL:             "/api/rate_sheets/rs_complete_all_001/activate",
			Headers:         map[string]string{"Authorization": jobToken},
			ExpectedStatus:  http.StatusOK,
			ExpectedContent: []string{`"message":"rate sheet activated"`},
			TestAppFactory:  testutils.SetupTestApp,
			AfterTestFunc:   assertActive("rs_complete_all_001", true),
		},
		{
			Name:            "a revision with the same rates as its active earlier revision is activated",
			Method:          http.MethodPost,
			URL:             "/api/rate_sheets/rs_revact_r1/activate",
			Headers:         map[string]string{"Authorization": jobToken},
			ExpectedStatus:  http.StatusOK,
			ExpectedContent: []string{`"message":"rate sheet activated"`},
			TestAppFactory:  testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				assertActive("rs_revact_r1", true)(t, app, res)
				assertActive("rs_revact_r0", false)(t, app, res)
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestCompareRateSheet(t *testing.T) {
	jobToken, err := testutils.GenerateRecordToken("users", "author@soup.com")
	if err != nil {
		t.Fatal(err)
	}
	noClaimsToken, err := testutils.GenerateRecordToken("users", "noclaims@example.com")
	if err != nil {
		t.Fatal(err)
	}

	identical := standardSheetEntries(t)
	changed := standardSheetEntries(t)
	changed[0].Rate++
	changed[1].OvertimeRate++

	compareBody := func(name string, entries []rateSheetBodyEntry) string {
		body, err := json.Marshal(map[string]any{"name": name, "entries": entries})
		if err != nil {
			t.Fatal(err)
		}
		return string(body)
	}

	scenarios := []tests.ApiScenario{
		{
			Name:           "identical entries are an exact match",
			Method:         http.MethodPost,
			URL:            "/api/rate_sheets/compare",
			Body:           strings.NewReader(compareBody("Copy of 2025 Standard Rates", identical)),
			Headers:        map[string]string{"Authorization": jobToken},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"exact":[{"id":"c41ofep525bcacj","name":"2025 Standard Rates","revision":0,"effective_date":"2025-01-07","differing_roles":0}]`,
			},
			ExpectedEvents: map[string]int{"OnRecordCreate": 0, "OnRecordUpdate": 0},
			TestAppFactory: testutils.SetupTestApp,
		},
		{
			Name:           "entries that differ in two roles are a near match",
			Method:         http.MethodPost,
			URL:            "/api/rate_sheets/compare",
			Body:           strings.NewReader(compareBody("Client X Rates", changed)),
			Headers:        map[string]string{"Authorization": jobToken},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"exact":[]`,
				`"near":[{"id":"c41ofep525bcacj","name":"2025 Standard Rates","revision":0,"effective_date":"2025-01-07","differing_roles":2}]`,
			},
			TestAppFactory: testutils.SetupTestApp,
		},
		{
			Name:               "sheets with the same name are not compared",
			Method:             http.MethodPost,
			URL:                "/api/rate_sheets/compare",
			Body:               strings.NewReader(compareBody("2025 Standard Rates", identical)),
			Headers:            map[string]string{"Authorization": jobToken},
			ExpectedStatus:     http.StatusOK,
			ExpectedContent:    []string{`"exact":[]`},
			NotExpectedContent: []string{`c41ofep525bcacj`},
			TestAppFactory:     testutils.SetupTestApp,
		},
		{
			Name:            "a user without the job claim cannot compare",
			Method:          http.MethodPost,
			URL:             "/api/rate_sheets/compare",
			Body:            strings.NewReader(compareBody("Client X Rates", identical)),
			Headers:         map[string]string{"Authorization": noClaimsToken},
			ExpectedStatus:  http.StatusForbidden,
			ExpectedContent: []string{`"status":403`},
			TestAppFactory:  testutils.SetupTestApp,
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
