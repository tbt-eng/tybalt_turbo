package migrations_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/core"
)

const jobsContactManagerOptionalMigration = "1793286261_jobs_contact_manager_optional.go"

// The down migration restores the required flags, and the up migration
// returns the jobs schema to the canonical one.
func TestJobsContactManagerOptionalMigrationRevert(t *testing.T) {
	app := core.NewBaseApp(core.BaseAppConfig{DataDir: t.TempDir()})
	if err := app.Bootstrap(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { app.ResetBootstrapState() })
	if err := app.RunAllMigrations(); err != nil {
		t.Fatal(err)
	}
	assertCanonicalCollections(t, app, "jobs")

	var migration *core.Migration
	for _, item := range core.AppMigrations.Items() {
		if item.File == jobsContactManagerOptionalMigration {
			migration = item
		}
	}
	if migration == nil {
		t.Fatal("jobs contact/manager migration must be registered")
	}

	required := func(want bool) {
		t.Helper()
		jobs := migrationCollection(t, app, "jobs")
		for _, name := range []string{"contact", "manager"} {
			if got := jobs.Fields.GetByName(name).(*core.RelationField).Required; got != want {
				t.Fatalf("jobs.%s required = %v, want %v", name, got, want)
			}
		}
	}

	required(false)
	if err := migration.Down(app); err != nil {
		t.Fatal(err)
	}
	required(true)
	if err := migration.Up(app); err != nil {
		t.Fatal(err)
	}
	required(false)
	assertCanonicalCollections(t, app, "jobs")
}
