<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_expenses_api.go","sha256":"21702d54d783a8e1e5bd35a3342122fe16cfa9ec8abd6626a34e86428f076f7a","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_expenses_api.go

## Job expenses list requires a job id

Lines 57-60

Asked as: id is required; job expenses no job id; missing job id error; expenses list without job

The job expenses list refuses a request that has no job id. Turbo shows the message "id is required".

## Job expenses list filters

Lines 62-98

Asked as: filter job expenses by division; filter expenses by branch; job expenses payment type filter; filter expenses by user; filter job expenses by category

The job expenses list accepts filters for division, branch, payment type (payment_type), user (uid) and category. The code passes the values to the job_expenses queries. The docstring omits the branch filter, but the code reads it. The effect of an empty filter depends on the SQL files, which are not in this file.

## Job expenses page number default and fallback

Lines 70-75

Asked as: expenses page number wrong; job expenses wrong page shown; page zero expenses list; expenses page not a number

The job expenses list shows page 1 by default. A page value that is not a whole number, or that is 0 or negative, falls back to page 1 with no error.

## Job expenses page size default and maximum

Lines 77-85

Asked as: how many expenses per page; expenses list page size; max expenses per page; limit job expenses 200; expenses limit not working

The job expenses list returns 50 expenses per page by default. A limit above 200 is cut to 200. A limit that is not a whole number, or that is 0 or negative, falls back to 50 with no error.

## Job expenses paging offset

Lines 87-87

Asked as: job expenses skip rows; expenses page offset; second page expenses wrong rows; expenses paging skips

The job expenses list skips (page - 1) × limit expenses to reach the requested page.

## Job expenses list response and total pages

Lines 100-124

Asked as: total pages expenses; expenses total count; job expenses how many pages; partial last page expenses; expenses list total number

The job expenses list returns the expenses for the page together with the page, the limit, the total count of matching expenses (total) and the total number of pages (total_pages). Turbo rounds the total pages up, so a partial last page counts as a page. The total count comes from a separate count query.

## Job expenses list query failures

Lines 102-112

Asked as: failed to execute query; expenses count query failed; job expenses database error; expenses list request refused; failed to execute count query

If the count query or the expenses query fails, Turbo refuses the request. It shows the message "failed to execute count query: " or "failed to execute query: " followed by the database error.
