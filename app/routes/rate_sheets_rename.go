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

	"tybalt/utilities"
)

// RenameRateSheetRequest is the body of POST /api/rate_sheets/{id}/rename.
type RenameRateSheetRequest struct {
	Name string `json:"name"`
}

// errRateSheetNameInUse is returned from the rename transaction when another
// sheet already has the new name.
var errRateSheetNameInUse = errors.New("rate sheet name in use")

// createRenameRateSheetHandler renames a rate sheet. Only admins can rename;
// rate_sheet_revise does not allow it.
//
// Revisions of a sheet are linked by name, so the handler changes name on
// every revision in one transaction and writes one rate_sheet_name_changes row
// whose anchor is the revision in the URL. No rates change and jobs keep their
// rate_sheet id, so nothing is repriced.
//
// Only current names are checked: a sheet may take a name that another sheet
// used before. The check ignores case, so two sheets can't differ only in
// capitals, but a sheet can change the capitals of its own name.
func createRenameRateSheetHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		id := e.Request.PathValue("id")

		hasAdminClaim, err := utilities.HasClaim(app, e.Auth, "admin")
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to check claims", err)
		}
		if !hasAdminClaim {
			return e.Error(http.StatusForbidden, "you do not have permission to rename rate sheets", nil)
		}

		var req RenameRateSheetRequest
		if err := json.NewDecoder(e.Request.Body).Decode(&req); err != nil {
			return e.Error(http.StatusBadRequest, "invalid request body", err)
		}
		newName := strings.TrimSpace(req.Name)
		if newName == "" {
			return apis.NewBadRequestError("Validation error", validation.Errors{
				"name": validation.NewError("value_required", "name is required"),
			})
		}

		var oldName string
		var renamed int
		err = app.RunInTransaction(func(txApp core.App) error {
			sheet, err := txApp.FindRecordById("rate_sheets", id)
			if err != nil {
				return err
			}
			oldName = sheet.GetString("name")
			if newName == oldName {
				return validation.Errors{
					"name": validation.NewError("name_unchanged", "the new name is the same as the current name"),
				}
			}

			var inUse int
			if err := txApp.DB().NewQuery(`
				SELECT COUNT(*) FROM rate_sheets
				WHERE LOWER(name) = LOWER({:newName}) AND name != {:oldName}
			`).Bind(dbx.Params{"newName": newName, "oldName": oldName}).Row(&inUse); err != nil {
				return err
			}
			if inUse > 0 {
				return errRateSheetNameInUse
			}

			revisions, err := txApp.FindRecordsByFilter("rate_sheets", "name = {:name}", "revision", 0, 0, dbx.Params{"name": oldName})
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
			if errors.Is(err, errRateSheetNameInUse) {
				return apis.NewApiError(http.StatusConflict, "another rate sheet already has this name", validation.Errors{
					"name": validation.NewError("validation_not_unique", "another rate sheet already has this name"),
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

		formerNames, err := RateSheetFormerNames(app, newName)
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to load former names", err)
		}

		return e.JSON(http.StatusOK, map[string]any{
			"id":                id,
			"name":              newName,
			"old_name":          oldName,
			"revisions_renamed": renamed,
			"former_names":      formerNames,
		})
	}
}

// createGetRateSheetFormerNamesHandler returns the earlier names of a rate
// sheet, newest first. Any signed-in user can read them.
func createGetRateSheetFormerNamesHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		id := e.Request.PathValue("id")
		sheet, err := app.FindRecordById("rate_sheets", id)
		if err != nil {
			return e.Error(http.StatusNotFound, "rate sheet not found", err)
		}
		name := sheet.GetString("name")
		formerNames, err := RateSheetFormerNames(app, name)
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to load former names", err)
		}
		return e.JSON(http.StatusOK, map[string]any{
			"id":           id,
			"name":         name,
			"former_names": formerNames,
		})
	}
}

// RateSheetFormerNames returns the earlier names of the sheet whose current
// name is name, newest first. They are the old_name values of the history rows
// whose anchor revision has that name now. The current name and repeats are
// left out, so a sheet renamed A to B and back to A shows only B.
func RateSheetFormerNames(app core.App, name string) ([]string, error) {
	var oldNames []string
	// rowid breaks ties between renames saved in the same millisecond.
	err := app.DB().NewQuery(`
		SELECT c.old_name
		FROM rate_sheet_name_changes c
		JOIN rate_sheets s ON s.id = c.rate_sheet
		WHERE s.name = {:name}
		ORDER BY c.created DESC, c.rowid DESC
	`).Bind(dbx.Params{"name": name}).Column(&oldNames)
	if err != nil {
		return nil, err
	}

	formerNames := []string{}
	seen := map[string]bool{name: true}
	for _, oldName := range oldNames {
		if seen[oldName] {
			continue
		}
		seen[oldName] = true
		formerNames = append(formerNames, oldName)
	}
	return formerNames, nil
}
