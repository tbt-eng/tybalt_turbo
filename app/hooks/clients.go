package hooks

import (
	"net/http"

	"tybalt/errs"
	"tybalt/utilities"

	"github.com/pocketbase/pocketbase/core"
)

// validateClient performs cross-field validation for the client record.
func validateClient(app core.App, clientRecord *core.Record) error {
	// The lead is optional, and the UI does not offer it for now. It is checked
	// only when it is set or changed, so a client whose lead has since been
	// deactivated can still be saved with that lead left as it is.
	leadID := clientRecord.GetString("business_development_lead")
	if leadID == "" {
		return nil
	}
	if !clientRecord.IsNew() && clientRecord.Original().GetString("business_development_lead") == leadID {
		return nil
	}

	// Check that the business development lead is an active user
	active, err := utilities.IsUserActive(app, leadID)
	if err != nil {
		return &errs.HookError{
			Status:  http.StatusInternalServerError,
			Message: "failed to check business development lead active status",
		}
	}
	if !active {
		return &errs.HookError{
			Status:  http.StatusBadRequest,
			Message: "business development lead must be an active user",
			Data: map[string]errs.CodeError{
				"business_development_lead": {
					Code:    "business_development_lead_not_active",
					Message: "the selected business development lead is not an active user",
				},
			},
		}
	}

	return nil
}

// ProcessClient enforces business rules for client create/update.
func ProcessClient(app core.App, e *core.RecordRequestEvent) error {
	clientRecord := e.Record

	if err := validateClient(app, clientRecord); err != nil {
		return err
	}

	return nil
}
