package main

import (
	"net/http"
	"strings"
	"testing"

	"tybalt/hooks"
	"tybalt/internal/testutils"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/tests"
)

// A PocketBase superuser may change a rate sheet code through the collection
// API; nobody else can change one by any route. The fixture superuser is
// test@example.com. Complete Hook Sheet (rs_complete_all_001, code
// 2026-03-00001) has no other revisions; Sheet A (rs_sheeta_r0 and
// rs_sheeta_r1, code 2026-03-00009) has two.

func superuserToken(t *testing.T) string {
	t.Helper()
	token, err := testutils.GenerateAdminToken("test@example.com")
	if err != nil {
		t.Fatal(err)
	}
	return token
}

func patchRateSheetScenario(name, token, id, body string, status int, content ...string) tests.ApiScenario {
	return tests.ApiScenario{
		Name:            name,
		Method:          http.MethodPatch,
		URL:             "/api/collections/rate_sheets/records/" + id,
		Body:            strings.NewReader(body),
		Headers:         map[string]string{"Authorization": token},
		ExpectedStatus:  status,
		ExpectedContent: content,
	}
}

func codeChanges(t testing.TB, app *tests.TestApp, sheet string) []string {
	t.Helper()
	rows, err := app.FindAllRecords("rate_sheet_code_changes", dbx.HashExp{"rate_sheet": sheet})
	if err != nil {
		t.Fatal(err)
	}
	result := []string{}
	for _, row := range rows {
		result = append(result, row.GetString("old_code")+">"+row.GetString("new_code"))
	}
	return result
}

func TestSuperuserCanChangeARateSheetCode(t *testing.T) {
	app := testutils.SetupTestApp(t)
	token := superuserToken(t)

	runOnApp(t, app, patchRateSheetScenario("superuser changes a code", token, "rs_complete_all_001",
		`{"code":"2026-03-00099"}`, http.StatusOK, `"code":"2026-03-00099"`))

	if got := rateSheetCodeOf(t, app, "rs_complete_all_001"); got != "2026-03-00099" {
		t.Fatalf("code = %s; want 2026-03-00099", got)
	}
	if got := codeChanges(t, app, "rs_complete_all_001"); len(got) != 1 || got[0] != "2026-03-00001>2026-03-00099" {
		t.Fatalf("code changes = %v", got)
	}
	change, err := app.FindFirstRecordByData("rate_sheet_code_changes", "rate_sheet", "rs_complete_all_001")
	if err != nil {
		t.Fatal(err)
	}
	if change.GetDateTime("created").IsZero() {
		t.Fatal("the change has no time")
	}
	// The new code is issued, so the generator goes on after it, and the old
	// code stays issued.
	for _, code := range []string{"2026-03-00099", "2026-03-00001"} {
		if issued, err := hooks.RateSheetCodeIssued(app, code); err != nil || !issued {
			t.Fatalf("%s is not recorded as issued (%v)", code, err)
		}
	}

	// A change of the name alone writes no code change.
	runOnApp(t, app, patchRateSheetScenario("superuser renames without a code change", token, "rs_complete_all_001",
		`{"name":"Complete Hook Sheet Renamed"}`, http.StatusOK, `"code":"2026-03-00099"`))
	if got := codeChanges(t, app, "rs_complete_all_001"); len(got) != 1 {
		t.Fatalf("a rename logged a code change: %v", got)
	}
}

func TestSuperuserCodeChangeRules(t *testing.T) {
	token := superuserToken(t)
	scenarios := []struct {
		name, id, body string
		status         int
		content        []string
	}{
		{"a code in use", "rs_complete_all_001", `{"code":"2026-03-00002"}`, http.StatusConflict, []string{"never used twice"}},
		{"a code that was issued to a revision family", "rs_complete_all_001", `{"code":"2026-03-00009"}`, http.StatusConflict, []string{"never used twice"}},
		{"a month that does not exist", "rs_complete_all_001", `{"code":"2026-13-00001"}`, http.StatusBadRequest, []string{"invalid_rate_sheet_code"}},
		{"a short serial", "rs_complete_all_001", `{"code":"2026-03-001"}`, http.StatusBadRequest, []string{"invalid_rate_sheet_code"}},
		{"serial zero", "rs_complete_all_001", `{"code":"2026-03-00000"}`, http.StatusBadRequest, []string{"invalid_rate_sheet_code"}},
		{"an empty code", "rs_complete_all_001", `{"code":""}`, http.StatusBadRequest, []string{"invalid_rate_sheet_code"}},
		{"a sheet with revisions", "rs_sheeta_r0", `{"code":"2026-03-00098"}`, http.StatusBadRequest, []string{"revisions"}},
	}
	for _, s := range scenarios {
		app := testutils.SetupTestApp(t)
		before := rateSheetCodeOf(t, app, s.id)
		runOnApp(t, app, patchRateSheetScenario(s.name, token, s.id, s.body, s.status, s.content...))
		if got := rateSheetCodeOf(t, app, s.id); got != before {
			t.Fatalf("%s: code changed to %s", s.name, got)
		}
		if got := codeChanges(t, app, s.id); len(got) != 0 {
			t.Fatalf("%s: logged %v", s.name, got)
		}
	}
}

// A code that was issued and then retired, by deleting its sheet, is never
// given out again, not even by a superuser.
func TestSuperuserCannotReuseARetiredCode(t *testing.T) {
	app := testutils.SetupTestApp(t)
	status, created := createRateSheetThroughRoute(t, app, "author@soup.com", map[string]any{"name": "Retired Code Sheet", "revision": 0})
	if status != http.StatusCreated {
		t.Fatalf("create status %d", status)
	}
	retired := created["code"].(string)
	sheet, err := app.FindRecordById("rate_sheets", created["id"].(string))
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
	runOnApp(t, app, patchRateSheetScenario("superuser reuses a retired code", superuserToken(t), "rs_complete_all_001",
		`{"code":"`+retired+`"}`, http.StatusConflict, "never used twice"))
}

// Nobody but a superuser can change a code, through the collection API or
// through server code, and no client can choose the code of a new sheet.
func TestOnlyASuperuserCanChangeARateSheetCode(t *testing.T) {
	adminToken, err := testutils.GenerateRecordToken("users", "author@soup.com")
	if err != nil {
		t.Fatal(err)
	}
	reviseToken, err := testutils.GenerateRecordToken("users", "time@test.com")
	if err != nil {
		t.Fatal(err)
	}

	app := testutils.SetupTestApp(t)
	for _, token := range []string{adminToken, reviseToken} {
		scenario := patchRateSheetScenario("a user changes a code", token, "rs_complete_all_001", `{"code":"2026-03-00097"}`, http.StatusForbidden, "Only superusers")
		runOnApp(t, app, scenario)
	}
	if got := rateSheetCodeOf(t, app, "rs_complete_all_001"); got != "2026-03-00001" {
		t.Fatalf("code = %s after refused changes", got)
	}

	// Server code that saves a sheet cannot change its code either.
	sheet, err := app.FindRecordById("rate_sheets", "rs_complete_all_001")
	if err != nil {
		t.Fatal(err)
	}
	sheet.Set("code", "2026-03-00096")
	if err := app.Save(sheet); err == nil {
		t.Fatal("a save without a superuser changed the code")
	}

	// The collection API refuses a code in a new sheet.
	runOnApp(t, app, tests.ApiScenario{
		Name:            "a user sends a code with a new sheet",
		Method:          http.MethodPost,
		URL:             "/api/collections/rate_sheets/records",
		Body:            strings.NewReader(`{"name":"Chosen Code","effective_date":"2026-11-01","revision":0,"code":"2026-11-00090"}`),
		Headers:         map[string]string{"Authorization": adminToken},
		ExpectedStatus:  http.StatusBadRequest,
		ExpectedContent: []string{"Failed to create record"},
	})
	if _, err := app.FindFirstRecordByData("rate_sheets", "name", "Chosen Code"); err == nil {
		t.Fatal("the collection API created a sheet with a chosen code")
	}
}
