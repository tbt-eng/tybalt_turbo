package migrations_test

import (
	"testing"

	"tybalt/internal/testseed"

	"github.com/pocketbase/pocketbase/core"
)

const revokeInactiveUserTokensMigration = "1795814551_revoke_inactive_user_tokens.go"

// The migration gives a new token key to every user whose admin profile is
// inactive and to every user with no admin profile, and leaves active users
// alone. It runs here against the seeded fixtures, which hold users of all
// three kinds; the seeded template ran it before the fixtures were loaded.
func TestRevokeInactiveUserTokensMigration(t *testing.T) {
	app := testseed.NewSeededTestApp(t)

	var migration *core.Migration
	for _, item := range core.AppMigrations.Items() {
		if item.File == revokeInactiveUserTokensMigration {
			migration = item
		}
	}
	if migration == nil {
		t.Fatal("revoke inactive user tokens migration must be registered")
	}

	var rows []struct {
		ID       string `db:"id"`
		TokenKey string `db:"tokenKey"`
		Kind     string `db:"kind"`
	}
	if err := app.DB().NewQuery(`
		SELECT u.id, u.tokenKey,
		       CASE WHEN ap.id IS NULL THEN 'none'
		            WHEN ap.active = 1 THEN 'active'
		            ELSE 'inactive' END AS kind
		FROM users u
		LEFT JOIN admin_profiles ap ON ap.uid = u.id
	`).All(&rows); err != nil {
		t.Fatal(err)
	}
	kinds := map[string]string{}
	for _, row := range rows {
		kinds[row.ID] = row.Kind
	}
	// u_inactive has an inactive profile, hruser000000001 has no profile and
	// f2j5a8vk006baub (author@soup.com) is active.
	for id, want := range map[string]string{"u_inactive": "inactive", "hruser000000001": "none", "f2j5a8vk006baub": "active"} {
		if kinds[id] != want {
			t.Fatalf("fixture user %s is %q, want %q", id, kinds[id], want)
		}
	}

	tokens := map[string]string{}
	for _, id := range []string{"u_inactive", "hruser000000001", "f2j5a8vk006baub"} {
		user, err := app.FindRecordById("users", id)
		if err != nil {
			t.Fatal(err)
		}
		token, err := user.NewAuthToken()
		if err != nil {
			t.Fatal(err)
		}
		tokens[id] = token
	}

	if err := migration.Up(app); err != nil {
		t.Fatal(err)
	}

	for _, row := range rows {
		user, err := app.FindRecordById("users", row.ID)
		if err != nil {
			t.Fatal(err)
		}
		changed := user.TokenKey() != row.TokenKey
		if want := row.Kind != "active"; changed != want {
			t.Errorf("user %s (%s profile): token key changed = %v, want %v", row.ID, row.Kind, changed, want)
		}
	}

	for id, token := range tokens {
		_, err := app.FindAuthRecordByToken(token, core.TokenTypeAuth)
		if valid, want := err == nil, kinds[id] == "active"; valid != want {
			t.Errorf("token of %s (%s profile) valid after the migration = %v, want %v", id, kinds[id], valid, want)
		}
	}
}
