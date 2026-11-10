<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/jobs.sql","sha256":"600e691cc47a3dc70ea5aa4f406b2b44061f0cc73793e00ddabbe4d3bbefe11e","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/jobs.sql

## Jobs list shows status, location and imported flag

Lines 9-11

Each job in the list shows its status (status), its location (location) and whether it was imported (imported, from _imported).

## Jobs list includes client name and client id

Lines 12-13

Each job in the list carries its client's name (client) and the client's id (client_id). The client join is a left join, so a job without a matching client still appears, with no client name.

## Missing branch shows as empty text

Lines 14-14

If a job has no branch, the jobs list shows an empty branch (branch) instead of a branch code.

## Missing manager shows as empty text

Lines 15-15

The job manager appears as given name and surname joined by a space. If the job has no manager profile, the manager (manager) is an empty string.

## Missing outstanding balance shows as 0

Lines 16-16

If a job has no outstanding balance (outstanding_balance), the jobs list shows 0.

## Missing outstanding balance date shows as empty text

Lines 17-17

If a job has no outstanding balance date (outstanding_balance_date), the jobs list shows an empty string.

## Jobs list returns all jobs or a single job

Lines 22-22

The jobs query returns every job when no job id (:id) is given or the id is empty. When an id is given, it returns only the job with that id.

## Jobs are sorted by number, highest first

Lines 23-23

The jobs list is ordered by job number (j.number) in descending order, so the highest numbers appear first.
