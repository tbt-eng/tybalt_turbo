<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/client_time_divisions.sql","sha256":"ac98f3fdb549c8c946b7108f98d38e8274460493027803cd6dc0c7908be0dd39","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/client_time_divisions.sql

## Client time divisions use jobs where the client is the client

Lines 1-15

The client's division report includes time on jobs whose client (j.client) is this client. Jobs where the client is only the job owner do not count.

## Client time divisions report hours per division in code order

Lines 6-19

The client's division report gives one row per division (id, code, name) with the sum of hours. Rows are sorted by division code.

## Client time divisions count only committed time sheets

Lines 13-17

The client's division report adds up hours only from time entries on time sheets that are committed (committed is not empty). Time on time sheets that are not committed is left out.

## Zero-hour entries are excluded from client division hours

Lines 16-16

The client's division report counts only time entries with more than 0 hours (te.hours > 0). Entries with 0 hours, or negative hours, are left out. A division whose entries all have 0 hours does not appear.
