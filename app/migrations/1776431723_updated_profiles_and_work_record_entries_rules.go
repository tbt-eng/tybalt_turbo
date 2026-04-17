package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/types"
)

const (
	profilesDefaultWorkRecordTypeFieldID  = "eycTZv1vG5C5MqH"
	workRecordConsumableEntriesCreateRule = "@request.auth.id != \"\" &&\n@request.body.work_record.approved = false &&\n(\n  @request.body.work_record.uid = @request.auth.id ||\n  @request.body.work_record.creator = @request.auth.id\n)"
	workRecordConsumableEntriesUpdateRule = "@request.auth.id != \"\" &&\nwork_record.approved = false &&\n(\n  work_record.uid = @request.auth.id ||\n  work_record.creator = @request.auth.id\n)"
)

func init() {
	m.Register(func(app core.App) error {
		profilesCollection, err := app.FindCollectionByNameOrId("profiles")
		if err != nil {
			return err
		}

		if profilesCollection.Fields.GetById(profilesDefaultWorkRecordTypeFieldID) == nil {
			if err := profilesCollection.Fields.AddMarshaledJSONAt(5, []byte(`{
				"cascadeDelete": false,
				"collectionId": "pbc_1776282101",
				"hidden": false,
				"id": "`+profilesDefaultWorkRecordTypeFieldID+`",
				"maxSelect": 1,
				"minSelect": 0,
				"name": "default_work_record_type",
				"presentable": false,
				"required": false,
				"system": false,
				"type": "relation"
			}`)); err != nil {
				return err
			}
		}

		if err := app.Save(profilesCollection); err != nil {
			return err
		}

		entriesCollection, err := app.FindCollectionByNameOrId("work_records_consumable_entries")
		if err != nil {
			return err
		}

		entriesCollection.CreateRule = types.Pointer(workRecordConsumableEntriesCreateRule)
		entriesCollection.UpdateRule = types.Pointer(workRecordConsumableEntriesUpdateRule)
		entriesCollection.DeleteRule = types.Pointer(workRecordConsumableEntriesUpdateRule)

		return app.Save(entriesCollection)
	}, func(app core.App) error {
		profilesCollection, err := app.FindCollectionByNameOrId("profiles")
		if err != nil {
			return err
		}
		profilesCollection.Fields.RemoveById(profilesDefaultWorkRecordTypeFieldID)
		if err := app.Save(profilesCollection); err != nil {
			return err
		}

		entriesCollection, err := app.FindCollectionByNameOrId("work_records_consumable_entries")
		if err != nil {
			return err
		}
		entriesCollection.CreateRule = nil
		entriesCollection.UpdateRule = nil
		entriesCollection.DeleteRule = nil

		return app.Save(entriesCollection)
	})
}
