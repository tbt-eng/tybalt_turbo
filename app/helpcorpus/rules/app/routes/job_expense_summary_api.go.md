<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_expense_summary_api.go","sha256":"359122fbdbf65a46133bc88a41fd82529f8464da05e3c162d27294c81fe0d28a","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_expense_summary_api.go

## Job expense summary requires a job id

Lines 37-40

The job expense summary refuses a request that has no job id. Turbo shows the message "id is required".

## Job expense summary optional filters

Lines 42-57

The job expense summary accepts optional filters for branch, division, payment type (payment_type), user (uid) and category. Turbo passes each filter to the summary query (job_expense_summary.sql) along with the job id. The file does not show how the query treats a filter that is left empty.

## Job expense summary is empty when no rows return

Lines 58-60

When the summary query returns no row, Turbo answers with an empty result and no error. The response has no totals or lists.

## Job expense summary query failure message

Lines 61-61

When the summary query fails for any other reason, Turbo refuses with the message "failed to execute query: " followed by the database error.

## Job expense summary returns blank text for missing values

Lines 64-89

The summary returns earliest expense, latest expense, branches, divisions, payment types, names and categories as text. When the query gives no value for one of them, Turbo returns an empty string.

## Job expense summary total falls back to zero

Lines 72-81

The summary reports total_amount as a number. When the total is missing or cannot be read as a number, Turbo reports 0.
