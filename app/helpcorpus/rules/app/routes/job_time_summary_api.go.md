<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_time_summary_api.go","sha256":"86e0d0547a16117edfd3dfb622f5124c06111de53c193eae09d76eff99cf42b8","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_time_summary_api.go

## Job time summary requires a job id

Lines 37-40

Asked as: id is required; job summary missing id; no job id given; summary request without job

The job time summary refuses a request that has no job id. Turbo shows the message "id is required".

## Job time summary filters

Lines 42-57

Asked as: filter job time summary; narrow summary by branch; summary by division or user; filter hours by time type

The job time summary can be narrowed by branch, division, time type, user (uid) and category. Turbo passes each filter to the job_time_summary.sql query along with the job id. The filters are optional. A filter that is left out is passed as an empty string, and this file does not show how the query treats an empty value.

## Job time summary returns empty result when no rows

Lines 58-60

Asked as: summary shows nothing; empty job summary; no rows in summary; summary blank no zeros

When the summary query returns no rows, Turbo returns an empty result with no fields, not zeros.

## Job time summary query failure message

Lines 61-61

Asked as: failed to execute query; summary query error; database error on summary; job summary crashed

When the summary query fails for any other reason, Turbo shows the message "failed to execute query: " followed by the database error.

## Summary text fields default to empty string

Lines 64-89

Asked as: blank summary fields; null entry in summary; no latest entry shown; empty branches or names

The job time summary returns earliest_entry, latest_entry, branches, divisions, time_types, names and categories as text. Any of these that is empty comes back as an empty string, not as null.

## Total hours falls back to zero

Lines 72-78

Asked as: total hours wrong; hours shows zero; total hours not a number; summary total missing

The job time summary reports total_hours as a number. It reports 0 when the total is empty or cannot be read as a number.
