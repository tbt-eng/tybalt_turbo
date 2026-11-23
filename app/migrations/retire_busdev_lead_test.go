package migrations_test

import (
	"testing"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

const retireBusdevLeadMigration = "1795463617_retire_busdev_lead.go"

// The down migration restores the required flag and the busdev claim. The up
// migration returns the clients schema to the canonical one and deletes the
// claim along with the records that give it to a user.
func TestRetireBusdevLeadMigrationRevert(t *testing.T) {
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
		if item.File == retireBusdevLeadMigration {
			migration = item
		}
	}
	if migration == nil {
		t.Fatal("retire busdev lead migration must be registered")
	}

	required := func(want bool) {
		t.Helper()
		clients := migrationCollection(t, app, "clients")
		if got := clients.Fields.GetByName("business_development_lead").(*core.RelationField).Required; got != want {
			t.Fatalf("clients.business_development_lead required = %v, want %v", got, want)
		}
	}
	claimExists := func(want bool) {
		t.Helper()
		_, err := app.FindFirstRecordByData("claims", "name", "busdev")
		if got := err == nil; got != want {
			t.Fatalf("busdev claim exists = %v, want %v", got, want)
		}
	}

	required(false)
	claimExists(false)

	if err := migration.Down(app); err != nil {
		t.Fatal(err)
	}
	required(true)
	claimExists(true)

	// This test runs on an empty migrated database, which the CSV fixtures do
	// not fill, so it creates one user who holds the restored claim here. The
	// up migration must take the claim from that user before it deletes it.
	claim, err := app.FindFirstRecordByData("claims", "name", "busdev")
	if err != nil {
		t.Fatal(err)
	}
	user := core.NewRecord(migrationCollection(t, app, "users"))
	user.Set("email", "holder@example.com")
	user.Set("username", "busdev_holder")
	user.SetPassword("a-long-test-password")
	if err := app.Save(user); err != nil {
		t.Fatal(err)
	}
	holder := core.NewRecord(migrationCollection(t, app, "user_claims"))
	holder.Set("uid", user.Id)
	holder.Set("cid", claim.Id)
	if err := app.Save(holder); err != nil {
		t.Fatal(err)
	}

	if err := migration.Up(app); err != nil {
		t.Fatal(err)
	}
	required(false)
	claimExists(false)
	if count, err := app.CountRecords("user_claims", dbx.HashExp{"uid": user.Id}); err != nil || count != 0 {
		t.Fatalf("user_claims of the holder = %d (err %v), want 0", count, err)
	}
	assertCanonicalCollections(t, app, "clients")
}
