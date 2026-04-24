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
		    work_record_subject_id = {:subjectId},
		    work_record = ''
		WHERE id = {:id}
	`).Bind(map[string]any{
		"tsid":      "pts000000000001",
		"subjectId": "wrsub260400003",
		"id":        "wrkrcdentry0001",
	}).Execute(); err != nil {
		t.Fatalf("failed to update seeded time entry: %v", err)
	}

	if _, err := app.DB().NewQuery(`
		UPDATE work_records_subjects
		SET approved = 0
		WHERE id = {:id}
	`).Bind(map[string]any{
		"id": "wrsub260400003",
	}).Execute(); err != nil {
		t.Fatalf("failed to reset work record subject approval: %v", err)
	}

	if err := SyncApprovalForTimeSheet(app, "pts000000000001"); err != nil {
		t.Fatalf("expected sync to succeed: %v", err)
	}

	var row struct {
		Approved bool `db:"approved"`
	}
	if err := app.DB().NewQuery(`
		SELECT approved
		FROM work_records_subjects
		WHERE id = {:id}
	`).Bind(map[string]any{
		"id": "wrsub260400003",
	}).One(&row); err != nil {
		t.Fatalf("failed to reload work record subject approval: %v", err)
	}
	if !row.Approved {
		t.Fatalf("expected linked work record subject to become approved")
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
		    work_record_subject_id = {:subjectId},
		    work_record = ''
		WHERE id = {:id}
	`).Bind(map[string]any{
		"tsid":      "pts000000000001",
		"subjectId": "wrsub260400003",
		"id":        "wrkrcdentry0001",
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
		FROM work_records_subjects
		WHERE id = {:id}
	`).Bind(map[string]any{
		"id": "wrsub260400003",
	}).One(&row); err != nil {
		t.Fatalf("failed to reload work record subject approval: %v", err)
	}
	if row.Approved {
		t.Fatalf("expected linked work record subject approval to clear after rejection")
	}
}

func TestSyncApprovalForWorkRecordSubjectIDs_ClearsApprovalWithoutApprovedLinks(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	t.Cleanup(app.Cleanup)
	if err := app.RunAllMigrations(); err != nil {
		t.Fatalf("failed to apply migrations in test app: %v", err)
	}

	if _, err := app.DB().NewQuery(`
		UPDATE work_records_subjects
		SET approved = 1
		WHERE id = {:id}
	`).Bind(map[string]any{
		"id": "wrsub260400003",
	}).Execute(); err != nil {
		t.Fatalf("failed to seed approved flag: %v", err)
	}

	if err := SyncApprovalForWorkRecordSubjectIDs(app, "wrsub260400003"); err != nil {
		t.Fatalf("expected sync to succeed: %v", err)
	}

	var row struct {
		Approved bool `db:"approved"`
	}
	if err := app.DB().NewQuery(`
		SELECT approved
		FROM work_records_subjects
		WHERE id = {:id}
	`).Bind(map[string]any{
		"id": "wrsub260400003",
	}).One(&row); err != nil {
		t.Fatalf("failed to reload work record subject approval: %v", err)
	}
	if row.Approved {
		t.Fatalf("expected approval to clear when there are no approved linked timesheets")
	}
}
