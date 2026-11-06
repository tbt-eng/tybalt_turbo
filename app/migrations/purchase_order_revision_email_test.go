package migrations_test

import (
	"encoding/csv"
	"os"
	"testing"

	"github.com/pocketbase/pocketbase/core"
)

const (
	purchaseOrderRevisionEmailMigration = "1793905177_purchase_order_revision_email.go"
	purchaseOrderRevisionEmailCode      = "po_revision_approval_required"
)

// The seeded test data must send the same revision email as a migrated database.
func TestPurchaseOrderRevisionEmailFixtureMatchesTheMigration(t *testing.T) {
	app := newMigratedApp(t)
	assertWorkflowEmails(t, app, []string{purchaseOrderRevisionEmailCode}, nil, nil)
	template, err := app.FindFirstRecordByData("notification_templates", "code", purchaseOrderRevisionEmailCode)
	if err != nil {
		t.Fatal(err)
	}
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
	found := false
	for _, row := range rows[1:] {
		if row[columns["code"]] != purchaseOrderRevisionEmailCode {
			continue
		}
		found = true
		for _, field := range []string{"id", "description", "subject", "text_email", "html_email"} {
			if row[columns[field]] != template.GetString(field) {
				t.Errorf("fixture %s = %q; migration has %q", field, row[columns[field]], template.GetString(field))
			}
		}
	}
	if !found {
		t.Errorf("fixtures are missing the %s template", purchaseOrderRevisionEmailCode)
	}
}

func TestPurchaseOrderRevisionEmailMigrationRevert(t *testing.T) {
	app := newMigratedApp(t)
	var migration *core.Migration
	for _, item := range core.AppMigrations.Items() {
		if item.File == purchaseOrderRevisionEmailMigration {
			migration = item
		}
	}
	if migration == nil {
		t.Fatal("purchase order revision email migration must be registered")
	}
	if err := migration.Down(app); err != nil {
		t.Fatal(err)
	}
	assertWorkflowEmails(t, app, nil, []string{purchaseOrderRevisionEmailCode}, nil)
	if err := migration.Up(app); err != nil {
		t.Fatal(err)
	}
	assertWorkflowEmails(t, app, []string{purchaseOrderRevisionEmailCode}, nil, nil)
}
