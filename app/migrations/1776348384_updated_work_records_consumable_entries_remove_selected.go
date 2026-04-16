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
