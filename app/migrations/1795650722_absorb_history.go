package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

const absorbHistoryCollectionID = "pbc_absorb_hist"

// Keep a history of committed absorbs, and record who absorbed.
//
// A commit deletes the absorb action, which until now left no trace of the
// absorb. The commit now copies the action into absorb_history first: the
// collection, the kept record and its name, the saved copies of the absorbed
// records, the references that moved, the proposal, who absorbed and when,
// and who committed (created is the commit time). A daily cron job deletes
// rows older than ABSORB_HISTORY_RETENTION_DAYS. The history route reads the
// rows, so the collection has no API rules.
//
// absorb_actions gains absorbed_by, set when the absorb is made, so that the
// history and the pending notice can name the person. See
// descriptions/Absorb.md.
func init() {
	m.Register(func(app core.App) error {
		users, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}

		actions, err := app.FindCollectionByNameOrId("absorb_actions")
		if err != nil {
			return err
		}
		actions.Fields.Add(&core.RelationField{Name: "absorbed_by", CollectionId: users.Id, MaxSelect: 1})
		if err := app.Save(actions); err != nil {
			return err
		}

		history := core.NewBaseCollection("absorb_history")
		history.Id = absorbHistoryCollectionID
		history.Fields.Add(
			&core.TextField{Name: "collection_name", Required: true},
			&core.TextField{Name: "target_id", Required: true},
			// The name of the kept record at the commit.
			&core.TextField{Name: "target_name"},
			&core.JSONField{Name: "absorbed_records", Required: true, MaxSize: 2000000},
			&core.JSONField{Name: "updated_references", Required: true, MaxSize: 2000000},
			&core.TextField{Name: "proposal"},
			&core.RelationField{Name: "absorbed_by", CollectionId: users.Id, MaxSelect: 1},
			&core.DateField{Name: "absorbed_at"},
			&core.RelationField{Name: "committed_by", CollectionId: users.Id, MaxSelect: 1},
			&core.AutodateField{Name: "created", OnCreate: true},
			&core.AutodateField{Name: "updated", OnCreate: true, OnUpdate: true},
		)
		history.AddIndex("idx_absorb_history_created", false, "`created`", "")
		return app.Save(history)
	}, func(app core.App) error {
		history, err := app.FindCollectionByNameOrId(absorbHistoryCollectionID)
		if err != nil {
			return err
		}
		if err := app.Delete(history); err != nil {
			return err
		}
		actions, err := app.FindCollectionByNameOrId("absorb_actions")
		if err != nil {
			return err
		}
		actions.Fields.RemoveByName("absorbed_by")
		return app.Save(actions)
	})
}
