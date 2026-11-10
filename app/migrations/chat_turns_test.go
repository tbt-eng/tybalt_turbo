package migrations_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/core"
)

// A user reads their own turns. Only server code writes them.
func TestChatTurnsRules(t *testing.T) {
	app := newMigratedApp(t)
	turns, err := app.FindCollectionByNameOrId("chat_turns")
	if err != nil {
		t.Fatal(err)
	}
	own := "@request.auth.id != '' && uid = @request.auth.id"
	if turns.ListRule == nil || *turns.ListRule != own || turns.ViewRule == nil || *turns.ViewRule != own {
		t.Fatalf("list and view must be for the owner, got %v and %v", turns.ListRule, turns.ViewRule)
	}
	if turns.CreateRule != nil || turns.UpdateRule != nil || turns.DeleteRule != nil {
		t.Fatal("the collection must have no create, update or delete rule")
	}
	uid, ok := turns.Fields.GetByName("uid").(*core.RelationField)
	if !ok || !uid.Required || !uid.CascadeDelete {
		t.Fatalf("uid must be a required relation that goes when the user goes, got %+v", uid)
	}
	for _, name := range []string{
		"chat", "turn", "question", "answer", "sources", "model", "effort", "build", "corpus", "claims", "tool_calls",
		"route", "record", "surface", "status", "gap", "stop_reason", "error", "request_id",
		"tokens_in", "tokens_cache_read", "tokens_cache_write", "tokens_out", "ms_first_token", "ms_total",
		"rating", "rating_note", "opened_sources", "copied", "created", "updated",
	} {
		if turns.Fields.GetByName(name) == nil {
			t.Errorf("chat_turns has no %s field", name)
		}
	}
}

func TestChatTurnsMigrationRollsBack(t *testing.T) {
	app := newMigratedApp(t)
	for _, item := range core.AppMigrations.Items() {
		if item.File == "1794265765_chat_turns.go" {
			if err := item.Down(app); err != nil {
				t.Fatal(err)
			}
			if _, err := app.FindCollectionByNameOrId("chat_turns"); err == nil {
				t.Fatal("expected the collection to be removed")
			}
			return
		}
	}
	t.Fatal("the chat_turns migration is not registered")
}
