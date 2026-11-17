<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_staff_summary.sql","sha256":"c39081581436a6236c43487ba753e547b2b90b745ca14d4604cf9b3ec9192a25","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_staff_summary.sql

## Job staff summary prices entries before grouping by person

Lines 1-3

Asked as: same person two roles; one row per person; person works two roles job; entries priced separately by role

The job staff summary prices each time entry first and then groups the entries by person (uid). A person who works in more than one role on the same job therefore gets one row, and each entry is priced at its own role.

## Rate sheets column lists each sheet used, with revision and name

Lines 4-16

Asked as: which rate sheets used; rate sheet list blank; rate sheet revision name; multiple rate sheets per person

For each person, the summary lists the distinct rate sheets that priced their entries, separated by '; ' and sorted by label. A label shows the code and the revision, as in '2026-11-00004 (rev. 1)'. When the sheet has a name, the name follows in quotes. Entries with no rate sheet code add no label, and a person with no sheet gets a blank list (rate_sheets).

## Staff without a profile show blank names

Lines 15-15

Asked as: missing name in summary; no profile name blank; staff row no name; person without profile

When a person has no profile, the job staff summary shows blank given and family names (given_name, surname) for that person. The row stays in the summary.

## Person's value is the sum of priced entries, zero if none

Lines 17-24

Asked as: person total zero; sum of hours per person; staff value adds up; hours totals per staff

A person's value is the sum of their entries' values. It is 0 when no entry has a value. Hours, meal hours, estimated hours and unpriced hours are summed per person in the same way.

## Percent of job total is blank when any hours are unpriced

Lines 19-35

Asked as: percent of job blank; no percentage shown; percent empty unpriced hours; job share missing

A person's percent is their value divided by the job total, times 100, rounded to 1 decimal place. Turbo leaves the percent blank when any hours on the job are unpriced. It also leaves it blank when the total value is 0, because the division by zero gives no result.

## Row status shows Incomplete, Estimated or Rate sheet

Lines 31-32

Asked as: row says incomplete; estimated status staff; why rate sheet status; unpriced hours status

A person's value status (value_status) is 'Incomplete' when any of their hours are unpriced. Otherwise it is 'Estimated' when any of their hours are estimated. Otherwise it is 'Rate sheet'. Unpriced hours take precedence over estimated hours.

## Total status uses the same rule across all staff

Lines 33-34

Asked as: total status incomplete; job total status wrong; overall job status; whole job estimated flag

The total status (total_status) is 'Incomplete' when any person on the job has unpriced hours. Otherwise it is 'Estimated' when any person has estimated hours. Otherwise it is 'Rate sheet'.

## Staff summary is sorted by surname, given name, then ID

Lines 37-37

Asked as: staff list order; sort staff by surname; blank names listed first; summary sorting people

The job staff summary lists people by surname, then given name, then uid. People without a profile have blank names, so they sort first.
