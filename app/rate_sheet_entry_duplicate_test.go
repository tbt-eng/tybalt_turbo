package main

import (
	"net/http"
	"strings"
	"testing"
	"tybalt/internal/testutils"

	"github.com/pocketbase/pocketbase/tests"
)

// Fixtures used here (app/test_seed_data/data), all rows of the seed CSV
// files. Nothing is changed at test time before the request.
//   - rsentrybase0001: "Entry Correction Base", active. Its rates are the
//     standard rates plus 5000, so that it matches no other sheet.
//   - rsentryactive01: "Entry Correction Active", active. It differs from the
//     base sheet in one role: rate 5161, where the base sheet has 5160.
//   - rsentryinact001: "Entry Correction Inactive", inactive, with the same
//     rates as rsentryactive01.
//
// A request that sets that rate to 5160 makes a sheet identical to the base
// sheet.

const (
	dupEntryActiveID   = "eentryactive01_0l1llckilwz1mqz"
	dupEntryInactiveID = "eentryinact001_0l1llckilwz1mqz"
)

func assertDupEntryRate(id string, want int) func(testing.TB, *tests.TestApp, *http.Response) {
	return func(t testing.TB, app *tests.TestApp, res *http.Response) {
		entry, err := app.FindRecordById("rate_sheet_entries", id)
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

	activeURL := "/api/rate_sheet_entries/" + dupEntryActiveID
	inactiveURL := "/api/rate_sheet_entries/" + dupEntryInactiveID

	scenarios := []tests.ApiScenario{
		{
			Name:           "a rate that makes an active sheet identical to another active sheet is not saved",
			Method:         http.MethodPut,
			URL:            activeURL,
			Body:           strings.NewReader(`{"rate": 5160}`),
			Headers:        map[string]string{"Authorization": adminToken},
			ExpectedStatus: http.StatusConflict,
			ExpectedContent: []string{
				`"status":409`,
				`"code":"duplicate_rate_sheet"`,
				`"entry":{`,
				`"id":"rsentrybase0001"`,
				`"name":"Entry Correction Base"`,
			},
			ExpectedEvents: map[string]int{"OnRecordUpdate": 0},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc:  assertDupEntryRate(dupEntryActiveID, 5161),
		},
		{
			Name:            "a rate that keeps an active sheet distinct is saved",
			Method:          http.MethodPut,
			URL:             activeURL,
			Body:            strings.NewReader(`{"rate": 5162}`),
			Headers:         map[string]string{"Authorization": adminToken},
			ExpectedStatus:  http.StatusOK,
			ExpectedContent: []string{`"message":"rate sheet entry updated"`, `"rate":5162`},
			TestAppFactory:  testutils.SetupTestApp,
			AfterTestFunc:   assertDupEntryRate(dupEntryActiveID, 5162),
		},
		{
			Name:            "a rate on an inactive sheet is saved even when it matches an active sheet",
			Method:          http.MethodPut,
			URL:             inactiveURL,
			Body:            strings.NewReader(`{"rate": 5160}`),
			Headers:         map[string]string{"Authorization": adminToken},
			ExpectedStatus:  http.StatusOK,
			ExpectedContent: []string{`"message":"rate sheet entry updated"`, `"rate":5160`},
			TestAppFactory:  testutils.SetupTestApp,
			AfterTestFunc:   assertDupEntryRate(dupEntryInactiveID, 5160),
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
