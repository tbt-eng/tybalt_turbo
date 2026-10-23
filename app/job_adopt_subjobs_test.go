package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/tests"
)

// =============================================================================
// Adopt Subjobs Endpoint Tests
// =============================================================================
//
// POST /api/jobs/{id}/adopt-subjobs links a project to the unlinked jobs
// numbered as its subjobs. The fixtures give project 26-7710 of Sepulchi & Co
// two jobs it can adopt, one with another client, one with a subjob of its
// own, one already linked to project 26-5555, and two that only look alike.
// Horace Silver (author@soup.com) is an admin.

const (
	adoptProject       = "adoptproject000"
	adoptClient        = "pqpd90fqd5ohjcs"
	adoptFirst         = "adoptsubjob0001"
	adoptSecond        = "adoptsubjob0002"
	adoptOtherClientID = "adoptsubjob0003"
	adoptLongerBase    = "adoptsubjob0004"
	adoptNested        = "adoptsubjob0005"
	adoptLinked        = "adoptsubjob0006"
	adoptWithChild     = "adoptsubjob0007"
	adoptChild         = "adoptsubjob0008"
	adoptOtherProject  = "adoptproject001"
	adoptStaleUpdated  = "2020-01-01 00:00:00.000Z"
)

func newAdoptSubjobsApp(t *testing.T) *tests.TestApp {
	t.Helper()
	return newWorkflowTestApp(t, false)
}

func postAdoptSubjobs(t *testing.T, app *tests.TestApp, job, email string, subjobs []string, status int) string {
	t.Helper()
	response := performJSONRequest(t, app, http.MethodPost, "/api/jobs/"+job+"/adopt-subjobs", workflowToken(t, app, email), map[string]any{"subjobs": subjobs})
	mustStatus(t, response, status)
	return response.Body.String()
}

func unlinkedSubjobsFromDetails(t *testing.T, app *tests.TestApp, job, email string) []map[string]string {
	t.Helper()
	response := performJSONRequest(t, app, http.MethodGet, "/api/jobs/"+job+"/details", workflowToken(t, app, email), nil)
	mustStatus(t, response, http.StatusOK)
	var body struct {
		UnlinkedSubjobs []map[string]string `json:"unlinked_subjobs"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	return body.UnlinkedSubjobs
}

func TestJobDetailsListsUnlinkedSubjobsForAdmins(t *testing.T) {
	app := newAdoptSubjobsApp(t)

	got := unlinkedSubjobsFromDetails(t, app, adoptProject, horaceEmail)
	want := []struct{ id, blocker string }{
		{adoptSecond, ""},
		{adoptFirst, ""},
		{adoptOtherClientID, "Its client is Mobilia SA. Change its client to Sepulchi & Co first."},
		{adoptWithChild, "It has subjobs of its own."},
	}
	if len(got) != len(want) {
		t.Fatalf("expected %d unlinked subjobs, got %v", len(want), got)
	}
	for i, w := range want {
		if got[i]["id"] != w.id || got[i]["blocker"] != w.blocker {
			t.Errorf("unlinked subjob %d: got %v, want id %s blocker %q", i, got[i], w.id, w.blocker)
		}
	}

	if got := unlinkedSubjobsFromDetails(t, app, adoptProject, "noclaims@example.com"); len(got) != 0 {
		t.Errorf("a user without admin was sent unlinked subjobs: %v", got)
	}
}

func TestAdoptSubjobsLinksAdoptableSubjobs(t *testing.T) {
	app := newAdoptSubjobsApp(t)

	body := postAdoptSubjobs(t, app, adoptProject, horaceEmail, []string{adoptFirst, adoptSecond}, http.StatusOK)
	if strings.TrimSpace(body) != `{"adopted":["26-7710-02","26-7710-1"]}` {
		t.Errorf("unexpected response %s", body)
	}

	for _, id := range []string{adoptFirst, adoptSecond} {
		job := workflowRecord(t, app, "jobs", id)
		if job.GetString("parent") != adoptProject {
			t.Errorf("job %s has parent %q", id, job.GetString("parent"))
		}
		// The legacy writeback sends jobs that are not imported and were updated.
		if job.GetBool("_imported") || job.GetString("updated") == adoptStaleUpdated {
			t.Errorf("job %s is not marked for writeback", id)
		}
	}
	for id, parent := range map[string]string{
		adoptOtherClientID: "",
		adoptLongerBase:    "",
		adoptNested:        "",
		adoptLinked:        adoptOtherProject,
		adoptWithChild:     "",
	} {
		if got := workflowRecord(t, app, "jobs", id).GetString("parent"); got != parent {
			t.Errorf("job %s parent changed to %q", id, got)
		}
	}

	note, err := app.FindFirstRecordByData("client_notes", "note", "Adopted unlinked subjobs: 26-7710-02, 26-7710-1.")
	if err != nil {
		t.Fatal("no note records the adoption")
	}
	if note.GetString("job") != adoptProject || note.GetString("client") != adoptClient {
		t.Error("the adoption note is not on the project")
	}

	got := unlinkedSubjobsFromDetails(t, app, adoptProject, horaceEmail)
	if len(got) != 2 {
		t.Errorf("expected only the blocked subjobs to remain unlinked, got %v", got)
	}
}

func TestAdoptSubjobsRequiresTheCurrentList(t *testing.T) {
	app := newAdoptSubjobsApp(t)

	for name, subjobs := range map[string][]string{
		"none confirmed":      nil,
		"one missing":         {adoptFirst},
		"a blocked subjob":    {adoptFirst, adoptSecond, adoptOtherClientID},
		"another project's":   {adoptFirst, adoptSecond, adoptLinked},
		"an unrelated number": {adoptFirst, adoptSecond, adoptLongerBase},
	} {
		t.Run(name, func(t *testing.T) {
			postAdoptSubjobs(t, app, adoptProject, horaceEmail, subjobs, http.StatusConflict)
		})
	}
	for _, id := range []string{adoptFirst, adoptSecond} {
		if got := workflowRecord(t, app, "jobs", id).GetString("parent"); got != "" {
			t.Errorf("job %s was adopted by a refused request", id)
		}
	}
}

func TestAdoptSubjobsRefusals(t *testing.T) {
	app := newAdoptSubjobsApp(t)

	postAdoptSubjobs(t, app, adoptProject, "noclaims@example.com", []string{adoptFirst, adoptSecond}, http.StatusForbidden)
	// Only a top-level project can adopt subjobs.
	postAdoptSubjobs(t, app, adoptWithChild, horaceEmail, []string{adoptChild}, http.StatusBadRequest)
	// 26-5555 has no unlinked subjobs of its own.
	postAdoptSubjobs(t, app, adoptOtherProject, horaceEmail, nil, http.StatusBadRequest)

	if got := workflowRecord(t, app, "jobs", adoptFirst).GetString("parent"); got != "" {
		t.Errorf("a refused request adopted %s", adoptFirst)
	}
}
