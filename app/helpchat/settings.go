// Package helpchat answers questions about Turbo in Help.
//
// A signed-in user posts a question to /api/help/ask. The server builds a
// prompt from the user docs and the user's claims, lets the model search the
// corpus (package helpcorpus), streams the answer back, and writes one row to
// chat_turns. The model has no tool that writes.
//
// The design is in descriptions/help_chat.md.
package helpchat

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"slices"

	"github.com/pocketbase/pocketbase/core"
)

// settingsKey is the app_config row that holds the settings.
const settingsKey = "help_chat"

// The models that the chat can use. The prompt and the limits are tuned for them.
var allowedModels = []string{"claude-haiku-5-5", "claude-sonnet-5-5"}

var allowedEfforts = []string{"low", "medium", "high"}

// Settings come from app_config. See descriptions/app_config.md.
type Settings struct {
	// Enabled turns the chat on. It is off until an admin adds the setting.
	Enabled bool   `json:"enabled"`
	Model   string `json:"model"`
	Effort  string `json:"effort"`
	// MaxOutputTokens limits one model reply, with its thinking.
	MaxOutputTokens int `json:"max_output_tokens"`
	// MaxQuestionChars limits a question.
	MaxQuestionChars int `json:"max_question_chars"`
	// QuestionsPerHour limits one user.
	QuestionsPerHour int `json:"questions_per_hour"`
	// MaxToolCalls limits the searches and reads in one turn.
	MaxToolCalls int `json:"max_tool_calls"`
	// HistoryTurns is how many earlier turns of the chat go into the prompt.
	HistoryTurns int `json:"history_turns"`
	// TurnTimeoutSeconds limits one turn, from the question to the last word.
	TurnTimeoutSeconds int `json:"turn_timeout_seconds"`
}

// DefaultSettings are the values for properties that the row doesn't set.
func DefaultSettings() Settings {
	return Settings{
		Enabled:            false,
		Model:              "claude-haiku-5-5",
		Effort:             "low",
		MaxOutputTokens:    8000,
		MaxQuestionChars:   2000,
		QuestionsPerHour:   30,
		MaxToolCalls:       6,
		HistoryTurns:       10,
		TurnTimeoutSeconds: 60,
	}
}

// LoadSettings reads the settings. With no row, the chat is off. A bad value is
// an error, so a typing mistake can't remove a limit.
func LoadSettings(app core.App) (Settings, error) {
	settings := DefaultSettings()
	record, err := app.FindFirstRecordByData("app_config", "key", settingsKey)
	if errors.Is(err, sql.ErrNoRows) {
		return settings, nil
	}
	if err != nil {
		return settings, err
	}
	if err := json.Unmarshal([]byte(record.GetString("value")), &settings); err != nil {
		return DefaultSettings(), fmt.Errorf("app_config %s: %w", settingsKey, err)
	}
	if err := settings.validate(); err != nil {
		return DefaultSettings(), fmt.Errorf("app_config %s: %w", settingsKey, err)
	}
	return settings, nil
}

func (s Settings) validate() error {
	if !slices.Contains(allowedModels, s.Model) {
		return fmt.Errorf("model must be one of %v", allowedModels)
	}
	if !slices.Contains(allowedEfforts, s.Effort) {
		return fmt.Errorf("effort must be one of %v", allowedEfforts)
	}
	limits := []struct {
		name     string
		value    int
		min, max int
	}{
		{"max_output_tokens", s.MaxOutputTokens, 500, 64000},
		{"max_question_chars", s.MaxQuestionChars, 20, 4000},
		{"questions_per_hour", s.QuestionsPerHour, 1, 1000},
		{"max_tool_calls", s.MaxToolCalls, 0, 20},
		{"history_turns", s.HistoryTurns, 0, 50},
		{"turn_timeout_seconds", s.TurnTimeoutSeconds, 5, 300},
	}
	for _, limit := range limits {
		if limit.value < limit.min || limit.value > limit.max {
			return fmt.Errorf("%s must be from %d to %d", limit.name, limit.min, limit.max)
		}
	}
	return nil
}
