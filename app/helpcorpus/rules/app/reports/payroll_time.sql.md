<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/reports/payroll_time.sql","sha256":"5dd394d3bea9b01947c705f17e3f2572b2b38b7814a0edc594418d502ed2e7fb","model":"claude-sonnet-5-5"} -->

# Rules in app/reports/payroll_time.sql

## Salaried staff hours over 44 are reported

Lines 12-15

Asked as: salaried overtime hours; hours over 44 salaried; salaried staff extra hours; salary hours over 44

For salaried staff (salary = 1), the salary hours over 44 column shows total hours worked minus 44 when hours worked exceed 44 in the week. Exactly 44 hours gives 0, and so does any non-salaried person.

## Adjusted hours worked for salaried staff

Lines 16-20

Asked as: adjusted hours salaried; stat holiday reduces hours; bereavement hours subtracted; hours over work week cap

For salaried staff, if hours worked plus Stat Holiday plus Bereavement is more than the work week hours (workWeekHours), adjusted hours worked is work week hours minus Stat and Bereavement. Exactly equal to the work week hours is not capped, and adjusted hours then equal hours worked. Missing Stat or Bereavement counts as 0.

## Adjusted hours worked for hourly staff capped at 44

Lines 21-25

Asked as: hourly hours capped at 44; hourly staff max 44 hours; hourly adjusted hours; non salaried hours cap

For non-salaried staff, adjusted hours worked is capped at 44. Hours of exactly 44 or fewer are reported unchanged.

## Overtime hours only for non-salaried staff over 44 hours

Lines 26-29

Asked as: overtime hours calculation; hourly overtime over 44; overtime hours not showing; total overtime hours wrong

Total overtime hours is hours worked minus 44 for non-salaried staff (salary = 0) with more than 44 hours worked. Salaried staff always show 0 here.

## Overtime hours to pay subtracts hours banked

Lines 30-36

Asked as: overtime to pay banked; banked hours subtracted from overtime; overtime owed minus bank; overtime payable after bank

For non-salaried staff, overtime hours to pay is the overtime over 44 hours minus the overtime hours banked (RB). The subtraction applies only when banked hours are above 0. Salaried staff always show 0.

## Overtime to pay is not limited to zero

Lines 32-32

Asked as: negative overtime to pay; overtime pay below zero; banked more than overtime; overtime to pay negative number

Overtime hours to pay is not floored at 0. If banked hours exceed the overtime over 44, the result is negative.

## Report shows salary as TRUE or FALSE

Lines 45-45

Asked as: salary column true false; salary shows true; why salary true or false; salary flag display

The salary column shows TRUE when the person's salary flag is 1 and FALSE otherwise.

## Week ending is displayed as year, month name and day

Lines 48-50

Asked as: week ending date format; week ending date looks wrong; date shows year month day; week ending display format

The week ending shows as year, three-letter month and day, such as 2025 Jan 05.

## Time entries take salary and work week hours from the time sheet

Lines 54-54

Asked as: where manager comes from; who is the manager on entries; timesheet approver manager; salary from timesheet

For regular time entries, salary and work week hours come from the time sheet (ts). The manager is the time sheet's approver.

## Days off rotation counts entries, not hours

Lines 56-56

Asked as: days off rotation count; rotation days off wrong; off day entries counted; OR entries rotation

Each time entry with type code OR adds 1 to days off rotation, whatever its hours.

## Hours worked counts only regular time types R and RT

Lines 57-57

Asked as: hours worked total wrong; which time types count hours; regular hours not adding up; R and RT hours total

Hours worked (hoursWorked) adds up the hours of entries with time type code R or RT. Missing hours count as 0.

## Time off and overtime bank hours come from time type codes

Lines 58-63

Asked as: bereavement column hours; sick hours in payroll; vacation hours missing from report; bank hours column source

Each time type code feeds its own payroll column. OB is Bereavement, OH is Stat Holiday, OP is PPTO, OS is Sick, OV is Vacation and RB is overtime hours to bank. Missing hours count as 0.

## Overtime payout requested comes from the entry's payout request amount

Lines 64-64

Asked as: overtime payout amount; payout request missing; overtime payout total wrong; payout requested not showing

Overtime Payout Requested sums the payout request amount of each entry. A missing amount counts as 0.

## Amendments flag the original weeks they changed

Lines 65-65

Asked as: amended weeks listed; which weeks were amended; amendment original week list; changed weeks json

The hasAmendmentsForWeeksEnding column lists, as a JSON array, the original week endings of the committed amendments for that person and commit week. It is empty for regular time entries.

## Payroll report only counts committed time sheets

Lines 74-77

Asked as: unsubmitted timesheet missing from payroll; uncommitted entries not in report; entries left out of payroll; timesheet not committed payroll

The payroll report for a week ending includes time entries only when they belong to a time sheet (tsid is not empty) that has been committed (committed is not empty). Entries on uncommitted time sheets are left out.

## Payroll report excludes payroll IDs matching the placeholder condition

Lines 78-78

Asked as: placeholder payroll id skipped; missing entries payroll id; why some entries dropped; payroll id excluded

Both time entries and amendments are skipped when they match the placeholder payroll ID condition (placeholder_payroll_id_condition). The code in this file does not say what that condition contains.

## Amendments are paid in the week they were committed

Lines 84-86

Asked as: amendment paid which week; late correction which pay week; amendment shows in wrong week; commit week amendment payroll

A committed amendment (committed is not empty) appears in the report for the week in which it was committed (committed_week_ending), not its original week. It is matched against the chosen week ending on that date.

## Amendment rows use admin profile salary settings and committer as manager

Lines 90-103

Asked as: amendment salary source; amendment manager shows committer; admin profile salary amendment; who is manager on amendment

For amendments, salary and work week hours come from the person's admin profile, not the time sheet. The manager column shows the person who committed the amendment.

## Amendment list of original weeks only uses committed amendments

Lines 110-123

Asked as: original weeks duplicates; amendment weeks list empty; list of corrected weeks; repeated week endings amendment

The list of original week endings for amendments is built from committed amendments only, with duplicates removed. It is grouped by person and commit week.

## Report is one row per payroll ID, sorted by ID length

Lines 128-129

Asked as: payroll rows order; payroll id sort order; duplicate rows same payroll id; combine rows per employee

Time entries and amendments for the same payroll ID are combined into one row. Rows are sorted by payroll ID length first, then by the ID, so shorter IDs come first.
