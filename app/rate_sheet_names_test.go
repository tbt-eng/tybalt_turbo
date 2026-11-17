package main

import (
	"net/http"
	"slices"
	"strings"
	"testing"

	"tybalt/internal/testutils"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/tests"
)

// The name of a rate sheet is an optional description (issue 207). Two
// sheets can't have the same name in any letter case; several may have none.
// The fixture draft rsdraftsheet001 holds "Proposed Fixture Rates"; "Sheet A"
// (rs_sheeta_r0, rs_sheeta_r1) and "Sheet B" are ordinary sheets.

// Finding 1 of the review of issue 187: a sheet could take a pending
// proposal's name with other capitals, because the create route compared
// exact text, and job creation then activated the proposal beside it. Now the
// create route, the collection API and the model hook all compare names
// without regard to capitals. The refusal says that a proposal has the name.
func TestANewSheetCannotTakeAProposalsNameInOtherCapitals(t *testing.T) {
	app := testutils.SetupTestApp(t)

	status, body := createRateSheetThroughRoute(t, app, "author@soup.com", map[string]any{"name": "proposed fixture rates", "revision": 0})
	nameErr, _ := body["data"].(map[string]any)["name"].(map[string]any)
	if status != http.StatusConflict || nameErr["code"] != "rate_sheet_name_proposed" ||
		!strings.Contains(nameErr["message"].(string), proposalHoldsNameMessage) {
		t.Fatalf("the create route accepted a proposal's name in other capitals: %d %v", status, body)
	}

	token, err := testutils.GenerateRecordToken("users", "author@soup.com")
	if err != nil {
		t.Fatal(err)
	}
	runOnApp(t, app, tests.ApiScenario{
		Name:            "the collection API refuses a proposal's name in other capitals",
		Method:          http.MethodPost,
		URL:             "/api/collections/rate_sheets/records",
		Body:            strings.NewReader(`{"name":"PROPOSED FIXTURE RATES","revision":0,"effective_date":"2026-11-20"}`),
		Headers:         map[string]string{"Authorization": token},
		ExpectedStatus:  http.StatusConflict,
		ExpectedContent: []string{`"rate_sheet_name_proposed"`, proposalHoldsNameMessage},
	})

	count, err := app.CountRecords("rate_sheets", dbx.NewExp("LOWER(name) = 'proposed fixture rates'"))
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("%d sheets have the proposal's name", count)
	}
}

// proposalHoldsNameMessage is the refusal of a name that a project request's
// proposed rates hold, without its first letter, which the API capitalizes in
// some responses.
const proposalHoldsNameMessage = "project request's proposed rates already have this name"

// Every path that names a sheet refuses a proposal's name with the same
// message, which says that a proposal has the name: a rename, another
// request's proposed rates, and the rename of proposed rates when the job is
// created. The fixture draft rsdraftsheet002 belongs to request
// rsdraftreq00002, which time@test.com created.
func TestAProposalsNameIsRefusedWithTheSameMessageEverywhere(t *testing.T) {
	app := newWorkflowTestApp(t, true)

	runOnApp(t, app, renameScenario(renameToken(t, "author@soup.com"), "rs_sheetb_r0", "PROPOSED FIXTURE RATES",
		http.StatusConflict, `"rate_sheet_name_proposed"`, proposalHoldsNameMessage))

	result := proposeRates(t, app, "rsdraftreq00002", draftFixtureCreator,
		proposedRates("proposed fixture rates", "Client rate schedule", standardRates(t, app, 265)), http.StatusBadRequest)
	nameErr, _ := result["data"].(map[string]any)["proposed_rates.name"].(map[string]any)
	if nameErr["code"] != "rate_sheet_name_proposed" || !strings.Contains(nameErr["message"].(string), proposalHoldsNameMessage) {
		t.Fatalf("the proposal save refused the name with %#v", result)
	}

	id, draftID := proposedRequest(t, app, 290, "Client MSA, schedule B")
	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/submit", workflowAuthor, requestActionData(t, app, id), http.StatusOK)
	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/approve", workflowApprover, ratesApprovalData(t, app, id), http.StatusOK)
	data := requestActionData(t, app, id)
	data["rate_sheet_name"] = "Proposed Fixture Rates"
	refused := requestJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/create", workflowAuthor, data, http.StatusBadRequest)
	if message, _ := refused["message"].(string); !strings.Contains(message, proposalHoldsNameMessage) {
		t.Fatalf("job creation refused the name with %#v", refused)
	}
	if workflowRecord(t, app, "rate_sheets", draftID).GetBool("active") {
		t.Fatal("the refused job creation activated the proposed rates")
	}
}

// The create route compares names of ordinary sheets without regard to
// capitals too; it compared exact text before.
func TestRateSheetNamesAreUniqueInAnyCapitals(t *testing.T) {
	app := testutils.SetupTestApp(t)
	for _, name := range []string{"sheet a", "SHEET B", "  Sheet A  "} {
		status, body := createRateSheetThroughRoute(t, app, "author@soup.com", map[string]any{"name": name, "revision": 0})
		if status != http.StatusConflict {
			t.Fatalf("%q: status %d %v; want 409", name, status, body)
		}
	}
	// A rename of a sheet's own capitals is allowed; another sheet's name
	// is not, in any capitals.
	adminToken := renameToken(t, "author@soup.com")
	runOnApp(t, app, renameScenario(adminToken, "rs_sheetb_r0", "SHEET B", http.StatusOK, `"name":"SHEET B"`))
	runOnApp(t, app, renameScenario(adminToken, "rs_sheetb_r0", "sheet a", http.StatusConflict, `"validation_not_unique"`))
}

// Several sheets may have no name. Each still has its own code.
func TestRateSheetsMayHaveNoName(t *testing.T) {
	app := testutils.SetupTestApp(t)
	codes := []string{}
	for i, rate := range []float64{7771, 7772} {
		roles, err := app.FindRecordsByFilter("rate_roles", "1=1", "name", 1, 0, nil)
		if err != nil {
			t.Fatal(err)
		}
		status, body := createRateSheetThroughRoute(t, app, "author@soup.com", map[string]any{
			"revision": 0,
			"entries":  []map[string]any{{"role": roles[0].Id, "rate": rate, "overtime_rate": rate}},
		})
		if status != http.StatusCreated {
			t.Fatalf("sheet %d without a name: status %d %v", i, status, body)
		}
		if body["name"] != "" {
			t.Fatalf("sheet %d has name %v", i, body["name"])
		}
		codes = append(codes, body["code"].(string))
	}
	if codes[0] == codes[1] {
		t.Fatalf("two sheets share code %s", codes[0])
	}

	// A rename to no name removes the description and keeps the old name
	// in the history.
	adminToken := renameToken(t, "author@soup.com")
	runOnApp(t, app, renameScenario(adminToken, "rs_sheetb_r0", "", http.StatusOK, `"name":""`, `"former_names":["Sheet B"]`))
	if got := formerNamesOf(t, app, "rs_sheetb_r0"); !slices.Equal(got, []string{"Sheet B"}) {
		t.Fatalf("former names = %v", got)
	}
	// It can then take a name again, and the empty name is not a former name.
	runOnApp(t, app, renameScenario(adminToken, "rs_sheetb_r0", "Sheet B Again", http.StatusOK, `"former_names":["Sheet B"]`))
}

// Proposed rates may have no name, and get a code on their first save.
func TestProposedRatesMayHaveNoNameAndGetACode(t *testing.T) {
	app := newWorkflowTestApp(t, true)
	id := createWorkflowRequest(t, app)
	completeWorkflowRequest(t, app, id)
	want := nextRateSheetCode(t, app)
	proposeRates(t, app, id, workflowAuthor, proposedRates("", "Client MSA, schedule B", standardRates(t, app, 270)), http.StatusOK)

	draft := workflowRecord(t, app, "rate_sheets", workflowRecord(t, app, "job_requests", id).GetString("rate_sheet"))
	if draft.GetString("name") != "" || draft.GetString("code") != want {
		t.Fatalf("draft name %q code %q; want no name and code %s", draft.GetString("name"), draft.GetString("code"), want)
	}

	// A later save keeps the code, with or without a name.
	proposeRates(t, app, id, workflowAuthor, proposedRates("Client Negotiated Rates", "Client MSA, schedule B", standardRates(t, app, 280)), http.StatusOK)
	if got := rateSheetCodeOf(t, app, draft.Id); got != want {
		t.Fatalf("a later save changed the code to %s", got)
	}
}
