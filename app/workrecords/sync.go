// Package workrecords contains work-record-specific domain logic that does not
// fit cleanly inside a single route, hook, or migration.
//
// This file exists to keep the work_records_subjects.approved flag synchronized
// with the real approval state of the linked time-entry/timesheet workflow.
// The feature added in this branch makes time_entries.work_record_subject_id the
// canonical one-to-one link between a time entry and a work record subject.
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
// small: they only need to identify which subject IDs or timesheet ID may
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

func shouldApproveWorkRecordSubject(app core.App, subjectID string) (bool, error) {
	var row approvalStateRow
	err := app.DB().NewQuery(`
		SELECT EXISTS(
			SELECT 1
			FROM time_entries te
			JOIN time_sheets ts
			  ON ts.id = te.tsid
			WHERE te.work_record_subject_id = {:subjectId}
			  AND te.tsid != ''
			  AND ts.approved != ''
			  AND ts.rejected = ''
		) AS approved
	`).Bind(dbx.Params{
		"subjectId": subjectID,
	}).One(&row)
	if err != nil {
		return false, err
	}

	return row.Approved, nil
}

// SyncApprovalForWorkRecordSubjectIDs recomputes and persists the approved flag
// for each referenced work record subject.
//
// Callers pass the set of subject IDs whose linkage or approval context may
// have changed. The function deduplicates IDs, ignores blanks, skips missing
// records defensively, and then asks shouldApproveWorkRecordSubject whether each row
// currently has at least one linked time entry on an approved, non-rejected
// timesheet. The stored approved flag is only updated when the derived state
// actually differs, which keeps no-op transitions cheap.
func SyncApprovalForWorkRecordSubjectIDs(app core.App, ids ...string) error {
	for _, id := range uniqueNonEmptyIDs(ids...) {
		approved, err := shouldApproveWorkRecordSubject(app, id)
		if err != nil {
			return err
		}

		if _, err := app.DB().NewQuery(`
			UPDATE work_records_subjects
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

// SyncApprovalForTimeSheet finds every work record subject linked through the
// specified timesheet and then delegates to SyncApprovalForWorkRecordSubjectIDs.
//
// This helper exists so routes that mutate a whole timesheet do not need to
// manually discover every distinct linked subject themselves. It treats the
// timesheet as the unit of change, queries the current linked
// work_record_subject_id
// values, and recomputes approval from that fresh database state.
func SyncApprovalForTimeSheet(app core.App, timeSheetID string) error {
	timeSheetID = strings.TrimSpace(timeSheetID)
	if timeSheetID == "" {
		return nil
	}

	type row struct {
		SubjectID string `db:"work_record_subject_id"`
	}

	var rows []row
	if err := app.DB().NewQuery(`
		SELECT DISTINCT te.work_record_subject_id
		FROM time_entries te
		WHERE te.tsid = {:timeSheetId}
		  AND te.work_record_subject_id != ''
	`).Bind(dbx.Params{
		"timeSheetId": timeSheetID,
	}).All(&rows); err != nil {
		return err
	}

	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.SubjectID)
	}

	return SyncApprovalForWorkRecordSubjectIDs(app, ids...)
}
