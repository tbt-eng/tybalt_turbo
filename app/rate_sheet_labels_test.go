package main

import (
	"net/http"
	"strings"
	"testing"

	"tybalt/hooks"
	"tybalt/internal/testutils"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/tests"
)

// The server names a sheet by its code and revision, with its name in quotes
// when it has one (issue 207).
func TestRateSheetLabel(t *testing.T) {
	if got := hooks.RateSheetLabel("2026-11-00004", "Client Rates", 2); got != `2026-11-00004 (rev. 2) "Client Rates"` {
		t.Fatalf("label = %s", got)
	}
	if got := hooks.RateSheetLabel("2026-11-00004", "  ", 0); got != "2026-11-00004 (rev. 0)" {
		t.Fatalf("label without a name = %s", got)
	}
}

// Job details, the rate sheet history and a duplicate refusal carry the code.
// jobvaluation001 uses rsvaluation0001 (code 2026-09-00001); the 2025 Standard
// Rates (c41ofep525bcacj) have code 2026-01-00001.
func TestRateSheetResponsesCarryTheCode(t *testing.T) {
	app := testutils.SetupTestApp(t)
	token, err := testutils.GenerateRecordToken("users", "author@soup.com")
	if err != nil {
		t.Fatal(err)
	}
	runOnApp(t, app, tests.ApiScenario{
		Name:            "job details",
		Method:          http.MethodGet,
		URL:             "/api/jobs/jobvaluation001/details",
		Headers:         map[string]string{"Authorization": token},
		ExpectedStatus:  http.StatusOK,
		ExpectedContent: []string{`"rate_sheet":{"id":"rsvaluation0001","code":"2026-09-00001","name":"Summary valuation fixture","revision":0}`},
	})
	runOnApp(t, app, tests.ApiScenario{
		Name:            "rate sheet history",
		Method:          http.MethodGet,
		URL:             "/api/jobs/jobvaluation001/rate_sheets",
		Headers:         map[string]string{"Authorization": token},
		ExpectedStatus:  http.StatusOK,
		ExpectedContent: []string{`"rate_sheet":{"id":"rsvaluation0001","code":"2026-09-00001","name":"Summary valuation fixture"`},
	})

	entries, err := app.FindAllRecords("rate_sheet_entries", dbx.HashExp{"rate_sheet": "c41ofep525bcacj"})
	if err != nil {
		t.Fatal(err)
	}
	copied := []map[string]any{}
	for _, entry := range entries {
		copied = append(copied, map[string]any{"role": entry.GetString("role"), "rate": entry.GetFloat("rate"), "overtime_rate": entry.GetFloat("overtime_rate")})
	}
	status, body := createRateSheetThroughRoute(t, app, "author@soup.com", map[string]any{"name": "A Copy", "revision": 0, "entries": copied})
	if status != http.StatusConflict {
		t.Fatalf("a copy of an active sheet: status %d", status)
	}
	message := body["data"].(map[string]any)["entries"].(map[string]any)["message"].(string)
	if !strings.Contains(message, `2026-01-00001 (rev. 0) "2025 Standard Rates"`) {
		t.Fatalf("the refusal does not name the sheet by its code: %s", message)
	}
}

// The review token leaves out the code: the request binds its sheet by id,
// and only a superuser can change a code, so an approval made before codes
// existed, or before a code changed, is not stale. This test changes the
// draft's code with SQL after the approval, without touching updated, to
// show that the code alone is not in the digest.
func TestReviewTokenLeavesOutTheRateSheetCode(t *testing.T) {
	app := newWorkflowTestApp(t, true)
	id, draftID := proposedRequest(t, app, 295, "Client MSA, schedule B")
	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/submit", workflowAuthor, requestActionData(t, app, id), http.StatusOK)
	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/approve", workflowApprover, ratesApprovalData(t, app, id), http.StatusOK)

	if _, err := app.DB().NewQuery("UPDATE rate_sheets SET code = '2030-01-00001' WHERE id = {:id}").
		Bind(dbx.Params{"id": draftID}).Execute(); err != nil {
		t.Fatal(err)
	}

	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/create", workflowAuthor, requestActionData(t, app, id), http.StatusOK)
	if !workflowRecord(t, app, "rate_sheets", draftID).GetBool("active") {
		t.Fatal("the proposed sheet was not activated")
	}
}
