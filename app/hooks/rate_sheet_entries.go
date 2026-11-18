package hooks

import (
	"math"
	"net/http"

	"tybalt/errs"

	"github.com/pocketbase/pocketbase/core"
)

// OvertimeRateDecimalsMessage is the refusal of an overtime rate with more
// than two decimal places, on every path that writes rate_sheet_entries.
const OvertimeRateDecimalsMessage = "overtime_rate must have at most two decimal places"

// OvertimeRateDecimalsCode is the error code of that refusal.
const OvertimeRateDecimalsCode = "validation_too_many_decimals"

// HasAtMostTwoDecimals reports whether value is a whole number of cents. A
// float is not exact (1.15 * 100 is 114.99999999999999), so a value within a
// millionth of a cent of a whole cent counts as one; 1.255 does not.
func HasAtMostTwoDecimals(value float64) bool {
	cents := value * 100
	return math.Abs(cents-math.Round(cents)) < 1e-6
}

// addRateSheetEntryHooks refuses an overtime rate with more than two decimal
// places in a create through the collection API. The create route, the entry
// correction route and proposed rates check it themselves. The collection has
// no update rule.
func addRateSheetEntryHooks(app core.App) {
	app.OnRecordCreateRequest("rate_sheet_entries").BindFunc(func(e *core.RecordRequestEvent) error {
		if !HasAtMostTwoDecimals(e.Record.GetFloat("overtime_rate")) {
			return AnnotateHookError(app, e, &errs.HookError{
				Status:  http.StatusBadRequest,
				Message: OvertimeRateDecimalsMessage,
				Data: map[string]errs.CodeError{
					"overtime_rate": {Code: OvertimeRateDecimalsCode, Message: OvertimeRateDecimalsMessage},
				},
			})
		}
		return e.Next()
	})
}
