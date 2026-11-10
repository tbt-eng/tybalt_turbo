<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_wip_missing_rates.go","sha256":"ba8e1d823f7a7f15036f63616ff97216a7af13f8c02a7e6e188693dd4c228999","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_wip_missing_rates.go

## Missing rates report defaults to today's date

Lines 31-35

When no as-of date (as_of) is given, the job's work-in-progress (WIP) missing rates report uses today's date in UTC.

## Missing rates as-of date must be valid and not in the future

Lines 36-38

Turbo refuses the missing rates request when the as-of date (as_of) is not a valid YYYY-MM-DD date or is later than today. The message is "as_of must be a valid date no later than today". Today's date itself is accepted.

## Missing rates counts time entries up to the as-of date

Lines 47-66

The report prices the job's time entries with the same pricing rules and date limit as WIP. It has no start date and includes entries up to the as-of date (end_date).

## Missing rates lists hours by role with no sheet rate

Lines 49-63

The report adds up the hours of time entries that have a role but no rate on the job's rate sheet (sheet_rate is null). It groups them by role and sorts the roles by name, then by ID. The role name falls back to the role ID when the name is empty.

## Job without a rate sheet has no missing rate roles

Lines 58-60

When the job has no rate sheet, the report lists no roles. Entries without a role are never listed either.

## Missing rates unknown job returns not found

Lines 70-72

Turbo answers "Job not found" when the query returns no rows, which means the job ID matches no job.

## Missing rates leaves out roles with zero hours

Lines 73-81

The report lists a role only when it has an ID and its hours are not zero. A job with no missing rates returns an empty roles list.

## Missing rates shows the job's rate sheet

Lines 74-76

The report includes the job's rate sheet (ID, name and revision) when the job has one. Otherwise the rate sheet is empty (null).
