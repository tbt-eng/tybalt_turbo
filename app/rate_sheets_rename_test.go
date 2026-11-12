package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"testing"
	"tybalt/internal/testutils"
	"tybalt/routes"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

// author@soup.com holds admin. time@test.com holds rate_sheet_revise but not
// admin. "Summary valuation fixture" has two revisions (rsvaluation0001 and
// rsvaluation0002), and jobs use both.
const (
	renameAdminUserID   = "f2j5a8vk006baub"
	renameValuationName = "Summary valuation fixture"
)

func renameToken(t *testing.T, email string) string {
	t.Helper()
	token, err := testutils.GenerateRecordToken("users", email)
	if err != nil {
		t.Fatal(err)
	}
	return token
}

func renameBody(name string) *strings.Reader {
	body, _ := json.Marshal(map[string]string{"name": name})
	return strings.NewReader(string(body))
}

// runOnApp runs a scenario against app and keeps app for the next scenario.
func runOnApp(t *testing.T, app *tests.TestApp, scenario tests.ApiScenario) {
	t.Helper()
	scenario.TestAppFactory = func(testing.TB) *tests.TestApp { return app }
	scenario.DisableTestAppCleanup = true
	scenario.Test(t)
}

func renameScenario(token, id, name string, status int, content ...string) tests.ApiScenario {
	return tests.ApiScenario{
		Name:            "rename " + id + " to " + name,
		Method:          http.MethodPost,
		URL:             "/api/rate_sheets/" + id + "/rename",
		Body:            renameBody(name),
		Headers:         map[string]string{"Authorization": token},
		ExpectedStatus:  status,
		ExpectedContent: content,
	}
}

func formerNamesOf(t *testing.T, app core.App, id string) []string {
	t.Helper()
	sheet, err := app.FindRecordById("rate_sheets", id)
	if err != nil {
		t.Fatalf("find rate sheet %s: %v", id, err)
	}
	names, err := routes.RateSheetFormerNames(app, sheet.GetString("name"))
	if err != nil {
		t.Fatalf("former names of %s: %v", id, err)
	}
	return names
}

func TestRenameRateSheet_Permissions(t *testing.T) {
	adminToken := renameToken(t, "author@soup.com")
	reviseToken := renameToken(t, "time@test.com")
	noClaimsToken := renameToken(t, "noclaims@example.com")

	scenarios := []tests.ApiScenario{
		{
			Name:           "unauthenticated request is refused",
			Method:         http.MethodPost,
			URL:            "/api/rate_sheets/rsvaluation0001/rename",
			Body:           renameBody("New Name"),
			ExpectedStatus: http.StatusUnauthorized,
			ExpectedContent: []string{
				`"status":401`,
			},
			TestAppFactory: testutils.SetupTestApp,
		},
		{
			Name:            "rate_sheet_revise does not allow a rename",
			Method:          http.MethodPost,
			URL:             "/api/rate_sheets/rsvaluation0001/rename",
			Body:            renameBody("New Name"),
			Headers:         map[string]string{"Authorization": reviseToken},
			ExpectedStatus:  http.StatusForbidden,
			ExpectedContent: []string{`You do not have permission to rename rate sheets`},
			TestAppFactory:  testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				sheet, err := app.FindRecordById("rate_sheets", "rsvaluation0001")
				if err != nil {
					t.Fatal(err)
				}
				if sheet.GetString("name") != renameValuationName {
					t.Fatalf("name changed to %q", sheet.GetString("name"))
				}
			},
		},
		{
			Name:            "user without claims is refused",
			Method:          http.MethodPost,
			URL:             "/api/rate_sheets/rsvaluation0001/rename",
			Body:            renameBody("New Name"),
			Headers:         map[string]string{"Authorization": noClaimsToken},
			ExpectedStatus:  http.StatusForbidden,
			ExpectedContent: []string{`You do not have permission to rename rate sheets`},
			TestAppFactory:  testutils.SetupTestApp,
		},
		{
			Name:            "unknown sheet returns 404",
			Method:          http.MethodPost,
			URL:             "/api/rate_sheets/nosuchsheet0000/rename",
			Body:            renameBody("New Name"),
			Headers:         map[string]string{"Authorization": adminToken},
			ExpectedStatus:  http.StatusNotFound,
			ExpectedContent: []string{`"status":404`},
			TestAppFactory:  testutils.SetupTestApp,
		},
	}
	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestRenameRateSheet_RenamesEveryRevisionAndWritesHistory(t *testing.T) {
	adminToken := renameToken(t, "author@soup.com")

	scenario := tests.ApiScenario{
		Name:           "admin renames all revisions",
		Method:         http.MethodPost,
		URL:            "/api/rate_sheets/rsvaluation0001/rename",
		Body:           renameBody("  Generic Valuation Rates  "),
		Headers:        map[string]string{"Authorization": adminToken},
		ExpectedStatus: http.StatusOK,
		ExpectedContent: []string{
			`"name":"Generic Valuation Rates"`,
			`"old_name":"Summary valuation fixture"`,
			`"revisions_renamed":2`,
			`"former_names":["Summary valuation fixture"]`,
		},
		TestAppFactory: testutils.SetupTestApp,
		BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
			var count int
			if err := app.DB().NewQuery("SELECT COUNT(*) FROM jobs WHERE rate_sheet IN ('rsvaluation0001', 'rsvaluation0002')").Row(&count); err != nil {
				t.Fatal(err)
			}
			if count == 0 {
				t.Fatal("fixture has no jobs on the sheet")
			}
		},
		AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
			for _, id := range []string{"rsvaluation0001", "rsvaluation0002"} {
				sheet, err := app.FindRecordById("rate_sheets", id)
				if err != nil {
					t.Fatal(err)
				}
				if got := sheet.GetString("name"); got != "Generic Valuation Rates" {
					t.Errorf("%s name = %q", id, got)
				}
			}
			old, err := app.FindRecordsByFilter("rate_sheets", "name = {:name}", "", 0, 0, dbx.Params{"name": renameValuationName})
			if err != nil {
				t.Fatal(err)
			}
			if len(old) != 0 {
				t.Errorf("%d sheets still have the old name", len(old))
			}

			// Jobs keep their rate_sheet ids and now show the new name.
			var jobNames []string
			if err := app.DB().NewQuery(`
				SELECT DISTINCT rs.name FROM jobs j JOIN rate_sheets rs ON rs.id = j.rate_sheet
				WHERE j.rate_sheet IN ('rsvaluation0001', 'rsvaluation0002')
			`).Column(&jobNames); err != nil {
				t.Fatal(err)
			}
			if !slices.Equal(jobNames, []string{"Generic Valuation Rates"}) {
				t.Errorf("job rate sheet names = %v", jobNames)
			}

			changes, err := app.FindAllRecords("rate_sheet_name_changes")
			if err != nil {
				t.Fatal(err)
			}
			if len(changes) != 1 {
				t.Fatalf("want 1 history row, got %d", len(changes))
			}
			change := changes[0]
			if change.GetString("rate_sheet") != "rsvaluation0001" ||
				change.GetString("old_name") != renameValuationName ||
				change.GetString("new_name") != "Generic Valuation Rates" ||
				change.GetString("renamed_by") != renameAdminUserID ||
				change.GetString("created") == "" {
				t.Errorf("unexpected history row: %v", change.FieldsData())
			}
		},
	}
	scenario.Test(t)
}

func TestRenameRateSheet_RefusesBadNames(t *testing.T) {
	adminToken := renameToken(t, "author@soup.com")

	noChange := func(t testing.TB, app *tests.TestApp, res *http.Response) {
		changes, err := app.FindAllRecords("rate_sheet_name_changes")
		if err != nil {
			t.Fatal(err)
		}
		if len(changes) != 0 {
			t.Errorf("want no history rows, got %d", len(changes))
		}
		sheet, err := app.FindRecordById("rate_sheets", "rsvaluation0002")
		if err != nil {
			t.Fatal(err)
		}
		if sheet.GetString("name") != renameValuationName {
			t.Errorf("name changed to %q", sheet.GetString("name"))
		}
	}

	cases := []struct {
		name    string
		newName string
		status  int
		content []string
	}{
		{"empty name", "   ", http.StatusBadRequest, []string{`"code":"value_required"`}},
		{"unchanged name", renameValuationName, http.StatusBadRequest, []string{`"code":"name_unchanged"`}},
		{"unchanged name with spaces", "  " + renameValuationName + " ", http.StatusBadRequest, []string{`"code":"name_unchanged"`}},
		{"name of another sheet", "Sheet A", http.StatusConflict, []string{`"status":409`, `"code":"validation_not_unique"`}},
		{"name of another sheet in other capitals", "sheet a", http.StatusConflict, []string{`"code":"validation_not_unique"`}},
	}
	for _, c := range cases {
		scenario := renameScenario(adminToken, "rsvaluation0002", c.newName, c.status, c.content...)
		scenario.Name = c.name
		scenario.TestAppFactory = testutils.SetupTestApp
		scenario.AfterTestFunc = noChange
		scenario.Test(t)
	}
}

func TestRenameRateSheet_CapitalsOfOwnName(t *testing.T) {
	adminToken := renameToken(t, "author@soup.com")
	scenario := renameScenario(adminToken, "rsvaluation0001", "Summary Valuation Fixture", http.StatusOK,
		`"revisions_renamed":2`, `"former_names":["Summary valuation fixture"]`)
	scenario.TestAppFactory = testutils.SetupTestApp
	scenario.Test(t)
}

func TestRenameRateSheet_HistoryAcrossRenamesAndRevisions(t *testing.T) {
	adminToken := renameToken(t, "author@soup.com")
	app := testutils.SetupTestApp(t)

	// "Test Deactivate Others" has revisions rs_deactivate_r0 and r1.
	runOnApp(t, app, renameScenario(adminToken, "rs_deactivate_r0", "Name B", http.StatusOK,
		`"former_names":["Test Deactivate Others"]`))
	runOnApp(t, app, renameScenario(adminToken, "rs_deactivate_r1", "Name C", http.StatusOK,
		`"old_name":"Name B"`, `"former_names":["Name B","Test Deactivate Others"]`))

	want := []string{"Name B", "Test Deactivate Others"}
	for _, id := range []string{"rs_deactivate_r0", "rs_deactivate_r1"} {
		if got := formerNamesOf(t, app, id); !slices.Equal(got, want) {
			t.Errorf("former names of %s = %v, want %v", id, got, want)
		}
	}

	// A revision made after the renames shares the history.
	collection, err := app.FindCollectionByNameOrId("rate_sheets")
	if err != nil {
		t.Fatal(err)
	}
	revision := core.NewRecord(collection)
	revision.Set("name", "Name C")
	revision.Set("revision", 2)
	revision.Set("effective_date", "2025-03-01")
	revision.Set("active", false)
	if err := app.Save(revision); err != nil {
		t.Fatal(err)
	}
	if got := formerNamesOf(t, app, revision.Id); !slices.Equal(got, want) {
		t.Errorf("former names of new revision = %v, want %v", got, want)
	}

	runOnApp(t, app, tests.ApiScenario{
		Name:           "former names route reads the history from the new revision",
		Method:         http.MethodGet,
		URL:            "/api/rate_sheets/" + revision.Id + "/former_names",
		Headers:        map[string]string{"Authorization": renameToken(t, "noclaims@example.com")},
		ExpectedStatus: http.StatusOK,
		ExpectedContent: []string{
			`"name":"Name C"`,
			`"former_names":["Name B","Test Deactivate Others"]`,
		},
	})

	// A rename renames the new revision too.
	runOnApp(t, app, renameScenario(adminToken, "rs_deactivate_r0", "Name D", http.StatusOK,
		`"revisions_renamed":3`, `"former_names":["Name C","Name B","Test Deactivate Others"]`))

	// Another sheet may take a name that this sheet used before. The history
	// stays with the renamed sheet.
	runOnApp(t, app, renameScenario(adminToken, "rs_sheetb_r0", "Name B", http.StatusOK,
		`"former_names":["Sheet B"]`))
	if got := formerNamesOf(t, app, "rs_sheetb_r0"); !slices.Equal(got, []string{"Sheet B"}) {
		t.Errorf("former names of Sheet B = %v", got)
	}
	if got := formerNamesOf(t, app, revision.Id); !slices.Equal(got, []string{"Name C", "Name B", "Test Deactivate Others"}) {
		t.Errorf("former names of Name D = %v", got)
	}
}

func TestRateSheetFormerNames_RouteWithoutHistory(t *testing.T) {
	token := renameToken(t, "noclaims@example.com")
	scenarios := []tests.ApiScenario{
		{
			Name:            "a sheet without renames has no former names",
			Method:          http.MethodGet,
			URL:             "/api/rate_sheets/rs_sheeta_r0/former_names",
			Headers:         map[string]string{"Authorization": token},
			ExpectedStatus:  http.StatusOK,
			ExpectedContent: []string{`"former_names":[]`},
			TestAppFactory:  testutils.SetupTestApp,
		},
		{
			Name:            "unknown sheet returns 404",
			Method:          http.MethodGet,
			URL:             "/api/rate_sheets/nosuchsheet0000/former_names",
			Headers:         map[string]string{"Authorization": token},
			ExpectedStatus:  http.StatusNotFound,
			ExpectedContent: []string{`"status":404`},
			TestAppFactory:  testutils.SetupTestApp,
		},
		{
			Name:            "unauthenticated request is refused",
			Method:          http.MethodGet,
			URL:             "/api/rate_sheets/rs_sheeta_r0/former_names",
			ExpectedStatus:  http.StatusUnauthorized,
			ExpectedContent: []string{`"status":401`},
			TestAppFactory:  testutils.SetupTestApp,
		},
	}
	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

// A sheet renamed back to an earlier name doesn't list its current name.
func TestRenameRateSheet_RenameBack(t *testing.T) {
	adminToken := renameToken(t, "author@soup.com")
	app := testutils.SetupTestApp(t)
	runOnApp(t, app, renameScenario(adminToken, "rs_sheeta_r0", "Sheet Q", http.StatusOK, `"name":"Sheet Q"`))
	runOnApp(t, app, renameScenario(adminToken, "rs_sheeta_r1", "Sheet A", http.StatusOK,
		`"former_names":["Sheet Q"]`))
}
