<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/expense_reimbursements.sql","sha256":"1bba86a0b8a141184c270cd0ee9ee3d7d4f9d5e81019f5b72efb21cac709872b","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/expense_reimbursements.sql

## Mileage amount is recalculated from tiered rates

Lines 10-13

For a Mileage expense, the amount is calculated from the tiered mileage rates (mileage_total), not read from the expense's total. This matches the payroll expense report. The stored total is fixed at commit time and can differ when mileage expenses that cross a tier boundary are committed out of date order.

## Mileage tiers use cumulative distance since the last reset date

Lines 29-52

Cumulative mileage is summed in date order across all the caller's committed mileage expenses, not just the requested pay period. The running total restarts at each mileage reset date. The rate table used is the latest one whose effective date is on or before the expense date.

## Mileage is split across tiers and rounded

Lines 54-77

Each mileage expense is charged at each tier's rate for the kilometres that fall inside that tier. The last tier has no upper limit. The total is rounded to 2 decimals and is 0 if no tier applies.

## Pay date is 13 days after pay period ending

Lines 83-84

Turbo sets the pay date (pay_date) to the date 13 days after the pay period ending. This is always the Friday after the Saturday on which the pay period ends.

## Allowance description is built from selected meal types

Lines 86-95

For Allowance and Meals expenses, the description lists the selected allowance types in the order Breakfast, Lunch, Dinner, Lodging, separated by commas. For other payment types, the description is the expense's own description.

## Missing job, client or distance show as blank or zero

Lines 96-99

When an expense has no job, the job number, job description and client name are empty strings. A missing distance shows as 0.

## Foreign currency total is 0 for CAD

Lines 100-104

The foreign currency total (foreign_currency_total) is 0 when the currency is CAD. For any other currency, it is the expense's original total. A missing currency counts as CAD.

## Allowance amount sums the daily rates for selected types

Lines 107-111

For an Allowance or Meals expense, the amount is the sum of the breakfast, lunch, dinner and lodging rates for the types selected. Turbo takes the rates from the latest expense rates whose effective date is on or before the expense date. It does not use the expense's stored total.

## Foreign currency expenses use settled total when present

Lines 112-113

For an Expense or PersonalReimbursement in a currency other than CAD, the amount is the settled total (settled_total). If the settled total is empty or 0, the amount falls back to the expense total. CAD expenses use the total.

## Only five payment types count as reimbursements

Lines 123-126

The reimbursement list includes only expenses with payment type Expense, PersonalReimbursement, Allowance, Meals (the old name for Allowance) or Mileage. OnAccount, CorporateCreditCard, FuelCard and FuelOnAccount expenses are left out because the company paid them.

## Only the caller's own committed expenses are listed

Lines 123-125

The list shows only expenses that belong to the caller (uid) and have been committed. An expense that has no pay period ending is also left out.

## Empty pay period filter lists all pay periods

Lines 127-127

When the pay period ending parameter is empty, Turbo lists reimbursements from every pay period. Otherwise it lists only expenses whose pay period ending equals the parameter.
