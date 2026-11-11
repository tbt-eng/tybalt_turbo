<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_priced_time_entries.sql","sha256":"58f04f789d0adfdb9afe0800ae67a31d9f25291296afa4024581bf6dc504df0a","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_priced_time_entries.sql

## Time entry priced with rate sheet in effect on its date

Lines 1-3

Asked as: wrong rate on old timesheet; rate changed after the date; which rate sheet applies to a date; time priced with the wrong rates; rate sheet history by date

Turbo prices each time entry with the rate sheet that was in effect on the entry's date. It takes the sheet from the project's rate sheet history (job_rate_sheets). This applies to WIP and to the staff and division summaries.

## Project without rate sheet history uses its own rate sheet

Lines 3-5

Asked as: project has no rate history; imported project wrong rate; job rate sheet not used; no rate sheet history

A project with no history rows uses the job's rate sheet (jobs.rate_sheet) for all of its time. Imported projects have no history until it is written. An empty rate sheet value on the job is treated as no sheet (null).

## Inactive rate sheet revisions still price past time

Lines 5-6

Asked as: inactive rate revision; rate change not updating past time; old revision still pricing hours; rates changed but value unchanged

Turbo uses the revision assigned to the rate sheet even when that revision is inactive. A change to a sheet's rates therefore recalculates the value of past time.

## Which time entries count toward value

Lines 6-8

Asked as: why amendments not in value; meal entries excluded from value; what counts toward job value; uncommitted hours in value

Value keeps the existing selection of committed and uncommitted entries. Amendments and meals are not part of the value.

## Rate sheet period boundaries and opening period

Lines 12-29

Asked as: entries before first effective date; rate changes on effective date; which rate applies on start day; time before rate sheet start; early entries no rate; period boundary rate sheet

Each rate sheet period runs from its effective date up to, but not including, the next period's effective date. The last period has no end. An entry on a period's start date uses that period, and an entry on the next effective date uses the next sheet. Every selected job gets an opening period with an empty start date that is before every date. It is added when the job has no history row with an empty effective date, so entries before the first effective date are still kept. Such early entries use no sheet when the job has history rows.

## Rate comes from the rate sheet by role, else employee default

Lines 33-51

Asked as: rate for role missing; fallback to employee default rate; role has no rate; charge out rate used for hours

The rate for an entry is the rate sheet's rate for the entry's role. If the sheet has no rate for that role, Turbo uses the employee's default charge-out rate from the admin profile (default_charge_out_rate). The value is hours times that rate.

## Zero or missing default rate is not used

Lines 34-35

Asked as: zero default rate; employee has no charge out rate; time shows no dollar value; missing default rate value

An employee's default charge-out rate counts only when it is above 0. A zero or missing default cannot supply a dollar value.

## Entries limited to the report date range, inclusive

Lines 47-48

Asked as: report date range included; start date end date report; entries missing from report dates; report includes last day

Only time entries dated from the start date through the end date count. Both dates are included.

## Estimated hours use the default rate

Lines 52-52

Asked as: estimated hours column; why hours marked estimated; estimate instead of rate; hours with no role rate

An entry counts as estimated hours when the rate sheet has no rate for its role but the employee has a usable default rate. Its estimated hours equal its hours. Otherwise estimated hours are 0.

## Unpriced hours when no rate is available

Lines 53-53

Asked as: unpriced hours shown; hours with no rate; time with zero value; no rate for hours

An entry counts as unpriced hours when neither the rate sheet nor the employee default gives a rate. Its unpriced hours equal its hours. Otherwise unpriced hours are 0.
