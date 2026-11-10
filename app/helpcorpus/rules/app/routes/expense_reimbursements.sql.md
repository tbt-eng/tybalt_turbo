<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/expense_reimbursements.sql","sha256":"1bba86a0b8a141184c270cd0ee9ee3d7d4f9d5e81019f5b72efb21cac709872b","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/expense_reimbursements.sql

## Mileage amount is recalculated from tiered rates

Lines 10-13

Asked as: mileage total differs from payroll; mileage amount wrong after commit; mileage tier recalculated; mileage total doesn't match report

For a Mileage expense, the amount is calculated from the tiered mileage rates (mileage_total), not read from the expense's total. This matches the payroll expense report. The stored total is fixed at commit time and can differ when mileage expenses that cross a tier boundary are committed out of date order.

## Mileage tiers use cumulative distance since the last reset date

Lines 29-52

Asked as: mileage reset date; cumulative mileage kilometres; mileage rate tiers running total; which mileage rate applies

Cumulative mileage is summed in date order across all the caller's committed mileage expenses, not just the requested pay period. The running total restarts at each mileage reset date. The rate table used is the latest one whose effective date is on or before the expense date.

## Mileage is split across tiers and rounded

Lines 54-77

Asked as: mileage split across tiers; mileage rounding to cents; mileage charged per tier; zero mileage amount

Each mileage expense is charged at each tier's rate for the kilometres that fall inside that tier. The last tier has no upper limit. The total is rounded to 2 decimals and is 0 if no tier applies.

## Pay date is 13 days after pay period ending

Lines 83-84

Asked as: when is pay date; payday after pay period; pay date calculation; how many days until payday

Turbo sets the pay date (pay_date) to the date 13 days after the pay period ending. This is always the Friday after the Saturday on which the pay period ends.

## Allowance description is built from selected meal types

Lines 86-95

Asked as: allowance description blank; meal types listed on expense; breakfast lunch dinner description; description of meals claim

For Allowance and Meals expenses, the description lists the selected allowance types in the order Breakfast, Lunch, Dinner, Lodging, separated by commas. For other payment types, the description is the expense's own description.

## Missing job, client or distance show as blank or zero

Lines 96-99

Asked as: job number is blank; no job on expense; client name missing; distance shows zero

When an expense has no job, the job number, job description and client name are empty strings. A missing distance shows as 0.

## Foreign currency total is 0 for CAD

Lines 100-104

Asked as: foreign currency total zero; USD amount showing 0; CAD currency foreign total; original total in other currency

The foreign currency total (foreign_currency_total) is 0 when the currency is CAD. For any other currency, it is the expense's original total. A missing currency counts as CAD.

## Allowance amount sums the daily rates for selected types

Lines 107-111

Asked as: meal allowance amount wrong; per diem rate calculation; allowance total ignores stored amount; daily meal rates amount

For an Allowance or Meals expense, the amount is the sum of the breakfast, lunch, dinner and lodging rates for the types selected. Turbo takes the rates from the latest expense rates whose effective date is on or before the expense date. It does not use the expense's stored total.

## Foreign currency expenses use settled total when present

Lines 112-113

Asked as: foreign currency reimbursement amount; settled total not used; US dollar expense amount; converted amount for reimbursement

For an Expense or PersonalReimbursement in a currency other than CAD, the amount is the settled total (settled_total). If the settled total is empty or 0, the amount falls back to the expense total. CAD expenses use the total.

## Only five payment types count as reimbursements

Lines 123-126

Asked as: corporate card expense not in reimbursements; why is my card charge missing; fuel card not reimbursed; on account expenses hidden

The reimbursement list includes only expenses with payment type Expense, PersonalReimbursement, Allowance, Meals (the old name for Allowance) or Mileage. OnAccount, CorporateCreditCard, FuelCard and FuelOnAccount expenses are left out because the company paid them.

## Only the caller's own committed expenses are listed

Lines 123-125

Asked as: other people's expenses in list; expense missing from my reimbursement list; uncommitted expense not showing; expense without pay period

The list shows only expenses that belong to the caller (uid) and have been committed. An expense that has no pay period ending is also left out.

## Empty pay period filter lists all pay periods

Lines 127-127

Asked as: show all pay periods; reimbursements filter blank period; list every pay period; filter by pay period ending

When the pay period ending parameter is empty, Turbo lists reimbursements from every pay period. Otherwise it lists only expenses whose pay period ending equals the parameter.
