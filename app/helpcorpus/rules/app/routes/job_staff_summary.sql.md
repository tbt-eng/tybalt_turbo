<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_staff_summary.sql","sha256":"ecd79a289141c1c56051e64fdc86c4005b51c9a72344ee241a1fb4f26d78cb89","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_staff_summary.sql

## Job staff summary prices entries before grouping by employee

Lines 1-17

The job staff summary prices each time entry first and then groups the entries by employee (uid). An employee who works in more than one role on the job gets one row, and each role's entries are priced separately.

## Staff without a profile show blank names in job summary

Lines 5-5

When an employee has no profile, the job staff summary shows the given name and surname as empty text, not as missing values. The row is still listed.

## Job staff summary sums hours, meals hours and estimated hours per employee

Lines 6-14

For each employee the job staff summary adds up the hours (hours), the meals hours (meals_hours), the estimated hours (estimated_hours) and the unpriced hours (unpriced_hours). It also shows the rate sheet name and revision (rate_sheet_name, rate_sheet_revision) from the priced entries.

## Employee value is zero when nothing is priced

Lines 8-8

The value of an employee's row is the sum of the priced entry values. It shows 0 when the sum is missing.

## Job staff summary computes totals across all employees

Lines 9-14

Each row also carries the job-wide total value (total), total estimated hours (total_estimated_hours) and total unpriced hours (total_unpriced_hours). These add up all the employees' rows.

## Employee value status: Incomplete, Estimated or Rate sheet

Lines 20-21

An employee's value status (value_status) is 'Incomplete' when the employee has any unpriced hours. Otherwise it is 'Estimated' when the employee has any estimated hours. Otherwise it is 'Rate sheet'. Incomplete takes precedence over Estimated.

## Total status is based on the whole job

Lines 22-23

The total status (total_status) uses the same rule as the value status but looks at the job-wide totals. It is 'Incomplete' when any hours on the job are unpriced. Otherwise it is 'Estimated' when any hours are estimated. Otherwise it is 'Rate sheet'.

## Percent of total shows only when no hours are unpriced

Lines 24-24

The percent column is the employee's value as a share of the job total, times 100, rounded to 1 decimal place. It is empty when any hours on the job are unpriced. It is also empty when the total is 0, because the division by zero is avoided.

## Job staff summary sort order

Lines 26-26

The job staff summary lists employees by surname, then given name, then uid. Employees without a profile have blank names and so sort first.
