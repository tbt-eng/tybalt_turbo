package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// Make jobs.contact and jobs.manager optional in the schema (issue #167).
//
// Many imported jobs have no contact or manager. A schema `required` flag is
// checked on every save, after the hooks, so it rejected even a status-only
// close of those jobs ("cannot be blank"). The requirement now lives in
// validateJob (hooks/jobs.go), which enforces it on creates and full edits and
// skips it for status-only updates.
//
// An optional relation lets PocketBase clear the reference when the referenced
// record is deleted, where a required one refused the delete. The delete hooks
// keep refusing instead: recordInUse for contacts, and the users delete hook in
// AddProjectWorkflowHooks for managers.
//
// The down migration restores the flags. PocketBase doesn't recheck existing
// rows when a collection's schema is saved, so it succeeds even with blank
// values present, but those jobs can no longer be saved afterwards.
func init() {
	m.Register(func(app core.App) error {
		return setJobsContactManagerRequired(app, false)
	}, func(app core.App) error {
		return setJobsContactManagerRequired(app, true)
	})
}

func setJobsContactManagerRequired(app core.App, required bool) error {
	jobs, err := app.FindCollectionByNameOrId("jobs")
	if err != nil {
		return err
	}
	for _, name := range []string{"contact", "manager"} {
		jobs.Fields.GetByName(name).(*core.RelationField).Required = required
	}
	return app.Save(jobs)
}
