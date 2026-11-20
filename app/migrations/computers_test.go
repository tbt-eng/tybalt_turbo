package migrations_test

import (
	"encoding/json"
	"reflect"
	"slices"
	"testing"

	"github.com/pocketbase/pocketbase/core"
)

const computersMigration = "1795200837_computers.go"

// The migration builds what the fixture _collections rows hold, so tests run
// against the same schema as a migrated database.
func TestComputersMigrationMatchesFixture(t *testing.T) {
	app := newMigratedApp(t)
	for _, name := range []string{"computers", "computer_logins", "computer_raw_logins", "admin_profiles", "machine_secrets"} {
		collection := migrationCollection(t, app, name)
		row := fixtureRow(t, "data/_collections.csv", collection.Id)
		for column, value := range map[string]any{"fields": collection.Fields, "indexes": collection.Indexes} {
			migrated, err := json.Marshal(value)
			if err != nil {
				t.Fatal(err)
			}
			var got, want any
			if err := json.Unmarshal(migrated, &got); err != nil {
				t.Fatal(err)
			}
			if err := json.Unmarshal([]byte(row[column]), &want); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("%s.%s differs from the fixture\n migrated: %s\n fixture:  %s", name, column, migrated, row[column])
			}
		}
		for column, rule := range map[string]*string{
			"listRule": collection.ListRule, "viewRule": collection.ViewRule, "createRule": collection.CreateRule,
			"updateRule": collection.UpdateRule, "deleteRule": collection.DeleteRule,
		} {
			want := row[column]
			got := `\N`
			if rule != nil {
				got = *rule
			}
			if got != want {
				t.Errorf("%s.%s = %q, fixture %q", name, column, got, want)
			}
		}
	}

	claim, err := app.FindFirstRecordByData("claims", "name", "computers")
	if err != nil {
		t.Fatal(err)
	}
	if want := fixtureRow(t, "data/claims.csv", claim.Id); want["description"] != claim.GetString("description") {
		t.Errorf("claim description = %q, fixture %q", claim.GetString("description"), want["description"])
	}
}

// Only server code writes the collections, and the down migration undoes it all.
func TestComputersMigrationRulesAndRollback(t *testing.T) {
	app := newMigratedApp(t)
	for _, name := range []string{"computers", "computer_logins", "computer_raw_logins"} {
		collection := migrationCollection(t, app, name)
		if collection.CreateRule != nil || collection.UpdateRule != nil || collection.DeleteRule != nil {
			t.Errorf("%s must have no create, update or delete rule", name)
		}
		if collection.ListRule == nil || collection.ViewRule == nil {
			t.Errorf("%s must be listable by a claim holder", name)
		}
	}
	computer, ok := migrationCollection(t, app, "computer_logins").Fields.GetByName("computer").(*core.RelationField)
	if !ok || !computer.Required || !computer.CascadeDelete {
		t.Errorf("computer_logins.computer must be a required relation that cascades, got %+v", computer)
	}

	var migration *core.Migration
	for _, item := range core.AppMigrations.Items() {
		if item.File == computersMigration {
			migration = item
		}
	}
	if migration == nil {
		t.Fatal("the computers migration is not registered")
	}
	if err := migration.Down(app); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"computers", "computer_logins", "computer_raw_logins"} {
		if _, err := app.FindCollectionByNameOrId(name); err == nil {
			t.Errorf("%s still exists after the down migration", name)
		}
	}
	if migrationCollection(t, app, "admin_profiles").Fields.GetByName("user_source_anchor") != nil {
		t.Error("admin_profiles.user_source_anchor still exists after the down migration")
	}
	role := migrationCollection(t, app, "machine_secrets").Fields.GetByName("role").(*core.SelectField)
	if slices.Contains(role.Values, "radiator") {
		t.Error("the radiator role still exists after the down migration")
	}
	if _, err := app.FindFirstRecordByData("claims", "name", "computers"); err == nil {
		t.Error("the computers claim still exists after the down migration")
	}
	if err := migration.Up(app); err != nil {
		t.Fatal(err)
	}
}
