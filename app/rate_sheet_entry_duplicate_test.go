package main

import (
	"net/http"
	"strings"
	"testing"
	"tybalt/internal/testutils"

	"github.com/pocketbase/pocketbase/tests"
)

// Fixtures used here (app/test_seed_data):
//   - c41ofep525bcacj: "2025 Standard Rates", active
//   - rs_dupact_001: "Duplicate Activation Sheet", inactive, the same entries
//     as c41ofep525bcacj. Its entry rsedupact_tbgoiwwwfj8cvju has rate 250
//     and overtime rate 325, as on the standard sheet.
//
// Each test app changes that entry's rate to 251, so the sheet differs from
// the standard sheet in one role. A request that sets it back to 250 makes the
// two sheets identical.

const dupEntryID = "rsedupact_tbgoiwwwfj8cvju"

// setupDuplicateEntryApp returns a test app factory where rs_dupact_001
// differs from the standard sheet in one rate, and is active when active is
// true. It saves without request hooks, so the duplicate check does not run.
func setupDuplicateEntryApp(active bool) func(testing.TB) *tests.TestApp {
	return func(t testing.TB) *tests.TestApp {
		app := testutils.SetupTestApp(t)
		entry, err := app.FindRecordById("rate_sheet_entries", dupEntryID)
		if err != nil {
			t.Fatalf("failed to load entry: %v", err)
		}
		entry.Set("rate", 251)
		if err := app.Save(entry); err != nil {
			t.Fatalf("failed to change entry: %v", err)
		}
		if active {
			sheet, err := app.FindRecordById("rate_sheets", "rs_dupact_001")
			if err != nil {
				t.Fatalf("failed to load sheet: %v", err)
			}
			sheet.Set("active", true)
			if err := app.Save(sheet); err != nil {
				t.Fatalf("failed to activate sheet: %v", err)
			}
		}
		return app
	}
}

func assertDupEntryRate(want int) func(testing.TB, *tests.TestApp, *http.Response) {
	return func(t testing.TB, app *tests.TestApp, res *http.Response) {
		entry, err := app.FindRecordById("rate_sheet_entries", dupEntryID)
		if err != nil {
			t.Fatalf("failed to load entry: %v", err)
		}
		if got := entry.GetInt("rate"); got != want {
			t.Errorf("expected rate %d, got %d", want, got)
		}
	}
}

func TestUpdateRateSheetEntry_DuplicateRates(t *testing.T) {
	adminToken, err := testutils.GenerateRecordToken("users", "author@soup.com")
	if err != nil {
		t.Fatal(err)
	}

	url := "/api/rate_sheet_entries/" + dupEntryID

	scenarios := []tests.ApiScenario{
		{
			Name:           "a rate that makes an active sheet identical to another active sheet is not saved",
			Method:         http.MethodPut,
			URL:            url,
			Body:           strings.NewReader(`{"rate": 250}`),
			Headers:        map[string]string{"Authorization": adminToken},
			ExpectedStatus: http.StatusConflict,
			ExpectedContent: []string{
				`"status":409`,
				`"code":"duplicate_rate_sheet"`,
				`"entry":{`,
				`"id":"c41ofep525bcacj"`,
				`"name":"2025 Standard Rates"`,
			},
			ExpectedEvents: map[string]int{"OnRecordUpdate": 0},
			TestAppFactory: setupDuplicateEntryApp(true),
			AfterTestFunc:  assertDupEntryRate(251),
		},
		{
			Name:            "a rate that keeps an active sheet distinct is saved",
			Method:          http.MethodPut,
			URL:             url,
			Body:            strings.NewReader(`{"rate": 252}`),
			Headers:         map[string]string{"Authorization": adminToken},
			ExpectedStatus:  http.StatusOK,
			ExpectedContent: []string{`"message":"rate sheet entry updated"`, `"rate":252`},
			TestAppFactory:  setupDuplicateEntryApp(true),
			AfterTestFunc:   assertDupEntryRate(252),
		},
		{
			Name:            "a rate on an inactive sheet is saved even when it matches an active sheet",
			Method:          http.MethodPut,
			URL:             url,
			Body:            strings.NewReader(`{"rate": 250}`),
			Headers:         map[string]string{"Authorization": adminToken},
			ExpectedStatus:  http.StatusOK,
			ExpectedContent: []string{`"message":"rate sheet entry updated"`, `"rate":250`},
			TestAppFactory:  setupDuplicateEntryApp(false),
			AfterTestFunc:   assertDupEntryRate(250),
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
