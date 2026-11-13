package hooks

import (
	"fmt"
	"sort"

	"tybalt/absorb"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

// Absorbing rate sheets (issue #191). The rule behind every check here: an
// absorb never reprices time. The jobs and project requests of an absorbed
// sheet move to the target, so the two sheets must have identical rates.
//
// Which sheets can be absorbed is decided on the name, never on the one
// revision that the user selects:
//   - A sheet that is the only revision of its name can be absorbed when its
//     rates are identical to the target's. The target's effective date must
//     not be later than the sheet's, because its jobs move.
//   - A name with revisions can be absorbed only when no job and no project
//     request uses any of its revisions. Then every revision is absorbed
//     together, and only the latest revision must match the target. Nothing
//     moves, so the date rule does not apply.
//   - Otherwise no revision of the name can be absorbed.
//
// Any revision of any sheet can be the target. An inactive target can't absorb
// an active sheet, because absorb does not change active and the active sheet
// would leave the Rate Sheet selector.

func init() {
	absorb.Register("rate_sheets", absorb.Config{
		RefConfigs: []absorb.RefConfig{
			{Table: "jobs", Column: "rate_sheet"},
			{Table: "job_requests", Column: "rate_sheet"},
		},
		// rate_sheet_name_changes also relates to rate_sheets. Its rows move
		// to the target on commit (commitRateSheetAbsorb), not on absorb.
		OwnedRows: []absorb.RefConfig{
			{Table: "rate_sheet_entries", Column: "rate_sheet"},
		},
		// An absorb moves jobs between contract rate sheets, so absorb alone
		// is not enough.
		ClaimSets: [][]string{
			{"admin"},
			{"absorb", "rate_sheet_revise"},
		},
		ChangesJobs: true,
		Prepare:     prepareRateSheetAbsorb,
		OnCommit:    commitRateSheetAbsorb,
	})
}

// Reasons why a sheet can't be absorbed. They are shown to the user.
const (
	rateSheetReasonSelf          = "A sheet can't absorb itself."
	rateSheetReasonSameName      = "It is a revision of the target sheet."
	rateSheetReasonTargetNoRates = "The target has no rates."
	rateSheetReasonRates         = "Its rates are not identical to the target's."
	rateSheetReasonLatestRates   = "Its latest revision's rates are not identical to the target's."
	rateSheetReasonNameInUse     = "Its name has revisions, and a job or project request uses one of them."
	rateSheetReasonInactive      = "An inactive target can't absorb an active sheet."
)

// RateSheetAbsorbSheet is one rate sheet as the absorb rules see it.
type RateSheetAbsorbSheet struct {
	ID            string `json:"id" db:"id"`
	Name          string `json:"name" db:"name"`
	Revision      int    `json:"revision" db:"revision"`
	EffectiveDate string `json:"effective_date" db:"effective_date"`
	Active        bool   `json:"active" db:"active"`
	JobCount      int    `json:"job_count" db:"job_count"`
	RequestCount  int    `json:"request_count" db:"request_count"`
	// Revisions is the number of revisions of the sheet's name.
	Revisions int `json:"revisions" db:"-"`
	ratesKey  string
}

// rateSheetAbsorbData holds every rate sheet with its rates and references,
// so that the rules can be checked for many pairs without more queries.
type rateSheetAbsorbData struct {
	sheets map[string]*RateSheetAbsorbSheet
	// byName holds the revisions of each name, oldest first.
	byName map[string][]*RateSheetAbsorbSheet
}

func loadRateSheetAbsorbData(app core.App) (*rateSheetAbsorbData, error) {
	rows := []*RateSheetAbsorbSheet{}
	err := app.DB().NewQuery(`
		SELECT
			rs.id AS id,
			rs.name AS name,
			rs.revision AS revision,
			rs.effective_date AS effective_date,
			rs.active AS active,
			(SELECT COUNT(*) FROM jobs j WHERE j.rate_sheet = rs.id) AS job_count,
			(SELECT COUNT(*) FROM job_requests r WHERE r.rate_sheet = rs.id) AS request_count
		FROM rate_sheets rs
	`).All(&rows)
	if err != nil {
		return nil, err
	}

	entryRows := []struct {
		RateSheet    string  `db:"rate_sheet"`
		Role         string  `db:"role"`
		Rate         float64 `db:"rate"`
		OvertimeRate float64 `db:"overtime_rate"`
	}{}
	err = app.DB().NewQuery(`
		SELECT rate_sheet, role, rate, overtime_rate FROM rate_sheet_entries
	`).All(&entryRows)
	if err != nil {
		return nil, err
	}
	entries := map[string][]RateSheetEntryValues{}
	for _, row := range entryRows {
		entries[row.RateSheet] = append(entries[row.RateSheet], RateSheetEntryValues{
			Role:         row.Role,
			Rate:         row.Rate,
			OvertimeRate: row.OvertimeRate,
		})
	}

	data := &rateSheetAbsorbData{
		sheets: map[string]*RateSheetAbsorbSheet{},
		byName: map[string][]*RateSheetAbsorbSheet{},
	}
	for _, sheet := range rows {
		sheet.ratesKey = RateSheetRatesKey(entries[sheet.ID])
		data.sheets[sheet.ID] = sheet
		data.byName[sheet.Name] = append(data.byName[sheet.Name], sheet)
	}
	for _, revisions := range data.byName {
		sort.Slice(revisions, func(i, j int) bool { return revisions[i].Revision < revisions[j].Revision })
		for _, sheet := range revisions {
			sheet.Revisions = len(revisions)
		}
	}
	return data, nil
}

// check returns why target can't absorb source, or "" and the ids that the
// absorb deletes: source alone, or every revision of its name.
func (d *rateSheetAbsorbData) check(target, source *RateSheetAbsorbSheet) (string, []string) {
	if source.ID == target.ID {
		return rateSheetReasonSelf, nil
	}
	if source.Name == target.Name {
		return rateSheetReasonSameName, nil
	}
	if target.ratesKey == "" {
		return rateSheetReasonTargetNoRates, nil
	}

	revisions := d.byName[source.Name]
	if len(revisions) == 1 {
		if source.ratesKey != target.ratesKey {
			return rateSheetReasonRates, nil
		}
		// The jobs of the sheet move, so the target must have been in effect
		// at least as early. Dates are YYYY-MM-DD, so strings compare.
		if target.EffectiveDate > source.EffectiveDate {
			return fmt.Sprintf("The target's effective date, %s, is later than this sheet's, %s.", target.EffectiveDate, source.EffectiveDate), nil
		}
		if !target.Active && source.Active {
			return rateSheetReasonInactive, nil
		}
		return "", []string{source.ID}
	}

	ids := make([]string, 0, len(revisions))
	anyActive := false
	for _, revision := range revisions {
		if revision.JobCount > 0 || revision.RequestCount > 0 {
			return rateSheetReasonNameInUse, nil
		}
		anyActive = anyActive || revision.Active
		ids = append(ids, revision.ID)
	}
	// The older revisions are not referenced, so their rates change nothing.
	if revisions[len(revisions)-1].ratesKey != target.ratesKey {
		return rateSheetReasonLatestRates, nil
	}
	if !target.Active && anyActive {
		return rateSheetReasonInactive, nil
	}
	return "", ids
}

// prepareRateSheetAbsorb is the absorb.Config Prepare of rate_sheets. It
// refuses any sheet that target can't absorb, and returns the ids to absorb,
// with every revision of a name that is absorbed as a whole.
func prepareRateSheetAbsorb(app core.App, target *core.Record, ids []string) ([]string, error) {
	data, err := loadRateSheetAbsorbData(app)
	if err != nil {
		return nil, err
	}
	targetSheet, ok := data.sheets[target.Id]
	if !ok {
		return nil, absorb.Refused("the target rate sheet was not found")
	}

	result := []string{}
	seen := map[string]bool{}
	for _, id := range ids {
		source, ok := data.sheets[id]
		if !ok {
			return nil, absorb.Refused("rate sheet %s was not found", id)
		}
		reason, absorbIDs := data.check(targetSheet, source)
		if reason != "" {
			return nil, absorb.Refused("%s (rev. %d) can't be absorbed into %s (rev. %d). %s",
				source.Name, source.Revision, targetSheet.Name, targetSheet.Revision, reason)
		}
		for _, absorbID := range absorbIDs {
			if !seen[absorbID] {
				seen[absorbID] = true
				result = append(result, absorbID)
			}
		}
	}
	return result, nil
}

// commitRateSheetAbsorb is the absorb.Config OnCommit of rate_sheets. Each
// absorbed name goes into the rate sheet name history once, so the target
// shows it after "Formerly". It is written on commit, not on absorb, so an
// undo has nothing to remove. The history rows of the absorbed sheets move to
// the target too, so their earlier names are kept.
func commitRateSheetAbsorb(app core.App, action *core.Record, auth *core.Record) error {
	if auth == nil || auth.Collection().Name != "users" {
		return absorb.Refused("a rate sheet absorb must be committed by a user")
	}
	target, err := app.FindRecordById("rate_sheets", action.GetString("target_id"))
	if err != nil {
		return fmt.Errorf("finding the target rate sheet: %w", err)
	}
	targetName := target.GetString("name")

	var absorbed []map[string]any
	if err := action.UnmarshalJSONField("absorbed_records", &absorbed); err != nil {
		return fmt.Errorf("reading absorbed records: %w", err)
	}
	ids := []any{}
	names := []string{}
	seen := map[string]bool{targetName: true}
	for _, record := range absorbed {
		if id, ok := record["id"].(string); ok && id != "" {
			ids = append(ids, id)
		}
		name, _ := record["name"].(string)
		if name != "" && !seen[name] {
			seen[name] = true
			names = append(names, name)
		}
	}

	if len(ids) > 0 {
		_, err = app.DB().Update(
			"rate_sheet_name_changes",
			dbx.Params{"rate_sheet": target.Id},
			dbx.In("rate_sheet", ids...),
		).Execute()
		if err != nil {
			return fmt.Errorf("moving the name history of absorbed rate sheets: %w", err)
		}
	}

	changes, err := app.FindCollectionByNameOrId("rate_sheet_name_changes")
	if err != nil {
		return err
	}
	for _, name := range names {
		change := core.NewRecord(changes)
		change.Set("rate_sheet", target.Id)
		change.Set("old_name", name)
		change.Set("new_name", targetName)
		change.Set("renamed_by", auth.Id)
		if err := app.Save(change); err != nil {
			return fmt.Errorf("saving the name history: %w", err)
		}
	}
	return nil
}

// RateSheetDuplicate is a sheet in a group of identical sheets.
type RateSheetDuplicate struct {
	RateSheetAbsorbSheet
	// Target marks the group's target: the sheet with the earliest
	// effective date.
	Target bool `json:"target"`
	// CanAbsorb is true when the sheet can absorb at least one other sheet
	// of its group.
	CanAbsorb bool `json:"can_absorb"`
	// Reason says why the group's target can't absorb this sheet. It is
	// empty for the target and for a sheet that it can absorb.
	Reason string `json:"reason,omitempty"`
	// Absorbs is the number of sheets that an absorb into the target
	// deletes: 1, or every revision of the name.
	Absorbs int `json:"absorbs,omitempty"`
}

// RateSheetDuplicateGroup is a group of sheets with identical rates, from at
// least two names, active and inactive.
type RateSheetDuplicateGroup struct {
	Sheets []RateSheetDuplicate `json:"sheets"`
	// Reason says why no sheet of the group can be absorbed into its target.
	Reason string `json:"reason,omitempty"`
}

// FindRateSheetDuplicates returns every group of sheets whose rates are
// identical. Sheets without entries are not in any group. The target of a
// group is first; it is the sheet with the earliest effective date, and an
// active sheet, then the sheet with the most jobs, goes first on a tie.
func FindRateSheetDuplicates(app core.App) ([]RateSheetDuplicateGroup, error) {
	data, err := loadRateSheetAbsorbData(app)
	if err != nil {
		return nil, err
	}

	byKey := map[string][]*RateSheetAbsorbSheet{}
	for _, sheet := range data.sheets {
		if sheet.ratesKey != "" {
			byKey[sheet.ratesKey] = append(byKey[sheet.ratesKey], sheet)
		}
	}

	groups := []RateSheetDuplicateGroup{}
	for _, sheets := range byKey {
		names := map[string]bool{}
		for _, sheet := range sheets {
			names[sheet.Name] = true
		}
		// Revisions of one name are not duplicates of each other.
		if len(names) < 2 {
			continue
		}
		sort.Slice(sheets, func(i, j int) bool {
			a, b := sheets[i], sheets[j]
			if a.EffectiveDate != b.EffectiveDate {
				return a.EffectiveDate < b.EffectiveDate
			}
			if a.Active != b.Active {
				return a.Active
			}
			if a.JobCount != b.JobCount {
				return a.JobCount > b.JobCount
			}
			if a.Name != b.Name {
				return a.Name < b.Name
			}
			return a.Revision < b.Revision
		})

		target := sheets[0]
		group := RateSheetDuplicateGroup{}
		absorbable := 0
		inactiveRefusal := false
		for _, sheet := range sheets {
			duplicate := RateSheetDuplicate{RateSheetAbsorbSheet: *sheet, Target: sheet == target}
			for _, other := range sheets {
				if reason, _ := data.check(sheet, other); reason == "" {
					duplicate.CanAbsorb = true
					break
				}
			}
			if sheet != target {
				reason, ids := data.check(target, sheet)
				duplicate.Reason = reason
				duplicate.Absorbs = len(ids)
				if reason == "" {
					absorbable++
				}
				inactiveRefusal = inactiveRefusal || reason == rateSheetReasonInactive
			}
			group.Sheets = append(group.Sheets, duplicate)
		}

		if absorbable == 0 {
			group.Reason = "No sheet in this group can be absorbed into the earliest sheet."
			if inactiveRefusal {
				newer := data.byName[target.Name]
				if newer[len(newer)-1].Revision > target.Revision {
					group.Reason = "The earliest sheet is inactive and has a newer revision. It can't absorb an active sheet, and it can't be activated."
				} else {
					group.Reason = "The earliest sheet is inactive and can't absorb an active sheet. Activate it first."
				}
			}
		}
		groups = append(groups, group)
	}

	sort.Slice(groups, func(i, j int) bool {
		a, b := groups[i].Sheets[0], groups[j].Sheets[0]
		if a.Name != b.Name {
			return a.Name < b.Name
		}
		return a.Revision < b.Revision
	})
	return groups, nil
}

// RateSheetAbsorbCandidate is a sheet with rates identical to a target, and
// whether the target can absorb it.
type RateSheetAbsorbCandidate struct {
	RateSheetAbsorbSheet
	// Reason says why the target can't absorb the sheet; empty when it can.
	Reason string `json:"reason,omitempty"`
	// Absorbs is the number of sheets that the absorb deletes.
	Absorbs int `json:"absorbs,omitempty"`
}

// RateSheetAbsorbCandidates returns the target and every sheet of another
// name with identical rates, with the reason when the target can't absorb
// it. ok is false when there is no sheet with id targetID.
func RateSheetAbsorbCandidates(app core.App, targetID string) (RateSheetAbsorbSheet, []RateSheetAbsorbCandidate, bool, error) {
	data, err := loadRateSheetAbsorbData(app)
	if err != nil {
		return RateSheetAbsorbSheet{}, nil, false, err
	}
	target, ok := data.sheets[targetID]
	if !ok {
		return RateSheetAbsorbSheet{}, nil, false, nil
	}

	candidates := []RateSheetAbsorbCandidate{}
	if target.ratesKey != "" {
		for _, sheet := range data.sheets {
			if sheet.Name == target.Name || sheet.ratesKey != target.ratesKey {
				continue
			}
			reason, ids := data.check(target, sheet)
			candidates = append(candidates, RateSheetAbsorbCandidate{
				RateSheetAbsorbSheet: *sheet,
				Reason:               reason,
				Absorbs:              len(ids),
			})
		}
	}
	sort.Slice(candidates, func(i, j int) bool {
		a, b := candidates[i], candidates[j]
		if a.Name != b.Name {
			return a.Name < b.Name
		}
		return a.Revision < b.Revision
	})
	return *target, candidates, true, nil
}
