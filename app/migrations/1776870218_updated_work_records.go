package migrations

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/types"
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
		for _, fieldID := range []string{
			"relation1776282104", // uid moved to work_records_subjects
			"number1776282103",   // hours_on_site moved to work_records_subjects
			"number1776282104",   // hours_travel_time moved to work_records_subjects
			"number1776282105",   // distance_travelled_km moved to work_records_subjects
			"select1776282103",   // vehicle_type moved to work_records_subjects
			"bool1776282103",     // is_passenger moved to work_records_subjects
			"text1776282112",     // company_vehicle_unit_number moved to work_records_subjects
			"bool1776282104",     // approval is now stored per subject
		} {
			collection.Fields.RemoveById(fieldID)
		}

		if err := json.Unmarshal([]byte(`{
			"indexes": [
				"CREATE UNIQUE INDEX idx_work_records_number ON work_records (number)",
				"CREATE INDEX idx_work_records_job_date ON work_records (job, date)",
				"CREATE UNIQUE INDEX `+"`"+`idx_zdRTW9LEgq`+"`"+` ON `+"`"+`work_records`+"`"+` (`+"`"+`attachment_hash`+"`"+`) WHERE `+"`"+`attachment_hash`+"`"+` != ''"
			]
		}`), &collection); err != nil {
			return err
		}
		collection.ListRule = types.Pointer("@request.auth.id != \"\" &&\n(\n  creator = @request.auth.id ||\n  work_records_subjects_via_work_record.uid ?= @request.auth.id ||\n  @request.auth.user_claims_via_uid.cid.name ?= 'report'\n)")
		collection.ViewRule = collection.ListRule

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(15, []byte(`{
			"hidden": false,
			"id": "number1776870218",
			"max": null,
			"min": null,
			"name": "lane_kms",
			"onlyInt": false,
			"presentable": false,
			"required": false,
			"system": false,
			"type": "number"
		}`)); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(16, []byte(`{
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
		if err := collection.Fields.AddMarshaledJSONAt(17, []byte(`{
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
		collection.ListRule = types.Pointer("@request.auth.id != \"\" &&\n(\n  uid = @request.auth.id ||\n  creator = @request.auth.id ||\n  @request.auth.user_claims_via_uid.cid.name ?= 'report'\n)")
		collection.ViewRule = collection.ListRule

		// remove field
		collection.Fields.RemoveById("number1776870218")

		// remove field
		collection.Fields.RemoveById("file2036324795")

		// remove field
		collection.Fields.RemoveById("text701200007")

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(3, []byte(`{
			"cascadeDelete": false,
			"collectionId": "_pb_users_auth_",
			"hidden": false,
			"id": "relation1776282104",
			"maxSelect": 1,
			"minSelect": 0,
			"name": "uid",
			"presentable": false,
			"required": true,
			"system": false,
			"type": "relation"
		}`)); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(10, []byte(`{
			"hidden": false,
			"id": "number1776282103",
			"max": null,
			"min": null,
			"name": "hours_on_site",
			"onlyInt": false,
			"presentable": false,
			"required": true,
			"system": false,
			"type": "number"
		}`)); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(11, []byte(`{
			"hidden": false,
			"id": "number1776282104",
			"max": null,
			"min": null,
			"name": "hours_travel_time",
			"onlyInt": false,
			"presentable": false,
			"required": true,
			"system": false,
			"type": "number"
		}`)); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(12, []byte(`{
			"hidden": false,
			"id": "number1776282105",
			"max": null,
			"min": null,
			"name": "distance_travelled_km",
			"onlyInt": false,
			"presentable": false,
			"required": true,
			"system": false,
			"type": "number"
		}`)); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(18, []byte(`{
			"hidden": false,
			"id": "select1776282103",
			"maxSelect": 1,
			"name": "vehicle_type",
			"presentable": false,
			"required": true,
			"system": false,
			"type": "select",
			"values": [
				"company",
				"personal"
			]
		}`)); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(19, []byte(`{
			"hidden": false,
			"id": "bool1776282103",
			"name": "is_passenger",
			"presentable": false,
			"required": false,
			"system": false,
			"type": "bool"
		}`)); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(20, []byte(`{
			"autogeneratePattern": "",
			"hidden": false,
			"id": "text1776282112",
			"max": 0,
			"min": 0,
			"name": "company_vehicle_unit_number",
			"pattern": "^(00[1-9]|0[1-9][0-9]|[1-9][0-9]{0,2})$",
			"presentable": false,
			"primaryKey": false,
			"required": false,
			"system": false,
			"type": "text"
		}`)); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(21, []byte(`{
			"hidden": false,
			"id": "bool1776282104",
			"name": "approved",
			"presentable": false,
			"required": false,
			"system": false,
			"type": "bool"
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	})
}
