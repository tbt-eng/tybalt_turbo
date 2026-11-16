package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/types"
)

// The approver of a project request approves its proposed rates with a
// statement of their own, apart from the checklist (issue 187). The draft rate
// sheet records it:
//
//   - rates_approver: the user who approved the creation of the sheet and its
//     rates.
//   - rates_approved: when they did.
//
// Recall, rejection and a new submission clear both. They stay on the sheet
// after job creation activates it, when the request is gone.
//
// See descriptions/job_rate_sheet_history.md, Proposed rates.

const rateSheetsCreateRuleNoApproval = "@request.auth.id != \"\" &&\n// A draft is made and approved only through a project request.\n@request.body.request:isset = false &&\n@request.body.based_on:isset = false &&\n@request.body.reason:isset = false &&\n@request.body.rates_approver:isset = false &&\n@request.body.rates_approved:isset = false && (\n  @request.auth.user_claims_via_uid.cid.name ?= 'rate_sheet_revise' ||\n  @request.auth.user_claims_via_uid.cid.name ?= 'admin' ||\n  (@request.auth.user_claims_via_uid.cid.name ?= 'job' && @request.body.revision = 0)\n)"

func init() {
	m.Register(func(app core.App) error {
		users, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}
		sheets, err := app.FindCollectionByNameOrId("rate_sheets")
		if err != nil {
			return err
		}
		sheets.Fields.Add(
			&core.RelationField{Name: "rates_approver", CollectionId: users.Id, MaxSelect: 1},
			&core.DateField{Name: "rates_approved"},
		)
		sheets.CreateRule = types.Pointer(rateSheetsCreateRuleNoApproval)
		return app.Save(sheets)
	}, func(app core.App) error {
		sheets, err := app.FindCollectionByNameOrId("rate_sheets")
		if err != nil {
			return err
		}
		sheets.Fields.RemoveByName("rates_approver")
		sheets.Fields.RemoveByName("rates_approved")
		sheets.CreateRule = types.Pointer(rateSheetsCreateRuleNoDrafts)
		return app.Save(sheets)
	})
}
