package routes

// The Logins and Raw Logins pages: the logins and the raw logins that
// POST /api/radiator writes, a page at a time, and the delete of one raw
// login. bindComputerRoutes binds them behind the same guard as the other
// computer routes: the computers claim or admin. See descriptions/computers.md.

import (
	"database/sql"
	_ "embed" // for go:embed
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"tybalt/errs"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

//go:embed computer_logins.sql
var computerLoginsQuery string

//go:embed computer_logins_count.sql
var computerLoginsCountQuery string

//go:embed computer_raw_logins.sql
var computerRawLoginsQuery string

//go:embed computer_raw_logins_count.sql
var computerRawLoginsCountQuery string

const (
	computerListDefaultPerPage = 20
	computerListMaxPerPage     = 100
)

// computerLoginRow is one sign-in as the pages show it. user_name is the
// matched user, which a holder of the computers claim can't read through the
// API.
type computerLoginRow struct {
	ID           string `db:"id" json:"id"`
	Created      string `db:"created" json:"created"`
	Computer     string `db:"computer" json:"computer"`
	ComputerName string `db:"computer_name" json:"computer_name"`
	Serial       string `db:"serial" json:"serial"`
	GivenName    string `db:"given_name" json:"given_name"`
	Surname      string `db:"surname" json:"surname"`
	User         string `db:"user" json:"user"`
	UserName     string `db:"user_name" json:"user_name"`
}

// computerRawLoginRow is one rejected report. The list leaves out report; the
// route for one raw login fills it.
type computerRawLoginRow struct {
	ID           string        `db:"id" json:"id"`
	Created      string        `db:"created" json:"created"`
	ComputerName string        `db:"computer_name" json:"computer_name"`
	Error        string        `db:"error" json:"error"`
	Report       types.JSONRaw `db:"-" json:"report,omitempty"`
}

// computerPage is one page of a list, in the shape of a PocketBase list
// response.
type computerPage[T any] struct {
	Items      []T `json:"items"`
	Page       int `json:"page"`
	PerPage    int `json:"perPage"`
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
}

// computerListPaging reads page and perPage. A missing or bad value gets the
// default; perPage is at most computerListMaxPerPage.
func computerListPaging(query url.Values) (page, perPage int) {
	page = 1
	if value, err := strconv.Atoi(query.Get("page")); err == nil && value > 0 {
		page = value
	}
	perPage = computerListDefaultPerPage
	if value, err := strconv.Atoi(query.Get("perPage")); err == nil && value > 0 {
		perPage = min(value, computerListMaxPerPage)
	}
	return page, perPage
}

// computerSearchPattern turns the search text into a LIKE pattern that
// matches it anywhere, with % _ and \ taken literally. Blank text gives "".
func computerSearchPattern(text string) string {
	text = strings.TrimSpace(text)
	if text == "" {
		return ""
	}
	escaped := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(text)
	return "%" + escaped + "%"
}

// loadComputerPage counts the rows, moves a page past the end back to the last
// page, and loads that page.
func loadComputerPage[T any](app core.App, listQuery, countQuery string, params dbx.Params, page, perPage int) (computerPage[T], error) {
	result := computerPage[T]{Items: []T{}, Page: page, PerPage: perPage}
	if err := app.DB().NewQuery(countQuery).Bind(params).Row(&result.TotalItems); err != nil {
		return result, err
	}
	result.TotalPages = (result.TotalItems + perPage - 1) / perPage
	if result.TotalPages > 0 && result.Page > result.TotalPages {
		result.Page = result.TotalPages
	}
	params["limit"] = perPage
	params["offset"] = (result.Page - 1) * perPage
	err := app.DB().NewQuery(listQuery).Bind(params).All(&result.Items)
	return result, err
}

// createGetComputerLoginsHandler lists sign-ins, newest first. ?computer=
// limits them to one computer and ?q= searches them.
func createGetComputerLoginsHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		query := e.Request.URL.Query()
		page, perPage := computerListPaging(query)
		result, err := loadComputerPage[computerLoginRow](app, computerLoginsQuery, computerLoginsCountQuery, dbx.Params{
			"computer": strings.TrimSpace(query.Get("computer")),
			"search":   computerSearchPattern(query.Get("q")),
		}, page, perPage)
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to load logins", err)
		}
		return e.JSON(http.StatusOK, result)
	}
}

// createGetComputerRawLoginsHandler lists rejected reports, newest first,
// without the reports. ?q= searches them.
func createGetComputerRawLoginsHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		query := e.Request.URL.Query()
		page, perPage := computerListPaging(query)
		result, err := loadComputerPage[computerRawLoginRow](app, computerRawLoginsQuery, computerRawLoginsCountQuery, dbx.Params{
			"search": computerSearchPattern(query.Get("q")),
		}, page, perPage)
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to load raw logins", err)
		}
		return e.JSON(http.StatusOK, result)
	}
}

// createGetComputerRawLoginHandler gives one rejected report with the report
// as it was received. The report is untrusted JSON from the network; the page
// shows it as text.
func createGetComputerRawLoginHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		record, err := app.FindRecordById("computer_raw_logins", e.Request.PathValue("id"))
		if errors.Is(err, sql.ErrNoRows) {
			return e.JSON(http.StatusNotFound, rawLoginNotFoundError())
		}
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to load the raw login", err)
		}
		row := computerRawLoginRow{
			ID:           record.Id,
			Created:      record.GetString("created"),
			ComputerName: record.GetString("computer_name"),
			Error:        record.GetString("error"),
		}
		if raw, ok := record.Get("report").(types.JSONRaw); ok && len(raw) > 0 {
			row.Report = raw
		} else {
			row.Report = types.JSONRaw("null")
		}
		return e.JSON(http.StatusOK, row)
	}
}

// createDeleteComputerRawLoginHandler deletes one rejected report. Nothing
// refers to a raw login, so the delete changes nothing else.
func createDeleteComputerRawLoginHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		record, err := app.FindRecordById("computer_raw_logins", e.Request.PathValue("id"))
		if errors.Is(err, sql.ErrNoRows) {
			return e.JSON(http.StatusNotFound, rawLoginNotFoundError())
		}
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to load the raw login", err)
		}
		if err := app.Delete(record); err != nil {
			return e.Error(http.StatusInternalServerError, "failed to delete the raw login", err)
		}
		return e.NoContent(http.StatusNoContent)
	}
}

func rawLoginNotFoundError() *errs.HookError {
	return computerError(http.StatusNotFound, "not_found", "Raw login not found.")
}
