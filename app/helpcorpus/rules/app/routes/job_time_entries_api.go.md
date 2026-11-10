<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_time_entries_api.go","sha256":"8f01b2fe89bd775f4ee8e1703a18e975a94df44182ad609bc976b2891488369a","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_time_entries_api.go

## Job time entries list requires a job id

Lines 58-61

Turbo refuses a request for a job's time entries when the job id is missing. The message is "id is required".

## Job time entries can be filtered by five fields

Lines 63-99

The job time entries list accepts the filters division, branch, time type (time_type), user (uid) and category. The code passes each filter to the query as given. An empty filter is passed as an empty string, and this file does not show how the query treats it. The header comment omits the branch filter, but the code reads it.

## Job time entries page number defaults to 1

Lines 71-76

The job time entries list starts on page 1 when no page is given. Turbo also uses page 1 when the page is not a whole number or is 0 or negative.

## Job time entries page size defaults to 50, maximum 200

Lines 78-86

The job time entries list shows 50 entries per page by default. A limit above 200 becomes 200. A limit that is not a whole number, or is 0 or negative, falls back to 50.

## Job time entries page offset calculation

Lines 88-88

Turbo skips (page - 1) × limit entries to reach the requested page.

## Job time entries response includes pagination totals

Lines 101-125

The response gives the entries together with the page, the limit, the total count of matching entries (total) and the total number of pages (total_pages). Turbo rounds the total pages up. A page beyond the last returns no entries, and the page number stays as requested.

## Job time entries list fails if a query fails

Lines 102-113

If the count query or the entries query fails, Turbo returns an error. The message is "failed to execute count query: " or "failed to execute query: ", followed by the database error.
