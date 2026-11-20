package main

import (
	"net/http"
	"strings"
	"testing"

	"tybalt/internal/testutils"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/tests"
)

// Delete of a rate sheet (DELETE /api/rate_sheets/{id}). The fixtures are rows
// of the seed CSV files (app/test_seed_data/data):
//   - rsdelunused0001: Delete Unused, inactive, two entries, nothing uses it
//   - rsdelrev0r00001, rsdelrev1r00001, rsdelrev2r00001: Delete Revisions
//     revisions 0, 1 and 2, inactive, nothing uses them
//   - rsdelrequse0001: inactive, the project request rsdelreq0000001 uses it
//   - rsdelbased00001: inactive, rsdelchild00001 is based on it
//   - rsdelnamechg001: inactive, with a row in rate_sheet_name_changes.csv
//   - rsdelcodechg001: inactive, with a row in rate_sheet_code_changes.csv
//
// Other sheets used here were seeded for other tests: rs_sheetb_r0 (active,
// unused), rsvaluation0001 (inactive; jobs use it, now and in their history),
// rsabhistrev0r01 (inactive; only an older period of a job uses it) and
// rsdraftsheet001 (a project request's proposed rates).
//
// Seeded users:
//   - admin.only@example.com holds admin only
//   - author@soup.com holds admin and job
//   - invoice-job@example.com holds job, not admin
//   - noclaims@example.com holds no claims

const (
	rsDelUnused   = "rsdelunused0001"
	rsDelUnusedCd = "2026-11-00038"
	rsDelRev0     = "rsdelrev0r00001"
	rsDelRev1     = "rsdelrev1r00001"
	rsDelRev2     = "rsdelrev2r00001"
)

func deleteRateSheetScenario(name, id, email string, status int, content []string, after func(t testing.TB, app *tests.TestApp)) tests.ApiScenario {
	headers := map[string]string{}
	if email != "" {
		token, err := testutils.GenerateRecordToken("users", email)
		if err != nil {
			panic(err)
		}
		headers["Authorization"] = token
	}
	return tests.ApiScenario{
		Name:            name,
		Method:          http.MethodDelete,
		URL:             "/api/rate_sheets/" + id,
		Headers:         headers,
		ExpectedStatus:  status,
		ExpectedContent: content,
		TestAppFactory:  testutils.SetupTestApp,
		AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
			if after != nil {
				after(t, app)
			}
		},
	}
}

func expectRateSheetKept(id string) func(t testing.TB, app *tests.TestApp) {
	return func(t testing.TB, app *tests.TestApp) {
		if !rateSheetExists(t, app, id) {
			t.Fatalf("rate sheet %s was deleted; want it kept", id)
		}
	}
}

func TestDeleteRateSheet_AdminDeletesUnusedSheet(t *testing.T) {
	scenario := deleteRateSheetScenario(
		"admin deletes an inactive sheet that nothing uses",
		rsDelUnused, "admin.only@example.com", http.StatusOK,
		[]string{`"message":"rate sheet deleted"`, `"id":"` + rsDelUnused + `"`},
		func(t testing.TB, app *tests.TestApp) {
			if rateSheetExists(t, app, rsDelUnused) {
				t.Fatalf("rate sheet %s still exists", rsDelUnused)
			}
			// rate_sheet_entries.rate_sheet has cascadeDelete.
			entries, err := app.CountRecords("rate_sheet_entries", dbx.HashExp{"rate_sheet": rsDelUnused})
			if err != nil {
				t.Fatal(err)
			}
			if entries != 0 {
				t.Fatalf("%d entries remain; want 0", entries)
			}
			// The code is never issued again, so its row stays.
			codes, err := app.CountRecords("rate_sheet_codes", dbx.HashExp{"code": rsDelUnusedCd})
			if err != nil {
				t.Fatal(err)
			}
			if codes != 1 {
				t.Fatalf("rate_sheet_codes has %d rows for %s; want 1", codes, rsDelUnusedCd)
			}
		},
	)
	scenario.Test(t)
}

// Any revision can be deleted, not only the latest.
func TestDeleteRateSheet_AnyRevision(t *testing.T) {
	keepOthers := func(deleted string, kept ...string) func(t testing.TB, app *tests.TestApp) {
		return func(t testing.TB, app *tests.TestApp) {
			if rateSheetExists(t, app, deleted) {
				t.Fatalf("rate sheet %s still exists", deleted)
			}
			for _, id := range kept {
				if !rateSheetExists(t, app, id) {
					t.Fatalf("rate sheet %s was deleted; want it kept", id)
				}
			}
		}
	}
	scenarios := []tests.ApiScenario{
		deleteRateSheetScenario("deletes the first revision while later revisions exist",
			rsDelRev0, "author@soup.com", http.StatusOK, []string{`"rate sheet deleted"`},
			keepOthers(rsDelRev0, rsDelRev1, rsDelRev2)),
		deleteRateSheetScenario("deletes a middle revision",
			rsDelRev1, "author@soup.com", http.StatusOK, []string{`"rate sheet deleted"`},
			keepOthers(rsDelRev1, rsDelRev0, rsDelRev2)),
		deleteRateSheetScenario("deletes the latest revision",
			rsDelRev2, "author@soup.com", http.StatusOK, []string{`"rate sheet deleted"`},
			keepOthers(rsDelRev2, rsDelRev0, rsDelRev1)),
	}
	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestDeleteRateSheet_RequiresAdmin(t *testing.T) {
	scenarios := []tests.ApiScenario{
		deleteRateSheetScenario("a job claim holder can't delete",
			rsDelUnused, "invoice-job@example.com", http.StatusForbidden,
			[]string{"You do not have permission to delete rate sheets."}, expectRateSheetKept(rsDelUnused)),
		deleteRateSheetScenario("a user without claims can't delete",
			rsDelUnused, "noclaims@example.com", http.StatusForbidden,
			[]string{"You do not have permission to delete rate sheets."}, expectRateSheetKept(rsDelUnused)),
		deleteRateSheetScenario("a guest can't delete",
			rsDelUnused, "", http.StatusUnauthorized, []string{`"data":{}`}, expectRateSheetKept(rsDelUnused)),
	}
	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestDeleteRateSheet_Refusals(t *testing.T) {
	cases := []struct {
		name   string
		id     string
		status int
		code   string
	}{
		{"an active sheet", "rs_sheetb_r0", http.StatusBadRequest, "rate_sheet_active"},
		{"a project request's proposed rates", "rsdraftsheet001", http.StatusBadRequest, "rate_sheet_draft"},
		{"a sheet that jobs use", "rsvaluation0001", http.StatusBadRequest, "rate_sheet_used_by_job"},
		{"a sheet that only a past period of a job uses", "rsabhistrev0r01", http.StatusBadRequest, "rate_sheet_used_by_job"},
		{"a sheet that a project request uses", "rsdelrequse0001", http.StatusBadRequest, "rate_sheet_used_by_request"},
		{"a sheet that another sheet is based on", "rsdelbased00001", http.StatusBadRequest, "rate_sheet_based_on"},
		{"a sheet with a former name", "rsdelnamechg001", http.StatusBadRequest, "rate_sheet_has_name_history"},
		{"a sheet with a former code", "rsdelcodechg001", http.StatusBadRequest, "rate_sheet_has_code_history"},
		{"a missing sheet", "rsdelmissing001", http.StatusNotFound, "rate_sheet_not_found"},
	}
	for _, c := range cases {
		after := expectRateSheetKept(c.id)
		if c.status == http.StatusNotFound {
			after = nil
		}
		scenario := deleteRateSheetScenario("refuses "+c.name, c.id, "admin.only@example.com", c.status,
			[]string{`"code":"` + c.code + `"`}, after)
		scenario.Test(t)
	}
}

// A pending rate sheet absorb blocks every delete until it is committed or
// undone. The test absorbs Absorb Early into Absorb Inactive Target through
// the API (the same pair TestRateSheetAbsorb uses), which writes an
// absorb_actions row at test time, then tries to delete the target, which is
// inactive and unused.
func TestDeleteRateSheet_PendingAbsorb(t *testing.T) {
	app := testutils.SetupTestApp(t)
	defer app.Cleanup()

	status, body := absorbRequest(t, app, http.MethodPost, "/api/rate_sheets/"+rsAbInactiveTgt+"/absorb", absorbAdminEmail, absorbBody(rsAbEarly))
	if status != http.StatusOK {
		t.Fatalf("absorb: status %d: %s", status, body)
	}

	status, body = absorbRequest(t, app, http.MethodDelete, "/api/rate_sheets/"+rsAbInactiveTgt, absorbAdminEmail, nil)
	if status != http.StatusConflict || !strings.Contains(body, `"code":"rate_sheet_absorb_pending"`) {
		t.Fatalf("delete during a pending absorb: status %d: %s", status, body)
	}
	if !rateSheetExists(t, app, rsAbInactiveTgt) {
		t.Fatalf("rate sheet %s was deleted during a pending absorb", rsAbInactiveTgt)
	}
}
