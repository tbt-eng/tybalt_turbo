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

// Rate sheet absorb (issue #191). The fixtures are made in each test from the
// seeded 2025 Standard Rates sheet (c41ofep525bcacj), which has an entry for
// every role. Seeded users:
//   - author@soup.com holds admin, not absorb
//   - book@keeper.com holds absorb, not admin or rate_sheet_revise
//   - time@test.com holds rate_sheet_revise, not admin or absorb
//
// absorbReviseEmail is book@keeper.com after a test gives it
// rate_sheet_revise, with grantRateSheetRevise.

const (
	absorbAdminEmail  = "author@soup.com"
	absorbOnlyEmail   = "book@keeper.com"
	absorbReviseEmail = "book@keeper.com"
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
	rsAbSingleJob    = "jobwip000000004"
	rsAbUsedJob      = "jobwip000000003"
	rsAbSingleReq    = "parequest000001"
	rsAbStandardRate = "c41ofep525bcacj"
)

type rateSheetEntryFixture struct {
	Role         string
	Rate         float64
	OvertimeRate float64
}

// standardEntries returns the entries of the seeded standard sheet, sorted by
// role.
func standardEntries(t testing.TB, app core.App) []rateSheetEntryFixture {
	t.Helper()
	records, err := app.FindAllRecords("rate_sheet_entries", dbx.HashExp{"rate_sheet": rsAbStandardRate})
	if err != nil {
		t.Fatal(err)
	}
	entries := []rateSheetEntryFixture{}
	for _, r := range records {
		entries = append(entries, rateSheetEntryFixture{r.GetString("role"), r.GetFloat("rate"), r.GetFloat("overtime_rate")})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Role < entries[j].Role })
	return entries
}

// ratesA is the standard sheet with 1000 added to each rate, so that it
// matches no seeded sheet. ratesB lacks one role of ratesA. ratesC differs
// from ratesA in one overtime rate.
func absorbRateSets(t testing.TB, app core.App) (a, b, c []rateSheetEntryFixture) {
	t.Helper()
	for _, e := range standardEntries(t, app) {
		a = append(a, rateSheetEntryFixture{e.Role, e.Rate + 1000, e.OvertimeRate + 1000})
	}
	b = append(b, a[1:]...)
	c = append(c, a...)
	c[0].OvertimeRate++
	return a, b, c
}

func createRateSheetFixture(t testing.TB, app core.App, id, name string, revision int, effectiveDate string, active bool, entries []rateSheetEntryFixture) {
	t.Helper()
	sheets, err := app.FindCollectionByNameOrId("rate_sheets")
	if err != nil {
		t.Fatal(err)
	}
	sheet := core.NewRecord(sheets)
	sheet.Id = id
	sheet.Set("name", name)
	sheet.Set("revision", revision)
	sheet.Set("effective_date", effectiveDate)
	sheet.Set("active", active)
	if err := app.Save(sheet); err != nil {
		t.Fatalf("saving rate sheet %s: %v", id, err)
	}
	entryCollection, err := app.FindCollectionByNameOrId("rate_sheet_entries")
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		entry := core.NewRecord(entryCollection)
		entry.Set("rate_sheet", id)
		entry.Set("role", e.Role)
		entry.Set("rate", e.Rate)
		entry.Set("overtime_rate", e.OvertimeRate)
		if err := app.Save(entry); err != nil {
			t.Fatalf("saving an entry of %s: %v", id, err)
		}
	}
}

func rateSheetAbsorbSQL(t testing.TB, app core.App, query string, params dbx.Params) {
	t.Helper()
	if _, err := app.DB().NewQuery(query).Bind(params).Execute(); err != nil {
		t.Fatalf("%s: %v", query, err)
	}
}

// newRateSheetAbsorbApp returns a seeded app with the fixtures listed in the
// constants above.
func newRateSheetAbsorbApp(t testing.TB) *tests.TestApp {
	t.Helper()
	app := testutils.SetupTestApp(t)
	a, b, c := absorbRateSets(t, app)
	createRateSheetFixture(t, app, rsAbTarget, "Absorb Target", 0, "2024-01-01", true, a)
	createRateSheetFixture(t, app, rsAbSingle, "Absorb Single", 0, "2024-06-01", true, a)
	createRateSheetFixture(t, app, rsAbEarly, "Absorb Early", 0, "2023-01-01", false, a)
	createRateSheetFixture(t, app, rsAbMissing, "Absorb Missing Role", 0, "2024-06-01", false, b)
	createRateSheetFixture(t, app, rsAbDiffRate, "Absorb Different Rate", 0, "2024-06-01", false, c)
	createRateSheetFixture(t, app, rsAbWhole0, "Absorb Whole", 0, "2022-01-01", false, c)
	createRateSheetFixture(t, app, rsAbWhole1, "Absorb Whole", 1, "2022-06-01", true, a)
	createRateSheetFixture(t, app, rsAbUsed0, "Absorb Used", 0, "2024-01-01", false, a)
	createRateSheetFixture(t, app, rsAbUsed1, "Absorb Used", 1, "2024-02-01", true, a)
	createRateSheetFixture(t, app, rsAbLateInactive, "Absorb Late Inactive", 0, "2025-01-01", false, a)
	createRateSheetFixture(t, app, rsAbLatest0, "Absorb Latest Differs", 0, "2024-06-01", false, a)
	createRateSheetFixture(t, app, rsAbLatest1, "Absorb Latest Differs", 1, "2024-07-01", true, c)
	createRateSheetFixture(t, app, rsAbInactiveTgt, "Absorb Inactive Target", 0, "2020-01-01", false, a)
	rateSheetAbsorbSQL(t, app, "UPDATE jobs SET rate_sheet = {:sheet} WHERE id = {:job}", dbx.Params{"sheet": rsAbSingle, "job": rsAbSingleJob})
	rateSheetAbsorbSQL(t, app, "UPDATE jobs SET rate_sheet = {:sheet} WHERE id = {:job}", dbx.Params{"sheet": rsAbUsed0, "job": rsAbUsedJob})
	rateSheetAbsorbSQL(t, app, "UPDATE job_requests SET rate_sheet = {:sheet} WHERE id = {:request}", dbx.Params{"sheet": rsAbSingle, "request": rsAbSingleReq})
	return app
}

// grantRateSheetRevise gives book@keeper.com the rate_sheet_revise claim, so
// that it holds absorb and rate_sheet_revise.
func grantRateSheetRevise(t testing.TB, app core.App) {
	t.Helper()
	user, err := app.FindAuthRecordByEmail("users", absorbReviseEmail)
	if err != nil {
		t.Fatal(err)
	}
	claim, err := app.FindFirstRecordByData("claims", "name", "rate_sheet_revise")
	if err != nil {
		t.Fatal(err)
	}
	collection, err := app.FindCollectionByNameOrId("user_claims")
	if err != nil {
		t.Fatal(err)
	}
	userClaim := core.NewRecord(collection)
	userClaim.Set("uid", user.Id)
	userClaim.Set("cid", claim.Id)
	if err := app.Save(userClaim); err != nil {
		t.Fatal(err)
	}
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
	grantRateSheetRevise(t, app)

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

	// The absorbed sheet was renamed before.
	history, err := app.FindCollectionByNameOrId("rate_sheet_name_changes")
	if err != nil {
		t.Fatal(err)
	}
	admin, err := app.FindAuthRecordByEmail("users", absorbAdminEmail)
	if err != nil {
		t.Fatal(err)
	}
	earlier := core.NewRecord(history)
	earlier.Set("rate_sheet", rsAbSingle)
	earlier.Set("old_name", "Absorb Single Before")
	earlier.Set("new_name", "Absorb Single")
	earlier.Set("renamed_by", admin.Id)
	if err := app.Save(earlier); err != nil {
		t.Fatal(err)
	}

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
	// check. The admin without absorb is refused by the collection rule.
	if status, body := absorbRequest(t, app, http.MethodDelete, path, absorbOnlyEmail, nil); status != http.StatusForbidden {
		t.Fatalf("commit with absorb only: status %d: %s", status, body)
	}
	if status, _ := absorbRequest(t, app, http.MethodDelete, path, absorbAdminEmail, nil); status != http.StatusNotFound {
		t.Fatalf("commit by an admin without absorb: status %d; want the collection rule's 404", status)
	}
	if names, _ := routes.RateSheetFormerNames(app, "Absorb Target"); len(names) != 0 {
		t.Fatalf("a refused commit wrote history: %v", names)
	}

	grantRateSheetRevise(t, app)
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
		grant  bool
		status int
	}{
		{"admin", absorbAdminEmail, false, http.StatusOK},
		{"absorb and rate_sheet_revise", absorbReviseEmail, true, http.StatusOK},
		{"absorb alone", absorbOnlyEmail, false, http.StatusForbidden},
		{"rate_sheet_revise alone", reviseOnlyEmail, false, http.StatusForbidden},
		{"signed out", "", false, http.StatusUnauthorized},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			app := newRateSheetAbsorbApp(t)
			defer app.Cleanup()
			if tc.grant {
				grantRateSheetRevise(t, app)
			}
			status, body := absorbRequest(t, app, http.MethodPost, "/api/rate_sheets/"+rsAbTarget+"/absorb", tc.email, absorbBody(rsAbSingle))
			if status != tc.status {
				t.Fatalf("absorb status %d; want %d: %s", status, tc.status, body)
			}
			if tc.status == http.StatusOK {
				if !strings.Contains(body, "Successfully absorbed 1 records") {
					t.Fatalf("absorb body: %s", body)
				}
				// Undo is refused for absorb alone and allowed for the user who
				// absorbed. book@keeper.com holds absorb alone until granted.
				if !tc.grant {
					if status, _ := absorbRequest(t, app, http.MethodPost, "/api/rate_sheets/undo_absorb", absorbOnlyEmail, nil); status != http.StatusForbidden {
						t.Fatalf("undo with absorb alone: status %d", status)
					}
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
	// Copy the job's sheet (rsvaluation0001, a revision of a name that jobs
	// use) to two single revision sheets with the same rates, and move the
	// job to one of them.
	records, err := app.FindAllRecords("rate_sheet_entries", dbx.HashExp{"rate_sheet": "rsvaluation0001"})
	if err != nil {
		t.Fatal(err)
	}
	entries := []rateSheetEntryFixture{}
	for _, r := range records {
		entries = append(entries, rateSheetEntryFixture{r.GetString("role"), r.GetFloat("rate"), r.GetFloat("overtime_rate")})
	}
	createRateSheetFixture(t, app, "rsabwipsource01", "WIP Source", 0, "2034-01-01", false, entries)
	createRateSheetFixture(t, app, "rsabwiptarget01", "WIP Target", 0, "2033-01-01", false, entries)
	rateSheetAbsorbSQL(t, app, "UPDATE jobs SET rate_sheet = 'rsabwipsource01' WHERE id = {:job}", dbx.Params{"job": job})

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
	if after := wip(); after != before {
		t.Fatalf("WIP after absorb = %+v; want %+v", after, before)
	}
	if status, body := absorbRequest(t, app, http.MethodPost, "/api/rate_sheets/undo_absorb", absorbAdminEmail, nil); status != http.StatusOK {
		t.Fatalf("undo status %d: %s", status, body)
	}
	if got := fieldOf(t, app, "jobs", job, "rate_sheet"); got != "rsabwipsource01" {
		t.Fatalf("undo left the job on %s", got)
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
	createRateSheetFixture(t, app, "rsabreqtarget01", "Request Target", 0, "2025-01-01", true, standardEntries(t, app))

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

	if status, body := absorbRequest(t, app, http.MethodPost, "/api/rate_sheets/undo_absorb", absorbAdminEmail, nil); status != http.StatusOK {
		t.Fatalf("undo status %d: %s", status, body)
	}
	if got := fieldOf(t, app, "jobs", jobID, "rate_sheet"); got != rsAbStandardRate {
		t.Fatalf("undo left the new job on %s", got)
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

	// A group whose earliest sheet is inactive and older than an active
	// revision says why nothing can be absorbed.
	app2 := testutils.SetupTestApp(t)
	defer app2.Cleanup()
	a, _, _ := absorbRateSets(t, app2)
	createRateSheetFixture(t, app2, "rsabgroupold001", "Group Old", 0, "2020-01-01", false, a)
	createRateSheetFixture(t, app2, "rsabgroupold002", "Group Old", 1, "2021-01-01", false, a)
	createRateSheetFixture(t, app2, "rsabgroupnew001", "Group New", 0, "2022-01-01", true, a)
	rateSheetAbsorbSQL(t, app2, "UPDATE jobs SET rate_sheet = 'rsabgroupold001' WHERE id = {:job}", dbx.Params{"job": rsAbSingleJob})
	status, body = absorbRequest(t, app2, http.MethodGet, "/api/rate_sheets/duplicates", reviseOnlyEmail, nil)
	if status != http.StatusOK {
		t.Fatalf("status %d: %s", status, body)
	}
	if !strings.Contains(body, "has a newer revision") {
		t.Fatalf("expected the group to say why it can't be absorbed: %s", body)
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
