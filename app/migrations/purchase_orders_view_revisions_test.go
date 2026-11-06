package migrations_test

import (
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/core"
)

const purchaseOrdersViewRevisionsMigration = "1793903864_purchase_orders_view_revisions.go"

// All migrations must produce the canonical purchase_orders list and view
// rules, which let any signed-in user see a purchase order under revision.
func TestPurchaseOrdersMigrationsMatchCanonicalRevisionVisibility(t *testing.T) {
	app := core.NewBaseApp(core.BaseAppConfig{DataDir: t.TempDir()})
	if err := app.Bootstrap(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { app.ResetBootstrapState() })
	if err := app.RunAllMigrations(); err != nil {
		t.Fatal(err)
	}
	assertCanonicalPurchaseOrdersReadRules(t, app)
}

func TestPurchaseOrdersViewRevisionsMigrationRevert(t *testing.T) {
	app := core.NewBaseApp(core.BaseAppConfig{DataDir: t.TempDir()})
	if err := app.Bootstrap(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { app.ResetBootstrapState() })
	if err := app.RunAllMigrations(); err != nil {
		t.Fatal(err)
	}
	var migration *core.Migration
	for _, item := range core.AppMigrations.Items() {
		if item.File == purchaseOrdersViewRevisionsMigration {
			migration = item
		}
	}
	if migration == nil {
		t.Fatal("purchase orders view revisions migration must be registered")
	}
	if err := migration.Down(app); err != nil {
		t.Fatal(err)
	}
	collection := migrationCollection(t, app, "purchase_orders")
	for name, rule := range map[string]*string{"listRule": collection.ListRule, "viewRule": collection.ViewRule} {
		if rule == nil || strings.Contains(*rule, `po_number != ""`) {
			t.Fatalf("down migration must remove the revision clause from %s; got %v", name, rule)
		}
	}
	if err := migration.Up(app); err != nil {
		t.Fatal(err)
	}
	assertCanonicalPurchaseOrdersReadRules(t, app)
}

// The clause must apply to a rule that doesn't have the shape the current rules
// have, such as one edited in the admin UI, and leave it unchanged when reverted.
func TestPurchaseOrdersViewRevisionsMigrationOnAnOtherRule(t *testing.T) {
	app := core.NewBaseApp(core.BaseAppConfig{DataDir: t.TempDir()})
	if err := app.Bootstrap(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { app.ResetBootstrapState() })
	if err := app.RunAllMigrations(); err != nil {
		t.Fatal(err)
	}
	var migration *core.Migration
	for _, item := range core.AppMigrations.Items() {
		if item.File == purchaseOrdersViewRevisionsMigration {
			migration = item
		}
	}
	if err := migration.Down(app); err != nil {
		t.Fatal(err)
	}
	collection := migrationCollection(t, app, "purchase_orders")
	edited := `@request.auth.id = uid`
	collection.ListRule = &edited
	if err := app.Save(collection); err != nil {
		t.Fatal(err)
	}
	if err := migration.Up(app); err != nil {
		t.Fatal(err)
	}
	rule := migrationCollection(t, app, "purchase_orders").ListRule
	if rule == nil || !strings.Contains(*rule, `po_number != ""`) || !strings.HasSuffix(*rule, edited) {
		t.Fatalf("up migration on an edited rule produced %v", rule)
	}
	if err := migration.Down(app); err != nil {
		t.Fatal(err)
	}
	rule = migrationCollection(t, app, "purchase_orders").ListRule
	if rule == nil || *rule != edited {
		t.Fatalf("down migration must restore the edited rule; got %v", rule)
	}
}

func assertCanonicalPurchaseOrdersReadRules(t *testing.T, app core.App) {
	t.Helper()
	collection := migrationCollection(t, app, "purchase_orders")
	row := fixtureRow(t, "data/_collections.csv", collection.Id)
	if collection.ListRule == nil || *collection.ListRule != row["listRule"] {
		t.Fatalf("listRule = %v, want fixture value %q", collection.ListRule, row["listRule"])
	}
	if collection.ViewRule == nil || *collection.ViewRule != row["viewRule"] {
		t.Fatalf("viewRule = %v, want fixture value %q", collection.ViewRule, row["viewRule"])
	}
}
