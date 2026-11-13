package hooks

import (
	"context"
	"testing"

	"tybalt/internal/testseed"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

// TestJobRateSheetRowsNeedHistoryActions verifies that rows cannot be written
// outside the rate sheet history actions, which keep jobs.rate_sheet in step.
func TestJobRateSheetRowsNeedHistoryActions(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()
	AddHooks(app)

	collection, err := app.FindCollectionByNameOrId("job_rate_sheets")
	if err != nil {
		t.Fatal(err)
	}
	row := core.NewRecord(collection)
	row.Set("job", "test_job_w_rs")
	row.Set("rate_sheet", "test_rs_alt_actv")
	row.Set("effective_date", "2026-01-02")
	if err := app.Save(row); err == nil {
		t.Fatal("expected a direct row write to be refused")
	}

	opening, err := app.FindFirstRecordByFilter("job_rate_sheets", "job = 'test_job_w_rs'")
	if err != nil {
		t.Fatal(err)
	}
	if err := app.Delete(opening); err == nil {
		t.Fatal("expected a direct row delete to be refused")
	}
}

// TestFirstRateSheetCreatesOpeningRow verifies that giving a project with no
// rate sheet its first one, by any job save, records the opening row.
func TestFirstRateSheetCreatesOpeningRow(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()
	AddHooks(app)

	job, err := app.FindRecordById("jobs", "u09fwwcg07y03m7")
	if err != nil {
		t.Fatal(err)
	}
	if job.GetString("rate_sheet") != "" {
		t.Fatal("test expects a project without a rate sheet")
	}
	job.Set("rate_sheet", "c41ofep525bcacj")
	ctx := WithProjectWorkflowMutation(context.Background(), ProjectWorkflowMutationRateSheet)
	if err := app.SaveNoValidateWithContext(ctx, job); err != nil {
		t.Fatal(err)
	}

	rows, err := app.FindRecordsByFilter("job_rate_sheets", "job = {:job}", "", 0, 0, dbx.Params{"job": job.Id})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].GetString("effective_date") != "" || rows[0].GetString("rate_sheet") != "c41ofep525bcacj" {
		t.Fatalf("rows = %v; want one opening row on c41ofep525bcacj", rows)
	}

	// A later save that does not change the sheet adds nothing.
	job.Set("description", job.GetString("description")+" (edited)")
	if err := app.SaveNoValidateWithContext(ctx, job); err != nil {
		t.Fatal(err)
	}
	if n, _ := app.CountRecords("job_rate_sheets", dbx.HashExp{"job": job.Id}); n != 1 {
		t.Fatalf("rows after a second save = %d; want 1", n)
	}
}
