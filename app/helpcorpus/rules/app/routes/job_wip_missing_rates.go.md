<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_wip_missing_rates.go","sha256":"f59dd3607664fb94d2fa567704bdfa4eb44a57308a859691309ca49b79c60ba1","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_wip_missing_rates.go

## Missing-rates report date defaults to today

Lines 35-39

Asked as: missing rates report date; report date defaults to today; no as_of date given; which date does the report use

The job's WIP missing-rates report uses the as_of date from the request. When as_of is empty, it uses today's date (UTC).

## Missing-rates report refuses invalid or future as_of date

Lines 40-42

Asked as: invalid as_of date; as_of in the future; future date rejected; report date format error; as_of must be a valid date

Turbo refuses the missing-rates report when as_of is not a valid YYYY-MM-DD date or is later than today. The message is "as_of must be a valid date no later than today". An as_of equal to today is accepted.

## Missing-rates report refuses unknown job

Lines 43-46

Asked as: job not found; unknown job report; wrong job number missing rates; job does not exist

When the job does not exist, Turbo refuses the missing-rates report with the message "Job not found".

## Missing-rates report lists roles with no rate on the priced sheet

Lines 56-68

Asked as: roles without a rate; missing rate on sheet; role has no rate; which roles are unpriced; no rate for role

For one job, the report lists each rate sheet that priced the job's time entries, with the roles that the sheet has no rate for (no sheet_rate). It uses the same pricing rules and date limit as WIP. Each entry is matched to the sheet of the rate sheet period its date falls in, so a role can be missing on one sheet and priced on another. The details of the pricing come from singleJobScopeQuery and jobPricedTimeEntriesQuery, which are defined in other files.

## Missing-rates report shows role name or falls back to role ID

Lines 61-66

Asked as: role shows as ID; role name missing; blank role name; role record not found; unnamed role on report

Each missing role shows its name with its total hours. When the role has no name, or the role record is not found, the report shows the role ID in place of the name.

## Missing-rates report skips entries without a role and zero-hour roles

Lines 67-69

Asked as: time entries without role; zero hours role hidden; negative hours role; role with no hours shown

The report ignores time entries that have no role. It sums the hours per rate sheet and role, and leaves out a role whose total hours equal 0. A role with negative total hours still appears.

## Missing-rates report order and grouping

Lines 70-89

Asked as: rate sheets order; sort missing rate sheets; empty rate sheet list; group roles by sheet; no missing rates shown

The report groups roles under their rate sheet (id, code, name, revision). It sorts sheets by code and then revision, and roles by name and then ID. When no role is missing a rate, the list of rate sheets is empty.

## Missing-rates report covers time up to as_of with no start date

Lines 71-73

Asked as: time since beginning; rates report start date; time up to as_of; all earlier time entries counted; report includes older hours

The report counts time entries from the beginning (the start date is empty) up to and including the as_of date.
