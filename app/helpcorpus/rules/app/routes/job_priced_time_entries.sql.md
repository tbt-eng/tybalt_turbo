<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_priced_time_entries.sql","sha256":"c0c8c852f6e44fb0d4a29e789973caeae9e2629e0466e626562bc8639f1042ff","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_priced_time_entries.sql

## Rate sheet revision assigned to the job is used even if inactive

Lines 2-4

Turbo prices with the rate sheet revision assigned to the job, even when that revision is inactive. A rate change therefore recalculates the values of past entries.

## Value uses entry hours only, not meals or amendments

Lines 3-4

Turbo prices the entry's hours (hours) only. Meals hours (meals_hours) and amendments are not part of the value. The selection of committed and uncommitted entries is not changed by these rules.

## Time entry value uses the job's rate sheet rate first

Lines 8-23

Turbo values each time entry as hours times the charge-out rate (sheet_rate) for the entry's role on the job's rate sheet. When the rate sheet has no rate for that role, it falls back to the employee's default charge-out rate (default_charge_out_rate).

## Zero or missing default charge-out rate gives no value

Lines 9-10

Turbo uses an employee's default charge-out rate only when it is greater than 0. A zero or missing default cannot supply a dollar value, so the entry has no default rate.

## Report counts entries only within the date range, inclusive

Lines 13-20

Turbo includes time entries dated on or after the start date and on or before the end date. Both the start and end dates are included. Only entries on the selected jobs are priced.

## Entries with no rate at all count as unpriced hours

Lines 23-25

When neither a rate sheet rate nor a usable default rate exists, Turbo adds the entry's hours to unpriced_hours. The entry's value is empty (NULL), not zero.

## Entries priced from the default rate count as estimated hours

Lines 24-24

When a role has no rate on the rate sheet but the employee has a usable default rate, Turbo adds the entry's hours to estimated_hours. Otherwise estimated_hours is 0 for that entry.
