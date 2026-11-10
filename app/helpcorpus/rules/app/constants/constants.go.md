<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/constants/constants.go","sha256":"daa1d728ef62354507e244f008825fb5db11cf4c89d897a28db492f503aa92cf","model":"claude-sonnet-5-5"} -->

# Rules in app/constants/constants.go

## Plus Code location format

Lines 18-21

A location Plus Code must be 8 characters, then '+', then 2 or 3 characters (LocationPlusCodeRegex). The allowed characters are 23456789CFGHJMPQRVWX.

## Default work week hours

Lines 25-25

The default work week for an admin profile is 40 hours (DEFAULT_WORK_WEEK_HOURS).

## Default charge-out rate

Lines 26-26

The default charge-out rate for an admin profile is 50 (DEFAULT_CHARGE_OUT_RATE).

## Default branch ID constant

Lines 28-30

The default branch ID constant (DEFAULT_BRANCH_ID) is "80875lm27v8wgi4". The code comment says it is the 'ThunderBay' branch from the test database.

## Recurring purchase order maximum span

Lines 32-34

A recurring purchase order can span at most 400 days between its start date and end date (RECURRING_MAX_DAYS).

## Default allowed excess on purchase order expenses

Lines 36-40

By default an expense on a purchase order can exceed it by 5% (MAX_PURCHASE_ORDER_EXCESS_PERCENT = 0.05) or by 100.0 (MAX_PURCHASE_ORDER_EXCESS_VALUE). An app_config setting "po_expense_allowed_excess" in the "expenses" domain can override these defaults. This file does not show how the two limits combine.

## Limit on expenses without a purchase order

Lines 42-45

The limit on expenses that have no purchase order is switched on (LIMIT_NON_PO_AMOUNTS = true). The limit is 100.0 (NO_PO_EXPENSE_LIMIT).

## Maximum total for second approver

Lines 47-49

A purchase order with an approval total (approval_total) above 1,000,000.0 cannot be approved by a second approver (MAX_APPROVAL_TOTAL).

## PO approver claim ID

Lines 51-52

The po_approver claim has the ID "5vh881k048bboim" (PO_APPROVER_CLAIM_ID).

## Legacy PO create/update claim

Lines 54-55

The claim named "legacy_po_create_update" gives access to the hidden legacy purchase order create and update flow (LEGACY_PO_CREATE_UPDATE_CLAIM_NAME).

## Payroll epoch date

Lines 58-59

The payroll epoch date is 2025-03-01 UTC (PAYROLL_EPOCH).

## Pay day after pay period ending

Lines 62-64

A pay day falls 13 days after the pay period ending date, which is a Saturday (PAY_DATE_DAYS_AFTER_PERIOD_ENDING). The pay day is a Friday.

## Uncommit window for payroll-paid expenses

Lines 65-67

A payroll-paid expense can be requested for uncommit until 7 days after its pay period ending date (UNCOMMIT_PAYROLL_LOCK_DAYS).

## Uncommit window for other expenses

Lines 68-70

Any other expense can be requested for uncommit until 90 days after it was committed (UNCOMMIT_NON_PAYROLL_LOCK_DAYS).
