package main

import (
	"net/http"
	"slices"
	"testing"
	"time"

	"tybalt/routes"

	"github.com/pocketbase/pocketbase/tests"
)

// The duplicate-review fixtures are the rows with ids that start with "dup" in
// clients.csv, jobs.csv, job_requests.csv and time_entries.csv (time entries
// "teduplicate*"). Request dupreq000000001 is a subjob request of 92-9001 for
// Duplicate Review Test Client, from proposal P26-9901, with client PO
// DUP-PO-77 and client reference DUP-REF-9.
const duplicateReviewRequest = "/api/job_requests/dupreq000000001/duplicates"

// duplicateReviewApp fixes the review's reference date, so "the last year"
// does not depend on how far today is from the fixture dates.
func duplicateReviewApp(t *testing.T, now time.Time) *tests.TestApp {
	t.Helper()
	app := newWorkflowTestApp(t, true)
	routes.SetJobRequestDuplicatesNowForTest(app, func() time.Time { return now })
	return app
}

type duplicateReviewRow struct {
	id           string
	tier         string
	reasons      []string
	lastActivity string
	raw          map[string]any
}

func duplicateReviewRows(t *testing.T, result map[string]any, key string) []duplicateReviewRow {
	t.Helper()
	items, ok := result[key].([]any)
	if !ok {
		t.Fatalf("response has no %s list: %#v", key, result)
	}
	rows := []duplicateReviewRow{}
	for _, item := range items {
		raw := item.(map[string]any)
		row := duplicateReviewRow{id: raw["id"].(string), tier: raw["tier"].(string), raw: raw}
		for _, reason := range raw["reasons"].([]any) {
			row.reasons = append(row.reasons, reason.(string))
		}
		if value, ok := raw["last_activity"].(string); ok {
			row.lastActivity = value
		}
		rows = append(rows, row)
	}
	return rows
}

func duplicateReviewIDs(rows []duplicateReviewRow) []string {
	ids := []string{}
	for _, row := range rows {
		ids = append(ids, row.id)
	}
	return ids
}

func TestJobRequestDuplicateReviewListsExactMatchesAndRecentOpenProjects(t *testing.T) {
	app := duplicateReviewApp(t, time.Date(2026, time.November, 20, 12, 0, 0, 0, time.UTC))
	result := clientSetupJSON(t, app, http.MethodGet, duplicateReviewRequest, workflowAuthor, nil, http.StatusOK)
	rows := duplicateReviewRows(t, result, "jobs")

	// Exact matches first, whatever their age or status, then the client's
	// recent open project with a similar description, then the client's other
	// recent open projects. Newest activity first within a tier.
	want := []struct{ id, tier string }{
		{"duppropjob00001", "exact"},   // same proposal, Cancelled, award 2020
		{"dupsibling00001", "exact"},   // same parent, Closed, award 2019-05
		{"dupparent000001", "exact"},   // the parent project, award 2019-03
		{"dupclientpo0001", "exact"},   // another client, same client PO, Closed, 2018
		{"dupclientref001", "exact"},   // another client, same client reference, Closed, 2017
		{"dupsimilar00001", "similar"}, // award 2026-03-02
		{"duprecentaward1", "client"},  // In Progress, award 2026-10-01
		{"duprecenttime01", "client"},  // Active, award 2023, time entry 2026-08-04
	}
	if len(rows) != len(want) {
		t.Fatalf("listed %v, want %d rows", duplicateReviewIDs(rows), len(want))
	}
	byID := map[string]duplicateReviewRow{}
	for i, row := range rows {
		if row.id != want[i].id || row.tier != want[i].tier {
			t.Fatalf("row %d is %s (%s), want %s (%s); list %v", i, row.id, row.tier, want[i].id, want[i].tier, duplicateReviewIDs(rows))
		}
		byID[row.id] = row
	}

	// Closed, Cancelled, not recent, or with no award and no time: left out
	// and counted. A similar description does not list an old project.
	for _, id := range []string{"dupoldactive001", "dupclosedrec001", "dupcancelrec001", "dupnoaward00001"} {
		if _, ok := byID[id]; ok {
			t.Fatalf("%s is listed but should be counted as hidden", id)
		}
	}
	if result["hidden_client_jobs"] != float64(4) {
		t.Fatalf("hidden_client_jobs = %v, want 4", result["hidden_client_jobs"])
	}

	// Location never lists another client's project, even with the same
	// description, and the proposal itself is never listed.
	for _, id := range []string{"duplocation0001", "dupprop00000001"} {
		if _, ok := byID[id]; ok {
			t.Fatalf("%s is listed", id)
		}
	}
	// Location is a supporting reason on a listed row.
	similar := byID["dupsimilar00001"]
	if !slices.Equal(similar.reasons, []string{"Same client", "Same location", "Similar description"}) {
		t.Fatalf("dupsimilar00001 reasons = %v", similar.reasons)
	}
	if got := byID["dupclientpo0001"].reasons; !slices.Equal(got, []string{"Same client PO"}) {
		t.Fatalf("dupclientpo0001 reasons = %v", got)
	}
	if got := byID["dupparent000001"].reasons; !slices.Equal(got, []string{"Same client", "Parent project"}) {
		t.Fatalf("dupparent000001 reasons = %v", got)
	}

	// Last activity is the later of the award date and the latest time entry.
	if got := byID["duprecenttime01"].lastActivity; got != "2026-08-04" {
		t.Fatalf("duprecenttime01 last_activity = %q", got)
	}
	if got := byID["duprecentaward1"].lastActivity; got != "2026-10-01" {
		t.Fatalf("duprecentaward1 last_activity = %q", got)
	}

	// Names come with the row.
	sibling := byID["dupsibling00001"].raw
	if sibling["parent_number"] != "92-9001" || sibling["client_name"] != "Duplicate Review Test Client" || sibling["manager_name"] != "Fakesy Manjor" || sibling["status"] != "Closed" {
		t.Fatalf("dupsibling00001 row = %#v", sibling)
	}
}

func TestJobRequestDuplicateReviewWindowFollowsTheReferenceDate(t *testing.T) {
	// A year later, the client's projects are no longer recent. The exact
	// matches stay.
	app := duplicateReviewApp(t, time.Date(2027, time.November, 20, 12, 0, 0, 0, time.UTC))
	result := clientSetupJSON(t, app, http.MethodGet, duplicateReviewRequest, workflowAuthor, nil, http.StatusOK)
	rows := duplicateReviewRows(t, result, "jobs")
	want := []string{"duppropjob00001", "dupsibling00001", "dupparent000001", "dupclientpo0001", "dupclientref001"}
	if !slices.Equal(duplicateReviewIDs(rows), want) {
		t.Fatalf("listed %v, want %v", duplicateReviewIDs(rows), want)
	}
	if result["hidden_client_jobs"] != float64(7) {
		t.Fatalf("hidden_client_jobs = %v, want 7", result["hidden_client_jobs"])
	}
}

func TestJobRequestDuplicateReviewListsPendingRequests(t *testing.T) {
	app := duplicateReviewApp(t, time.Date(2026, time.November, 20, 12, 0, 0, 0, time.UTC))
	result := clientSetupJSON(t, app, http.MethodGet, duplicateReviewRequest, workflowAuthor, nil, http.StatusOK)
	rows := duplicateReviewRows(t, result, "requests")
	// dupreq000000003: another client with the same client PO (exact).
	// dupreq000000002: the same client with a similar description.
	// dupreq000000004 (another client, same location and description) and
	// the request itself are not listed.
	if got := duplicateReviewIDs(rows); !slices.Equal(got, []string{"dupreq000000003", "dupreq000000002"}) {
		t.Fatalf("pending requests %v", got)
	}
	if rows[0].tier != "exact" || !slices.Equal(rows[0].reasons, []string{"Same client PO"}) {
		t.Fatalf("dupreq000000003 = %s %v", rows[0].tier, rows[0].reasons)
	}
	if rows[1].tier != "similar" || !slices.Equal(rows[1].reasons, []string{"Same client", "Similar description"}) {
		t.Fatalf("dupreq000000002 = %s %v", rows[1].tier, rows[1].reasons)
	}
}
