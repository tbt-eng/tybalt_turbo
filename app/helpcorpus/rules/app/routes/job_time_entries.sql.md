<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_time_entries.sql","sha256":"e14e176998499b6b3db2d55b7312a93d59827e46d6c24d194968e722eada0799","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_time_entries.sql

## Missing names and codes show blank, and missing category shows 'No Category'

Lines 13-18

When an entry has no branch, division or time type, its code shows as an empty string. When it has no matching profile, the surname and given name show as empty strings. When it has no category, the category name shows 'No Category'.

## Job time entries list only committed time sheets

Lines 25-26

The job time entries list includes only entries whose time sheet has been committed (time_sheets.committed is not empty). Entries on uncommitted time sheets do not appear. An entry with no matching time sheet is also left out.

## Job time entries list skips entries with zero hours

Lines 27-27

The job time entries list shows only entries with hours greater than 0. Entries with 0 hours are left out.

## Job time entries list is limited to one job

Lines 28-28

The list returns only time entries for the job whose id is given (te.job = id).

## Optional filters by branch, division, time type and person

Lines 29-32

The list can be filtered by branch, division, time type (time_type) and person (uid). A filter that is missing or empty is ignored and does not narrow the list. A filter with a value shows only entries that match it exactly.

## Category filter and the 'none' value

Lines 33-39

A missing or empty category filter shows entries of every category. The category value 'none' shows only entries with no category (te.category is empty). Any other category value shows only entries in that category.

## Job time entries are sorted newest first and paged

Lines 40-41

The list is sorted by entry date, newest first. It returns the number of rows given by the limit, starting at the given offset.
