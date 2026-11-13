package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"tybalt/hooks"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/tests"
)

// Proposed rates in a project request (issue 187). The fixture request
// rsdraftreq00001, made by time@test.com, owns the draft rate sheet
// rsdraftsheet001: the 2025 Standard Rates (c41ofep525bcacj) with the
// Principal at 260 instead of 250. Other tests start a request through the API
// and propose rates the same way the request editor does.

const (
	standardRateSheet   = "c41ofep525bcacj"
	standardSheetName   = "2025 Standard Rates"
	principalRole       = "tbgoiwwwfj8cvju"
	draftFixtureRequest = "rsdraftreq00001"
	draftFixtureSheet   = "rsdraftsheet001"
	draftFixtureCreator = "time@test.com"
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
	proposeRates(t, app, id, workflowAuthor, proposedRates("Client Negotiated Rates", "", bad), http.StatusBadRequest)
	// Another sheet's name, in any case, is refused.
	proposeRates(t, app, id, workflowAuthor, proposedRates(strings.ToUpper(standardSheetName), "", standardRates(t, app, 280)), http.StatusBadRequest)
	// Choosing a sheet and proposing rates at once is refused.
	r = workflowRecord(t, app, "job_requests", id)
	requestJSON(t, app, http.MethodPatch, "/api/job_requests/"+id, workflowAuthor, map[string]any{"version": r.GetInt("version"), "request": map[string]any{"rate_sheet": "rs_sheetb_r0"}, "proposed_rates": proposedRates("Client Negotiated Rates", "", standardRates(t, app, 280))}, http.StatusBadRequest)

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

	// The collection API can't make a draft or add rates to one.
	mustStatus(t, performJSONRequest(t, app, http.MethodPost, "/api/collections/rate_sheets/records", author, map[string]any{"name": "Sneaky Draft", "effective_date": "2026-11-13", "revision": 0, "request": draftFixtureRequest}), http.StatusBadRequest)
	mustStatus(t, performJSONRequest(t, app, http.MethodPost, "/api/collections/rate_sheet_entries/records", author, map[string]any{"rate_sheet": draftFixtureSheet, "role": entry.GetString("role"), "rate": 1, "overtime_rate": 1}), http.StatusBadRequest)
	// The same writes without a draft pass, so the refusals above come from
	// the draft rules.
	plain := performJSONRequest(t, app, http.MethodPost, "/api/collections/rate_sheets/records", author, map[string]any{"name": "Plain Sheet", "effective_date": "2026-11-13", "revision": 0})
	mustStatus(t, plain, http.StatusOK)
	plainSheet := map[string]any{}
	if err := json.Unmarshal(plain.Body.Bytes(), &plainSheet); err != nil {
		t.Fatal(err)
	}
	mustStatus(t, performJSONRequest(t, app, http.MethodPost, "/api/collections/rate_sheet_entries/records", author, map[string]any{"rate_sheet": plainSheet["id"], "role": entry.GetString("role"), "rate": 1, "overtime_rate": 1}), http.StatusOK)

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
	inactive := proposedRates("Proposed Fixture Rates", "", standardRates(t, app, 265))
	inactive["based_on"] = "rs_hook_inactive_001"
	proposeRates(t, app, draftFixtureRequest, draftFixtureCreator, inactive, http.StatusBadRequest)
}

func TestProposedRatesNeedAReasonToSubmit(t *testing.T) {
	app := newWorkflowTestApp(t, true)
	id, _ := proposedRequest(t, app, 260, "")
	response := performJSONRequest(t, app, http.MethodPost, "/api/job_requests/"+id+"/submit", workflowToken(t, app, workflowAuthor), requestActionData(t, app, id))
	mustStatus(t, response, http.StatusBadRequest)
	if !strings.Contains(response.Body.String(), "reason") {
		t.Fatalf("the refusal does not ask for the reason: %s", response.Body.String())
	}
	proposeRates(t, app, id, workflowAuthor, proposedRates("Client Negotiated Rates", "   ", standardRates(t, app, 260)), http.StatusOK)
	mustStatus(t, performJSONRequest(t, app, http.MethodPost, "/api/job_requests/"+id+"/submit", workflowToken(t, app, workflowAuthor), requestActionData(t, app, id)), http.StatusBadRequest)
	long := strings.Repeat("x", 501)
	proposeRates(t, app, id, workflowAuthor, proposedRates("Client Negotiated Rates", long, standardRates(t, app, 260)), http.StatusBadRequest)
	proposeRates(t, app, id, workflowAuthor, proposedRates("Client Negotiated Rates", "Client MSA, schedule B", standardRates(t, app, 260)), http.StatusOK)
	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/submit", workflowAuthor, requestActionData(t, app, id), http.StatusOK)
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
	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/approve", workflowApprover, requestActionData(t, app, id), http.StatusOK)

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

	// The approver's rate sheet statement is about the proposed rates.
	detail := clientSetupJSON(t, app, http.MethodGet, "/api/job_requests/"+id, workflowApprover, nil, http.StatusOK)
	found := false
	for _, question := range detail["checklist"].(map[string]any)["questions"].([]any) {
		q := question.(map[string]any)
		if q["key"] == "rate_sheet_correct" {
			found = q["text"] == "The proposed rates are approved for this client."
		}
	}
	if !found {
		t.Fatalf("approval checklist: %#v", detail["checklist"])
	}
	if detail["related"].(map[string]any)["rate_sheet_proposed"] != true {
		t.Fatal("the approver does not see the proposal")
	}
	// Only the selected approver approves, and a job holder can't create first.
	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/create", workflowAuthor, requestActionData(t, app, id), http.StatusBadRequest)
	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/approve", workflowAuthor, requestActionData(t, app, id), http.StatusForbidden)
	clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+id+"/approve", workflowApprover, requestActionData(t, app, id), http.StatusOK)
	if workflowRecord(t, app, "rate_sheets", draftID).GetBool("active") {
		t.Fatal("approval activated the draft")
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
		sheet.GetString("effective_date") != hooks.RateSheetToday() || sheet.GetString("based_on") != standardRateSheet {
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
