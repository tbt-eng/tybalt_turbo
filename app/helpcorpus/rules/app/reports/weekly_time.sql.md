<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/reports/weekly_time.sql","sha256":"59763461432316486ab6db9b526c53a66ae6a58a10444e2e0acb854b2cac94c7","model":"claude-sonnet-5-5"} -->

# Rules in app/reports/weekly_time.sql

## Client shows company short name when entry has no client

Lines 2-2

Asked as: client shows company name; client blank on entry; entry with no job client; client column wrong name

The client column shows the client's name. If the entry has no job or the job has no client, it shows the company short name ({:company_short_name}) instead.

## Missing job or division shows blank in report

Lines 3-4

Asked as: blank job number in report; missing division empty report; no job description showing; empty job code column

If an entry has no job or no division, the job number and division code show as empty text. The project description also shows as empty text when there is no job.

## Report date, month and year come from the entry date text

Lines 6-8

Asked as: wrong date in report; report day month year; month name from date; date columns wrong

The report takes the day, the year and the three-letter month name from the entry date. The day is characters 9-10, the year is characters 1-4 and the month name comes from the month number.

## Entries with a job count as hours, entries with no job as non-chargeable

Lines 9-11

Asked as: hours in wrong column; non-chargeable hours showing; billable vs non-chargeable; hours not charged to job

If a time entry has a job (te.job != ''), its hours go in the qty column and nc is 0. If it has no job, its hours go in nc (non-chargeable) and qty is 0. The unit is always 'hours'.

## Report shows staff name, meals hours and work record

Lines 12-19

Asked as: staff name missing report; meals hours in report; work record missing; comments column empty

Each row shows the staff member's full name (given name, then surname), plus separate surname and given name. It also shows the meals hours (meals), the work record (ref) and the entry description. The comments column is always empty.

## Weekly time report only counts committed time

Lines 32-36

Asked as: uncommitted entries missing; why entries not in report; timesheet not committed report; entries without timesheet

The weekly time report includes time entries only when their time sheet has a committed value (ts.committed != ''). Entries on uncommitted time sheets, or with no time sheet (tsid), are left out.

## Weekly time report filters by a chosen date column and value

Lines 36-51

Asked as: choose date range report; filter report by date; which date column report; report date picker

The report picks entries and amendments by comparing a configurable date column to one date value ({:date_column_entries}, {:date_column_amendments}, {:date_column_value}). The query does not set which column this is. The caller does.

## Weekly time report includes committed time amendments, flagged as amended

Lines 38-52

Asked as: amended entries in report; amendment flagged amended; changed hours showing report; amendments missing weekly report

The weekly time report adds committed time amendments to the time entries. Each amendment row has amended = 'true'. Regular entries have an empty amended value. Uncommitted amendments are left out.

## Weekly time report row order

Lines 60-60

Asked as: report rows out of order; sort weekly time report; row order wrong report; report sorting by date

Rows are sorted by date, time type code, job, division, qty, nc, surname and given name, in that order.
