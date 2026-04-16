package workrecords

import (
	"testing"
	"tybalt/internal/testseed"
)

func TestSyncApprovalForTimeSheet_SetsApprovedForLinkedWorkRecord(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	t.Cleanup(app.Cleanup)
	if err := app.RunAllMigrations(); err != nil {
		t.Fatalf("failed to apply migrations in test app: %v", err)
	}

	if _, err := app.DB().NewQuery(`
		UPDATE time_entries
		SET tsid = {:tsid},
		    work_record_id = {:workRecordId},
		    work_record = ''
		WHERE id = {:id}
	`).Bind(map[string]any{
		"tsid":         "pts000000000001",
		"workRecordId": "wr260400000003a",
		"id":           "wrkrcdentry0001",
	}).Execute(); err != nil {
		t.Fatalf("failed to update seeded time entry: %v", err)
	}

	if _, err := app.DB().NewQuery(`
		UPDATE work_records
		SET approved = 0
		WHERE id = {:id}
	`).Bind(map[string]any{
		"id": "wr260400000003a",
	}).Execute(); err != nil {
		t.Fatalf("failed to reset work record approval: %v", err)
	}

	if err := SyncApprovalForTimeSheet(app, "pts000000000001"); err != nil {
		t.Fatalf("expected sync to succeed: %v", err)
	}

	var row struct {
		Approved bool `db:"approved"`
	}
	if err := app.DB().NewQuery(`
		SELECT approved
		FROM work_records
		WHERE id = {:id}
	`).Bind(map[string]any{
		"id": "wr260400000003a",
	}).One(&row); err != nil {
		t.Fatalf("failed to reload work record approval: %v", err)
	}
	if !row.Approved {
		t.Fatalf("expected linked work record to become approved")
	}
}

func TestSyncApprovalForTimeSheet_ClearsApprovalAfterRejection(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	t.Cleanup(app.Cleanup)
	if err := app.RunAllMigrations(); err != nil {
		t.Fatalf("failed to apply migrations in test app: %v", err)
	}

	if _, err := app.DB().NewQuery(`
		UPDATE time_entries
		SET tsid = {:tsid},
		    work_record_id = {:workRecordId},
		    work_record = ''
		WHERE id = {:id}
	`).Bind(map[string]any{
		"tsid":         "pts000000000001",
		"workRecordId": "wr260400000003a",
		"id":           "wrkrcdentry0001",
	}).Execute(); err != nil {
		t.Fatalf("failed to update seeded time entry: %v", err)
	}

	if err := SyncApprovalForTimeSheet(app, "pts000000000001"); err != nil {
		t.Fatalf("expected initial sync to succeed: %v", err)
	}

	timeSheet, err := app.FindRecordById("time_sheets", "pts000000000001")
	if err != nil {
		t.Fatalf("failed to load timesheet: %v", err)
	}
	timeSheet.Set("rejected", "2026-04-26 12:00:00.000Z")
	if err := app.Save(timeSheet); err != nil {
		t.Fatalf("failed to reject timesheet in test: %v", err)
	}

	if err := SyncApprovalForTimeSheet(app, "pts000000000001"); err != nil {
		t.Fatalf("expected rejection sync to succeed: %v", err)
	}

	var row struct {
		Approved bool `db:"approved"`
	}
	if err := app.DB().NewQuery(`
		SELECT approved
		FROM work_records
		WHERE id = {:id}
	`).Bind(map[string]any{
		"id": "wr260400000003a",
	}).One(&row); err != nil {
		t.Fatalf("failed to reload work record approval: %v", err)
	}
	if row.Approved {
		t.Fatalf("expected linked work record approval to clear after rejection")
	}
}

func TestSyncApprovalForWorkRecordIDs_ClearsApprovalWithoutApprovedLinks(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	t.Cleanup(app.Cleanup)
	if err := app.RunAllMigrations(); err != nil {
		t.Fatalf("failed to apply migrations in test app: %v", err)
	}

	if _, err := app.DB().NewQuery(`
		UPDATE work_records
		SET approved = 1
		WHERE id = {:id}
	`).Bind(map[string]any{
		"id": "wr260400000003a",
	}).Execute(); err != nil {
		t.Fatalf("failed to seed approved flag: %v", err)
	}

	if err := SyncApprovalForWorkRecordIDs(app, "wr260400000003a"); err != nil {
		t.Fatalf("expected sync to succeed: %v", err)
	}

	var row struct {
		Approved bool `db:"approved"`
	}
	if err := app.DB().NewQuery(`
		SELECT approved
		FROM work_records
		WHERE id = {:id}
	`).Bind(map[string]any{
		"id": "wr260400000003a",
	}).One(&row); err != nil {
		t.Fatalf("failed to reload work record approval: %v", err)
	}
	if row.Approved {
		t.Fatalf("expected approval to clear when there are no approved linked timesheets")
	}
}
