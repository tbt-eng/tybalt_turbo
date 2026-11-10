<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/jobs.sql","sha256":"600e691cc47a3dc70ea5aa4f406b2b44061f0cc73793e00ddabbe4d3bbefe11e","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/jobs.sql

## Jobs list shows status, location and imported flag

Lines 9-11

Asked as: job status column; job location missing; imported flag jobs; which jobs were imported

Each job in the list shows its status (status), its location (location) and whether it was imported (imported, from _imported).

## Jobs list includes client name and client id

Lines 12-13

Asked as: client name missing on job; job shows no client; client id on jobs list; job without client

Each job in the list carries its client's name (client) and the client's id (client_id). The client join is a left join, so a job without a matching client still appears, with no client name.

## Missing branch shows as empty text

Lines 14-14

Asked as: branch blank on job; no branch code showing; job branch empty

If a job has no branch, the jobs list shows an empty branch (branch) instead of a branch code.

## Missing manager shows as empty text

Lines 15-15

Asked as: manager blank on job; job has no manager name; who manages this job; manager name empty

The job manager appears as given name and surname joined by a space. If the job has no manager profile, the manager (manager) is an empty string.

## Missing outstanding balance shows as 0

Lines 16-16

Asked as: outstanding amount shows zero; job balance is 0; owed amount blank; unpaid balance missing

If a job has no outstanding balance (outstanding_balance), the jobs list shows 0.

## Missing outstanding balance date shows as empty text

Lines 17-17

Asked as: outstanding date blank; balance due date empty; no date on balance; unpaid date missing

If a job has no outstanding balance date (outstanding_balance_date), the jobs list shows an empty string.

## Jobs list returns all jobs or a single job

Lines 22-22

Asked as: show only one job; list every job; search job by id; filter jobs by id

The jobs query returns every job when no job id (:id) is given or the id is empty. When an id is given, it returns only the job with that id.

## Jobs are sorted by number, highest first

Lines 23-23

Asked as: jobs order wrong; newest job numbers first; job list sorting; latest job at top

The jobs list is ordered by job number (j.number) in descending order, so the highest numbers appear first.
