package main

import (
	"net/http"
	"strings"
	"testing"

	"tybalt/internal/testutils"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/tests"
)

// Revisions of a rate sheet are linked by code (issue 207). The draft
// rsdraftsheet001 ("Proposed Fixture Rates") belongs to request
// rsdraftreq00001. time@test.com holds rate_sheet_revise.

// Finding 2 of the review of issue 187: the collection API let a holder of
// rate_sheet_revise create revision 1 of the name a proposal held, and job
// creation then activated the proposal (revision 0) over it. With the code as
// identity no client can choose a code, and nothing can revise a draft, so
// the proposal can't get a newer revision on any path.
func TestNothingCanReviseAProposal(t *testing.T) {
	app := testutils.SetupTestApp(t)
	token, err := testutils.GenerateRecordToken("users", "time@test.com")
	if err != nil {
		t.Fatal(err)
	}
	draftCode := rateSheetCodeOf(t, app, "rsdraftsheet001")

	scenarios := []tests.ApiScenario{
		{
			Name:            "the collection API refuses a code",
			Method:          http.MethodPost,
			URL:             "/api/collections/rate_sheets/records",
			Body:            strings.NewReader(`{"name":"Proposed Fixture Rates","code":"` + draftCode + `","revision":1,"effective_date":"2026-11-20"}`),
			Headers:         map[string]string{"Authorization": token},
			ExpectedStatus:  http.StatusBadRequest,
			ExpectedContent: []string{"Failed to create record"},
		},
		{
			Name:            "the collection API refuses a revision without a code",
			Method:          http.MethodPost,
			URL:             "/api/collections/rate_sheets/records",
			Body:            strings.NewReader(`{"name":"Another Revision","revision":1,"effective_date":"2026-11-20"}`),
			Headers:         map[string]string{"Authorization": token},
			ExpectedStatus:  http.StatusBadRequest,
			ExpectedContent: []string{"Failed to create record"},
		},
		{
			Name:            "the create route refuses a revision of a draft by code",
			Method:          http.MethodPost,
			URL:             "/api/rate_sheets",
			Body:            strings.NewReader(`{"code":"` + draftCode + `","revision":1,"effective_date":"2026-11-20","entries":[{"role":"x","rate":1,"overtime_rate":1}]}`),
			Headers:         map[string]string{"Authorization": token},
			ExpectedStatus:  http.StatusBadRequest,
			ExpectedContent: []string{"rate_sheet_draft"},
		},
		{
			Name:            "the create route refuses a revision of a draft by name",
			Method:          http.MethodPost,
			URL:             "/api/rate_sheets",
			Body:            strings.NewReader(`{"name":"Proposed Fixture Rates","revision":1,"effective_date":"2026-11-20","entries":[{"role":"x","rate":1,"overtime_rate":1}]}`),
			Headers:         map[string]string{"Authorization": token},
			ExpectedStatus:  http.StatusBadRequest,
			ExpectedContent: []string{"rate_sheet_draft"},
		},
	}
	for _, scenario := range scenarios {
		runOnApp(t, app, scenario)
	}
	count, err := app.CountRecords("rate_sheets", dbx.HashExp{"code": draftCode})
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("the proposal's code has %d sheets; want only the proposal", count)
	}
	if _, err := app.FindFirstRecordByData("rate_sheets", "name", "Another Revision"); err == nil {
		t.Fatal("the collection API created a revision without a code")
	}
}

// A new sheet created through the collection API gets the next code.
func TestCollectionAPICreateGetsACode(t *testing.T) {
	app := testutils.SetupTestApp(t)
	token, err := testutils.GenerateRecordToken("users", "author@soup.com")
	if err != nil {
		t.Fatal(err)
	}
	want := nextRateSheetCode(t, app)
	runOnApp(t, app, tests.ApiScenario{
		Name:            "a job holder creates a sheet",
		Method:          http.MethodPost,
		URL:             "/api/collections/rate_sheets/records",
		Body:            strings.NewReader(`{"name":"Collection Sheet","revision":0,"effective_date":"2026-11-20"}`),
		Headers:         map[string]string{"Authorization": token},
		ExpectedStatus:  http.StatusOK,
		ExpectedContent: []string{`"code":"` + want + `"`},
	})
}

// A revision found by code keeps the sheet's code even after a rename, and
// the create route matches a client that sends only the name to the sheet
// with that name.
func TestRevisionByCodeAfterARename(t *testing.T) {
	app := testutils.SetupTestApp(t)
	code := rateSheetCodeOf(t, app, "rsrevbase0001")
	adminToken := renameToken(t, "author@soup.com")
	runOnApp(t, app, renameScenario(adminToken, "rsrevbase0001", "Renamed Revise Sheet", http.StatusOK, `"code":"`+code+`"`))

	status, revision := createRateSheetThroughRoute(t, app, "time@test.com", map[string]any{
		"code": code, "revision": 1, "effective_date": "2026-02-01",
	})
	if status != http.StatusCreated {
		t.Fatalf("create status %d", status)
	}
	if revision["code"] != code || revision["name"] != "Renamed Revise Sheet" {
		t.Fatalf("revision = %v; want code %s and the sheet's name", revision, code)
	}

	// A new sheet may not send a code.
	status, _ = createRateSheetThroughRoute(t, app, "author@soup.com", map[string]any{
		"code": code, "name": "Chosen Code", "revision": 0,
	})
	if status != http.StatusBadRequest {
		t.Fatalf("a new sheet with a code: status %d; want 400", status)
	}

	// A revision may not take another name.
	status, _ = createRateSheetThroughRoute(t, app, "time@test.com", map[string]any{
		"code": code, "name": "Some Other Name", "revision": 2, "effective_date": "2026-03-01",
	})
	if status != http.StatusBadRequest {
		t.Fatalf("a revision with another name: status %d; want 400", status)
	}
}
