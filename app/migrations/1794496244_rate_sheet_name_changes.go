package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/types"
)

const rateSheetNameChangesCollectionID = "pbc_rs_renames"

// Add rate_sheet_name_changes, the history of rate sheet names.
//
// An admin renames a sheet through /api/rate_sheets/{id}/rename, which changes
// name on every revision and writes one row here in the same transaction.
// rate_sheet is any one revision of the sheet. It is only an anchor: a sheet's
// earlier names are the old_name values of the rows whose anchor has the
// sheet's current name. That works for every revision, including revisions
// made after the rename, and no row changes on a later rename.
//
// Every signed-in user can list and view the rows. Only server code writes
// them: the collection has no create, update or delete rule.
//
// The anchor relation does not cascade. PocketBase refuses to delete a
// revision while a history row points at it, so the history can't vanish
// with one revision. The same holds for renamed_by and the user.
func init() {
	m.Register(func(app core.App) error {
		rateSheets, err := app.FindCollectionByNameOrId("rate_sheets")
		if err != nil {
			return err
		}
		users, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}

		changes := core.NewBaseCollection("rate_sheet_name_changes")
		changes.Id = rateSheetNameChangesCollectionID
		changes.ListRule = types.Pointer("@request.auth.id != ''")
		changes.ViewRule = types.Pointer("@request.auth.id != ''")
		changes.Fields.Add(
			&core.RelationField{Name: "rate_sheet", CollectionId: rateSheets.Id, MaxSelect: 1, Required: true},
			&core.TextField{Name: "old_name", Required: true},
			&core.TextField{Name: "new_name", Required: true},
			&core.RelationField{Name: "renamed_by", CollectionId: users.Id, MaxSelect: 1, Required: true},
			&core.AutodateField{Name: "created", OnCreate: true},
		)
		changes.AddIndex("idx_rate_sheet_name_changes_rate_sheet", false, "`rate_sheet`", "")
		return app.Save(changes)
	}, func(app core.App) error {
		changes, err := app.FindCollectionByNameOrId(rateSheetNameChangesCollectionID)
		if err != nil {
			return nil
		}
		return app.Delete(changes)
	})
}
