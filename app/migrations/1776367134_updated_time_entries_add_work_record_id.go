// Package migrations registers schema changes for PocketBase collections.
//
// This migration adds the time_entries.work_record_subject_id relation that the work
// records feature depends on. Before this change, time entries only carried the
// legacy free-text work_record field. That text value is fine for display, but
// it cannot enforce that a referenced work record actually exists, cannot
// guarantee one-to-one linkage, and cannot support relational queries and rules
// that join time entries back to the first-class work_records_subjects
// collection.
//
// Concretely, the forward migration does three things:
//
//  1. It loads the existing time_entries collection definition by its stable
//     PocketBase collection id.
//  2. It inserts a new optional relation field named work_record_subject_id at
//     field position 9. The relation targets the work_records_subjects
//     collection (pbc_1776282107), allows at most one selected record, and is
//     intentionally nullable so legacy rows without a relational match remain
//     valid.
//  3. It appends a partial unique index:
//     CREATE UNIQUE INDEX idx_time_entries_work_record_subject_id
//     ON time_entries (work_record_subject_id)
//     WHERE work_record_subject_id != ''
//     That index enforces the business rule that a single worker row may be
//     linked to at most one time entry, while still permitting any number of rows
//     with an empty relation value.
//
// The rollback does the exact inverse: it removes the relation field from the
// collection metadata and removes the partial unique index definition.
//
// The migration looks small in Go, but app.Save(collection) is expensive here.
// PocketBase persists collection changes by reconciling the underlying SQLite
// table schema and index set. When PocketBase sees a collection schema change
// that makes the indexes "dirty" (for example a field list change, collection
// rename, or index list change), it drops all existing indexes for that
// collection and recreates the full set afterward rather than applying only a
// narrowly targeted delta. Because this migration modifies the large
// time_entries collection, SQLite ends up adding the new column and then
// dropping and recreating the table's existing indexes so the stored schema
// matches the updated collection definition. In practice, the runtime cost is
// dominated by rebuilding the pre-existing time_entries indexes rather than by
// adding work_record_subject_id itself. In local testing captured in
// app/wr_migration_logs, this migration took about 5.3 seconds end-to-end, and
// about 5.279 seconds of that time was spent in the block that drops and
// recreates time_entries indexes. As of PocketBase v0.37.3, the collection
// schema sync code still appears to use this same drop-all-and-recreate-all
// index strategy when a collection update dirties the index set.
//
// The helper functions in this file keep the migration idempotent with respect
// to the new index definition so repeated schema loads do not append duplicate
// CREATE INDEX statements.
package migrations

import (
	"slices"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

const (
	timeEntriesWorkRecordSubjectFieldID = "relation1776368000"
	timeEntriesWorkRecordSubjectIndex   = "CREATE UNIQUE INDEX idx_time_entries_work_record_subject_id ON time_entries (work_record_subject_id) WHERE work_record_subject_id != ''"
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

		if collection.Fields.GetById(timeEntriesWorkRecordSubjectFieldID) == nil {
			if err := collection.Fields.AddMarshaledJSONAt(9, []byte(`{
				"cascadeDelete": false,
				"collectionId": "pbc_1776282107",
				"hidden": false,
				"id": "`+timeEntriesWorkRecordSubjectFieldID+`",
				"maxSelect": 1,
				"minSelect": 0,
				"name": "work_record_subject_id",
				"presentable": false,
				"required": false,
				"system": false,
				"type": "relation"
			}`)); err != nil {
				return err
			}
		}

		collection.Indexes = ensureIndex(collection.Indexes, timeEntriesWorkRecordSubjectIndex)
		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("ranctx5xgih6n3a")
		if err != nil {
			return err
		}

		collection.Fields.RemoveById(timeEntriesWorkRecordSubjectFieldID)
		collection.Indexes = removeIndex(collection.Indexes, timeEntriesWorkRecordSubjectIndex)
		return app.Save(collection)
	})
}
