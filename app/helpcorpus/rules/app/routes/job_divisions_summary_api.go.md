<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_divisions_summary_api.go","sha256":"b13d112e007d06f1eafe1d656479fcbc663975721d2506675934d8e637cf8b25","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_divisions_summary_api.go

## Job divisions summary content and scope

Lines 14-49

The job divisions summary covers one job and the date range given. It returns one row per division with the job number, division code, division name and hours. Each row also carries the time valuation fields of JobTimeValuation. The figures come from the single-job scope, the priced time entries and the divisions summary queries, which are defined in other files.

## Job divisions summary requires a job id

Lines 26-29

The job divisions summary refuses a request that has no job id (id). Turbo shows the message "id is required".

## Job divisions summary requires start and end dates

Lines 31-36

The job divisions summary refuses a request that lacks a start date (start_date) or an end date (end_date). Turbo shows the message "start_date and end_date are required".

## Job divisions summary fails when the query fails

Lines 45-47

When the summary query fails, Turbo shows the message "failed to execute query: " followed by the database error.
