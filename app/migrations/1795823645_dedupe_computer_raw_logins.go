package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// This migration keeps only the newest raw login of each computer name.
//
// Why this exists:
//   - The radiator endpoint (routes/radiator.go) now keeps one raw login per
//     computer name: storing a reject deletes the older ones with the same
//     name. Before, a computer rejected at every sign-in added one row per
//     sign-in, and the retention job deleted them after 30 days. The
//     retention job no longer deletes a raw login with a name, so the rows
//     stored before the change would stay for ever.
//
// What it deletes:
//   - For each non-blank computer_name, every raw login except the newest by
//     created, the tie broken by the larger id. The name matches exactly, as
//     the endpoint matches it.
//   - Raw logins with a blank computer_name are left alone. They can't be
//     tied to a computer; the retention job still deletes them after 30 days.
//
// It changes no schema. Running it again deletes nothing more.
//
// Rollback note: the down migration is a no-op. The deleted rows can't come
// back.
func init() {
	m.Register(func(app core.App) error {
		_, err := app.DB().NewQuery(`
			DELETE FROM computer_raw_logins
			WHERE computer_name != ''
			  AND EXISTS (
			    SELECT 1
			    FROM computer_raw_logins newer
			    WHERE newer.computer_name = computer_raw_logins.computer_name
			      AND (
			        newer.created > computer_raw_logins.created
			        OR (newer.created = computer_raw_logins.created AND newer.id > computer_raw_logins.id)
			      )
			  )
		`).Execute()
		return err
	}, func(app core.App) error {
		// Nothing to undo: the deleted raw logins can't come back.
		return nil
	})
}
