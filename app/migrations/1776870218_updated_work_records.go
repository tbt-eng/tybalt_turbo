package migrations

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_1776282103")
		if err != nil {
			return err
		}

		// Add the paper-copy attachment fields and use a partial unique index for
		// attachment_hash so multiple rows can keep the default empty hash while
		// real uploaded files remain deduplicated across the collection.
		//
		// Intentionally do not add a collection-level updateRule here. Attachment
		// upload/remove is handled by a dedicated custom route under
		// /api/work_records/{id}/attachment so normal users do not gain a broad
		// direct PATCH path for arbitrary work_records fields.
		if err := json.Unmarshal([]byte(`{
			"indexes": [
				"CREATE UNIQUE INDEX idx_work_records_number ON work_records (number)",
				"CREATE INDEX idx_work_records_uid_date ON work_records (uid, date)",
				"CREATE INDEX idx_work_records_job_date ON work_records (job, date)",
				"CREATE UNIQUE INDEX idx_work_records_uid_job_date ON work_records (uid, job, date)",
				"CREATE UNIQUE INDEX `+"`"+`idx_zdRTW9LEgq`+"`"+` ON `+"`"+`work_records`+"`"+` (`+"`"+`attachment_hash`+"`"+`) WHERE `+"`"+`attachment_hash`+"`"+` != ''"
			]
		}`), &collection); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(23, []byte(`{
			"hidden": false,
			"id": "file2036324795",
			"maxSelect": 1,
			"maxSize": 10485760,
			"mimeTypes": [
				"application/pdf"
			],
			"name": "attachment",
			"presentable": false,
			"protected": false,
			"required": false,
			"system": false,
			"thumbs": [],
			"type": "file"
		}`)); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(24, []byte(`{
			"autogeneratePattern": "",
			"hidden": false,
			"id": "text701200007",
			"max": 0,
			"min": 0,
			"name": "attachment_hash",
			"pattern": "",
			"presentable": false,
			"primaryKey": false,
			"required": false,
			"system": false,
			"type": "text"
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_1776282103")
		if err != nil {
			return err
		}

		// Roll back the attachment dedupe index change so the collection returns
		// to its pre-attachment behavior. The migration never opened a collection
		// updateRule, so there is no direct-write rule state to unwind here.
		if err := json.Unmarshal([]byte(`{
			"indexes": [
				"CREATE UNIQUE INDEX idx_work_records_number ON work_records (number)",
				"CREATE INDEX idx_work_records_uid_date ON work_records (uid, date)",
				"CREATE INDEX idx_work_records_job_date ON work_records (job, date)",
				"CREATE UNIQUE INDEX idx_work_records_uid_job_date ON work_records (uid, job, date)"
			]
		}`), &collection); err != nil {
			return err
		}

		// remove field
		collection.Fields.RemoveById("file2036324795")

		// remove field
		collection.Fields.RemoveById("text701200007")

		return app.Save(collection)
	})
}
