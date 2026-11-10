<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/work_records_details.sql","sha256":"669188aaac7d603668be92b8dcdf193df6afd77b541f231ea4806beeaaca23c4","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/work_records_details.sql

## Missing hours show as 0

Lines 6-6

Turbo reports hours (hours) as a number. When an entry has no hours, it shows 0.

## Job number is blank when entry has no matching job

Lines 7-8

Turbo shows the job number (job_number) from the linked job. It shows an empty string when the entry has no job or the job is not found. The job ID (job_id) is also an empty string when the entry has no job.

## Missing description shows as empty text

Lines 9-9

When a time entry has no description, Turbo shows an empty string.

## Staff name is blank when profile is missing

Lines 10-11

Turbo takes the surname and given name from the profile of the entry's user (uid). When no profile matches, both show as empty strings.

## Time sheet is found by the entry's time sheet ID

Lines 12-19

When a time entry has a time sheet ID (tsid), Turbo uses the time sheet with that ID as the entry's time sheet (timesheet_id).

## Entry without time sheet ID falls back to same user and week

Lines 12-23

When a time entry has an empty time sheet ID (tsid), Turbo uses the time sheet with the same user (uid) and week ending (week_ending). If no time sheet matches, timesheet_id is an empty string. If several time sheets match, the entry can appear once per match.

## Work record details list time entries for one work record

Lines 13-26

The query returns the time entries (time_entries) whose work record (work_record) equals the one requested. Turbo sorts them by week ending (week_ending), newest first.

## Empty work record returns no time entries

Lines 24-25

The query excludes time entries whose work record is an empty string. A request with an empty work record returns no rows.
