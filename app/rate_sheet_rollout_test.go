package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"tybalt/routes"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

// Revision rollout (issue 205): an admin moves every active project on an
// earlier revision of a rate sheet to its latest revision, from one date. The
// fixtures are rows of the seed CSV files (app/test_seed_data/data). The
// Rollout Rates revisions are the 2025 Standard Rates sheet (c41ofep525bcacj)
// with 8000 (rev 0) or 9000 (rev 1) added to each rate, so that they match no
// other seeded sheet.
//
// Sheets:
//   - Rollout Rates rev 0 (rsrollout0r0001): 2024-01-01, inactive
//   - Rollout Rates rev 1 (rsrollout1r0001): 2025-03-01, active, the latest
//
// Projects, with their current sheet:
//   - jobrollout00001: Active, rev 0, opening period on rev 0. Time on
//     2025-02-03 (2 h) and 2025-04-07 (3 h) in role 0l1llckilwz1mqz, which is
//     8160 on rev 0 and 9160 on rev 1.
//   - jobrollout00002: Active, rev 0, no history rows
//   - jobrollout00003: Closed, rev 0
//   - jobrollout00004: Cancelled, rev 0
//   - jobrollout00005: Active, rev 0: opening period on the standard sheet,
//     from 2025-06-02 on rev 0
//   - jobrollout00006: Active, already on rev 1
//
// admin.only@example.com (u_admin_only) holds only admin.
// history.revise@example.com holds job and rate_sheet_revise, not admin.
// Only TestRateSheetRolloutPendingAbsorb changes data other than through the
// rollout routes, and it says so.

const (
	rsRollout0     = "rsrollout0r0001"
	rsRollout1     = "rsrollout1r0001"
	rolloutAdmin   = "admin.only@example.com"
	rolloutAdminID = "u_admin_only"
	rolloutURL     = "/api/rate_sheets/" + rsRollout1 + "/rollout"
	rolloutJobWIP  = "jobrollout00001"
	rolloutNoRows  = "jobrollout00002"
)

// The projects that the rollout must not change: closed, cancelled, a later
// period, already on rev 1, and a project on another sheet.
var rolloutUnmoved = []string{"jobrollout00003", "jobrollout00004", "jobrollout00005", "jobrollout00006", "test_job_w_rs"}

type rolloutResult struct {
	Message       string `json:"message"`
	EffectiveDate string `json:"effective_date"`
	Move          int    `json:"move"`
	Closed        int    `json:"closed"`
	Cancelled     int    `json:"cancelled"`
	LaterPeriod   int    `json:"later_period"`
}

func rolloutRequest(t *testing.T, app *tests.TestApp, method, path, email string, data any) (int, rolloutResult, string) {
	t.Helper()
	status, body := absorbRequest(t, app, method, path, email, data)
	var result rolloutResult
	_ = json.Unmarshal([]byte(body), &result)
	return status, result, body
}

func rolloutBody(date, reason string, expected int) map[string]any {
	return map[string]any{"effective_date": date, "reason": reason, "expected_count": expected}
}

// rolloutSnapshot returns the job fields that the writeback reads and every
// field of the history of each project.
func rolloutSnapshot(t testing.TB, app core.App, jobs ...string) string {
	t.Helper()
	parts := []string{}
	for _, job := range jobs {
		parts = append(parts, jobWritebackState(t, app, job))
		parts = append(parts, periodSnapshot(t, app, job)...)
	}
	return strings.Join(parts, "\n")
}

func allRolloutJobs() []string {
	return append([]string{rolloutJobWIP, rolloutNoRows}, rolloutUnmoved...)
}

func TestRateSheetRolloutPreview(t *testing.T) {
	app := newRateSheetAbsorbApp(t)
	defer app.Cleanup()

	counts := []struct {
		name, query string
		want        rolloutResult
	}{
		{"the date defaults to the revision's effective date", "",
			rolloutResult{EffectiveDate: "2025-03-01", Move: 2, Closed: 1, Cancelled: 1, LaterPeriod: 1}},
		{"a project whose latest period starts on the date is skipped", "?effective_date=2025-06-02",
			rolloutResult{EffectiveDate: "2025-06-02", Move: 2, Closed: 1, Cancelled: 1, LaterPeriod: 1}},
		{"a later date moves that project too", "?effective_date=2025-07-01",
			rolloutResult{EffectiveDate: "2025-07-01", Move: 3, Closed: 1, Cancelled: 1}},
	}
	for _, tc := range counts {
		status, got, body := rolloutRequest(t, app, http.MethodGet, rolloutURL+tc.query, rolloutAdmin, nil)
		if status != http.StatusOK || got != tc.want {
			t.Errorf("%s: status %d, %+v; want %+v (%s)", tc.name, status, got, tc.want, body)
		}
	}

	refusals := []struct {
		name, path, email string
		status            int
		contains          string
	}{
		{"a user without admin is refused", rolloutURL, histRevise, http.StatusForbidden, "admin"},
		{"authentication is required", rolloutURL, "", http.StatusUnauthorized, "{"},
		{"an earlier revision is refused", "/api/rate_sheets/" + rsRollout0 + "/rollout", rolloutAdmin, http.StatusBadRequest, "latest revision"},
		{"an inactive latest revision is refused", "/api/rate_sheets/" + rsAbLateInactive + "/rollout", rolloutAdmin, http.StatusBadRequest, "Activate this revision"},
		{"an unknown sheet is not found", "/api/rate_sheets/missing/rollout", rolloutAdmin, http.StatusNotFound, "not found"},
		{"a future date is refused", rolloutURL + "?effective_date=2999-01-01", rolloutAdmin, http.StatusBadRequest, `"future"`},
		{"a date before the revision is refused", rolloutURL + "?effective_date=2025-02-28", rolloutAdmin, http.StatusBadRequest, "before_revision"},
		{"an invalid date is refused", rolloutURL + "?effective_date=2025-13-01", rolloutAdmin, http.StatusBadRequest, `"invalid"`},
	}
	for _, tc := range refusals {
		status, _, body := rolloutRequest(t, app, http.MethodGet, tc.path, tc.email, nil)
		if status != tc.status || !strings.Contains(body, tc.contains) {
			t.Errorf("%s: status %d: %s; want %d with %q", tc.name, status, body, tc.status, tc.contains)
		}
	}
}

// An admin without the job and rate_sheet_revise claims moves the two active
// projects. Each gets a dated period on rev 1, its job follows, and its time
// from the date is priced by rev 1. Nothing else changes. A second run moves
// nothing.
func TestRateSheetRolloutMovesActiveProjects(t *testing.T) {
	app := newRateSheetAbsorbApp(t)
	defer app.Cleanup()

	before := histWIP(t, app, rolloutJobWIP)
	if before.Hours != 5 || before.TimeValue != 5*8160 {
		t.Fatalf("WIP before = %+v; want 5 h worth %d", before, 5*8160)
	}
	moved := map[string]string{}
	for _, job := range []string{rolloutJobWIP, rolloutNoRows} {
		moved[job] = jobWritebackState(t, app, job)
	}
	unmoved := rolloutSnapshot(t, app, rolloutUnmoved...)

	const reason = "Clients accepted Rollout Rates rev 1"
	status, got, body := rolloutRequest(t, app, http.MethodPost, rolloutURL, rolloutAdmin, rolloutBody("2025-03-01", reason, 2))
	want := rolloutResult{EffectiveDate: "2025-03-01", Move: 2, Closed: 1, Cancelled: 1, LaterPeriod: 1}
	if status != http.StatusOK || got != want {
		t.Fatalf("rollout: status %d, %+v; want %+v (%s)", status, got, want, body)
	}

	for job, state := range moved {
		if got := fieldOf(t, app, "jobs", job, "rate_sheet"); got != rsRollout1 {
			t.Errorf("%s is on %s; want rev 1", job, got)
		}
		if jobWritebackState(t, app, job) == state {
			t.Errorf("%s was not saved for the writeback", job)
		}
		rows, err := app.FindRecordsByFilter("job_rate_sheets", "job = {:job}", "effective_date", 0, 0, dbx.Params{"job": job})
		if err != nil {
			t.Fatal(err)
		}
		if len(rows) != 2 {
			t.Fatalf("%s has %d periods; want 2", job, len(rows))
		}
		opening, added := rows[0], rows[1]
		if opening.GetString("effective_date") != "" || opening.GetString("rate_sheet") != rsRollout0 {
			t.Errorf("%s opening period: %v", job, recordAllFields(opening))
		}
		if added.GetString("effective_date") != "2025-03-01" || added.GetString("rate_sheet") != rsRollout1 ||
			added.GetString("reason") != reason || added.GetString("created_by") != rolloutAdminID ||
			added.GetString("updated_by") != rolloutAdminID {
			t.Errorf("%s new period: %v", job, recordAllFields(added))
		}
	}
	// The project with no history got its opening row from its sheet.
	if got := periodSnapshot(t, app, rolloutNoRows)[0]; !strings.Contains(got, `"created_by":""`) {
		t.Errorf("the opening row of %s: %s", rolloutNoRows, got)
	}
	if got := rolloutSnapshot(t, app, rolloutUnmoved...); got != unmoved {
		t.Fatalf("projects that do not move changed:\n%s\nwant\n%s", got, unmoved)
	}

	// 2 h before the date keep the rev 0 rate; 3 h after it use rev 1.
	after := histWIP(t, app, rolloutJobWIP)
	if after.Hours != 5 || after.TimeValue != 2*8160+3*9160 {
		t.Fatalf("WIP after = %+v; want 5 h worth %d", after, 2*8160+3*9160)
	}

	// A second run finds nothing to move.
	done := rolloutSnapshot(t, app, allRolloutJobs()...)
	status, got, body = rolloutRequest(t, app, http.MethodGet, rolloutURL, rolloutAdmin, nil)
	want = rolloutResult{EffectiveDate: "2025-03-01", Closed: 1, Cancelled: 1, LaterPeriod: 1}
	if status != http.StatusOK || got != want {
		t.Fatalf("preview after the rollout: status %d, %+v; want %+v (%s)", status, got, want, body)
	}
	status, got, body = rolloutRequest(t, app, http.MethodPost, rolloutURL, rolloutAdmin, rolloutBody("2025-03-01", reason, 2))
	if status != http.StatusConflict || got.Move != 0 {
		t.Fatalf("a repeat with the old count: status %d: %s; want 409 with move 0", status, body)
	}
	status, got, body = rolloutRequest(t, app, http.MethodPost, rolloutURL, rolloutAdmin, rolloutBody("2025-03-01", reason, 0))
	if status != http.StatusOK || got.Move != 0 {
		t.Fatalf("a second run: status %d: %s; want 200 with move 0", status, body)
	}
	if got := rolloutSnapshot(t, app, allRolloutJobs()...); got != done {
		t.Fatalf("a second run changed projects:\n%s\nwant\n%s", got, done)
	}
}

// Every refusal writes nothing.
func TestRateSheetRolloutRefusals(t *testing.T) {
	app := newRateSheetAbsorbApp(t)
	defer app.Cleanup()

	before := rolloutSnapshot(t, app, allRolloutJobs()...)
	const reason = "Clients accepted Rollout Rates rev 1"
	cases := []struct {
		name, path, email string
		body              map[string]any
		status            int
		contains          string
	}{
		{"a user with job and rate_sheet_revise but not admin", rolloutURL, histRevise,
			rolloutBody("2025-03-01", reason, 2), http.StatusForbidden, "admin"},
		{"no authentication", rolloutURL, "",
			rolloutBody("2025-03-01", reason, 2), http.StatusUnauthorized, "{"},
		{"an earlier revision", "/api/rate_sheets/" + rsRollout0 + "/rollout", rolloutAdmin,
			rolloutBody("2025-03-01", reason, 2), http.StatusBadRequest, "latest revision"},
		{"an inactive latest revision", "/api/rate_sheets/" + rsAbLateInactive + "/rollout", rolloutAdmin,
			rolloutBody("2025-03-01", reason, 0), http.StatusBadRequest, "Activate this revision"},
		{"a future date", rolloutURL, rolloutAdmin,
			rolloutBody("2999-01-01", reason, 2), http.StatusBadRequest, `"future"`},
		{"a date before the revision", rolloutURL, rolloutAdmin,
			rolloutBody("2025-02-28", reason, 2), http.StatusBadRequest, "before_revision"},
		{"no date", rolloutURL, rolloutAdmin,
			rolloutBody("", reason, 2), http.StatusBadRequest, `"invalid"`},
		{"no reason", rolloutURL, rolloutAdmin,
			rolloutBody("2025-03-01", "  ", 2), http.StatusBadRequest, `"reason"`},
		{"a reason over 500 characters", rolloutURL, rolloutAdmin,
			rolloutBody("2025-03-01", strings.Repeat("a", 501), 2), http.StatusBadRequest, "too_long"},
		{"no expected count", rolloutURL, rolloutAdmin,
			map[string]any{"effective_date": "2025-03-01", "reason": reason}, http.StatusBadRequest, "expected_count"},
		{"a count that differs from the projects to move", rolloutURL, rolloutAdmin,
			rolloutBody("2025-03-01", reason, 3), http.StatusConflict, `"move":2,"closed":1,"cancelled":1,"later_period":1`},
	}
	for _, tc := range cases {
		status, _, body := rolloutRequest(t, app, http.MethodPost, tc.path, tc.email, tc.body)
		if status != tc.status || !strings.Contains(body, tc.contains) {
			t.Errorf("%s: status %d: %s; want %d with %q", tc.name, status, body, tc.status, tc.contains)
		}
		if got := rolloutSnapshot(t, app, allRolloutJobs()...); got != before {
			t.Fatalf("%s changed projects:\n%s\nwant\n%s", tc.name, got, before)
		}
	}
}

// No project moves while a rate sheet absorb is pending, even one that the
// absorb did not touch.
func TestRateSheetRolloutPendingAbsorb(t *testing.T) {
	app := newRateSheetAbsorbApp(t)
	defer app.Cleanup()

	// Test-time change to fixture data: absorb Absorb History Source into
	// Absorb History Target, which leaves a pending rate sheet absorb. It
	// moves no Rollout Rates project.
	if err := routes.AbsorbRecords(app, "rate_sheets", rsHistTarget, []string{rsHistSource}); err != nil {
		t.Fatal(err)
	}
	before := rolloutSnapshot(t, app, allRolloutJobs()...)
	status, _, body := rolloutRequest(t, app, http.MethodPost, rolloutURL, rolloutAdmin, rolloutBody("2025-03-01", "Clients accepted rev 1", 2))
	if status != http.StatusConflict || !strings.Contains(body, "absorb is waiting") {
		t.Fatalf("rollout while an absorb is pending: status %d: %s", status, body)
	}
	if got := rolloutSnapshot(t, app, allRolloutJobs()...); got != before {
		t.Fatalf("a refused rollout changed projects:\n%s\nwant\n%s", got, before)
	}
}
