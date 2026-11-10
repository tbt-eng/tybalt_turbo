<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_time_entries_count.sql","sha256":"c64ada2c815d2fe8a8a0001da8eb2e1edbf0a96925d8af10543fd779680e5449","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_time_entries_count.sql

## Job time entry count includes only committed time sheets

Lines 4-5

The count of time entries for a job covers only entries whose time sheet has been committed (a non-empty committed value). Entries on time sheets that are not committed are not counted. An entry with no matching time sheet is not counted either.

## Job time entry count skips entries with zero hours

Lines 6-6

The count includes only time entries with more than 0 hours (hours > 0). Entries with exactly 0 hours are not counted.

## Job time entry count is for one job

Lines 7-7

The count covers only the time entries that belong to the job whose id is given (te.job = id).

## Optional filters on job time entry count

Lines 8-11

The count can be narrowed by branch, division, time type and staff member (uid). A filter that is missing or empty is ignored. A filter with a value counts only entries that match it exactly.

## Category filter on job time entry count, including 'none'

Lines 12-18

When the category filter is missing or empty, entries of every category are counted. When the filter is 'none', only entries with no category (an empty category) are counted. Any other category value counts only entries with exactly that category.
