package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/types"
)

// The collection API creates only a new sheet (revision 0), and only for a
// holder of the job claim, as the create route does (issue 207). Since rate
// sheet codes, a create there can't be a revision: it can't send a code, and
// a revision needs the code of the sheet it revises. The rule still let a
// holder of rate_sheet_revise or admin without the job claim create a new
// sheet, which the route refuses. Revisions go through POST /api/rate_sheets.
//
// See descriptions/job_rate_sheet_history.md, Rate sheet codes.

const rateSheetsCreateRuleJobOnly = "@request.auth.id != \"\" &&\n// A draft is made and approved only through a project request.\n@request.body.request:isset = false &&\n@request.body.based_on:isset = false &&\n@request.body.reason:isset = false &&\n@request.body.rates_approver:isset = false &&\n@request.body.rates_approved:isset = false &&\n// The server gives a new sheet its code.\n@request.body.code:isset = false &&\n// Only a new sheet, by Procurement. A revision goes through the create route.\n@request.auth.user_claims_via_uid.cid.name ?= 'job' &&\n@request.body.revision = 0"

func init() {
	m.Register(func(app core.App) error {
		return setRateSheetsCreateRule(app, rateSheetsCreateRuleJobOnly)
	}, func(app core.App) error {
		return setRateSheetsCreateRule(app, rateSheetsCreateRuleNoCode)
	})
}

func setRateSheetsCreateRule(app core.App, rule string) error {
	sheets, err := app.FindCollectionByNameOrId("rate_sheets")
	if err != nil {
		return err
	}
	sheets.CreateRule = types.Pointer(rule)
	return app.Save(sheets)
}
