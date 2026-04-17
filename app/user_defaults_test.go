package main

import (
	"encoding/json"
	"net/http"
	"testing"
	"tybalt/internal/testutils"

	"github.com/pocketbase/pocketbase/tests"
)

type userDefaultsTestResponse struct {
	DefaultDivision       string `json:"default_division"`
	DefaultWorkRecordType string `json:"default_work_record_type"`
	DefaultRole           string `json:"default_role"`
}

func TestUserDefaultsRoute(t *testing.T) {
	workRecordToken, err := testutils.GenerateRecordToken("users", "author@soup.com")
	if err != nil {
		t.Fatal(err)
	}

	scenario := tests.ApiScenario{
		Name:   "user defaults include default work record type",
		Method: http.MethodGet,
		URL:    "/api/users/defaults",
		Headers: map[string]string{
			"Authorization": workRecordToken,
		},
		ExpectedStatus: http.StatusOK,
		ExpectedContent: []string{
			`"default_work_record_type":"wrtype000000001"`,
		},
		TestAppFactory: testutils.SetupTestApp,
		AfterTestFunc: func(t testing.TB, _ *tests.TestApp, res *http.Response) {
			defer res.Body.Close()

			var payload userDefaultsTestResponse
			if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
				t.Fatalf("failed decoding user defaults response: %v", err)
			}
			if payload.DefaultWorkRecordType != "wrtype000000001" {
				t.Fatalf("expected seeded default work record type, got %#v", payload)
			}
		},
	}

	scenario.Test(t)
}
