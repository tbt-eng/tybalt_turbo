package migrations_test

import (
	"encoding/csv"
	"os"
	"testing"

	"github.com/pocketbase/pocketbase/core"
)

var expenseUncommitEmailCodes = []string{"expense_uncommit_approved", "expense_uncommit_declined", "expense_uncommit_reminder"}

func TestExpenseUncommitMigrationSeedsReasonsWithoutCumulative(t *testing.T) {
	app := newMigratedApp(t)
	reasons, err := app.FindRecordsByFilter("expense_uncommit_reasons", "active = true", "sort", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(reasons) != 12 || reasons[len(reasons)-1].GetString("name") != "Other" {
		t.Fatalf("unexpected reasons: %d", len(reasons))
	}
}

// A rejected expense has to stay deletable, so the relation can't be required.
func TestExpenseUncommitRequestExpenseRelationIsOptional(t *testing.T) {
	app := newMigratedApp(t)
	requests, err := app.FindCollectionByNameOrId("expense_uncommit_requests")
	if err != nil {
		t.Fatal(err)
	}
	field, ok := requests.Fields.GetByName("expense").(*core.RelationField)
	if !ok || field.Required || field.CascadeDelete {
		t.Fatalf("expected an optional relation that doesn't cascade, got %+v", field)
	}
	if requests.ListRule != nil || requests.ViewRule != nil || requests.CreateRule != nil || requests.UpdateRule != nil || requests.DeleteRule != nil {
		t.Fatal("expected the collection to have no API rules")
	}
}

func TestExpenseUncommitMigrationRollsBack(t *testing.T) {
	app := newMigratedApp(t)
	for _, item := range core.AppMigrations.Items() {
		if item.File == "1793893333_expense_uncommit_requests.go" {
			if err := item.Down(app); err != nil {
				t.Fatal(err)
			}
			if _, err := app.FindCollectionByNameOrId("expense_uncommit_requests"); err == nil {
				t.Fatal("expected the requests collection to be removed")
			}
			if _, err := app.FindCollectionByNameOrId("expense_uncommit_reasons"); err == nil {
				t.Fatal("expected the reasons collection to be removed")
			}
			assertWorkflowEmails(t, app, nil, expenseUncommitEmailCodes, nil)
			return
		}
	}
	t.Fatal("expense uncommit requests migration is not registered")
}

// The seeded test data must send the same emails as a migrated database.
func TestExpenseUncommitEmailFixturesMatchTheMigration(t *testing.T) {
	app := newMigratedApp(t)
	file, err := os.Open("../test_seed_data/data/notification_templates.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	rows, err := csv.NewReader(file).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	columns := map[string]int{}
	for i, name := range rows[0] {
		columns[name] = i
	}
	for _, code := range expenseUncommitEmailCodes {
		template, err := app.FindFirstRecordByData("notification_templates", "code", code)
		if err != nil {
			t.Fatal(err)
		}
		found := false
		for _, row := range rows[1:] {
			if row[columns["code"]] != code {
				continue
			}
			found = true
			if row[columns["id"]] != template.Id || row[columns["subject"]] != template.GetString("subject") || row[columns["text_email"]] != template.GetString("text_email") {
				t.Fatalf("fixture for %s doesn't match the migration", code)
			}
		}
		if !found {
			t.Fatalf("no fixture for %s", code)
		}
	}
}
