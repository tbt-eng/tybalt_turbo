<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_time_summary.sql","sha256":"3700f6c61835e20ed96a3434e81ddd1003aefa51d5cddb0d0ae3cf372a581c9b","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_time_summary.sql

## Job time summary totals hours and date range

Lines 7-9

Turbo adds up the hours of the matching entries (total_hours). It gives the earliest and latest entry dates (earliest_entry, latest_entry).

## Job time summary lists distinct branches, divisions, time types and people

Lines 10-21

Turbo lists each distinct branch, division and time type (id and code) found in the matching entries. It also lists each distinct person (uid and full name as given name plus surname). Each value appears once, however many entries use it.

## Entries with no category show as 'No Category'

Lines 22-24

In the category list, an entry with no category appears with the id 'none' and the name 'No Category'. Other categories appear with their own id and name.

## Entries still count when lookup records are missing

Lines 25-32

The summary joins divisions, time types, profiles, categories and branches with left joins. An entry whose related record is missing is not dropped for that reason. The time sheet join is also a left join, but the committed condition removes entries with no matching time sheet.

## Job time summary counts only committed time sheets

Lines 31-32

The job time summary includes only time entries whose time sheet has a commit date (ts.committed is not empty). Entries on time sheets that are not yet committed are left out of the totals and lists.

## Job time summary ignores entries with zero hours

Lines 33-33

The job time summary leaves out time entries whose hours are 0 or less. Only entries with hours greater than 0 count.

## Job time summary is for one job

Lines 34-34

The job time summary covers only the time entries recorded against the job whose id is passed in (te.job = id).

## Optional filters by branch, division, time type and person

Lines 35-38

The summary can be narrowed by branch, division, time type or person (uid). A filter that is empty or not given is ignored. A filter that is set keeps only entries that match it exactly.

## Category filter treats 'none' as uncategorized entries

Lines 39-45

An empty or missing category filter applies no category filter. The value 'none' keeps only entries with no category. Any other value keeps only entries in that category.
