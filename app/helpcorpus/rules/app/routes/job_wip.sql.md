<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_wip.sql","sha256":"ae096a7c292a311053190734097e696bb98b2f76102d913f6d2a7fd23218f086","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_wip.sql

## Job WIP time value and hour totals by job

Lines 3-7

The job WIP report sums each job's priced time entries into time value, hours, estimated hours (estimated_hours) and unpriced hours (unpriced_hours). A job with no time entries shows 0 for each. The time value is rounded to 2 decimals.

## Hours with a rate sheet but no role are counted separately

Lines 8-14

The WIP report counts hours as missing role hours (missing_role_hours) when the entry's rate sheet is set, the sheet gives no rate (sheet_rate is null) and the entry has no role. Hours with a role but no rate on the sheet are counted as missing role rate hours (missing_role_rate_hours). The counts include hours estimated at the employee rate and hours that stay unpriced.

## Which expenses count toward job WIP

Lines 15-30

The WIP report counts only expenses that are committed and not rejected. The expense date must be on or before the end date; an expense dated on the end date counts. An expense counts if it belongs to a selected job, or to an Active purchase order of a selected job.

## How an expense's CAD value is set for WIP

Lines 17-23

The WIP report uses an expense's settled total (settled_total) as its CAD value when that total is above 0. Otherwise, with no currency or with currency CAD, it uses the expense total. A zero-total expense in another currency is worth 0. Any other foreign-currency expense has no CAD value and is counted as unpriced.

## Expense value and unpriced expenses per job

Lines 31-34

The WIP report sums the CAD values of a job's counted expenses into the expense value, rounded to 2 decimals. Expenses without a CAD value add nothing to the sum. They are counted in unpriced expenses (unpriced_expenses).

## PO spend is totalled in the PO's own currency

Lines 35-39

The WIP report totals committed expenses per purchase order in their native currency (total), not their settled CAD payments. It subtracts this spend from the PO amount before it converts the balance to CAD.

## Remaining balance of a purchase order in WIP

Lines 40-46

The WIP report counts only Active purchase orders of the selected jobs. The remaining balance is the PO total minus spend, and never goes below 0. A Recurring PO uses its approval total (approval_total) in place of its total.

## Recurring PO with no approval total has unknown balance

Lines 42-43

A Recurring purchase order with an approval total of 0 or less and a total above 0 has no balance in the WIP report (null). It is counted as an unpriced PO.

## Currency rate used to convert a PO balance to CAD

Lines 47-52

The WIP report uses a rate of 1 for a PO with no currency or with currency CAD. For another currency it uses the currency's rate if the rate is above 0. If the rate is 0 or missing, the PO has no CAD rate.

## PO value, unpriced POs and estimated POs per job

Lines 57-61

The WIP report sums each job's PO balance times its CAD rate, rounded to 2 decimals per PO, into the PO value. A PO is unpriced (unpriced_pos) when its balance is unknown, or when its balance is above 0 and it has no CAD rate. A PO is estimated (estimated_pos) when its balance is above 0, it has a CAD rate and its currency is not CAD.

## WIP row for each selected job

Lines 62-81

The WIP report gives one row for each selected job, with all its values filled with 0 when it has none. The project value (project_value) is 0 when not set. The as-of date is the end date. The row is flagged no rate sheet (no_rate_sheet) when the job has no rate sheet. Time, expense and PO values are rounded to 2 decimals.
