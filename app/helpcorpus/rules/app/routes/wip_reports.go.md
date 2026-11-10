<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/wip_reports.go","sha256":"496d952e98d74a08502c60c890cbcb8a0cbc934daf5f9dd6eb9df52f538faaa5","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/wip_reports.go

## WIP report rows are ordered by job number

Lines 16-25

Each report row shows the job number, description, client name, manager's full name, branch ID and branch name. Client, manager and branch are empty text when the job has none. Rows are sorted by job number, then by job ID.

## Jobs without a positive project value are excluded from the WIP report

Lines 51-58

The WIP report counts a job as excluded when its project value (project_value) is 0, negative or empty. The report returns this count as excluded_project_value. Only jobs with a project value above 0 are included and priced.

## WIP report skips pricing when no job has a value

Lines 59-61

When no job in the report scope has a project value above 0, Turbo returns the report straight away with no rows. It does not price time entries, expenses or POs in that case. The excluded count is still given.

## Who may view management WIP reports

Lines 66-87

For every mode except "my", Turbo allows the report only to people with the kpi claim or the admin claim, or to a manager of at least one branch. Anyone else gets: "You need to be a branch manager or hold the kpi or admin claim to view this report".

## The "my" WIP report has no claim check

Lines 79-87

In "my" mode Turbo skips the management access check, so any signed-in user can request the report. The report uses the user's own ID as the scope.

## Division WIP report requires a division

Lines 90-92

In division mode, a request with no division is refused with "Select a division".

## Selected branch or division must exist

Lines 93-103

In branch mode Turbo checks that the chosen branch exists, and in division mode that the chosen division exists. If it does not, Turbo refuses with "Unknown report selection". An empty branch is not checked, and a branch given in division mode, or a division given in branch mode, is not checked.

## require_time filter accepts only true or false in division mode

Lines 104-107

In division mode, a require_time value other than empty, "true" or "false" is refused with "require_time must be true or false". In other modes the value is not validated.

## WIP report is as of today, UTC, from the start

Lines 108-113

The report is dated today in UTC (as_of, format YYYY-MM-DD). It uses an empty start date and today as the end date, so it covers everything up to today.

## Division only selects jobs; each row keeps its full WIP

Lines 109-114

A division chooses which jobs appear in the report. Each job's row still includes all of that job's WIP factors.

## require_time defaults to on

Lines 113-113

Turbo sets require_time to true unless the request says "false". An empty or missing value counts as true.
