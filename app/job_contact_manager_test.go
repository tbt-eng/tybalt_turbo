package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"tybalt/internal/testutils"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

// Imported jobs can lack a contact or manager (issue #167). validateJob
// requires both on creates and full edits, but a status-only change, the close
// action, and the close action's proposal auto-award must still work.

const (
	// An Active project with incomplete billing data.
	contactManagerProjectID = "test_job_w_rs"
	// An In Progress, time and materials proposal, so a status-only change to
	// Submitted needs no proposal value.
	contactManagerProposalID = "pcproposal00001"
)

// setupJobWithBlankField clears field on jobID, as imported data can.
func setupJobWithBlankField(factory func(testing.TB) *tests.TestApp, jobID, field string) func(testing.TB) *tests.TestApp {
	return setupJobWithFields(factory, jobID, dbx.Params{field: ""})
}

// setupJobWithFields writes values to jobID directly, bypassing validation.
func setupJobWithFields(factory func(testing.TB) *tests.TestApp, jobID string, values dbx.Params) func(testing.TB) *tests.TestApp {
	return func(t testing.TB) *tests.TestApp {
		app := factory(t)
		if _, err := app.DB().Update("jobs", values, dbx.HashExp{"id": jobID}).Execute(); err != nil {
			t.Fatalf("failed to update job %s: %v", jobID, err)
		}
		return app
	}
}

func TestJobsStatusOnly_AllowsMissingContactOrManager(t *testing.T) {
	recordToken, err := testutils.GenerateRecordToken("users", "author@soup.com")
	if err != nil {
		t.Fatal(err)
	}

	for _, field := range []string{"contact", "manager"} {
		projectFactory := setupJobWithBlankField(testutils.SetupTestApp, contactManagerProjectID, field)
		scenarios := []tests.ApiScenario{
			{
				Name:            "close action closes a project with no " + field,
				Method:          http.MethodPost,
				URL:             "/api/jobs/" + contactManagerProjectID + "/close",
				Headers:         map[string]string{"Authorization": recordToken},
				ExpectedStatus:  http.StatusOK,
				ExpectedContent: []string{`"status":"Closed"`},
				TestAppFactory:  projectFactory,
			},
			{
				Name:            "status-only close of a project with no " + field,
				Method:          http.MethodPatch,
				URL:             "/api/collections/jobs/records/" + contactManagerProjectID,
				Body:            strings.NewReader(`{"status":"Closed"}`),
				Headers:         map[string]string{"Authorization": recordToken},
				ExpectedStatus:  http.StatusOK,
				ExpectedContent: []string{`"status":"Closed"`},
				TestAppFactory:  projectFactory,
			},
			{
				Name:            "status-only cancel of a project with no " + field,
				Method:          http.MethodPatch,
				URL:             "/api/collections/jobs/records/" + contactManagerProjectID,
				Body:            strings.NewReader(`{"status":"Cancelled"}`),
				Headers:         map[string]string{"Authorization": recordToken},
				ExpectedStatus:  http.StatusOK,
				ExpectedContent: []string{`"status":"Cancelled"`},
				TestAppFactory:  projectFactory,
			},
			{
				// The fixture proposal carries a rate sheet, which cleanJob
				// clears from proposals. Clear it first so the PATCH is a
				// status-only update (see the caveat in processJobCoreWithMode).
				Name:            "status-only change of a proposal with no " + field,
				Method:          http.MethodPatch,
				URL:             "/api/collections/jobs/records/" + contactManagerProposalID,
				Body:            strings.NewReader(`{"status":"Submitted"}`),
				Headers:         map[string]string{"Authorization": recordToken},
				ExpectedStatus:  http.StatusOK,
				ExpectedContent: []string{`"status":"Submitted"`},
				TestAppFactory:  setupJobWithFields(testutils.SetupTestApp, contactManagerProposalID, dbx.Params{field: "", "rate_sheet": ""}),
			},
			{
				Name:            "set-status cancels a proposal with no " + field,
				Method:          http.MethodPost,
				URL:             "/api/jobs/" + contactManagerProposalID + "/set-status",
				Body:            strings.NewReader(`{"status":"Cancelled","comment":"No longer pursued"}`),
				Headers:         map[string]string{"Authorization": recordToken},
				ExpectedStatus:  http.StatusOK,
				ExpectedContent: []string{`"status":"Cancelled"`},
				TestAppFactory:  setupJobWithBlankField(testutils.SetupTestApp, contactManagerProposalID, field),
				AfterTestFunc: func(tb testing.TB, app *tests.TestApp, _ *http.Response) {
					proposal, err := app.FindRecordById("jobs", contactManagerProposalID)
					if err != nil {
						tb.Fatal(err)
					}
					if got := proposal.GetString("status"); got != "Cancelled" {
						tb.Fatalf("expected proposal status Cancelled, got %q", got)
					}
				},
			},
		}
		for _, scenario := range scenarios {
			scenario.Test(t)
		}
	}
}

// Closing a project auto-awards its imported proposal with a plain Save, which
// the schema requirement used to reject when the proposal lacked a contact.
func TestCloseJob_AutoAwardsProposalMissingContactOrManager(t *testing.T) {
	recordToken, err := testutils.GenerateRecordToken("users", "author@soup.com")
	if err != nil {
		t.Fatal(err)
	}

	const projectID = "fcprojimpauto01"
	const proposalID = "fcpropimpinprog"

	for _, field := range []string{"contact", "manager"} {
		scenario := tests.ApiScenario{
			Name:            "close auto-awards an imported proposal with no " + field,
			Method:          http.MethodPost,
			URL:             "/api/jobs/" + projectID + "/close",
			Headers:         map[string]string{"Authorization": recordToken},
			ExpectedStatus:  http.StatusOK,
			ExpectedContent: []string{`"status":"Closed"`, `"auto_awarded":true`},
			TestAppFactory:  setupJobWithBlankField(testutils.SetupTestApp, proposalID, field),
			AfterTestFunc: func(tb testing.TB, app *tests.TestApp, _ *http.Response) {
				proposal, err := app.FindRecordById("jobs", proposalID)
				if err != nil {
					tb.Fatal(err)
				}
				if got := proposal.GetString("status"); got != "Awarded" {
					tb.Fatalf("expected proposal status Awarded, got %q", got)
				}
			},
		}
		scenario.Test(t)
	}
}

func TestJobsFullEdit_RequiresContactAndManager(t *testing.T) {
	recordToken, err := testutils.GenerateRecordToken("users", "author@soup.com")
	if err != nil {
		t.Fatal(err)
	}

	// A project that passes full validation once billing fixtures are added.
	const projectID = "cjf0kt0defhq480"

	for _, field := range []string{"contact", "manager"} {
		createJob := map[string]any{
			"description":           "Test job via API",
			"client":                "ee3xvodl583b61o",
			"contact":               "235g6k01xx3sdjk",
			"invoicing_information": "painvoice000003",
			"manager":               "f2j5a8vk006baub",
			"branch":                "80875lm27v8wgi4",
			"location":              "87Q8H976+2M",
			"project_award_date":    "2025-02-01",
			"rate_sheet":            "c41ofep525bcacj",
		}
		delete(createJob, field)
		createBody, err := json.Marshal(map[string]any{
			"job":         createJob,
			"allocations": []map[string]any{{"division": "fy4i9poneukvq9u", "hours": 10}},
		})
		if err != nil {
			t.Fatal(err)
		}

		requiredError := `"` + field + `":{"code":"required"`
		scenarios := []tests.ApiScenario{
			{
				Name:            "full project edit clearing " + field + " is rejected",
				Method:          http.MethodPatch,
				URL:             "/api/collections/jobs/records/" + projectID,
				Body:            strings.NewReader(`{"description":"Edited","` + field + `":""}`),
				Headers:         map[string]string{"Authorization": recordToken},
				ExpectedStatus:  http.StatusBadRequest,
				ExpectedContent: []string{requiredError},
				TestAppFactory:  setupJobsWithBillingFixtures,
			},
			{
				Name:            "full edit of a project already missing " + field + " is rejected",
				Method:          http.MethodPatch,
				URL:             "/api/collections/jobs/records/" + projectID,
				Body:            strings.NewReader(`{"description":"Edited"}`),
				Headers:         map[string]string{"Authorization": recordToken},
				ExpectedStatus:  http.StatusBadRequest,
				ExpectedContent: []string{requiredError},
				TestAppFactory:  setupJobWithBlankField(setupJobsWithBillingFixtures, projectID, field),
			},
			{
				Name:            "full proposal edit clearing " + field + " is rejected",
				Method:          http.MethodPatch,
				URL:             "/api/collections/jobs/records/" + contactManagerProposalID,
				Body:            strings.NewReader(`{"description":"Edited","` + field + `":""}`),
				Headers:         map[string]string{"Authorization": recordToken},
				ExpectedStatus:  http.StatusBadRequest,
				ExpectedContent: []string{requiredError},
				TestAppFactory:  testutils.SetupTestApp,
			},
			{
				Name:            "create without " + field + " is rejected",
				Method:          http.MethodPost,
				URL:             "/api/jobs",
				Body:            strings.NewReader(string(createBody)),
				Headers:         map[string]string{"Authorization": recordToken},
				ExpectedStatus:  http.StatusBadRequest,
				ExpectedContent: []string{requiredError},
				TestAppFactory:  setupJobsWithBillingFixtures,
			},
		}
		for _, scenario := range scenarios {
			scenario.Test(t)
		}
	}
}

// jobs.manager is optional in the schema, so PocketBase would clear it on a
// proposal when the manager is deleted (jobGuard already refuses that write for
// a project). The users delete hook refuses the delete instead.
func TestUsersDelete_RefusedWhileManagingJobs(t *testing.T) {
	app := testutils.SetupTestApp(t)
	defer app.Cleanup()

	users, err := app.FindCollectionByNameOrId("users")
	if err != nil {
		t.Fatal(err)
	}
	user := core.NewRecord(users)
	user.Set("email", "manager-delete-guard@example.com")
	user.Set("username", "managerdeleteguard")
	user.SetPassword("password123456")
	if err := app.SaveNoValidate(user); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	if _, err := app.DB().Update("jobs", dbx.Params{"manager": user.Id}, dbx.HashExp{"id": contactManagerProposalID}).Execute(); err != nil {
		t.Fatalf("failed to assign manager: %v", err)
	}

	err = app.Delete(user)
	if err == nil || !strings.Contains(err.Error(), "manager of one or more jobs") {
		t.Fatalf("expected delete to be refused for a proposal manager, got %v", err)
	}
	proposal, err := app.FindRecordById("jobs", contactManagerProposalID)
	if err != nil {
		t.Fatal(err)
	}
	if proposal.GetString("manager") != user.Id {
		t.Fatalf("expected manager to be kept, got %q", proposal.GetString("manager"))
	}
}
