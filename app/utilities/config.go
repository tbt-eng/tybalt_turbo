package utilities

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"tybalt/constants"
	"tybalt/errs"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

// GetConfigValue retrieves a config domain record and returns the parsed JSON value.
// Returns nil if the domain key doesn't exist.
func GetConfigValue(app core.App, domainKey string) (map[string]any, error) {
	record, err := app.FindFirstRecordByData("app_config", "key", domainKey)
	if err != nil {
		return nil, nil // Domain not found
	}

	valueStr := record.GetString("value")
	var value map[string]any
	if err := json.Unmarshal([]byte(valueStr), &value); err != nil {
		return nil, err
	}
	return value, nil
}

// GetConfigBool retrieves a boolean from a domain config.
// domainKey is the row key (e.g., "jobs"), property is the JSON property (e.g., "create_edit_absorb")
func GetConfigBool(app core.App, domainKey string, property string, defaultValue bool) (bool, error) {
	config, err := GetConfigValue(app, domainKey)
	if err != nil {
		return defaultValue, err
	}
	if config == nil {
		return defaultValue, nil
	}

	val, ok := config[property]
	if !ok {
		return defaultValue, nil
	}
	boolVal, ok := val.(bool)
	if !ok {
		return defaultValue, nil
	}
	return boolVal, nil
}

// InDeployGroup says whether a feature that names the deploy group is on for
// the user. group is the id of a deploy_groups row.
//
// A deploy group gives a feature to some users before the rest. A feature's
// app_config row names the group in its property deploy_group, and the feature
// is on for a user when its own switch is on AND this function returns true:
//
//   - no group ("") returns true, so the switch alone decides;
//   - a group returns true only for a user with a row in
//     deploy_group_memberships.
//
// An id that matches no group has no members, so it returns false for every
// user. A mistake in the id can turn a feature off, but can't give it to all.
func InDeployGroup(app core.App, group string, uid string) (bool, error) {
	if group == "" {
		return true, nil
	}
	if uid == "" {
		return false, nil
	}
	_, err := app.FindFirstRecordByFilter(
		"deploy_group_memberships",
		"user={:user} && deploy_group={:group}",
		dbx.Params{"user": uid, "group": group},
	)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// DeployGroupAllows reads the property deploy_group of an app_config row and
// says whether that group lets the user in. See InDeployGroup. Use it with the
// feature's own switch:
//
//	on, _ := GetConfigBool(app, "some_feature", "enabled", false)
//	allowed, _ := DeployGroupAllows(app, "some_feature", uid)
//	if on && allowed { ... }
//
// With no row or no deploy_group, it returns true. A deploy_group that is not
// text returns false and an error.
func DeployGroupAllows(app core.App, domainKey string, uid string) (bool, error) {
	// GetConfigValue treats every error as "no row". Here that would let all
	// users in when the database can't be read, so this reads the row itself.
	record, err := app.FindFirstRecordByData("app_config", "key", domainKey)
	if errors.Is(err, sql.ErrNoRows) {
		return true, nil
	}
	if err != nil {
		return false, err
	}
	var config map[string]any
	if err := json.Unmarshal([]byte(record.GetString("value")), &config); err != nil {
		return false, err
	}
	raw, exists := config["deploy_group"]
	if !exists || raw == nil {
		return true, nil
	}
	group, valid := raw.(string)
	if !valid {
		return false, fmt.Errorf("%s.deploy_group must be the id of a deploy group", domainKey)
	}
	return InDeployGroup(app, group, uid)
}

// IsJobsEditingEnabled checks if job creation/editing/absorb is allowed.
// Reads from app_config where key="jobs", checks value.create_edit_absorb.
// Defaults to true (fail-open) if config is missing.
func IsJobsEditingEnabled(app core.App) (bool, error) {
	return GetConfigBool(app, "jobs", "create_edit_absorb", true)
}

// IsJobRequestApprovalRequired reads the project creation flag.
func IsJobRequestApprovalRequired(app core.App) (bool, error) {
	record, err := app.FindFirstRecordByData("app_config", "key", "jobs")
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	var config map[string]any
	if err := json.Unmarshal([]byte(record.GetString("value")), &config); err != nil {
		return false, err
	}
	raw, exists := config["require_job_request_approval"]
	if !exists {
		return false, nil
	}
	enabled, valid := raw.(bool)
	if !valid {
		return false, fmt.Errorf("jobs.require_job_request_approval must be Boolean")
	}
	return enabled, nil
}

// ErrJobsEditingDisabled is returned when job editing is disabled
var ErrJobsEditingDisabled = &errs.HookError{
	Status:  http.StatusForbidden,
	Message: "job editing is currently disabled",
	Data: map[string]errs.CodeError{
		"global": {Code: "jobs_editing_disabled", Message: "job creation and editing is disabled during transition"},
	},
}

// IsExpensesEditingEnabled checks if expense/vendor/PO creation/editing/absorb is allowed.
// Reads from app_config where key="expenses", checks value.create_edit_absorb.
// Defaults to true (fail-open) if config is missing.
func IsExpensesEditingEnabled(app core.App) (bool, error) {
	return GetConfigBool(app, "expenses", "create_edit_absorb", true)
}

// IsTimeEditingEnabled checks if time entry creation/editing and selected
// timesheet mutations are allowed.
// Reads from app_config where key="time", checks value.create_edit.
// Defaults to true (fail-open) if config is missing.
func IsTimeEditingEnabled(app core.App) (bool, error) {
	return GetConfigBool(app, "time", "create_edit", true)
}

// IsNotificationFeatureEnabled checks whether a notification feature/template is enabled.
// Reads from app_config where key="notifications", and uses templateCode as the JSON key.
// Defaults to false (fail-closed) when config is missing.
func IsNotificationFeatureEnabled(app core.App, templateCode string) (bool, error) {
	return GetConfigBool(app, "notifications", templateCode, false)
}

func GetPurchaseOrderSecondStageTimeoutHours(app core.App) float64 {
	const defaultTimeoutHours = 24.0

	config, err := GetConfigValue(app, "purchase_orders")
	if err != nil || config == nil {
		return defaultTimeoutHours
	}

	rawValue, ok := config["second_stage_timeout_hours"]
	if !ok {
		return defaultTimeoutHours
	}

	switch v := rawValue.(type) {
	case float64:
		if v > 0 {
			return v
		}
	case int:
		if v > 0 {
			return float64(v)
		}
	case int64:
		if v > 0 {
			return float64(v)
		}
	}

	return defaultTimeoutHours
}

// IsLegacyPOCreateUpdateEnabled checks whether the hidden legacy PO
// create/update flow is enabled via app_config.purchase_orders.
// Defaults to false (fail-closed).
func IsLegacyPOCreateUpdateEnabled(app core.App) (bool, error) {
	return GetConfigBool(app, "purchase_orders", "enable_legacy_po_create_update", false)
}

// POExpenseExcessConfig holds the configuration for how much expenses can
// exceed a purchase order total.
type POExpenseExcessConfig struct {
	Percent float64 // fractional, e.g. 0.05 = 5%
	Value   float64 // fixed amount in the home currency (CAD)
	Mode    string  // "lesser_of" or "greater_of"
}

// POExpenseCurrency is the currency of a PO as the excess limit needs it.
type POExpenseCurrency struct {
	Code string // ISO code of the PO's currency; "" means the home currency
	// HomeRate is home-currency units per unit of Code. It is not used for the
	// home currency.
	HomeRate float64
}

// POExpenseLimitResult holds the computed total limit and a human-readable
// description of the excess amount for use in error messages.
type POExpenseLimitResult struct {
	TotalLimitCents int64  // largest allowed total, in whole cents of the PO's currency
	ExcessText      string // e.g. "5.00%", "$100.00", or "100.00 CAD (74.07 USD)" on a foreign PO
	POTotalText     string // the PO total, e.g. "$1500.00", or "1,500.00 USD" on a foreign PO
}

// GetPOExpenseExcessConfig reads the po_expense_allowed_excess config from the
// "expenses" domain in app_config. Returns defaults matching the existing
// hardcoded constants when the config is missing or invalid.
func GetPOExpenseExcessConfig(app core.App) POExpenseExcessConfig {
	defaults := POExpenseExcessConfig{
		Percent: constants.MAX_PURCHASE_ORDER_EXCESS_PERCENT,
		Value:   constants.MAX_PURCHASE_ORDER_EXCESS_VALUE,
		Mode:    "lesser_of",
	}

	config, err := GetConfigValue(app, "expenses")
	if err != nil || config == nil {
		return defaults
	}

	rawExcess, ok := config["po_expense_allowed_excess"]
	if !ok {
		return defaults
	}

	excessMap, ok := rawExcess.(map[string]any)
	if !ok {
		return defaults
	}

	result := defaults // start from defaults, override only valid fields

	// "percent" is stored as a human-readable percentage (e.g. 5 means 5%)
	// and normalized to a fraction for internal use.
	if percent, err := CoerceFloat64(excessMap["percent"]); err == nil && percent >= 0 && percent <= 100 {
		result.Percent = percent / 100
	}

	if value, err := CoerceFloat64(excessMap["value"]); err == nil && value >= 0 {
		result.Value = value
	}

	if mode, ok := excessMap["mode"].(string); ok {
		if mode == "lesser_of" || mode == "greater_of" {
			result.Mode = mode
		}
	}

	return result
}

// GetPOExpenseCurrency returns the currency of a PO and the rate between it and
// the home currency, for CalculatePOExpenseTotalLimit. The rate is the one the
// PO's stored home total implies: approval_total_home / approval_total. Both
// are totals over all periods of a Recurring PO, so the rate holds for total
// too. A PO in the home currency is never converted.
func GetPOExpenseCurrency(app core.App, po *core.Record) POExpenseCurrency {
	info, err := ResolveCurrencyInfo(app, po.GetString("currency"))
	if err != nil || IsHomeCurrencyInfo(info) {
		// A currency that can't be found is treated as the home currency, as
		// EffectiveCurrencyCode does.
		return POExpenseCurrency{}
	}
	code := CurrencyCodeOrHome(info)
	approvalTotal := po.GetFloat("approval_total")
	approvalTotalHome := po.GetFloat("approval_total_home")
	if approvalTotal > 0 && approvalTotalHome > 0 {
		return POExpenseCurrency{Code: code, HomeRate: approvalTotalHome / approvalTotal}
	}
	// Fallback: an old or imported PO with no stored home total. Use the
	// currency's rate today; if it has none, use 1, so the fixed amount is
	// added as is, as it was before it was a home-currency amount.
	rate := CurrencyRateOrOne(info)
	if rate <= 0 {
		rate = 1
	}
	return POExpenseCurrency{Code: code, HomeRate: rate}
}

// CalculatePOExpenseTotalLimit computes the maximum allowed expense total for a
// purchase order, based on the excess configuration. poTotal is in the PO's
// currency and so is the limit. The fixed amount (cfg.Value) is in the home
// currency: on a foreign PO it is divided by currency.HomeRate before the mode
// compares it with the percent and before it is added, so that all of the
// arithmetic is in the PO's currency. On a home-currency PO it is used as is.
// The limit is floored to the cent (CurrencyLimitToCents). The ExcessText field
// describes which limit was applied (percent or fixed amount).
//
// Note: when poTotal is 0 the percent excess is also 0, so the percent branch
// is always selected and ExcessText will report a percentage even though the
// effective allowed excess is $0. This is acceptable because $0 POs are not
// expected in practice.
func CalculatePOExpenseTotalLimit(poTotal float64, currency POExpenseCurrency, cfg POExpenseExcessConfig) POExpenseLimitResult {
	isHome := IsHomeCurrencyInfo(CurrencyInfo{Code: currency.Code})
	value := cfg.Value
	if !isHome && currency.HomeRate > 0 {
		value = cfg.Value / currency.HomeRate
	}

	percentLimit := poTotal * (1.0 + cfg.Percent)
	valueLimit := poTotal + value

	var useValue bool
	switch cfg.Mode {
	case "greater_of":
		useValue = value >= poTotal*cfg.Percent
	default: // "lesser_of"
		useValue = value < poTotal*cfg.Percent
	}

	totalLimit := percentLimit
	excessText := fmt.Sprintf("%0.2f%%", cfg.Percent*100)
	if useValue {
		totalLimit = valueLimit
		if isHome {
			excessText = fmt.Sprintf("$%0.2f", cfg.Value)
		} else {
			// The amount in the PO's currency is floored to the cent, as the
			// limit is.
			excessText = fmt.Sprintf("%s (%s)",
				FormatCurrencyAmount(cfg.Value, HomeCurrencyCode),
				FormatCurrencyAmount(float64(CurrencyLimitToCents(value))/100, currency.Code))
		}
	}

	poTotalText := fmt.Sprintf("$%0.2f", poTotal)
	if !isHome {
		poTotalText = FormatCurrencyAmount(poTotal, currency.Code)
	}

	return POExpenseLimitResult{
		TotalLimitCents: CurrencyLimitToCents(totalLimit),
		ExcessText:      excessText,
		POTotalText:     poTotalText,
	}
}

// CoerceFloat64 converts a JSON numeric value (which may be float64, int, or
// int64 after json.Unmarshal) to float64.
func CoerceFloat64(v any) (float64, error) {
	if v == nil {
		return 0, fmt.Errorf("nil value")
	}
	switch n := v.(type) {
	case float64:
		return n, nil
	case int:
		return float64(n), nil
	case int64:
		return float64(n), nil
	default:
		return 0, fmt.Errorf("unsupported type %T", v)
	}
}

// GetNoPOExpenseLimit reads the no-PO expense limit from the "expenses" domain
// in app_config. Returns the default from constants when the config is missing,
// the key is absent, or the value is negative.
// A value of 0 means all non-exempt expenses require a PO.
func GetNoPOExpenseLimit(app core.App) float64 {
	config, err := GetConfigValue(app, "expenses")
	if err != nil || config == nil {
		return constants.NO_PO_EXPENSE_LIMIT
	}
	if limit, err := CoerceFloat64(config["no_po_expense_limit"]); err == nil && limit >= 0 {
		return limit
	}
	return constants.NO_PO_EXPENSE_LIMIT
}

// ErrExpensesEditingDisabled is returned when expense editing is disabled
var ErrExpensesEditingDisabled = &errs.HookError{
	Status:  http.StatusForbidden,
	Message: "expense editing is currently disabled",
	Data: map[string]errs.CodeError{
		"global": {Code: "expenses_editing_disabled", Message: "expense, purchase order, and vendor creation and editing is disabled during transition"},
	},
}

// ErrTimeEditingDisabled is returned when time editing is disabled
var ErrTimeEditingDisabled = &errs.HookError{
	Status:  http.StatusForbidden,
	Message: "time editing is currently disabled",
	Data: map[string]errs.CodeError{
		"global": {Code: "time_editing_disabled", Message: "time entry and time amendment creation/editing, plus timesheet submission/approval, are disabled during transition"},
	},
}
