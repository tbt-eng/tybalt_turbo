<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/timesheet_approver.go","sha256":"c301e23cdd6ec776cf4e33af226f567613112689808b736bd458421d5258ebd7","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/timesheet_approver.go

## Time sheet approver info shows approver, committer and rejector

Lines 19-33

Asked as: who approved my timesheet; timesheet approver name; who rejected timesheet; approval date on timesheet

Turbo returns the approver's name and approval date (approver_name, approved_date) for a time sheet. It also returns the committer's name and commit date, and the rejector's name and rejection date. Each name is the person's given name and surname joined by a space.

## Missing approver, committer or rejector shows an empty name

Lines 21-31

Asked as: blank approver on timesheet; timesheet approver shows empty; no name for committer; timesheet not approved blank fields

When a time sheet has no approver, committer or rejector, Turbo shows an empty name for that role. This happens because the profile lookup falls back to an empty string (COALESCE). A time sheet that has not been approved, committed or rejected therefore shows blank fields for that step.

## Unknown time sheet returns blank approver info, not an error

Lines 35-50

Asked as: timesheet approver info empty; timesheet not found no error; approval history missing timesheet; blank timesheet details no error

When the approver info lookup fails, for example because no time sheet has that id, Turbo returns empty names and dates for all six fields. It does not report an error. Staff cannot tell a missing time sheet from one that has no approval history.

## Approver info request needs a time sheet id

Lines 56-59

Asked as: time sheet id is required; approver info missing id; timesheet id error message; approver lookup no id

When the request has no time sheet id, Turbo refuses it with the message "time sheet id is required".
