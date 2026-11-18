package migrations_test

import (
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/core"
)

const rateSheetsCreateRuleMigration = "1795012495_rate_sheets_create_rule_job_only.go"

// The migrations give rate_sheets the create rule of the fixture: only a job
// holder, and only revision 0. The down function restores the rule that let
// rate_sheet_revise and admin create a sheet.
func TestRateSheetsCreateRuleMigration(t *testing.T) {
	app := core.NewBaseApp(core.BaseAppConfig{DataDir: t.TempDir()})
	if err := app.Bootstrap(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { app.ResetBootstrapState() })
	if err := app.RunAllMigrations(); err != nil {
		t.Fatal(err)
	}
	assertFixtureRateSheetsCreateRule(t, app)

	var migration *core.Migration
	for _, item := range core.AppMigrations.Items() {
		if item.File == rateSheetsCreateRuleMigration {
			migration = item
		}
	}
	if migration == nil {
		t.Fatal("rate sheets create rule migration must be registered")
	}
	if err := migration.Down(app); err != nil {
		t.Fatal(err)
	}
	rule := migrationCollection(t, app, "rate_sheets").CreateRule
	if rule == nil || !strings.Contains(*rule, "?= 'rate_sheet_revise'") || !strings.Contains(*rule, "@request.body.code:isset = false") {
		t.Fatalf("down migration must restore the previous createRule; got %v", rule)
	}
	if err := migration.Up(app); err != nil {
		t.Fatal(err)
	}
	assertFixtureRateSheetsCreateRule(t, app)
}

func assertFixtureRateSheetsCreateRule(t *testing.T, app core.App) {
	t.Helper()
	collection := migrationCollection(t, app, "rate_sheets")
	want := fixtureRow(t, "data/_collections.csv", collection.Id)["createRule"]
	if collection.CreateRule == nil || *collection.CreateRule != want {
		t.Fatalf("createRule = %v, want fixture value %q", collection.CreateRule, want)
	}
}
