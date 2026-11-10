<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/reports/payroll_time_branch_hours.sql","sha256":"01f984eee3111ac754627cd5bcde8a0987db525e852ba7acef1b81ccfce99b1f","model":"claude-sonnet-5-5"} -->

# Rules in app/reports/payroll_time_branch_hours.sql

## Default branch comes from the employee's admin profile

Lines 9-17

Asked as: default branch name; employee home branch; branch from admin profile; which branch is the default

The report shows the employee's default branch name (defaultBranchName) from their admin profile. The comments say reports.go reduces this branch first when a salaried employee's regular branch hours exceed the cap.

## Payroll branch hours report: which time types are counted

Lines 18-42

Asked as: which time types are counted; time type codes in payroll; missing hours in payroll report; entries left out of report

The report adds up hours only for time types R, RT, RB, OH and OB. R and RT hours count as regular hours (regularHours), RB as banked hours (bankedHours), OH as stat hours (statHours) and OB as bereavement hours (bereavementHours). Entries of any other time type are left out.

## Hours are totalled per payroll ID, branch, salary, week hours and default branch

Lines 18-68

Asked as: employee appears twice; same person on two lines; split rows per employee; payroll totals grouped by branch

The report sums regular, banked, stat and bereavement hours across time entries and amendments. It groups the sums by payroll ID, regular branch, reduction branch, salary, work week hours and default branch. An employee whose rows differ in any of these appears on more than one line. Missing hours count as 0.

## Regular hours without a branch go to 'Unassigned'

Lines 24-24

Asked as: hours with no branch; blank branch regular hours; unassigned branch; time entry missing branch

For R and RT time entries, the report groups hours under the entry's branch name. When the entry has no branch, it uses 'Unassigned' as the branch name (regularBranchName).

## Banked hours are grouped by branch only when a branch is set

Lines 25-25

Asked as: banked hours reduction branch; empty branch for banked time; banked entry no branch; reduction branch blank

For RB (banked) entries, the report records a reduction branch name (reductionBranchName) only when the entry has a branch. Without a branch, the reduction branch name is empty.

## Salary and work week hours for regular entries come from the time sheet

Lines 26-27

Asked as: salary wrong on payroll; work week hours from timesheet; timesheet salary not profile; regular entry salary source

For normal time entries, the salary setting and the work week hours (workWeekHours) come from the committed time sheet snapshot. They do not fall back to the employee's admin profile.

## Regular time entries count only if the time sheet is committed

Lines 33-42

Asked as: uncommitted timesheet hours missing; time entry not in payroll; unsubmitted timesheet payroll; why entry is left out

The report includes a time entry only when it belongs to a time sheet (tsid is not empty) and that time sheet is committed. It selects entries by the entry's week ending date. Entries on uncommitted time sheets are left out.

## Payroll ID filter excludes some employees

Lines 43-43

Asked as: employee excluded from payroll; payroll id filter; hide employee from payroll report; payroll ID condition

Both time entries and amendments are filtered by a payroll ID condition that is passed in as a placeholder. Rows that match that condition are left out. The condition itself is not defined in this file.

## Amendments fall back to the admin profile for salary and work week hours

Lines 50-51

Asked as: amendment uses profile salary; amendment no timesheet salary; amendment work week hours wrong; fallback to admin profile

For a time amendment, salary and work week hours come from its linked time sheet when there is one. When the amendment has no linked time sheet value, Turbo uses the employee's admin profile values instead.

## Amendments count only when committed, by their committed week ending

Lines 57-65

Asked as: amendment not in report; adjustment wrong week payroll; amendment committed week date; uncommitted adjustment missing

The report includes a time amendment only when it is committed. It selects amendments by their committed week ending (committed_week_ending), not by the original week. An uncommitted amendment is left out.

## Report row order

Lines 69-69

Asked as: payroll report sort order; row order payroll; why payroll IDs out of order; sorting of report lines

Rows are sorted by the length of the payroll ID first, then by the payroll ID, then by regular branch name and reduction branch name. Shorter payroll IDs come first.
