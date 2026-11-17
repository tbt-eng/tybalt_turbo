package routes

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/pocketbase/dbx"
	validation "github.com/pocketbase/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"

	"tybalt/errs"
	"tybalt/hooks"
	"tybalt/utilities"
)

// RenameRateSheetRequest is the body of POST /api/rate_sheets/{id}/rename.
type RenameRateSheetRequest struct {
	Name string `json:"name"`
}

// errRateSheetNameInUse is returned from the rename transaction when another
// sheet already has the new name. The handler names the holder's kind in the
// message.
var errRateSheetNameInUse = errors.New("rate sheet name in use")

// createRenameRateSheetHandler renames a rate sheet. Holders of admin or
// rate_sheet_revise can rename. The name is an optional description, so a
// rename to an empty name removes it.
//
// Revisions of a sheet share its code and its name, so the handler changes
// name on every revision with the code in one transaction and writes one
// rate_sheet_name_changes row whose anchor is the revision in the URL. The
// code does not change, no rates change and jobs keep their rate_sheet id, so
// nothing is repriced.
//
// Only current names are checked: a sheet may take a name that another sheet
// used before. The check ignores case, so two sheets can't differ only in
// capitals, but a sheet can change the capitals of its own name.
func createRenameRateSheetHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		id := e.Request.PathValue("id")

		allowed := false
		for _, claim := range []string{"admin", "rate_sheet_revise"} {
			has, err := utilities.HasClaim(app, e.Auth, claim)
			if err != nil {
				return e.Error(http.StatusInternalServerError, "failed to check claims", err)
			}
			allowed = allowed || has
		}
		if !allowed {
			return e.Error(http.StatusForbidden, "you do not have permission to rename rate sheets", nil)
		}

		var req RenameRateSheetRequest
		if err := json.NewDecoder(e.Request.Body).Decode(&req); err != nil {
			return e.Error(http.StatusBadRequest, "invalid request body", err)
		}
		newName := strings.TrimSpace(req.Name)

		var oldName string
		var code string
		var renamed int
		var heldBy *hooks.RateSheetNameUse
		err := app.RunInTransaction(func(txApp core.App) error {
			sheet, err := txApp.FindRecordById("rate_sheets", id)
			if err != nil {
				return err
			}
			// The requester names proposed rates in the request.
			if hooks.IsRateSheetDraft(sheet) {
				return hooks.RateSheetDraftError("name")
			}
			oldName = sheet.GetString("name")
			if newName == oldName {
				return validation.Errors{
					"name": validation.NewError("name_unchanged", "the new name is the same as the current name"),
				}
			}

			code = sheet.GetString("code")
			holder, err := hooks.RateSheetNameHolder(txApp, newName, code)
			if err != nil {
				return err
			}
			if holder != nil {
				heldBy = holder
				return errRateSheetNameInUse
			}

			revisions, err := hooks.RateSheetFamily(txApp, code)
			if err != nil {
				return err
			}
			for _, revision := range revisions {
				revision.Set("name", newName)
				if err := txApp.Save(revision); err != nil {
					return err
				}
			}
			renamed = len(revisions)

			changes, err := txApp.FindCollectionByNameOrId("rate_sheet_name_changes")
			if err != nil {
				return err
			}
			change := core.NewRecord(changes)
			change.Set("rate_sheet", sheet.Id)
			change.Set("old_name", oldName)
			change.Set("new_name", newName)
			change.Set("renamed_by", e.Auth.Id)
			return txApp.Save(change)
		})
		if err != nil {
			var hookErr *errs.HookError
			if errors.As(err, &hookErr) {
				return e.JSON(hookErr.Status, hookErr)
			}
			if errors.Is(err, errRateSheetNameInUse) {
				code, message := hooks.RateSheetNameHeld(heldBy)
				return apis.NewApiError(http.StatusConflict, message, validation.Errors{
					"name": validation.NewError(code, message),
				})
			}
			if ve, ok := asValidationErrors(err); ok {
				if hasUniqueValidationError(ve) {
					return apis.NewApiError(http.StatusConflict, "another rate sheet already has this name", ve)
				}
				return apis.NewBadRequestError("Validation error", ve)
			}
			if errors.Is(err, sql.ErrNoRows) {
				return e.Error(http.StatusNotFound, "rate sheet not found", err)
			}
			return e.Error(http.StatusInternalServerError, "failed to rename rate sheet", err)
		}

		formerNames, err := RateSheetFormerNames(app, code, newName)
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to load former names", err)
		}

		return e.JSON(http.StatusOK, map[string]any{
			"id":                id,
			"code":              code,
			"name":              newName,
			"old_name":          oldName,
			"revisions_renamed": renamed,
			"former_names":      formerNames,
		})
	}
}

// createGetRateSheetFormerNamesHandler returns the earlier names and the
// earlier codes of a rate sheet, newest first. Any signed-in user can read
// them.
func createGetRateSheetFormerNamesHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		id := e.Request.PathValue("id")
		sheet, err := app.FindRecordById("rate_sheets", id)
		if err != nil {
			return e.Error(http.StatusNotFound, "rate sheet not found", err)
		}
		name := sheet.GetString("name")
		code := sheet.GetString("code")
		formerNames, err := RateSheetFormerNames(app, code, name)
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to load former names", err)
		}
		formerCodes, err := RateSheetFormerCodes(app, code)
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to load former codes", err)
		}
		return e.JSON(http.StatusOK, map[string]any{
			"id":           id,
			"code":         code,
			"name":         name,
			"former_names": formerNames,
			"former_codes": formerCodes,
		})
	}
}

// RateSheetFormerNames returns the earlier names of the sheet with code,
// whose current name is name, newest first. They are the old_name values of
// the history rows whose anchor is a revision with that code. The current
// name, empty names and repeats are left out, so a sheet renamed A to B and
// back to A shows only B.
func RateSheetFormerNames(app core.App, code, name string) ([]string, error) {
	var oldNames []string
	// rowid breaks ties between renames saved in the same millisecond.
	err := app.DB().NewQuery(`
		SELECT c.old_name
		FROM rate_sheet_name_changes c
		JOIN rate_sheets s ON s.id = c.rate_sheet
		WHERE s.code = {:code}
		ORDER BY c.created DESC, c.rowid DESC
	`).Bind(dbx.Params{"code": code}).Column(&oldNames)
	if err != nil {
		return nil, err
	}

	formerNames := []string{}
	seen := map[string]bool{name: true, "": true}
	for _, oldName := range oldNames {
		if seen[oldName] {
			continue
		}
		seen[oldName] = true
		formerNames = append(formerNames, oldName)
	}
	return formerNames, nil
}

// RateSheetFormerCodes returns the earlier codes of the sheet with code,
// newest first: the old_code values of the code changes whose anchor has that
// code now. Printed documents and the legacy system keep an old code, and this
// is the way back from it to the sheet.
func RateSheetFormerCodes(app core.App, code string) ([]string, error) {
	var oldCodes []string
	err := app.DB().NewQuery(`
		SELECT c.old_code
		FROM rate_sheet_code_changes c
		JOIN rate_sheets s ON s.id = c.rate_sheet
		WHERE s.code = {:code}
		ORDER BY c.created DESC, c.rowid DESC
	`).Bind(dbx.Params{"code": code}).Column(&oldCodes)
	if err != nil {
		return nil, err
	}
	formerCodes := []string{}
	seen := map[string]bool{code: true}
	for _, oldCode := range oldCodes {
		if !seen[oldCode] {
			seen[oldCode] = true
			formerCodes = append(formerCodes, oldCode)
		}
	}
	return formerCodes, nil
}
