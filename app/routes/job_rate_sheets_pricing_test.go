package routes

import (
	"context"
	"strings"
	"testing"

	"tybalt/hooks"
	"tybalt/internal/testseed"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

// copyTimeEntry inserts a copy of a time entry on another date. It writes the
// row directly, since pricing reads the table and the time entry hooks would
// check unrelated rules.
func copyTimeEntry(t *testing.T, app core.App, sourceID, newID, date string) {
	t.Helper()
	var cols []struct {
		Name string `db:"name"`
	}
	if err := app.DB().NewQuery("SELECT name FROM pragma_table_info('time_entries')").All(&cols); err != nil {
		t.Fatal(err)
	}
	names := []string{}
	values := []string{}
	for _, c := range cols {
		names = append(names, "["+c.Name+"]")
		switch c.Name {
		case "id":
			values = append(values, "{:id}")
		case "date":
			values = append(values, "{:date}")
		default:
			values = append(values, "["+c.Name+"]")
		}
	}
	q := "INSERT INTO time_entries (" + strings.Join(names, ", ") + ") SELECT " + strings.Join(values, ", ") + " FROM time_entries WHERE id = {:source}"
	if _, err := app.DB().NewQuery(q).Bind(dbx.Params{"id": newID, "date": date, "source": sourceID}).Execute(); err != nil {
		t.Fatal(err)
	}
}

func addRateSheetRow(t *testing.T, app core.App, job, sheet, date string) {
	t.Helper()
	collection, err := app.FindCollectionByNameOrId("job_rate_sheets")
	if err != nil {
		t.Fatal(err)
	}
	row := core.NewRecord(collection)
	row.Set("job", job)
	row.Set("rate_sheet", sheet)
	row.Set("effective_date", date)
	ctx := hooks.WithProjectWorkflowMutation(context.Background(), hooks.ProjectWorkflowMutationRateSheet)
	if err := app.SaveWithContext(ctx, row); err != nil {
		t.Fatalf("add %s from %q: %v", sheet, date, err)
	}
}

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
// Fixtures: jobwipreason001 is on the inactive sheet rsvaluation0001, which
// rates role 0l1llckilwz1mqz at 100 and has no rate for 2bd6pxxof7p649x.
// c41ofep525bcacj (active, effective 2025-01-07) rates them at 160 and 90.
func TestPricingFollowsRateSheetHistory(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()
	hooks.AddHooks(app)

	const job = "jobwipreason001"
	const admin, field = "2bd6pxxof7p649x", "0l1llckilwz1mqz"
	var source struct {
		ID string `db:"id"`
	}
	if err := app.DB().NewQuery("SELECT id FROM time_entries WHERE job = {:job} AND role = {:role} AND date = '2026-01-02'").
		Bind(dbx.Params{"job": job, "role": field}).One(&source); err != nil {
		t.Fatal(err)
	}
	copyTimeEntry(t, app, source.ID, "tecopyearly0001", "2026-01-01")

	// No history, as for an imported project: every entry uses
	// jobs.rate_sheet. The migration gave this fixture an opening row.
	if _, err := app.DB().NewQuery("DELETE FROM job_rate_sheets WHERE job = {:job}").Bind(dbx.Params{"job": job}).Execute(); err != nil {
		t.Fatal(err)
	}
	before := pricedEntries(t, app, job)
	if got := before["2026-01-01 "+field]; got.RateSheet != "rsvaluation0001" || got.SheetRate == nil || *got.SheetRate != 100 {
		t.Fatalf("without history, 2026-01-01 entry = %+v; want rsvaluation0001 at 100", got)
	}
	if got := before["2026-01-02 "+admin]; got.SheetRate != nil {
		t.Fatalf("without history, admin role should have no sheet rate, got %+v", got)
	}

	// The opening row keeps the inactive sheet; c41ofep525bcacj applies from
	// 2026-01-02.
	if err := hooks.EnsureOpeningRateSheetRow(app, context.Background(), mustFindJob(t, app, job), ""); err != nil {
		t.Fatalf("opening row on an inactive sheet: %v", err)
	}
	addRateSheetRow(t, app, job, "c41ofep525bcacj", "2026-01-02")

	after := pricedEntries(t, app, job)
	for key, want := range map[string]struct {
		sheet string
		rate  float64
	}{
		"2026-01-01 " + field: {"rsvaluation0001", 100},
		"2026-01-02 " + field: {"c41ofep525bcacj", 160},
		"2026-01-02 " + admin: {"c41ofep525bcacj", 90},
	} {
		got := after[key]
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
