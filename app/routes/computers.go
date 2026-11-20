package routes

// The Computers pages: the list, the retired list, one computer, and the
// assign, retire and un-retire actions. The computers collection has no
// update rule, so the actions are server routes. A holder of the computers
// claim or of admin may use them. See descriptions/computers.md.

import (
	"database/sql"
	_ "embed" // for go:embed
	"errors"
	"io"
	"net/http"
	"strings"

	"tybalt/errs"
	"tybalt/utilities"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

//go:embed computers.sql
var computersQuery string

// computerRow is one computer as the pages show it. The names come from the
// profiles of the last user and of the assignee, which a holder of the
// computers claim can't read through the API.
type computerRow struct {
	ID              string        `db:"id" json:"id"`
	ComputerName    string        `db:"computer_name" json:"computer_name"`
	Serial          string        `db:"serial" json:"serial"`
	Mfg             string        `db:"mfg" json:"mfg"`
	Model           string        `db:"model" json:"model"`
	SystemType      float64       `db:"system_type" json:"system_type"`
	RAM             float64       `db:"ram" json:"ram"`
	OSVersion       string        `db:"os_version" json:"os_version"`
	OSSku           float64       `db:"os_sku" json:"os_sku"`
	OSArch          string        `db:"os_arch" json:"os_arch"`
	BootDrive       string        `db:"boot_drive" json:"boot_drive"`
	BootDriveFS     string        `db:"boot_drive_fs" json:"boot_drive_fs"`
	BootDriveCap    float64       `db:"boot_drive_cap" json:"boot_drive_cap"`
	BootDriveFree   float64       `db:"boot_drive_free" json:"boot_drive_free"`
	RadiatorVersion float64       `db:"radiator_version" json:"radiator_version"`
	NetworkConfig   types.JSONRaw `db:"network_config" json:"network_config"`
	UserGivenName   string        `db:"user_given_name" json:"user_given_name"`
	UserSurname     string        `db:"user_surname" json:"user_surname"`
	Email           string        `db:"email" json:"email"`
	UPN             string        `db:"upn" json:"upn"`
	LastUser        string        `db:"last_user" json:"last_user"`
	LastUserName    string        `db:"last_user_name" json:"last_user_name"`
	LastLogin       string        `db:"last_login" json:"last_login"`
	AssignedTo      string        `db:"assigned_to" json:"assigned_to"`
	AssignedToName  string        `db:"assigned_to_name" json:"assigned_to_name"`
	AssignedTime    string        `db:"assigned_time" json:"assigned_time"`
	Retired         string        `db:"retired" json:"retired"`
	Created         string        `db:"created" json:"created"`
	Updated         string        `db:"updated" json:"updated"`
}

// assignComputerRequest is the optional body of an assign. When user is
// given, the assign is refused unless it is still the last user, so a report
// that arrives while the page is open can't assign the computer to someone
// the person didn't see.
type assignComputerRequest struct {
	User string `json:"user"`
}

func bindComputerRoutes(app core.App, se *core.ServeEvent) {
	group := se.Router.Group("/api/computers")
	group.Bind(apis.RequireAuth("users"))
	group.BindFunc(func(e *core.RequestEvent) error {
		if err := requireComputersClaim(app, e.Auth); err != nil {
			return err
		}
		return e.Next()
	})
	group.GET("", createGetComputersHandler(app))
	group.GET("/{id}", createGetComputerHandler(app))
	group.POST("/{id}/assign", createAssignComputerHandler(app))
	group.POST("/{id}/retire", createRetireComputerHandler(app))
	group.POST("/{id}/unretire", createUnretireComputerHandler(app))
}

// requireComputersClaim allows a holder of the computers claim or of admin.
func requireComputersClaim(app core.App, auth *core.Record) error {
	for _, name := range []string{"computers", "admin"} {
		has, err := utilities.HasClaim(app, auth, name)
		if err != nil {
			return err
		}
		if has {
			return nil
		}
	}
	return apis.NewForbiddenError("you do not have permission to manage computers", nil)
}

// loadComputerRows runs computers.sql. A blank id lists the retired computers
// or the others.
func loadComputerRows(app core.App, id string, retired bool) ([]computerRow, error) {
	retiredParam := 0
	if retired {
		retiredParam = 1
	}
	rows := []computerRow{}
	err := app.DB().NewQuery(computersQuery).Bind(dbx.Params{
		"id":      id,
		"retired": retiredParam,
	}).All(&rows)
	return rows, err
}

func loadComputerRow(app core.App, id string) (*computerRow, error) {
	rows, err := loadComputerRows(app, id, false)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, sql.ErrNoRows
	}
	return &rows[0], nil
}

func createGetComputersHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		retired := e.Request.URL.Query().Get("retired")
		rows, err := loadComputerRows(app, "", retired == "1" || retired == "true")
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to load computers", err)
		}
		return e.JSON(http.StatusOK, rows)
	}
}

func createGetComputerHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		row, err := loadComputerRow(app, e.Request.PathValue("id"))
		if errors.Is(err, sql.ErrNoRows) {
			return e.JSON(http.StatusNotFound, computerNotFoundError())
		}
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to load the computer", err)
		}
		return e.JSON(http.StatusOK, row)
	}
}

// createAssignComputerHandler sets assigned_to to the last user of the
// computer and assigned_time to now. It refuses a retired computer and one
// whose last report matched no user.
func createAssignComputerHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		var req assignComputerRequest
		// The body is optional.
		if err := e.BindBody(&req); err != nil && !errors.Is(err, io.EOF) {
			return e.JSON(http.StatusBadRequest, computerError(http.StatusBadRequest, "invalid_request_body", "The request body is not valid."))
		}
		expected := strings.TrimSpace(req.User)

		return changeComputer(app, e, func(computer *core.Record) error {
			if computer.GetString("retired") != "" {
				return computerError(http.StatusBadRequest, "computer_retired",
					"This computer is retired. Un-retire it before you assign it.")
			}
			lastUser := computer.GetString("last_user")
			if lastUser == "" {
				return computerError(http.StatusBadRequest, "no_last_user",
					"The last sign-in on this computer matched no user in Turbo, so it can't be assigned.")
			}
			if expected != "" && expected != lastUser {
				return computerError(http.StatusConflict, "last_user_changed",
					"Someone else signed in to this computer since the page loaded. Reload the page and try again.")
			}
			computer.Set("assigned_to", lastUser)
			computer.Set("assigned_time", types.NowDateTime())
			return nil
		})
	}
}

// createRetireComputerHandler sets retired to now.
func createRetireComputerHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		return changeComputer(app, e, func(computer *core.Record) error {
			if computer.GetString("retired") != "" {
				return computerError(http.StatusBadRequest, "already_retired", "This computer is already retired.")
			}
			computer.Set("retired", types.NowDateTime())
			return nil
		})
	}
}

// createUnretireComputerHandler clears retired.
func createUnretireComputerHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		return changeComputer(app, e, func(computer *core.Record) error {
			if computer.GetString("retired") == "" {
				return computerError(http.StatusBadRequest, "not_retired", "This computer is not retired.")
			}
			computer.Set("retired", "")
			return nil
		})
	}
}

// changeComputer loads the computer of the request, lets change check and
// modify it, and saves it, in one transaction. It answers with the computer
// as the pages show it.
func changeComputer(app core.App, e *core.RequestEvent, change func(computer *core.Record) error) error {
	id := e.Request.PathValue("id")
	err := app.RunInTransaction(func(txApp core.App) error {
		computer, err := txApp.FindRecordById("computers", id)
		if errors.Is(err, sql.ErrNoRows) {
			return computerNotFoundError()
		}
		if err != nil {
			return err
		}
		if err := change(computer); err != nil {
			return err
		}
		return txApp.Save(computer)
	})
	if err != nil {
		var hookErr *errs.HookError
		if errors.As(err, &hookErr) {
			return e.JSON(hookErr.Status, hookErr)
		}
		return e.Error(http.StatusInternalServerError, "failed to save the computer", err)
	}

	row, err := loadComputerRow(app, id)
	if err != nil {
		return e.Error(http.StatusInternalServerError, "failed to load the computer", err)
	}
	return e.JSON(http.StatusOK, row)
}

func computerError(status int, code, message string) *errs.HookError {
	return &errs.HookError{
		Status:  status,
		Message: message,
		Data:    map[string]errs.CodeError{"global": {Code: code, Message: message}},
	}
}

func computerNotFoundError() *errs.HookError {
	return computerError(http.StatusNotFound, "not_found", "Computer not found.")
}
