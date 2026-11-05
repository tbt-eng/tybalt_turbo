package utilities

import (
	"time"
	"tybalt/constants"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

const (
	ExpenseUncommitBlockPayrollClosed = "payroll_closed"
	ExpenseUncommitBlockTooOld        = "too_old"
)

// ExpenseUncommitBlock says why a committed expense can no longer be requested
// for uncommit. The dates let the UI explain it.
type ExpenseUncommitBlock struct {
	Code            string `json:"code"`
	Message         string `json:"message"`
	PayPeriodEnding string `json:"pay_period_ending,omitempty"`
	PayDate         string `json:"pay_date,omitempty"`
	LimitDays       int    `json:"limit_days,omitempty"`
}

// IsPayrollExpensePaymentType reports whether an expense of this payment type
// is repaid to the employee on payroll. Keep this in step with
// routes/expense_reimbursements.sql.
func IsPayrollExpensePaymentType(paymentType string) bool {
	switch paymentType {
	case "Expense", "PersonalReimbursement", "Allowance", "Mileage", "Meals":
		return true
	}
	return false
}

func expenseUncommitLockDays(app core.App, property string, defaultValue int) int {
	config, err := GetConfigValue(app, "expenses")
	if err != nil || config == nil {
		return defaultValue
	}
	if days, err := CoerceFloat64(config[property]); err == nil && days >= 0 {
		return int(days)
	}
	return defaultValue
}

// FindExpenseUncommitBlock returns why an expense can no longer be requested
// for uncommit, or nil when it still can.
//
// A payroll type is blocked once payroll may have run. Commits still land in a
// pay period until the Saturday one week after it ends (see
// GenerateCommittedPayPeriodEnding), and pay day is the Friday after that, so
// payroll can't be final before pay_period_ending + 7 days. Other types are
// blocked a fixed number of days after commit.
func FindExpenseUncommitBlock(app core.App, paymentType string, payPeriodEnding string, committed types.DateTime, now time.Time) *ExpenseUncommitBlock {
	today := now.Format(time.DateOnly)

	if IsPayrollExpensePaymentType(paymentType) {
		block := &ExpenseUncommitBlock{
			Code:            ExpenseUncommitBlockPayrollClosed,
			Message:         "payroll for this expense has closed",
			PayPeriodEnding: payPeriodEnding,
		}
		// Without a pay period there is no way to tell whether it was paid.
		ending, err := time.Parse(time.DateOnly, payPeriodEnding)
		if err != nil {
			return block
		}
		block.PayDate = ending.AddDate(0, 0, constants.PAY_DATE_DAYS_AFTER_PERIOD_ENDING).Format(time.DateOnly)
		lockDays := expenseUncommitLockDays(app, "uncommit_payroll_lock_days", constants.UNCOMMIT_PAYROLL_LOCK_DAYS)
		if today > ending.AddDate(0, 0, lockDays).Format(time.DateOnly) {
			return block
		}
		return nil
	}

	lockDays := expenseUncommitLockDays(app, "uncommit_non_payroll_lock_days", constants.UNCOMMIT_NON_PAYROLL_LOCK_DAYS)
	if today > committed.Time().In(now.Location()).AddDate(0, 0, lockDays).Format(time.DateOnly) {
		return &ExpenseUncommitBlock{
			Code:      ExpenseUncommitBlockTooOld,
			Message:   "this expense was committed too long ago",
			LimitDays: lockDays,
		}
	}
	return nil
}
