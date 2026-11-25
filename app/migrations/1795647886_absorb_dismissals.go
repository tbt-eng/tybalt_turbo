package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

const absorbDismissalsCollectionID = "pbc_absorb_dism"

// Add absorb_dismissals.
//
// On the absorb suggestions page an absorb holder marks a bucket of probable
// duplicates "Not duplicates". The row keeps the records of the bucket at
// that moment, and it is shared by every user. A bucket is hidden while one
// row holds all of its records, so it shows again when a new record joins
// it. The suggestion routes read and write the rows, so the collection has no
// API rules. See descriptions/Absorb.md.
func init() {
	m.Register(func(app core.App) error {
		users, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}

		dismissals := core.NewBaseCollection("absorb_dismissals")
		dismissals.Id = absorbDismissalsCollectionID
		dismissals.Fields.Add(
			&core.SelectField{Name: "collection_name", Values: []string{"clients", "vendors"}, MaxSelect: 1, Required: true},
			// The ids of the records marked as not duplicates, as a JSON array.
			&core.JSONField{Name: "record_ids", Required: true, MaxSize: 100000},
			&core.RelationField{Name: "dismissed_by", CollectionId: users.Id, MaxSelect: 1},
			&core.AutodateField{Name: "created", OnCreate: true},
			&core.AutodateField{Name: "updated", OnCreate: true, OnUpdate: true},
		)
		dismissals.AddIndex("idx_absorb_dismissals_collection", false, "`collection_name`", "")
		return app.Save(dismissals)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId(absorbDismissalsCollectionID)
		if err != nil {
			return err
		}
		return app.Delete(collection)
	})
}
