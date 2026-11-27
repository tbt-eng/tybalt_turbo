<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/jobs_stale.sql","sha256":"1ff7cbd32007c1bb37fd932debe81846ad3bf9191fe7ccdf489e169f1c26fa8d","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/jobs_stale.sql

## Age of 0 or less lists unused jobs

Lines 3-4

Asked as: jobs with no time entries; unused jobs list; jobs never used; empty last reference; jobs with no activity

When the age (:age) is 0 or less, the report lists jobs with no time entries, time amendments, expenses or purchase orders at all. Such jobs show an empty last reference and an empty last reference type.

## Stale and unused job report only lists Active jobs

Lines 9-15

Asked as: inactive jobs missing from report; closed jobs not showing; stale report missing jobs; only active jobs shown; job status filter

The stale/unused jobs report only considers jobs whose status is 'Active'. Jobs in any other status never appear in it.

## Job number prefix filter matches the start of the number

Lines 13-15

Asked as: job number prefix search; jobs starting with prefix; job number filter not working; search job number start; prefix matches middle of number

The report only includes jobs whose number starts with the prefix given (j.number LIKE prefix || '%'). The prefix is not matched anywhere else in the number.

## What counts as a reference to a job

Lines 16-48

Asked as: what is last reference; same date reference type; which entry counts as last use; job last activity date; time entry vs expense priority

Turbo looks at the date of every time entry, time amendment, expense and purchase order that points to the job. The most recent date is the job's last reference. When two references have the same date, the type ranks in this order: time entry, time amendment, expense, purchase order. The first in that order is reported as the last reference type.

## Report fills in blank values for missing data

Lines 56-65

Asked as: blank client alias in report; missing branch code empty; zero outstanding balance missing; job missing client not shown; empty fields in stale report

The report shows an empty string for a missing client alias, branch code, outstanding balance date, last reference and last reference type. It shows 0 for a missing outstanding balance. A job still appears if its client or branch cannot be found.

## Positive age lists stale jobs, with the boundary day included

Lines 66-67

Asked as: how many days stale; older than days threshold; jobs not used for days; boundary day included stale; stale age cutoff

When the age (:age) is above 0, the report lists jobs that have at least one reference and whose last reference date is on or before today minus that many days (<=). A job whose last reference is exactly that many days old is included. Jobs with no references at all are not listed as stale.

## Report sort order and row limit

Lines 68-69

Asked as: stale report sort order; stale jobs oldest first; report row limit; how many rows returned; unused jobs sorted by number

For stale jobs, the report sorts by last reference date, oldest first, then by job number. For unused jobs, it sorts by job number only. It returns at most the number of rows given by the limit parameter.
