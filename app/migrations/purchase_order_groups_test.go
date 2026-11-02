package migrations_test

import (
	"encoding/csv"
	"os"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/core"
)

const purchaseOrderGroupsMigration = "1793635548_purchase_order_groups.go"

var purchaseOrderGroupEmailCodes = []string{"po_group_approval_required", "po_group_second_approval_required"}

func purchaseOrderGroupsMigrationItem(t *testing.T) *core.Migration {
	t.Helper()
	for _, item := range core.AppMigrations.Items() {
		if item.File == purchaseOrderGroupsMigration {
			return item
		}
	}
	t.Fatal("purchase order groups migration is not registered")
	return nil
}

func ruleText(rule *string) string {
	if rule == nil {
		return "<nil>"
	}
	return *rule
}

// groupRules returns the rules the migration changes, keyed by collection and
// rule, and the profiles list rule, which it must leave alone.
func groupRules(t *testing.T, app core.App) map[string]string {
	t.Helper()
	profiles, err := app.FindCollectionByNameOrId("profiles")
	if err != nil {
		t.Fatal(err)
	}
	purchaseOrders, err := app.FindCollectionByNameOrId("purchase_orders")
	if err != nil {
		t.Fatal(err)
	}
	return map[string]string{
		"profiles list":        ruleText(profiles.ListRule),
		"purchase_orders list": ruleText(purchaseOrders.ListRule),
		"purchase_orders view": ruleText(purchaseOrders.ViewRule),
	}
}

// The seeded test data must send the same purchase order group emails as a
// migrated database.
func TestPurchaseOrderGroupEmailFixturesMatchTheMigration(t *testing.T) {
	app := newMigratedApp(t)
	assertWorkflowEmails(t, app, purchaseOrderGroupEmailCodes, nil, nil)
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
	for _, code := range purchaseOrderGroupEmailCodes {
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
			for _, field := range []string{"id", "description", "subject", "text_email", "html_email"} {
				if row[columns[field]] != template.GetString(field) {
					t.Errorf("%s fixture %s = %q; migration has %q", code, field, row[columns[field]], template.GetString(field))
				}
			}
		}
		if !found {
			t.Errorf("fixtures are missing the %s template", code)
		}
	}
}

// The seeded collections must have the rules a migrated database has.
func TestPurchaseOrderGroupRuleFixturesMatchTheMigration(t *testing.T) {
	app := newMigratedApp(t)
	migrated := groupRules(t, app)
	file, err := os.Open("../test_seed_data/data/_collections.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.LazyQuotes = true
	rows, err := reader.ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	columns := map[string]int{}
	for i, name := range rows[0] {
		columns[name] = i
	}
	fixture := map[string]string{}
	for _, row := range rows[1:] {
		switch row[columns["name"]] {
		case "profiles":
			fixture["profiles list"] = row[columns["listRule"]]
		case "purchase_orders":
			fixture["purchase_orders list"] = row[columns["listRule"]]
			fixture["purchase_orders view"] = row[columns["viewRule"]]
		case "purchase_order_groups":
			if row[columns["listRule"]] != `\N` || row[columns["viewRule"]] != `\N` {
				t.Error("purchase_order_groups fixture has list or view rules")
			}
		}
	}
	for key, rule := range migrated {
		if fixture[key] != rule {
			t.Errorf("%s fixture rule = %q; migration has %q", key, fixture[key], rule)
		}
	}
}

// Undoing the migration removes the groups collection, the po_group field,
// the claim and its holders, the rule clauses and the emails, and leaves the
// other notification flags and the original rules as they were.
func TestPurchaseOrderGroupsMigrationRoundTrip(t *testing.T) {
	app := newMigratedApp(t)
	migration := purchaseOrderGroupsMigrationItem(t)

	// Someone holds the claim, which must not stop the rollback.
	users, err := app.FindCollectionByNameOrId("users")
	if err != nil {
		t.Fatal(err)
	}
	user := core.NewRecord(users)
	user.SetEmail("holder@example.com")
	user.SetPassword("holder-password")
	if err := app.Save(user); err != nil {
		t.Fatal(err)
	}
	userClaims, err := app.FindCollectionByNameOrId("user_claims")
	if err != nil {
		t.Fatal(err)
	}
	holder := core.NewRecord(userClaims)
	holder.Set("uid", user.Id)
	holder.Set("cid", "superpo00000001")
	if err := app.Save(holder); err != nil {
		t.Fatal(err)
	}

	before := notificationFlags(t, app)
	widened := groupRules(t, app)
	if err := migration.Down(app); err != nil {
		t.Fatal(err)
	}
	assertWorkflowEmails(t, app, nil, purchaseOrderGroupEmailCodes, nil)
	if _, err := app.FindCollectionByNameOrId("purchase_order_groups"); err == nil {
		t.Error("purchase_order_groups still exists")
	}
	purchaseOrders, err := app.FindCollectionByNameOrId("purchase_orders")
	if err != nil {
		t.Fatal(err)
	}
	if purchaseOrders.Fields.GetByName("po_group") != nil {
		t.Error("purchase_orders.po_group still exists")
	}
	if _, err := app.FindFirstRecordByData("claims", "name", "super_po"); err == nil {
		t.Error("the super_po claim still exists")
	}
	if _, err := app.FindRecordById("user_claims", holder.Id); err == nil {
		t.Error("the super_po holder's user_claims record still exists")
	}
	original := groupRules(t, app)
	for key, rule := range original {
		if strings.Contains(rule, "super_po") || strings.Contains(rule, "po_group") {
			t.Errorf("%s rule still has the group clause: %q", key, rule)
		}
	}

	if err := migration.Up(app); err != nil {
		t.Fatal(err)
	}
	assertWorkflowEmails(t, app, purchaseOrderGroupEmailCodes, nil, nil)
	after := notificationFlags(t, app)
	for code, value := range before {
		if after[code] != value {
			t.Errorf("notification flag %s changed: %v -> %v", code, value, after[code])
		}
	}
	for key, rule := range groupRules(t, app) {
		if rule != widened[key] {
			t.Errorf("%s rule after redo = %q; want %q", key, rule, widened[key])
		}
	}

	if err := migration.Down(app); err != nil {
		t.Fatal(err)
	}
	for key, rule := range groupRules(t, app) {
		if rule != original[key] {
			t.Errorf("%s rule after second undo = %q; want %q", key, rule, original[key])
		}
	}
}

// The migration must not depend on how the existing rules are written. A
// production rule can differ from the one in the snapshot, and a failed
// migration stops the app from starting.
func TestPurchaseOrderGroupsMigrationWidensAnyRule(t *testing.T) {
	cases := map[string]*string{
		"one line":              ptr(`@request.auth.id != ""`),
		"ends with a comment":   ptr("@request.auth.id != \"\" // signed in"),
		"starts with a comment": ptr("// signed in\n@request.auth.id != \"\""),
		"superusers only":       nil,
		"public":                ptr(""),
	}
	for name, rule := range cases {
		t.Run(name, func(t *testing.T) {
			app := newMigratedApp(t)
			migration := purchaseOrderGroupsMigrationItem(t)
			if err := migration.Down(app); err != nil {
				t.Fatal(err)
			}
			collection, err := app.FindCollectionByNameOrId("purchase_orders")
			if err != nil {
				t.Fatal(err)
			}
			collection.ListRule = rule
			collection.ViewRule = rule
			if err := app.Save(collection); err != nil {
				t.Fatal(err)
			}
			original := groupRules(t, app)

			if err := migration.Up(app); err != nil {
				t.Fatal(err)
			}
			widened := groupRules(t, app)
			if rule != nil && *rule == "" {
				if widened["purchase_orders list"] != "" {
					t.Errorf("a public rule was changed to %q", widened["purchase_orders list"])
				}
			} else if !strings.Contains(widened["purchase_orders list"], "po_group") {
				t.Errorf("purchase_orders list rule was not widened: %q", widened["purchase_orders list"])
			}

			if err := migration.Down(app); err != nil {
				t.Fatal(err)
			}
			for key, restored := range groupRules(t, app) {
				if restored != original[key] {
					t.Errorf("%s rule after undo = %q; want %q", key, restored, original[key])
				}
			}
		})
	}
}

func ptr(value string) *string {
	return &value
}
