package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

const requestChecklistsScopeVersionIndex = "idx_request_checklists_scope_version"

// Scope each request checklist to the stage that answers it. Approval and job
// creation then keep separate lists with separate version numbers, so a change
// to one never asks pages open on the other to reload.
//
// Every saved version becomes a creation checklist, which leaves no approval
// checklist until one is added. Requests cannot be approved until then, and no
// approval refers to these versions because the project request flag has not
// been on. The down migration drops any approval checklists, whose version
// numbers could collide under the single-column index.
func init() {
	m.Register(func(app core.App) error {
		checklists, err := app.FindCollectionByNameOrId(requestChecklistsCollectionID)
		if err != nil {
			return err
		}
		checklists.Fields.Add(&core.SelectField{
			Name:      "scope",
			Values:    []string{"approval", "creation"},
			MaxSelect: 1,
			Required:  true,
		})
		checklists.RemoveIndex(requestChecklistsVersionIndex)
		checklists.AddIndex(requestChecklistsScopeVersionIndex, true, "`scope`, `version`", "")
		if err := app.Save(checklists); err != nil {
			return err
		}
		// Raw SQL because hooks refuse changes to a saved checklist.
		_, err = app.DB().NewQuery("UPDATE request_checklists SET scope = 'creation'").Execute()
		return err
	}, func(app core.App) error {
		if _, err := app.DB().NewQuery("DELETE FROM request_checklists WHERE scope != 'creation'").Execute(); err != nil {
			return err
		}
		checklists, err := app.FindCollectionByNameOrId(requestChecklistsCollectionID)
		if err != nil {
			return err
		}
		checklists.Fields.RemoveByName("scope")
		checklists.RemoveIndex(requestChecklistsScopeVersionIndex)
		checklists.AddIndex(requestChecklistsVersionIndex, true, "`version`", "")
		return app.Save(checklists)
	})
}
