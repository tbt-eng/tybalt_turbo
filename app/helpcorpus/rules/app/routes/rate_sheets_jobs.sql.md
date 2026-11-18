<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/rate_sheets_jobs.sql","sha256":"e22993e09518d6344df2377e22d9ce515974d21f205c58706a570634001d995b","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/rate_sheets_jobs.sql

## Current flag shows whether the job uses the revision now

Lines 2-10

Asked as: job still on this rate sheet; current flag false; job uses revision now; rate sheet not current

Each job in the list has a current flag (current). It is true when the job's rate sheet (jobs.rate_sheet) is this revision. It is false when only a period in the job's rate sheet history uses the revision.

## Jobs without a client still appear with blank client

Lines 8-12

Asked as: job with no client; blank client name in list; missing client on job; client empty in rate sheet list

A job that has no client still appears in the list for a rate sheet revision. Its client ID (client_id) and client name (client_name) are empty text.

## Jobs listed for a rate sheet revision include current and past users

Lines 11-14

Asked as: jobs that used this rate sheet before; past jobs on rate sheet; job missing from rate sheet list; duplicate job in rate sheet list

The list for one rate sheet revision (rate_sheets_jobs) contains every job that uses the revision now (jobs.rate_sheet) and every job whose rate sheet history (job_rate_sheets) has a period that uses it. A job that matches both ways appears once.

## Jobs are sorted by client name then job number

Lines 15-15

Asked as: job list order; sort jobs by client; why are jobs out of order; jobs without client listed first

The list orders jobs by client name, ignoring upper and lower case, and then by job number. Jobs with no client have an empty name and come first.
