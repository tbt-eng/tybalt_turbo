<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_staff_summary_api.go","sha256":"a57199b7a7da99f2a069f7a2c5ac3dccdaf9e6e900b5f53d2469f6e6f0ee4b47","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_staff_summary_api.go

## Job staff summary requires a job id

Lines 28-31

Asked as: job id is required; missing job id; empty job id error; job staff summary no id

The job staff summary request is refused with the message "id is required" when the job id is empty.

## Job staff summary requires start and end dates

Lines 33-38

Asked as: start_date and end_date required; missing date range error; staff summary needs dates; end date missing

The job staff summary request is refused with the message "start_date and end_date are required" when the start date or the end date is missing or empty. Turbo checks only that both are present. It does not check their format or their order here.

## Job staff summary lists hours per staff member for one job

Lines 40-51

Asked as: hours per staff on a job; staff hours for job; who worked on this job; job staff hours list

Turbo returns one row per result of the job staff summary query for the job and the date range. Each row has the job number, the staff member's given name, surname and uid, their hours and their meals hours (meals_hours), plus the job time valuation fields (JobTimeValuation). The result is built from the single job scope query, the job priced time entries query and the job staff summary query. This file does not show how the hours are calculated or priced. If no rows match, Turbo returns an empty list.

## Job staff summary fails when the query fails

Lines 47-49

Asked as: failed to execute query; staff summary query error; job summary database error; staff summary failed

When the summary query fails, Turbo refuses the request with the message "failed to execute query: " followed by the database error text.
