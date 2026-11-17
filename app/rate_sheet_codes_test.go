package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"tybalt/hooks"
	"tybalt/internal/testutils"
	"tybalt/migrations"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

// Rate sheet codes (issue 207). Every sheet in the seed CSV files has the
// code that the migration's backfill gives it (rate_sheets.csv), and
// rate_sheet_codes.csv holds every one of those codes as issued. The months
// 2025-01, 2026-01, 2026-03, 2026-09 and 2026-11 have codes; 2026-03 has
// twelve and 2026-11 twenty-seven.

// rateSheetCodeOf returns the code of a sheet.
func rateSheetCodeOf(t testing.TB, app core.App, id string) string {
	t.Helper()
	sheet, err := app.FindRecordById("rate_sheets", id)
	if err != nil {
		t.Fatalf("find rate sheet %s: %v", id, err)
	}
	return sheet.GetString("code")
}

// nextRateSheetCode is the code that the next sheet created now gets.
func nextRateSheetCode(t testing.TB, app core.App) string {
	t.Helper()
	month := time.Now().UTC().Format("2006-01")
	var issued int
	if err := app.DB().NewQuery("SELECT COALESCE(MAX(serial), 0) FROM rate_sheet_codes WHERE month = {:month}").
		Bind(dbx.Params{"month": month}).Row(&issued); err != nil {
		t.Fatal(err)
	}
	return hooks.FormatRateSheetCode(month, issued+1)
}

// createRateSheetThroughRoute creates a sheet with one entry through POST
// /api/rate_sheets and returns the response.
func createRateSheetThroughRoute(t *testing.T, app *tests.TestApp, email string, body map[string]any) (int, map[string]any) {
	t.Helper()
	if _, ok := body["entries"]; !ok {
		roles, err := app.FindRecordsByFilter("rate_roles", "1=1", "name", 1, 0, nil)
		if err != nil || len(roles) == 0 {
			t.Fatalf("find a rate role: %v", err)
		}
		// A rate no other sheet has, so the duplicate check passes.
		body["entries"] = []map[string]any{{"role": roles[0].Id, "rate": 7777, "overtime_rate": 9999}}
	}
	if _, ok := body["effective_date"]; !ok {
		body["effective_date"] = "2026-11-01"
	}
	status, raw := absorbRequest(t, app, http.MethodPost, "/api/rate_sheets", email, body)
	result := map[string]any{}
	_ = json.Unmarshal([]byte(raw), &result)
	if status != http.StatusCreated {
		t.Logf("create rate sheet: %d %s", status, raw)
	}
	return status, result
}

// Codes run from 00001 in each month, one after another, and a new month
// starts again.
func TestRateSheetCodesRunPerMonth(t *testing.T) {
	app := testutils.SetupTestApp(t)

	issue := func(at string) string {
		t.Helper()
		when, err := time.Parse("2006-01-02", at)
		if err != nil {
			t.Fatal(err)
		}
		var code string
		if err := app.RunInTransaction(func(tx core.App) error {
			code, err = hooks.IssueRateSheetCode(tx, when)
			return err
		}); err != nil {
			t.Fatal(err)
		}
		return code
	}

	if got := issue("2031-05-10"); got != "2031-05-00001" {
		t.Fatalf("first code of a month = %s; want 2031-05-00001", got)
	}
	if got := issue("2031-05-31"); got != "2031-05-00002" {
		t.Fatalf("second code of a month = %s; want 2031-05-00002", got)
	}
	if got := issue("2031-06-01"); got != "2031-06-00001" {
		t.Fatalf("first code of the next month = %s; want 2031-06-00001", got)
	}
	// 2026-03 has twelve codes in the fixtures.
	if got := issue("2026-03-20"); got != "2026-03-00013" {
		t.Fatalf("next code after the fixtures = %s; want 2026-03-00013", got)
	}
	for _, code := range []string{"2031-05-00001", "2031-05-00002", "2031-06-00001", "2026-03-00013"} {
		issued, err := hooks.RateSheetCodeIssued(app, code)
		if err != nil || !issued {
			t.Fatalf("%s is not recorded as issued (%v)", code, err)
		}
	}
}

// A new sheet gets the next code of the current month, and a deleted sheet's
// code is never issued again.
func TestNewRateSheetGetsTheNextCodeAndNoCodeIsReused(t *testing.T) {
	app := testutils.SetupTestApp(t)

	want := nextRateSheetCode(t, app)
	status, first := createRateSheetThroughRoute(t, app, "author@soup.com", map[string]any{"name": "Code Test First", "revision": 0})
	if status != http.StatusCreated {
		t.Fatalf("create status %d", status)
	}
	if first["code"] != want {
		t.Fatalf("code = %v; want %s", first["code"], want)
	}
	if got := rateSheetCodeOf(t, app, first["id"].(string)); got != want {
		t.Fatalf("stored code = %s; want %s", got, want)
	}

	// Delete it: its code stays issued and the next sheet gets a new one.
	sheet, err := app.FindRecordById("rate_sheets", first["id"].(string))
	if err != nil {
		t.Fatal(err)
	}
	entries, err := app.FindAllRecords("rate_sheet_entries", dbx.HashExp{"rate_sheet": sheet.Id})
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if err := app.Delete(entry); err != nil {
			t.Fatal(err)
		}
	}
	if err := app.Delete(sheet); err != nil {
		t.Fatal(err)
	}
	issued, err := hooks.RateSheetCodeIssued(app, want)
	if err != nil || !issued {
		t.Fatalf("the deleted sheet's code is no longer issued (%v)", err)
	}

	status, second := createRateSheetThroughRoute(t, app, "author@soup.com", map[string]any{"name": "Code Test Second", "revision": 0})
	if status != http.StatusCreated {
		t.Fatalf("create status %d", status)
	}
	if second["code"] == want {
		t.Fatalf("the deleted sheet's code %s was issued again", want)
	}
	if second["code"] != nextAfter(t, want) {
		t.Fatalf("code = %v; want %s", second["code"], nextAfter(t, want))
	}
}

// nextAfter returns the code after code in the same month.
func nextAfter(t testing.TB, code string) string {
	t.Helper()
	month, serial, ok := hooks.ParseRateSheetCode(code)
	if !ok {
		t.Fatalf("bad code %s", code)
	}
	return hooks.FormatRateSheetCode(month, serial+1)
}

// A revision has the code of the sheet it revises.
func TestRateSheetRevisionKeepsTheCode(t *testing.T) {
	app := testutils.SetupTestApp(t)
	base := rateSheetCodeOf(t, app, "rsrevbase0001")

	status, revision := createRateSheetThroughRoute(t, app, "time@test.com", map[string]any{
		"name": "Revise Test Sheet", "revision": 1, "effective_date": "2026-02-01",
	})
	if status != http.StatusCreated {
		t.Fatalf("create status %d", status)
	}
	if revision["code"] != base {
		t.Fatalf("revision code = %v; want %s", revision["code"], base)
	}
	if got := rateSheetCodeOf(t, app, revision["id"].(string)); got != base {
		t.Fatalf("stored revision code = %s; want %s", got, base)
	}
}

// An absorb keeps the target's code. The absorbed sheet's code is retired:
// undo puts the sheet back with it, and after a commit it is never issued
// again.
func TestAbsorbedRateSheetCodeIsRetired(t *testing.T) {
	app := newRateSheetAbsorbApp(t)
	targetCode := rateSheetCodeOf(t, app, rsAbTarget)
	absorbedCode := rateSheetCodeOf(t, app, rsAbSingle)

	absorb := func() {
		t.Helper()
		status, body := absorbRequest(t, app, http.MethodPost, "/api/rate_sheets/"+rsAbTarget+"/absorb", absorbAdminEmail, absorbBody(rsAbSingle))
		if status != http.StatusOK {
			t.Fatalf("absorb status %d: %s", status, body)
		}
	}

	absorb()
	status, body := absorbRequest(t, app, http.MethodPost, "/api/rate_sheets/undo_absorb", absorbAdminEmail, nil)
	if status != http.StatusOK {
		t.Fatalf("undo status %d: %s", status, body)
	}
	if got := rateSheetCodeOf(t, app, rsAbSingle); got != absorbedCode {
		t.Fatalf("undo restored code %s; want %s", got, absorbedCode)
	}

	absorb()
	// Deleting the absorb action commits the absorb.
	action := rateSheetAction(t, app)
	status, body = absorbRequest(t, app, http.MethodDelete, "/api/collections/absorb_actions/records/"+action.Id, absorbAdminEmail, nil)
	if status != http.StatusNoContent {
		t.Fatalf("commit status %d: %s", status, body)
	}
	if got := rateSheetCodeOf(t, app, rsAbTarget); got != targetCode {
		t.Fatalf("target code = %s; want %s", got, targetCode)
	}
	issued, err := hooks.RateSheetCodeIssued(app, absorbedCode)
	if err != nil || !issued {
		t.Fatalf("the absorbed code is no longer recorded as issued (%v)", err)
	}
	status, created := createRateSheetThroughRoute(t, app, "author@soup.com", map[string]any{"name": "After Absorb", "revision": 0})
	if status != http.StatusCreated {
		t.Fatalf("create status %d", status)
	}
	if created["code"] == absorbedCode {
		t.Fatalf("the absorbed code %s was issued again", absorbedCode)
	}
}

// The migration's backfill gives the fixture sheets exactly the codes in the
// seed CSV files. The test changes fixture data at test time: it drops the
// unique index on code and revision, empties every code and empties
// rate_sheet_codes, so that the backfill starts from the state before the
// migration.
func TestRateSheetCodeBackfill(t *testing.T) {
	app := testutils.SetupTestApp(t)

	type row struct {
		ID   string `db:"id"`
		Code string `db:"code"`
	}
	before := []row{}
	if err := app.DB().NewQuery("SELECT id, code FROM rate_sheets ORDER BY id").All(&before); err != nil {
		t.Fatal(err)
	}
	type issuedRow struct {
		Code   string `db:"code"`
		Month  string `db:"month"`
		Serial int    `db:"serial"`
	}
	issuedBefore := []issuedRow{}
	if err := app.DB().NewQuery("SELECT code, month, serial FROM rate_sheet_codes ORDER BY code").All(&issuedBefore); err != nil {
		t.Fatal(err)
	}

	for _, query := range []string{
		"DROP INDEX idx_rate_sheets_code_revision",
		"UPDATE rate_sheets SET code = ''",
		"DELETE FROM rate_sheet_codes",
	} {
		if _, err := app.DB().NewQuery(query).Execute(); err != nil {
			t.Fatalf("%s: %v", query, err)
		}
	}
	if err := migrations.BackfillRateSheetCodes(app); err != nil {
		t.Fatal(err)
	}

	after := []row{}
	if err := app.DB().NewQuery("SELECT id, code FROM rate_sheets ORDER BY id").All(&after); err != nil {
		t.Fatal(err)
	}
	if fmt.Sprint(after) != fmt.Sprint(before) {
		t.Fatalf("backfill codes differ from the fixtures:\n got %v\nwant %v", after, before)
	}
	issuedAfter := []issuedRow{}
	if err := app.DB().NewQuery("SELECT code, month, serial FROM rate_sheet_codes ORDER BY code").All(&issuedAfter); err != nil {
		t.Fatal(err)
	}
	if fmt.Sprint(issuedAfter) != fmt.Sprint(issuedBefore) {
		t.Fatalf("backfill issued codes differ from the fixtures:\n got %v\nwant %v", issuedAfter, issuedBefore)
	}

	codes := map[string]string{}
	for _, r := range after {
		codes[r.ID] = r.Code
	}
	checks := []struct {
		what, id, want string
	}{
		// The first sheet of 2026-01.
		{"a sheet's month is its created month", "c41ofep525bcacj", "2026-01-00001"},
		// Both revisions were created on the same day; they share a code.
		{"revision 0", "rs_sheeta_r0", "2026-03-00009"},
		{"revision 1 shares the code", "rs_sheeta_r1", "2026-03-00009"},
		// No created date: the effective date gives the month.
		{"no created date", "test_rs_alt_actv", "2025-01-00001"},
		// Only revision 1 exists; its created date has no Z.
		{"a sheet without revision 0", "test_empty_sheet", "2025-01-00002"},
		// A draft has its own code.
		{"a draft", "rsdraftsheet001", "2026-11-00025"},
		{"the latest sheet of 2026-11", "rsdraftsheet002", "2026-11-00027"},
	}
	for _, check := range checks {
		if codes[check.id] != check.want {
			t.Errorf("%s: %s has code %s; want %s", check.what, check.id, codes[check.id], check.want)
		}
	}
}

// Families come from names, a draft is always its own family, and the month
// falls back from the first revision's created date to the family's earliest
// created date, then to the effective date. These rows are the input of a
// pure function, not stored data.
func TestPlanRateSheetCodes(t *testing.T) {
	now := time.Date(2030, 7, 4, 12, 0, 0, 0, time.UTC)
	rows := []migrations.RateSheetCodeBackfillRow{
		{ID: "b", Name: "Shared", Revision: 1, Created: "2026-05-02 00:00:00.000Z"},
		{ID: "a", Name: "Shared", Revision: 0, Created: ""},
		{ID: "c", Name: "shared", Revision: 0, Created: "2026-05-01 00:00:00.000Z"},
		{ID: "d", Name: "Shared", Revision: 0, Created: "2026-04-01 00:00:00.000Z", Request: "req1"},
		{ID: "e", Name: "Nothing", Revision: 0, EffectiveDate: "not a date"},
	}
	got := map[string]string{}
	for _, assignment := range migrations.PlanRateSheetCodes(rows, now) {
		for _, id := range assignment.IDs {
			got[id] = assignment.Code
		}
	}
	want := map[string]string{
		// "shared" differs in capitals, so it was a sheet of its own.
		"c": "2026-05-00001",
		// Revision 0 of "Shared" has no created date; revision 1's is used.
		"a": "2026-05-00002",
		"b": "2026-05-00002",
		"d": "2026-04-00001",
		"e": "2030-07-00001",
	}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("plan = %v; want %v", got, want)
	}
}
