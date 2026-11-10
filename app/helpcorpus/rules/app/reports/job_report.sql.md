<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/reports/job_report.sql","sha256":"8e82334d1782473f41ba8aae3726b913012def7b60f67d9bdb97e29975f37e83","model":"claude-sonnet-5-5"} -->

# Rules in app/reports/job_report.sql

## Full Time Report covers one job with no date range

Lines 1-8

The Full Time Report on the Job Details page (job_report.sql) lists time for the single job given by job_id. It applies no date range, so all committed time for that job appears.

## Missing client, job or division show as fallback values

Lines 9-11

If the job has no client name, the report shows the company short name as the client. A missing job number, division code or job description shows as an empty string.

## Report date, month and year come from the entry date text

Lines 13-15

The report splits the entry date into day (characters 9-10), month as a three-letter name (Jan to Dec) and year (characters 1-4). Day and year are numbers.

## Hours count as billable quantity only when the entry has a job

Lines 16-18

In the Full Time Report, an entry with a job puts its hours in qty and 0 in nc. An entry with an empty job puts its hours in nc (non-chargeable) and 0 in qty. The unit is always 'hours'.

## Report shows staff name, meals hours and work record

Lines 19-24

Each row shows the staff member's name as 'given name surname', plus surname and given name separately. It also shows meals hours (meals), the work record (ref) and the entry description. The comments column is always empty.

## Full Time Report only includes committed time sheet entries

Lines 39-43

The Full Time Report takes time entries for the job only from time sheets that have been committed (ts.committed is not empty). Entries on time sheets that are not yet committed do not appear. Entries with no time sheet (tsid empty) are left out.

## Full Time Report adds committed time amendments, marked amended

Lines 44-57

The Full Time Report adds time amendments for the job that have been committed (ta.committed is not empty). Each amendment row is flagged amended = 'true'. Rows from regular time entries have amended empty. Uncommitted amendments are left out.

## Full Time Report rows are sorted in a fixed order

Lines 65-66

The report sorts rows by date (day of month), time type, job, division, qty, nc, surname and given name. The day of month is sorted before month and year, so rows from different months are interleaved by day number.
