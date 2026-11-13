package routes

import (
	"context"
	"testing"

	"tybalt/hooks"
	"tybalt/internal/testseed"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type pricedEntry struct {
	Date      string   `db:"date"`
	Role      string   `db:"role"`
	RateSheet string   `db:"rate_sheet"`
	SheetRate *float64 `db:"sheet_rate"`
}

func pricedEntries(t *testing.T, app core.App, job string) map[string]pricedEntry {
	t.Helper()
	var rows []pricedEntry
	err := app.DB().NewQuery(singleJobScopeQuery + jobPricedTimeEntriesQuery + `
		SELECT date, COALESCE(role, '') AS role, COALESCE(rate_sheet, '') AS rate_sheet, sheet_rate
		FROM priced_entries WHERE COALESCE(role, '') != ''
	`).Bind(dbx.Params{"job_id": job, "start_date": "", "end_date": "2026-12-31"}).All(&rows)
	if err != nil {
		t.Fatal(err)
	}
	result := map[string]pricedEntry{}
	for _, row := range rows {
		result[row.Date+" "+row.Role] = row
	}
	return result
}

// TestPricingFollowsRateSheetHistory checks that each entry is priced by the
// sheet in effect on its date, and that a project with no history uses its
// jobs.rate_sheet for all time.
//
// Fixtures, all rows of the seed CSV files (app/test_seed_data/data):
//   - rsvaluation0001: inactive sheet that rates role 0l1llckilwz1mqz at 100
//     and has no rate for 2bd6pxxof7p649x.
//   - c41ofep525bcacj: active, effective 2025-01-07, rates them at 160 and 90.
//   - jobpricenohist1: on rsvaluation0001, with no rows in job_rate_sheets, as
//     for an imported project.
//   - jobpricehist001: on c41ofep525bcacj. Its opening period is on
//     rsvaluation0001, and c41ofep525bcacj applies from 2025-06-03.
//   - Each project has one entry for 0l1llckilwz1mqz on 2025-06-02, and one
//     for each role on 2025-06-03.
//
// The only write is the opening row that the test asks the hook to make,
// which is the behaviour under test.
func TestPricingFollowsRateSheetHistory(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()
	hooks.AddHooks(app)

	const noHistory, withHistory = "jobpricenohist1", "jobpricehist001"
	const admin, field = "2bd6pxxof7p649x", "0l1llckilwz1mqz"

	// No history: every entry uses jobs.rate_sheet.
	checkNoHistory := func(when string) {
		t.Helper()
		priced := pricedEntries(t, app, noHistory)
		for _, date := range []string{"2025-06-02", "2025-06-03"} {
			if got := priced[date+" "+field]; got.RateSheet != "rsvaluation0001" || got.SheetRate == nil || *got.SheetRate != 100 {
				t.Fatalf("%s, %s entry = %+v; want rsvaluation0001 at 100", when, date, got)
			}
		}
		if got := priced["2025-06-03 "+admin]; got.RateSheet != "rsvaluation0001" || got.SheetRate != nil {
			t.Fatalf("%s, the admin role should have no sheet rate on rsvaluation0001, got %+v", when, got)
		}
	}
	checkNoHistory("without history")

	// The opening row keeps the inactive sheet, and pricing does not change.
	if err := hooks.EnsureOpeningRateSheetRow(app, context.Background(), mustFindJob(t, app, noHistory), ""); err != nil {
		t.Fatalf("opening row on an inactive sheet: %v", err)
	}
	if n, err := app.CountRecords("job_rate_sheets", dbx.HashExp{"job": noHistory}); err != nil || n != 1 {
		t.Fatalf("opening rows = %d (%v); want 1", n, err)
	}
	checkNoHistory("with only an opening row")

	// With history: the opening period prices 2025-06-02, and
	// c41ofep525bcacj prices 2025-06-03.
	priced := pricedEntries(t, app, withHistory)
	for key, want := range map[string]struct {
		sheet string
		rate  float64
	}{
		"2025-06-02 " + field: {"rsvaluation0001", 100},
		"2025-06-03 " + field: {"c41ofep525bcacj", 160},
		"2025-06-03 " + admin: {"c41ofep525bcacj", 90},
	} {
		got := priced[key]
		if got.RateSheet != want.sheet || got.SheetRate == nil || *got.SheetRate != want.rate {
			t.Errorf("%s = %+v; want %s at %v", key, got, want.sheet, want.rate)
		}
	}
}

func mustFindJob(t *testing.T, app core.App, id string) *core.Record {
	t.Helper()
	job, err := app.FindRecordById("jobs", id)
	if err != nil {
		t.Fatal(err)
	}
	return job
}
