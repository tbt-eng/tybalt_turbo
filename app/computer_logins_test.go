package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"tybalt/internal/testutils"

	"github.com/pocketbase/pocketbase/tests"
)

// Fixtures used here (app/test_seed_data/data/), besides those listed in
// radiator_test.go and computers_test.go:
//   - computer_logins: cmplogin0000002 to cmplogin0000023, 22 sign-ins on
//     cmpunassigned01 (FAKE-PC-0003) six hours apart, the newest at
//     2026-11-10 13:00. They matched time@test.com (Tester Time), except
//     cmplogin0000004 and cmplogin0000011 (Guest Account), which matched no
//     user. cmplogin0000024 is on cmpnolastuser01 (FAKE-PC-0004, no user) and
//     cmplogin0000025, the newest of all, on cmpreassign0001 (FAKE-PC-0005).
//     With cmplogin0000001 there are 25.
//   - cmplogin0000026, the oldest, on cmpnolastuser01. Its matched user,
//     uabsorbrevise01, has no profile and no name, so the row shows the email.
//   - computer_raw_logins: cmprawlogin0003 (FAKE-PC-0098, its report holds
//     markup) and cmprawlogin0004 (no computer name). With the two of
//     radiator_test.go there are 4.
const (
	fixtureLogins        = 26
	fixtureRawLogins     = 4
	unassignedLogins     = 22
	newestLogin          = "cmplogin0000025"
	oldestLogin          = "cmplogin0000026"
	newestUnassigned     = "cmplogin0000023"
	markupRawLogin       = "cmprawlogin0003"
	namelessRawLogin     = "cmprawlogin0004"
	firstFixtureRawLogin = "cmprawlogin0001"
)

type loginsPage struct {
	Items []struct {
		ID           string `json:"id"`
		Computer     string `json:"computer"`
		ComputerName string `json:"computer_name"`
		GivenName    string `json:"given_name"`
		Surname      string `json:"surname"`
		User         string `json:"user"`
		UserName     string `json:"user_name"`
		Created      string `json:"created"`
		Error        string `json:"error"`
		// Present only to check that the list leaves it out.
		Report json.RawMessage `json:"report"`
	} `json:"items"`
	Page       int `json:"page"`
	PerPage    int `json:"perPage"`
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
}

func readLoginsPage(t testing.TB, res *http.Response) loginsPage {
	t.Helper()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	var page loginsPage
	if err := json.Unmarshal(body, &page); err != nil {
		t.Fatalf("%v: %s", err, body)
	}
	return page
}

func (p loginsPage) ids() []string {
	ids := make([]string, 0, len(p.Items))
	for _, item := range p.Items {
		ids = append(ids, item.ID)
	}
	return ids
}

// checkPage fails unless the page numbers are as given.
func checkPage(t testing.TB, got loginsPage, page, perPage, totalItems, totalPages, items int) {
	t.Helper()
	if got.Page != page || got.PerPage != perPage || got.TotalItems != totalItems ||
		got.TotalPages != totalPages || len(got.Items) != items {
		t.Errorf("page %d perPage %d totalItems %d totalPages %d items %d, want %d %d %d %d %d",
			got.Page, got.PerPage, got.TotalItems, got.TotalPages, len(got.Items),
			page, perPage, totalItems, totalPages, items)
	}
}

// Every route is allowed for the computers claim and for admin, and refused
// for other users and for a request without a token.
func TestComputerLoginRoutesAuthorization(t *testing.T) {
	tokens := computerTestTokens(t)

	routes := []struct {
		name    string
		method  string
		url     string
		status  int
		content []string
	}{
		{"sign-ins", http.MethodGet, "/api/computers/logins", http.StatusOK, []string{`"totalItems":26`}},
		{"rejected reports", http.MethodGet, "/api/computers/raw_logins", http.StatusOK, []string{`"totalItems":4`}},
		{"one rejected report", http.MethodGet, "/api/computers/raw_logins/" + markupRawLogin, http.StatusOK, []string{`"id":"` + markupRawLogin + `"`}},
		{"delete a rejected report", http.MethodDelete, "/api/computers/raw_logins/" + markupRawLogin, http.StatusNoContent, nil},
	}

	rawLoginKept := func(t testing.TB, app *tests.TestApp, res *http.Response) {
		if _, err := app.FindRecordById("computer_raw_logins", markupRawLogin); err != nil {
			t.Errorf("the rejected report is gone: %v", err)
		}
	}

	for _, route := range routes {
		allowed := func(name, token string) tests.ApiScenario {
			scenario := tests.ApiScenario{
				Name:            route.name + ": " + name + " is allowed",
				Method:          route.method,
				URL:             route.url,
				Headers:         map[string]string{"Authorization": token},
				ExpectedStatus:  route.status,
				ExpectedContent: route.content,
				TestAppFactory:  testutils.SetupTestApp,
			}
			if route.method == http.MethodDelete {
				scenario.AfterTestFunc = func(t testing.TB, app *tests.TestApp, res *http.Response) {
					if _, err := app.FindRecordById("computer_raw_logins", markupRawLogin); err == nil {
						t.Error("the rejected report was not deleted")
					}
				}
			}
			return scenario
		}
		scenarios := []tests.ApiScenario{
			allowed("computers claim", tokens.computers),
			allowed("admin", tokens.admin),
			{
				Name:            route.name + ": a user without the claim is refused",
				Method:          route.method,
				URL:             route.url,
				Headers:         map[string]string{"Authorization": tokens.other},
				ExpectedStatus:  http.StatusForbidden,
				ExpectedContent: []string{`"status":403`},
				TestAppFactory:  testutils.SetupTestApp,
				AfterTestFunc:   rawLoginKept,
			},
			{
				Name:            route.name + ": a request without a token is refused",
				Method:          route.method,
				URL:             route.url,
				ExpectedStatus:  http.StatusUnauthorized,
				ExpectedContent: []string{`"status":401`},
				TestAppFactory:  testutils.SetupTestApp,
				AfterTestFunc:   rawLoginKept,
			},
		}
		for _, scenario := range scenarios {
			scenario.Test(t)
		}
	}
}

// The sign-ins come a page at a time, newest first, with the computer name
// and the name of the matched user.
func TestComputerLoginsPaging(t *testing.T) {
	tokens := computerTestTokens(t)

	scenarios := []tests.ApiScenario{
		{
			// /logins must not reach the route of one computer, which would
			// answer 404 for a computer with the id "logins".
			Name:            "the first page holds the newest sign-ins",
			Method:          http.MethodGet,
			URL:             "/api/computers/logins?perPage=10",
			Headers:         map[string]string{"Authorization": tokens.computers},
			ExpectedStatus:  http.StatusOK,
			ExpectedContent: []string{`"totalItems":26`},
			TestAppFactory:  testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				page := readLoginsPage(t, res)
				checkPage(t, page, 1, 10, fixtureLogins, 3, 10)
				if len(page.Items) == 0 {
					return
				}
				first := page.Items[0]
				if first.ID != newestLogin || first.ComputerName != "FAKE-PC-0005" ||
					first.Computer != reassignComputer || first.User != timeUserID ||
					first.UserName != "Tester Time" {
					t.Errorf("first item = %+v", first)
				}
				for i := 1; i < len(page.Items); i++ {
					if page.Items[i-1].Created < page.Items[i].Created {
						t.Errorf("item %d is newer than item %d", i, i-1)
					}
				}
			},
		},
		{
			Name:            "the last page holds the rest",
			Method:          http.MethodGet,
			URL:             "/api/computers/logins?perPage=10&page=3",
			Headers:         map[string]string{"Authorization": tokens.admin},
			ExpectedStatus:  http.StatusOK,
			ExpectedContent: []string{`"totalItems":26`},
			TestAppFactory:  testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				page := readLoginsPage(t, res)
				checkPage(t, page, 3, 10, fixtureLogins, 3, 6)
				if ids := page.ids(); len(ids) > 0 && ids[len(ids)-1] != oldestLogin {
					t.Errorf("last id = %s, want %s", ids[len(ids)-1], oldestLogin)
				}
			},
		},
		{
			Name:            "a page past the end gives the last page",
			Method:          http.MethodGet,
			URL:             "/api/computers/logins?perPage=10&page=99",
			Headers:         map[string]string{"Authorization": tokens.computers},
			ExpectedStatus:  http.StatusOK,
			ExpectedContent: []string{`"page":3`},
			TestAppFactory:  testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				checkPage(t, readLoginsPage(t, res), 3, 10, fixtureLogins, 3, 6)
			},
		},
		{
			Name:            "bad paging values get the defaults and perPage is at most 100",
			Method:          http.MethodGet,
			URL:             "/api/computers/logins?perPage=1000&page=-2",
			Headers:         map[string]string{"Authorization": tokens.computers},
			ExpectedStatus:  http.StatusOK,
			ExpectedContent: []string{`"perPage":100`},
			TestAppFactory:  testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				checkPage(t, readLoginsPage(t, res), 1, 100, fixtureLogins, 1, fixtureLogins)
			},
		},
		{
			Name:            "the default page size is 20",
			Method:          http.MethodGet,
			URL:             "/api/computers/logins?perPage=abc",
			Headers:         map[string]string{"Authorization": tokens.computers},
			ExpectedStatus:  http.StatusOK,
			ExpectedContent: []string{`"perPage":20`},
			TestAppFactory:  testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				checkPage(t, readLoginsPage(t, res), 1, 20, fixtureLogins, 2, 20)
			},
		},
	}
	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

// The details page shows the newest 20 sign-ins of one computer.
func TestComputerLoginsOfOneComputer(t *testing.T) {
	tokens := computerTestTokens(t)

	scenarios := []tests.ApiScenario{
		{
			Name:            "the newest 20 sign-ins of one computer",
			Method:          http.MethodGet,
			URL:             "/api/computers/logins?computer=" + unassignedComputer + "&perPage=20",
			Headers:         map[string]string{"Authorization": tokens.computers},
			ExpectedStatus:  http.StatusOK,
			ExpectedContent: []string{`"totalItems":22`},
			TestAppFactory:  testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				page := readLoginsPage(t, res)
				checkPage(t, page, 1, 20, unassignedLogins, 2, 20)
				unmatched := 0
				for _, item := range page.Items {
					if item.Computer != unassignedComputer {
						t.Errorf("%s is on %s", item.ID, item.Computer)
					}
					if item.User == "" {
						unmatched++
						if item.UserName != "" || item.GivenName != "Guest" {
							t.Errorf("unmatched item = %+v", item)
						}
					} else if item.UserName != "Tester Time" {
						t.Errorf("user_name = %q, want Tester Time", item.UserName)
					}
				}
				if unmatched != 2 {
					t.Errorf("unmatched = %d, want 2", unmatched)
				}
				if ids := page.ids(); len(ids) > 0 && ids[0] != newestUnassigned {
					t.Errorf("first id = %s, want %s", ids[0], newestUnassigned)
				}
			},
		},
		{
			Name:            "an unknown computer has no sign-ins",
			Method:          http.MethodGet,
			URL:             "/api/computers/logins?computer=nosuchcomputer1",
			Headers:         map[string]string{"Authorization": tokens.computers},
			ExpectedStatus:  http.StatusOK,
			ExpectedContent: []string{`"items":[]`, `"totalItems":0`, `"page":1`},
			TestAppFactory:  testutils.SetupTestApp,
		},
	}
	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

// The search matches the computer, the reported name and the matched user,
// and takes % and _ literally.
func TestComputerLoginsSearch(t *testing.T) {
	tokens := computerTestTokens(t)

	search := func(q string, want int) tests.ApiScenario {
		return tests.ApiScenario{
			Name:            "search " + q,
			Method:          http.MethodGet,
			URL:             "/api/computers/logins?perPage=100&q=" + q,
			Headers:         map[string]string{"Authorization": tokens.computers},
			ExpectedStatus:  http.StatusOK,
			ExpectedContent: []string{`"totalItems":`},
			TestAppFactory:  testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				if got := readLoginsPage(t, res).TotalItems; got != want {
					t.Errorf("totalItems = %d, want %d", got, want)
				}
			},
		}
	}
	scenarios := []tests.ApiScenario{
		search("guest%20account", 2),
		search("FAKE-PC-0004", 2),
		search("fakeserial0005", 1),
		// Fatty Maclean is the matched user of cmplogin0000001, which reported
		// the name Fake Assignee.
		search("maclean", 1),
		search("tester%20time", 21),
		search("%25", 0),
		// The matched user of cmplogin0000026 has no name: the row shows the
		// email, and the search finds it.
		search("absorb.revise", 1),
		search("_", 0),
	}
	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

// The rejected reports come a page at a time, newest first, without the
// report. One rejected report carries the report as it was received.
func TestComputerRawLogins(t *testing.T) {
	tokens := computerTestTokens(t)

	scenarios := []tests.ApiScenario{
		{
			Name:            "the list leaves out the reports",
			Method:          http.MethodGet,
			URL:             "/api/computers/raw_logins",
			Headers:         map[string]string{"Authorization": tokens.computers},
			ExpectedStatus:  http.StatusOK,
			ExpectedContent: []string{`"error":"upn: is required"`},
			NotExpectedContent: []string{
				`"report"`,
				`FAKESERIAL`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				page := readLoginsPage(t, res)
				checkPage(t, page, 1, 20, fixtureRawLogins, 1, fixtureRawLogins)
				got := strings.Join(page.ids(), ",")
				want := strings.Join([]string{namelessRawLogin, markupRawLogin, "cmprawlogin0002", firstFixtureRawLogin}, ",")
				if got != want {
					t.Errorf("ids = %s, want %s", got, want)
				}
			},
		},
		{
			Name:            "a page of one",
			Method:          http.MethodGet,
			URL:             "/api/computers/raw_logins?perPage=1&page=2",
			Headers:         map[string]string{"Authorization": tokens.computers},
			ExpectedStatus:  http.StatusOK,
			ExpectedContent: []string{`"id":"` + markupRawLogin + `"`, `"totalPages":4`},
			TestAppFactory:  testutils.SetupTestApp,
		},
		{
			Name:               "search the computer name and the error",
			Method:             http.MethodGet,
			URL:                "/api/computers/raw_logins?q=usersurname",
			Headers:            map[string]string{"Authorization": tokens.computers},
			ExpectedStatus:     http.StatusOK,
			ExpectedContent:    []string{`"totalItems":1`, `"computer_name":"FAKE-PC-0098"`},
			NotExpectedContent: []string{"FAKE-PC-0099"},
			TestAppFactory:     testutils.SetupTestApp,
		},
		{
			Name:            "one rejected report holds the report",
			Method:          http.MethodGet,
			URL:             "/api/computers/raw_logins/" + markupRawLogin,
			Headers:         map[string]string{"Authorization": tokens.admin},
			ExpectedStatus:  http.StatusOK,
			ExpectedContent: []string{`"computer_name":"FAKE-PC-0098"`, `"error":"userSurname: is required"`},
			TestAppFactory:  testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				body, err := io.ReadAll(res.Body)
				if err != nil {
					t.Fatal(err)
				}
				var row struct {
					Report map[string]any `json:"report"`
				}
				if err := json.Unmarshal(body, &row); err != nil {
					t.Fatalf("%v: %s", err, body)
				}
				if got := row.Report["userGivenName"]; got != "<b>Not bold</b>" {
					t.Errorf("report userGivenName = %v", got)
				}
				if got := row.Report["serial"]; got != "FAKESERIAL0098" {
					t.Errorf("report serial = %v", got)
				}
			},
		},
		{
			Name:            "an unknown rejected report",
			Method:          http.MethodGet,
			URL:             "/api/computers/raw_logins/nosuchrawlogin1",
			Headers:         map[string]string{"Authorization": tokens.computers},
			ExpectedStatus:  http.StatusNotFound,
			ExpectedContent: []string{`"code":"not_found"`},
			TestAppFactory:  testutils.SetupTestApp,
		},
	}
	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

// A delete removes one rejected report and nothing else.
func TestDeleteComputerRawLogin(t *testing.T) {
	tokens := computerTestTokens(t)

	scenarios := []tests.ApiScenario{
		{
			Name:           "delete removes only that rejected report",
			Method:         http.MethodDelete,
			URL:            "/api/computers/raw_logins/" + namelessRawLogin,
			Headers:        map[string]string{"Authorization": tokens.computers},
			ExpectedStatus: http.StatusNoContent,
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				if _, err := app.FindRecordById("computer_raw_logins", namelessRawLogin); err == nil {
					t.Error("the rejected report was not deleted")
				}
				if got := radiatorCount(t, app, "computer_raw_logins"); got != fixtureRawLogins-1 {
					t.Errorf("computer_raw_logins = %d, want %d", got, fixtureRawLogins-1)
				}
				if got := radiatorCount(t, app, "computer_logins"); got != fixtureLogins {
					t.Errorf("computer_logins = %d, want %d", got, fixtureLogins)
				}
			},
		},
		{
			Name:            "delete of an unknown rejected report",
			Method:          http.MethodDelete,
			URL:             "/api/computers/raw_logins/nosuchrawlogin1",
			Headers:         map[string]string{"Authorization": tokens.admin},
			ExpectedStatus:  http.StatusNotFound,
			ExpectedContent: []string{`"code":"not_found"`},
			TestAppFactory:  testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				if got := radiatorCount(t, app, "computer_raw_logins"); got != fixtureRawLogins {
					t.Errorf("computer_raw_logins = %d, want %d", got, fixtureRawLogins)
				}
			},
		},
		{
			// A sign-in is not a rejected report: the id of a login is unknown
			// to the delete.
			Name:            "delete refuses the id of a sign-in",
			Method:          http.MethodDelete,
			URL:             "/api/computers/raw_logins/" + oldestLogin,
			Headers:         map[string]string{"Authorization": tokens.admin},
			ExpectedStatus:  http.StatusNotFound,
			ExpectedContent: []string{`"code":"not_found"`},
			TestAppFactory:  testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				if got := radiatorCount(t, app, "computer_logins"); got != fixtureLogins {
					t.Errorf("computer_logins = %d, want %d", got, fixtureLogins)
				}
			},
		},
	}
	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
