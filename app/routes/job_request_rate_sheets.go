package routes

import (
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"tybalt/hooks"

	"github.com/pocketbase/dbx"
	validation "github.com/pocketbase/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
)

// A project request either chooses an active rate sheet or proposes new rates.
// Proposed rates are saved as a draft rate sheet that the request owns (see
// hooks/rate_sheet_drafts.go). Only the request's creator changes the draft,
// and only through the request save, under the same edit rules as the rest of
// the request. A draft is never listed or used elsewhere until job creation
// activates it.

// proposedRatesWrite is the proposed_rates object of a request save. Name and
// entries are the new sheet's; based_on is the active sheet the requester
// started from; reason says why the project needs new rates.
type proposedRatesWrite struct {
	Name    string                 `json:"name"`
	Reason  string                 `json:"reason"`
	BasedOn string                 `json:"based_on"`
	Entries []CreateRateSheetEntry `json:"entries"`
}

// proposedRatesPlan is a checked proposal, ready to save in the request's
// transaction.
type proposedRatesPlan struct {
	draft   *core.Record // nil for a new draft
	input   proposedRatesWrite
	changed bool
}

// proposedRatesReasonMax is the length of rate_sheets.reason.
const proposedRatesReasonMax = 500

// requestDraft returns the draft that the request owns, or nil.
func requestDraft(app core.App, r *core.Record) (*core.Record, error) {
	drafts, err := hooks.RequestRateSheetDrafts(app, r.Id)
	if err != nil || len(drafts) == 0 {
		return nil, err
	}
	return drafts[0], nil
}

// checkChosenRateSheet refuses a rate sheet chosen in the request that is a
// draft, unless it is the request's own draft.
func checkChosenRateSheet(app core.App, id string, own *core.Record) error {
	if id == "" || (own != nil && own.Id == id) {
		return nil
	}
	sheet, err := app.FindRecordById("rate_sheets", id)
	if err != nil {
		return workflowError("Select a valid rate sheet.")
	}
	if hooks.IsRateSheetDraft(sheet) {
		return workflowError("Select an active rate sheet.")
	}
	return nil
}

// prepareProposedRates checks proposed rates against the rate sheet rules and
// the active sheets. It changes nothing. Every save needs a reason and the
// sheet to start from; submission checks the reason again, for a draft saved
// before the reason was required. The name is an optional description, as for
// any sheet; the draft has its code from its first save.
func prepareProposedRates(app core.App, r *core.Record, own *core.Record, input proposedRatesWrite) (*proposedRatesPlan, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.Reason = strings.TrimSpace(input.Reason)
	input.BasedOn = strings.TrimSpace(input.BasedOn)

	errs := validation.Errors{}
	prefix := hooks.ProposedRatesField + "."
	if input.Reason == "" {
		setValidationError(errs, prefix+"reason", validation.NewError("value_required", "reason is required"))
	} else if utf8.RuneCountInString(input.Reason) > proposedRatesReasonMax {
		setValidationError(errs, prefix+"reason", validation.NewError("validation_length_too_long", "the reason can have at most 500 characters"))
	}
	if input.BasedOn == "" {
		setValidationError(errs, prefix+"based_on", validation.NewError("value_required", "choose the rate sheet to start from"))
	}
	validateRateSheetEntryInputs(input.Entries, prefix, errs)
	if len(errs) > 0 {
		return nil, errs
	}

	// The proposal starts from an active sheet. A draft keeps its source
	// after that sheet is deactivated, for the approver's comparison.
	source, err := app.FindRecordById("rate_sheets", input.BasedOn)
	if err != nil || hooks.IsRateSheetDraft(source) {
		return nil, workflowError("Start the proposed rates from an active rate sheet.")
	}
	if !source.GetBool("active") && (own == nil || own.GetString("based_on") != source.Id) {
		return nil, workflowError("Start the proposed rates from an active rate sheet.")
	}
	for _, entry := range input.Entries {
		if _, err := app.FindRecordById("rate_roles", entry.Role); err != nil {
			return nil, workflowError("A proposed rate has an unknown role.")
		}
	}

	// Any other sheet holds its name, in any letter case, as for a rename.
	ownCode := ""
	if own != nil {
		ownCode = own.GetString("code")
	}
	holder, err := hooks.RateSheetNameHolder(app, input.Name, ownCode)
	if err != nil {
		return nil, err
	}
	if holder != nil {
		code, message := hooks.RateSheetNameHeld(holder)
		return nil, validation.Errors{prefix + "name": validation.NewError(code, message)}
	}

	// Turbo refuses rates that an active sheet already has, and names it, so
	// the requester can choose that sheet instead. A draft is inactive and
	// has its own code, so no active sheet is one of its revisions.
	values := rateSheetEntryValues(input.Entries)
	comparison, err := hooks.CompareWithActiveRateSheets(app, values, ownCode)
	if err != nil {
		return nil, err
	}
	if len(comparison.Exact) > 0 {
		return nil, hooks.DuplicateRateSheetError(hooks.ProposedRatesField, comparison.Exact)
	}

	plan := &proposedRatesPlan{draft: own, input: input, changed: own == nil}
	if own != nil {
		stored, err := hooks.LoadRateSheetEntryValues(app, own.Id)
		if err != nil {
			return nil, err
		}
		plan.changed = own.GetString("name") != input.Name ||
			own.GetString("reason") != input.Reason ||
			own.GetString("based_on") != input.BasedOn ||
			hooks.RateSheetRatesKey(stored) != hooks.RateSheetRatesKey(values)
	}
	return plan, nil
}

// saveProposedRates writes the draft and its entries and points the request at
// the draft. The request must already be saved, since the draft refers to it.
// The caller saves the request afterwards.
func saveProposedRates(app core.App, r *core.Record, plan *proposedRatesPlan) error {
	draft := plan.draft
	if draft == nil {
		collection, err := app.FindCollectionByNameOrId("rate_sheets")
		if err != nil {
			return err
		}
		draft = core.NewRecord(collection)
		draft.Set("revision", 0)
		draft.Set("active", false)
		draft.Set("request", r.Id)
	}
	if !plan.changed {
		r.Set("rate_sheet", draft.Id)
		return nil
	}
	draft.Set("name", plan.input.Name)
	draft.Set("reason", plan.input.Reason)
	draft.Set("based_on", plan.input.BasedOn)
	// Changed rates need a new approval.
	clearRatesApproval(draft)
	// Job creation sets the date the sheet becomes active.
	draft.Set("effective_date", hooks.RateSheetToday())
	if err := app.Save(draft); err != nil {
		return err
	}

	stored, err := app.FindRecordsByFilter("rate_sheet_entries", "rate_sheet = {:sheet}", "", 0, 0, dbx.Params{"sheet": draft.Id})
	if err != nil {
		return err
	}
	byRole := map[string]*core.Record{}
	for _, entry := range stored {
		byRole[entry.GetString("role")] = entry
	}
	entries, err := app.FindCollectionByNameOrId("rate_sheet_entries")
	if err != nil {
		return err
	}
	for _, input := range plan.input.Entries {
		entry, ok := byRole[input.Role]
		delete(byRole, input.Role)
		if !ok {
			entry = core.NewRecord(entries)
			entry.Set("rate_sheet", draft.Id)
			entry.Set("role", input.Role)
		} else if entry.GetFloat("rate") == input.Rate && entry.GetFloat("overtime_rate") == input.OvertimeRate {
			continue
		}
		entry.Set("rate", int(input.Rate))
		entry.Set("overtime_rate", input.OvertimeRate)
		if err := app.Save(entry); err != nil {
			return err
		}
	}
	for _, entry := range byRole {
		if err := app.Delete(entry); err != nil {
			return err
		}
	}
	r.Set("rate_sheet", draft.Id)
	return nil
}

// deleteRequestDraft deletes a draft that the request no longer uses, after the
// request was saved with its new rate sheet.
func deleteRequestDraft(app core.App, draft *core.Record) error {
	if draft == nil {
		return nil
	}
	return app.Delete(draft)
}

// renameProposedRates gives the request's proposed rates a new name as the job
// is created. The job holder who creates the job keeps the rate sheet list, so
// they can correct the requester's name before the sheet becomes active. The
// approver approved the rates, not the name, so the approval stands. The old
// name is kept in rate_sheet_name_changes, as for an admin's rename. An empty
// or unchanged name does nothing. The caller checks the job claim.
func renameProposedRates(app core.App, r *core.Record, name, renamedBy string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil
	}
	draft, err := hooks.RequestRateSheetDraft(app, r)
	if err != nil {
		return err
	}
	if draft == nil {
		return workflowError("This request has no proposed rates to rename.")
	}
	oldName := draft.GetString("name")
	if name == oldName {
		return nil
	}
	holder, err := hooks.RateSheetNameHolder(app, name, draft.GetString("code"))
	if err != nil {
		return err
	}
	if holder != nil {
		_, message := hooks.RateSheetNameHeld(holder)
		return workflowError(message)
	}
	draft.Set("name", name)
	if err := app.Save(draft); err != nil {
		return err
	}
	changes, err := app.FindCollectionByNameOrId("rate_sheet_name_changes")
	if err != nil {
		return err
	}
	change := core.NewRecord(changes)
	change.Set("rate_sheet", draft.Id)
	change.Set("old_name", oldName)
	change.Set("new_name", name)
	change.Set("renamed_by", renamedBy)
	return app.Save(change)
}

// requestProposedRates adds the proposal to a request's related records, so
// the editor and the approver see it: the draft's reason, the sheet it
// started from, and each rate beside the source's, with changes marked.
// rate_sheet_missing_roles names the roles without a rate, by the check that
// submission makes, so the printed authorization is signable only when it is
// empty. rate_sheet_missing_rates has a row for each of them, with the role's
// id and the source sheet's rates when it has them, so the review can list
// them. Both are for display and stay out of the review token. It does
// nothing when the request uses an active sheet.
func requestProposedRates(app core.App, r *core.Record, result map[string]any, rates []map[string]any) error {
	draft, err := hooks.RequestRateSheetDraft(app, r)
	if err != nil || draft == nil {
		return err
	}
	result["rate_sheet_proposed"] = true
	missing, err := hooks.ValidateRateSheetComplete(app, draft.Id)
	if err != nil {
		return err
	}
	if missing == nil {
		missing = []string{}
	}
	result["rate_sheet_missing_roles"] = missing
	if id := draft.GetString("rates_approver"); id != "" {
		if approver, err := app.FindRecordById("users", id); err == nil {
			result["rate_sheet_rates_approver"] = workflowPerson(app, approver)
		}
	}
	sourceRates := map[string]hooks.RateSheetEntryValues{}
	if id := draft.GetString("based_on"); id != "" {
		source, err := app.FindRecordById("rate_sheets", id)
		if err == nil {
			result["rate_sheet_source"] = recordPublic(source)
			values, err := hooks.LoadRateSheetEntryValues(app, source.Id)
			if err != nil {
				return err
			}
			for _, value := range values {
				sourceRates[value.Role] = value
			}
		}
	}
	missingRates, err := requestMissingRates(app, draft.Id, sourceRates)
	if err != nil {
		return err
	}
	result["rate_sheet_missing_rates"] = missingRates
	for _, rate := range rates {
		role, _ := rate["role_id"].(string)
		source, ok := sourceRates[role]
		if !ok {
			rate["changed"] = true
			continue
		}
		rate["source_rate"] = source.Rate
		rate["source_overtime_rate"] = source.OvertimeRate
		rate["changed"] = source.Rate != rate["rate"] || source.OvertimeRate != rate["overtime_rate"]
	}
	return nil
}

// requestMissingRates returns a row for each rate role that the draft has no
// rate for, in the order of ValidateRateSheetComplete: role_id, role_name and,
// when the source sheet has a rate for the role, source_rate and
// source_overtime_rate.
func requestMissingRates(app core.App, draftID string, sourceRates map[string]hooks.RateSheetEntryValues) ([]map[string]any, error) {
	roles := []struct {
		ID   string `db:"id"`
		Name string `db:"name"`
	}{}
	if err := app.DB().NewQuery(`
		SELECT r.id, r.name
		FROM rate_roles r
		WHERE r.id NOT IN (
			SELECT rse.role FROM rate_sheet_entries rse WHERE rse.rate_sheet = {:sheet}
		)
		ORDER BY r.name
	`).Bind(dbx.Params{"sheet": draftID}).All(&roles); err != nil {
		return nil, err
	}
	rows := []map[string]any{}
	for _, role := range roles {
		row := map[string]any{"role_id": role.ID, "role_name": role.Name}
		if source, ok := sourceRates[role.ID]; ok {
			row["source_rate"] = source.Rate
			row["source_overtime_rate"] = source.OvertimeRate
		}
		rows = append(rows, row)
	}
	return rows, nil
}

// The approver of a request with proposed rates approves the creation of the
// new rate sheet and its rates. This is a statement of its own, apart from the
// checklist, and the draft records who made it and when. The record stays on
// the sheet after job creation, when the request is gone.

func clearRatesApproval(draft *core.Record) {
	draft.Set("rates_approver", "")
	draft.Set("rates_approved", "")
}

// approveProposedRates records the approval on the request's draft. It refuses
// when the approver did not make the statement, and does nothing when the
// request uses an active sheet. The caller checks that approver is the
// request's approver.
func approveProposedRates(app core.App, r *core.Record, approved bool, approver string) error {
	draft, err := hooks.RequestRateSheetDraft(app, r)
	if err != nil || draft == nil {
		return err
	}
	if !approved {
		return workflowError("Approve the creation of the new rate sheet and its rates.")
	}
	draft.Set("rates_approver", approver)
	draft.Set("rates_approved", time.Now().UTC())
	return app.Save(draft)
}

// clearProposedRatesApproval ends the approval of the request's draft, so a
// new submission needs a new one.
func clearProposedRatesApproval(app core.App, r *core.Record) error {
	draft, err := hooks.RequestRateSheetDraft(app, r)
	if err != nil || draft == nil {
		return err
	}
	if draft.GetString("rates_approver") == "" && draft.GetString("rates_approved") == "" {
		return nil
	}
	clearRatesApproval(draft)
	return app.Save(draft)
}

// requireProposedRatesApproval refuses job creation when the request's draft
// has no approval by the request's approver.
func requireProposedRatesApproval(app core.App, r *core.Record) error {
	draft, err := hooks.RequestRateSheetDraft(app, r)
	if err != nil || draft == nil {
		return err
	}
	if draft.GetString("rates_approver") == "" || draft.GetString("rates_approver") != r.GetString("approver") || draft.GetString("rates_approved") == "" {
		return workflowError("The proposed rates are not approved. Recall the request for a new approval.")
	}
	return nil
}

// requestCompareRatesHandler compares unsaved proposed rates with the active
// rate sheets, for POST /api/job_requests/compare_rates. It returns the exact
// matches, which a save would refuse, and a few close sheets, so the requester
// can choose one of them instead of proposing a new sheet. The sheet in
// based_on is left out of the close sheets: the requester started from it and
// already sees how the rates differ from it. A role with no rate is left out
// of the entries, and the close sheets do not count it as a difference
// (IgnoringRolesLeftOut); the exact matches are those a save refuses. It
// changes nothing.
// Anyone who can start a request can call it; the rate sheets compare route
// needs a claim that most requesters lack.
func requestCompareRatesHandler(app core.App) func(*core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		if err := requestEnabled(app); err != nil {
			return err
		}
		if err := workflowClaim(app, e.Auth, "time"); err != nil {
			return err
		}
		var input proposedRatesWrite
		if err := e.BindBody(&input); err != nil {
			return workflowError("Invalid request data.")
		}
		// The proposal is a new sheet, so every active sheet is compared.
		comparison, err := hooks.CompareWithActiveRateSheets(app, rateSheetEntryValues(input.Entries), "")
		if err != nil {
			return err
		}
		comparison = comparison.IgnoringRolesLeftOut()
		near := []hooks.RateSheetMatch{}
		for _, match := range comparison.Near {
			if match.ID != strings.TrimSpace(input.BasedOn) {
				near = append(near, match)
			}
		}
		comparison.Near = near
		return e.JSON(http.StatusOK, map[string]any{
			"exact": comparison.Exact,
			"near":  comparison.Closest(hooks.RateSheetNearMatchMaxRoles, hooks.RateSheetNearMatchLimit),
		})
	}
}

// findRatesEntry is one role's rate that the requester knows, in the body of
// POST /api/job_requests/find_rates. OvertimeRate is optional: without it, a
// sheet matches on the rate alone.
type findRatesEntry struct {
	Role         string   `json:"role"`
	Rate         *float64 `json:"rate"`
	OvertimeRate *float64 `json:"overtime_rate"`
}

// foundRateSheet is an active sheet that has every rate searched for.
// JobCount counts the jobs that use the revision now or in their rate sheet
// history, as job_count in rate_sheets_augmented does.
type foundRateSheet struct {
	ID            string `json:"id"`
	Code          string `json:"code"`
	Name          string `json:"name"`
	Revision      int    `json:"revision"`
	EffectiveDate string `json:"effective_date"`
	JobCount      int    `json:"job_count"`
}

// validateFindRatesEntries checks the rates of a find: at least one, each
// role once and known, a whole rate of at least 1, and an overtime rate, when
// given, of at least 1 with at most two decimal places, as for a new sheet. An
// overtime rate needs the role's rate. Keys are entries.N.field.
func validateFindRatesEntries(app core.App, entries []findRatesEntry) (validation.Errors, error) {
	errs := validation.Errors{}
	if len(entries) == 0 {
		setValidationError(errs, "entries", validation.NewError("value_required", "enter at least one rate"))
		return errs, nil
	}
	roles := []struct {
		ID string `db:"id"`
	}{}
	if err := app.DB().NewQuery("SELECT id FROM rate_roles").All(&roles); err != nil {
		return nil, err
	}
	known := map[string]bool{}
	for _, role := range roles {
		known[role.ID] = true
	}
	seen := map[string]bool{}
	for i, entry := range entries {
		prefix := fmt.Sprintf("entries.%d.", i)
		switch {
		case entry.Role == "":
			setValidationError(errs, prefix+"role", validation.NewError("value_required", "entry role is required"))
		case !known[entry.Role]:
			setValidationError(errs, prefix+"role", validation.NewError("validation_invalid_role", "entry role is unknown"))
		case seen[entry.Role]:
			setValidationError(errs, prefix+"role", validation.NewError("validation_not_unique", "entry role must be unique"))
		default:
			seen[entry.Role] = true
		}
		switch {
		case entry.Rate == nil && entry.OvertimeRate != nil:
			setValidationError(errs, prefix+"rate", validation.NewError("value_required", "enter the rate of a role that has an overtime rate"))
		case entry.Rate == nil:
			setValidationError(errs, prefix+"rate", validation.NewError("value_required", "entry rate is required"))
		case *entry.Rate < 1:
			setValidationError(errs, prefix+"rate", validation.NewError("validation_negative_number", "entry rate must be at least 1"))
		case *entry.Rate != math.Floor(*entry.Rate):
			setValidationError(errs, prefix+"rate", validation.NewError("validation_not_whole_number", "entry rate must be a whole number"))
		}
		if entry.OvertimeRate != nil {
			if *entry.OvertimeRate < 1 {
				setValidationError(errs, prefix+"overtime_rate", validation.NewError("validation_negative_number", "entry overtime_rate must be at least 1"))
			} else if !hooks.HasAtMostTwoDecimals(*entry.OvertimeRate) {
				setValidationError(errs, prefix+"overtime_rate", validation.NewError(hooks.OvertimeRateDecimalsCode, "entry "+hooks.OvertimeRateDecimalsMessage))
			}
		}
	}
	return errs, nil
}

// rateSheetJobCounts returns the number of jobs of each sheet in ids that a
// job uses now or in its rate sheet history, in one query. A sheet with no
// job is left out.
func rateSheetJobCounts(app core.App, ids []string) (map[string]int, error) {
	counts := map[string]int{}
	if len(ids) == 0 {
		return counts, nil
	}
	params := dbx.Params{}
	placeholders := make([]string, 0, len(ids))
	for i, id := range ids {
		name := fmt.Sprintf("sheet%d", i)
		params[name] = id
		placeholders = append(placeholders, "{:"+name+"}")
	}
	rows := []struct {
		ID       string `db:"id"`
		JobCount int    `db:"job_count"`
	}{}
	err := app.DB().NewQuery(`
		SELECT u.rate_sheet AS id, COUNT(u.job) AS job_count
		FROM (
			SELECT j.rate_sheet AS rate_sheet, j.id AS job FROM jobs j WHERE j.rate_sheet != ''
			UNION
			SELECT h.rate_sheet AS rate_sheet, h.job AS job FROM job_rate_sheets h
		) u
		WHERE u.rate_sheet IN (` + strings.Join(placeholders, ", ") + `)
		GROUP BY u.rate_sheet
	`).Bind(params).All(&rows)
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		counts[row.ID] = row.JobCount
	}
	return counts, nil
}

// requestFindRatesHandler finds the active rate sheets that have the rates a
// requester knows, for POST /api/job_requests/find_rates. The body is
// {entries: [{role, rate, overtime_rate}]}, with overtime_rate optional. A
// sheet matches when, for every entry, it has the role with the same rate,
// and the same overtime rate when one is given; roles the body leaves out are
// ignored (hooks.FindActiveRateSheetsByRates). It returns {matches: [{id,
// code, name, revision, effective_date, job_count}]}, newest code first. A
// request's draft is never returned. It changes nothing. It has the
// permission checks of the compare route: project requests on, and the time
// claim that starting a request needs.
func requestFindRatesHandler(app core.App) func(*core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		if err := requestEnabled(app); err != nil {
			return err
		}
		if err := workflowClaim(app, e.Auth, "time"); err != nil {
			return err
		}
		var input struct {
			Entries []findRatesEntry `json:"entries"`
		}
		if err := e.BindBody(&input); err != nil {
			return workflowError("Invalid request data.")
		}
		for i := range input.Entries {
			input.Entries[i].Role = strings.TrimSpace(input.Entries[i].Role)
		}
		errs, err := validateFindRatesEntries(app, input.Entries)
		if err != nil {
			return err
		}
		if len(errs) > 0 {
			return errs
		}
		wanted := make([]hooks.RateSheetRateQuery, 0, len(input.Entries))
		for _, entry := range input.Entries {
			wanted = append(wanted, hooks.RateSheetRateQuery{Role: entry.Role, Rate: *entry.Rate, OvertimeRate: entry.OvertimeRate})
		}
		matches, err := hooks.FindActiveRateSheetsByRates(app, wanted)
		if err != nil {
			return err
		}
		ids := make([]string, 0, len(matches))
		for _, match := range matches {
			ids = append(ids, match.ID)
		}
		counts, err := rateSheetJobCounts(app, ids)
		if err != nil {
			return err
		}
		found := make([]foundRateSheet, 0, len(matches))
		for _, match := range matches {
			found = append(found, foundRateSheet{
				ID:            match.ID,
				Code:          match.Code,
				Name:          match.Name,
				Revision:      match.Revision,
				EffectiveDate: match.EffectiveDate,
				JobCount:      counts[match.ID],
			})
		}
		return e.JSON(http.StatusOK, map[string]any{"matches": found})
	}
}
