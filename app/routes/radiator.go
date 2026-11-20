package routes

// POST /api/radiator receives the login reports of the radiator script, which
// runs on each Windows computer at sign-in. It is a port of the rawLogins
// function of the legacy system. See descriptions/computers.md.
//
// The order of the checks matters. The method, the machine secret and the
// Content-Type are checked by radiatorGuard before any of the body is read.
// apis.BodyLimit then caps the body at radiatorMaxBodySize: it refuses a
// larger Content-Length outright and stops reading a body without one as soon
// as it passes the limit.

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"mime"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"tybalt/utilities"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/pocketbase/dbx"
	validation "github.com/pocketbase/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/hook"
	"github.com/pocketbase/pocketbase/tools/types"
)

const (
	radiatorPath          = "/api/radiator"
	radiatorMachineRole   = "radiator"
	radiatorMaxBodySize   = 64 << 10
	radiatorMaxErrorRunes = 4000
)

// radiatorRouteMethods are the methods the route answers. Only POST is
// accepted; the others get 405 with Allow: POST. They are registered so that
// they don't fall through to the static file route.
var radiatorRouteMethods = []string{
	http.MethodPost,
	http.MethodGet,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
}

// bindRadiatorRoutes registers POST /api/radiator. It takes no user auth.
func bindRadiatorRoutes(app core.App, se *core.ServeEvent) {
	handler := createRadiatorHandler(app)
	for _, method := range radiatorRouteMethods {
		se.Router.Route(method, radiatorPath, handler).
			Bind(radiatorGuard(app)).
			// Replaces the default limit of PocketBase (32 MB) for this route.
			Bind(apis.BodyLimit(radiatorMaxBodySize))
	}
}

// radiatorGuard checks the method, the machine secret and the Content-Type.
// It runs just before the body limit middleware, so a request that fails here
// is refused before any of its body is read.
func radiatorGuard(app core.App) *hook.Handler[*core.RequestEvent] {
	return &hook.Handler[*core.RequestEvent]{
		Id:       "radiatorGuard",
		Priority: apis.DefaultBodyLimitMiddlewarePriority - 1,
		Func: func(e *core.RequestEvent) error {
			if e.Request.Method != http.MethodPost {
				e.Response.Header().Set("Allow", http.MethodPost)
				return e.Error(http.StatusMethodNotAllowed, "method not allowed", nil)
			}

			authHeader := e.Request.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				return e.Error(http.StatusUnauthorized, "unauthorized", nil)
			}
			// TrimSpace handles trailing newlines from secret managers
			token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
			if token == "" || !utilities.ValidateMachineToken(app, token, radiatorMachineRole) {
				return e.Error(http.StatusUnauthorized, "unauthorized", nil)
			}

			mediaType, _, err := mime.ParseMediaType(e.Request.Header.Get("Content-Type"))
			if err != nil || mediaType != "application/json" {
				return e.Error(http.StatusUnsupportedMediaType, "Content-Type must be application/json", nil)
			}

			return e.Next()
		},
	}
}

func createRadiatorHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		body, err := io.ReadAll(e.Request.Body)
		if err != nil {
			if errors.Is(err, apis.ErrRequestEntityTooLarge) {
				return apis.ErrRequestEntityTooLarge
			}
			return e.Error(http.StatusBadRequest, "could not read the request body", nil)
		}
		// Old PowerShell may start the body with a UTF-8 byte order mark.
		body = bytes.TrimPrefix(body, []byte("\xef\xbb\xbf"))

		if len(bytes.TrimSpace(body)) == 0 {
			return e.Error(http.StatusBadRequest, "empty submission, not saving", nil)
		}

		var parsed any
		decoder := json.NewDecoder(bytes.NewReader(body))
		decoder.UseNumber()
		if err := decoder.Decode(&parsed); err != nil {
			return e.Error(http.StatusBadRequest, "the body is not valid JSON", nil)
		}
		if _, err := decoder.Token(); err != io.EOF {
			return e.Error(http.StatusBadRequest, "the body is not valid JSON", nil)
		}
		object, ok := parsed.(map[string]any)
		if !ok {
			if parsed == nil {
				return e.Error(http.StatusBadRequest, "empty submission, not saving", nil)
			}
			if list, isList := parsed.([]any); isList && len(list) == 0 {
				return e.Error(http.StatusBadRequest, "empty submission, not saving", nil)
			}
			return e.Error(http.StatusBadRequest, "the report must be a JSON object", nil)
		}
		if len(object) == 0 {
			return e.Error(http.StatusBadRequest, "empty submission, not saving", nil)
		}

		report, problems := validateRadiatorReport(object)
		if len(problems) == 0 {
			err = storeRadiatorLogin(app, report)
			var invalid *radiatorInvalidError
			if errors.As(err, &invalid) {
				problems = []string{invalid.Error()}
			} else if err != nil {
				app.Logger().Error("radiator: failed to store a login", "error", err)
				return e.Error(http.StatusInternalServerError, "failed to store the report", nil)
			}
		}

		if len(problems) > 0 {
			if err := storeRadiatorRawLogin(app, body, report, problems); err != nil {
				app.Logger().Error("radiator: failed to store a raw login", "error", err)
				return e.Error(http.StatusInternalServerError, "failed to store the report", nil)
			}
		}

		return e.NoContent(http.StatusAccepted)
	}
}

// radiatorInvalidError marks a report that passed the schema but still can't
// be stored, such as a serial that is too short once cleaned. It is stored as
// a raw login, not answered with 500.
type radiatorInvalidError struct{ msg string }

func (e *radiatorInvalidError) Error() string { return e.msg }

// storeRadiatorLogin writes a valid report in one transaction: it upserts the
// computer, adds a login and deletes the raw logins of the computer name.
func storeRadiatorLogin(app core.App, report map[string]any) error {
	serial, _ := report["serial"].(string)
	mfg, _ := report["mfg"].(string)
	slug, err := makeSlug(serial, mfg)
	if err != nil {
		return &radiatorInvalidError{msg: err.Error()}
	}

	anchor := strings.ToLower(report["userSourceAnchor"].(string))
	upn := strings.ToLower(report["upn"].(string))
	email := ""
	if value, ok := report["email"].(string); ok {
		email = strings.ToLower(value)
	}
	givenName, _ := report["userGivenName"].(string)
	surname, _ := report["userSurname"].(string)
	computerName, _ := report["computerName"].(string)

	err = app.RunInTransaction(func(txApp core.App) error {
		uid, err := resolveRadiatorUser(txApp, anchor, email, upn)
		if err != nil {
			return err
		}

		computer, err := txApp.FindFirstRecordByData("computers", "slug", slug)
		if errors.Is(err, sql.ErrNoRows) {
			collection, err := txApp.FindCollectionByNameOrId("computers")
			if err != nil {
				return err
			}
			computer = core.NewRecord(collection)
			computer.Set("slug", slug)
		} else if err != nil {
			return err
		}

		// The computer fields are merged as legacy did: a field missing from
		// the report keeps its stored value. assigned_to, assigned_time and
		// retired are never written here.
		for _, field := range radiatorComputerFields {
			if value, ok := report[field.json]; ok {
				computer.Set(field.column, value)
			}
		}
		// The user fields describe the last report, so they are all replaced.
		computer.Set("user_given_name", givenName)
		computer.Set("user_surname", surname)
		computer.Set("email", email)
		computer.Set("upn", upn)
		computer.Set("user_source_anchor", anchor)
		computer.Set("last_user", uid)
		computer.Set("last_login", types.NowDateTime())
		if err := txApp.Save(computer); err != nil {
			return err
		}

		loginsCollection, err := txApp.FindCollectionByNameOrId("computer_logins")
		if err != nil {
			return err
		}
		login := core.NewRecord(loginsCollection)
		login.Set("computer", computer.Id)
		login.Set("given_name", givenName)
		login.Set("surname", surname)
		login.Set("user_source_anchor", anchor)
		login.Set("user", uid)
		if err := txApp.Save(login); err != nil {
			return err
		}

		// A valid report makes the earlier rejects of the computer redundant.
		if computerName != "" {
			if _, err := txApp.DB().NewQuery(
				"DELETE FROM computer_raw_logins WHERE computer_name = {:name}",
			).Bind(dbx.Params{"name": computerName}).Execute(); err != nil {
				return err
			}
		}
		return nil
	})

	// A value the schema allows but a field refuses (a very long UPN, say)
	// makes the report invalid, not the server broken.
	var fieldErrors validation.Errors
	if errors.As(err, &fieldErrors) {
		return &radiatorInvalidError{msg: "the report could not be saved: " + fieldErrors.Error()}
	}
	return err
}

// resolveRadiatorUser finds the user of a report: first the admin profile
// with the source anchor, then a user whose email equals the report's email
// or UPN, ignoring case. It returns "" when nothing matches; it never creates
// a user.
func resolveRadiatorUser(app core.App, anchor, email, upn string) (string, error) {
	var row struct {
		UID string `db:"uid"`
	}

	if anchor != "" {
		err := app.DB().NewQuery(
			"SELECT uid FROM admin_profiles WHERE user_source_anchor = {:anchor} LIMIT 1",
		).Bind(dbx.Params{"anchor": anchor}).One(&row)
		if err == nil && row.UID != "" {
			return row.UID, nil
		}
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return "", err
		}
	}

	for _, address := range []string{email, upn} {
		if address == "" {
			continue
		}
		row.UID = ""
		err := app.DB().NewQuery(
			"SELECT id AS uid FROM users WHERE lower(email) = {:address} ORDER BY id LIMIT 1",
		).Bind(dbx.Params{"address": address}).One(&row)
		if err == nil && row.UID != "" {
			return row.UID, nil
		}
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return "", err
		}
	}

	return "", nil
}

// storeRadiatorRawLogin stores a rejected report as it was received, with the
// reasons. computer_name comes from the report when it is valid.
func storeRadiatorRawLogin(app core.App, body []byte, report map[string]any, problems []string) error {
	collection, err := app.FindCollectionByNameOrId("computer_raw_logins")
	if err != nil {
		return err
	}
	computerName, _ := report["computerName"].(string)
	message := strings.Join(problems, "; ")
	if utf8.RuneCountInString(message) > radiatorMaxErrorRunes {
		message = string([]rune(message)[:radiatorMaxErrorRunes])
	}

	record := core.NewRecord(collection)
	record.Set("report", types.JSONRaw(body))
	record.Set("computer_name", computerName)
	record.Set("error", message)
	return app.Save(record)
}

// radiatorComputerFields maps the computer fields of a report to columns of
// computers.
var radiatorComputerFields = []struct{ json, column string }{
	{"serial", "serial"},
	{"mfg", "mfg"},
	{"model", "model"},
	{"computerName", "computer_name"},
	{"systemType", "system_type"},
	{"ram", "ram"},
	{"osVersion", "os_version"},
	{"osSku", "os_sku"},
	{"osArch", "os_arch"},
	{"bootDrive", "boot_drive"},
	{"bootDriveFS", "boot_drive_fs"},
	{"bootDriveCap", "boot_drive_cap"},
	{"bootDriveFree", "boot_drive_free"},
	{"radiatorVersion", "radiator_version"},
	{"networkConfig", "network_config"},
}

// ---------------------------------------------------------------------------
// Validation: a port of RawLogins.schema.json as ajv 6 applied it in legacy,
// with removeAdditional, coerceTypes and allErrors, and the custom
// removeIfFails keyword.

type radiatorKind int

const (
	radiatorString radiatorKind = iota
	radiatorInteger
	radiatorObject
)

type radiatorProperty struct {
	name      string
	kind      radiatorKind
	minLength int // -1: no limit
	maxLength int // -1: no limit
	pattern   *regexp.Regexp
	email     bool
	// removeIfFails drops the property when it fails, instead of failing the
	// report.
	removeIfFails bool
}

// jsLineChar is what . matches in a JavaScript regular expression.
const jsLineChar = `[^\n\r\x{2028}\x{2029}]`

var (
	radiatorMACPattern = regexp.MustCompile(`^([0-9A-Fa-f]{2}:){5}[0-9A-Fa-f]{2}$`)
	// The "fast" email format of ajv 6.
	radiatorEmailPattern = regexp.MustCompile("(?i)^[a-z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?(?:\\.[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?)*$")
	// A JavaScript decimal literal, as Number() reads it.
	jsDecimalPattern = regexp.MustCompile(`^[+-]?(?:[0-9]+\.?[0-9]*|\.[0-9]+)(?:[eE][+-]?[0-9]+)?$`)
)

// radiatorProperties are the properties of the schema, in its order.
var radiatorProperties = []radiatorProperty{
	{name: "bootDriveFS", kind: radiatorString, minLength: 2, maxLength: 32},
	{name: "serial", kind: radiatorString, minLength: 4, maxLength: 48},
	{name: "computerName", kind: radiatorString, minLength: 1, maxLength: 48, removeIfFails: true},
	{name: "model", kind: radiatorString, minLength: 0, maxLength: 48},
	{name: "systemType", kind: radiatorInteger},
	{name: "osSku", kind: radiatorInteger},
	{name: "bootDriveFree", kind: radiatorInteger},
	{name: "userGivenName", kind: radiatorString, minLength: 1, maxLength: 24},
	{name: "upn", kind: radiatorString, minLength: -1, maxLength: -1, pattern: regexp.MustCompile(`^(` + jsLineChar + `+?@` + jsLineChar + `+?)$`)},
	{name: "userSourceAnchor", kind: radiatorString, minLength: -1, maxLength: -1, pattern: regexp.MustCompile(`^([a-fA-F0-9]{32})$`)},
	{name: "osArch", kind: radiatorString, minLength: 2, maxLength: 16},
	{name: "bootDrive", kind: radiatorString, minLength: -1, maxLength: -1, pattern: regexp.MustCompile(`^[A-Za-z]:$`)},
	{name: "osVersion", kind: radiatorString, minLength: -1, maxLength: -1, pattern: regexp.MustCompile(`^[0-9]+\.[0-9]+(\.[0-9]+)*$`)},
	{name: "radiatorVersion", kind: radiatorInteger},
	{name: "mfg", kind: radiatorString, minLength: 2, maxLength: 32},
	{name: "email", kind: radiatorString, minLength: -1, maxLength: -1, email: true, removeIfFails: true},
	{name: "ram", kind: radiatorInteger},
	{name: "userSurname", kind: radiatorString, minLength: 2, maxLength: 35},
	{name: "bootDriveCap", kind: radiatorInteger},
	{name: "networkConfig", kind: radiatorObject},
}

// radiatorRequired are the required properties of the schema, in its order.
var radiatorRequired = []string{
	"serial",
	"upn",
	"networkConfig",
	"userSourceAnchor",
	"radiatorVersion",
	"mfg",
	"userGivenName",
	"userSurname",
}

// validateRadiatorReport returns the report with unknown properties dropped,
// types coerced and failing removeIfFails properties dropped, and the reasons
// it is invalid. An empty list means the report is valid.
func validateRadiatorReport(input map[string]any) (map[string]any, []string) {
	report := map[string]any{}
	var problems []string

	for _, property := range radiatorProperties {
		raw, present := input[property.name]
		if !present {
			continue
		}
		value, propertyProblems := validateRadiatorProperty(property, raw)
		if len(propertyProblems) > 0 {
			if property.removeIfFails {
				continue
			}
			problems = append(problems, propertyProblems...)
			// ajv coerces in place even when a later keyword fails.
			report[property.name] = value
			continue
		}
		report[property.name] = value
	}

	for _, name := range radiatorRequired {
		if _, present := input[name]; !present {
			problems = append(problems, name+": is required")
		}
	}

	return report, problems
}

func validateRadiatorProperty(property radiatorProperty, raw any) (any, []string) {
	name := property.name
	switch property.kind {
	case radiatorString:
		value, ok := coerceToJSString(raw)
		if !ok {
			return raw, []string{name + ": must be a string"}
		}
		var problems []string
		length := utf8.RuneCountInString(value)
		if property.minLength > 0 && length < property.minLength {
			problems = append(problems, fmt.Sprintf("%s: must be at least %d characters", name, property.minLength))
		}
		if property.maxLength >= 0 && length > property.maxLength {
			problems = append(problems, fmt.Sprintf("%s: must be at most %d characters", name, property.maxLength))
		}
		if property.pattern != nil && !property.pattern.MatchString(value) {
			problems = append(problems, fmt.Sprintf("%s: must match the pattern %s", name, property.pattern.String()))
		}
		if property.email && !radiatorEmailPattern.MatchString(value) {
			problems = append(problems, name+": must be an email address")
		}
		return value, problems

	case radiatorInteger:
		value, ok := coerceToJSInteger(raw)
		if !ok {
			return raw, []string{name + ": must be an integer"}
		}
		return value, nil

	case radiatorObject:
		value, ok := raw.(map[string]any)
		if !ok {
			return raw, []string{name + ": must be an object"}
		}
		var problems []string
		for key := range value {
			if !radiatorMACPattern.MatchString(key) {
				problems = append(problems, fmt.Sprintf("%s: the property name %q must be a MAC address", name, key))
			}
		}
		if len(value) < 1 {
			problems = append(problems, name+": must have at least 1 property")
		}
		return value, problems
	}
	return raw, []string{name + ": unknown type"}
}

// coerceToJSString converts a JSON value to a string as ajv coerceTypes does:
// a number or boolean becomes its JavaScript text, null becomes "".
func coerceToJSString(raw any) (string, bool) {
	switch value := raw.(type) {
	case string:
		return value, true
	case nil:
		return "", true
	case bool:
		if value {
			return "true", true
		}
		return "false", true
	case json.Number:
		f, err := strconv.ParseFloat(value.String(), 64)
		if err != nil {
			return "", false
		}
		return jsNumberString(f), true
	}
	return "", false
}

// coerceToJSInteger converts a JSON value to an integer as ajv coerceTypes
// does: a boolean becomes 1 or 0, null becomes 0, and a string that
// JavaScript reads as a whole number becomes that number.
func coerceToJSInteger(raw any) (float64, bool) {
	var f float64
	switch value := raw.(type) {
	case json.Number:
		parsed, err := strconv.ParseFloat(value.String(), 64)
		if err != nil {
			return 0, false
		}
		f = parsed
	case bool:
		if value {
			return 1, true
		}
		return 0, true
	case nil:
		return 0, true
	case string:
		if value == "" {
			return 0, false
		}
		parsed, ok := jsStringToNumber(value)
		if !ok {
			return 0, false
		}
		f = parsed
	default:
		return 0, false
	}
	if math.IsInf(f, 0) || math.IsNaN(f) || f != math.Trunc(f) {
		return 0, false
	}
	return f, true
}

// jsStringToNumber reads a string as JavaScript Number() does. Infinity is
// refused.
func jsStringToNumber(s string) (float64, bool) {
	s = strings.TrimFunc(s, jsIsSpace)
	if s == "" {
		return 0, true
	}
	if len(s) > 2 && s[0] == '0' {
		base := 0
		switch s[1] {
		case 'x', 'X':
			base = 16
		case 'o', 'O':
			base = 8
		case 'b', 'B':
			base = 2
		}
		if base != 0 {
			digits := s[2:]
			if strings.ContainsAny(digits, "+-_") {
				return 0, false
			}
			n, err := strconv.ParseUint(digits, base, 64)
			if err != nil {
				return 0, false
			}
			return float64(n), true
		}
	}
	if !jsDecimalPattern.MatchString(s) {
		return 0, false
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil && !errors.Is(err, strconv.ErrRange) {
		return 0, false
	}
	return f, true
}

// jsNumberString formats a number as JavaScript String() does.
func jsNumberString(f float64) string {
	if f == 0 {
		return "0"
	}
	abs := math.Abs(f)
	if abs >= 1e-6 && abs < 1e21 {
		return strconv.FormatFloat(f, 'f', -1, 64)
	}
	mantissa, exponent, _ := strings.Cut(strconv.FormatFloat(f, 'e', -1, 64), "e")
	sign, digits := exponent[:1], strings.TrimLeft(exponent[1:], "0")
	return mantissa + "e" + sign + digits
}

// jsIsSpace reports whether r is white space or a line terminator to
// JavaScript (\s, String.prototype.trim and Number()).
func jsIsSpace(r rune) bool {
	switch r {
	case '\t', '\n', '\v', '\f', '\r', ' ', 0x00A0, 0x1680, 0x2028, 0x2029, 0x202F, 0x205F, 0x3000, 0xFEFF:
		return true
	}
	return r >= 0x2000 && r <= 0x200A
}

var slugMfgRemovePattern = regexp.MustCompile(`(?i)/|\.|,|inc|ltd`)

// makeSlug returns the key of a computer from its serial and manufacturer. It
// is a port of makeSlug in the legacy utilities.ts: the serial loses white
// space, slashes and commas; the manufacturer is lowercased, loses slashes,
// periods, commas, "inc" and "ltd", is trimmed and has spaces replaced by
// underscores. Lengths are counted in UTF-16 units, as JavaScript does.
func makeSlug(serial, mfg string) (string, error) {
	cleanSerial := strings.Map(func(r rune) rune {
		if jsIsSpace(r) || r == '/' || r == ',' {
			return -1
		}
		return r
	}, serial)
	cleanMfg := strings.ToLower(mfg)
	cleanMfg = slugMfgRemovePattern.ReplaceAllString(cleanMfg, "")
	cleanMfg = strings.TrimFunc(cleanMfg, jsIsSpace)
	cleanMfg = strings.ReplaceAll(cleanMfg, " ", "_")

	if utf16Length(cleanSerial) >= 4 && utf16Length(cleanMfg) >= 2 {
		return cleanSerial + "," + cleanMfg, nil
	}
	return "", fmt.Errorf("serial %s or manufacturer %s too short", cleanSerial, cleanMfg)
}

func utf16Length(s string) int {
	return len(utf16.Encode([]rune(s)))
}
