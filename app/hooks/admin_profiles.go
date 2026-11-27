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
// route, the admin profile editor and saves made in code. A profile whose uid
// changes no longer makes its old user active, so that user is signed out too.
// The API and the save routes refuse a uid change, but a superuser can make one.
//
// The key changes after the profile is written, so a sign-in or token refresh
// that runs between the two writes already sees the inactive profile and is
// refused. Inside a transaction both writes commit or roll back together. The
// collection API saves without a transaction; if the key change fails there,
// the save returns the error but the profile stays inactive.
func RevokeTokensOnDeactivation(e *core.RecordEvent) error {
	original := e.Record.Original()
	uid := original.GetString("uid")
	if !original.GetBool("active") || (e.Record.GetBool("active") && e.Record.GetString("uid") == uid) {
		return e.Next()
	}
	if err := e.Next(); err != nil {
		return err
	}
	return revokeUserTokens(e.App, uid)
}

// RevokeTokensOnInactiveProfileCreate signs a person out when an admin profile
// is created for them as inactive. A user with no admin profile already counts
// as inactive and cannot sign in with a password, so this matters only for a
// token they got another way, such as a superuser's impersonation. A profile
// created as active changes nothing.
func RevokeTokensOnInactiveProfileCreate(e *core.RecordEvent) error {
	if e.Record.GetBool("active") {
		return e.Next()
	}
	if err := e.Next(); err != nil {
		return err
	}
	return revokeUserTokens(e.App, e.Record.GetString("uid"))
}

// RevokeTokensOnProfileDelete signs a person out when their admin profile is
// deleted, because a user with no admin profile counts as inactive. Only a
// superuser can delete a profile. A user who has a profile cannot be deleted,
// because the profile's required uid relation does not cascade.
func RevokeTokensOnProfileDelete(e *core.RecordEvent) error {
	if err := e.Next(); err != nil {
		return err
	}
	return revokeUserTokens(e.App, e.Record.GetString("uid"))
}

// revokeUserTokens gives the user a new token key. A missing user holds no
// token, so it is not an error.
func revokeUserTokens(app core.App, uid string) error {
	if uid == "" {
		return nil
	}
	user, err := app.FindRecordById("users", uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}
	user.RefreshTokenKey()
	return app.SaveNoValidate(user)
}
