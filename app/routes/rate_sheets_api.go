package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strings"

	"github.com/pocketbase/dbx"
	validation "github.com/pocketbase/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"

	"tybalt/errs"
	"tybalt/hooks"
	"tybalt/utilities"
)

// createActivateRateSheetHandler creates a handler that activates a rate sheet.
// Only users with the 'job' claim can activate rate sheets.
// The rate sheet must have entries for all roles before it can be activated.
func createActivateRateSheetHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		id := e.Request.PathValue("id")

		// Check for job claim
		hasJobClaim, err := utilities.HasClaim(app, e.Auth, "job")
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to check claims", err)
		}
		if !hasJobClaim {
			return e.Error(http.StatusForbidden, "you do not have permission to activate rate sheets", nil)
		}

		// Get the rate sheet
		record, err := app.FindRecordById("rate_sheets", id)
		if err != nil {
			return e.Error(http.StatusNotFound, "rate sheet not found", err)
		}

		// A request's proposed rates are activated when the job is created.
		if hooks.IsRateSheetDraft(record) {
			return e.JSON(http.StatusBadRequest, hooks.RateSheetDraftError("active"))
		}

		// Check if already active
		if record.GetBool("active") {
			return e.JSON(http.StatusOK, map[string]any{
				"message": "rate sheet is already active",
				"id":      id,
				"active":  true,
			})
		}

		// Check if a newer revision exists (cannot activate older revisions)
		code := record.GetString("code")
		revision := record.GetInt("revision")
		newerExists, err := hooks.CheckNewerRevisionExists(app, code, revision)
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to check for newer revisions", err)
		}
		if newerExists {
			return e.Error(http.StatusBadRequest, "cannot activate - a newer revision exists", nil)
		}

		// Validate completeness using the existing hook logic
		missingRoles, err := hooks.ValidateRateSheetComplete(app, id)
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to validate rate sheet", err)
		}
		if len(missingRoles) > 0 {
			return e.Error(http.StatusBadRequest, "rate sheet is incomplete - missing entries for some roles", map[string]any{
				"missing_roles": missingRoles,
			})
		}

		// Check for a duplicate, deactivate the other revisions and activate this
		// one in a single transaction, so that two identical sheets can't both
		// pass the check and be activated at the same time.
		err = app.RunInTransaction(func(txApp core.App) error {
			return hooks.ActivateRateSheet(txApp, record)
		})
		if err != nil {
			var hookErr *errs.HookError
			if errors.As(err, &hookErr) {
				return e.JSON(hookErr.Status, hookErr)
			}
			return e.Error(http.StatusInternalServerError, "failed to activate rate sheet", err)
		}

		return e.JSON(http.StatusOK, map[string]any{
			"message": "rate sheet activated",
			"id":      id,
			"active":  true,
		})
	}
}

// createDeactivateRateSheetHandler creates a handler that deactivates a rate sheet.
// Only users with the 'job' claim can deactivate rate sheets.
func createDeactivateRateSheetHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		id := e.Request.PathValue("id")

		// Check for job claim
		hasJobClaim, err := utilities.HasClaim(app, e.Auth, "job")
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to check claims", err)
		}
		if !hasJobClaim {
			return e.Error(http.StatusForbidden, "you do not have permission to deactivate rate sheets", nil)
		}

		// Get the rate sheet
		record, err := app.FindRecordById("rate_sheets", id)
		if err != nil {
			return e.Error(http.StatusNotFound, "rate sheet not found", err)
		}

		// Check if already inactive
		if !record.GetBool("active") {
			return e.JSON(http.StatusOK, map[string]any{
				"message": "rate sheet is already inactive",
				"id":      id,
				"active":  false,
			})
		}

		// Update the record
		record.Set("active", false)
		if err := app.Save(record); err != nil {
			return e.Error(http.StatusInternalServerError, "failed to deactivate rate sheet", err)
		}

		return e.JSON(http.StatusOK, map[string]any{
			"message": "rate sheet deactivated",
			"id":      id,
			"active":  false,
		})
	}
}

// UpdateRateSheetEntryRequest defines the allowed fields for updating a rate sheet entry
type UpdateRateSheetEntryRequest struct {
	Rate         *float64 `json:"rate"`
	OvertimeRate *float64 `json:"overtime_rate"`
}

// createUpdateRateSheetEntryHandler creates a handler that updates a rate sheet entry.
// Only users with the 'admin' claim can update entries.
// Only rate and overtime_rate fields can be updated.
func createUpdateRateSheetEntryHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		id := e.Request.PathValue("id")

		// Check for admin claim
		hasAdminClaim, err := utilities.HasClaim(app, e.Auth, "admin")
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to check claims", err)
		}
		if !hasAdminClaim {
			return e.Error(http.StatusForbidden, "you do not have permission to update rate sheet entries", nil)
		}

		// Parse request body
		var req UpdateRateSheetEntryRequest
		if err := json.NewDecoder(e.Request.Body).Decode(&req); err != nil {
			return e.Error(http.StatusBadRequest, "invalid request body", err)
		}

		// Validate at least one field is being updated
		if req.Rate == nil && req.OvertimeRate == nil {
			return e.Error(http.StatusBadRequest, "at least one of rate or overtime_rate must be provided", nil)
		}

		// Get the entry
		record, err := app.FindRecordById("rate_sheet_entries", id)
		if err != nil {
			return e.Error(http.StatusNotFound, "rate sheet entry not found", err)
		}

		// Update only the provided fields
		if req.Rate != nil {
			if *req.Rate < 1 {
				return e.Error(http.StatusBadRequest, "rate must be at least 1", nil)
			}
			// Validate rate is a whole number
			if *req.Rate != math.Floor(*req.Rate) {
				return e.Error(http.StatusBadRequest, "rate must be a whole number", nil)
			}
			record.Set("rate", int(*req.Rate))
		}
		if req.OvertimeRate != nil {
			if *req.OvertimeRate < 1 {
				return e.Error(http.StatusBadRequest, "overtime_rate must be at least 1", nil)
			}
			if !hooks.HasAtMostTwoDecimals(*req.OvertimeRate) {
				return e.Error(http.StatusBadRequest, hooks.OvertimeRateDecimalsMessage, nil)
			}
			record.Set("overtime_rate", *req.OvertimeRate)
		}

		// Check for a duplicate and save in a single transaction. A change may
		// not make an active sheet's rates equal to another active sheet's. An
		// inactive sheet is not checked here, because activation checks it.
		err = app.RunInTransaction(func(txApp core.App) error {
			sheet, err := txApp.FindRecordById("rate_sheets", record.GetString("rate_sheet"))
			if err != nil {
				return err
			}
			// The requester changes proposed rates through the request.
			if hooks.IsRateSheetDraft(sheet) {
				return hooks.RateSheetDraftError("entry")
			}
			if sheet.GetBool("active") {
				entries, err := hooks.LoadRateSheetEntryValues(txApp, sheet.Id)
				if err != nil {
					return err
				}
				// Apply the change. Each role is on a sheet at most once.
				role := record.GetString("role")
				for i := range entries {
					if entries[i].Role == role {
						entries[i].Rate = record.GetFloat("rate")
						entries[i].OvertimeRate = record.GetFloat("overtime_rate")
					}
				}
				// Revisions of this sheet are excluded by code.
				comparison, err := hooks.CompareWithActiveRateSheets(txApp, entries, sheet.GetString("code"))
				if err != nil {
					return err
				}
				if len(comparison.Exact) > 0 {
					return hooks.DuplicateRateSheetError("entry", comparison.Exact)
				}
			}
			return txApp.Save(record)
		})
		if err != nil {
			var hookErr *errs.HookError
			if errors.As(err, &hookErr) {
				return e.JSON(hookErr.Status, hookErr)
			}
			return e.Error(http.StatusInternalServerError, "failed to update rate sheet entry", err)
		}

		return e.JSON(http.StatusOK, map[string]any{
			"message":       "rate sheet entry updated",
			"id":            id,
			"rate":          record.GetInt("rate"),
			"overtime_rate": record.GetFloat("overtime_rate"),
		})
	}
}

// CreateRateSheetEntry defines a single entry in the rate sheet creation request
type CreateRateSheetEntry struct {
	Role         string  `json:"role"`
	Rate         float64 `json:"rate"`
	OvertimeRate float64 `json:"overtime_rate"`
}

// CreateRateSheetRequest defines the request body for creating a rate sheet with entries.
// Code names the sheet that a revision (revision > 0) revises; the new revision takes that
// code. A new sheet (revision 0) gets the next code when it is saved, so it may not send one.
// A revision without a code is matched to its sheet by name, as before codes existed.
type CreateRateSheetRequest struct {
	Code          string                 `json:"code"`
	Name          string                 `json:"name"`
	EffectiveDate string                 `json:"effective_date"`
	Revision      int                    `json:"revision"`
	Entries       []CreateRateSheetEntry `json:"entries"`
}

type entryValidationError struct {
	index  int
	errors validation.Errors
}

func (e *entryValidationError) Error() string {
	return "entry validation error"
}

type codeError interface {
	Code() string
}

func setValidationError(errs validation.Errors, key string, err error) {
	if err == nil {
		return
	}
	if _, exists := errs[key]; !exists {
		errs[key] = err
	}
}

func rateSheetUniqueValidationErrors() validation.Errors {
	return validation.Errors{
		"revision": validation.NewError("validation_not_unique", "code and revision must be unique"),
	}
}

func hasUniqueValidationError(errs validation.Errors) bool {
	for _, fieldErr := range errs {
		if fieldErr == nil {
			continue
		}
		if ce, ok := fieldErr.(codeError); ok && ce.Code() == "validation_not_unique" {
			return true
		}
		if strings.Contains(strings.ToLower(fieldErr.Error()), "unique") {
			return true
		}
	}
	return false
}

func asValidationErrors(err error) (validation.Errors, bool) {
	if err == nil {
		return nil, false
	}
	if ve, ok := err.(validation.Errors); ok {
		return ve, true
	}
	if ve, ok := err.(*validation.Errors); ok && ve != nil {
		return *ve, true
	}
	return nil, false
}

func prefixEntryErrors(index int, errs validation.Errors) validation.Errors {
	prefixed := validation.Errors{}
	for field, fieldErr := range errs {
		prefixed[fmt.Sprintf("entries.%d.%s", index, field)] = fieldErr
	}
	return prefixed
}

// validateRateSheetEntryInputs adds an error to validationErrors for each
// invalid entry of a new sheet or of a request's proposed rates: at least one
// entry, each role once, a whole rate of at least 1 and an overtime rate of at
// least 1 with at most two decimal places. prefix goes before each key, such as "proposed_rates.".
func validateRateSheetEntryInputs(entries []CreateRateSheetEntry, prefix string, validationErrors validation.Errors) {
	if len(entries) == 0 {
		setValidationError(validationErrors, prefix+"entries", validation.NewError("value_required", "at least one entry is required"))
	}
	seenRoles := map[string]bool{}
	for i, entry := range entries {
		entryPrefix := fmt.Sprintf("%sentries.%d.", prefix, i)
		if entry.Role == "" {
			setValidationError(validationErrors, entryPrefix+"role", validation.NewError("value_required", "entry role is required"))
		} else if seenRoles[entry.Role] {
			setValidationError(validationErrors, entryPrefix+"role", validation.NewError("validation_not_unique", "entry role must be unique"))
		} else {
			seenRoles[entry.Role] = true
		}
		if entry.Rate < 1 {
			setValidationError(validationErrors, entryPrefix+"rate", validation.NewError("validation_negative_number", "entry rate must be at least 1"))
		} else if entry.Rate != math.Floor(entry.Rate) {
			setValidationError(validationErrors, entryPrefix+"rate", validation.NewError("validation_not_whole_number", "entry rate must be a whole number"))
		}
		if entry.OvertimeRate < 1 {
			setValidationError(validationErrors, entryPrefix+"overtime_rate", validation.NewError("validation_negative_number", "entry overtime_rate must be at least 1"))
		} else if !hooks.HasAtMostTwoDecimals(entry.OvertimeRate) {
			setValidationError(validationErrors, entryPrefix+"overtime_rate", validation.NewError(hooks.OvertimeRateDecimalsCode, "entry "+hooks.OvertimeRateDecimalsMessage))
		}
	}
}

// createRateSheetHandler creates a handler that atomically creates a rate sheet
// and all its entries in a single transaction. This prevents partial-save states
// where the rate sheet exists but some entries are missing.
//
// Permission model:
//   - revision == 0 (new/copy): requires "job" claim
//   - revision > 0 (revise): requires "rate_sheet_revise" or "admin" claim
func createRateSheetHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		// Parse request body
		var req CreateRateSheetRequest
		if err := json.NewDecoder(e.Request.Body).Decode(&req); err != nil {
			return e.Error(http.StatusBadRequest, "invalid request body", err)
		}

		// Check claims based on revision
		if req.Revision == 0 {
			// New rate sheet or copy: requires "job" claim
			hasJobClaim, err := utilities.HasClaim(app, e.Auth, "job")
			if err != nil {
				return e.Error(http.StatusInternalServerError, "failed to check claims", err)
			}
			if !hasJobClaim {
				return e.Error(http.StatusForbidden, "you do not have permission to create rate sheets", nil)
			}
		} else {
			// Revision: requires "rate_sheet_revise" or "admin" claim
			hasReviseClaim, err := utilities.HasClaim(app, e.Auth, "rate_sheet_revise")
			if err != nil {
				return e.Error(http.StatusInternalServerError, "failed to check claims", err)
			}
			hasAdminClaim, err := utilities.HasClaim(app, e.Auth, "admin")
			if err != nil {
				return e.Error(http.StatusInternalServerError, "failed to check claims", err)
			}
			if !hasReviseClaim && !hasAdminClaim {
				return e.Error(http.StatusForbidden, "you do not have permission to revise rate sheets", nil)
			}
		}

		validationErrors := validation.Errors{}

		// The name is an optional description. A revision keeps its sheet's
		// name.
		req.Name = strings.TrimSpace(req.Name)

		// Validate effective_date format (YYYY-MM-DD)
		req.EffectiveDate = strings.TrimSpace(req.EffectiveDate)
		datePattern := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
		if req.EffectiveDate == "" {
			setValidationError(validationErrors, "effective_date", validation.NewError("value_required", "effective date is required"))
		} else if !datePattern.MatchString(req.EffectiveDate) {
			setValidationError(validationErrors, "effective_date", validation.NewError("validation_invalid_date", "effective date must be in YYYY-MM-DD format"))
		} else if err := utilities.IsValidDate(req.EffectiveDate); err != nil {
			setValidationError(validationErrors, "effective_date", err)
		}

		// Validate revision
		if req.Revision < 0 {
			setValidationError(validationErrors, "revision", validation.NewError("validation_negative_number", "revision must be 0 or greater"))
		}

		validateRateSheetEntryInputs(req.Entries, "", validationErrors)

		if len(validationErrors) > 0 {
			return apis.NewBadRequestError("Validation error", validationErrors)
		}

		// A revision takes the code of the sheet it revises.
		req.Code = strings.TrimSpace(req.Code)
		if req.Revision == 0 && req.Code != "" {
			return apis.NewBadRequestError("Validation error", validation.Errors{
				"code": validation.NewError("validation_not_allowed", "a new rate sheet gets its code when it is saved"),
			})
		}
		var revised *core.Record
		if req.Revision > 0 {
			var err error
			revised, err = revisedRateSheet(app, req.Code, req.Name)
			if err != nil {
				return e.Error(http.StatusInternalServerError, "failed to find the rate sheet to revise", err)
			}
			if revised == nil {
				return apis.NewBadRequestError("Validation error", validation.Errors{
					"code": validation.NewError("validation_missing_rate_sheet", "there is no rate sheet to revise"),
				})
			}
			// A request's proposed rates have no revisions until the job is
			// created and they become an ordinary sheet.
			if hooks.IsRateSheetDraft(revised) {
				return e.JSON(http.StatusBadRequest, hooks.RateSheetDraftError("name"))
			}
			req.Code = revised.GetString("code")
			// A revision keeps the name of its sheet.
			if req.Name == "" {
				req.Name = revised.GetString("name")
			}
		}

		// Validate revision effective date (for revisions, must be >= previous revision's date)
		if req.Revision > 0 {
			prevEffectiveDate, err := hooks.CheckRevisionEffectiveDate(app, req.Code, req.Revision, req.EffectiveDate)
			if err != nil {
				return e.Error(http.StatusInternalServerError, "failed to validate effective date", err)
			}
			if prevEffectiveDate != "" {
				return e.Error(http.StatusBadRequest, "effective date must be on or after the previous revision's effective date", map[string]any{
					"previous_effective_date": prevEffectiveDate,
				})
			}
		}

		// Run in transaction for atomicity
		var createdSheetId string
		var createdCode string
		var entriesCreated int

		err := app.RunInTransaction(func(txApp core.App) error {
			// Get the rate_sheets collection
			rateSheetsCollection, err := txApp.FindCollectionByNameOrId("rate_sheets")
			if err != nil {
				return err
			}

			if req.Revision > 0 {
				// Check for existing revision to provide a consistent conflict response
				existing, err := txApp.CountRecords("rate_sheets", dbx.HashExp{"code": req.Code, "revision": req.Revision})
				if err != nil {
					return err
				}
				if existing > 0 {
					return &uniqueConstraintError{
						message: "this rate sheet already has this revision",
						data:    rateSheetUniqueValidationErrors(),
					}
				}
			} else {
				// Two sheets can't have the same name, in any letter case.
				// Several sheets may have none. A request's proposed rates
				// hold their name too, and the message says so.
				holder, err := hooks.RateSheetNameHolder(txApp, req.Name, "")
				if err != nil {
					return err
				}
				if holder != nil {
					code, message := hooks.RateSheetNameHeld(holder)
					return &uniqueConstraintError{
						message: message,
						data:    validation.Errors{"name": validation.NewError(code, message)},
					}
				}
			}

			// Block a sheet whose rates are identical to another active sheet.
			// The revisions of a revised sheet are excluded, so a revision that
			// changes only the effective date is allowed.
			comparison, err := hooks.CompareWithActiveRateSheets(txApp, rateSheetEntryValues(req.Entries), req.Code)
			if err != nil {
				return err
			}
			if len(comparison.Exact) > 0 {
				return hooks.DuplicateRateSheetError("entries", comparison.Exact)
			}

			// Create the rate sheet record
			rateSheet := core.NewRecord(rateSheetsCollection)
			rateSheet.Set("name", req.Name)
			rateSheet.Set("effective_date", req.EffectiveDate)
			rateSheet.Set("revision", req.Revision)
			rateSheet.Set("active", false)
			// A revision takes the code of the sheet it revises. A new sheet
			// gets the next code when it is saved.
			rateSheet.Set("code", req.Code)

			if err := txApp.Save(rateSheet); err != nil {
				return err
			}

			createdSheetId = rateSheet.Id
			createdCode = rateSheet.GetString("code")

			// Get the rate_sheet_entries collection
			entriesCollection, err := txApp.FindCollectionByNameOrId("rate_sheet_entries")
			if err != nil {
				return err
			}

			// Create all entries
			for i, entry := range req.Entries {
				entryRecord := core.NewRecord(entriesCollection)
				entryRecord.Set("rate_sheet", createdSheetId)
				entryRecord.Set("role", entry.Role)
				entryRecord.Set("rate", int(entry.Rate))
				entryRecord.Set("overtime_rate", entry.OvertimeRate)

				if err := txApp.Save(entryRecord); err != nil {
					if ve, ok := asValidationErrors(err); ok {
						return &entryValidationError{index: i, errors: ve}
					}
					if strings.Contains(err.Error(), "UNIQUE constraint failed") {
						return &entryValidationError{
							index: i,
							errors: validation.Errors{
								"role": validation.NewError("validation_not_unique", "entry role must be unique"),
							},
						}
					}
					return err
				}
				entriesCreated++
			}

			return nil
		})

		if err != nil {
			// Check if it's our custom unique constraint error
			if uce, ok := err.(*uniqueConstraintError); ok {
				return apis.NewApiError(http.StatusConflict, uce.message, uce.data)
			}
			var hookErr *errs.HookError
			if errors.As(err, &hookErr) {
				return e.JSON(hookErr.Status, hookErr)
			}
			if eve, ok := err.(*entryValidationError); ok {
				return apis.NewBadRequestError("Validation error", prefixEntryErrors(eve.index, eve.errors))
			}
			if ve, ok := asValidationErrors(err); ok {
				if hasUniqueValidationError(ve) {
					return apis.NewApiError(http.StatusConflict, "a rate sheet with this code and revision already exists", ve)
				}
				return apis.NewBadRequestError("Validation error", ve)
			}
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				return apis.NewApiError(http.StatusConflict, "a rate sheet with this code and revision already exists", rateSheetUniqueValidationErrors())
			}
			return e.Error(http.StatusInternalServerError, "failed to create rate sheet", err)
		}

		return e.JSON(http.StatusCreated, map[string]any{
			"id":              createdSheetId,
			"code":            createdCode,
			"name":            req.Name,
			"revision":        req.Revision,
			"entries_created": entriesCreated,
		})
	}
}

// uniqueConstraintError is a custom error type for unique constraint violations
type uniqueConstraintError struct {
	message string
	data    validation.Errors
}

func (e *uniqueConstraintError) Error() string {
	return e.message
}

// rateSheetEntryValues converts request entries to the form that
// hooks.CompareWithActiveRateSheets takes.
func rateSheetEntryValues(entries []CreateRateSheetEntry) []hooks.RateSheetEntryValues {
	values := make([]hooks.RateSheetEntryValues, 0, len(entries))
	for _, entry := range entries {
		values = append(values, hooks.RateSheetEntryValues{
			Role:         entry.Role,
			Rate:         entry.Rate,
			OvertimeRate: entry.OvertimeRate,
		})
	}
	return values
}

// CompareRateSheetRequest is the body of POST /api/rate_sheets/compare. Code
// is the code of the sheet being revised; its active revision is not compared.
// A request without a code is matched to a sheet by name, as before codes
// existed. A new sheet sends neither.
type CompareRateSheetRequest struct {
	Code    string                 `json:"code"`
	Name    string                 `json:"name"`
	Entries []CreateRateSheetEntry `json:"entries"`
}

// createCompareRateSheetHandler compares unsaved entries with the active rate
// sheets. It returns the exact matches, which a save or an activation would
// block, and a few close sheets as information. It changes nothing. Only users
// who can save a sheet can call it: the 'job' claim creates new sheets and
// copies, and the 'rate_sheet_revise' or 'admin' claim creates revisions.
func createCompareRateSheetHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		allowed := false
		for _, claim := range []string{"job", "rate_sheet_revise", "admin"} {
			hasClaim, err := utilities.HasClaim(app, e.Auth, claim)
			if err != nil {
				return e.Error(http.StatusInternalServerError, "failed to check claims", err)
			}
			if hasClaim {
				allowed = true
				break
			}
		}
		if !allowed {
			return e.Error(http.StatusForbidden, "you do not have permission to compare rate sheets", nil)
		}

		var req CompareRateSheetRequest
		if err := json.NewDecoder(e.Request.Body).Decode(&req); err != nil {
			return e.Error(http.StatusBadRequest, "invalid request body", err)
		}

		excludeCode := ""
		revised, err := revisedRateSheet(app, strings.TrimSpace(req.Code), strings.TrimSpace(req.Name))
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to compare rate sheets", err)
		}
		if revised != nil {
			excludeCode = revised.GetString("code")
		}
		comparison, err := hooks.CompareWithActiveRateSheets(app, rateSheetEntryValues(req.Entries), excludeCode)
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to compare rate sheets", err)
		}

		return e.JSON(http.StatusOK, map[string]any{
			"exact": comparison.Exact,
			"near":  comparison.Closest(hooks.RateSheetNearMatchMaxRoles, hooks.RateSheetNearMatchLimit),
		})
	}
}

// revisedRateSheet returns the latest revision of the sheet with code, or,
// when code is empty, of the sheet named name, so that a client that sends
// only the name still finds the sheet it revises. It returns nil when there
// is none.
func revisedRateSheet(app core.App, code, name string) (*core.Record, error) {
	if code == "" {
		if name == "" {
			return nil, nil
		}
		named, err := app.FindRecordsByFilter("rate_sheets", "name = {:name}", "-revision", 1, 0, dbx.Params{"name": name})
		if err != nil || len(named) == 0 {
			return nil, err
		}
		code = named[0].GetString("code")
	}
	family, err := hooks.RateSheetFamily(app, code)
	if err != nil || len(family) == 0 {
		return nil, err
	}
	return family[len(family)-1], nil
}
