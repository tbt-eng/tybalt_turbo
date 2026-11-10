<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/utilities/expense_uncommit.go","sha256":"8c1e7f6da4e3d7c2cefb4adf8cac8aa2641bc06f09bc0701ba5b84b0e80276ea","model":"claude-sonnet-5-5"} -->

# Rules in app/utilities/expense_uncommit.go

## Payment types repaid through payroll

Lines 26-35

Turbo treats an expense as a payroll expense when its payment type (payment_type) is Expense, PersonalReimbursement, Allowance, Mileage or Meals. Every other payment type is a non-payroll expense. The two kinds follow different uncommit deadlines.

## Payroll uncommit lock days setting and default

Lines 37-46

The payroll lock is the number of days after pay period ending set by 'uncommit_payroll_lock_days' in the expenses config. Turbo uses the default UNCOMMIT_PAYROLL_LOCK_DAYS when the config is missing, unreadable, not a number or negative. A value of 0 is accepted.

## Non-payroll uncommit lock days setting and default

Lines 37-46

The non-payroll lock is the number of days after commit set by 'uncommit_non_payroll_lock_days' in the expenses config. Turbo uses the default UNCOMMIT_NON_PAYROLL_LOCK_DAYS when the config is missing, unreadable, not a number or negative. A value of 0 is accepted.

## Uncommit allowed when no block applies

Lines 56-87

Turbo returns no block, so an uncommit request is still possible, when the expense is inside its payroll lock or inside its non-payroll commit age limit.

## Payroll expense uncommit request blocked after payroll lock

Lines 59-76

A committed payroll expense can no longer be requested for uncommit once today is later than its pay period ending plus the lock days. Turbo then gives the reason 'payroll for this expense has closed' (payroll_closed). Today equal to the last allowed day is still allowed.

## Payroll block shows pay date

Lines 60-70

For a payroll expense with a valid pay period ending, the block includes the pay date. Turbo sets the pay date to the pay period ending plus PAY_DATE_DAYS_AFTER_PERIOD_ENDING days. The block also shows the pay period ending.

## Payroll expense without a pay period is always blocked

Lines 65-69

When a payroll expense has no valid pay period ending date, Turbo blocks the uncommit request with 'payroll for this expense has closed'. Without a pay period Turbo cannot tell whether the expense was paid. In this case the block has no pay date.

## Non-payroll expense uncommit blocked after commit age limit

Lines 78-86

A committed non-payroll expense can no longer be requested for uncommit once today is later than its commit date plus the lock days. Turbo then gives the reason 'this expense was committed too long ago' (too_old) and the limit in days (limit_days). The commit date is read in the current time zone, and the last allowed day itself is still allowed.
