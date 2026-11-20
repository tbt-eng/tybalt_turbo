package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
	"tybalt/internal/testutils"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

// Fixtures used here (app/test_seed_data/data/computers.csv), besides those
// listed in radiator_test.go:
//   - cmpunassigned01 (FAKE-PC-0003): last user time@test.com, not assigned.
//   - cmpnolastuser01 (FAKE-PC-0004): its last report matched no user.
//   - cmpreassign0001 (FAKE-PC-0005): last user time@test.com, assigned to
//     fatt@mac.com.
//
// it.identity@example.com holds the computers claim, author@soup.com holds
// admin, and time@test.com holds neither.
const (
	unassignedComputer = "cmpunassigned01"
	noLastUserComputer = "cmpnolastuser01"
	reassignComputer   = "cmpreassign0001"
)

type computerTokens struct {
	computers, admin, other string
}

func computerTestTokens(t testing.TB) computerTokens {
	t.Helper()
	var tokens computerTokens
	for _, item := range []struct {
		dst   *string
		email string
	}{
		{&tokens.computers, "it.identity@example.com"},
		{&tokens.admin, "author@soup.com"},
		{&tokens.other, "time@test.com"},
	} {
		token, err := testutils.GenerateRecordToken("users", item.email)
		if err != nil {
			t.Fatal(err)
		}
		*item.dst = token
	}
	return tokens
}

func findComputer(t testing.TB, app core.App, id string) *core.Record {
	t.Helper()
	record, err := app.FindRecordById("computers", id)
	if err != nil {
		t.Fatalf("computer %s: %v", id, err)
	}
	return record
}

// assertRecent fails unless value is a datetime within the last minute.
func assertRecent(t testing.TB, field string, value string) {
	t.Helper()
	parsed, err := time.Parse("2006-01-02 15:04:05.000Z", value)
	if err != nil {
		t.Fatalf("%s = %q: %v", field, value, err)
	}
	if age := time.Since(parsed); age < 0 || age > time.Minute {
		t.Errorf("%s = %s, want now", field, value)
	}
}

// Each action is allowed for the computers claim and for admin, and refused
// for other users and for a request without a token.
func TestComputerActionsAuthorization(t *testing.T) {
	tokens := computerTestTokens(t)

	actions := []struct {
		path     string
		computer string
		check    func(t testing.TB, record *core.Record)
	}{
		{"assign", unassignedComputer, func(t testing.TB, record *core.Record) {
			if got := record.GetString("assigned_to"); got != timeUserID {
				t.Errorf("assigned_to = %q, want %q", got, timeUserID)
			}
		}},
		{"retire", unassignedComputer, func(t testing.TB, record *core.Record) {
			if record.GetString("retired") == "" {
				t.Error("retired is empty")
			}
		}},
		{"unretire", retiredComputer, func(t testing.TB, record *core.Record) {
			if got := record.GetString("retired"); got != "" {
				t.Errorf("retired = %q, want empty", got)
			}
		}},
	}

	for _, action := range actions {
		url := "/api/computers/" + action.computer + "/" + action.path
		allowed := func(name, token string) tests.ApiScenario {
			return tests.ApiScenario{
				Name:            action.path + ": " + name + " is allowed",
				Method:          http.MethodPost,
				URL:             url,
				Headers:         map[string]string{"Authorization": token},
				ExpectedStatus:  http.StatusOK,
				ExpectedContent: []string{`"id":"` + action.computer + `"`},
				TestAppFactory:  testutils.SetupTestApp,
				AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
					action.check(t, findComputer(t, app, action.computer))
				},
			}
		}
		unchanged := func(t testing.TB, app *tests.TestApp, res *http.Response) {
			record := findComputer(t, app, action.computer)
			if action.computer == unassignedComputer {
				if record.GetString("assigned_to") != "" || record.GetString("retired") != "" {
					t.Error("the computer changed")
				}
			} else if record.GetString("retired") == "" {
				t.Error("the computer changed")
			}
		}
		scenarios := []tests.ApiScenario{
			allowed("computers claim", tokens.computers),
			allowed("admin", tokens.admin),
			{
				Name:            action.path + ": a user without the claim is refused",
				Method:          http.MethodPost,
				URL:             url,
				Headers:         map[string]string{"Authorization": tokens.other},
				ExpectedStatus:  http.StatusForbidden,
				ExpectedContent: []string{`"status":403`},
				TestAppFactory:  testutils.SetupTestApp,
				AfterTestFunc:   unchanged,
			},
			{
				Name:            action.path + ": a request without a token is refused",
				Method:          http.MethodPost,
				URL:             url,
				ExpectedStatus:  http.StatusUnauthorized,
				ExpectedContent: []string{`"status":401`},
				TestAppFactory:  testutils.SetupTestApp,
				AfterTestFunc:   unchanged,
			},
		}
		for _, scenario := range scenarios {
			scenario.Test(t)
		}
	}
}

// Assign sets assigned_to to the last user, not to anyone the caller names,
// and assigned_time to now. It replaces an assignment to someone else.
func TestAssignComputerUsesLastUser(t *testing.T) {
	tokens := computerTestTokens(t)

	scenarios := []tests.ApiScenario{
		{
			Name:            "assign replaces the assignee with the last user",
			Method:          http.MethodPost,
			URL:             "/api/computers/" + reassignComputer + "/assign",
			Headers:         map[string]string{"Authorization": tokens.computers},
			ExpectedStatus:  http.StatusOK,
			ExpectedContent: []string{`"assigned_to":"` + timeUserID + `"`, `"assigned_to_name":"Tester Time"`},
			TestAppFactory:  testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				record := findComputer(t, app, reassignComputer)
				if got := record.GetString("assigned_to"); got != timeUserID {
					t.Errorf("assigned_to = %q, want the last user %q", got, timeUserID)
				}
				assertRecent(t, "assigned_time", record.GetString("assigned_time"))
			},
		},
		{
			Name:            "assign with the user the page showed",
			Method:          http.MethodPost,
			URL:             "/api/computers/" + unassignedComputer + "/assign",
			Body:            strings.NewReader(`{"user":"` + timeUserID + `"}`),
			Headers:         map[string]string{"Authorization": tokens.computers, "Content-Type": "application/json"},
			ExpectedStatus:  http.StatusOK,
			ExpectedContent: []string{`"assigned_to":"` + timeUserID + `"`},
			TestAppFactory:  testutils.SetupTestApp,
		},
		{
			Name:            "assign refuses when the last user is not the one the page showed",
			Method:          http.MethodPost,
			URL:             "/api/computers/" + unassignedComputer + "/assign",
			Body:            strings.NewReader(`{"user":"` + fattUserID + `"}`),
			Headers:         map[string]string{"Authorization": tokens.computers, "Content-Type": "application/json"},
			ExpectedStatus:  http.StatusConflict,
			ExpectedContent: []string{`"code":"last_user_changed"`},
			TestAppFactory:  testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				if got := findComputer(t, app, unassignedComputer).GetString("assigned_to"); got != "" {
					t.Errorf("assigned_to = %q, want empty", got)
				}
			},
		},
		{
			Name:            "assign refuses a computer whose last report matched no user",
			Method:          http.MethodPost,
			URL:             "/api/computers/" + noLastUserComputer + "/assign",
			Headers:         map[string]string{"Authorization": tokens.computers},
			ExpectedStatus:  http.StatusBadRequest,
			ExpectedContent: []string{`"code":"no_last_user"`, "matched no user"},
			TestAppFactory:  testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				record := findComputer(t, app, noLastUserComputer)
				if record.GetString("assigned_to") != "" || record.GetString("assigned_time") != "" {
					t.Error("the computer was assigned")
				}
			},
		},
		{
			Name:            "assign refuses a retired computer",
			Method:          http.MethodPost,
			URL:             "/api/computers/" + retiredComputer + "/assign",
			Headers:         map[string]string{"Authorization": tokens.admin},
			ExpectedStatus:  http.StatusBadRequest,
			ExpectedContent: []string{`"code":"computer_retired"`},
			TestAppFactory:  testutils.SetupTestApp,
		},
		{
			Name:            "assign refuses an unknown computer",
			Method:          http.MethodPost,
			URL:             "/api/computers/nosuchcomputer1/assign",
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

func TestRetireAndUnretireComputer(t *testing.T) {
	tokens := computerTestTokens(t)

	scenarios := []tests.ApiScenario{
		{
			Name:            "retire sets retired to now",
			Method:          http.MethodPost,
			URL:             "/api/computers/" + assignedComputer + "/retire",
			Headers:         map[string]string{"Authorization": tokens.computers},
			ExpectedStatus:  http.StatusOK,
			ExpectedContent: []string{`"id":"` + assignedComputer + `"`},
			TestAppFactory:  testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				record := findComputer(t, app, assignedComputer)
				assertRecent(t, "retired", record.GetString("retired"))
				// Retiring keeps the assignment.
				if got := record.GetString("assigned_to"); got != fattUserID {
					t.Errorf("assigned_to = %q, want %q", got, fattUserID)
				}
			},
		},
		{
			Name:            "retire refuses a retired computer",
			Method:          http.MethodPost,
			URL:             "/api/computers/" + retiredComputer + "/retire",
			Headers:         map[string]string{"Authorization": tokens.computers},
			ExpectedStatus:  http.StatusBadRequest,
			ExpectedContent: []string{`"code":"already_retired"`},
			TestAppFactory:  testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				if got := findComputer(t, app, retiredComputer).GetString("retired"); got != "2026-06-01 00:00:00.000Z" {
					t.Errorf("retired = %q, want it unchanged", got)
				}
			},
		},
		{
			Name:            "un-retire refuses a computer that is not retired",
			Method:          http.MethodPost,
			URL:             "/api/computers/" + unassignedComputer + "/unretire",
			Headers:         map[string]string{"Authorization": tokens.computers},
			ExpectedStatus:  http.StatusBadRequest,
			ExpectedContent: []string{`"code":"not_retired"`},
			TestAppFactory:  testutils.SetupTestApp,
		},
		{
			Name:            "retire refuses an unknown computer",
			Method:          http.MethodPost,
			URL:             "/api/computers/nosuchcomputer1/retire",
			Headers:         map[string]string{"Authorization": tokens.admin},
			ExpectedStatus:  http.StatusNotFound,
			ExpectedContent: []string{`"code":"not_found"`},
			TestAppFactory:  testutils.SetupTestApp,
		},
	}
	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

// The lists split the computers on retired, and carry the names of the last
// user and of the assignee, which the claim holder can't read elsewhere.
func TestComputersListAndDetails(t *testing.T) {
	tokens := computerTestTokens(t)

	listIDs := func(t testing.TB, res *http.Response) []string {
		t.Helper()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}
		var rows []struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal(body, &rows); err != nil {
			t.Fatalf("%v: %s", err, body)
		}
		ids := make([]string, 0, len(rows))
		for _, row := range rows {
			ids = append(ids, row.ID)
		}
		return ids
	}

	scenarios := []tests.ApiScenario{
		{
			Name:           "the list holds the computers that are not retired",
			Method:         http.MethodGet,
			URL:            "/api/computers",
			Headers:        map[string]string{"Authorization": tokens.computers},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"assigned_to_name":"Fatty Maclean"`,
				`"last_user_name":"Tester Time"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				got := strings.Join(listIDs(t, res), ",")
				want := strings.Join([]string{assignedComputer, unassignedComputer, noLastUserComputer, reassignComputer}, ",")
				if got != want {
					t.Errorf("ids = %s, want %s", got, want)
				}
			},
		},
		{
			Name:            "the retired list holds the retired computers",
			Method:          http.MethodGet,
			URL:             "/api/computers?retired=1",
			Headers:         map[string]string{"Authorization": tokens.admin},
			ExpectedStatus:  http.StatusOK,
			ExpectedContent: []string{`"retired":"2026-06-01 00:00:00.000Z"`},
			TestAppFactory:  testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				if got := strings.Join(listIDs(t, res), ","); got != retiredComputer {
					t.Errorf("ids = %s, want %s", got, retiredComputer)
				}
			},
		},
		{
			Name:           "details include the network config",
			Method:         http.MethodGet,
			URL:            "/api/computers/" + unassignedComputer,
			Headers:        map[string]string{"Authorization": tokens.computers},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"computer_name":"FAKE-PC-0003"`,
				`"00:11:22:33:44:67":{`,
				`"ram":34359738368`,
			},
			TestAppFactory: testutils.SetupTestApp,
		},
		{
			Name:            "details of a retired computer",
			Method:          http.MethodGet,
			URL:             "/api/computers/" + retiredComputer,
			Headers:         map[string]string{"Authorization": tokens.computers},
			ExpectedStatus:  http.StatusOK,
			ExpectedContent: []string{`"retired":"2026-06-01 00:00:00.000Z"`},
			TestAppFactory:  testutils.SetupTestApp,
		},
		{
			Name:            "details of an unknown computer",
			Method:          http.MethodGet,
			URL:             "/api/computers/nosuchcomputer1",
			Headers:         map[string]string{"Authorization": tokens.computers},
			ExpectedStatus:  http.StatusNotFound,
			ExpectedContent: []string{`"code":"not_found"`},
			TestAppFactory:  testutils.SetupTestApp,
		},
		{
			Name:            "a user without the claim can't list",
			Method:          http.MethodGet,
			URL:             "/api/computers",
			Headers:         map[string]string{"Authorization": tokens.other},
			ExpectedStatus:  http.StatusForbidden,
			ExpectedContent: []string{`"status":403`},
			TestAppFactory:  testutils.SetupTestApp,
		},
		{
			Name:            "a user without the claim can't see details",
			Method:          http.MethodGet,
			URL:             "/api/computers/" + unassignedComputer,
			Headers:         map[string]string{"Authorization": tokens.other},
			ExpectedStatus:  http.StatusForbidden,
			ExpectedContent: []string{`"status":403`},
			TestAppFactory:  testutils.SetupTestApp,
		},
		{
			Name:            "a request without a token can't list",
			Method:          http.MethodGet,
			URL:             "/api/computers",
			ExpectedStatus:  http.StatusUnauthorized,
			ExpectedContent: []string{`"status":401`},
			TestAppFactory:  testutils.SetupTestApp,
		},
	}
	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
