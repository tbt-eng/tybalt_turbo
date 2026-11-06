package utilities

import (
	"errors"
	"fmt"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

// RequestChecklistQuestion is one statement of a project request checklist. An
// approver or job holder must answer it true to continue.
type RequestChecklistQuestion struct {
	Key  string `json:"key"`
	Text string `json:"text"`
}

// RequestChecklistQuestions reads a question set and checks its shape.
func RequestChecklistQuestions(record *core.Record) ([]RequestChecklistQuestion, error) {
	var questions []RequestChecklistQuestion
	if err := record.UnmarshalJSONField("questions", &questions); err != nil {
		return nil, fmt.Errorf("request checklist questions must be a list: %w", err)
	}
	if len(questions) == 0 {
		return nil, errors.New("a request checklist needs at least one question")
	}
	seen := map[string]bool{}
	for _, question := range questions {
		if strings.TrimSpace(question.Key) == "" || strings.TrimSpace(question.Text) == "" {
			return nil, errors.New("each request checklist question needs a key and text")
		}
		if seen[question.Key] {
			return nil, fmt.Errorf("request checklist question %s appears twice", question.Key)
		}
		seen[question.Key] = true
	}
	return questions, nil
}

// LatestRequestChecklistVersion returns the highest saved version, or zero when
// none is saved.
func LatestRequestChecklistVersion(app core.App) (int, error) {
	records, err := app.FindRecordsByFilter("request_checklists", "", "-version", 1, 0)
	if err != nil {
		return 0, err
	}
	if len(records) == 0 {
		return 0, nil
	}
	return records[0].GetInt("version"), nil
}

// CurrentRequestChecklist returns the questions approvers and job holders answer
// now, which are those of the highest version.
func CurrentRequestChecklist(app core.App) (int, []RequestChecklistQuestion, error) {
	records, err := app.FindRecordsByFilter("request_checklists", "", "-version", 1, 0)
	if err != nil {
		return 0, nil, err
	}
	if len(records) == 0 {
		return 0, nil, errors.New("no request checklist is saved")
	}
	questions, err := RequestChecklistQuestions(records[0])
	if err != nil {
		return 0, nil, err
	}
	return records[0].GetInt("version"), questions, nil
}

// RequestChecklistVersionExists reports whether a version has saved questions.
func RequestChecklistVersionExists(app core.App, version int) (bool, error) {
	count, err := app.CountRecords("request_checklists", dbx.HashExp{"version": version})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
