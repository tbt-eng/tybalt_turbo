package migrations

import (
	"fmt"
	"regexp"
	"sort"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/types"
)

// Each rate sheet gets a code, such as 2026-11-00004 (issue 207). The code is
// the sheet's identity: every revision of a sheet has the same code, and the
// unique index on code and revision replaces the one on name and revision.
// The month is the creation month (UTC) of the sheet's first revision, and
// the serial runs from 00001 in each month. The name becomes an optional
// description.
//
//   - rate_sheet_codes holds every code ever issued, one row each. It is the
//     counter: the next serial of a month is one more than the highest issued
//     in that month. Its rows are never deleted, so a code is never issued
//     twice, even after its sheet is deleted or absorbed. Only server code
//     and superusers can see it.
//   - rate_sheet_code_changes logs every edit of a code by a superuser: the
//     sheet, the old code, the new code and the time. A superuser is not a
//     users record, so the log has no user. The sheet relation is optional,
//     so a deleted sheet leaves its log rows.
//   - rate_sheet_name_changes no longer requires old_name and new_name, since
//     a sheet may have no name.
//
// Existing sheets get codes from their created date (see
// BackfillRateSheetCodes). The update writes code with SQL and does not touch
// updated, so the review token of a project request does not change.
//
// See descriptions/job_rate_sheet_history.md, Rate sheet codes.

const (
	rateSheetCodesCollectionID       = "pbc_rs_codes"
	rateSheetCodeChangesCollectionID = "pbc_rs_codechg"
	rateSheetNameRevisionIndex       = "idx_dHsmkOUme0"
	rateSheetCodeRevisionIndex       = "idx_rate_sheets_code_revision"
	// RateSheetCodePattern is the form of a code: YYYY-MM- and a five digit
	// serial.
	RateSheetCodePattern = `^\d{4}-(0[1-9]|1[0-2])-\d{5}$`
)

const rateSheetsCreateRuleNoCode = "@request.auth.id != \"\" &&\n// A draft is made and approved only through a project request.\n@request.body.request:isset = false &&\n@request.body.based_on:isset = false &&\n@request.body.reason:isset = false &&\n@request.body.rates_approver:isset = false &&\n@request.body.rates_approved:isset = false &&\n// The server gives a new sheet its code.\n@request.body.code:isset = false && (\n  @request.auth.user_claims_via_uid.cid.name ?= 'rate_sheet_revise' ||\n  @request.auth.user_claims_via_uid.cid.name ?= 'admin' ||\n  (@request.auth.user_claims_via_uid.cid.name ?= 'job' && @request.body.revision = 0)\n)"

// The Rate Sheets list shows the code.
const rateSheetsAugmentedCodeQuery = "SELECT \n  rs.id AS id,\n  rs.code AS code,\n  rs.name AS name,\n  rs.effective_date AS effective_date,\n  rs.revision AS revision,\n  rs.active AS active,\n  rs.created AS created,\n  rs.updated AS updated,\n  COUNT(u.job) AS job_count\nFROM rate_sheets rs\nLEFT JOIN (\n  SELECT j.rate_sheet AS rate_sheet, j.id AS job FROM jobs j WHERE j.rate_sheet != ''\n  UNION\n  SELECT h.rate_sheet AS rate_sheet, h.job AS job FROM job_rate_sheets h\n) u ON u.rate_sheet = rs.id\nWHERE COALESCE(rs.request, '') = ''\nGROUP BY rs.id, rs.code, rs.name, rs.effective_date, rs.revision, rs.active, rs.created, rs.updated"

const (
	rateSheetReviseDescriptionBefore = "Can create a new revision of a rate sheet, and change the rate sheet history of a project that already has a sheet. To change the history, the holder also needs the job claim."
	rateSheetReviseDescriptionCodes  = "Can create a new revision of a rate sheet, rename a rate sheet, and change the rate sheet history of a project that already has a sheet. To change the history, the holder also needs the job claim."
)

func init() {
	m.Register(func(app core.App) error {
		sheets, err := app.FindCollectionByNameOrId("rate_sheets")
		if err != nil {
			return err
		}
		sheets.Fields.Add(&core.TextField{Name: "code", Required: true, Pattern: RateSheetCodePattern})
		if name, ok := sheets.Fields.GetByName("name").(*core.TextField); ok {
			name.Required = false
		}
		sheets.RemoveIndex(rateSheetNameRevisionIndex)
		sheets.CreateRule = types.Pointer(rateSheetsCreateRuleNoCode)
		if err := app.Save(sheets); err != nil {
			return err
		}

		codes := core.NewBaseCollection("rate_sheet_codes")
		codes.Id = rateSheetCodesCollectionID
		codes.Fields.Add(
			&core.TextField{Name: "code", Required: true, Pattern: RateSheetCodePattern},
			&core.TextField{Name: "month", Required: true, Pattern: `^\d{4}-(0[1-9]|1[0-2])$`},
			&core.NumberField{Name: "serial", Required: true, OnlyInt: true, Min: types.Pointer(1.0), Max: types.Pointer(99999.0)},
			&core.AutodateField{Name: "created", OnCreate: true},
		)
		codes.AddIndex("idx_rate_sheet_codes_code", true, "`code`", "")
		codes.AddIndex("idx_rate_sheet_codes_month_serial", true, "`month`, `serial`", "")
		if err := app.Save(codes); err != nil {
			return err
		}

		changes := core.NewBaseCollection("rate_sheet_code_changes")
		changes.Id = rateSheetCodeChangesCollectionID
		changes.ListRule = types.Pointer("@request.auth.id != ''")
		changes.ViewRule = types.Pointer("@request.auth.id != ''")
		changes.Fields.Add(
			&core.RelationField{Name: "rate_sheet", CollectionId: sheets.Id, MaxSelect: 1},
			&core.TextField{Name: "old_code", Required: true},
			&core.TextField{Name: "new_code", Required: true},
			&core.AutodateField{Name: "created", OnCreate: true},
		)
		changes.AddIndex("idx_rate_sheet_code_changes_rate_sheet", false, "`rate_sheet`", "")
		if err := app.Save(changes); err != nil {
			return err
		}

		nameChanges, err := app.FindCollectionByNameOrId(rateSheetNameChangesCollectionID)
		if err != nil {
			return err
		}
		for _, field := range []string{"old_name", "new_name"} {
			if text, ok := nameChanges.Fields.GetByName(field).(*core.TextField); ok {
				text.Required = false
			}
		}
		if err := app.Save(nameChanges); err != nil {
			return err
		}

		if err := BackfillRateSheetCodes(app); err != nil {
			return err
		}

		sheets, err = app.FindCollectionByNameOrId("rate_sheets")
		if err != nil {
			return err
		}
		sheets.AddIndex(rateSheetCodeRevisionIndex, true, "`code`, `revision`", "")
		if err := app.Save(sheets); err != nil {
			return err
		}

		if _, err := app.DB().NewQuery(`UPDATE claims SET description = {:description} WHERE name = 'rate_sheet_revise'`).
			Bind(dbx.Params{"description": rateSheetReviseDescriptionCodes}).Execute(); err != nil {
			return err
		}

		return setRateSheetsAugmentedQuery(app, rateSheetsAugmentedCodeQuery)
	}, func(app core.App) error {
		if err := setRateSheetsAugmentedQuery(app, rateSheetsAugmentedNoDraftsQuery); err != nil {
			return err
		}
		if _, err := app.DB().NewQuery(`UPDATE claims SET description = {:description} WHERE name = 'rate_sheet_revise'`).
			Bind(dbx.Params{"description": rateSheetReviseDescriptionBefore}).Execute(); err != nil {
			return err
		}
		for _, id := range []string{rateSheetCodeChangesCollectionID, rateSheetCodesCollectionID} {
			if collection, err := app.FindCollectionByNameOrId(id); err == nil {
				if err := app.Delete(collection); err != nil {
					return err
				}
			}
		}
		nameChanges, err := app.FindCollectionByNameOrId(rateSheetNameChangesCollectionID)
		if err != nil {
			return err
		}
		for _, field := range []string{"old_name", "new_name"} {
			if text, ok := nameChanges.Fields.GetByName(field).(*core.TextField); ok {
				text.Required = true
			}
		}
		if err := app.Save(nameChanges); err != nil {
			return err
		}
		sheets, err := app.FindCollectionByNameOrId("rate_sheets")
		if err != nil {
			return err
		}
		sheets.RemoveIndex(rateSheetCodeRevisionIndex)
		sheets.Fields.RemoveByName("code")
		if name, ok := sheets.Fields.GetByName("name").(*core.TextField); ok {
			name.Required = true
		}
		sheets.AddIndex(rateSheetNameRevisionIndex, true, "`name`, `revision`", "")
		sheets.CreateRule = types.Pointer(rateSheetsCreateRuleNoApproval)
		return app.Save(sheets)
	})
}

var rateSheetCodeMonthPattern = regexp.MustCompile(`^(\d{4}-(0[1-9]|1[0-2]))`)

// rateSheetCodeMonth returns the YYYY-MM at the start of a date or timestamp,
// or "" when it has none.
func rateSheetCodeMonth(value string) string {
	match := rateSheetCodeMonthPattern.FindStringSubmatch(value)
	if match == nil {
		return ""
	}
	return match[1]
}

// RateSheetCodeBackfillRow is one rate sheet as BackfillRateSheetCodes reads
// it.
type RateSheetCodeBackfillRow struct {
	ID            string `db:"id"`
	Name          string `db:"name"`
	Revision      int    `db:"revision"`
	Created       string `db:"created"`
	EffectiveDate string `db:"effective_date"`
	Request       string `db:"request"`
}

// RateSheetCodeAssignment is a code given to a sheet by the backfill. Month
// and Serial make up Code.
type RateSheetCodeAssignment struct {
	Code   string
	Month  string
	Serial int
	IDs    []string
}

// PlanRateSheetCodes groups sheets into families and gives each family a
// code. Revisions were linked by name, so the sheets of one name are one
// family, except a draft of a project request, which is always a family of
// its own. The month is that of the created date of the family's first
// revision (the lowest revision). When it has none, the earliest created date
// of the family is used, then the first revision's effective date, then now.
// Within a month, families are numbered in the order of that date, then the
// first revision's id.
func PlanRateSheetCodes(rows []RateSheetCodeBackfillRow, now time.Time) []RateSheetCodeAssignment {
	type family struct {
		first   RateSheetCodeBackfillRow
		month   string
		sortKey string
		ids     []string
		created string
	}
	families := map[string]*family{}
	keys := []string{}
	for _, row := range rows {
		key := "name:" + row.Name
		if row.Request != "" {
			key = "draft:" + row.ID
		}
		f, ok := families[key]
		if !ok {
			f = &family{first: row}
			families[key] = f
			keys = append(keys, key)
		} else if row.Revision < f.first.Revision || (row.Revision == f.first.Revision && row.ID < f.first.ID) {
			f.first = row
		}
		f.ids = append(f.ids, row.ID)
		if month := rateSheetCodeMonth(row.Created); month != "" && (f.created == "" || row.Created < f.created) {
			f.created = row.Created
		}
	}

	byMonth := map[string][]*family{}
	months := []string{}
	for _, key := range keys {
		f := families[key]
		date := f.first.Created
		if rateSheetCodeMonth(date) == "" {
			date = f.created
		}
		if rateSheetCodeMonth(date) == "" {
			date = f.first.EffectiveDate
		}
		if rateSheetCodeMonth(date) == "" {
			date = now.UTC().Format("2006-01-02 15:04:05.000Z")
		}
		f.month = rateSheetCodeMonth(date)
		f.sortKey = date
		sort.Strings(f.ids)
		if _, ok := byMonth[f.month]; !ok {
			months = append(months, f.month)
		}
		byMonth[f.month] = append(byMonth[f.month], f)
	}
	sort.Strings(months)

	result := []RateSheetCodeAssignment{}
	for _, month := range months {
		group := byMonth[month]
		sort.SliceStable(group, func(i, j int) bool {
			if group[i].sortKey != group[j].sortKey {
				return group[i].sortKey < group[j].sortKey
			}
			return group[i].first.ID < group[j].first.ID
		})
		for i, f := range group {
			serial := i + 1
			result = append(result, RateSheetCodeAssignment{
				Code:   fmt.Sprintf("%s-%05d", month, serial),
				Month:  month,
				Serial: serial,
				IDs:    f.ids,
			})
		}
	}
	return result
}

// BackfillRateSheetCodes gives every rate sheet without a code a code (see
// PlanRateSheetCodes) and records each code in rate_sheet_codes. It expects
// no sheet to have a code yet and rate_sheet_codes to be empty, as when the
// migration runs. It writes code with SQL, so updated does not change and no
// hook runs.
func BackfillRateSheetCodes(app core.App) error {
	rows := []RateSheetCodeBackfillRow{}
	if err := app.DB().NewQuery(`
		SELECT id, COALESCE(name, '') AS name, COALESCE(revision, 0) AS revision,
			COALESCE(created, '') AS created, COALESCE(effective_date, '') AS effective_date,
			COALESCE(request, '') AS request
		FROM rate_sheets
		WHERE COALESCE(code, '') = ''
	`).All(&rows); err != nil {
		return err
	}
	if len(rows) == 0 {
		return nil
	}
	registry, err := app.FindCollectionByNameOrId("rate_sheet_codes")
	if err != nil {
		return err
	}
	for _, assignment := range PlanRateSheetCodes(rows, time.Now()) {
		for _, id := range assignment.IDs {
			if _, err := app.DB().NewQuery(`UPDATE rate_sheets SET code = {:code} WHERE id = {:id}`).
				Bind(dbx.Params{"code": assignment.Code, "id": id}).Execute(); err != nil {
				return err
			}
		}
		issued := core.NewRecord(registry)
		issued.Set("code", assignment.Code)
		issued.Set("month", assignment.Month)
		issued.Set("serial", assignment.Serial)
		if err := app.Save(issued); err != nil {
			return fmt.Errorf("recording rate sheet code %s: %w", assignment.Code, err)
		}
	}
	return nil
}
