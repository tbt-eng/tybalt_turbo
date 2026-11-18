<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/client_rate_sheets.sql","sha256":"8d806433ac843eb74a2069cc6c8914cbb60c86ad9abf11bde271d437abb129be","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/client_rate_sheets.sql

## Current jobs and past jobs counts for each rate sheet

Lines 15-16

Asked as: current jobs count; past jobs count; job counted twice; rate sheet job numbers

current_jobs counts the client's jobs that use the revision now. past_jobs counts the client's jobs that only a period of their rate sheet history puts on the revision. Turbo counts each job once, and a job on the revision now is never counted in past_jobs.

## Client rate sheets list: which revisions appear

Lines 17-27

Asked as: which rate sheets show for client; client rate sheet list missing; job owner not client; rate sheets listed for a client

For a client, Turbo lists each rate sheet revision that one of the client's jobs uses now (jobs.rate_sheet) or that a period in a job's rate sheet history (job_rate_sheets) names. Turbo counts only jobs whose client (jobs.client) is this client. A job where the client is only the job owner does not count.

## Jobs with no current rate sheet are not counted as current

Lines 18-20

Asked as: job with no rate sheet; empty rate sheet job; job missing from current; job without rate sheet counted

A job counts as current only when its rate sheet (jobs.rate_sheet) is not empty. A job with an empty rate sheet appears only through its rate sheet history, as a past job.

## Proposed rates from a project request are left out

Lines 28-28

Asked as: proposed rates hidden; project request rate sheet; rates missing from client list; proposal rates not showing

Turbo leaves out any rate sheet that belongs to a project request (rate_sheets.request is not empty). Proposed rates therefore do not appear in the client's rate sheet list. This matches the read rule on rate_sheets.

## Client rate sheets sort order

Lines 30-30

Asked as: rate sheet list order; sort rate sheets newest; why rate sheets in this order; current rate sheets first

Turbo lists the rate sheets that have current jobs first. Within each group it puts the newest code first, then the newest revision.
