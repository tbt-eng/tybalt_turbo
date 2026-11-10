<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/active_jobs_report.go","sha256":"08001acf5fbb1da330d87edbd5a992804b9ef5bc2422f03a54a6cd319e121dc1","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/active_jobs_report.go

## Only staff with the report claim can view the active jobs report

Lines 22-46

Turbo checks that the user has the report claim before it builds the active jobs report. A user without it gets the message "you are not authorized to view this report".

## Active jobs report columns and blank values

Lines 50-58

Each row of the active jobs report has five columns: status, job_number, manager_name, branch_name and time_and_materials. The manager name is the manager's given name and surname joined by a space. When a job has no manager or no branch, the manager name or branch name is blank. A job with no time and materials value shows as false.

## Active jobs report lists only Active jobs, not proposals

Lines 59-60

The active jobs report includes only jobs with status 'Active'. It leaves out every job whose number starts with P, because those are proposals.

## Active jobs report sorts by job number, highest first

Lines 61-61

The active jobs report sorts its rows by job number in descending order (ORDER BY j.number DESC).

## Active jobs report downloads as active_jobs.csv

Lines 69-99

Turbo returns the active jobs report as a CSV file named active_jobs.csv, with a header row first. The time_and_materials column holds true or false.
