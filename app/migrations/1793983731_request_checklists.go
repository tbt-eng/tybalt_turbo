package migrations

import (
	"encoding/json"
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/types"
)

const (
	requestChecklistsCollectionID = "pbc_req_checklst"
	requestChecklistsVersionIndex = "idx_request_checklists_version"
	checklistAnswersFieldID       = "pa4006f9ce47288"
)

type requestChecklistQuestion struct {
	Key  string `json:"key"`
	Text string `json:"text"`
}

var (
	requestChecklistFirstVersion = []requestChecklistQuestion{
		{"fn_agreement_reviewed", "I reviewed the First Nations agreement answer."},
		{"no_known_conflict", "I am not aware of a company conflict of interest for this project."},
		{"no_executive_restriction", "There is no executive restriction or moratorium that prevents this work."},
		{"no_duplicate_work", "I checked existing jobs and this request does not duplicate work already represented by a job."},
		{"pdf_sufficient", "The attached PDF is sufficient to authorize this project and its stated value."},
		{"allocations_complete", "All necessary division allocations are present."},
	}
	requestChecklistSecondVersion = []requestChecklistQuestion{
		requestChecklistFirstVersion[0],
		requestChecklistFirstVersion[1],
		requestChecklistFirstVersion[2],
		requestChecklistFirstVersion[3],
		requestChecklistFirstVersion[4],
		{"rate_sheet_correct", "The selected rate sheet is correct for this client and project."},
		requestChecklistFirstVersion[5],
	}
)

// Keep each version of the request checklist's questions in the database. An
// approval records only the version, written after every question is answered
// true. A saved answer could only ever be true, so checklist_answers recorded
// nothing the version doesn't. The job holder who creates the job answers the
// same questions, recorded in creation_checklist_version.
//
// The up migration refuses to drop checklist_answers if any saved answer isn't
// true. The down migration rebuilds the answers from each row's version, and
// loses creation_checklist_version.
func init() {
	m.Register(func(app core.App) error {
		checklists := core.NewBaseCollection("request_checklists")
		checklists.Id = requestChecklistsCollectionID
		checklists.Fields.Add(
			&core.NumberField{Name: "version", OnlyInt: true, Min: types.Pointer(1.0), Required: true},
			&core.JSONField{Name: "questions", Required: true},
			&core.AutodateField{Name: "created", OnCreate: true},
			&core.AutodateField{Name: "updated", OnCreate: true, OnUpdate: true},
		)
		checklists.AddIndex(requestChecklistsVersionIndex, true, "`version`", "")
		if err := app.Save(checklists); err != nil {
			return err
		}
		for version, questions := range [][]requestChecklistQuestion{requestChecklistFirstVersion, requestChecklistSecondVersion} {
			record := core.NewRecord(checklists)
			record.Set("version", version+1)
			record.Set("questions", questions)
			if err := app.Save(record); err != nil {
				return err
			}
		}

		if err := requireTrueChecklistAnswers(app); err != nil {
			return err
		}
		attachments, err := app.FindCollectionByNameOrId("job_attachments")
		if err != nil {
			return err
		}
		attachments.Fields.RemoveByName("checklist_answers")
		attachments.Fields.Add(&core.NumberField{
			Name:    "creation_checklist_version",
			OnlyInt: true,
			Min:     types.Pointer(0.0),
		})
		return app.Save(attachments)
	}, func(app core.App) error {
		attachments, err := app.FindCollectionByNameOrId("job_attachments")
		if err != nil {
			return err
		}
		attachments.Fields.RemoveByName("creation_checklist_version")
		answers := &core.JSONField{Name: "checklist_answers"}
		answers.SetId(checklistAnswersFieldID)
		attachments.Fields.Add(answers)
		if err := app.Save(attachments); err != nil {
			return err
		}
		for version, questions := range [][]requestChecklistQuestion{requestChecklistFirstVersion, requestChecklistSecondVersion} {
			rebuilt := map[string]bool{}
			for _, question := range questions {
				rebuilt[question.Key] = true
			}
			encoded, err := json.Marshal(rebuilt)
			if err != nil {
				return err
			}
			if _, err := app.DB().NewQuery("UPDATE job_attachments SET checklist_answers = {:answers} WHERE checklist_version = {:version}").
				Bind(dbx.Params{"answers": string(encoded), "version": version + 1}).Execute(); err != nil {
				return err
			}
		}
		checklists, err := app.FindCollectionByNameOrId(requestChecklistsCollectionID)
		if err != nil {
			return err
		}
		return app.Delete(checklists)
	})
}

// requireTrueChecklistAnswers stops the migration unless every saved approval
// answered every question of its version with true.
func requireTrueChecklistAnswers(app core.App) error {
	var rows []struct {
		ID      string `db:"id"`
		Version int    `db:"checklist_version"`
		Answers string `db:"checklist_answers"`
	}
	if err := app.DB().NewQuery("SELECT id, checklist_version, COALESCE(checklist_answers, '') AS checklist_answers FROM job_attachments WHERE checklist_version > 0").All(&rows); err != nil {
		return err
	}
	for _, row := range rows {
		questions := requestChecklistSecondVersion
		switch row.Version {
		case 1:
			questions = requestChecklistFirstVersion
		case 2:
		default:
			return fmt.Errorf("job_attachments %s has unknown checklist version %d", row.ID, row.Version)
		}
		answers := map[string]any{}
		if err := json.Unmarshal([]byte(row.Answers), &answers); err != nil {
			return fmt.Errorf("job_attachments %s has invalid checklist answers: %w", row.ID, err)
		}
		for _, question := range questions {
			if answers[question.Key] != true {
				return fmt.Errorf("job_attachments %s did not answer %s with true", row.ID, question.Key)
			}
		}
	}
	return nil
}
