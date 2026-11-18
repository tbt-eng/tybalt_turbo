package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"tybalt/internal/testutils"

	"github.com/pocketbase/pocketbase/tests"
)

// The rate sheets of a client's jobs (issue 207). The fixtures are rows of the
// seed CSV files (app/test_seed_data/data); no test changes them.
//
//   - BN Rail Corporation (eldtxi3i4h00k8r) has jobs on twelve revisions:
//     Rollout Rates (2026-11-00026) rev 0 is the current sheet of 91-8930 to
//     91-8934 and rev 1 of 91-8935; 2025 Standard Rates (c41ofep525bcacj) is
//     the current sheet of 91-8914, 91-8916 and 91-8918, and only a past one
//     of 91-8934; only a past period of 91-8914 is on Absorb History Source
//     (rsabhistsrc0001); 91-8914 is also on the standard sheet now and in two
//     of its periods.
//   - Lakehead Surveys, Inc. (clientcomma0001) has one job, 91-8940, on the
//     standard sheet, which jobs of three other clients also use.
//   - City of Ipperville (iq7ge8osf6kuj3l) has jobs, none with a rate sheet.

type clientRateSheetRow struct {
	ID            string `json:"id"`
	Code          string `json:"code"`
	Name          string `json:"name"`
	Revision      int    `json:"revision"`
	EffectiveDate string `json:"effective_date"`
	Active        bool   `json:"active"`
	CurrentJobs   int    `json:"current_jobs"`
	PastJobs      int    `json:"past_jobs"`
}

// clientRateSheets returns the route's sheets in order, and fails if a
// revision is listed twice.
func clientRateSheets(t *testing.T, app *tests.TestApp, client string) []clientRateSheetRow {
	t.Helper()
	status, body := absorbRequest(t, app, http.MethodGet, "/api/clients/"+client+"/rate_sheets", "author@soup.com", nil)
	if status != http.StatusOK {
		t.Fatalf("%s status = %d; body %s", client, status, body)
	}
	var rows []clientRateSheetRow
	if err := json.Unmarshal([]byte(body), &rows); err != nil {
		t.Fatal(err)
	}
	seen := map[string]bool{}
	for _, row := range rows {
		if seen[row.ID] {
			t.Errorf("%s lists %s twice", client, row.ID)
		}
		seen[row.ID] = true
	}
	return rows
}

func clientRateSheetByID(rows []clientRateSheetRow) map[string]clientRateSheetRow {
	byID := map[string]clientRateSheetRow{}
	for _, row := range rows {
		byID[row.ID] = row
	}
	return byID
}

func TestClientRateSheetsRoute(t *testing.T) {
	app := testutils.SetupTestApp(t)
	defer app.Cleanup()

	rows := clientRateSheets(t, app, "eldtxi3i4h00k8r")
	bn := clientRateSheetByID(rows)

	for id, want := range map[string]struct{ current, past int }{
		// Current only, on each of two revisions of one code: two rows.
		"rsrollout0r0001": {5, 0},
		"rsrollout1r0001": {1, 0},
		// Past only: a period of the job's history.
		"rsabhistsrc0001": {0, 1},
		// Both on one sheet. 91-8914 is on it now and in two periods, and
		// counts once, as current.
		"c41ofep525bcacj": {3, 1},
	} {
		row, ok := bn[id]
		if !ok {
			t.Errorf("BN Rail does not list %s", id)
			continue
		}
		if row.CurrentJobs != want.current || row.PastJobs != want.past {
			t.Errorf("%s = %d now, %d past; want %d now, %d past", id, row.CurrentJobs, row.PastJobs, want.current, want.past)
		}
	}

	// The row carries what the section shows.
	rollout := bn["rsrollout1r0001"]
	if rollout.Code != "2026-11-00026" || rollout.Revision != 1 || rollout.Name != "Rollout Rates" || rollout.EffectiveDate == "" || !rollout.Active {
		t.Errorf("rev 1 = %+v; want 2026-11-00026 rev 1 Rollout Rates, active, with a date", rollout)
	}
	if bn["rsrollout0r0001"].Active {
		t.Error("rev 0 is active; want inactive")
	}

	// Sheets with current jobs come first, then code and revision descending.
	seenPast := false
	for i, row := range rows {
		if row.CurrentJobs == 0 {
			seenPast = true
		} else if seenPast {
			t.Errorf("%s has current jobs after a sheet with none", row.ID)
		}
		if i == 0 || (rows[i-1].CurrentJobs > 0) != (row.CurrentJobs > 0) {
			continue
		}
		prev := rows[i-1]
		if prev.Code < row.Code || (prev.Code == row.Code && prev.Revision < row.Revision) {
			t.Errorf("%s (%s rev %d) before %s (%s rev %d); want code and revision descending", prev.ID, prev.Code, prev.Revision, row.ID, row.Code, row.Revision)
		}
	}

	// Another client's jobs are not counted: the standard sheet has jobs of
	// four clients, and Lakehead Surveys has one of them.
	lakehead := clientRateSheets(t, app, "clientcomma0001")
	if len(lakehead) != 1 || lakehead[0].ID != "c41ofep525bcacj" || lakehead[0].CurrentJobs != 1 || lakehead[0].PastJobs != 0 {
		t.Errorf("Lakehead Surveys sheets = %+v; want only the standard sheet, 1 now", lakehead)
	}

	// A client whose jobs have no rate sheet gets an empty list, not null.
	status, body := absorbRequest(t, app, http.MethodGet, "/api/clients/iq7ge8osf6kuj3l/rate_sheets", "author@soup.com", nil)
	if status != http.StatusOK || strings.TrimSpace(body) != "[]" {
		t.Errorf("City of Ipperville = %d %s; want 200 []", status, body)
	}
}

func TestClientRateSheetsRouteAuth(t *testing.T) {
	app := testutils.SetupTestApp(t)
	defer app.Cleanup()
	if status, _ := absorbRequest(t, app, http.MethodGet, "/api/clients/clientcomma0001/rate_sheets", "", nil); status != http.StatusUnauthorized {
		t.Errorf("signed out status = %d; want 401", status)
	}
	if status, _ := absorbRequest(t, app, http.MethodGet, "/api/clients/nosuchclient001/rate_sheets", "author@soup.com", nil); status != http.StatusNotFound {
		t.Errorf("missing client status = %d; want 404", status)
	}
}
