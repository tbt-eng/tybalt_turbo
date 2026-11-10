<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/reports/employee_branch_hours.sql","sha256":"fb0e3b7656750842b3ec5f1c18b78e04f6d0084ed01c5bc79876756a32e77df6","model":"claude-sonnet-5-5"} -->

# Rules in app/reports/employee_branch_hours.sql

## Branch hours report counts only time types R and RT

Lines 8-12

The employee branch hours report adds up only time entries whose time type code (tt.code) is 'R' or 'RT'. Entries with any other time type code are left out of the totals.

## Branch hours report skips entries without a valid branch

Lines 10-11

The report keeps a time entry only if its branch matches a real branch (branches.id). An entry with no branch, or with a branch that is not found, is left out of the totals.

## Branch hours report date range includes both end dates

Lines 13-14

The report includes time entries dated on or after the start date ({:start_date}) and on or before the end date ({:end_date}). Entries on the start date and on the end date both count.

## Branch hours report gives one row per employee with total hours

Lines 16-35

The report gives one row per employee. The row shows payroll ID, surname, given name, default branch name and the total (sum of hours). The total is 0 when there are no hours to add.

## Missing names, payroll ID or default branch show as blank

Lines 17-29

When an employee has no payroll ID, surname, given name or default branch, the report shows an empty text in that column instead of a null value. A default branch that is not set does not remove the employee from the report.

## Branch hours report skips staff without admin profile or profile

Lines 24-27

The report lists only staff who have both an admin profile (admin_profiles) and a profile (profiles). Time entries from staff who lack either record are left out.

## Branch hours report sorts by name, then payroll ID length

Lines 36-40

The report sorts rows by surname, then given name. Employees with the same name are sorted by the length of the payroll ID (shorter first), then by the payroll ID. Numeric payroll IDs therefore come out in number order.
