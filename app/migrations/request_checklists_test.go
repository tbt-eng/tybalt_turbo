package migrations_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

const requestChecklistsMigration = "request_checklists.go"

func requestChecklistsItem(t *testing.T) *core.Migration {
	t.Helper()
	for _, item := range core.AppMigrations.Items() {
		if strings.HasSuffix(item.File, "_"+requestChecklistsMigration) {
			return item
		}
	}
	t.Fatal("request checklists migration is not registered")
	return nil
}

// Each version keeps the exact statements approvers were shown.
func TestRequestChecklistsMigrationSavesEachVersion(t *testing.T) {
	app := newMigratedApp(t)
	records, err := app.FindRecordsByFilter("request_checklists", "", "version", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 2 || records[0].GetInt("version") != 1 || records[1].GetInt("version") != 2 {
		t.Fatalf("got %d checklist versions; want versions 1 and 2", len(records))
	}
	var first, second []struct{ Key, Text string }
	if err := records[0].UnmarshalJSONField("questions", &first); err != nil {
		t.Fatal(err)
	}
	if err := records[1].UnmarshalJSONField("questions", &second); err != nil {
		t.Fatal(err)
	}
	if len(first) != 6 || len(second) != 7 || second[5].Key != "rate_sheet_correct" || second[5].Text != "The selected rate sheet is correct for this client and project." {
		t.Fatalf("unexpected questions: %v %v", first, second)
	}
	for _, question := range first {
		if question.Key == "rate_sheet_correct" {
			t.Fatal("version 1 did not ask about the rate sheet")
		}
	}
	attachments, err := app.FindCollectionByNameOrId("job_attachments")
	if err != nil {
		t.Fatal(err)
	}
	if attachments.Fields.GetByName("checklist_answers") != nil || attachments.Fields.GetByName("creation_checklist_version") == nil {
		t.Fatal("job_attachments should hold versions, not answers")
	}
}

// Going down and up keeps each approval's version. Only an all-true saved answer
// can be dropped, and anything else stops the migration.
func TestRequestChecklistsMigrationDownAndUp(t *testing.T) {
	app := newMigratedApp(t)
	item := requestChecklistsItem(t)
	if err := item.Down(app); err != nil {
		t.Fatal(err)
	}
	if _, err := app.FindCollectionByNameOrId("request_checklists"); err == nil {
		t.Fatal("expected the checklists collection to be removed")
	}
	attachments, err := app.FindCollectionByNameOrId("job_attachments")
	if err != nil {
		t.Fatal(err)
	}
	if attachments.Fields.GetByName("checklist_answers") == nil || attachments.Fields.GetByName("creation_checklist_version") != nil {
		t.Fatal("down did not restore the answers field")
	}
	insert := func(id string, version int, answers map[string]bool) {
		t.Helper()
		encoded, err := json.Marshal(answers)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := app.DB().NewQuery("INSERT INTO job_attachments (id, document, document_hash, uploader, status, checklist_version, checklist_answers) VALUES ({:id}, 'a.pdf', {:id}, 'uploader', 'Approved', {:version}, {:answers})").
			Bind(dbx.Params{"id": id, "version": version, "answers": string(encoded)}).Execute(); err != nil {
			t.Fatal(err)
		}
	}
	complete := map[string]bool{"fn_agreement_reviewed": true, "no_known_conflict": true, "no_executive_restriction": true, "no_duplicate_work": true, "pdf_sufficient": true, "allocations_complete": true}
	insert("migattachment01", 1, complete)
	incomplete := map[string]bool{"fn_agreement_reviewed": true, "no_known_conflict": false}
	insert("migattachment02", 1, incomplete)
	if err := item.Up(app); err == nil || !strings.Contains(err.Error(), "migattachment02") {
		t.Fatalf("expected the migration to stop at the false answer, got %v", err)
	}
	if _, err := app.DB().NewQuery("DELETE FROM job_attachments WHERE id='migattachment02'").Execute(); err != nil {
		t.Fatal(err)
	}
	// The stopped run left the collection, so start from the state before it.
	if collection, err := app.FindCollectionByNameOrId("request_checklists"); err == nil {
		if err := app.Delete(collection); err != nil {
			t.Fatal(err)
		}
	}
	if err := item.Up(app); err != nil {
		t.Fatal(err)
	}
	var version int
	if err := app.DB().NewQuery("SELECT checklist_version FROM job_attachments WHERE id='migattachment01'").Row(&version); err != nil || version != 1 {
		t.Fatalf("up changed the saved version to %d: %v", version, err)
	}
}
