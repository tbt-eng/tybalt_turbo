<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_expenses_api.go","sha256":"21702d54d783a8e1e5bd35a3342122fe16cfa9ec8abd6626a34e86428f076f7a","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_expenses_api.go

## Job expenses list requires a job id

Lines 57-60

The job expenses list refuses a request that has no job id. Turbo shows the message "id is required".

## Job expenses list filters

Lines 62-98

The job expenses list accepts filters for division, branch, payment type (payment_type), user (uid) and category. The code passes the values to the job_expenses queries. The docstring omits the branch filter, but the code reads it. The effect of an empty filter depends on the SQL files, which are not in this file.

## Job expenses page number default and fallback

Lines 70-75

The job expenses list shows page 1 by default. A page value that is not a whole number, or that is 0 or negative, falls back to page 1 with no error.

## Job expenses page size default and maximum

Lines 77-85

The job expenses list returns 50 expenses per page by default. A limit above 200 is cut to 200. A limit that is not a whole number, or that is 0 or negative, falls back to 50 with no error.

## Job expenses paging offset

Lines 87-87

The job expenses list skips (page - 1) × limit expenses to reach the requested page.

## Job expenses list response and total pages

Lines 100-124

The job expenses list returns the expenses for the page together with the page, the limit, the total count of matching expenses (total) and the total number of pages (total_pages). Turbo rounds the total pages up, so a partial last page counts as a page. The total count comes from a separate count query.

## Job expenses list query failures

Lines 102-112

If the count query or the expenses query fails, Turbo refuses the request. It shows the message "failed to execute count query: " or "failed to execute query: " followed by the database error.
