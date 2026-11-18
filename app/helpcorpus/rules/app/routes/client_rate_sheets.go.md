<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/client_rate_sheets.go","sha256":"e9410ea867f921a926a3325b390dcfd38fd4b738c4632d13ce02aaa8fcbd02a0","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/client_rate_sheets.go

## Client rate sheet list shows revisions with current and past job counts

Lines 14-44

Asked as: rate sheet revisions for a client; current and past jobs count; which rate sheet version a client uses; client rate sheet history

For a client, Turbo lists each rate sheet revision with its ID, code, name, revision number, effective date and whether it is active (active). It also gives current_jobs, the client's jobs on that revision now. It gives past_jobs, the client's jobs that only a period of their rate sheet history put on that revision. The exact selection and counting are in client_rate_sheets.sql, which this file does not show.

## Client rate sheets list is open to any signed-in user

Lines 30-35

Asked as: who can see client rate sheets; non procurement access to rates; any user view client pricing; permission to read rate sheets

Any signed-in user can read the list of rate sheet revisions that a client's jobs use, current and past. It is not limited to procurement. Turbo allows it because any signed-in user can already list rate sheets, jobs and clients.

## Unknown client gives 'client not found'

Lines 36-39

Asked as: client not found error; wrong client ID rate sheets; can't find client rate sheets; client lookup fails

When the client ID does not match a client, Turbo refuses the request to list the client's rate sheets. It shows the message 'client not found'.

## Client with no rate sheets returns an empty list

Lines 40-44

Asked as: client has no rate sheets; empty rate sheet list; client rate sheets blank; no revisions shown for client

When the client exists but the query finds no revisions, Turbo returns an empty list and not an error.
