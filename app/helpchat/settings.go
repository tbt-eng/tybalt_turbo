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
	"strings"

	"github.com/pocketbase/pocketbase/core"
)

// settingsKey is the app_config row that holds the settings.
const settingsKey = "help_chat"

// Settings come from app_config. See descriptions/app_config.md.
type Settings struct {
	// Enabled turns the chat on. It is off until an admin adds the setting.
	Enabled bool `json:"enabled"`
	// DeployGroup is the id of a deploy_groups row. With one, the chat is on
	// only for the users in that group. Empty means every user. See
	// utilities.InDeployGroup.
	DeployGroup string `json:"deploy_group"`
	// Provider is the company whose API runs the model. See providers.go.
	Provider string `json:"provider"`
	// Model is a model id of that provider.
	Model string `json:"model"`
	// Effort is how much the model thinks before it answers.
	Effort string `json:"effort"`
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
		Provider:           providerAnthropic,
		Model:              providers[providerAnthropic].defaultModel,
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
	// The default model belongs to the default provider. A row that names
	// another provider gets that provider's default, or must name a model.
	settings.Model = ""
	// A property that the chat doesn't know is an error. Otherwise a mistake in
	// the name deploy_group would give the chat to every user.
	decoder := json.NewDecoder(strings.NewReader(record.GetString("value")))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&settings); err != nil {
		return DefaultSettings(), fmt.Errorf("app_config %s: %w", settingsKey, err)
	}
	if settings.Model == "" {
		settings.Model = providers[settings.Provider].defaultModel
	}
	if err := settings.validate(); err != nil {
		return DefaultSettings(), fmt.Errorf("app_config %s: %w", settingsKey, err)
	}
	return settings, nil
}

func (s Settings) validate() error {
	chosen, ok := providers[s.Provider]
	if !ok {
		return fmt.Errorf("provider must be one of %v", providerNames())
	}
	if s.Model == "" {
		return fmt.Errorf("model is required for the provider %s", s.Provider)
	}
	if err := chosen.check(s.Model, s.Effort); err != nil {
		return err
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
