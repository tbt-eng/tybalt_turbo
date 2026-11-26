package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"tybalt/routes"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/tests"
)

// These tests use the fixture draft parequest000001 (client pqpd90fqd5ohjcs,
// created by time@test.com, who lacks the job claim, with invoicing profile
// painvoice000001 and attachment paattachment001) and painvoice000004, the
// client's second profile. painvoice000002 belongs to another client.
// author@soup.com holds the job claim.
const invoicingRequest = "parequest000001"

// setJobsEditingFlag changes only the jobs editing switch, keeping the other
// values of the jobs app_config row.
func setJobsEditingFlag(t *testing.T, app *tests.TestApp, enabled bool) {
	t.Helper()
	config, err := app.FindFirstRecordByData("app_config", "key", "jobs")
	if err != nil {
		t.Fatal(err)
	}
	values := map[string]any{}
	if err = json.Unmarshal([]byte(config.GetString("value")), &values); err != nil {
		t.Fatal(err)
	}
	values["create_edit_absorb"] = enabled
	config.Set("value", values)
	if err = app.Save(config); err != nil {
		t.Fatal(err)
	}
}

// approveInvoicingRequest submits and approves the fixture draft, so it waits
// for job creation.
func approveInvoicingRequest(t *testing.T, app *tests.TestApp) {
	t.Helper()
	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+invoicingRequest+"/submit", "time@test.com", requestActionData(t, app, invoicingRequest), http.StatusOK)
	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+invoicingRequest+"/approve", workflowApprover, requestActionData(t, app, invoicingRequest), http.StatusOK)
}

// setRequestProfile posts the procurement profile change and returns the body.
func setRequestProfile(t *testing.T, app *tests.TestApp, email, profile string, version, status int) string {
	t.Helper()
	token := ""
	if email != "" {
		token = workflowToken(t, app, email)
	}
	response := performJSONRequest(t, app, http.MethodPost, "/api/job_requests/"+invoicingRequest+"/invoicing_profile", token, map[string]any{"version": version, "invoicing_information": profile})
	mustStatus(t, response, status)
	return response.Body.String()
}

// requireApprovalKept checks that the request still waits for creation with
// its saved approval, version and review digest.
func requireApprovalKept(t *testing.T, app *tests.TestApp, before map[string]any) {
	t.Helper()
	r := workflowRecord(t, app, "job_requests", invoicingRequest)
	if r.GetString("status") != "Awaiting creation" || r.GetInt("version") != before["version"] {
		t.Fatalf("request status %s, version %d; want Awaiting creation and version %v", r.GetString("status"), r.GetInt("version"), before["version"])
	}
	a := workflowRecord(t, app, "job_attachments", "paattachment001")
	if a.GetString("status") != "Approved" || a.GetInt("request_version") != r.GetInt("version") {
		t.Fatalf("attachment status %s, request_version %d; want the approval kept", a.GetString("status"), a.GetInt("request_version"))
	}
	if after := requestActionData(t, app, invoicingRequest); after["review_token"] != before["review_token"] {
		t.Fatal("the review digest changed")
	}
}

// createInvoicingJob creates the job and returns its invoicing profile.
func createInvoicingJob(t *testing.T, app *tests.TestApp) string {
	t.Helper()
	result := clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+invoicingRequest+"/create", workflowAuthor, requestActionData(t, app, invoicingRequest), http.StatusOK)
	jobID, _ := result["job_id"].(string)
	return workflowRecord(t, app, "jobs", jobID).GetString("invoicing_information")
}

// A request without an invoicing profile can be submitted and approved, and
// the job is created without one.
func TestJobRequestInvoicingProfileIsOptional(t *testing.T) {
	app := newWorkflowTestApp(t, true)
	r := workflowRecord(t, app, "job_requests", invoicingRequest)
	clientSetupJSON(t, app, http.MethodPatch, "/api/job_requests/"+invoicingRequest, "time@test.com", map[string]any{"version": r.GetInt("version"), "request": map[string]any{"invoicing_information": ""}}, http.StatusOK)
	approveInvoicingRequest(t, app)
	detail := clientSetupJSON(t, app, http.MethodGet, "/api/job_requests/"+invoicingRequest, workflowAuthor, nil, http.StatusOK)
	related := detail["related"].(map[string]any)
	if _, ok := related["billing_contact"]; ok {
		t.Fatal("a request without a profile has a billing contact")
	}
	if profile := createInvoicingJob(t, app); profile != "" {
		t.Fatalf("job profile = %q; want none", profile)
	}
}

// Procurement sets, changes and clears the profile of a request that waits for
// job creation. The request keeps its status, version and approval, and the
// job takes the last profile.
func TestJobRequestProcurementChangesInvoicingProfile(t *testing.T) {
	app := newWorkflowTestApp(t, true)
	approveInvoicingRequest(t, app)
	before := requestActionData(t, app, invoicingRequest)
	version := before["version"].(int)
	for _, profile := range []string{"painvoice000004", "", "painvoice000004"} {
		body := setRequestProfile(t, app, workflowAuthor, profile, version, http.StatusOK)
		if !strings.Contains(body, `"can_set_invoicing":true`) {
			t.Fatalf("detail after the change does not offer it again: %s", body)
		}
		if got := workflowRecord(t, app, "job_requests", invoicingRequest).GetString("invoicing_information"); got != profile {
			t.Fatalf("profile = %q; want %q", got, profile)
		}
		requireApprovalKept(t, app, before)
	}
	if profile := createInvoicingJob(t, app); profile != "painvoice000004" {
		t.Fatalf("job profile = %q; want painvoice000004", profile)
	}
}

// A refused change leaves the profile as it was.
func TestJobRequestProcurementInvoicingProfileRefusals(t *testing.T) {
	app := newWorkflowTestApp(t, true)
	// Not while the request is a draft or awaits approval.
	setRequestProfile(t, app, workflowAuthor, "painvoice000004", 1, http.StatusBadRequest)
	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+invoicingRequest+"/submit", "time@test.com", requestActionData(t, app, invoicingRequest), http.StatusOK)
	setRequestProfile(t, app, workflowAuthor, "painvoice000004", 1, http.StatusBadRequest)
	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+invoicingRequest+"/approve", workflowApprover, requestActionData(t, app, invoicingRequest), http.StatusOK)
	detail := clientSetupJSON(t, app, http.MethodGet, "/api/job_requests/"+invoicingRequest, "time@test.com", nil, http.StatusOK)
	if detail["permissions"].(map[string]any)["can_set_invoicing"] != false {
		t.Fatal("the creator without the job claim is offered the change")
	}
	for _, tc := range []struct {
		name, email, profile string
		version, status      int
		want                 string
	}{
		{"creator without the job claim", "time@test.com", "painvoice000004", 1, http.StatusForbidden, "job claim"},
		{"no claims", "u_no_claims@example.com", "painvoice000004", 1, http.StatusForbidden, ""},
		{"unauthenticated", "", "painvoice000004", 1, http.StatusUnauthorized, ""},
		{"profile of another client", workflowAuthor, "painvoice000002", 1, http.StatusBadRequest, `"invoicing_information":{"code":"client_mismatch"`},
		{"unknown profile", workflowAuthor, "missingprofile0", 1, http.StatusBadRequest, `"invoicing_information":{"code":"invalid_reference"`},
		{"outdated page", workflowAuthor, "painvoice000004", 0, http.StatusConflict, "Reload"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if body := setRequestProfile(t, app, tc.email, tc.profile, tc.version, tc.status); !strings.Contains(body, tc.want) {
				t.Fatalf("refusal %s does not contain %s", body, tc.want)
			}
			if got := workflowRecord(t, app, "job_requests", invoicingRequest).GetString("invoicing_information"); got != "painvoice000001" {
				t.Fatalf("a refused change set the profile to %q", got)
			}
		})
	}
	// With job editing off, the change is refused like other job writes.
	setJobsEditingFlag(t, app, false)
	setRequestProfile(t, app, workflowAuthor, "painvoice000004", 1, http.StatusForbidden)
	setJobsEditingFlag(t, app, true)
	// Test-time change to fixture data: author@soup.com, who holds job, is
	// made inactive.
	execSQL(t, app, "UPDATE admin_profiles SET active = 0 WHERE uid = 'f2j5a8vk006baub'")
	setRequestProfile(t, app, workflowAuthor, "painvoice000004", 1, http.StatusForbidden)
	if got := workflowRecord(t, app, "job_requests", invoicingRequest).GetString("invoicing_information"); got != "painvoice000001" {
		t.Fatalf("a refused change set the profile to %q", got)
	}
}

// An edit to the contents of the profile that a request uses keeps the
// request's approval, before and after approval.
func TestJobRequestProfileEditKeepsApproval(t *testing.T) {
	app := newWorkflowTestApp(t, true)
	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+invoicingRequest+"/submit", "time@test.com", requestActionData(t, app, invoicingRequest), http.StatusOK)
	reviewed := requestActionData(t, app, invoicingRequest)
	editProfile := func(fields map[string]any) {
		clientSetupJSON(t, app, http.MethodPatch, "/api/collections/client_invoicing_information/records/painvoice000001", workflowAuthor, fields, http.StatusOK)
	}
	editProfile(map[string]any{"invoicing_instructions": "Send invoices to the project office.", "fax": "807-555-0100"})
	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+invoicingRequest+"/approve", workflowApprover, reviewed, http.StatusOK)
	before := requestActionData(t, app, invoicingRequest)
	editProfile(map[string]any{"contact": "pacontact000002", "name": "Head office"})
	requireApprovalKept(t, app, before)
	if profile := createInvoicingJob(t, app); profile != "painvoice000001" {
		t.Fatalf("job profile = %q; want painvoice000001", profile)
	}
}

// Replace with… moves an approved request to another profile and keeps its
// approval and version.
func TestJobRequestProfileReplaceKeepsApproval(t *testing.T) {
	app := newWorkflowTestApp(t, true)
	approveInvoicingRequest(t, app)
	before := requestActionData(t, app, invoicingRequest)
	response := performJSONRequest(t, app, http.MethodPost, "/api/clients/pqpd90fqd5ohjcs/invoicing_profiles/painvoice000001/replace", workflowToken(t, app, workflowAuthor), map[string]any{"replacement": "painvoice000004"})
	mustStatus(t, response, http.StatusOK)
	if got := workflowRecord(t, app, "job_requests", invoicingRequest).GetString("invoicing_information"); got != "painvoice000004" {
		t.Fatalf("profile = %q; want painvoice000004", got)
	}
	requireApprovalKept(t, app, before)
	if profile := createInvoicingJob(t, app); profile != "painvoice000004" {
		t.Fatalf("job profile = %q; want painvoice000004", profile)
	}
}

// saveLegacyProfileDigest replaces the saved approval digest with the form
// that held the invoicing profile, computed by the server code, to stand in
// for an approval saved before the profile left the digest. It writes the
// database directly, as such an approval was saved.
func saveLegacyProfileDigest(t *testing.T, app *tests.TestApp) string {
	t.Helper()
	legacy, err := routes.RequestLegacyReviewTokenForTest(app, invoicingRequest, "profile")
	if err != nil {
		t.Fatal(err)
	}
	execSQL(t, app, "UPDATE job_attachments SET review_token = '"+legacy+"' WHERE id = 'paattachment001'")
	return legacy
}

// An approval saved before the profile left the digest holds the profile. The
// upgrade that runs when the server starts gives it the current digest, so a
// later profile edit, change or replace does not make it stale.
func TestJobRequestLegacyApprovalSurvivesProfileChange(t *testing.T) {
	t.Run("upgrade then profile edit", func(t *testing.T) {
		app := newWorkflowTestApp(t, true)
		approveInvoicingRequest(t, app)
		legacy := saveLegacyProfileDigest(t, app)
		// Every test request also runs the upgrade, as the server start does,
		// so the saved digest is read before the next request.
		if err := routes.UpgradeLegacyRequestApprovalsForTest(app); err != nil {
			t.Fatal(err)
		}
		saved := workflowRecord(t, app, "job_attachments", "paattachment001").GetString("review_token")
		current := requestActionData(t, app, invoicingRequest)["review_token"]
		if legacy == current {
			t.Fatal("the legacy digest equals the current digest, so this test proves nothing")
		}
		if saved != current {
			t.Fatal("the upgrade did not save the current digest")
		}
		// Edited in the database, so no request runs before the job creation.
		execSQL(t, app, "UPDATE client_invoicing_information SET invoicing_instructions = 'Changed after the upgrade' WHERE id = 'painvoice000001'")
		createInvoicingJob(t, app)
	})
	t.Run("procurement change", func(t *testing.T) {
		app := newWorkflowTestApp(t, true)
		approveInvoicingRequest(t, app)
		saveLegacyProfileDigest(t, app)
		setRequestProfile(t, app, workflowAuthor, "painvoice000004", 1, http.StatusOK)
		if profile := createInvoicingJob(t, app); profile != "painvoice000004" {
			t.Fatalf("job profile = %q; want painvoice000004", profile)
		}
	})
	t.Run("stale approval stays stale", func(t *testing.T) {
		app := newWorkflowTestApp(t, true)
		approveInvoicingRequest(t, app)
		legacy := saveLegacyProfileDigest(t, app)
		// A change the approval covers, made after it, in the database.
		if _, err := app.DB().NewQuery("UPDATE job_requests SET description = 'Changed after approval' WHERE id = {:id}").Bind(dbx.Params{"id": invoicingRequest}).Execute(); err != nil {
			t.Fatal(err)
		}
		if err := routes.UpgradeLegacyRequestApprovalsForTest(app); err != nil {
			t.Fatal(err)
		}
		if got := workflowRecord(t, app, "job_attachments", "paattachment001").GetString("review_token"); got != legacy {
			t.Fatal("the upgrade replaced a digest that no longer matches")
		}
		clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+invoicingRequest+"/create", workflowAuthor, requestActionData(t, app, invoicingRequest), http.StatusConflict)
	})
}
