package migrations_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/core"
)

const clientsBusdevLeadOptionalMigration = "1795463617_clients_busdev_lead_optional.go"

// The down migration restores the required flag, and the up migration
// returns the clients schema to the canonical one.
func TestClientsBusdevLeadOptionalMigrationRevert(t *testing.T) {
	app := core.NewBaseApp(core.BaseAppConfig{DataDir: t.TempDir()})
	if err := app.Bootstrap(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { app.ResetBootstrapState() })
	if err := app.RunAllMigrations(); err != nil {
		t.Fatal(err)
	}
	assertCanonicalCollections(t, app, "clients")

	var migration *core.Migration
	for _, item := range core.AppMigrations.Items() {
		if item.File == clientsBusdevLeadOptionalMigration {
			migration = item
		}
	}
	if migration == nil {
		t.Fatal("clients business development lead migration must be registered")
	}

	required := func(want bool) {
		t.Helper()
		clients := migrationCollection(t, app, "clients")
		if got := clients.Fields.GetByName("business_development_lead").(*core.RelationField).Required; got != want {
			t.Fatalf("clients.business_development_lead required = %v, want %v", got, want)
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
	assertCanonicalCollections(t, app, "clients")
}
