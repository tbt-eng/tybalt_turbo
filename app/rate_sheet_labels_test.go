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

// The review token leaves out the sheet's code, name and updated time: the
// request binds its sheet by id, so a code change or a rename after the
// approval does not make it stale. Nobody can rename a proposal through the
// rename route, and only a superuser can change a code, so this test changes
// the draft's code, name and updated time with SQL after the approval.
func TestReviewTokenLeavesOutTheRateSheetCodeAndName(t *testing.T) {
	app := newWorkflowTestApp(t, true)
	id, draftID := proposedRequest(t, app, 295, "Client MSA, schedule B")
	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/submit", workflowAuthor, requestActionData(t, app, id), http.StatusOK)
	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/approve", workflowApprover, ratesApprovalData(t, app, id), http.StatusOK)

	if _, err := app.DB().NewQuery("UPDATE rate_sheets SET code = '2030-01-00001', name = 'Renamed Proposal', updated = '2030-01-01 00:00:00.000Z' WHERE id = {:id}").
		Bind(dbx.Params{"id": draftID}).Execute(); err != nil {
		t.Fatal(err)
	}

	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/create", workflowAuthor, requestActionData(t, app, id), http.StatusOK)
	if !workflowRecord(t, app, "rate_sheets", draftID).GetBool("active") {
		t.Fatal("the proposed sheet was not activated")
	}
}

// approvedRequestOnTheStandardSheet makes a request that chose the active 2025
// Standard Rates (c41ofep525bcacj), has it approved, and returns its id.
func approvedRequestOnTheStandardSheet(t *testing.T, app *tests.TestApp) string {
	t.Helper()
	id := createWorkflowRequest(t, app)
	completeWorkflowRequest(t, app, id)
	if got := workflowRecord(t, app, "job_requests", id).GetString("rate_sheet"); got != standardRateSheet {
		t.Fatalf("the request uses %s", got)
	}
	uploadRequestDocument(t, app, id)
	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/submit", workflowAuthor, requestActionData(t, app, id), http.StatusOK)
	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/approve", workflowApprover, requestActionData(t, app, id), http.StatusOK)
	return id
}

// A rename of the chosen active sheet after the approval leaves the review
// token as it was, and the job is created. This test renames the fixture sheet
// 2025 Standard Rates through the rename route (author@soup.com holds admin).
func TestApprovalSurvivesARenameOfTheChosenSheet(t *testing.T) {
	app := newWorkflowTestApp(t, true)
	id := approvedRequestOnTheStandardSheet(t, app)
	before := requestActionData(t, app, id)["review_token"]

	clientSetupJSON(t, app, http.MethodPost, "/api/rate_sheets/"+standardRateSheet+"/rename", workflowAuthor, map[string]any{"name": "2025 Standard Rates Renamed"}, http.StatusOK)
	if name := workflowRecord(t, app, "rate_sheets", standardRateSheet).GetString("name"); name != "2025 Standard Rates Renamed" {
		t.Fatalf("the sheet was not renamed: %q", name)
	}
	if after := requestActionData(t, app, id)["review_token"]; after != before {
		t.Fatalf("the rename changed the review token: %v to %v", before, after)
	}

	result := clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/create", workflowAuthor, requestActionData(t, app, id), http.StatusOK)
	if job := workflowRecord(t, app, "jobs", result["job_id"].(string)); job.GetString("rate_sheet") != standardRateSheet {
		t.Fatalf("job uses %s", job.GetString("rate_sheet"))
	}
}

// A changed rate on the chosen active sheet after the approval changes what was
// approved, so job creation refuses the stale approval. This test changes a
// fixture rate of 2025 Standard Rates through the entry route
// (author@soup.com holds admin).
func TestApprovalIsStaleAfterAChosenSheetRateChanges(t *testing.T) {
	app := newWorkflowTestApp(t, true)
	id := approvedRequestOnTheStandardSheet(t, app)

	entry, err := app.FindFirstRecordByFilter("rate_sheet_entries", "rate_sheet = {:sheet} && role = {:role}", dbx.Params{"sheet": standardRateSheet, "role": principalRole})
	if err != nil {
		t.Fatal(err)
	}
	clientSetupJSON(t, app, http.MethodPut, "/api/rate_sheet_entries/"+entry.Id, workflowAuthor, map[string]any{"rate": entry.GetFloat("rate") + 1}, http.StatusOK)

	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/create", workflowAuthor, requestActionData(t, app, id), http.StatusConflict)
	if r := workflowRecord(t, app, "job_requests", id); r.GetString("status") != "Awaiting creation" {
		t.Fatalf("a refused creation changed the request: %s", r.GetString("status"))
	}
}
