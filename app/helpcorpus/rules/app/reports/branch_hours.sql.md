<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/reports/branch_hours.sql","sha256":"0601c9eb17e77c96d8f075c5dec38f7192e102510dee9ed6f65d6e8d73799a0d","model":"claude-sonnet-5-5"} -->

# Rules in app/reports/branch_hours.sql

## Entries without a job still appear with blank job fields

Lines 11-11

Asked as: time entry without job; blank job fields; hours with no job; entries missing job number

Time entries with no job are kept and grouped under an empty job ID. Their job number, description and client show as blank.

## Hours are summed per job, employee, work branch and staff branch

Lines 12-56

Asked as: total hours per employee; sum hours by job; hours grouped by branch; decimal hours totals

Each row totals the hours (hours) for one job, one employee, one work branch and one staff branch. Hours are treated as decimal numbers.

## Branch Hours counts only time types R and RT

Lines 14-20

Asked as: which time types show; branch hours missing entries; time type R only; filter by time type code

The Branch Hours report includes only time entries whose time type code (tt.code) is 'R' or 'RT'. Time entries of other time types do not appear in the report.

## Entries without a work branch or admin profile are excluded

Lines 16-19

Asked as: missing entries in report; entries dropped from report; no admin profile hours; branch hours not showing hours

Time entries drop out of the report when their branch does not match a branch record, or when the employee has no admin profile. Both are required joins.

## Report includes time worked in or by staff of the branch

Lines 18-23

Asked as: who counts for my branch; staff from my branch hours; default branch hours; entries lent to other branch

An entry qualifies when its work branch is the selected branch, or when the employee's current default branch (admin_profiles.default_branch) is the selected branch. The staff branch is the employee's current default branch, not a branch stored on the entry.

## Branch Hours date range includes the start and end dates

Lines 21-22

Asked as: report date boundaries; start and end day included; hours on last day; date range inclusive

The report includes time entries dated on or after the start date and on or before the end date. Entries on both boundary dates are counted.

## Direction labels internal, supplied or received

Lines 26-30

Asked as: internal supplied received meaning; what does received mean; staff lent out label; direction column explained

A row is 'internal' when both the employee's default branch and the work branch are the selected branch. It is 'supplied' when only the staff branch is the selected branch (staff lent to another branch). Every other row is 'received' (another branch's staff worked in the selected branch).

## Blank values for missing job, client or staff branch

Lines 32-50

Asked as: empty client name; blank surname in report; missing staff branch shown; employee row disappeared

Job number, description, client name, surname, given name and staff branch show as empty text when no matching record exists. The employee must have a profile, otherwise the row is dropped.

## Branch Hours row sort order

Lines 57-64

Asked as: branch hours ordering; why rows out of order; sort order of report; jobs listed first

Rows sort by direction (alphabetically), then rows with a job number come before rows without one. After that they sort by job number, work branch name, surname, given name and employee ID.
