<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/time_entry_branch_mismatches.go","sha256":"74cb676a2f1ac53cebc805fd6ad329f919d2bdfa9d9d0efef3ba3ff0fbab30a1","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/time_entry_branch_mismatches.go

## Only report viewers can run the branch mismatch report

Lines 58-60

Asked as: cannot run branch report; branch mismatch access denied; who can view reports; report permission error; no csv for branch report

Turbo checks that the user may view reports (requireReportViewer) before it runs the time entry branch mismatch report. A user who fails the check gets the error that the check returns, and no CSV is produced. The file does not show who counts as a report viewer.

## Branch mismatch report lists only rows where both branches are set and differ

Lines 62-90

Asked as: branch report missing entries; blank branch not shown; different branches only; why entry not in mismatch report; empty branch hidden report

The report lists time entries that have a job, where the time entry's branch and the job's branch are both non-blank and are different. It does not list entries with a blank time entry branch. It does not list entries whose job has a blank branch. Blank-versus-filled differences are therefore invisible in this report.

## Blank values appear as empty text or zero in the report

Lines 65-79

Asked as: zero hours in report; blank cells csv; employee name empty report; missing job number report; empty fields branch export

Missing hours show as 0. Missing job number, job description, branch id, code, name and time entry description show as empty text. The employee name is the given name and surname joined by a space and trimmed. It is blank if the profile is missing.

## Branch mismatch report skips entries whose job record is missing

Lines 80-81

Asked as: entry missing from report; deleted job branch report; time entry without job; job not found report; orphan entries branch report

The report joins each time entry to its job with an inner join. A time entry whose job cannot be found does not appear in the report.

## Branch mismatch report sort order

Lines 89-89

Asked as: branch report row order; report sorted how; newest entries first report; mismatch report ordering; which date first csv

The report sorts rows by date with the newest first. Ties go by employee name A to Z, then job number, then time entry id.

## Branch mismatch report downloads as a CSV file

Lines 97-149

Asked as: branch report download; export branch mismatches; time_entry_branch_mismatches file; branch csv columns; download mismatch spreadsheet

The report downloads as the file time_entry_branch_mismatches.csv. It has a header row and one row per mismatched time entry. The columns are time_entry_id, date, week_ending, hours, uid, employee_name, job_id, job_number, job_description, the time entry branch id, code and name, the job branch id, code and name, and description. Hours print without trailing zeros.
