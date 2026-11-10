package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/types"
)

const chatTurnsCollectionID = "pbc_chat_turns"

// Add chat_turns, the log of Help chat.
//
// One row is one question and its answer. The row records how the answer was
// made (provider, model, corpus, claims, tool calls), what it cost, and what the user did
// with it, so that we can improve the chat later.
//
// A user can list and view their own rows. Only server code writes rows: the
// collection has no create, update or delete rule. Feedback goes through
// /api/help/turns/{id}/feedback.
//
// The log holds pointers to the corpus (chunk ids), not its text. The corpus
// hash and the build say which text those pointers mean.
func init() {
	m.Register(func(app core.App) error {
		users, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}

		turns := core.NewBaseCollection("chat_turns")
		turns.Id = chatTurnsCollectionID
		turns.ListRule = types.Pointer("@request.auth.id != '' && uid = @request.auth.id")
		turns.ViewRule = types.Pointer("@request.auth.id != '' && uid = @request.auth.id")
		turns.Fields.Add(
			// Identity.
			&core.RelationField{Name: "uid", CollectionId: users.Id, MaxSelect: 1, Required: true, CascadeDelete: true},
			&core.TextField{Name: "chat", Required: true, Max: 15},
			&core.NumberField{Name: "turn", OnlyInt: true},

			// Content. sources lists what the answer cites.
			&core.TextField{Name: "question", Required: true, Max: 4000},
			&core.TextField{Name: "answer", Max: 60000},
			&core.JSONField{Name: "sources"},

			// How the answer was made.
			&core.TextField{Name: "provider", Max: 50},
			&core.TextField{Name: "model", Max: 100},
			&core.TextField{Name: "effort", Max: 20},
			&core.TextField{Name: "build", Max: 100},
			&core.TextField{Name: "corpus", Max: 64},
			&core.JSONField{Name: "claims"},
			&core.JSONField{Name: "tool_calls"},

			// Where the user was. The model gets route. record stays on the server.
			&core.TextField{Name: "route", Max: 200},
			&core.TextField{Name: "record", Max: 100},
			&core.SelectField{Name: "surface", Values: []string{"help", "popover"}, MaxSelect: 1},

			// Result.
			&core.SelectField{Name: "status", Values: []string{"running", "answered", "gap", "refused", "cancelled", "error"}, MaxSelect: 1, Required: true},
			&core.JSONField{Name: "gap"},
			&core.TextField{Name: "stop_reason", Max: 50},
			&core.TextField{Name: "error", Max: 2000},
			&core.TextField{Name: "request_id", Max: 100},

			// Cost and speed.
			&core.NumberField{Name: "tokens_in", OnlyInt: true},
			&core.NumberField{Name: "tokens_cache_read", OnlyInt: true},
			&core.NumberField{Name: "tokens_cache_write", OnlyInt: true},
			&core.NumberField{Name: "tokens_out", OnlyInt: true},
			&core.NumberField{Name: "ms_first_token", OnlyInt: true},
			&core.NumberField{Name: "ms_total", OnlyInt: true},

			// Feedback from the user.
			&core.SelectField{Name: "rating", Values: []string{"up", "down"}, MaxSelect: 1},
			&core.TextField{Name: "rating_note", Max: 1000},
			&core.JSONField{Name: "opened_sources"},
			&core.BoolField{Name: "copied"},

			&core.AutodateField{Name: "created", OnCreate: true},
			&core.AutodateField{Name: "updated", OnCreate: true, OnUpdate: true},
		)
		turns.AddIndex("idx_chat_turns_chat_turn", true, "`chat`, `turn`", "")
		turns.AddIndex("idx_chat_turns_uid_created", false, "`uid`, `created`", "")
		return app.Save(turns)
	}, func(app core.App) error {
		turns, err := app.FindCollectionByNameOrId(chatTurnsCollectionID)
		if err != nil {
			return nil
		}
		return app.Delete(turns)
	})
}
