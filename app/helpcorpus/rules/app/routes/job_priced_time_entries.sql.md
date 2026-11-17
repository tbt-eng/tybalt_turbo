<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_priced_time_entries.sql","sha256":"b5b2b93b59e5166f679462249ef9d3908e8d859be44ecc5d6f38fc9d783c8b56","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_priced_time_entries.sql

## Time entry priced with the rate sheet in effect on its date

Lines 1-3

Asked as: wrong rate on old time; rate sheet by date; which rate applies to entry; time priced with old rates

For WIP and the staff and division summaries, Turbo prices each time entry with the rate sheet that is in effect on the entry's date. It takes the sheet from the project's rate sheet history (job_rate_sheets). A period runs from its effective date up to, but not including, the next period's effective date.

## Project with no rate sheet history uses the job's rate sheet

Lines 3-5

Asked as: no rate history for project; imported project rates; job rate sheet used; project has no rate sheet

When a project has no rows in its rate sheet history, Turbo prices all its time with the job's own rate sheet (jobs.rate_sheet). Imported projects have no history until it is written. An empty rate sheet value on the job counts as no sheet.

## Pricing uses the assigned rate sheet revision even if inactive

Lines 4-5

Asked as: inactive rate sheet still used; old revision recalculated; changing rates changes past value; rate change updates old time

Turbo prices time with the revision of the rate sheet that is assigned, even when that revision is inactive. A change to a sheet's rates therefore recalculates the value of past time.

## Amendments and meals are not part of the time value

Lines 5-8

Asked as: expenses in time value; meals counted in value; why value differs from hours; amendments included in total

Turbo keeps its existing selection of committed and uncommitted time entries when it prices time. Amendments and meals are not part of the value. The value is the entry's hours multiplied by the rate.

## Opening period covers entries dated before the first rate sheet

Lines 16-24

Asked as: entries before first rate sheet; missing rate before start date; time dated before sheet; opening rate period

Turbo adds an opening period with an empty start date to each selected job that has no history row with an empty effective date. This period is before every date, so no entry of the job is dropped. When the job has history rows, the opening period has no rate sheet, and entries before the first effective date get no sheet rate.

## Empty job rate sheet is treated as no rate sheet

Lines 17-19

Asked as: blank rate sheet on job; job has no sheet; empty rate sheet field

For a job with no history rows, an empty rate sheet on the job (jobs.rate_sheet) counts as none. The opening period then has no rate sheet.

## Each rate sheet period runs until the next period starts

Lines 25-42

Asked as: rate change on boundary date; which period an entry falls in; end date of rate period; rate sheet switch day

A project's rate sheet history is split into periods. A period starts on its effective date and ends when the next one starts. An entry dated on the next start date belongs to the next period, and the last period has no end.

## Rate sheet code and name are blank when no sheet applies

Lines 32-32

Asked as: rate sheet name missing; blank rate sheet code; no sheet shown on entry

Turbo shows an empty code and name for the rate sheet when an entry's period has no rate sheet or the sheet is not found.

## Employee default rate fills in when the sheet has no rate

Lines 34-35

Asked as: default charge out rate; employee rate used instead; no rate so default applied; fallback hourly rate

When no sheet rate exists, Turbo prices the entry with the employee's default charge-out rate from the admin profile. A default rate of zero or less, or a missing one, is not used.

## Sheet rate is matched by the entry's role

Lines 44-45

Asked as: rate missing for role; wrong role rate; role has no price; rate by job role

Turbo takes the rate from the period's rate sheet for the role on the time entry. When the sheet has no rate for that role, there is no sheet rate.

## Pricing covers only entries inside the report dates, both ends included

Lines 47-48

Asked as: entries outside report dates; time missing from WIP report; report date range excluded; last day of report

Turbo prices only the time entries dated from the start date to the end date of the report. Both dates are included.

## Entries priced with the default rate count as estimated hours

Lines 51-52

Asked as: estimated hours shown; default rate hours estimate; why hours marked estimated

When an entry has no sheet rate but has a default rate, its hours count as estimated hours (estimated_hours). Other entries count zero estimated hours.

## Entries with no rate at all are unpriced hours

Lines 51-53

Asked as: unpriced hours appear; entry has no value; hours with no rate; null value on time

When an entry has neither a sheet rate nor a usable default rate, its hours count as unpriced hours (unpriced_hours) and its value is empty (null).
