// Package migrations registers schema changes for PocketBase collections.
//
// This migration exists to tighten the initial work-record schema after the
// first collection creation landed. Two follow-up adjustments are bundled here:
// the work record number field is made optional so backend numbering hooks can
// assign it during save, and the transient "selected" flag is removed from work
// record consumable entries because it represented UI state rather than durable
// domain data.
//
// The migration updates the existing collections in place rather than replacing
// them wholesale. That keeps the change small, preserves collection IDs that
// other code already references, and provides a clean rollback path that
// restores the previous number-field requirement and selected flag if the
// migration is reversed.
package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		workRecordsCollection, err := app.FindCollectionByNameOrId("pbc_1776282103")
		if err != nil {
			return err
		}

		workRecordsCollection.Fields.RemoveById("text1776282103")
		if err := workRecordsCollection.Fields.AddMarshaledJSONAt(1, []byte(`{
			"autogeneratePattern": "",
			"hidden": false,
			"id": "text1776282103",
			"max": 0,
			"min": 0,
			"name": "number",
			"pattern": "^W[0-9]{4}-[0-9]{4}$",
			"presentable": true,
			"primaryKey": false,
			"required": false,
			"system": false,
			"type": "text"
		}`)); err != nil {
			return err
		}
		if err := app.Save(workRecordsCollection); err != nil {
			return err
		}

		consumableEntriesCollection, err := app.FindCollectionByNameOrId("pbc_1776282106")
		if err != nil {
			return err
		}

		consumableEntriesCollection.Fields.RemoveById("bool1776282106")

		return app.Save(consumableEntriesCollection)
	}, func(app core.App) error {
		workRecordsCollection, err := app.FindCollectionByNameOrId("pbc_1776282103")
		if err != nil {
			return err
		}

		workRecordsCollection.Fields.RemoveById("text1776282103")
		if err := workRecordsCollection.Fields.AddMarshaledJSONAt(1, []byte(`{
			"autogeneratePattern": "",
			"hidden": false,
			"id": "text1776282103",
			"max": 0,
			"min": 0,
			"name": "number",
			"pattern": "^W[0-9]{4}-[0-9]{4}$",
			"presentable": true,
			"primaryKey": false,
			"required": true,
			"system": false,
			"type": "text"
		}`)); err != nil {
			return err
		}
		if err := app.Save(workRecordsCollection); err != nil {
			return err
		}

		consumableEntriesCollection, err := app.FindCollectionByNameOrId("pbc_1776282106")
		if err != nil {
			return err
		}

		if err := consumableEntriesCollection.Fields.AddMarshaledJSONAt(4, []byte(`{
			"hidden": false,
			"id": "bool1776282106",
			"name": "selected",
			"presentable": false,
			"required": false,
			"system": false,
			"type": "bool"
		}`)); err != nil {
			return err
		}

		return app.Save(consumableEntriesCollection)
	})
}
