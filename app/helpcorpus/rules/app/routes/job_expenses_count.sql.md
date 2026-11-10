<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_expenses_count.sql","sha256":"0e749bdd69d660240c3a2dabf6bb523a55406d3a21ff7c8cae2418828306b2f1","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_expenses_count.sql

## Job expense count includes only committed expenses

Lines 2-5

The count of a job's expenses counts only expenses that have been committed (committed is not empty). Uncommitted expenses are not counted. It counts only expenses for the requested job (e.job).

## Job expense count optional filters ignore blank values

Lines 6-9

The count can be narrowed by branch, division, payment type and staff member (uid). A filter that is missing or an empty string is ignored, so it does not narrow the count. A filter with a value counts only expenses that match it exactly.

## Category filter 'none' counts uncategorized expenses

Lines 10-16

The category filter works in three ways. A missing or empty category does not narrow the count. The special value 'none' counts only expenses with no category (category is empty). Any other value counts only expenses in that exact category.
