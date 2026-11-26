<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/constants/constants.go","sha256":"ba15e682b407a81fef75ae7f9f5c9c6b418422699a91f84be266dc7371aadfac","model":"claude-sonnet-5-5"} -->

# Rules in app/constants/constants.go

## Location Plus Code format

Lines 18-21

Asked as: plus code format; location plus code invalid; plus code characters allowed; open location code length

A location Plus Code must be 8 characters, then '+', then 2 or 3 characters (LocationPlusCodeRegex). Only the characters 23456789CFGHJMPQRVWX are allowed.

## Default work week hours and charge-out rate

Lines 24-26

Asked as: default work week hours; new admin profile defaults; default charge out rate; admin default hourly rate

An admin profile starts with a work week of 40 hours (DEFAULT_WORK_WEEK_HOURS). Its default charge-out rate is 50 (DEFAULT_CHARGE_OUT_RATE).

## Recurring purchase order maximum length

Lines 32-34

Asked as: recurring purchase order max days; recurring PO too long; recurring purchase order date range limit; how long can recurring PO run

A recurring purchase order can have at most 400 days between its start date and its end date (RECURRING_MAX_DAYS).

## Allowed excess on purchase order expenses

Lines 36-40

Asked as: expense over purchase order amount; PO expense excess allowed; expense exceeds PO percentage; override PO excess setting

By default an expense can exceed its purchase order by 5% (MAX_PURCHASE_ORDER_EXCESS_PERCENT = 0.05) or by 100.00 (MAX_PURCHASE_ORDER_EXCESS_VALUE). An admin can override these in app_config, domain "expenses", property "po_expense_allowed_excess". This file does not show how Turbo combines the two values.

## Limit on expenses without a purchase order

Lines 42-45

Asked as: expenses without purchase order limit; non PO expense cap; expense without PO too much; no PO spending limit

Turbo limits the total of expenses that have no purchase order (LIMIT_NON_PO_AMOUNTS is true). The limit is 100.00 (NO_PO_EXPENSE_LIMIT).

## Second approver purchase order ceiling

Lines 47-49

Asked as: second approver PO limit; second approver approval total max; PO over million approval; can't approve large purchase order

A second approver can approve a purchase order only when its approval total is at most 1,000,000 (MAX_APPROVAL_TOTAL).

## Legacy PO create/update claim

Lines 51-55

Asked as: legacy PO create access; hidden purchase order editing; old PO form permission; po approver claim id

The claim named legacy_po_create_update gives access to the hidden legacy purchase order create and update flow (LEGACY_PO_CREATE_UPDATE_CLAIM_NAME). The purchase order approver claim (po_approver) has a fixed claim ID.

## Pay date follows pay period ending

Lines 58-64

Asked as: when is pay day; pay date after period ending; payday calculation; pay date not Friday

A pay day is 13 days after a Saturday pay period ending, and it falls on a Friday (PAY_DATE_DAYS_AFTER_PERIOD_ENDING). Payroll counts from the epoch date of 1 March 2025 (PAYROLL_EPOCH).

## Uncommit request window for expenses

Lines 65-70

Asked as: uncommit deadline expense; how long to request uncommit; expense can't be uncommitted anymore; uncommit window payroll expense

A payroll-paid expense can be requested for uncommit until 7 days after its pay period ending (UNCOMMIT_PAYROLL_LOCK_DAYS). Any other expense can be requested for uncommit until 90 days after it was committed (UNCOMMIT_NON_PAYROLL_LOCK_DAYS).

## Computer login retention

Lines 73-75

Asked as: how long login history kept; computer login deleted; raw login retention; login records auto delete

Turbo keeps a computer login or raw login for 30 days (COMPUTER_LOGIN_RETENTION_DAYS). A daily job deletes older rows.

## Absorb proposal retention

Lines 77-80

Asked as: absorb proposal deleted after; declined absorb proposal gone; pending absorb proposal never deleted; how long absorb proposals kept

Turbo keeps a decided absorb proposal (accepted, declined, withdrawn or obsolete) for 30 days (ABSORB_PROPOSAL_RETENTION_DAYS). A daily job then deletes it. Turbo never deletes a pending proposal.

## Absorb history retention

Lines 82-84

Asked as: absorb history deleted; committed absorb history kept how long; absorb history row removed; absorb record retention period

Turbo keeps the history row of a committed absorb for 90 days (ABSORB_HISTORY_RETENTION_DAYS). A daily job then deletes older rows.
