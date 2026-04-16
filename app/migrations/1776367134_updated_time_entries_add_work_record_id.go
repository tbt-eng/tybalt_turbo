// Package migrations registers schema changes for PocketBase collections.
//
// This migration exists to add the relational link that the work-record branch
// is built around: time_entries.work_record_id. Historically, time entries only
// stored a free-text work_record number, which was enough for display but not
// for enforcing one-to-one linkage, validating ownership/job/date invariants,
// or deriving approval state through the first-class work_records table.
//
// The migration adds an optional relation field pointing at work_records and a
// unique partial index so a work record can be linked to at most one time
// entry, while still allowing many legacy or otherwise unlinked rows to remain
// valid. The helper functions in this file keep the migration idempotent by
// checking whether the index already exists before appending or removing it.
package migrations

import (
	"slices"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

const (
	timeEntriesWorkRecordFieldID = "relation1776368000"
	timeEntriesWorkRecordIndex   = "CREATE UNIQUE INDEX idx_time_entries_work_record_id ON time_entries (work_record_id) WHERE work_record_id != ''"
)

// ensureIndex appends the supplied index definition only when it is not already
// present.
func ensureIndex(indexes []string, value string) []string {
	if slices.Contains(indexes, value) {
		return indexes
	}

	return append(indexes, value)
}

// removeIndex returns the index slice without the supplied definition.
func removeIndex(indexes []string, value string) []string {
	filtered := make([]string, 0, len(indexes))
	for _, index := range indexes {
		if index == value {
			continue
		}
		filtered = append(filtered, index)
	}
	return filtered
}

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("ranctx5xgih6n3a")
		if err != nil {
			return err
		}

		if collection.Fields.GetById(timeEntriesWorkRecordFieldID) == nil {
			if err := collection.Fields.AddMarshaledJSONAt(9, []byte(`{
				"cascadeDelete": false,
				"collectionId": "pbc_1776282103",
				"hidden": false,
				"id": "`+timeEntriesWorkRecordFieldID+`",
				"maxSelect": 1,
				"minSelect": 0,
				"name": "work_record_id",
				"presentable": false,
				"required": false,
				"system": false,
				"type": "relation"
			}`)); err != nil {
				return err
			}
		}

		collection.Indexes = ensureIndex(collection.Indexes, timeEntriesWorkRecordIndex)
		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("ranctx5xgih6n3a")
		if err != nil {
			return err
		}

		collection.Fields.RemoveById(timeEntriesWorkRecordFieldID)
		collection.Indexes = removeIndex(collection.Indexes, timeEntriesWorkRecordIndex)
		return app.Save(collection)
	})
}
