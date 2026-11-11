package main

import (
	"net/http"
	"testing"
	"tybalt/internal/testutils"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

// A time sheet rejected after approval stays on Approved By Me. Its row must
// carry the rejection, with the rejector's name, so that the screen can show it.
func TestTimeSheetTalliesApprovedIncludesRejection(t *testing.T) {
	approverToken, err := testutils.GenerateRecordToken("users", "fakemanager@fakesite.xyz")
	if err != nil {
		t.Fatal(err)
	}

	scenario := tests.ApiScenario{
		Name:   "approved tallies include the rejection and the rejector name",
		Method: http.MethodGet,
		URL:    "/api/time_sheets/tallies/approved",
		Headers: map[string]string{
			"Authorization": approverToken,
		},
		BeforeTestFunc: func(t testing.TB, app *tests.TestApp, _ *core.ServeEvent) {
			// ctw2sheet000001 is approved by wegviunlyr2jjjv (Fakesy Manjor). Make it
			// rejected after approval and not committed. SQL skips the hooks.
			_, err := app.DB().NewQuery(`UPDATE time_sheets
				SET committed = '', committer = '',
					rejected = '2026-04-26 09:00:00.000Z',
					rejection_reason = 'Wrong job on Tuesday',
					rejector = 'wegviunlyr2jjjv'
				WHERE id = 'ctw2sheet000001'`).Execute()
			if err != nil {
				t.Fatal(err)
			}
		},
		ExpectedStatus: http.StatusOK,
		ExpectedContent: []string{
			`"id":"ctw2sheet000001"`,
			`"rejected":"2026-04-26 09:00:00.000Z"`,
			`"rejection_reason":"Wrong job on Tuesday"`,
			`"rejector_name":"Fakesy Manjor"`,
		},
		TestAppFactory: testutils.SetupTestApp,
	}
	scenario.Test(t)
}
