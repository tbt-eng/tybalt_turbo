<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_request_duplicates.sql","sha256":"3ede5bff07e0025096c9b22ceaee28468be534d8f739a48a830e5ee3739cf1e3","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_request_duplicates.sql

## Which jobs are candidates for the duplicate-work review

Lines 8-19

Asked as: duplicate work review; jobs related to this request; same client jobs; which jobs match my request; related jobs under parent

For a project request, Turbo lists as candidates the jobs that share the request's client, proposal, parent, client PO or client reference number. The request's parent job is also a candidate, as are the other jobs under that parent (j.id = parent or j.parent = parent). The Go code then decides which candidates are listed.

## Proposals are never duplicate-work candidates

Lines 11-11

Asked as: proposals missing from duplicate check; why proposal not listed; P jobs excluded; proposal not showing in duplicates; leave out proposals

Turbo leaves out of the duplicate-work review every job whose number starts with P, which means every proposal. A proposal is left out even when it matches the request's client, parent, client PO or client reference.

## Empty request fields do not match any job

Lines 13-17

Asked as: empty field matching jobs; blank client PO no match; request with no client reference; empty fields in duplicate search; missing PO matches everything

Turbo checks a field of the request (client, proposal, parent, client PO, client reference) only when the field is not empty. A request with an empty field never matches jobs on that field.

## Client PO and client reference matching ignores case and spaces

Lines 16-17

Asked as: case sensitive PO match; client reference spaces not matching; PO number capitals mismatch; trailing space reference not found; lowercase matching duplicate

Turbo compares the client PO and the client reference number after trimming spaces and lowercasing the job's value. The request's value is compared as it is passed in, so it is assumed to be already lowercase and trimmed. Client and proposal matches are exact.

## Latest time entry date for each candidate job

Lines 20-25

Asked as: last time entry date; when was job last worked on; latest timesheet date; last logged hours for job; job with no time entries

For each candidate job, Turbo shows the date of its most recent time entry (last_time). A job with no time entries shows an empty last_time.

## Names shown for each candidate job

Lines 39-41

Asked as: client name blank in duplicates; manager name missing; parent job number empty; who manages this candidate job; manager name fallback user name

Each candidate row shows the client name, the parent job's number and the manager's name. They are empty when the job has no client, parent or manager. The manager's name is the manager's profile given name and surname joined by a space. If both are empty, Turbo falls back to the user's name, and then to empty.
