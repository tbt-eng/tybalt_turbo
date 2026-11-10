<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/constants/constants.go","sha256":"daa1d728ef62354507e244f008825fb5db11cf4c89d897a28db492f503aa92cf","model":"claude-sonnet-5-5"} -->

# Rules in app/constants/constants.go

## Plus Code location format

Lines 18-21

Asked as: plus code format; location plus code wrong; open location code length; plus code characters allowed

A location Plus Code must be 8 characters, then '+', then 2 or 3 characters (LocationPlusCodeRegex). The allowed characters are 23456789CFGHJMPQRVWX.

## Default work week hours

Lines 25-25

Asked as: default work week hours; standard hours per week; weekly hours admin profile; how many hours a week

The default work week for an admin profile is 40 hours (DEFAULT_WORK_WEEK_HOURS).

## Default charge-out rate

Lines 26-26

Asked as: default charge out rate; default hourly rate admin; standard charge rate; admin profile rate

The default charge-out rate for an admin profile is 50 (DEFAULT_CHARGE_OUT_RATE).

## Default branch ID constant

Lines 28-30

Asked as: default branch id; which branch is default; thunderbay branch id; default branch constant

The default branch ID constant (DEFAULT_BRANCH_ID) is "80875lm27v8wgi4". The code comment says it is the 'ThunderBay' branch from the test database.

## Recurring purchase order maximum span

Lines 32-34

Asked as: recurring po too long; recurring purchase order max days; how long can recurring po run; recurring po date range limit

A recurring purchase order can span at most 400 days between its start date and end date (RECURRING_MAX_DAYS).

## Default allowed excess on purchase order expenses

Lines 36-40

Asked as: po expense over budget allowed; expense exceeds purchase order amount; how much over po allowed; po excess percentage; purchase order overspend limit

By default an expense on a purchase order can exceed it by 5% (MAX_PURCHASE_ORDER_EXCESS_PERCENT = 0.05) or by 100.0 (MAX_PURCHASE_ORDER_EXCESS_VALUE). An app_config setting "po_expense_allowed_excess" in the "expenses" domain can override these defaults. This file does not show how the two limits combine.

## Limit on expenses without a purchase order

Lines 42-45

Asked as: expense without purchase order limit; non po expense cap; expense over 100 no po; no po amount limit

The limit on expenses that have no purchase order is switched on (LIMIT_NON_PO_AMOUNTS = true). The limit is 100.0 (NO_PO_EXPENSE_LIMIT).

## Maximum total for second approver

Lines 47-49

Asked as: second approver limit; po approval total too high; cannot approve large purchase order; approval total maximum

A purchase order with an approval total (approval_total) above 1,000,000.0 cannot be approved by a second approver (MAX_APPROVAL_TOTAL).

## PO approver claim ID

Lines 51-52

Asked as: po approver claim id; which claim approves purchase orders; po approver permission code; purchase order approver claim

The po_approver claim has the ID "5vh881k048bboim" (PO_APPROVER_CLAIM_ID).

## Legacy PO create/update claim

Lines 54-55

Asked as: legacy po form; old purchase order screen; hidden po create page; legacy purchase order access

The claim named "legacy_po_create_update" gives access to the hidden legacy purchase order create and update flow (LEGACY_PO_CREATE_UPDATE_CLAIM_NAME).

## Payroll epoch date

Lines 58-59

Asked as: payroll start date; payroll epoch; when payroll began; first payroll date

The payroll epoch date is 2025-03-01 UTC (PAYROLL_EPOCH).

## Pay day after pay period ending

Lines 62-64

Asked as: when do we get paid; pay day after period end; payday date calculation; payment date after pay period

A pay day falls 13 days after the pay period ending date, which is a Saturday (PAY_DATE_DAYS_AFTER_PERIOD_ENDING). The pay day is a Friday.

## Uncommit window for payroll-paid expenses

Lines 65-67

Asked as: uncommit payroll expense deadline; cannot uncommit paid expense; how long to uncommit payroll; payroll lock uncommit days

A payroll-paid expense can be requested for uncommit until 7 days after its pay period ending date (UNCOMMIT_PAYROLL_LOCK_DAYS).

## Uncommit window for other expenses

Lines 68-70

Asked as: uncommit expense time limit; how long to uncommit expense; expense commit lock period; cannot uncommit old expense

Any other expense can be requested for uncommit until 90 days after it was committed (UNCOMMIT_NON_PAYROLL_LOCK_DAYS).
