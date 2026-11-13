package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"strings"
	"testing"

	"tybalt/absorb"
	"tybalt/internal/testutils"
	"tybalt/routes"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

// Rate sheet absorb (issue #191). The fixtures are rows of the seed CSV files
// (app/test_seed_data/data): the rate sheets below with their entries, the
// jobs and the request that use them, one earlier name of Absorb Single
// (rate_sheet_name_changes.csv), and the user absorb.revise@example.com.
//
// Rate sets, each made from the seeded 2025 Standard Rates sheet
// (c41ofep525bcacj), which has an entry for every role:
//   - rates A: 1000 added to each rate, so that it matches no other sheet
//   - rates B: rates A without one role
//   - rates C: rates A with one overtime rate 1 higher
//   - rates G: 3000 added to each rate (the Group sheets)
//
// Seeded users:
//   - author@soup.com holds admin, not absorb
//   - book@keeper.com holds absorb, not admin or rate_sheet_revise
//   - absorb.revise@example.com holds absorb and rate_sheet_revise, not admin
//   - time@test.com holds rate_sheet_revise, not admin or absorb
//
// Two tests change fixture data at test time, and each says why at the
// change: TestRateSheetAbsorbKeepsJobWIP and
// TestRateSheetAbsorbRequestBecomesJob.

const (
	absorbAdminEmail  = "author@soup.com"
	absorbOnlyEmail   = "book@keeper.com"
	absorbReviseEmail = "absorb.revise@example.com"
	reviseOnlyEmail   = "time@test.com"

	rsAbTarget       = "rsabtarget00001" // Absorb Target, rev 0, 2024-01-01, active, rates A
	rsAbSingle       = "rsabsingle00001" // Absorb Single, rev 0, 2024-06-01, active, rates A; a job and a request use it
	rsAbEarly        = "rsabearly000001" // Absorb Early, rev 0, 2023-01-01, inactive, rates A
	rsAbMissing      = "rsabmissing0001" // Absorb Missing Role, rates A without one role
	rsAbDiffRate     = "rsabdiffrate001" // Absorb Different Rate, rates C (one rate differs)
	rsAbWhole0       = "rsabwhole0r0001" // Absorb Whole rev 0, 2022-01-01, inactive, rates C
	rsAbWhole1       = "rsabwhole0r1001" // Absorb Whole rev 1, 2022-06-01, active, rates A; no references
	rsAbUsed0        = "rsabused0r00001" // Absorb Used rev 0, 2024-01-01, inactive, rates A; a job uses it
	rsAbUsed1        = "rsabused0r10001" // Absorb Used rev 1, 2024-02-01, active, rates A; nothing uses it
	rsAbLateInactive = "rsablateinact01" // Absorb Late Inactive, rev 0, 2025-01-01, inactive, rates A
	rsAbLatest0      = "rsablatest0r001" // Absorb Latest Differs rev 0, inactive, rates A
	rsAbLatest1      = "rsablatest1r001" // Absorb Latest Differs rev 1, active, rates C; no references
	rsAbInactiveTgt  = "rsabinactive001" // Absorb Inactive Target, 2020-01-01, inactive, rates A
	rsAbSingleJob    = "jobabsingle0001"
	rsAbUsedJob      = "jobabused000001"
	rsAbSingleReq    = "parequestabs001"
	rsAbStandardRate = "c41ofep525bcacj"
)

func rateSheetAbsorbSQL(t testing.TB, app core.App, query string, params dbx.Params) {
	t.Helper()
	if _, err := app.DB().NewQuery(query).Bind(params).Execute(); err != nil {
		t.Fatalf("%s: %v", query, err)
	}
}

// newRateSheetAbsorbApp returns a seeded app. The fixtures listed in the
// constants above are rows of the seed CSV files.
func newRateSheetAbsorbApp(t testing.TB) *tests.TestApp {
	t.Helper()
	return testutils.SetupTestApp(t)
}

// absorbRequest sends a request as the user with email and returns the
// status and body.
func absorbRequest(t *testing.T, app *tests.TestApp, method, path, email string, data any) (int, string) {
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
	return response.Code, response.Body.String()
}

func absorbBody(ids ...string) map[string]any {
	return map[string]any{"ids_to_absorb": ids}
}

func rateSheetExists(t testing.TB, app core.App, id string) bool {
	t.Helper()
	_, err := app.FindRecordById("rate_sheets", id)
	return err == nil
}

// entrySnapshot returns the id, role, rate and overtime rate of each entry
// of a sheet, sorted by id.
func entrySnapshot(t testing.TB, app core.App, sheet string) []string {
	t.Helper()
	records, err := app.FindAllRecords("rate_sheet_entries", dbx.HashExp{"rate_sheet": sheet})
	if err != nil {
		t.Fatal(err)
	}
	result := []string{}
	for _, r := range records {
		raw, _ := json.Marshal([]any{r.Id, r.GetString("role"), r.GetFloat("rate"), r.GetFloat("overtime_rate")})
		result = append(result, string(raw))
	}
	sort.Strings(result)
	return result
}

func fieldOf(t testing.TB, app core.App, collection, id, field string) string {
	t.Helper()
	record, err := app.FindRecordById(collection, id)
	if err != nil {
		t.Fatal(err)
	}
	return record.GetString(field)
}

func rateSheetAction(t testing.TB, app core.App) *core.Record {
	t.Helper()
	action, err := app.FindFirstRecordByData("absorb_actions", "collection_name", "rate_sheets")
	if err != nil {
		t.Fatalf("no pending rate sheet absorb: %v", err)
	}
	return action
}

// Every Decided rule that refuses an absorb, checked in AbsorbRecords so
// that nothing changes, and the happy paths that the same rules allow.
func TestRateSheetAbsorbRules(t *testing.T) {
	cases := []struct {
		name        string
		target      string
		ids         []string
		refusal     string   // part of the refusal message; "" when allowed
		wantDeleted []string // sheets that an allowed absorb deletes
	}{
		{"a missing role is a difference", rsAbTarget, []string{rsAbMissing}, "rates are not identical", nil},
		{"a different rate is a difference", rsAbTarget, []string{rsAbDiffRate}, "rates are not identical", nil},
		{"the target's date can't be later than an absorbed sheet's", rsAbTarget, []string{rsAbEarly}, "effective date, 2024-01-01, is later than this sheet's, 2023-01-01", nil},
		{"no revision of a name that a job uses: the used one", rsAbTarget, []string{rsAbUsed0}, "a job or project request uses one of them", nil},
		{"no revision of a name that a job uses: an unused one", rsAbTarget, []string{rsAbUsed1}, "a job or project request uses one of them", nil},
		{"the latest revision must match", rsAbTarget, []string{rsAbLatest0}, "latest revision's rates are not identical", nil},
		{"an inactive target can't absorb an active sheet", rsAbInactiveTgt, []string{rsAbSingle}, "inactive target can't absorb an active sheet", nil},
		{"an inactive target can't absorb a name with an active revision", rsAbInactiveTgt, []string{rsAbWhole0}, "inactive target can't absorb an active sheet", nil},
		{"a revision of the target's own name", rsAbWhole1, []string{rsAbWhole0}, "revision of the target sheet", nil},
		{"one refused sheet refuses the whole request", rsAbTarget, []string{rsAbSingle, rsAbMissing}, "rates are not identical", nil},

		{"a single revision sheet with a later date", rsAbTarget, []string{rsAbSingle}, "", []string{rsAbSingle}},
		{"an inactive target can absorb an inactive sheet", rsAbInactiveTgt, []string{rsAbEarly}, "", []string{rsAbEarly}},
		{"an older revision of a used name can be the target", rsAbUsed0, []string{rsAbLateInactive}, "", []string{rsAbLateInactive}},
		{"the latest revision of a used name can be the target", rsAbUsed1, []string{rsAbSingle}, "", []string{rsAbSingle}},
		{"an unused name is absorbed whole, without the date rule", rsAbTarget, []string{rsAbWhole0}, "", []string{rsAbWhole0, rsAbWhole1}},
		{"selecting the latest revision absorbs the whole name too", rsAbTarget, []string{rsAbWhole1}, "", []string{rsAbWhole0, rsAbWhole1}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			app := newRateSheetAbsorbApp(t)
			defer app.Cleanup()

			err := routes.AbsorbRecords(app, "rate_sheets", tc.target, tc.ids)
			if tc.refusal != "" {
				var refused *absorb.RefusedError
				if !errors.As(err, &refused) {
					t.Fatalf("expected a refusal, got %v", err)
				}
				if !strings.Contains(refused.Message, tc.refusal) {
					t.Fatalf("refusal %q does not contain %q", refused.Message, tc.refusal)
				}
				for _, id := range tc.ids {
					if !rateSheetExists(t, app, id) {
						t.Fatalf("a refused absorb deleted %s", id)
					}
				}
				if _, err := app.FindFirstRecordByData("absorb_actions", "collection_name", "rate_sheets"); err == nil {
					t.Fatal("a refused absorb left an action")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			for _, id := range tc.wantDeleted {
				if rateSheetExists(t, app, id) {
					t.Fatalf("expected %s to be absorbed", id)
				}
				if len(entrySnapshot(t, app, id)) != 0 {
					t.Fatalf("the entries of %s were not deleted", id)
				}
			}
			var absorbed []map[string]any
			if err := rateSheetAction(t, app).UnmarshalJSONField("absorbed_records", &absorbed); err != nil {
				t.Fatal(err)
			}
			if len(absorbed) != len(tc.wantDeleted) {
				t.Fatalf("the action holds %d sheets; want %d", len(absorbed), len(tc.wantDeleted))
			}
		})
	}
}

// Absorb and undo restore the sheet, its entries with their ids, and the job
// and request references.
func TestRateSheetAbsorbUndoRoundTrip(t *testing.T) {
	app := newRateSheetAbsorbApp(t)
	defer app.Cleanup()

	entriesBefore := entrySnapshot(t, app, rsAbSingle)
	if len(entriesBefore) == 0 {
		t.Fatal("the fixture has no entries")
	}
	targetEntries := entrySnapshot(t, app, rsAbTarget)

	status, body := absorbRequest(t, app, http.MethodPost, "/api/rate_sheets/"+rsAbTarget+"/absorb", absorbAdminEmail, absorbBody(rsAbSingle))
	if status != http.StatusOK {
		t.Fatalf("absorb status %d: %s", status, body)
	}
	if rateSheetExists(t, app, rsAbSingle) || len(entrySnapshot(t, app, rsAbSingle)) != 0 {
		t.Fatal("the absorbed sheet or its entries remain")
	}
	if got := fieldOf(t, app, "jobs", rsAbSingleJob, "rate_sheet"); got != rsAbTarget {
		t.Fatalf("job rate sheet = %s; want the target", got)
	}
	if got := fieldOf(t, app, "job_requests", rsAbSingleReq, "rate_sheet"); got != rsAbTarget {
		t.Fatalf("request rate sheet = %s; want the target", got)
	}
	// The moved job is written back to the legacy system.
	job, err := app.FindRecordById("jobs", rsAbSingleJob)
	if err != nil {
		t.Fatal(err)
	}
	if job.GetBool("_imported") {
		t.Fatal("the moved job was not marked for writeback")
	}
	if strings.Join(entrySnapshot(t, app, rsAbTarget), "|") != strings.Join(targetEntries, "|") {
		t.Fatal("the absorb changed the target's entries")
	}

	status, body = absorbRequest(t, app, http.MethodPost, "/api/rate_sheets/undo_absorb", absorbAdminEmail, nil)
	if status != http.StatusOK {
		t.Fatalf("undo status %d: %s", status, body)
	}
	if !rateSheetExists(t, app, rsAbSingle) {
		t.Fatal("undo did not restore the sheet")
	}
	if got := entrySnapshot(t, app, rsAbSingle); strings.Join(got, "|") != strings.Join(entriesBefore, "|") {
		t.Fatalf("undo restored entries %v; want %v", got, entriesBefore)
	}
	sheet, err := app.FindRecordById("rate_sheets", rsAbSingle)
	if err != nil {
		t.Fatal(err)
	}
	if sheet.GetString("name") != "Absorb Single" || !sheet.GetBool("active") || sheet.GetString("effective_date") != "2024-06-01" {
		t.Fatalf("undo restored a different sheet: %v", sheet.PublicExport())
	}
	if got := fieldOf(t, app, "jobs", rsAbSingleJob, "rate_sheet"); got != rsAbSingle {
		t.Fatalf("undo left the job on %s", got)
	}
	if got := fieldOf(t, app, "job_requests", rsAbSingleReq, "rate_sheet"); got != rsAbSingle {
		t.Fatalf("undo left the request on %s", got)
	}
	if _, err := app.FindFirstRecordByData("absorb_actions", "collection_name", "rate_sheets"); err == nil {
		t.Fatal("undo left the action")
	}
}

// A name with revisions is absorbed, undone and committed as a whole, and its
// name goes into the history once.
func TestRateSheetAbsorbWholeName(t *testing.T) {
	app := newRateSheetAbsorbApp(t)
	defer app.Cleanup()

	before := map[string][]string{
		rsAbWhole0: entrySnapshot(t, app, rsAbWhole0),
		rsAbWhole1: entrySnapshot(t, app, rsAbWhole1),
	}

	// The older revision is selected; every revision goes.
	if status, body := absorbRequest(t, app, http.MethodPost, "/api/rate_sheets/"+rsAbTarget+"/absorb", absorbReviseEmail, absorbBody(rsAbWhole0)); status != http.StatusOK {
		t.Fatalf("absorb status %d: %s", status, body)
	}
	for id := range before {
		if rateSheetExists(t, app, id) {
			t.Fatalf("revision %s was not absorbed", id)
		}
	}
	if status, body := absorbRequest(t, app, http.MethodPost, "/api/rate_sheets/undo_absorb", absorbReviseEmail, nil); status != http.StatusOK {
		t.Fatalf("undo status %d: %s", status, body)
	}
	for id, entries := range before {
		if !rateSheetExists(t, app, id) {
			t.Fatalf("undo did not restore revision %s", id)
		}
		if got := entrySnapshot(t, app, id); strings.Join(got, "|") != strings.Join(entries, "|") {
			t.Fatalf("undo restored other entries for %s", id)
		}
	}

	if status, body := absorbRequest(t, app, http.MethodPost, "/api/rate_sheets/"+rsAbTarget+"/absorb", absorbReviseEmail, absorbBody(rsAbWhole1)); status != http.StatusOK {
		t.Fatalf("absorb status %d: %s", status, body)
	}
	action := rateSheetAction(t, app)
	if status, body := absorbRequest(t, app, http.MethodDelete, "/api/collections/absorb_actions/records/"+action.Id, absorbReviseEmail, nil); status != http.StatusNoContent {
		t.Fatalf("commit status %d: %s", status, body)
	}
	changes, err := app.FindAllRecords("rate_sheet_name_changes", dbx.HashExp{"rate_sheet": rsAbTarget})
	if err != nil {
		t.Fatal(err)
	}
	if len(changes) != 1 || changes[0].GetString("old_name") != "Absorb Whole" || changes[0].GetString("new_name") != "Absorb Target" {
		t.Fatalf("name history = %d rows; want one row for Absorb Whole", len(changes))
	}
}

// Commit writes each absorbed name to the history, moves the absorbed
// sheet's own history to the target, and needs the same claims as absorb.
func TestRateSheetAbsorbCommitWritesNameHistory(t *testing.T) {
	app := newRateSheetAbsorbApp(t)
	defer app.Cleanup()

	// The absorbed sheet was renamed before: rate_sheet_name_changes.csv has
	// its earlier name, Absorb Single Before.
	if err := routes.AbsorbRecords(app, "rate_sheets", rsAbTarget, []string{rsAbSingle}); err != nil {
		t.Fatal(err)
	}
	// Nothing is written on absorb.
	if names, err := routes.RateSheetFormerNames(app, "Absorb Target"); err != nil || len(names) != 0 {
		t.Fatalf("former names after absorb = %v, %v; want none", names, err)
	}
	action := rateSheetAction(t, app)
	path := "/api/collections/absorb_actions/records/" + action.Id

	// absorb alone passes the collection rule and is refused by the claim
	// check. The admin without absorb can see the rate sheet action.
	if status, body := absorbRequest(t, app, http.MethodDelete, path, absorbOnlyEmail, nil); status != http.StatusForbidden {
		t.Fatalf("commit with absorb only: status %d: %s", status, body)
	}
	if status, body := absorbRequest(t, app, http.MethodGet, path, absorbAdminEmail, nil); status != http.StatusOK {
		t.Fatalf("view by an admin without absorb: status %d: %s", status, body)
	}
	if names, _ := routes.RateSheetFormerNames(app, "Absorb Target"); len(names) != 0 {
		t.Fatalf("a refused commit wrote history: %v", names)
	}

	if status, body := absorbRequest(t, app, http.MethodDelete, path, absorbReviseEmail, nil); status != http.StatusNoContent {
		t.Fatalf("commit status %d: %s", status, body)
	}
	names, err := routes.RateSheetFormerNames(app, "Absorb Target")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(names, "|") != "Absorb Single|Absorb Single Before" {
		t.Fatalf("former names = %v; want the absorbed name, then its earlier name", names)
	}
	committer, err := app.FindAuthRecordByEmail("users", absorbReviseEmail)
	if err != nil {
		t.Fatal(err)
	}
	change, err := app.FindFirstRecordByData("rate_sheet_name_changes", "old_name", "Absorb Single")
	if err != nil {
		t.Fatal(err)
	}
	if change.GetString("renamed_by") != committer.Id || change.GetString("rate_sheet") != rsAbTarget {
		t.Fatalf("history row = %v", change.PublicExport())
	}
	if rateSheetExists(t, app, rsAbSingle) {
		t.Fatal("commit restored the absorbed sheet")
	}
}

// Absorb and undo need admin, or both absorb and rate_sheet_revise. absorb
// alone is not enough for rate sheets, and is still enough for vendors.
func TestRateSheetAbsorbClaims(t *testing.T) {
	cases := []struct {
		name   string
		email  string
		status int
	}{
		{"admin", absorbAdminEmail, http.StatusOK},
		{"absorb and rate_sheet_revise", absorbReviseEmail, http.StatusOK},
		{"absorb alone", absorbOnlyEmail, http.StatusForbidden},
		{"rate_sheet_revise alone", reviseOnlyEmail, http.StatusForbidden},
		{"signed out", "", http.StatusUnauthorized},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			app := newRateSheetAbsorbApp(t)
			defer app.Cleanup()
			status, body := absorbRequest(t, app, http.MethodPost, "/api/rate_sheets/"+rsAbTarget+"/absorb", tc.email, absorbBody(rsAbSingle))
			if status != tc.status {
				t.Fatalf("absorb status %d; want %d: %s", status, tc.status, body)
			}
			if tc.status == http.StatusOK {
				if !strings.Contains(body, "Successfully absorbed 1 records") {
					t.Fatalf("absorb body: %s", body)
				}
				// Undo is refused for absorb alone and allowed for the user who
				// absorbed.
				if status, _ := absorbRequest(t, app, http.MethodPost, "/api/rate_sheets/undo_absorb", absorbOnlyEmail, nil); status != http.StatusForbidden {
					t.Fatalf("undo with absorb alone: status %d", status)
				}
				if status, body := absorbRequest(t, app, http.MethodPost, "/api/rate_sheets/undo_absorb", tc.email, nil); status != http.StatusOK {
					t.Fatalf("undo status %d: %s", status, body)
				}
			} else if !rateSheetExists(t, app, rsAbSingle) {
				t.Fatal("a refused absorb deleted the sheet")
			}
		})
	}

	// Other collections keep the absorb-only check.
	app := newRateSheetAbsorbApp(t)
	defer app.Cleanup()
	if status, body := absorbRequest(t, app, http.MethodPost, "/api/clients/lb0fnenkeyitsny/absorb", absorbAdminEmail, absorbBody("eldtxi3i4h00k8r")); status != http.StatusForbidden {
		t.Fatalf("client absorb by an admin without absorb: status %d: %s", status, body)
	}
}

// An admin without the absorb claim can see and commit a rate sheet action.
// The collection rule still hides the actions of the other collections from
// that admin, and refuses to commit them.
func TestRateSheetAbsorbAdminWithoutAbsorbClaim(t *testing.T) {
	app := newRateSheetAbsorbApp(t)
	defer app.Cleanup()

	const listPath = "/api/collections/absorb_actions/records"
	countActions := func(email string) int {
		t.Helper()
		status, body := absorbRequest(t, app, http.MethodGet, listPath, email, nil)
		if status != http.StatusOK {
			t.Fatalf("list as %s: status %d: %s", email, status, body)
		}
		var page struct {
			Items []map[string]any `json:"items"`
		}
		if err := json.Unmarshal([]byte(body), &page); err != nil {
			t.Fatal(err)
		}
		return len(page.Items)
	}

	// A client action, made by an absorb holder, is hidden from the admin.
	if status, body := absorbRequest(t, app, http.MethodPost, "/api/clients/lb0fnenkeyitsny/absorb", absorbOnlyEmail, absorbBody("eldtxi3i4h00k8r")); status != http.StatusOK {
		t.Fatalf("client absorb: status %d: %s", status, body)
	}
	clientAction, err := app.FindFirstRecordByData("absorb_actions", "collection_name", "clients")
	if err != nil {
		t.Fatal(err)
	}
	clientPath := listPath + "/" + clientAction.Id
	if n := countActions(absorbAdminEmail); n != 0 {
		t.Fatalf("the admin sees %d actions; want none before a rate sheet absorb", n)
	}
	if status, _ := absorbRequest(t, app, http.MethodGet, clientPath, absorbAdminEmail, nil); status != http.StatusNotFound {
		t.Fatalf("view of a client action by the admin: status %d; want 404", status)
	}
	if status, _ := absorbRequest(t, app, http.MethodDelete, clientPath, absorbAdminEmail, nil); status != http.StatusNotFound {
		t.Fatalf("commit of a client action by the admin: status %d; want 404", status)
	}

	// The admin absorbs a rate sheet, sees only that action and commits it.
	if status, body := absorbRequest(t, app, http.MethodPost, "/api/rate_sheets/"+rsAbTarget+"/absorb", absorbAdminEmail, absorbBody(rsAbSingle)); status != http.StatusOK {
		t.Fatalf("absorb status %d: %s", status, body)
	}
	if n := countActions(absorbAdminEmail); n != 1 {
		t.Fatalf("the admin sees %d actions; want only the rate sheet action", n)
	}
	if n := countActions(absorbOnlyEmail); n != 2 {
		t.Fatalf("the absorb holder sees %d actions; want both", n)
	}
	path := listPath + "/" + rateSheetAction(t, app).Id
	if status, body := absorbRequest(t, app, http.MethodDelete, path, absorbAdminEmail, nil); status != http.StatusNoContent {
		t.Fatalf("commit by the admin: status %d: %s", status, body)
	}
	if _, err := app.FindRecordById("absorb_actions", clientAction.Id); err != nil {
		t.Fatalf("the client action is gone: %v", err)
	}
}

// The route answers a refused rule with 400 and the reason.
func TestRateSheetAbsorbRouteRefusal(t *testing.T) {
	app := newRateSheetAbsorbApp(t)
	defer app.Cleanup()
	status, body := absorbRequest(t, app, http.MethodPost, "/api/rate_sheets/"+rsAbTarget+"/absorb", absorbAdminEmail, absorbBody(rsAbMissing))
	if status != http.StatusBadRequest || !strings.Contains(body, "rates are not identical") {
		t.Fatalf("status %d: %s", status, body)
	}
}

// A moved job's WIP is the same before the absorb, after it and after undo.
func TestRateSheetAbsorbKeepsJobWIP(t *testing.T) {
	app := testutils.SetupTestApp(t)
	defer app.Cleanup()

	const job = "jobwip000000001"
	// WIP Source and WIP Target are seeded single revision sheets with the
	// rates of the job's own sheet (rsvaluation0001, a revision of a name
	// that jobs use, so it can't be absorbed itself).
	//
	// Test-time change to fixture data: the job is moved to WIP Source here,
	// with its opening rate sheet period (job_rate_sheets.csv), which is what
	// prices its time. A seeded job on WIP Source would need its own time
	// entries, and those would change the totals that the WIP and time sheet
	// tests assert.
	rateSheetAbsorbSQL(t, app, "UPDATE jobs SET rate_sheet = 'rsabwipsource01' WHERE id = {:job}", dbx.Params{"job": job})
	rateSheetAbsorbSQL(t, app, "UPDATE job_rate_sheets SET rate_sheet = 'rsabwipsource01' WHERE job = {:job}", dbx.Params{"job": job})
	periodSheet := func() string {
		t.Helper()
		row, err := app.FindFirstRecordByData("job_rate_sheets", "job", job)
		if err != nil {
			t.Fatal(err)
		}
		return row.GetString("rate_sheet")
	}

	wip := func() routes.JobWIP {
		t.Helper()
		status, body := absorbRequest(t, app, http.MethodGet, "/api/jobs/"+job+"/wip", "valuation1@example.com", nil)
		if status != http.StatusOK {
			t.Fatalf("WIP status %d: %s", status, body)
		}
		var result routes.JobWIP
		if err := json.Unmarshal([]byte(body), &result); err != nil {
			t.Fatal(err)
		}
		return result
	}

	before := wip()
	// The seeded value of this job (TestJobWIP), so the copy prices as before.
	if before.TimeValue != 849.5 || before.Hours != 3.5 || before.MissingRoleHours != .5 {
		t.Fatalf("WIP before = %+v", before)
	}
	if err := routes.AbsorbRecords(app, "rate_sheets", "rsabwiptarget01", []string{"rsabwipsource01"}); err != nil {
		t.Fatal(err)
	}
	if got := fieldOf(t, app, "jobs", job, "rate_sheet"); got != "rsabwiptarget01" {
		t.Fatalf("the job is on %s", got)
	}
	if got := periodSheet(); got != "rsabwiptarget01" {
		t.Fatalf("the job's period is on %s", got)
	}
	if after := wip(); after != before {
		t.Fatalf("WIP after absorb = %+v; want %+v", after, before)
	}
	if status, body := absorbRequest(t, app, http.MethodPost, "/api/rate_sheets/undo_absorb", absorbAdminEmail, nil); status != http.StatusOK {
		t.Fatalf("undo status %d: %s", status, body)
	}
	if got := fieldOf(t, app, "jobs", job, "rate_sheet"); got != "rsabwipsource01" {
		t.Fatalf("undo left the job on %s", got)
	}
	if got := periodSheet(); got != "rsabwipsource01" {
		t.Fatalf("undo left the job's period on %s", got)
	}
	if undone := wip(); undone != before {
		t.Fatalf("WIP after undo = %+v; want %+v", undone, before)
	}
}

// A request that becomes a job while the absorb is pending carries the
// pending undo data to the job, as for clients.
func TestRateSheetAbsorbRequestBecomesJob(t *testing.T) {
	app := newWorkflowTestApp(t, true)
	const requestID = "parequest000001"
	// Request Target is a seeded sheet with the standard rates.
	//
	// Test-time change to fixture data: Request Target is activated here. The
	// request of this workflow fixture is on the active standard sheet, and
	// only an active target can absorb an active sheet. Request Target is
	// inactive in the seed, because an active copy of the standard sheet
	// would be a duplicate for every test of the duplicate check.
	rateSheetAbsorbSQL(t, app, "UPDATE rate_sheets SET active = true WHERE id = 'rsabreqtarget01'", nil)

	if err := routes.AbsorbRecords(app, "rate_sheets", "rsabreqtarget01", []string{rsAbStandardRate}); err != nil {
		t.Fatal(err)
	}
	if got := fieldOf(t, app, "job_requests", requestID, "rate_sheet"); got != "rsabreqtarget01" {
		t.Fatalf("the request is on %s", got)
	}
	approveMergeFixtureRequest(t, app)
	result := clientSetupJSON(t, app, http.MethodPost, "/api/job_requests/"+requestID+"/create", workflowAuthor, requestActionData(t, app, requestID), http.StatusOK)
	jobID := result["job_id"].(string)
	if got := fieldOf(t, app, "jobs", jobID, "rate_sheet"); got != "rsabreqtarget01" {
		t.Fatalf("the new job is on %s", got)
	}

	var updates map[string]map[string]map[string]string
	if err := rateSheetAction(t, app).UnmarshalJSONField("updated_references", &updates); err != nil {
		t.Fatal(err)
	}
	if updates["jobs"]["rate_sheet"][jobID] != rsAbStandardRate {
		t.Fatal("the pending undo data did not move to the new job")
	}
	if _, tracked := updates["job_requests"]["rate_sheet"][requestID]; tracked {
		t.Fatal("the deleted request is still in the undo data")
	}
	// The new project's opening period was written on the target with the
	// job, so its undo data moves too.
	opening, err := app.FindFirstRecordByFilter("job_rate_sheets", "job = {:job} && effective_date = ''", dbx.Params{"job": jobID})
	if err != nil {
		t.Fatalf("the new job has no opening period: %v", err)
	}
	if opening.GetString("rate_sheet") != "rsabreqtarget01" {
		t.Fatalf("the opening period is on %s", opening.GetString("rate_sheet"))
	}
	if updates["job_rate_sheets"]["rate_sheet"][opening.Id] != rsAbStandardRate {
		t.Fatal("the pending undo data does not cover the new job's opening period")
	}

	if status, body := absorbRequest(t, app, http.MethodPost, "/api/rate_sheets/undo_absorb", absorbAdminEmail, nil); status != http.StatusOK {
		t.Fatalf("undo status %d: %s", status, body)
	}
	if got := fieldOf(t, app, "jobs", jobID, "rate_sheet"); got != rsAbStandardRate {
		t.Fatalf("undo left the new job on %s", got)
	}
	if got := fieldOf(t, app, "job_rate_sheets", opening.Id, "rate_sheet"); got != rsAbStandardRate {
		t.Fatalf("undo left the new job's opening period on %s", got)
	}
}

// The Duplicates view groups identical sheets, active and inactive, marks the
// earliest as the target, and says why a sheet can't be absorbed.
func TestRateSheetDuplicatesView(t *testing.T) {
	app := newRateSheetAbsorbApp(t)
	defer app.Cleanup()

	status, body := absorbRequest(t, app, http.MethodGet, "/api/rate_sheets/duplicates", reviseOnlyEmail, nil)
	if status != http.StatusOK {
		t.Fatalf("status %d: %s", status, body)
	}
	var result struct {
		Groups []struct {
			Sheets []struct {
				ID        string `json:"id"`
				Name      string `json:"name"`
				Target    bool   `json:"target"`
				CanAbsorb bool   `json:"can_absorb"`
				Reason    string `json:"reason"`
				Absorbs   int    `json:"absorbs"`
				JobCount  int    `json:"job_count"`
				Active    bool   `json:"active"`
			} `json:"sheets"`
			Reason string `json:"reason"`
		} `json:"groups"`
	}
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		t.Fatal(err)
	}

	// The seeded standard sheet and its inactive copy are one group.
	foundSeeded := false
	for _, group := range result.Groups {
		ids := map[string]bool{}
		for _, sheet := range group.Sheets {
			ids[sheet.ID] = true
		}
		if ids[rsAbStandardRate] && ids["rs_dupact_001"] {
			foundSeeded = true
		}
		if !ids[rsAbTarget] {
			continue
		}
		// Rates A: the earliest sheet is the inactive target of 2020.
		first := group.Sheets[0]
		if first.ID != rsAbInactiveTgt || !first.Target {
			t.Fatalf("target = %s; want %s", first.ID, rsAbInactiveTgt)
		}
		want := map[string]string{
			rsAbEarly:        "",
			rsAbLateInactive: "",
			rsAbTarget:       "inactive target",
			rsAbSingle:       "inactive target",
			rsAbWhole1:       "inactive target",
			rsAbUsed0:        "a job or project request uses",
			rsAbUsed1:        "a job or project request uses",
			rsAbLatest0:      "latest revision's rates",
		}
		for _, sheet := range group.Sheets[1:] {
			wantReason, ok := want[sheet.ID]
			if !ok {
				t.Fatalf("unexpected sheet %s in the group", sheet.ID)
			}
			if wantReason == "" && sheet.Reason != "" || !strings.Contains(sheet.Reason, wantReason) {
				t.Fatalf("%s reason = %q; want %q", sheet.ID, sheet.Reason, wantReason)
			}
			delete(want, sheet.ID)
		}
		if len(want) != 0 {
			t.Fatalf("missing sheets %v", want)
		}
		for _, sheet := range group.Sheets {
			// The active target can absorb the single sheet; the used name
			// can be a target but can't be absorbed.
			if sheet.ID == rsAbTarget && !sheet.CanAbsorb {
				t.Fatal("the active target should be able to absorb")
			}
			if sheet.ID == rsAbSingle && sheet.JobCount != 1 {
				t.Fatalf("job count = %d", sheet.JobCount)
			}
		}
		if group.Reason != "" {
			t.Fatalf("the group has an absorbable sheet but says %q", group.Reason)
		}
	}
	if !foundSeeded {
		t.Fatal("the seeded duplicate group is missing")
	}

	// The seeded Group sheets: the earliest sheet (Group Old rev 0, which a
	// job uses) is inactive and has a newer revision, and Group New is
	// active. The group says why nothing can be absorbed.
	if !strings.Contains(body, "has a newer revision") {
		t.Fatalf("expected the Group sheets to say why they can't be absorbed: %s", body)
	}
}

// The absorb page lists the target's identical sheets with the reasons.
func TestRateSheetAbsorbCandidates(t *testing.T) {
	app := newRateSheetAbsorbApp(t)
	defer app.Cleanup()
	status, body := absorbRequest(t, app, http.MethodGet, "/api/rate_sheets/"+rsAbTarget+"/absorb_candidates", reviseOnlyEmail, nil)
	if status != http.StatusOK {
		t.Fatalf("status %d: %s", status, body)
	}
	var result struct {
		Target     struct{ ID string } `json:"target"`
		Candidates []struct {
			ID      string `json:"id"`
			Reason  string `json:"reason"`
			Absorbs int    `json:"absorbs"`
		} `json:"candidates"`
	}
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		t.Fatal(err)
	}
	got := map[string]string{}
	absorbs := map[string]int{}
	for _, c := range result.Candidates {
		got[c.ID] = c.Reason
		absorbs[c.ID] = c.Absorbs
	}
	if got[rsAbSingle] != "" || absorbs[rsAbSingle] != 1 {
		t.Fatalf("single sheet: %q, %d", got[rsAbSingle], absorbs[rsAbSingle])
	}
	if got[rsAbWhole1] != "" || absorbs[rsAbWhole1] != 2 {
		t.Fatalf("whole name: %q, %d", got[rsAbWhole1], absorbs[rsAbWhole1])
	}
	if !strings.Contains(got[rsAbEarly], "effective date") {
		t.Fatalf("early sheet: %q", got[rsAbEarly])
	}
	if _, listed := got[rsAbMissing]; listed {
		t.Fatal("a sheet with other rates is listed")
	}
	if status, _ := absorbRequest(t, app, http.MethodGet, "/api/rate_sheets/missingsheet001/absorb_candidates", reviseOnlyEmail, nil); status != http.StatusNotFound {
		t.Fatalf("unknown sheet: status %d", status)
	}
}
