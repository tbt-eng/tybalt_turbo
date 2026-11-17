package main

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"imports/load"
)

// The rate sheets of the legacy export, as planRateSheetCodesForImport reads
// them: "Standard" has two revisions, "Client" one.
var importedRateSheets = []load.RateSheet{
	{Id: "rsstandard00001", Name: "Standard", Revision: 1},
	{Id: "rsstandard00000", Name: "Standard", Revision: 0},
	{Id: "rsclient0000000", Name: "Client", Revision: 0},
}

// importTestDB makes a database with the two tables that the code plan
// reads and writes, in the form the app's migrations give them.
func importTestDB(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "data.db")
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	for _, statement := range []string{
		"CREATE TABLE rate_sheets (id TEXT PRIMARY KEY, code TEXT DEFAULT '' NOT NULL, name TEXT DEFAULT '' NOT NULL, effective_date TEXT, revision NUMERIC, active BOOLEAN)",
		"CREATE UNIQUE INDEX idx_rate_sheets_code_revision ON rate_sheets (code, revision)",
		"CREATE TABLE rate_sheet_codes (id TEXT PRIMARY KEY, code TEXT NOT NULL, month TEXT NOT NULL, serial NUMERIC NOT NULL, created TEXT)",
		"CREATE UNIQUE INDEX idx_rate_sheet_codes_code ON rate_sheet_codes (code)",
	} {
		if _, err := db.Exec(statement); err != nil {
			t.Fatal(err)
		}
	}
	return path
}

// runRateSheetImport does what the import does with --clear_rates: read the
// codes, clear rate_sheets, plan the codes, insert the sheets and record the
// codes.
func runRateSheetImport(t *testing.T, path string, items []load.RateSheet, now time.Time) map[string]string {
	t.Helper()
	snapshot, err := snapshotRateSheetCodes(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := deleteAllFromTables(path, []string{"rate_sheets"}); err != nil {
		t.Fatal(err)
	}
	codes, err := planRateSheetCodesForImport(path, items, snapshot, now)
	if err != nil {
		t.Fatal(err)
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	for _, item := range items {
		if _, err := db.Exec("INSERT INTO rate_sheets (id, code, name, revision) VALUES (?, ?, ?, ?)", item.Id, codes[item.Id], item.Name, item.Revision); err != nil {
			t.Fatal(err)
		}
	}
	if err := recordIssuedRateSheetCodes(path, codes); err != nil {
		t.Fatal(err)
	}
	return codes
}

func issuedCodes(t *testing.T, path string) string {
	t.Helper()
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT code, month, serial FROM rate_sheet_codes ORDER BY code")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	result := ""
	for rows.Next() {
		var code, month string
		var serial int
		if err := rows.Scan(&code, &month, &serial); err != nil {
			t.Fatal(err)
		}
		result += fmt.Sprintf("%s/%s/%d;", code, month, serial)
	}
	return result
}

// A first import gives each sheet a code in the month of the import, shared
// by its revisions. A second import, in a later month, keeps every code, and
// a sheet new to it gets the next code of that month.
func TestImportGivesRateSheetsCodesAndKeepsThem(t *testing.T) {
	path := importTestDB(t)

	first := runRateSheetImport(t, path, importedRateSheets, time.Date(2026, 11, 20, 0, 0, 0, 0, time.UTC))
	want := map[string]string{
		"rsclient0000000": "2026-11-00001",
		"rsstandard00000": "2026-11-00002",
		"rsstandard00001": "2026-11-00002",
	}
	if fmt.Sprint(first) != fmt.Sprint(want) {
		t.Fatalf("first import codes = %v; want %v", first, want)
	}

	second := append([]load.RateSheet{}, importedRateSheets...)
	second = append(second,
		load.RateSheet{Id: "rsstandard00002", Name: "Standard", Revision: 2},
		load.RateSheet{Id: "rsnewsheet00000", Name: "New Sheet", Revision: 0},
	)
	codes := runRateSheetImport(t, path, second, time.Date(2027, 1, 5, 0, 0, 0, 0, time.UTC))
	want["rsstandard00002"] = "2026-11-00002"
	want["rsnewsheet00000"] = "2027-01-00001"
	if fmt.Sprint(codes) != fmt.Sprint(want) {
		t.Fatalf("second import codes = %v; want %v", codes, want)
	}
	if got := issuedCodes(t, path); got != "2026-11-00001/2026-11/1;2026-11-00002/2026-11/2;2027-01-00001/2027-01/1;" {
		t.Fatalf("issued codes = %s", got)
	}

	// A sheet that is gone from the export leaves its code issued, so a
	// later new sheet does not get it.
	third := runRateSheetImport(t, path, []load.RateSheet{{Id: "rsother00000000", Name: "Other", Revision: 0}}, time.Date(2027, 1, 6, 0, 0, 0, 0, time.UTC))
	if third["rsother00000000"] != "2027-01-00002" {
		t.Fatalf("a new sheet after a removed one got %s; want 2027-01-00002", third["rsother00000000"])
	}
}
