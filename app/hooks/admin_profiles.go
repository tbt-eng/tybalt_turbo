package hooks

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"tybalt/errs"
	"tybalt/utilities"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

// EnforceAdminProfileRequestPermissions keeps direct collection creates
// admin-only. Limited editors update through the custom route layer.
func EnforceAdminProfileRequestPermissions(app core.App, e *core.RecordRequestEvent) error {
	hasAdminClaim, err := utilities.HasClaim(app, e.Auth, "admin")
	if err != nil {
		return err
	}
	if hasAdminClaim {
		return nil
	}

	return apis.NewForbiddenError("you do not have permission to edit admin profiles directly", nil)
}

// ProcessAdminProfile enforces business rules for admin_profiles create/update.
func ProcessAdminProfile(app core.App, e *core.RecordRequestEvent) error {
	record := e.Record
	uid := strings.TrimSpace(record.GetString("uid"))
	if uid == "" {
		return &errs.HookError{
			Status:  http.StatusBadRequest,
			Message: "admin profile validation error",
			Data: map[string]errs.CodeError{
				"uid": {
					Code:    "required",
					Message: "uid is required",
				},
			},
		}
	}

	openingDate := strings.TrimSpace(record.GetString("opening_date"))
	openingOP := record.GetFloat("opening_op")
	openingOV := record.GetFloat("opening_ov")
	if (openingOP != 0 || openingOV != 0) && openingDate == "" {
		return &errs.HookError{
			Status:  http.StatusBadRequest,
			Message: "admin profile validation error",
			Data: map[string]errs.CodeError{
				"opening_date": {
					Code:    "invalid_opening_date",
					Message: "opening_date is required when opening balances are non-zero",
				},
			},
		}
	}

	if openingDate != "" {
		if err := utilities.ValidateTimeOffOpeningDate(openingDate); err != nil {
			return &errs.HookError{
				Status:  http.StatusBadRequest,
				Message: "admin profile validation error",
				Data: map[string]errs.CodeError{
					"opening_date": {
						Code:    "invalid_opening_date",
						Message: "opening_date must be the Sunday after a pay period ending date",
					},
				},
			}
		}
	}

	return utilities.EnsureUserCanUseBranch(app, record.GetString("default_branch"), uid, "default_branch")
}

// RevokeTokensOnDeactivation signs a person out when their admin profile goes
// from active to inactive. It gives the user a new token key, which makes every
// token they hold invalid at once. Sign-in and token refresh already refuse an
// inactive account (see OnRecordAuthRequest in hooks.go), so without this the
// person could still use a token they already held until it expired.
//
// It is a model hook, so it covers every save of an admin profile: the active
// route, the admin profile editor and saves made in code.
func RevokeTokensOnDeactivation(e *core.RecordEvent) error {
	if !e.Record.Original().GetBool("active") || e.Record.GetBool("active") {
		return e.Next()
	}
	// The key changes before the profile saves. If the profile save then fails
	// outside a transaction, the person is only signed out.
	user, err := e.App.FindRecordById("users", e.Record.GetString("uid"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return e.Next()
		}
		return err
	}
	user.RefreshTokenKey()
	if err := e.App.SaveNoValidate(user); err != nil {
		return err
	}
	return e.Next()
}
