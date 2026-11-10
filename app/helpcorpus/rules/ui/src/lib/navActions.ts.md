<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/navActions.ts","sha256":"983f2622f5d741823f3ffc32d1a4397ad93b14bad9bcb19f76173535cf5d6de1","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/navActions.ts

## Time entry branch mismatches report downloads as CSV

Lines 4-9

The time entry branch mismatches download fetches the report from the /api/reports/time_entry_branch_mismatches endpoint. Turbo saves it as "time_entry_branch_mismatches.csv" (downloadTimeEntryBranchMismatchesCsv). This file does not show who may run the report or what it contains.

## Active jobs report downloads as CSV

Lines 11-13

The active jobs download fetches the report from the /api/reports/active_jobs endpoint. Turbo saves it as "active_jobs.csv" (downloadActiveJobsCsv). This file does not show who may run the report or which jobs count as active.
