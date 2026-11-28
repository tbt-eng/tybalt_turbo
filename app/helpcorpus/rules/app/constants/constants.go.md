<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/constants/constants.go","sha256":"6f2cf66e4e1ac467230359510db716d77d2f4a75c40db4ed79330ae5dc6c720d","model":"claude-sonnet-5-5"} -->

# Rules in app/constants/constants.go

## Location Plus Code format

Lines 18-21

Asked as: plus code format; location plus code invalid; open location code characters; plus code rejected

A location Plus Code (LocationPlusCodeRegex) has 8 characters, then '+', then 2 or 3 characters. Only the characters 23456789CFGHJMPQRVWX are valid.

## Default work week hours and charge out rate

Lines 24-26

Asked as: admin default hours; default work week; default charge out rate; new admin profile defaults

An admin profile defaults to a work week of 40 hours (DEFAULT_WORK_WEEK_HOURS). Its default charge out rate is 50 (DEFAULT_CHARGE_OUT_RATE).

## Default branch ID constant

Lines 28-30

Asked as: default branch id; thunderbay branch id; test database branch; branch id constant

The constant DEFAULT_BRANCH_ID is "80875lm27v8wgi4". It is the ID of the 'ThunderBay' branch in the test database.

## Recurring purchase order maximum length

Lines 32-34

Asked as: recurring po too long; recurring purchase order max days; 400 day limit; recurring po date range

A recurring purchase order can span at most 400 days between its start date and end date (RECURRING_MAX_DAYS).

## Allowed excess on purchase order expenses

Lines 36-40

Asked as: po expense over amount; expense exceeds purchase order; allowed excess percent; purchase order overspend tolerance

By default an expense on a purchase order may exceed it by 5% (MAX_PURCHASE_ORDER_EXCESS_PERCENT = 0.05) or by 100.0 (MAX_PURCHASE_ORDER_EXCESS_VALUE). An admin can override these in the app_config "expenses" domain, property "po_expense_allowed_excess".

## Limit on expenses without a purchase order

Lines 42-45

Asked as: expense without po limit; non po expense cap; expense over 100 no po; no purchase order limit

The limit on expenses that have no purchase order is switched on (LIMIT_NON_PO_AMOUNTS = true). The limit is 100.0 (NO_PO_EXPENSE_LIMIT).

## Maximum total for second approver

Lines 47-49

Asked as: second approver limit; approval total too high; second approver million; po approval cap

A second approver can approve a purchase order only when its approval_total is at most 1,000,000.0 (MAX_APPROVAL_TOTAL).

## PO approver claim ID

Lines 51-52

Asked as: po approver claim id; po approver claim; approver claim number; claim for po approval

The po_approver claim has the ID "5vh881k048bboim" (PO_APPROVER_CLAIM_ID).

## Legacy PO create/update claim

Lines 54-55

Asked as: legacy po access; legacy po screen; hidden po create page; old po update claim

The claim named "legacy_po_create_update" gives access to the hidden legacy PO create/update flow (LEGACY_PO_CREATE_UPDATE_CLAIM_NAME).

## Payroll epoch date

Lines 58-59

Asked as: payroll start date; payroll epoch date; when payroll began; first payroll date

The payroll epoch (PAYROLL_EPOCH) is 1 March 2025 UTC.

## Pay day after pay period ending

Lines 62-64

Asked as: pay day friday; when is payday; pay date after period; friday after saturday

The pay day is 13 days after a Saturday pay_period_ending, which gives a Friday (PAY_DATE_DAYS_AFTER_PERIOD_ENDING).

## Uncommit window for payroll-paid expenses

Lines 65-67

Asked as: uncommit payroll expense; payroll expense uncommit deadline; locked payroll expense; uncommit after pay period

A payroll-paid expense can be requested for uncommit until 7 days after its pay_period_ending (UNCOMMIT_PAYROLL_LOCK_DAYS).

## Uncommit window for other expenses

Lines 68-70

Asked as: uncommit expense deadline; how long to uncommit; expense uncommit 90 days; uncommit time limit

Any other expense can be requested for uncommit until 90 days after its commit (UNCOMMIT_NON_PAYROLL_LOCK_DAYS).

## Computer login retention

Lines 73-77

Asked as: old computer logins deleted; login history cleanup; raw login expiry; computer login removed

Turbo keeps a computer login row, or a raw login without a computer name, for 30 days. A daily cron job deletes older rows. A raw login with a name never expires, and Turbo keeps only the newest one per name.

## Absorb proposal retention

Lines 79-82

Asked as: absorb proposal deleted; old absorb proposals gone; proposal cleanup after decision; declined proposal removed

Turbo keeps an absorb proposal for 30 days after its decision (accepted, declined, withdrawn or obsolete). A daily cron job then deletes it. Turbo never deletes a pending proposal.

## Absorb history retention

Lines 84-86

Asked as: absorb history deleted; absorb history cleanup; old absorb records removed; absorb log retention

Turbo keeps the history row of a committed absorb (absorb_history) for 90 days. A daily cron job deletes older rows.
