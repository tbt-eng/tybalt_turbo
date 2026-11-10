<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/reports/reports.go","sha256":"55d63fb74b507adbf0efe47d91c81b782735e4a45e4b3633ac7a8deea9988826","model":"claude-sonnet-5-5"} -->

# Rules in app/reports/reports.go

## Payroll reports skip placeholder payroll IDs by condition

Lines 53-96

The payroll time report and the expense report add a SQL condition built by utilities.GeneratedPayrollPlaceholderSQLCondition on ap.payroll_id. The condition replaces the placeholder token {:placeholder_payroll_id_condition} in their queries. This file does not show which payroll IDs the condition excludes.

## Hourly overtime cap is 44 regular hours

Lines 57-57

Hourly staff (not salary) have a regular cap of 44 hours in total across all branches (hourlyPayrollBranchRegularCap). Turbo moves hours above the cap from the branch columns to the OT branch columns.

## Payroll time report is CSV with fixed columns then branches

Lines 59-59

The payroll time report is returned as CSV. Its fixed columns come first (payrollId through salary), followed by the branch columns and then the OT branch columns.

## Unassigned branch column is always present

Lines 98-124

The payroll branch columns always include a branch named "Unassigned". Turbo adds it at the end when no branch has that name. Branches with a blank name are left out.

## Branch columns come from the branches table, sorted by name

Lines 98-110

Turbo reads the branch names from the branches table in name order and makes one payroll report column for each. It leaves out branches whose name is empty or only spaces.

## Overtime branch columns carry an OT prefix

Lines 104-110

Each branch gets a matching overtime column named "OT " followed by the branch name, for example "OT Unassigned". Branches with a blank name are left out of the list.

## Overtime columns repeat every branch with an OT prefix

Lines 126-136

The payroll report adds an overtime column for every branch column. The column is named "OT " followed by the branch name, such as "OT Unassigned".

## Branch hours are summed per payroll ID

Lines 147-177

Turbo adds up regular hours by branch for each payroll ID. A person counts as salaried if any of their rows is salaried. The first work week hours above 0 and the first default branch found are used. Stat and bereavement hours are totalled too.

## Banked overtime reduces branch hours of hourly staff

Lines 165-174

For hourly staff, banked overtime hours reduce the regular hours of the branch that the hours are tied to. Banked hours with no branch reduce the branch with the most hours. When the named branch has too few hours, the rest comes from the branch with the most hours, ties by name. Salaried staff are not changed.

## Banked hours below the epsilon are ignored

Lines 168-174

Banked hours count only when they exceed 0.000000001. Banked hours with a reduction branch reduce that branch. Banked hours with no reduction branch, or a blank one, go to an unassigned pool that reduces the largest branches.

## Branch hours are adjusted in a fixed order

Lines 179-191

For each payroll ID, Turbo takes banked overtime out of the branch hours first. Then it splits hourly overtime above 44 hours. Last, it caps salaried hours to the work week. Then it formats the totals for the report.

## Banked hours reduce before the 44 hour overtime split

Lines 180-183

Turbo first subtracts banked overtime from the branch hours, then splits the overtime above 44 hours.

## Banked overtime is removed from its own branch first

Lines 196-211

For hourly staff, Turbo takes each branch's banked hours out of the regular hours of that same branch, with the branches in alphabetical order. If that branch has too few hours, Turbo takes the rest from the branch with the most hours. Salaried staff get no banked reduction.

## Hourly staff overtime above 44 regular branch hours

Lines 213-244

For hourly staff, Turbo moves branch hours above 44 into the "OT" columns. It counts hours in half-hour units. Salaried staff get no overtime split.

## Overtime hours reduce the branch columns

Lines 236-243

For hourly staff, the overtime hours go in the OT branch column. Turbo subtracts them from that branch's regular column, and the column never goes below 0.

## Overtime is spread evenly across branches

Lines 246-289

Turbo takes the overage from branches in turn, starting with the branch with the most hours, ties by name. Each branch gives an equal share. Leftover half-hours go to the branches first in that order. A branch with fewer hours gives all it has and the rest is shared by the others.

## Overtime shares are shared in a fixed way

Lines 246-316

Turbo splits the overage across the branches that have hours, in half-hour units. When the overage does not divide equally, the largest branches get the extra half hours first, and ties go in alphabetical order. A branch with too few hours gives all it has, and Turbo shares the remainder among the other branches.

## Overtime split uses half-hour units

Lines 318-324

Turbo rounds branch hours to the nearest half hour when it splits overtime. The 44 hour cap is therefore 88 half-hour units. Each overtime branch amount is reported in hours.

## Salaried branch hours are capped to the work week

Lines 326-355

For salaried staff, Turbo caps total branch hours at the work week hours. Stat holiday and bereavement hours lower the cap when the total plus those hours is over the work week. The cap never goes below 0. The excess comes off the default branch first, then off the branch with the most hours.

## Missing work week hours skips the salary cap

Lines 326-333

When a salaried person has no work week hours, or the value is effectively zero, Turbo does not cap their branch hours.

## Salaried branch hours reduced for stat and bereavement time

Lines 340-354

For a salaried person, Turbo lowers the target branch hours when stat holiday plus bereavement hours added to the branch total exceed the work week hours. The target becomes work week hours minus those paid time off hours, and never goes below 0. Any excess over the target comes off the person's default branch first.

## Excess salary hours reduce default branch first, then largest

Lines 357-413

When Turbo cuts branch hours, it first reduces the named branch (for salaried staff, the default branch). If hours still remain to cut, it takes them from the branch with the most hours, and a tie goes to the name that sorts first. It stops when no branch has hours left.

## Branch columns default to 0

Lines 415-436

In the payroll time report every branch and OT branch column starts at "0" for each row. Turbo fills in hours only for branches the person has hours in, matched by payroll ID. Hours near zero show as 0.

## Small branch hour values show as 0

Lines 415-420

Branch hour values that are less than 0.000000001 hours away from zero show as 0. Other values show with no extra trailing digits.

## Payroll time report accepts only week 1 or 2

Lines 444-447

The payroll time report refuses any week value other than 1 or 2. Turbo shows the message "week must be either 1 or 2".

## Payroll week error message

Lines 445-447

A payroll time report with a week other than 1 or 2 shows "week must be either 1 or 2".

## Payroll week 1 reports the week before the pay period end

Lines 454-457

The payroll time report for week 1 uses the Saturday that is 7 days before the pay period ending date in the request. Week 2 uses the date as given.

## Payroll time report lists branches and overtime branches as columns

Lines 469-492

The payroll time report is a CSV. After the base payroll columns it adds one hours column for each branch in name order, then one "OT " column for each branch. Every branch column shows 0 when the person has no hours there.

## Weekly time report uses TBTE and chosen date columns

Lines 497-526

The weekly time report is a CSV for the Saturday in the request. It filters by the company short name "TBTE". It uses the entries and amendments date columns chosen for that report.

## Time reports fail with query or CSV errors

Lines 513-522

If a report query fails, Turbo shows "failed to execute query: " with the reason. If the CSV cannot be built, it shows "failed to generate CSV report: " with the reason.

## Job time report requires a job id

Lines 534-558

The job time report for a job refuses a missing id with the message "id is required". It returns a CSV of all time for that job, for company "TBTE".

## Job report error message

Lines 535-537

A job time report with no job id shows "id is required".

## Time summary by employee report

Lines 563-588

The time summary by employee report is a CSV for one Saturday. It shows each employee's divisions, chargeable (c) and non-chargeable (nc) hours, total, percent chargeable, percent non-chargeable worked, and worked and unworked non-chargeable hours.

## Expense and receipts reports check the two-week cycle only for pay period ending

Lines 594-596

The expense report and the receipts report apply the two-week payroll cycle check only when the date column is pay_period_ending. Other date columns need only a valid Saturday. The weekly time reports and the time summary never apply it.

## Expense report has no payroll cycle check for other date columns

Lines 594-604

The expense report applies the two-week payroll check only when the date column is pay_period_ending. It also leaves out placeholder payroll IDs through a generated condition on the payroll ID.

## Expense report columns and currency

Lines 615-625

The expense report is a CSV. Total stays the CAD amount. The currency and foreign_currency_total columns come last.

## Reports return CSV with a text/csv header

Lines 626-626

The payroll, time, job, summary and expense reports come back as CSV text. The content type is text/csv. They return an error instead if the query or the CSV generation fails.

## Receipts report fetches every row without a limit shown here

Lines 641-653

The receipts report queries receipts for the chosen date value. It builds its list from the rows that come back. This file sets no limit on the number of rows.

## Receipts missing required fields are skipped

Lines 654-676

The receipts zip skips any receipt row that lacks an id, source path, filename, zip filename or hash. A row without a collection id is treated as an expense receipt.

## Receipts without a collection ID use the expense collection

Lines 673-675

If a receipt row has the other required fields but no collection ID, Turbo uses the expense collection for it.

## Receipts report returns a zip link, reusing a cached zip

Lines 678-724

The receipts report returns a URL to a zip of receipts. It reuses a cached zip when one matches the date, report type and receipts list. Otherwise it builds a new zip and caches it.

## Zip failures show an error

Lines 681-684

If Turbo cannot look up the zip cache, it shows "failed to lookup zip cache: " with the reason. If it cannot build the zip, it shows "failed to generate zip archive: " with the reason.

## Date error messages

Lines 730-754

A missing date shows "date_column_value is required". An invalid date shows "date_column_value must be a valid date". A date that is not a Saturday shows "date_column_value must be a Saturday". A payroll date off the two-week cycle shows "date_column_value must be a multiple of 2 weeks before or after 2025-03-01 when forPayroll is true".

## Report dates must be given and valid

Lines 731-740

Reports refuse a missing date with "date_column_value is required". They refuse a date that does not parse as YYYY-MM-DD with "date_column_value must be a valid date".

## Report date must be a Saturday

Lines 743-745

Every report that takes a date refuses a date that is not a Saturday. Turbo shows "date_column_value must be a Saturday".

## Payroll dates must fall on a two-week cycle

Lines 749-754

For payroll, the Saturday must be a multiple of 14 days from the payroll epoch (2025-03-01). Otherwise Turbo shows "date_column_value must be a multiple of 2 weeks before or after 2025-03-01 when forPayroll is true".
