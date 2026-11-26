package main

import (
	"encoding/json"
	"net/http"
	"testing"
	"tybalt/internal/testutils"

	"github.com/pocketbase/pocketbase/tests"
)

// Fixtures for the absorb history (absorb_history.csv):
//   - absorbhist00001: vendors, committed 2026-08-20, two records absorbed
//     and three references moved, absorbed and committed by tqqf7q0f3378rvp.
//   - absorbhist00002: clients, committed 2026-11-20, no committer.
//   - absorbhist00003: rate_sheets, committed 2026-11-22, no absorber.
//
// author@soup.com holds admin and not absorb (see rate_sheet_absorb_test.go).
const (
	histVendors    = "absorbhist00001"
	histClients    = "absorbhist00002"
	histRateSheets = "absorbhist00003"
)

type historyTestRow struct {
	ID              string          `json:"id"`
	CollectionName  string          `json:"collection_name"`
	TargetID        string          `json:"target_id"`
	TargetName      string          `json:"target_name"`
	Absorbed        []string        `json:"absorbed"`
	ReferencesMoved int             `json:"references_moved"`
	AbsorbedBy      string          `json:"absorbed_by"`
	AbsorbedByName  string          `json:"absorbed_by_name"`
	CommittedBy     string          `json:"committed_by"`
	CommittedByName string          `json:"committed_by_name"`
	CommittedAt     string          `json:"committed_at"`
	AbsorbedRecords json.RawMessage `json:"absorbed_records"`
}

func getHistory(t testing.TB, app *tests.TestApp, email string) []historyTestRow {
	t.Helper()
	res := performTestAPIRequest(t, app, http.MethodGet, "/api/absorb_history", nil, uncommitToken(t, email))
	mustStatus(t, res, http.StatusOK)
	var out []historyTestRow
	if err := json.Unmarshal(res.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	return out
}

func TestAbsorbHistoryList(t *testing.T) {
	app := testutils.SetupTestApp(t)

	rows := getHistory(t, app, proposalHolderEmail)
	if len(rows) != 3 || rows[0].ID != histRateSheets || rows[1].ID != histClients || rows[2].ID != histVendors {
		t.Fatalf("holder rows = %+v, want the three fixture rows, newest first", rows)
	}
	vendors := rows[2]
	if vendors.TargetName != "Kept Vendor (KV)" || vendors.ReferencesMoved != 3 {
		t.Errorf("vendors row = %+v", vendors)
	}
	if len(vendors.Absorbed) != 2 || vendors.Absorbed[0] != "Old Vendor (OV)" || vendors.Absorbed[1] != "Older Vendor" {
		t.Errorf("absorbed names = %v", vendors.Absorbed)
	}
	if vendors.AbsorbedBy != "tqqf7q0f3378rvp" || vendors.AbsorbedByName == "" || vendors.CommittedByName == "" {
		t.Errorf("vendors row names nobody: %+v", vendors)
	}
	if rows[1].CommittedBy != "" || rows[1].CommittedByName != "" || rows[1].ReferencesMoved != 0 {
		t.Errorf("clients row = %+v", rows[1])
	}

	// An admin without the absorb claim sees the rate sheet rows only.
	adminRows := getHistory(t, app, absorbAdminEmail)
	if len(adminRows) != 1 || adminRows[0].ID != histRateSheets {
		t.Fatalf("admin rows = %+v, want the rate sheet row only", adminRows)
	}

	cases := []struct {
		name   string
		url    string
		email  string
		status int
		code   string
	}{
		{"list without a claim", "/api/absorb_history", proposalRequesterEmail, 403, "unauthorized"},
		{"one without a claim", "/api/absorb_history/" + histVendors, proposalRequesterEmail, 403, "unauthorized"},
		{"a vendor row for an admin", "/api/absorb_history/" + histVendors, absorbAdminEmail, 404, "history_not_found"},
		{"a missing row", "/api/absorb_history/doesnotexist000", proposalHolderEmail, 404, "history_not_found"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res := performTestAPIRequest(t, app, http.MethodGet, c.url, nil, uncommitToken(t, c.email))
			if res.Code != c.status || proposalErrorCode(t, res.Body.Bytes()) != c.code {
				t.Fatalf("got %d %s, want %d %s", res.Code, res.Body.String(), c.status, c.code)
			}
		})
	}
	if res := performTestAPIRequest(t, app, http.MethodGet, "/api/absorb_history", nil, nil); res.Code != http.StatusUnauthorized {
		t.Fatalf("signed out: got %d", res.Code)
	}

	res := performTestAPIRequest(t, app, http.MethodGet, "/api/absorb_history/"+histVendors, nil, uncommitToken(t, proposalHolderEmail))
	mustStatus(t, res, http.StatusOK)
	var detail historyTestRow
	if err := json.Unmarshal(res.Body.Bytes(), &detail); err != nil {
		t.Fatal(err)
	}
	saved := []map[string]any{}
	if err := json.Unmarshal(detail.AbsorbedRecords, &saved); err != nil || len(saved) != 2 || saved[0]["id"] != "histvendgone001" {
		t.Errorf("detail absorbed_records = %s (%v)", detail.AbsorbedRecords, err)
	}
	if res := performTestAPIRequest(t, app, http.MethodGet, "/api/absorb_history/"+histRateSheets, nil, uncommitToken(t, absorbAdminEmail)); res.Code != http.StatusOK {
		t.Errorf("admin reads a rate sheet row: got %d", res.Code)
	}
}

// A commit writes the history row with who absorbed and who committed; an
// undo writes none. Both tests absorb sugclient000001 into BN Rail
// Corporation through the routes, so they change the test database: the
// absorbed client is deleted, and the commit adds a row to absorb_history.
func TestAbsorbCommitWritesHistory(t *testing.T) {
	app := testutils.SetupTestApp(t)
	holder, err := app.FindAuthRecordByEmail("users", proposalHolderEmail)
	if err != nil {
		t.Fatal(err)
	}
	token := uncommitToken(t, proposalHolderEmail)

	res := performTestAPIRequest(t, app, http.MethodPost, "/api/clients/"+clientBNRail+"/absorb",
		proposalJSON(t, map[string]any{"ids_to_absorb": []string{sugBNRailCorp}}), token)
	mustStatus(t, res, http.StatusOK)

	action, err := app.FindFirstRecordByData("absorb_actions", "collection_name", "clients")
	if err != nil {
		t.Fatal(err)
	}
	if action.GetString("absorbed_by") != holder.Id {
		t.Fatalf("action absorbed_by = %q, want %q", action.GetString("absorbed_by"), holder.Id)
	}

	// The pending action names the person on the suggestions page.
	pending := performTestAPIRequest(t, app, http.MethodGet, "/api/absorb_suggestions/clients", nil, token)
	mustStatus(t, pending, http.StatusOK)
	var suggestions struct {
		PendingAction *struct {
			AbsorbedBy     string `json:"absorbed_by"`
			AbsorbedByName string `json:"absorbed_by_name"`
		} `json:"pending_action"`
	}
	if err := json.Unmarshal(pending.Body.Bytes(), &suggestions); err != nil {
		t.Fatal(err)
	}
	if suggestions.PendingAction == nil || suggestions.PendingAction.AbsorbedBy != holder.Id || suggestions.PendingAction.AbsorbedByName == "" {
		t.Fatalf("pending action = %+v, want it to name the holder", suggestions.PendingAction)
	}

	// Nothing is in the history until the commit.
	if rows := getHistory(t, app, proposalHolderEmail); len(rows) != 3 {
		t.Fatalf("history before the commit has %d rows, want 3", len(rows))
	}

	res = performTestAPIRequest(t, app, http.MethodDelete, "/api/collections/absorb_actions/records/"+action.Id, nil, token)
	mustStatus(t, res, http.StatusNoContent)

	rows := getHistory(t, app, proposalHolderEmail)
	if len(rows) != 4 {
		t.Fatalf("history after the commit has %d rows, want 4", len(rows))
	}
	row := rows[0]
	if row.CollectionName != "clients" || row.TargetID != clientBNRail || row.TargetName == "" {
		t.Errorf("history row = %+v", row)
	}
	if row.AbsorbedBy != holder.Id || row.CommittedBy != holder.Id || row.AbsorbedByName == "" || row.CommittedByName == "" {
		t.Errorf("history row does not name the holder: %+v", row)
	}
	if len(row.Absorbed) != 1 || row.Absorbed[0] == "" || row.CommittedAt == "" {
		t.Errorf("history row = %+v", row)
	}
	if _, err := app.FindRecordById("absorb_actions", action.Id); err == nil {
		t.Error("the action is still there after the commit")
	}
}

func TestAbsorbUndoWritesNoHistory(t *testing.T) {
	app := testutils.SetupTestApp(t)
	token := uncommitToken(t, proposalHolderEmail)

	res := performTestAPIRequest(t, app, http.MethodPost, "/api/clients/"+clientBNRail+"/absorb",
		proposalJSON(t, map[string]any{"ids_to_absorb": []string{sugBNRailCorp}}), token)
	mustStatus(t, res, http.StatusOK)
	res = performTestAPIRequest(t, app, http.MethodPost, "/api/clients/undo_absorb", nil, token)
	mustStatus(t, res, http.StatusOK)

	if rows := getHistory(t, app, proposalHolderEmail); len(rows) != 3 {
		t.Fatalf("history after an undo has %d rows, want 3", len(rows))
	}
}
