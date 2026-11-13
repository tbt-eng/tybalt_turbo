package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/types"
)

// job_rate_sheets records which rate sheet prices a project's time from which
// date. A project has one opening row with an empty effective_date, which
// covers all time before its first dated row. See
// descriptions/job_rate_sheet_history.md.
func init() {
	m.Register(func(app core.App) error {
		jobs, err := app.FindCollectionByNameOrId("jobs")
		if err != nil {
			return err
		}
		rateSheets, err := app.FindCollectionByNameOrId("rate_sheets")
		if err != nil {
			return err
		}
		users, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}

		c := core.NewBaseCollection("job_rate_sheets")
		// Anyone signed in can read jobs, so anyone signed in can read the
		// history. Writes go through the history routes only.
		c.ListRule = types.Pointer("@request.auth.id != \"\"")
		c.ViewRule = types.Pointer("@request.auth.id != \"\"")
		c.Fields.Add(
			&core.RelationField{Name: "job", CollectionId: jobs.Id, Required: true, MaxSelect: 1, CascadeDelete: true},
			&core.RelationField{Name: "rate_sheet", CollectionId: rateSheets.Id, Required: true, MaxSelect: 1},
			// YYYY-MM-DD, the format of time_entries.date. Empty on the opening row.
			&core.TextField{Name: "effective_date", Pattern: `^(\d{4}-\d{2}-\d{2})?$`},
			&core.TextField{Name: "reason", Max: 500},
			&core.RelationField{Name: "created_by", CollectionId: users.Id, MaxSelect: 1},
			&core.RelationField{Name: "updated_by", CollectionId: users.Id, MaxSelect: 1},
			&core.AutodateField{Name: "created", OnCreate: true},
			&core.AutodateField{Name: "updated", OnCreate: true, OnUpdate: true},
		)
		c.AddIndex("idx_job_rate_sheets_job_date", true, "`job`, `effective_date`", "")
		c.AddIndex("idx_job_rate_sheets_opening", true, "`job`", "`effective_date` = ''")
		if err := app.Save(c); err != nil {
			return err
		}

		// rate_sheet_revise now governs the history instead of jobs.rate_sheet.
		if _, err := app.DB().NewQuery(`UPDATE claims SET description = {:description} WHERE name = 'rate_sheet_revise'`).
			Bind(dbx.Params{"description": "Can create a new revision of a rate sheet, and change the rate sheet history of a project that already has a sheet. To change the history, the holder also needs the job claim."}).Execute(); err != nil {
			return err
		}

		// Each project with a rate sheet gets an opening row on its current
		// sheet, so pricing is unchanged by the migration. The ids come from
		// PocketBase, so they have the same form as every other record id
		// (15 characters from a-z and 0-9).
		var projects []struct {
			Id        string `db:"id"`
			RateSheet string `db:"rate_sheet"`
		}
		if err := app.DB().NewQuery(`
			SELECT j.id, j.rate_sheet
			FROM jobs j
			WHERE j.number NOT LIKE 'P%' AND COALESCE(j.rate_sheet, '') != ''
			ORDER BY j.id
		`).All(&projects); err != nil {
			return err
		}
		insert := app.DB().NewQuery(`
			INSERT INTO job_rate_sheets (id, job, rate_sheet, effective_date, reason, created_by, updated_by, created, updated)
			VALUES ({:id}, {:job}, {:rate_sheet}, '', '', '', '',
				strftime('%Y-%m-%d %H:%M:%fZ', 'now'), strftime('%Y-%m-%d %H:%M:%fZ', 'now'))
		`)
		for _, project := range projects {
			if _, err := insert.Bind(dbx.Params{
				"id":         core.GenerateDefaultRandomId(),
				"job":        project.Id,
				"rate_sheet": project.RateSheet,
			}).Execute(); err != nil {
				return err
			}
		}
		return nil
	}, func(app core.App) error {
		c, err := app.FindCollectionByNameOrId("job_rate_sheets")
		if err != nil {
			return err
		}
		return app.Delete(c)
	})
}
