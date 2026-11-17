package hooks

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"tybalt/errs"

	"github.com/pocketbase/dbx"
	validation "github.com/pocketbase/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"
)

// Every rate sheet has a code, such as 2026-11-00004: the creation month (UTC)
// of its first revision and a five digit serial that starts at 00001 in each
// month. The code is the sheet's identity: all revisions of a sheet share it,
// and code and revision are unique together. The name is an optional
// description. See the Rate sheet codes section of
// descriptions/job_rate_sheet_history.md.
//
// rate_sheet_codes holds every code ever issued and is never cleared, so it is
// the per-month counter and a code is never issued twice, even after its sheet
// is deleted or absorbed. A code is issued inside the transaction that saves
// the sheet; PocketBase runs one write transaction at a time, and the unique
// indexes on rate_sheet_codes.code and rate_sheets (code, revision) are the
// backstop.
//
// Nobody chooses a code. A new sheet with no code gets the next one; a new
// revision takes the code of the sheet it revises. Only a PocketBase superuser
// can change a code, through the collection API, and only to a code that was
// never issued, in the same form, on a sheet without other revisions. Every
// such change is logged in rate_sheet_code_changes.

var rateSheetCodePattern = regexp.MustCompile(`^(\d{4}-(0[1-9]|1[0-2]))-(\d{5})$`)

// RateSheetCodeMaxSerial is the highest serial in a month.
const RateSheetCodeMaxSerial = 99999

// ParseRateSheetCode splits a code into its month (YYYY-MM) and serial. ok is
// false when code is not in the form YYYY-MM-serial or the serial is 0.
func ParseRateSheetCode(code string) (month string, serial int, ok bool) {
	match := rateSheetCodePattern.FindStringSubmatch(code)
	if match == nil {
		return "", 0, false
	}
	serial, err := strconv.Atoi(match[3])
	if err != nil || serial < 1 {
		return "", 0, false
	}
	return match[1], serial, true
}

// FormatRateSheetCode joins a month and a serial into a code.
func FormatRateSheetCode(month string, serial int) string {
	return fmt.Sprintf("%s-%05d", month, serial)
}

// RateSheetLabel is how the server names a sheet in labels and messages: the
// code and revision, and the name in quotes when the sheet has one, such as
// 2026-11-00004 (rev. 2) "Client Rates".
func RateSheetLabel(code, name string, revision int) string {
	label := fmt.Sprintf("%s (rev. %d)", code, revision)
	if name = strings.TrimSpace(name); name != "" {
		label += fmt.Sprintf(" %q", name)
	}
	return label
}

// RateSheetRecordLabel is RateSheetLabel for a record.
func RateSheetRecordLabel(sheet *core.Record) string {
	return RateSheetLabel(sheet.GetString("code"), sheet.GetString("name"), sheet.GetInt("revision"))
}

// IssueRateSheetCode issues the next code of the month of now (UTC) and
// records it in rate_sheet_codes. Call it inside the transaction that saves
// the sheet.
func IssueRateSheetCode(app core.App, now time.Time) (string, error) {
	month := now.UTC().Format("2006-01")
	var issued, used int
	if err := app.DB().NewQuery(`
		SELECT COALESCE(MAX(serial), 0) FROM rate_sheet_codes WHERE month = {:month}
	`).Bind(dbx.Params{"month": month}).Row(&issued); err != nil {
		return "", err
	}
	// A code on a sheet is always issued, but a sheet written outside
	// Turbo's hooks could carry one that is not; never issue it again.
	if err := app.DB().NewQuery(`
		SELECT COALESCE(MAX(CAST(SUBSTR(code, 9) AS INTEGER)), 0) FROM rate_sheets
		WHERE code LIKE {:prefix}
	`).Bind(dbx.Params{"prefix": month + "-%"}).Row(&used); err != nil {
		return "", err
	}
	serial := max(issued, used) + 1
	if serial > RateSheetCodeMaxSerial {
		return "", fmt.Errorf("no rate sheet code is left in %s", month)
	}
	code := FormatRateSheetCode(month, serial)
	if err := recordRateSheetCode(app, code); err != nil {
		return "", err
	}
	return code, nil
}

// recordRateSheetCode adds code to rate_sheet_codes. It fails when the code
// was issued before.
func recordRateSheetCode(app core.App, code string) error {
	month, serial, ok := ParseRateSheetCode(code)
	if !ok {
		return fmt.Errorf("invalid rate sheet code %q", code)
	}
	registry, err := app.FindCollectionByNameOrId("rate_sheet_codes")
	if err != nil {
		return err
	}
	issued := core.NewRecord(registry)
	issued.Set("code", code)
	issued.Set("month", month)
	issued.Set("serial", serial)
	if err := app.Save(issued); err != nil {
		return fmt.Errorf("recording rate sheet code %s: %w", code, err)
	}
	return nil
}

// RateSheetCodeIssued reports whether code was ever issued.
func RateSheetCodeIssued(app core.App, code string) (bool, error) {
	var count int
	err := app.DB().NewQuery(`SELECT COUNT(*) FROM rate_sheet_codes WHERE code = {:code}`).
		Bind(dbx.Params{"code": code}).Row(&count)
	return count > 0, err
}

// RateSheetFamily returns the revisions of the sheet with code, lowest
// revision first. It is empty when no sheet has the code.
func RateSheetFamily(app core.App, code string) ([]*core.Record, error) {
	if code == "" {
		return nil, nil
	}
	return app.FindRecordsByFilter("rate_sheets", "code = {:code}", "revision", 0, 0, dbx.Params{"code": code})
}

// RateSheetNameUse says which sheet holds a name.
type RateSheetNameUse struct {
	Code  string
	Draft bool
}

// RateSheetNameHolder returns the sheet of another code whose name equals
// name without regard to capitals or surrounding spaces, or nil when there is
// none. An empty name is never held: several sheets may have no name.
// exceptCode is the code of the sheet being named, whose revisions share its
// name; pass "" for a sheet that has no code yet.
func RateSheetNameHolder(app core.App, name, exceptCode string) (*RateSheetNameUse, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, nil
	}
	rows := []struct {
		Code    string `db:"code"`
		Name    string `db:"name"`
		Request string `db:"request"`
	}{}
	// SQLite LOWER changes only ASCII letters, so the names are compared
	// in Go. There are few sheets.
	if err := app.DB().NewQuery(`
		SELECT code, name, COALESCE(request, '') AS request FROM rate_sheets
		WHERE COALESCE(name, '') != '' AND code != {:code}
	`).Bind(dbx.Params{"code": exceptCode}).All(&rows); err != nil {
		return nil, err
	}
	for _, row := range rows {
		if strings.EqualFold(strings.TrimSpace(row.Name), name) {
			return &RateSheetNameUse{Code: row.Code, Draft: row.Request != ""}, nil
		}
	}
	return nil, nil
}

// RateSheetNameInUseError refuses a name that another sheet has. field is the
// key the UI shows the error under.
func RateSheetNameInUseError(field string) *errs.HookError {
	message := "another rate sheet already has this name"
	return &errs.HookError{
		Status:  http.StatusConflict,
		Message: message,
		Data: map[string]errs.CodeError{
			field: {Code: "validation_not_unique", Message: message},
		},
	}
}

func rateSheetCodeError(status int, message string) *errs.HookError {
	return &errs.HookError{
		Status:  status,
		Message: message,
		Data: map[string]errs.CodeError{
			"code": {Code: "invalid_rate_sheet_code", Message: message},
		},
	}
}

// restoringAbsorbedRateSheet reports whether code belongs to a sheet that a
// pending rate sheet absorb deleted, so that undo can recreate it with its
// code.
func restoringAbsorbedRateSheet(app core.App, code string) (bool, error) {
	var count int
	err := app.DB().NewQuery(`
		SELECT COUNT(*) FROM absorb_actions a, json_each(a.absorbed_records) j
		WHERE a.collection_name = 'rate_sheets'
			AND json_valid(a.absorbed_records)
			AND json_extract(j.value, '$.code') = {:code}
	`).Bind(dbx.Params{"code": code}).Row(&count)
	return count > 0, err
}

// rateSheetCodeEditMarker marks a record that a superuser saves through the
// collection API, which is the only way to change a code. Its private type
// cannot come from a request body.
const rateSheetCodeEditMarker = "@rateSheetCodeEdit"

type rateSheetCodeEditState struct{}

func superuserCodeEdit(record *core.Record) bool {
	_, ok := record.Get(rateSheetCodeEditMarker).(rateSheetCodeEditState)
	return ok
}

// checkNewRateSheet gives a new sheet its code, or checks the code it carries,
// and checks its name. A sheet with no code must be revision 0 and gets the
// next code. A sheet with a code is a new revision of the sheet with that
// code: it may not revise a draft, and it keeps the sheet's name. The only
// other sheet that may carry a code is one that the undo of a rate sheet
// absorb recreates.
func checkNewRateSheet(app core.App, record *core.Record) error {
	code := strings.TrimSpace(record.GetString("code"))
	record.Set("name", strings.TrimSpace(record.GetString("name")))
	name := record.GetString("name")

	if code == "" {
		if record.GetInt("revision") != 0 {
			return rateSheetCodeError(http.StatusBadRequest, "a new revision takes the code of the sheet it revises")
		}
		issued, err := IssueRateSheetCode(app, time.Now())
		if err != nil {
			return err
		}
		record.Set("code", issued)
		if holder, err := RateSheetNameHolder(app, name, issued); err != nil {
			return err
		} else if holder != nil {
			return RateSheetNameInUseError("name")
		}
		return nil
	}

	if _, _, ok := ParseRateSheetCode(code); !ok {
		return rateSheetCodeError(http.StatusBadRequest, "a rate sheet code has the form YYYY-MM-00001")
	}
	record.Set("code", code)
	restoring, err := restoringAbsorbedRateSheet(app, code)
	if err != nil {
		return err
	}
	family, err := RateSheetFamily(app, code)
	if err != nil {
		return err
	}
	if len(family) == 0 {
		if restoring {
			return nil
		}
		return rateSheetCodeError(http.StatusBadRequest, "a new rate sheet gets its code when it is saved")
	}
	if restoring {
		// The undo puts the sheet back as it was.
		return nil
	}
	for _, sheet := range family {
		if IsRateSheetDraft(sheet) {
			return RateSheetDraftError("code")
		}
	}
	latest := family[len(family)-1]
	if name != latest.GetString("name") {
		return &errs.HookError{
			Status:  http.StatusBadRequest,
			Message: "a revision keeps the name of its rate sheet; rename the sheet instead",
			Data: map[string]errs.CodeError{
				"name": {Code: "revision_name", Message: "a revision keeps the name of its rate sheet; rename the sheet instead"},
			},
		}
	}
	return nil
}

// checkRateSheetUpdate checks a change of code or name on a saved sheet. A
// code changes only in a superuser's save, to a code that was never issued, on
// a sheet without other revisions; changed reports whether it changed. A new
// name must not be another sheet's.
func checkRateSheetUpdate(app core.App, record *core.Record) (changed bool, oldCode string, err error) {
	oldCode = record.Original().GetString("code")
	newCode := strings.TrimSpace(record.GetString("code"))
	if newCode != oldCode {
		if !superuserCodeEdit(record) {
			return false, "", rateSheetCodeError(http.StatusForbidden, "only a superuser can change a rate sheet code")
		}
		if _, _, ok := ParseRateSheetCode(newCode); !ok {
			return false, "", rateSheetCodeError(http.StatusBadRequest, "a rate sheet code has the form YYYY-MM-00001")
		}
		issued, err := RateSheetCodeIssued(app, newCode)
		if err != nil {
			return false, "", err
		}
		if issued {
			return false, "", rateSheetCodeError(http.StatusConflict, "this code was issued before; a code is never used twice")
		}
		var others int
		if err := app.DB().NewQuery(`SELECT COUNT(*) FROM rate_sheets WHERE code = {:code} AND id != {:id}`).
			Bind(dbx.Params{"code": oldCode, "id": record.Id}).Row(&others); err != nil {
			return false, "", err
		}
		if others > 0 {
			return false, "", rateSheetCodeError(http.StatusBadRequest, "the code of a rate sheet with revisions can't change")
		}
		record.Set("code", newCode)
		changed = true
	}
	name := strings.TrimSpace(record.GetString("name"))
	if name != record.Original().GetString("name") {
		record.Set("name", name)
		holder, err := RateSheetNameHolder(app, name, newCode)
		if err != nil {
			return false, "", err
		}
		if holder != nil {
			return false, "", RateSheetNameInUseError("name")
		}
	}
	return changed, oldCode, nil
}

// logRateSheetCodeChange records a superuser's change of a code and issues the
// new code, so that it is never issued again.
func logRateSheetCodeChange(app core.App, sheet *core.Record, oldCode string) error {
	newCode := sheet.GetString("code")
	if err := recordRateSheetCode(app, newCode); err != nil {
		return err
	}
	changes, err := app.FindCollectionByNameOrId("rate_sheet_code_changes")
	if err != nil {
		return err
	}
	change := core.NewRecord(changes)
	change.Set("rate_sheet", sheet.Id)
	change.Set("old_code", oldCode)
	change.Set("new_code", newCode)
	return app.Save(change)
}

// addRateSheetCodeHooks gives every new sheet its code and guards changes of
// code and name, on every path that saves a sheet: the routes, the collection
// API and a superuser's edit.
func addRateSheetCodeHooks(app core.App) {
	app.OnRecordCreate("rate_sheets").BindFunc(func(e *core.RecordEvent) error {
		return e.App.RunInTransaction(func(tx core.App) error {
			if err := checkNewRateSheet(tx, e.Record); err != nil {
				return err
			}
			original := e.App
			e.App = tx
			defer func() { e.App = original }()
			return e.Next()
		})
	})
	app.OnRecordUpdate("rate_sheets").BindFunc(func(e *core.RecordEvent) error {
		return e.App.RunInTransaction(func(tx core.App) error {
			changed, oldCode, err := checkRateSheetUpdate(tx, e.Record)
			if err != nil {
				return err
			}
			original := e.App
			e.App = tx
			defer func() { e.App = original }()
			if err := e.Next(); err != nil {
				return err
			}
			if changed {
				return logRateSheetCodeChange(tx, e.Record, oldCode)
			}
			return nil
		})
	})
	// The collection API has no update rule, so only a superuser reaches
	// this hook.
	app.OnRecordUpdateRequest("rate_sheets").BindFunc(func(e *core.RecordRequestEvent) error {
		if e.HasSuperuserAuth() {
			e.Record.Set(rateSheetCodeEditMarker, rateSheetCodeEditState{})
			defer e.Record.Set(rateSheetCodeEditMarker, nil)
		}
		return e.Next()
	})
}

// rateSheetRequestError turns a HookError from the rate_sheets model hooks
// into the API error of a collection API request. PocketBase wraps an error of
// the save in a generic "Failed to ..." error; the HookError is its raw data.
func rateSheetRequestError(err error) error {
	if err == nil {
		return nil
	}
	var hookErr *errs.HookError
	if !errors.As(err, &hookErr) {
		var apiErr *router.ApiError
		if !errors.As(err, &apiErr) {
			return err
		}
		raw, ok := apiErr.RawData().(error)
		if !ok || !errors.As(raw, &hookErr) {
			return err
		}
	}
	data := validation.Errors{}
	for key, value := range hookErr.Data {
		data[key] = validation.NewError(value.Code, value.Message)
	}
	return apis.NewApiError(hookErr.Status, hookErr.Message, data)
}
