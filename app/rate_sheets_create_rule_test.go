package main

import (
	"net/http"
	"strings"
	"testing"

	"tybalt/internal/testutils"

	"github.com/pocketbase/pocketbase/tests"
)

// The collection create rule of rate_sheets matches POST /api/rate_sheets:
// only a holder of the job claim creates a sheet there, and only a new one
// (revision 0). A revision needs a code, which the collection API can't
// send. time@test.com holds rate_sheet_revise and admin.only@example.com
// holds admin, neither with job; author@soup.com holds job and admin.
func TestRateSheetsCollectionCreateRule(t *testing.T) {
	app := testutils.SetupTestApp(t)
	token := func(email string) string {
		t.Helper()
		value, err := testutils.GenerateRecordToken("users", email)
		if err != nil {
			t.Fatal(err)
		}
		return value
	}
	create := func(name, email, body string, status int, content ...string) tests.ApiScenario {
		return tests.ApiScenario{
			Name:            name,
			Method:          http.MethodPost,
			URL:             "/api/collections/rate_sheets/records",
			Body:            strings.NewReader(body),
			Headers:         map[string]string{"Authorization": token(email)},
			ExpectedStatus:  status,
			ExpectedContent: content,
		}
	}

	scenarios := []tests.ApiScenario{
		create("rate_sheet_revise without job is refused a new sheet", "time@test.com",
			`{"name":"Rule Revise Sheet","revision":0,"effective_date":"2026-11-20"}`,
			http.StatusBadRequest, "Failed to create record"),
		create("admin without job is refused a new sheet", "admin.only@example.com",
			`{"name":"Rule Admin Sheet","revision":0,"effective_date":"2026-11-20"}`,
			http.StatusBadRequest, "Failed to create record"),
		create("rate_sheet_revise is refused a revision", "time@test.com",
			`{"name":"Rule Revise Revision","revision":1,"effective_date":"2026-11-20"}`,
			http.StatusBadRequest, "Failed to create record"),
		create("admin without job is refused a revision", "admin.only@example.com",
			`{"name":"Rule Admin Revision","revision":1,"effective_date":"2026-11-20"}`,
			http.StatusBadRequest, "Failed to create record"),
		create("job is refused a revision", "author@soup.com",
			`{"name":"Rule Job Revision","revision":1,"effective_date":"2026-11-20"}`,
			http.StatusBadRequest, "Failed to create record"),
		create("job creates a new sheet", "author@soup.com",
			`{"name":"Rule Job Sheet","revision":0,"effective_date":"2026-11-20"}`,
			http.StatusOK, `"name":"Rule Job Sheet"`, `"revision":0`),
	}
	for _, scenario := range scenarios {
		runOnApp(t, app, scenario)
	}

	for _, name := range []string{"Rule Revise Sheet", "Rule Admin Sheet", "Rule Revise Revision", "Rule Admin Revision", "Rule Job Revision"} {
		if _, err := app.FindFirstRecordByData("rate_sheets", "name", name); err == nil {
			t.Fatalf("the collection API created %q", name)
		}
	}
	if _, err := app.FindFirstRecordByData("rate_sheets", "name", "Rule Job Sheet"); err != nil {
		t.Fatalf("the job holder's sheet was not created: %v", err)
	}
}
