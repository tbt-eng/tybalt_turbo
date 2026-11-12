package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"tybalt/internal/testutils"

	"github.com/pocketbase/pocketbase/tests"
)

// TestCompareRateSheet_Claims checks who may call POST /api/rate_sheets/compare:
// a holder of job, rate_sheet_revise or admin. Seeded users (app/test_seed_data):
//   - time@test.com holds rate_sheet_revise but not job or admin
//   - admin.only@example.com holds admin and nothing else
//   - noclaims@example.com holds no claims
func TestCompareRateSheet_Claims(t *testing.T) {
	reviseToken, err := testutils.GenerateRecordToken("users", "time@test.com")
	if err != nil {
		t.Fatal(err)
	}
	adminToken, err := testutils.GenerateRecordToken("users", "admin.only@example.com")
	if err != nil {
		t.Fatal(err)
	}
	noClaimsToken, err := testutils.GenerateRecordToken("users", "noclaims@example.com")
	if err != nil {
		t.Fatal(err)
	}

	body, err := json.Marshal(map[string]any{
		"name":    "Copy of 2025 Standard Rates",
		"entries": standardSheetEntries(t),
	})
	if err != nil {
		t.Fatal(err)
	}
	exact := `"exact":[{"id":"c41ofep525bcacj","name":"2025 Standard Rates","revision":0,"effective_date":"2025-01-07","differing_roles":0}]`

	scenarios := []tests.ApiScenario{
		{
			Name:            "a user with only rate_sheet_revise can compare",
			Method:          http.MethodPost,
			URL:             "/api/rate_sheets/compare",
			Body:            strings.NewReader(string(body)),
			Headers:         map[string]string{"Authorization": reviseToken},
			ExpectedStatus:  http.StatusOK,
			ExpectedContent: []string{exact},
			ExpectedEvents:  map[string]int{"OnRecordCreate": 0, "OnRecordUpdate": 0},
			TestAppFactory:  testutils.SetupTestApp,
		},
		{
			Name:            "a user with only admin can compare",
			Method:          http.MethodPost,
			URL:             "/api/rate_sheets/compare",
			Body:            strings.NewReader(string(body)),
			Headers:         map[string]string{"Authorization": adminToken},
			ExpectedStatus:  http.StatusOK,
			ExpectedContent: []string{exact},
			ExpectedEvents:  map[string]int{"OnRecordCreate": 0, "OnRecordUpdate": 0},
			TestAppFactory:  testutils.SetupTestApp,
		},
		{
			Name:            "a user without job, rate_sheet_revise or admin cannot compare",
			Method:          http.MethodPost,
			URL:             "/api/rate_sheets/compare",
			Body:            strings.NewReader(string(body)),
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
