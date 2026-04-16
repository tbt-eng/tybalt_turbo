// Package workrecords contains work-record-specific domain logic that does not
// fit cleanly inside a single route, hook, or migration.
//
// This file exists to keep the derived work_records.approved flag synchronized
// with the real approval state of the linked time-entry/timesheet workflow.
// The feature added in this branch makes time_entries.work_record_id the
// canonical one-to-one link between a time entry and a work record, but a work
// record's "approved" state is still determined indirectly: a linked time entry
// must belong to a timesheet, that timesheet must be approved, and it must not
// be rejected.
//
// Multiple code paths can change that derived truth:
// creating or editing a linked time entry, bundling and unbundling a timesheet,
// approving or rejecting a timesheet, and any later recall-style flow that
// reopens approval. Instead of re-encoding the approval rules in every caller,
// this file gives the rest of the application one place to ask "recompute the
// linked work record state from the database as it exists right now."
//
// The implementation intentionally re-queries the database rather than trusting
// route-local assumptions. That keeps the behavior correct across standalone and
// bundled entries, avoids stale in-memory decisions, and lets callers stay
// small: they only need to identify which work record IDs or timesheet ID may
// have changed, then delegate the recomputation here.
package workrecords

import (
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type approvalStateRow struct {
	Approved bool `db:"approved"`
}

func uniqueNonEmptyIDs(ids ...string) []string {
	seen := map[string]struct{}{}
	result := make([]string, 0, len(ids))

	for _, id := range ids {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		result = append(result, id)
	}

	return result
}

func shouldApproveWorkRecord(app core.App, workRecordID string) (bool, error) {
	var row approvalStateRow
	err := app.DB().NewQuery(`
		SELECT EXISTS(
			SELECT 1
			FROM time_entries te
			JOIN time_sheets ts
			  ON ts.id = te.tsid
			WHERE te.work_record_id = {:workRecordId}
			  AND te.tsid != ''
			  AND ts.approved != ''
			  AND ts.rejected = ''
		) AS approved
	`).Bind(dbx.Params{
		"workRecordId": workRecordID,
	}).One(&row)
	if err != nil {
		return false, err
	}

	return row.Approved, nil
}

// SyncApprovalForWorkRecordIDs recomputes and persists the approved flag for
// each referenced work record.
//
// Callers pass the set of work record IDs whose linkage or approval context may
// have changed. The function deduplicates IDs, ignores blanks, skips missing
// records defensively, and then asks shouldApproveWorkRecord whether each record
// currently has at least one linked time entry on an approved, non-rejected
// timesheet. The stored approved flag is only updated when the derived state
// actually differs, which keeps no-op transitions cheap.
func SyncApprovalForWorkRecordIDs(app core.App, ids ...string) error {
	for _, id := range uniqueNonEmptyIDs(ids...) {
		approved, err := shouldApproveWorkRecord(app, id)
		if err != nil {
			return err
		}

		if _, err := app.DB().NewQuery(`
			UPDATE work_records
			SET approved = {:approved}
			WHERE id = {:id}
		`).Bind(dbx.Params{
			"approved": approved,
			"id":       id,
		}).Execute(); err != nil {
			return err
		}
	}

	return nil
}

// SyncApprovalForTimeSheet finds every work record linked through the specified
// timesheet and then delegates to SyncApprovalForWorkRecordIDs.
//
// This helper exists so routes that mutate a whole timesheet do not need to
// manually discover every distinct linked work record themselves. It treats the
// timesheet as the unit of change, queries the current linked work_record_id
// values, and recomputes approval from that fresh database state.
func SyncApprovalForTimeSheet(app core.App, timeSheetID string) error {
	timeSheetID = strings.TrimSpace(timeSheetID)
	if timeSheetID == "" {
		return nil
	}

	type row struct {
		WorkRecordID string `db:"work_record_id"`
	}

	var rows []row
	if err := app.DB().NewQuery(`
		SELECT DISTINCT te.work_record_id
		FROM time_entries te
		WHERE te.tsid = {:timeSheetId}
		  AND te.work_record_id != ''
	`).Bind(dbx.Params{
		"timeSheetId": timeSheetID,
	}).All(&rows); err != nil {
		return err
	}

	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.WorkRecordID)
	}

	return SyncApprovalForWorkRecordIDs(app, ids...)
}
