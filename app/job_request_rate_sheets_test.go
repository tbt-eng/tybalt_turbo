package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"tybalt/hooks"
	"tybalt/utilities"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/tests"
)

// Proposed rates in a project request (issue 187). The fixture request
// rsdraftreq00001, made by time@test.com, owns the draft rate sheet
// rsdraftsheet001: the 2025 Standard Rates (c41ofep525bcacj) with the
// Principal at 260 instead of 250. The fixture request rsdraftreq00002 owns
// rsdraftsheet002, a copy of that draft without a rate for the Survey Drone
// Pilot, as a draft saved before that role was added. Other tests start a
// request through the API and propose rates the same way the request editor
// does.

const (
	standardRateSheet   = "c41ofep525bcacj"
	standardSheetName   = "2025 Standard Rates"
	principalRole       = "tbgoiwwwfj8cvju"
	draftFixtureRequest = "rsdraftreq00001"
	draftFixtureSheet   = "rsdraftsheet001"
	draftFixtureCreator = "time@test.com"
	// The draft rsdraftsheet002 has no rate for this role.
	incompleteDraftSheet = "rsdraftsheet002"
	missingDraftRole     = "1qortezd8bfc513"
	missingDraftRoleName = "Survey Drone Pilot"
)

// standardRates returns the standard sheet's entries, with the Principal's
// rate replaced when principal is not zero.
func standardRates(t *testing.T, app *tests.TestApp, principal float64) []map[string]any {
	t.Helper()
	values, err := hooks.LoadRateSheetEntryValues(app, standardRateSheet)
	if err != nil {
		t.Fatal(err)
	}
	entries := []map[string]any{}
	for _, value := range values {
		rate, overtime := value.Rate, value.OvertimeRate
		if value.Role == principalRole && principal != 0 {
			rate, overtime = principal, principal*1.3
		}
		entries = append(entries, map[string]any{"role": value.Role, "rate": rate, "overtime_rate": overtime})
	}
	return entries
}

func proposedRates(name, reason string, entries []map[string]any) map[string]any {
	return map[string]any{"name": name, "reason": reason, "based_on": standardRateSheet, "entries": entries}
}

// proposeRates saves proposed rates on a request as its creator and returns
// the response.
func proposeRates(t *testing.T, app *tests.TestApp, id, email string, rates map[string]any, status int) map[string]any {
	t.Helper()
	r := workflowRecord(t, app, "job_requests", id)
	response := performJSONRequest(t, app, http.MethodPatch, "/api/job_requests/"+id, workflowToken(t, app, email), map[string]any{"version": r.GetInt("version"), "request": map[string]any{}, "proposed_rates": rates})
	mustStatus(t, response, status)
	result := map[string]any{}
	_ = json.Unmarshal(response.Body.Bytes(), &result)
	return result
}

// proposedRequest makes a complete request whose rate sheet is a new draft
// with the Principal at principal, and returns the request and draft ids.
func proposedRequest(t *testing.T, app *tests.TestApp, principal float64, reason string) (string, string) {
	t.Helper()
	id := createWorkflowRequest(t, app)
	completeWorkflowRequest(t, app, id)
	proposeRates(t, app, id, workflowAuthor, proposedRates("Client Negotiated Rates", reason, standardRates(t, app, principal)), http.StatusOK)
	uploadRequestDocument(t, app, id)
	draft := workflowRecord(t, app, "job_requests", id).GetString("rate_sheet")
	if draft == standardRateSheet || draft == "" {
		t.Fatalf("the request does not use its proposed rates: %q", draft)
	}
	return id, draft
}

// ratesApprovalData is the approver's action data with the statement that
// approves the new rate sheet and its rates.
func ratesApprovalData(t *testing.T, app *tests.TestApp, id string) map[string]any {
	t.Helper()
	data := requestActionData(t, app, id)
	data["rates_approved"] = true
	return data
}

func requestJSON(t *testing.T, app *tests.TestApp, method, path, email string, body any, status int) map[string]any {
	t.Helper()
	response := performJSONRequest(t, app, method, path, workflowToken(t, app, email), body)
	mustStatus(t, response, status)
	result := map[string]any{}
	_ = json.Unmarshal(response.Body.Bytes(), &result)
	return result
}

func duplicateMatchNames(t *testing.T, result map[string]any, field string) []string {
	t.Helper()
	data, _ := result["data"].(map[string]any)
	fieldErr, _ := data[field].(map[string]any)
	if fieldErr["code"] != "duplicate_rate_sheet" {
		t.Fatalf("not a duplicate refusal under %s: %#v", field, result)
	}
	names := []string{}
	for _, match := range fieldErr["data"].(map[string]any)["matches"].([]any) {
		names = append(names, match.(map[string]any)["name"].(string))
	}
	return names
}

func TestProposedRatesDraftBelongsToTheRequest(t *testing.T) {
	app := newWorkflowTestApp(t, true)
	id := createWorkflowRequest(t, app)
	completeWorkflowRequest(t, app, id)
	before := workflowRecord(t, app, "job_requests", id).GetInt("version")
	detail := proposeRates(t, app, id, workflowAuthor, proposedRates("  Client Negotiated Rates ", "Client MSA, schedule B", standardRates(t, app, 270)), http.StatusOK)

	r := workflowRecord(t, app, "job_requests", id)
	if r.GetInt("version") != before+1 {
		t.Fatalf("proposing rates did not count as an edit: version %d", r.GetInt("version"))
	}
	draft := workflowRecord(t, app, "rate_sheets", r.GetString("rate_sheet"))
	if !hooks.IsRateSheetDraft(draft) || draft.GetString("request") != id || draft.GetBool("active") ||
		draft.GetString("name") != "Client Negotiated Rates" || draft.GetString("reason") != "Client MSA, schedule B" ||
		draft.GetString("based_on") != standardRateSheet || draft.GetInt("revision") != 0 {
		t.Fatalf("draft fields: %#v", draft.PublicExport())
	}
	if count, _ := app.CountRecords("rate_sheet_entries", dbx.HashExp{"rate_sheet": draft.Id}); count != 33 {
		t.Fatalf("draft has %d entries", count)
	}

	// The detail shows the proposal beside the sheet it started from.
	related := detail["related"].(map[string]any)
	if related["rate_sheet_proposed"] != true || related["rate_sheet_source"].(map[string]any)["id"] != standardRateSheet {
		t.Fatalf("detail does not show the proposal: %#v", related["rate_sheet_source"])
	}
	changed := 0
	for _, row := range related["rate_sheet_entries"].([]any) {
		entry := row.(map[string]any)
		if entry["changed"] == true {
			changed++
			if entry["role_id"] != principalRole || entry["source_rate"] != float64(250) || entry["rate"] != float64(270) {
				t.Fatalf("wrong change: %#v", entry)
			}
		}
	}
	if changed != 1 {
		t.Fatalf("%d rates marked changed, want 1", changed)
	}

	// The same proposal again changes nothing.
	proposeRates(t, app, id, workflowAuthor, proposedRates("Client Negotiated Rates", "Client MSA, schedule B", standardRates(t, app, 270)), http.StatusOK)
	if v := workflowRecord(t, app, "job_requests", id).GetInt("version"); v != before+1 {
		t.Fatalf("an unchanged proposal changed the version to %d", v)
	}
	// A change updates the same draft.
	proposeRates(t, app, id, workflowAuthor, proposedRates("Client Negotiated Rates", "Client MSA, schedule B", standardRates(t, app, 280)), http.StatusOK)
	if got := workflowRecord(t, app, "job_requests", id).GetString("rate_sheet"); got != draft.Id {
		t.Fatalf("a change made a second draft: %s", got)
	}
	if drafts, _ := hooks.RequestRateSheetDrafts(app, id); len(drafts) != 1 {
		t.Fatalf("request has %d drafts", len(drafts))
	}
	// Rates must be whole and each role appears once.
	bad := standardRates(t, app, 280)
	bad[0]["rate"] = 100.5
	proposeRates(t, app, id, workflowAuthor, proposedRates("Client Negotiated Rates", "Client MSA, schedule B", bad), http.StatusBadRequest)
	// Another sheet's name, in any case, is refused.
	proposeRates(t, app, id, workflowAuthor, proposedRates(strings.ToUpper(standardSheetName), "Client MSA, schedule B", standardRates(t, app, 280)), http.StatusBadRequest)
	// Choosing a sheet and proposing rates at once is refused.
	r = workflowRecord(t, app, "job_requests", id)
	requestJSON(t, app, http.MethodPatch, "/api/job_requests/"+id, workflowAuthor, map[string]any{"version": r.GetInt("version"), "request": map[string]any{"rate_sheet": "rs_sheetb_r0"}, "proposed_rates": proposedRates("Client Negotiated Rates", "Client MSA, schedule B", standardRates(t, app, 280))}, http.StatusBadRequest)

	// Choosing an active sheet instead deletes the draft and its entries.
	requestJSON(t, app, http.MethodPatch, "/api/job_requests/"+id, workflowAuthor, map[string]any{"version": r.GetInt("version"), "request": map[string]any{"rate_sheet": standardRateSheet}}, http.StatusOK)
	if _, err := app.FindRecordById("rate_sheets", draft.Id); err == nil {
		t.Fatal("choosing a sheet kept the draft")
	}
	if count, _ := app.CountRecords("rate_sheet_entries", dbx.HashExp{"rate_sheet": draft.Id}); count != 0 {
		t.Fatalf("the draft left %d entries", count)
	}
}

// The draft is listed nowhere and no rate sheet route or job accepts it.
func TestProposedRatesDraftIsHiddenAndUnusable(t *testing.T) {
	app := newWorkflowTestApp(t, true)
	author := workflowToken(t, app, workflowAuthor) // job, admin

	list := performJSONRequest(t, app, http.MethodGet, "/api/collections/rate_sheets/records?perPage=500", author, nil)
	mustStatus(t, list, http.StatusOK)
	if strings.Contains(list.Body.String(), draftFixtureSheet) {
		t.Fatal("the rate sheet list shows the draft")
	}
	mustStatus(t, performJSONRequest(t, app, http.MethodGet, "/api/collections/rate_sheets/records/"+draftFixtureSheet, author, nil), http.StatusNotFound)
	entries := performJSONRequest(t, app, http.MethodGet, "/api/collections/rate_sheet_entries/records?filter=rate_sheet%3D%22"+draftFixtureSheet+"%22", author, nil)
	mustStatus(t, entries, http.StatusOK)
	if strings.Contains(entries.Body.String(), draftFixtureSheet) {
		t.Fatal("the entries list shows the draft's rates")
	}
	augmented := performJSONRequest(t, app, http.MethodGet, "/api/collections/rate_sheets_augmented/records?perPage=500", author, nil)
	mustStatus(t, augmented, http.StatusOK)
	if strings.Contains(augmented.Body.String(), draftFixtureSheet) {
		t.Fatal("the Rate Sheets list shows the draft")
	}
	options := clientSetupJSON(t, app, http.MethodGet, "/api/job_requests/options", workflowAuthor, nil, http.StatusOK)
	for _, row := range options["rate_sheets"].([]any) {
		if row.(map[string]any)["id"] == draftFixtureSheet {
			t.Fatal("the request options list the draft")
		}
	}

	// Rate sheet routes refuse it.
	mustStatus(t, performJSONRequest(t, app, http.MethodPost, "/api/rate_sheets/"+draftFixtureSheet+"/activate", author, nil), http.StatusBadRequest)
	mustStatus(t, performJSONRequest(t, app, http.MethodPost, "/api/rate_sheets/"+draftFixtureSheet+"/rename", author, map[string]any{"name": "Renamed Draft"}), http.StatusBadRequest)
	entry, err := app.FindFirstRecordByData("rate_sheet_entries", "rate_sheet", draftFixtureSheet)
	if err != nil {
		t.Fatal(err)
	}
	mustStatus(t, performJSONRequest(t, app, http.MethodPut, "/api/rate_sheet_entries/"+entry.Id, author, map[string]any{"rate": 999}), http.StatusBadRequest)
	// time@test.com holds rate_sheet_revise.
	revise := performJSONRequest(t, app, http.MethodPost, "/api/rate_sheets", workflowToken(t, app, draftFixtureCreator), map[string]any{"name": "Proposed Fixture Rates", "effective_date": "2026-12-01", "revision": 1, "entries": standardRates(t, app, 0)})
	mustStatus(t, revise, http.StatusBadRequest)

	// The collection API can't make a draft or add rates to one. The rate is
	// for a role the draft lacks, so only the create rule can refuse it, not
	// the unique index on (rate_sheet, role).
	mustStatus(t, performJSONRequest(t, app, http.MethodPost, "/api/collections/rate_sheets/records", author, map[string]any{"name": "Sneaky Draft", "effective_date": "2026-11-13", "revision": 0, "request": draftFixtureRequest}), http.StatusBadRequest)
	if !hooks.IsRateSheetDraft(workflowRecord(t, app, "rate_sheets", incompleteDraftSheet)) {
		t.Fatal("the fixture sheet is not a draft")
	}
	if _, err := app.FindFirstRecordByFilter("rate_sheet_entries", "rate_sheet = {:sheet} && role = {:role}", dbx.Params{"sheet": incompleteDraftSheet, "role": missingDraftRole}); err == nil {
		t.Fatal("the fixture draft already has a rate for the role")
	}
	mustStatus(t, performJSONRequest(t, app, http.MethodPost, "/api/collections/rate_sheet_entries/records", author, map[string]any{"rate_sheet": incompleteDraftSheet, "role": missingDraftRole, "rate": 1, "overtime_rate": 1}), http.StatusBadRequest)
	if _, err := app.FindFirstRecordByFilter("rate_sheet_entries", "rate_sheet = {:sheet} && role = {:role}", dbx.Params{"sheet": incompleteDraftSheet, "role": missingDraftRole}); err == nil {
		t.Fatal("the collection API added a rate to a draft")
	}
	// The same writes without a draft pass, so the refusals above come from
	// the draft rules.
	plain := performJSONRequest(t, app, http.MethodPost, "/api/collections/rate_sheets/records", author, map[string]any{"name": "Plain Sheet", "effective_date": "2026-11-13", "revision": 0})
	mustStatus(t, plain, http.StatusOK)
	plainSheet := map[string]any{}
	if err := json.Unmarshal(plain.Body.Bytes(), &plainSheet); err != nil {
		t.Fatal(err)
	}
	mustStatus(t, performJSONRequest(t, app, http.MethodPost, "/api/collections/rate_sheet_entries/records", author, map[string]any{"rate_sheet": plainSheet["id"], "role": missingDraftRole, "rate": 1, "overtime_rate": 1}), http.StatusOK)

	// Another request can't choose it.
	other := createWorkflowRequest(t, app)
	requestJSON(t, app, http.MethodPatch, "/api/job_requests/"+other, workflowAuthor, map[string]any{"version": 1, "request": map[string]any{"rate_sheet": draftFixtureSheet}}, http.StatusBadRequest)

	// A project can't use it: it is inactive.
	job := workflowRecord(t, app, "jobs", "pcactivedated01")
	if job.GetString("rate_sheet") == "" {
		t.Fatal("fixture project has no rate sheet")
	}
	// history.revise@example.com holds job and rate_sheet_revise.
	history := performJSONRequest(t, app, http.MethodPost, "/api/jobs/"+job.Id+"/rate_sheets", workflowToken(t, app, "history.revise@example.com"), map[string]any{"rate_sheet": draftFixtureSheet, "effective_date": hooks.RateSheetToday(), "reason": "Use the draft"})
	mustStatus(t, history, http.StatusBadRequest)
	if !strings.Contains(history.Body.String(), "active") {
		t.Fatalf("a project period refused the draft for another reason: %s", history.Body.String())
	}

	// A rate sheet absorb sees no draft.
	duplicates, err := hooks.FindRateSheetDuplicates(app)
	if err != nil {
		t.Fatal(err)
	}
	for _, group := range duplicates {
		raw, _ := json.Marshal(group)
		if strings.Contains(string(raw), draftFixtureSheet) {
			t.Fatal("the duplicates view shows the draft")
		}
	}
}

// Only the creator edits the draft, through the request, while the request is
// editable.
func TestProposedRatesPermissions(t *testing.T) {
	app := newWorkflowTestApp(t, true)
	rates := proposedRates("Proposed Fixture Rates", "Changed reason", standardRates(t, app, 265))
	proposeRates(t, app, draftFixtureRequest, workflowAuthor, rates, http.StatusForbidden)
	proposeRates(t, app, draftFixtureRequest, workflowApprover, rates, http.StatusForbidden)
	setRequestStatus(t, app, draftFixtureRequest, "Awaiting approval")
	proposeRates(t, app, draftFixtureRequest, draftFixtureCreator, rates, http.StatusBadRequest)
	setRequestStatus(t, app, draftFixtureRequest, "Draft")
	proposeRates(t, app, draftFixtureRequest, draftFixtureCreator, rates, http.StatusOK)
	if reason := workflowRecord(t, app, "rate_sheets", draftFixtureSheet).GetString("reason"); reason != "Changed reason" {
		t.Fatalf("creator edit not saved: %q", reason)
	}
	// Only the creator sees the draft through the request.
	clientSetupJSON(t, app, http.MethodGet, "/api/job_requests/"+draftFixtureRequest, "corp.noclaim@example.com", nil, http.StatusForbidden)
	// The draft must start from an active sheet.
	inactive := proposedRates("Proposed Fixture Rates", "Changed reason", standardRates(t, app, 265))
	inactive["based_on"] = "rs_hook_inactive_001"
	proposeRates(t, app, draftFixtureRequest, draftFixtureCreator, inactive, http.StatusBadRequest)
}

// Every save of proposed rates needs a reason, as it needs a name and the sheet
// to start from. The refusal is a field error under proposed_rates.reason.
func TestProposedRatesNeedAReasonToSave(t *testing.T) {
	app := newWorkflowTestApp(t, true)
	id := createWorkflowRequest(t, app)
	completeWorkflowRequest(t, app, id)
	for _, reason := range []string{"", "   "} {
		result := proposeRates(t, app, id, workflowAuthor, proposedRates("Client Negotiated Rates", reason, standardRates(t, app, 260)), http.StatusBadRequest)
		data, _ := result["data"].(map[string]any)
		fieldErr, _ := data[hooks.ProposedRatesField+".reason"].(map[string]any)
		if fieldErr["code"] != "value_required" {
			t.Fatalf("reason %q: no field error for the reason: %#v", reason, result)
		}
		if drafts, _ := hooks.RequestRateSheetDrafts(app, id); len(drafts) != 0 {
			t.Fatalf("reason %q: a refused save made a draft", reason)
		}
	}
	long := strings.Repeat("x", 501)
	proposeRates(t, app, id, workflowAuthor, proposedRates("Client Negotiated Rates", long, standardRates(t, app, 260)), http.StatusBadRequest)
	proposeRates(t, app, id, workflowAuthor, proposedRates("Client Negotiated Rates", "Client MSA, schedule B", standardRates(t, app, 260)), http.StatusOK)

	// A saved draft can't lose its reason either.
	draftID := workflowRecord(t, app, "job_requests", id).GetString("rate_sheet")
	proposeRates(t, app, id, workflowAuthor, proposedRates("Client Negotiated Rates", "", standardRates(t, app, 260)), http.StatusBadRequest)
	if reason := workflowRecord(t, app, "rate_sheets", draftID).GetString("reason"); reason != "Client MSA, schedule B" {
		t.Fatalf("a refused save changed the reason to %q", reason)
	}
}

// Submission checks the reason again, for a draft saved before every save
// needed one. This test clears the reason of its draft directly, since no
// save can.
func TestProposedRatesNeedAReasonToSubmit(t *testing.T) {
	app := newWorkflowTestApp(t, true)
	id, draftID := proposedRequest(t, app, 260, "Client MSA, schedule B")
	draft := workflowRecord(t, app, "rate_sheets", draftID)
	draft.Set("reason", "")
	if err := app.Save(draft); err != nil {
		t.Fatal(err)
	}
	response := performJSONRequest(t, app, http.MethodPost, "/api/job_requests/"+id+"/submit", workflowToken(t, app, workflowAuthor), requestActionData(t, app, id))
	mustStatus(t, response, http.StatusBadRequest)
	if !strings.Contains(response.Body.String(), "reason") {
		t.Fatalf("the refusal does not ask for the reason: %s", response.Body.String())
	}
	proposeRates(t, app, id, workflowAuthor, proposedRates("Client Negotiated Rates", "Client MSA, schedule B", standardRates(t, app, 260)), http.StatusOK)
	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/submit", workflowAuthor, requestActionData(t, app, id), http.StatusOK)
}

// A proposal can lack a rate for a role: the source sheet may have none, or
// the draft was saved before the role was added. The request form then sends
// the rates it has, and the save works; submission, approval and job creation
// name the role. The fixture draft rsdraftsheet002 has no rate for the
// Survey Drone Pilot.
func TestProposedRatesWithoutEveryRoleNameTheMissingRole(t *testing.T) {
	app := newWorkflowTestApp(t, true)
	const request = "rsdraftreq00002"
	incomplete := func() {
		t.Helper()
		err := hooks.ValidateRequestRateSheetDraft(app, workflowRecord(t, app, "rate_sheets", incompleteDraftSheet))
		if err == nil || !strings.Contains(err.Error(), "no rate for "+missingDraftRoleName) {
			t.Fatalf("the check does not name the missing role: %v", err)
		}
	}
	incomplete()

	// The request detail lists the rates the draft has, so the form adds an
	// empty row for the missing role.
	detail := clientSetupJSON(t, app, http.MethodGet, "/api/job_requests/"+request, draftFixtureCreator, nil, http.StatusOK)
	for _, row := range detail["related"].(map[string]any)["rate_sheet_entries"].([]any) {
		if row.(map[string]any)["role_id"] == missingDraftRole {
			t.Fatal("the detail lists a rate the draft lacks")
		}
	}

	// A save without the role works and keeps the draft incomplete.
	entries := []map[string]any{}
	for _, entry := range standardRates(t, app, 265) {
		if entry["role"] != missingDraftRole {
			entries = append(entries, entry)
		}
	}
	proposeRates(t, app, request, draftFixtureCreator, proposedRates("Proposed Fixture Rates Missing A Role", "Client rate schedule", entries), http.StatusOK)
	incomplete()

	// With the role, the draft is complete.
	proposeRates(t, app, request, draftFixtureCreator, proposedRates("Proposed Fixture Rates Missing A Role", "Client rate schedule", standardRates(t, app, 265)), http.StatusOK)
	if err := hooks.ValidateRequestRateSheetDraft(app, workflowRecord(t, app, "rate_sheets", incompleteDraftSheet)); err != nil {
		t.Fatalf("a complete draft is refused: %v", err)
	}
}

// Turbo refuses proposed rates that an active sheet already has and names it,
// and the requester can choose that sheet instead.
func TestProposedRatesThatMatchAnActiveSheetAreRefused(t *testing.T) {
	app := newWorkflowTestApp(t, true)
	id := createWorkflowRequest(t, app)
	result := proposeRates(t, app, id, workflowAuthor, proposedRates("Same As Standard", "Client MSA", standardRates(t, app, 0)), http.StatusConflict)
	names := duplicateMatchNames(t, result, hooks.ProposedRatesField)
	if len(names) != 1 || names[0] != standardSheetName {
		t.Fatalf("the refusal names %v", names)
	}
	if drafts, _ := hooks.RequestRateSheetDrafts(app, id); len(drafts) != 0 {
		t.Fatal("a refused proposal left a draft")
	}
	requestJSON(t, app, http.MethodPatch, "/api/job_requests/"+id, workflowAuthor, map[string]any{"version": workflowRecord(t, app, "job_requests", id).GetInt("version"), "request": map[string]any{"rate_sheet": standardRateSheet}}, http.StatusOK)
	if got := workflowRecord(t, app, "job_requests", id).GetString("rate_sheet"); got != standardRateSheet {
		t.Fatalf("the request did not take the matching sheet: %s", got)
	}
}

// An active sheet with the same rates can appear after the draft is saved. The
// check runs again at submission, approval and creation, and creation leaves
// the request and its draft as they were.
func TestProposedRatesAreCheckedAgainWhenTheJobIsCreated(t *testing.T) {
	app := newWorkflowTestApp(t, true)
	// A job holder makes an inactive sheet with the same rates as the
	// proposal. Inactive sheets are not compared, so both can be saved.
	directID := requestJSON(t, app, http.MethodPost, "/api/rate_sheets", workflowAuthor, map[string]any{"name": "Direct Copy", "effective_date": "2026-11-01", "revision": 0, "entries": standardRates(t, app, 275)}, http.StatusCreated)["id"].(string)

	id, draftID := proposedRequest(t, app, 275, "Client MSA")
	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/submit", workflowAuthor, requestActionData(t, app, id), http.StatusOK)
	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/approve", workflowApprover, ratesApprovalData(t, app, id), http.StatusOK)

	// The draft is not active, so the direct sheet can still be activated.
	clientSetupJSON(t, app, http.MethodPost, "/api/rate_sheets/"+directID+"/activate", workflowAuthor, nil, http.StatusOK)
	result := requestJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/create", workflowAuthor, requestActionData(t, app, id), http.StatusConflict)
	if names := duplicateMatchNames(t, result, hooks.ProposedRatesField); len(names) != 1 || names[0] != "Direct Copy" {
		t.Fatalf("the refusal names %v", names)
	}
	draft := workflowRecord(t, app, "rate_sheets", draftID)
	if draft.GetBool("active") || draft.GetString("request") != id {
		t.Fatal("a refused creation changed the draft")
	}
	if r := workflowRecord(t, app, "job_requests", id); r.GetString("status") != "Awaiting creation" {
		t.Fatalf("a refused creation changed the request: %s", r.GetString("status"))
	}
}

// The approver approves the proposed rates with the request, and job creation
// activates them and gives the job its opening rate sheet period.
func TestProposedRatesApprovedAndActivatedWithTheJob(t *testing.T) {
	app := newWorkflowTestApp(t, true)
	id, draftID := proposedRequest(t, app, 290, "Client MSA, schedule B")
	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/submit", workflowAuthor, requestActionData(t, app, id), http.StatusOK)

	// The checklist keeps its stored statements for proposed rates.
	detail := clientSetupJSON(t, app, http.MethodGet, "/api/job_requests/"+id, workflowApprover, nil, http.StatusOK)
	_, stored, err := utilities.CurrentRequestChecklist(app, utilities.RequestChecklistScopeApproval)
	if err != nil {
		t.Fatal(err)
	}
	shown := detail["checklist"].(map[string]any)["questions"].([]any)
	if len(shown) != len(stored) {
		t.Fatalf("approval checklist: %#v", detail["checklist"])
	}
	for i, question := range shown {
		if q := question.(map[string]any); q["key"] != stored[i].Key || q["text"] != stored[i].Text {
			t.Fatalf("checklist statement %d: %#v", i, q)
		}
	}
	if detail["related"].(map[string]any)["rate_sheet_proposed"] != true {
		t.Fatal("the approver does not see the proposal")
	}
	// Only the selected approver approves, and a job holder can't create first.
	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/create", workflowAuthor, requestActionData(t, app, id), http.StatusBadRequest)
	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/approve", workflowAuthor, requestActionData(t, app, id), http.StatusForbidden)
	// The approver must approve the new sheet and its rates as well.
	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/approve", workflowApprover, requestActionData(t, app, id), http.StatusBadRequest)
	if workflowRecord(t, app, "job_requests", id).GetString("status") != "Awaiting approval" {
		t.Fatal("the request was approved without the rates")
	}
	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/approve", workflowApprover, ratesApprovalData(t, app, id), http.StatusOK)
	approver := workflowRecord(t, app, "job_requests", id).GetString("approver")
	draft := workflowRecord(t, app, "rate_sheets", draftID)
	if draft.GetBool("active") {
		t.Fatal("approval activated the draft")
	}
	if draft.GetString("rates_approver") != approver || draft.GetString("rates_approved") == "" {
		t.Fatalf("rates approval: %#v", draft.PublicExport())
	}
	// The job creator sees who approved the rates.
	detail = clientSetupJSON(t, app, http.MethodGet, "/api/job_requests/"+id, workflowAuthor, nil, http.StatusOK)
	if person, _ := detail["related"].(map[string]any)["rate_sheet_rates_approver"].(map[string]any); person["id"] != approver {
		t.Fatalf("rates approver in the detail: %#v", person)
	}
	// A job holder makes a sheet with the same rates directly. The draft is
	// inactive, so the sheet can be saved, inactive.
	lateID := requestJSON(t, app, http.MethodPost, "/api/rate_sheets", workflowAuthor, map[string]any{"name": "Late Copy", "effective_date": "2026-11-01", "revision": 0, "entries": standardRates(t, app, 290)}, http.StatusCreated)["id"].(string)
	// Approval needs no job claim, and creation needs one.
	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/create", workflowApprover, requestActionData(t, app, id), http.StatusForbidden)

	result := clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/create", workflowAuthor, requestActionData(t, app, id), http.StatusOK)
	jobID := result["job_id"].(string)
	sheet := workflowRecord(t, app, "rate_sheets", draftID)
	if !sheet.GetBool("active") || sheet.GetString("request") != "" || sheet.GetString("reason") != "Client MSA, schedule B" ||
		sheet.GetString("effective_date") != hooks.RateSheetToday() || sheet.GetString("based_on") != standardRateSheet ||
		sheet.GetString("rates_approver") != approver || sheet.GetString("rates_approved") == "" {
		t.Fatalf("activated sheet: %#v", sheet.PublicExport())
	}
	job := workflowRecord(t, app, "jobs", jobID)
	if job.GetString("rate_sheet") != draftID {
		t.Fatalf("job uses %s", job.GetString("rate_sheet"))
	}
	rows, err := hooks.JobRateSheetRows(app, jobID)
	if err != nil || len(rows) != 1 || rows[0].GetString("rate_sheet") != draftID || rows[0].GetString("effective_date") != "" {
		t.Fatalf("opening period: %v %d", err, len(rows))
	}
	if _, err := app.FindRecordById("job_requests", id); err == nil {
		t.Fatal("the request was not deleted")
	}
	// The sheet is now listed like any other.
	list := performJSONRequest(t, app, http.MethodGet, "/api/collections/rate_sheets/records/"+draftID, workflowToken(t, app, workflowAuthor), nil)
	mustStatus(t, list, http.StatusOK)

	// Direct activation of the job holder's sheet with the same rates is
	// refused and names the activated proposal.
	refused := requestJSON(t, app, http.MethodPost, "/api/rate_sheets/"+lateID+"/activate", workflowAuthor, nil, http.StatusConflict)
	if names := duplicateMatchNames(t, refused, "active"); len(names) != 1 || names[0] != "Client Negotiated Rates" {
		t.Fatalf("activation refusal names %v", names)
	}
}

// The job holder who creates the job can rename the proposed sheet. A name that
// another sheet has is refused, and the old name is kept as a former name.
func TestJobCreatorCanRenameProposedRates(t *testing.T) {
	app := newWorkflowTestApp(t, true)
	id, draftID := proposedRequest(t, app, 290, "Client MSA, schedule B")
	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/submit", workflowAuthor, requestActionData(t, app, id), http.StatusOK)
	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/approve", workflowApprover, ratesApprovalData(t, app, id), http.StatusOK)

	withName := func(name string) map[string]any {
		data := requestActionData(t, app, id)
		data["rate_sheet_name"] = name
		return data
	}
	// A name in use, in any letter case, is refused and nothing is created.
	taken := workflowRecord(t, app, "rate_sheets", standardRateSheet).GetString("name")
	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/create", workflowAuthor, withName(strings.ToUpper(taken)), http.StatusBadRequest)
	if sheet := workflowRecord(t, app, "rate_sheets", draftID); sheet.GetBool("active") || sheet.GetString("name") != "Client Negotiated Rates" {
		t.Fatalf("refused rename changed the draft: %#v", sheet.PublicExport())
	}

	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/create", workflowAuthor, withName("  Sepulchi MSA Rates  "), http.StatusOK)
	sheet := workflowRecord(t, app, "rate_sheets", draftID)
	if !sheet.GetBool("active") || sheet.GetString("name") != "Sepulchi MSA Rates" {
		t.Fatalf("renamed sheet: %#v", sheet.PublicExport())
	}
	changes, err := app.FindRecordsByFilter("rate_sheet_name_changes", "rate_sheet = {:sheet}", "", 0, 0, dbx.Params{"sheet": draftID})
	if err != nil || len(changes) != 1 || changes[0].GetString("old_name") != "Client Negotiated Rates" || changes[0].GetString("new_name") != "Sepulchi MSA Rates" {
		t.Fatalf("name change: %v %d", err, len(changes))
	}
}

// A returned request keeps its draft, and the creator can still change it.
func TestRejectedRequestKeepsAnEditableDraft(t *testing.T) {
	app := newWorkflowTestApp(t, true)
	id, draftID := proposedRequest(t, app, 285, "Client MSA")
	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/submit", workflowAuthor, requestActionData(t, app, id), http.StatusOK)
	reject := requestActionData(t, app, id)
	reject["rejection_reason"] = "Rates need a second look"
	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/reject", workflowApprover, reject, http.StatusOK)

	r := workflowRecord(t, app, "job_requests", id)
	if r.GetString("status") != "Rejected" || r.GetString("rate_sheet") != draftID {
		t.Fatalf("rejected request: %s %s", r.GetString("status"), r.GetString("rate_sheet"))
	}
	if draft := workflowRecord(t, app, "rate_sheets", draftID); draft.GetString("request") != id || draft.GetBool("active") {
		t.Fatal("rejection changed the draft")
	}
	proposeRates(t, app, id, workflowAuthor, proposedRates("Client Negotiated Rates", "Client MSA, corrected", standardRates(t, app, 286)), http.StatusOK)
	r = workflowRecord(t, app, "job_requests", id)
	if r.GetString("status") != "Draft" || r.GetString("rate_sheet") != draftID {
		t.Fatalf("edited request: %s %s", r.GetString("status"), r.GetString("rate_sheet"))
	}
	if reason := workflowRecord(t, app, "rate_sheets", draftID).GetString("reason"); reason != "Client MSA, corrected" {
		t.Fatalf("draft not edited: %q", reason)
	}
}

// Deleting a request deletes its draft and the draft's rates.
func TestDeletingARequestDeletesItsDraft(t *testing.T) {
	app := newWorkflowTestApp(t, true)
	clientSetupJSON(t, app, http.MethodDelete, "/api/job_requests/"+draftFixtureRequest, workflowAuthor, nil, http.StatusBadRequest)
	if _, err := app.FindRecordById("rate_sheets", draftFixtureSheet); err != nil {
		t.Fatal("a refused delete removed the draft")
	}
	clientSetupJSON(t, app, http.MethodDelete, "/api/job_requests/"+draftFixtureRequest, draftFixtureCreator, nil, http.StatusNoContent)
	if _, err := app.FindRecordById("rate_sheets", draftFixtureSheet); err == nil {
		t.Fatal("the draft outlived its request")
	}
	if count, _ := app.CountRecords("rate_sheet_entries", dbx.HashExp{"rate_sheet": draftFixtureSheet}); count != 0 {
		t.Fatalf("the draft left %d entries", count)
	}
	// The sheet it started from is untouched.
	if count, _ := app.CountRecords("rate_sheet_entries", dbx.HashExp{"rate_sheet": standardRateSheet}); count != 33 {
		t.Fatalf("the source sheet has %d entries", count)
	}
}

// Recall ends the approval of the rates, so the next approval states it again,
// and nobody can set the approval through the collection API.
func TestRecallClearsTheRatesApproval(t *testing.T) {
	app := newWorkflowTestApp(t, true)
	id, draftID := proposedRequest(t, app, 290, "Client MSA, schedule B")
	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/submit", workflowAuthor, requestActionData(t, app, id), http.StatusOK)
	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/approve", workflowApprover, ratesApprovalData(t, app, id), http.StatusOK)
	if workflowRecord(t, app, "rate_sheets", draftID).GetString("rates_approver") == "" {
		t.Fatal("the approval was not recorded")
	}

	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/recall", workflowAuthor, requestActionData(t, app, id), http.StatusOK)
	draft := workflowRecord(t, app, "rate_sheets", draftID)
	if draft.GetString("rates_approver") != "" || draft.GetString("rates_approved") != "" {
		t.Fatalf("recall kept the approval: %#v", draft.PublicExport())
	}
	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/submit", workflowAuthor, requestActionData(t, app, id), http.StatusOK)
	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/approve", workflowApprover, requestActionData(t, app, id), http.StatusBadRequest)

	// A sheet made directly can't claim an approval.
	approver := workflowRecord(t, app, "job_requests", id).GetString("approver")
	requestJSON(t, app, http.MethodPost, "/api/collections/rate_sheets/records", workflowAuthor, map[string]any{"name": "Claimed Approval", "effective_date": "2026-11-01", "revision": 0, "rates_approver": approver}, http.StatusBadRequest)
}

// The requester sees the active sheets close to the proposed rates before the
// save. The standard sheet with one changed rate differs from it in one role;
// with no changed rate it is an exact match. The sheet the requester started
// from is not listed as close. noclaims@example.com lacks the
// time claim, which starting a request needs.
func TestRequesterComparesProposedRatesWithActiveSheets(t *testing.T) {
	app := newWorkflowTestApp(t, true)
	body := func(principal float64) map[string]any {
		return map[string]any{"name": "Client Negotiated Rates", "entries": standardRates(t, app, principal)}
	}
	matches := func(result map[string]any, key string) []map[string]any {
		list := []map[string]any{}
		for _, match := range result[key].([]any) {
			list = append(list, match.(map[string]any))
		}
		return list
	}

	result := requestJSON(t, app, http.MethodPost, "/api/job_requests/compare_rates", workflowAuthor, body(260), http.StatusOK)
	if exact := matches(result, "exact"); len(exact) != 0 {
		t.Fatalf("changed rates have exact matches: %v", exact)
	}
	near := matches(result, "near")
	if len(near) == 0 || len(near) > hooks.RateSheetNearMatchLimit {
		t.Fatalf("near matches: %v", near)
	}
	if near[0]["id"] != standardRateSheet || near[0]["differing_roles"] != float64(1) {
		t.Fatalf("the closest sheet is %v", near[0])
	}
	for _, match := range near {
		if match["id"] == draftFixtureSheet {
			t.Fatal("a draft is shown as a close sheet")
		}
	}

	fromStandard := body(260)
	fromStandard["based_on"] = standardRateSheet
	result = requestJSON(t, app, http.MethodPost, "/api/job_requests/compare_rates", workflowAuthor, fromStandard, http.StatusOK)
	for _, match := range matches(result, "near") {
		if match["id"] == standardRateSheet {
			t.Fatal("the source sheet is shown as a close sheet")
		}
	}

	result = requestJSON(t, app, http.MethodPost, "/api/job_requests/compare_rates", workflowAuthor, body(0), http.StatusOK)
	if exact := matches(result, "exact"); len(exact) != 1 || exact[0]["id"] != standardRateSheet {
		t.Fatalf("exact matches: %v", exact)
	}

	requestJSON(t, app, http.MethodPost, "/api/job_requests/compare_rates", "noclaims@example.com", body(260), http.StatusForbidden)
	if drafts, _ := hooks.RequestRateSheetDrafts(app, draftFixtureRequest); len(drafts) != 1 {
		t.Fatal("the comparison changed the drafts")
	}
}

func TestComparingProposedRatesNeedsProjectRequests(t *testing.T) {
	app := newWorkflowTestApp(t, false)
	requestJSON(t, app, http.MethodPost, "/api/job_requests/compare_rates", workflowAuthor, map[string]any{"name": "X", "entries": standardRates(t, app, 260)}, http.StatusForbidden)
}

// An approval of proposed rates is bound to the approver's id, the source
// sheet's id and the rates, not to display data. This test changes fixture
// data after the approval: it edits the approver's profile and renames the
// source sheet (2025 Standard Rates), and the job is still created.
func TestRatesApprovalSurvivesAProfileEditAndASourceRename(t *testing.T) {
	app := newWorkflowTestApp(t, true)
	id, draftID := proposedRequest(t, app, 290, "Client MSA, schedule B")
	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/submit", workflowAuthor, requestActionData(t, app, id), http.StatusOK)
	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/approve", workflowApprover, ratesApprovalData(t, app, id), http.StatusOK)

	approver := workflowRecord(t, app, "job_requests", id).GetString("approver")
	profile, err := app.FindFirstRecordByData("profiles", "uid", approver)
	if err != nil {
		t.Fatal(err)
	}
	profile.Set("given_name", "Renamed")
	if err := app.Save(profile); err != nil {
		t.Fatal(err)
	}
	source := workflowRecord(t, app, "rate_sheets", standardRateSheet)
	source.Set("name", "2025 Standard Rates Renamed")
	if err := app.Save(source); err != nil {
		t.Fatal(err)
	}

	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/create", workflowAuthor, requestActionData(t, app, id), http.StatusOK)
	if !workflowRecord(t, app, "rate_sheets", draftID).GetBool("active") {
		t.Fatal("the proposed sheet was not activated")
	}
}

// A corrected rate on the source sheet changes the comparison that the
// approver saw, so the approval is stale. This test changes a fixture rate of
// the source sheet (2025 Standard Rates) after the approval.
func TestRatesApprovalIsStaleAfterASourceRateChanges(t *testing.T) {
	app := newWorkflowTestApp(t, true)
	id, draftID := proposedRequest(t, app, 290, "Client MSA, schedule B")
	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/submit", workflowAuthor, requestActionData(t, app, id), http.StatusOK)
	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/approve", workflowApprover, ratesApprovalData(t, app, id), http.StatusOK)

	entries, err := app.FindRecordsByFilter("rate_sheet_entries", "rate_sheet = {:sheet} && role != {:role}", "", 1, 0, dbx.Params{"sheet": standardRateSheet, "role": principalRole})
	if err != nil || len(entries) != 1 {
		t.Fatalf("source entry: %v %d", err, len(entries))
	}
	entries[0].Set("rate", entries[0].GetFloat("rate")+1)
	if err := app.UnsafeWithoutHooks().Save(entries[0]); err != nil {
		t.Fatal(err)
	}

	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/create", workflowAuthor, requestActionData(t, app, id), http.StatusConflict)
	if workflowRecord(t, app, "rate_sheets", draftID).GetBool("active") {
		t.Fatal("a stale approval activated the proposed sheet")
	}
}
