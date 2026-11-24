<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/constants/constants.go","sha256":"faf00b3fae171ee1660e3519bdedbada410363fd9010994d24d9d2f0a43da4b8","model":"claude-sonnet-5-5"} -->

# Rules in app/constants/constants.go

## Plus Code location format

Lines 18-21

Asked as: plus code format; location plus code invalid; open location code characters; plus code length wrong

A location Plus Code (LocationPlusCodeRegex) must be 8 characters, then a '+', then 2 or 3 characters. Only the characters 23456789CFGHJMPQRVWX are valid.

## Default work week hours and charge-out rate

Lines 24-26

Asked as: default work week hours; admin default charge out rate; standard weekly hours admin; default hourly rate admin

An admin profile defaults to a work week of 40 hours (DEFAULT_WORK_WEEK_HOURS). Its default charge-out rate is 50 (DEFAULT_CHARGE_OUT_RATE).

## Default branch ID constant

Lines 28-30

Asked as: default branch id; thunderbay branch id; branch id constant; which branch is default

The constant DEFAULT_BRANCH_ID is "80875lm27v8wgi4". The code comment says it is the ThunderBay branch in the test database.

## Recurring purchase order maximum span

Lines 32-34

Asked as: recurring purchase order max days; recurring po too long; recurring po span limit; how long can recurring po run

A recurring purchase order can have at most 400 days between its start date and its end date (RECURRING_MAX_DAYS).

## Default allowed excess on PO expenses

Lines 36-40

Asked as: po expense over amount; expense exceeds purchase order; allowed excess on po; change po excess percent

By default, an expense on a purchase order may exceed it by 5% (MAX_PURCHASE_ORDER_EXCESS_PERCENT = 0.05) or by 100.0 (MAX_PURCHASE_ORDER_EXCESS_VALUE). An admin can override these values in app_config, in the "expenses" domain, under "po_expense_allowed_excess".

## Limit on expenses without a purchase order

Lines 42-45

Asked as: expense without purchase order limit; non po expense cap; expense no po 100 limit; why non po expense rejected

The limit on the total of expenses that have no purchase order is switched on (LIMIT_NON_PO_AMOUNTS = true). The limit is 100.0 (NO_PO_EXPENSE_LIMIT).

## Maximum total for a second approver

Lines 47-49

Asked as: second approver limit; approval total too high; co approver max amount; po over million approval

A second approver can approve a purchase order only if its approval_total is at most 1,000,000.0 (MAX_APPROVAL_TOTAL).

## PO approver claim ID

Lines 51-52

Asked as: po approver claim id; po_approver claim code; which id is po approver; po approver permission id

The ID of the po_approver claim is "5vh881k048bboim" (PO_APPROVER_CLAIM_ID).

## Legacy PO create/update claim

Lines 54-55

Asked as: legacy po form; hidden po create screen; old purchase order editor access; legacy purchase order claim

The claim named "legacy_po_create_update" gives access to the hidden legacy flow that creates and updates purchase orders.

## Payroll epoch date

Lines 58-59

Asked as: payroll start date; when payroll began; payroll epoch date; first payroll date 2025

The payroll epoch (PAYROLL_EPOCH) is 1 March 2025 UTC.

## Pay day after pay period ending

Lines 62-64

Asked as: when do I get paid; pay date after period; payday after pay period end; pay day calculation

The pay day is 13 days after a Saturday pay_period_ending, which is a Friday (PAY_DATE_DAYS_AFTER_PERIOD_ENDING).

## Uncommit request window for payroll-paid expenses

Lines 65-67

Asked as: uncommit payroll expense deadline; can't uncommit paid expense; how late to uncommit payroll; uncommit window payroll

A person can request an uncommit of a payroll-paid expense up to 7 days after its pay_period_ending (UNCOMMIT_PAYROLL_LOCK_DAYS).

## Uncommit request window for other expenses

Lines 68-70

Asked as: uncommit expense time limit; how long to uncommit expense; old expense uncommit request; uncommit deadline 90 days

A person can request an uncommit of any other expense up to 90 days after its commit (UNCOMMIT_NON_PAYROLL_LOCK_DAYS).

## Computer login retention

Lines 73-75

Asked as: how long login history kept; computer login deleted; login records retention; old computer logins removed

Turbo keeps a row of computer_logins or computer_raw_logins for 30 days (COMPUTER_LOGIN_RETENTION_DAYS). A daily cron job deletes older rows.

## Absorb proposal retention after decision

Lines 77-80

Asked as: absorb proposal deleted; how long absorb proposals kept; old absorb proposal removed; proposal auto delete after decision

Turbo keeps an absorb proposal for 30 days after its decision: accepted, declined, withdrawn or obsolete (ABSORB_PROPOSAL_RETENTION_DAYS). A daily cron job deletes older ones. A pending proposal is never deleted.
