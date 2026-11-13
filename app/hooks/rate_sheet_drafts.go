package hooks

import (
	"fmt"
	"net/http"
	"strings"

	"tybalt/errs"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

// A project request can propose new rates instead of choosing an active rate
// sheet. The proposal is a draft: a rate sheet whose request field names the
// request that owns it. A draft is inactive, the API lists it nowhere, and no
// job or rate sheet route accepts it. When a job holder creates the job from
// the approved request, ActivateRequestRateSheet activates the draft in the
// same transaction and clears request, so the sheet becomes an ordinary active
// sheet that keeps its reason. Deleting the request deletes its draft. See the
// Proposed rates section of descriptions/job_rate_sheet_history.md.

// ProposedRatesField is the key under which the request routes and their
// errors carry the proposed rates.
const ProposedRatesField = "proposed_rates"

// IsRateSheetDraft reports whether a rate sheet is a draft that a project
// request owns.
func IsRateSheetDraft(sheet *core.Record) bool {
	return sheet != nil && sheet.GetString("request") != ""
}

// RateSheetDraftError refuses a direct use of a draft. field is the key the UI
// shows the error under.
func RateSheetDraftError(field string) *errs.HookError {
	message := "this rate sheet is proposed in a project request; it is activated when the job is created"
	return &errs.HookError{
		Status:  http.StatusBadRequest,
		Message: message,
		Data: map[string]errs.CodeError{
			field: {Code: "rate_sheet_draft", Message: message},
		},
	}
}

// RequestRateSheetDrafts returns the drafts that a request owns. The request
// routes keep at most one.
func RequestRateSheetDrafts(app core.App, requestID string) ([]*core.Record, error) {
	if requestID == "" {
		return nil, nil
	}
	return app.FindRecordsByFilter("rate_sheets", "request = {:request}", "", 0, 0, dbx.Params{"request": requestID})
}

// RequestRateSheetDraft returns the request's rate sheet when it is a draft
// that this request owns, or nil when the request uses any other sheet or none.
func RequestRateSheetDraft(app core.App, request *core.Record) (*core.Record, error) {
	id := request.GetString("rate_sheet")
	if id == "" || request.Id == "" {
		return nil, nil
	}
	sheet, err := app.FindRecordById("rate_sheets", id)
	if err != nil {
		// A missing sheet is reported by the normal rate sheet checks.
		return nil, nil
	}
	if sheet.GetString("request") != request.Id {
		return nil, nil
	}
	return sheet, nil
}

func proposedRatesError(message string) *errs.HookError {
	return &errs.HookError{
		Status:  http.StatusBadRequest,
		Message: message,
		Data: map[string]errs.CodeError{
			ProposedRatesField: {Code: "invalid_proposed_rates", Message: message},
		},
	}
}

// ValidateRequestRateSheetDraft checks a draft for submission, approval and job
// creation: it has a reason, an entry for every role, and rates that differ
// from every active sheet. Rates are compared again at each step because an
// active sheet can change after the draft is saved.
func ValidateRequestRateSheetDraft(app core.App, draft *core.Record) error {
	if strings.TrimSpace(draft.GetString("reason")) == "" {
		return proposedRatesError("enter the reason why this project needs new rates")
	}
	missing, err := ValidateRateSheetComplete(app, draft.Id)
	if err != nil {
		return err
	}
	if len(missing) > 0 {
		return proposedRatesError("the proposed rates have no rate for " + strings.Join(missing, ", "))
	}
	entries, err := LoadRateSheetEntryValues(app, draft.Id)
	if err != nil {
		return err
	}
	comparison, err := CompareWithActiveRateSheets(app, entries, draft.GetString("name"))
	if err != nil {
		return err
	}
	if len(comparison.Exact) > 0 {
		return DuplicateRateSheetError(ProposedRatesField, comparison.Exact)
	}
	return nil
}

// ActivateRateSheet activates sheet after checking that no active sheet of
// another name has the same rates, and deactivates the other revisions of its
// name first, so a name never has two active revisions. Call it in a
// transaction. The caller checks permissions, completeness and newer
// revisions.
func ActivateRateSheet(app core.App, sheet *core.Record) error {
	// Entries can be added between create and activate, so compare the stored
	// entries with the other active sheets. Revisions of this sheet are
	// excluded by name.
	entries, err := LoadRateSheetEntryValues(app, sheet.Id)
	if err != nil {
		return err
	}
	comparison, err := CompareWithActiveRateSheets(app, entries, sheet.GetString("name"))
	if err != nil {
		return err
	}
	if len(comparison.Exact) > 0 {
		return DuplicateRateSheetError("active", comparison.Exact)
	}
	if err := DeactivateOtherRevisions(app, sheet.GetString("name"), sheet.Id); err != nil {
		return err
	}
	sheet.Set("active", true)
	return app.Save(sheet)
}

// ActivateRequestRateSheet activates the request's draft when the job is
// created from it. It does nothing when the request uses an active sheet. The
// draft's effective date becomes the activation date and request is cleared,
// so deleting the request afterwards keeps the sheet.
func ActivateRequestRateSheet(app core.App, request *core.Record) error {
	draft, err := RequestRateSheetDraft(app, request)
	if err != nil || draft == nil {
		return err
	}
	// This also makes the duplicate check that ActivateRateSheet repeats, so
	// a duplicate is reported under the proposed rates.
	if err := ValidateRequestRateSheetDraft(app, draft); err != nil {
		return err
	}
	draft.Set("request", "")
	draft.Set("effective_date", RateSheetToday())
	return ActivateRateSheet(app, draft)
}

// addRateSheetDraftHooks deletes a request's draft with the request, on every
// path: the creator's delete, an admin's delete of an abandoned request, or
// any other delete. Job creation clears the draft's request first, so the
// activated sheet stays.
func addRateSheetDraftHooks(app core.App) {
	app.OnRecordDelete("job_requests").BindFunc(func(e *core.RecordEvent) error {
		drafts, err := RequestRateSheetDrafts(e.App, e.Record.Id)
		if err != nil {
			return err
		}
		// The request goes first: the draft is its rate sheet, and the
		// request may only be saved through the project workflow.
		if err := e.Next(); err != nil {
			return err
		}
		for _, draft := range drafts {
			current, err := e.App.FindRecordById("rate_sheets", draft.Id)
			if err != nil {
				continue // already gone
			}
			if err := e.App.Delete(current); err != nil {
				return fmt.Errorf("deleting the proposed rate sheet of request %s: %w", e.Record.Id, err)
			}
		}
		return nil
	})
}
