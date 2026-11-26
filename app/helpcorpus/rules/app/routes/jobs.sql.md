<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/jobs.sql","sha256":"bb84838383769b6dcf5140df1599924da7b5783f7494ef38ba2b3c8327b4de71","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/jobs.sql

## Job list includes status, location and imported flag

Lines 6-12

Asked as: job list columns; job status and location; imported job flag; which fields show on jobs

Each job row includes its status, its location and whether it was imported (_imported, shown as imported). It also includes the client id (client_id), the number and the description.

## Jobs without a client still appear in the list

Lines 13-23

Asked as: job missing client; job without client; client name blank; jobs disappearing from list

The jobs query joins clients, branches and managers with LEFT JOINs, so a job stays in the list when any of them is missing. A job with no client has an empty client name (client). The query does not replace that empty value with a default.

## Job list shows blank branch and manager when none is set

Lines 14-15

Asked as: blank branch on job; job with no manager; empty manager name; no branch showing

A job with no branch shows an empty branch (branch code). A job with no manager shows an empty manager name. The manager name is the manager's given name and surname with a space between them.

## Outstanding balance defaults to 0 and its date to blank

Lines 16-17

Asked as: outstanding balance shows zero; balance date empty; job amount owed; no balance date

A job with no outstanding balance (outstanding_balance) shows 0. A job with no outstanding balance date shows an empty value.

## Flag for jobs without an invoicing profile

Lines 18-19

Asked as: no invoicing profile flag; jobs missing invoicing details; procurement invoicing check; invoice profile not set

The no_invoicing_profile flag is true when the job's invoicing profile (invoicing_information) is empty or not set. Procurement uses it to find jobs with no invoicing profile. The query does not filter by status. The comment says the list is of active projects, but the flag itself is only set from the invoicing profile.

## Jobs list returns all jobs or one job by id

Lines 24-24

Asked as: show one job; filter jobs by id; all jobs query; look up job by id

The jobs query returns every job when no job id is given (empty or null). When an id is given, it returns only the job with that id (j.id).

## Jobs are listed newest number first

Lines 25-25

Asked as: job list order; newest jobs first; sort jobs by number; job numbers order wrong

The jobs list is sorted by job number in descending order (ORDER BY j.number DESC). It is a text sort of the number, not a sort by date created.
