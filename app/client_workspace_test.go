package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"tybalt/internal/testutils"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/tests"
)

const workspaceClient = "cwclient0000001"

func newClientWorkspaceTestApp(t *testing.T) *tests.TestApp {
	t.Helper()
	app := testutils.SetupTestApp(t)
	t.Cleanup(app.Cleanup)
	// These CSV rows test historical use in isolation. Loading them here keeps
	// them out of the project queues and totals used by unrelated tests.
	// Insert the stored rows directly, as the main CSV seed loader does.
	for _, table := range []string{"jobs", "job_requests"} {
		f, err := os.Open("test_seed_data/client_workspace/" + table + ".csv")
		if err != nil {
			t.Fatal(err)
		}
		rows, err := csv.NewReader(f).ReadAll()
		f.Close()
		if err != nil {
			t.Fatal(err)
		}
		for _, values := range rows[1:] {
			row := dbx.Params{}
			for i, key := range rows[0] {
				row[key] = values[i]
			}
			if _, err := app.DB().Insert(table, row).Execute(); err != nil {
				t.Fatal(err)
			}
		}
	}
	return app
}

func TestClientWorkspaceDetailsAndUsage(t *testing.T) {
	app := newClientWorkspaceTestApp(t)
	for _, email := range []string{"author@soup.com", "time@test.com"} {
		t.Run(email, func(t *testing.T) {
			result := clientSetupJSON(t, app, http.MethodGet, "/api/clients/"+workspaceClient, email, nil, http.StatusOK)
			for field, want := range map[string]any{
				"id": workspaceClient, "name": "Workspace Test Client", "address": "100 Test Street",
				"city": "Thunder Bay", "province_state": "ON", "postal_code": "P7A 1A1", "country": "Canada", "phone": "807-555-0100",
				"business_development_lead": "f2j5a8vk006baub", "referencing_jobs_count": float64(2), "outstanding_balance": float64(46),
			} {
				if result[field] != want {
					t.Errorf("%s = %#v; want %#v", field, result[field], want)
				}
			}
			contacts := result["contacts"].([]any)
			if len(contacts) != 3 {
				t.Fatalf("contacts = %d; want only this client's 3 contacts", len(contacts))
			}
			first := contacts[0].(map[string]any)
			for field, want := range map[string]any{
				"id": "cwcontact000001", "client": workspaceClient, "given_name": "Alex", "surname": "Workspace",
				"email": "alex.workspace@example.com", "phone": "807-555-0101", "address": "101 Test Street", "city": "Dryden",
				"province_state": "ON", "postal_code": "P8N 1A1", "country": "Canada",
				"job_count": float64(2), "request_count": float64(2), "profile_count": float64(2),
			} {
				if first[field] != want {
					t.Errorf("contact %s = %#v; want %#v", field, first[field], want)
				}
			}
			// Alex is used directly and through invoicing on the first record,
			// and through invoicing alone on the second. Neither is counted twice.
			second := contacts[1].(map[string]any)
			if second["job_count"] != float64(1) || second["request_count"] != float64(1) || second["profile_count"] != float64(1) {
				t.Fatalf("wrong direct contact counts: %#v", second)
			}
			third := contacts[2].(map[string]any)
			if third["job_count"] != float64(0) || third["request_count"] != float64(0) || third["profile_count"] != float64(0) {
				t.Fatalf("unused contact has references: %#v", third)
			}
			profiles := result["invoicing_profiles"].([]any)
			if len(profiles) != 3 {
				t.Fatalf("profiles = %d; want only this client's 3 setups", len(profiles))
			}
			for _, value := range profiles {
				profile := value.(map[string]any)
				if profile["id"] == "cwinvoice000001" {
					for field, want := range map[string]any{"name": "Capital projects", "client": workspaceClient, "contact": "cwcontact000001", "fax": "807-555-0199", "invoicing_instructions": "Include the purchase order.", "job_count": float64(2), "request_count": float64(2)} {
						if profile[field] != want {
							t.Errorf("profile %s = %#v; want %#v", field, profile[field], want)
						}
					}
				} else if profile["job_count"] != float64(0) || profile["request_count"] != float64(0) {
					t.Errorf("unused setup has references: %#v", profile)
				}
			}
		})
	}
}

func TestClientWorkspaceJobDetailsIncludesInvoicing(t *testing.T) {
	app := newClientWorkspaceTestApp(t)
	for _, email := range []string{"author@soup.com", "time@test.com"} {
		t.Run(email, func(t *testing.T) {
			result := clientSetupJSON(t, app, http.MethodGet, "/api/jobs/cwjob0000000001/details", email, nil, http.StatusOK)
			for field, want := range map[string]string{
				"invoicing_information":  "cwinvoice000001",
				"invoicing_profile_name": "Capital projects",
				"invoice_contact_name":   "Alex Workspace",
				"invoice_contact_email":  "alex.workspace@example.com",
				"invoice_instructions":   "Include the purchase order.",
			} {
				if result[field] != want {
					t.Errorf("%s = %#v; want %q", field, result[field], want)
				}
			}
		})
	}
}

func TestClientWorkspaceJobDetailsRequiresSignIn(t *testing.T) {
	app := newClientWorkspaceTestApp(t)
	clientSetupJSON(t, app, http.MethodGet, "/api/jobs/cwjob0000000001/details", "", nil, http.StatusUnauthorized)
}

func TestClientWorkspaceEmptyMissingAndUnauthenticated(t *testing.T) {
	app := newClientWorkspaceTestApp(t)
	empty := clientSetupJSON(t, app, http.MethodGet, "/api/clients/cwclient0000002", "author@soup.com", nil, http.StatusOK)
	for _, key := range []string{"contacts", "invoicing_profiles"} {
		items, ok := empty[key].([]any)
		if !ok || len(items) != 0 {
			t.Errorf("%s must be an empty array, got %#v", key, empty[key])
		}
	}
	clientSetupJSON(t, app, http.MethodGet, "/api/clients/missingclient00", "author@soup.com", nil, http.StatusNotFound)
	clientSetupJSON(t, app, http.MethodGet, "/api/clients/"+workspaceClient, "", nil, http.StatusUnauthorized)
}

// BN Rail's seeded jobs have committed time in CB and CI. Its E entries are
// not on a timesheet, and some CB entries are on uncommitted timesheets, so
// neither counts.
func TestClientWorkspaceTimeDivisions(t *testing.T) {
	app := newClientWorkspaceTestApp(t)
	divisions := func(client string) string {
		auth, err := app.FindAuthRecordByEmail("users", "author@soup.com")
		if err != nil {
			t.Fatal(err)
		}
		token, err := auth.NewAuthToken()
		if err != nil {
			t.Fatal(err)
		}
		res := performJSONRequest(t, app, http.MethodGet, "/api/clients/"+client+"/time_divisions", token, nil)
		mustStatus(t, res, http.StatusOK)
		var items []map[string]any
		if err := json.Unmarshal(res.Body.Bytes(), &items); err != nil || items == nil {
			t.Fatalf("time divisions must be an array, got %s", res.Body.String())
		}
		got := []string{}
		for _, item := range items {
			got = append(got, fmt.Sprintf("%s=%v", item["code"], item["hours"]))
		}
		return strings.Join(got, ",")
	}
	if got, want := divisions("eldtxi3i4h00k8r"), "CB=2016,CI=3"; got != want {
		t.Errorf("time divisions = %s; want %s", got, want)
	}
	if got := divisions("cwclient0000002"); got != "" {
		t.Errorf("client without time has divisions %s", got)
	}
	// A client that only owns the jobs does not get their divisions.
	if _, err := app.DB().NewQuery("UPDATE jobs SET job_owner = 'cwclient0000002' WHERE client = 'eldtxi3i4h00k8r'").Execute(); err != nil {
		t.Fatal(err)
	}
	if got := divisions("cwclient0000002"); got != "" {
		t.Errorf("job owner time divisions = %s; want none", got)
	}
}

func TestClientWorkspaceInvoicingNamesAndPermissions(t *testing.T) {
	app := newClientWorkspaceTestApp(t)
	path := "/api/collections/client_invoicing_information/records/"
	unnamed := clientSetupJSON(t, app, http.MethodPost, strings.TrimSuffix(path, "/"), "author@soup.com", map[string]any{"client": workspaceClient, "contact": "cwcontact000002"}, http.StatusOK)
	if unnamed["name"] != "" {
		t.Fatalf("name should be optional: %#v", unnamed)
	}
	id := unnamed["id"].(string)
	for _, name := range []string{"Field services", ""} {
		updated := clientSetupJSON(t, app, http.MethodPatch, path+id, "author@soup.com", map[string]any{"name": name}, http.StatusOK)
		if updated["name"] != name {
			t.Fatalf("name was not saved: %#v", updated)
		}
	}
	clientSetupJSON(t, app, http.MethodPatch, path+id, "author@soup.com", map[string]any{"name": strings.Repeat("x", 121)}, http.StatusBadRequest)
	clientSetupJSON(t, app, http.MethodPost, strings.TrimSuffix(path, "/"), "author@soup.com", map[string]any{"client": workspaceClient, "contact": "184qbpajciwrrp2"}, http.StatusBadRequest)
	for _, collection := range []string{"clients", "client_contacts", "client_invoicing_information"} {
		id := map[string]string{"clients": workspaceClient, "client_contacts": "cwcontact000001", "client_invoicing_information": "cwinvoice000001"}[collection]
		base := "/api/collections/" + collection + "/records/"
		validCreate := map[string]map[string]any{
			"clients":                      {"name": "Unauthorized workspace client", "business_development_lead": "f2j5a8vk006baub"},
			"client_contacts":              {"client": workspaceClient, "given_name": "Unauthorized", "surname": "Contact", "email": "unauthorized@example.com"},
			"client_invoicing_information": {"client": workspaceClient, "contact": "cwcontact000001", "name": "Unauthorized setup"},
		}[collection]
		// PocketBase returns 400 when a create rule is not satisfied, and 404
		// when the same rule prevents selecting a record to update or delete.
		createStatus := http.StatusBadRequest
		// Active time holders can create invoicing profiles from existing contacts.
		if collection == "client_invoicing_information" {
			createStatus = http.StatusOK
		}
		clientSetupJSON(t, app, http.MethodPost, strings.TrimSuffix(base, "/"), "time@test.com", validCreate, createStatus)
		clientSetupJSON(t, app, http.MethodPatch, base+id, "time@test.com", map[string]any{}, http.StatusNotFound)
		clientSetupJSON(t, app, http.MethodDelete, base+id, "time@test.com", nil, http.StatusNotFound)
	}
	clientSetupJSON(t, app, http.MethodDelete, path+"cwinvoice000001", "author@soup.com", nil, http.StatusNotFound)
	clientSetupJSON(t, app, http.MethodDelete, "/api/collections/client_contacts/records/cwcontact000001", "author@soup.com", nil, http.StatusNotFound)
	clientSetupJSON(t, app, http.MethodDelete, path+id, "author@soup.com", nil, http.StatusNoContent)
}

func clientSetupJSON(t *testing.T, app *tests.TestApp, method, path, email string, data any, status int) map[string]any {
	t.Helper()
	token := ""
	if email != "" {
		auth, err := app.FindAuthRecordByEmail("users", email)
		if err != nil {
			t.Fatal(err)
		}
		token, err = auth.NewAuthToken()
		if err != nil {
			t.Fatal(err)
		}
	}
	response := performJSONRequest(t, app, method, path, token, data)
	mustStatus(t, response, status)
	result := map[string]any{}
	if status == http.StatusOK {
		if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
			t.Fatal(err)
		}
	}
	return result
}

// Client details report which contacts are marked as having no email address.
func TestClientWorkspaceReportsContactsWithoutEmail(t *testing.T) {
	app := newClientWorkspaceTestApp(t)
	contact, err := app.FindRecordById("client_contacts", "cwcontact000003")
	if err != nil {
		t.Fatal(err)
	}
	contact.Set("email", "")
	contact.Set("no_email", true)
	if err := app.Save(contact); err != nil {
		t.Fatal(err)
	}
	result := clientSetupJSON(t, app, http.MethodGet, "/api/clients/"+workspaceClient, "author@soup.com", nil, http.StatusOK)
	marks := map[string]any{}
	for _, value := range result["contacts"].([]any) {
		row := value.(map[string]any)
		marks[row["id"].(string)] = row["no_email"]
	}
	if marks["cwcontact000003"] != true || marks["cwcontact000001"] != false {
		t.Fatalf("contact no-email marks = %#v", marks)
	}

	// Entering an email address clears the mark.
	updated := clientSetupJSON(t, app, http.MethodPatch, "/api/collections/client_contacts/records/cwcontact000003", "author@soup.com", map[string]any{"email": "found@example.com"}, http.StatusOK)
	if updated["no_email"] != false {
		t.Fatalf("a contact with an email address kept the no-email mark: %#v", updated)
	}
}

// The usage lists name the jobs and requests behind the counts in the client
// details. They use the rows of test_seed_data/client_workspace: Alex is the
// project contact and invoice recipient of the first job and request, and the
// invoice recipient alone of the second; Blair is the project contact of the
// second.
func TestClientWorkspaceUsageLists(t *testing.T) {
	app := newClientWorkspaceTestApp(t)
	base := "/api/clients/" + workspaceClient
	type use struct {
		id                               string
		projectContact, invoiceRecipient bool
	}
	check := func(t *testing.T, kind string, rows any, want []use) {
		t.Helper()
		list := rows.([]any)
		if len(list) != len(want) {
			t.Fatalf("%s = %d; want %d: %#v", kind, len(list), len(want), list)
		}
		for i, value := range list {
			row := value.(map[string]any)
			if row["id"] != want[i].id || row["project_contact"] != want[i].projectContact || row["invoice_recipient"] != want[i].invoiceRecipient {
				t.Errorf("%s[%d] = %#v; want %#v", kind, i, row, want[i])
			}
		}
	}

	for _, tc := range []struct {
		name, path     string
		jobs, requests []use
	}{
		{"contact in both roles", "/contacts/cwcontact000001/usage",
			[]use{{"cwjob0000000002", false, true}, {"cwjob0000000001", true, true}},
			[]use{{"cwrequest000001", true, true}, {"cwrequest000002", false, true}}},
		{"project contact only", "/contacts/cwcontact000002/usage",
			[]use{{"cwjob0000000002", true, false}},
			[]use{{"cwrequest000002", true, false}}},
		{"unused contact", "/contacts/cwcontact000003/usage", nil, nil},
		{"profile", "/invoicing_profiles/cwinvoice000001/usage",
			[]use{{"cwjob0000000002", false, true}, {"cwjob0000000001", false, true}},
			[]use{{"cwrequest000001", false, true}, {"cwrequest000002", false, true}}},
		{"unused profile", "/invoicing_profiles/cwinvoice000002/usage", nil, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result := clientSetupJSON(t, app, http.MethodGet, base+tc.path, "author@soup.com", nil, http.StatusOK)
			check(t, "jobs", result["jobs"], tc.jobs)
			check(t, "requests", result["requests"], tc.requests)
			if result["hidden_requests"] != float64(0) {
				t.Errorf("hidden_requests = %v; want 0 for procurement", result["hidden_requests"])
			}
		})
	}
	first := clientSetupJSON(t, app, http.MethodGet, base+"/contacts/cwcontact000001/usage", "author@soup.com", nil, http.StatusOK)["jobs"].([]any)[0].(map[string]any)
	if first["number"] != "88-9992" || first["description"] != "Client workspace usage fixture" || first["status"] != "Closed" {
		t.Errorf("job fields = %#v", first)
	}

	// Requests follow the job_requests read rule; jobs are open to everyone.
	for email, visible := range map[string]int{
		"time@test.com":        2, // the creator
		"tier2@poapprover.com": 2, // the approver, with the project_authorization claim
		"noclaims@example.com": 0,
	} {
		result := clientSetupJSON(t, app, http.MethodGet, base+"/invoicing_profiles/cwinvoice000001/usage", email, nil, http.StatusOK)
		if got := len(result["requests"].([]any)); got != visible || result["hidden_requests"] != float64(2-visible) {
			t.Errorf("%s sees %d requests, %v hidden; want %d and %d", email, got, result["hidden_requests"], visible, 2-visible)
		}
		if len(result["jobs"].([]any)) != 2 {
			t.Errorf("%s sees %d jobs; want 2", email, len(result["jobs"].([]any)))
		}
	}

	// A record of another client, a contact id on the profile route, a missing
	// id and a signed-out caller get no list.
	clientSetupJSON(t, app, http.MethodGet, "/api/clients/cwclient0000002/contacts/cwcontact000001/usage", "author@soup.com", nil, http.StatusNotFound)
	clientSetupJSON(t, app, http.MethodGet, base+"/invoicing_profiles/cwcontact000001/usage", "author@soup.com", nil, http.StatusNotFound)
	clientSetupJSON(t, app, http.MethodGet, base+"/contacts/missingcontact0/usage", "author@soup.com", nil, http.StatusNotFound)
	clientSetupJSON(t, app, http.MethodGet, base+"/contacts/cwcontact000001/usage", "", nil, http.StatusUnauthorized)
}

const replaceProfilePath = "/api/clients/" + workspaceClient + "/invoicing_profiles/"

// replaceProfile posts a replace request and returns the response body.
func replaceProfile(t *testing.T, app *tests.TestApp, profile, email string, body any, status int) string {
	t.Helper()
	token := ""
	if email != "" {
		auth, err := app.FindAuthRecordByEmail("users", email)
		if err != nil {
			t.Fatal(err)
		}
		if token, err = auth.NewAuthToken(); err != nil {
			t.Fatal(err)
		}
	}
	response := performJSONRequest(t, app, http.MethodPost, replaceProfilePath+profile+"/replace", token, body)
	mustStatus(t, response, status)
	return response.Body.String()
}

// profileUses returns the ids of the jobs and requests that use a profile.
func profileUses(t *testing.T, app *tests.TestApp, table, profile string) []string {
	t.Helper()
	ids := []string{}
	if err := app.DB().Select("id").From(table).Where(dbx.HashExp{"invoicing_information": profile}).OrderBy("id").Column(&ids); err != nil {
		t.Fatal(err)
	}
	return ids
}

// Replacing cwinvoice000001 with cwinvoice000002 moves both Closed jobs and
// both requests of test_seed_data/client_workspace, and deletes the old profile.
func TestClientWorkspaceReplaceInvoicingProfile(t *testing.T) {
	app := newClientWorkspaceTestApp(t)
	// Test-time changes to fixture data: cwrequest000002 is put in Awaiting
	// approval, so a request past Draft moves too, and cwjob0000000001 is
	// marked imported, to show that a moved job is marked for writeback.
	execSQL(t, app,
		"UPDATE job_requests SET status = 'Awaiting approval' WHERE id = 'cwrequest000002'",
		"UPDATE jobs SET _imported = 1 WHERE id = 'cwjob0000000001'",
	)

	body := replaceProfile(t, app, "cwinvoice000001", "author@soup.com", map[string]any{"replacement": "cwinvoice000002"}, http.StatusOK)
	result := map[string]any{}
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		t.Fatal(err)
	}
	if result["jobs"] != float64(2) || result["requests"] != float64(2) {
		t.Fatalf("moved counts = %#v; want 2 jobs and 2 requests", result)
	}
	if _, err := app.FindRecordById("client_invoicing_information", "cwinvoice000001"); err == nil {
		t.Fatal("the replaced profile was not deleted")
	}
	if got := profileUses(t, app, "jobs", "cwinvoice000002"); fmt.Sprint(got) != "[cwjob0000000001 cwjob0000000002]" {
		t.Errorf("jobs on the replacement = %v", got)
	}
	if got := profileUses(t, app, "job_requests", "cwinvoice000002"); fmt.Sprint(got) != "[cwrequest000001 cwrequest000002]" {
		t.Errorf("requests on the replacement = %v", got)
	}
	for _, id := range []string{"cwjob0000000001", "cwjob0000000002"} {
		job, err := app.FindRecordById("jobs", id)
		if err != nil {
			t.Fatal(err)
		}
		if job.GetBool("_imported") || job.GetString("status") != "Closed" || !job.GetDateTime("updated").Time().After(job.GetDateTime("created").Time()) {
			t.Errorf("job %s: _imported %v, status %s, updated %s; want marked for writeback, still Closed, and updated", id, job.GetBool("_imported"), job.GetString("status"), job.GetString("updated"))
		}
	}
	// The profile is not part of an approval, so a submitted request keeps
	// its version. A draft gets a new one, so an open edit form reloads.
	for id, want := range map[string]struct {
		status  string
		version int
	}{"cwrequest000001": {"Draft", 2}, "cwrequest000002": {"Awaiting approval", 1}} {
		request, err := app.FindRecordById("job_requests", id)
		if err != nil {
			t.Fatal(err)
		}
		if request.GetInt("version") != want.version || request.GetString("status") != want.status {
			t.Errorf("request %s: version %d, status %s; want version %d and %s", id, request.GetInt("version"), request.GetString("status"), want.version, want.status)
		}
	}
	// The client details count the moved records on the replacement.
	details := clientSetupJSON(t, app, http.MethodGet, "/api/clients/"+workspaceClient, "author@soup.com", nil, http.StatusOK)
	profiles := details["invoicing_profiles"].([]any)
	if len(profiles) != 2 {
		t.Fatalf("profiles = %d; want 2 after the replace", len(profiles))
	}
	for _, value := range profiles {
		profile := value.(map[string]any)
		if profile["id"] == "cwinvoice000002" && (profile["job_count"] != float64(2) || profile["request_count"] != float64(2)) {
			t.Errorf("replacement counts = %#v", profile)
		}
	}

	// An unused profile can be replaced too; nothing moves.
	body = replaceProfile(t, app, "cwinvoice000003", "author@soup.com", map[string]any{"replacement": "cwinvoice000002"}, http.StatusOK)
	if !strings.Contains(body, `"jobs":0`) || !strings.Contains(body, `"requests":0`) {
		t.Errorf("replace of an unused profile = %s", body)
	}
}

// A refused replace changes nothing.
func TestClientWorkspaceReplaceInvoicingProfileRejections(t *testing.T) {
	app := newClientWorkspaceTestApp(t)
	for _, tc := range []struct {
		name, profile string
		body          any
		status        int
		want          string
	}{
		{"no replacement", "cwinvoice000001", map[string]any{}, http.StatusBadRequest, `"replacement":{"code":"required"`},
		{"same profile", "cwinvoice000001", map[string]any{"replacement": "cwinvoice000001"}, http.StatusBadRequest, `"replacement":{"code":"unchanged"`},
		{"profile of another client", "cwinvoice000001", map[string]any{"replacement": "painvoice000001"}, http.StatusBadRequest, `"replacement":{"code":"client_mismatch"`},
		{"missing replacement", "cwinvoice000001", map[string]any{"replacement": "missingprofile0"}, http.StatusBadRequest, `"replacement":{"code":"invalid_reference"`},
		{"missing profile", "missingprofile0", map[string]any{"replacement": "cwinvoice000002"}, http.StatusNotFound, "not found for this client"},
		{"profile of another client in the path", "painvoice000001", map[string]any{"replacement": "cwinvoice000002"}, http.StatusNotFound, "not found for this client"},
		{"contact id as the profile", "cwcontact000001", map[string]any{"replacement": "cwinvoice000002"}, http.StatusNotFound, "not found for this client"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if body := replaceProfile(t, app, tc.profile, "author@soup.com", tc.body, tc.status); !strings.Contains(body, tc.want) {
				t.Errorf("response %s does not contain %s", body, tc.want)
			}
		})
	}

	// Only holders of job may replace a profile.
	replacement := map[string]any{"replacement": "cwinvoice000002"}
	replaceProfile(t, app, "cwinvoice000001", "time@test.com", replacement, http.StatusForbidden)
	replaceProfile(t, app, "cwinvoice000001", "", replacement, http.StatusUnauthorized)

	if got := profileUses(t, app, "jobs", "cwinvoice000001"); len(got) != 2 {
		t.Errorf("a refused replace moved jobs: %v remain", got)
	}
	if got := profileUses(t, app, "job_requests", "cwinvoice000001"); len(got) != 2 {
		t.Errorf("a refused replace moved requests: %v remain", got)
	}
	if _, err := app.FindRecordById("client_invoicing_information", "cwinvoice000001"); err != nil {
		t.Errorf("a refused replace deleted the profile: %v", err)
	}
}

// While a contact merge that moved a job of the profile can still be undone,
// the profile can't be replaced: undo would restore the job's references.
func TestClientWorkspaceReplaceInvoicingProfileWaitsForMerge(t *testing.T) {
	app := newClientWorkspaceTestApp(t)
	// Test-time change to fixture data: book@keeper.com merges Blair
	// (cwcontact000002) into Alex (cwcontact000001), which moves
	// cwjob0000000002, cwrequest000002 and cwinvoice000003 to Alex.
	auth, err := app.FindAuthRecordByEmail("users", "book@keeper.com")
	if err != nil {
		t.Fatal(err)
	}
	token, err := auth.NewAuthToken()
	if err != nil {
		t.Fatal(err)
	}
	mustStatus(t, performJSONRequest(t, app, http.MethodPost, "/api/client_contacts/cwcontact000001/absorb", token, map[string]any{"ids_to_absorb": []string{"cwcontact000002"}}), http.StatusOK)

	body := replaceProfile(t, app, "cwinvoice000001", "author@soup.com", map[string]any{"replacement": "cwinvoice000002"}, http.StatusConflict)
	if !strings.Contains(body, "waiting to be committed or undone") {
		t.Errorf("response %s does not name the pending merge", body)
	}
	// The merge moved cwinvoice000003 itself, so it can't be the replacement either.
	replaceProfile(t, app, "cwinvoice000002", "author@soup.com", map[string]any{"replacement": "cwinvoice000003"}, http.StatusConflict)
	if _, err := app.FindRecordById("client_invoicing_information", "cwinvoice000001"); err != nil {
		t.Errorf("a refused replace deleted the profile: %v", err)
	}
	if got := profileUses(t, app, "jobs", "cwinvoice000001"); len(got) != 2 {
		t.Errorf("a refused replace moved jobs: %v remain", got)
	}
}
