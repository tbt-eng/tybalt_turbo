<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/jobs.sql","sha256":"14098862edeecee0c4a2e0d798399836db8b2181cc42c3d1ceb8780f19e30d64","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/jobs.sql

## Job list includes status, location and imported flag

Lines 9-11

Asked as: job status and location; job list missing status; imported job flag; which jobs were imported; job location column

Each job row includes the job's status, its location, and whether the job was imported (_imported, returned as imported). The query returns all jobs regardless of status.

## Job list shows client name, alias and client id

Lines 12-14

Asked as: client alias missing; client name on job list; job without client; client id on jobs; client record not showing

Each job row shows its client's name (client) and the client's record id (client_id). The client alias (client_alias) is an empty string when the client has none. A job without a matching client record still appears, because the client is joined with a left join.

## Job branch and manager default to empty text

Lines 15-16

Asked as: job has no manager; blank branch on job; manager name empty; job missing branch code; no manager shown

The job's branch shows the branch code (branch), and the manager shows the manager's given name and surname joined by a space. Each is an empty string when the job has no branch or no manager, or when no matching branch or profile exists. The job is still listed.

## Outstanding balance defaults to 0 and empty date

Lines 17-18

Asked as: outstanding balance shows zero; empty balance date; job balance blank; outstanding amount missing; balance date not set

A job without an outstanding balance (outstanding_balance) shows 0. A job without an outstanding balance date (outstanding_balance_date) shows an empty string.

## No-invoicing-profile flag is true for empty invoicing information

Lines 19-20

Asked as: no invoicing profile flag; missing invoicing information; procurement active projects without invoicing; invoicing profile missing; job lacks invoicing info

The no_invoicing_profile flag is true when the job's invoicing information (invoicing_information) is empty or missing. Otherwise it is false. The query does not filter by status. It sets the flag for every job, and the flag is not limited to active projects. Procurement uses the flag to find active projects that lack an invoicing profile.

## Jobs list returns all jobs or one job by id

Lines 25-26

Asked as: show all jobs; find job by id; jobs sorted by number; single job lookup; job list order

The jobs query returns every job when no id is given (null or empty). When an id is given, it returns only the job with that id. Jobs are ordered by job number, highest first.
