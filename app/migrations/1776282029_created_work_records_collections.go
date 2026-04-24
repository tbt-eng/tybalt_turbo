// Package migrations registers schema changes for PocketBase collections.
//
// This migration exists because work records were promoted from legacy text
// references embedded in time entries into a first-class data model with their
// own collections, related child tables, and access rules. The work_record_id
// work in this branch depends on these collections already existing: time
// entries can only link relationally once work records, work record types, and
// their supporting child collections have concrete schema definitions.
//
// The migration creates the initial work-record collection family from JSON
// snapshots rather than constructing every field procedurally. That keeps the
// PocketBase schema definition explicit, reproducible, and easy to compare with
// later collection snapshots. The helper at the top exists purely to decode and
// save those JSON definitions in one place so the init registration stays
// focused on the set of collections being installed.
package migrations

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func saveCollectionFromJSON(app core.App, jsonData string) error {
	collection := &core.Collection{}
	if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
		return err
	}

	return app.Save(collection)
}

func init() {
	m.Register(func(app core.App) error {
		workRecordsTypesJSON := `{
			"createRule": null,
			"deleteRule": null,
			"fields": [
				{
					"autogeneratePattern": "[a-z0-9]{15}",
					"hidden": false,
					"id": "text3208210256",
					"max": 15,
					"min": 15,
					"name": "id",
					"pattern": "^[a-z0-9]+$",
					"presentable": false,
					"primaryKey": true,
					"required": true,
					"system": true,
					"type": "text"
				},
				{
					"autogeneratePattern": "",
					"hidden": false,
					"id": "text1776282101",
					"max": 0,
					"min": 0,
					"name": "name",
					"pattern": "",
					"presentable": true,
					"primaryKey": false,
					"required": true,
					"system": false,
					"type": "text"
				},
				{
					"hidden": false,
					"id": "number1776282115",
					"max": null,
					"min": 1,
					"name": "max_subjects",
					"onlyInt": true,
					"presentable": false,
					"required": true,
					"system": false,
					"type": "number"
				},
				{
					"hidden": false,
					"id": "bool1776282115",
					"name": "allow_lane_kms",
					"presentable": false,
					"required": false,
					"system": false,
					"type": "bool"
				},
				{
					"hidden": false,
					"id": "bool1776282101",
					"name": "active",
					"presentable": false,
					"required": false,
					"system": false,
					"type": "bool"
				},
				{
					"hidden": false,
					"id": "number1776282101",
					"max": null,
					"min": null,
					"name": "sort_order",
					"onlyInt": true,
					"presentable": false,
					"required": false,
					"system": false,
					"type": "number"
				}
			],
			"id": "pbc_1776282101",
			"indexes": [
				"CREATE UNIQUE INDEX idx_work_records_types_name ON work_records_types (name)"
			],
			"listRule": "@request.auth.id != \"\"",
			"name": "work_records_types",
			"system": false,
			"type": "base",
			"updateRule": null,
			"viewRule": "@request.auth.id != \"\""
		}`

		workRecordsConsumablesJSON := `{
			"createRule": null,
			"deleteRule": null,
			"fields": [
				{
					"autogeneratePattern": "[a-z0-9]{15}",
					"hidden": false,
					"id": "text3208210256",
					"max": 15,
					"min": 15,
					"name": "id",
					"pattern": "^[a-z0-9]+$",
					"presentable": false,
					"primaryKey": true,
					"required": true,
					"system": true,
					"type": "text"
				},
				{
					"autogeneratePattern": "",
					"hidden": false,
					"id": "text1776282102",
					"max": 0,
					"min": 0,
					"name": "name",
					"pattern": "",
					"presentable": true,
					"primaryKey": false,
					"required": true,
					"system": false,
					"type": "text"
				},
				{
					"hidden": false,
					"id": "select1776282102",
					"maxSelect": 1,
					"name": "input_kind",
					"presentable": false,
					"required": true,
					"system": false,
					"type": "select",
					"values": [
						"number",
						"boolean"
					]
				},
				{
					"hidden": false,
					"id": "bool1776282102",
					"name": "active",
					"presentable": false,
					"required": false,
					"system": false,
					"type": "bool"
				}
			],
			"id": "pbc_1776282102",
			"indexes": [
				"CREATE UNIQUE INDEX idx_work_records_consumables_name ON work_records_consumables (name)"
			],
			"listRule": "@request.auth.id != \"\"",
			"name": "work_records_consumables",
			"system": false,
			"type": "base",
			"updateRule": null,
			"viewRule": "@request.auth.id != \"\""
		}`

		workRecordsJSON := `{
			"createRule": null,
			"deleteRule": null,
			"fields": [
				{
					"autogeneratePattern": "[a-z0-9]{15}",
					"hidden": false,
					"id": "text3208210256",
					"max": 15,
					"min": 15,
					"name": "id",
					"pattern": "^[a-z0-9]+$",
					"presentable": false,
					"primaryKey": true,
					"required": true,
					"system": true,
					"type": "text"
				},
				{
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
				},
				{
					"autogeneratePattern": "",
					"hidden": false,
					"id": "text1776282104",
					"max": 0,
					"min": 0,
					"name": "date",
					"pattern": "^\\d{4}-\\d{2}-\\d{2}$",
					"presentable": false,
					"primaryKey": false,
					"required": true,
					"system": false,
					"type": "text"
				},
				{
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
				},
				{
					"cascadeDelete": false,
					"collectionId": "_pb_users_auth_",
					"hidden": false,
					"id": "relation1776282105",
					"maxSelect": 1,
					"minSelect": 0,
					"name": "creator",
					"presentable": false,
					"required": true,
					"system": false,
					"type": "relation"
				},
				{
					"cascadeDelete": false,
					"collectionId": "yovqzrnnomp0lkx",
					"hidden": false,
					"id": "relation1776282106",
					"maxSelect": 1,
					"minSelect": 0,
					"name": "job",
					"presentable": false,
					"required": true,
					"system": false,
					"type": "relation"
				},
				{
					"cascadeDelete": false,
					"collectionId": "pbc_1776282101",
					"hidden": false,
					"id": "relation1776282107",
					"maxSelect": 1,
					"minSelect": 0,
					"name": "type",
					"presentable": false,
					"required": true,
					"system": false,
					"type": "relation"
				},
				{
					"autogeneratePattern": "",
					"hidden": false,
					"id": "text1776282105",
					"max": 0,
					"min": 0,
					"name": "location",
					"pattern": "",
					"presentable": false,
					"primaryKey": false,
					"required": true,
					"system": false,
					"type": "text"
				},
				{
					"autogeneratePattern": "",
					"hidden": false,
					"id": "text1776282106",
					"max": 0,
					"min": 0,
					"name": "work_description",
					"pattern": "",
					"presentable": false,
					"primaryKey": false,
					"required": true,
					"system": false,
					"type": "text"
				},
				{
					"autogeneratePattern": "",
					"hidden": false,
					"id": "text1776282107",
					"max": 0,
					"min": 0,
					"name": "report_to",
					"pattern": "",
					"presentable": false,
					"primaryKey": false,
					"required": false,
					"system": false,
					"type": "text"
				},
				{
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
				},
				{
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
				},
				{
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
				},
				{
					"autogeneratePattern": "",
					"hidden": false,
					"id": "text1776282108",
					"max": 0,
					"min": 0,
					"name": "field_book_number",
					"pattern": "^[0-9]{2}-(00[1-9]|0[1-9][0-9]|[1-9][0-9]{2})$",
					"presentable": false,
					"primaryKey": false,
					"required": false,
					"system": false,
					"type": "text"
				},
				{
					"hidden": false,
					"id": "number1776282106",
					"max": null,
					"min": 1,
					"name": "field_book_page_number",
					"onlyInt": true,
					"presentable": false,
					"required": false,
					"system": false,
					"type": "number"
				},
				{
					"autogeneratePattern": "",
					"hidden": false,
					"id": "text1776282109",
					"max": 0,
					"min": 0,
					"name": "sub_contractor",
					"pattern": "",
					"presentable": false,
					"primaryKey": false,
					"required": false,
					"system": false,
					"type": "text"
				},
				{
					"autogeneratePattern": "",
					"hidden": false,
					"id": "text1776282110",
					"max": 0,
					"min": 0,
					"name": "equipment",
					"pattern": "",
					"presentable": false,
					"primaryKey": false,
					"required": false,
					"system": false,
					"type": "text"
				},
				{
					"autogeneratePattern": "",
					"hidden": false,
					"id": "text1776282111",
					"max": 0,
					"min": 0,
					"name": "supplies",
					"pattern": "",
					"presentable": false,
					"primaryKey": false,
					"required": false,
					"system": false,
					"type": "text"
				},
				{
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
				},
				{
					"hidden": false,
					"id": "bool1776282103",
					"name": "is_passenger",
					"presentable": false,
					"required": false,
					"system": false,
					"type": "bool"
				},
				{
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
				},
				{
					"hidden": false,
					"id": "bool1776282104",
					"name": "approved",
					"presentable": false,
					"required": false,
					"system": false,
					"type": "bool"
				},
				{
					"hidden": false,
					"id": "autodate2990389176",
					"name": "created",
					"onCreate": true,
					"onUpdate": false,
					"presentable": false,
					"system": false,
					"type": "autodate"
				},
				{
					"hidden": false,
					"id": "autodate3332085495",
					"name": "updated",
					"onCreate": true,
					"onUpdate": true,
					"presentable": false,
					"system": false,
					"type": "autodate"
				}
			],
			"id": "pbc_1776282103",
			"indexes": [
				"CREATE UNIQUE INDEX idx_work_records_number ON work_records (number)",
				"CREATE INDEX idx_work_records_uid_date ON work_records (uid, date)",
				"CREATE INDEX idx_work_records_job_date ON work_records (job, date)",
				"CREATE UNIQUE INDEX idx_work_records_uid_job_date ON work_records (uid, job, date)"
			],
			"listRule": "@request.auth.id != \"\" &&\n(\n  uid = @request.auth.id ||\n  creator = @request.auth.id ||\n  @request.auth.user_claims_via_uid.cid.name ?= 'report'\n)",
			"name": "work_records",
			"system": false,
			"type": "base",
			"updateRule": null,
			"viewRule": "@request.auth.id != \"\" &&\n(\n  uid = @request.auth.id ||\n  creator = @request.auth.id ||\n  @request.auth.user_claims_via_uid.cid.name ?= 'report'\n)"
		}`

		workRecordsSubjectsJSON := `{
			"createRule": null,
			"deleteRule": null,
			"fields": [
				{
					"autogeneratePattern": "[a-z0-9]{15}",
					"hidden": false,
					"id": "text3208210256",
					"max": 15,
					"min": 15,
					"name": "id",
					"pattern": "^[a-z0-9]+$",
					"presentable": false,
					"primaryKey": true,
					"required": true,
					"system": true,
					"type": "text"
				},
				{
					"cascadeDelete": true,
					"collectionId": "pbc_1776282103",
					"hidden": false,
					"id": "relation1776282114",
					"maxSelect": 1,
					"minSelect": 0,
					"name": "work_record",
					"presentable": false,
					"required": true,
					"system": false,
					"type": "relation"
				},
				{
					"cascadeDelete": false,
					"collectionId": "_pb_users_auth_",
					"hidden": false,
					"id": "relation1776282115",
					"maxSelect": 1,
					"minSelect": 0,
					"name": "uid",
					"presentable": false,
					"required": true,
					"system": false,
					"type": "relation"
				},
				{
					"hidden": false,
					"id": "number1776282116",
					"max": null,
					"min": null,
					"name": "hours_on_site",
					"onlyInt": false,
					"presentable": false,
					"required": true,
					"system": false,
					"type": "number"
				},
				{
					"hidden": false,
					"id": "number1776282117",
					"max": null,
					"min": null,
					"name": "hours_travel_time",
					"onlyInt": false,
					"presentable": false,
					"required": true,
					"system": false,
					"type": "number"
				},
				{
					"hidden": false,
					"id": "number1776282118",
					"max": null,
					"min": null,
					"name": "distance_travelled_km",
					"onlyInt": false,
					"presentable": false,
					"required": false,
					"system": false,
					"type": "number"
				},
				{
					"hidden": false,
					"id": "select1776282114",
					"maxSelect": 1,
					"name": "vehicle_type",
					"presentable": false,
					"required": false,
					"system": false,
					"type": "select",
					"values": [
						"company",
						"personal"
					]
				},
				{
					"hidden": false,
					"id": "bool1776282116",
					"name": "is_passenger",
					"presentable": false,
					"required": false,
					"system": false,
					"type": "bool"
				},
				{
					"autogeneratePattern": "",
					"hidden": false,
					"id": "text1776282115",
					"max": 0,
					"min": 0,
					"name": "company_vehicle_unit_number",
					"pattern": "^(00[1-9]|0[1-9][0-9]|[1-9][0-9]{0,2})$",
					"presentable": false,
					"primaryKey": false,
					"required": false,
					"system": false,
					"type": "text"
				},
				{
					"hidden": false,
					"id": "bool1776282117",
					"name": "approved",
					"presentable": false,
					"required": false,
					"system": false,
					"type": "bool"
				},
				{
					"hidden": false,
					"id": "autodate2990389176",
					"name": "created",
					"onCreate": true,
					"onUpdate": false,
					"presentable": false,
					"system": false,
					"type": "autodate"
				},
				{
					"hidden": false,
					"id": "autodate3332085495",
					"name": "updated",
					"onCreate": true,
					"onUpdate": true,
					"presentable": false,
					"system": false,
					"type": "autodate"
				}
			],
			"id": "pbc_1776282107",
			"indexes": [
				"CREATE UNIQUE INDEX idx_work_records_subjects_work_record_uid ON work_records_subjects (work_record, uid)",
				"CREATE INDEX idx_work_records_subjects_uid_created ON work_records_subjects (uid, created)"
			],
			"listRule": "@request.auth.id != \"\" &&\n(\n  uid = @request.auth.id ||\n  work_record.creator = @request.auth.id ||\n  @request.auth.user_claims_via_uid.cid.name ?= 'report'\n)",
			"name": "work_records_subjects",
			"system": false,
			"type": "base",
			"updateRule": null,
			"viewRule": "@request.auth.id != \"\" &&\n(\n  uid = @request.auth.id ||\n  work_record.creator = @request.auth.id ||\n  @request.auth.user_claims_via_uid.cid.name ?= 'report'\n)"
		}`

		workRecordsTypeConsumablesJSON := `{
			"createRule": null,
			"deleteRule": null,
			"fields": [
				{
					"autogeneratePattern": "[a-z0-9]{15}",
					"hidden": false,
					"id": "text3208210256",
					"max": 15,
					"min": 15,
					"name": "id",
					"pattern": "^[a-z0-9]+$",
					"presentable": false,
					"primaryKey": true,
					"required": true,
					"system": true,
					"type": "text"
				},
				{
					"cascadeDelete": true,
					"collectionId": "pbc_1776282101",
					"hidden": false,
					"id": "relation1776282108",
					"maxSelect": 1,
					"minSelect": 0,
					"name": "type",
					"presentable": false,
					"required": true,
					"system": false,
					"type": "relation"
				},
				{
					"cascadeDelete": true,
					"collectionId": "pbc_1776282102",
					"hidden": false,
					"id": "relation1776282109",
					"maxSelect": 1,
					"minSelect": 0,
					"name": "consumable",
					"presentable": false,
					"required": true,
					"system": false,
					"type": "relation"
				},
				{
					"autogeneratePattern": "",
					"hidden": false,
					"id": "text1776282113",
					"max": 0,
					"min": 0,
					"name": "unit_label",
					"pattern": "",
					"presentable": false,
					"primaryKey": false,
					"required": false,
					"system": false,
					"type": "text"
				},
				{
					"hidden": false,
					"id": "number1776282107",
					"max": null,
					"min": null,
					"name": "sort_order",
					"onlyInt": true,
					"presentable": false,
					"required": false,
					"system": false,
					"type": "number"
				},
				{
					"hidden": false,
					"id": "bool1776282105",
					"name": "active",
					"presentable": false,
					"required": false,
					"system": false,
					"type": "bool"
				}
			],
			"id": "pbc_1776282104",
			"indexes": [
				"CREATE UNIQUE INDEX idx_work_records_type_consumables_type_consumable ON work_records_type_consumables (type, consumable)"
			],
			"listRule": "@request.auth.id != \"\"",
			"name": "work_records_type_consumables",
			"system": false,
			"type": "base",
			"updateRule": null,
			"viewRule": "@request.auth.id != \"\""
		}`

		workRecordsNotesJSON := `{
			"createRule": null,
			"deleteRule": null,
			"fields": [
				{
					"autogeneratePattern": "[a-z0-9]{15}",
					"hidden": false,
					"id": "text3208210256",
					"max": 15,
					"min": 15,
					"name": "id",
					"pattern": "^[a-z0-9]+$",
					"presentable": false,
					"primaryKey": true,
					"required": true,
					"system": true,
					"type": "text"
				},
				{
					"cascadeDelete": true,
					"collectionId": "pbc_1776282103",
					"hidden": false,
					"id": "relation1776282110",
					"maxSelect": 1,
					"minSelect": 0,
					"name": "work_record",
					"presentable": false,
					"required": true,
					"system": false,
					"type": "relation"
				},
				{
					"cascadeDelete": false,
					"collectionId": "_pb_users_auth_",
					"hidden": false,
					"id": "relation1776282111",
					"maxSelect": 1,
					"minSelect": 0,
					"name": "uid",
					"presentable": false,
					"required": true,
					"system": false,
					"type": "relation"
				},
				{
					"autogeneratePattern": "",
					"hidden": false,
					"id": "text1776282114",
					"max": 0,
					"min": 0,
					"name": "note",
					"pattern": "",
					"presentable": false,
					"primaryKey": false,
					"required": true,
					"system": false,
					"type": "text"
				},
				{
					"hidden": false,
					"id": "autodate2990389176",
					"name": "created",
					"onCreate": true,
					"onUpdate": false,
					"presentable": false,
					"system": false,
					"type": "autodate"
				}
			],
			"id": "pbc_1776282105",
			"indexes": [
				"CREATE INDEX idx_work_records_notes_work_record_created ON work_records_notes (work_record, created)"
			],
			"listRule": null,
			"name": "work_records_notes",
			"system": false,
			"type": "base",
			"updateRule": null,
			"viewRule": null
		}`

		workRecordsConsumableEntriesJSON := `{
			"createRule": null,
			"deleteRule": null,
			"fields": [
				{
					"autogeneratePattern": "[a-z0-9]{15}",
					"hidden": false,
					"id": "text3208210256",
					"max": 15,
					"min": 15,
					"name": "id",
					"pattern": "^[a-z0-9]+$",
					"presentable": false,
					"primaryKey": true,
					"required": true,
					"system": true,
					"type": "text"
				},
				{
					"cascadeDelete": true,
					"collectionId": "pbc_1776282103",
					"hidden": false,
					"id": "relation1776282112",
					"maxSelect": 1,
					"minSelect": 0,
					"name": "work_record",
					"presentable": false,
					"required": true,
					"system": false,
					"type": "relation"
				},
				{
					"cascadeDelete": false,
					"collectionId": "pbc_1776282102",
					"hidden": false,
					"id": "relation1776282113",
					"maxSelect": 1,
					"minSelect": 0,
					"name": "consumable",
					"presentable": false,
					"required": true,
					"system": false,
					"type": "relation"
				},
				{
					"hidden": false,
					"id": "number1776282108",
					"max": null,
					"min": null,
					"name": "quantity_number",
					"onlyInt": false,
					"presentable": false,
					"required": false,
					"system": false,
					"type": "number"
				},
				{
					"hidden": false,
					"id": "bool1776282106",
					"name": "selected",
					"presentable": false,
					"required": false,
					"system": false,
					"type": "bool"
				},
				{
					"hidden": false,
					"id": "autodate2990389176",
					"name": "created",
					"onCreate": true,
					"onUpdate": false,
					"presentable": false,
					"system": false,
					"type": "autodate"
				},
				{
					"hidden": false,
					"id": "autodate3332085495",
					"name": "updated",
					"onCreate": true,
					"onUpdate": true,
					"presentable": false,
					"system": false,
					"type": "autodate"
				}
			],
			"id": "pbc_1776282106",
			"indexes": [
				"CREATE UNIQUE INDEX idx_work_records_consumable_entries_work_record_consumable ON work_records_consumable_entries (work_record, consumable)"
			],
			"listRule": null,
			"name": "work_records_consumable_entries",
			"system": false,
			"type": "base",
			"updateRule": null,
			"viewRule": null
		}`

		if err := saveCollectionFromJSON(app, workRecordsTypesJSON); err != nil {
			return err
		}

		if err := saveCollectionFromJSON(app, workRecordsConsumablesJSON); err != nil {
			return err
		}

		if err := saveCollectionFromJSON(app, workRecordsJSON); err != nil {
			return err
		}

		workRecordsCollection, err := app.FindCollectionByNameOrId("pbc_1776282103")
		if err != nil {
			return err
		}

		// PocketBase validates relation targets while creating a collection, so
		// the self-referencing parent_work_record field has to be added only after
		// work_records itself has been saved once.
		if err := workRecordsCollection.Fields.AddMarshaledJSONAt(2, []byte(`{
			"cascadeDelete": false,
			"collectionId": "pbc_1776282103",
			"hidden": false,
			"id": "relation1776282103",
			"maxSelect": 1,
			"minSelect": 0,
			"name": "parent_work_record",
			"presentable": false,
			"required": false,
			"system": false,
			"type": "relation"
		}`)); err != nil {
			return err
		}

		workRecordsCollection.UpdateRule = nil

		if err := app.Save(workRecordsCollection); err != nil {
			return err
		}

		if err := saveCollectionFromJSON(app, workRecordsSubjectsJSON); err != nil {
			return err
		}

		if err := saveCollectionFromJSON(app, workRecordsTypeConsumablesJSON); err != nil {
			return err
		}

		if err := saveCollectionFromJSON(app, workRecordsNotesJSON); err != nil {
			return err
		}

		if err := saveCollectionFromJSON(app, workRecordsConsumableEntriesJSON); err != nil {
			return err
		}

		return nil
	}, func(app core.App) error {
		for _, collectionID := range []string{
			"pbc_1776282106",
			"pbc_1776282105",
			"pbc_1776282104",
			"pbc_1776282107",
			"pbc_1776282103",
			"pbc_1776282102",
			"pbc_1776282101",
		} {
			collection, err := app.FindCollectionByNameOrId(collectionID)
			if err != nil {
				return err
			}

			if err := app.Delete(collection); err != nil {
				return err
			}
		}

		return nil
	})
}
