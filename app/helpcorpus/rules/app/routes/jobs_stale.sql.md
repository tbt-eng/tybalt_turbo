<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/jobs_stale.sql","sha256":"c656d71c1a71f0944242ddaa263ed0066960b40bc313374c31bcff437e06dd8d","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/jobs_stale.sql

## Stale and unused job report only covers Active jobs

Lines 13-13

The stale/unused jobs query only looks at jobs whose status is 'Active'. Jobs with any other status never appear in the results.

## Job number prefix filter matches the start of the number

Lines 14-14

The report returns only jobs whose number starts with the prefix given (j.number LIKE prefix || '%'). The match is on the start of the number, not on text elsewhere in it.

## What counts as a reference to a job

Lines 16-35

Turbo counts a job as referenced when a time entry, a time amendment, an expense or a purchase order points to it. It uses the date of each such record as the reference date.

## Last reference is the most recent date; ties go to the first type in the order

Lines 36-48

A job's last reference is the record with the latest date. If several records share that date, Turbo picks by type order: time entry, then time amendment, then expense, then purchase order. The report shows that record's date and type (last_reference, last_reference_type).

## Blank values shown when a job lacks data

Lines 56-60

The report shows an empty branch if the job has no branch. It shows an outstanding balance of 0 if none is set. It shows an empty outstanding balance date, last reference and last reference type if none exist.

## Age of 0 or less lists unused jobs

Lines 65-65

When the age (:age) is 0 or less, the report lists only jobs with no reference at all: no time entries, time amendments, expenses or purchase orders.

## Age above 0 lists stale jobs; the cutoff day is included

Lines 66-66

When the age (:age) is above 0, the report lists only jobs that have at least one reference and whose last reference is on or before today minus that many days. A last reference exactly that many days old counts as stale (<=). Jobs with no references at all are left out.

## Report sort order and row limit

Lines 67-68

For stale jobs (age above 0), the report sorts by last reference date, oldest first, then by job number. For unused jobs it sorts by job number alone. It returns at most the limit (:limit) number of rows.
