package migrations_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/core"
)

const purchaseOrdersDeleteRuleMigration = "1793673271_purchase_orders_delete_rule_unnumbered.go"

// All migrations must produce the canonical purchase_orders deleteRule, which
// refuses to delete purchase orders that have a po_number. Only the rule is
// compared because the fixture lists the fields in a different order.
func TestPurchaseOrdersMigrationsMatchCanonicalDeleteRule(t *testing.T) {
	app := core.NewBaseApp(core.BaseAppConfig{DataDir: t.TempDir()})
	if err := app.Bootstrap(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { app.ResetBootstrapState() })
	if err := app.RunAllMigrations(); err != nil {
		t.Fatal(err)
	}
	assertCanonicalPurchaseOrdersDeleteRule(t, app)
}

func TestPurchaseOrdersDeleteRuleMigrationRevert(t *testing.T) {
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
		if item.File == purchaseOrdersDeleteRuleMigration {
			migration = item
		}
	}
	if migration == nil {
		t.Fatal("purchase orders delete rule migration must be registered")
	}
	if err := migration.Down(app); err != nil {
		t.Fatal(err)
	}
	rule := migrationCollection(t, app, "purchase_orders").DeleteRule
	if rule == nil || *rule != "@request.auth.id = uid && status = 'Unapproved'" {
		t.Fatalf("down migration must restore the previous deleteRule; got %v", rule)
	}
	if err := migration.Up(app); err != nil {
		t.Fatal(err)
	}
	assertCanonicalPurchaseOrdersDeleteRule(t, app)
}

func assertCanonicalPurchaseOrdersDeleteRule(t *testing.T, app core.App) {
	t.Helper()
	collection := migrationCollection(t, app, "purchase_orders")
	want := fixtureRow(t, "data/_collections.csv", collection.Id)["deleteRule"]
	if collection.DeleteRule == nil || *collection.DeleteRule != want {
		t.Fatalf("deleteRule = %v, want fixture value %q", collection.DeleteRule, want)
	}
}
