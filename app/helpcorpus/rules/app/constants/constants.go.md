<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/constants/constants.go","sha256":"7b69fd2afd7b410f88abcc2d68b5e0a1776534c0693b5a13cc7507f37fb7b053","model":"claude-sonnet-5-5"} -->

# Rules in app/constants/constants.go

## Plus Code location format

Lines 18-21

Asked as: plus code format; location code invalid; open location code wrong; plus code characters allowed

A location Plus Code (LocationPlusCodeRegex) must be 8 characters, then '+', then 2 or 3 characters. Only the characters 23456789CFGHJMPQRVWX are valid.

## Default work week hours on admin profile

Lines 24-25

Asked as: default work week hours; standard weekly hours admin; 40 hours default week; admin profile work week

The default work week on an admin profile is 40 hours (DEFAULT_WORK_WEEK_HOURS).

## Default charge out rate on admin profile

Lines 24-26

Asked as: default charge out rate; admin profile hourly rate; standard charge rate; default rate 50

The default charge out rate on an admin profile is 50 (DEFAULT_CHARGE_OUT_RATE).

## Default branch ID constant

Lines 28-30

Asked as: default branch id; ThunderBay branch code; which branch is default; branch id test database

The default branch ID (DEFAULT_BRANCH_ID) is 80875lm27v8wgi4. The code comment says it is the ThunderBay branch in the test database.

## Recurring purchase order maximum length

Lines 32-34

Asked as: recurring po too long; recurring purchase order 400 days; max length recurring order; recurring po date range limit

A recurring purchase order can span at most 400 days between its start and end dates (RECURRING_MAX_DAYS).

## Allowed excess on purchase order expenses

Lines 36-40

Asked as: expense over po amount; expense exceeds purchase order; how much over po allowed; po excess tolerance percent

By default an expense on a purchase order may exceed the PO by 5% (MAX_PURCHASE_ORDER_EXCESS_PERCENT = 0.05) or by 100.0 (MAX_PURCHASE_ORDER_EXCESS_VALUE). An admin can override these in app_config, domain "expenses", property "po_expense_allowed_excess".

## Limit on expenses without a purchase order

Lines 42-45

Asked as: expense without po limit; non po expenses cap; expense no purchase order max; 100 limit no po

Turbo limits the total of expenses that have no purchase order (LIMIT_NON_PO_AMOUNTS is true). The limit is 100.0 (NO_PO_EXPENSE_LIMIT).

## Maximum total for second approver

Lines 47-49

Asked as: second approver limit; approve large po; po total approval cap; second approver million

A second approver can approve a purchase order whose approval_total is at most 1,000,000.0 (MAX_APPROVAL_TOTAL).

## PO approver claim ID

Lines 51-52

Asked as: po approver claim; which claim approves po; po approver permission id; claim id for po approving

The po_approver claim has the ID 5vh881k048bboim (PO_APPROVER_CLAIM_ID).

## Legacy PO create/update claim

Lines 54-55

Asked as: hidden legacy po screen; old po create form; legacy po permission; access old purchase order editor

The claim named legacy_po_create_update gives access to the hidden legacy purchase order create and update flow (LEGACY_PO_CREATE_UPDATE_CLAIM_NAME).

## Payroll epoch date

Lines 58-59

Asked as: payroll start date; when payroll began; payroll epoch date; first payroll date 2025

The payroll epoch (PAYROLL_EPOCH) is 1 March 2025 UTC.

## Pay day after pay period ending

Lines 62-64

Asked as: when is payday; pay date after period end; days until pay day; friday pay date calculation

The pay day is 13 days after a Saturday pay_period_ending, which is a Friday (PAY_DATE_DAYS_AFTER_PERIOD_ENDING).

## Uncommit window for payroll-paid expenses

Lines 65-67

Asked as: uncommit payroll expense deadline; cannot uncommit paid expense; how long to uncommit payroll; uncommit lock after pay period

A payroll-paid expense can be requested for uncommit until 7 days after its pay_period_ending (UNCOMMIT_PAYROLL_LOCK_DAYS).

## Uncommit window for other expenses

Lines 68-70

Asked as: uncommit expense time limit; how long can I uncommit; 90 day uncommit window; expense uncommit deadline

Any other expense can be requested for uncommit until 90 days after its commit (UNCOMMIT_NON_PAYROLL_LOCK_DAYS).

## Computer login retention

Lines 73-75

Asked as: how long are computer logins kept; login history deleted; computer login log retention; old login records removed

Turbo keeps a row of computer_logins or computer_raw_logins for 30 days (COMPUTER_LOGIN_RETENTION_DAYS). A daily cron job deletes older rows.
