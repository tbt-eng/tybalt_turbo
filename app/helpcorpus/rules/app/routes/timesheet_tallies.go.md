<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/timesheet_tallies.go","sha256":"f200434f62db5d7604a9ad7b3011dd21a67198d2548f936b5811ebf008d0d646","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/timesheet_tallies.go

## Time sheet tallies are listed for the signed-in user by role

Lines 56-69

The time sheet tally handler runs the tallies query (tallies.sql) for the signed-in user (uid) and returns the matching time sheets as a list. It passes the role, the pending-only flag (pendingOnly) and the approved-only flag (approvedOnly) to the query. This file does not show which time sheets those filters select; the query decides that.
