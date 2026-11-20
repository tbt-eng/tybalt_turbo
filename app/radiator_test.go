package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"tybalt/internal/testutils"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

// Fixtures used here (app/test_seed_data/data/):
//   - machine_secrets: radiatorkey0001 (role radiator, token below, unexpired),
//     radiatorkey0002 (role radiator, expired) and machinekey00001
//     (legacy_writeback, token test-secret-123).
//   - admin_profiles: user_source_anchor aaaa... on time@test.com
//     (rzr98oadsp9qc11) and bbbb... on hal@2005.com (66ct66w380ob6w8).
//   - computers: cmpassigned0001 (FAKE-PC-0001, assigned to fatt@mac.com) and
//     cmpretired00001 (FAKE-PC-0002, retired). The tests of computers_test.go
//     add FAKE-PC-0003 to FAKE-PC-0005, which these tests only count.
//   - computer_logins: one login of cmpassigned0001.
//   - computer_raw_logins: one reject for FAKE-PC-0001, one for FAKE-PC-0099.
const (
	radiatorToken        = "radiator-test-secret"
	radiatorExpiredToken = "radiator-expired-secret"
	legacyWritebackToken = "test-secret-123"

	radiatorFixtureComputers = 5
	radiatorFixtureLogins    = 1
	radiatorFixtureRawLogins = 2

	anchorTimeUser   = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	timeUserID       = "rzr98oadsp9qc11"
	fattUserID       = "etysnrlup2f6bak"
	halUserID        = "66ct66w380ob6w8"
	unmatchedAnchor  = "cccccccccccccccccccccccccccccccc"
	assignedComputer = "cmpassigned0001"
	retiredComputer  = "cmpretired00001"
)

// radiatorReport returns a valid report, shaped as radiator 9 sends it, for a
// computer and a user that are not in the fixtures. overrides replaces or adds
// properties; a nil value in overrides removes the property.
func radiatorReport(overrides map[string]any) string {
	report := map[string]any{
		"radiatorVersion":  9,
		"userGivenName":    "Fake",
		"userSurname":      "Person",
		"email":            "Nobody.Here@Example.com",
		"upn":              "Nobody.Here@Example.com",
		"userSourceAnchor": strings.ToUpper(unmatchedAnchor),
		"serial":           "FAKENEW0001",
		"osVersion":        "10.0.19045",
		"osSku":            48,
		"osArch":           "64-bit",
		"mfg":              "New Maker Ltd.",
		"model":            "Fake Model 3",
		"computerName":     "FAKE-PC-NEW1",
		"systemType":       2,
		"ram":              17179869184,
		"bootDrive":        "C:",
		"bootDriveFS":      "NTFS",
		"bootDriveCap":     510770802688,
		"bootDriveFree":    212345678848,
		"networkConfig": map[string]any{
			"AA:BB:CC:DD:EE:01": map[string]any{
				"dhcpEnabled": true,
				"ips":         []string{"10.0.0.20"},
			},
		},
	}
	for key, value := range overrides {
		if value == nil {
			delete(report, key)
			continue
		}
		report[key] = value
	}
	body, err := json.Marshal(report)
	if err != nil {
		panic(err)
	}
	return string(body)
}

func radiatorHeaders() map[string]string {
	return map[string]string{"Authorization": "Bearer " + radiatorToken}
}

func radiatorCount(t testing.TB, app core.App, table string) int {
	t.Helper()
	var row struct {
		N int `db:"n"`
	}
	if err := app.DB().NewQuery("SELECT COUNT(*) AS n FROM " + table).One(&row); err != nil {
		t.Fatal(err)
	}
	return row.N
}

// assertNothingStored checks that the request left the three collections as
// the fixtures have them.
func assertNothingStored(t testing.TB, app core.App) {
	t.Helper()
	if got := radiatorCount(t, app, "computers"); got != radiatorFixtureComputers {
		t.Errorf("computers = %d, want %d", got, radiatorFixtureComputers)
	}
	if got := radiatorCount(t, app, "computer_logins"); got != radiatorFixtureLogins {
		t.Errorf("computer_logins = %d, want %d", got, radiatorFixtureLogins)
	}
	if got := radiatorCount(t, app, "computer_raw_logins"); got != radiatorFixtureRawLogins {
		t.Errorf("computer_raw_logins = %d, want %d", got, radiatorFixtureRawLogins)
	}
}

func findComputerBySlug(t testing.TB, app core.App, slug string) *core.Record {
	t.Helper()
	record, err := app.FindFirstRecordByData("computers", "slug", slug)
	if err != nil {
		t.Fatalf("computer %s: %v", slug, err)
	}
	return record
}

func loginsOf(t testing.TB, app core.App, computerID string) []*core.Record {
	t.Helper()
	records, err := app.FindRecordsByFilter("computer_logins", "computer = {:c}", "-created", 0, 0, dbx.Params{"c": computerID})
	if err != nil {
		t.Fatal(err)
	}
	return records
}

// countingReader counts the bytes the server reads from a request body. It is
// not a *strings.Reader, so the request has no Content-Length.
type countingReader struct {
	r    io.Reader
	read int
}

func (c *countingReader) Read(p []byte) (int, error) {
	n, err := c.r.Read(p)
	c.read += n
	return n, err
}

func TestRadiatorRefusesBadTokens(t *testing.T) {
	for _, tc := range []struct {
		name   string
		header string
	}{
		{"missing Authorization header", ""},
		{"legacy TYBALT scheme", "TYBALT " + radiatorToken},
		{"wrong token", "Bearer not-the-secret"},
		{"empty bearer token", "Bearer "},
		{"expired radiator token", "Bearer " + radiatorExpiredToken},
		{"legacy_writeback token", "Bearer " + legacyWritebackToken},
	} {
		body := &countingReader{r: strings.NewReader(radiatorReport(nil))}
		headers := map[string]string{}
		if tc.header != "" {
			headers["Authorization"] = tc.header
		}
		scenario := tests.ApiScenario{
			Name:            tc.name + " returns 401 without reading the body",
			Method:          http.MethodPost,
			URL:             "/api/radiator",
			Body:            body,
			Headers:         headers,
			ExpectedStatus:  http.StatusUnauthorized,
			ExpectedContent: []string{`"status":401`},
			TestAppFactory:  testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				if body.read != 0 {
					t.Errorf("read %d bytes of the body before refusing the token", body.read)
				}
				assertNothingStored(t, app)
			},
		}
		scenario.Test(t)
	}
}

func TestRadiatorRefusesOtherMethods(t *testing.T) {
	for _, method := range []string{http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodDelete} {
		scenario := tests.ApiScenario{
			Name:            method + " returns 405 with Allow: POST",
			Method:          method,
			URL:             "/api/radiator",
			Headers:         radiatorHeaders(),
			ExpectedStatus:  http.StatusMethodNotAllowed,
			ExpectedContent: []string{`"status":405`},
			TestAppFactory:  testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				if got := res.Header.Get("Allow"); got != http.MethodPost {
					t.Errorf("Allow = %q, want POST", got)
				}
			},
		}
		scenario.Test(t)
	}
}

func TestRadiatorRefusesOtherContentTypes(t *testing.T) {
	for _, contentType := range []string{"", "text/plain", "multipart/form-data; boundary=x"} {
		headers := radiatorHeaders()
		headers["Content-Type"] = contentType
		scenario := tests.ApiScenario{
			Name:            "Content-Type " + contentType + " returns 415",
			Method:          http.MethodPost,
			URL:             "/api/radiator",
			Body:            strings.NewReader(radiatorReport(nil)),
			Headers:         headers,
			ExpectedStatus:  http.StatusUnsupportedMediaType,
			ExpectedContent: []string{`"status":415`},
			TestAppFactory:  testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				assertNothingStored(t, app)
			},
		}
		scenario.Test(t)
	}

	// A charset parameter is fine.
	headers := radiatorHeaders()
	headers["Content-Type"] = "application/json; charset=utf-8"
	scenario := tests.ApiScenario{
		Name:           "application/json with a charset returns 202",
		Method:         http.MethodPost,
		URL:            "/api/radiator",
		Body:           strings.NewReader(radiatorReport(nil)),
		Headers:        headers,
		ExpectedStatus: http.StatusAccepted,
		TestAppFactory: testutils.SetupTestApp,
	}
	scenario.Test(t)
}

func TestRadiatorRefusesOversizedBodies(t *testing.T) {
	// Valid JSON, padded past 64 KB with an unknown property.
	oversized := radiatorReport(map[string]any{"padding": strings.Repeat("x", 70<<10)})

	scenarios := []tests.ApiScenario{
		{
			Name:            "oversized body with Content-Length returns 413",
			Method:          http.MethodPost,
			URL:             "/api/radiator",
			Body:            strings.NewReader(oversized),
			Headers:         radiatorHeaders(),
			ExpectedStatus:  http.StatusRequestEntityTooLarge,
			ExpectedContent: []string{`"status":413`},
			TestAppFactory:  testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				assertNothingStored(t, app)
			},
		},
	}
	for _, scenario := range scenarios {
		scenario.Test(t)
	}

	// Without a Content-Length the server stops reading soon after the limit.
	huge := &countingReader{r: strings.NewReader(radiatorReport(map[string]any{"padding": strings.Repeat("x", 4<<20)}))}
	scenario := tests.ApiScenario{
		Name:            "oversized body without Content-Length returns 413 without reading it all",
		Method:          http.MethodPost,
		URL:             "/api/radiator",
		Body:            huge,
		Headers:         radiatorHeaders(),
		ExpectedStatus:  http.StatusRequestEntityTooLarge,
		ExpectedContent: []string{`"status":413`},
		TestAppFactory:  testutils.SetupTestApp,
		AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
			if huge.read > 1<<20 {
				t.Errorf("read %d bytes of a 4 MB body; want the read to stop near 64 KB", huge.read)
			}
			assertNothingStored(t, app)
		},
	}
	scenario.Test(t)
}

func TestRadiatorRefusesEmptyBodies(t *testing.T) {
	for _, tc := range []struct{ name, body string }{
		{"no body", ""},
		{"white space", "  \n"},
		{"empty object", "{}"},
		{"null", "null"},
		{"empty array", "[]"},
		{"invalid JSON", `{"serial": `},
		{"a JSON string", `"report"`},
	} {
		scenario := tests.ApiScenario{
			Name:            tc.name + " returns 400 and stores nothing",
			Method:          http.MethodPost,
			URL:             "/api/radiator",
			Body:            strings.NewReader(tc.body),
			Headers:         radiatorHeaders(),
			ExpectedStatus:  http.StatusBadRequest,
			ExpectedContent: []string{`"status":400`},
			TestAppFactory:  testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				assertNothingStored(t, app)
			},
		}
		scenario.Test(t)
	}
}

func TestRadiatorStoresInvalidReportsAsRawLogins(t *testing.T) {
	body := radiatorReport(map[string]any{
		"upn":           nil,
		"networkConfig": map[string]any{},
		"serial":        "abc",
	})
	scenario := tests.ApiScenario{
		Name:           "invalid report returns 202 and is stored with its errors",
		Method:         http.MethodPost,
		URL:            "/api/radiator",
		Body:           strings.NewReader(body),
		Headers:        radiatorHeaders(),
		ExpectedStatus: http.StatusAccepted,
		TestAppFactory: testutils.SetupTestApp,
		AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
			if got := radiatorCount(t, app, "computers"); got != radiatorFixtureComputers {
				t.Errorf("computers = %d, want %d", got, radiatorFixtureComputers)
			}
			if got := radiatorCount(t, app, "computer_logins"); got != radiatorFixtureLogins {
				t.Errorf("computer_logins = %d, want %d", got, radiatorFixtureLogins)
			}
			raw, err := app.FindFirstRecordByData("computer_raw_logins", "computer_name", "FAKE-PC-NEW1")
			if err != nil {
				t.Fatalf("raw login not stored: %v", err)
			}
			errorText := raw.GetString("error")
			for _, want := range []string{
				"upn: is required",
				"networkConfig: must have at least 1 property",
				"serial: must be at least 4 characters",
			} {
				if !strings.Contains(errorText, want) {
					t.Errorf("error %q does not contain %q", errorText, want)
				}
			}
			var stored map[string]any
			if err := json.Unmarshal([]byte(raw.GetString("report")), &stored); err != nil {
				t.Fatal(err)
			}
			if stored["serial"] != "abc" || stored["mfg"] != "New Maker Ltd." {
				t.Errorf("report not stored as received: %v", stored)
			}
		},
	}
	scenario.Test(t)

	// The computer key is cleaned after validation. A serial that is long
	// enough for the schema but too short once cleaned is a reject, not a
	// server error.
	scenario = tests.ApiScenario{
		Name:           "serial too short once cleaned is stored as a raw login",
		Method:         http.MethodPost,
		URL:            "/api/radiator",
		Body:           strings.NewReader(radiatorReport(map[string]any{"serial": "a/ /b"})),
		Headers:        radiatorHeaders(),
		ExpectedStatus: http.StatusAccepted,
		TestAppFactory: testutils.SetupTestApp,
		AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
			raw, err := app.FindFirstRecordByData("computer_raw_logins", "computer_name", "FAKE-PC-NEW1")
			if err != nil {
				t.Fatalf("raw login not stored: %v", err)
			}
			if !strings.Contains(raw.GetString("error"), "too short") {
				t.Errorf("error = %q", raw.GetString("error"))
			}
			if got := radiatorCount(t, app, "computers"); got != radiatorFixtureComputers {
				t.Errorf("computers = %d, want %d", got, radiatorFixtureComputers)
			}
		},
	}
	scenario.Test(t)
}

func TestRadiatorStoresValidReportForNewComputer(t *testing.T) {
	scenario := tests.ApiScenario{
		Name:           "valid report creates a computer and a login, with no user match",
		Method:         http.MethodPost,
		URL:            "/api/radiator",
		Body:           strings.NewReader(radiatorReport(nil)),
		Headers:        radiatorHeaders(),
		ExpectedStatus: http.StatusAccepted,
		TestAppFactory: testutils.SetupTestApp,
		AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
			// makeSlug: "New Maker Ltd." loses "ltd" and "." and is trimmed.
			computer := findComputerBySlug(t, app, "FAKENEW0001,new_maker")
			checks := map[string]any{
				"serial":             "FAKENEW0001",
				"mfg":                "New Maker Ltd.",
				"model":              "Fake Model 3",
				"computer_name":      "FAKE-PC-NEW1",
				"os_version":         "10.0.19045",
				"os_arch":            "64-bit",
				"boot_drive":         "C:",
				"boot_drive_fs":      "NTFS",
				"user_given_name":    "Fake",
				"user_surname":       "Person",
				"email":              "nobody.here@example.com",
				"upn":                "nobody.here@example.com",
				"user_source_anchor": unmatchedAnchor,
				"last_user":          "",
				"assigned_to":        "",
			}
			for field, want := range checks {
				if got := computer.GetString(field); got != want {
					t.Errorf("%s = %q, want %q", field, got, want)
				}
			}
			numbers := map[string]float64{
				"system_type":      2,
				"os_sku":           48,
				"ram":              17179869184,
				"boot_drive_cap":   510770802688,
				"boot_drive_free":  212345678848,
				"radiator_version": 9,
			}
			for field, want := range numbers {
				if got := computer.GetFloat(field); got != want {
					t.Errorf("%s = %v, want %v", field, got, want)
				}
			}
			if computer.GetDateTime("last_login").IsZero() {
				t.Error("last_login not set")
			}
			if !computer.GetDateTime("retired").IsZero() || !computer.GetDateTime("assigned_time").IsZero() {
				t.Error("a report must not set retired or assigned_time")
			}
			if !strings.Contains(computer.GetString("network_config"), "AA:BB:CC:DD:EE:01") {
				t.Errorf("network_config = %s", computer.GetString("network_config"))
			}

			logins := loginsOf(t, app, computer.Id)
			if len(logins) != 1 {
				t.Fatalf("logins = %d, want 1", len(logins))
			}
			login := logins[0]
			if login.GetString("given_name") != "Fake" || login.GetString("surname") != "Person" ||
				login.GetString("user_source_anchor") != unmatchedAnchor || login.GetString("user") != "" {
				t.Errorf("login = %v", login.PublicExport())
			}
			if got := radiatorCount(t, app, "computer_raw_logins"); got != radiatorFixtureRawLogins {
				t.Errorf("computer_raw_logins = %d, want %d", got, radiatorFixtureRawLogins)
			}
		},
	}
	scenario.Test(t)
}

func TestRadiatorUpdatesExistingComputers(t *testing.T) {
	scenario := tests.ApiScenario{
		Name:   "report for an assigned computer updates it, keeps the assignment and clears its raw logins",
		Method: http.MethodPost,
		URL:    "/api/radiator",
		Body: strings.NewReader(radiatorReport(map[string]any{
			"serial":       "FAKESERIAL0001",
			"mfg":          "Fake Computers Inc.",
			"model":        "Fake Laptop 1 Rev B",
			"computerName": "FAKE-PC-0001",
		})),
		Headers:        radiatorHeaders(),
		ExpectedStatus: http.StatusAccepted,
		TestAppFactory: testutils.SetupTestApp,
		AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
			if got := radiatorCount(t, app, "computers"); got != radiatorFixtureComputers {
				t.Errorf("computers = %d, want %d (no new computer)", got, radiatorFixtureComputers)
			}
			computer := findComputerBySlug(t, app, "FAKESERIAL0001,fake_computers")
			if computer.Id != assignedComputer {
				t.Fatalf("matched %s, want %s", computer.Id, assignedComputer)
			}
			if got := computer.GetString("model"); got != "Fake Laptop 1 Rev B" {
				t.Errorf("model = %q", got)
			}
			if got := computer.GetString("assigned_to"); got != fattUserID {
				t.Errorf("assigned_to = %q, want %q", got, fattUserID)
			}
			if got := computer.GetDateTime("assigned_time").String(); got != "2026-10-01 12:00:00.000Z" {
				t.Errorf("assigned_time = %q", got)
			}
			if !computer.GetDateTime("retired").IsZero() {
				t.Error("retired must stay empty")
			}
			// The user fields describe the last report; this one matched nobody.
			if computer.GetString("last_user") != "" || computer.GetString("user_surname") != "Person" {
				t.Errorf("last user fields = %q %q", computer.GetString("last_user"), computer.GetString("user_surname"))
			}
			if got := computer.GetDateTime("last_login").String(); got <= "2026-11-01 13:00:00.000Z" {
				t.Errorf("last_login = %q, want now", got)
			}
			if got := len(loginsOf(t, app, computer.Id)); got != 2 {
				t.Errorf("logins of the computer = %d, want 2", got)
			}

			// The reject for FAKE-PC-0001 is gone; the one for FAKE-PC-0099 stays.
			if _, err := app.FindRecordById("computer_raw_logins", "cmprawlogin0001"); err == nil {
				t.Error("the raw login of FAKE-PC-0001 was not deleted")
			}
			if _, err := app.FindRecordById("computer_raw_logins", "cmprawlogin0002"); err != nil {
				t.Error("the raw login of another computer was deleted")
			}
		},
	}
	scenario.Test(t)

	scenario = tests.ApiScenario{
		Name:   "report for a retired computer keeps it retired and keeps fields the report leaves out",
		Method: http.MethodPost,
		URL:    "/api/radiator",
		Body: strings.NewReader(radiatorReport(map[string]any{
			"serial":       "FAKESERIAL0002",
			"mfg":          "Fake Computers Inc.",
			"computerName": "FAKE-PC-0002",
			"model":        nil,
			"ram":          nil,
		})),
		Headers:        radiatorHeaders(),
		ExpectedStatus: http.StatusAccepted,
		TestAppFactory: testutils.SetupTestApp,
		AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
			computer := findComputerBySlug(t, app, "FAKESERIAL0002,fake_computers")
			if computer.Id != retiredComputer {
				t.Fatalf("matched %s, want %s", computer.Id, retiredComputer)
			}
			if got := computer.GetDateTime("retired").String(); got != "2026-06-01 00:00:00.000Z" {
				t.Errorf("retired = %q", got)
			}
			if got := computer.GetString("model"); got != "Fake Desktop 2" {
				t.Errorf("model = %q, want the stored value", got)
			}
			if got := computer.GetFloat("ram"); got != 8589934592 {
				t.Errorf("ram = %v, want the stored value", got)
			}
			if got := computer.GetFloat("boot_drive_cap"); got != 510770802688 {
				t.Errorf("boot_drive_cap = %v, want the reported value", got)
			}
		},
	}
	scenario.Test(t)
}

func TestRadiatorMatchesUsers(t *testing.T) {
	for _, tc := range []struct {
		name      string
		overrides map[string]any
		wantUser  string
	}{
		{
			// The anchor wins over an email of another user.
			name: "by source anchor, in any case",
			overrides: map[string]any{
				"userSourceAnchor": strings.ToUpper(anchorTimeUser),
				"email":            "fatt@mac.com",
			},
			wantUser: timeUserID,
		},
		{
			name:      "by email, ignoring case",
			overrides: map[string]any{"email": "FATT@Mac.com"},
			wantUser:  fattUserID,
		},
		{
			name:      "by UPN when there is no email",
			overrides: map[string]any{"email": nil, "upn": "Hal@2005.com"},
			wantUser:  halUserID,
		},
		{
			name:      "no match",
			overrides: nil,
			wantUser:  "",
		},
	} {
		scenario := tests.ApiScenario{
			Name:           "user matched " + tc.name,
			Method:         http.MethodPost,
			URL:            "/api/radiator",
			Body:           strings.NewReader(radiatorReport(tc.overrides)),
			Headers:        radiatorHeaders(),
			ExpectedStatus: http.StatusAccepted,
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				computer := findComputerBySlug(t, app, "FAKENEW0001,new_maker")
				if got := computer.GetString("last_user"); got != tc.wantUser {
					t.Errorf("last_user = %q, want %q", got, tc.wantUser)
				}
				logins := loginsOf(t, app, computer.Id)
				if len(logins) != 1 || logins[0].GetString("user") != tc.wantUser {
					t.Errorf("login user, want %q", tc.wantUser)
				}
				anchor := computer.GetString("user_source_anchor")
				if anchor != strings.ToLower(anchor) {
					t.Errorf("anchor not lowercased: %q", anchor)
				}
				var users struct {
					N int `db:"n"`
				}
				if err := app.DB().NewQuery("SELECT COUNT(*) AS n FROM users WHERE lower(email) = 'nobody.here@example.com'").One(&users); err != nil {
					t.Fatal(err)
				}
				if users.N != 0 {
					t.Error("the endpoint must not create users")
				}
			},
		}
		scenario.Test(t)
	}
}

func TestRadiatorCoercesTypesAndDropsBadOptionalValues(t *testing.T) {
	// Old PowerShell sends some numbers as strings; legacy coerced them.
	scenario := tests.ApiScenario{
		Name:   "numbers sent as strings and a numeric serial are coerced",
		Method: http.MethodPost,
		URL:    "/api/radiator",
		Body: strings.NewReader(radiatorReport(map[string]any{
			"systemType":      " 5.",
			"osSku":           "48",
			"ram":             "17179869184",
			"radiatorVersion": "9",
			"bootDriveFree":   true,
			"serial":          12345678,
			"unknownProperty": "dropped",
		})),
		Headers:        radiatorHeaders(),
		ExpectedStatus: http.StatusAccepted,
		TestAppFactory: testutils.SetupTestApp,
		AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
			computer := findComputerBySlug(t, app, "12345678,new_maker")
			if got := computer.GetString("serial"); got != "12345678" {
				t.Errorf("serial = %q", got)
			}
			for field, want := range map[string]float64{
				"system_type":      5,
				"os_sku":           48,
				"ram":              17179869184,
				"radiator_version": 9,
				"boot_drive_free":  1,
			} {
				if got := computer.GetFloat(field); got != want {
					t.Errorf("%s = %v, want %v", field, got, want)
				}
			}
			if got := radiatorCount(t, app, "computer_raw_logins"); got != radiatorFixtureRawLogins {
				t.Errorf("computer_raw_logins = %d, want %d", got, radiatorFixtureRawLogins)
			}
		},
	}
	scenario.Test(t)

	scenario = tests.ApiScenario{
		Name:   "a string that is not a whole number is a reject",
		Method: http.MethodPost,
		URL:    "/api/radiator",
		Body: strings.NewReader(radiatorReport(map[string]any{
			"osSku": "4.5",
		})),
		Headers:        radiatorHeaders(),
		ExpectedStatus: http.StatusAccepted,
		TestAppFactory: testutils.SetupTestApp,
		AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
			raw, err := app.FindFirstRecordByData("computer_raw_logins", "computer_name", "FAKE-PC-NEW1")
			if err != nil {
				t.Fatalf("raw login not stored: %v", err)
			}
			if !strings.Contains(raw.GetString("error"), "osSku: must be an integer") {
				t.Errorf("error = %q", raw.GetString("error"))
			}
		},
	}
	scenario.Test(t)

	for _, tc := range []struct {
		name  string
		email any
	}{
		{"empty email", ""},
		{"null email", json.RawMessage("null")},
		{"invalid email", "not an email"},
	} {
		body := radiatorReport(map[string]any{
			"email":        tc.email,
			"computerName": strings.Repeat("N", 49),
		})
		scenario := tests.ApiScenario{
			Name:           tc.name + " and a too long computer name are dropped, not rejected",
			Method:         http.MethodPost,
			URL:            "/api/radiator",
			Body:           strings.NewReader(body),
			Headers:        radiatorHeaders(),
			ExpectedStatus: http.StatusAccepted,
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				computer := findComputerBySlug(t, app, "FAKENEW0001,new_maker")
				if got := computer.GetString("email"); got != "" {
					t.Errorf("email = %q, want empty", got)
				}
				if got := computer.GetString("computer_name"); got != "" {
					t.Errorf("computer_name = %q, want empty", got)
				}
				if got := radiatorCount(t, app, "computer_raw_logins"); got != radiatorFixtureRawLogins {
					t.Errorf("computer_raw_logins = %d, want %d", got, radiatorFixtureRawLogins)
				}
			},
		}
		scenario.Test(t)
	}
}

// Only server code writes the three collections. A holder of the computers
// claim (it.identity@example.com) or of admin (author@soup.com) can read them;
// others see nothing.
func TestComputersCollectionRules(t *testing.T) {
	computersToken, err := testutils.GenerateRecordToken("users", "it.identity@example.com")
	if err != nil {
		t.Fatal(err)
	}
	adminToken, err := testutils.GenerateRecordToken("users", "author@soup.com")
	if err != nil {
		t.Fatal(err)
	}
	otherToken, err := testutils.GenerateRecordToken("users", "time@test.com")
	if err != nil {
		t.Fatal(err)
	}

	var scenarios []tests.ApiScenario
	for _, collection := range []string{"computers", "computer_logins", "computer_raw_logins"} {
		url := "/api/collections/" + collection + "/records"
		scenarios = append(scenarios,
			tests.ApiScenario{
				Name:            collection + ": computers claim can list",
				Method:          http.MethodGet,
				URL:             url,
				Headers:         map[string]string{"Authorization": computersToken},
				ExpectedStatus:  http.StatusOK,
				ExpectedContent: []string{`"totalItems":`},
				NotExpectedContent: []string{
					`"totalItems":0`,
				},
				TestAppFactory: testutils.SetupTestApp,
			},
			tests.ApiScenario{
				Name:               collection + ": admin can list",
				Method:             http.MethodGet,
				URL:                url,
				Headers:            map[string]string{"Authorization": adminToken},
				ExpectedStatus:     http.StatusOK,
				ExpectedContent:    []string{`"totalItems":`},
				NotExpectedContent: []string{`"totalItems":0`},
				TestAppFactory:     testutils.SetupTestApp,
			},
			tests.ApiScenario{
				Name:            collection + ": other user sees nothing",
				Method:          http.MethodGet,
				URL:             url,
				Headers:         map[string]string{"Authorization": otherToken},
				ExpectedStatus:  http.StatusOK,
				ExpectedContent: []string{`"totalItems":0`},
				TestAppFactory:  testutils.SetupTestApp,
			},
			tests.ApiScenario{
				Name:            collection + ": admin cannot create through the API",
				Method:          http.MethodPost,
				URL:             url,
				Body:            strings.NewReader(`{}`),
				Headers:         map[string]string{"Authorization": adminToken},
				ExpectedStatus:  http.StatusForbidden,
				ExpectedContent: []string{`"status":403`},
				TestAppFactory:  testutils.SetupTestApp,
			},
		)
	}
	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
