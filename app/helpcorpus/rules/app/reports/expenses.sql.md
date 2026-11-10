<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/reports/expenses.sql","sha256":"3bdff1c35b9312f2cb29366bf737ea5c3a58cf13959248fa45dadc2aa579f9aa","model":"claude-sonnet-5-5"} -->

# Rules in app/reports/expenses.sql

## Report shows employee, approver, job, client, PO and vendor names

Lines 5-24

Asked as: blank employee name on report; missing job number in report; empty columns in expense report; client name not showing

Each report row shows the employee's payroll ID and name, the approver's name, and the division code. It also shows the job number and description, client, PO number and vendor. Missing links leave these columns empty, but the expense still appears.

## Report Month column comes from the expense date

Lines 11-14

Asked as: wrong month on report; date shows wrong day; year column in report; month column wrong

The expense report shows the expense's Date as the day of the month (two digits), Month as a three-letter abbreviation (Jan to Dec) and Year as four digits. All three come from the expense date (e2.date).

## Report splits Total into subtotal and Ontario HST at 13/113

Lines 15-17

Asked as: HST amount on report; subtotal and tax wrong; why is tax taken out; Ontario tax in total

Turbo treats the Total as already including 13% Ontario HST. It calculates the HST (calculatedOntarioHST) as Total × 13 / 113, rounded to 2 decimals. It calculates the subtotal (calculatedSubtotal) as Total minus that HST, rounded to 2 decimals. The employee reimbursements page repeats this calculation, and the two must stay in sync.

## Currency defaults to CAD and foreign amount shown separately

Lines 41-45

Asked as: expense in US dollars; foreign currency column empty; no currency on expense; which currency is reported

An expense with no currency is reported as CAD. For any other currency, the report shows the original amount in the foreign currency total (foreign_currency_total). For CAD it leaves this column empty.

## Mileage expense total uses the tiered mileage calculation

Lines 46-47

Asked as: mileage amount wrong; km claim total; mileage pay calculation; mileage total incorrect

For a Mileage expense, the report Total is the mileage total (mileage_total) calculated from the distance claimed. The foreign-currency and expense total (e.total) are not used for Mileage.

## Allowance and Meals total comes from the allowance rates

Lines 48-58

Asked as: meal allowance amount; per diem total wrong; allowance amount on report; meals total incorrect

For an Allowance or Meals expense, the report Total is the allowance total (allowance_total). The report Description lists the allowance types instead of the expense description.

## Foreign-currency expenses report Total in CAD using settled total

Lines 49-54

Asked as: converted total in CAD; foreign expense CAD total; settled amount in report; USD expense wrong dollars

For an expense in a currency other than CAD, the report Total is the settled total (settled_total, in CAD). When settled_total is 0 or empty, Total falls back to the expense total (e.total). CAD expenses always use the expense total.

## Foreign-currency Total falls back to the original amount when settled total is empty or 0

Lines 52-52

Asked as: settled total is zero; foreign expense shows zero; total shows original amount; foreign total missing

For a foreign-currency expense that is not mileage, allowance or meals, Total is the settled total (CAD). If the settled total is empty or 0, Total uses the expense's original total instead.

## Mileage rates come from the latest rate table in effect

Lines 71-83

Asked as: which mileage rate applies; old mileage rate used; mileage rate changed date; mileage rate by date

Each mileage expense uses the latest expense rates (expense_rates) whose effective date is on or before the expense date. Rates are set per distance tier.

## Mileage tier upper bound is the next tier's lower bound; top tier is open-ended

Lines 71-79

Asked as: highest mileage tier limit; no maximum km tier; top mileage tier cap

Each mileage tier runs from its starting distance to the next tier's starting distance. The highest tier has no upper bound, and Turbo treats it as 1e9 km when pricing.

## Mileage cumulative distance is counted per employee between reset dates, ordered by date

Lines 86-107

Asked as: running kilometres total; km count per employee; employee distance total wrong

The cumulative distance (end_distance) for a mileage expense is the running sum of distance for the same employee (uid), ordered by expense date. The sum restarts at each mileage reset date. The cumulative total decides which mileage tier a claim falls in.

## Mileage resets on reset dates, so cumulative distance restarts

Lines 93-107

Asked as: mileage counter reset; kilometres start over; when mileage resets; yearly km reset

Cumulative mileage for each staff member (uid) adds up the distances of their committed Mileage expenses in date order. The running total restarts at the latest mileage reset date (mileage_reset_dates) on or before the expense date.

## Mileage cumulative distance counts only committed expenses

Lines 108-115

Asked as: uncommitted trips counted; draft km included; unsubmitted mileage in total

Only committed Mileage expenses count toward a staff member's cumulative distance. The report calculates their mileage from the same pay period or week filter as other expenses.

## Mileage tier rate set is chosen by latest effective date on or before expense date

Lines 108-112

Asked as: mileage rates effective date; which rate table used; rate row for trip

Each mileage expense uses the expense rates row with the latest effective date that is on or before the expense date. Tier rates come from the mileage JSON of that row.

## Mileage cumulative distance is computed from the employee's committed mileage in the same period filter

Lines 113-134

Asked as: earlier km from last period; previous pay period mileage counted; old trips in mileage total

Mileage totals are computed for employees who have at least one mileage expense in the selected period or week. The cumulative sum includes only their committed mileage expenses, even those outside the period. Only the expenses in the selected period or week get a row.

## Expense report is filtered to one pay period or week

Lines 116-120

Asked as: report for one pay period; expenses for one week; filter report by week

The report covers only expenses whose chosen date column equals the chosen value, such as a pay period ending date. Mileage totals are calculated only for staff who have a mileage expense in that period.

## Mileage claim crossing a tier is priced piecewise

Lines 144-198

Asked as: mileage split across tiers; part of trip at other rate; trip crosses tier boundary

When a mileage claim crosses a tier boundary, each part of the distance is priced at its own tier rate. The expense total is the sum of the kilometres in each tier times that tier's rate, rounded to 2 decimals. A claim with no overlapping tier totals 0. The top tier has no upper limit.

## Mileage tier overlap includes only tiers the claim reaches

Lines 144-157

Asked as: km exactly at tier boundary; kilometres on tier edge; trip ends at tier limit

A tier applies to a claim only when the claim's cumulative end distance is greater than the tier's lower bound and the claim's start distance is less than the tier's upper bound. A claim that ends exactly at a tier's lower bound gets no kilometres in that tier.

## Mileage total is rounded to 2 decimals and is 0 when no tier applies

Lines 172-198

Asked as: mileage total is zero; mileage rounded to cents; no tier mileage zero

The mileage total for an expense is the sum of kilometres in each tier times that tier's rate, rounded to 2 decimals. If no tier overlaps, the total is 0.

## Allowance total adds one rate per selected meal or lodging type

Lines 207-232

Asked as: add breakfast lunch dinner; lodging allowance added; meal rates added together

An Allowance or Meals expense totals the rates for each type it selects: Breakfast, Lunch, Dinner and Lodging. The rates come from the latest allowance rates in effect on or before the expense date. The description lists the selected types in the fixed order Breakfast, Lunch, Dinner, Lodging.

## Allowance description lists Breakfast, Lunch, Dinner, Lodging in fixed order

Lines 223-229

Asked as: allowance description order; meal types listed wrong; description replaced by meals

For Allowance and Meals expenses, the description is built from the selected types in the fixed order Breakfast, Lunch, Dinner, Lodging, separated by spaces. It replaces the expense's own description in the report.

## Allowance rates use the latest rate row on or before the expense date

Lines 231-231

Asked as: allowance blank no rate; meal total empty; no allowance rate found

Allowance and Meals expenses use the expense rates row with the latest effective date on or before the expense date. If no such row exists, the total is empty.

## Expense report only includes committed expenses

Lines 236-237

Asked as: expense not showing in report; missing expense from report; draft expense not in report

The expenses report lists only expenses that have been committed (committed is not empty). Expenses that are not yet committed do not appear in it.

## Report excludes rows matching the payroll ID condition

Lines 247-247

Asked as: employee left out of report; payroll ID excluded from report; row missing for payroll

The report leaves out any expense row that matches the payroll ID exclusion condition (placeholder_payroll_id_condition). The condition itself is set outside this file.

## Report rows are sorted by date, payroll ID and total

Lines 248-248

Asked as: report row order; rows sorted by date; sort expense report

The expense report sorts rows by expense date, then payroll ID, then Total.
