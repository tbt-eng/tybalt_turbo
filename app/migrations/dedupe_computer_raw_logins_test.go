package migrations_test

import (
	"reflect"
	"testing"

	"tybalt/internal/testseed"

	"github.com/pocketbase/pocketbase/core"
)

const dedupeComputerRawLoginsMigration = "1795823645_dedupe_computer_raw_logins.go"

// The migration keeps the newest raw login of each computer name and leaves
// the raw logins without a name alone. It runs here against the seeded
// fixtures; the seeded template ran it before the fixtures were loaded, so the
// fixtures still hold the duplicates of data stored before the change:
// FAKE-PC-0097 has cmprawlogin0005 (2026-10-19 13:00) and cmprawlogin0006
// and cmprawlogin0007 (both 2026-10-21 13:00, so the larger id wins).
// cmprawlogin0004 and cmprawlogin0008 have no name.
func TestDedupeComputerRawLoginsMigration(t *testing.T) {
	app := testseed.NewSeededTestApp(t)

	var migration *core.Migration
	for _, item := range core.AppMigrations.Items() {
		if item.File == dedupeComputerRawLoginsMigration {
			migration = item
		}
	}
	if migration == nil {
		t.Fatal("dedupe computer raw logins migration must be registered")
	}

	want := []string{
		"cmprawlogin0001", // FAKE-PC-0001, the only one
		"cmprawlogin0002", // FAKE-PC-0099, the only one
		"cmprawlogin0003", // FAKE-PC-0098, the only one
		"cmprawlogin0004", // no name
		"cmprawlogin0007", // FAKE-PC-0097, the newest
		"cmprawlogin0008", // no name
	}
	before := rawLoginIDs(t, app)
	if len(before) != len(want)+2 {
		t.Fatalf("fixture raw logins = %v, want %d", before, len(want)+2)
	}

	if err := migration.Up(app); err != nil {
		t.Fatal(err)
	}
	if got := rawLoginIDs(t, app); !reflect.DeepEqual(got, want) {
		t.Fatalf("raw logins after the migration = %v, want %v", got, want)
	}

	// Running it again deletes nothing more.
	if err := migration.Up(app); err != nil {
		t.Fatal(err)
	}
	if got := rawLoginIDs(t, app); !reflect.DeepEqual(got, want) {
		t.Fatalf("raw logins after a second run = %v, want %v", got, want)
	}

	// The down migration does nothing.
	if err := migration.Down(app); err != nil {
		t.Fatal(err)
	}
	if got := rawLoginIDs(t, app); !reflect.DeepEqual(got, want) {
		t.Fatalf("raw logins after the down migration = %v, want %v", got, want)
	}
}

func rawLoginIDs(t *testing.T, app core.App) []string {
	t.Helper()
	var rows []struct {
		ID string `db:"id"`
	}
	if err := app.DB().NewQuery("SELECT id FROM computer_raw_logins ORDER BY id").All(&rows); err != nil {
		t.Fatal(err)
	}
	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.ID)
	}
	return ids
}
