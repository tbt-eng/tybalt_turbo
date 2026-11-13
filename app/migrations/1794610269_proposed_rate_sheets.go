package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/types"
)

// A project request can propose new rates (issue 187). The proposal is a draft
// rate sheet that belongs to the request:
//
//   - request: the request that owns the draft. A sheet with a request is a
//     draft. It is inactive, and the API lists it nowhere. Job creation
//     activates it and clears request; deleting the request deletes it.
//   - reason: why the requester needs a custom sheet. It stays after the
//     request is gone.
//   - based_on: the active sheet that the requester started from, which the
//     approver compares the proposed rates with.
//
// See descriptions/job_rate_sheet_history.md, Proposed rates.

const (
	rateSheetsReadAllRule          = "@request.auth.id != \"\""
	rateSheetsReadNoDraftsRule     = "@request.auth.id != \"\" && request = \"\""
	rateSheetEntriesReadNoDrafts   = "@request.auth.id != \"\" && rate_sheet.request = \"\""
	rateSheetsCreateRuleBefore     = "@request.auth.id != \"\" && (\n  @request.auth.user_claims_via_uid.cid.name ?= 'rate_sheet_revise' ||\n  @request.auth.user_claims_via_uid.cid.name ?= 'admin' ||\n  (@request.auth.user_claims_via_uid.cid.name ?= 'job' && @request.body.revision = 0)\n)"
	rateSheetsCreateRuleNoDrafts   = "@request.auth.id != \"\" &&\n// A draft is made only through a project request.\n@request.body.request:isset = false &&\n@request.body.based_on:isset = false &&\n@request.body.reason:isset = false && (\n  @request.auth.user_claims_via_uid.cid.name ?= 'rate_sheet_revise' ||\n  @request.auth.user_claims_via_uid.cid.name ?= 'admin' ||\n  (@request.auth.user_claims_via_uid.cid.name ?= 'job' && @request.body.revision = 0)\n)"
	rateSheetEntriesCreateBefore   = "@request.auth.id != \"\" &&\n@request.auth.user_claims_via_uid.cid.name ?= 'job'"
	rateSheetEntriesCreateNoDrafts = "@request.auth.id != \"\" &&\n@request.auth.user_claims_via_uid.cid.name ?= 'job' &&\n// The requester changes a draft through the project request.\n@request.body.rate_sheet.request = \"\""
)

// The Rate Sheets list leaves out drafts.
const rateSheetsAugmentedNoDraftsQuery = "SELECT \n  rs.id AS id,\n  rs.name AS name,\n  rs.effective_date AS effective_date,\n  rs.revision AS revision,\n  rs.active AS active,\n  rs.created AS created,\n  rs.updated AS updated,\n  COUNT(u.job) AS job_count\nFROM rate_sheets rs\nLEFT JOIN (\n  SELECT j.rate_sheet AS rate_sheet, j.id AS job FROM jobs j WHERE j.rate_sheet != ''\n  UNION\n  SELECT h.rate_sheet AS rate_sheet, h.job AS job FROM job_rate_sheets h\n) u ON u.rate_sheet = rs.id\nWHERE COALESCE(rs.request, '') = ''\nGROUP BY rs.id, rs.name, rs.effective_date, rs.revision, rs.active, rs.created, rs.updated"

func init() {
	m.Register(func(app core.App) error {
		requests, err := app.FindCollectionByNameOrId("job_requests")
		if err != nil {
			return err
		}
		sheets, err := app.FindCollectionByNameOrId("rate_sheets")
		if err != nil {
			return err
		}
		sheets.Fields.Add(
			&core.RelationField{Name: "request", CollectionId: requests.Id, MaxSelect: 1},
			&core.TextField{Name: "reason", Max: 500},
			&core.RelationField{Name: "based_on", CollectionId: sheets.Id, MaxSelect: 1},
		)
		sheets.ListRule = types.Pointer(rateSheetsReadNoDraftsRule)
		sheets.ViewRule = types.Pointer(rateSheetsReadNoDraftsRule)
		sheets.CreateRule = types.Pointer(rateSheetsCreateRuleNoDrafts)
		sheets.AddIndex("idx_rate_sheets_request", false, "`request`", "")
		if err := app.Save(sheets); err != nil {
			return err
		}

		entries, err := app.FindCollectionByNameOrId("rate_sheet_entries")
		if err != nil {
			return err
		}
		entries.ListRule = types.Pointer(rateSheetEntriesReadNoDrafts)
		entries.ViewRule = types.Pointer(rateSheetEntriesReadNoDrafts)
		entries.CreateRule = types.Pointer(rateSheetEntriesCreateNoDrafts)
		if err := app.Save(entries); err != nil {
			return err
		}

		return setRateSheetsAugmentedQuery(app, rateSheetsAugmentedNoDraftsQuery)
	}, func(app core.App) error {
		if err := setRateSheetsAugmentedQuery(app, rateSheetsAugmentedHistoryQuery); err != nil {
			return err
		}
		entries, err := app.FindCollectionByNameOrId("rate_sheet_entries")
		if err != nil {
			return err
		}
		entries.ListRule = types.Pointer(rateSheetsReadAllRule)
		entries.ViewRule = types.Pointer(rateSheetsReadAllRule)
		entries.CreateRule = types.Pointer(rateSheetEntriesCreateBefore)
		if err := app.Save(entries); err != nil {
			return err
		}
		sheets, err := app.FindCollectionByNameOrId("rate_sheets")
		if err != nil {
			return err
		}
		sheets.RemoveIndex("idx_rate_sheets_request")
		sheets.Fields.RemoveByName("request")
		sheets.Fields.RemoveByName("reason")
		sheets.Fields.RemoveByName("based_on")
		sheets.ListRule = types.Pointer(rateSheetsReadAllRule)
		sheets.ViewRule = types.Pointer(rateSheetsReadAllRule)
		sheets.CreateRule = types.Pointer(rateSheetsCreateRuleBefore)
		return app.Save(sheets)
	})
}
