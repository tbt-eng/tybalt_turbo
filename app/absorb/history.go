package absorb

import (
	"fmt"

	"github.com/pocketbase/pocketbase/core"
)

// HistoryCollection keeps one row for each committed absorb.
const HistoryCollection = "absorb_history"

// RecordHistory copies an absorb action into absorb_history. The commit calls
// it in the transaction that deletes the action, so a committed absorb always
// has its row and an absorb that was undone never has one. committedBy may be
// nil when no user made the request.
func RecordHistory(txApp core.App, action *core.Record, committedBy *core.Record) error {
	collection, err := txApp.FindCollectionByNameOrId(HistoryCollection)
	if err != nil {
		return fmt.Errorf("finding %s: %w", HistoryCollection, err)
	}
	collectionName := action.GetString("collection_name")
	targetID := action.GetString("target_id")

	history := core.NewRecord(collection)
	history.Set("collection_name", collectionName)
	history.Set("target_id", targetID)
	// The kept record can change or be absorbed later, so keep its name.
	if target, err := txApp.FindRecordById(collectionName, targetID); err == nil {
		history.Set("target_name", RecordName(target))
	}
	history.Set("absorbed_records", action.GetString("absorbed_records"))
	history.Set("updated_references", action.GetString("updated_references"))
	history.Set("proposal", action.GetString("proposal"))
	history.Set("absorbed_by", action.GetString("absorbed_by"))
	history.Set("absorbed_at", action.GetDateTime("created"))
	if committedBy != nil {
		history.Set("committed_by", committedBy.Id)
	}
	if err := txApp.Save(history); err != nil {
		return fmt.Errorf("saving %s: %w", HistoryCollection, err)
	}
	return nil
}
