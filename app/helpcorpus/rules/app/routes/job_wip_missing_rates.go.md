<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_wip_missing_rates.go","sha256":"8e4ca5b958d611b1c2d0609d485d8537dc6e7a08ef1149d3b510327e74ffacfb","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_wip_missing_rates.go

## Missing rates report date defaults to today

Lines 34-38

Asked as: missing rates default date; report date today; as of date blank; which date is used

The job's WIP missing rates report uses an as-of date (as_of). If none is given, it uses today's date in UTC. The report counts time entries up to and including that date.

## Missing rates report refuses invalid or future dates

Lines 39-41

Asked as: invalid as_of date; future date rejected; as of date error; date later than today

Turbo refuses the missing rates request when the as-of date (as_of) is not a valid YYYY-MM-DD date or is later than today. The message is "as_of must be a valid date no later than today". Today itself is accepted.

## Missing rates report needs an existing job

Lines 42-45

Asked as: job not found error; wrong job id missing rates; unknown job for rates report; job does not exist

Turbo refuses the missing rates request when the job does not exist. The message is "Job not found". Turbo checks the date before it checks the job.

## Missing rates grouped by rate sheet and role

Lines 54-71

Asked as: roles without rate; no rate for role; which rate sheet lacks role; hours missing a rate

The report lists each rate sheet that priced the job's time, with the roles that sheet has no rate for. Each role shows its total hours. A role missing a rate in one period of the job's rate sheet history can still be priced in another period. Each entry is placed under the rate sheet of the period its date falls in. The report uses the same pricing rules and date limit as WIP.

## Missing rates ordering and role name fallback

Lines 57-68

Asked as: missing rates sort order; role shows id not name; unnamed role in report; rate sheet listing order

The report sorts rate sheets by name, then revision. Roles within a sheet sort by role name, then role ID. When a role has no name, the report shows the role ID as its name.

## Only entries with a role and no sheet rate count

Lines 65-65

Asked as: entries without role; why entry not in missing rates; time with no role; role blank in report

A time entry shows up as missing a rate only when it has a role and the rate sheet has no rate for that role. Entries without a role are left out of the report.

## Roles with zero net hours are hidden

Lines 66-67

Asked as: zero hours role hidden; role missing from report; zero net hours; negative hours missing rate

A role on a rate sheet is left out of the missing rates report when its hours add up to exactly 0. Negative totals still appear.

## Missing rates report with nothing missing is empty

Lines 75-88

Asked as: no missing rates shown; empty missing rates report; all roles have rates; report has nothing

When no role lacks a rate, the report returns the as-of date and an empty list of rate sheets.
