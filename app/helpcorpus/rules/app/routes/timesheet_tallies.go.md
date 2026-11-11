<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/timesheet_tallies.go","sha256":"6322905390c9954565419d3be1e8394e48fd0762c6f7c7bffa34ca2246ddcffc","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/timesheet_tallies.go

## Each time sheet tally lists hours, status, approver and committer details

Lines 15-53

Asked as: time sheet tally columns; what does a tally row show; timesheet hours by type; who approved my timesheet; rejected timesheet reason

Each tally row describes one time sheet. It has the week ending, the person's name, and hours by type: work, job, week, meals, non-work, and the OB, OP, OS and OV hours. It also has payout request amount and dates, bank entry dates, off-rotation and off-week dates, divisions, job numbers and time types. Status fields show submitted, approved, rejected with the reason, and committed. Names of the approver, committer and rejector are included.

## Time sheet tallies depend on the signed-in user, role and status filters

Lines 57-65

Asked as: which timesheets are listed; pending only timesheets filter; approved only timesheets; tally list missing my timesheets; timesheet list depends on role

The time sheet tally handler runs the embedded tallies query (tallies.sql) for the signed-in user (uid), a role, and the pendingOnly and approvedOnly flags. The route that creates the handler sets the role and the flags. This file does not show which time sheets the query returns for them.

## Tally request fails with an error message when the query fails

Lines 66-70

Asked as: failed to execute query error; timesheet tally error message; tally request failed; no tallies returned error; database error timesheet list

When the tally query fails, Turbo returns the error "failed to execute query: " followed by the database error. It returns no tallies in that case. On success it returns the list of tallies as JSON.
