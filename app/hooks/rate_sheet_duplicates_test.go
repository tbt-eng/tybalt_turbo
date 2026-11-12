package hooks

import (
	"testing"

	"tybalt/internal/testseed"

	"github.com/pocketbase/pocketbase/core"
)

// Fixtures used here (app/test_seed_data):
//   - c41ofep525bcacj: "2025 Standard Rates", active, an entry for every role
//   - rs_dupact_001: "Duplicate Activation Sheet", inactive, the same entries
//   - rs_complete_all_001: "Complete Hook Sheet", inactive, every role at 100/150
const standardSheetID = "c41ofep525bcacj"
const standardSheetName = "2025 Standard Rates"

func loadEntries(t *testing.T, app core.App, sheetID string) []RateSheetEntryValues {
	t.Helper()
	entries, err := LoadRateSheetEntryValues(app, sheetID)
	if err != nil {
		t.Fatalf("load entries of %s: %v", sheetID, err)
	}
	if len(entries) == 0 {
		t.Fatalf("expected entries on %s", sheetID)
	}
	return entries
}

func findMatch(matches []RateSheetMatch, id string) (RateSheetMatch, bool) {
	for _, m := range matches {
		if m.ID == id {
			return m, true
		}
	}
	return RateSheetMatch{}, false
}

func TestCompareWithActiveRateSheets_Identical(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()

	entries := loadEntries(t, app, standardSheetID)
	got, err := CompareWithActiveRateSheets(app, entries, "Some New Sheet")
	if err != nil {
		t.Fatalf("compare: %v", err)
	}

	match, ok := findMatch(got.Exact, standardSheetID)
	if !ok {
		t.Fatalf("expected %s in exact matches, got %+v", standardSheetID, got.Exact)
	}
	if match.Name != standardSheetName || match.Revision != 0 || match.DifferingRoles != 0 {
		t.Errorf("unexpected match: %+v", match)
	}
	if _, ok := findMatch(got.Near, standardSheetID); ok {
		t.Errorf("an exact match must not be listed as near")
	}
}

func TestCompareWithActiveRateSheets_DiffersByRate(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()

	entries := loadEntries(t, app, standardSheetID)
	entries[0].Rate++
	entries[1].Rate++

	got, err := CompareWithActiveRateSheets(app, entries, "Some New Sheet")
	if err != nil {
		t.Fatalf("compare: %v", err)
	}
	if len(got.Exact) != 0 {
		t.Fatalf("expected no exact match, got %+v", got.Exact)
	}
	match, ok := findMatch(got.Near, standardSheetID)
	if !ok || match.DifferingRoles != 2 {
		t.Fatalf("expected %s near with 2 differing roles, got %+v", standardSheetID, got.Near)
	}
	if got.Near[0].ID != standardSheetID {
		t.Errorf("expected the closest sheet first, got %+v", got.Near[0])
	}
}

func TestCompareWithActiveRateSheets_DiffersByOvertimeRate(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()

	entries := loadEntries(t, app, standardSheetID)
	entries[0].OvertimeRate += 0.5

	got, err := CompareWithActiveRateSheets(app, entries, "Some New Sheet")
	if err != nil {
		t.Fatalf("compare: %v", err)
	}
	if len(got.Exact) != 0 {
		t.Fatalf("expected no exact match, got %+v", got.Exact)
	}
	if match, ok := findMatch(got.Near, standardSheetID); !ok || match.DifferingRoles != 1 {
		t.Fatalf("expected %s near with 1 differing role, got %+v", standardSheetID, got.Near)
	}
}

func TestCompareWithActiveRateSheets_DiffersByRoleSet(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()

	entries := loadEntries(t, app, standardSheetID)

	// A role missing from the new entries counts as a difference.
	got, err := CompareWithActiveRateSheets(app, entries[1:], "Some New Sheet")
	if err != nil {
		t.Fatalf("compare: %v", err)
	}
	if len(got.Exact) != 0 {
		t.Fatalf("expected no exact match with a role missing, got %+v", got.Exact)
	}
	if match, ok := findMatch(got.Near, standardSheetID); !ok || match.DifferingRoles != 1 {
		t.Fatalf("expected 1 differing role with a role missing, got %+v", got.Near)
	}

	// A role that only the new entries have counts as a difference too.
	extra := append(append([]RateSheetEntryValues{}, entries...), RateSheetEntryValues{Role: "not_on_the_sheet", Rate: 100, OvertimeRate: 130})
	got, err = CompareWithActiveRateSheets(app, extra, "Some New Sheet")
	if err != nil {
		t.Fatalf("compare: %v", err)
	}
	if len(got.Exact) != 0 {
		t.Fatalf("expected no exact match with an extra role, got %+v", got.Exact)
	}
	if match, ok := findMatch(got.Near, standardSheetID); !ok || match.DifferingRoles != 1 {
		t.Fatalf("expected 1 differing role with an extra role, got %+v", got.Near)
	}
}

func TestCompareWithActiveRateSheets_ExcludesSameName(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()

	entries := loadEntries(t, app, standardSheetID)
	got, err := CompareWithActiveRateSheets(app, entries, standardSheetName)
	if err != nil {
		t.Fatalf("compare: %v", err)
	}
	if _, ok := findMatch(got.Exact, standardSheetID); ok {
		t.Errorf("a sheet with the excluded name must not be an exact match")
	}
	if _, ok := findMatch(got.Near, standardSheetID); ok {
		t.Errorf("a sheet with the excluded name must not be a near match")
	}
}

func TestCompareWithActiveRateSheets_IgnoresInactive(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()

	// rs_dupact_001 has the same entries as the standard sheet but is inactive.
	entries := loadEntries(t, app, "rs_dupact_001")
	got, err := CompareWithActiveRateSheets(app, entries, "Some New Sheet")
	if err != nil {
		t.Fatalf("compare: %v", err)
	}
	if len(got.Exact) != 1 || got.Exact[0].ID != standardSheetID {
		t.Fatalf("expected only the active standard sheet as an exact match, got %+v", got.Exact)
	}

	// rs_complete_all_001 is inactive, so its identical entries never match.
	entries = loadEntries(t, app, "rs_complete_all_001")
	got, err = CompareWithActiveRateSheets(app, entries, "Some New Sheet")
	if err != nil {
		t.Fatalf("compare: %v", err)
	}
	if len(got.Exact) != 0 {
		t.Fatalf("expected no exact match, got %+v", got.Exact)
	}
	for _, m := range append(got.Exact, got.Near...) {
		if m.ID == "rs_complete_all_001" || m.ID == "rs_dupact_001" {
			t.Errorf("inactive sheet %s must not be compared", m.ID)
		}
	}
}

func TestCompareWithActiveRateSheets_InTransaction(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()

	err := app.RunInTransaction(func(txApp core.App) error {
		entries := loadEntries(t, txApp, standardSheetID)
		got, err := CompareWithActiveRateSheets(txApp, entries, "Some New Sheet")
		if err != nil {
			return err
		}
		if _, ok := findMatch(got.Exact, standardSheetID); !ok {
			t.Errorf("expected an exact match inside a transaction, got %+v", got.Exact)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("transaction: %v", err)
	}
}

func TestRateSheetComparison_Closest(t *testing.T) {
	c := RateSheetComparison{Near: []RateSheetMatch{
		{ID: "a", DifferingRoles: 1},
		{ID: "b", DifferingRoles: 2},
		{ID: "c", DifferingRoles: 2},
		{ID: "d", DifferingRoles: 4},
		{ID: "e", DifferingRoles: 9},
	}}

	got := c.Closest(5, 3)
	if len(got) != 3 || got[0].ID != "a" || got[2].ID != "c" {
		t.Errorf("expected the 3 closest, got %+v", got)
	}

	got = c.Closest(3, 10)
	if len(got) != 3 {
		t.Errorf("expected only sheets within 3 roles, got %+v", got)
	}

	if got := (RateSheetComparison{}).Closest(5, 3); got == nil || len(got) != 0 {
		t.Errorf("expected an empty, non-nil slice, got %#v", got)
	}
}

func TestDuplicateRateSheetError(t *testing.T) {
	err := DuplicateRateSheetError("entries", []RateSheetMatch{{ID: "x1", Name: "Client X", Revision: 2}})
	if err.Status != 409 {
		t.Errorf("expected status 409, got %d", err.Status)
	}
	field, ok := err.Data["entries"]
	if !ok || field.Code != "duplicate_rate_sheet" {
		t.Fatalf("expected a duplicate_rate_sheet error on entries, got %+v", err.Data)
	}
	data, ok := field.Data.(map[string]any)
	if !ok {
		t.Fatalf("expected map data, got %T", field.Data)
	}
	matches, ok := data["matches"].([]RateSheetMatch)
	if !ok || len(matches) != 1 || matches[0].ID != "x1" {
		t.Errorf("expected the match in the data, got %+v", data)
	}
}
