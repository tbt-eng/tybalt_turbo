package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

const (
	workRecordsConsumablesCollectionName     = "work_records_consumables"
	workRecordsTypeConsumablesCollectionName = "work_records_type_consumables"
	workRecordTypeUnitLabelFieldID           = "text1776282113"
	workRecordConsumableUnitLabelFieldID     = "text1776719879"
)

func init() {
	m.Register(func(app core.App) error {
		typeConsumablesCollection, err := app.FindCollectionByNameOrId(workRecordsTypeConsumablesCollectionName)
		if err != nil {
			return err
		}
		typeConsumablesCollection.Fields.RemoveById(workRecordTypeUnitLabelFieldID)
		if err := app.Save(typeConsumablesCollection); err != nil {
			return err
		}

		consumablesCollection, err := app.FindCollectionByNameOrId(workRecordsConsumablesCollectionName)
		if err != nil {
			return err
		}
		if err := consumablesCollection.Fields.AddMarshaledJSONAt(3, []byte(`{
			"autogeneratePattern": "",
			"hidden": false,
			"id": "`+workRecordConsumableUnitLabelFieldID+`",
			"max": 0,
			"min": 0,
			"name": "unit_label",
			"pattern": "",
			"presentable": false,
			"primaryKey": false,
			"required": false,
			"system": false,
			"type": "text"
		}`)); err != nil {
			return err
		}

		return app.Save(consumablesCollection)
	}, func(app core.App) error {
		consumablesCollection, err := app.FindCollectionByNameOrId(workRecordsConsumablesCollectionName)
		if err != nil {
			return err
		}
		consumablesCollection.Fields.RemoveById(workRecordConsumableUnitLabelFieldID)
		if err := app.Save(consumablesCollection); err != nil {
			return err
		}

		typeConsumablesCollection, err := app.FindCollectionByNameOrId(workRecordsTypeConsumablesCollectionName)
		if err != nil {
			return err
		}
		if err := typeConsumablesCollection.Fields.AddMarshaledJSONAt(3, []byte(`{
			"autogeneratePattern": "",
			"hidden": false,
			"id": "`+workRecordTypeUnitLabelFieldID+`",
			"max": 0,
			"min": 0,
			"name": "unit_label",
			"pattern": "",
			"presentable": false,
			"primaryKey": false,
			"required": false,
			"system": false,
			"type": "text"
		}`)); err != nil {
			return err
		}

		return app.Save(typeConsumablesCollection)
	})
}
