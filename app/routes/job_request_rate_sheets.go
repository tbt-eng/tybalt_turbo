package routes

import (
	"strings"
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
// the active sheets. It changes nothing. A draft can be saved without a reason,
// like any incomplete request; submission requires one.
func prepareProposedRates(app core.App, r *core.Record, own *core.Record, input proposedRatesWrite) (*proposedRatesPlan, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.Reason = strings.TrimSpace(input.Reason)
	input.BasedOn = strings.TrimSpace(input.BasedOn)

	errs := validation.Errors{}
	prefix := hooks.ProposedRatesField + "."
	if input.Name == "" {
		setValidationError(errs, prefix+"name", validation.NewError("value_required", "name is required"))
	}
	if utf8.RuneCountInString(input.Reason) > proposedRatesReasonMax {
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

	// Any revision of any sheet holds its name, in any letter case, as a
	// rename does.
	ownID := ""
	if own != nil {
		ownID = own.Id
	}
	var named int
	if err := app.DB().NewQuery(`
		SELECT COUNT(*) FROM rate_sheets
		WHERE LOWER(name) = LOWER({:name}) AND id != {:own}
	`).Bind(dbx.Params{"name": input.Name, "own": ownID}).Row(&named); err != nil {
		return nil, err
	}
	if named > 0 {
		return nil, validation.Errors{prefix + "name": validation.NewError("validation_not_unique", "another rate sheet already has this name")}
	}

	// Turbo refuses rates that an active sheet already has, and names it, so
	// the requester can choose that sheet instead.
	values := rateSheetEntryValues(input.Entries)
	comparison, err := hooks.CompareWithActiveRateSheets(app, values, input.Name)
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

// requestProposedRates adds the proposal to a request's related records, so
// the editor and the approver see it: the draft's reason, the sheet it
// started from, and each rate beside the source's, with changes marked.
// It does nothing when the request uses an active sheet.
func requestProposedRates(app core.App, r *core.Record, result map[string]any, rates []map[string]any) error {
	draft, err := hooks.RequestRateSheetDraft(app, r)
	if err != nil || draft == nil {
		return err
	}
	result["rate_sheet_proposed"] = true
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

// proposedRatesChecklistText replaces the rate sheet statement of the approval
// checklist for proposed rates. The approver approves the rates themselves, not
// the choice of an existing sheet. The key and the stored version stay the
// same, so approvals need no new checklist version.
const proposedRatesChecklistText = "The proposed rates are approved for this client."

// requestRateSheetChecklistKey is the approval statement about the rate sheet.
const requestRateSheetChecklistKey = "rate_sheet_correct"
