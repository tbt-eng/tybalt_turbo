// Package migrations registers schema changes for PocketBase collections.
//
// This migration exists because work records and modern purchase orders both
// need backend-reserved numbers that remain correct under concurrency. Scanning
// existing records for the current maximum on every create is not enough once
// multiple writers can race, so the branch introduces a dedicated
// number_sequences collection to act as the durable counter store.
//
// The collection is intentionally minimal: a unique scope identifies which
// sequence is being reserved and an integer value stores the most recently
// issued number. The numbering package uses this table transactionally so
// callers can reserve the next identifier while saving related business records
// in the same transaction.
package migrations

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func saveNumberSequencesCollection(app core.App, jsonData string) error {
	collection := &core.Collection{}
	if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
		return err
	}

	return app.Save(collection)
}

func init() {
	m.Register(func(app core.App) error {
		numberSequencesJSON := `{
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
					"id": "text1776350691",
					"max": 0,
					"min": 0,
					"name": "scope",
					"pattern": "",
					"presentable": false,
					"primaryKey": false,
					"required": true,
					"system": false,
					"type": "text"
				},
				{
					"hidden": false,
					"id": "number1776350690",
					"max": null,
					"min": null,
					"name": "value",
					"onlyInt": true,
					"presentable": false,
					"required": true,
					"system": false,
					"type": "number"
				}
			],
			"id": "pbc_1776350690",
			"indexes": [
				"CREATE UNIQUE INDEX idx_number_sequences_scope ON number_sequences (scope)"
			],
			"listRule": null,
			"name": "number_sequences",
			"system": false,
			"type": "base",
			"updateRule": null,
			"viewRule": null
		}`

		return saveNumberSequencesCollection(app, numberSequencesJSON)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("number_sequences")
		if err != nil {
			return err
		}

		return app.Delete(collection)
	})
}
