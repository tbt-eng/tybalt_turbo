<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/reports/payroll_time.sql","sha256":"5dd394d3bea9b01947c705f17e3f2572b2b38b7814a0edc594418d502ed2e7fb","model":"claude-sonnet-5-5"} -->

# Rules in app/reports/payroll_time.sql

## Salaried staff hours over 44 are reported

Lines 12-15

For salaried staff (salary = 1), the salary hours over 44 column shows total hours worked minus 44 when hours worked exceed 44 in the week. Exactly 44 hours gives 0, and so does any non-salaried person.

## Adjusted hours worked for salaried staff

Lines 16-20

For salaried staff, if hours worked plus Stat Holiday plus Bereavement is more than the work week hours (workWeekHours), adjusted hours worked is work week hours minus Stat and Bereavement. Exactly equal to the work week hours is not capped, and adjusted hours then equal hours worked. Missing Stat or Bereavement counts as 0.

## Adjusted hours worked for hourly staff capped at 44

Lines 21-25

For non-salaried staff, adjusted hours worked is capped at 44. Hours of exactly 44 or fewer are reported unchanged.

## Overtime hours only for non-salaried staff over 44 hours

Lines 26-29

Total overtime hours is hours worked minus 44 for non-salaried staff (salary = 0) with more than 44 hours worked. Salaried staff always show 0 here.

## Overtime hours to pay subtracts hours banked

Lines 30-36

For non-salaried staff, overtime hours to pay is the overtime over 44 hours minus the overtime hours banked (RB). The subtraction applies only when banked hours are above 0. Salaried staff always show 0.

## Overtime to pay is not limited to zero

Lines 32-32

Overtime hours to pay is not floored at 0. If banked hours exceed the overtime over 44, the result is negative.

## Report shows salary as TRUE or FALSE

Lines 45-45

The salary column shows TRUE when the person's salary flag is 1 and FALSE otherwise.

## Week ending is displayed as year, month name and day

Lines 48-50

The week ending shows as year, three-letter month and day, such as 2025 Jan 05.

## Time entries take salary and work week hours from the time sheet

Lines 54-54

For regular time entries, salary and work week hours come from the time sheet (ts). The manager is the time sheet's approver.

## Days off rotation counts entries, not hours

Lines 56-56

Each time entry with type code OR adds 1 to days off rotation, whatever its hours.

## Hours worked counts only regular time types R and RT

Lines 57-57

Hours worked (hoursWorked) adds up the hours of entries with time type code R or RT. Missing hours count as 0.

## Time off and overtime bank hours come from time type codes

Lines 58-63

Each time type code feeds its own payroll column. OB is Bereavement, OH is Stat Holiday, OP is PPTO, OS is Sick, OV is Vacation and RB is overtime hours to bank. Missing hours count as 0.

## Overtime payout requested comes from the entry's payout request amount

Lines 64-64

Overtime Payout Requested sums the payout request amount of each entry. A missing amount counts as 0.

## Amendments flag the original weeks they changed

Lines 65-65

The hasAmendmentsForWeeksEnding column lists, as a JSON array, the original week endings of the committed amendments for that person and commit week. It is empty for regular time entries.

## Payroll report only counts committed time sheets

Lines 74-77

The payroll report for a week ending includes time entries only when they belong to a time sheet (tsid is not empty) that has been committed (committed is not empty). Entries on uncommitted time sheets are left out.

## Payroll report excludes payroll IDs matching the placeholder condition

Lines 78-78

Both time entries and amendments are skipped when they match the placeholder payroll ID condition (placeholder_payroll_id_condition). The code in this file does not say what that condition contains.

## Amendments are paid in the week they were committed

Lines 84-86

A committed amendment (committed is not empty) appears in the report for the week in which it was committed (committed_week_ending), not its original week. It is matched against the chosen week ending on that date.

## Amendment rows use admin profile salary settings and committer as manager

Lines 90-103

For amendments, salary and work week hours come from the person's admin profile, not the time sheet. The manager column shows the person who committed the amendment.

## Amendment list of original weeks only uses committed amendments

Lines 110-123

The list of original week endings for amendments is built from committed amendments only, with duplicates removed. It is grouped by person and commit week.

## Report is one row per payroll ID, sorted by ID length

Lines 128-129

Time entries and amendments for the same payroll ID are combined into one row. Rows are sorted by payroll ID length first, then by the ID, so shorter IDs come first.
