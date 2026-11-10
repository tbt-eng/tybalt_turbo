<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_divisions_summary.sql","sha256":"b933feef14babacdab0b03dba8dbfc5ef949676ae491915ef95e1321af893701","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_divisions_summary.sql

## Job divisions summary groups priced time entries by division

Lines 1-17

The job divisions summary prices each time entry first and then groups the entries by division (te.division). It prices before grouping because an employee can use more than one role on the same job. It sums hours, value, estimated hours and unpriced hours for each division.

## Entries with no division show blank division code and name

Lines 2-8

When a time entry has no division, the summary shows a blank division code and a blank division name instead of NULL. The value is also set to 0 when there is nothing to sum.

## Division value status: Incomplete, Estimated or Rate sheet

Lines 19-20

A division's value status is 'Incomplete' when it has any unpriced hours (unpriced_hours > 0). Otherwise it is 'Estimated' when it has any estimated hours (estimated_hours > 0). Otherwise it is 'Rate sheet'.

## Total status uses the same rule across all divisions

Lines 21-22

The total status is 'Incomplete' when the job has any unpriced hours in any division. Otherwise it is 'Estimated' when any division has estimated hours. Otherwise it is 'Rate sheet'.

## Division percent of total value is blank when any hours are unpriced

Lines 23-23

Turbo shows a division's percent as its value divided by the job total value, times 100, rounded to 1 decimal place. It shows no percent at all, for every division, when the job has any unpriced hours. When the total is 0, the percent is also blank, because Turbo does not divide by zero.

## Divisions are sorted by division code, then division

Lines 25-25

The summary lists divisions in order of division code, then division id. Entries with no division have a blank code and so come first.
